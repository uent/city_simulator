## Why

The project currently uses Go 1.24.1. Updating to Go 1.26.1 ensures we have the latest language features, standard library improvements, performance gains, and security patches. Keeping the toolchain current reduces technical debt and ensures compatibility with the most recent Go ecosystem.

## What Changes

- Update the `go` directive in `go.mod` from `1.24.1` to `1.26.1`
- Verify that all existing dependencies and code compile cleanly under Go 1.26.1

## Capabilities

### Modified Capabilities

- `go.mod`: The module file's Go version directive is bumped to `1.26.1`, enabling any new language or standard library features introduced between 1.24.1 and 1.26.1

### New Capabilities

<!-- None — this is a toolchain upgrade only -->

## Impact

- Single-file change: `go.mod`
- No API or behavioral changes to the application
- All existing packages (`internal/character`, `internal/llm`, `internal/simulation`, `internal/world`, `internal/conversation`, `cmd/simulator`) must compile without errors under the new toolchain
- Developers will need Go 1.26.1 installed locally to build the project
