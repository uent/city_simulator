package simulation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/director"
	"github.com/jnn-z/city_simulator/internal/llm"
	"github.com/jnn-z/city_simulator/internal/messaging"
	"github.com/jnn-z/city_simulator/internal/scenario"
	"github.com/jnn-z/city_simulator/internal/summary"
	"github.com/jnn-z/city_simulator/internal/world"
)

// Config holds everything the engine needs to run.
type Config struct {
	Scenario     scenario.Scenario
	LLMClient    *llm.Client        // used only by the Game Director; character actors own their own client ref
	Bus          *messaging.MessageBus
	Turns        int
	Seed         int64
	OutputWriter io.Writer
	Language     string
}

// Engine drives the simulation tick loop.
type Engine struct {
	cfg             Config
	chars           []*character.Character
	directorChar    *character.Character // nil if scenario has no Game Director
	registry        *director.Registry
	world           *world.State
	scheduler       *Scheduler
	locationNames   []string
	registeredChars map[string]bool // tracks which character IDs have bus actors
}

// NewEngine validates config and returns a ready Engine.
// The Bus must have one actor registered per character in cfg.Scenario.Characters.
func NewEngine(cfg Config) (*Engine, error) {
	chars := make([]*character.Character, len(cfg.Scenario.Characters))
	for i := range cfg.Scenario.Characters {
		chars[i] = &cfg.Scenario.Characters[i]
	}
	if len(chars) < 2 {
		return nil, fmt.Errorf("simulation requires at least 2 characters, got %d", len(chars))
	}
	w := world.NewState(cfg.Scenario.World)

	locationNames := make([]string, len(cfg.Scenario.World.Locations))
	for i, loc := range cfg.Scenario.World.Locations {
		locationNames[i] = loc.Name
	}

	registered := make(map[string]bool, len(chars))
	for _, c := range chars {
		registered[c.ID] = true
	}

	return &Engine{
		cfg:             cfg,
		chars:           chars,
		directorChar:    cfg.Scenario.GameDirector,
		registry:        director.NewRegistry(),
		world:           w,
		scheduler:       NewScheduler(chars, locationNames, cfg.Seed),
		locationNames:   locationNames,
		registeredChars: registered,
	}, nil
}

// registerSpawnedChars creates bus actors for any characters in e.chars that
// don't yet have a registered actor, starts them with ctx, and adds pairs to
// the scheduler. Should be called after the director runs each tick.
func (e *Engine) registerSpawnedChars(ctx context.Context) {
	for _, c := range e.chars {
		if e.registeredChars[c.ID] {
			continue
		}
		actor := character.NewCharacterActor(c, e.cfg.LLMClient)
		e.cfg.Bus.Register(actor)
		actor.Start(ctx)
		// Add pairs between this character and all already-known characters.
		known := make([]*character.Character, 0, len(e.chars)-1)
		for _, existing := range e.chars {
			if existing.ID != c.ID && e.registeredChars[existing.ID] {
				known = append(known, existing)
			}
		}
		e.scheduler.AddCharacter(c, known, e.locationNames)
		e.registeredChars[c.ID] = true
		log.Printf("[spawn] created character %s (%s)", c.ID, c.Name)
	}
}

// runDirector invokes the Game Director for the current tick using the tool-use
// dispatch loop. Errors per action are logged and skipped (fail-open).
func (e *Engine) runDirector(ctx context.Context, tick int) {
	systemPrompt := director.BuildDirectorPrompt(e.world, e.chars, tick, e.cfg.Language)
	raw, err := e.cfg.LLMClient.Chat([]llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("Generate world events for tick %d.", tick)},
	})
	if err != nil {
		log.Printf("tick %d: game director LLM error (skipping): %v", tick, err)
		return
	}

	calls, _ := director.ParseToolCalls(raw)
	for _, call := range calls {
		if err := e.registry.Dispatch(call.Name, call.Args, e.world, &e.chars); err != nil {
			log.Printf("tick %d: director action %q error (skipping): %v", tick, call.Name, err)
			continue
		}
		fmt.Printf("  [Director] %s\n", e.registry.Summarize(call.Name, call.Args))
	}
}

