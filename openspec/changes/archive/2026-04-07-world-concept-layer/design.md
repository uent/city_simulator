## Context

The simulation engine currently passes a `PublicSummary()` of world state to every character's system prompt, and a `LocalContext()` for their current location. Neither conveys the *rules of the world itself* — the premise that gives the scenario its dramatic tension. In Honey Heist, bears are pretending to be humans; without a "World Rules" section in the prompt, characters have no reason to worry about exposure, slip-ups, or maintaining a cover identity. The change adds two new optional YAML blocks — `concept` in `world.yaml` and `cover_identity` in `characters.yaml` — and surfaces them in existing prompt-building paths.

## Goals / Non-Goals

**Goals:**
- `WorldConcept` struct loadable from `world.yaml` with `premise`, `rules`, and `flavor` fields.
- `CoverIdentity` struct loadable from `characters.yaml` with `alias`, `role`, `backstory`, and `weaknesses`.
- `PublicSummary()` includes a "World Rules" block when `WorldConfig.Concept` is non-empty.
- Character system prompt includes a "Cover Identity" block when `Character.CoverIdentity` is non-nil.
- Honey Heist scenario updated with concrete authored content for both.

**Non-Goals:**
- Enforcing cover identity at runtime (no automatic "exposure" mechanic — the LLM narrates consequences naturally).
- A UI or admin tool to edit world concepts.
- Migration of existing scenarios (default and test-scenario remain valid with omitted fields).

## Decisions

**1. Optional structs via pointer for `CoverIdentity`, value for `WorldConcept`**

`CoverIdentity` is a pointer (`*CoverIdentity`) on `Character` — nil means "no cover, character is who they say they are." `WorldConcept` is a value struct on `WorldConfig` with all-empty fields as the zero state, which cleanly reads as "no concept defined." This avoids nil-checking in `PublicSummary()` while still giving prompt builders a simple `if c.CoverIdentity != nil` guard.

Alternatives considered: a boolean `HasConcept` flag — rejected as redundant; checking `Concept.Premise != ""` is sufficient.

**2. `PublicSummary()` renders `WorldConcept`, not `LocalContext()`**

The world concept applies to every character globally — all bears know they're in a human convention. Local context is the right place for location-specific secrets, not world-level rules. Putting rules in `PublicSummary()` means every character's system prompt always includes the premise without any extra plumbing.

**3. No new prompt builder function — inject into `BuildSystemPrompt`**

The existing `BuildSystemPrompt(c Character, worldCtx string) string` already composes the full system prompt. Rather than adding a new function, a `cover_identity` block is appended inside `BuildSystemPrompt` when `c.CoverIdentity != nil`. This keeps the call sites unchanged.

**4. `rules` as a `[]string` in `WorldConcept`**

A list of short sentences is more LLM-friendly than a paragraph — the model can scan discrete constraints. Formatting them as a bulleted list in `PublicSummary()` output is straightforward.

## Risks / Trade-offs

- **Prompt length increase** → Mitigation: `rules` should stay under 6 items; `backstory` in cover identity is a single sentence. Authors are responsible for brevity.
- **Cover identity leaks between characters** → Non-issue: `CoverIdentity` is in the per-character system prompt, not in world context shared between characters.
- **Scenarios that partially fill `WorldConcept`** → `Premise` is the only meaningful field; if empty, the block is omitted from `PublicSummary()` entirely even if `Rules` is non-empty.

## Migration Plan

No migration needed. All new fields are optional and zero-valued by default. Existing `simulations/default/` and `simulations/test-scenario/` require no changes.

Honey Heist update is additive YAML — existing keys are untouched.
