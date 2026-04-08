package scenario

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
	"gopkg.in/yaml.v3"
)

// RuntimeOverrides holds optional per-scenario config from scenario.yaml.
// All fields are pointers; nil means "not set by this scenario".
// Model is intentionally excluded: model selection belongs to the user's
// environment (.env / OLLAMA_MODEL), not to individual scenarios.
type RuntimeOverrides struct {
	Turns  *int    `yaml:"turns"`
	Seed   *int64  `yaml:"seed"`
	Output *string `yaml:"output"`
}

// SimConfig is the final resolved runtime configuration after merging all sources.
type SimConfig struct {
	Model  string
	Turns  int
	Seed   int64
	Output string
}

// CLIFlags captures values explicitly provided via CLI flags.
// Nil means the flag was not set by the user (default was used).
type CLIFlags struct {
	Model  *string
	Turns  *int
	Seed   *int64
	Output *string
}

// Scenario bundles everything needed to run one simulation.
type Scenario struct {
	Name         string
	Dir          string
	Characters   []character.Character   // regular characters only (Type != "game_director")
	GameDirector *character.Character    // nil if not defined in this scenario
	World        world.WorldConfig
	Overrides    RuntimeOverrides
}

// Load resolves dirOrName to a scenario directory and reads all config files.
// If dirOrName is an absolute path it is used directly. If it contains no path
// separators it is looked up under simulations/<dirOrName> relative to the
// current working directory.
func Load(dirOrName string) (Scenario, error) {
	dir := dirOrName
	if !filepath.IsAbs(dirOrName) && !strings.ContainsAny(dirOrName, "/\\") {
		dir = filepath.Join("simulations", dirOrName)
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return Scenario{}, fmt.Errorf("scenario %q not found at %q", dirOrName, dir)
	}

	sc := Scenario{
		Name: filepath.Base(dir),
		Dir:  dir,
	}

	// Load characters.yaml (required)
	charsPath := filepath.Join(dir, "characters.yaml")
	if _, err := os.Stat(charsPath); err != nil {
		return Scenario{}, fmt.Errorf("scenario %q: missing required file \"characters.yaml\" in %s", dirOrName, dir)
	}
	chars, err := character.LoadCharacters(charsPath)
	if err != nil {
		return Scenario{}, fmt.Errorf("scenario %q: %w", dirOrName, err)
	}

	// Separate Game Director(s) from regular characters.
	var regular []character.Character
	directorCount := 0
	for i := range chars {
		if chars[i].Type == "game_director" {
			directorCount++
			if directorCount == 1 {
				c := chars[i]
				sc.GameDirector = &c
			} else if directorCount == 2 {
				log.Printf("scenario %q: more than one game_director entry found; using the first", dirOrName)
			}
		} else {
			regular = append(regular, chars[i])
		}
	}
	sc.Characters = regular

	// Load world.yaml (required)
	worldPath := filepath.Join(dir, "world.yaml")
	if _, err := os.Stat(worldPath); err != nil {
		return Scenario{}, fmt.Errorf("scenario %q: missing required file \"world.yaml\" in %s", dirOrName, dir)
	}
	wdata, err := os.ReadFile(worldPath)
	if err != nil {
		return Scenario{}, fmt.Errorf("scenario %q: cannot read world.yaml: %w", dirOrName, err)
	}
	if err := yaml.Unmarshal(wdata, &sc.World); err != nil {
		return Scenario{}, fmt.Errorf("scenario %q: cannot parse world.yaml: %w", dirOrName, err)
	}

	// Load scenario.yaml (optional)
	scenarioPath := filepath.Join(dir, "scenario.yaml")
	if data, err := os.ReadFile(scenarioPath); err == nil {
		if err := yaml.Unmarshal(data, &sc.Overrides); err != nil {
			return Scenario{}, fmt.Errorf("scenario %q: cannot parse scenario.yaml: %w", dirOrName, err)
		}
	}

	return sc, nil
}

// MergeConfig produces a final SimConfig applying priority: CLI flags > scenario.yaml overrides > defaults.
// Model is not overridable by scenarios; it always comes from env vars or CLI flags.
func MergeConfig(overrides RuntimeOverrides, flags CLIFlags, defaults SimConfig) SimConfig {
	cfg := defaults

	if overrides.Turns != nil {
		cfg.Turns = *overrides.Turns
	}
	if overrides.Seed != nil {
		cfg.Seed = *overrides.Seed
	}
	if overrides.Output != nil {
		cfg.Output = *overrides.Output
	}

	// CLI flags override scenario.yaml
	if flags.Model != nil {
		cfg.Model = *flags.Model
	}
	if flags.Turns != nil {
		cfg.Turns = *flags.Turns
	}
	if flags.Seed != nil {
		cfg.Seed = *flags.Seed
	}
	if flags.Output != nil {
		cfg.Output = *flags.Output
	}

	return cfg
}
