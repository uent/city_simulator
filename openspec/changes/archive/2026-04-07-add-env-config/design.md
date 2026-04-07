## Context

The simulator binary configures itself entirely via `flag` defaults in `cmd/simulator/main.go`. Developers running the simulator locally always need to remember and re-type flag values. The Go `flag` package doesn't natively support env-var fallbacks, so we need a thin helper that reads env vars and feeds them as default values before `flag.Parse()` runs.

## Goals / Non-Goals

**Goals:**
- Read six env vars (`OLLAMA_URL`, `OLLAMA_MODEL`, `SIM_TURNS`, `SIM_SEED`, `SIM_OUTPUT`, `SIM_CHARACTERS`) and use them as flag defaults.
- Ship `.env.example` documenting every variable with comments.
- Keep `.env` out of version control via `.gitignore`.

**Non-Goals:**
- Shipping a dotenv loader library (`.env` files are for the developer's shell to source, not auto-loaded by the binary).
- Validating or transforming env var values beyond what `strconv` already does for numeric types.
- Supporting nested/complex config (YAML config file is already `configs/characters.yaml`).

## Decisions

**Env vars as flag defaults, not overrides**
CLI flags always win. We call `flag.StringVar` / `flag.IntVar` / `flag.Int64Var` with `os.Getenv(...)` (or a parsed numeric equivalent) as the default value. This keeps the ergonomics of `--help` accurate and the precedence model simple.

Alternative considered: read env vars after `flag.Parse()` and override zero-value flags. Rejected because it would make `--help` show incorrect defaults and silently ignore explicit `--flag=""` invocations.

**No third-party dotenv library**
`godotenv` and similar packages auto-load `.env` at startup, which is magic and can mask misconfiguration in CI. Instead, `.env.example` is provided for developers to source manually (`source .env` or `export $(cat .env | xargs)`). The Makefile can provide a `make run` target that sources it.

**Helper function in `main.go`**
A small `envOr*` set of helpers (e.g., `envOrString`, `envOrInt`, `envOrInt64`) keeps the flag registration readable without adding a new package.

## Risks / Trade-offs

- `SIM_SEED=abc` will cause the binary to fall back to the zero value silently → Mitigation: log a warning on parse failure.
- `.env` accidentally committed → Mitigation: `.gitignore` entry; `.env.example` is the committed template.

## Migration Plan

1. Add `.env` to `.gitignore` (create the file if absent).
2. Add `.env.example` to project root.
3. Add `envOr*` helpers and update `flag.*` defaults in `cmd/simulator/main.go`.
4. No database migrations or binary protocol changes; fully backward compatible.
