## Context

The project is a Go-based city simulator (`github.com/jnn-z/city_simulator`) that currently declares `go 1.24.1` in its `go.mod`. This is a minimal, focused toolchain upgrade with no application logic changes.

## Goals / Non-Goals

**Goals:**
- Bump `go.mod` to `go 1.26.1`
- Confirm the project builds cleanly under Go 1.26.1
- Confirm existing tests (if any) pass

**Non-Goals:**
- Adopting new Go 1.26.x language features in existing code
- Upgrading third-party dependencies (e.g., `gopkg.in/yaml.v3`)
- Refactoring any existing code

## Decisions

### 1. Edit `go.mod` directly

Change the single `go` directive line:

```
// Before
go 1.24.1

// After
go 1.26.1
```

**Rationale:** This is the only required change. The Go toolchain reads this directive to enforce minimum version compatibility. No other files reference the Go version.

**Alternative considered:** Running `go mod tidy` after the edit — this is safe to do and may update checksums, but is not required unless dependencies change. Include it as a verification step.

### 2. Verification via `go build`

After editing `go.mod`, run `go build ./...` to confirm all packages compile without errors under Go 1.26.1.

**Rationale:** Ensures no API breakage was introduced between 1.24.x and 1.26.x that affects the codebase.

## Risks / Trade-offs

- **Go 1.26.1 may not be released yet** — as of the proposal date, verify the version exists before applying. If unavailable, use the latest stable (e.g., 1.24.2 or 1.25.x).
- **Toolchain mismatch on CI** — any CI pipeline must also use Go 1.26.1 or later after this change.
- **Minimal risk** — the codebase has no complex CGO, no assembly, and a single small dependency, making breakage unlikely.

## Migration Plan

1. Edit `go 1.24.1` → `go 1.26.1` in `go.mod`
2. Run `go mod tidy` to update `go.sum` if needed
3. Run `go build ./...` to verify compilation
4. Run `go vet ./...` to check for any new vet warnings

## Open Questions

- Is Go 1.26.1 the intended version, or should we target the latest stable release available at time of implementation?
