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

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		options  []sx.CaseOption
	}{
		{
			name:     "PascalCase to camelCase",
			input:    "PascalCase",
			expected: "pascalCase",
		},
		{
			name:     "kebab-case to camelCase",
			input:    "kebab-case",
			expected: "kebabCase",
		},
		{
			name:     "snake_case to camelCase",
			input:    "snake_case",
			expected: "snakeCase",
		},
		{
			name:     "XMLHttpRequest",
			input:    "XMLHttpRequest",
			expected: "xMLHttpRequest",
		},
		{
			name:     "XMLHttpRequest normalized",
			input:    "XMLHttpRequest",
			expected: "xmlHttpRequest",
			options:  []sx.CaseOption{sx.WithNormalize(true)},
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "Word",
			expected: "word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.CamelCase(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("CamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCamelCaseWithSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
		options  []sx.CaseOption
	}{
		{
			name:     "string slice",
			input:    []string{"hello", "world", "test"},
			expected: "helloWorldTest",
		},
		{
			name:     "string slice normalized",
			input:    []string{"HELLO", "WORLD", "TEST"},
			expected: "helloWorldTest",
			options:  []sx.CaseOption{sx.WithNormalize(true)},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single item slice",
			input:    []string{"Word"},
			expected: "word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.CamelCase(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("CamelCase(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		separator string
	}{
		{
			name:     "camelCase to kebab-case",
			input:    "camelCase",
			expected: "camel-case",
		},
		{
			name:     "PascalCase to kebab-case",
			input:    "PascalCase",
			expected: "pascal-case",
		},
		{
			name:     "snake_case to kebab-case",
			input:    "snake_case",
			expected: "snake-case",
		},
		{
			name:     "XMLHttpRequest to kebab-case",
			input:    "XMLHttpRequest",
			expected: "xml-http-request",
		},
		{
			name:      "custom separator",
			input:     "camelCase",
			expected:  "camel|case",
			separator: "|",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "word",
			expected: "word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string

			if tt.separator != "" {
				result = sx.KebabCase(tt.input, tt.separator)
			} else {
				result = sx.KebabCase(tt.input)
			}

			if result != tt.expected {
				t.Errorf("KebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase to snake_case",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "PascalCase to snake_case",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "kebab-case to snake_case",
			input:    "kebab-case",
			expected: "kebab_case",
		},
		{
			name:     "XMLHttpRequest to snake_case",
			input:    "XMLHttpRequest",
			expected: "xml_http_request",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "word",
			expected: "word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.SnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTrainCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		options  []sx.CaseOption
	}{
		{
			name:     "camelCase to Train-Case",
			input:    "camelCase",
			expected: "Camel-Case",
		},
		{
			name:     "snake_case to Train-Case",
			input:    "snake_case",
			expected: "Snake-Case",
		},
		{
			name:     "XMLHttpRequest to Train-Case",
			input:    "XMLHttpRequest",
			expected: "XML-Http-Request",
		},
		{
			name:     "XMLHttpRequest normalized",
			input:    "XMLHttpRequest",
			expected: "Xml-Http-Request",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.TrainCase(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("TrainCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFlatCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase to flatcase",
			input:    "camelCase",
			expected: "camelcase",
		},
		{
			name:     "PascalCase to flatcase",
			input:    "PascalCase",
			expected: "pascalcase",
		},
		{
			name:     "kebab-case to flatcase",
			input:    "kebab-case",
			expected: "kebabcase",
		},
		{
			name:     "XMLHttpRequest to flatcase",
			input:    "XMLHttpRequest",
			expected: "xmlhttprequest",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "Word",
			expected: "word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sx.FlatCase(tt.input)
			if result != tt.expected {
				t.Errorf("FlatCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		function func(string) string
		expected string
	}{
		{
			name:     "unicode characters",
			input:    "helloWörld",
			function: func(s string) string { return sx.CamelCase(s) },
			expected: "helloWörld",
		},
		{
			name:     "numbers in string",
			input:    "html5Parser",
			function: func(s string) string { return sx.PascalCase(s) },
			expected: "Html5Parser",
		},
		{
			name:     "consecutive uppercase",
			input:    "HTTPSConnection",
			function: func(s string) string { return sx.KebabCase(s) },
			expected: "https-connection",
		},
		{
			name:     "mixed separators",
			input:    "hello_world-test.case/example",
			function: func(s string) string { return sx.CamelCase(s) },
			expected: "helloWorldTestCaseExample",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.input)
			if result != tt.expected {
				t.Errorf("Function(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
