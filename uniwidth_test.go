package uniwidth

import (
	"testing"
	"unicode"
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

// runeWidthViaBinarySearch computes the full RuneWidth using the legacy
// binary search path (Tier 1-3 hot paths + binary search fallback).
// This is a reference implementation for verifying the table lookup.
func runeWidthViaBinarySearch(r rune) int {
	// Tier 1: ASCII
	if r < 0x80 {
		if r < 0x20 {
			return 0
		}
		if r == 0x7F {
			return 0
		}
		return 1
	}

	// Tier 2: CJK
	if r >= 0x4E00 && r <= 0x9FFF {
		return 2
	}
	if r >= 0xAC00 && r <= 0xD7AF {
		return 2
	}
	if r >= 0x3040 && r <= 0x312F {
		return 2
	}
	if r >= 0xF900 && r <= 0xFAFF {
		return 2
	}

	// Tier 3: Emoji
	if r >= 0x1F600 && r <= 0x1F64F {
		return 2
	}
	if r >= 0x1F300 && r <= 0x1F5FF {
		return 2
	}
	if r >= 0x1F680 && r <= 0x1F6FF {
		return 2
	}
	if r >= 0x1F900 && r <= 0x1F9FF {
		return 2
	}
	if r >= 0x2600 && r <= 0x26FF {
		return 2
	}
	if r >= 0x2700 && r <= 0x27BF {
		return 2
	}

	// Zero-width format characters (ZWSP, ZWNJ, ZWJ, LRM, RLM)
	if r >= 0x200B && r <= 0x200F {
		return 0
	}
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}
	if r >= 0xE0100 && r <= 0xE01EF {
		return 0
	}

	// Combining marks (same as RuneWidth uses unicode.In)
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
		return 0
	}

	// Tier 4: Legacy binary search
	return binarySearchWidth(r)
}

// TestTableLookup_ExhaustiveVerification iterates ALL valid Unicode codepoints
// (0x0000-0x10FFFF, skipping surrogates 0xD800-0xDFFF) and verifies that
// RuneWidth (which uses tableLookupWidth in Tier 4) returns the same result
// as the reference implementation using binarySearchWidth in Tier 4.
//
// This ensures the multi-stage table produces identical results to the legacy
// binary search tables when called through the full RuneWidth path.
func TestTableLookup_ExhaustiveVerification(t *testing.T) {
	mismatches := 0
	const maxMismatchLog = 20

	for cp := rune(0); cp <= 0x10FFFF; cp++ {
		// Skip surrogates (not valid Unicode scalar values)
		if cp >= 0xD800 && cp <= 0xDFFF {
			continue
		}

		tableW := RuneWidth(cp)                 // uses tableLookupWidth in Tier 4
		binaryW := runeWidthViaBinarySearch(cp) // uses binarySearchWidth in Tier 4

		if tableW != binaryW {
			mismatches++
			if mismatches <= maxMismatchLog {
				t.Errorf("U+%04X: RuneWidth(table)=%d, runeWidthViaBinarySearch=%d", cp, tableW, binaryW)
			}
		}
	}

	if mismatches > maxMismatchLog {
		t.Errorf("... and %d more mismatches (total: %d)", mismatches-maxMismatchLog, mismatches)
	}

	if mismatches == 0 {
		t.Logf("Verified %d codepoints: RuneWidth matches reference implementation for all", 0x10FFFF+1-(0xDFFF-0xD800+1))
	}
}

// TestTableLookupInternal_ExhaustiveVerification verifies that the internal
// table lookup (used by Options API) matches the legacy binary search internal
// for ALL codepoints that reach Tier 4 (after Tier 1-3 hot paths).
func TestTableLookupInternal_ExhaustiveVerification(t *testing.T) {
	mismatches := 0
	const maxMismatchLog = 20

	for cp := rune(0); cp <= 0x10FFFF; cp++ {
		// Skip surrogates
		if cp >= 0xD800 && cp <= 0xDFFF {
			continue
		}

		// Compare the full runeWidthInternal path (which uses tableLookupWidthInternal)
		// against a reference that uses binarySearchWidthInternal.
		// runeWidthInternal handles Tier 1-3 and zero-width checks before Tier 4,
		// so we test the full path for consistency.
		tableW := runeWidthInternal(cp) // uses tableLookupWidthInternal in Tier 4

		// Reference: replicate runeWidthInternal logic but with binary search
		binaryW := runeWidthInternalViaBinarySearch(cp)

		if tableW != binaryW {
			mismatches++
			if mismatches <= maxMismatchLog {
				t.Errorf("U+%04X: runeWidthInternal(table)=%d, runeWidthInternalViaBinarySearch=%d", cp, tableW, binaryW)
			}
		}
	}

	if mismatches > maxMismatchLog {
		t.Errorf("... and %d more mismatches (total: %d)", mismatches-maxMismatchLog, mismatches)
	}

	if mismatches == 0 {
		t.Logf("Verified %d codepoints: runeWidthInternal matches reference for all", 0x10FFFF+1-(0xDFFF-0xD800+1))
	}
}

