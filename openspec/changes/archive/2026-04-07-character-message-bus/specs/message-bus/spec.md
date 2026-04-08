## ADDED Requirements

### Requirement: Message type with routing and reply semantics
The system SHALL define a `Message` struct in `internal/messaging/` with the following fields:
- `Type MessageType` — discriminator for routing and processing (e.g., `CharChat`, `DirectorDirective`, `CharReport`, `Broadcast`)
- `From string` — sender actor ID
- `To string` — recipient actor ID; empty string means broadcast
- `Tick int` — simulation tick when the message was sent
- `Payload any` — type-specific data (character state snapshot, dialogue text, event list, etc.)
- `ReplyChan chan Message` — buffered channel (capacity 1) for the recipient to write exactly one reply; nil for fire-and-forget messages

#### Scenario: Reply channel capacity
- **WHEN** a `Message` is constructed with a non-nil `ReplyChan`
- **THEN** the channel SHALL have a buffer capacity of at least 1 so the recipient can write the reply without blocking regardless of whether the sender has started reading

#### Scenario: Fire-and-forget message
- **WHEN** a `Message` is constructed with a nil `ReplyChan`
- **THEN** the recipient SHALL process the message without attempting to write a reply

### Requirement: MessageBus routes messages to registered actor inboxes
The system SHALL provide a `MessageBus` type that maintains a registry of all active actors keyed by actor ID. The bus SHALL expose:
- `Register(id string, inbox chan Message)` — adds an actor inbox to the registry
- `Send(msg Message) error` — writes `msg` to the inbox of the actor identified by `msg.To`; returns an error if the actor is not registered
- `Broadcast(msg Message, excludeIDs ...string) (<-chan Message, error)` — sends `msg` concurrently to every registered actor except those listed in `excludeIDs`; returns a receive-only channel that emits one reply per recipient as replies arrive; returns an error if no recipients exist

#### Scenario: Send to registered actor
- **WHEN** `Send` is called with a `msg.To` that matches a registered actor ID
- **THEN** the message SHALL be written to that actor's inbox channel without blocking (inbox is buffered)

#### Scenario: Send to unknown actor
- **WHEN** `Send` is called with a `msg.To` that is not in the registry
- **THEN** `Send` SHALL return a non-nil error and the message SHALL NOT be delivered

#### Scenario: Broadcast to all characters
- **WHEN** `Broadcast` is called with no `excludeIDs`
- **THEN** the message SHALL be delivered to every registered actor's inbox concurrently via one goroutine per actor, and the returned channel SHALL emit exactly N reply messages (one per actor)

#### Scenario: Broadcast with exclusions
- **WHEN** `Broadcast` is called with one or more `excludeIDs`
- **THEN** the listed actors SHALL NOT receive the message, and the reply channel SHALL emit replies only from the non-excluded actors

#### Scenario: Broadcast with no eligible recipients
- **WHEN** `Broadcast` is called and all registered actors are in `excludeIDs`
- **THEN** `Broadcast` SHALL return a nil channel and a non-nil error

### Requirement: Concurrent broadcast does not deadlock
The `Broadcast` implementation SHALL use one goroutine per recipient to write to each inbox. The returned fan-in channel SHALL be closed after all N replies have been forwarded, so the caller can range over it safely.

#### Scenario: All actors reply before caller reads
- **WHEN** all N actor goroutines write their replies before the engine starts reading from the fan-in channel
- **THEN** no goroutine SHALL block and all N replies SHALL be available on the fan-in channel

#### Scenario: Caller reads replies out of arrival order
- **WHEN** the engine reads from the fan-in channel
- **THEN** replies SHALL arrive in the order actors complete their processing (non-deterministic); the engine SHALL NOT assume a fixed order
