## Why

The simulation currently only drives characters via their own agency and LLM prompts, but has no autonomous force to generate uncontrolled world events (weather shifts, random encounters, environmental changes, etc.). A Game Director character type with near-universal knowledge fills this gap, making simulations feel alive and reactive beyond just NPC dialogue.

## What Changes

- Introduce a new `GameDirector` character type distinct from regular characters
- Game Director has read access to full world state (all characters, locations, time, scenario context)
- Game Director autonomously generates world events on each simulation tick (or at configurable intervals)
- World events are injected into the shared world state and visible to all characters in subsequent turns
- Regular characters cannot produce world events — only the Game Director can
- Game Director's LLM prompt is structured around event generation, not conversation or personal goals
- Event types include: weather, environmental hazards, random encounters, news/rumors, time-based triggers

## Capabilities

### New Capabilities
- `game-director`: A privileged character type with universal world knowledge that autonomously generates uncontrolled world events each simulation tick

### Modified Capabilities
- `character-engine`: Must distinguish between regular characters and the Game Director, routing them through different prompt/response pipelines
- `world-state`: Must support a structured event log where the Game Director can append world events that other characters can observe
- `simulation-engine`: Must invoke the Game Director on each tick before processing regular characters, so events are visible when characters act

## Impact

- `internal/character/character.go`: New `CharacterType` enum and `GameDirector` struct/logic
- `internal/llm/prompt.go`: New prompt builder for Game Director (world event generation)
- `internal/scenario/scenario.go`: Load and initialize Game Director from scenario config
- `simulations/*/characters.yaml`: Optional Game Director entry in scenario character files
- World state model needs an `events` or `world_events` field to store generated events
