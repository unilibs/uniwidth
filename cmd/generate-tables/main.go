// generate-tables generates Unicode width tables from official Unicode 16.0 data.
//
// This tool downloads and parses:
// - EastAsianWidth.txt - East Asian Width property assignments
// - emoji-data.txt - Emoji presentation properties
//
// It generates optimized tables for uniwidth's tiered lookup strategy:
//   - Tier 1-3 (hot paths) are hardcoded in uniwidth.go for O(1) lookup
//   - This generates Tier 4 tables: both legacy binary search tables and
//     a 3-stage multi-stage lookup table for O(1) fallback
//
// Usage:
//
//	go run cmd/generate-tables/main.go
//
// Output:
//
//	tables_generated.go - Generated Unicode width tables
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	unicodeVersion    = "16.0.0"
	eastAsianWidthURL = "https://www.unicode.org/Public/16.0.0/ucd/EastAsianWidth.txt"
	emojiDataURL      = "https://www.unicode.org/Public/16.0.0/ucd/emoji/emoji-data.txt"
	outputFile        = "tables_generated.go"

	// maxCodepoint is the maximum valid Unicode codepoint (U+10FFFF).
	maxCodepoint = 0x10FFFF

	// 2-bit width encoding for multi-stage table leaves.
	widthZero      = 0 // width 0: control, combining, zero-width
	widthNarrow    = 1 // width 1: narrow (default)
	widthWide      = 2 // width 2: wide (CJK, emoji, fullwidth)
	widthAmbiguous = 3 // width 1 in neutral context, 2 in East Asian
)

// runeRange represents a contiguous range of runes with the same property.
type runeRange struct {
	first rune
	last  rune
}

// category represents different width categories
type category int

const (
	catWide category = iota
	catZeroWidth
	catAmbiguous
	catNarrow
)

func main() {
	log.Println("Generating Unicode 16.0 width tables...")

	// Download and parse Unicode data
	log.Println("Downloading EastAsianWidth.txt...")
	eawData, err := downloadFile(eastAsianWidthURL)
	if err != nil {
		log.Fatalf("Failed to download EastAsianWidth.txt: %v", err)
	}

	log.Println("Parsing East Asian Width data...")
	wideRanges, ambiguousRanges := parseEastAsianWidth(eawData)

	log.Println("Downloading emoji-data.txt...")
	emojiData, err := downloadFile(emojiDataURL)
	if err != nil {
		log.Fatalf("Failed to download emoji-data.txt: %v", err)
	}

	log.Println("Parsing Emoji data...")
	emojiRanges := parseEmojiData(emojiData)

	// Build multi-stage table from UNFILTERED ranges (covers all codepoints)
	log.Println("Building multi-stage lookup table...")
	root, middle, leaves := buildMultiStageTable(wideRanges, ambiguousRanges, emojiRanges)
	log.Printf("  - Root table: %d entries", len(root))
	log.Printf("  - Middle tables: %d unique sub-tables", len(middle))
	log.Printf("  - Leaf tables: %d unique sub-tables", len(leaves))
	totalBytes := len(root) + len(middle)*64 + len(leaves)*32
	log.Printf("  - Total size: %d bytes (%.1f KiB)", totalBytes, float64(totalBytes)/1024)

	// Merge emoji into wide ranges for legacy tables
	wideRanges = mergeRanges(wideRanges, emojiRanges)

	// Generate zero-width tables (control chars, combining marks, format chars)
	log.Println("Generating zero-width tables...")
	zeroWidthRanges := generateZeroWidthRanges()

	// Filter out hot path ranges (already handled in uniwidth.go) for legacy tables
	log.Println("Filtering hot path ranges (Tier 1-3)...")
	wideRanges = filterHotPaths(wideRanges)
	zeroWidthRanges = filterZeroWidthHotPaths(zeroWidthRanges)
	ambiguousRanges = filterAmbiguousHotPaths(ambiguousRanges)

	// Optimize ranges (merge adjacent ranges)
	log.Println("Optimizing range tables...")
	wideRanges = optimizeRanges(wideRanges)
	zeroWidthRanges = optimizeRanges(zeroWidthRanges)
	ambiguousRanges = optimizeRanges(ambiguousRanges)

	// Generate output file
	log.Println("Generating tables_generated.go...")
	err = generateGoFile(wideRanges, zeroWidthRanges, ambiguousRanges, &root, middle, leaves)
	if err != nil {
		log.Fatalf("Failed to generate Go file: %v", err)
	}

	log.Printf("Successfully generated %s with:", outputFile)
	log.Printf("  - Wide characters: %d ranges", len(wideRanges))
	log.Printf("  - Zero-width characters: %d ranges", len(zeroWidthRanges))
	log.Printf("  - Ambiguous characters: %d ranges", len(ambiguousRanges))
	log.Printf("  - Multi-stage table: root=%d, middle=%d, leaves=%d", len(root), len(middle), len(leaves))
	log.Println("Done!")
}

