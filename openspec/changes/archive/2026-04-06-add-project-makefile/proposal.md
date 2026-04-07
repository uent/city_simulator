## Why

The project lacks a Makefile, forcing developers to remember long `go` CLI commands with multiple flags. A Makefile centralizes all common operations and self-documents the project's workflow.

## What Changes

- Add a `Makefile` at the project root with targets for building, running, testing, formatting, and cleaning
- Include inline examples (via `@echo` comments) in the Makefile itself so developers can learn usage without reading external docs
- Provide convenience targets with sensible defaults matching the simulator's existing flag defaults

## Capabilities

### New Capabilities

- `project-makefile`: A root-level Makefile exposing all common dev/run/build/test/clean tasks with inline usage examples

### Modified Capabilities

<!-- No existing spec-level requirements change -->

## Impact

- New file: `Makefile` at project root
- No code changes to Go source files
- No API or dependency changes
- Improves DX for local development and CI usage
