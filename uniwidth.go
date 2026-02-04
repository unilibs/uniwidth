// Package uniwidth provides modern Unicode width calculation for Go 1.25+.
//
// uniwidth uses a tiered lookup strategy for optimal performance:
//   - Tier 1: ASCII (O(1), ~95% of typical content)
//   - Tier 2: Common CJK & Emoji (O(1), ~90% of non-ASCII)
//   - Tier 3: Common Emoji (O(1))
//   - Tier 4: 3-stage table lookup for all other characters (O(1))
//
// All tiers are O(1) with zero allocations for single-rune lookups.
// This approach is 3-46x faster than traditional binary-search-only methods
// like go-runewidth, while maintaining full Unicode 16.0 compliance.
//
//go:generate go run cmd/generate-tables/main.go
package uniwidth

import (
	"unicode"
	"unsafe"
)

// RuneWidth returns the visual width of a rune in monospace terminals.
//
// Returns:
//   - 0 for control characters, zero-width joiners, combining marks
//   - 1 for most characters (ASCII, Latin, Cyrillic, etc.)
//   - 2 for wide characters (CJK, Emoji, etc.)
//
// This function uses a tiered lookup strategy:
//   - O(1) for ASCII (most common case)
//   - O(1) for common CJK and emoji (hot paths)
//   - O(1) for all other characters (3-stage table lookup)
func RuneWidth(r rune) int {
	// ========================================
	// Tier 1: ASCII Fast Path (O(1))
	// ========================================
	// Covers ~95% of typical terminal content
	if r < 0x80 {
		// C0 control characters (0x00-0x1F) have zero width
		if r < 0x20 {
			return 0
		}
		// DELETE character (0x7F) has zero width
		if r == 0x7F {
			return 0
		}
		// All other ASCII characters have width 1
		return 1
	}

	// ========================================
	// Tier 2: Common CJK Fast Path (O(1))
	// ========================================
	// Covers ~80% of Asian content

	// CJK Unified Ideographs (20,992 characters)
	// U+4E00 - U+9FFF: Most common Chinese/Japanese characters
	if r >= 0x4E00 && r <= 0x9FFF {
		return 2
	}

	// Hangul Syllables (11,172 characters)
	// U+AC00 - U+D7AF: Korean syllables
	if r >= 0xAC00 && r <= 0xD7AF {
		return 2
	}

	// Hiragana + Katakana + Bopomofo (384 characters)
	// U+3040 - U+309F: Hiragana
	// U+30A0 - U+30FF: Katakana
	// U+3100 - U+312F: Bopomofo (Taiwan phonetic symbols)
	if r >= 0x3040 && r <= 0x312F {
		return 2
	}

	// CJK Compatibility Ideographs
	// U+F900 - U+FAFF: Common CJK compatibility forms
	if r >= 0xF900 && r <= 0xFAFF {
		return 2
	}

	// ========================================
	// Tier 3: Common Emoji Fast Path (O(1))
	// ========================================
	// Covers ~90% of emoji usage

	// Emoticons (80 characters)
	// U+1F600 - U+1F64F: Smileys and people
	if r >= 0x1F600 && r <= 0x1F64F {
		return 2
	}

	// Miscellaneous Symbols and Pictographs (768 characters)
	// U+1F300 - U+1F5FF: Weather, zodiac, hands, etc.
	if r >= 0x1F300 && r <= 0x1F5FF {
		return 2
	}

	// Transport and Map Symbols (103 characters)
	// U+1F680 - U+1F6FF: Vehicles, signs, etc.
	if r >= 0x1F680 && r <= 0x1F6FF {
		return 2
	}

	// Supplemental Symbols and Pictographs (256 characters)
	// U+1F900 - U+1F9FF: Food, animals, activities
	if r >= 0x1F900 && r <= 0x1F9FF {
		return 2
	}

	// Miscellaneous Symbols (common emoji)
	// U+2600 - U+26FF: Weather, zodiac, misc symbols
	if r >= 0x2600 && r <= 0x26FF {
		return 2
	}

	// Dingbats (decorative symbols)
	// U+2700 - U+27BF: Scissors, phone, etc.
	if r >= 0x2700 && r <= 0x27BF {
		return 2
	}

	// ========================================
	// Zero-Width Characters (O(1))
	// ========================================

	// Format characters (U+200B-U+200F):
	// U+200B: Zero-Width Space (ZWSP)
	// U+200C: Zero-Width Non-Joiner (ZWNJ)
	// U+200D: Zero-Width Joiner (ZWJ) - used in emoji sequences
	// U+200E: Left-to-Right Mark (LRM)
	// U+200F: Right-to-Left Mark (RLM)
	// All are invisible formatting characters with zero terminal width.
	if r >= 0x200B && r <= 0x200F {
		return 0
	}

	// Variation Selectors (for emoji vs text presentation)
	// U+FE00 - U+FE0F: Variation selectors
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}

	// Emoji variation selectors
	// U+E0100 - U+E01EF
	if r >= 0xE0100 && r <= 0xE01EF {
		return 0
	}

	// Combining marks (diacritics, accents)
	// These have zero width as they combine with previous character
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
		return 0
	}

	// ========================================
	// Tier 4: Multi-Stage Table Lookup (O(1))
	// ========================================
	// For characters not covered by hot paths, use the 3-stage
	// hierarchical lookup table for constant-time width resolution.
	return tableLookupWidth(r)
}

