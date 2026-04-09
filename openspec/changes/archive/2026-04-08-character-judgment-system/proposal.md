## Why

Characters currently know other characters only by name and location — they have no prior opinion, no instinctive read, no subjective interpretation. Every conversation starts from zero, which flattens the richness of inter-character dynamics. A paranoid detective should already distrust the bookshop owner before they speak; a con artist should have already profiled every mark in the room.

## What Changes

- Add an `appearance` field to the Character schema — a single authored sentence describing how this person presents to the world, visible to others when forming impressions.
- Add a `CharacterJudgment` struct holding a character's subjective opinion of another: a first-person impression, trust level, interest level, and perceived threat level.
- Add a `Judgments map[string]CharacterJudgment` field to Character, keyed by character ID.
- Before the simulation starts, each character forms a judgment of every other character via a targeted LLM call — using only observable information (name, age, occupation, emotional state, appearance, location) filtered through the judging character's full psychological profile.
- When two characters converse, inject the relevant prior judgment into each character's system prompt.
- After every 10 conversations between the same pair, refresh both judgments using accumulated conversation history.
- When a new character is spawned by the director, immediately form judgments in both directions (new → all existing, all existing → new).
- Characters with a `CoverIdentity` are judged by their alias and role — never by their true identity.
- Add `appearance` to all existing scenario `characters.yaml` files.

## Capabilities

### New Capabilities

- `character-judgment`: Subjective per-character opinions formed before simulation start, injected into conversation prompts, and refreshed every 10 interactions between the same pair.

### Modified Capabilities

- `character-schema`: Adds `appearance` (string, optional) and `judgments` (runtime map, not persisted to YAML) to the Character struct.

## Impact

- `internal/character/character.go` — new fields and struct
- `internal/character/judgment.go` — new file: judgment formation, update, observable snapshot, prompt builders, format helpers
- `internal/simulation/engine.go` — pre-simulation judgment phase, per-pair conversation counter, post-conversation update trigger, judgment injection into system prompts
- `simulations/*/characters.yaml` — all existing scenarios gain `appearance` on each character
- No breaking changes to the simulation's external interface or output format; new prompt content is additive
