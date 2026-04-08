## REMOVED Requirements

### Requirement: Conversation thread per character pair
**Reason**: Per-pair conversation history is now owned by each `CharacterActor` in `internal/character/actor.go`. The central `Manager.threads` map is eliminated to remove shared mutable state and simplify ownership.
**Migration**: Any code accessing `Manager.AddTurn` or `Manager.threads` must be updated to send a `CharChat` message to the relevant `CharacterActor` via the `MessageBus`. History retrieval is handled internally by the actor.

### Requirement: Format conversation history for LLM prompt
**Reason**: History formatting (turn → `llm.Message` conversion) moves into `CharacterActor.buildHistory()`. The `Manager.History()` method is no longer needed as a public API.
**Migration**: Remove all call sites of `manager.History(fromID, toID, maxTurns)`. The actor handles this internally when processing a `CharChat` message.

### Requirement: Run a full exchange between two characters
**Reason**: `Manager.RunExchange` is replaced by the `CharChat` message handled by `CharacterActor`. The engine dispatches exchanges by calling `bus.Send(CharChat{...})` and awaiting the reply, rather than calling the manager directly.
**Migration**: Replace all `manager.RunExchange(ctx, initiator, responder, world, tick)` call sites with:
```go
reply := make(chan messaging.Message, 1)
bus.Send(messaging.Message{
    Type:      messaging.CharChat,
    From:      initiator.ID,
    To:        responder.ID,
    Tick:      tick,
    Payload:   worldCtxPayload,
    ReplyChan: reply,
})
result := <-reply
```
