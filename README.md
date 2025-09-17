# sx

[![Go Reference](https://pkg.go.dev/badge/github.com/gomantics/sx.svg)](https://pkg.go.dev/github.com/gomantics/sx) [![CI](https://github.com/gomantics/sx/actions/workflows/ci.yml/badge.svg)](https://github.com/gomantics/sx/actions/workflows/ci.yml)

A simple Go library for string case conversion. Converts between camelCase, PascalCase, kebab-case, snake_case, and more.

## Installation

```bash
go get github.com/gomantics/sx
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/gomantics/sx"
)

func main() {
	// Convert to different cases
	fmt.Println(sx.CamelCase("hello-world"))  // helloWorld
	fmt.Println(sx.PascalCase("hello_world")) // HelloWorld
	fmt.Println(sx.KebabCase("HelloWorld"))   // hello-world
	fmt.Println(sx.SnakeCase("HelloWorld"))   // hello_world
	fmt.Println(sx.TrainCase("hello-world"))  // Hello-World
	fmt.Println(sx.FlatCase("hello-world"))   // helloworld

	// Works with mixed separators and cases
	fmt.Println(sx.CamelCase("mixed_caseWith-different.separators")) // mixedCaseWithDifferentSeparators

	// Handle complex acronyms
	fmt.Println(sx.KebabCase("XMLHttpRequest"))                      // xml-http-request
	fmt.Println(sx.CamelCase("HTML5Parser", sx.WithNormalize(true))) // html5Parser
}
```

## Case Functions

- `CamelCase()` - converts to camelCase
- `PascalCase()` - converts to PascalCase
- `KebabCase()` - converts to kebab-case
- `SnakeCase()` - converts to snake_case
- `TrainCase()` - converts to Train-Case
- `FlatCase()` - converts to flatcase (no separators)

## Utilities

- `SplitByCase()` - splits strings into words by case changes and separators
- `UpperFirst()` - capitalizes first character
- `LowerFirst()` - lowercases first character

## Options

Some functions support options for customization:

```go
// Normalize case for strict PascalCase/camelCase
sx.PascalCase("XMLHttpRequest", sx.WithNormalize(true))  // XmlHttpRequest

// Custom separators for splitting
sx.SplitByCase("hello.world", sx.WithSeparators('.'))  // ["hello", "world"]

// Custom separator for kebab case
sx.KebabCase("hello world", "|")  // hello|world
```

## Contributing

Contributors are always welcome! Feel free to raise a PR or create an issue first.

## Acknowledgements

This library is highly inspired by [scule](https://github.com/unjs/scule) - a fantastic JavaScript string case utility library by the UnJS team.

## License

MIT
