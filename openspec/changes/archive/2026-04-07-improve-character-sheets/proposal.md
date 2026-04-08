## Why

The current character schema uses flat trait lists and biography-style backstories that cause LLMs to produce averaged, generic behavior — particularly problematic for smaller local models (7B–13B). Characters need psychological anchors that let the model derive consistent behavior in novel situations without per-scene scripting.

## What Changes

- **Replace** `Personality []string` and `Backstory string` fields on `Character` with a structured psychological core: `motivation`, `fear`, `core_belief`, `internal_tension`, `formative_events`
- **Add** a `Voice` sub-struct capturing linguistic patterns: formality, verbal tics, response length tendency, humor type, and communication style
- **Add** a `RelationalDefaults` sub-struct defining default stance toward strangers, authority figures, and vulnerable people
- **Add** `DialogueExamples []string` — 3–4 representative lines that anchor the model's voice more concretely than prose descriptions
- **Update** `simulations/default/characters.yaml` and `simulations/honey-heist/characters.yaml` to the new schema
- **Update** `LoadCharacters` to handle validation and defaults for the new fields
- **Add** generation rules document at `simulations/CHARACTER_RULES.md` so another LLM can create new characters following the schema

## Capabilities

### New Capabilities

- `character-schema`: Defines the new structured character profile format — fields, YAML layout, Go struct, validation rules, and the generative rules an LLM must follow to create coherent characters

### Modified Capabilities

- `honey-heist-scenario`: Character roster requirements change from `name`/`role`/`personality` string to the new structured schema (motivation, fear, core_belief, formative_events, voice, relational_defaults, dialogue_examples)

## Impact

- `internal/character/character.go`: `Character` struct gains new fields; old `Personality` and `Backstory` fields are removed
- `simulations/default/characters.yaml`: Rewritten with new schema (Elena, Marcus, Nadia)
- `simulations/honey-heist/characters.yaml`: Rewritten with new schema (all 6 bears)
- `simulations/CHARACTER_RULES.md`: New file — generation rules for future LLM character creation
- No changes to the loader interface (`Load`, `LoadCharacters` signatures unchanged)
- No changes to world, scenario, or simulation engine packages
