# uniwidth Performance Comparison Benchmarks

This directory contains **performance comparison benchmarks** between `uniwidth` and `go-runewidth`.

## 🎯 Purpose

Demonstrate the **3.9-46x performance improvement** achieved by uniwidth's tiered lookup strategy compared to go-runewidth's traditional binary search approach.

## 📦 Why Separate Module?

**Best Practice 2025**: Keep benchmark dependencies separate from the production library.

**Benefits**:
- ✅ **Main module**: ZERO dependencies
- ✅ **Clean go.mod**: Users don't see competitor library
- ✅ **Professional**: Industry-standard approach (fasthttp, gjson, sonic)
- ✅ **Optional**: Comparison benchmarks are not required for library usage

## 🚀 Running Benchmarks

### Quick Comparison
```bash
cd bench
go test -bench=. -benchmem
```

### Full Results
```bash
cd bench
go test -bench=. -benchmem -run=^$ | tee results.txt
```

### Compare Specific Categories
```bash
# ASCII strings
go test -bench=ASCII -benchmem

# CJK strings
go test -bench=CJK -benchmem

# Emoji strings
go test -bench=Emoji -benchmem

# Real-world TUI scenarios
go test -bench=TUI -benchmem
```

## 📊 Expected Results

**ASCII Strings** (15-46x faster):
```
BenchmarkStringWidth_ASCII_Short_Uniwidth      149590729    9.5 ns/op    0 B/op   0 allocs/op
BenchmarkStringWidth_ASCII_Short_GoRunewidth    10065044  150.1 ns/op    0 B/op   0 allocs/op
                                                            ^^^^^^^^^^
                                                            15.8x faster!
```

**CJK Strings** (4-14x faster):
```
BenchmarkStringWidth_CJK_Short_Uniwidth         19064941   63.6 ns/op    0 B/op   0 allocs/op
BenchmarkStringWidth_CJK_Short_GoRunewidth       2771077  368.0 ns/op    0 B/op   0 allocs/op
                                                            ^^^^^^^^^^
                                                            5.8x faster!
```

**Emoji Strings** (6-8x faster):
```
BenchmarkStringWidth_Emoji_Short_Uniwidth       12384722   96.2 ns/op    0 B/op   0 allocs/op
BenchmarkStringWidth_Emoji_Short_GoRunewidth     1854066  646.8 ns/op    0 B/op   0 allocs/op
                                                            ^^^^^^^^^^
                                                            6.7x faster!
```

## 📁 Structure

```
bench/
├── go.mod               # Separate module with go-runewidth dependency
├── go.sum               # Dependencies checksums
├── comparison_test.go   # Comparison benchmarks (uniwidth vs go-runewidth)
└── README.md            # This file
```

## 🔗 Dependencies

This module depends on:
- `github.com/unilibs/uniwidth` (parent module, via replace directive)
- `github.com/mattn/go-runewidth` (competitor, for comparison only)

**Note**: The main `uniwidth` library has ZERO dependencies. These dependencies exist only in this benchmark module for performance comparison purposes.

## 📝 Notes

- Benchmarks are isolated from the main library
- Main `uniwidth` module remains dependency-free
- Comparison benchmarks prove performance claims (marketing)
- Users can verify performance independently

## 🎓 Learn More

- [Main Documentation](../README.md)
- [Architecture Guide](../docs/ARCHITECTURE.md)
- [PoC Results](../docs/POC_RESULTS.md)

---

*These benchmarks demonstrate why uniwidth is 3.9-46x faster than go-runewidth.*
