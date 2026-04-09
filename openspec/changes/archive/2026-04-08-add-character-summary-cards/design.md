## Context

The simulation produces a narrative summary via `internal/summary/summary.go`. `GenerateSummary` sends world events and character states to the LLM and returns a prose narrative. That text is then printed to the terminal and saved to a `.md` file via `SaveSummary`.

Currently, all character data is only used as LLM input context — it is never surfaced directly to the reader. After the simulation, there is no structured record of who each character was.

## Goals / Non-Goals

**Goals:**
- Append a "Character Cards" section after the LLM narrative in the final summary string.
- Each card renders a non-director character's attributes in readable Markdown.
- Empty/nil fields are omitted. Game Director characters are excluded.

**Non-Goals:**
- No changes to the LLM prompt or the narrative text itself.
- No new persistence format — the cards are part of the same `.md` file.
- No interactive or dynamic card rendering.

## Decisions

### Decision: Render cards in `internal/summary`, not at the call site

The `GenerateSummary` function already receives `[]*character.Character`. Adding a `renderCharacterCards` helper in the same package keeps the feature self-contained and requires no changes to callers.

**Alternative considered**: Render cards in the engine/main and concatenate after `GenerateSummary`. Rejected — it scatters summary-building logic across layers and forces every caller to know about cards.

### Decision: Append cards to the returned string, not write a separate file

The existing `SaveSummary` writes whatever string it receives. Embedding the cards in the narrative string means both the terminal output and the saved file include the cards with zero extra coordination.

**Alternative considered**: Separate `SaveCharacterCards` file. Rejected — unnecessary file proliferation and harder to read as a self-contained document.

### Decision: Markdown format with `---` separator and `###` headers per character

Matches the existing `.md` output style. Uses a horizontal rule to visually separate the narrative from the cards section, then one `###` heading per character name.

## Risks / Trade-offs

- [Long card sections for large casts] → Cards are data-only (no LLM call), so length scales linearly and predictably. Omitting empty fields keeps cards compact.
- [Game Director confusion] → Excluded by checking `c.Type == "game_director"` so it never appears as a "character" card.