// generateAndSaveSummary produces a narrative summary of the completed simulation
// and writes it to a timestamped file. Errors are logged and suppressed (fail-open).
func (e *Engine) generateAndSaveSummary(ctx context.Context) {
	text, err := summary.GenerateSummary(ctx, e.cfg.LLMClient, e.world, e.chars, e.cfg.Scenario, e.cfg.Language)
	if err != nil {
		log.Printf("summary: generation failed (skipping): %v", err)
		return
	}
	path, err := summary.SaveSummary(e.cfg.Scenario.Name, text)
	if err != nil {
		log.Printf("summary: save failed (skipping): %v", err)
		return
	}
	fmt.Printf("\n── Summary saved: %s ──\n", path)
}

// logEntry is written as one JSONL line per tick.
type logEntry struct {
	Tick              int    `json:"tick"`
	Initiator         string `json:"initiator"`
	InitiatorLocation string `json:"initiator_location"`
	InitiatorSpeech   string `json:"initiator_speech"`
	InitiatorAction   string `json:"initiator_action"`
	Responder         string `json:"responder"`
	ResponderLocation string `json:"responder_location"`
	ResponderSpeech   string `json:"responder_speech"`
	ResponderAction   string `json:"responder_action"`
}

// Run executes the simulation for the configured number of turns.
func (e *Engine) Run(ctx context.Context) error {
	// Start all character actors; they stop when ctx is cancelled.
	e.cfg.Bus.StartAll(ctx)

	for tick := 1; tick <= e.cfg.Turns; tick++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Game Director generates autonomous world events before characters act.
		if e.directorChar != nil {
			e.runDirector(ctx, tick)

			// Register any characters spawned by the director this tick.
			e.registerSpawnedChars(ctx)

			// Broadcast DirectorDirective so actors acknowledge tick start.
			replies, err := e.cfg.Bus.Broadcast(messaging.NewRequest(
				messaging.DirectorDirective, "director", "", tick, nil,
			))
			if err == nil {
				for range replies {
				} // drain
			}
		}

		pair := e.scheduler.Next()

		// Compute zone roster snapshot once per tick.
		zoneRoster := character.BuildZoneRoster(e.chars)

		// Build world contexts for this exchange.
		initiatorPeers := peersAt(zoneRoster, pair.Initiator.Location, pair.Initiator.Name)
		initiatorWorldCtx := e.world.PublicSummary()
		if local := e.world.LocalContext(pair.Initiator.Location, initiatorPeers); local != "" {
			initiatorWorldCtx += "\n" + local
		}
		initiatorSystem := character.BuildSystemPrompt(*pair.Initiator, e.cfg.Language) +
			character.FlushInbox(pair.Initiator) +
			"\n\nWorld context:\n" + initiatorWorldCtx +
			character.BuildZoneContext(zoneRoster)

		responderPeers := peersAt(zoneRoster, pair.Responder.Location, pair.Responder.Name)
		responderWorldCtx := e.world.PublicSummary()
		if local := e.world.LocalContext(pair.Responder.Location, responderPeers); local != "" {
			responderWorldCtx += "\n" + local
		}
		responderSystem := character.BuildSystemPrompt(*pair.Responder, e.cfg.Language) +
			character.FlushInbox(pair.Responder) +
			"\n\nWorld context:\n" + responderWorldCtx +
			character.BuildZoneContext(zoneRoster)

		// Send CharChat to the responder actor; it generates both sides of the exchange.
		chatMsg := messaging.NewRequest(
			messaging.CharChat,
			pair.Initiator.ID,
			pair.Responder.ID,
			tick,
			messaging.CharChatPayload{
				InitiatorID:     pair.Initiator.ID,
				InitiatorName:   pair.Initiator.Name,
				InitiatorSystem: initiatorSystem,
				ResponderSystem: responderSystem,
			},
		)
		if err := e.cfg.Bus.Send(chatMsg); err != nil {
			log.Printf("tick %d: send CharChat: %v (skipping)", tick, err)
			e.world.AdvanceTick()
			continue
		}
		chatReply := <-chatMsg.ReplyChan
		result, ok := chatReply.Payload.(messaging.CharChatReply)
		if !ok || result.Err != nil {
			if ok {
				log.Printf("tick %d: exchange error (skipping): %v", tick, result.Err)
			} else {
				log.Printf("tick %d: unexpected reply payload (skipping)", tick)
			}
			e.world.AdvanceTick()
			continue
		}

		// Print dialogue.
		fmt.Printf("\n── Tick %d ── %s [%s] → %s [%s] ──\n",
			tick,
			pair.Initiator.Name, pair.Initiator.Location,
			pair.Responder.Name, pair.Responder.Location,
		)
		if result.InitiatorAction != "" {
			fmt.Printf("*%s*\n", result.InitiatorAction)
		}
		fmt.Printf("%s: %s\n", pair.Initiator.Name, result.InitiatorSpeech)
		if result.ResponderAction != "" {
			fmt.Printf("*%s*\n", result.ResponderAction)
		}
		fmt.Printf("%s: %s\n", pair.Responder.Name, result.ResponderSpeech)

		// Append conversation event before advancing tick.
		e.world.AppendEvent(world.Event{
			Tick:         tick,
			Type:         "conversation",
			Description:  fmt.Sprintf("%s spoke to %s", pair.Initiator.Name, pair.Responder.Name),
			Participants: []string{pair.Initiator.ID, pair.Responder.ID},
		})
		e.world.AdvanceTick()

		// Ask both characters where to move next.
		for _, c := range []*character.Character{pair.Initiator, pair.Responder} {
			moveMsg := messaging.NewRequest(
				messaging.MoveDecision,
				"engine",
				c.ID,
				tick,
				messaging.MoveDecisionPayload{
					SystemPrompt: character.BuildMovementPrompt(*c, e.locationNames, zoneRoster, e.cfg.Language),
					Locations:    e.locationNames,
				},
			)
			if err := e.cfg.Bus.Send(moveMsg); err != nil {
				log.Printf("tick %d: send MoveDecision to %s: %v", tick, c.ID, err)
				continue
			}
			moveReply := <-moveMsg.ReplyChan
			if move, ok := moveReply.Payload.(messaging.MoveDecisionReply); ok {
				if move.Location != "stay" && move.Location != c.Location {
					prev := c.Location
					c.Location = move.Location
					fmt.Printf("  → %s moves from %s to %s\n", c.Name, prev, c.Location)
				}
			}
		}

		if e.cfg.OutputWriter != nil {
			entry := logEntry{
				Tick:              tick,
				Initiator:         pair.Initiator.ID,
				InitiatorLocation: pair.Initiator.Location,
				InitiatorSpeech:   result.InitiatorSpeech,
				InitiatorAction:   result.InitiatorAction,
				Responder:         pair.Responder.ID,
				ResponderLocation: pair.Responder.Location,
				ResponderSpeech:   result.ResponderSpeech,
				ResponderAction:   result.ResponderAction,
			}
			line, _ := json.Marshal(entry)
			fmt.Fprintf(e.cfg.OutputWriter, "%s\n", line)
		}
	}

	if ctx.Err() == nil {
		e.generateAndSaveSummary(ctx)
	}
	return nil
}

// peersAt returns the names of characters at locationID from roster, excluding selfName.
func peersAt(roster map[string][]string, locationID, selfName string) []string {
	all := roster[locationID]
	if len(all) == 0 {
		return nil
	}
	peers := make([]string, 0, len(all))
	for _, name := range all {
		if name != selfName {
			peers = append(peers, name)
		}
	}
	return peers
}
