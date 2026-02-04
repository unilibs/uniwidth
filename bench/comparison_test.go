package bench_test

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
	"github.com/unilibs/uniwidth"
)

// ============================================================================
// Comparison Benchmarks: uniwidth vs go-runewidth vs uniseg
//
// Three-way performance comparison between Unicode width calculation libraries:
//   - uniwidth:      Tiered fast-path lookup (O(1) for common characters)
//   - go-runewidth:  Binary search over Unicode tables
//   - uniseg:        Grapheme cluster segmentation with width calculation
//
// Run all comparison benchmarks:
//
//	cd bench
//	go test -bench=. -benchmem
//
// Filter by library:
//
//	go test -bench=Uniwidth -benchmem
//	go test -bench=GoRunewidth -benchmem
//	go test -bench=Uniseg -benchmem
//
// Filter by category:
//
//	go test -bench=ASCII -benchmem
//	go test -bench=CJK -benchmem
//	go test -bench=Emoji -benchmem
//	go test -bench=TUI -benchmem
//	go test -bench=Flags -benchmem
//	go test -bench=ZWJ -benchmem
//
// ============================================================================

// ============================================================================
// RuneWidth Benchmarks
//
// Note: uniseg does not expose a public RuneWidth function; it operates on
// grapheme clusters via StringWidth and iterator APIs. RuneWidth comparison
// is limited to uniwidth vs go-runewidth.
// ============================================================================

func BenchmarkRuneWidth_ASCII_Uniwidth(b *testing.B) {
	r := 'a'
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_ASCII_GoRunewidth(b *testing.B) {
	r := 'a'
	b.ResetTimer()
	for range b.N {
		_ = runewidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_CJK_Uniwidth(b *testing.B) {
	r := '世' // Chinese character
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_CJK_GoRunewidth(b *testing.B) {
	r := '世'
	b.ResetTimer()
	for range b.N {
		_ = runewidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_Emoji_Uniwidth(b *testing.B) {
	r := '😀' // Smiling face
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.RuneWidth(r)
	}
}

func BenchmarkRuneWidth_Emoji_GoRunewidth(b *testing.B) {
	r := '😀'
	b.ResetTimer()
	for range b.N {
		_ = runewidth.RuneWidth(r)
	}
}

// ============================================================================
// StringWidth Benchmarks - ASCII
// ============================================================================

func BenchmarkStringWidth_ASCII_Short_Uniwidth(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Short_GoRunewidth(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Short_Uniseg(b *testing.B) {
	s := "Hello"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium_Uniwidth(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium_GoRunewidth(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Medium_Uniseg(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long_Uniwidth(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long_GoRunewidth(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ASCII_Long_Uniseg(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - CJK
// ============================================================================

func BenchmarkStringWidth_CJK_Short_Uniwidth(b *testing.B) {
	s := "你好世界" // Hello World in Chinese
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Short_GoRunewidth(b *testing.B) {
	s := "你好世界"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Short_Uniseg(b *testing.B) {
	s := "你好世界"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium_Uniwidth(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium_GoRunewidth(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_CJK_Medium_Uniseg(b *testing.B) {
	s := "これは日本語のテキストです。漢字とひらがなとカタカナが含まれています。"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - Mixed (ASCII + CJK)
// ============================================================================

func BenchmarkStringWidth_Mixed_Short_Uniwidth(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Short_GoRunewidth(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Short_Uniseg(b *testing.B) {
	s := "Hello 世界 World"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium_Uniwidth(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium_GoRunewidth(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Mixed_Medium_Uniseg(b *testing.B) {
	s := "User: John Doe (管理者) | Status: Active | 日本語対応"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ============================================================================
// StringWidth Benchmarks - Emoji
// ============================================================================

func BenchmarkStringWidth_Emoji_Short_Uniwidth(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Short_GoRunewidth(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Short_Uniseg(b *testing.B) {
	s := "Hello 👋 World 😀"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium_Uniwidth(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium_GoRunewidth(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Emoji_Medium_Uniseg(b *testing.B) {
	s := "Status: ✅ Success | Error: ❌ Failed | Progress: 🚀 Loading..."
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ============================================================================
// Real-world TUI Scenarios
// ============================================================================

func BenchmarkStringWidth_TUI_Prompt_Uniwidth(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_Prompt_GoRunewidth(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_Prompt_Uniseg(b *testing.B) {
	s := "❯ Enter command:"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_TableHeader_Uniwidth(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_TableHeader_GoRunewidth(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_TableHeader_Uniseg(b *testing.B) {
	s := "│ ID │ Name │ Status │ Created At │"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_StatusLine_Uniwidth(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_StatusLine_GoRunewidth(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_TUI_StatusLine_Uniseg(b *testing.B) {
	s := "✅ 12 passed | ❌ 3 failed | ⏭️  5 skipped | ⏱️  1.234s"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ============================================================================
// Complex Unicode Sequences
//
// These benchmarks test handling of multi-codepoint sequences that require
// context-aware processing:
//   - Flag emoji (regional indicator pairs)
//   - ZWJ sequences (family emoji, profession emoji)
//   - Combined complex strings mixing all sequence types
//
// Width results may differ between libraries for complex sequences.
// uniseg performs full grapheme cluster segmentation (UAX #29) which produces
// the most accurate results for ZWJ sequences. uniwidth and go-runewidth use
// simpler per-rune or limited lookahead approaches optimized for speed.
// ============================================================================

// Flag emoji: regional indicator pairs forming country flags.
// Each flag is two regional indicator codepoints rendered as a single glyph.
func BenchmarkStringWidth_Flags_Uniwidth(b *testing.B) {
	s := "🇺🇸🇩🇪🇯🇵🇬🇧🇫🇷"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Flags_GoRunewidth(b *testing.B) {
	s := "🇺🇸🇩🇪🇯🇵🇬🇧🇫🇷"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Flags_Uniseg(b *testing.B) {
	s := "🇺🇸🇩🇪🇯🇵🇬🇧🇫🇷"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// ZWJ sequences: emoji composed with Zero Width Joiner (U+200D).
// These form complex glyphs like family groups and gendered professions.
func BenchmarkStringWidth_ZWJ_Uniwidth(b *testing.B) {
	s := "👨‍👩‍👧‍👦 👩‍💻 🏳️‍🌈"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ZWJ_GoRunewidth(b *testing.B) {
	s := "👨‍👩‍👧‍👦 👩‍💻 🏳️‍🌈"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_ZWJ_Uniseg(b *testing.B) {
	s := "👨‍👩‍👧‍👦 👩‍💻 🏳️‍🌈"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}

// Combined: real-world complex string mixing ASCII, CJK, flags, and ZWJ.
// Represents a realistic worst-case scenario for width calculation.
func BenchmarkStringWidth_Combined_Uniwidth(b *testing.B) {
	s := "Hello 🇺🇸 世界 👨‍👩‍👧‍👦 café"
	b.ResetTimer()
	for range b.N {
		_ = uniwidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Combined_GoRunewidth(b *testing.B) {
	s := "Hello 🇺🇸 世界 👨‍👩‍👧‍👦 café"
	b.ResetTimer()
	for range b.N {
		_ = runewidth.StringWidth(s)
	}
}

func BenchmarkStringWidth_Combined_Uniseg(b *testing.B) {
	s := "Hello 🇺🇸 世界 👨‍👩‍👧‍👦 café"
	b.ResetTimer()
	for range b.N {
		_ = uniseg.StringWidth(s)
	}
}
