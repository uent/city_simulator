## Context

The simulation has two LLM-driven agents: characters (via `BuildSystemPrompt` / `BuildMovementPrompt` in `internal/character/prompt.go`) and the Game Director (via `BuildDirectorPrompt` in `internal/director/prompt.go`). Both build system prompts that include instructional text. Currently the character prompts use Spanish labels while the director prompt uses English — neither is configurable without touching source code.

All other runtime parameters (`OLLAMA_MODEL`, `SIM_TURNS`, etc.) are read from `.env` / CLI flags and flow through `SimConfig`. Language follows the same path.

## Goals / Non-Goals

**Goals:**
- Add `SIM_LANGUAGE` env variable and `--language` CLI flag with the same priority semantics as existing flags (CLI > env > scenario > default).
- Thread the resolved language value through `SimConfig` into both prompt builders.
- Each prompt builder appends a single language instruction line to the LLM prompt.
- Document the variable in `.env.example` and update the `env-config` spec table.

**Non-Goals:**
- Translating hardcoded field labels in the prompts (e.g. "Motivación:", "Miedo:") — these are field identifiers, not prose the LLM must replicate.
- Validating that the language value is a known BCP-47 tag — free-form strings like "Spanish", "English", or "es" are all accepted.
- Per-character language overrides.
- Translating scenario YAML content (world descriptions, character backstories, etc.).

## Decisions

### Language is a plain string, not an enum
**Decision**: Accept any free-form string (e.g. `"Spanish"`, `"English"`, `"Español"`, `"fr"`).
**Rationale**: The value is forwarded verbatim to the LLM prompt (`"Respond in <language>."`), so correctness is the LLM's concern, not ours. An enum would need maintenance for every language and provides little value.
**Alternatives considered**: BCP-47 validation — rejected as unnecessary coupling.

### Inject language as a trailing instruction, not a prefix
**Decision**: Append `"Respond in <language>."` as the final line of both the character system prompt and the director prompt.
**Rationale**: LLMs follow trailing instructions reliably. Prefixing risks being overridden by later prompt content.

### Language is not overridable per-scenario
**Decision**: `RuntimeOverrides` does not gain a `Language` field.
**Rationale**: The proposal spec states language is an operator-level concern (`.env`/CLI), not a per-scenario concern, consistent with the existing treatment of `OLLAMA_MODEL`.

### Prompt builder signature change
**Decision**: `BuildSystemPrompt`, `BuildMovementPrompt`, and `BuildDirectorPrompt` each gain a `language string` parameter.
**Rationale**: Pure functions are easier to test than global state. Callers already have `SimConfig` available.
**Alternatives considered**: Package-level global — rejected, harder to test and thread-unsafe.

## Risks / Trade-offs

- [Prompt builder signature change] → All call sites must be updated; the compiler enforces this so no runtime risk.
- [Free-form language value] → A typo in `.env` (e.g. `Spansih`) produces unexpected LLM output with no error. Mitigation: `.env.example` documents common values clearly.

## Migration Plan

1. Add `Language` field to `SimConfig` and `CLIFlags` in `internal/scenario/scenario.go`. Update `MergeConfig` (language is not in `RuntimeOverrides`).
2. Update `cmd/simulator/main.go` to read `SIM_LANGUAGE`, add `--language` flag, pass through to `SimConfig`.
3. Update `BuildSystemPrompt` and `BuildMovementPrompt` in `internal/character/prompt.go` — add `language string` param, append instruction when non-empty.
4. Update `BuildDirectorPrompt` in `internal/director/prompt.go` — same pattern.
5. Update all call sites in `internal/character/actor.go` and `internal/simulation/engine.go`.
6. Update `.env.example` and `openspec/specs/env-config/spec.md`.

Rollback: revert the commit. No data migrations or external system changes.
