# uniwidth Architecture

**Version**: v0.2.0
**Date**: 2026-02-05
**Unicode**: 16.0.0

---

## Design Goals

1. **Performance**: 3-46x faster than existing solutions (proven in benchmarks)
2. **Correctness**: Full Unicode 16.0 compliance, including ZWJ emoji sequences
3. **Zero Allocations**: No GC pressure on ASCII paths
4. **Modern Go**: Leverage Go 1.25+ compiler optimizations and SWAR techniques
5. **Simple API**: Drop-in replacement for go-runewidth

---

## Core Architecture: 4-Tier O(1) Lookup

uniwidth uses a **4-tier lookup system** where all tiers operate in constant time O(1). This optimizes for the 90-95% common case while maintaining full Unicode coverage.

### Tier 1: ASCII Fast Path (O(1))

**Coverage**: ~95% of typical terminal content
**Performance**: 15-46x faster than go-runewidth

```go
// Tier 1: ASCII (0x00-0x7F)
if r < 0x80 {
    if r < 0x20 {
        return 0 // C0 controls
    }
    if r == 0x7F {
        return 0 // DELETE
    }
    return 1 // Printable ASCII
}
```

**Why this works**:
- Most TUI content is ASCII (English text, numbers, punctuation)
- Single comparison `r < 0x80` is incredibly fast
- Go compiler optimizes this to a few CPU instructions
- No memory lookups, no cache misses

For strings, the ASCII fast path uses SWAR optimization (see below).

### Tier 2: Common CJK Fast Path (O(1))

**Coverage**: ~80-90% of non-ASCII content
**Performance**: 4-14x faster than go-runewidth

```go
// CJK Unified Ideographs (20,992 characters)
if r >= 0x4E00 && r <= 0x9FFF {
    return 2
}

// Hangul Syllables (11,172 characters)
if r >= 0xAC00 && r <= 0xD7AF {
    return 2
}

// Hiragana + Katakana + Bopomofo (384 characters)
if r >= 0x3040 && r <= 0x312F {
    return 2
}

// CJK Compatibility Ideographs
if r >= 0xF900 && r <= 0xFAFF {
    return 2
}
```

**Why this works**:
- CJK characters cluster in large contiguous ranges
- Range checks (`>=` and `<=`) are O(1) operations
- Covers 99% of Japanese, Chinese, Korean text

### Tier 3: Common Emoji Fast Path (O(1))

**Coverage**: ~90% of commonly used emoji
**Performance**: 6-8x faster than go-runewidth

```go
// Emoticons (U+1F600-U+1F64F)
if r >= 0x1F600 && r <= 0x1F64F { return 2 }

// Misc Symbols and Pictographs (U+1F300-U+1F5FF)
if r >= 0x1F300 && r <= 0x1F5FF { return 2 }

// Transport and Map (U+1F680-U+1F6FF)
if r >= 0x1F680 && r <= 0x1F6FF { return 2 }

// Supplemental Symbols (U+1F900-U+1F9FF)
if r >= 0x1F900 && r <= 0x1F9FF { return 2 }

// Misc Symbols (U+2600-U+26FF), Dingbats (U+2700-U+27BF)
```

### Tier 4: 3-Stage Hierarchical Table (O(1))

**Coverage**: All remaining Unicode codepoints
**Performance**: O(1) — 3 array lookups + bit extraction

Replaced the previous O(log n) binary search with a compact 3-stage hierarchical table that encodes every Unicode codepoint (U+0000-U+10FFFF) as a 2-bit width value.

#### Table Structure

```
ROOT[256] → MIDDLE[17×64] → LEAVES[78×32]
```

- **ROOT**: 256 entries, indexed by `cp >> 13` (top 8 bits of plane + block)
- **MIDDLE**: 17 pages × 64 entries each, indexed by `(cp >> 7) & 0x3F`
- **LEAVES**: 78 pages × 32 bytes each, packed 2-bit encoding, indexed by `(cp >> 2) & 0x1F`

