## 1. Makefile cleanup

- [x] 1.1 Remove Make variable declarations: `MODEL`, `TURNS`, `SEED`, `OUTPUT`, `SCENARIO`, `OLLAMA_URL`, `BINARY`
- [x] 1.2 Simplify the `run` target to invoke `./city-simulator` with no flags
- [x] 1.3 Update `make help` to remove the "Configurable variables" section and add a line referencing `.env.example`
- [x] 1.4 Update `make clean` if it references removed variables (e.g., `$(OUTPUT)` or `$(BINARY)`)

## 2. Verification

- [x] 2.1 Run `make run` with `.env` sourced and confirm the binary picks up env var values (not hardcoded defaults)
- [x] 2.2 Run `OLLAMA_MODEL=mistral make run` and confirm the override is respected
- [x] 2.3 Run `make help` and confirm output no longer shows stale Make variable values
