## MODIFIED Requirements

### Requirement: Honey Heist character roster

The `characters.yaml` file SHALL define exactly 6 bear characters, each with a unique `name`, an `occupation` describing their criminal specialty, a `motivation` string, a `fear` string, a `core_belief` string, an `internal_tension` string, a `formative_events` list of 2–3 causal bullets, a `voice` block with at least `formality` and `verbal_tics`, a `relational_defaults` block with `strangers`, `authority`, and `vulnerable`, and a `dialogue_examples` list of 3–4 representative lines. The legacy `personality` list and `backstory` prose fields SHALL NOT be present.

#### Scenario: All characters present with new schema
- **WHEN** the honey-heist scenario is loaded
- **THEN** `Scenario.Characters` SHALL contain 6 entries, each with non-empty `Name`, `Occupation`, `Motivation`, `Fear`, `CoreBelief`, and `InternalTension` fields

#### Scenario: Character names are unique
- **WHEN** the honey-heist scenario is loaded
- **THEN** no two characters in the roster SHALL share the same `Name`

#### Scenario: Each character has formative events
- **WHEN** the honey-heist scenario is loaded
- **THEN** each character's `FormativeEvents` slice SHALL have length 2 or 3

#### Scenario: Each character has dialogue examples
- **WHEN** the honey-heist scenario is loaded
- **THEN** each character's `DialogueExamples` slice SHALL have length 3 or 4

#### Scenario: Legacy fields absent from YAML
- **WHEN** `simulations/honey-heist/characters.yaml` is read as raw text
- **THEN** it SHALL NOT contain the keys `personality` or `backstory` at the top level of any character entry
