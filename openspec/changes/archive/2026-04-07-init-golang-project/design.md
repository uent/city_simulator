## Context

This is a greenfield Go project. There is no existing code. The goal is to stand up a working simulator where N characters (agents) live in a shared city world and interact using LLMs running locally via Ollama. Each character has a persona, memory, and goals; the simulation drives turn-by-turn exchanges and records the narrative. The project must work offline without any cloud dependencies.

## Goals / Non-Goals

**Goals:**
- Clean, idiomatic Go project layout (`cmd/`, `internal/`, `pkg/`)
- Characters defined as data (YAML/JSON configs), not hardcoded
- All LLM calls go through a single `llm` package wrapping the Ollama REST API
- Simulation loop is deterministic given the same seed (reproducible runs)
- Output (dialogue, events, state snapshots) written to a structured log file
- CLI flags to select model, number of turns, characters file, and Ollama base URL

**Non-Goals:**
- HTTP server or web UI (pure CLI)
- Persistent database (in-memory only, optional JSON dump at end)
- Multi-model routing or cloud LLM providers
- Real-time visualization or frontend
- Concurrency/parallelism in the simulation loop (single-threaded tick loop for v1)

## Decisions

### 1. Project Layout — Standard Go layout with `internal/`

```
city_simulator/
├── cmd/
│   └── simulator/
│       └── main.go          # CLI entry point, flag parsing, wires everything
├── internal/
│   ├── character/
│   │   ├── character.go     # Character struct, loader from YAML
│   │   └── memory.go        # Rolling memory buffer per character
│   ├── llm/
│   │   ├── client.go        # Ollama HTTP client (generate + chat endpoints)
│   │   └── prompt.go        # Prompt template builder
│   ├── simulation/
│   │   ├── engine.go        # Main tick loop, interaction scheduler
│   │   └── scheduler.go     # Decides who interacts with whom each tick
│   ├── world/
│   │   └── state.go         # City state: locations, time, event log
│   └── conversation/
│       └── manager.go       # Per-pair dialogue thread, history formatting
├── configs/
│   └── characters.yaml      # Example character definitions
├── go.mod
└── go.sum
```

**Rationale:** `internal/` prevents accidental import from outside the module. Each subdirectory maps directly to one of the five capabilities in the proposal. `cmd/simulator/main.go` is the only binary; a clean `main` that just wires packages together.

**Alternative considered:** Flat package layout — rejected because it collapses concerns and makes the codebase harder to reason about as it grows.

### 2. LLM Interface — Ollama `/api/chat` (chat completions format)

Use Ollama's `/api/chat` endpoint with the messages array format rather than `/api/generate`. This lets us pass the full conversation history naturally and supports system prompts per character.

```go
type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
}
```

**Rationale:** `/api/chat` maps directly to the conversation-history pattern we need. `/api/generate` requires manual prompt stitching and has no native message roles.

**Alternative considered:** Using the `ollama` Go SDK — decided against it to keep dependencies minimal and control the exact request shape.

### 3. Character Memory — Sliding window of recent messages

Each character holds a `[]MemoryEntry` capped at a configurable `MaxMemory` count (default 20). Older entries are dropped. No vector store, no embeddings.

**Rationale:** Sufficient for short simulation runs. Simple, zero dependencies. Can be swapped for a retrieval-augmented approach later without changing the interface.

### 4. Interaction Scheduling — Round-robin with optional randomness

The scheduler maintains a queue of character pairs. Each tick, it pops the next pair, runs an exchange (A speaks → B responds → optional A reply), then pushes the pair back. A `--random-seed` flag enables shuffling for variety.

**Rationale:** Deterministic by default makes debugging easier. Randomness is opt-in.

### 5. Configuration — YAML files for characters, CLI flags for runtime params

Characters are defined in a `configs/characters.yaml` file. Runtime parameters (model name, tick count, Ollama URL) come from CLI flags with sensible defaults.

**Rationale:** Separating character data from code makes it easy to experiment with new personas without recompiling. YAML is human-readable and has good Go library support (`gopkg.in/yaml.v3`).

### 6. Output — Structured JSONL log + human-readable stdout

Each interaction event is written as a JSON line to an output file (default `simulation_output.jsonl`) AND as formatted text to stdout. This gives both machine-readable history and a readable narrative in the terminal.

## Risks / Trade-offs

- **Ollama must be running before the simulation starts** → The LLM client checks connectivity on startup and exits with a clear error message if Ollama is unreachable.
- **Context window overflow for long simulations** → Memory is capped per character; the scheduler limits history passed per prompt to the last N messages. This truncates very old context.
- **LLM response latency makes long runs slow** → Acceptable for v1; streaming output to stdout improves perceived responsiveness. No mitigation for actual throughput.
- **Character coherence degrades as memory slides out** → Known limitation of the sliding-window approach. A future retrieval layer could address this.
- **Non-determinism from the LLM** → Even with `temperature=0`, Ollama/llama.cpp can produce varying outputs. Reproducibility is best-effort, not guaranteed.

## Migration Plan

Not applicable — this is a new project. Steps to bootstrap:

1. `go mod init github.com/<user>/city_simulator`
2. Add `gopkg.in/yaml.v3` dependency
3. Create package skeletons in order: `world` → `character` → `llm` → `conversation` → `simulation` → `cmd/simulator`
4. Add `configs/characters.yaml` with 2–3 sample characters
5. Verify end-to-end with `go run ./cmd/simulator --turns 5`

## Open Questions

- Should characters be able to have different Ollama models (e.g., one character uses `llama3`, another uses `mistral`)? Currently all characters share one model from the CLI flag.
- Should the world emit "events" that characters react to, or only character↔character interactions for v1?
- What's the desired output format for the narrative log — plain text story, dialogue script, or structured JSON?
