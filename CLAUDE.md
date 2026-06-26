# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`fatcontext` is a Go linter that detects nested ("fat") contexts created inside
loops, function literals, and (optionally) struct pointers — a known performance
pitfall (https://gabnotes.org/fat-contexts/). It is distributed standalone and is
also vendored into `golangci-lint` (since v1.58.0).

## Commands

- Run tests: `mise run test` (`go test -race -v ./...`)
- Run a single test: `go test -race -run TestSuggestedFixes ./pkg/analyzer/`
- Lint (golangci-lint with autofix): `mise run lint`
- Run the full hook suite (file hooks + golangci-lint + tests): `prek run --all-files`
- Build the CLI: `go build ./cmd/fatcontext`
- Release: push a git tag (`git tag -a vX.Y.Z -m "vX.Y.Z" && git push --follow-tags`),
  which triggers GoReleaser via GitHub Actions.

`CGO_ENABLED=1` is set via `mise.toml` for local dev (needed by the `cgo` test fixture).

## Architecture

The linter is built on the standard `golang.org/x/tools/go/analysis` framework.

- `cmd/fatcontext/main.go` — CLI entrypoint; wraps the analyzer in `singlechecker.Main`.
- `pkg/analyzer/analyzer.go` — the entire analysis logic:
  - `NewAnalyzer()` builds the `*analysis.Analyzer`, registering the
    `check-struct-pointers` flag (constant `FlagCheckStructPointers`, exported so
    golangci-lint can set it).
  - `runner.run` uses the `inspect` pass to walk `ForStmt`, `RangeStmt`, `FuncLit`,
    and `FuncDecl` nodes. For each, it scans the body for a context reassignment.
  - `findNestedContext` recurses through nested statement lists (if/switch/select/
    blocks via `getStmtList`) looking for an `*ast.AssignStmt` to a
    `context.Context` using `=` (not `:=`). It tracks variables reset to an empty
    context (`isEmptyContext`: `context.Background()`/`TODO()`, or `testing.{T,B,TB}.Context()`)
    to avoid false positives, and skips non-pointer values declared within the loop.
  - `getCategory` classifies findings: in-loop, in-func-literal, in-struct-pointer,
    or unsupported. Struct-pointer reports are suppressed unless the flag is on.
  - `getSuggestedFixes` proposes replacing `=` with `:=` (no fix for struct-pointer
    or unsupported categories).

## Testing

Tests use `analysistest` with fixtures under `pkg/analyzer/testdata/src/<dir>/`.
Expected diagnostics are encoded as `// want "..."` comments in the fixture files;
suggested fixes are validated against `.golden` files via `RunWithSuggestedFixes`.
To add a case: create/extend a fixture dir, add `// want` comments, and register
the dir in the `testCases` table in `analyzer_test.go` (with options if the
struct-pointer flag is needed).