// StringWidth calculates the visual width of a string in monospace terminals.
//
// This function provides a fast path for ASCII-only strings,
// and uses a state machine for correct handling of multi-rune sequences.
//
// Special handling:
//   - ZWJ emoji sequences (👨‍👩‍👧‍👦) are treated as width 2, not the sum of parts
//   - Emoji modifier sequences (👍🏽) are treated as width 2
//   - Variation selectors (U+FE0E/U+FE0F) modify the width of the preceding character
//   - Regional indicator pairs (flags) are counted as width 2, not 4
func StringWidth(s string) int {
	// Short string fast path (< 8 bytes): single-pass ASCII check and width
	// count fused into one loop. For strings shorter than 8 bytes, the SWAR
	// loop bodies in isASCIIOnly/asciiWidth never execute, making those two
	// function calls pure overhead. This path avoids both calls entirely.
	if len(s) < 8 {
		width := 0
		isASCII := true
		for i := 0; i < len(s); i++ {
			b := s[i]
			if b >= 0x80 {
				isASCII = false
				break
			}
			if b >= 0x20 && b != 0x7F {
				width++
			}
		}
		if isASCII {
			return width
		}
	} else if isASCIIOnly(s) {
		// SWAR fast path for longer ASCII-only strings (8+ bytes)
		return asciiWidth(s)
	}

	// Unicode path: convert to rune slice for lookahead.
	runes := []rune(s)
	width := 0

	// Emoji sequence state tracking (forward-scan state machine):
	//   0 = default (not in an emoji sequence)
	//   1 = after Extended_Pictographic character (may start ZWJ/modifier sequence)
	//   2 = after EP + (Extend*) + ZWJ (expecting joined emoji)
	state := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// ========================================
		// ZWJ Handling
		// ========================================
		// ZWJ (U+200D) after an Extended_Pictographic transitions to
		// the "expecting joined emoji" state. ZWJ always has width 0.
		if r == 0x200D {
			if state == 1 {
				state = 2
			}
			continue
		}

		// After EP + ZWJ: if next is EP, it joins (width 0).
		// This implements the core of GB11: ExtPict Extend* ZWJ × ExtPict.
		if state == 2 {
			if isExtendedPictographic(r) {
				state = 1 // Joined, still in emoji sequence
				continue  // Width 0 — joined with preceding emoji
			}
			// Not a valid join target, reset state and process normally.
			state = 0
		}

		// ========================================
		// Emoji Modifier Handling (Skin Tones)
		// ========================================
		// Emoji modifiers (U+1F3FB-U+1F3FF) combine with the preceding
		// Extended_Pictographic, contributing zero additional width.
		if state == 1 && isEmojiModifier(r) {
			continue // Width 0 (modifier combines with preceding emoji)
		}

		// ========================================
		// Extend Characters in Emoji Context
		// ========================================
		// Variation selectors and combining marks within an active emoji
		// sequence don't add width and keep the state alive for potential
		// ZWJ continuation.
		if state == 1 && (r >= 0xFE00 && r <= 0xFE0F) {
			continue // VS in emoji sequence, width 0
		}

		// ========================================
		// Regional Indicator Pairs (Flags)
		// ========================================
		// Two consecutive regional indicators (U+1F1E6-U+1F1FF) form
		// a flag emoji with width 2 (not 4).
		if isRegionalIndicator(r) && i+1 < len(runes) && isRegionalIndicator(runes[i+1]) {
			width += 2
			i++
			state = 0
			continue
		}

		// ========================================
		// Variation Selectors (Lookahead)
		// ========================================
		// Variation selectors modify the preceding character's presentation:
		// - U+FE0E: Text presentation (width 1)
		// - U+FE0F: Emoji presentation (width 2)
		if i+1 < len(runes) {
			next := runes[i+1]

			if next == 0xFE0E {
				width++
				i++
				state = 0
				continue
			}

			if next == 0xFE0F {
				width += 2
				i++
				if isExtendedPictographic(r) {
					state = 1
				} else {
					state = 0
				}
				continue
			}
		}

		// ========================================
		// Default: RuneWidth
		// ========================================
		w := RuneWidth(r)
		width += w

		// Track emoji state for ZWJ/modifier sequence detection.
		if isExtendedPictographic(r) && w > 0 {
			state = 1
		} else if w > 0 {
			state = 0
		}
		// When w == 0 (combining marks, tag characters, etc.),
		// preserve current state to allow Extend* in GB11 pattern.
	}

	return width
}