// runeWidthInternalViaBinarySearch is a reference implementation using binary search
// for verifying the table-based runeWidthInternal.
func runeWidthInternalViaBinarySearch(r rune) int {
	// Tier 1: ASCII
	if r < 0x80 {
		if r < 0x20 {
			return 0
		}
		if r == 0x7F {
			return 0
		}
		return 1
	}

	// Tier 2: CJK
	if r >= 0x4E00 && r <= 0x9FFF {
		return 2
	}
	if r >= 0xAC00 && r <= 0xD7AF {
		return 2
	}
	if r >= 0x3040 && r <= 0x30FF {
		return 2
	}
	if r >= 0xF900 && r <= 0xFAFF {
		return 2
	}

	// Tier 3: Emoji
	if r >= 0x1F600 && r <= 0x1F64F {
		return 2
	}
	if r >= 0x1F300 && r <= 0x1F5FF {
		return 2
	}
	if r >= 0x1F680 && r <= 0x1F6FF {
		return 2
	}
	if r >= 0x1F900 && r <= 0x1F9FF {
		return 2
	}
	if r >= 0x2600 && r <= 0x26FF {
		return 2
	}
	if r >= 0x2700 && r <= 0x27BF {
		return 2
	}

	// Zero-width format characters (ZWSP, ZWNJ, ZWJ, LRM, RLM)
	if r >= 0x200B && r <= 0x200F {
		return 0
	}
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}
	if r >= 0xE0100 && r <= 0xE01EF {
		return 0
	}

	// Combining marks
	if (r >= 0x0300 && r <= 0x036F) ||
		(r >= 0x1AB0 && r <= 0x1AFF) ||
		(r >= 0x1DC0 && r <= 0x1DFF) ||
		(r >= 0x20D0 && r <= 0x20FF) ||
		(r >= 0xFE20 && r <= 0xFE2F) {
		return 0
	}

	// Tier 4: Legacy binary search
	return binarySearchWidthInternal(r)
}

// TestTableLookup_SpecificCodepoints tests the table lookup for specific
// important codepoints to ensure correctness of the 2-bit encoding.
func TestTableLookup_SpecificCodepoints(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		// Width 0: control characters
		{"NUL", 0x0000, 0},
		{"TAB", 0x0009, 0},
		{"LF", 0x000A, 0},
		{"DEL", 0x007F, 0},
		{"C1 control", 0x0080, 0},
		{"Soft hyphen", 0x00AD, 0},

		// Width 0: combining marks
		{"Combining grave", 0x0300, 0},
		{"Combining acute", 0x0301, 0},
		{"Combining marks extended", 0x1AB0, 0},
		{"Combining marks extended end", 0x1AFF, 0},
		{"Combining marks supplement", 0x1DC0, 0},

		// Width 0: zero-width characters
		{"ZWSP", 0x200B, 0},
		{"ZWNJ", 0x200C, 0},
		{"ZWJ", 0x200D, 0},
		{"Variation selector 1", 0xFE00, 0},
		{"Variation selector 16", 0xFE0F, 0},
		{"BOM", 0xFEFF, 0},

		// Width 1: ASCII printable
		{"Space", 0x0020, 1},
		{"Letter A", 0x0041, 1},
		{"Tilde", 0x007E, 1},

		// Width 1: Latin extended
		{"e-acute", 0x00E9, 1},

		// Width 2: CJK
		{"CJK ideograph", 0x4E00, 2},
		{"Hangul syllable", 0xAC00, 2},
		{"Hiragana A", 0x3042, 2},
		{"Katakana A", 0x30A2, 2},

		// Width 2: Emoji
		{"Grinning face", 0x1F600, 2},
		{"Rocket", 0x1F680, 2},
		{"Sun", 0x2600, 2},

		// Width 2: Fullwidth
		{"Fullwidth A", 0xFF21, 2},
		{"Fullwidth 0", 0xFF10, 2},
		{"Ideographic space", 0x3000, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tableLookupWidth(tt.r)
			if got != tt.want {
				t.Errorf("tableLookupWidth(%U) = %d, want %d", tt.r, got, tt.want)
			}
		})
	}
}

