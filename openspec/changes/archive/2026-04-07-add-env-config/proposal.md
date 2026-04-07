## Why

Currently all runtime parameters (Ollama URL, model name, turns, seed, output path, characters file) must be passed as CLI flags on every run, which is error-prone and inconvenient for local development. A `.env` file lets developers persist their local defaults without modifying source code or repeating long flag strings.

## What Changes

- Add `.env.example` documenting every supported environment variable with sensible defaults and comments.
- Add `.env` to `.gitignore` so local overrides are never committed.
- The binary reads environment variables as fallback defaults before applying CLI flags (flags still take precedence).

## Capabilities

### New Capabilities

- `env-config`: Load runtime configuration from environment variables (`OLLAMA_URL`, `OLLAMA_MODEL`, `SIM_TURNS`, `SIM_SEED`, `SIM_OUTPUT`, `SIM_CHARACTERS`) with CLI flags retaining highest priority.

### Modified Capabilities

<!-- No existing spec-level requirements are changing. -->

## Impact

- `cmd/simulator/main.go`: replace hardcoded flag defaults with values read from env vars via `os.Getenv`.
- New file: `.env.example` at project root.
- `.gitignore`: add `.env` entry (create file if it doesn't exist).
- No new dependencies required (`os` stdlib is sufficient).
