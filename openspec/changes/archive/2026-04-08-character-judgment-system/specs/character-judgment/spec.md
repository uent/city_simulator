# Character Judgment

## ADDED Requirements

### Requirement: CharacterJudgment struct

The system SHALL define a `CharacterJudgment` struct in the `character` package with the following fields:
- `About string` — ID of the character being judged
- `Name string` — name snapshot at judgment time (cover alias if applicable)
- `Impression string` — first-person internal narrative opinion (2–3 sentences)
- `Trust string` — one of: `"high"`, `"medium"`, `"low"`, `"none"`
- `Interest string` — one of: `"high"`, `"medium"`, `"low"`
- `Threat string` — one of: `"high"`, `"medium"`, `"low"`, `"none"`
- `FormedTick int` — tick at which the judgment was formed (0 = pre-simulation)
- `UpdatedTick int` — tick at which the judgment was last updated (0 = never updated)

The `Character` struct SHALL include a `Judgments map[string]CharacterJudgment` field (YAML tag `"-"`, runtime-only, never persisted) keyed by character ID.

#### Scenario: Judgment fields hold correct enum values
- **WHEN** a judgment is parsed from a valid LLM response with `trust: "low"`, `interest: "high"`, `threat: "none"`
- **THEN** the resulting `CharacterJudgment` SHALL have `Trust == "low"`, `Interest == "high"`, `Threat == "none"`

#### Scenario: Invalid enum values fall back to defaults
- **WHEN** the LLM response contains an unrecognized value for `trust` (e.g., `"unsure"`)
- **THEN** `ParseJudgmentResponse` SHALL replace it with `"medium"` and return no error

---

### Requirement: Observable snapshot respects CoverIdentity

`ObservableSnapshot(c Character) ObservableProfile` SHALL return a filtered view of a character containing only what another character could realistically observe on first encounter:
- If `CoverIdentity` is non-nil: `Name = CoverIdentity.Alias`, `Occupation = CoverIdentity.Role`
- If `CoverIdentity` is nil: `Name = Character.Name`, `Occupation = Character.Occupation`
- Always included: `Age`, `EmotionalState`, `Appearance`, `Location`
- Never included: `Motivation`, `Fear`, `CoreBelief`, `InternalTension`, `FormativeEvents`, `Goals`, `CoverIdentity` internals

#### Scenario: Cover identity character produces alias-based snapshot
- **WHEN** `ObservableSnapshot` is called on a character with `CoverIdentity.Alias = "Don Gregorio Wald"` and `CoverIdentity.Role = "Coleccionista privado"`
- **THEN** the returned profile SHALL have `Name = "Don Gregorio Wald"` and `Occupation = "Coleccionista privado"`, with `Motivation` absent

#### Scenario: Character without cover identity uses real fields
- **WHEN** `ObservableSnapshot` is called on a character with no `CoverIdentity`
- **THEN** the returned profile SHALL have `Name = Character.Name` and `Occupation = Character.Occupation`

---

### Requirement: Initial judgment formation before tick 1

`FormInitialJudgments(ctx, chars, llmClient, language)` SHALL be called by the engine before the tick loop begins. For every ordered pair (A, B) where A ≠ B, it SHALL execute one LLM call to form A's judgment of B, using:
- A's full psychological profile as the judging lens
- B's `ObservableProfile` as the observable input
- A JSON response format: `{ "impression": string, "trust": enum, "interest": enum, "threat": enum }`

All N×(N-1) LLM calls SHALL execute concurrently. Results SHALL be written to `char.Judgments[targetID]` on each character. `FormedTick` SHALL be set to 0 for all initial judgments.

#### Scenario: All pairs receive judgments before tick 1
- **WHEN** a simulation with 3 characters (A, B, C) starts
- **THEN** before tick 1, `A.Judgments` SHALL contain entries for B and C, `B.Judgments` for A and C, and `C.Judgments` for A and B (6 total judgments)

#### Scenario: LLM failure for one pair does not block others
- **WHEN** the LLM call for pair (A→B) returns an error
- **THEN** the simulation SHALL continue with a neutral fallback judgment for that pair, and all other pairs SHALL be unaffected

---

### Requirement: Judgment injected into conversation system prompts

When building the system prompt for each participant in a conversation between A and B, the engine SHALL append A's judgment of B to A's system prompt, and B's judgment of A to B's system prompt.

The injected block SHALL be formatted as:
```
Your prior impression of [Name] ([occupation]):
"[impression]"
Trust: [trust] | Interest: [interest] | Perceived threat: [threat]
```

Only the judgment for the current conversation partner SHALL be injected — not the full `Judgments` map.

If no judgment exists for the partner (e.g., judgment formation failed), no block SHALL be injected.

#### Scenario: Judgment block appears in initiator system prompt
- **WHEN** A initiates a conversation with B and `A.Judgments[B.ID]` exists
- **THEN** A's system prompt SHALL contain the formatted judgment block for B

#### Scenario: Missing judgment is silently skipped
- **WHEN** `A.Judgments[B.ID]` does not exist
- **THEN** A's system prompt SHALL be built without any judgment block, with no error

---

### Requirement: Judgment refreshed every 10 conversations per pair

The engine SHALL maintain a `pairConversations map[string]int` counter, keyed by `sortedID(A) + ":" + sortedID(B)`. After each conversation between A and B, the counter SHALL be incremented.

When `count % 10 == 0 && count > 0`, the engine SHALL call `UpdateJudgment` for both A→B and B→A concurrently. The update LLM call SHALL include:
- The judging character's full psychological profile
- The target's current `ObservableProfile`
- The prior judgment text and levels
- The last 5 public conversation events from `world.EventLog` involving both participants as contextual history

The resulting judgment SHALL overwrite `char.Judgments[targetID]` with `UpdatedTick` set to the current tick.

#### Scenario: Judgment update triggers at exactly 10 conversations
- **WHEN** A and B have had 10 conversations
- **THEN** `UpdateJudgment` SHALL be called for both A→B and B→A after that 10th conversation

#### Scenario: Judgment update does not trigger at 9 conversations
- **WHEN** A and B have had 9 conversations
- **THEN** `UpdateJudgment` SHALL NOT be called

#### Scenario: Update preserves judgment for unrelated pairs
- **WHEN** A and B's judgment is updated
- **THEN** A's judgments of C, D, etc. SHALL remain unchanged

---

### Requirement: Judgments formed for director-spawned characters

When the engine registers a new character spawned by the game director (in `registerSpawnedChars`), it SHALL:
1. Form judgments from the new character to all previously registered characters (new → existing)
2. Form judgments from all previously registered characters to the new character (existing → new)

Both sets of calls SHALL complete before the simulation tick proceeds. Existing `pairConversations` counters for the new character SHALL start at 0.

#### Scenario: Spawned character receives judgments from existing cast
- **WHEN** the director spawns a new character X during tick 5
- **THEN** before tick 6, every existing character SHALL have an entry in their `Judgments` map for X

#### Scenario: Spawned character judges all existing characters
- **WHEN** the director spawns a new character X during tick 5
- **THEN** before tick 6, `X.Judgments` SHALL contain entries for every previously existing character
