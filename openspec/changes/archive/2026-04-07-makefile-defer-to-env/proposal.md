## Why

The Makefile always passes CLI flags explicitly (e.g., `-model llama3`), which overrides any env vars set via `.env` — making the `.env` file ineffective when using `make run`. This contradicts the documented contract: "CLI flags always take precedence over these variables."

## What Changes

- Remove explicit flag arguments from the `run` target in the Makefile
- The binary is invoked with no flags; it reads configuration exclusively from env vars (or its own hardcoded defaults)
- Overrides are still possible via `make run` by setting env vars inline: `make run OLLAMA_MODEL=mistral`

## Capabilities

### New Capabilities

<!-- None introduced -->

### Modified Capabilities

- `makefile-run-target`: The `run` target no longer passes `-model`, `-ollama-url`, `-turns`, `-seed`, `-output`, or `-scenario` as CLI flags. Configuration is delegated to env vars read by the binary.

## Impact

- `Makefile`: `run` target simplified — flags removed, Make variables (`MODEL`, `TURNS`, etc.) removed or repurposed as env var passthrough
- `.env` / `.env.example`: Now the authoritative source for local configuration when using `make run`
- No Go source changes needed — `main.go` already reads env vars via `envOrString`/`envOrInt`/`envOrInt64`
- `README.md` or help text may need updating to reflect new usage pattern
