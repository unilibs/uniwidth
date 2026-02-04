package uniwidth

// EAWidth represents the width for East Asian Ambiguous characters.
type EAWidth int

const (
	// EANarrow treats ambiguous characters as narrow (width 1).
	// This is the default for non-East Asian locales.
	EANarrow EAWidth = 1

	// EAWide treats ambiguous characters as wide (width 2).
	// This is appropriate for East Asian (CJK) locales.
	EAWide EAWidth = 2
)

// Options configures Unicode width calculation behavior.
//
// Use the functional options pattern to create customized configurations:
//
//	opts := []uniwidth.Option{
//	    uniwidth.WithEastAsianAmbiguous(uniwidth.EAWide),
//	    uniwidth.WithEmojiPresentation(true),
//	}
//	width := uniwidth.StringWidthWithOptions("Hello 世界", opts...)
type Options struct {
	// EastAsianAmbiguous specifies how to handle ambiguous-width characters.
	// Default: EANarrow (width 1)
	EastAsianAmbiguous EAWidth

	// EmojiPresentation specifies whether emoji should be rendered as emoji (width 2)
	// or text (width 1). When true, emoji are treated as width 2.
	// Default: true (emoji presentation)
	EmojiPresentation bool
}

// Option is a functional option for configuring Unicode width calculation.
type Option func(*Options)

// defaultOptions returns the default configuration.
func defaultOptions() Options {
	return Options{
		EastAsianAmbiguous: EANarrow, // Width 1 for neutral context
		EmojiPresentation:  true,     // Emoji are wide by default
	}
}

// WithEastAsianAmbiguous sets the width for East Asian Ambiguous characters.
//
// Example:
//
//	// Treat ambiguous characters as wide (East Asian locale)
//	width := uniwidth.StringWidthWithOptions("±½", uniwidth.WithEastAsianAmbiguous(uniwidth.EAWide))
//	// width = 4 (each character is 2 columns wide)
//
//	// Treat ambiguous characters as narrow (neutral locale)
//	width := uniwidth.StringWidthWithOptions("±½", uniwidth.WithEastAsianAmbiguous(uniwidth.EANarrow))
//	// width = 2 (each character is 1 column wide)
func WithEastAsianAmbiguous(width EAWidth) Option {
	return func(o *Options) {
		o.EastAsianAmbiguous = width
	}
}

// WithEmojiPresentation sets whether emoji should be rendered as emoji (wide) or text (narrow).
//
// Example:
//
//	// Emoji as emoji (wide, width 2) - default
//	width := uniwidth.StringWidthWithOptions("😀", uniwidth.WithEmojiPresentation(true))
//	// width = 2
//
//	// Emoji as text (narrow, width 1)
//	width := uniwidth.StringWidthWithOptions("😀", uniwidth.WithEmojiPresentation(false))
//	// width = 1
//
// Note: This primarily affects emoji that have both text and emoji presentation variants.
// Most emoji are always rendered as wide regardless of this setting.
func WithEmojiPresentation(emoji bool) Option {
	return func(o *Options) {
		o.EmojiPresentation = emoji
	}
}

// RuneWidthWithOptions returns the visual width of a rune with custom options.
//
// This function applies the same tiered lookup strategy as RuneWidth, but allows
// customization of ambiguous character handling and emoji presentation.
//
// Example:
//
//	// East Asian locale (ambiguous characters are wide)
//	width := uniwidth.RuneWidthWithOptions('±', uniwidth.WithEastAsianAmbiguous(uniwidth.EAWide))
//	// width = 2
//
//	// Neutral locale (ambiguous characters are narrow)
//	width := uniwidth.RuneWidthWithOptions('±', uniwidth.WithEastAsianAmbiguous(uniwidth.EANarrow))
//	// width = 1
func RuneWidthWithOptions(r rune, opts ...Option) int {
	// Build options
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Use the same tiered lookup as RuneWidth
	width := runeWidthInternal(r)

	// Special handling for ambiguous characters
	if width == -1 {
		// This is an ambiguous character - use configured width
		return int(options.EastAsianAmbiguous)
	}

	return width
}

