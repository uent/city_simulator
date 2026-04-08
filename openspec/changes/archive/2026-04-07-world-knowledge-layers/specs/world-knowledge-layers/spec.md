## ADDED Requirements

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

The `Character` struct SHALL expose a `Location string` field (YAML key `location`) indicating the character's current position in the world, expressed as the `Name` of a `Location` from the scenario's world config. This field is set from `characters.yaml` and used by the conversation manager to determine which `LocalContext` to inject into the character's system prompt.

#### Scenario: Location field parsed from YAML
- **WHEN** a `characters.yaml` entry includes `location: "Town Square"`
- **THEN** `Character.Location` SHALL equal `"Town Square"` after loading

#### Scenario: Missing location field defaults to empty string
- **WHEN** a `characters.yaml` entry omits `location`
- **THEN** `Character.Location` SHALL be an empty string and no error returned

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

### Requirement: State.Summary removed

The existing `State.Summary() string` method SHALL be removed. All callers SHALL use `PublicSummary()` and/or `LocalContext()` instead.

**Reason**: `Summary()` collapses public and local knowledge into a single global view, which is the root of the omniscience problem this change solves.

**Migration**: Replace `w.Summary()` calls with `w.PublicSummary()` for global context. For per-character context, use `w.PublicSummary() + "\n" + w.LocalContext(character.Location)`.

#### Scenario: Summary method does not exist after change
- **WHEN** `go build ./...` is run after the change is applied
- **THEN** there SHALL be no compilation errors referencing `State.Summary`
