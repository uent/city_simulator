## Context

The project has 23 capability specs, an actor-model architecture, a tool-call-based Game Director, an OpenSpec development workflow, and non-obvious patterns (fail-open LLM calls, MessageBus-only communication, XML tool-call format). None of this is documented in a single entry point. New agent sessions re-derive it from reading code and specs — inconsistently and incompletely.

`AGENTS.md` is a documentation-only change. No code is modified.

## Goals / Non-Goals

**Goals:**
- Provide a single file that gives any AI agent or new contributor an accurate mental model of the project in one read
- Document the five non-obvious architecture rules that specs don't surface explicitly
- Explain the extension points (new scenarios, new Director actions) with concrete steps
- Reference the OpenSpec workflow so agents know how to propose and implement changes
- Document key data flows and output formats as quick reference

**Non-Goals:**
- Replacing or duplicating the individual capability specs in `openspec/specs/`
- Documenting every field of every struct (that belongs in code comments or specs)
- Serving as a tutorial or user guide for running simulations
- Being exhaustive — brevity and accuracy matter more than completeness

## Decisions

### Decision: Single root-level file, not a docs/ directory

`AGENTS.md` at the project root is picked up automatically by most AI agent frameworks (Claude Code, GitHub Copilot Workspace, Cursor) without any configuration. A `docs/` directory requires explicit pointing.

**Alternative considered**: `docs/ARCHITECTURE.md` — rejected because it's not auto-loaded by agents and adds indirection for no benefit given the project's current size.

### Decision: Five architecture rules, not a full architecture doc

The rules section distills patterns that are spread across 23 specs into five actionable constraints. An agent that follows these five rules will not make the most common mistakes (bypassing MessageBus, propagating LLM errors, writing Director actions without event log entries).

**Alternative considered**: Inline links to individual specs for each rule — rejected because it makes the document dependent on spec file stability and adds noise for a quick-read document.

### Decision: Include data flow diagram in ASCII, not mermaid

The project has no rendering pipeline for mermaid. ASCII diagrams work in any terminal, any markdown renderer, and in LLM context windows without additional tooling.

### Decision: AGENTS.md, not CLAUDE.md

`AGENTS.md` is the emerging convention for multi-agent-compatible projects. `CLAUDE.md` is Claude Code-specific. Since the project may be used with other agents, `AGENTS.md` is more appropriate — Claude Code reads both.

## Risks / Trade-offs

- [Staleness] The file can drift from the actual code over time → Mitigation: Keep the file high-level (patterns and rules, not specific function signatures). Low-level details stay in specs and code.
- [Duplication] Some content overlaps with individual specs → Mitigation: The file summarizes patterns, not requirements. It points to specs for authoritative detail.
- [Scope creep] Future contributors may keep adding sections → Mitigation: The Non-Goals section explicitly excludes exhaustive documentation.
