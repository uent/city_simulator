## Context

All inter-character and director communication currently flows through two direct Go call chains:
- `conversation.Manager.RunExchange(initiator, responder, ...)` — synchronous, tightly coupled, only supports 1-to-1 exchanges initiated by the engine
- `Engine.Run` selects pairs via `Scheduler.Next()` and passes character pointers directly into Manager

The `conversation.Manager` acts as a shared mutable state bag (owns all `Thread` objects for all pairs) and an LLM dispatcher — two responsibilities that belong to separate layers. There is no mechanism for a director to address all characters simultaneously, and adding new routing patterns (e.g., directed whisper, broadcast rumor) requires changes deep inside the engine and manager.

The `game-director` change (already designed) introduces a `GameDirector` actor that must both receive messages and broadcast events. The current architecture cannot support this without bespoke engine code for each routing mode.

## Goals / Non-Goals

**Goals:**
- Every actor (Character, GameDirector) has a dedicated `chan Message` inbox that is the sole entry point for receiving communications
- Four routing modes are supported: **direct** (char→char), **directive** (director→char), **report** (char→director), **broadcast** (director→all chars simultaneously)
- Request-reply semantics: any `Message` can carry a reply channel so the sender can await a typed response in the same tick
- `conversation.Manager` is fully dissolved; per-pair conversation history moves into `CharacterActor`
- `Engine.Run` drives the tick loop by sending messages through the bus and collecting replies — no direct calls to LLM or character state
- Goroutine lifecycle is tied to `context.Context`; all actor goroutines exit cleanly on cancellation

**Non-Goals:**
- Persistent message queues or message replay across simulation restarts
- Multi-director support (deferred to a future change)
- Asynchronous fire-and-forget (all messages in v1 require a reply to preserve tick ordering)
- Runtime-configurable channel buffer sizes

## Decisions

### 1. Buffered reply channels (capacity 1) instead of unbuffered

Each `Message.ReplyChan` is a `make(chan Message, 1)`. During broadcast fan-out the director sends to N actor inboxes concurrently (N goroutines), each actor processes the message and writes exactly one reply. A capacity-1 reply channel means the actor goroutine never blocks on the write even if the director hasn't started `<-` yet. An unbuffered reply channel would deadlock when all N actors reply before the director's range loop reads them.

**Alternatives considered:**
- Unbuffered `ReplyChan` — requires the director to read replies in the exact order actors write them; deadlock-prone
- Single shared results channel — harder to correlate which reply belongs to which actor

### 2. Each actor owns its conversation history (map[string][]Turn)

`CharacterActor` holds `history map[string][]Turn` keyed by peer character ID. This replaces `Manager.threads map[string]*Thread`. Benefits: no global mutable state, actors are independently testable, history naturally partitions by actor.

**Alternatives considered:**
- Keep Manager as a shared history store accessed by actors — actors would need a reference back to the bus for history lookups, creating a circular dependency
- Central event log only (no per-pair history) — degrades LLM context quality

### 3. Actor goroutines run for the lifetime of the simulation context

Each `CharacterActor.Start(ctx context.Context)` spawns a goroutine with a `select { case msg := <-a.inbox: ... case <-ctx.Done(): return }` loop. The engine calls `Start` on all actors before `Run` and cancels the context to stop them.

**Alternatives considered:**
- On-demand goroutines per message — high goroutine churn; LLM calls inside actors become harder to cancel
- No goroutines (channels as synchronous queues) — removes the ability to do broadcast fan-out without blocking

### 4. MessageBus owns channel creation and routing; actors own processing

`MessageBus` knows all registered actors by ID. It exposes:
- `Send(msg Message)` — writes to `msg.To` actor's inbox
- `Broadcast(msg Message, excludeIDs ...string)` — spawns one goroutine per actor to write to each inbox concurrently, returns a `<-chan Message` receive-only fan-in channel that drains all N replies

The bus does NOT process messages; it only routes. Actors do NOT know the bus; they only know their inbox channel. This keeps coupling minimal.

**Alternatives considered:**
- Actors register handlers with the bus (observer pattern) — bus becomes a dispatcher, harder to reason about goroutine ownership
- Direct channel references between actors (no bus) — actors must hold references to all peers; adding/removing actors requires O(N) updates

### 5. Engine sends one CharChat message per tick pair, awaits reply

`Engine.Run` keeps the `Scheduler` for pair selection. Each tick:
1. If director exists: `bus.Broadcast(DirectorDirective{...})`, drain all replies
2. `pair := scheduler.Next()`
3. `bus.Send(CharChat{From: initiator, To: responder, ...})`
4. `reply := <-msg.ReplyChan` — receives both sides' dialogue in the reply payload
5. Log and advance tick

The engine still owns the tick cadence and output — only the dispatch mechanism changes.

**Alternatives considered:**
- Engine sends two separate messages (initiator speaks, then responder) — splits exchange logic across two round-trips; sequencing becomes engine-side bookkeeping
- Engine sends pair into a "conversation channel" consumed by a conversation goroutine — an extra goroutine layer that adds latency without benefit

### 6. conversation.Manager is deleted, not refactored

The manager's two responsibilities (thread storage + LLM dispatch) both move into `CharacterActor`. Keeping an empty shell would create confusion. Callers that previously held a `*Manager` reference are updated to use the bus.

**Alternatives considered:**
- Keep Manager as internal implementation of CharacterActor — still a shared struct; gains nothing
- Turn Manager into a pure history store — partial dissolution, leaves unclear ownership

## Risks / Trade-offs

- **Goroutine leak on panic inside actor** → Mitigation: wrap actor goroutine body in `recover`; log and exit gracefully
- **Deadlock if engine never reads ReplyChan** → Mitigation: buffered reply channels (cap 1) ensure actors never block; engine always reads replies before advancing the tick
- **Race on `Character` fields (Location, EmotionalState)** → Mitigation: only the owning `CharacterActor` goroutine mutates these fields; engine reads state via a `StateReply` message, not direct struct access
- **Broadcast adds per-tick latency proportional to N actors** → Accepted trade-off; N is small (≤20 in typical scenarios); all N LLM calls run concurrently
- **Test complexity increases** → Mitigation: bus accepts an `ActorFactory` interface; tests inject stub actors that respond deterministically without LLM calls

## Migration Plan

1. Create `internal/messaging/` package: `Message`, `MessageType`, `MessageBus`
2. Create `internal/character/actor.go`: `Actor` interface + `CharacterActor`
3. Move thread history logic from `Manager` into `CharacterActor.history`
4. Refactor `Engine`: replace `Manager` field with `*messaging.MessageBus`; update `NewEngine`, `Run`
5. Update `cmd/simulator/main.go`: instantiate `MessageBus`, register actors, remove `Manager` construction
6. Delete `internal/conversation/manager.go`
7. Verify all existing tests pass; update or replace Manager-centric tests
