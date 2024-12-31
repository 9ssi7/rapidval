# RapidVal

[![Go Reference](https://pkg.go.dev/badge/github.com/9ssi7/rapidval.svg)](https://pkg.go.dev/github.com/9ssi7/rapidval)
[![Go Report Card](https://goreportcard.com/badge/github.com/9ssi7/rapidval)](https://goreportcard.com/report/github.com/9ssi7/rapidval)
[![Coverage Status](https://coveralls.io/repos/github/9ssi7/rapidval/badge.svg?branch=main)](https://coveralls.io/github/9ssi7/rapidval?branch=main)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

RapidVal is a high-performance, zero-dependency validation library for Go. It focuses on simplicity, extensibility, and performance while providing a clean and intuitive API.

## Philosophy

- **No Reflection**: Uses compile-time type safety instead of runtime reflection
- **No Dependencies**: Zero external dependencies for core functionality
- **Performance First**: Designed with performance in mind
- **Developer Friendly**: Clean API with excellent IDE support
- **Extensible**: Easy to add custom validation rules
- **i18n Ready**: Built-in translation support with customizable messages

## Features

- Type-safe validation rules
- Composable validation chains
- Custom validation rules support
- Built-in translation system
- Template-friendly error messages
- Comprehensive test coverage
- Detailed error reporting

## Installation

```bash
go get github.com/9ssi7/rapidval
```

## Quick Start

```go
package main

import (
	"time"

	"github.com/9ssi7/rapidval"
)

type Business struct {
	Title       string
	Description string
	FoundAt     time.Time
}

func (b *Business) Validations() rapidval.P {
	now := time.Now()
	return rapidval.P{
		rapidval.Required("Title", b.Title),
		rapidval.Required("Description", b.Description),
		rapidval.Required("FoundAt", b.FoundAt),
		rapidval.MinLength("Title", b.Title, 3),
		rapidval.MaxLength("Title", b.Title, 100),
		rapidval.MinLength("Description", b.Description, 10),
		rapidval.MaxLength("Description", b.Description, 1000),
		rapidval.DateGreaterThan("FoundAt", b.FoundAt, now),
		rapidval.DateLessThan("FoundAt", b.FoundAt, now.Add(24*time.Hour)),
	}
}

func main() {
	business := &Business{
		Title:       "RapidVal",
		Description: "RapidVal is a high-performance, zero-dependency validation library for Go.",
		FoundAt:     time.Now(),
	}
	v := rapidval.New()
	if err := v.Validate(business); err != nil {
		fmt.Println(err)
	}
}
```

## Translation Support

RapidVal comes with a built-in translation system that allows you to customize error messages. You can use the `NewTranslator` function to create a new translator with your own messages or use the `NewTranslatorWithMessages` function to create a new translator with predefined messages.

```go
tr := rapidval.NewTranslator()
```

You can also use the `Translate` method to translate an error message.

```go
err := &rapidval.ValidationError{
	Field:      "Title",
	MessageKey: rapidval.MsgRequired,
}
translated := tr.Translate(err)
fmt.Println(translated)
```

## Examples

You can find more examples in the [examples](examples) directory.

## Benchmarks

RapidVal is designed to be extremely fast and memory-efficient. Here are some benchmark results showing the performance of different operations:

```bash
goos: darwin
goarch: arm64
pkg: github.com/9ssi7/rapidval
cpu: Apple M3 Pro

# Single Validation Rules
BenchmarkValidations/Required-12                    10880550     129.0 ns/op    416 B/op    4 allocs/op
BenchmarkValidations/Email-12                      177856407       6.877 ns/op     0 B/op    0 allocs/op
BenchmarkValidations/MinLength-12                 1000000000       0.2787 ns/op    0 B/op    0 allocs/op
BenchmarkValidations/MaxLength-12                 1000000000       0.2796 ns/op    0 B/op    0 allocs/op
BenchmarkValidations/Between-12                   1000000000       0.2796 ns/op    0 B/op    0 allocs/op
BenchmarkValidations/DateGreaterThan-12            421514941       2.965 ns/op     0 B/op    0 allocs/op
BenchmarkValidations/DateLessThan-12               421147266       2.836 ns/op     0 B/op    0 allocs/op

# Multiple Validations
BenchmarkValidator/Multiple_Validations-12           2097082     590.5 ns/op   1447 B/op   16 allocs/op

# Translation Performance
BenchmarkTranslator/Translate-12                     3366672     356.2 ns/op    272 B/op    8 allocs/op

# Zero Value Checks
BenchmarkIsZero/string-12                         1000000000       0.7256 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/empty_string-12                   1000000000       0.7075 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/int-12                            1000000000       0.9177 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/zero_int-12                       1000000000       0.8901 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/bool-12                           1000000000       0.7769 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/time-12                           1000000000       0.8737 ns/op    0 B/op    0 allocs/op
BenchmarkIsZero/nil-12                            1000000000       0.5586 ns/op    0 B/op    0 allocs/op

PASS
ok      github.com/9ssi7/rapidval       17.819s
```

Key observations from the benchmarks:
- Most validation rules have zero allocations
- Single validations complete in nanoseconds
- Multiple validations are efficient even with 16 allocations
- Translation has minimal overhead
- Zero value checks are extremely fast with no allocations

These benchmarks were run on an Apple M3 Pro processor. Your results may vary depending on your hardware.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

```
Copyright 2024 RapidVal Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
