package vm

import (
	"sync"
	"testing"
)

// TestNewInstructionCache tests cache creation
func TestNewInstructionCache(t *testing.T) {
	cache := NewInstructionCache(100)

	if cache == nil {
		t.Fatal("Expected cache to not be nil")
	}

	if cache.cache == nil {
		t.Error("Expected cache.cache to be initialized")
	}

	stats := cache.GetStats()
	if stats.Hits != 0 {
		t.Errorf("Expected 0 hits, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected 0 misses, got %d", stats.Misses)
	}
	if stats.Size != 0 {
		t.Errorf("Expected 0 size, got %d", stats.Size)
	}
}

// TestInstructionCachePutGet tests basic put and get operations
func TestInstructionCachePutGet(t *testing.T) {
	cache := NewInstructionCache(100)
	instructions := []byte{1, 2, 3, 4, 5}

	// Test cache miss
	seq, hit := cache.Get(instructions)
	if hit {
		t.Error("Expected cache miss for new instructions")
	}
	if seq != nil {
		t.Error("Expected nil sequence for cache miss")
	}

	stats := cache.GetStats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	// Put instructions in cache
	cache.Put(instructions)

	stats = cache.GetStats()
	if stats.Size != 1 {
		t.Errorf("Expected size 1, got %d", stats.Size)
	}

	// Test cache hit
	seq, hit = cache.Get(instructions)
	if !hit {
		t.Error("Expected cache hit for stored instructions")
	}
	if seq == nil {
		t.Fatal("Expected non-nil sequence for cache hit")
	}
	if seq.Size != len(instructions) {
		t.Errorf("Expected size %d, got %d", len(instructions), seq.Size)
	}
	if seq.HitCount != 1 {
		t.Errorf("Expected hit count 1, got %d", seq.HitCount)
	}

	stats = cache.GetStats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
}

// TestInstructionCacheMultipleOperations tests multiple cache operations
func TestInstructionCacheMultipleOperations(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions1 := []byte{1, 2, 3}
	instructions2 := []byte{4, 5, 6}
	instructions3 := []byte{1, 2, 3} // Same as instructions1

	// Add first instruction sequence
	cache.Put(instructions1)
	cache.Put(instructions2)

	// Test that cache contains both sequences
	stats := cache.GetStats()
	if stats.Size != 2 {
		t.Errorf("Expected size 2, got %d", stats.Size)
	}

	// Test cache hits
	seq1, hit1 := cache.Get(instructions1)
	seq2, hit2 := cache.Get(instructions2)
	seq3, hit3 := cache.Get(instructions3) // Should hit same as instructions1

	if !hit1 || !hit2 || !hit3 {
		t.Error("Expected all cache operations to be hits")
	}

	if seq1 != seq3 {
		t.Error("Expected same instruction sequences to return same cached object")
	}

	if seq1.HitCount != 2 { // Hit by both instructions1 and instructions3
		t.Errorf("Expected hit count 2 for seq1, got %d", seq1.HitCount)
	}

	if seq2.HitCount != 1 {
		t.Errorf("Expected hit count 1 for seq2, got %d", seq2.HitCount)
	}

	stats = cache.GetStats()
	if stats.Hits != 3 {
		t.Errorf("Expected 3 hits, got %d", stats.Hits)
	}
}

// TestInstructionCacheHashing tests that different instructions have different hashes
func TestInstructionCacheHashing(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions1 := []byte{1, 2, 3}
	instructions2 := []byte{3, 2, 1}

	hash1 := cache.hashInstructions(instructions1)
	hash2 := cache.hashInstructions(instructions2)

	if hash1 == hash2 {
		t.Error("Expected different instructions to have different hashes")
	}

	// Test same instructions have same hash
	hash3 := cache.hashInstructions(instructions1)
	if hash1 != hash3 {
		t.Error("Expected same instructions to have same hash")
	}
}

// TestInstructionCacheClear tests cache clearing
func TestInstructionCacheClear(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions := []byte{1, 2, 3, 4, 5}
	cache.Put(instructions)

	// Verify cache has content
	stats := cache.GetStats()
	if stats.Size != 1 {
		t.Errorf("Expected size 1 before clear, got %d", stats.Size)
	}

	seq, hit := cache.Get(instructions)
	if !hit {
		t.Error("Expected cache hit before clear")
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	stats = cache.GetStats()
	if stats.Size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", stats.Size)
	}
	if stats.Hits != 0 {
		t.Errorf("Expected 0 hits after clear, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected 0 misses after clear, got %d", stats.Misses)
	}

	// Verify cache miss after clear
	seq, hit = cache.Get(instructions)
	if hit {
		t.Error("Expected cache miss after clear")
	}
	if seq != nil {
		t.Error("Expected nil sequence after clear")
	}
}

