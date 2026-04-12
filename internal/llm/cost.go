package llm

import "sync"

// CostAccumulator aggregates token usage across multiple LLM calls.
// It is safe for concurrent use.
type CostAccumulator struct {
	mu               sync.Mutex
	promptTokens     int
	completionTokens int
	estimatedCostUSD float64
}

// Add adds the given usage to the running totals.
func (a *CostAccumulator) Add(u Usage) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.promptTokens += u.PromptTokens
	a.completionTokens += u.CompletionTokens
	a.estimatedCostUSD += u.EstimatedCostUSD
}

// Total returns a snapshot of the accumulated totals.
func (a *CostAccumulator) Total() Usage {
	a.mu.Lock()
	defer a.mu.Unlock()
	return Usage{
		PromptTokens:     a.promptTokens,
		CompletionTokens: a.completionTokens,
		EstimatedCostUSD: a.estimatedCostUSD,
	}
}