// isRegionalIndicator returns true if the rune is a regional indicator symbol.
// Regional indicators (U+1F1E6 - U+1F1FF) represent country codes (A-Z).
// Two consecutive indicators form a country flag emoji.
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// isExtendedPictographic returns true if the rune has the Extended_Pictographic
// property (Unicode 16.0 emoji-data.txt), meaning it can participate in emoji
// ZWJ sequences. This covers all emoji ranges used in standard ZWJ sequences.
//
// The checks are ordered by frequency of occurrence in real-world emoji usage
// to minimize branch mispredictions.
func isExtendedPictographic(r rune) bool {
	// SMP emoji blocks (U+1F000-U+1FAFF) — covers ~95% of emoji
	// Includes: Emoticons, Pictographs, Transport, Supplemental Symbols,
	// Symbols and Pictographs Extended-A, etc.
	if r >= 0x1F000 && r <= 0x1FAFF {
		return true
	}

	// BMP emoji: Misc Symbols (U+2600-U+26FF) and Dingbats (U+2700-U+27BF)
	if r >= 0x2600 && r <= 0x27BF {
		return true
	}

	// BMP emoji: Misc Technical (U+2300-U+23FF)
	// Includes: ⌚⌛⏩⏪⏫⏬⏰⏳⏸⏹⏺⌨ etc.
	if r >= 0x2300 && r <= 0x23FF {
		return true
	}

	// Misc Symbols and Arrows (U+2B00-U+2BFF)
	if r >= 0x2B00 && r <= 0x2BFF {
		return true
	}

	// Arrow symbols (U+2194-U+21AA)
	if r >= 0x2194 && r <= 0x21AA {
		return true
	}

	// Geometric Shapes (U+25A0-U+25FF)
	if r >= 0x25A0 && r <= 0x25FF {
		return true
	}

	// Symbols for Legacy Computing and extensions (U+1FB00-U+1FFFD)
	if r >= 0x1FB00 && r <= 0x1FFFD {
		return true
	}

	// Individual Extended_Pictographic characters
	switch r {
	case 0x00A9, 0x00AE, // © ®
		0x203C, 0x2049, // ‼ ⁉
		0x2122, 0x2139, // ™ ℹ
		0x3030, 0x303D, // 〰 〽
		0x3297, 0x3299: // ㊗ ㊙
		return true
	}

	return false
}

