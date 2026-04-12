## Requirements

### Requirement: World state struct
The system SHALL define a `State` struct containing: current tick (int), time-of-day label (string, e.g., "morning", "afternoon", "evening", "night"), weather (string, default `""`), atmosphere (string, default `""`), tension (int, 0–10, default 0), a `Concept WorldConcept` field mirroring `WorldConfig.Concept`, list of locations (each with `Name`, `Description`, and `Details` string fields), and a shared event log (slice of `Event`).

An `Event` SHALL have: tick number, event type (string), description (string), optional participant character IDs (slice of strings), `Visibility` (string, `"public"` or `"local"`, default `"public"`), `Location` (string, name of the location where the event occurred, optional), `Target` (string, optional character or location ID the event is "about"), and `PrivateRecipient` (string, optional — if set, only that character should receive this event in their inbox).

#### Scenario: New world state creation from WorldConfig
- **WHEN** `NewState(cfg scenario.WorldConfig) *State` is called with a valid `WorldConfig`
- **THEN** the function SHALL return a `*State` with tick 0, time-of-day "morning", weather `""`, atmosphere `""`, tension 0, `Concept` copied from `cfg.Concept`, locations from `cfg.Locations`, and an event log pre-populated with `cfg.InitialEvents`

#### Scenario: Concept copied from WorldConfig into State
- **WHEN** `WorldConfig.Concept.Premise == "Bears pretending to be humans"`
- **THEN** `NewState(cfg).Concept.Premise` SHALL equal `"Bears pretending to be humans"`

#### Scenario: Initial events appear in log
- **WHEN** `WorldConfig.InitialEvents` contains one event
- **THEN** `State.EventLog` SHALL contain that event at index 0 immediately after `NewState` returns

#### Scenario: Location Details preserved in State
- **WHEN** `WorldConfig.Locations` contains a location with a non-empty `Details` field
- **THEN** `State.Locations` SHALL contain that location with `Details` intact

---

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

The `Event` struct SHALL expose two fields: `Visibility string` (YAML key `visibility`, values `"public"` or `"local"`) and `Location string` (YAML key `location`, the name of the location where the event occurred). Events with `visibility: "public"` appear in every character's world context. Events with `visibility: "local"` appear only in the `LocalContext` of characters at the matching location. The default visibility when omitted SHALL be `"public"`.

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

`State` SHALL expose a `PublicSummary() string` method that returns the universal world context available to all characters. It SHALL include: the current time of day, the list of all location names with their `Description` (not `Details`), weather and atmosphere when non-empty, tension level when non-zero (as a descriptor, e.g., "Tension: 7/10"), the last 5 events whose `Visibility` is `"public"`, and — when `State.Concept.Premise` is non-empty — a "World Rules" block containing the premise and, if present, the rules as a bulleted list. If there are no public events, only time of day, location list, and (if set) world rules are included.

#### Scenario: Returns time and location names
- **WHEN** `PublicSummary()` is called on a `State` with two locations and no events
- **THEN** the returned string SHALL contain the time-of-day label and both location names

#### Scenario: Includes only public events
- **WHEN** the event log contains one public event and one local event
- **THEN** `PublicSummary()` SHALL include only the public event

#### Scenario: Limits to 5 most recent public events
- **WHEN** the event log contains 8 public events
- **THEN** `PublicSummary()` SHALL include only the 5 most recent ones

#### Scenario: PublicSummary includes weather when set
- **WHEN** `state.Weather == "storm"`
- **THEN** `PublicSummary()` SHALL contain the word `"storm"`

#### Scenario: PublicSummary includes tension when non-zero
- **WHEN** `state.Tension == 7`
- **THEN** `PublicSummary()` SHALL contain `"7"` in the tension description

#### Scenario: PublicSummary omits weather when empty
- **WHEN** `state.Weather == ""`
- **THEN** `PublicSummary()` SHALL NOT contain a weather line

#### Scenario: World Rules block present when Concept.Premise is set
- **WHEN** `State.Concept.Premise == "Bears disguised as humans at a honey convention"`
- **THEN** `PublicSummary()` SHALL contain that premise string in a "World Rules" section

#### Scenario: World Rules block includes rules list when non-empty
- **WHEN** `State.Concept.Rules` contains `["Do not walk on all fours", "Never eat honey directly from a jar"]`
- **THEN** `PublicSummary()` SHALL contain both rule strings

#### Scenario: World Rules block absent when Concept is zero value
- **WHEN** `WorldConfig.Concept.Premise == ""`
- **THEN** `PublicSummary()` SHALL NOT contain the heading "World Rules"

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

### Requirement: State.Summary removed

The `State.Summary() string` method SHALL NOT exist. All callers SHALL use `PublicSummary()` and/or `LocalContext()` instead.

#### Scenario: Summary method does not exist after change
- **WHEN** `go build ./...` is run
- **THEN** there SHALL be no compilation errors referencing `State.Summary`

---

### Requirement: Zone roster computation
The system SHALL provide a `BuildZoneRoster(chars []*character.Character) map[string][]string` function in `internal/character/` that returns a map from location name to the list of character **names** (not IDs) currently at that location. Characters with an empty `Location` field SHALL be omitted.

