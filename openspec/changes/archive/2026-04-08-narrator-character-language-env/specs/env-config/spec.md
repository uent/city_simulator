## MODIFIED Requirements

### Requirement: Environment variable fallback for CLI flags
The binary SHALL read environment variables as default values for CLI flags. The following variables SHALL be supported:

| Variable | Flag | Type | Default |
|---|---|---|---|
| `OLLAMA_URL` | `--ollama-url` | string | `http://localhost:11434` |
| `OLLAMA_MODEL` | `--model` | string | `llama3` |
| `SIM_TURNS` | `--turns` | int | `10` |
| `SIM_SEED` | `--seed` | int64 | `0` |
| `SIM_OUTPUT` | `--output` | string | `simulation_output.jsonl` |
| `SIM_SCENARIO` | `--scenario` | string | `default` |
| `SIM_LANGUAGE` | `--language` | string | `""` (empty — no instruction injected) |

#### Scenario: Env var sets default when flag not provided
- **WHEN** `OLLAMA_MODEL=mistral` is set in the environment and `--model` is not passed
- **THEN** the binary SHALL use `mistral` as the model name

#### Scenario: CLI flag overrides env var
- **WHEN** `OLLAMA_MODEL=mistral` is set and `--model=llama3` is also passed
- **THEN** the binary SHALL use `llama3`

#### Scenario: Invalid numeric env var falls back to hardcoded default
- **WHEN** `SIM_TURNS=abc` is set (non-numeric)
- **THEN** the binary SHALL log a warning and use the hardcoded default value (`10`)

#### Scenario: SIM_LANGUAGE env var sets language for prompts
- **WHEN** `SIM_LANGUAGE=Spanish` is set and `--language` is not passed
- **THEN** the binary SHALL use `Spanish` as the language value for all LLM prompts

### Requirement: .env.example file at project root
The repository SHALL contain a `.env.example` file documenting every supported environment variable with comments describing its purpose and the hardcoded default value. This includes `SIM_LANGUAGE`.

#### Scenario: Developer can copy example to start
- **WHEN** a developer runs `cp .env.example .env`
- **THEN** the resulting `.env` file SHALL contain all supported variables with their default values, ready to customise
