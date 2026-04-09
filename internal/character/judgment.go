package character

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jnn-z/city_simulator/internal/llm"
)

// ObservableProfile is the filtered view of a character that others can perceive
// on first encounter. Private psychological fields are never included.
// If the character has a CoverIdentity, the alias and role are used instead of
// the real name and occupation.
type ObservableProfile struct {
	ID             string
	Name           string
	Age            int
	Occupation     string
	EmotionalState string
	Appearance     string
	Location       string
}

// ObservableSnapshot returns the observable profile of c. If c has a CoverIdentity,
// its Alias and Role are used instead of the real Name and Occupation.
func ObservableSnapshot(c Character) ObservableProfile {
	name := c.Name
	occupation := c.Occupation
	if c.CoverIdentity != nil {
		if c.CoverIdentity.Alias != "" {
			name = c.CoverIdentity.Alias
		}
		if c.CoverIdentity.Role != "" {
			occupation = c.CoverIdentity.Role
		}
	}
	return ObservableProfile{
		ID:             c.ID,
		Name:           name,
		Age:            c.Age,
		Occupation:     occupation,
		EmotionalState: c.EmotionalState,
		Appearance:     c.Appearance,
		Location:       c.Location,
	}
}

// BuildJudgmentPrompt constructs the LLM prompt used to form a first impression.
// The judging character's full psychological profile is the interpretive lens;
// only the target's observable profile is exposed.
func BuildJudgmentPrompt(judge Character, target ObservableProfile, language string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s, a %d-year-old %s.\n", judge.Name, judge.Age, judge.Occupation))
	if judge.Motivation != "" {
		sb.WriteString(fmt.Sprintf("Your motivation: %s\n", judge.Motivation))
	}
	if judge.Fear != "" {
		sb.WriteString(fmt.Sprintf("Your fear: %s\n", judge.Fear))
	}
	if judge.CoreBelief != "" {
		sb.WriteString(fmt.Sprintf("Your core belief: %s\n", judge.CoreBelief))
	}
	if judge.InternalTension != "" {
		sb.WriteString(fmt.Sprintf("Your internal tension: %s\n", judge.InternalTension))
	}

	sb.WriteString(fmt.Sprintf("\nYou are forming a first impression of %s", target.Name))
	if target.Occupation != "" {
		sb.WriteString(fmt.Sprintf(", a %d-year-old %s", target.Age, target.Occupation))
	}
	sb.WriteString(".\n")

	if target.EmotionalState != "" && target.EmotionalState != "neutral" {
		sb.WriteString(fmt.Sprintf("Their demeanor: %s.\n", target.EmotionalState))
	}
	if target.Appearance != "" {
		sb.WriteString(fmt.Sprintf("What you observe: %s\n", target.Appearance))
	}
	if target.Location != "" {
		sb.WriteString(fmt.Sprintf("Where you encounter them: %s\n", target.Location))
	}

	sb.WriteString("\nBased entirely on your own psychology, fears, and personality — not objective fact — what is your gut reaction to this person?\n")
	sb.WriteString("Respond with a JSON object only, no other text:\n")
	sb.WriteString(`{"impression": "<2-3 sentences, first-person internal thought>", "trust": "<high|medium|low|none>", "interest": "<high|medium|low>", "threat": "<high|medium|low|none>"}`)
	if language != "" {
		sb.WriteString(fmt.Sprintf("\nWrite the impression field in %s.", language))
	}
	return sb.String()
}

