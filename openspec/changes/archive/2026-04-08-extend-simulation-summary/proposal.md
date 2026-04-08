## Why

The current summary prompt instructs the LLM to produce only "two to four paragraphs," resulting in a superficial narrative that doesn't capture the richness of the simulation. A more extensive summary would make the output feel like a proper after-action chronicle rather than a brief recap.

## What Changes

- Update the LLM system prompt to request a longer, more detailed narrative (five to eight paragraphs or more)
- Expand the character context sent to the LLM to include goals, cover identity hints, and any notable events tied to that character
- Include the scenario description in the prompt so the LLM can ground the narrative in the original setup
- Raise the event cap from 100 to 200 to give the LLM more material to work with

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `simulation-summary`: Extend prompt construction to include richer context (scenario description, higher event cap, more character fields) and instruct the LLM to produce a longer narrative

## Impact

- `internal/summary/summary.go`: `buildPrompt` function updated (system prompt wording, event cap constant, character state block)
- `openspec/specs/simulation-summary/spec.md`: requirements updated to reflect new event cap and richer prompt content
