## Requirements

### Requirement: Location details field

The `Location` struct SHALL expose a `Details string` field (YAML key `details`). `Details` contains rich, private information about the location — atmosphere, secrets, specific observations — visible only to characters currently present there. It is distinct from `Description`, which remains the public one-liner shown to all characters.

#### Scenario: Details field parsed from YAML
- **WHEN** a `world.yaml` location entry includes a `details` key
- **THEN** `Location.Details` SHALL be populated with that value and `Location.Description` SHALL remain unchanged

#### Scenario: Missing details field defaults to empty string
- **WHEN** a `world.yaml` location entry omits `details`
- **THEN** `Location.Details` SHALL be an empty string and `Load` SHALL return no error

---

### Requirement: Event visibility and location fields

The `Event` struct SHALL expose two new fields: `Visibility string` (YAML key `visibility`, values `"public"` or `"local"`) and `Location string` (YAML key `location`, the name of the location where the event occurred). Events with `visibility: "public"` appear in every character's world context. Events with `visibility: "local"` appear only in the `LocalContext` of characters at the matching location. The default visibility when omitted SHALL be `"public"`.

#### Scenario: Public event visible in PublicSummary
- **WHEN** an event has `visibility: "public"` or omits `visibility`
- **THEN** `State.PublicSummary()` SHALL include that event in its recent-events section

#### Scenario: Local event not visible in PublicSummary
- **WHEN** an event has `visibility: "local"`
- **THEN** `State.PublicSummary()` SHALL NOT include that event

#### Scenario: Local event visible in matching LocalContext
- **WHEN** an event has `visibility: "local"` and `location: "Tavern"`
- **THEN** `State.LocalContext("Tavern")` SHALL include that event in its recent-events section

#### Scenario: Local event not visible in non-matching LocalContext
- **WHEN** an event has `visibility: "local"` and `location: "Tavern"`
- **THEN** `State.LocalContext("Market")` SHALL NOT include that event

#### Scenario: Omitted visibility defaults to public
- **WHEN** an event in `world.yaml` omits the `visibility` key
- **THEN** the event SHALL behave as if `visibility: "public"` were set

---

### Requirement: PublicSummary method

`State` SHALL expose a `PublicSummary() string` method that returns the universal world context available to all characters. It SHALL include: the current time of day, the list of all location names with their `Description` (not `Details`), and the last 5 events whose `Visibility` is `"public"`. If there are no public events, only time of day and location list are included.

#### Scenario: Returns time and location names
- **WHEN** `PublicSummary()` is called on a `State` with two locations and no events
- **THEN** the returned string SHALL contain the time-of-day label and both location names

#### Scenario: Includes only public events
- **WHEN** the event log contains one public event and one local event
- **THEN** `PublicSummary()` SHALL include only the public event

#### Scenario: Limits to 5 most recent public events
- **WHEN** the event log contains 8 public events
- **THEN** `PublicSummary()` SHALL include only the 5 most recent ones

---

### Requirement: LocalContext method

`State` SHALL expose a `LocalContext(locationID string) string` method that returns private context for a specific location. It SHALL include: the `Details` of the matching location (if non-empty) and the last 5 events whose `Visibility` is `"local"` and `Location` matches `locationID`. If `locationID` does not match any known location name, `LocalContext` SHALL return an empty string and emit a warning log.

#### Scenario: Returns location details for matching location
- **WHEN** `LocalContext("Tavern")` is called and `Tavern` has a non-empty `Details` field
- **THEN** the returned string SHALL contain the `Details` text

#### Scenario: Returns empty string for unknown location
- **WHEN** `LocalContext("Nonexistent")` is called and no location with that name exists
- **THEN** the returned string SHALL be empty

#### Scenario: Includes only local events for that location
- **WHEN** the event log contains a local event at "Tavern" and a local event at "Market"
- **THEN** `LocalContext("Tavern")` SHALL include only the Tavern event

#### Scenario: Limits to 5 most recent local events for location
- **WHEN** the event log contains 8 local events all at "Tavern"
- **THEN** `LocalContext("Tavern")` SHALL include only the 5 most recent ones

---

### Requirement: Character location field

The `Character` struct SHALL expose a `Location string` runtime field indicating the character's current position in the world, expressed as the `Name` of a `Location` from the scenario's world config. This field is NOT loaded from `characters.yaml`; it is assigned and updated entirely at runtime by the simulation (initially by the `Scheduler`, then updated by movement decisions). The conversation manager uses it to determine which `LocalContext` to inject into the character's system prompt.

#### Scenario: Location field is populated before first tick
- **WHEN** `NewEngine` initialises the simulation
- **THEN** every character SHALL have a non-empty `Location` set to a valid location name before `Run` is called