#### 2-Bit Width Encoding

```
0b00 = width 0 (control, combining, zero-width)
0b01 = width 1 (narrow, default)
0b10 = width 2 (wide: CJK, emoji, fullwidth)
0b11 = ambiguous (treated as width 1 in neutral context)
```

#### Lookup Code

```go
func tableLookupWidth(r rune) int {
    cp := uint32(r)
    rootIdx := widthRoot[cp>>13]
    midIdx  := widthMiddle[rootIdx][cp>>7&0x3F]
    packed  := widthLeaves[midIdx][cp>>2&0x1F]
    width   := (packed >> (2 * (cp & 0x03))) & 0x03
    if width == 3 {
        return 1 // ambiguous → narrow in neutral context
    }
    return int(width)
}
```

#### Memory Footprint

| Component | Size |
|-----------|------|
| ROOT | 256 bytes |
| MIDDLE | 1,088 bytes (17 × 64) |
| LEAVES | 2,496 bytes (78 × 32) |
| **Total** | **3,840 bytes (3.8 KB)** |

Compare to go-runewidth: ~500KB of tables.

---

## SWAR Optimization (SIMD Within A Register)

### isASCIIOnly() — ASCII Detection

Processes 8 bytes at a time by loading them into a uint64 and checking all high bits simultaneously:

```go
func isASCIIOnly(s string) bool {
    p := unsafe.StringData(s)
    const asciiMask = uint64(0x8080808080808080)

    // Process 8 bytes at a time
    for ; i+8 <= n; i += 8 {
        word := *(*uint64)(unsafe.Add(unsafe.Pointer(p), i))
        if word & asciiMask != 0 {
            return false // Non-ASCII byte found
        }
    }
    // Scalar tail for remaining 0-7 bytes
    ...
}
```

If any byte has its high bit set (>= 0x80), the AND with `0x8080808080808080` produces a non-zero result. Works regardless of endianness.

### asciiWidth() — Control Character Detection

Uses Daniel Lemire's SWAR underflow trick to detect control characters in 8-byte chunks:

```go
// Detect bytes < 0x20 (C0 controls):
// Subtracting 0x20 from a byte < 0x20 causes unsigned underflow,
// setting the high bit. AND with ~word isolates genuine underflows.
hasLow := (word - 0x2020202020202020) & ^word & 0x8080808080808080

// Detect byte == 0x7F (DELETE):
// XOR with 0x7F zeros out any 0x7F bytes, then zero-byte detection
// finds them via the underflow pattern.
xored := word ^ 0x7F7F7F7F7F7F7F7F
has7F := (xored - 0x0101010101010101) & ^xored & 0x8080808080808080
```

If neither `hasLow` nor `has7F` is set, the entire 8-byte chunk has no control characters and `width += 8` directly.

### Short String Optimization

Strings shorter than 8 bytes use a fused single-pass loop that combines ASCII detection and width counting, avoiding the overhead of calling both `isASCIIOnly()` and `asciiWidth()` separately:

```go
if len(s) < 8 {
    width, isASCII := 0, true
    for i := 0; i < len(s); i++ {
        b := s[i]
        if b >= 0x80 { isASCII = false; break }
        if b >= 0x20 && b != 0x7F { width++ }
    }
    if isASCII { return width }
}
```

---

## ZWJ State Machine

StringWidth uses a forward-scan state machine for correct handling of multi-rune emoji sequences. Inspired by Ghostty's approach, adapted for width calculation.

### States

```
State 0: Default (not in an emoji sequence)
State 1: After Extended_Pictographic character (may start ZWJ/modifier sequence)
State 2: After EP + (Extend*) + ZWJ (expecting joined emoji)
```

### Transitions

```
[Default] ──EP(w>0)──→ [Emoji]
[Emoji]   ──ZWJ──────→ [EmojiZWJ]
[EmojiZWJ]──EP────────→ [Emoji]     (joined, width 0)
[EmojiZWJ]──other─────→ [Default]   (broken sequence)
[Emoji]   ──modifier──→ [Emoji]     (skin tone, width 0)
[Emoji]   ──VS────────→ [Emoji]     (variation selector, width 0)
[any]     ──w==0──────→ [preserve]  (combining marks keep state for Extend*)
```

