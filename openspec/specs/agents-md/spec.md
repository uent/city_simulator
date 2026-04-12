## ADDED Requirements

### Requirement: AGENTS.md exists at the project root
The repository SHALL contain an `AGENTS.md` file at its root. The file SHALL be written in GitHub-flavored Markdown and SHALL be readable without rendering (plain text friendly).

#### Scenario: File present at root
- **WHEN** the repository is cloned
- **THEN** `AGENTS.md` SHALL exist at the top level alongside `go.mod`, `Makefile`, and `README.md`

---

### Requirement: Project identity section
`AGENTS.md` SHALL open with a concise description of what the project is — its purpose, LLM backend, language, config format, and output format — in no more than one short paragraph.

#### Scenario: Identity section provides essential context
- **WHEN** an agent reads the first section of `AGENTS.md`
- **THEN** it SHALL understand that the project is a Go-based LLM simulation engine using Ollama locally, configured via YAML, and producing JSONL and markdown output

---

### Requirement: Build, run, and test reference
`AGENTS.md` SHALL document the primary `make` targets (`build`, `run`, `test`, `fmt`, `vet`), the required `.env` setup step, and the key environment variables (`OLLAMA_BASE_URL`, `OLLAMA_MODEL`, `SIMULATION_NAME`, `SIMULATION_TURNS`).

#### Scenario: Agent can build and run without reading Makefile
- **WHEN** an agent reads the build/run section
- **THEN** it SHALL have sufficient information to compile, configure, and run the simulator without consulting the `Makefile` or `.env.example`

---

### Requirement: Package map section
`AGENTS.md` SHALL include a package map listing every package under `cmd/` and `internal/` with a one-line description of its responsibility.

#### Scenario: All internal packages are documented
- **WHEN** the package map is read
- **THEN** it SHALL include entries for: `character`, `director`, `llm`, `messaging`, `scenario`, `simulation`, `summary`, `world`, and `cmd/simulator`

---

### Requirement: Architecture rules section
`AGENTS.md` SHALL document the following five architecture rules explicitly:

1. All inter-actor communication goes through `MessageBus` — never direct method calls on actors
2. `CharacterActor` owns its per-pair conversation history — no other component reads or writes it
3. The Director uses XML tool-call format (`<tool_calls>[...]</tool_calls>`) — not standard function-calling APIs
4. All LLM failures are fail-open — errors are logged, the simulation continues, movement defaults to `"stay"`
5. World state is mutated exclusively by Director actions via `Execute()` — each action also appends to `state.EventLog`

#### Scenario: Architecture rules prevent common agent mistakes
- **WHEN** an agent reads the architecture rules section before making changes
- **THEN** it SHALL have enough information to avoid bypassing the MessageBus, propagating LLM errors, or mutating world state outside of registered Director actions

---

### Requirement: Extension guides
`AGENTS.md` SHALL include step-by-step instructions for the two most common extension tasks:

- **Adding a new scenario**: directory structure, required files, optional files, how to run it, reference to `CHARACTER_RULES.md`
- **Adding a new Director action**: implement the `Action` interface, register in the registry, update the prompt builder

#### Scenario: Agent can add a scenario without reading scenario loader spec
- **WHEN** an agent reads the scenario extension section
- **THEN** it SHALL know the exact files to create (`characters.yaml`, `world.yaml`, optional `scenario.yaml`), where to place them (`simulations/<name>/`), and how to run the new scenario

#### Scenario: Agent can add a Director action without reading director spec
- **WHEN** an agent reads the Director action extension section
- **THEN** it SHALL know the three required steps: implement the `Action` interface, register it in `registry.go`, and update `BuildDirectorPrompt` if the tool schema needs updating

---

### Requirement: OpenSpec workflow section
`AGENTS.md` SHALL document the OpenSpec development workflow used in this project, including the four slash commands (`/opsx:propose`, `/opsx:apply`, `/opsx:archive`, `/opsx:explore`) and where specs live (`openspec/specs/`).

#### Scenario: Agent follows OpenSpec workflow before implementing features
- **WHEN** an agent reads the OpenSpec section
- **THEN** it SHALL know to check `openspec/specs/` for existing specs, use `/opsx:propose` before implementing, and use `/opsx:apply` to work through implementation tasks

---

### Requirement: Key data flow section
`AGENTS.md` SHALL document the per-tick loop and end-of-simulation data flows using ASCII notation, covering: Director turn, character exchange (CharChat), movement decisions, JSONL logging, and summary generation.

#### Scenario: Agent understands tick execution order
- **WHEN** an agent reads the data flow section
- **THEN** it SHALL know that the Director always runs before character exchanges in each tick, and that movement decisions follow each exchange

---

### Requirement: Output format reference
`AGENTS.md` SHALL document the console output format (tick header, action lines, speech lines) and the JSONL log format (field names and types) as a quick reference.

#### Scenario: Agent generating or parsing output uses correct format
- **WHEN** an agent needs to produce or consume simulation output
- **THEN** it SHALL find the exact field names and formatting rules in `AGENTS.md` without needing to read the engine source
