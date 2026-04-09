# Character Schema

## Requirements

### Requirement: Psychological core fields on Character

The `Character` struct SHALL expose four string fields that form its psychological core: `Motivation` (what the character wants and why), `Fear` (what they avoid at all costs), `CoreBelief` (their foundational view of how the world works), and `InternalTension` (a single contradiction that defines their complexity). The YAML keys SHALL be `motivation`, `fear`, `core_belief`, and `internal_tension`. The struct SHALL also expose a `CoverIdentity *CoverIdentity` field (YAML key `cover_identity`) as specified in the character-cover-identity capability; a nil pointer means the character has no cover identity.

#### Scenario: All core fields parsed from YAML
- **WHEN** a `characters.yaml` entry includes `motivation`, `fear`, `core_belief`, and `internal_tension` keys
- **THEN** `LoadCharacters` SHALL populate the corresponding struct fields with the verbatim string values and return no error

#### Scenario: Missing psychological core fields default to empty string
- **WHEN** a `characters.yaml` entry omits one or more of `motivation`, `fear`, `core_belief`, `internal_tension`
- **THEN** `LoadCharacters` SHALL still succeed, leaving the missing fields as empty strings

#### Scenario: CoverIdentity nil when cover_identity omitted
- **WHEN** a `characters.yaml` entry omits the `cover_identity` key
- **THEN** `Character.CoverIdentity` SHALL be nil after loading and no error returned

#### Scenario: CoverIdentity populated when cover_identity present
- **WHEN** a `characters.yaml` entry contains a `cover_identity:` block with at minimum an `alias` field
- **THEN** `Character.CoverIdentity` SHALL be non-nil after loading

---

### Requirement: Formative events field on Character

The `Character` struct SHALL expose `FormativeEvents []string` (YAML key `formative_events`). Each entry SHALL be a single causal sentence in the format "event → consequence" that explains why the character behaves a certain way. The field SHOULD contain 2–3 entries; it is not enforced at load time.

#### Scenario: Formative events parsed from YAML list
- **WHEN** `characters.yaml` contains a `formative_events` list with 3 entries
- **THEN** `Character.FormativeEvents` SHALL have length 3 with entries matching the YAML values in order

#### Scenario: Formative events absent defaults to nil slice
- **WHEN** `characters.yaml` omits `formative_events`
- **THEN** `Character.FormativeEvents` SHALL be nil (or empty slice) and no error returned

---

### Requirement: Voice profile sub-struct on Character

The `Character` struct SHALL expose a `Voice VoiceProfile` field (YAML key `voice`). The `VoiceProfile` struct SHALL contain: `Formality string` (`formality`), `VerbalTics string` (`verbal_tics`), `ResponseLength string` (`response_length`), `HumorType string` (`humor_type`), and `CommunicationStyle string` (`communication_style`).

#### Scenario: Voice sub-struct parsed from nested YAML
- **WHEN** `characters.yaml` contains a nested `voice:` block with all five sub-keys
- **THEN** `Character.Voice` SHALL be populated with all five field values

#### Scenario: Partial voice block is accepted
- **WHEN** `characters.yaml` contains a `voice:` block with only `formality` set
- **THEN** `Character.Voice.Formality` SHALL be set and all other Voice fields SHALL be empty strings, with no error

---

### Requirement: Relational defaults sub-struct on Character

The `Character` struct SHALL expose a `RelationalDefaults RelationalProfile` field (YAML key `relational_defaults`). The `RelationalProfile` struct SHALL contain: `Strangers string` (`strangers`), `Authority string` (`authority`), and `Vulnerable string` (`vulnerable`). Each describes the character's default behavioral stance toward that category of person.

#### Scenario: Relational defaults parsed from nested YAML
- **WHEN** `characters.yaml` contains a `relational_defaults:` block with `strangers`, `authority`, and `vulnerable` keys
- **THEN** `Character.RelationalDefaults` SHALL be populated with all three values

#### Scenario: Relational defaults absent defaults to zero-value struct
- **WHEN** `characters.yaml` omits `relational_defaults`
- **THEN** `Character.RelationalDefaults` fields SHALL all be empty strings and no error returned

---

### Requirement: Dialogue examples field on Character

