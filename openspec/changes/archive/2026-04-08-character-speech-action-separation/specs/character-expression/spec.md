## ADDED Requirements

### Requirement: Expression type
The system SHALL define an `Expression` struct in `internal/character/expression.go` with two string fields: `Speech` (what the character says aloud) and `Action` (a brief third-person description of what the character does physically, may be empty).

#### Scenario: Expression with both fields populated
- **WHEN** a raw LLM response contains `*picks up the envelope* I wasn't expecting this.`
- **THEN** `Expression.Action` SHALL equal `picks up the envelope` and `Expression.Speech` SHALL equal `I wasn't expecting this.`

#### Scenario: Expression with no action markers
- **WHEN** a raw LLM response contains no `*...*` markers
- **THEN** `Expression.Action` SHALL be an empty string and `Expression.Speech` SHALL equal the full trimmed response

#### Scenario: Expression with action marker only
- **WHEN** a raw LLM response is only `*walks away silently*` with no remaining text
- **THEN** `Expression.Action` SHALL equal `walks away silently` and `Expression.Speech` SHALL be an empty string

---

### Requirement: ParseExpression function
The system SHALL provide a `ParseExpression(raw string) Expression` function that:
1. Finds the first `*...*` substring in `raw` (open `*` to first closing `*`).
2. Assigns the content between the markers (trimmed) to `Expression.Action`.
3. Concatenates the text before and after the marker block (both trimmed), separated by a single space if both are non-empty, and assigns the result to `Expression.Speech`.
4. If no `*...*` marker is found, assigns the full trimmed `raw` to `Expression.Speech` and leaves `Expression.Action` empty.

#### Scenario: Action marker at the start of text
- **WHEN** raw is `*looks around nervously* You again?`
- **THEN** `Action` is `looks around nervously` and `Speech` is `You again?`

#### Scenario: Action marker in the middle of text
- **WHEN** raw is `I told you *slams fist on table* I don't know anything.`
- **THEN** `Action` is `slams fist on table` and `Speech` is `I told you I don't know anything.`

#### Scenario: Action marker at the end of text
- **WHEN** raw is `Fine. Have it your way. *turns and walks out*`
- **THEN** `Action` is `turns and walks out` and `Speech` is `Fine. Have it your way.`

#### Scenario: Nested asterisks use the first open and first close
- **WHEN** raw contains `*says "yes" firmly*`
- **THEN** the parser SHALL NOT produce incorrect output by treating inner quotes as markers; `Action` SHALL equal `says "yes" firmly`

#### Scenario: Empty or whitespace-only input
- **WHEN** raw is an empty string or contains only whitespace
- **THEN** both `Speech` and `Action` SHALL be empty strings

---

### Requirement: Expression format instruction in character system prompt
The system SHALL append an expression format instruction to every character system prompt built by `BuildSystemPrompt`. The instruction SHALL explain the `*action*` convention and tell the character to use it when performing a physical action alongside speech.

The instruction text SHALL be consistent and placed after the existing "Stay in character" closing line.

#### Scenario: System prompt includes expression instruction
- **WHEN** `BuildSystemPrompt` is called for any character
- **THEN** the returned string SHALL contain a line instructing the character to prefix physical actions with `*...*` markers

#### Scenario: Instruction is language-agnostic
- **WHEN** `BuildSystemPrompt` is called with a non-empty `language` parameter
- **THEN** the expression format instruction SHALL appear in the prompt regardless of the language setting

---

### Requirement: FormatExpression helper
The system SHALL provide a `FormatExpression(e Expression) string` helper that recombines an `Expression` into the `*action* speech` wire format used when injecting one character's turn into another character's context:
- If `Action` is non-empty: returns `*{action}* {speech}` (trimmed).
- If `Action` is empty: returns `Speech` as-is.
- If `Speech` is empty and `Action` is non-empty: returns `*{action}*`.

#### Scenario: Both fields present
- **WHEN** `Action` is `picks up the card` and `Speech` is `Interesting.`
- **THEN** `FormatExpression` returns `*picks up the card* Interesting.`

#### Scenario: Action empty
- **WHEN** `Action` is empty and `Speech` is `What do you want?`
- **THEN** `FormatExpression` returns `What do you want?`

#### Scenario: Speech empty
- **WHEN** `Action` is `nods silently` and `Speech` is empty
- **THEN** `FormatExpression` returns `*nods silently*`
