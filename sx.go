package sx

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Common separators used for splitting strings
var defaultSeparators = []rune{'-', '_', '/', '.', ' ', '\\'}

// isSeparator checks if a rune is a common separator
func isSeparator(r rune) bool {
	return slices.Contains(defaultSeparators, r)
}

// isSeparatorCustom checks if a rune is in the custom separator list
func isSeparatorCustom(r rune, separators []rune) bool {
	return slices.Contains(separators, r)
}

// isLetterCaseChange detects case transitions (like camelCase -> camel Case)
func isLetterCaseChange(prev, curr, next rune) bool {
	// Handle letter-to-letter case changes
	if unicode.IsLetter(prev) && unicode.IsLetter(curr) {
		// Lower to Upper transition (camelCase -> camel Case)
		if unicode.IsLower(prev) && unicode.IsUpper(curr) {
			return true
		}

		// Upper to Lower with next also lower (XMLHttpRequest -> XML Http Request)
		// But for cases like "FooBARb", we want: Foo-BA-Rb
		// So we split before the last uppercase in a sequence when followed by lowercase
		if unicode.IsUpper(prev) && unicode.IsUpper(curr) && unicode.IsLetter(next) && unicode.IsLower(next) {
			return true
		}
	}

	// Handle number to letter transitions (HTML5Parser -> HTML5 Parser)
	if (unicode.IsDigit(prev) || unicode.IsLetter(prev)) && unicode.IsLetter(curr) {
		if unicode.IsDigit(prev) && unicode.IsUpper(curr) {
			return true
		}
	}

	return false
}

// splitByCaseWithCustomSeparators splits a string into words with optional custom separators
func splitByCaseWithCustomSeparators(s string, customSeparators []rune) []string {
	if s == "" {
		return []string{}
	}

	runes := []rune(s)
	var words []string
	var currentWord strings.Builder

	for i, r := range runes {
		var prevRune, nextRune rune
		if i > 0 {
			prevRune = runes[i-1]
		}
		if i < len(runes)-1 {
			nextRune = runes[i+1]
		}

		// Check if we should start a new word
		shouldSplit := false
		skipCurrentRune := false

		// Check separators (custom if provided, otherwise default)
		var isSep bool
		if customSeparators != nil {
			// Custom separators specified - only split on those (could be empty list)
			isSep = isSeparatorCustom(r, customSeparators)
		} else {
			// No custom separators specified - use defaults
			isSep = isSeparator(r)
		}

		if isSep {
			// Skip separator and start new word
			shouldSplit = true
			skipCurrentRune = true
		} else if i > 0 && isLetterCaseChange(prevRune, r, nextRune) {
			// Case change detected
			shouldSplit = true
		}

		// If we should split, finalize current word (even if empty to handle consecutive separators)
		if shouldSplit {
			word := strings.TrimSpace(currentWord.String())
			// Always add the word, even if empty (for consecutive separators)
			words = append(words, word)
			currentWord.Reset()
		}

		// Add current rune to the word unless it's a separator we're skipping
		if !skipCurrentRune {
			currentWord.WriteRune(r)
		}
	}

	// Add the last word
	if currentWord.Len() > 0 {
		word := strings.TrimSpace(currentWord.String())
		words = append(words, word)
	}

	return words
}

// SplitOption configures how SplitByCase splits strings
type SplitOption func(*SplitConfig)

// SplitConfig holds the configuration for splitting behavior
type SplitConfig struct {
	Separators []rune
}

// defaultSplitConfig returns the default configuration
func defaultSplitConfig() *SplitConfig {
	return &SplitConfig{
		Separators: nil, // nil means use defaults
	}
}

// WithSeparators sets custom separator runes (replaces defaults)
func WithSeparators(separators ...rune) SplitOption {
	return func(c *SplitConfig) {
		c.Separators = make([]rune, len(separators))
		copy(c.Separators, separators)
	}
}

// SplitByCase splits a string into words based on case changes and separators
// Accepts optional configuration via functional options
func SplitByCase(s string, opts ...SplitOption) []string {
	config := defaultSplitConfig()
	for _, opt := range opts {
		opt(config)
	}

	return splitByCaseWithCustomSeparators(s, config.Separators)
}

// normalizeWord normalizes a word's case if needed
func normalizeWord(word string, normalize bool) string {
	if normalize {
		return strings.ToLower(word)
	}
	return word
}

// capitalizeWord capitalizes the first letter of a word
func capitalizeWord(word string) string {
	if word == "" {
		return word
	}

	r, size := utf8.DecodeRuneInString(word)
	if size == 0 {
		return word
	}

	return string(unicode.ToUpper(r)) + word[size:]
}

// joinWords joins words with a separator
func joinWords(words []string, separator string, transform func(string, int) string) string {
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i > 0 && separator != "" {
			result.WriteString(separator)
		}
		result.WriteString(transform(word, i))
	}

	return result.String()
}

type CaseOption func(*CaseConfig)

// CaseConfig configures case conversion behavior
type CaseConfig struct {
	Normalize bool // If an uppercase letter is followed by other uppercase letters (like FooBAR), they are preserved. You can use { normalize: true } for strictly following pascalCase convention.
}

// WithNormalize sets the normalize option
func WithNormalize(normalize bool) CaseOption {
	return func(c *CaseConfig) {
		c.Normalize = normalize
	}
}

// StringOrStringSlice represents input that can be either a string or slice of strings
type StringOrStringSlice interface {
	string | []string
}

// PascalCase converts input to PascalCase
func PascalCase[T StringOrStringSlice](input T, opts ...CaseOption) string {
	options := CaseConfig{}
	for _, opt := range opts {
		opt(&options)
	}

	switch v := any(input).(type) {
	case string:
		words := splitByCaseWithCustomSeparators(v, nil)
		result := joinWords(words, "", func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})

		return result
	case []string:
		result := joinWords(v, "", func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})

		return result
	default:
		return ""
	}
}
