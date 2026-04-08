## Why

Character LLM responses are currently a single unstructured text blob — it is impossible to distinguish what a character *says* from what they *do*. Separating speech from action enables richer narrative output, makes each turn more expressive, and ensures that when characters interact the receiving character gets a complete picture of what happened (not just words).

## What Changes

- Character prompts will instruct the LLM to use asterisk-delimited markers for physical actions: `*walks to the door* I wasn't expecting this.` — text inside `*...*` is the action, remaining text is speech.
- A parser extracts `action` and `speech` from the raw LLM response. If no asterisk markers are present, the entire response is treated as speech (safe fallback for weaker models).
- `CharChatReply` will carry separate `Speech` and `Action` fields per character (initiator and responder).
- When characters interact, the message sent to the receiving character includes both the action and speech of the sender, so the responder has full context.
- The console output and JSONL log will clearly distinguish action from speech.

## Capabilities

### New Capabilities

- `character-expression`: Defines the asterisk-based expression format, the parser contract (`ParseExpression(raw) → {Speech, Action}`), fallback behavior, and how expression fields flow through messages and logs.

### Modified Capabilities

- `character-actor`: The actor must pass both action and speech to the responder's context when building the CharChat exchange; `CharChatReply` payload gains `InitiatorSpeech`, `InitiatorAction`, `ResponderSpeech`, `ResponderAction` fields.
- `character-engine`: Engine rendering updated to display action (italics-style markers) and speech on separate lines; JSONL log entries updated to include `initiator_speech`, `initiator_action`, `responder_speech`, `responder_action`.

## Impact

- `internal/character/prompt.go` — system prompt updated to instruct asterisk format
- `internal/character/expression.go` (new) — `ParseExpression` function
- `internal/messaging/` — `CharChatReply` struct fields
- `internal/character/actor.go` — parsing and context injection
- `internal/simulation/engine.go` — rendering loop and log entry struct
