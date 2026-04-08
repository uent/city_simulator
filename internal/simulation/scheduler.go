package simulation

import (
	"math/rand"

	"github.com/jnn-z/city_simulator/internal/character"
)

// Pair is an ordered (initiator, responder) pairing.
type Pair struct {
	Initiator *character.Character
	Responder *character.Character
}

// Scheduler manages character pairing and initial location assignment.
// Movement decisions after each exchange are delegated to the LLM via Manager.
type Scheduler struct {
	pairs []Pair
	index int
	rng   *rand.Rand
}

// NewScheduler builds all unique character pairs, assigns random initial
// locations to all characters, and shuffles pair order.
// A zero seed uses a random source.
func NewScheduler(characters []*character.Character, locations []string, seed int64) *Scheduler {
	src := seed
	if src == 0 {
		src = rand.Int63()
	}
	rng := rand.New(rand.NewSource(src))

	var pairs []Pair
	for i := 0; i < len(characters); i++ {
		for j := i + 1; j < len(characters); j++ {
			pairs = append(pairs, Pair{Initiator: characters[i], Responder: characters[j]})
		}
	}
	rng.Shuffle(len(pairs), func(i, k int) { pairs[i], pairs[k] = pairs[k], pairs[i] })

	s := &Scheduler{
		pairs: pairs,
		rng:   rng,
	}

	// Assign a random starting location to every character.
	if len(locations) > 0 {
		for _, c := range characters {
			c.Location = locations[rng.Intn(len(locations))]
		}
	}

	return s
}

// Next returns the next pair, cycling indefinitely.
func (s *Scheduler) Next() Pair {
	p := s.pairs[s.index%len(s.pairs)]
	s.index++
	return p
}

