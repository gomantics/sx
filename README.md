# sx

[![Go Reference](https://pkg.go.dev/badge/github.com/gomantics/sx.svg)](https://pkg.go.dev/github.com/gomantics/sx) [![CI](https://github.com/gomantics/sx/actions/workflows/ci.yml/badge.svg)](https://github.com/gomantics/sx/actions/workflows/ci.yml)

A simple Go library for string case conversion. Converts between camelCase, PascalCase, kebab-case, snake_case, and more.

## Installation

```bash
go get github.com/gomantics/sx
```

## Documentation

For complete documentation, usage examples, and API reference, visit:

**https://gomantics.dev/sx**

## Quick Example

```go
import "github.com/gomantics/sx"

sx.CamelCase("hello-world")   // helloWorld
sx.PascalCase("hello_world")  // HelloWorld
sx.KebabCase("HelloWorld")    // hello-world
sx.SnakeCase("HelloWorld")    // hello_world
```

## Acknowledgements

This library is highly inspired by [scule](https://github.com/unjs/scule) - a fantastic JavaScript string case utility library by the UnJS team.

## License

MIT
