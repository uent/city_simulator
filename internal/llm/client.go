package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message is a single chat message with a role and content.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the payload sent to /api/chat.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse is the response from /api/chat (non-streaming).
type ChatResponse struct {
	Message Message `json:"message"`
}

// Option allows optional overrides on a Chat call.
type Option func(*ChatRequest)

// WithModel overrides the model for a single call.
func WithModel(model string) Option {
	return func(r *ChatRequest) { r.Model = model }
}

// Client wraps the Ollama HTTP API.
type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewClient creates an Ollama client. No network call is made.
func NewClient(baseURL, model string) *Client {
	return &Client{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Ping checks that Ollama is reachable by calling /api/tags.
func (c *Client) Ping() error {
	resp, err := c.httpClient.Get(c.baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("cannot reach Ollama at %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama at %s returned status %d", c.baseURL, resp.StatusCode)
	}
	return nil
}

// Generate sends a system prompt and a user prompt to Ollama and returns the
// assistant's reply. It is a convenience wrapper around Chat.
func (c *Client) Generate(systemPrompt, userPrompt string, opts ...Option) (string, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	return c.Chat(messages, opts...)
}

// Chat sends messages to /api/chat and returns the assistant's reply text.
func (c *Client) Chat(messages []Message, opts ...Option) (string, error) {
	req := &ChatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}
	for _, o := range opts {
		o(req)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/api/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("chat request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama /api/chat returned status %d: %s", resp.StatusCode, string(raw))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	return chatResp.Message.Content, nil
}