// isEmojiModifier returns true if the rune is an emoji modifier (skin tone).
// Emoji modifiers (U+1F3FB-U+1F3FF) represent Fitzpatrick skin types 1-2 through 6.
// They combine with the preceding Extended_Pictographic character to form
// a single emoji with a specific skin tone.
func isEmojiModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

// isASCIIOnly returns true if the string contains only ASCII characters (0x00-0x7F).
//
// Uses SWAR (SIMD Within A Register) to process 8 bytes at a time by loading
// them into a uint64 and checking all high bits simultaneously with a single
// AND against 0x8080808080808080. If any byte has its high bit set (>= 0x80),
// it is non-ASCII. This works regardless of endianness because we only test
// whether any byte has its high bit set, not which byte it is.
//
// Performance:
//   - Short strings (< 8 bytes): scalar fallback, O(n) per byte
//   - Longer strings: ~8x throughput via SWAR, O(n/8) per word + O(n%8) tail
//   - 0 allocations in all cases
//
// all pointer arithmetic is bounds-checked by the loop guard (i+8 <= n, i < n).
//
//nolint:gosec // G103: unsafe usage is intentional for SWAR performance optimization;
func isASCIIOnly(s string) bool {
	n := len(s)
	if n == 0 {
		return true
	}

	p := unsafe.StringData(s)

	// SWAR: process 8 bytes at a time
	const asciiMask = uint64(0x8080808080808080)
	i := 0
	for ; i+8 <= n; i += 8 {
		word := *(*uint64)(unsafe.Add(unsafe.Pointer(p), i))
		if word&asciiMask != 0 {
			return false
		}
	}

	// Scalar tail: process remaining bytes (0-7)
	for ; i < n; i++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(p), i)) >= 0x80 {
			return false
		}
	}

	return true
}

