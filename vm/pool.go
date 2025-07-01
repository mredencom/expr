package vm

import (
	"sync"

	"github.com/mredencom/expr/types"
)

// ValuePool provides object pooling for Value types to reduce allocations
type ValuePool struct {
	intPool    sync.Pool
	floatPool  sync.Pool
	stringPool sync.Pool
	boolPool   sync.Pool

	// Value caches for common values
	intCache    map[int64]*types.IntValue
	floatCache  map[float64]*types.FloatValue
	stringCache map[string]*types.StringValue
	cacheMutex  sync.RWMutex
}

// Enhanced value caching
const (
	maxIntCacheSize    = 1024 // Cache integers from -512 to 511
	maxFloatCacheSize  = 256  // Cache common float values
	maxStringCacheSize = 512  // Cache common strings
)

// VMPool provides object pooling for VM instances to reduce allocations
type VMPool struct {
	pool sync.Pool
}

// Global VM pool instance
var GlobalVMPool = NewVMPool()

// NewVMPool creates a new VM pool
func NewVMPool() *VMPool {
	return &VMPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &VM{
					stack:          make([]types.Value, StackSize),
					sp:             0,
					globals:        make([]types.Value, GlobalsSize),
					pool:           NewValuePool(),
					cache:          NewInstructionCache(1000),
					customBuiltins: make(map[string]interface{}),
					safeJumpTable:  NewSafeJumpTable(),
				}
			},
		},
	}
}

// Get retrieves a VM from the pool
func (p *VMPool) Get() *VM {
	vm := p.pool.Get().(*VM)
	return vm
}

// Put returns a VM to the pool after resetting it
func (p *VMPool) Put(vm *VM) {
	if vm != nil {
		vm.Reset()
		p.pool.Put(vm)
	}
}

// NewValuePool creates a new value pool with caching
func NewValuePool() *ValuePool {
	pool := &ValuePool{
		intPool: sync.Pool{
			New: func() interface{} {
				return &types.IntValue{}
			},
		},
		floatPool: sync.Pool{
			New: func() interface{} {
				return &types.FloatValue{}
			},
		},
		stringPool: sync.Pool{
			New: func() interface{} {
				return &types.StringValue{}
			},
		},
		boolPool: sync.Pool{
			New: func() interface{} {
				return &types.BoolValue{}
			},
		},
		intCache:    make(map[int64]*types.IntValue),
		floatCache:  make(map[float64]*types.FloatValue),
		stringCache: make(map[string]*types.StringValue),
	}

	// Pre-populate common values
	pool.prePopulateCache()
	return pool
}

// prePopulateCache populates the cache with common values
func (p *ValuePool) prePopulateCache() {
	// Cache common integers (-512 to 511)
	for i := int64(-512); i <= 511; i++ {
		p.intCache[i] = types.NewInt(i)
	}

	// Cache common floats
	commonFloats := []float64{0.0, 1.0, -1.0, 0.5, -0.5, 2.0, -2.0, 10.0, -10.0, 100.0, -100.0}
	for _, f := range commonFloats {
		p.floatCache[f] = types.NewFloat(f)
	}

	// Cache common strings
	commonStrings := []string{"", " ", "0", "1", "true", "false", "null", "nil"}
	for _, s := range commonStrings {
		p.stringCache[s] = types.NewString(s)
	}
}

// GetInt gets an int value from cache or creates new one - optimized for performance
func (p *ValuePool) GetInt(value int64) *types.IntValue {
	// Only use cache for very common small values, no locking
	if value >= -512 && value <= 511 {
		p.cacheMutex.RLock()
		cached := p.intCache[value] // This will always exist due to prePopulateCache
		p.cacheMutex.RUnlock()
		return cached
	}
	// For other values, direct creation is faster than cache lookup
	return types.NewInt(value)
}

// PutInt returns an int value to the pool (no-op for cached values)
func (p *ValuePool) PutInt(v *types.IntValue) {
	// No-op for immutable cached values
}

// GetFloat gets a float value from cache or creates new one - simplified for performance
func (p *ValuePool) GetFloat(value float64) *types.FloatValue {
	// Only check cache for very common values (0.0, 1.0, -1.0)
	if value == 0.0 || value == 1.0 || value == -1.0 {
		p.cacheMutex.RLock()
		if cached, exists := p.floatCache[value]; exists {
			p.cacheMutex.RUnlock()
			return cached
		}
		p.cacheMutex.RUnlock()
	}
	// Direct creation for all other values
	return types.NewFloat(value)
}

// PutFloat returns a float value to the pool
func (p *ValuePool) PutFloat(v *types.FloatValue) {
	// No-op for immutable values
}

// GetString gets a string value from cache or creates new one - performance focused
func (p *ValuePool) GetString(value string) *types.StringValue {
	// Only cache very common short strings, and skip cache write for performance
	if len(value) <= 8 {
		p.cacheMutex.RLock()
		if cached, exists := p.stringCache[value]; exists {
			p.cacheMutex.RUnlock()
			return cached
		}
		p.cacheMutex.RUnlock()
	}
	// Direct creation is faster than managing cache for most strings
	return types.NewString(value)
}

// PutString returns a string value to the pool
func (p *ValuePool) PutString(v *types.StringValue) {
	// No-op for immutable values
}

// GetBool gets a bool value (always cached)
func (p *ValuePool) GetBool(value bool) *types.BoolValue {
	if value {
		return types.NewBool(true)
	}
	return types.NewBool(false)
}

// PutBool returns a bool value to the pool
func (p *ValuePool) PutBool(v *types.BoolValue) {
	// No-op for immutable values
}

// ClearCache clears the value caches (for memory management)
func (p *ValuePool) ClearCache() {
	p.cacheMutex.Lock()
	defer p.cacheMutex.Unlock()

	// Clear caches but keep pre-populated common values
	p.stringCache = make(map[string]*types.StringValue)
	p.floatCache = make(map[float64]*types.FloatValue)

	// Re-populate with common values
	p.prePopulateCache()
}

// Stats returns pool statistics
type PoolStats struct {
	IntCacheSize    int
	FloatCacheSize  int
	StringCacheSize int
	IntCacheHits    int64 // Would need atomic counters in practice
	FloatCacheHits  int64
	StringCacheHits int64
}

// GetStats returns current pool statistics
func (p *ValuePool) GetStats() PoolStats {
	p.cacheMutex.RLock()
	defer p.cacheMutex.RUnlock()

	return PoolStats{
		IntCacheSize:    len(p.intCache),
		FloatCacheSize:  len(p.floatCache),
		StringCacheSize: len(p.stringCache),
		// Cache hits would need atomic counters for thread safety
	}
}
