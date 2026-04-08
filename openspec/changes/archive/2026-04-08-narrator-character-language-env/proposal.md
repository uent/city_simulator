## Why

The narrator (Game Director) and character prompts are currently hardcoded to a fixed language mix — Spanish labels in character prompts and English in director prompts — with no way to change this without modifying source code. Adding a `SIM_LANGUAGE` environment variable lets operators configure the simulation language at runtime, consistent with how all other runtime parameters are managed.

## What Changes

- Add `SIM_LANGUAGE` env variable and `--language` CLI flag following the existing `env-config` pattern.
- Pass the resolved language value through `SimConfig` to both the character prompt builder (`BuildSystemPrompt`) and the director prompt builder (`BuildDirectorPrompt`).
- Both builders inject a language instruction into the prompts so the LLM responds in the configured language.
- Update `.env.example` to document `SIM_LANGUAGE`.
- Update the `env-config` spec table with the new variable.

## Capabilities

### New Capabilities

- `simulation-language`: Configures the language used by the narrator and characters via a `SIM_LANGUAGE` env variable and `--language` CLI flag.

### Modified Capabilities

- `env-config`: The supported variable table gains a new entry for `SIM_LANGUAGE` / `--language`.

## Impact

- `cmd/simulator/main.go` — reads `SIM_LANGUAGE`, adds `--language` flag, threads value into `SimConfig`.
- `internal/scenario/scenario.go` — `SimConfig` and related merge/override structs gain a `Language` field.
- `internal/character/prompt.go` — `BuildSystemPrompt` and `BuildMovementPrompt` accept a language parameter and inject it.
- `internal/director/prompt.go` — `BuildDirectorPrompt` accepts a language parameter and injects it.
- `.env.example` — new documented variable.
- `openspec/specs/env-config/spec.md` — updated variable table.
- No external API or dependency changes.
