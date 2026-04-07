package scenario

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jnn-z/city_simulator/internal/character"
	"github.com/jnn-z/city_simulator/internal/world"
	"gopkg.in/yaml.v3"
)

// RuntimeOverrides holds optional per-scenario config from scenario.yaml.
// All fields are pointers; nil means "not set by this scenario".
type RuntimeOverrides struct {
	Model  *string `yaml:"model"`
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
	Name       string
	Dir        string
	Characters []character.Character
	World      world.WorldConfig
	Overrides  RuntimeOverrides
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
	sc.Characters = chars

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
func MergeConfig(overrides RuntimeOverrides, flags CLIFlags, defaults SimConfig) SimConfig {
	cfg := defaults

	if overrides.Model != nil {
		cfg.Model = *overrides.Model
	}
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
