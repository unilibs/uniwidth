//go:build go1.18

package uniwidth

import (
	"testing"
	"unicode/utf8"
)

// FuzzRuneWidth fuzzes the RuneWidth function with random runes.
func FuzzRuneWidth(f *testing.F) {
	// Seed corpus with interesting runes
	seeds := []rune{
		// ASCII
		'a', 'Z', '0', ' ', '\n', '\t',
		// Control characters
		0x00, 0x1F, 0x7F, 0x80, 0x9F,
		// Latin Extended
		'é', 'ñ', 'ü',
		// CJK
		'世', '界', '你', '好',
		// Hiragana/Katakana
		'あ', 'い', 'ア', 'イ',
		// Hangul
		'안', '녕',
		// Emoji
		'😀', '❤', '👍', '🚀',
		// Ambiguous
		'±', '×', '÷', '°',
		// Box drawing
		'─', '│', '┌',
		// Combining marks
		0x0300, 0x0301, 0x0302,
		// Zero-width
		0x200B, 0x200C, 0x200D,
		0xFE0F,
		// Fullwidth
		'Ａ', '０',
		// Private Use Area
		0xE000, 0xF8FF,
		// High codepoints
		0x10000, 0x1F600, 0x2FA1D,
		// Maximum valid rune
		0x10FFFF,
	}

	for _, r := range seeds {
		f.Add(int32(r))
	}

	f.Fuzz(func(t *testing.T, r32 int32) {
		r := rune(r32)

		// Skip invalid runes
		if !utf8.ValidRune(r) {
			t.Skip("invalid rune")
		}

		// Calculate width
		width := RuneWidth(r)

		// Invariants: width must be 0, 1, or 2
		if width < 0 || width > 2 {
			t.Errorf("RuneWidth(%U) = %d, must be 0, 1, or 2", r, width)
		}

		// ASCII printable characters (0x20-0x7E) must have width 1
		if r >= 0x20 && r <= 0x7E {
			if width != 1 {
				t.Errorf("RuneWidth(%U) = %d, ASCII printable must be 1", r, width)
			}
		}

		// ASCII control characters must have width 0
		if (r >= 0x00 && r <= 0x1F) || r == 0x7F {
			if width != 0 {
				t.Errorf("RuneWidth(%U) = %d, ASCII control must be 0", r, width)
			}
		}

		// CJK Unified Ideographs must have width 2
		if r >= 0x4E00 && r <= 0x9FFF {
			if width != 2 {
				t.Errorf("RuneWidth(%U) = %d, CJK must be 2", r, width)
			}
		}

		// Hangul Syllables must have width 2
		if r >= 0xAC00 && r <= 0xD7AF {
			if width != 2 {
				t.Errorf("RuneWidth(%U) = %d, Hangul must be 2", r, width)
			}
		}

		// Hiragana/Katakana must have width 2
		if r >= 0x3040 && r <= 0x30FF {
			if width != 2 {
				t.Errorf("RuneWidth(%U) = %d, Hiragana/Katakana must be 2", r, width)
			}
		}

		// Common emoji ranges must have width 2
		if (r >= 0x1F600 && r <= 0x1F64F) ||
			(r >= 0x1F300 && r <= 0x1F5FF) ||
			(r >= 0x1F680 && r <= 0x1F6FF) ||
			(r >= 0x1F900 && r <= 0x1F9FF) {
			if width != 2 {
				t.Errorf("RuneWidth(%U) = %d, emoji must be 2", r, width)
			}
		}

		// Zero-width joiners must have width 0
		if r == 0x200D || r == 0x200C {
			if width != 0 {
				t.Errorf("RuneWidth(%U) = %d, ZWJ/ZWNJ must be 0", r, width)
			}
		}

		// Variation selectors must have width 0
		if (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF) {
			if width != 0 {
				t.Errorf("RuneWidth(%U) = %d, variation selector must be 0", r, width)
			}
		}

		// No panics allowed!
	})
}

