## MODIFIED Requirements

### Requirement: CharChatReply payload fields
The system SHALL define `CharChatReply` in `internal/messaging/message.go` with the following fields:
- `InitiatorSpeech string` — what the initiator character said aloud
- `InitiatorAction string` — physical action performed by the initiator (may be empty)
- `ResponderSpeech string` — what the responder character said aloud
- `ResponderAction string` — physical action performed by the responder (may be empty)
- `Err error` — non-nil if the exchange failed

The previous `InitiatorText` and `ResponderText` fields SHALL be removed.

#### Scenario: Reply with action and speech
- **WHEN** the actor generates an initiator response with an action and speech
- **THEN** `CharChatReply.InitiatorAction` SHALL be non-empty and `CharChatReply.InitiatorSpeech` SHALL contain only the spoken text

#### Scenario: Reply with speech only
- **WHEN** the actor generates a response with no `*...*` markers
- **THEN** the corresponding `Action` field SHALL be an empty string and the `Speech` field SHALL contain the full response text

---

### Requirement: Engine renders action and speech separately
The simulation engine SHALL display character action and speech on separate lines per turn. The format SHALL be:

```
── Tick N ── InitiatorName [location] → ResponderName [location] ──
*initiator action*
InitiatorName: initiator speech
*responder action*
ResponderName: responder speech
```

Action lines are only printed when the `Action` field is non-empty. A missing action results in the action line being skipped entirely (no blank line placeholder).

#### Scenario: Tick with both action and speech for both characters
- **WHEN** `CharChatReply` has non-empty Action and Speech for both initiator and responder
- **THEN** the console SHALL print four lines: initiator action, initiator speech, responder action, responder speech

#### Scenario: Tick where one character has no action
- **WHEN** `CharChatReply.InitiatorAction` is empty but `ResponderAction` is non-empty
- **THEN** the console SHALL omit the initiator action line and still print the responder action line

#### Scenario: Tick where neither character has an action
- **WHEN** both `Action` fields are empty
- **THEN** the console SHALL print only the two speech lines (one per character), with no action lines

---

### Requirement: JSONL log entry includes action and speech fields
The `logEntry` struct written to `OutputWriter` SHALL include the following additional fields:
- `initiator_speech string` — spoken text from the initiator
- `initiator_action string` — action text from the initiator (empty string if none)
- `responder_speech string` — spoken text from the responder
- `responder_action string` — action text from the responder (empty string if none)

#### Scenario: Log entry with action present
- **WHEN** `CharChatReply.InitiatorAction` is `slams the table`
- **THEN** the JSONL line for that tick SHALL contain `"initiator_action":"slams the table"`

#### Scenario: Log entry with no action
- **WHEN** `CharChatReply.ResponderAction` is empty
- **THEN** the JSONL line SHALL contain `"responder_action":""` (empty string, not omitted)
