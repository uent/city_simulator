## Context

The Makefile defines Make variables (`MODEL`, `TURNS`, `SEED`, `OUTPUT`, `OLLAMA_URL`, `SCENARIO`) with hardcoded defaults, then passes them as explicit CLI flags to the binary. Since CLI flags win over env vars in `main.go`, the `.env` file is silently bypassed whenever `make run` is used.

The binary already supports env var configuration via `envOrString`/`envOrInt`/`envOrInt64` helpers. No Go changes are needed — only the Makefile changes.

## Goals / Non-Goals

**Goals:**
- `make run` with no arguments reads config from `.env` (or the binary's defaults)
- Users can still override individual values at the shell: `OLLAMA_MODEL=mistral make run`
- The Makefile stays useful as a task runner without being a second config layer

**Non-Goals:**
- Auto-loading `.env` inside the Makefile (avoids `include` complexity and shell quoting edge cases)
- Changing any Go source code
- Adding new Make targets

## Decisions

### Remove all flag passthrough from the `run` target

The `run` target becomes simply:

```makefile
run: build
    ./$(BINARY)
```

The binary reads env vars on its own. Users who want to override pass env vars inline:

```bash
OLLAMA_MODEL=mistral SIM_TURNS=20 make run
# or
make run OLLAMA_MODEL=mistral SIM_TURNS=20   # Make exports them into the subprocess env
```

**Alternative considered — `-include .env` in Makefile**: Would auto-load `.env` into Make's variable namespace, but requires mapping `OLLAMA_MODEL` → flag, handling quoting, and breaks when `.env` has comments or shell syntax. Rejected for fragility.

**Alternative considered — keep flags, but use `?=` with env var names**: e.g., `MODEL ?= $(OLLAMA_MODEL)`. Works but adds indirection with two variable names for the same thing. Rejected for complexity.

### Remove the now-unused Make variables

`MODEL`, `TURNS`, `SEED`, `OUTPUT`, `SCENARIO` are only used to build the flag string. Once flags are removed, these variables serve no purpose and should be deleted to avoid confusion.

`OLLAMA_URL` is the exception — it matches the env var name exactly and is shown in `make help`. It can be removed too since the binary reads `OLLAMA_URL` directly.

### Update `make help` output

The help target currently prints the current values of Make variables as a convenience. After the change, replace that section with a reference to `.env` / `.env.example` for configuration.

## Risks / Trade-offs

- **Discoverability** — Users who relied on `make help` to see current config values lose that view. Mitigation: update help to point to `.env.example`.
- **Muscle memory** — `make run MODEL=mistral` no longer works; must use `OLLAMA_MODEL=mistral make run`. Mitigation: document in help output and README.
- **Silent `.env` not found** — If user forgets to create `.env`, binary falls back to its hardcoded defaults (llama3, 10 turns, etc.). This was already the case when running the binary directly; no regression.
