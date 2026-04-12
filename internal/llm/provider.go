package llm

import "fmt"

// Usage holds token consumption and cost data from a single LLM call.
// Providers that do not report token counts (e.g. Ollama) return the zero value.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	EstimatedCostUSD float64
}

// Provider is the common interface for all LLM backends.
// Every Generate and Chat call returns a Usage alongside the response text.
type Provider interface {
	// Generate sends a system prompt and a user prompt and returns the assistant's reply.
	Generate(systemPrompt, userPrompt string, opts ...Option) (string, Usage, error)
	// Chat sends an arbitrary message slice and returns the assistant's reply.
	Chat(messages []Message, opts ...Option) (string, Usage, error)
	// Ping checks that the backend is reachable and credentials are valid.
	Ping() error
}

// ProviderConfig carries all fields needed to initialise either backend.
type ProviderConfig struct {
	ProviderName string // "ollama" or "openrouter"

	// Ollama fields
	OllamaURL   string
	OllamaModel string

	// OpenRouter fields
	OpenRouterAPIKey  string
	OpenRouterModel   string
	OpenRouterBaseURL string
}

// NewProvider returns the Provider implementation requested by cfg.ProviderName.
// Supported values: "ollama", "openrouter".
func NewProvider(cfg ProviderConfig) (Provider, error) {
	switch cfg.ProviderName {
	case "ollama":
		client := NewClient(cfg.OllamaURL, cfg.OllamaModel)
		return &OllamaProvider{client: client}, nil

	case "openrouter":
		if cfg.OpenRouterAPIKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY is required when LLM_PROVIDER=openrouter")
		}
		baseURL := cfg.OpenRouterBaseURL
		if baseURL == "" {
			baseURL = "https://openrouter.ai/api/v1"
		}
		model := cfg.OpenRouterModel
		if model == "" {
			model = "openai/gpt-4o-mini"
		}
		return NewOpenRouterClient(cfg.OpenRouterAPIKey, model, baseURL), nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider %q: supported values are \"ollama\" and \"openrouter\"", cfg.ProviderName)
	}
}

// OllamaProvider wraps the Ollama *Client and implements Provider.
// All calls return zero Usage because Ollama does not expose token counts.
type OllamaProvider struct {
	client *Client
}

// Generate delegates to the underlying Client and returns zero Usage.
func (p *OllamaProvider) Generate(systemPrompt, userPrompt string, opts ...Option) (string, Usage, error) {
	text, usage, err := p.client.Generate(systemPrompt, userPrompt, opts...)
	return text, usage, err
}

// Chat delegates to the underlying Client and returns zero Usage.
func (p *OllamaProvider) Chat(messages []Message, opts ...Option) (string, Usage, error) {
	text, usage, err := p.client.Chat(messages, opts...)
	return text, usage, err
}

// Ping delegates to the underlying Client.
func (p *OllamaProvider) Ping() error {
	return p.client.Ping()
}