// downloadFile downloads a file from a URL and returns its content as a string.
//
//nolint:gosec // URL is hardcoded constant from Unicode.org
func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// parseEastAsianWidth parses EastAsianWidth.txt and returns wide and ambiguous ranges.
func parseEastAsianWidth(data string) (wide, ambiguous []runeRange) {
	// Regex to match lines like:
	// 0020          ; N        # Zs       SPACE
	// 3000..303F    ; W        # So  [64] IDEOGRAPHIC SPACE..IDEOGRAPHIC HALF FILL SPACE
	lineRe := regexp.MustCompile(`^([0-9A-F]+)(?:\.\.([0-9A-F]+))?\s*;\s*([A-Z]+)`)

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := lineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		first, err := strconv.ParseInt(matches[1], 16, 64)
		if err != nil {
			continue
		}

		last := first
		if matches[2] != "" {
			l, err := strconv.ParseInt(matches[2], 16, 64)
			if err != nil {
				continue
			}
			last = l
		}

		width := matches[3]

		rr := runeRange{first: rune(first), last: rune(last)}

		switch width {
		case "W", "F": // Wide or Fullwidth
			wide = append(wide, rr)
		case "A": // Ambiguous
			ambiguous = append(ambiguous, rr)
		}
	}

	return wide, ambiguous
}

// parseEmojiData parses emoji-data.txt and returns emoji ranges.
func parseEmojiData(data string) []runeRange {
	// Regex to match lines like:
	// 0023          ; Emoji                # E0.0   [1] (#)       number sign
	// 1F600..1F64F  ; Emoji                # E0.6  [80] (...)    grinning face..folded hands
	lineRe := regexp.MustCompile(`^([0-9A-F]+)(?:\.\.([0-9A-F]+))?\s*;\s*Emoji\s`)

	var ranges []runeRange

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := lineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		first, err := strconv.ParseInt(matches[1], 16, 64)
		if err != nil {
			continue
		}

		last := first
		if matches[2] != "" {
			l, err := strconv.ParseInt(matches[2], 16, 64)
			if err != nil {
				continue
			}
			last = l
		}

		ranges = append(ranges, runeRange{first: rune(first), last: rune(last)})
	}

	return ranges
}

// generateZeroWidthRanges generates ranges for zero-width characters.
func generateZeroWidthRanges() []runeRange {
	// These are well-known zero-width character ranges
	return []runeRange{
		// C0 control characters
		{0x0000, 0x001F},
		// DELETE
		{0x007F, 0x007F},
		// C1 control characters
		{0x0080, 0x009F},
		// Soft hyphen
		{0x00AD, 0x00AD},
		// Combining Diacritical Marks
		{0x0300, 0x036F},
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
		// Combining Diacritical Marks Extended (U+1AB0-U+1AFF)
		{0x1AB0, 0x1AFF},
		// Combining Diacritical Marks Supplement (U+1DC0-U+1DFF)
		{0x1DC0, 0x1DFF},
		// Format characters (ZWS, ZWNJ, ZWJ, LRM, RLM, etc.)
		{0x200B, 0x200F},
		// Combining marks for symbols
		{0x20D0, 0x20FF},
		// Variation selectors
		{0xFE00, 0xFE0F},
		// Arabic presentation forms
		{0xFE20, 0xFE2F},
		// BOM and other specials
		{0xFEFF, 0xFEFF},
		// Emoji variation selectors
		{0xE0100, 0xE01EF},
	}
}

