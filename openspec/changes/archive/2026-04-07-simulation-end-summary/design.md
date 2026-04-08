## Context

The simulation engine (`internal/simulation/engine.go`) runs a tick loop and currently writes only JSONL log entries per tick to an `OutputWriter`. When the loop finishes there is no narrative record of what happened. Users wanting to review or compare runs have no persistent artifact beyond stdout output that scrolls past.

The world state (`internal/world/state.go`) accumulates events and character states throughout the run. The conversation manager holds per-exchange dialogue. Both are available inside `Engine.Run` at loop end — making the end of `Run()` the natural injection point.

## Goals / Non-Goals

**Goals:**
- Generate a human-readable narrative summary at the end of each simulation run.
- Persist the summary to a timestamped file so multiple runs never overwrite each other.
- Use the existing LLM client to produce the narrative (no new dependency).
- Place summary files alongside the scenario data (e.g., `simulations/<scenario>/summary-<timestamp>.md`).

**Non-Goals:**
- Real-time or per-tick summaries.
- Structured/machine-readable summary format (markdown prose is enough).
- Retroactive summarization of past runs from JSONL logs.
- Configurable summary format or language (English prose is the default).

## Decisions

### 1. Injection point: end of `Engine.Run()`

After the tick loop the engine has full access to world state and character list. Appending a `generateSummary(ctx)` call at the end of `Run()` is the least invasive change and keeps all summary logic encapsulated.

**Alternative considered**: a post-run hook passed in via `Config`. Rejected — adds API surface for a single use case.

### 2. Output path: `simulations/<scenario-name>/summary-<RFC3339-timestamp>.md`

The scenario name is already available in `cfg.Scenario`. Using `time.Now().Format(time.RFC3339)` (with `:` replaced by `-` for Windows filesystem compatibility) guarantees uniqueness per run without needing a counter or UUID dependency.

**Alternative considered**: output next to the JSONL log using the same base path. Rejected — the log path is an `io.Writer` with no filename, so the scenario directory is the only reliable anchor.

### 3. Summary prompt: world events + character final states

The prompt passes:
- All world events (from `world.State.Events`)
- Each character's name, role, and final emotional state
- The scenario description and total tick count

This gives the LLM enough narrative material without including raw dialogue (which would bloat the prompt and risk token-limit issues on long runs).

**Alternative considered**: include full conversation transcripts. Rejected — too large for typical runs; events + emotional arc are sufficient for a coherent summary.

### 4. New package: `internal/summary`

A dedicated `summary` package (`GenerateSummary(ctx, client, world, chars, scenario) (string, error)` + `SaveSummary(scenarioName, content string) (string, error)`) keeps the engine free of file I/O and prompt-building concerns.

**Alternative considered**: add directly to `engine.go`. Rejected — mixes concerns and makes the logic harder to test independently.

## Risks / Trade-offs

- **LLM latency at shutdown** → The summary call adds one extra LLM round-trip at the end of every run. Mitigation: log and skip if it fails (fail-open, consistent with director behavior).
- **Token limits on long runs** → Passing all events verbatim may exceed context on very long simulations. Mitigation: truncate event list to last N events (configurable constant, default 100).
- **Filesystem permissions on Windows** → Colons in RFC3339 timestamps are illegal on Windows paths. Mitigation: replace `:` with `-` when formatting the filename.