The `Character` struct SHALL expose `DialogueExamples []string` (YAML key `dialogue_examples`). Each entry is a representative spoken line that anchors the character's voice. The field SHOULD contain 3–4 entries; it is not enforced at load time.

#### Scenario: Dialogue examples parsed from YAML list
- **WHEN** `characters.yaml` contains a `dialogue_examples` list with 4 entries
- **THEN** `Character.DialogueExamples` SHALL have length 4 with entries matching the YAML values

#### Scenario: Dialogue examples absent defaults to nil slice
- **WHEN** `characters.yaml` omits `dialogue_examples`
- **THEN** `Character.DialogueExamples` SHALL be nil (or empty slice) and no error returned

---

### Requirement: Inventory and initial state fields on Character

The `Character` struct SHALL expose `Inventory []string` (YAML key `inventory`) — an ordered list of objects the character carries at the start of the simulation — and `InitialState string` (YAML key `initial_state`) — a short description of the character's tactical or narrative state at simulation start. Both fields are optional; omitting them SHALL leave `Inventory` as nil and `InitialState` as an empty string.

#### Scenario: Inventory parsed from YAML list
- **WHEN** a `characters.yaml` entry contains an `inventory` list with two items
- **THEN** `Character.Inventory` SHALL have length 2 with entries in YAML order and no error returned

#### Scenario: Inventory absent defaults to nil
- **WHEN** a `characters.yaml` entry omits the `inventory` key
- **THEN** `Character.Inventory` SHALL be nil and loading SHALL return no error

#### Scenario: InitialState parsed from YAML string
- **WHEN** a `characters.yaml` entry contains `initial_state: "ready to infiltrate"`
- **THEN** `Character.InitialState` SHALL equal `"ready to infiltrate"` after loading

#### Scenario: InitialState absent defaults to empty string
- **WHEN** a `characters.yaml` entry omits the `initial_state` key
- **THEN** `Character.InitialState` SHALL be an empty string and loading SHALL return no error

---

### Requirement: Gender field on Character

The `Character` struct SHALL expose a `Gender string` field (YAML key `gender`). The field is optional; omitting it SHALL leave `Gender` as an empty string and SHALL not cause a load error. No validation of allowed values is performed — the field is free-form.

#### Scenario: Gender parsed from YAML
- **WHEN** a `characters.yaml` entry includes `gender: "femenino"`
- **THEN** `LoadCharacters` SHALL populate `Character.Gender` with `"femenino"` and return no error

#### Scenario: Gender absent defaults to empty string
- **WHEN** a `characters.yaml` entry omits the `gender` key
- **THEN** `Character.Gender` SHALL be an empty string and loading SHALL return no error

---

### Requirement: Removal of Personality and Backstory fields

The `Character` struct SHALL NOT contain `Personality []string` or `Backstory string` fields. YAML files containing these keys SHALL have them silently ignored by the loader.

#### Scenario: Old YAML with personality and backstory loads without error
- **WHEN** a `characters.yaml` entry includes `personality` and `backstory` keys that are no longer in the struct
- **THEN** `LoadCharacters` SHALL return a populated `Character` with no error, ignoring the unknown keys

---

### Requirement: BuildSystemPrompt uses the structured psychological template

`BuildSystemPrompt(c character.Character) string` SHALL produce a system prompt using the following section order, omitting any section whose fields are all empty:

1. Identity line: `"You are {Name}, a {Age}-year-old {Gender} {Occupation}."` when `Gender` is non-empty, or `"You are {Name}, a {Age}-year-old {Occupation}."` when `Gender` is empty.
2. `Motivación:` line using `Motivation`
3. `Miedo:` line using `Fear`
4. `Creencia central:` line using `CoreBelief`
5. `Tensión interna:` line using `InternalTension`
6. `Eventos formativos:` bulleted block using `FormativeEvents`
7. `Voz:` block using Voice sub-fields (formality, verbal tics, response length, humor, style)
8. `Relaciones default:` block using RelationalDefaults (strangers, authority, vulnerable)
9. `Objetivos:` bulleted block using `Goals`
10. `Estado emocional actual:` line using `EmotionalState`
11. `Ejemplos de diálogo:` quoted block using `DialogueExamples`
12. Closing instruction: `"Stay in character at all times. Respond as this person would. Keep responses concise."`