// filterHotPaths removes ranges already handled by Tier 1-3 hot paths in uniwidth.go.
func filterHotPaths(ranges []runeRange) []runeRange {
	// Hot paths (hardcoded in uniwidth.go for O(1) lookup):
	hotPaths := []runeRange{
		// ASCII (Tier 1)
		{0x0000, 0x007F},
		// CJK Unified Ideographs (Tier 2)
		{0x4E00, 0x9FFF},
		// Hangul Syllables (Tier 2)
		{0xAC00, 0xD7AF},
		// Hiragana + Katakana (Tier 2)
		{0x3040, 0x30FF},
		// CJK Compatibility Ideographs (Tier 2)
		{0xF900, 0xFAFF},
		// Emoticons (Tier 3)
		{0x1F600, 0x1F64F},
		// Miscellaneous Symbols and Pictographs (Tier 3)
		{0x1F300, 0x1F5FF},
		// Transport and Map Symbols (Tier 3)
		{0x1F680, 0x1F6FF},
		// Supplemental Symbols and Pictographs (Tier 3)
		{0x1F900, 0x1F9FF},
		// Miscellaneous Symbols (Tier 3)
		{0x2600, 0x26FF},
		// Dingbats (Tier 3)
		{0x2700, 0x27BF},
	}

	return removeOverlappingRanges(ranges, hotPaths)
}

// filterZeroWidthHotPaths removes zero-width ranges already handled in uniwidth.go.
func filterZeroWidthHotPaths(ranges []runeRange) []runeRange {
	hotPaths := []runeRange{
		// ASCII control chars (handled in Tier 1)
		{0x0000, 0x001F},
		{0x007F, 0x007F},
		// Format characters: ZWSP, ZWNJ, ZWJ, LRM, RLM (handled explicitly)
		{0x200B, 0x200F},
		// Variation selectors (handled explicitly)
		{0xFE00, 0xFE0F},
		{0xE0100, 0xE01EF},
	}

	return removeOverlappingRanges(ranges, hotPaths)
}

// filterAmbiguousHotPaths removes ambiguous ranges that overlap with hot paths.
func filterAmbiguousHotPaths(ranges []runeRange) []runeRange {
	// No specific hot paths for ambiguous, but we should remove ASCII range
	hotPaths := []runeRange{
		{0x0000, 0x007F}, // ASCII
	}

	return removeOverlappingRanges(ranges, hotPaths)
}

// removeOverlappingRanges removes ranges that overlap with hotPaths.
func removeOverlappingRanges(ranges, hotPaths []runeRange) []runeRange {
	var result []runeRange

	for _, rr := range ranges {
		overlaps := false
		for _, hp := range hotPaths {
			if rangesOverlap(rr, hp) {
				overlaps = true
				break
			}
		}
		if !overlaps {
			result = append(result, rr)
		}
	}

	return result
}

// rangesOverlap returns true if two ranges overlap.
func rangesOverlap(a, b runeRange) bool {
	return a.first <= b.last && b.first <= a.last
}

// mergeRanges merges two sets of ranges.
func mergeRanges(a, b []runeRange) []runeRange {
	result := append([]runeRange{}, a...)
	result = append(result, b...)
	return result
}

// optimizeRanges merges adjacent ranges and sorts them.
func optimizeRanges(ranges []runeRange) []runeRange {
	if len(ranges) == 0 {
		return ranges
	}

	// Sort by first rune
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].first < ranges[j].first
	})

	// Merge adjacent or overlapping ranges
	result := []runeRange{ranges[0]}

	for i := 1; i < len(ranges); i++ {
		last := &result[len(result)-1]
		current := ranges[i]

		// If current range is adjacent or overlapping, merge it
		if current.first <= last.last+1 {
			if current.last > last.last {
				last.last = current.last
			}
		} else {
			result = append(result, current)
		}
	}

	return result
}

