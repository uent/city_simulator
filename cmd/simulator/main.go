package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jnn-z/city_simulator/internal/conversation"
	"github.com/jnn-z/city_simulator/internal/llm"
	"github.com/jnn-z/city_simulator/internal/scenario"
	"github.com/jnn-z/city_simulator/internal/simulation"
)

func envOrString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("warning: %s=%q is not a valid integer, using default %d", key, v, fallback)
		return fallback
	}
	return n
}

func envOrInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Printf("warning: %s=%q is not a valid integer, using default %d", key, v, fallback)
		return fallback
	}
	return n
}

func main() {
	scenarioFlag := flag.String("scenario", envOrString("SIM_SCENARIO", "default"), "Scenario name under simulations/ or absolute path")
	model := flag.String("model", envOrString("OLLAMA_MODEL", "llama3"), "Ollama model name")
	ollamaURL := flag.String("ollama-url", envOrString("OLLAMA_URL", "http://localhost:11434"), "Ollama base URL")
	turns := flag.Int("turns", envOrInt("SIM_TURNS", 10), "Number of simulation ticks")
	seed := flag.Int64("seed", envOrInt64("SIM_SEED", 0), "Random seed (0 = deterministic round-robin)")
	output := flag.String("output", envOrString("SIM_OUTPUT", "simulation_output.jsonl"), "JSONL output file path")

	// Detect removed flag before parsing so we can give a clear error.
	for _, arg := range os.Args[1:] {
		if arg == "--characters" || arg == "-characters" {
			log.Fatal("--characters has been removed, use --scenario instead")
		}
	}

	flag.Parse()

	// Determine which flags were explicitly set by the user.
	explicitFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { explicitFlags[f.Name] = true })

	cliFlags := scenario.CLIFlags{}
	if explicitFlags["model"] {
		cliFlags.Model = model
	}
	if explicitFlags["turns"] {
		cliFlags.Turns = turns
	}
	if explicitFlags["seed"] {
		cliFlags.Seed = seed
	}
	if explicitFlags["output"] {
		cliFlags.Output = output
	}

	// Load scenario.
	sc, err := scenario.Load(*scenarioFlag)
	if err != nil {
		log.Fatalf("load scenario: %v", err)
	}

	// Merge config: CLI flags > scenario.yaml overrides > compiled defaults.
	defaults := scenario.SimConfig{
		Model:  *model,
		Turns:  *turns,
		Seed:   *seed,
		Output: *output,
	}
	simCfg := scenario.MergeConfig(sc.Overrides, cliFlags, defaults)

	// Ping Ollama.
	client := llm.NewClient(*ollamaURL, simCfg.Model)
	if err := client.Ping(); err != nil {
		log.Fatalf("cannot connect to Ollama: %v\nMake sure Ollama is running at %s", err, *ollamaURL)
	}
	fmt.Printf("Connected to Ollama at %s (model: %s)\n\n", *ollamaURL, simCfg.Model)

	// Open output file for JSONL log.
	outFile, err := os.Create(simCfg.Output)
	if err != nil {
		log.Fatalf("create output file %s: %v", simCfg.Output, err)
	}
	defer outFile.Close()

	// Wire engine.
	mgr := conversation.NewManager(client)
	engine, err := simulation.NewEngine(simulation.Config{
		Scenario:     sc,
		Manager:      mgr,
		Turns:        simCfg.Turns,
		Seed:         simCfg.Seed,
		OutputWriter: outFile,
	})
	if err != nil {
		log.Fatalf("create engine: %v", err)
	}

	// Handle Ctrl+C gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Printf("Starting simulation: scenario=%s, %d characters, %d turns\n", sc.Name, len(sc.Characters), simCfg.Turns)
	fmt.Printf("Output log: %s\n", simCfg.Output)

	if err := engine.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("simulation error: %v", err)
	}

	fmt.Printf("\nSimulation complete. Output written to %s\n", simCfg.Output)
}
