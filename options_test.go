package uniwidth

import (
	"testing"
)

// TestRuneWidthWithOptions_EastAsianAmbiguous tests handling of ambiguous characters.
func TestRuneWidthWithOptions_EastAsianAmbiguous(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		eaWidth  EAWidth
		expected int
	}{
		// Ambiguous characters - should respect EAWidth setting
		{"± narrow", '±', EANarrow, 1},
		{"± wide", '±', EAWide, 2},
		{"½ narrow", '½', EANarrow, 1},
		{"½ wide", '½', EAWide, 2},
		{"° narrow", '°', EANarrow, 1},
		{"° wide", '°', EAWide, 2},
		{"× narrow", '×', EANarrow, 1},
		{"× wide", '×', EAWide, 2},
		{"÷ narrow", '÷', EANarrow, 1},
		{"÷ wide", '÷', EAWide, 2},

		// Non-ambiguous characters - should be unaffected
		{"ASCII a narrow", 'a', EANarrow, 1},
		{"ASCII a wide", 'a', EAWide, 1},
		{"CJK 世 narrow", '世', EANarrow, 2},
		{"CJK 世 wide", '世', EAWide, 2},
		{"Emoji 😀 narrow", '😀', EANarrow, 2},
		{"Emoji 😀 wide", '😀', EAWide, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidthWithOptions(tt.r, WithEastAsianAmbiguous(tt.eaWidth))
			if got != tt.expected {
				t.Errorf("RuneWidthWithOptions(%U, EAWidth=%d) = %d, want %d", tt.r, tt.eaWidth, got, tt.expected)
			}
		})
	}
}

// TestStringWidthWithOptions_EastAsianAmbiguous tests string width with ambiguous characters.
func TestStringWidthWithOptions_EastAsianAmbiguous(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		eaWidth  EAWidth
		expected int
	}{
		// Neutral locale (narrow)
		{"Hello narrow", "Hello", EANarrow, 5},
		{"±½ narrow", "±½", EANarrow, 2},
		{"Hello ±½ narrow", "Hello ±½", EANarrow, 8}, // Hello=5, space=1, ±=1, ½=1

		// East Asian locale (wide)
		{"Hello wide", "Hello", EAWide, 5},
		{"±½ wide", "±½", EAWide, 4},
		{"Hello ±½ wide", "Hello ±½", EAWide, 10}, // Hello=5, space=1, ±=2, ½=2

		// Mixed content
		{"CJK + ambiguous narrow", "你好±", EANarrow, 5}, // 你=2, 好=2, ±=1
		{"CJK + ambiguous wide", "你好±", EAWide, 6},     // 你=2, 好=2, ±=2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidthWithOptions(tt.s, WithEastAsianAmbiguous(tt.eaWidth))
			if got != tt.expected {
				t.Errorf("StringWidthWithOptions(%q, EAWidth=%d) = %d, want %d", tt.s, tt.eaWidth, got, tt.expected)
			}
		})
	}
}

// TestOptions_Default tests default option values.
func TestOptions_Default(t *testing.T) {
	// Test that defaults match non-options functions
	ambiguous := '±'

	defaultWidth := RuneWidth(ambiguous)
	optionsWidth := RuneWidthWithOptions(ambiguous) // No options = use defaults

	if defaultWidth != optionsWidth {
		t.Errorf("Default options differ from RuneWidth: RuneWidth=%d, RuneWidthWithOptions=%d", defaultWidth, optionsWidth)
	}

	// Default should be EANarrow (width 1 for ambiguous)
	expected := 1
	if optionsWidth != expected {
		t.Errorf("Default ambiguous width should be %d, got %d", expected, optionsWidth)
	}
}

// TestOptions_MultipleOptions tests combining multiple options.
func TestOptions_MultipleOptions(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		opts     []Option
		expected int
	}{
		{
			name: "Both options",
			s:    "Hello ±",
			opts: []Option{
				WithEastAsianAmbiguous(EAWide),
				WithEmojiPresentation(true),
			},
			expected: 8, // Hello=5, space=1, ±=2
		},
		{
			name:     "No options",
			s:        "Hello ±",
			opts:     []Option{},
			expected: 7, // Hello=5, space=1, ±=1 (default narrow)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidthWithOptions(tt.s, tt.opts...)
			if got != tt.expected {
				t.Errorf("StringWidthWithOptions(%q, %d opts) = %d, want %d", tt.s, len(tt.opts), got, tt.expected)
			}
		})
	}
}

