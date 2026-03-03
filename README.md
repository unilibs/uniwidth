# uniwidth - Modern Unicode Width Calculation for Go

[![Go Version](https://img.shields.io/github/go-mod/go-version/unilibs/uniwidth?label=Go)](https://go.dev/dl/)
[![CI Status](https://github.com/unilibs/uniwidth/actions/workflows/ci.yml/badge.svg)](https://github.com/unilibs/uniwidth/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/unilibs/uniwidth)](https://goreportcard.com/report/github.com/unilibs/uniwidth)
[![codecov](https://codecov.io/gh/unilibs/uniwidth/branch/main/graph/badge.svg)](https://codecov.io/gh/unilibs/uniwidth)
[![Go Reference](https://pkg.go.dev/badge/github.com/unilibs/uniwidth.svg)](https://pkg.go.dev/github.com/unilibs/uniwidth)
[![License](https://img.shields.io/github/license/unilibs/uniwidth)](LICENSE)
[![Release](https://img.shields.io/github/v/release/unilibs/uniwidth)](https://github.com/unilibs/uniwidth/releases)
[![Stars](https://img.shields.io/github/stars/unilibs/uniwidth)](https://github.com/unilibs/uniwidth/stargazers)

**uniwidth** is a modern, high-performance Unicode width calculation library for Go 1.25+. It provides **3-46x faster** width calculation compared to existing solutions through a 4-tier O(1) lookup architecture, SWAR optimization, and a ZWJ-aware emoji state machine.

## Performance

Based on comprehensive benchmarks vs `go-runewidth`:

- **ASCII strings**: 15-46x faster (SWAR, 8 bytes/iter)
- **CJK strings**: 4-14x faster (O(1) table lookup)
- **Mixed/Emoji strings**: 6-8x faster
- **ZWJ emoji**: Correct width (👨‍👩‍👧‍👦 = 2, ~95 ns)
- **Zero allocations**: 0 B/op, 0 allocs/op for ASCII paths

Run benchmarks yourself: `cd bench && go test -bench=. -benchmem`

## Features

- **3-46x faster** than go-runewidth (proven in benchmarks)
- **All tiers O(1)** — 4-tier lookup with 3-stage hierarchical table (3.8KB)
- **ZWJ-aware** — family emoji, skin tones, flags handled correctly
- **SWAR optimized** — ASCII detection and width counting at 8 bytes/iter
- **Zero allocations** for ASCII strings (no GC pressure)
- **Thread-safe** (immutable design, no global state)
- **Unicode 16.0** support
- **Modern API** (Go 1.25+, functional options pattern)

## Installation

```bash
go get github.com/unilibs/uniwidth
```

**Requirements**: Go 1.25 or later

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/unilibs/uniwidth"
)

func main() {
    // Calculate width of a string
    width := uniwidth.StringWidth("Hello 世界")
    fmt.Println(width) // Output: 10 (Hello=5, space=1, 世界=4)

    // Calculate width of a single rune
    w := uniwidth.RuneWidth('世')
    fmt.Println(w) // Output: 2

    // ASCII-only strings are super fast!
    width = uniwidth.StringWidth("Hello, World!")
    fmt.Println(width) // Output: 13
}
```

### ZWJ Emoji Sequences

```go
// ZWJ family emoji — correctly returns 2, not 8
width := uniwidth.StringWidth("👨‍👩‍👧‍👦")
fmt.Println(width) // Output: 2

// Skin tone modifiers — correctly returns 2, not 4
width = uniwidth.StringWidth("👍🏽")
fmt.Println(width) // Output: 2

// Rainbow flag
width = uniwidth.StringWidth("🏳️‍🌈")
fmt.Println(width) // Output: 2

// Country flags
width = uniwidth.StringWidth("🇺🇸")
fmt.Println(width) // Output: 2
```

### Options API

Configure handling of ambiguous-width characters:

```go
import "github.com/unilibs/uniwidth"

// East Asian locale (ambiguous characters are wide)
opts := []uniwidth.Option{
    uniwidth.WithEastAsianAmbiguous(uniwidth.EAWide),
}
width := uniwidth.StringWidthWithOptions("±½", opts...)
fmt.Println(width) // Output: 4 (each character is 2 columns)

// Neutral locale (ambiguous characters are narrow) - DEFAULT
opts = []uniwidth.Option{
    uniwidth.WithEastAsianAmbiguous(uniwidth.EANarrow),
}
width = uniwidth.StringWidthWithOptions("±½", opts...)
fmt.Println(width) // Output: 2 (each character is 1 column)
```

### Real-World TUI Examples

```go
// Terminal prompt
prompt := "❯ Enter command: "
width := uniwidth.StringWidth(prompt)
fmt.Printf("Prompt width: %d columns\n", width)

// Table cell padding
text := "Hello 世界"
padding := 20 - uniwidth.StringWidth(text)
fmt.Printf("%s%s\n", text, strings.Repeat(" ", padding))

// Truncate to fit terminal width
func truncate(s string, maxWidth int) string {
    width := 0
    for i, r := range s {
        w := uniwidth.RuneWidth(r)
        if width+w > maxWidth {
            return s[:i] + "…"
        }
        width += w
    }
    return s
}
```

## Architecture

### 4-Tier O(1) Lookup

uniwidth uses a multi-tier approach where **all tiers are O(1)**:

1. **Tier 1: ASCII Fast Path** (O(1))
   - Covers ~95% of typical terminal content
   - SWAR `isASCIIOnly()` + `asciiWidth()` process 8 bytes/iter
   - Short strings (< 8 bytes) use fused single-pass loop

2. **Tier 2: Common CJK** (O(1))
   - CJK Unified Ideographs, Hangul Syllables, Hiragana/Katakana
   - Simple range checks for 32,000+ characters

3. **Tier 3: Common Emoji** (O(1))
   - Emoticons, Pictographs, Dingbats, Symbols
   - Range checks for ~1,200 emoji codepoints

4. **Tier 4: 3-Stage Table** (O(1))
   - ROOT[256] → MIDDLE[17×64] → LEAVES[78×32]
   - 2-bit width encoding, 3.8KB total
   - Covers all remaining Unicode codepoints in 3 array lookups

### ZWJ State Machine

Forward-scan state machine for correct emoji sequence handling:
- **3 states**: default → emoji → emojiZWJ
- Handles: ZWJ sequences, skin tone modifiers, variation selectors, flag pairs
- Inspired by Ghostty's approach, adapted for width calculation

### SWAR Optimization

ASCII paths use SIMD Within A Register (SWAR) for high throughput:
- `isASCIIOnly()`: uint64 word AND with `0x8080808080808080` mask
- `asciiWidth()`: Daniel Lemire's underflow trick for control character detection
- Both process 8 bytes per iteration with zero allocations

## Benchmarks

```
goos: windows
goarch: amd64

BenchmarkStringWidth_ASCII_Short     ~7 ns/op     0 B/op   0 allocs/op
BenchmarkStringWidth_ASCII_Medium   ~20 ns/op     0 B/op   0 allocs/op
BenchmarkStringWidth_CJK_Short     ~25 ns/op     0 B/op   0 allocs/op
BenchmarkStringWidth_ZWJ_Family    ~95 ns/op     0 B/op   0 allocs/op
BenchmarkStringWidth_EmojiModifier ~40 ns/op     0 B/op   0 allocs/op
```

Run benchmarks yourself:
```bash
go test -bench=. -benchmem
```

## Use Cases

Perfect for:
- **TUI frameworks** (terminal rendering hot paths)
- **Terminal emulators** (text layout calculations)
- **CLI tools** (table alignment, formatting)
- **Text editors** (cursor positioning, column calculation)
- **Any high-performance text width calculation**

## Migration from go-runewidth

uniwidth provides a compatible API for easy migration:

```go
// Before (go-runewidth)
import "github.com/mattn/go-runewidth"
width := runewidth.StringWidth(s)

// After (uniwidth) - drop-in replacement!
import "github.com/unilibs/uniwidth"
width := uniwidth.StringWidth(s)
```

**Performance improvement**: 3-46x faster, zero code changes!

## Documentation

- [API Reference](https://pkg.go.dev/github.com/unilibs/uniwidth) - Full godoc documentation
- [Benchmark Comparisons](bench/README.md) - Performance comparison vs go-runewidth
- [Architecture Design](docs/ARCHITECTURE.md) - Technical deep dive & design decisions
- [Changelog](CHANGELOG.md) - Version history & upgrade guide
- [Roadmap](ROADMAP.md) - What's next for uniwidth

## Testing

```bash
# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Run with coverage
go test -cover
```

Current test coverage: **100%** (library package)

## Development Status

**Current**: v0.2.0

> This library is stable and production-ready. The API is backward-compatible across minor versions. ZWJ emoji sequences, skin tone modifiers, variation selectors, and flag emoji are all handled correctly.

**v0.2.0 Highlights**:
- All 4 lookup tiers are now O(1) (3-stage table replaced binary search)
- SWAR ASCII optimization (8 bytes/iter)
- ZWJ emoji state machine (👨‍👩‍👧‍👦 = width 2)
- Emoji modifier support (👍🏽 = width 2)
- 100% test coverage (library package)
- Automated benchmark CI (regression detection + library comparison)

**Roadmap** (v0.3.0+):
- Non-ASCII StringWidth path optimization
- Profile-Guided Optimization (PGO)
- Explicit SIMD via Go assembly and `archsimd`
- Unicode 17.0 tables

## Contributing

Contributions welcome! This is part of the [unilibs](https://github.com/unilibs) organization - modern Unicode libraries for Go.

## License

MIT License - see [LICENSE](LICENSE) file

## Related Projects

Built by the [Phoenix TUI Framework](https://github.com/phoenix-tui/phoenix) team.

Part of the **unilibs** ecosystem:
- **uniwidth** - Unicode width calculation (this project)
- **unigrapheme** - Grapheme clustering (planned)
- More Unicode utilities coming soon!

## Support

- Issues: [GitHub Issues](https://github.com/unilibs/uniwidth/issues)
- Discussions: [GitHub Discussions](https://github.com/unilibs/uniwidth/discussions)

---

## Special Thanks

**Professor Ancha Baranova** - This project would not have been possible without her invaluable help and support. Her assistance was crucial in bringing uniwidth to life.

---

**Made with care by the Phoenix team** | **Powered by Go 1.25+**
