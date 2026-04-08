# Spec: simulation-language

## Requirements

### Requirement: SIM_LANGUAGE configures narrator and character response language
The binary SHALL read a `SIM_LANGUAGE` environment variable and expose a `--language` CLI flag. The resolved value SHALL be injected into the system prompt of every character actor and the Game Director, instructing the LLM to respond in that language.

#### Scenario: Language instruction appears in character prompt
- **WHEN** `SIM_LANGUAGE=Spanish` is set
- **THEN** each character's system prompt SHALL end with `"Respond in Spanish."`

#### Scenario: Language instruction appears in director prompt
- **WHEN** `SIM_LANGUAGE=Spanish` is set
- **THEN** the Game Director's system prompt SHALL contain `"Respond in Spanish."`

#### Scenario: Language instruction appears in movement prompt
- **WHEN** `SIM_LANGUAGE=French` is set
- **THEN** each character's movement prompt SHALL end with `"Respond in French."`

#### Scenario: Language instruction appears in simulation summary
- **WHEN** `SIM_LANGUAGE=Spanish` is set
- **THEN** the simulation summary system prompt SHALL contain `"Respond in Spanish."`

#### Scenario: No instruction injected when language is unset
- **WHEN** `SIM_LANGUAGE` is not set and `--language` is not passed
- **THEN** neither the character, director, nor summary prompts SHALL contain a language instruction line

#### Scenario: CLI flag overrides env var
- **WHEN** `SIM_LANGUAGE=Spanish` is set and `--language=English` is also passed
- **THEN** the binary SHALL use `"English"` as the language value

### Requirement: Language value is accepted as a free-form string
The binary SHALL accept any non-empty string as the `SIM_LANGUAGE` value and forward it verbatim to the LLM prompts. The binary SHALL NOT validate or normalise the language identifier.

#### Scenario: BCP-47 tag accepted
- **WHEN** `SIM_LANGUAGE=es` is set
- **THEN** prompts SHALL contain `"Respond in es."` without error

#### Scenario: Full language name accepted
- **WHEN** `SIM_LANGUAGE=Español` is set
- **THEN** prompts SHALL contain `"Respond in Español."` without error