// buildWidthMap builds a complete width map for all Unicode codepoints.
// The map uses the 2-bit encoding: 0=zero, 1=narrow, 2=wide, 3=ambiguous.
//
// The application order matters: later assignments override earlier ones,
// matching the priority logic in RuneWidth() and binarySearchWidth().
func buildWidthMap(wide, ambiguous, emoji []runeRange) []byte {
	// Allocate map for all codepoints (0x000000 - 0x10FFFF)
	widthMap := make([]byte, maxCodepoint+1)

	// Step 1: Default everything to narrow (width 1)
	for i := range widthMap {
		widthMap[i] = widthNarrow
	}

	// Step 2: Apply zero-width ranges
	zeroWidthRanges := generateZeroWidthRanges()
	for _, rr := range zeroWidthRanges {
		for cp := rr.first; cp <= rr.last; cp++ {
			widthMap[cp] = widthZero
		}
	}

	// Step 3: Mark surrogates as zero-width (they are invalid in Go strings)
	for cp := rune(0xD800); cp <= 0xDFFF; cp++ {
		widthMap[cp] = widthZero
	}

	// Step 4: Apply ambiguous ranges (encoded as 3)
	// Must be applied BEFORE wide ranges so that characters that are both
	// ambiguous AND wide (due to emoji overlap) get the correct wide width.
	for _, rr := range ambiguous {
		for cp := rr.first; cp <= rr.last; cp++ {
			widthMap[cp] = widthAmbiguous
		}
	}

	// Step 5: Apply wide ranges from EastAsianWidth (W, F)
	for _, rr := range wide {
		for cp := rr.first; cp <= rr.last; cp++ {
			widthMap[cp] = widthWide
		}
	}

	// Step 6: Apply emoji ranges (width 2)
	// Emoji override ambiguous (e.g., U+2600-U+26FF are both ambiguous and emoji)
	for _, rr := range emoji {
		for cp := rr.first; cp <= rr.last; cp++ {
			widthMap[cp] = widthWide
		}
	}

	// Step 7: Re-apply zero-width overrides that must take precedence
	// (variation selectors, ZWJ, ZWNJ, combining marks, control chars, etc.)
	// These are zero-width regardless of any other property.
	for _, rr := range zeroWidthRanges {
		for cp := rr.first; cp <= rr.last; cp++ {
			widthMap[cp] = widthZero
		}
	}

	// Step 8: Ensure ASCII is correct
	// C0 control characters (0x00-0x1F): width 0
	for cp := rune(0x00); cp <= 0x1F; cp++ {
		widthMap[cp] = widthZero
	}
	// Printable ASCII (0x20-0x7E): width 1
	for cp := rune(0x20); cp <= 0x7E; cp++ {
		widthMap[cp] = widthNarrow
	}
	// DELETE (0x7F): width 0
	widthMap[0x7F] = widthZero

	return widthMap
}

