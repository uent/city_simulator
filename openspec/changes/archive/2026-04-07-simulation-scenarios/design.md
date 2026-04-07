## Context

The initial `init-golang-project` change loads characters from a single `configs/characters.yaml` and world parameters from CLI flags. This works for one fixed simulation but makes experimenting with different scenarios (different worlds, casts, tones) tedious. We need a scenario abstraction that bundles all per-simulation data into a directory so switching scenarios is a one-flag change.

## Goals / Non-Goals

**Goals:**
- One `--scenario <name>` flag selects an entire simulation (characters + world + optional overrides)
- Scenario directories are self-contained and portable (copy a folder → new scenario)
- A new `internal/scenario` package owns all loading/validation logic; no other package reads YAML directly
- Existing engine and conversation manager stay unchanged except where `Config` struct fields are updated
- Ship `simulations/default/` as a working example out of the box

**Non-Goals:**
- Hot-reloading scenarios mid-run
- Remote or URL-based scenario sources
- A scenario registry or discovery mechanism beyond the `simulations/` directory convention
- Versioning or migration of scenario file formats

## Decisions

### 1. Scenario directory layout

```
simulations/
└── <scenario-name>/
    ├── characters.yaml      # required — character definitions (same format as before)
    ├── world.yaml           # required — locations and optional initial events
    └── scenario.yaml        # optional — runtime overrides (model, turns, seed, etc.)
```

**Rationale:** Three files keeps each concern separate and makes partial overrides natural — you can swap just `world.yaml` without touching characters. A single fat config file would conflate character data with world data and make diffs noisy.

**Alternative considered:** Nested keys in one `simulation.yaml` — rejected because merging YAML maps for overrides is fiddly and the three-file split is self-documenting.

### 2. `WorldConfig` struct loaded from `world.yaml`

```yaml
# world.yaml
locations:
  - name: Town Square
    description: The bustling center of the city
  - name: Tavern
    description: A dimly lit gathering place
initial_events:
  - description: "A mysterious stranger arrived at dawn"
    event_type: "arrival"
```

`world.State` is initialized from a `WorldConfig` instead of a raw `[]Location`. This is a small constructor signature change but makes the world data fully file-driven.

**Rationale:** World parameters logically belong in the scenario, not in CLI flags. The initial events list lets scenario authors seed narrative context before the first tick.

### 3. `scenario.yaml` as optional runtime overrides

```yaml
# scenario.yaml (all fields optional)
model: mistral
turns: 30
seed: 42
output: my_run.jsonl
```

When present, fields in `scenario.yaml` override the corresponding CLI flags. CLI flags override `scenario.yaml`. Priority: CLI > scenario.yaml > compiled defaults.

**Rationale:** Allows scenarios to declare their intended model and run length without forcing users to remember flags. CLI flags still win so users can experiment freely.

### 4. `internal/scenario` package as the single loading boundary

```go
type Scenario struct {
    Name       string
    Dir        string
    Characters []character.Character
    World      world.WorldConfig
    Overrides  RuntimeOverrides  // from scenario.yaml, all fields are pointers (nil = not set)
}

func Load(dirOrName string) (Scenario, error)
```

`Load` resolves the path (if not absolute, looks under `simulations/<dirOrName>`), reads and validates all three files, and returns a fully populated `Scenario`. The engine and `main.go` only ever touch `Scenario` — they never read YAML themselves.

**Rationale:** Centralizing IO here means tests for character loading, world loading, and override merging are all in one place. The engine remains pure logic.

### 5. `simulations/` tracked in git, `configs/` removed

The old `configs/characters.yaml` is superseded by `simulations/default/characters.yaml`. The `configs/` directory is deleted.

**Rationale:** Scenarios are first-class project assets, not throwaway config. Committing `simulations/` makes it trivial to share scenarios and diff changes to world/characters.

## Risks / Trade-offs

- **Breaking change to `--characters` flag** → Mitigated by clear error message on startup: "The --characters flag has been removed. Use --scenario instead."
- **`world-state` constructor signature change** → Only called from `main.go` and `simulation/engine.go`; both are updated as part of this change.
- **Users who store characters outside `simulations/`** → `Load` accepts an absolute path as fallback, so custom paths still work with `--scenario /path/to/dir`.

## Migration Plan

1. Create `simulations/default/` from existing `configs/characters.yaml`
2. Add `simulations/default/world.yaml` with default city locations
3. Implement `internal/scenario` package
4. Update `internal/world.NewState` to accept `WorldConfig`
5. Update `internal/simulation.Config` to embed `scenario.Scenario`
6. Update `cmd/simulator/main.go`: replace `--characters` with `--scenario`, add override merge logic
7. Delete `configs/` directory

## Open Questions

- Should `scenario.yaml` support a `description` field displayed at simulation start? Nice for documentation but not blocking.
- Should the CLI have a `list-scenarios` subcommand that prints available names from `simulations/`? Could be a v2 addition.
