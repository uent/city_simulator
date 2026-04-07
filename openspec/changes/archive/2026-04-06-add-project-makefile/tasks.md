## 1. Makefile Creation

- [x] 1.1 Create `Makefile` at the project root with variables block: `BINARY`, `MODEL`, `OLLAMA_URL`, `CHARACTERS`, `TURNS`, `SEED`, `OUTPUT` — all with defaults matching app flags
- [x] 1.2 Add `.PHONY` declaration listing all targets: `help build run test fmt vet clean`

## 2. Targets Implementation

- [x] 2.1 Add `help` target that prints all available targets, their descriptions, and example invocations via `@echo`
- [x] 2.2 Add `build` target: `go build -o $(BINARY) ./cmd/simulator/` with `## Example: make build` comment
- [x] 2.3 Add `run` target: depends on `build`, runs the binary with all variable flags; include `## Example: make run MODEL=mistral TURNS=20` comment
- [x] 2.4 Add `test` target: `go test -v ./...` with `## Example: make test` comment
- [x] 2.5 Add `fmt` target: `go fmt ./...` with `## Example: make fmt` comment
- [x] 2.6 Add `vet` target: `go vet ./...` with `## Example: make vet` comment
- [x] 2.7 Add `clean` target: removes `$(BINARY)`, `$(BINARY).exe`, and `$(OUTPUT)` with `## Example: make clean` comment

## 3. Verification

- [x] 3.1 Run `make help` and confirm all targets and examples are printed correctly
- [x] 3.2 Run `make build` and confirm `./city-simulator` binary is produced
- [x] 3.3 Run `make fmt` and `make vet` and confirm they complete without errors
- [x] 3.4 Run `make clean` and confirm binary and output file are removed
