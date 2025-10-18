package bench_test

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/unilibs/uniwidth"
)

// ============================================================================
// Comparison Benchmarks: uniwidth vs go-runewidth
//
// This package contains performance comparison benchmarks between uniwidth
// and the go-runewidth library. These benchmarks demonstrate the 3.9-46x
// performance improvement achieved by uniwidth's tiered lookup strategy.
//
// Run comparison benchmarks:
//   cd bench
//   go test -bench=. -benchmem
//
// Compare results:
//   go test -bench=. -benchmem | tee results.txt
// ============================================================================

// ============================================================================
// RuneWidth Benchmarks
// ============================================================================

func BenchmarkRuneWidth_ASCII_Uniwidth(b *testing.B) {
	r := 'a'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_ASCII_GoRunewidth(b *testing.B) {
	r := 'a'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_CJK_Uniwidth(b *testing.B) {
	r := '世' // Chinese character
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_CJK_GoRunewidth(b *testing.B) {
	r := '世'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_Emoji_Uniwidth(b *testing.B) {
	r := '😀' // Smiling face
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_Emoji_GoRunewidth(b *testing.B) {
	r := '😀'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.RuneWidth(r)
	}
}

// ============================================================================
// StringWidth Benchmarks - ASCII
// ============================================================================

func BenchmarkStringWidth_ASCII_Short_Uniwidth(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Short_GoRunewidth(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium_Uniwidth(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium_GoRunewidth(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long_Uniwidth(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long_GoRunewidth(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - CJK
// ============================================================================

func BenchmarkStringWidth_CJK_Short_Uniwidth(b *testing.B) {
	s := "你好世界" // Hello World in Chinese
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Short_GoRunewidth(b *testing.B) {
	s := "你好世界"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium_Uniwidth(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium_GoRunewidth(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - Mixed
// ============================================================================

func BenchmarkStringWidth_Mixed_Short_Uniwidth(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Short_GoRunewidth(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium_Uniwidth(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium_GoRunewidth(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - Emoji
// ============================================================================

func BenchmarkStringWidth_Emoji_Short_Uniwidth(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Short_GoRunewidth(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium_Uniwidth(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium_GoRunewidth(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

// ============================================================================
// Real-world TUI Scenarios
// ============================================================================

func BenchmarkStringWidth_TUI_Prompt_Uniwidth(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_Prompt_GoRunewidth(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_TableHeader_Uniwidth(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_TableHeader_GoRunewidth(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_StatusLine_Uniwidth(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_StatusLine_GoRunewidth(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runewidth.StringWidth(s)
	}
}