func TestStringWidth_ZWJSequences(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Family ZWJ sequences
		{
			name: "Family: man+woman+girl+boy",
			s:    "👨\u200D👩\u200D👧\u200D👦", // 👨‍👩‍👧‍👦
			want: 2,
		},
		{
			name: "Family: man+woman+girl",
			s:    "👨\u200D👩\u200D👧", // 👨‍👩‍👧
			want: 2,
		},
		{
			name: "Couple with heart",
			s:    "👩\u200D\u2764\uFE0F\u200D👨", // 👩‍❤️‍👨
			want: 2,
		},
		{
			name: "Kiss: woman+man",
			s:    "👩\u200D\u2764\uFE0F\u200D\U0001F48B\u200D👨",
			want: 2,
		},
		// Profession ZWJ sequences
		{
			name: "Woman scientist",
			s:    "👩\u200D🔬", // 👩‍🔬
			want: 2,
		},
		{
			name: "Man firefighter",
			s:    "👨\u200D🚒", // 👨‍🚒
			want: 2,
		},
		{
			name: "Woman technologist",
			s:    "👩\u200D💻", // 👩‍💻
			want: 2,
		},
		// Gendered ZWJ sequences
		{
			name: "Man with probing cane",
			s:    "👨\u200D🦯", // 👨‍🦯
			want: 2,
		},
		// Heart sequences
		{
			name: "Heart on fire",
			s:    "\u2764\uFE0F\u200D🔥", // ❤️‍🔥
			want: 2,
		},
		{
			name: "Mending heart",
			s:    "\u2764\uFE0F\u200D\U0001FA79", // ❤️‍🩹
			want: 2,
		},
		// Rainbow flag
		{
			name: "Rainbow flag",
			s:    "🏳\uFE0F\u200D🌈", // 🏳️‍🌈
			want: 2,
		},
		// Transgender flag
		{
			name: "Transgender flag",
			s:    "🏳\uFE0F\u200D\u26A7\uFE0F", // 🏳️‍⚧️
			want: 2,
		},
		// Pirate flag
		{
			name: "Pirate flag",
			s:    "🏴\u200D\u2620\uFE0F", // 🏴‍☠️
			want: 2,
		},
		// Multiple ZWJ emoji in a string
		{
			name: "Multiple ZWJ sequences",
			s:    "👨\u200D👩\u200D👧 and 👩\u200D💻",
			want: 9, // family(2) + " and "(5) + technologist(2)
		},
		// ZWJ in mixed content
		{
			name: "Mixed: ASCII + ZWJ family",
			s:    "Family: 👨\u200D👩\u200D👧\u200D👦!",
			want: 11, // "Family: "(8) + family(2) + "!"(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
				t.Logf("Runes: %U", []rune(tt.s))
			}
		})
	}
}

func TestStringWidth_EmojiModifiers(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Skin tone modifiers
		{
			name: "Thumbs up + light skin",
			s:    "👍🏻", // U+1F44D + U+1F3FB
			want: 2,
		},
		{
			name: "Thumbs up + medium skin",
			s:    "👍🏽", // U+1F44D + U+1F3FD
			want: 2,
		},
		{
			name: "Thumbs up + dark skin",
			s:    "👍🏿", // U+1F44D + U+1F3FF
			want: 2,
		},
		{
			name: "Wave + medium-light skin",
			s:    "👋🏼", // U+1F44B + U+1F3FC
			want: 2,
		},
		// Skin tone + ZWJ (profession with skin tone)
		{
			name: "Woman scientist medium skin",
			s:    "👩🏽\u200D🔬", // 👩🏽‍🔬
			want: 2,
		},
		{
			name: "Man firefighter dark skin",
			s:    "👨🏿\u200D🚒", // 👨🏿‍🚒
			want: 2,
		},
		// Multiple modified emoji
		{
			name: "Two skin-toned emoji",
			s:    "👍🏻👋🏿",
			want: 4, // 2 + 2
		},
		// Modified emoji in mixed text
		{
			name: "Mixed text with modified emoji",
			s:    "Hi 👍🏽!",
			want: 6, // H(1)+i(1)+space(1)+thumbs(2)+!(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
				t.Logf("Runes: %U", []rune(tt.s))
			}
		})
	}
}

