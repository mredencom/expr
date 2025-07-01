package vm

import (
	"hash/fnv"
	"sync"
)

// InstructionCache caches frequently executed instruction sequences
type InstructionCache struct {
	cache map[uint64]*CachedSequence
	mutex sync.RWMutex
	stats CacheStats
}

// CachedSequence represents a cached instruction sequence
type CachedSequence struct {
	Instructions []byte
	HitCount     int64
	Size         int
}

// CacheStats tracks cache performance
type CacheStats struct {
	Hits   int64
	Misses int64
	Evictions int64
	Size   int
}

// NewInstructionCache creates a new instruction cache
func NewInstructionCache(maxSize int) *InstructionCache {
	return &InstructionCache{
		cache: make(map[uint64]*CachedSequence),
		stats: CacheStats{},
	}
}

// Get retrieves a cached instruction sequence
func (c *InstructionCache) Get(instructions []byte) (*CachedSequence, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	hash := c.hashInstructions(instructions)
	if seq, exists := c.cache[hash]; exists {
		seq.HitCount++
		c.stats.Hits++
		return seq, true
	}

	c.stats.Misses++
	return nil, false
}

// Put stores an instruction sequence in the cache
func (c *InstructionCache) Put(instructions []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	hash := c.hashInstructions(instructions)
	
	// Create cached sequence
	seq := &CachedSequence{
		Instructions: make([]byte, len(instructions)),
		HitCount:     0,
		Size:         len(instructions),
	}
	copy(seq.Instructions, instructions)

	c.cache[hash] = seq
	c.stats.Size = len(c.cache)
}

// hashInstructions creates a hash of the instruction sequence
func (c *InstructionCache) hashInstructions(instructions []byte) uint64 {
	h := fnv.New64a()
	h.Write(instructions)
	return h.Sum64()
}

// Clear clears the cache
func (c *InstructionCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[uint64]*CachedSequence)
	c.stats = CacheStats{}
}

// GetStats returns cache statistics
func (c *InstructionCache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	statsCopy := c.stats
	statsCopy.Size = len(c.cache)
	return statsCopy
}

// HitRate returns the cache hit rate
func (c *InstructionCache) HitRate() float64 {
	stats := c.GetStats()
	total := stats.Hits + stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(stats.Hits) / float64(total)
} 