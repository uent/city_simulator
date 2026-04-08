## MODIFIED Requirements

### Requirement: Psychological core fields on Character

The `Character` struct SHALL expose four string fields that form its psychological core: `Motivation` (what the character wants and why), `Fear` (what they avoid at all costs), `CoreBelief` (their foundational view of how the world works), and `InternalTension` (a single contradiction that defines their complexity). The YAML keys SHALL be `motivation`, `fear`, `core_belief`, and `internal_tension`. The struct SHALL also expose a `CoverIdentity *CoverIdentity` field (YAML key `cover_identity`) as specified in the character-cover-identity capability.

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