// TestOptions_BackwardCompatibility ensures default functions still work.
func TestOptions_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		s        string
		expected int
	}{
		{"Hello", 5},
		{"Hello 世界", 10},
		{"😀", 2},
		{"Hello ± World", 13}, // With default narrow ambiguous (± is width 1)
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			// Test that StringWidth still works (backward compatibility)
			got := StringWidth(tt.s)
			if got != tt.expected {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.expected)
			}

			// Test that it matches WithOptions with no options
			gotWithOptions := StringWidthWithOptions(tt.s)
			if gotWithOptions != got {
				t.Errorf("StringWidthWithOptions(%q) = %d, want %d (to match StringWidth)", tt.s, gotWithOptions, got)
			}
		})
	}
}

// TestOptions_GreekAndCyrillic tests ambiguous Greek and Cyrillic characters.
func TestOptions_GreekAndCyrillic(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		eaWidth  EAWidth
		expected int
	}{
		// Greek characters (ambiguous in East Asian Width)
		{"α narrow", 'α', EANarrow, 1},
		{"α wide", 'α', EAWide, 2},
		{"β narrow", 'β', EANarrow, 1},
		{"β wide", 'β', EAWide, 2},

		// Cyrillic characters (ambiguous)
		{"А narrow", 'А', EANarrow, 1}, // Cyrillic A
		{"А wide", 'А', EAWide, 2},
		{"Я narrow", 'Я', EANarrow, 1}, // Cyrillic Ya
		{"Я wide", 'Я', EAWide, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidthWithOptions(tt.r, WithEastAsianAmbiguous(tt.eaWidth))
			if got != tt.expected {
				t.Errorf("RuneWidthWithOptions(%U %c, EAWidth=%d) = %d, want %d", tt.r, tt.r, tt.eaWidth, got, tt.expected)
			}
		})
	}
}

// TestOptions_BoxDrawing tests box-drawing characters (often ambiguous).
func TestOptions_BoxDrawing(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		eaWidth  EAWidth
		expected int
	}{
		// Box-drawing characters
		{"─ narrow", "─", EANarrow, 1},
		{"─ wide", "─", EAWide, 2},
		{"│ narrow", "│", EANarrow, 1},
		{"│ wide", "│", EAWide, 2},
		{"┌ narrow", "┌", EANarrow, 1},
		{"┌ wide", "┌", EAWide, 2},

		// Box-drawing table
		{"table narrow", "┌─┐", EANarrow, 3},
		{"table wide", "┌─┐", EAWide, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringWidthWithOptions(tt.s, WithEastAsianAmbiguous(tt.eaWidth))
			if got != tt.expected {
				t.Errorf("StringWidthWithOptions(%q, EAWidth=%d) = %d, want %d", tt.s, tt.eaWidth, got, tt.expected)
			}
		})
	}
}

// BenchmarkRuneWidthWithOptions benchmarks the options API performance.
func BenchmarkRuneWidthWithOptions(b *testing.B) {
	opts := []Option{WithEastAsianAmbiguous(EAWide)}

	b.Run("ASCII", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RuneWidthWithOptions('a', opts...)
		}
	})

	b.Run("Ambiguous", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RuneWidthWithOptions('±', opts...)
		}
	})

	b.Run("CJK", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RuneWidthWithOptions('世', opts...)
		}
	})
}

// BenchmarkStringWidthWithOptions benchmarks string width with options.
func BenchmarkStringWidthWithOptions(b *testing.B) {
	opts := []Option{WithEastAsianAmbiguous(EAWide)}

	b.Run("ASCII", func(b *testing.B) {
		s := "Hello, World!"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			StringWidthWithOptions(s, opts...)
		}
	})

	b.Run("Ambiguous", func(b *testing.B) {
		s := "Hello ±½"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			StringWidthWithOptions(s, opts...)
		}
	})

	b.Run("Mixed", func(b *testing.B) {
		s := "Hello 世界 ±½"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			StringWidthWithOptions(s, opts...)
		}
	})
}
