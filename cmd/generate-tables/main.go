// generate-tables generates Unicode width tables from official Unicode 16.0 data.
//
// This tool downloads and parses:
// - EastAsianWidth.txt - East Asian Width property assignments
// - emoji-data.txt - Emoji presentation properties
//
// It generates optimized tables for uniwidth's tiered lookup strategy:
// - Tier 1-3 (hot paths) are hardcoded in uniwidth.go for O(1) lookup
// - This generates Tier 4 (binary search fallback) tables
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
	unicodeVersion       = "16.0.0"
	eastAsianWidthURL    = "https://www.unicode.org/Public/16.0.0/ucd/EastAsianWidth.txt"
	emojiDataURL         = "https://www.unicode.org/Public/16.0.0/ucd/emoji/emoji-data.txt"
	outputFile           = "tables_generated.go"
	outputFileWithHeader = "tables_generated.go"
)

// runeRange represents a contiguous range of runes with the same property.
type runeRange struct {
	first rune
	last  rune
}

// Category represents different width categories
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

	// Merge emoji into wide ranges
	wideRanges = mergeRanges(wideRanges, emojiRanges)

	// Generate zero-width tables (control chars, combining marks, format chars)
	log.Println("Generating zero-width tables...")
	zeroWidthRanges := generateZeroWidthRanges()

	// Filter out hot path ranges (already handled in uniwidth.go)
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
	err = generateGoFile(wideRanges, zeroWidthRanges, ambiguousRanges)
	if err != nil {
		log.Fatalf("Failed to generate Go file: %v", err)
	}

	log.Printf("Successfully generated %s with:", outputFile)
	log.Printf("  - Wide characters: %d ranges", len(wideRanges))
	log.Printf("  - Zero-width characters: %d ranges", len(zeroWidthRanges))
	log.Printf("  - Ambiguous characters: %d ranges", len(ambiguousRanges))
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
	// 0023          ; Emoji                # E0.0   [1] (#️)       number sign
	// 1F600..1F64F  ; Emoji                # E0.6  [80] (😀..🙏)    grinning face..folded hands
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
		// ZWJ (handled explicitly)
		{0x200D, 0x200D},
		// ZWNJ (handled explicitly)
		{0x200C, 0x200C},
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

// generateGoFile generates the Go source file with tables.
func generateGoFile(wide, zeroWidth, ambiguous []runeRange) error {
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

// This file contains Unicode width tables for characters NOT covered by
// the hot path tiers (Tier 1-3) in uniwidth.go.
//
// These tables are used as a fallback for rare characters that need
// binary search (Tier 4).

`, unicodeVersion); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	// Write wide table
	if _, err := fmt.Fprintf(w, "// wideTableGenerated contains wide characters (width 2) not covered by hot paths.\n"); err != nil {
		return fmt.Errorf("failed to write wide table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "// These are characters with East Asian Width property W (Wide) or F (Fullwidth),\n"); err != nil {
		return fmt.Errorf("failed to write wide table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "// plus emoji characters not in the common emoji fast path.\n"); err != nil {
		return fmt.Errorf("failed to write wide table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "var wideTableGenerated = []runeRange{\n"); err != nil {
		return fmt.Errorf("failed to write wide table declaration: %w", err)
	}
	for _, rr := range wide {
		if _, err := fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last); err != nil {
			return fmt.Errorf("failed to write wide table entry: %w", err)
		}
	}
	if _, err := fmt.Fprintf(w, "}\n\n"); err != nil {
		return fmt.Errorf("failed to close wide table: %w", err)
	}

	// Write zero-width table
	if _, err := fmt.Fprintf(w, "// zeroWidthTableGenerated contains zero-width characters not covered by hot paths.\n"); err != nil {
		return fmt.Errorf("failed to write zero-width table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "// These are control characters, combining marks, and format characters.\n"); err != nil {
		return fmt.Errorf("failed to write zero-width table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "var zeroWidthTableGenerated = []runeRange{\n"); err != nil {
		return fmt.Errorf("failed to write zero-width table declaration: %w", err)
	}
	for _, rr := range zeroWidth {
		if _, err := fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last); err != nil {
			return fmt.Errorf("failed to write zero-width table entry: %w", err)
		}
	}
	if _, err := fmt.Fprintf(w, "}\n\n"); err != nil {
		return fmt.Errorf("failed to close zero-width table: %w", err)
	}

	// Write ambiguous table
	if _, err := fmt.Fprintf(w, "// ambiguousTableGenerated contains ambiguous-width characters.\n"); err != nil {
		return fmt.Errorf("failed to write ambiguous table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "// These are characters with East Asian Width property A (Ambiguous).\n"); err != nil {
		return fmt.Errorf("failed to write ambiguous table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "// Width depends on context: 2 in East Asian locales, 1 in neutral context.\n"); err != nil {
		return fmt.Errorf("failed to write ambiguous table comment: %w", err)
	}
	if _, err := fmt.Fprintf(w, "var ambiguousTableGenerated = []runeRange{\n"); err != nil {
		return fmt.Errorf("failed to write ambiguous table declaration: %w", err)
	}
	for _, rr := range ambiguous {
		if _, err := fmt.Fprintf(w, "\t{0x%04X, 0x%04X},\n", rr.first, rr.last); err != nil {
			return fmt.Errorf("failed to write ambiguous table entry: %w", err)
		}
	}
	if _, err := fmt.Fprintf(w, "}\n"); err != nil {
		return fmt.Errorf("failed to close ambiguous table: %w", err)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}