// TestInstructionCacheStats tests statistics tracking
func TestInstructionCacheStats(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions1 := []byte{1, 2, 3}
	instructions2 := []byte{4, 5, 6}

	// Initial stats
	stats := cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Size != 0 {
		t.Error("Expected initial stats to be zero")
	}

	// Generate misses
	cache.Get(instructions1) // Miss 1
	cache.Get(instructions2) // Miss 2

	stats = cache.GetStats()
	if stats.Misses != 2 {
		t.Errorf("Expected 2 misses, got %d", stats.Misses)
	}
	if stats.Hits != 0 {
		t.Errorf("Expected 0 hits, got %d", stats.Hits)
	}

	// Add to cache
	cache.Put(instructions1)
	cache.Put(instructions2)

	stats = cache.GetStats()
	if stats.Size != 2 {
		t.Errorf("Expected size 2, got %d", stats.Size)
	}

	// Generate hits
	cache.Get(instructions1) // Hit 1
	cache.Get(instructions1) // Hit 2
	cache.Get(instructions2) // Hit 3

	stats = cache.GetStats()
	if stats.Hits != 3 {
		t.Errorf("Expected 3 hits, got %d", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("Expected 2 misses, got %d", stats.Misses)
	}
}

// TestInstructionCacheHitRate tests hit rate calculation
func TestInstructionCacheHitRate(t *testing.T) {
	cache := NewInstructionCache(100)

	// Test hit rate with no operations
	hitRate := cache.HitRate()
	if hitRate != 0.0 {
		t.Errorf("Expected hit rate 0.0 for empty cache, got %f", hitRate)
	}

	instructions := []byte{1, 2, 3}

	// Generate some misses and hits
	cache.Get(instructions) // Miss
	cache.Put(instructions)
	cache.Get(instructions) // Hit
	cache.Get(instructions) // Hit

	hitRate = cache.HitRate()
	expected := float64(2) / float64(3) // 2 hits out of 3 total operations
	if hitRate != expected {
		t.Errorf("Expected hit rate %f, got %f", expected, hitRate)
	}
}

// TestInstructionCacheConcurrency tests concurrent access
func TestInstructionCacheConcurrency(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions := []byte{1, 2, 3, 4, 5}

	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := 100

	// Concurrent puts and gets
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Create unique instructions for this goroutine
			localInstructions := make([]byte, len(instructions)+1)
			copy(localInstructions, instructions)
			localInstructions[len(instructions)] = byte(id)

			for j := 0; j < operationsPerGoroutine; j++ {
				if j%2 == 0 {
					cache.Put(localInstructions)
				} else {
					cache.Get(localInstructions)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional
	stats := cache.GetStats()
	if stats.Size == 0 {
		t.Error("Expected cache to have some entries after concurrent operations")
	}

	// Test that we can still use the cache
	cache.Put(instructions)
	seq, hit := cache.Get(instructions)
	if !hit {
		t.Error("Expected cache hit after concurrent operations")
	}
	if seq == nil {
		t.Error("Expected non-nil sequence after concurrent operations")
	}
}

// TestCachedSequenceOperations tests CachedSequence functionality
func TestCachedSequenceOperations(t *testing.T) {
	cache := NewInstructionCache(100)

	instructions := []byte{1, 2, 3, 4, 5}
	cache.Put(instructions)

	seq, hit := cache.Get(instructions)
	if !hit {
		t.Fatal("Expected cache hit")
	}

	// Test initial values
	if seq.HitCount != 1 {
		t.Errorf("Expected initial hit count 1, got %d", seq.HitCount)
	}
	if seq.Size != len(instructions) {
		t.Errorf("Expected size %d, got %d", len(instructions), seq.Size)
	}
	if len(seq.Instructions) != len(instructions) {
		t.Errorf("Expected instructions length %d, got %d", len(instructions), len(seq.Instructions))
	}

	// Verify instructions are copied correctly
	for i, b := range instructions {
		if seq.Instructions[i] != b {
			t.Errorf("Expected instruction[%d] = %d, got %d", i, b, seq.Instructions[i])
		}
	}

	// Test hit count increment
	cache.Get(instructions)
	if seq.HitCount != 2 {
		t.Errorf("Expected hit count 2 after second get, got %d", seq.HitCount)
	}
}

// TestInstructionCacheEdgeCases tests edge cases
func TestInstructionCacheEdgeCases(t *testing.T) {
	cache := NewInstructionCache(100)

	t.Run("EmptyInstructions", func(t *testing.T) {
		emptyInstructions := []byte{}
		cache.Put(emptyInstructions)

		seq, hit := cache.Get(emptyInstructions)
		if !hit {
			t.Error("Expected cache hit for empty instructions")
		}
		if seq.Size != 0 {
			t.Errorf("Expected size 0 for empty instructions, got %d", seq.Size)
		}
	})

	t.Run("LargeInstructions", func(t *testing.T) {
		largeInstructions := make([]byte, 10000)
		for i := range largeInstructions {
			largeInstructions[i] = byte(i % 256)
		}

		cache.Put(largeInstructions)
		seq, hit := cache.Get(largeInstructions)
		if !hit {
			t.Error("Expected cache hit for large instructions")
		}
		if seq.Size != len(largeInstructions) {
			t.Errorf("Expected size %d for large instructions, got %d", len(largeInstructions), seq.Size)
		}
	})

	t.Run("NilInstructions", func(t *testing.T) {
		// This tests behavior with nil slice
		var nilInstructions []byte = nil
		cache.Put(nilInstructions)

		seq, hit := cache.Get(nilInstructions)
		if !hit {
			t.Error("Expected cache hit for nil instructions")
		}
		if seq.Size != 0 {
			t.Errorf("Expected size 0 for nil instructions, got %d", seq.Size)
		}
	})
}