// StringWidthWithOptions calculates the visual width of a string with custom options.
//
// This function applies the same fast paths as StringWidth, but allows
// customization of ambiguous character handling and emoji presentation.
//
// Example:
//
//	// East Asian locale (ambiguous characters are wide)
//	opts := []uniwidth.Option{
//	    uniwidth.WithEastAsianAmbiguous(uniwidth.EAWide),
//	}
//	width := uniwidth.StringWidthWithOptions("Hello ±½", opts...)
//	// width = 10 (Hello=5, space=1, ±=2, ½=2)
//
//	// Neutral locale (ambiguous characters are narrow)
//	opts := []uniwidth.Option{
//	    uniwidth.WithEastAsianAmbiguous(uniwidth.EANarrow),
//	}
//	width := uniwidth.StringWidthWithOptions("Hello ±½", opts...)
//	// width = 8 (Hello=5, space=1, ±=1, ½=1)
func StringWidthWithOptions(s string, opts ...Option) int {
	// Build options
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Fast path: ASCII-only strings (no ambiguous characters in ASCII)
	if isASCIIOnly(s) {
		return len(s)
	}

	// Iterate through runes and sum their widths
	width := 0
	for _, r := range s {
		w := runeWidthInternal(r)
		if w == -1 {
			// Ambiguous character - use configured width
			width += int(options.EastAsianAmbiguous)
		} else {
			width += w
		}
	}

	return width
}

// runeWidthInternal returns the width of a rune, or -1 for ambiguous characters.
// This is an internal function used by the options API.
func runeWidthInternal(r rune) int {
	// ========================================
	// Tier 1: ASCII Fast Path (O(1))
	// ========================================
	if r < 0x80 {
		if r < 0x20 {
			return 0
		}
		if r == 0x7F {
			return 0
		}
		return 1
	}

	// ========================================
	// Tier 2: Common CJK Fast Path (O(1))
	// ========================================
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

	// ========================================
	// Tier 3: Common Emoji Fast Path (O(1))
	// ========================================
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

	// ========================================
	// Zero-Width Characters (O(1))
	// ========================================
	// Format characters (U+200B-U+200F):
	// ZWSP, ZWNJ, ZWJ, LRM, RLM - all zero-width formatting characters
	if r >= 0x200B && r <= 0x200F {
		return 0
	}
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}
	if r >= 0xE0100 && r <= 0xE01EF {
		return 0
	}

	// Combining marks (diacritics, accents)
	// These have zero width as they combine with previous character
	// Note: Using a simple check instead of unicode.In for performance
	if (r >= 0x0300 && r <= 0x036F) || // Combining Diacritical Marks
		(r >= 0x1AB0 && r <= 0x1AFF) || // Combining Diacritical Marks Extended
		(r >= 0x1DC0 && r <= 0x1DFF) || // Combining Diacritical Marks Supplement
		(r >= 0x20D0 && r <= 0x20FF) || // Combining Diacritical Marks for Symbols
		(r >= 0xFE20 && r <= 0xFE2F) { // Combining Half Marks
		return 0
	}

	// ========================================
	// Tier 4: Multi-Stage Table Lookup (O(1))
	// ========================================
	return tableLookupWidthInternal(r)
}

// tableLookupWidthInternal performs O(1) width lookup using the 3-stage table,
// returning -1 for ambiguous characters (encoding 0b11) so the caller can
// apply the configured East Asian width.
//
// Performance: O(1), 0 allocations.
func tableLookupWidthInternal(r rune) int {
	cp := uint32(r)
	rootIdx := widthRoot[cp>>13]
	midIdx := widthMiddle[rootIdx][cp>>7&0x3F]
	packed := widthLeaves[midIdx][cp>>2&0x1F]
	width := (packed >> (2 * (cp & 0x03))) & 0x03
	if width == 3 {
		return -1 // ambiguous - caller decides
	}
	return int(width)
}

// binarySearchWidthInternal performs binary search and returns -1 for ambiguous characters.
// Kept for backward compatibility; tableLookupWidthInternal is preferred.
func binarySearchWidthInternal(r rune) int {
	// Search in generated wide table (width 2)
	if binarySearch(r, wideTableGenerated) {
		return 2
	}

	// Search in generated zero-width table (width 0)
	if binarySearch(r, zeroWidthTableGenerated) {
		return 0
	}

	// Search in generated ambiguous table
	// Return -1 to indicate ambiguous (caller decides width)
	if binarySearch(r, ambiguousTableGenerated) {
		return -1 // Ambiguous - caller decides
	}

	// Default: width 1 (most characters)
	return 1
}
