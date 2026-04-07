package simulation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/conversation"
	"github.com/jnn-z/city_simulator/internal/scenario"
	"github.com/jnn-z/city_simulator/internal/world"
)

// Config holds everything the engine needs to run.
type Config struct {
	Scenario     scenario.Scenario
	Manager      *conversation.Manager
	Turns        int
	Seed         int64
	OutputWriter io.Writer
}

// Engine drives the simulation tick loop.
type Engine struct {
	cfg       Config
	chars     []*character.Character
	world     *world.State
	scheduler *Scheduler
}

// NewEngine validates config and returns a ready Engine.
func NewEngine(cfg Config) (*Engine, error) {
	chars := make([]*character.Character, len(cfg.Scenario.Characters))
	for i := range cfg.Scenario.Characters {
		chars[i] = &cfg.Scenario.Characters[i]
	}
	if len(chars) < 2 {
		return nil, fmt.Errorf("simulation requires at least 2 characters, got %d", len(chars))
	}
	w := world.NewState(cfg.Scenario.World)
	return &Engine{
		cfg:       cfg,
		chars:     chars,
		world:     w,
		scheduler: NewScheduler(chars, cfg.Seed),
	}, nil
}

// logEntry is written as one JSONL line per tick.
type logEntry struct {
	Tick      int    `json:"tick"`
	Initiator string `json:"initiator"`
	Responder string `json:"responder"`
}

// Run executes the simulation for the configured number of turns.
func (e *Engine) Run(ctx context.Context) error {
	for tick := 1; tick <= e.cfg.Turns; tick++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		pair := e.scheduler.Next()

		result, err := e.cfg.Manager.RunExchange(ctx, pair.Initiator, pair.Responder, e.world, tick)
		if err != nil {
			log.Printf("tick %d: exchange error (skipping): %v", tick, err)
			e.world.AdvanceTick()
			continue
		}

		// Print dialogue to stdout
		fmt.Printf("\n── Tick %d ── %s → %s ──\n", tick, pair.Initiator.Name, pair.Responder.Name)
		fmt.Printf("%s: %s\n", pair.Initiator.Name, result.InitiatorText)
		fmt.Printf("%s: %s\n", pair.Responder.Name, result.ResponderText)

		e.world.AdvanceTick()

		if e.cfg.OutputWriter != nil {
			entry := logEntry{
				Tick:      tick,
				Initiator: pair.Initiator.ID,
				Responder: pair.Responder.ID,
			}
			line, _ := json.Marshal(entry)
			fmt.Fprintf(e.cfg.OutputWriter, "%s\n", line)
		}
	}
	return nil
}
