## 1. Prerequisites

- [x] 1.1 Verify that the `simulation-scenarios` change has been fully implemented (scenario loader, `--scenario` flag, and `simulations/` directory convention are in place)
- [x] 1.2 Confirm that `scenario.yaml` supports a `turns` override field by checking `openspec/changes/simulation-scenarios/specs/scenario-loader/spec.md`

## 2. Create scenario directory

- [x] 2.1 Create the directory `simulations/honey-heist/`

## 3. Author characters.yaml

- [x] 3.1 Write `simulations/honey-heist/characters.yaml` with 6 bear characters: Grizwald (Mastermind), Honeydrop (Safe-cracker), Claws McGee (Muscle), Lady Marmalade (Grifter), Patches (Wheelman), Dr. Snuffles (Tech Bear)
- [x] 3.2 Ensure each character entry has `name`, `role`, and `personality` fields with non-empty values

## 4. Author world.yaml

- [x] 4.1 Write `simulations/honey-heist/world.yaml` with 6 locations: Convention Lobby, Vendor Hall, Security Office, Vault Antechamber, Vault, Alley (Exit)
- [x] 4.2 Ensure each location has a non-empty `name` and `description`
- [x] 4.3 Add at least one `initial_event` describing the scene (e.g., guards changed shift, HoneyCon at peak attendance)

## 5. Author scenario.yaml

- [x] 5.1 Write `simulations/honey-heist/scenario.yaml` with `turns: 20`
- [x] 5.2 If `system_prompt_prefix` is supported by the scenario loader, add a 1–3 sentence heist-tone prefix; otherwise omit the field

## 6. Verify

- [x] 6.1 Run `go run ./cmd/simulator --scenario honey-heist` and confirm the simulation starts without errors
- [x] 6.2 Confirm all 6 characters appear in the first tick's output
- [x] 6.3 Confirm the simulation stops at turn 20 by default
- [x] 6.4 Run `go run ./cmd/simulator --scenario honey-heist --turns 5` and confirm it stops at turn 5