#### Scenario: Location updates after movement decision
- **WHEN** a character's movement decision resolves to a location name different from their current one
- **THEN** `Character.Location` SHALL be updated to that name before the next tick

#### Scenario: Character with empty location receives only public summary
- **WHEN** a character's `Location` field is empty at exchange time
- **THEN** that character's world context SHALL consist of only `PublicSummary()` with no local context appended

---

### Requirement: Per-character world context in RunExchange

`RunExchange` SHALL build a distinct world context string for each character by combining `State.PublicSummary()` with `State.LocalContext(character.Location)`. This combined string SHALL replace the global `w.Summary()` currently injected into each character's system prompt. If `character.Location` is empty, only `PublicSummary()` is used.

#### Scenario: Initiator and responder receive different local context
- **WHEN** initiator is at "Tavern" and responder is at "Market" and each location has local events
- **THEN** the initiator's system prompt SHALL contain the Tavern local context and NOT the Market local context, and vice versa for the responder

#### Scenario: Character with no location receives only public summary
- **WHEN** a character's `Location` field is empty
- **THEN** that character's world context SHALL consist of only `PublicSummary()` with no local context appended

---

### Requirement: Scheduler assigns initial character locations

The `Scheduler` SHALL assign a random starting location to every character during initialisation, drawing from the list of location names provided by the scenario's `WorldConfig`. No character SHALL begin a simulation with an empty `Location`. The `characters.yaml` file SHALL NOT be the source of location data; the YAML `location` key on characters is ignored.

#### Scenario: All characters have a location after NewEngine
- **WHEN** `NewEngine` is called with a valid scenario containing at least one location
- **THEN** every character in `Engine.chars` SHALL have a non-empty `Location` equal to one of the world's location names

#### Scenario: Distribution across locations
- **WHEN** a scenario has 6 characters and 6 locations and a fixed seed
- **THEN** the initial location assignment SHALL be deterministic and repeatable for that seed

---

### Requirement: LLM-driven movement decision after each exchange

After each tick's exchange completes, the simulation SHALL ask the LLM to decide where each participant moves next. The system SHALL provide `Manager.DecideMovement(ctx, character, locations) string` which: builds a movement prompt using the character's motivation, fear, current location, and available locations; calls the LLM; and returns either an exact location name from the available list or `"stay"`. If the LLM returns an unrecognisable response or errors, the result SHALL default to `"stay"`. The engine SHALL update `Character.Location` when the decision differs from the current location.

#### Scenario: Valid location name returned by LLM
- **WHEN** the LLM responds with a string that exactly matches a location name
- **THEN** `DecideMovement` SHALL return that location name

#### Scenario: Case-insensitive match accepted
- **WHEN** the LLM responds with a location name in a different case than stored
- **THEN** `DecideMovement` SHALL return the canonical location name

#### Scenario: Location name embedded in longer response
- **WHEN** the LLM responds with a sentence that contains a valid location name as a substring
- **THEN** `DecideMovement` SHALL extract and return that location name

#### Scenario: Unrecognisable response defaults to stay
- **WHEN** the LLM responds with a string that does not match any known location or "stay"
- **THEN** `DecideMovement` SHALL return `"stay"` and the character's location SHALL not change

#### Scenario: LLM error defaults to stay
- **WHEN** the LLM call returns an error
- **THEN** `DecideMovement` SHALL return `"stay"` without propagating the error

---

### Requirement: Movement prompt structure

`llm.BuildMovementPrompt(c character.Character, locations []string) string` SHALL construct a concise prompt that includes the character's name, motivation, fear, current location, and the list of available locations. It SHALL instruct the LLM to respond with exactly one location name from the list or the word `"stay"`.

#### Scenario: Prompt contains character context
- **WHEN** `BuildMovementPrompt` is called with a character who has a non-empty `Motivation`
- **THEN** the returned string SHALL contain that motivation text

#### Scenario: Prompt lists all available locations
- **WHEN** `BuildMovementPrompt` is called with three location names
- **THEN** the returned string SHALL contain all three location names

---

### Requirement: State.Summary removed

The existing `State.Summary() string` method SHALL be removed. All callers SHALL use `PublicSummary()` and/or `LocalContext()` instead.

**Reason**: `Summary()` collapses public and local knowledge into a single global view, which is the root of the omniscience problem this change solves.

**Migration**: Replace `w.Summary()` calls with `w.PublicSummary()` for global context. For per-character context, use `w.PublicSummary() + "\n" + w.LocalContext(character.Location)`.

#### Scenario: Summary method does not exist after change
- **WHEN** `go build ./...` is run after the change is applied
- **THEN** there SHALL be no compilation errors referencing `State.Summary`
