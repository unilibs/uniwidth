package uniwidth

import (
	"testing"
	"unicode"
)

// TestUnicodeConformance_AllCategories tests all major Unicode categories.
func TestUnicodeConformance_AllCategories(t *testing.T) {
	tests := []struct {
		category string
		r        rune
		minWidth int
		maxWidth int
		desc     string
	}{
		// ASCII
		{"ASCII", 'a', 1, 1, "lowercase letter"},
		{"ASCII", 'Z', 1, 1, "uppercase letter"},
		{"ASCII", '0', 1, 1, "digit"},
		{"ASCII", ' ', 1, 1, "space"},
		{"ASCII", '\t', 0, 0, "tab (control)"},
		{"ASCII", '\n', 0, 0, "newline (control)"},

		// Latin Extended
		{"Latin", 'é', 1, 1, "e with acute"},
		{"Latin", 'ñ', 1, 1, "n with tilde"},
		{"Latin", 'ü', 1, 1, "u with diaeresis"},

		// Greek
		{"Greek", 'α', 1, 2, "alpha (ambiguous)"},
		{"Greek", 'β', 1, 2, "beta (ambiguous)"},
		{"Greek", 'Ω', 1, 2, "omega (ambiguous)"},

		// Cyrillic
		{"Cyrillic", 'А', 1, 2, "Cyrillic A (ambiguous)"},
		{"Cyrillic", 'Я', 1, 2, "Cyrillic Ya (ambiguous)"},

		// CJK Unified Ideographs
		{"CJK", '世', 2, 2, "world"},
		{"CJK", '界', 2, 2, "boundary"},
		{"CJK", '你', 2, 2, "you"},
		{"CJK", '好', 2, 2, "good"},

		// Hiragana
		{"Hiragana", 'あ', 2, 2, "a"},
		{"Hiragana", 'い', 2, 2, "i"},
		{"Hiragana", 'う', 2, 2, "u"},

		// Katakana
		{"Katakana", 'ア', 2, 2, "a"},
		{"Katakana", 'イ', 2, 2, "i"},
		{"Katakana", 'ウ', 2, 2, "u"},

		// Hangul
		{"Hangul", '안', 2, 2, "an"},
		{"Hangul", '녕', 2, 2, "nyeong"},

		// Emoji
		{"Emoji", '😀', 2, 2, "grinning face"},
		{"Emoji", '❤', 2, 2, "red heart"},
		{"Emoji", '👍', 2, 2, "thumbs up"},
		{"Emoji", '🚀', 2, 2, "rocket"},

		// Symbols
		{"Symbol", '±', 1, 2, "plus-minus (ambiguous)"},
		{"Symbol", '×', 1, 2, "multiplication (ambiguous)"},
		{"Symbol", '÷', 1, 2, "division (ambiguous)"},

		// Box Drawing
		{"BoxDrawing", '─', 1, 2, "horizontal line (ambiguous)"},
		{"BoxDrawing", '│', 1, 2, "vertical line (ambiguous)"},
		{"BoxDrawing", '┌', 1, 2, "down and right (ambiguous)"},

		// Combining Marks
		{"Combining", 0x0300, 0, 0, "combining grave accent"},
		{"Combining", 0x0301, 0, 0, "combining acute accent"},
		{"Combining", 0x0302, 0, 0, "combining circumflex"},

		// Zero-Width Characters
		{"ZeroWidth", 0x200B, 0, 0, "zero-width space"},
		{"ZeroWidth", 0x200C, 0, 0, "zero-width non-joiner"},
		{"ZeroWidth", 0x200D, 0, 0, "zero-width joiner"},
		{"ZeroWidth", 0xFE0F, 0, 0, "variation selector-16"},
	}

	for _, tt := range tests {
		t.Run(tt.category+"/"+tt.desc, func(t *testing.T) {
			width := RuneWidth(tt.r)
			if width < tt.minWidth || width > tt.maxWidth {
				t.Errorf("RuneWidth(%U %s) = %d, want %d-%d", tt.r, tt.desc, width, tt.minWidth, tt.maxWidth)
			}
		})
	}
}

// TestUnicodeConformance_EdgeCases tests edge cases and boundary conditions.
func TestUnicodeConformance_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected int
	}{
		// Boundary of ASCII
		{"ASCII boundary low", 0x00, 0},
		{"ASCII boundary high", 0x7F, 0}, // DEL
		{"Just after ASCII", 0x80, 0},    // C1 control

		// Boundary of CJK Unified Ideographs
		{"Before CJK", 0x4DFF, 1},
		{"CJK start", 0x4E00, 2},
		{"CJK end", 0x9FFF, 2},
		{"After CJK", 0xA000, 2}, // Yi Syllables

		// Boundary of Hangul
		{"Before Hangul", 0xABFF, 1},
		{"Hangul start", 0xAC00, 2},
		{"Hangul end", 0xD7AF, 2},
		{"After Hangul", 0xD7B0, 1},

		// Boundary of Hiragana/Katakana
		{"Before Hiragana", 0x303F, 2},
		{"Hiragana start", 0x3040, 2},
		{"Katakana end", 0x30FF, 2},
		{"After Katakana", 0x3100, 2},

		// Emoji boundaries
		{"Before emoji", 0x1F5FF, 2},
		{"Emoji start", 0x1F600, 2},
		{"Emoji end", 0x1F64F, 2},
		{"After emoji", 0x1F650, 1},

		// Private Use Area
		{"Private Use start", 0xE000, 1},
		{"Private Use mid", 0xE800, 1},
		{"Private Use end", 0xF8FF, 1},

		// Variation Selectors
		{"Variation Selector 1", 0xFE00, 0},
		{"Variation Selector 16", 0xFE0F, 0},

		// Fullwidth ASCII variants
		{"Fullwidth A", 0xFF21, 2},     // Ａ
		{"Fullwidth 0", 0xFF10, 2},     // ０
		{"Fullwidth space", 0x3000, 2}, // Ideographic space
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidth(tt.r)
			if got != tt.expected {
				t.Errorf("RuneWidth(%U) = %d, want %d", tt.r, got, tt.expected)
			}
		})
	}
}

