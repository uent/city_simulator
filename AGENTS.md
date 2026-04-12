# City Simulator — Agent Guide

## What this project is

An LLM-powered character simulation engine. Characters are system prompts given deep psychological structure; the simulator pairs them up, has them converse and move around a world, and narrates the result. A Game Director LLM steers each tick via tool-call actions (weather, tension, NPC spawning, narration). LLM backend: Ollama (local, no external APIs required). Language: Go. Config: YAML. Output: JSONL per-tick log + markdown narrative summaries.

---

## Build, run, test

```bash
make build   # compile → ./city-simulator
make run     # build + run (reads .env for config)
make test    # go test -v ./...
make fmt     # go fmt ./...
make vet     # go vet ./...
```

**Setup**: copy `.env.example` → `.env` before running.

Key environment variables:

| Variable            | Default                      | Description                        |
|---------------------|------------------------------|------------------------------------|
| `OLLAMA_BASE_URL`   | `http://localhost:11434`     | Ollama server URL                  |
| `OLLAMA_MODEL`      | (required)                   | Default model name (e.g. `llama3`) |
| `SIMULATION_NAME`   | `default`                    | Scenario to load from `simulations/` |
| `SIMULATION_TURNS`  | varies                       | Number of simulation ticks         |

Override inline: `OLLAMA_MODEL=mistral make run`

---

## Package map

```
cmd/simulator/       Entry point: CLI flags, dependency wiring, startup ping
internal/
  character/         Character struct, actor, per-pair memory, zone roster
  director/          Game Director: action interface, registry, prompt builder, parser
  llm/               Ollama HTTP client (Chat/Generate), prompt builders
  messaging/         Actor interface, MessageBus, Message types and payloads
  scenario/          YAML loaders: characters, world, scenario overrides + merge
  simulation/        Engine (tick loop), Scheduler (pair selection)
  summary/           End-of-simulation narrative generation and file saving
  world/             World state, events, locations, PublicSummary/LocalContext
simulations/         Built-in scenarios (YAML files — the actual content)
openspec/            Spec-driven development artifacts (specs, changes, config)
summaries/           Generated narrative summaries (git-ignored)
```

---

## Architecture rules

**1. All inter-actor communication goes through `MessageBus`.**
Never call methods on actors directly. Use `bus.Send()` for point-to-point messages and `bus.Broadcast()` for director directives. Bypassing the bus breaks the actor lifecycle and deadlock guarantees.

**2. `CharacterActor` owns its per-pair conversation history.**
The `map[string][]Turn` keyed by peer ID is private to the actor. No other component reads or writes it. History is capped at `MaxHistory` turns and evicted oldest-first.

**3. The Game Director uses XML tool-call format — not function-calling APIs.**
The Director LLM returns a `<tool_calls>[{"name":"...","args":{...}}]</tool_calls>` block inside free text. The parser scans for those tags; malformed or missing blocks produce an empty slice, not an error. Don't replace this with OpenAI-style function calling.

**4. All LLM failures are fail-open.**
LLM errors are logged; the simulation never stops for a network or model error. Movement decisions default to `"stay"` on any failure. Never propagate LLM errors up to the engine caller.

**5. World state is mutated exclusively by Director actions.**
The engine reads `world.State`; only registered `director.Action` implementations write it via `Execute()`. Every action that mutates state MUST also append a public event to `state.EventLog` describing the change.

---

## Adding a new scenario

1. Create `simulations/<name>/` with two required files:
   - `characters.yaml` — list of character definitions
   - `world.yaml` — locations, initial events, world concept/premise
2. Optionally add `scenario.yaml` to override model, turns, or seed
3. Optionally add a `type: game_director` entry in `characters.yaml`
4. Run: `SIMULATION_NAME=<name> make run`

For character authoring rules and the YAML schema, read `simulations/CHARACTER_RULES.md`. Key points: use psychological anchors (`motivation`, `fear`, `core_belief`, `internal_tension`) — not trait lists. Include `dialogue_examples` as actual lines of speech, not descriptions.

---

## Adding a new Director action

1. **Implement the `Action` interface** in `internal/director/` (or a new file there):
   ```go
   type Action interface {
       Name() string
       Execute(args map[string]any, state *world.State, chars *[]*character.Character) error
       Summary(args map[string]any) string
   }
   ```
   `Execute` MUST append a public event to `state.EventLog`. Return a non-nil error (without mutating state) if required args are missing or invalid.

2. **Register it** in `internal/director/registry.go`.

3. **Update the `<tools>` block** in `internal/director/prompt.go` (`BuildDirectorPrompt`) so the Director LLM knows the action exists and what args it expects.

---

## OpenSpec workflow

This project uses spec-driven development. Before implementing any feature:

1. Check `openspec/specs/` for existing capability specs — they are the authoritative contracts
2. Use `/opsx:explore` to think through a problem before proposing
3. Use `/opsx:propose` to create a change with proposal, design, specs, and tasks
4. Use `/opsx:apply` to implement from the task list
5. Use `/opsx:archive` to finalize when all tasks are done

Specs use BDD format: `### Requirement:` with `#### Scenario: WHEN/THEN` blocks. Each scenario is a potential test case.

---

## Key data flows

**Per-tick loop:**
```
Director turn (if Game Director is configured)
  → BuildDirectorPrompt(state, chars, tick)
  → LLM call
  → ParseToolCalls → registry.Dispatch (each action, fail-open per action)
  → "[Director] <Summary()>" printed to stdout

Character exchange
  → Scheduler.Next() → (initiator, responder) pair
  → bus.Send(CharChat to responder actor)
      → responder actor: 2× LLM calls (initiator voice, then responder voice)
      → CharChatReply{InitiatorSpeech, InitiatorAction, ResponderSpeech, ResponderAction}
  → bus.Send(MoveDecision) × 2 → await both replies → update Character.Location

End of tick
  → JSONL log entry written to output writer
  → world state tick advanced
```

**End of simulation:**
```
summary.GenerateSummary(llmClient, state, chars, scenario)
  → LLM narrative call
  → summary.SaveSummary(scenarioName, text)
  → prints saved file path to stdout
  (fail-open: errors are logged, simulation return value is still nil)
```

---

## Output format reference

**Console — per tick:**
```
── Tick N ── InitiatorName [location] → ResponderName [location] ──
*initiator action*          ← omitted when InitiatorAction is empty
InitiatorName: speech
*responder action*          ← omitted when ResponderAction is empty
ResponderName: speech
```

**Director actions:**
```
  [Director] set_weather: storm
  [Director] move_npc: alice → market
```

**JSONL** (`simulation_output.jsonl`) — one JSON object per tick:

| Field              | Type   | Notes                          |
|--------------------|--------|--------------------------------|
| `tick`             | int    | Tick number                    |
| `initiator`        | string | Initiator character name       |
| `responder`        | string | Responder character name       |
| `initiator_speech` | string | Spoken text from initiator     |
| `initiator_action` | string | Physical action (empty if none)|
| `responder_speech` | string | Spoken text from responder     |
| `responder_action` | string | Physical action (empty if none)|
