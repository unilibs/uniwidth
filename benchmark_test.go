package uniwidth

import (
	"testing"
)

// ============================================================================
// Benchmark: RuneWidth - Single Rune Performance
// ============================================================================

func BenchmarkRuneWidth_ASCII(b *testing.B) {
	r := 'a'
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RuneWidth(r)
	}
}

func BenchmarkRuneWidth_CJK(b *testing.B) {
	r := '世' // Chinese character
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RuneWidth(r)
	}
}

func BenchmarkRuneWidth_Emoji(b *testing.B) {
	r := '😀' // Smiling face
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RuneWidth(r)
	}
}

// ============================================================================
// Benchmark: StringWidth - Full String Performance
// ============================================================================

// ASCII-only strings (most common case ~95% of content)
func BenchmarkStringWidth_ASCII_Short(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// CJK strings (common in Asian locales)
func BenchmarkStringWidth_CJK_Short(b *testing.B) {
	s := "你好世界" // Hello World in Chinese
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// Mixed ASCII + CJK (typical TUI content)
func BenchmarkStringWidth_Mixed_Short(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// Emoji strings (growing usage in modern terminals)
func BenchmarkStringWidth_Emoji_Short(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// ZWJ emoji sequences
func BenchmarkStringWidth_ZWJ_Family(b *testing.B) {
	s := "👨\u200D👩\u200D👧\u200D👦" // Family emoji
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_ZWJ_CoupleHeart(b *testing.B) {
	s := "👩\u200D\u2764\uFE0F\u200D👨" // Couple with heart
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_EmojiModifier(b *testing.B) {
	s := "👍🏽" // Thumbs up with skin tone
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

func BenchmarkStringWidth_ZWJ_Mixed(b *testing.B) {
	s := "Hello 👨\u200D👩\u200D👧 World 👍🏽 Test 🏳\uFE0F\u200D🌈"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// ============================================================================
// Benchmark: isASCIIOnly - Fast Path Detection
// ============================================================================

func BenchmarkIsASCIIOnly_Short_ASCII(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isASCIIOnly(s)
	}
}

func BenchmarkIsASCIIOnly_Medium_ASCII(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isASCIIOnly(s)
	}
}

func BenchmarkIsASCIIOnly_Long_ASCII(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isASCIIOnly(s)
	}
}

func BenchmarkIsASCIIOnly_Short_NonASCII(b *testing.B) {
	s := "Hello 世界"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isASCIIOnly(s)
	}
}

// ============================================================================
// Benchmark: Real-world TUI scenarios
// ============================================================================

// Typical TUI prompt
func BenchmarkStringWidth_TUI_Prompt(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// Typical TUI table header
func BenchmarkStringWidth_TUI_TableHeader(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}

// Typical TUI status line with emojis
func BenchmarkStringWidth_TUI_StatusLine(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringWidth(s)
	}
}
