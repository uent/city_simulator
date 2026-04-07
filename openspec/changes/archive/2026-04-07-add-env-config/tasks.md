## 1. Version Control

- [x] 1.1 Create or update `.gitignore` at project root with a `.env` entry

## 2. Environment Example File

- [x] 2.1 Create `.env.example` at project root with all six variables, default values, and inline comments

## 3. Binary: env-var helpers

- [x] 3.1 Add `envOrString(key, fallback string) string` helper in `cmd/simulator/main.go`
- [x] 3.2 Add `envOrInt(key string, fallback int) int` helper with warning log on parse failure
- [x] 3.3 Add `envOrInt64(key string, fallback int64) int64` helper with warning log on parse failure

## 4. Binary: wire env vars into flag defaults

- [x] 4.1 Update `--characters` flag default to use `envOrString("SIM_CHARACTERS", "configs/characters.yaml")`
- [x] 4.2 Update `--model` flag default to use `envOrString("OLLAMA_MODEL", "llama3")`
- [x] 4.3 Update `--ollama-url` flag default to use `envOrString("OLLAMA_URL", "http://localhost:11434")`
- [x] 4.4 Update `--turns` flag default to use `envOrInt("SIM_TURNS", 10)`
- [x] 4.5 Update `--seed` flag default to use `envOrInt64("SIM_SEED", 0)`
- [x] 4.6 Update `--output` flag default to use `envOrString("SIM_OUTPUT", "simulation_output.jsonl")`
