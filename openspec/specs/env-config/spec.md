# Spec: env-config

## Requirements

### Requirement: Configuration priority order
The binary SHALL apply configuration values in the following priority order (highest to lowest):

1. **CLI flags** — explicitly passed at invocation time
2. **Environment variables** — set in the shell or loaded from `.env`
3. **scenario.yaml overrides** — per-scenario defaults in the loaded scenario
4. **Hardcoded defaults** — compiled-in fallback values

#### Scenario: Env var beats scenario.yaml
- **WHEN** `SIM_TURNS=20` is set and the loaded `scenario.yaml` defines `turns: 5`
- **THEN** the binary SHALL use `20`

#### Scenario: CLI flag beats env var
- **WHEN** `SIM_TURNS=20` is set and `--turns=15` is also passed
- **THEN** the binary SHALL use `15`

#### Scenario: scenario.yaml beats hardcoded default
- **WHEN** no env var or CLI flag for `turns` is set and `scenario.yaml` defines `turns: 5`
- **THEN** the binary SHALL use `5`

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

### Requirement: .env excluded from version control
The `.gitignore` file SHALL contain an entry for `.env` so that local environment files are never committed.

#### Scenario: .env is ignored by git
- **WHEN** a `.env` file exists at the project root
- **THEN** `git status` SHALL NOT list it as a tracked or untracked file
