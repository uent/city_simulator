## Why

The simulation worlds lack a layer that communicates the *fundamental nature* of the world — its premise, its rules, and the underlying tension that characters must navigate. In Honey Heist, for example, the premise is that bears are disguised as humans in a human world; without that layer, characters can't reason about their cover identities, exposure risk, or what it means to "slip up." This layer is essential for any scenario where characters are operating under a hidden identity or against a social constraint they must not violate.

## What Changes

- Add a `WorldConcept` struct to `world.yaml` that captures the world's premise, its rules (what is normal vs. what would expose a character), and a short `flavor` description for narrative tone.
- Add a `cover_identity` block to individual character entries in `characters.yaml` that describes how a specific character presents themselves within the world (their alias, role, and any known weaknesses in their cover).
- Expose `WorldConcept` through `PublicSummary()` so every character receives the world's rules in their system prompt context.
- Expose per-character `CoverIdentity` in character prompts so each character reasons about maintaining their cover when they act.
- Update the honey-heist scenario YAML files to demonstrate the feature with a fully-authored world concept and cover identities for all characters.

## Capabilities

### New Capabilities

- `world-concept`: A `WorldConcept` struct in `world.yaml` — `premise`, `rules` (list of constraints that define normalcy), and `flavor` (tone/mood string). Surfaced in `PublicSummary()` as a "World Rules" section.
- `character-cover-identity`: A `CoverIdentity` struct on `Character` — `alias`, `role`, `backstory`, and `weaknesses` (list). Loaded from `characters.yaml` and injected into character system prompts.

### Modified Capabilities

- `world`: `WorldConfig` gains a `Concept WorldConcept` field; `PublicSummary()` gains a "World Rules" block when `Concept` is non-empty.
- `character-schema`: `Character` struct gains a `CoverIdentity *CoverIdentity` field loaded from YAML; `BuildSystemPrompt` includes cover identity block when present.

## Impact

- `internal/world/state.go` — add `WorldConcept` struct and field to `WorldConfig`; update `PublicSummary()`.
- `internal/character/character.go` — add `CoverIdentity` struct and field to `Character`.
- `internal/character/prompt.go` — inject cover identity section into system prompt builder.
- `simulations/honey-heist/world.yaml` — add `concept:` block.
- `simulations/honey-heist/characters.yaml` — add `cover_identity:` to each character.
- No breaking changes to existing scenarios that omit these fields (both are optional).
