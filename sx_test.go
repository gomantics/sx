package sx_test

import (
	"reflect"
	"testing"

	"github.com/gomantics/sx"
)

func TestSplitByCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "camelCase",
			input:    "camelCase",
			expected: []string{"camel", "Case"},
		},
		{
			name:     "PascalCase",
			input:    "PascalCase",
			expected: []string{"Pascal", "Case"},
		},
		{
			name:     "kebab-case",
			input:    "kebab-case",
			expected: []string{"kebab", "case"},
		},
		{
			name:     "snake_case",
			input:    "snake_case",
			expected: []string{"snake", "case"},
		},
		{
			name:     "dot.case",
			input:    "dot.case",
			expected: []string{"dot", "case"},
		},
		{
			name:     "slash/case",
			input:    "slash/case",
			expected: []string{"slash", "case"},
		},
		{
			name:     "XMLHttpRequest",
			input:    "XMLHttpRequest",
			expected: []string{"XML", "Http", "Request"},
		},
		{
			name:     "IOError",
			input:    "IOError",
			expected: []string{"IO", "Error"},
		},
		{
			name:     "iPhone",
			input:    "iPhone",
			expected: []string{"i", "Phone"},
		},
		{
			name:     "hello--world",
			input:    "hello--world",
			expected: []string{"hello", "", "world"},
		},
		{
			name:     "hello\\World.Foo-Barb",
			input:    "hello\\World.Foo-Bar",
			expected: []string{"hello", "World", "Foo", "Bar"},
		},
		{
			name:     "mixed_caseWith-different.separators",
			input:    "mixed_caseWith-different.separators",
			expected: []string{"mixed", "case", "With", "different", "separators"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single word",
			input:    "word",
			expected: []string{"word"},
		},
		{
			name:     "ALLCAPS",
			input:    "ALLCAPS",
			expected: []string{"ALLCAPS"},
		},
		{
			name:     "alllowercase",
			input:    "alllowercase",
			expected: []string{"alllowercase"},
		},
		{
			name:     "HTML5Parser",
			input:    "HTML5Parser",
			expected: []string{"HTML5", "Parser"},
		},
		{
			name:     "spaces in string",
			input:    "spaces in string",
			expected: []string{"spaces", "in", "string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.SplitByCase(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitByCase(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitByCase_CustomSeparators(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		separators []rune
		expected   []string
	}{
		{
			name:       "custom separators - only comma",
			input:      "hello,world-test_case",
			separators: []rune{','},
			expected:   []string{"hello", "world-test_case"},
		},
		{
			name:       "custom separators - pipe and ampersand",
			input:      "hello|world&test",
			separators: []rune{'|', '&'},
			expected:   []string{"hello", "world", "test"},
		},
		{
			name:       "no separators - only case splitting",
			input:      "hello-world_testCase",
			separators: []rune{},
			expected:   []string{"hello-world_test", "Case"},
		},
		{
			name:       "custom separators - empty list",
			input:      "hello-world_testCase",
			separators: []rune{},
			expected:   []string{"hello-world_test", "Case"},
		},
		{
			name:       "default behavior - no custom separators",
			input:      "camelCaseExample-with_separators",
			separators: nil,
			expected:   []string{"camel", "Case", "Example", "with", "separators"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []string
			if tt.separators == nil {
				result = sx.SplitByCase(tt.input)
			} else {
				result = sx.SplitByCase(tt.input, sx.WithSeparators(tt.separators...))
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitByCase(%q, %v) = %v, want %v",
					tt.input, tt.separators, result, tt.expected)
			}
		})
	}
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		options  []sx.CaseOption
	}{
		{
			name:     "camelCase to PascalCase",
			input:    "camelCase",
			expected: "CamelCase",
		},
		{
			name:     "kebab-case to PascalCase",
			input:    "kebab-case",
			expected: "KebabCase",
		},
		{
			name:     "snake_case to PascalCase",
			input:    "snake_case",
			expected: "SnakeCase",
		},
		{
			name:     "mixed.case_with-separators",
			input:    "mixed.case_with-separators",
			expected: "MixedCaseWithSeparators",
		},
		{
			name:     "XMLHttpRequest",
			input:    "XMLHttpRequest",
			expected: "XMLHttpRequest",
		},
		{
			name:     "XMLHttpRequest normalized",
			input:    "XMLHttpRequest",
			expected: "XmlHttpRequest",
			options:  []sx.CaseOption{sx.WithNormalize(true)},
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "word",
			expected: "Word",
		},
		{
			name:     "hello--world-42",
			input:    "hello--world-42",
			expected: "HelloWorld42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.PascalCase(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("PascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPascalCaseWithSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
		options  []sx.CaseOption
	}{
		{
			name:     "string slice",
			input:    []string{"hello", "world", "test"},
			expected: "HelloWorldTest",
		},
		{
			name:     "string slice normalized",
			input:    []string{"HELLO", "WORLD", "TEST"},
			expected: "HelloWorldTest",
			options:  []sx.CaseOption{sx.WithNormalize(true)},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single item slice",
			input:    []string{"word"},
			expected: "Word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.PascalCase(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("PascalCase(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
