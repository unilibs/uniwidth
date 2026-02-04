# Changelog

All notable changes to uniwidth will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Profile-Guided Optimization (PGO) support
- Unicode 17.0 preparation
- Benchmark CI for regression detection
- Explicit SIMD via Go assembly and `archsimd` (Go 1.26+)
- API stability review based on community feedback

## [0.2.0] - 2026-02-05

Major performance and emoji correctness release. All four lookup tiers are now O(1), ZWJ emoji sequences are handled correctly, and ASCII paths use SWAR for 8 bytes/iter throughput.

### Added
- **ZWJ emoji sequence support**: Forward-scan state machine with 3 states (default/emoji/emojiZWJ). Family emoji 👨‍👩‍👧‍👦 now correctly returns width 2, not 8.
- **Emoji modifier (skin tone) support**: U+1F3FB-U+1F3FF (Fitzpatrick types) combine with preceding emoji. 👍🏽 now correctly returns width 2, not 4.
- **`isExtendedPictographic()` helper**: Range-based Extended_Pictographic detection, frequency-ordered for minimal branch mispredictions.
- **`isEmojiModifier()` helper**: Fitzpatrick skin tone modifier detection.
- **48 new test cases**: ZWJ sequences (15), emoji modifiers (8), edge cases (11), Extended_Pictographic validation (18), emoji modifier validation (9).
- **ZWJ benchmarks**: Family (~95 ns), couple with heart (~82 ns), skin tone modifier (~40 ns), mixed ZWJ text (~357 ns). All zero allocations.
- **Three-way benchmark suite** (`bench/`): uniwidth vs go-runewidth vs rivo/uniseg.

### Changed
- **Tier 4 lookup**: Replaced O(log n) binary search with O(1) 3-stage hierarchical table. ROOT[256] → MIDDLE[17×64] → LEAVES[78×32], 3.8KB total. All Unicode codepoints resolved in 3 array lookups.
- **ASCII detection**: SWAR `isASCIIOnly()` processes 8 bytes/iter via uint64 word with `0x8080808080808080` mask. No unsafe pointer escapes.
- **ASCII width counting**: SWAR `asciiWidth()` uses Daniel Lemire's underflow trick for control character detection in 8-byte chunks.
- **Short string optimization**: Strings < 8 bytes use a fused single-pass loop that combines ASCII check and width counting, avoiding SWAR function call overhead.
- **Test coverage**: 87.1% → 96.4% (+9.3%).

### Performance
- **ASCII**: 3-46x faster than go-runewidth (SWAR fast paths)
  - Short (5 chars): ~7 ns, 0 allocs
  - Medium (44 chars): ~20 ns, 0 allocs
  - Long (234 chars): ~50 ns, 0 allocs
- **CJK**: 30-35% faster from O(1) table lookup (previously O(log n))
- **ZWJ sequences**: New capability, ~95 ns for family emoji, 0 allocs
- **Emoji modifiers**: New capability, ~40 ns for skin tone, 0 allocs

## [0.1.0] - 2025-11-22

First stable release after beta testing. Variation selector and flag emoji bugs fixed.

### Added
- Variation selector handling: U+FE0E (text, width 1) and U+FE0F (emoji, width 2)
- Regional indicator pair handling: Flag emoji 🇺🇸 = width 2, not 4
- `isRegionalIndicator()` helper function
- Edge case tests: variation selectors (6), regional indicators (5), helper validation (7)
- Project structure: `docs/` for public docs, `docs/dev/` for dev docs (gitignored)

### Changed
- `StringWidth()` now converts to `[]rune` for lookahead (variation selectors require it)
  - Trade-off: 1 allocation for Unicode strings (correctness > performance)
  - ASCII fast path still has 0 allocations
- Test coverage: 84.6% → 87.1%

### Fixed (from beta)
- Combining marks edge cases (U+1AD7, U+1AFF)
- Boundary issues (U+4DFF, U+303F, U+3100)
- Surrogate pair handling (U+10000 - Linear B Syllable)

## [0.1.0-beta] - 2025-10-15

Initial public beta. Core architecture proven with 3.9-46x speedup over go-runewidth.

### Added
- 4-tier lookup architecture: ASCII O(1) → CJK/Emoji O(1) → Zero-width O(1) → Binary search O(log n)
- Core API: `RuneWidth()`, `StringWidth()`
- Options API: `WithEastAsianAmbiguous()`, `WithEmojiPresentation()`
- Full Unicode 16.0 support via generated tables from official data
- Zero allocation design for all code paths
- Table generation from EastAsianWidth.txt and emoji-data.txt
- Comprehensive test suite (84.6% coverage)
- Conformance tests for Unicode categories
- Fuzzing tests (Go native)
- Benchmarks vs go-runewidth (3-46x speedup proven)

### Performance
- ASCII strings: 15-46x faster than go-runewidth
- CJK strings: 4-14x faster than go-runewidth
- Emoji strings: 6-8x faster than go-runewidth
- Zero allocations: 0 B/op, 0 allocs/op

### API
- `RuneWidth(r rune) int`
- `StringWidth(s string) int`
- `RuneWidthWithOptions(r rune, opts ...Option) int`
- `StringWidthWithOptions(s string, opts ...Option) int`
- `WithEastAsianAmbiguous(width EAWidth) Option`
- `WithEmojiPresentation(emoji bool) Option`

### Known Limitations
- Grapheme clusters not supported (complex emoji ZWJ sequences counted as sum of parts)
- Some combining marks edge cases at boundaries
- Test coverage 84.6% (target 90%+)

### Requirements
- Go 1.25.0 or later
- No external dependencies

---

## Version History

| Version | Date | Highlights |
|---------|------|------------|
| 0.2.0 | 2026-02-05 | ZWJ emoji, SWAR, O(1) 3-stage table |
| 0.1.0 | 2025-11-22 | Stable release, variation selectors, flags |
| 0.1.0-beta | 2025-10-15 | Initial beta, 4-tier architecture |

## Upgrade Guide

### From v0.1.0 to v0.2.0
- No breaking API changes
- ZWJ sequences now return correct width (e.g., 👨‍👩‍👧‍👦 = 2, was 8)
- Emoji modifiers now return correct width (e.g., 👍🏽 = 2, was 4)
- All tiers are now O(1) (Tier 4 upgraded from binary search to table lookup)
- ASCII paths are significantly faster (SWAR optimization)

### From go-runewidth to uniwidth
Drop-in replacement:

```go
// Before
import "github.com/mattn/go-runewidth"
width := runewidth.StringWidth(s)

// After
import "github.com/unilibs/uniwidth"
width := uniwidth.StringWidth(s)
```

---

*For architecture details, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)*
