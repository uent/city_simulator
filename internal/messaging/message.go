package messaging

// MessageType discriminates the intent of a Message.
type MessageType int

const (
	// CharChat asks the receiving actor to run a full exchange:
	// generate the initiator's line and then the receiver's reply.
	CharChat MessageType = iota
	// DirectorDirective signals that the game director has run this tick;
	// characters should acknowledge before proceeding with exchanges.
	DirectorDirective
	// CharReport is reserved for character→director communication (future use).
	CharReport
	// MoveDecision asks the receiving character where to move next.
	MoveDecision
)

// Message is the atomic unit of communication between simulation actors.
type Message struct {
	Type      MessageType
	From      string      // sender actor ID
	To        string      // recipient actor ID; empty for broadcast
	Tick      int
	Payload   any
	ReplyChan chan Message // buffered (cap 1); nil for fire-and-forget
}

// NewRequest creates a Message with a buffered reply channel (capacity 1).
func NewRequest(msgType MessageType, from, to string, tick int, payload any) Message {
	return Message{
		Type:      msgType,
		From:      from,
		To:        to,
		Tick:      tick,
		Payload:   payload,
		ReplyChan: make(chan Message, 1),
	}
}

// Reply writes a reply to msg.ReplyChan. No-op if ReplyChan is nil.
func Reply(msg Message, payload any) {
	if msg.ReplyChan == nil {
		return
	}
	msg.ReplyChan <- Message{
		Type:    msg.Type,
		From:    msg.To,
		To:      msg.From,
		Tick:    msg.Tick,
		Payload: payload,
	}
}

// --- Payload types ---

// CharChatPayload is the payload for a CharChat message sent to the responder actor.
// The engine pre-builds both system prompts so the actor avoids importing llm.
type CharChatPayload struct {
	InitiatorID     string // ID of the initiating character
	InitiatorName   string // display name of the initiator
	InitiatorSystem string // pre-built system prompt for initiator (persona + inbox + world ctx)
	ResponderSystem string // pre-built system prompt for responder (persona + inbox + world ctx)
}

// CharChatReply is the payload returned by the responder actor after a CharChat.
type CharChatReply struct {
	InitiatorSpeech string
	InitiatorAction string
	ResponderSpeech string
	ResponderAction string
	Err             error
}

// MoveDecisionPayload is the payload for a MoveDecision message.
type MoveDecisionPayload struct {
	SystemPrompt string   // pre-built movement prompt
	Locations    []string // valid location names for response matching
}

// MoveDecisionReply is the payload returned after a MoveDecision message.
type MoveDecisionReply struct {
	Location string // chosen location name, or "stay"
}
