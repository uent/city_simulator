## 1. Zone roster helpers (internal/character/)

- [x] 1.1 Add `BuildZoneRoster(chars []*character.Character) map[string][]string` to `internal/character/prompt.go`
- [x] 1.2 Add `BuildZoneContext(roster map[string][]string) string` to `internal/character/prompt.go` (returns empty string for empty roster)
- [x] 1.3 Update `BuildMovementPrompt` signature to accept `zoneRoster map[string][]string` and append "Who is where" section when roster is non-empty

## 2. World state LocalContext update (internal/world/)

- [x] 2.1 Update `LocalContext(locationID string, presentNames []string) string` signature to accept a `presentNames []string` parameter
- [x] 2.2 Render a "Characters present" section in `LocalContext` when `presentNames` is non-empty; omit the section when nil or empty
- [x] 2.3 Update all call sites of `LocalContext` in the engine and director to pass the appropriate name list

## 3. Engine roster computation and prompt injection (internal/simulation/)

- [x] 3.1 In `engine.go`, compute zone roster via `character.BuildZoneRoster(e.chars)` once per tick before dispatching MoveDecision and CharChat messages
- [x] 3.2 Pass the roster to `BuildMovementPrompt` when constructing `MoveDecisionPayload.SystemPrompt`
- [x] 3.3 Append `character.BuildZoneContext(rosterExcludingSelf)` to `CharChatPayload.InitiatorSystem` and `CharChatPayload.ResponderSystem`, excluding each character's own name from their location's listing

## 4. Verification

- [x] 4.1 Run `go build ./...` and confirm no compilation errors
- [x] 4.2 Run a short simulation (`make run` or equivalent) and inspect output to confirm zone presence appears in movement and chat prompts