func TestStringWidth_ZWJEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Standalone ZWJ
		{
			name: "Standalone ZWJ",
			s:    "\u200D",
			want: 0,
		},
		// ZWJ between non-emoji characters
		{
			name: "ZWJ between ASCII",
			s:    "a\u200Db",
			want: 2, // a(1) + ZWJ(0) + b(1)
		},
		// Emoji + ZWJ + non-emoji (invalid ZWJ sequence)
		{
			name: "Emoji + ZWJ + ASCII",
			s:    "😀\u200Da",
			want: 3, // emoji(2) + ZWJ(0) + a(1)
		},
		// Multiple ZWJs without emoji
		{
			name: "Multiple standalone ZWJs",
			s:    "\u200D\u200D\u200D",
			want: 0,
		},
		// Emoji without ZWJ (should be normal)
		{
			name: "Two emoji without ZWJ",
			s:    "😀🚀",
			want: 4, // 2 + 2
		},
		// Single emoji modifier without base
		{
			name: "Orphan skin tone modifier",
			s:    "🏽", // U+1F3FD alone
			want: 2,   // Not preceded by EP, so normal width
		},
		// ZWJ at string boundaries
		{
			name: "Leading ZWJ + emoji",
			s:    "\u200D😀",
			want: 2, // ZWJ(0) + emoji(2)
		},
		{
			name: "Emoji + trailing ZWJ",
			s:    "😀\u200D",
			want: 2, // emoji(2) + ZWJ(0)
		},
		// Very long ZWJ chain
		{
			name: "Long ZWJ chain (3 joins)",
			s:    "👨\u200D👩\u200D👧\u200D👦",
			want: 2,
		},
		// ZWJ sequence followed by regular emoji
		{
			name: "ZWJ family + regular emoji",
			s:    "👨\u200D👩\u200D👧🚀",
			want: 4, // family(2) + rocket(2)
		},
		// Keycap sequences (should still work)
		{
			name: "Keycap 1",
			s:    "1\uFE0F\u20E3",
			want: 2, // 1+VS16 → width 2, combining keycap → width 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidth(tt.s)
			if got != tt.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
				t.Logf("Runes: %U", []rune(tt.s))
			}
		})
	}
}

