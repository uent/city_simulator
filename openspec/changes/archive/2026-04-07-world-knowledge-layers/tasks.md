## 1. Extend World Data Structs

- [x] 1.1 Add `Details string` field (yaml tag `details`) to `Location` struct in `internal/world/state.go`
- [x] 1.2 Add `Visibility string` field (yaml tag `visibility`) and `Location string` field (yaml tag `location`) to `Event` struct in `internal/world/state.go`
- [x] 1.3 Add `Location string` field (yaml tag `location`) to `Character` struct in `internal/character/character.go`

## 2. Replace State.Summary with Two-Layer Methods

- [x] 2.1 Remove `Summary() string` method from `internal/world/state.go`
- [x] 2.2 Add `PublicSummary() string` — returns time of day + location names with public descriptions + last 5 public events (visibility == "public" or "")
- [x] 2.3 Add `LocalContext(locationID string) string` — returns `Details` of the matching location + last 5 local events at that location; returns empty string (with warning log) if `locationID` is unknown or empty

## 3. Update Conversation Manager

- [x] 3.1 Update `RunExchange` in `internal/conversation/manager.go` to build per-character world context: `w.PublicSummary() + "\n" + w.LocalContext(character.Location)` for each character (initiator and responder separately)
- [x] 3.2 Remove the `w.Summary()` calls and replace with the per-character context strings

## 4. Update Default Scenario World

- [x] 4.1 Add `details` field to each location in `simulations/default/world.yaml` (private atmospheric details visible only when present)
- [x] 4.2 Add `visibility` and `location` fields to `initial_events` in `simulations/default/world.yaml` (inciting_incident = local at Market, rumor = public)
- [x] 4.3 Add `location` field to each character in `simulations/default/characters.yaml` (Elena → North Quarter, Marcus → Back Alley, Nadia → Market)

## 5. Update Honey Heist Scenario World

- [x] 5.1 Add `details` field to each location in `simulations/honey-heist/world.yaml`
- [x] 5.2 Add `visibility` and `location` fields to events in `simulations/honey-heist/world.yaml`
- [x] 5.3 Add `location` field to each character in `simulations/honey-heist/characters.yaml` (Grizwald → Convention Lobby, Honeydrop → Vault Antechamber, Claws McGee → Vendor Hall, Lady Marmalade → Convention Lobby, Patches → Alley (Exit), Dr. Snuffles → Security Office)

## 6. Verify

- [x] 6.1 Run `go build ./...` and confirm no compilation errors
- [x] 6.2 Run `--scenario default --turns 2` and verify two characters at different locations receive different world context (visible in stdout output or logs)
- [x] 6.3 Run `--scenario honey-heist --turns 1` and confirm all 6 characters load and Patches (at Alley) does not see Honeydrop's vault-local context
