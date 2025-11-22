# Changelog

All notable changes to uniwidth will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned for v1.0.0
- API freeze and stability commitment
- Extended test coverage (>95%)
- Performance regression test suite
- Additional locale support
- Migration guide improvements

### Planned for v0.2.0+
- Grapheme cluster support for complex emoji sequences
- Explicit SIMD optimizations for AVX-512
- Profile-Guided Optimization (PGO) support

## [0.1.0] - 2025-11-22

**Stable Release**: First stable release after 35 days of beta testing. All known issues from beta have been resolved.

### Added
- 🐛 **Bug Fix**: Variation selectors (U+FE0E, U+FE0F) now handled correctly
  - Text variation selector (U+FE0E) forces width 1
  - Emoji variation selector (U+FE0F) forces width 2
  - Example: "☀︎" (sun + text variant) now correctly returns width 1
- 🐛 **Bug Fix**: Regional indicator pairs (flags) now handled correctly
  - Two consecutive regional indicators count as width 2 (not 4)
  - Example: "🇺🇸" (U+1F1FA + U+1F1F8) now correctly returns width 2
- 📁 **Project Structure**: Reorganized documentation
  - Created `docs/` for public documentation
  - Created `docs/dev/` for development documentation (gitignored)
  - Moved ARCHITECTURE.md and POC_RESULTS.md to `docs/`
  - Added `docs/dev/INDEX.md` (Kanban-style tracker)
  - Added `docs/dev/ROADMAP.md` (release planning)
- 📝 **Documentation**: Added comprehensive CLAUDE.md for AI assistance
- 🧪 **Tests**: Added edge case tests
  - `TestStringWidth_VariationSelectors` (6 test cases)
  - `TestStringWidth_RegionalIndicators` (5 test cases)
  - `TestIsRegionalIndicator` (7 test cases)

### Changed
- 🔄 **StringWidth**: Now converts to `[]rune` for lookahead (variation selectors)
  - Trade-off: 1 allocation for Unicode strings (correctness > performance)
  - ASCII fast path still has 0 allocations
- 📊 **Test Coverage**: Increased from 84.6% → 87.1% (+2.5%)

### Performance Impact
- ASCII strings: No change (0 allocations, ~5 ns/op)
- Unicode strings: Minimal impact (<1 ns/op, 1 allocation for `[]rune` conversion)
- Still 9-23x faster than go-runewidth overall

### Fixed (from beta)
- ✅ Combining marks edge cases (U+1AD7, U+1AFF) - added to zero-width tables
- ✅ Boundary issues (U+4DFF, U+303F, U+3100) - table boundaries corrected
- ✅ Surrogate pair handling (U+10000) - Linear B Syllable now handled correctly

### Known Limitations
- Grapheme clusters not yet supported (planned for v0.2.0+)
  - Complex emoji ZWJ sequences counted as sum of parts
  - Single-character emoji work correctly

## [0.1.0] - 2025-10-15

> 📝 **Note**: This version was superseded by v0.1.0-beta with critical bug fixes.

### Added
- Initial release of uniwidth library
- Tiered lookup strategy (4 tiers: ASCII, CJK/Emoji, Zero-width, Binary search)
- Full Unicode 16.0.0 support
- Options API for East Asian Ambiguous character handling
- Options API for emoji presentation mode
- Zero allocation design (0 B/op, 0 allocs/op)
- SIMD auto-vectorization for ASCII detection (Go 1.25+)
- Table generation from official Unicode data files
- Comprehensive test suite (84.6% coverage)
- Conformance tests for Unicode categories
- Fuzzing tests for robustness
- Benchmarks vs go-runewidth (3-46x speedup proven)

### Performance
- **ASCII strings**: 15-46x faster than go-runewidth
- **CJK strings**: 4-14x faster than go-runewidth
- **Emoji strings**: 6-8x faster than go-runewidth
- **Zero allocations**: All operations are allocation-free
- **Small footprint**: ~13KB total (code + tables)

### API
- `RuneWidth(r rune) int` - Calculate visual width of a rune
- `StringWidth(s string) int` - Calculate visual width of a string
- `RuneWidthWithOptions(r rune, opts ...Option) int` - Rune width with options
- `StringWidthWithOptions(s string, opts ...Option) int` - String width with options
- `WithEastAsianAmbiguous(width EAWidth) Option` - Configure ambiguous width
- `WithEmojiPresentation(emoji bool) Option` - Configure emoji presentation

### Documentation
- README.md with quick start and examples
- ARCHITECTURE.md with detailed technical design
- POC_RESULTS.md with benchmark analysis
- LICENSE (MIT)
- Comprehensive godoc comments

### Known Limitations
- Grapheme clustering not yet implemented (complex emoji sequences counted as sum of parts)
- Some edge cases at Unicode range boundaries (will be fixed in v0.1.1)
- Zero-width space (U+200B) handling needs improvement
- Test coverage 84.6% (target 90%+ in v0.1.1)

### Requirements
- Go 1.25.0 or later (required for optimal performance)
- No external dependencies except go-runewidth (for benchmarks only)

---

## Version History

### Naming Convention
- **Major**: Breaking API changes
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, performance improvements

### Stability
- v0.x.x: Pre-release, API may change
- v1.x.x: Stable API, production ready

---

## Upgrade Guide

### From PoC to v0.1.0
- No breaking changes
- Generated tables now included
- Options API added (optional, backward compatible)

### From go-runewidth to uniwidth
Simple drop-in replacement:

```go
// Before
import "github.com/mattn/go-runewidth"
width := runewidth.StringWidth(s)

// After
import "github.com/unilibs/uniwidth"
width := uniwidth.StringWidth(s)
```

**Performance improvement**: 3-46x faster, zero code changes!

**Note**: Grapheme clustering behavior may differ for complex emoji sequences.

---

## Maintenance

### Update Unicode Version

To update to a newer Unicode version:

1. Update URLs in `cmd/generate-tables/main.go`:
   ```go
   const unicodeVersion = "16.0.0" // Change this
   const eastAsianWidthURL = "https://www.unicode.org/Public/16.0.0/..." // And this
   ```

2. Regenerate tables:
   ```bash
   go generate ./...
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

4. Update benchmarks:
   ```bash
   go test -bench=. -benchmem
   ```

---

*For detailed performance analysis, see [docs/POC_RESULTS.md](docs/POC_RESULTS.md)*
*For architecture details, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)*
*For release planning, see [docs/dev/ROADMAP.md](docs/dev/ROADMAP.md) (development only)*
