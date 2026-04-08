## Why

Communication between characters and the Game Director currently happens through direct Go function calls (`Manager.RunExchange`, direct struct field access), tightly coupling the engine to specific caller/callee pairs and making it impossible to add new communication patterns (director→all broadcast, asynchronous events) without threading more parameters through the call stack. Replacing this with a Go-channel message bus gives every actor a typed inbox, decouples senders from receivers, and makes the communication topology explicit and extensible.

## What Changes

- Introduce a `MessageBus` that routes typed `Message` values between named actors via buffered Go channels
- Each `Character` and the `GameDirector` is wrapped in an `Actor` that owns a goroutine and an `Inbox chan Message`
- Support four routing modes: **character→character** (direct), **director→character** (directive), **character→director** (report), **director→all** (broadcast fan-out)
- Request-reply semantics: every `Message` carries a `ReplyChan chan Message` so the sender can `<-msg.ReplyChan` to await a response within the same tick
- **BREAKING**: `conversation.Manager` is dissolved; conversation thread history and LLM dispatch move into each `CharacterActor`
- **BREAKING**: `Engine.Run` no longer calls `Manager.RunExchange` directly; it drives the tick by sending messages through the bus and collecting replies

## Capabilities

### New Capabilities
- `message-bus`: Typed Go-channel message bus with per-actor inbox channels, four routing modes, and request-reply semantics
- `character-actor`: Actor wrapper around `Character` that owns a goroutine, processes inbox messages, maintains per-pair conversation history, and generates LLM responses

### Modified Capabilities
- `simulation-engine`: Engine now sends `Message` values through the bus and collects replies instead of calling `Manager.RunExchange`; tick orchestration moves to message dispatch
- `conversation-manager`: Requirements change — thread history and LLM dispatch responsibilities move into `CharacterActor`; the `Manager` type is removed

## Impact

- `internal/messaging/` — new package: `Message`, `MessageBus`, `Actor` interface
- `internal/character/actor.go` — new file: `CharacterActor` wrapping `Character`
- `internal/conversation/manager.go` — **deleted**
- `internal/simulation/engine.go` — refactored: use bus instead of Manager
- `cmd/simulator/main.go` — wiring updated: instantiate bus + actors instead of Manager
