## Context

The project is a Go CLI application (`cmd/simulator/main.go`) that runs an LLM-powered city simulation via Ollama. Currently there is no Makefile, so developers must run raw `go build`, `go run`, and `go test` commands manually with multiple flags. The Makefile will be the single entry point for all dev tasks.

## Goals / Non-Goals

**Goals:**
- Provide `make` targets for: build, run, test, lint, format, clean, and help
- Embed inline examples inside the Makefile (via `@echo` or comment blocks) so the file itself is self-documenting
- Use sensible defaults that match the app's existing flag defaults (`llama3`, `http://localhost:11434`, `configs/characters.yaml`, `10` turns)
- Keep the Makefile compatible with both Linux/macOS and Git Bash on Windows

**Non-Goals:**
- Docker / container orchestration targets
- CI/CD pipeline integration
- Cross-compilation or release packaging

## Decisions

### Use GNU Make with `.PHONY` declarations
All targets are phony (no file outputs match target names). Declaring them `.PHONY` avoids false-positive "up to date" skips and is standard Go project practice.

**Alternative considered**: Task runner (Taskfile.yml) — rejected to avoid adding a new tool dependency when `make` is universally available.

### Inline examples via `@echo` in a `help` target
A `help` target prints available commands and example invocations at runtime. Additionally, `## Example:` comment blocks above each target serve as in-file documentation visible to anyone reading the Makefile.

**Alternative considered**: A separate `USAGE.md` — rejected because it drifts out of sync. Keeping examples in the Makefile ensures they are always co-located with the commands.

### Variables for all tuneable flags
All CLI flags are exposed as Makefile variables with defaults (e.g., `MODEL ?= llama3`), allowing one-line overrides: `make run MODEL=mistral TURNS=20`.

## Risks / Trade-offs

- [Windows compatibility] `make` may not be installed on Windows outside of Git Bash/WSL → Mitigation: document requirement in help text; the project already targets Unix-style tooling.
- [Default model drift] If the app changes its default model, the Makefile variable must be updated manually → Mitigation: note the coupling in a comment.

## Open Questions

- Should a `docker-run` target be added later for Ollama? Left out of scope for this change.