### Supported Sequences

| Sequence | Example | Width |
|----------|---------|-------|
| ZWJ family | 👨‍👩‍👧‍👦 | 2 |
| Skin tone | 👍🏽 | 2 |
| Professional | 👩🏽‍🔬 | 2 |
| Rainbow flag | 🏳️‍🌈 | 2 |
| Heart + fire | ❤️‍🔥 | 2 |
| Country flag | 🇺🇸 | 2 |
| VS-16 emoji | ☀️ | 2 |
| VS-15 text | ☀︎ | 1 |

### Extended_Pictographic Detection

`isExtendedPictographic()` uses range checks ordered by frequency of occurrence in real-world emoji usage:

1. SMP emoji blocks (U+1F000-U+1FAFF) — covers ~95% of emoji
2. BMP: Misc Symbols (U+2600-U+27BF)
3. BMP: Misc Technical (U+2300-U+23FF)
4. BMP: Misc Symbols and Arrows (U+2B00-U+2BFF)
5. BMP: Arrow symbols (U+2194-U+21AA)
6. BMP: Geometric Shapes (U+25A0-U+25FF)
7. SMP: Legacy Computing (U+1FB00-U+1FFFD)
8. Individual characters: ©, ®, ‼, ⁉, ™, ℹ, 〰, 〽, ㊗, ㊙

---

## String Width Calculation

### Flow

```
StringWidth(s)
    │
    ├─ len < 8? → Fused ASCII check + width count
    │
    ├─ isASCIIOnly(s)? → asciiWidth(s) [SWAR]
    │
    └─ Unicode path:
        ├─ Convert to []rune (1 allocation)
        └─ State machine loop:
            ├─ ZWJ handling (state transitions)
            ├─ Emoji modifier handling
            ├─ VS in emoji context
            ├─ Regional indicator pairs
            ├─ Variation selector lookahead
            └─ Default: RuneWidth(r)
```

### Allocation Behavior

| Input | Allocations | Reason |
|-------|-------------|--------|
| ASCII-only, any length | 0 | SWAR fast path, no rune conversion |
| Unicode, short (< ~32 runes) | 0 | Go stack-allocates small `[]rune` slices |
| Unicode, long | 1 | `[]rune` heap allocation for lookahead |

---

## Options API

### Functional Options Pattern

```go
type Options struct {
    EastAsianAmbiguous EAWidth // 1 or 2
    EmojiPresentation  bool    // true or false
}

type Option func(*Options)

func WithEastAsianAmbiguous(width EAWidth) Option { ... }
func WithEmojiPresentation(emoji bool) Option { ... }
```

Ambiguous characters (width encoding `0b11` in the table) return width based on the configured option. Default: narrow (width 1).

---

## Table Generation

### Source Data

Tables are generated from official Unicode 16.0 data:
- `EastAsianWidth.txt` — East Asian Width property
- `emoji-data.txt` — Emoji presentation data

### Process

```bash
go generate ./...
# or
go run cmd/generate-tables/main.go
```

1. Download Unicode 16.0 data files
2. Parse East Asian Width (W, F, N, A properties)
3. Parse Emoji data
4. Build full codepoint-to-width mapping (U+0000-U+10FFFF)
5. Compress into 3-stage hierarchical table via page deduplication
6. Generate `tables_generated.go`

### Hot Path Filtering

The 3-stage table encodes ALL codepoints, but Tiers 1-3 short-circuit before reaching the table for common characters. The table primarily serves rare characters that don't fall into the hot paths.

---

## Performance Characteristics

### Time Complexity

| Operation | All Tiers |
|-----------|-----------|
| `RuneWidth()` | O(1) |
| `StringWidth()` ASCII-only | O(n/8) via SWAR |
| `StringWidth()` Unicode | O(n) per rune |

