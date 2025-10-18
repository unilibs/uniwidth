package uniwidth

// This file contains Unicode width tables generated from Unicode 16.0 data.
// These tables are used as fallback for characters not covered by fast path tiers.
//
// Table generation strategy:
// - Hot path characters (ASCII, common CJK, emoji) are handled by range checks in uniwidth.go
// - These tables contain remaining characters that need binary search
// - This minimizes table size while maximizing performance

// wideTable contains ranges of characters with East Asian Width property W (Wide) or F (Fullwidth).
// These characters occupy 2 terminal columns.
var wideTable = []runeRange{
	// CJK Symbols and Punctuation (partial, not covered by fast path)
	{0x3000, 0x303F}, // Ideographic space, CJK symbols, Ideographic Half Fill Space

	// CJK Radicals Supplement
	{0x2E80, 0x2E99},
	{0x2E9B, 0x2EF3},

	// Kangxi Radicals
	{0x2F00, 0x2FD5},

	// CJK Strokes
	{0x31C0, 0x31E3},

	// Enclosed CJK Letters and Months
	{0x3200, 0x321E},
	{0x3220, 0x3247},
	{0x3250, 0x4DBE}, // Fixed: U+4DBF-U+4DFF are unassigned

	// CJK Unified Ideographs Extension A
	// (Already covered by fast path: 0x4E00-0x9FFF)

	// CJK Compatibility Forms
	{0xFE30, 0xFE4F},

	// Halfwidth and Fullwidth Forms (fullwidth part)
	{0xFF01, 0xFF60}, // Fullwidth ASCII variants
	{0xFFE0, 0xFFE6}, // Fullwidth currency signs

	// Kana Supplement
	{0x1B000, 0x1B0FF},

	// CJK Unified Ideographs Extension B-G (not covered by fast path)
	{0x20000, 0x2A6DF}, // Extension B
	{0x2A700, 0x2B73F}, // Extension C
	{0x2B740, 0x2B81F}, // Extension D
	{0x2B820, 0x2CEAF}, // Extension E
	{0x2CEB0, 0x2EBEF}, // Extension F
	{0x30000, 0x3134F}, // Extension G

	// Additional emoji ranges not in fast path
	{0x2600, 0x26FF},   // Miscellaneous Symbols
	{0x2700, 0x27BF},   // Dingbats
	{0x1F000, 0x1F02F}, // Mahjong Tiles
	{0x1F0A0, 0x1F0FF}, // Playing Cards
	{0x1FA00, 0x1FA6F}, // Chess Symbols
	{0x1FA70, 0x1FAFF}, // Symbols and Pictographs Extended-A

	// Ancient scripts (supplementary plane)
	{0x10000, 0x1007F}, // Linear B Syllabary (Ancient Greek)
}

// zeroWidthTable contains ranges of characters with zero width.
// These are control characters, combining marks, and format characters.
var zeroWidthTable = []runeRange{
	// C0 control characters (already handled in fast path)
	// {0x0000, 0x001F},

	// C1 control characters
	{0x0080, 0x009F},

	// Combining Diacritical Marks (partial, rest handled by unicode.In check)
	{0x0300, 0x036F},

	// Combining Diacritical Marks Extended
	{0x1AB0, 0x1AFF},

	// Hebrew combining marks
	{0x0591, 0x05BD},
	{0x05BF, 0x05BF},
	{0x05C1, 0x05C2},
	{0x05C4, 0x05C5},
	{0x05C7, 0x05C7},

	// Arabic combining marks
	{0x0610, 0x061A},
	{0x064B, 0x065F},
	{0x0670, 0x0670},
	{0x06D6, 0x06DC},
	{0x06DF, 0x06E4},
	{0x06E7, 0x06E8},
	{0x06EA, 0x06ED},

	// Devanagari combining marks
	{0x0901, 0x0902},
	{0x093A, 0x093A},
	{0x093C, 0x093C},
	{0x0941, 0x0948},
	{0x094D, 0x094D},
	{0x0951, 0x0957},
	{0x0962, 0x0963},

	// Soft hyphen
	{0x00AD, 0x00AD},

	// Format characters
	{0x200B, 0x200F}, // Zero-width space, LRM, RLM, etc.
	// ZWJ and ZWNJ already handled in fast path

	// Combining marks for symbols
	{0x20D0, 0x20FF},

	// Variation selectors (partial, rest in fast path)
	// {0xFE00, 0xFE0F}, // Already in fast path

	// Arabic presentation forms (zero-width)
	{0xFE20, 0xFE2F},

	// Combining Half Marks
	{0xFE30, 0xFE2F},

	// Specials (BOM, etc.)
	{0xFEFF, 0xFEFF},
}

