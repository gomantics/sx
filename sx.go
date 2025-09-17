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
func joinWords(words []string, separator string, preserveEmpty bool, transform func(string, int) string) string {
	if len(words) == 0 {
		return ""
	}

	// Filter out empty words if not preserving them
	wordsToUse := words
	if !preserveEmpty {
		var filteredWords []string
		for _, word := range words {
			if word != "" {
				filteredWords = append(filteredWords, word)
			}
		}
		wordsToUse = filteredWords
	}

	var result strings.Builder
	for i, word := range wordsToUse {
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
	// If an uppercase letter is followed by other uppercase letters (like FooBAR), they are preserved. You can use sx.WithNormalize(true) for strictly following PascalCase convention.
	Normalize bool
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
		result := joinWords(words, "", false, func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})

		return result
	case []string:
		result := joinWords(v, "", false, func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})

		return result
	default:
		return ""
	}
}

// lowercaseWord converts the first letter to lowercase
func lowercaseWord(word string) string {
	if word == "" {
		return word
	}

	r, size := utf8.DecodeRuneInString(word)
	if size == 0 {
		return word
	}

	return string(unicode.ToLower(r)) + word[size:]
}

// CamelCase converts input to camelCase
func CamelCase[T StringOrStringSlice](input T, opts ...CaseOption) string {
	switch v := any(input).(type) {
	case string:
		pascalCase := PascalCase(v, opts...)
		return lowercaseWord(pascalCase)
	case []string:
		if len(v) == 0 {
			return ""
		}

		options := CaseConfig{}
		for _, opt := range opts {
			opt(&options)
		}

		result := joinWords(v, "", false, func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			if i == 0 {
				return lowercaseWord(normalized)
			}

			return capitalizeWord(normalized)
		})
		return result
	default:
		return ""
	}
}

// KebabCase converts input to kebab-case
func KebabCase[T StringOrStringSlice](input T, separator ...string) string {
	sep := "-"
	if len(separator) > 0 {
		sep = separator[0]
	}

	switch v := any(input).(type) {
	case string:
		words := splitByCaseWithCustomSeparators(v, nil)
		result := joinWords(words, sep, true, func(word string, i int) string {
			return strings.ToLower(word)
		})
		return result
	case []string:
		result := joinWords(v, sep, true, func(word string, i int) string {
			return strings.ToLower(word)
		})
		return result
	default:
		return ""
	}
}

// SnakeCase converts input to snake_case
func SnakeCase[T StringOrStringSlice](input T) string {
	return KebabCase(input, "_")
}

// TrainCase converts input to Train-Case
func TrainCase[T StringOrStringSlice](input T, opts ...CaseOption) string {
	options := CaseConfig{}
	for _, opt := range opts {
		opt(&options)
	}

	switch v := any(input).(type) {
	case string:
		words := splitByCaseWithCustomSeparators(v, nil)
		result := joinWords(words, "-", false, func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})
		return result
	case []string:
		result := joinWords(v, "-", false, func(word string, i int) string {
			normalized := normalizeWord(word, options.Normalize)
			return capitalizeWord(normalized)
		})
		return result
	default:
		return ""
	}
}

// FlatCase converts input to flatcase (no separators)
func FlatCase[T StringOrStringSlice](input T) string {
	return KebabCase(input, "")
}

// UpperFirst converts the first character to uppercase
func UpperFirst(s string) string {
	return capitalizeWord(s)
}

// LowerFirst converts the first character to lowercase
func LowerFirst(s string) string {
	return lowercaseWord(s)
}