#### Scenario: Multiple characters at the same location
- **WHEN** characters Alice and Bob both have `Location == "Market"`
- **THEN** `BuildZoneRoster` SHALL return a map entry `"Market": ["Alice", "Bob"]`

#### Scenario: Character with empty location is omitted
- **WHEN** a character has `Location == ""`
- **THEN** that character SHALL NOT appear in any roster entry

#### Scenario: Each location appears at most once as a key
- **WHEN** three characters occupy two different locations
- **THEN** the returned map SHALL have exactly two keys

---

### Requirement: Zone context prompt section
The system SHALL provide a `BuildZoneContext(roster map[string][]string) string` function in `internal/character/` that renders the roster as a human-readable block suitable for appending to LLM system prompts. The output SHALL list every location with its occupants; an empty roster SHALL return an empty string.

#### Scenario: Non-empty roster renders all zones
- **WHEN** roster has two entries: `"Park": ["Carlos"]` and `"Market": ["Alice", "Bob"]`
- **THEN** `BuildZoneContext` SHALL return a string containing "Park", "Carlos", "Market", "Alice", and "Bob"

#### Scenario: Empty roster returns empty string
- **WHEN** roster is an empty map
- **THEN** `BuildZoneContext` SHALL return `""`

---

### Requirement: Engine computes roster once per tick before message dispatch
The simulation engine SHALL compute the zone roster by calling `BuildZoneRoster` once per tick, before any `MoveDecision` or `CharChat` messages are dispatched in that tick. The same snapshot SHALL be used for all messages within the tick.

#### Scenario: Roster snapshot is consistent across a tick
- **WHEN** the engine dispatches multiple MoveDecision messages in one tick
- **THEN** all those messages SHALL be built using the same zone roster snapshot computed at the start of the tick

---

### Requirement: Movement prompt includes full zone roster
`BuildMovementPrompt` SHALL accept a `zoneRoster map[string][]string` parameter and include a "Who is where" section in the rendered prompt listing all zones and their occupants. The character's own name SHALL still appear at their current location (they are present).

#### Scenario: Movement prompt shows other characters' locations
- **WHEN** `BuildMovementPrompt` is called with a roster where location "Bar" has ["Maria", "Luis"]
- **THEN** the returned prompt string SHALL contain "Bar" and both "Maria" and "Luis"

#### Scenario: Empty roster renders no zone section
- **WHEN** `BuildMovementPrompt` is called with an empty roster
- **THEN** the returned prompt SHALL NOT contain a "Who is where" section

---

### Requirement: CharChat system prompt includes zone presence
When building `CharChatPayload.InitiatorSystem` and `CharChatPayload.ResponderSystem`, the engine SHALL append a zone-presence block (via `BuildZoneContext`) that shows who is at each location. The character's own name SHALL be excluded from the listing of their current location so the block reads as "who else is here".

#### Scenario: System prompt mentions co-located characters
- **WHEN** initiator and a third character Charlie are both at "Plaza" and the engine builds the initiator's system prompt
- **THEN** `InitiatorSystem` SHALL contain "Charlie" in the zone-presence block

#### Scenario: Character's own name is excluded from their location listing
- **WHEN** the roster at "Plaza" is ["Alice", "Charlie"] and Alice is the initiator
- **THEN** `InitiatorSystem` SHALL NOT list "Alice" under "Plaza" in the zone-presence block

---

### Requirement: Character location field

The `Character` struct SHALL expose a `Location string` runtime field indicating the character's current position, expressed as the `Name` of a `Location` from the scenario's world config. This field is NOT loaded from `characters.yaml`; it is assigned and updated at runtime by the simulation (initially by the `Scheduler`, then updated by movement decisions).

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

`RunExchange` SHALL build a distinct world context string for each character by combining `State.PublicSummary()` with `State.LocalContext(character.Location)`. If `character.Location` is empty, only `PublicSummary()` is used.

#### Scenario: Initiator and responder receive different local context
- **WHEN** initiator is at "Tavern" and responder is at "Market" and each location has local events
- **THEN** the initiator's system prompt SHALL contain the Tavern local context and NOT the Market local context, and vice versa for the responder

#### Scenario: Character with no location receives only public summary
- **WHEN** a character's `Location` field is empty
- **THEN** that character's world context SHALL consist of only `PublicSummary()` with no local context appended

---

### Requirement: Scheduler assigns initial character locations

The `Scheduler` SHALL assign a starting location to every character during initialisation. If `WorldConfig.InitialLocation` is non-empty, all characters SHALL be placed at that location. If `WorldConfig.InitialLocation` is empty, the Scheduler SHALL assign a random starting location to each character, drawing from the list of location names in `WorldConfig`. No character SHALL begin a simulation with an empty `Location`.

#### Scenario: All characters start at initial_location when set
- **WHEN** `NewEngine` is called with a `WorldConfig` where `InitialLocation == "Convention Lobby"`
- **THEN** every character in `Engine.chars` SHALL have `Location == "Convention Lobby"` before `Run` is called

