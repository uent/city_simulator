package conversation

import (
	"context"
	"fmt"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/llm"
	"github.com/jnn-z/city_simulator/internal/world"
)

// Turn records one utterance in a conversation thread.
type Turn struct {
	Tick    int
	Speaker string
	Text    string
}

// Thread is the ordered history of turns between two characters.
type Thread struct {
	Turns []Turn
}

// Manager owns all conversation threads and drives exchanges.
type Manager struct {
	threads   map[string]*Thread
	llmClient *llm.Client
}

// NewManager creates a Manager backed by the given LLM client.
func NewManager(llmClient *llm.Client) *Manager {
	return &Manager{
		threads:   make(map[string]*Thread),
		llmClient: llmClient,
	}
}

func threadKey(fromID, toID string) string {
	return fmt.Sprintf("%s→%s", fromID, toID)
}

// AddTurn appends a turn to the thread for fromID→toID, creating it if needed.
func (m *Manager) AddTurn(fromID, toID, text string, tick int) {
	key := threadKey(fromID, toID)
	if m.threads[key] == nil {
		m.threads[key] = &Thread{}
	}
	m.threads[key].Turns = append(m.threads[key].Turns, Turn{
		Tick:    tick,
		Speaker: fromID,
		Text:    text,
	})
}

// History returns the last maxTurns turns as LLM messages.
// Role "user" is assigned to the initiator (fromID), "assistant" to the responder.
func (m *Manager) History(fromID, toID string, maxTurns int) []llm.Message {
	key := threadKey(fromID, toID)
	thread := m.threads[key]
	if thread == nil {
		return nil
	}
	turns := thread.Turns
	if maxTurns > 0 && len(turns) > maxTurns {
		turns = turns[len(turns)-maxTurns:]
	}
	msgs := make([]llm.Message, 0, len(turns))
	for _, t := range turns {
		role := "assistant"
		if t.Speaker == fromID {
			role = "user"
		}
		msgs = append(msgs, llm.Message{Role: role, Content: t.Text})
	}
	return msgs
}

// ExchangeResult holds the generated dialogue from a single exchange.
type ExchangeResult struct {
	InitiatorText string
	ResponderText string
}

// RunExchange drives a full two-turn interaction between initiator and responder.
// Steps: build prompts → initiator speaks → responder replies → store memories → log event.
func (m *Manager) RunExchange(
	ctx context.Context,
	initiator, responder *character.Character,
	w *world.State,
	tick int,
) (ExchangeResult, error) {
	// 1. Build initiator system prompt (persona + world context)
	initiatorSystem := llm.BuildSystemPrompt(*initiator) +
		"\n\nWorld context:\n" + w.Summary()

	// 2. Assemble history for this pair
	history := m.History(initiator.ID, responder.ID, 10)

	// 3. Build initiator messages: system + history + prompt to speak
	initiatorMsgs := []llm.Message{
		{Role: "system", Content: initiatorSystem},
	}
	initiatorMsgs = append(initiatorMsgs, history...)
	initiatorMsgs = append(initiatorMsgs, llm.Message{
		Role:    "user",
		Content: fmt.Sprintf("You encounter %s. What do you say?", responder.Name),
	})

	// 4. Generate initiator's message
	initiatorText, err := m.llmClient.Chat(initiatorMsgs)
	if err != nil {
		return ExchangeResult{}, fmt.Errorf("LLM call for initiator %s: %w", initiator.Name, err)
	}

	// 5. Store initiator message in both memories and thread
	entry := character.MemoryEntry{Tick: tick, Speaker: initiator.Name, Text: initiatorText}
	initiator.AddMemory(entry)
	responder.AddMemory(entry)
	m.AddTurn(initiator.ID, responder.ID, initiatorText, tick)

	// 6. Build responder system prompt
	responderSystem := llm.BuildSystemPrompt(*responder) +
		"\n\nWorld context:\n" + w.Summary()

	responderMsgs := []llm.Message{
		{Role: "system", Content: responderSystem},
		{Role: "user", Content: initiatorText},
	}

	// 7. Generate responder's reply
	responderText, err := m.llmClient.Chat(responderMsgs)
	if err != nil {
		return ExchangeResult{}, fmt.Errorf("LLM call for responder %s: %w", responder.Name, err)
	}

	// 8. Store responder reply in both memories and thread
	replyEntry := character.MemoryEntry{Tick: tick, Speaker: responder.Name, Text: responderText}
	initiator.AddMemory(replyEntry)
	responder.AddMemory(replyEntry)
	m.AddTurn(responder.ID, initiator.ID, responderText, tick)

	// Append world event
	w.AppendEvent(world.Event{
		Tick:         tick,
		Type:         "conversation",
		Description:  fmt.Sprintf("%s spoke to %s", initiator.Name, responder.Name),
		Participants: []string{initiator.ID, responder.ID},
	})

	return ExchangeResult{InitiatorText: initiatorText, ResponderText: responderText}, nil
}