func TestIsExtendedPictographic(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		// SMP emoji
		{"Grinning face", 0x1F600, true},
		{"Rocket", 0x1F680, true},
		{"Thumbs up", 0x1F44D, true},
		{"Woman", 0x1F469, true},
		{"Man", 0x1F468, true},
		{"Microscope", 0x1F52C, true},

		// BMP emoji
		{"Sun", 0x2600, true},
		{"Heart", 0x2764, true},
		{"Scissors", 0x2702, true},
		{"Watch", 0x231A, true},

		// Individual EP characters
		{"Copyright", 0x00A9, true},
		{"Registered", 0x00AE, true},
		{"Trademark", 0x2122, true},

		// Non-EP characters
		{"ASCII a", 'a', false},
		{"CJK ideograph", 0x4E00, false},
		{"Hangul", 0xAC00, false},
		{"Latin extended", 0x00E9, false},
		{"Combining mark", 0x0300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExtendedPictographic(tt.r)
			if got != tt.want {
				t.Errorf("isExtendedPictographic(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestIsEmojiModifier(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"Light skin tone", 0x1F3FB, true},
		{"Medium-light", 0x1F3FC, true},
		{"Medium", 0x1F3FD, true},
		{"Medium-dark", 0x1F3FE, true},
		{"Dark skin tone", 0x1F3FF, true},
		{"Before range", 0x1F3FA, false},
		{"After range", 0x1F400, false},
		{"Regular emoji", 0x1F600, false},
		{"ASCII", 'a', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmojiModifier(tt.r)
			if got != tt.want {
				t.Errorf("isEmojiModifier(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// TestRuneWidth_UncommonRanges tests coverage for less common Unicode ranges
func TestRuneWidth_UncommonRanges(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		// CJK Compatibility Ideographs (U+F900-U+FAFF) - Tier 2
		{"CJK Compat 豈", '\uF900', 2},
		{"CJK Compat 舘", '\uFAFF', 2},
		{"CJK Compat 福", '\uFA10', 2},

		// Additional emoji ranges - Tier 3
		{"Emoji Transport 🚀", '\U0001F680', 2},
		{"Emoji Transport 🛸", '\U0001F6FF', 2},
		{"Emoji Misc 🔧", '\U0001F527', 2},
		{"Emoji Supplemental 🤗", '\U0001F917', 2},
		{"Emoji Extended 🥳", '\U0001F973', 2},
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

// =============================================================================
// isExtendedPictographic — exhaustive branch coverage for all Unicode ranges
// =============================================================================

func TestIsExtendedPictographic_AllRanges(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		// Range: Misc Symbols and Arrows (U+2B00-U+2BFF)
		{"Misc Symbols/Arrows start U+2B00", 0x2B00, true},
		{"Up arrow U+2B06", 0x2B06, true},
		{"Star U+2B50", 0x2B50, true},
		{"Misc Symbols/Arrows end U+2BFF", 0x2BFF, true},
		{"Below Misc Symbols/Arrows U+2AFF", 0x2AFF, false},

		// Range: Arrow symbols (U+2194-U+21AA)
		{"Left-right arrow U+2194", 0x2194, true},
		{"Rightwards arrow with hook U+21AA", 0x21AA, true},
		{"Mid arrow range U+219E", 0x219E, true},
		{"Below arrow range U+2193", 0x2193, false},

		// Range: Geometric Shapes (U+25A0-U+25FF)
		{"Black square U+25A0", 0x25A0, true},
		{"White circle U+25CB", 0x25CB, true},
		{"Geometric end U+25FF", 0x25FF, true},
		{"Below geometric range U+259F", 0x259F, false},

		// Range: Legacy Computing (U+1FB00-U+1FFFD)
		{"Legacy Computing start U+1FB00", 0x1FB00, true},
		{"Legacy Computing mid U+1FC00", 0x1FC00, true},
		{"Legacy Computing end U+1FFFD", 0x1FFFD, true},
		{"Above legacy range U+1FFFE", 0x1FFFE, false},

		// Verify existing ranges still work
		{"SMP emoji start U+1F000", 0x1F000, true},
		{"SMP emoji end U+1FAFF", 0x1FAFF, true},
		{"Misc Symbols start U+2600", 0x2600, true},
		{"Dingbats end U+27BF", 0x27BF, true},
		{"Misc Technical start U+2300", 0x2300, true},
		{"Misc Technical end U+23FF", 0x23FF, true},

		// Individual EP characters (switch statement)
		{"Double exclamation U+203C", 0x203C, true},
		{"Exclamation question U+2049", 0x2049, true},
		{"Info source U+2139", 0x2139, true},
		{"Wavy dash U+3030", 0x3030, true},
		{"Part alternation U+303D", 0x303D, true},
		{"Circled congratulation U+3297", 0x3297, true},
		{"Circled secret U+3299", 0x3299, true},

		// Negatives: characters in gaps between EP ranges
		{"After arrow, before misc tech U+21AB", 0x21AB, false},
		{"Control Pictures block U+2400", 0x2400, false},
		{"Box Drawing block U+2500", 0x2500, false},
		{"After dingbats U+27C0", 0x27C0, false},
		{"BMP high U+FFFF", 0xFFFF, false},
		{"Regular Latin", 'Z', false},
		{"Null character", 0x0000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExtendedPictographic(tt.r)
			if got != tt.want {
				t.Errorf("isExtendedPictographic(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// =============================================================================
// asciiWidth — direct unit tests for SWAR control character detection
// =============================================================================

func TestAsciiWidth(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// Empty string (n == 0 early return)
		{"empty string", "", 0},

		// Short strings (< 8 bytes, scalar tail only)
		{"single printable", "A", 1},
		{"7 printable chars", "abcdefg", 7},

		// Exactly 8 bytes (one SWAR chunk, no tail)
		{"8 printable chars", "abcdefgh", 8},
		{"8 spaces (0x20 boundary)", "        ", 8},
		{"8 tilde (0x7E boundary)", "~~~~~~~~", 8},

		// Longer strings (multiple SWAR chunks + tail)
		{"16 chars (2 chunks)", "0123456789abcdef", 16},
		{"17 chars (2 chunks + 1 tail)", "0123456789abcdefg", 17},
		{"24 chars (3 chunks)", "abcdefghijklmnopqrstuvwx", 24},

		// SWAR slow path: control character forces byte-by-byte fallback
		{"null in 8-byte chunk", "abcd\x00efg", 7},
		{"tab in 8-byte chunk", "abcdefg\t", 7},
		{"mixed CR/LF in chunk", "abc\ndef\rg", 7},
		{"DEL (0x7F) in chunk", "abcdefg\x7F", 7},
		{"BEL (0x07) in chunk", "abc\x07defg", 7},
		{"multiple controls in chunk", "\t\n\r\x00ABCD", 4},
		{"all control chars in chunk", "\x00\x01\x02\x03\x04\x05\x06\x07", 0},

		// Mixed fast path + slow path across chunks
		{"ctrl first chunk, clean second", "\x01bcdefghijklmnop", 15},
		{"clean first chunk, ctrl second", "abcdefgh\x00jklmnop", 15},
		{"ctrl in both chunks", "\x01bcdefgh\x02jklmnop", 14},

		// Scalar tail with control characters
		{"9 chars with tab at tail", "abcdefgh\t", 8},
		{"10 chars with null at tail", "abcdefghi\x00", 9},
		{"15 chars with DEL at tail", "abcdefghijklmn\x7F", 14},

		// Printable boundary: 0x1F is control, 0x20 is printable
		{"unit separator 0x1F in chunk", "abcdefg\x1F", 7},
		{"space 0x20 in chunk", "abcdefg ", 8},

		// Realistic TUI content (8+ bytes ASCII with mixed content)
		{"terminal prompt", "user@host:~$ ", 13},
		{"table separator", "+--------+--------+", 19},
		{"progress bar", "[=====>          ]", 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := asciiWidth(tt.s)
			if got != tt.want {
				t.Errorf("asciiWidth(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

// TestAsciiWidth_SWARControlAtEveryPosition exercises the SWAR slow path
// with a control character placed at each byte offset within an 8-byte chunk,
// ensuring correct byte-by-byte fallback regardless of position.
func TestAsciiWidth_SWARControlAtEveryPosition(t *testing.T) {
	// SOH (0x01) at each position in an 8-byte chunk
	for pos := 0; pos < 8; pos++ {
		s := []byte("ABCDEFGH")
		s[pos] = 0x01
		t.Run("soh_at_"+string(rune('0'+pos)), func(t *testing.T) {
			got := asciiWidth(string(s))
			if got != 7 {
				t.Errorf("asciiWidth(%q) = %d, want 7 (SOH at pos %d)", s, got, pos)
			}
		})
	}

	// DEL (0x7F) at each position in an 8-byte chunk
	for pos := 0; pos < 8; pos++ {
		s := []byte("ABCDEFGH")
		s[pos] = 0x7F
		t.Run("del_at_"+string(rune('0'+pos)), func(t *testing.T) {
			got := asciiWidth(string(s))
			if got != 7 {
				t.Errorf("asciiWidth(%q) = %d, want 7 (DEL at pos %d)", s, got, pos)
			}
		})
	}
}

// TestStringWidth_ASCIIControlMix verifies that StringWidth correctly delegates
// to asciiWidth for long ASCII strings with embedded control characters.
func TestStringWidth_ASCIIControlMix(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		// 8+ byte ASCII strings that hit the isASCIIOnly→asciiWidth path
		{"16 printable", "Hello, World!!! ", 16},
		{"tab in long string", "Hello\tWorld!!!", 13},
		{"newline in long string", "Hello\nWorld!!!", 13},
		{"DEL in long string", "Hello\x7FWorld!!!", 13},
		{"multiple newlines", "line1\nline2\nline3\n", 15},
		{"null bytes scattered", "abc\x00defg\x00ijklmnop", 15},
		{"TUI box drawing ASCII", "+----------+----------+", 23},

		// Boundary: exactly 8 bytes
		{"8 bytes all printable", "12345678", 8},
		{"8 bytes with tab", "1234567\t", 7},
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
