## ADDED Requirements

### Requirement: LLM provider selection variable
The binario SHALL read `LLM_PROVIDER` from the environment to select the LLM backend. Accepted values are `ollama` (default) and `openrouter`. Any other value SHALL cause the binary to exit with a descriptive error before the simulation starts.

#### Scenario: Default provider is Ollama
- **WHEN** `LLM_PROVIDER` is not set
- **THEN** the binary SHALL initialize the Ollama provider using `OLLAMA_URL` and `OLLAMA_MODEL`

#### Scenario: OpenRouter selected via env var
- **WHEN** `LLM_PROVIDER=openrouter` is set
- **THEN** the binary SHALL initialize the OpenRouter provider using `OPENROUTER_API_KEY` and `OPENROUTER_MODEL`

#### Scenario: Unknown provider value causes startup error
- **WHEN** `LLM_PROVIDER=bedrock` is set (unsupported value)
- **THEN** the binary SHALL exit with a non-zero status and an error message naming the unsupported provider

---

### Requirement: OpenRouter environment variables
The binary SHALL read the following environment variables when `LLM_PROVIDER=openrouter`:

| Variable | Purpose | Default |
|---|---|---|
| `OPENROUTER_API_KEY` | Bearer token for OpenRouter API authentication | (required — no default) |
| `OPENROUTER_MODEL` | Model identifier to use (e.g. `openai/gpt-4o-mini`) | `openai/gpt-4o-mini` |
| `OPENROUTER_BASE_URL` | Base URL override for the OpenRouter API | `https://openrouter.ai/api/v1` |

If `LLM_PROVIDER=openrouter` and `OPENROUTER_API_KEY` is empty, the binary SHALL exit with a descriptive error before making any network calls.

#### Scenario: API key missing causes startup error
- **WHEN** `LLM_PROVIDER=openrouter` is set and `OPENROUTER_API_KEY` is empty
- **THEN** the binary SHALL exit with a non-zero status and an error explaining that `OPENROUTER_API_KEY` is required

#### Scenario: Default model used when OPENROUTER_MODEL not set
- **WHEN** `LLM_PROVIDER=openrouter` is set, `OPENROUTER_API_KEY` is set, and `OPENROUTER_MODEL` is not set
- **THEN** the binary SHALL use `openai/gpt-4o-mini` as the model

#### Scenario: Custom base URL override
- **WHEN** `OPENROUTER_BASE_URL=https://proxy.example.com/v1` is set
- **THEN** all OpenRouter API calls SHALL be made to that base URL instead of the default

---

### Requirement: Updated .env.example with provider variables
The `.env.example` file SHALL document all new provider variables alongside the existing ones, with comments explaining their purpose and valid values.

#### Scenario: Developer copies example and sees provider options
- **WHEN** a developer opens `.env.example`
- **THEN** they SHALL see `LLM_PROVIDER`, `OPENROUTER_API_KEY`, `OPENROUTER_MODEL`, and `OPENROUTER_BASE_URL` with descriptive comments

## MODIFIED Requirements

### Requirement: Variables de entorno soportadas
El binario SHALL leer las siguientes variables de entorno como valores default para sus flags:

| Variable | Flag | Tipo | Default |
|---|---|---|---|
| `LLM_PROVIDER` | `--provider` | string | `ollama` |
| `OLLAMA_URL` | `--ollama-url` | string | `http://localhost:11434` |
| `OLLAMA_MODEL` | `--model` | string | `llama3` |
| `OPENROUTER_API_KEY` | (no flag — secrets not exposed as flags) | string | `""` |
| `OPENROUTER_MODEL` | `--model` (shared) | string | `openai/gpt-4o-mini` |
| `OPENROUTER_BASE_URL` | `--openrouter-base-url` | string | `https://openrouter.ai/api/v1` |
| `SIM_TURNS` | `--turns` | int | `10` |
| `SIM_SEED` | `--seed` | int64 | `0` |
| `SIM_OUTPUT` | `--output` | string | `simulation_output.jsonl` |
| `SIM_SCENARIO` | `--scenario` | string | `default` |
| `SIM_LANGUAGE` | `--language` | string | `""` (vacío — sin instrucción de idioma) |

#### Scenario: Env var setea el default cuando el flag no se provee
- **WHEN** `OLLAMA_MODEL=mistral` está en el entorno y `--model` no se pasa
- **THEN** el binario SHALL usar `mistral` como nombre del modelo

#### Scenario: CLI flag sobreescribe env var
- **WHEN** `OLLAMA_MODEL=mistral` está seteada y se pasa también `--model=llama3`
- **THEN** el binario SHALL usar `llama3`

#### Scenario: Env var numérica inválida usa el default hardcodeado
- **WHEN** `SIM_TURNS=abc` está seteada (no-numérica)
- **THEN** el binario SHALL loguear un warning y usar el valor default hardcodeado (`10`)

#### Scenario: SIM_LANGUAGE setea el idioma para los prompts
- **WHEN** `SIM_LANGUAGE=Spanish` está seteada y `--language` no se pasa
- **THEN** el binario SHALL usar `"Spanish"` como valor de idioma para todos los prompts LLM

#### Scenario: LLM_PROVIDER flag sobreescribe env var
- **WHEN** `LLM_PROVIDER=ollama` está seteada y se pasa `--provider=openrouter`
- **THEN** el binario SHALL usar `openrouter`