### Space Complexity

| Component | Size |
|-----------|------|
| Code (uniwidth.go + options.go) | ~10 KB |
| 3-stage table (tables_generated.go) | 3.8 KB |
| Binary search tables (legacy, for Options API) | ~3 KB |
| Runtime (ASCII path) | 0 bytes |
| **Total** | **~17 KB** |

Compare to go-runewidth: ~500KB.

### Benchmark Results

| Category | Time | Allocs | vs go-runewidth |
|----------|------|--------|-----------------|
| ASCII short (5 chars) | ~7 ns | 0 | 15-22x faster |
| ASCII medium (44 chars) | ~20 ns | 0 | 30-46x faster |
| CJK short (4 chars) | ~25 ns | 0 | 5-14x faster |
| ZWJ family (👨‍👩‍👧‍👦) | ~95 ns | 0 | New capability |
| Emoji modifier (👍🏽) | ~40 ns | 0 | New capability |
| Mixed (ASCII + CJK + emoji) | ~65 ns | 0 | 6-8x faster |

---

## Comparison with go-runewidth

### Architectural Differences

| Aspect | uniwidth | go-runewidth |
|--------|----------|--------------|
| Lookup strategy | 4-tier O(1) | Binary search O(log n) |
| Table size | 3.8 KB | ~500 KB |
| ASCII path | SWAR (8 bytes/iter) | Grapheme + binary search |
| ZWJ emoji | Forward-scan state machine | Delegates to uax29 |
| Allocations (ASCII) | 0 | 0 |
| Go version | 1.25+ | 1.9+ |

### Trade-offs

**uniwidth wins**: Performance (3-46x), memory (130x smaller tables), ZWJ correctness with minimal overhead.

**go-runewidth wins**: Mature ecosystem (10+ years), wider Go version support, full UAX #29 grapheme clustering via uax29.

---

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| 4-tier lookup | 95% of content is ASCII; O(1) >> O(log n) |
| 3-stage table (Tier 4) | O(1) for all codepoints, only 3.8KB |
| Forward-scan ZWJ state machine | Simpler than reverse iteration, covers 99%+ of emoji |
| SWAR over auto-vectorization | Explicit uint64 word processing, portable, predictable |
| Functional options | Clean, extensible, backward compatible, zero alloc when unused |
| Generate tables from Unicode data | Easy version updates, correctness guaranteed, reproducible |
| Defer full UAX #29 | 2-5x performance cost, <1% real-world demand in terminals |

---

## Future Optimizations

### Explicit SIMD (Later)
- **Go assembly** (Plan 9 `.s` files): Hand-written SSE2/AVX2/NEON for `isASCIIOnly()` and `asciiWidth()`. Potential 16-32 bytes/iter (2-4x over current SWAR).
- **`archsimd` package** (Go 1.26+): Portable SIMD intrinsics when `GOEXPERIMENT=simd` stabilizes.

### PGO (Profile-Guided Optimization)
- Collect real-world profiles from TUI applications
- Feed to Go compiler for better inlining and branch prediction
- Expected 10-20% improvement on hot paths

---

## Testing Strategy

### Test Categories

| Category | Tests | Coverage |
|----------|-------|----------|
| Core unit tests | ASCII, CJK, Emoji, Zero-width | RuneWidth, StringWidth |
| ZWJ sequences | 15 test cases | Family, professions, flags, modifiers |
| Emoji modifiers | 8 test cases | Skin tones, combined sequences |
| Edge cases | 11 test cases | Standalone ZWJ, orphan modifiers, boundaries |
| Conformance | All Unicode categories | Categories, combining marks, controls, fullwidth |
| Fuzzing | Go native | No panics, valid widths (0-2) |
| Benchmarks | 20+ scenarios | ASCII, CJK, Emoji, ZWJ, TUI |

**Coverage**: 96.4% (target: >90%)

---

*Architecture document for uniwidth v0.2.0*
*Updated: 2026-02-05*
*Unicode Version: 16.0.0*
