package character

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jnn-z/city_simulator/internal/llm"
	"github.com/jnn-z/city_simulator/internal/messaging"
)

// Turn records one utterance in a conversation thread between two characters.
type Turn struct {
	Tick    int
	Speaker string // character ID
	Text    string
}

// CharacterActor wraps a Character and processes inbox messages via a goroutine.
// It satisfies messaging.Actor.
type CharacterActor struct {
	char       *Character
	llmClient  llm.Provider
	acc        *llm.CostAccumulator
	inbox      chan messaging.Message
	history    map[string][]Turn // keyed by peer character ID
	maxHistory int
	once       sync.Once
	cancel     context.CancelFunc
}

// NewCharacterActor creates a CharacterActor backed by the given character and LLM provider.
func NewCharacterActor(c *Character, llmClient llm.Provider, acc *llm.CostAccumulator) *CharacterActor {
	return &CharacterActor{
		char:       c,
		llmClient:  llmClient,
		acc:        acc,
		inbox:      make(chan messaging.Message, 4),
		history:    make(map[string][]Turn),
		maxHistory: 20,
	}
}

func (a *CharacterActor) ID() string                { return a.char.ID }
func (a *CharacterActor) Inbox() chan messaging.Message { return a.inbox }

// Start spawns the processing goroutine exactly once.
func (a *CharacterActor) Start(ctx context.Context) {
	a.once.Do(func() {
		ctx, a.cancel = context.WithCancel(ctx)
		go a.run(ctx)
	})
}

// Stop signals the goroutine to exit.
func (a *CharacterActor) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *CharacterActor) run(ctx context.Context) {
	for {
		select {
		case msg := <-a.inbox:
			a.process(msg)
		case <-ctx.Done():
			return
		}
	}
}

func (a *CharacterActor) process(msg messaging.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("actor %s: panic: %v", a.char.ID, r)
			messaging.Reply(msg, nil)
		}
	}()
	switch msg.Type {
	case messaging.CharChat:
		a.handleCharChat(msg)
	case messaging.MoveDecision:
		a.handleMoveDecision(msg)
	case messaging.DirectorDirective:
		// Director has already updated world state; just acknowledge.
		messaging.Reply(msg, nil)
	default:
		log.Printf("actor %s: unhandled message type %v", a.char.ID, msg.Type)
		messaging.Reply(msg, nil)
	}
}

// handleCharChat generates the initiator's opening line and the responder's reply
// using two sequential LLM calls. Both texts are returned in CharChatReply.
func (a *CharacterActor) handleCharChat(msg messaging.Message) {
	payload, ok := msg.Payload.(messaging.CharChatPayload)
	if !ok {
		messaging.Reply(msg, messaging.CharChatReply{Err: fmt.Errorf("invalid CharChat payload")})
		return
	}

	history := a.buildLLMHistory(msg.From)

	// Step 1: Generate initiator's opening line (using pre-built initiator system prompt).
	initiatorMsgs := []llm.Message{
		{Role: "system", Content: payload.InitiatorSystem},
	}
	initiatorMsgs = append(initiatorMsgs, history...)
	initiatorMsgs = append(initiatorMsgs, llm.Message{
		Role:    "user",
		Content: fmt.Sprintf("You encounter %s. What do you say?", a.char.Name),
	})

	initiatorRaw, usage, err := a.llmClient.Chat(initiatorMsgs)
	if err != nil {
		messaging.Reply(msg, messaging.CharChatReply{
			Err: fmt.Errorf("LLM call for initiator %s: %w", payload.InitiatorName, err),
		})
		return
	}
	if a.acc != nil {
		a.acc.Add(usage)
	}
	initiatorExpr := ParseExpression(initiatorRaw)

	// Step 2: Generate responder's reply (using pre-built responder system prompt).
	// Pass the formatted expression so the responder sees both action and speech.
	responderMsgs := []llm.Message{
		{Role: "system", Content: payload.ResponderSystem},
		{Role: "user", Content: FormatExpression(initiatorExpr)},
	}

	responderRaw, usage, err := a.llmClient.Chat(responderMsgs)
	if err != nil {
		messaging.Reply(msg, messaging.CharChatReply{
			Err: fmt.Errorf("LLM call for responder %s: %w", a.char.Name, err),
		})
		return
	}
	if a.acc != nil {
		a.acc.Add(usage)
	}
	responderExpr := ParseExpression(responderRaw)

	// Step 3: Update per-pair history and character memory (store raw text including markers).
	a.appendHistory(msg.From, Turn{Tick: msg.Tick, Speaker: msg.From, Text: initiatorRaw})
	a.appendHistory(msg.From, Turn{Tick: msg.Tick, Speaker: a.char.ID, Text: responderRaw})
	a.char.AddMemory(MemoryEntry{Tick: msg.Tick, Speaker: payload.InitiatorName, Text: initiatorRaw})
	a.char.AddMemory(MemoryEntry{Tick: msg.Tick, Speaker: a.char.Name, Text: responderRaw})

	messaging.Reply(msg, messaging.CharChatReply{
		InitiatorSpeech: initiatorExpr.Speech,
		InitiatorAction: initiatorExpr.Action,
		ResponderSpeech: responderExpr.Speech,
		ResponderAction: responderExpr.Action,
	})
}