// buildMultiStageTable constructs a 3-stage hierarchical lookup table from Unicode data.
//
// The 3-stage table splits a 21-bit Unicode codepoint into 3 parts:
//
//	Codepoint: [20...13][12...7][6...0]
//	             8 bits   6 bits  7 bits
//	Stage:       ROOT    MIDDLE   LEAF
//
// ROOT (256 entries): indexes into MIDDLE
// MIDDLE (N x 64 entries): indexes into LEAVES
// LEAVES (M x 32 entries): packed 2-bit width values, 4 per byte
//
// Deduplication of identical sub-tables is critical for compact size.
func buildMultiStageTable(wide, ambiguous, emoji []runeRange) (root [256]byte, middle [][64]byte, leaves [][32]byte) {
	widthMap := buildWidthMap(wide, ambiguous, emoji)

	// Maps for deduplication: serialized sub-table -> index
	leafIndex := make(map[[32]byte]byte)
	midIndex := make(map[[64]byte]byte)

	// Iterate over root blocks (each covers 2^13 = 8192 codepoints)
	for rootBlock := 0; rootBlock < 256; rootBlock++ {
		var midTable [64]byte

		baseCP := rootBlock << 13

		// Iterate over middle entries within this root block
		// Each middle entry covers 2^7 = 128 codepoints
		for midEntry := 0; midEntry < 64; midEntry++ {
			var leafTable [32]byte

			midBaseCP := baseCP + (midEntry << 7)

			// Pack 128 codepoints into 32 bytes (4 codepoints per byte, 2 bits each)
			for leafByte := 0; leafByte < 32; leafByte++ {
				var packed byte
				for bit := 0; bit < 4; bit++ {
					cp := midBaseCP + (leafByte << 2) + bit
					var w byte
					if cp <= maxCodepoint {
						w = widthMap[cp]
					}
					packed |= w << (2 * uint(bit))
				}
				leafTable[leafByte] = packed
			}

			// Deduplicate leaf table
			idx, ok := leafIndex[leafTable]
			if !ok {
				if len(leaves) > 255 {
					log.Fatalf("Too many unique leaf tables (%d > 255), cannot fit in uint8", len(leaves))
				}
				idx = byte(len(leaves))
				leafIndex[leafTable] = idx
				leaves = append(leaves, leafTable)
			}
			midTable[midEntry] = idx
		}

		// Deduplicate middle table
		idx, ok := midIndex[midTable]
		if !ok {
			if len(middle) > 255 {
				log.Fatalf("Too many unique middle tables (%d > 255), cannot fit in uint8", len(middle))
			}
			idx = byte(len(middle))
			midIndex[midTable] = idx
			middle = append(middle, midTable)
		}
		root[rootBlock] = idx
	}

	return root, middle, leaves
}