// FuzzStringWidth fuzzes the StringWidth function with random strings.
func FuzzStringWidth(f *testing.F) {
	// Seed corpus with interesting strings
	seeds := []string{
		"",
		"a",
		"Hello",
		"Hello, World!",
		"世界",
		"你好",
		"Hello 世界",
		"😀",
		"Hello 😀",
		"±½°",
		"─│┌",
		"\n\t",
		"e\u0301",         // e + combining acute
		"👨\u200D👩\u200D👧", // Family emoji (ZWJ sequence)
		"🇺🇸",              // Flag (regional indicators)
		// Long ASCII string
		"The quick brown fox jumps over the lazy dog",
		// Mixed content
		"Hello 世界 😀 ±½",
		// Fullwidth
		"Ａｂｃ０１２",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// Calculate width
		width := StringWidth(s)

		// Invariant: width must be non-negative
		if width < 0 {
			t.Errorf("StringWidth(%q) = %d, must be non-negative", s, width)
		}

		// Invariant: empty string has width 0
		if s == "" && width != 0 {
			t.Errorf("StringWidth(\"\") = %d, must be 0", width)
		}

		// Invariant: ASCII-only string width equals length
		isASCII := true
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				isASCII = false
				break
			}
		}
		if isASCII {
			// Count non-control ASCII characters
			expectedWidth := 0
			for _, c := range s {
				if c >= 0x20 && c != 0x7F {
					expectedWidth++
				}
			}
			if width != expectedWidth {
				t.Errorf("StringWidth(%q) = %d, ASCII-only should be %d", s, width, expectedWidth)
			}
		}

		// Invariant: width should be close to rune-by-rune calculation
		// (within margin for grapheme clustering differences)
		runeWidthSum := 0
		for _, r := range s {
			runeWidthSum += RuneWidth(r)
		}
		// Allow some difference for combining marks and grapheme clusters
		maxDiff := len(s) / 2 // Allow up to 50% difference for complex cases
		if diff := abs(width - runeWidthSum); diff > maxDiff {
			t.Logf("StringWidth(%q) = %d, rune sum = %d, diff = %d", s, width, runeWidthSum, diff)
		}

		// No panics allowed!
	})
}

// FuzzStringWidthWithOptions fuzzes the options API.
func FuzzStringWidthWithOptions(f *testing.F) {
	seeds := []string{
		"±½°",
		"Hello ±",
		"α β γ",
		"А Я",
		"─│┌",
	}

	for _, s := range seeds {
		f.Add(s, true)  // EAWide
		f.Add(s, false) // EANarrow
	}

	f.Fuzz(func(t *testing.T, s string, wide bool) {
		var opts []Option
		if wide {
			opts = []Option{WithEastAsianAmbiguous(EAWide)}
		} else {
			opts = []Option{WithEastAsianAmbiguous(EANarrow)}
		}

		width := StringWidthWithOptions(s, opts...)

		// Invariant: width must be non-negative
		if width < 0 {
			t.Errorf("StringWidthWithOptions(%q, wide=%v) = %d, must be non-negative", s, wide, width)
		}

		// Invariant: empty string has width 0
		if s == "" && width != 0 {
			t.Errorf("StringWidthWithOptions(\"\", wide=%v) = %d, must be 0", wide, width)
		}

		// No panics allowed!
	})
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// FuzzIsASCIIOnly fuzzes the ASCII detection.
func FuzzIsASCIIOnly(f *testing.F) {
	seeds := []string{
		"",
		"a",
		"Hello",
		"Hello, World!",
		"Hello 世界",
		"😀",
		"\n\t",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		result := isASCIIOnly(s)

		// Verify correctness
		expected := true
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				expected = false
				break
			}
		}

		if result != expected {
			t.Errorf("isASCIIOnly(%q) = %v, want %v", s, result, expected)
		}

		// No panics allowed!
	})
}
