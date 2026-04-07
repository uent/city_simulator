package character

// MemoryEntry records a single utterance in a character's memory.
type MemoryEntry struct {
	Tick    int
	Speaker string
	Text    string
}

// AddMemory appends an entry, evicting the oldest if at capacity.
func (c *Character) AddMemory(entry MemoryEntry) {
	if c.MaxMemory <= 0 {
		c.MaxMemory = 20
	}
	if len(c.Memory) >= c.MaxMemory {
		c.Memory = c.Memory[1:]
	}
	c.Memory = append(c.Memory, entry)
}

// RecentMemory returns up to the last n entries in chronological order.
func (c *Character) RecentMemory(n int) []MemoryEntry {
	if n <= 0 {
		return nil
	}
	if n >= len(c.Memory) {
		result := make([]MemoryEntry, len(c.Memory))
		copy(result, c.Memory)
		return result
	}
	result := make([]MemoryEntry, n)
	copy(result, c.Memory[len(c.Memory)-n:])
	return result
}