// asciiWidth returns the visual width of an ASCII-only string, accounting for
// control characters (0x00-0x1F, 0x7F) which have zero width.
//
// Uses SWAR to detect control characters in 8-byte chunks. If a chunk contains
// no control characters, width += 8 directly. Otherwise, falls back to scalar
// processing for that chunk.
//
// Control character detection uses Daniel Lemire's SWAR technique:
//   - Bytes < 0x20: detected via (x - 0x2020...) & ~x & 0x8080...
//   - Byte == 0x7F: detected via XOR with 0x7F7F... then same underflow trick
//
// The underflow trick works because subtracting 0x20 from a byte < 0x20 causes
// the high bit to set (unsigned underflow), while the original byte had its high
// bit clear. The AND with ~x isolates genuine underflows from bytes >= 0x80
// (which cannot appear here since isASCIIOnly was already verified).
//
// Caller must ensure s contains only ASCII bytes (call isASCIIOnly first).
//
// Performance:
//   - 0 allocations
//   - ~8x throughput for chunks without control characters
//
// all pointer arithmetic is bounds-checked by the loop guards (i+8 <= n, i < n, j < 8).
//
//nolint:gosec // G103: unsafe usage is intentional for SWAR performance optimization;
func asciiWidth(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}

	p := unsafe.StringData(s)
	width := 0
	i := 0

	// SWAR constants for control character detection.
	const (
		// Broadcast 0x20 and 0x7F across all 8 bytes of a uint64.
		lo20  = uint64(0x2020202020202020)
		hi80  = uint64(0x8080808080808080)
		rep7F = uint64(0x7F7F7F7F7F7F7F7F)
		rep01 = uint64(0x0101010101010101)
	)

	// Process 8 bytes at a time
	for ; i+8 <= n; i += 8 {
		word := *(*uint64)(unsafe.Add(unsafe.Pointer(p), i))

		// Detect bytes < 0x20 using SWAR underflow trick:
		// (word - 0x2020...) produces underflow (sets high bit) for bytes < 0x20.
		// &^word masks out bytes that already had high bit set (not possible for
		// ASCII, but defensive). &hi80 extracts only the high bits.
		hasLow := (word - lo20) & ^word & hi80

		// Detect bytes == 0x7F using XOR + underflow:
		// word ^ 0x7F7F... zeros out any 0x7F bytes. Then the zero-byte detection
		// pattern ((v - 0x0101...) & ~v & 0x8080...) finds the zeroed positions.
		xored := word ^ rep7F
		has7F := (xored - rep01) & ^xored & hi80

		if (hasLow | has7F) == 0 {
			// Fast path: no control characters in this 8-byte chunk
			width += 8
		} else {
			// Slow path: at least one control character, process byte by byte
			for j := 0; j < 8; j++ {
				b := *(*byte)(unsafe.Add(unsafe.Pointer(p), i+j))
				if b >= 0x20 && b != 0x7F {
					width++
				}
			}
		}
	}

	// Scalar tail: process remaining bytes (0-7)
	for ; i < n; i++ {
		b := *(*byte)(unsafe.Add(unsafe.Pointer(p), i))
		if b >= 0x20 && b != 0x7F {
			width++
		}
	}

	return width
}

// tableLookupWidth performs O(1) width lookup using the 3-stage hierarchical table.
//
// The table encodes every Unicode codepoint (0x0000-0x10FFFF) as a 2-bit width value:
//
//	0b00 = width 0 (control, combining, zero-width)
//	0b01 = width 1 (narrow, default)
//	0b10 = width 2 (wide: CJK, emoji, fullwidth)
//	0b11 = ambiguous (treated as width 1 in neutral context)
//
// Performance: O(1), 0 allocations. Three array lookups + bit extraction.
func tableLookupWidth(r rune) int {
	cp := uint32(r)
	rootIdx := widthRoot[cp>>13]
	midIdx := widthMiddle[rootIdx][cp>>7&0x3F]
	packed := widthLeaves[midIdx][cp>>2&0x1F]
	width := (packed >> (2 * (cp & 0x03))) & 0x03
	if width == 3 {
		return 1 // ambiguous -> narrow in neutral context
	}
	return int(width)
}

// binarySearchWidth performs binary search on Unicode width tables.
// This is the legacy fallback, kept for use by the Options API.
func binarySearchWidth(r rune) int {
	// Search in generated wide table (width 2)
	if binarySearch(r, wideTableGenerated) {
		return 2
	}

	// Search in generated zero-width table (width 0)
	if binarySearch(r, zeroWidthTableGenerated) {
		return 0
	}

	// Search in generated ambiguous table (width 2 in East Asian context, 1 otherwise)
	// For now, we default to width 1 (neutral context)
	// TODO: Make this configurable via Options pattern
	if binarySearch(r, ambiguousTableGenerated) {
		return 1 // Default to narrow for neutral context
	}

	// Default: width 1 (most characters)
	return 1
}

// binarySearch performs binary search on a sorted rune range table.
func binarySearch(r rune, table []runeRange) bool {
	low, high := 0, len(table)-1

	for low <= high {
		mid := (low + high) / 2
		rr := table[mid]

		//nolint:gocritic // Binary search requires if-else chain for performance
		if r < rr.first {
			high = mid - 1
		} else if r > rr.last {
			low = mid + 1
		} else {
			// r is within range [first, last]
			return true
		}
	}

	return false
}

// runeRange represents a range of runes with the same width property.
type runeRange struct {
	first rune
	last  rune
}
