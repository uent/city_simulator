## Why

There is no project yet. We need to bootstrap a Go-based city simulator where autonomous characters interact with each other using locally-hosted LLMs via Ollama — enabling experimentation with emergent social dynamics, narrative generation, and agent behavior without requiring any cloud API or external service.

## What Changes

- Create the entire Go project from scratch (module, directory layout, entry point)
- Introduce a **character engine** that defines agent personas, memory, and goals
- Introduce an **LLM client** that talks to a running Ollama instance to generate character thoughts and dialogue
- Introduce a **simulation engine** that drives turn-based or tick-based interactions between characters
- Introduce a **world state** module that tracks locations, time, and shared context
- Introduce a **conversation manager** that routes dialogue exchanges and stores history
- Provide a CLI entry point to configure and launch a simulation run

## Capabilities

### New Capabilities

- `character-engine`: Defines the Character type with persona, memory, goals, emotional state, and current context; handles character initialization from config files
- `llm-client`: Wraps Ollama HTTP API (`/api/generate`, `/api/chat`), supporting model selection, prompt templating, streaming, and retry logic
- `simulation-engine`: Orchestrates the simulation loop — selects which characters interact each tick, routes prompts to the LLM client, advances world time, and decides when to stop
- `world-state`: Tracks the city's locations, time-of-day, active events, and the shared narrative log visible to all characters
- `conversation-manager`: Manages multi-turn dialogue threads between two or more characters, maintains message history per pair, and formats context windows for LLM prompts

### Modified Capabilities

<!-- No existing capabilities — this is a greenfield project -->

## Impact

- Creates `go.mod` and all source files under a new package structure
- Runtime dependency: Ollama running locally (no network calls to external services)
- No database — state is in-memory and optionally persisted to JSON/YAML files at the end of a run
- No HTTP server exposed — pure CLI/REPL tool
