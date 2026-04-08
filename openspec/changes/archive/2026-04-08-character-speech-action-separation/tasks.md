## 1. Expression Parser

- [x] 1.1 Create `internal/character/expression.go` with the `Expression` struct (`Speech`, `Action` string fields)
- [x] 1.2 Implement `ParseExpression(raw string) Expression` — extracts first `*...*` block as `Action`, rest as `Speech`; falls back to full text as `Speech` when no markers present
- [x] 1.3 Implement `FormatExpression(e Expression) string` — recombines to `*action* speech` wire format (or plain speech if `Action` is empty)

## 2. System Prompt Update

- [x] 2.1 Append expression format instruction to `BuildSystemPrompt` in `internal/character/prompt.go` — instruct the character to use `*action*` markers for physical actions, placed after the existing "Stay in character" line

## 3. Messaging Contract

- [x] 3.1 Replace `InitiatorText` and `ResponderText` fields in `CharChatReply` (`internal/messaging/message.go`) with `InitiatorSpeech`, `InitiatorAction`, `ResponderSpeech`, `ResponderAction string`

## 4. Character Actor

- [x] 4.1 In `handleCharChat` (`internal/character/actor.go`), after generating the initiator's raw LLM response, call `ParseExpression` to extract `InitiatorSpeech` and `InitiatorAction`
- [x] 4.2 Format the initiator's expression with `FormatExpression` and use it as the user-role message in the responder LLM call (so the responder sees both action and speech)
- [x] 4.3 After generating the responder's raw response, call `ParseExpression` to extract `ResponderSpeech` and `ResponderAction`
- [x] 4.4 Store raw LLM texts (including `*...*` markers) in `Turn.Text` and character memory as before (no change to storage format)
- [x] 4.5 Return `CharChatReply` with the four parsed fields instead of the two raw text fields

## 5. Engine Rendering and Log

- [x] 5.1 Update the tick print loop in `internal/simulation/engine.go` to render action lines (`*action*`) only when non-empty, followed by the speech line (`Name: speech`), for both initiator and responder
- [x] 5.2 Expand `logEntry` struct to include `initiator_speech`, `initiator_action`, `responder_speech`, `responder_action` string fields and populate them from `CharChatReply`
- [x] 5.3 Fix any compile errors caused by removing `InitiatorText` / `ResponderText` references in `engine.go`