#### Scenario: Gender included in identity line when present
- **WHEN** `BuildSystemPrompt` is called with a `Character` where `Gender` is `"femenino"` and `Occupation` is `"Detective"`
- **THEN** the returned string SHALL contain `"femenino Detective"` in the identity line

#### Scenario: Gender omitted from identity line when empty
- **WHEN** `BuildSystemPrompt` is called with a `Character` where `Gender` is `""`
- **THEN** the identity line SHALL follow the original format without a gender token

#### Scenario: Full character produces all sections
- **WHEN** `BuildSystemPrompt` is called with a `Character` where all fields are populated
- **THEN** the returned string SHALL contain the identity line, all labeled sections, and the closing instruction

#### Scenario: Empty fields are silently omitted
- **WHEN** `BuildSystemPrompt` is called with a `Character` where `FormativeEvents` is nil and `InternalTension` is empty
- **THEN** the returned string SHALL NOT contain `"Tensión interna:"` or `"Eventos formativos:"` sections

#### Scenario: Minimal character (name, age, occupation only) produces a valid prompt
- **WHEN** `BuildSystemPrompt` is called with a `Character` that has only `Name`, `Age`, and `Occupation` set
- **THEN** the returned string SHALL contain only the identity line and the closing instruction, with no empty section headers

---

### Requirement: CHARACTER_RULES.md generation rulebook

The system SHALL provide `simulations/CHARACTER_RULES.md` — a Markdown document usable as a prompt context for LLMs tasked with creating new characters. It SHALL include: the complete YAML template with all fields, field-by-field descriptions of what each field should contain, anti-pattern guidance (what NOT to do), and at least one fully worked example character in the new schema.

#### Scenario: Rulebook file exists and is readable
- **WHEN** the simulator repository is cloned
- **THEN** `simulations/CHARACTER_RULES.md` SHALL exist as a non-empty Markdown file

#### Scenario: Rulebook contains all required sections
- **WHEN** `CHARACTER_RULES.md` is read
- **THEN** it SHALL contain sections covering: the YAML template, field descriptions, anti-patterns, and a worked example

---

### Requirement: Appearance field on Character

The `Character` struct SHALL expose an `Appearance string` field (YAML key `appearance`). This field contains a single authored sentence describing how this character presents to others on first encounter — their visible manner, posture, or energy — without revealing internal psychology. The field is optional; omitting it in YAML SHALL leave the field as an empty string.

`ObservableSnapshot` SHALL include `Appearance` in the observable profile passed to judgment formation. `BuildSystemPrompt` SHALL NOT include `Appearance` in a character's own system prompt (it describes how others see them, not how they see themselves).

#### Scenario: Appearance parsed from YAML
- **WHEN** a `characters.yaml` entry contains `appearance: "Carries herself with the controlled stillness of someone who observes before acting"`
- **THEN** `Character.Appearance` SHALL equal that string after loading and no error returned

#### Scenario: Appearance absent defaults to empty string
- **WHEN** a `characters.yaml` entry omits the `appearance` key
- **THEN** `Character.Appearance` SHALL be an empty string and `LoadCharacters` SHALL return no error

#### Scenario: Appearance included in observable snapshot
- **WHEN** `ObservableSnapshot` is called on a character with a non-empty `Appearance`
- **THEN** the returned `ObservableProfile.Appearance` SHALL equal the character's `Appearance` value

#### Scenario: Appearance absent from character's own system prompt
- **WHEN** `BuildSystemPrompt` is called with a character who has a non-empty `Appearance`
- **THEN** the returned system prompt SHALL NOT contain the `Appearance` text

---

### Requirement: Appearance added to all existing scenario characters

All `characters.yaml` files in the `simulations/` directory SHALL include an `appearance` field on every non-director character entry. Director characters (type `game_director`) SHALL NOT have an `appearance` field.

#### Scenario: Existing scenario characters have appearance field
- **WHEN** any `characters.yaml` in `simulations/default/`, `simulations/honey-heist/`, `simulations/doom-hell-crusade/`, or `simulations/test-scenario/` is loaded
- **THEN** every non-director character SHALL have a non-empty `Appearance` after loading
