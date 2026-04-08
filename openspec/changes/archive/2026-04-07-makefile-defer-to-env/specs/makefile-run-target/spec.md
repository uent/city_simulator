## MODIFIED Requirements

### Requirement: `make run` invokes the binary without explicit CLI flags
The `run` target SHALL invoke the binary with no configuration flags. All runtime configuration SHALL be sourced from environment variables (read by the binary itself) or the binary's hardcoded defaults. The Makefile SHALL NOT pass `-model`, `-ollama-url`, `-turns`, `-seed`, `-output`, or `-scenario` flags.

#### Scenario: Run with no arguments uses .env values
- **WHEN** the user has exported env vars from `.env` (e.g., `OLLAMA_MODEL=hermes3`) and runs `make run`
- **THEN** the binary starts using the value from the env var (`hermes3`), not a hardcoded Makefile default

#### Scenario: Run with inline env override
- **WHEN** the user runs `OLLAMA_MODEL=mistral make run`
- **THEN** the binary starts using `mistral` as the model

#### Scenario: Run with no .env falls back to binary defaults
- **WHEN** the user runs `make run` with no env vars set
- **THEN** the binary starts using its own hardcoded defaults (e.g., `llama3`, 10 turns)

### Requirement: Makefile removes unused Make variables
The Make variables `MODEL`, `TURNS`, `SEED`, `OUTPUT`, `SCENARIO` SHALL be removed from the Makefile. `OLLAMA_URL` SHALL also be removed. Their presence after the flag passthrough is removed would be misleading with no effect.

#### Scenario: No stale variables in Makefile
- **WHEN** a developer reads the Makefile
- **THEN** there are no variable declarations that have no effect on any target

### Requirement: `make help` documents env-based configuration
The `help` target SHALL direct users to `.env.example` for configuration instead of printing Make variable values.

#### Scenario: Help output references .env.example
- **WHEN** the user runs `make help`
- **THEN** the output mentions `.env` or `.env.example` as the configuration source
