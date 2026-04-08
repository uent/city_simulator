## Why

The simulator needs a new scenario that demonstrates its dramatic range beyond comedy/heist: a single-character survival-horror in the DOOM universe, where the pre-loaded protagonist (Doom Guy) must navigate a multi-zone Hell campaign to stop the Demon Prince before he activates the Hellstone Convergence to terraform Earth.

## What Changes

- New simulation scenario: `simulations/doom-hell-crusade/` with `world.yaml`, `characters.yaml`, and `scenario.yaml`
- Doom Guy is the only pre-loaded character; all other characters are spawned dynamically by the director
- The director spawns characters from a curated archetype pool: demon renegades, trapped human souls, the Prince's heralds
- The world spans 5 Hell locations across three acts (infiltration → ascent → confrontation)
- The scenario runs for 30 turns by default to allow narrative arc to develop

## Capabilities

### New Capabilities

- `doom-hell-crusade-scenario`: World layout, character roster, director config, and scenario.yaml for the DOOM hell crusade simulation

### Modified Capabilities

<!-- No existing spec-level requirements change — new scenario files only -->

## Impact

- New files under `simulations/doom-hell-crusade/`
- No changes to existing Go code or specs
- Demonstrates director's `spawn_character` action in a single-protagonist scenario with dynamic ensemble