// BuildUpdatePrompt constructs the LLM prompt for refreshing a judgment after
// 10 conversations. It includes the prior judgment and recent conversation history.
func BuildUpdatePrompt(judge Character, target ObservableProfile, prior CharacterJudgment, recentHistory []string, language string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s, a %d-year-old %s.\n", judge.Name, judge.Age, judge.Occupation))
	if judge.Motivation != "" {
		sb.WriteString(fmt.Sprintf("Your motivation: %s\n", judge.Motivation))
	}
	if judge.Fear != "" {
		sb.WriteString(fmt.Sprintf("Your fear: %s\n", judge.Fear))
	}
	if judge.CoreBelief != "" {
		sb.WriteString(fmt.Sprintf("Your core belief: %s\n", judge.CoreBelief))
	}
	if judge.InternalTension != "" {
		sb.WriteString(fmt.Sprintf("Your internal tension: %s\n", judge.InternalTension))
	}

	sb.WriteString(fmt.Sprintf("\nYou have now had multiple conversations with %s", target.Name))
	if target.Occupation != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", target.Occupation))
	}
	sb.WriteString(".\n")

	if target.Appearance != "" {
		sb.WriteString(fmt.Sprintf("What you observe about them: %s\n", target.Appearance))
	}

	sb.WriteString("\nYour prior impression of them:\n")
	sb.WriteString(fmt.Sprintf("%q\n", prior.Impression))
	sb.WriteString(fmt.Sprintf("Trust: %s | Interest: %s | Perceived threat: %s\n", prior.Trust, prior.Interest, prior.Threat))

	if len(recentHistory) > 0 {
		sb.WriteString("\nRecent exchanges between you:\n")
		for _, h := range recentHistory {
			sb.WriteString(fmt.Sprintf("- %s\n", h))
		}
	}

	sb.WriteString("\nBased on your interactions so far, has your view of them evolved? Update your impression.\n")
	sb.WriteString("Respond with a JSON object only, no other text:\n")
	sb.WriteString(`{"impression": "<2-3 sentences, first-person internal thought>", "trust": "<high|medium|low|none>", "interest": "<high|medium|low>", "threat": "<high|medium|low|none>"}`)
	if language != "" {
		sb.WriteString(fmt.Sprintf("\nWrite the impression field in %s.", language))
	}
	return sb.String()
}

// judgmentJSON is the expected LLM response shape.
type judgmentJSON struct {
	Impression string `json:"impression"`
	Trust      string `json:"trust"`
	Interest   string `json:"interest"`
	Threat     string `json:"threat"`
}

var validTrust    = map[string]bool{"high": true, "medium": true, "low": true, "none": true}
var validInterest = map[string]bool{"high": true, "medium": true, "low": true}
var validThreat   = map[string]bool{"high": true, "medium": true, "low": true, "none": true}

// ParseJudgmentResponse extracts a CharacterJudgment from a raw LLM response.
// It tolerates surrounding text by scanning for the first '{' and last '}'.
// Invalid enum values are replaced with safe defaults.
func ParseJudgmentResponse(raw string, about string, name string) CharacterJudgment {
	fallback := CharacterJudgment{
		About:      about,
		Name:       name,
		Impression: "No strong impression yet.",
		Trust:      "medium",
		Interest:   "medium",
		Threat:     "none",
	}

	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return fallback
	}

	var j judgmentJSON
	if err := json.Unmarshal([]byte(raw[start:end+1]), &j); err != nil {
		return fallback
	}

	if j.Impression == "" {
		j.Impression = fallback.Impression
	}
	if !validTrust[j.Trust] {
		j.Trust = "medium"
	}
	if !validInterest[j.Interest] {
		j.Interest = "medium"
	}
	if !validThreat[j.Threat] {
		j.Threat = "none"
	}

	return CharacterJudgment{
		About:      about,
		Name:       name,
		Impression: j.Impression,
		Trust:      j.Trust,
		Interest:   j.Interest,
		Threat:     j.Threat,
	}
}

// FormatJudgmentForPrompt renders a judgment as a human-readable block suitable
// for appending to an LLM system prompt. Returns empty string for zero-value judgment.
func FormatJudgmentForPrompt(j CharacterJudgment) string {
	if j.About == "" {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nYour prior impression of %s:\n", j.Name))
	sb.WriteString(fmt.Sprintf("%q\n", j.Impression))
	sb.WriteString(fmt.Sprintf("Trust: %s | Interest: %s | Perceived threat: %s\n", j.Trust, j.Interest, j.Threat))
	return sb.String()
}

// llmCaller is a minimal interface so judgment.go doesn't import the full engine.
type llmCaller interface {
	Chat(messages []llm.Message, opts ...llm.Option) (string, error)
}

// FormInitialJudgments fires N×(N-1) parallel LLM calls to populate every
// character's Judgments map before tick 1. Errors for individual pairs are
// logged and replaced with a neutral fallback judgment (non-fatal).
func FormInitialJudgments(ctx context.Context, chars []*Character, client llmCaller, language string) {
	if len(chars) < 2 {
		return
	}

	type work struct {
		judge  *Character
		target *Character
	}

	jobs := make([]work, 0, len(chars)*(len(chars)-1))
	for _, judge := range chars {
		for _, target := range chars {
			if judge.ID != target.ID {
				jobs = append(jobs, work{judge, target})
			}
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		job := job
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}

			snapshot := ObservableSnapshot(*job.target)
			prompt := BuildJudgmentPrompt(*job.judge, snapshot, language)
			raw, err := client.Chat([]llm.Message{
				{Role: "user", Content: prompt},
			})

			var judgment CharacterJudgment
			if err != nil {
				log.Printf("[judgment] %s→%s: LLM error (using fallback): %v", job.judge.ID, job.target.ID, err)
				judgment = CharacterJudgment{
					About: job.target.ID, Name: snapshot.Name,
					Impression: "No strong impression yet.", Trust: "medium", Interest: "medium", Threat: "none",
				}
			} else {
				judgment = ParseJudgmentResponse(raw, job.target.ID, snapshot.Name)
			}

			mu.Lock()
			job.judge.Judgments[job.target.ID] = judgment
			mu.Unlock()
		}()
	}

	wg.Wait()
	log.Printf("[judgment] initial judgments formed for %d characters (%d pairs)", len(chars), len(jobs))
}

