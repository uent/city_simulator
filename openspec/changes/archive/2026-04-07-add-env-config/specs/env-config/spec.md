## ADDED Requirements

### Requirement: Environment variable fallback for CLI flags
The binary SHALL read environment variables as default values for CLI flags. CLI flags SHALL take precedence over environment variables. The following variables SHALL be supported:

| Variable | Flag | Type | Default |
|---|---|---|---|
| `OLLAMA_URL` | `--ollama-url` | string | `http://localhost:11434` |
| `OLLAMA_MODEL` | `--model` | string | `llama3` |
| `SIM_TURNS` | `--turns` | int | `10` |
| `SIM_SEED` | `--seed` | int64 | `0` |
| `SIM_OUTPUT` | `--output` | string | `simulation_output.jsonl` |
| `SIM_CHARACTERS` | `--characters` | string | `configs/characters.yaml` |

#### Scenario: Env var sets default when flag not provided
- **WHEN** `OLLAMA_MODEL=mistral` is set in the environment and `--model` is not passed
- **THEN** the binary SHALL use `mistral` as the model name

#### Scenario: CLI flag overrides env var
- **WHEN** `OLLAMA_MODEL=mistral` is set and `--model=llama3` is also passed
- **THEN** the binary SHALL use `llama3`

#### Scenario: Invalid numeric env var falls back to hardcoded default
- **WHEN** `SIM_TURNS=abc` is set (non-numeric)
- **THEN** the binary SHALL log a warning and use the hardcoded default value (`10`)

### Requirement: .env.example file at project root
The repository SHALL contain a `.env.example` file documenting every supported environment variable with comments describing its purpose and the hardcoded default value.

#### Scenario: Developer can copy example to start
- **WHEN** a developer runs `cp .env.example .env`
- **THEN** the resulting `.env` file SHALL contain all supported variables with their default values, ready to customise

### Requirement: .env excluded from version control
The `.gitignore` file SHALL contain an entry for `.env` so that local environment files are never committed.

#### Scenario: .env is ignored by git
- **WHEN** a `.env` file exists at the project root
- **THEN** `git status` SHALL NOT list it as a tracked or untracked file
