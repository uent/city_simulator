## Context

Character LLM responses are currently a single `string` (`InitiatorText` / `ResponderText` in `CharChatReply`). The engine prints them verbatim and the JSONL log records them as opaque text. There is no way to distinguish spoken dialogue from physical actions.

The goal is to add an **expression layer** that parses `*action*` markers out of raw LLM responses, making speech and action independently usable by the engine, the log, and the context passed between characters.

This change touches four files: a new parser (`expression.go`), the messaging contract (`message.go`), the actor (`actor.go`), the engine (`engine.go`), and the character system prompt (`prompt.go`).

## Goals / Non-Goals

**Goals:**
- Parse `*...*` markers from raw LLM text to extract `Action` and `Speech` fields.
- Safe fallback: if no markers are present, the entire text becomes `Speech`; `Action` is empty.
- Pass both fields to the other character as context during an exchange.
- Render action and speech on separate, clearly labelled lines in the console.
- Include `action` and `speech` as distinct fields in the JSONL log.

**Non-Goals:**
- Changing the LLM client or adding structured output / JSON mode — the format stays plain text.
- Modifying movement decision prompts or the director prompt.
- Persisting expression data to YAML or summaries.
- Supporting multiple action blocks per response (only the first `*...*` block is extracted as action; any remaining markers become part of speech).

## Decisions

### D1 — Asterisk convention over JSON or prefix tags

**Decision**: Use `*action text*` markers embedded in plain text.

**Rationale**: LLMs in roleplay contexts produce this format naturally without instruction — it is a learned convention from training data. JSON output is fragile on weaker models (extra prose before/after the JSON block breaks parsers). Prefix tags (`ACTION:`) require exact keyword matching and are inconsistently capitalized. Asterisks degrade gracefully: a model that ignores the instruction produces plain text, which the fallback handles correctly.

**Alternatives considered**: JSON with a fallback re-parse pass (too brittle), two separate LLM calls for action and speech (doubles latency and cost), prefix-based format (fragile keyword matching).

---

### D2 — Parser extracts only the first `*...*` block as Action

**Decision**: `ParseExpression` finds the first `*...*` substring, assigns it to `Action`, and everything else (before + after, trimmed) becomes `Speech`.

**Rationale**: Characters typically perform one dominant action per turn. Extracting multiple blocks would require a more complex rendering strategy and complicates what gets passed to the other character. Keeping it to one block keeps the output clean.

---

### D3 — Raw full text stored in history and memory; parsed fields used only for output and context injection

**Decision**: `Turn.Text` continues to store the raw LLM response (including `*...*` markers). The parsed `Speech` and `Action` fields are used for: (a) rendering, (b) the JSONL log, and (c) the formatted context injected into the next character's prompt.

**Rationale**: History is fed back to the LLM as prior turns. Keeping raw text (with markers) preserves the LLM's own format and allows it to maintain consistent style across turns. Stripping markers from history would lose the behavioral signal.

---

### D4 — Context injection format when passing initiator's expression to responder

**Decision**: The user-role message sent to the responder is formatted as:
- If action is non-empty: `*{action}* {speech}` (combined on one line).
- If action is empty: just `{speech}`.

**Rationale**: This mirrors the format the LLM itself produces, so the responder naturally interprets it as a complete turn and can react to both the action and speech. Sending only speech would lose behavioral context.

---

### D5 — `CharChatReply` replaces monolithic text fields with four named fields

**Decision**: Replace `InitiatorText string` and `ResponderText string` with `InitiatorSpeech`, `InitiatorAction`, `ResponderSpeech`, `ResponderAction string`.

**Rationale**: The engine is the consumer of `CharChatReply` and needs both fields separately to render and log them. Returning a pre-parsed struct avoids parsing the same text twice and makes the contract explicit.

## Risks / Trade-offs

- **[Risk] Weak models ignore the `*action*` instruction** → Mitigation: fallback treats full response as speech; simulation continues normally, just without action output.
- **[Risk] Model produces nested asterisks (`*he said *yes* firmly*`)** → Mitigation: parser uses the first `*` open and first `*` close after it (greedy-safe approach). Nested markers are unlikely in practice.
- **[Risk] `CharChatReply` field rename is a breaking change to the engine** → Mitigation: both fields are only used in `engine.go` (`result.InitiatorText`, `result.ResponderText`), so the blast radius is one file.

## Migration Plan

1. Add `internal/character/expression.go` with `ParseExpression`.
2. Update `internal/messaging/message.go`: replace `CharChatReply` text fields.
3. Update `internal/character/prompt.go`: append expression format instruction to `BuildSystemPrompt`.
4. Update `internal/character/actor.go`: parse both LLM responses; format context for responder.
5. Update `internal/simulation/engine.go`: render action/speech separately; expand `logEntry` struct.
6. No migration needed — this is a single-binary simulation with no persistent state to migrate.
