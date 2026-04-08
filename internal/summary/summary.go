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

const maxEvents = 200

// GenerateSummary asks the LLM to produce a narrative summary of the simulation.
// It uses world events (capped at the last 100) and each character's final state.
// ctx is accepted for API consistency but the underlying LLM client is not context-aware.
func GenerateSummary(_ context.Context, client *llm.Client, w *world.State, chars []*character.Character, sc scenario.Scenario, language string) (string, error) {
	system, user := buildPrompt(w, chars, sc, language)
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

func buildPrompt(w *world.State, chars []*character.Character, sc scenario.Scenario, language string) (system, user string) {
	system = "You are a narrative chronicler. Given the events and character states of a simulation, write a rich, immersive story summary in prose. Cover the arc of the simulation, key turning points, how characters evolved, and the final outcome. Write at least six detailed paragraphs — do not summarize briefly."
	if language != "" {
		system += fmt.Sprintf(" Respond in %s.", language)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scenario: %s\n", sc.Name))
	sb.WriteString(fmt.Sprintf("Total ticks: %d\n", w.Tick))

	// World concept — premise, flavor, rules
	concept := sc.World.Concept
	if concept.Premise != "" || concept.Flavor != "" || len(concept.Rules) > 0 {
		sb.WriteString("\nWorld concept:\n")
		if concept.Premise != "" {
			sb.WriteString(fmt.Sprintf("- Premise: %s\n", concept.Premise))
		}
		if concept.Flavor != "" {
			sb.WriteString(fmt.Sprintf("- Tone: %s\n", concept.Flavor))
		}
		if len(concept.Rules) > 0 {
			sb.WriteString("- Rules:\n")
			for _, r := range concept.Rules {
				sb.WriteString(fmt.Sprintf("  - %s\n", r))
			}
		}
	}

	// World atmosphere and weather
	if sc.World.Atmosphere != "" || sc.World.Weather != "" {
		sb.WriteString("\nSetting:\n")
		if sc.World.Atmosphere != "" {
			sb.WriteString(fmt.Sprintf("- Atmosphere: %s\n", sc.World.Atmosphere))
		}
		if sc.World.Weather != "" {
			sb.WriteString(fmt.Sprintf("- Weather: %s\n", sc.World.Weather))
		}
	}

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
			line := fmt.Sprintf("- %s (%s)", c.Name, c.Occupation)
			if c.Location != "" {
				line += fmt.Sprintf(", at %s", c.Location)
			}
			line += fmt.Sprintf(", feeling %s", c.EmotionalState)
			sb.WriteString(line + "\n")
			if c.Motivation != "" {
				sb.WriteString(fmt.Sprintf("  Motivation: %s\n", c.Motivation))
			}
			if c.Fear != "" {
				sb.WriteString(fmt.Sprintf("  Fear: %s\n", c.Fear))
			}
			if len(c.Goals) > 0 {
				sb.WriteString(fmt.Sprintf("  Goals: %s\n", strings.Join(c.Goals, "; ")))
			}
		}
	}

	sb.WriteString("\nWrite the narrative summary now.")
	return system, sb.String()
}
