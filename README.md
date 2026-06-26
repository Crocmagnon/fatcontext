# fatcontext

[![Go Reference](https://pkg.go.dev/badge/github.com/Crocmagnon/fatcontext.svg)](https://pkg.go.dev/github.com/Crocmagnon/fatcontext)
[![Go Report Card](https://goreportcard.com/badge/github.com/Crocmagnon/fatcontext)](https://goreportcard.com/report/github.com/Crocmagnon/fatcontext)
[![Go Coverage](https://github.com/Crocmagnon/fatcontext/wiki/coverage.svg)](https://github.com/Crocmagnon/fatcontext/wiki/Coverage)

`fatcontext` is a Go linter which detects potential fat contexts in loops or function literals.
They can lead to performance issues, as documented here: https://gabnotes.org/fat-contexts/

## Installation / usage

`fatcontext` is available in `golangci-lint` since v1.58.0.

```bash
go install go.augendre.info/fatcontext/cmd/fatcontext@latest
fatcontext ./...
```

## Example

```go
package main

import "context"

func ok() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx := context.WithValue(ctx, "key", i)
		_ = ctx
	}
}

func notOk() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i) // "nested context in loop"
		_ = ctx
	}
}
```

## Configuration

The linter exposes the following options, all available on the command line:

| Option                     | Default    | Description                                                              |
|----------------------------|------------|--------------------------------------------------------------------------|
| `check-loops`              | `true`     | Detect fat contexts created inside loops.                                |
| `check-function-literals`  | `true`     | Detect fat contexts created inside function literals.                    |
| `check-struct-pointers`    | `false`    | Detect potential fat contexts created through struct pointers.           |

```bash
# Disable loop detection (enabled by default)
fatcontext -check-loops=false ./...

# Disable function literal detection (enabled by default)
fatcontext -check-function-literals=false ./...

# Enable struct pointer detection (disabled by default)
fatcontext -check-struct-pointers ./...
```

When used through `golangci-lint`, these options are configured in the linter's
settings rather than on the command line.

## Development

Setup pre-commit locally:
```bash
pre-commit install
```

Run tests & linter:
```bash
make lint test
```

To release, just publish a git tag:
```bash
git tag -a v0.1.0 -m "v0.1.0"
git push --follow-tags
```
