## Why

The initial design hardcodes a single `configs/characters.yaml` and CLI flags for world parameters. This makes it cumbersome to switch between different simulations (a medieval town, a space station, a detective noir city, etc.) — the user would need to overwrite files or remember long flag combinations. We need a first-class concept of a **scenario**: a self-contained directory that bundles everything needed to run one simulation.

## What Changes

- Introduce a `simulations/` root directory where each subdirectory is a named scenario
- Each scenario folder replaces the global `configs/characters.yaml` with its own `characters.yaml`, `world.yaml`, and optional `scenario.yaml` override file
- The CLI gains a `--scenario` flag (accepts a name under `simulations/` or an absolute path) replacing the old `--characters` flag
- Add a `scenario-loader` package (`internal/scenario`) responsible for loading and validating all three config files from a scenario directory
- Ship one built-in example scenario: `simulations/default/`
- **BREAKING**: `--characters` CLI flag removed in favor of `--scenario`

## Capabilities

### New Capabilities

- `scenario-loader`: Reads a scenario directory, validates and assembles a `Scenario` struct (characters + world config + runtime overrides) for injection into the engine

### Modified Capabilities

- `world-state`: `NewState` now accepts a `WorldConfig` loaded from `world.yaml` instead of a raw `[]Location` slice — requires updating how the engine initializes world state
- `simulation-engine`: `Config` struct replaces the `Characters []character.Character` field with a `Scenario scenario.Scenario` field; the engine delegates character and world loading to the scenario loader

## Impact

- New `simulations/` directory at project root (tracked in git, serves as scenario library)
- `internal/scenario/` package added
- `internal/world/state.go` constructor signature changes
- `cmd/simulator/main.go` flag wiring updated
- Old `configs/` directory removed (replaced by `simulations/default/`)