// generateGoFile generates the Go source file with both legacy and multi-stage tables.
func generateGoFile(wide, zeroWidth, ambiguous []runeRange, root *[256]byte, middle [][64]byte, leaves [][32]byte) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("warning: failed to close output file: %v", err)
		}
	}()

	w := bufio.NewWriter(file)

	// Write file header
	if _, err := fmt.Fprintf(w, `// Code generated by go generate; DO NOT EDIT.
// Generated from Unicode %s data files:
// - EastAsianWidth.txt
// - emoji-data.txt
//
// To regenerate:
//   go generate ./...

package uniwidth

// This file contains Unicode width tables for character width lookup.
//
// Two table formats are provided:
// 1. Legacy runeRange tables (used by Options API for ambiguous character handling)
// 2. Multi-stage lookup tables (used by tableLookupWidth for O(1) fallback)

`, unicodeVersion); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	// Write legacy wide table
	writeComment(w, "wideTableGenerated contains wide characters (width 2) not covered by hot paths.")
	writeComment(w, "These are characters with East Asian Width property W (Wide) or F (Fullwidth),")
	writeComment(w, "plus emoji characters not in the common emoji fast path.")
	writeComment(w, "Used by the Options API (binarySearchWidthInternal).")
	fmt.Fprint(w, "var wideTableGenerated = []runeRange{\n")
	for _, rr := range wide {
		fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last)
	}
	fmt.Fprint(w, "}\n\n")

	// Write legacy zero-width table
	writeComment(w, "zeroWidthTableGenerated contains zero-width characters not covered by hot paths.")
	writeComment(w, "These are control characters, combining marks, and format characters.")
	writeComment(w, "Used by the Options API (binarySearchWidthInternal).")
	fmt.Fprint(w, "var zeroWidthTableGenerated = []runeRange{\n")
	for _, rr := range zeroWidth {
		fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last)
	}
	fmt.Fprint(w, "}\n\n")

	// Write legacy ambiguous table
	writeComment(w, "ambiguousTableGenerated contains ambiguous-width characters.")
	writeComment(w, "These are characters with East Asian Width property A (Ambiguous).")
	writeComment(w, "Width depends on context: 2 in East Asian locales, 1 in neutral context.")
	writeComment(w, "Used by the Options API (binarySearchWidthInternal).")
	fmt.Fprint(w, "var ambiguousTableGenerated = []runeRange{\n")
	for _, rr := range ambiguous {
		fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last)
	}
	fmt.Fprint(w, "}\n\n")

	// Write multi-stage table documentation
	writeComment(w, "3-Stage Multi-Stage Lookup Table")
	writeComment(w, "")
	writeComment(w, "Splits a 21-bit Unicode codepoint into 3 parts:")
	writeComment(w, "  Codepoint: [20...13][12...7][6...0]")
	writeComment(w, "               8 bits   6 bits  7 bits")
	writeComment(w, "  Stage:       ROOT    MIDDLE   LEAF")
	writeComment(w, "")
	writeComment(w, "Lookup: widthRoot[cp>>13] -> midIdx")
	writeComment(w, "        widthMiddle[midIdx][cp>>7 & 0x3F] -> leafIdx")
	writeComment(w, "        widthLeaves[leafIdx][cp>>2 & 0x1F] >> (2*(cp&0x03)) & 0x03 -> width")
	writeComment(w, "")
	writeComment(w, "2-bit width encoding:")
	writeComment(w, "  0b00 = width 0 (control, combining, zero-width)")
	writeComment(w, "  0b01 = width 1 (narrow, default)")
	writeComment(w, "  0b10 = width 2 (wide: CJK, emoji, fullwidth)")
	writeComment(w, "  0b11 = ambiguous (width 1 in neutral context; 2 in East Asian)")
	fmt.Fprint(w, "\n")

	// Write root table
	fmt.Fprintf(w, "// widthRoot maps the top 8 bits of a codepoint (cp >> 13) to a middle table index.\n")
	fmt.Fprintf(w, "// Size: 256 bytes.\n")
	fmt.Fprint(w, "var widthRoot = [256]uint8{\n")
	for i := 0; i < 256; i += 16 {
		fmt.Fprint(w, "\t")
		for j := 0; j < 16; j++ {
			if j > 0 {
				fmt.Fprint(w, " ")
			}
			fmt.Fprintf(w, "0x%02X,", root[i+j])
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "}\n\n")

	// Write middle tables
	fmt.Fprintf(w, "// widthMiddle contains %d unique middle sub-tables.\n", len(middle))
	fmt.Fprintf(w, "// Each sub-table has 64 entries mapping bits [12:7] to a leaf table index.\n")
	fmt.Fprintf(w, "// Size: %d bytes.\n", len(middle)*64)
	fmt.Fprintf(w, "var widthMiddle = [%d][64]uint8{\n", len(middle))
	for i, mt := range middle {
		fmt.Fprintf(w, "\t// Middle table %d\n", i)
		fmt.Fprint(w, "\t{\n")
		for row := 0; row < 64; row += 16 {
			fmt.Fprint(w, "\t\t")
			end := row + 16
			if end > 64 {
				end = 64
			}
			for j := row; j < end; j++ {
				if j > row {
					fmt.Fprint(w, " ")
				}
				fmt.Fprintf(w, "0x%02X,", mt[j])
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\t},\n")
	}
	fmt.Fprint(w, "}\n\n")

	// Write leaf tables
	fmt.Fprintf(w, "// widthLeaves contains %d unique leaf sub-tables.\n", len(leaves))
	fmt.Fprintf(w, "// Each sub-table has 32 bytes of packed 2-bit width values (128 codepoints).\n")
	fmt.Fprintf(w, "// Size: %d bytes.\n", len(leaves)*32)
	fmt.Fprintf(w, "var widthLeaves = [%d][32]uint8{\n", len(leaves))
	for i, lt := range leaves {
		fmt.Fprintf(w, "\t// Leaf table %d\n", i)
		fmt.Fprint(w, "\t{\n")
		for row := 0; row < 32; row += 16 {
			fmt.Fprint(w, "\t\t")
			end := row + 16
			if end > 32 {
				end = 32
			}
			for j := row; j < end; j++ {
				if j > row {
					fmt.Fprint(w, " ")
				}
				fmt.Fprintf(w, "0x%02X,", lt[j])
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\t},\n")
	}
	fmt.Fprint(w, "}\n")

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

// writeComment writes a single-line Go comment to the writer.
func writeComment(w *bufio.Writer, line string) {
	if line == "" {
		fmt.Fprint(w, "//\n")
	} else {
		fmt.Fprintf(w, "// %s\n", line)
	}
}
