## Why

AI agents working on this codebase have no single document that explains its architecture, conventions, and non-obvious patterns — so each session re-derives them from scratch, often incompletely. An `AGENTS.md` at the repo root gives any agent (or new contributor) an immediate, accurate mental model of the system and prevents common mistakes like bypassing the MessageBus, propagating LLM errors, or ignoring the OpenSpec workflow.

## What Changes

- Add `AGENTS.md` at the project root
- The file covers: what the project is, how to build/run/test, package map, five core architecture rules, how to extend the system (new scenarios, new Director actions), the OpenSpec workflow, key data flows, and output format reference

## Capabilities

### New Capabilities

- `agents-md`: A root-level `AGENTS.md` document that serves as the authoritative onboarding guide for AI agents and new contributors working on this codebase

### Modified Capabilities

_(none — no existing specs change)_

## Impact

- One new file: `AGENTS.md` at the project root
- No code changes, no dependency changes, no breaking changes
- Improves agent consistency and reduces per-session context-building overhead
