# Specification: project-makefile

## Requirement: Makefile exists at project root
The project SHALL provide a `Makefile` at the repository root that exposes all common developer tasks as named targets.

### Scenario: Developer runs make without arguments
- **WHEN** a developer runs `make` in the project root
- **THEN** the `help` target SHALL execute and print all available targets with descriptions and example invocations

## Requirement: Build target compiles the simulator binary
The Makefile SHALL provide a `build` target that compiles `cmd/simulator/main.go` into a binary named `city-simulator` (or `city-simulator.exe` on Windows).

### Scenario: Successful build
- **WHEN** a developer runs `make build`
- **THEN** the Go binary SHALL be compiled to `./city-simulator` with no errors

## Requirement: Run target invokes the binary without explicit CLI flags
The Makefile SHALL provide a `run` target that builds and runs the simulator. The `run` target SHALL NOT pass any CLI flags to the binary — all runtime configuration SHALL be sourced from environment variables (loaded from `.env` by the binary) or the binary's hardcoded defaults.

### Scenario: Run with .env configuration
- **WHEN** a developer runs `make run` with a `.env` file present
- **THEN** the simulator SHALL start using values from `.env` (e.g., `OLLAMA_MODEL`, `SIM_TURNS`)

### Scenario: Run with inline env override
- **WHEN** a developer runs `OLLAMA_MODEL=mistral make run`
- **THEN** the simulator SHALL start with model `mistral`

### Scenario: Run with no .env falls back to binary defaults
- **WHEN** a developer runs `make run` with no env vars set and no `.env` file
- **THEN** the simulator SHALL start using its hardcoded defaults

## Requirement: Test target runs the Go test suite
The Makefile SHALL provide a `test` target that runs `go test ./...` with verbose output.

### Scenario: Tests pass
- **WHEN** a developer runs `make test`
- **THEN** all Go tests SHALL execute and results SHALL be printed to stdout

## Requirement: Fmt target formats Go source code
The Makefile SHALL provide a `fmt` target that runs `go fmt ./...` across the entire module.

### Scenario: Format applied
- **WHEN** a developer runs `make fmt`
- **THEN** all `.go` files SHALL be formatted in place using `gofmt` conventions

## Requirement: Vet target runs static analysis
The Makefile SHALL provide a `vet` target that runs `go vet ./...`.

### Scenario: Vet passes
- **WHEN** a developer runs `make vet`
- **THEN** the Go vet tool SHALL analyze all packages and report any issues to stdout

## Requirement: Clean target removes build artifacts
The Makefile SHALL provide a `clean` target that removes the compiled binary and the default JSONL output file.

### Scenario: Artifacts removed
- **WHEN** a developer runs `make clean`
- **THEN** `./city-simulator` (and `./city-simulator.exe`) and `simulation_output.jsonl` SHALL be deleted if they exist

## Requirement: Inline examples embedded in the Makefile
The Makefile SHALL include `## Example:` comment blocks above each target and a `help` target that prints usage examples via `@echo`, making the file self-documenting without any external docs.

### Scenario: Examples visible in source
- **WHEN** a developer opens the Makefile
- **THEN** each target SHALL have an `## Example:` comment showing a sample invocation

### Scenario: Help target prints examples
- **WHEN** a developer runs `make help`
- **THEN** all targets and example invocations SHALL be printed to stdout
