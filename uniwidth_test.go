package uniwidth

import (
	"testing"
)

func TestRuneWidth_ASCII(t *testing.T) {
	tests := []struct {
		r    rune
		want int
	}{
		// Printable ASCII
		{'a', 1},
		{'A', 1},
		{'0', 1},
		{'!', 1},
		{' ', 1},
		{'~', 1},

		// Control characters
		{'\n', 0},
		{'\t', 0},
		{'\r', 0},
		{0x00, 0}, // NUL
		{0x1F, 0}, // US (Unit Separator)
		{0x7F, 0}, // DEL
	}

	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			got := RuneWidth(tt.r)
			if got != tt.want {
				t.Errorf("RuneWidth(%U) = %d, want %d", tt.r, got, tt.want)
			}
		})
	}
}

func TestRuneWidth_CJK(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		// CJK Unified Ideographs
		{"Chinese 世", '世', 2},
		{"Chinese 界", '界', 2},
		{"Chinese 你", '你', 2},
		{"Chinese 好", '好', 2},

		// Hiragana
		{"Hiragana あ", 'あ', 2},
		{"Hiragana い", 'い', 2},

		// Katakana
		{"Katakana ア", 'ア', 2},
		{"Katakana イ", 'イ', 2},

		// Hangul
		{"Korean 안", '안', 2},
		{"Korean 녕", '녕', 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidth(tt.r)
			if got != tt.want {
				t.Errorf("RuneWidth(%U %s) = %d, want %d", tt.r, tt.name, got, tt.want)
			}
		})
	}
}

func TestRuneWidth_Emoji(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		// Simple emoji
		{"Smiling face 😀", '😀', 2},
		{"Heart ❤", '❤', 2},
		{"Thumbs up 👍", '👍', 2},
		{"Wave 👋", '👋', 2},

		// Weather/symbols
		{"Sun ☀", '☀', 2},
		{"Cloud ☁", '☁', 2},

		// Transport
		{"Rocket 🚀", '🚀', 2},
		{"Car 🚗", '🚗', 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidth(tt.r)
			if got != tt.want {
				t.Errorf("RuneWidth(%U %s) = %d, want %d", tt.r, tt.name, got, tt.want)
			}
		})
	}
}

func TestRuneWidth_ZeroWidth(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		// Zero-width joiners
		{"ZWJ", 0x200D, 0},
		{"ZWNJ", 0x200C, 0},

		// Variation selectors
		{"Variation selector", 0xFE0F, 0},

		// Some combining marks (unicode.Mn category handled separately)
		{"Combining acute accent", 0x0301, 0},
		{"Combining grave accent", 0x0300, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidth(tt.r)
			if got != tt.want {
				t.Errorf("RuneWidth(%U %s) = %d, want %d", tt.r, tt.name, got, tt.want)
			}
		})
	}
}

func TestStringWidth_ASCII(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"Hello", 5},
		{"Hello, World!", 13},
		{"12345", 5},
		{"ASCII only content", 18},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestStringWidth_Mixed(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		// ASCII + CJK
		{"Hello 世界", 10}, // "Hello " (6) + "世界" (4) = 10
		{"你好", 4},        // 你(2) 好(2)

		// ASCII + Emoji
		{"Hello 👋", 8}, // "Hello " (6) + 👋 (2)
		{"Test 😀", 7},  // "Test " (5) + 😀 (2)

		// CJK + Emoji
		{"世界 👋", 7}, // 世(2) 界(2) space(1) 👋(2)
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestStringWidth_VariationSelectors(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Text variation selector (U+FE0E) - forces narrow width
		{
			name: "Sun with text variation",
			s:    "☀︎", // U+2600 + U+FE0E
			want: 1,    // Text presentation = width 1
		},
		// Emoji variation selector (U+FE0F) - forces wide width
		{
			name: "Sun with emoji variation",
			s:    "☀️", // U+2600 + U+FE0F
			want: 2,    // Emoji presentation = width 2
		},
		// Shield with emoji variation
		{
			name: "Shield with emoji variation",
			s:    "🛡️", // U+1F6E1 + U+FE0F
			want: 2,    // Emoji presentation = width 2
		},
		// No variation selector
		{
			name: "Clock (no variation selector)",
			s:    "⏰", // U+23F0
			want: 2,   // Default width 2
		},
		// Heart with variation selector
		{
			name: "Heart with emoji variation",
			s:    "❤️", // U+2764 + U+FE0F
			want: 2,    // Emoji presentation = width 2
		},
		// Multiple characters with variation selectors
		{
			name: "Multiple with variations",
			s:    "☀︎❤️", // U+2600+U+FE0E + U+2764+U+FE0F
			want: 3,      // 1 (text sun) + 2 (emoji heart)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
				// Debug output
				t.Logf("Runes: %U", []rune(tt.s))
			}
		})
	}
}

func TestStringWidth_RegionalIndicators(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Country flags (2 regional indicators = 1 flag)
		{
			name: "US flag",
			s:    "🇺🇸", // U+1F1FA + U+1F1F8
			want: 2,    // Flag = width 2 (not 4!)
		},
		{
			name: "Japan flag",
			s:    "🇯🇵", // U+1F1EF + U+1F1F5
			want: 2,    // Flag = width 2
		},
		{
			name: "UK flag",
			s:    "🇬🇧", // U+1F1EC + U+1F1E7
			want: 2,    // Flag = width 2
		},
		// Multiple flags
		{
			name: "Two flags",
			s:    "🇺🇸🇯🇵", // US + Japan
			want: 4,      // 2 + 2
		},
		// Flag with other emoji
		{
			name: "Flag with emoji",
			s:    "🇺🇸👋", // US flag + wave
			want: 4,     // 2 (flag) + 2 (wave)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
				// Debug output
				t.Logf("Runes: %U", []rune(tt.s))
			}
		})
	}
}

func TestIsRegionalIndicator(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"First regional indicator (A)", 0x1F1E6, true},
		{"Last regional indicator (Z)", 0x1F1FF, true},
		{"Middle regional indicator (U)", 0x1F1FA, true},
		{"Before range", 0x1F1E5, false},
		{"After range", 0x1F200, false},
		{"Regular emoji", '😀', false},
		{"ASCII", 'A', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRegionalIndicator(tt.r)
			if got != tt.want {
				t.Errorf("isRegionalIndicator(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestIsASCIIOnly(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"a", true},
		{"Hello", true},
		{"Hello, World!", true},
		{"12345", true},
		{"Hello 世界", false},
		{"你好", false},
		{"Hello 👋", false},
		{"Test\n\t", true}, // Control chars are still ASCII
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := isASCIIOnly(tt.s)
			if got != tt.want {
				t.Errorf("isASCIIOnly(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
