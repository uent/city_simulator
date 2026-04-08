# ==============================================================================
# City Simulator — Makefile
# ==============================================================================
# Configuration is read from .env (via env vars). Copy .env.example to get started.
# Override individual vars inline: OLLAMA_MODEL=mistral make run
# ==============================================================================

.PHONY: help build run test fmt vet clean

# ------------------------------------------------------------------------------
# Default target
# ------------------------------------------------------------------------------

## Example: make
## Example: make help
help:
	@echo ""
	@echo "City Simulator — available targets:"
	@echo ""
	@echo "  make help    Show this help message"
	@echo "  make build   Compile the simulator binary"
	@echo "  make run     Build and run the simulation"
	@echo "  make test    Run the full test suite"
	@echo "  make fmt     Format all Go source files"
	@echo "  make vet     Run static analysis on all packages"
	@echo "  make clean   Remove the compiled binary and output file"
	@echo ""
	@echo "Configuration: copy .env.example to .env and set values there."
	@echo "Overrides:     OLLAMA_MODEL=mistral make run"
	@echo ""

# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

## Example: make build
build:
	go build -o city-simulator ./cmd/simulator/

# ------------------------------------------------------------------------------
# Run
# ------------------------------------------------------------------------------

## Example: make run
## Example: OLLAMA_MODEL=mistral make run
run: build
	./city-simulator

# ------------------------------------------------------------------------------
# Test
# ------------------------------------------------------------------------------

## Example: make test
test:
	go test -v ./...

# ------------------------------------------------------------------------------
# Format
# ------------------------------------------------------------------------------

## Example: make fmt
fmt:
	go fmt ./...

# ------------------------------------------------------------------------------
# Vet (static analysis)
# ------------------------------------------------------------------------------

## Example: make vet
vet:
	go vet ./...

# ------------------------------------------------------------------------------
# Clean
# ------------------------------------------------------------------------------

## Example: make clean
clean:
	rm -f city-simulator city-simulator.exe simulation_output.jsonl
