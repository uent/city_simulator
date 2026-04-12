package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// openRouterRequest is the OpenAI-compatible payload sent to /chat/completions.
type openRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// openRouterResponse is the response from /chat/completions.
type openRouterResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage *openRouterUsage `json:"usage,omitempty"`
}

// openRouterUsage contains the token counts and cost from OpenRouter.
type openRouterUsage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	Cost             float64 `json:"cost"` // OpenRouter-specific USD cost field
}

// OpenRouterClient calls the OpenRouter chat-completions endpoint and tracks usage.
type OpenRouterClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// NewOpenRouterClient creates an OpenRouter client. No network call is made.
func NewOpenRouterClient(apiKey, model, baseURL string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Ping validates that the API key is accepted by calling GET /models.
func (c *OpenRouterClient) Ping() error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/models", nil)
	if err != nil {
		return fmt.Errorf("openrouter ping: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot reach OpenRouter at %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("OpenRouter authentication failed (status 401): check OPENROUTER_API_KEY")
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenRouter /models returned status %d: %s", resp.StatusCode, string(raw))
	}
	return nil
}

// Generate sends a system prompt and a user prompt to OpenRouter and returns the
// assistant's reply and token usage. It is a convenience wrapper around Chat.
func (c *OpenRouterClient) Generate(systemPrompt, userPrompt string, opts ...Option) (string, Usage, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	return c.Chat(messages, opts...)
}

// Chat sends messages to /chat/completions and returns the assistant's reply and usage.
func (c *OpenRouterClient) Chat(messages []Message, opts ...Option) (string, Usage, error) {
	// Build request; apply options via a temporary ChatRequest to reuse the Option type.
	tmp := &ChatRequest{Model: c.model}
	for _, o := range opts {
		o(tmp)
	}

	payload := openRouterRequest{
		Model:    tmp.Model,
		Messages: messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", Usage{}, fmt.Errorf("marshal openrouter request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", Usage{}, fmt.Errorf("openrouter: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", Usage{}, fmt.Errorf("chat request to OpenRouter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", Usage{}, fmt.Errorf("OpenRouter /chat/completions returned status %d: %s", resp.StatusCode, string(raw))
	}

	var orResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return "", Usage{}, fmt.Errorf("decode openrouter response: %w", err)
	}

	if len(orResp.Choices) == 0 {
		return "", Usage{}, fmt.Errorf("openrouter response contained no choices")
	}

	text := orResp.Choices[0].Message.Content

	var usage Usage
	if orResp.Usage != nil {
		usage = Usage{
			PromptTokens:     orResp.Usage.PromptTokens,
			CompletionTokens: orResp.Usage.CompletionTokens,
			EstimatedCostUSD: orResp.Usage.Cost,
		}
	}

	return text, usage, nil
}
