# uniwidth - Modern Unicode Width Calculation for Go

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://go.dev/dl/)
[![CI Status](https://github.com/unilibs/uniwidth/actions/workflows/ci.yml/badge.svg)](https://github.com/unilibs/uniwidth/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/unilibs/uniwidth)](https://goreportcard.com/report/github.com/unilibs/uniwidth)
[![Coverage](https://img.shields.io/badge/Coverage-87.1%25-brightgreen)](https://github.com/unilibs/uniwidth)
[![Go Reference](https://pkg.go.dev/badge/github.com/unilibs/uniwidth.svg)](https://pkg.go.dev/github.com/unilibs/uniwidth)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Beta-orange)](https://github.com/unilibs/uniwidth)
[![Version](https://img.shields.io/badge/Version-0.1.0--beta-orange)](CHANGELOG.md)

**uniwidth** is a modern, high-performance Unicode width calculation library for Go 1.25+. It provides **3.9-46x faster** width calculation compared to existing solutions through tiered lookup optimization and Go 1.25+ compiler features.

## 🚀 Performance

Based on comprehensive benchmarks vs `go-runewidth`:

- **ASCII strings**: 15-46x faster
- **CJK strings**: 4-14x faster
- **Mixed/Emoji strings**: 6-8x faster
- **Zero allocations**: 0 B/op, 0 allocs/op

Run benchmarks yourself: `cd bench && go test -bench=. -benchmem`

## ✨ Features

- 🚀 **3.9-46x faster** than go-runewidth (proven in benchmarks)
- 💎 **Zero allocations** (no GC pressure)
- 🧵 **Thread-safe** (immutable design, no global state)
- 🎯 **Unicode 16.0** support
- 🔧 **Modern API** (Go 1.25+, clean design)
- 📊 **Tiered lookup** (O(1) for 90-95% of cases)

## 📦 Installation

```bash
go get github.com/unilibs/uniwidth
```

**Requirements**: Go 1.25 or later

## 🔧 Usage

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

### Options API (NEW!)

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

### Performance-Critical Code

```go
// ASCII fast path (46x faster than go-runewidth!)
text := "Hello, World!"
width := uniwidth.StringWidth(text) // ~4.6 ns/op

// CJK fast path (14x faster!)
text := "你好世界"
width := uniwidth.StringWidth(text) // ~33.7 ns/op

// Mixed content (8x faster!)
text := "Hello 👋 World"
width := uniwidth.StringWidth(text) // ~65.9 ns/op

// All with zero allocations!
```

## 🏗️ Architecture

### Tiered Lookup Strategy

uniwidth uses a multi-tier approach for optimal performance:

1. **Tier 1: ASCII Fast Path** (O(1))
   - Covers ~95% of typical terminal content
   - Uses simple `len(s)` for ASCII-only strings
   - 15-46x faster than binary search

2. **Tier 2: Common CJK & Emoji** (O(1))
   - Range checks for frequent characters
   - CJK Unified Ideographs: 20,992 characters
   - Common emoji ranges
   - 4-14x faster than binary search

3. **Tier 3: Binary Search Fallback** (O(log n))
   - For rare characters not in hot paths
   - Minimal overhead (~5-10% of cases)

### Go 1.25+ Optimizations

- **SIMD Auto-Vectorization**: ASCII detection uses SSE2/AVX2
- **Aggressive Inlining**: Hot paths compile to minimal instructions
- **Zero Allocations**: No heap allocations, no GC pressure

## 📊 Benchmarks

```
BenchmarkStringWidth_ASCII_Short_Uniwidth-12     149590729   9.500 ns/op   0 B/op   0 allocs/op
BenchmarkStringWidth_ASCII_Short_GoRunewidth-12   10065044  150.1 ns/op   0 B/op   0 allocs/op
                                                             ^^^^^^^^^^
                                                             15.8x faster!

BenchmarkStringWidth_CJK_Short_Uniwidth-12        19064941   63.64 ns/op   0 B/op   0 allocs/op
BenchmarkStringWidth_CJK_Short_GoRunewidth-12      2771077  368.0 ns/op   0 B/op   0 allocs/op
                                                             ^^^^^^^^^^^
                                                             5.8x faster!
```

Run benchmarks yourself:
```bash
go test -bench=. -benchmem
```

## 🎯 Use Cases

Perfect for:
- **TUI frameworks** (terminal rendering hot paths)
- **Terminal emulators** (text layout calculations)
- **CLI tools** (table alignment, formatting)
- **Text editors** (cursor positioning, column calculation)
- **Any high-performance text width calculation**

## 🔄 Migration from go-runewidth

uniwidth provides a compatible API for easy migration:

```go
// Before (go-runewidth)
import "github.com/mattn/go-runewidth"
width := runewidth.StringWidth(s)

// After (uniwidth) - drop-in replacement!
import "github.com/unilibs/uniwidth"
width := uniwidth.StringWidth(s)
```

**Performance improvement**: 3.9-46x faster, zero code changes!

## 📚 Documentation

- [API Reference](https://pkg.go.dev/github.com/unilibs/uniwidth) - Full godoc documentation
- [Benchmark Comparisons](bench/README.md) - Performance comparison vs go-runewidth
- [Architecture Design](docs/ARCHITECTURE.md) - Technical deep dive & design decisions
- [Changelog](CHANGELOG.md) - Version history & upgrade guide

## 🧪 Testing

```bash
# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Run with coverage
go test -cover
```

Current test coverage: 87.1% (target 90%+ for RC)

## 🚀 Development Status

**Current**: v0.1.0-beta (Beta Testing Phase)

> ⚠️ **Beta Notice**: This library is in active beta testing. The API may change before v1.0.0. We encourage testing and feedback, but be prepared for potential breaking changes until we reach Release Candidate (RC) status.

**What Beta Means**:
- ✅ Feature-complete for core functionality
- ✅ Production-quality code and performance
- ⚠️ API may evolve based on community feedback
- ⚠️ Edge cases still being discovered and fixed
- 🎯 Goal: API freeze before v1.0.0-rc

**Completed**:
- ✅ PoC (3 days) - 3.9-46x speedup proven
- ✅ Complete Unicode 16.0 tables - Generated from official data
- ✅ Options API - East Asian Width & emoji configuration
- ✅ Comprehensive testing - 84.6% coverage, fuzzing, conformance tests
- ✅ Bug fixes - Variation selectors, regional indicator flags
- ✅ Documentation - README, ARCHITECTURE, CHANGELOG

**Beta Goals** (Before RC):
- [ ] Community feedback integration
- [ ] Edge case coverage >95%
- [ ] API stability validation
- [ ] Performance regression testing
- [ ] Documentation refinement

**Future Roadmap** (v1.0+):
- [ ] Grapheme cluster support (for complex emoji ZWJ sequences)
- [ ] Additional locale support
- [ ] Extended SIMD optimizations
- [ ] Profile-Guided Optimization (PGO)

## 🤝 Contributing

Contributions welcome! This is part of the [unilibs](https://github.com/unilibs) organization - modern Unicode libraries for Go.

## 📄 License

MIT License - see [LICENSE](LICENSE) file

## 🌟 Related Projects

Built by the [Phoenix TUI Framework](https://github.com/phoenix-tui/phoenix) team.

Part of the **unilibs** ecosystem:
- **uniwidth** - Unicode width calculation (this project)
- **unigrapheme** - Grapheme clustering (planned)
- More Unicode utilities coming soon!

## 📞 Support

- Issues: [GitHub Issues](https://github.com/unilibs/uniwidth/issues)
- Discussions: [GitHub Discussions](https://github.com/unilibs/uniwidth/discussions)

---

## 🙏 Special Thanks

**Professor Ancha Baranova** - This project would not have been possible without her invaluable help and support. Her assistance was crucial in bringing uniwidth to life.

---

**Made with ❤️ by the Phoenix team** | **Powered by Go 1.25+**
