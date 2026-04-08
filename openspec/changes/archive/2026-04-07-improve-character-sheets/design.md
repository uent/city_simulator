## Context

The city simulator builds a system prompt per character via `BuildSystemPrompt` in `internal/llm/prompt.go`. Currently it joins `Personality []string` as a comma list and dumps a prose `Backstory`. Local LLMs (7B–13B) average across trait lists and produce generic behavior because they lack psychological anchors: they know a character is "sharp, guarded, dry humor" but have nothing to resolve ambiguity when that character faces a situation not covered by the backstory.

The new schema replaces flat trait lists with structured psychological anchors — motivation, fear, core belief, internal tension — plus a voice profile and relational defaults. This gives the model generative rules rather than descriptions to parrot.

Two concrete files need changes: the `Character` struct and the `BuildSystemPrompt` function. All scenario YAML files are internal and fully owned, so there is no migration concern for third-party consumers.

## Goals / Non-Goals

**Goals:**
- Replace `Personality []string` and `Backstory string` with structured psychological fields
- Add `Voice`, `RelationalDefaults`, `FormativeEvents`, `InternalTension`, and `DialogueExamples` to `Character`
- Rewrite `BuildSystemPrompt` to emit the new structured template
- Rewrite both `simulations/default/characters.yaml` and `simulations/honey-heist/characters.yaml`
- Add `simulations/CHARACTER_RULES.md` — an LLM-readable rulebook for creating characters in this schema
- Update the honey-heist spec to reflect the new required fields

**Non-Goals:**
- Backwards compatibility with old YAML format — all scenario files are internal
- Lorebook / retrieval for deep historical context (mentioned as future path; not in scope)
- Runtime mutation of psychological core fields (those are static; only `EmotionalState` is mutable)
- Changes to the simulation engine, conversation manager, or scheduler

## Decisions

### 1. Hard replacement of `Personality` and `Backstory` — no deprecation shim

**Decision**: Remove both fields entirely from the `Character` struct and YAML.

**Rationale**: Keeping deprecated fields alongside new ones creates ambiguity about which the LLM prompt should use. Since all scenario files are internal, we can rewrite them atomically with the struct change. There are no external consumers.

**Alternative considered**: Mark fields as `yaml:",omitempty"` and keep them as optional. Rejected because it leaves two competing character models in the codebase with no clear winner.

---

### 2. Nested structs for `Voice` and `RelationalDefaults`

**Decision**: Introduce two sub-structs that map to YAML nested objects:

```go
type VoiceProfile struct {
    Formality         string `yaml:"formality"`
    VerbalTics        string `yaml:"verbal_tics"`
    ResponseLength    string `yaml:"response_length"`
    HumorType         string `yaml:"humor_type"`
    CommunicationStyle string `yaml:"communication_style"`
}

type RelationalDefaults struct {
    Strangers  string `yaml:"strangers"`
    Authority  string `yaml:"authority"`
    Vulnerable string `yaml:"vulnerable"`
}
```

**Rationale**: Nested structs make the YAML readable and self-documenting. Flat fields with underscores (`voice_formality`) would bloat the struct and obscure grouping.

**Alternative considered**: Single prose `Voice string` field. Rejected because prose is harder for the LLM to parse consistently and harder for scenario authors to fill in correctly.

---

### 3. Keep `Goals []string` alongside new `Motivation string`

**Decision**: Retain `Goals` as a distinct field. `Motivation` answers *why* the character acts; `Goals` lists *what* they concretely pursue. Both are surfaced in the system prompt.

**Rationale**: A character's motivation ("wants to believe the world is fair") and their goals ("find the person who killed her partner") are complementary, not redundant. The model uses them differently — motivation resolves value conflicts, goals drive plot-level behavior.

---

### 4. `BuildSystemPrompt` renders the new structured template verbatim

**Decision**: Rewrite the function to emit sections matching the user-designed template:
```
Motivación / Miedo / Creencia central / Tensión interna /
Eventos formativos / Voz / Relaciones default / Ejemplos de diálogo
```
Sections with empty values are omitted silently.

**Rationale**: The prompt structure is the interface between the character data and the LLM. Keeping it close to the design template makes it easy to reason about.

---

### 5. `CHARACTER_RULES.md` as a prose rulebook, not a schema file

**Decision**: Add `simulations/CHARACTER_RULES.md` as human- and LLM-readable Markdown containing: the template, field-by-field guidance, anti-patterns, and a worked example.

**Rationale**: The rules need to be accessible to another LLM creating characters (e.g., via a prompt that says "read CHARACTER_RULES.md and create a character"). A YAML schema is machine-checkable but not instructive. A prose rulebook is.

## Risks / Trade-offs

- **YAML authoring burden increases** → All new fields are optional strings; missing fields produce empty prompt sections (silently skipped). Scenario authors can migrate incrementally.
- **`FormativeEvents` may be too short for complex characters** → Capped at 3 bullets by convention, not enforcement. Authors who need more depth should use the lorebook pattern (future work).
- **System prompt length increases** → The new structured template is longer than the old flat format. On 4K context models this is a real constraint. Mitigation: `DialogueExamples` can be trimmed to 2 entries for smaller models; the template keeps each section terse.
- **`EmotionalState` vs `InternalTension` conceptual overlap** → `EmotionalState` is runtime/mutable ("weary but determined at start of scene"); `InternalTension` is static psychology ("values honesty but lies to protect loved ones"). These are documented clearly in `CHARACTER_RULES.md`.
