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
	return text + renderCharacterCards(chars), nil
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

// renderCharacterCards produces a Markdown "Character Cards" section for all
// non-director characters. Returns empty string if there are no eligible characters.
func renderCharacterCards(chars []*character.Character) string {
	var cards []string
	for _, c := range chars {
		if c.Type == "game_director" {
			continue
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("### %s\n\n", c.Name))

		if c.Age != 0 {
			sb.WriteString(fmt.Sprintf("- **Age:** %d\n", c.Age))
		}
		if c.Occupation != "" {
			sb.WriteString(fmt.Sprintf("- **Occupation:** %s\n", c.Occupation))
		}
		if c.Appearance != "" {
			sb.WriteString(fmt.Sprintf("- **Appearance:** %s\n", c.Appearance))
		}
		if c.Location != "" {
			sb.WriteString(fmt.Sprintf("- **Location:** %s\n", c.Location))
		}
		if c.EmotionalState != "" {
			sb.WriteString(fmt.Sprintf("- **Emotional State:** %s\n", c.EmotionalState))
		}
		if len(c.Goals) > 0 {
			sb.WriteString(fmt.Sprintf("- **Goals:** %s\n", strings.Join(c.Goals, "; ")))
		}
		if c.Motivation != "" {
			sb.WriteString(fmt.Sprintf("- **Motivation:** %s\n", c.Motivation))
		}
		if c.Fear != "" {
			sb.WriteString(fmt.Sprintf("- **Fear:** %s\n", c.Fear))
		}
		if c.CoreBelief != "" {
			sb.WriteString(fmt.Sprintf("- **Core Belief:** %s\n", c.CoreBelief))
		}
		if c.InternalTension != "" {
			sb.WriteString(fmt.Sprintf("- **Internal Tension:** %s\n", c.InternalTension))
		}

		if len(c.FormativeEvents) > 0 {
			sb.WriteString("- **Formative Events:**\n")
			for _, e := range c.FormativeEvents {
				sb.WriteString(fmt.Sprintf("  - %s\n", e))
			}
		}

		v := c.Voice
		if v.Formality != "" || v.VerbalTics != "" || v.ResponseLength != "" || v.HumorType != "" || v.CommunicationStyle != "" {
			sb.WriteString("- **Voice:**\n")
			if v.Formality != "" {
				sb.WriteString(fmt.Sprintf("  - Formality: %s\n", v.Formality))
			}
			if v.VerbalTics != "" {
				sb.WriteString(fmt.Sprintf("  - Verbal Tics: %s\n", v.VerbalTics))
			}
			if v.ResponseLength != "" {
				sb.WriteString(fmt.Sprintf("  - Response Length: %s\n", v.ResponseLength))
			}
			if v.HumorType != "" {
				sb.WriteString(fmt.Sprintf("  - Humor Type: %s\n", v.HumorType))
			}
			if v.CommunicationStyle != "" {
				sb.WriteString(fmt.Sprintf("  - Communication Style: %s\n", v.CommunicationStyle))
			}
		}

		r := c.RelationalDefaults
		if r.Strangers != "" || r.Authority != "" || r.Vulnerable != "" {
			sb.WriteString("- **Relational Defaults:**\n")
			if r.Strangers != "" {
				sb.WriteString(fmt.Sprintf("  - Strangers: %s\n", r.Strangers))
			}
			if r.Authority != "" {
				sb.WriteString(fmt.Sprintf("  - Authority: %s\n", r.Authority))
			}
			if r.Vulnerable != "" {
				sb.WriteString(fmt.Sprintf("  - Vulnerable: %s\n", r.Vulnerable))
			}
		}

		if c.CoverIdentity != nil {
			ci := c.CoverIdentity
			sb.WriteString("- **Cover Identity:**\n")
			if ci.Alias != "" {
				sb.WriteString(fmt.Sprintf("  - Alias: %s\n", ci.Alias))
			}
			if ci.Role != "" {
				sb.WriteString(fmt.Sprintf("  - Role: %s\n", ci.Role))
			}
			if ci.Backstory != "" {
				sb.WriteString(fmt.Sprintf("  - Backstory: %s\n", ci.Backstory))
			}
			if len(ci.Weaknesses) > 0 {
				sb.WriteString("  - Weaknesses:\n")
				for _, w := range ci.Weaknesses {
					sb.WriteString(fmt.Sprintf("    - %s\n", w))
				}
			}
		}

		if len(c.DialogueExamples) > 0 {
			sb.WriteString("- **Dialogue Examples:**\n")
			for _, d := range c.DialogueExamples {
				sb.WriteString(fmt.Sprintf("  - %s\n", d))
			}
		}

		if len(c.Judgments) > 0 {
			sb.WriteString("- **Relationships:**\n")
			for _, j := range c.Judgments {
				sb.WriteString(fmt.Sprintf("  - **%s** — Trust: %s | Interest: %s | Threat: %s\n", j.Name, j.Trust, j.Interest, j.Threat))
				if j.Impression != "" {
					sb.WriteString(fmt.Sprintf("    > %s\n", j.Impression))
				}
			}
		}

		cards = append(cards, sb.String())
	}

	if len(cards) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString("\n\n---\n\n## Character Cards\n\n")
	result.WriteString(strings.Join(cards, "\n"))
	return result.String()
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
