## 1. Data Model — WorldConcept

- [x] 1.1 Add `WorldConcept` struct (`Premise string`, `Rules []string`, `Flavor string`) to `internal/world/state.go`
- [x] 1.2 Add `Concept WorldConcept` field (YAML key `concept`) to `WorldConfig` struct in `internal/world/state.go`
- [x] 1.3 Add `Concept WorldConcept` field to `State` struct and copy it from `WorldConfig` in `NewState()`

## 2. Data Model — CoverIdentity

- [x] 2.1 Add `CoverIdentity` struct (`Alias`, `Role`, `Backstory string`, `Weaknesses []string`) to `internal/character/character.go`
- [x] 2.2 Add `CoverIdentity *CoverIdentity` field (YAML key `cover_identity`) to `Character` struct

## 3. PublicSummary — World Rules block

- [x] 3.1 Update `State.PublicSummary()` in `internal/world/state.go` to append a "World Rules" section when `State.Concept.Premise != ""`
- [x] 3.2 Render `Concept.Rules` as a bulleted list within the "World Rules" section when non-empty

## 4. Character Prompt — Cover Identity block

- [x] 4.1 Update `BuildSystemPrompt` in `internal/character/prompt.go` to append a "Cover Identity" section when `c.CoverIdentity != nil`
- [x] 4.2 Include alias, role, backstory (if non-empty), and weaknesses as a bulleted list (if non-empty) in the cover identity section

## 5. Honey Heist — Scenario Update

- [x] 5.1 Add `concept:` block to `simulations/honey-heist/world.yaml` with `premise`, `rules`, and `flavor` describing the bears-in-human-world setup
- [x] 5.2 Add `cover_identity:` block to each character in `simulations/honey-heist/characters.yaml` (Grizwald, Honeydrop, Claws McGee, Lady Marmalade, Patches, Dr. Snuffles)

## 6. Validation

- [x] 6.1 Run `go build ./...` and confirm no compilation errors
- [x] 6.2 Run `go test ./...` and confirm existing tests pass
- [x] 6.3 Run a short honey-heist simulation (`make run SCENARIO=honey-heist TURNS=3`) and verify "World Rules" appears in the output and characters reference their cover identities
