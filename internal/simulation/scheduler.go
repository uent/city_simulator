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

// Scheduler returns character pairs in a round-robin order,
// optionally shuffled with a seed.
type Scheduler struct {
	pairs []Pair
	index int
}

// NewScheduler builds all unique pairs from the character list.
// A non-zero seed shuffles the order.
func NewScheduler(characters []*character.Character, seed int64) *Scheduler {
	var pairs []Pair
	for i := 0; i < len(characters); i++ {
		for j := i + 1; j < len(characters); j++ {
			pairs = append(pairs, Pair{Initiator: characters[i], Responder: characters[j]})
		}
	}
	if seed != 0 {
		r := rand.New(rand.NewSource(seed))
		r.Shuffle(len(pairs), func(i, k int) { pairs[i], pairs[k] = pairs[k], pairs[i] })
	}
	return &Scheduler{pairs: pairs}
}

// Next returns the next pair, cycling indefinitely.
func (s *Scheduler) Next() Pair {
	p := s.pairs[s.index%len(s.pairs)]
	s.index++
	return p
}