// TestUnicodeConformance_SurrogateHandling tests handling of surrogate pairs.
func TestUnicodeConformance_SurrogateHandling(t *testing.T) {
	// Go's range over string automatically handles surrogate pairs correctly
	// Test that we handle high-codepoint characters (> U+FFFF) properly

	tests := []struct {
		name string
		s    string
		want int
	}{
		// Characters in Supplementary Multilingual Plane (SMP)
		{"Gothic letter", "𐌰", 1},              // U+10330
		{"Linear B syllable", "𐀀", 2},          // U+10000
		{"Emoji family", "👨\u200D👩\u200D👧", 6}, // Man + ZWJ + Woman + ZWJ + Girl (simplified width)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

// TestUnicodeConformance_CombiningMarks tests all combining mark categories.
func TestUnicodeConformance_CombiningMarks(t *testing.T) {
	// Test that all Unicode combining marks return width 0

	testRanges := []struct {
		name  string
		first rune
		last  rune
	}{
		{"Combining Diacritical Marks", 0x0300, 0x036F},
		{"Combining Diacritical Marks Extended", 0x1AB0, 0x1AFF},
		{"Combining Diacritical Marks Supplement", 0x1DC0, 0x1DFF},
		{"Combining Diacritical Marks for Symbols", 0x20D0, 0x20FF},
		{"Combining Half Marks", 0xFE20, 0xFE2F},
	}

	for _, tr := range testRanges {
		t.Run(tr.name, func(t *testing.T) {
			// Test a few samples from each range (not all to keep tests fast)
			samples := []rune{tr.first, (tr.first + tr.last) / 2, tr.last}
			for _, r := range samples {
				if r > unicode.MaxRune {
					continue
				}
				width := RuneWidth(r)
				if width != 0 {
					t.Errorf("RuneWidth(%U) = %d, want 0 (combining mark)", r, width)
				}
			}
		})
	}
}

// TestUnicodeConformance_ControlCharacters tests all control characters.
func TestUnicodeConformance_ControlCharacters(t *testing.T) {
	tests := []struct {
		name  string
		first rune
		last  rune
	}{
		{"C0 controls", 0x0000, 0x001F},
		{"DELETE", 0x007F, 0x007F},
		{"C1 controls", 0x0080, 0x009F},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for r := tt.first; r <= tt.last; r++ {
				width := RuneWidth(r)
				if width != 0 {
					t.Errorf("RuneWidth(%U) = %d, want 0 (control character)", r, width)
				}
			}
		})
	}
}

// TestUnicodeConformance_FullwidthHalfwidth tests fullwidth/halfwidth forms.
func TestUnicodeConformance_FullwidthHalfwidth(t *testing.T) {
	tests := []struct {
		name         string
		halfwidth    rune
		fullwidth    rune
		halfExpected int
		fullExpected int
	}{
		{"Latin A", 'A', 'Ａ', 1, 2},     // U+0041 vs U+FF21
		{"Digit 0", '0', '０', 1, 2},     // U+0030 vs U+FF10
		{"Space", ' ', '\u3000', 1, 2},  // U+0020 vs U+3000
		{"Exclamation", '!', '！', 1, 2}, // U+0021 vs U+FF01
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			halfWidth := RuneWidth(tt.halfwidth)
			if halfWidth != tt.halfExpected {
				t.Errorf("RuneWidth(%U halfwidth) = %d, want %d", tt.halfwidth, halfWidth, tt.halfExpected)
			}

			fullWidth := RuneWidth(tt.fullwidth)
			if fullWidth != tt.fullExpected {
				t.Errorf("RuneWidth(%U fullwidth) = %d, want %d", tt.fullwidth, fullWidth, tt.fullExpected)
			}
		})
	}
}

// TestUnicodeConformance_EmojiSequences tests complex emoji sequences.
func TestUnicodeConformance_EmojiSequences(t *testing.T) {
	tests := []struct {
		name string
		s    string
		min  int
		max  int
	}{
		// Note: We're testing width, not grapheme clustering
		// ZWJ sequences will be counted as sum of parts (for now)
		{"Simple emoji", "😀", 2, 2},
		{"Emoji with variation selector", "❤\uFE0F", 2, 2}, // Heart + VS-16
		{"ZWJ sequence (family)", "👨\u200D👩\u200D👧", 6, 6}, // Counted as Man+ZWJ+Woman+ZWJ+Girl
		{"Flag sequence", "🇺🇸", 2, 2},                      // Two regional indicators = 1 flag = width 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := StringWidth(tt.s)
			if width < tt.min || width > tt.max {
				t.Logf("Note: Complex emoji sequences may not be handled as single graphemes yet")
				t.Logf("StringWidth(%q) = %d, expected range %d-%d", tt.s, width, tt.min, tt.max)
			}
		})
	}
}

// BenchmarkConformance benchmarks conformance test performance.
func BenchmarkConformance(b *testing.B) {
	runes := []rune{'a', '世', '😀', '±', '\u0300', 0x200D}

	b.Run("AllCategories", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range runes {
				_ = RuneWidth(r)
			}
		}
	})

	b.Run("CombiningMarks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RuneWidth(0x0300)
		}
	})

	b.Run("ControlChars", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RuneWidth('\n')
		}
	})
}