#### Scenario: All characters have a location after NewEngine when initial_location is empty
- **WHEN** `NewEngine` is called with a valid scenario where `WorldConfig.InitialLocation` is empty
- **THEN** every character in `Engine.chars` SHALL have a non-empty `Location` equal to one of the world's location names

#### Scenario: Random distribution when initial_location is empty and seed is fixed
- **WHEN** a scenario has 6 characters and 6 locations, `InitialLocation` is empty, and a fixed seed is used
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

`llm.BuildMovementPrompt(c character.Character, locations []string, zoneRoster map[string][]string) string` SHALL construct a concise prompt that includes the character's name, motivation, fear, current location, the list of available locations, and a "Who is where" section from the roster. It SHALL instruct the LLM to respond with exactly one location name from the list or the word `"stay"`.

#### Scenario: Prompt contains character context
- **WHEN** `BuildMovementPrompt` is called with a character who has a non-empty `Motivation`
- **THEN** the returned string SHALL contain that motivation text

#### Scenario: Prompt lists all available locations
- **WHEN** `BuildMovementPrompt` is called with three location names
- **THEN** the returned string SHALL contain all three location names

---

## ADDED Requirements

### Requirement: WorldConcept struct en WorldConfig
El sistema SHALL definir un struct `WorldConcept` con cinco campos:
- `Premise string` (yaml: `premise`) — una sola oración describiendo la naturaleza fundamental de este mundo y la verdad oculta que los personajes deben preservar
- `Rules []string` (yaml: `rules`) — lista de restricciones que definen qué es "normal" en este mundo
- `Flavor string` (yaml: `flavor`) — string corto de tono/mood (e.g., "absurdist heist comedy")
- `CharacterSpawnRule string` (yaml: `character_spawn_rule`) — regla que describe cómo deben diseñarse los personajes creados dinámicamente; string vacío deshabilita el `spawn_character` del director
- `MaxSpawnedCharacters int` (yaml: `max_spawned_characters`) — número máximo de personajes que el director puede spawnear en runtime (`0` significa ilimitado)

`WorldConfig` SHALL exponer un campo `Concept WorldConcept` (yaml: `concept`) y un campo `InitialLocation string` (yaml: `initial_location`). Todos los sub-campos son opcionales; omitir el bloque `concept:` completo SHALL dejar `WorldConcept` en su valor cero.

#### Scenario: Bloque concept completo parseado de world.yaml
- **WHEN** un `world.yaml` contiene un bloque `concept:` con `premise`, `rules`, `flavor`, `character_spawn_rule` y `max_spawned_characters` seteados
- **THEN** los cinco campos SHALL estar poblados tras la carga

#### Scenario: Bloque concept parcial aceptado
- **WHEN** un `world.yaml` contiene `concept: { premise: "Bears disguised as humans" }` sin otros campos
- **THEN** `WorldConfig.Concept.Premise` SHALL ser `"Bears disguised as humans"`, todos los demás campos SHALL ser valor cero, y la carga SHALL retornar nil error

#### Scenario: Bloque concept ausente resulta en valor cero
- **WHEN** un `world.yaml` omite la clave `concept:` completamente
- **THEN** `WorldConfig.Concept` SHALL igualar el `WorldConcept{}` de valor cero y la carga SHALL retornar nil error

#### Scenario: Rules parseadas como lista ordenada
- **WHEN** `world.yaml` contiene `concept.rules` con tres entradas
- **THEN** `WorldConfig.Concept.Rules` SHALL tener longitud 3 con entradas en orden YAML

#### Scenario: character_spawn_rule parseado de world.yaml
- **WHEN** un `world.yaml` contiene `concept.character_spawn_rule: "All characters must be bears in human disguise"`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL ser ese string tras la carga

#### Scenario: character_spawn_rule ausente resulta en string vacío
- **WHEN** un `world.yaml` omite `character_spawn_rule` bajo `concept`
- **THEN** `WorldConfig.Concept.CharacterSpawnRule` SHALL ser string vacío y la carga SHALL retornar nil error

#### Scenario: max_spawned_characters parseado de world.yaml
- **WHEN** un `world.yaml` contiene `concept.max_spawned_characters: 3`
- **THEN** `WorldConfig.Concept.MaxSpawnedCharacters` SHALL ser `3` tras la carga

#### Scenario: max_spawned_characters ausente resulta en cero (ilimitado)
- **WHEN** un `world.yaml` omite `max_spawned_characters` bajo `concept`
- **THEN** `WorldConfig.Concept.MaxSpawnedCharacters` SHALL ser `0` y la carga SHALL retornar nil error

#### Scenario: initial_location parseado de world.yaml
- **WHEN** un `world.yaml` contiene `initial_location: "Convention Lobby"`
- **THEN** `WorldConfig.InitialLocation` SHALL ser `"Convention Lobby"` tras la carga

#### Scenario: initial_location ausente resulta en string vacío
- **WHEN** un `world.yaml` omite la clave `initial_location`
- **THEN** `WorldConfig.InitialLocation` SHALL ser string vacío y la carga SHALL retornar nil error
