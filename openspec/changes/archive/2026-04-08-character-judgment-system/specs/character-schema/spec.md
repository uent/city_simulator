# Character Schema — Delta

## ADDED Requirements

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
