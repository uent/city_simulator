# ==============================================================================
# City Simulator — Makefile
# ==============================================================================
# Variables (override any of these on the command line)
#   make run MODEL=mistral TURNS=20 SEED=42
# ==============================================================================

BINARY      ?= city-simulator
MODEL       ?= llama3
OLLAMA_URL  ?= http://localhost:11434
CHARACTERS  ?= configs/characters.yaml
TURNS       ?= 10
SEED        ?= 0
OUTPUT      ?= simulation_output.jsonl

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
	@echo "  make help                          Show this help message"
	@echo "  make build                         Compile the simulator binary"
	@echo "  make run                           Build and run the simulation (default flags)"
	@echo "  make run MODEL=mistral TURNS=20    Run with a custom model and turn count"
	@echo "  make run SEED=42 OUTPUT=out.jsonl  Run with a fixed seed and custom output file"
	@echo "  make test                          Run the full test suite"
	@echo "  make fmt                           Format all Go source files"
	@echo "  make vet                           Run static analysis on all packages"
	@echo "  make clean                         Remove the compiled binary and output file"
	@echo ""
	@echo "Configurable variables (current values):"
	@echo "  BINARY     = $(BINARY)"
	@echo "  MODEL      = $(MODEL)"
	@echo "  OLLAMA_URL = $(OLLAMA_URL)"
	@echo "  CHARACTERS = $(CHARACTERS)"
	@echo "  TURNS      = $(TURNS)"
	@echo "  SEED       = $(SEED)"
	@echo "  OUTPUT     = $(OUTPUT)"
	@echo ""

# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

## Example: make build
build:
	go build -o $(BINARY) ./cmd/simulator/

# ------------------------------------------------------------------------------
# Run
# ------------------------------------------------------------------------------

## Example: make run
## Example: make run MODEL=mistral TURNS=20
## Example: make run SEED=42 OLLAMA_URL=http://localhost:11434 OUTPUT=run.jsonl
run: build
	./$(BINARY) \
		-characters $(CHARACTERS) \
		-model      $(MODEL) \
		-ollama-url $(OLLAMA_URL) \
		-turns      $(TURNS) \
		-seed       $(SEED) \
		-output     $(OUTPUT)

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
	rm -f $(BINARY) $(BINARY).exe $(OUTPUT)
