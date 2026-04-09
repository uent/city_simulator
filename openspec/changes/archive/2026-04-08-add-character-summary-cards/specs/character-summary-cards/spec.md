## ADDED Requirements

### Requirement: Character card rendering
The `internal/summary` package SHALL expose a `renderCharacterCards(chars []*character.Character) string` function that produces a Markdown block containing one card per non-director character. The block SHALL start with a horizontal rule (`---`) followed by a `## Character Cards` heading, then one `### <Name>` section per character. Game Director characters (`Type == "game_director"`) SHALL be excluded. Fields that are empty strings, zero values, or nil SHALL be omitted from the card. The card SHALL display the following fields when present: Name (as the section heading), Age, Occupation, Appearance, Motivation, Fear, Core Belief, Internal Tension, Formative Events, Location, Emotional State, Goals, Voice (Formality, Verbal Tics, Response Length, Humor Type, Communication Style), Relational Defaults (Strangers, Authority, Vulnerable), Dialogue Examples, Cover Identity (Alias, Role, Backstory, Weaknesses), and Relationships (one entry per `CharacterJudgment` in `Character.Judgments`, showing the known character's name, trust/interest/threat levels, and impression).

#### Scenario: All fields populated
- **WHEN** `renderCharacterCards` is called with a character that has all fields set
- **THEN** the output SHALL contain a `### <Name>` section with each non-empty field rendered as a labeled Markdown list item or sub-list

#### Scenario: Empty fields are omitted
- **WHEN** a character has empty `Fear`, `CoreBelief`, and `InternalTension`
- **THEN** those labels SHALL NOT appear in that character's card

#### Scenario: Game Director excluded
- **WHEN** the character slice contains one character with `Type == "game_director"`
- **THEN** no card SHALL be rendered for that character

#### Scenario: No non-director characters
- **WHEN** all characters in the slice are game directors or the slice is empty
- **THEN** `renderCharacterCards` SHALL return an empty string

#### Scenario: Cover identity present
- **WHEN** a character has a non-nil `CoverIdentity` with `Alias`, `Role`, and `Backstory` set
- **THEN** the card SHALL include a "Cover Identity" sub-section with those fields

#### Scenario: Nil cover identity omitted
- **WHEN** `Character.CoverIdentity` is nil
- **THEN** no "Cover Identity" section SHALL appear in that character's card

#### Scenario: Relationships rendered when judgments exist
- **WHEN** a character's `Judgments` map contains one or more entries
- **THEN** the card SHALL include a "Relationships" sub-section with one line per judgment showing the known character's name, trust, interest, and threat levels, followed by the impression as a blockquote

#### Scenario: No relationships when judgments map is empty
- **WHEN** a character's `Judgments` map is empty
- **THEN** no "Relationships" section SHALL appear in that character's card
