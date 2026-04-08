## Requirements

### Requirement: CoverIdentity struct on Character

The system SHALL define a `CoverIdentity` struct with four fields: `Alias string` (YAML key `alias`) — the name the character uses in this world; `Role string` (YAML key `role`) — their claimed occupation or social position; `Backstory string` (YAML key `backstory`) — a single sentence explaining their invented personal history; and `Weaknesses []string` (YAML key `weaknesses`) — a list of specific behaviours, topics, or situations that could expose their true nature. The `Character` struct SHALL expose a `CoverIdentity *CoverIdentity` field (YAML key `cover_identity`). A nil pointer means the character has no cover — they are who they appear to be.

#### Scenario: Full cover_identity block parsed from characters.yaml
- **WHEN** a `characters.yaml` character entry contains a `cover_identity:` block with `alias`, `role`, `backstory`, and `weaknesses` set
- **THEN** `Character.CoverIdentity` SHALL be non-nil and all four fields SHALL be populated with those values

#### Scenario: Partial cover_identity block accepted
- **WHEN** a `characters.yaml` entry contains `cover_identity: { alias: "Gerald", role: "Honey sommelier" }` with no `backstory` or `weaknesses`
- **THEN** `Character.CoverIdentity` SHALL be non-nil, `Alias` and `Role` SHALL be set, `Backstory` SHALL be empty, `Weaknesses` SHALL be nil/empty, and loading SHALL return no error

#### Scenario: Missing cover_identity block leaves pointer nil
- **WHEN** a `characters.yaml` entry omits `cover_identity` entirely
- **THEN** `Character.CoverIdentity` SHALL be nil and loading SHALL return no error

#### Scenario: Weaknesses parsed as ordered list
- **WHEN** `cover_identity.weaknesses` contains two entries
- **THEN** `CoverIdentity.Weaknesses` SHALL have length 2 in YAML order

---

### Requirement: Cover identity injected into character system prompt

`BuildSystemPrompt` SHALL append a "Cover Identity" section when `Character.CoverIdentity` is non-nil. The section SHALL include the alias, role, backstory (when non-empty), and weaknesses as a bulleted list (when non-empty). When `CoverIdentity` is nil, no cover identity section SHALL appear in the prompt.

#### Scenario: Cover identity present in system prompt when set
- **WHEN** a character has a non-nil `CoverIdentity` with `Alias: "Gerald"` and `Role: "Honey sommelier"`
- **THEN** the string returned by `BuildSystemPrompt` SHALL contain both `"Gerald"` and `"Honey sommelier"`

#### Scenario: Weaknesses listed in system prompt
- **WHEN** `CoverIdentity.Weaknesses` contains `["Gets excited by raw honey", "Uses bear paw metaphors"]`
- **THEN** the system prompt SHALL contain both weakness strings

#### Scenario: No cover identity section when nil
- **WHEN** `Character.CoverIdentity` is nil
- **THEN** `BuildSystemPrompt` SHALL NOT contain the words "Cover Identity" or "alias"