// ambiguousTable contains ranges of characters with East Asian Width property A (Ambiguous).
// Width depends on context (East Asian: 2, neutral: 1).
// For now, we default to width 1 (neutral context).
var ambiguousTable = []runeRange{
	// Greek and Coptic (partial)
	{0x00A1, 0x00A1}, // Inverted exclamation mark
	{0x00A4, 0x00A4}, // Currency sign
	{0x00A7, 0x00A8}, // Section sign, diaeresis
	{0x00AA, 0x00AA}, // Feminine ordinal indicator
	{0x00AD, 0x00AE}, // Soft hyphen, registered sign
	{0x00B0, 0x00B4}, // Degree sign, acute accent, etc.
	{0x00B6, 0x00BA}, // Pilcrow sign, middle dot, etc.
	{0x00BC, 0x00BF}, // Fractions, inverted question mark
	{0x00C6, 0x00C6}, // Latin capital letter AE
	{0x00D0, 0x00D0}, // Latin capital letter Eth
	{0x00D7, 0x00D8}, // Multiplication sign, O with stroke
	{0x00DE, 0x00E1}, // Thorn, a with acute, etc.
	{0x00E6, 0x00E6}, // Latin small letter ae
	{0x00E8, 0x00EA}, // e with grave, acute, circumflex
	{0x00EC, 0x00ED}, // i with grave, acute
	{0x00F0, 0x00F0}, // Latin small letter eth
	{0x00F2, 0x00F3}, // o with grave, acute
	{0x00F7, 0x00FA}, // Division sign, o with stroke, etc.
	{0x00FC, 0x00FC}, // u with diaeresis
	{0x00FE, 0x00FE}, // Latin small letter thorn
	{0x0101, 0x0101}, // a with macron
	{0x0111, 0x0111}, // d with stroke
	{0x0113, 0x0113}, // e with macron
	{0x011B, 0x011B}, // e with caron
	{0x0126, 0x0127}, // H with stroke
	{0x012B, 0x012B}, // i with macron
	{0x0131, 0x0133}, // Dotless i, IJ ligature
	{0x0138, 0x0138}, // Kra
	{0x013F, 0x0142}, // L with middle dot, l with stroke
	{0x0144, 0x0144}, // n with acute
	{0x0148, 0x014B}, // n with caron, Eng
	{0x014D, 0x014D}, // o with macron
	{0x0152, 0x0153}, // OE ligature
	{0x0166, 0x0167}, // T with stroke
	{0x016B, 0x016B}, // u with macron
	{0x01CE, 0x01CE}, // a with caron
	{0x01D0, 0x01D0}, // i with caron
	{0x01D2, 0x01D2}, // o with caron
	{0x01D4, 0x01D4}, // u with caron
	{0x01D6, 0x01D6}, // u with diaeresis and macron
	{0x01D8, 0x01D8}, // u with diaeresis and acute
	{0x01DA, 0x01DA}, // u with diaeresis and caron
	{0x01DC, 0x01DC}, // u with diaeresis and grave

	// Box Drawing
	{0x2500, 0x257F},

	// Block Elements
	{0x2580, 0x259F},

	// Geometric Shapes
	{0x25A0, 0x25FF},
}
