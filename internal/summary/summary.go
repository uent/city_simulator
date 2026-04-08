package summary

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/llm"
	"github.com/jnn-z/city_simulator/internal/scenario"
	"github.com/jnn-z/city_simulator/internal/world"
)

const maxEvents = 100

// GenerateSummary asks the LLM to produce a narrative summary of the simulation.
// It uses world events (capped at the last 100) and each character's final state.
// ctx is accepted for API consistency but the underlying LLM client is not context-aware.
func GenerateSummary(_ context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario) (string, error) {
	system, user := buildPrompt(w, chars, sc)
	text, err := client.Generate(system, user)
	if err != nil {
		return "", fmt.Errorf("summary LLM call failed: %w", err)
	}
	return text, nil
}

// SaveSummary writes content to summaries/<scenarioName>-<timestamp>.md.
// The directory is created if it does not exist.
// Returns the path of the written file.
func SaveSummary(scenarioName, content string) (string, error) {
	if err := os.MkdirAll("summaries", 0o755); err != nil {
		return "", fmt.Errorf("create summary dir: %w", err)
	}

	ts := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	filename := fmt.Sprintf("%s-%s.md", scenarioName, ts)
	path := filepath.Join("summaries", filename)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write summary file: %w", err)
	}
	return path, nil
}

func buildPrompt(w *world.State, chars []*character.Character, sc scenario.Scenario) (system, user string) {
	system = "You are a narrative chronicler. Given the events and character states of a simulation, write a cohesive story summary in prose. Be vivid and concise — two to four paragraphs."

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scenario: %s\n", sc.Name))
	sb.WriteString(fmt.Sprintf("Total ticks: %d\n\n", w.Tick))

	// Events — cap at last maxEvents
	events := w.EventLog
	if len(events) > maxEvents {
		events = events[len(events)-maxEvents:]
	}
	if len(events) > 0 {
		sb.WriteString("Events (chronological):\n")
		for _, e := range events {
			line := fmt.Sprintf("- [tick %d] %s", e.Tick, e.Description)
			if e.Location != "" {
				line += fmt.Sprintf(" (at %s)", e.Location)
			}
			sb.WriteString(line + "\n")
		}
	}

	// Character final states
	if len(chars) > 0 {
		sb.WriteString("\nCharacter states at end of simulation:\n")
		for _, c := range chars {
			sb.WriteString(fmt.Sprintf("- %s (%s): currently at %s, feeling %s\n",
				c.Name, c.Occupation, c.Location, c.EmotionalState))
		}
	}

	sb.WriteString("\nWrite the narrative summary now.")
	return system, sb.String()
}
