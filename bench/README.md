# uniwidth Performance Comparison Benchmarks

This directory contains **three-way performance comparison benchmarks** between Unicode width calculation libraries for Go.

## Libraries Compared

| Library | Approach | Strengths |
|---------|----------|-----------|
| [**uniwidth**](https://github.com/unilibs/uniwidth) | Tiered fast-path lookup | Extreme speed, zero allocs for ASCII |
| [**go-runewidth**](https://github.com/mattn/go-runewidth) | Binary search over tables | Established, widely adopted |
| [**uniseg**](https://github.com/rivo/uniseg) | Grapheme cluster segmentation | Full UAX #29 compliance, ZWJ support |

## Why Separate Module?

**Best Practice**: Benchmark dependencies are isolated from the production library.

- Main `uniwidth` module has **ZERO** dependencies
- Competitor libraries appear only in this benchmark module
- Industry-standard approach (used by fasthttp, gjson, sonic)
- Users can verify performance claims independently

## Running Benchmarks

### All Benchmarks
```bash
cd bench
go test -bench=. -benchmem
```

### Full Results (saved to file)
```bash
cd bench
go test -bench=. -benchmem -count=5 -run=^$ | tee results.txt
```

### Filter by Library
```bash
go test -bench=Uniwidth -benchmem      # uniwidth only
go test -bench=GoRunewidth -benchmem    # go-runewidth only
go test -bench=Uniseg -benchmem         # uniseg only
```

### Filter by Category
```bash
go test -bench=ASCII -benchmem          # ASCII strings
go test -bench=CJK -benchmem            # CJK strings
go test -bench=Emoji -benchmem          # Emoji strings
go test -bench=Mixed -benchmem          # Mixed ASCII + CJK
go test -bench=TUI -benchmem            # Real-world TUI scenarios
go test -bench=Flags -benchmem          # Flag emoji (regional indicators)
go test -bench=ZWJ -benchmem            # ZWJ emoji sequences
go test -bench=Combined -benchmem       # Complex mixed strings
```

## Benchmark Categories

### Core Categories
- **RuneWidth** - Single rune width (uniwidth vs go-runewidth only; uniseg does not expose RuneWidth)
- **ASCII** - Pure ASCII strings (short / medium / long)
- **CJK** - Chinese, Japanese, Korean characters
- **Mixed** - ASCII + CJK combinations
- **Emoji** - Emoji-containing strings

### Real-world Scenarios
- **TUI** - Terminal UI patterns (prompts, table headers, status lines)

### Complex Unicode
- **Flags** - Regional indicator pairs (e.g. `🇺🇸🇩🇪🇯🇵`)
- **ZWJ** - Zero Width Joiner sequences (e.g. `👨‍👩‍👧‍👦`, `👩‍💻`, `🏳️‍🌈`)
- **Combined** - All sequence types mixed in a single string

> **Note**: Width results may differ between libraries for ZWJ sequences. uniseg performs full grapheme cluster segmentation (UAX #29), while uniwidth and go-runewidth use simpler approaches optimized for speed.

## Results

Measured on Intel Core i7-1255U (Windows, amd64). Run benchmarks yourself to get results for your platform.

### RuneWidth (single rune)

uniseg does not expose a public `RuneWidth` function.

| Input | uniwidth | go-runewidth | Speedup |
|-------|----------|--------------|---------|
| ASCII (`'a'`) | 2.1 ns/op | 3.7 ns/op | **1.7x** |
| CJK (`'世'`) | 1.9 ns/op | 37.6 ns/op | **~20x** |
| Emoji (`'😀'`) | 3.3 ns/op | 21.7 ns/op | **~6.5x** |

### StringWidth

| Input | uniwidth | go-runewidth | uniseg | vs go-runewidth | vs uniseg |
|-------|----------|--------------|--------|-----------------|-----------|
| ASCII Short (`"Hello"`) | 9 ns | 107 ns | 165 ns | **12x** | **18x** |
| ASCII Medium (43 chars) | 71 ns | 832 ns | 1,224 ns | **12x** | **17x** |
| ASCII Long (228 chars) | 340 ns | 5,380 ns | 8,058 ns | **16x** | **24x** |
| CJK Short (`"你好世界"`) | 96 ns | 347 ns | 379 ns | **3.6x** | **3.9x** |
| CJK Medium (30 chars) | 1,034 ns | 2,790 ns | 3,424 ns | **2.7x** | **3.3x** |
| Mixed Short | 173 ns | 469 ns | 632 ns | **2.7x** | **3.7x** |
| Mixed Medium | 635 ns | 1,603 ns | 2,172 ns | **2.5x** | **3.4x** |
| Emoji Short | 158 ns | 380 ns | 534 ns | **2.4x** | **3.4x** |
| Emoji Medium | 677 ns | 1,749 ns | 2,259 ns | **2.6x** | **3.3x** |

### Real-world TUI Scenarios

| Input | uniwidth | go-runewidth | uniseg | vs go-runewidth | vs uniseg |
|-------|----------|--------------|--------|-----------------|-----------|
| Prompt (`"❯ Enter command:"`) | 156 ns | 456 ns | 638 ns | **2.9x** | **4.1x** |
| Table Header (box-drawing) | 949 ns | 1,281 ns | 1,708 ns | **1.3x** | **1.8x** |
| Status Line (emoji-rich) | 664 ns | 1,614 ns | 2,174 ns | **2.4x** | **3.3x** |

### Complex Unicode Sequences

| Input | uniwidth | go-runewidth | uniseg | vs go-runewidth | vs uniseg |
|-------|----------|--------------|--------|-----------------|-----------|
| Flags (`🇺🇸🇩🇪🇯🇵🇬🇧🇫🇷`) | 201 ns | 391 ns | 455 ns | **1.9x** | **2.3x** |
| ZWJ (`👨‍👩‍👧‍👦 👩‍💻 🏳️‍🌈`) | 323 ns | 326 ns | 691 ns | **~1x** | **2.1x** |
| Combined (all types) | 374 ns | 505 ns | 1,175 ns | **1.4x** | **3.1x** |

All libraries: **0 allocs/op** for short strings. uniwidth allocates 1 `[]rune` for medium/long Unicode strings (needed for lookahead on variation selectors and regional indicators).

## Structure

```
bench/
├── go.mod               # Separate module with benchmark dependencies
├── go.sum               # Dependency checksums
├── comparison_test.go   # Three-way comparison benchmarks
└── README.md            # This file
```

## Dependencies

This benchmark module depends on:
- `github.com/unilibs/uniwidth` (parent module, via replace directive)
- `github.com/mattn/go-runewidth` (comparison baseline)
- `github.com/rivo/uniseg` (comparison baseline)

The main `uniwidth` library has **ZERO** dependencies. These exist only in this benchmark module.

## See Also

- [Main Documentation](../README.md)
- [Architecture Guide](../docs/ARCHITECTURE.md)
- [Changelog](../CHANGELOG.md)
