## Why

The simulator ships with no built-in scenarios beyond `simulations/default/`. A concrete, flavourful example scenario demonstrates the scenario system in action and gives contributors a reference for writing their own — while also being fun to run. Honey Heist is a comedic bear-heist premise that is self-contained, tonally distinct from the default city, and small enough to implement in one pass.

## What Changes

- Add `simulations/honey-heist/` with all three required files: `characters.yaml`, `world.yaml`, and `scenario.yaml`
- `characters.yaml` defines a cast of bears (each with a criminal specialty) attempting to steal a legendary pot of honey from a heavily guarded convention centre
- `world.yaml` defines the locations of the HoneyCon convention centre and its surroundings (lobby, vault, roof, alley, etc.)
- `scenario.yaml` sets a shorter default run (20 turns) and a fun system-prompt preamble that primes the LLM with the heist tone

## Capabilities

### New Capabilities

- `honey-heist-scenario`: The `simulations/honey-heist/` directory bundle — characters, world, and runtime overrides — that can be launched with `--scenario honey-heist`

### Modified Capabilities

<!-- No existing spec-level requirements change; this is purely additive data. -->

## Impact

- New `simulations/honey-heist/` directory (3 YAML files, no code changes)
- No Go source changes required; relies entirely on the scenario-loader infrastructure being in place
- Depends on the `simulation-scenarios` change being implemented first
