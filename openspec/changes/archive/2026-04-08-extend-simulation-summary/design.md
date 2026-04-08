## Context

The summary is generated in `internal/summary/summary.go` by `buildPrompt`, which feeds a system prompt and a user message to the LLM. Currently the system prompt asks for "two to four paragraphs" and the user message contains only the scenario name, tick count, a capped event list (max 100), and each character's name, occupation, location, and emotional state.

The result is a brief, surface-level recap. The LLM has access to far richer data — world concept, atmosphere, character psychology — but none of it is surfaced.

## Goals / Non-Goals

**Goals:**
- Instruct the LLM to produce a longer, richer narrative (target: 6+ paragraphs)
- Surface world concept (premise, flavor, rules) and atmosphere in the prompt
- Include character psychological fields (motivation, fear, goals) in the character block
- Raise the event cap from 100 to 200

**Non-Goals:**
- Changing function signatures in `summary.go`
- Adding new LLM calls or multi-step summarization
- Changing how the summary is saved or displayed

## Decisions

**1. Extend system prompt target length**

Change "two to four paragraphs" → "a rich narrative of at least six paragraphs". More explicit length guidance directly controls LLM output length without any code architecture change.

**2. Include `WorldConcept` and atmosphere in user message**

`sc.World.Concept` (Premise, Flavor, Rules) and `sc.World.Atmosphere`/`sc.World.Weather` are already on the `Scenario` argument. Adding them to the user message gives the LLM the genre/tone anchor it needs to write a coherent chronicle. No new parameters needed.

**3. Richer character block**

The current block only shows name, occupation, location, and emotional state. Adding `Motivation`, `Fear`, and `Goals` (already on `Character`) lets the LLM reflect on whether characters achieved their objectives, which is the most narratively interesting question.

**4. Raise event cap to 200**

`maxEvents` is a package-level constant — one-line change. 200 events still fits comfortably within typical LLM context windows.

## Risks / Trade-offs

- **Higher token usage** — longer prompt + longer response increases API cost. Acceptable given this is a one-off end-of-simulation call. → No mitigation needed.
- **LLM may ignore length guidance** — "at least six paragraphs" is a soft constraint. → Acceptable; the richer context will naturally produce more output regardless.

## Migration Plan

All changes are in `internal/summary/summary.go` (`buildPrompt` function body and `maxEvents` constant). No API, config, or schema changes. No migration steps required.