// handleMoveDecision uses the pre-built movement prompt from the payload.
func (a *CharacterActor) handleMoveDecision(msg messaging.Message) {
	payload, ok := msg.Payload.(messaging.MoveDecisionPayload)
	if !ok || len(payload.Locations) == 0 {
		messaging.Reply(msg, messaging.MoveDecisionReply{Location: "stay"})
		return
	}

	msgs := []llm.Message{
		{Role: "system", Content: payload.SystemPrompt},
		{Role: "user", Content: "Where do you go next?"},
	}

	raw, usage, err := a.llmClient.Chat(msgs)
	if err != nil {
		log.Printf("actor %s: movement LLM error: %v", a.char.ID, err)
		messaging.Reply(msg, messaging.MoveDecisionReply{Location: "stay"})
		return
	}
	if a.acc != nil {
		a.acc.Add(usage)
	}

	decision := strings.TrimSpace(raw)
	for _, loc := range payload.Locations {
		if decision == loc {
			messaging.Reply(msg, messaging.MoveDecisionReply{Location: loc})
			return
		}
	}
	for _, loc := range payload.Locations {
		if strings.EqualFold(decision, loc) {
			messaging.Reply(msg, messaging.MoveDecisionReply{Location: loc})
			return
		}
	}
	for _, loc := range payload.Locations {
		if strings.Contains(decision, loc) {
			messaging.Reply(msg, messaging.MoveDecisionReply{Location: loc})
			return
		}
	}
	messaging.Reply(msg, messaging.MoveDecisionReply{Location: "stay"})
}

// buildLLMHistory returns the last ≤10 turns with peerID as llm.Message values.
func (a *CharacterActor) buildLLMHistory(peerID string) []llm.Message {
	turns := a.history[peerID]
	if len(turns) > 10 {
		turns = turns[len(turns)-10:]
	}
	msgs := make([]llm.Message, 0, len(turns))
	for _, t := range turns {
		role := "assistant"
		if t.Speaker == peerID {
			role = "user"
		}
		msgs = append(msgs, llm.Message{Role: role, Content: t.Text})
	}
	return msgs
}

// appendHistory adds a turn, evicting the oldest entry if over maxHistory.
func (a *CharacterActor) appendHistory(peerID string, t Turn) {
	a.history[peerID] = append(a.history[peerID], t)
	if len(a.history[peerID]) > a.maxHistory {
		a.history[peerID] = a.history[peerID][1:]
	}
}
