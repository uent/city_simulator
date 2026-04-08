## 1. Message Bus Package

- [x] 1.1 Create `internal/messaging/` package with `MessageType` constants: `CharChat`, `DirectorDirective`, `CharReport`, `Broadcast`, `MoveDecision`
- [x] 1.2 Define `Message` struct with fields: `Type`, `From`, `To`, `Tick`, `Payload any`, `ReplyChan chan Message`; add `NewRequest(...)` and `NewReply(...)` constructor helpers that create `ReplyChan` with buffer capacity 1
- [x] 1.3 Define `Actor` interface: `ID() string`, `Inbox() chan Message`, `Start(ctx context.Context)`, `Stop()`
- [x] 1.4 Implement `MessageBus` struct with `Register(id string, inbox chan Message)`, `Send(msg Message) error`, `Broadcast(msg Message, excludeIDs ...string) (<-chan Message, error)`, and `StartAll(ctx context.Context)` (calls Start on all registered Actor implementations stored alongside their inboxes)
- [x] 1.5 Implement `Broadcast` fan-out: spawn one goroutine per recipient, each writes `msg` to inbox; collect all replies on a buffered fan-in channel that is closed after N replies; return the fan-in channel

## 2. Character Actor

- [x] 2.1 Create `internal/character/actor.go` with `CharacterActor` struct holding: `*Character`, `*llm.Client`, `inbox chan Message`, `history map[string][]Turn`, `maxHistory int` (default 20), `once sync.Once`, `cancel context.CancelFunc`
- [x] 2.2 Implement `CharacterActor.Start(ctx context.Context)`: use `sync.Once` to spawn the processing goroutine; goroutine runs `select { case msg := <-a.inbox; case <-ctx.Done() }` loop
- [x] 2.3 Implement `CharacterActor.process(msg Message)`: dispatch on `msg.Type` — `CharChat` → `handleCharChat`, `MoveDecision` → `handleMoveDecision`; default → log unhandled type
- [x] 2.4 Implement `handleCharChat`: build system prompt from character persona + world context in payload; retrieve last ≤10 turns from `history[msg.From]`; call `llmClient.Chat`; append exchange to `history[msg.From]`; update `Character.Memory`; evict history beyond `maxHistory`; write reply to `msg.ReplyChan` (text in payload, nil error on success)
- [x] 2.5 Implement `handleMoveDecision`: build movement prompt from character + locations list in payload; call `llmClient.Chat`; parse location (exact → case-insensitive → substring → "stay"); write reply with chosen location to `msg.ReplyChan`
- [x] 2.6 Move `Turn` type definition from `internal/conversation/` to `internal/character/actor.go`; update any existing references

## 3. Simulation Engine Refactor

- [x] 3.1 Replace `Manager *conversation.Manager` field in `simulation.Config` with `Bus *messaging.MessageBus`
- [x] 3.2 Update `NewEngine`: remove Manager validation; verify that the bus has at least 2 character actors registered (excluding any director actor)
- [x] 3.3 In `Engine.Run`: call `e.cfg.Bus.StartAll(ctx)` before the tick loop; defer stop of all actors via context cancellation
- [x] 3.4 In the tick loop: if director registered, call `e.cfg.Bus.Broadcast(DirectorDirective{...})` and drain the fan-in channel before `scheduler.Next()`
- [x] 3.5 Replace `manager.RunExchange(...)` call with: build `CharChat` message with world context payload, `bus.Send`, read reply from `msg.ReplyChan`; extract dialogue texts from reply payload
- [x] 3.6 Replace `manager.DecideMovement(...)` calls with: build `MoveDecision` messages for both pair characters, send via bus, read replies; apply location changes from reply payloads
- [x] 3.7 Remove `Manager` import from `engine.go`; update `logEntry` struct if needed

## 4. Wiring (main.go)

- [x] 4.1 In `cmd/simulator/main.go`: instantiate `messaging.NewMessageBus()`
- [x] 4.2 For each character in `scenario.Characters`: create `character.NewCharacterActor(c, llmClient)`, call `bus.Register(actor.ID(), actor.Inbox())`, store actors in a slice
- [x] 4.3 Pass `Bus: bus` in `simulation.Config` instead of `Manager: manager`
- [x] 4.4 Remove `conversation.NewManager(...)` instantiation

## 5. Cleanup

- [x] 5.1 Delete `internal/conversation/manager.go`
- [x] 5.2 Remove the `internal/conversation/` package directory if no other files remain; update any import in tests
- [x] 5.3 Verify `go build ./...` passes with no errors after cleanup
