## 1. Config Layer

- [x] 1.1 Add `Language string` field to `SimConfig` and `CLIFlags` in `internal/scenario/scenario.go`
- [x] 1.2 Update `MergeConfig` in `internal/scenario/scenario.go` to apply `flags.Language` (language is not a scenario override)

## 2. Entry Point

- [x] 2.1 Add `envOrString("SIM_LANGUAGE", "")` read and `--language` flag in `cmd/simulator/main.go`
- [x] 2.2 Populate `cliFlags.Language` when `SIM_LANGUAGE` is set or `--language` is passed explicitly
- [x] 2.3 Pass `simCfg.Language` through to wherever prompt builders are called

## 3. Prompt Builders

- [x] 3.1 Add `language string` parameter to `BuildSystemPrompt` in `internal/character/prompt.go`; append `"Respond in <language>."` when non-empty
- [x] 3.2 Add `language string` parameter to `BuildMovementPrompt` in `internal/character/prompt.go`; append `"Respond in <language>."` when non-empty
- [x] 3.3 Add `language string` parameter to `BuildDirectorPrompt` in `internal/director/prompt.go`; append `"Respond in <language>."` when non-empty

## 4. Call Sites

- [x] 4.1 Update `internal/character/actor.go` to pass `language` to `BuildSystemPrompt` and `BuildMovementPrompt`
- [x] 4.2 Update `internal/simulation/engine.go` to pass `language` to `BuildDirectorPrompt`

## 5. Documentation

- [x] 5.1 Add `SIM_LANGUAGE` entry to `.env.example` with a comment describing its purpose and default
- [x] 5.2 Update `openspec/specs/env-config/spec.md` variable table to include `SIM_LANGUAGE` / `--language`
