## ADDED Requirements

### Requirement: Gender field on Character

The `Character` struct SHALL expose a `Gender string` field (YAML key `gender`). The field is optional; omitting it SHALL leave `Gender` as an empty string and SHALL not cause a load error. No validation of allowed values is performed — the field is free-form.

#### Scenario: Gender parsed from YAML

- **WHEN** a `characters.yaml` entry includes `gender: "femenino"`
- **THEN** `LoadCharacters` SHALL populate `Character.Gender` with `"femenino"` and return no error

#### Scenario: Gender absent defaults to empty string

- **WHEN** a `characters.yaml` entry omits the `gender` key
- **THEN** `Character.Gender` SHALL be an empty string and loading SHALL return no error

## MODIFIED Requirements

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