// FormJudgmentsForNew forms judgments from newChar toward all existing characters.
// Used when the director spawns a new character mid-simulation.
func FormJudgmentsForNew(ctx context.Context, newChar *Character, existing []*Character, client llmCaller, language string) {
	if newChar.Judgments == nil {
		newChar.Judgments = make(map[string]CharacterJudgment)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(existing))

	for _, target := range existing {
		target := target
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}

			snapshot := ObservableSnapshot(*target)
			prompt := BuildJudgmentPrompt(*newChar, snapshot, language)
			raw, err := client.Chat([]llm.Message{{Role: "user", Content: prompt}})

			var judgment CharacterJudgment
			if err != nil {
				log.Printf("[judgment] %s→%s: LLM error (using fallback): %v", newChar.ID, target.ID, err)
				judgment = CharacterJudgment{
					About: target.ID, Name: snapshot.Name,
					Impression: "No strong impression yet.", Trust: "medium", Interest: "medium", Threat: "none",
				}
			} else {
				judgment = ParseJudgmentResponse(raw, target.ID, snapshot.Name)
			}

			mu.Lock()
			newChar.Judgments[target.ID] = judgment
			mu.Unlock()
		}()
	}

	wg.Wait()
}

// FormJudgmentsOfNew forms judgments from all existing characters toward newChar.
// Used when the director spawns a new character mid-simulation.
func FormJudgmentsOfNew(ctx context.Context, existing []*Character, newChar *Character, client llmCaller, language string) {
	snapshot := ObservableSnapshot(*newChar)

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(existing))

	for _, judge := range existing {
		judge := judge
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}

			prompt := BuildJudgmentPrompt(*judge, snapshot, language)
			raw, err := client.Chat([]llm.Message{{Role: "user", Content: prompt}})

			var judgment CharacterJudgment
			if err != nil {
				log.Printf("[judgment] %s→%s: LLM error (using fallback): %v", judge.ID, newChar.ID, err)
				judgment = CharacterJudgment{
					About: newChar.ID, Name: snapshot.Name,
					Impression: "No strong impression yet.", Trust: "medium", Interest: "medium", Threat: "none",
				}
			} else {
				judgment = ParseJudgmentResponse(raw, newChar.ID, snapshot.Name)
			}

			mu.Lock()
			if judge.Judgments == nil {
				judge.Judgments = make(map[string]CharacterJudgment)
			}
			judge.Judgments[newChar.ID] = judgment
			mu.Unlock()
		}()
	}

	wg.Wait()
}

// UpdateJudgment refreshes a single judgment using recent conversation history.
// Called by the engine after every 10 conversations between the same pair.
func UpdateJudgment(ctx context.Context, judge *Character, target *Character, recentHistory []string, tick int, client llmCaller, language string) {
	if judge.Judgments == nil {
		judge.Judgments = make(map[string]CharacterJudgment)
	}

	snapshot := ObservableSnapshot(*target)
	prior := judge.Judgments[target.ID]

	prompt := BuildUpdatePrompt(*judge, snapshot, prior, recentHistory, language)
	raw, err := client.Chat([]llm.Message{{Role: "user", Content: prompt}})
	if err != nil {
		log.Printf("[judgment] update %s→%s: LLM error (keeping prior): %v", judge.ID, target.ID, err)
		return
	}

	updated := ParseJudgmentResponse(raw, target.ID, snapshot.Name)
	updated.FormedTick = prior.FormedTick
	updated.UpdatedTick = tick

	judge.Judgments[target.ID] = updated
}
