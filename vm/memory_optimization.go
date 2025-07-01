package vm

import (
	"sync"

	"github.com/mredencom/expr/types"
)

// MemoryOptimizer implements P1 optimizations from PERFORMANCE_SUMMARY.md
// - Memory pre-allocation
// - Instruction buffer pre-allocation
// - String content pool
// - Expression parsing cache
// - Variable lookup cache
type MemoryOptimizer struct {
	// Pre-allocated pools
	stackPool       sync.Pool
	globalsPool     sync.Pool
	instructionPool sync.Pool

	// Caching mechanisms
	exprCache  *ExpressionCache
	varCache   *VariableLookupCache
	stringPool *StringPool

	// Performance metrics
	cacheHits   int64
	cacheMisses int64
	poolHits    int64
	poolMisses  int64
}

// ExpressionCache caches compiled expressions
type ExpressionCache struct {
	cache   map[string]*CachedExpression
	mutex   sync.RWMutex
	maxSize int
}

// CachedExpression represents a cached compiled expression
type CachedExpression struct {
	Instructions []byte
	Constants    []types.Value
	UsageCount   int64
	LastUsed     int64
}

// VariableLookupCache accelerates variable resolution
type VariableLookupCache struct {
	cache map[string]int // variable name -> global index
	mutex sync.RWMutex
}

// StringPool manages string value allocation
type StringPool struct {
	smallStrings  map[string]*types.StringValue // strings <= 64 chars
	mediumStrings map[string]*types.StringValue // strings <= 1024 chars
	mutex         sync.RWMutex
	smallCount    int
	mediumCount   int
	maxSmallSize  int
	maxMediumSize int
}

// Global memory optimizer instance
var GlobalMemoryOptimizer = NewMemoryOptimizer()

// NewMemoryOptimizer creates a new memory optimizer with P1 optimizations
func NewMemoryOptimizer() *MemoryOptimizer {
	return &MemoryOptimizer{
		stackPool: sync.Pool{
			New: func() interface{} {
				return make([]types.Value, StackSize)
			},
		},
		globalsPool: sync.Pool{
			New: func() interface{} {
				return make([]types.Value, GlobalsSize)
			},
		},
		instructionPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024) // Pre-allocate 1KB capacity
			},
		},
		exprCache: &ExpressionCache{
			cache:   make(map[string]*CachedExpression),
			maxSize: 1000, // Cache up to 1000 expressions
		},
		varCache: &VariableLookupCache{
			cache: make(map[string]int),
		},
		stringPool: &StringPool{
			smallStrings:  make(map[string]*types.StringValue),
			mediumStrings: make(map[string]*types.StringValue),
			maxSmallSize:  64,
			maxMediumSize: 1024,
		},
	}
}

// GetOptimizedStack returns a pre-allocated stack
func (mo *MemoryOptimizer) GetOptimizedStack() []types.Value {
	stack := mo.stackPool.Get().([]types.Value)
	// ✅ 延迟清理：仅在放回池时清理使用部分
	// 不在获取时清理，避免巨大开销
	mo.poolHits++
	return stack
}

// PutOptimizedStack returns a stack to the pool
func (mo *MemoryOptimizer) PutOptimizedStack(stack []types.Value) {
	if len(stack) == StackSize {
		// ✅ 智能清理：仅清理可能使用的前256个位置
		// 大多数表达式栈深度 < 256
		clearLimit := 256
		if clearLimit > len(stack) {
			clearLimit = len(stack)
		}
		for i := 0; i < clearLimit; i++ {
			stack[i] = nil
		}
		mo.stackPool.Put(stack)
	}
}

// GetOptimizedGlobals returns pre-allocated globals array
func (mo *MemoryOptimizer) GetOptimizedGlobals() []types.Value {
	globals := mo.globalsPool.Get().([]types.Value)
	// ✅ 延迟清理：仅在放回池时清理使用部分
	mo.poolHits++
	return globals
}

// PutOptimizedGlobals returns globals to the pool
func (mo *MemoryOptimizer) PutOptimizedGlobals(globals []types.Value) {
	if len(globals) == GlobalsSize {
		// ✅ 智能清理：仅清理可能使用的前64个全局变量位置
		// 大多数表达式全局变量 < 64
		clearLimit := 64
		if clearLimit > len(globals) {
			clearLimit = len(globals)
		}
		for i := 0; i < clearLimit; i++ {
			globals[i] = nil
		}
		mo.globalsPool.Put(globals)
	}
}

// GetOptimizedInstructionBuffer returns a pre-allocated instruction buffer
func (mo *MemoryOptimizer) GetOptimizedInstructionBuffer() []byte {
	buffer := mo.instructionPool.Get().([]byte)
	// Reset buffer but keep capacity
	buffer = buffer[:0]
	mo.poolHits++
	return buffer
}

// PutOptimizedInstructionBuffer returns instruction buffer to pool
func (mo *MemoryOptimizer) PutOptimizedInstructionBuffer(buffer []byte) {
	if cap(buffer) >= 1024 && cap(buffer) <= 16384 { // Keep reasonable sized buffers
		mo.instructionPool.Put(buffer)
	}
}

// Expression caching methods

// GetCachedExpression retrieves a cached expression
func (ec *ExpressionCache) GetCachedExpression(key string) (*CachedExpression, bool) {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	if expr, exists := ec.cache[key]; exists {
		expr.UsageCount++
		expr.LastUsed = getCurrentTimestamp()
		return expr, true
	}

	return nil, false
}

// CacheExpression stores a compiled expression
func (ec *ExpressionCache) CacheExpression(key string, instructions []byte, constants []types.Value) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	// Check cache size limit
	if len(ec.cache) >= ec.maxSize {
		ec.evictLRU()
	}

	ec.cache[key] = &CachedExpression{
		Instructions: instructions,
		Constants:    constants,
		UsageCount:   1,
		LastUsed:     getCurrentTimestamp(),
	}
}

// evictLRU removes least recently used entries
func (ec *ExpressionCache) evictLRU() {
	var oldestKey string
	var oldestTime int64 = 9223372036854775807 // max int64

	for key, expr := range ec.cache {
		if expr.LastUsed < oldestTime {
			oldestTime = expr.LastUsed
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(ec.cache, oldestKey)
	}
}

// Variable lookup caching methods

// GetVariableIndex retrieves cached variable index
func (vlc *VariableLookupCache) GetVariableIndex(name string) (int, bool) {
	vlc.mutex.RLock()
	defer vlc.mutex.RUnlock()

	if index, exists := vlc.cache[name]; exists {
		return index, true
	}

	return -1, false
}

// CacheVariableIndex stores variable name -> index mapping
func (vlc *VariableLookupCache) CacheVariableIndex(name string, index int) {
	vlc.mutex.Lock()
	defer vlc.mutex.Unlock()

	vlc.cache[name] = index
}

// ClearVariableCache clears the variable cache
func (vlc *VariableLookupCache) ClearVariableCache() {
	vlc.mutex.Lock()
	defer vlc.mutex.Unlock()

	vlc.cache = make(map[string]int)
}

// String pool methods for optimized string allocation

// GetOptimizedString returns a string value from the pool or creates new one
func (sp *StringPool) GetOptimizedString(value string) *types.StringValue {
	length := len(value)

	sp.mutex.RLock()

	// Check small strings cache
	if length <= sp.maxSmallSize {
		if cached, exists := sp.smallStrings[value]; exists {
			sp.mutex.RUnlock()
			return cached
		}
	} else if length <= sp.maxMediumSize {
		// Check medium strings cache
		if cached, exists := sp.mediumStrings[value]; exists {
			sp.mutex.RUnlock()
			return cached
		}
	}

	sp.mutex.RUnlock()

	// Create new string value
	stringVal := types.NewString(value)

	// Cache it if appropriate
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if length <= sp.maxSmallSize && sp.smallCount < 1000 {
		sp.smallStrings[value] = stringVal
		sp.smallCount++
	} else if length <= sp.maxMediumSize && sp.mediumCount < 500 {
		sp.mediumStrings[value] = stringVal
		sp.mediumCount++
	}

	return stringVal
}

// ClearStringPool clears the string pool
func (sp *StringPool) ClearStringPool() {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	sp.smallStrings = make(map[string]*types.StringValue)
	sp.mediumStrings = make(map[string]*types.StringValue)
	sp.smallCount = 0
	sp.mediumCount = 0
}

// Performance monitoring methods

// GetOptimizationStats returns optimization performance statistics
func (mo *MemoryOptimizer) GetOptimizationStats() OptimizationStats {
	return OptimizationStats{
		CacheHits:   mo.cacheHits,
		CacheMisses: mo.cacheMisses,
		PoolHits:    mo.poolHits,
		PoolMisses:  mo.poolMisses,
		HitRatio:    float64(mo.cacheHits) / float64(mo.cacheHits+mo.cacheMisses),
	}
}

// OptimizationStats represents optimization performance metrics
type OptimizationStats struct {
	CacheHits   int64
	CacheMisses int64
	PoolHits    int64
	PoolMisses  int64
	HitRatio    float64
}

// ResetStats resets performance counters
func (mo *MemoryOptimizer) ResetStats() {
	mo.cacheHits = 0
	mo.cacheMisses = 0
	mo.poolHits = 0
	mo.poolMisses = 0
}

// Optimized VM factory that uses memory optimization

// NewOptimizedVMWithMemoryPool creates a new VM with P1 memory optimizations
func NewOptimizedVMWithMemoryPool(bytecode *Bytecode) *VM {
	optimizer := GlobalMemoryOptimizer

	vm := &VM{
		bytecode:       bytecode,
		constants:      bytecode.Constants,
		stack:          optimizer.GetOptimizedStack(),
		sp:             0,
		globals:        optimizer.GetOptimizedGlobals(),
		pool:           NewValuePool(),
		cache:          NewInstructionCache(1000),
		customBuiltins: make(map[string]interface{}),
		safeJumpTable:  NewSafeJumpTable(),
	}

	return vm
}

// ReleaseOptimizedVM releases VM resources back to memory pools
func ReleaseOptimizedVM(vm *VM) {
	optimizer := GlobalMemoryOptimizer

	// Return stack to pool if it came from pool
	if vm.stack != nil {
		optimizer.ReleaseOptimizedStack(vm.stack)
	}

	// Return globals to pool if they came from pool
	if vm.globals != nil {
		optimizer.ReleaseOptimizedGlobals(vm.globals)
	}

	// Clear references for GC
	vm.stack = nil
	vm.globals = nil
}

// ReleaseOptimizedStack returns a stack slice to the pool
func (mo *MemoryOptimizer) ReleaseOptimizedStack(stack []types.Value) {
	mo.PutOptimizedStack(stack)
}

// ReleaseOptimizedGlobals returns a globals slice to the pool
func (mo *MemoryOptimizer) ReleaseOptimizedGlobals(globals []types.Value) {
	mo.PutOptimizedGlobals(globals)
}

// Helper function to get current timestamp
func getCurrentTimestamp() int64 {
	// Simple timestamp for LRU - could use time.Now().UnixNano() for real implementation
	// Using simple counter for now to avoid time overhead
	return globalCounter.getNext()
}

// Simple counter for timestamps
type counter struct {
	value int64
	mutex sync.Mutex
}

var globalCounter = &counter{}

func (c *counter) getNext() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
	return c.value
}

// Constants and compilation caching integration

// CompilationCache integrates with compiler for cached compilation
type CompilationCache struct {
	optimizer *MemoryOptimizer
}

// NewCompilationCache creates a new compilation cache
func NewCompilationCache() *CompilationCache {
	return &CompilationCache{
		optimizer: GlobalMemoryOptimizer,
	}
}

// GetOrCompile retrieves cached compilation or compiles new expression
func (cc *CompilationCache) GetOrCompile(expression string, compileFunc func(string) ([]byte, []types.Value, error)) ([]byte, []types.Value, error) {
	if cached, exists := cc.optimizer.exprCache.GetCachedExpression(expression); exists {
		cc.optimizer.cacheHits++
		return cached.Instructions, cached.Constants, nil
	}

	cc.optimizer.cacheMisses++

	instructions, constants, err := compileFunc(expression)
	if err != nil {
		return nil, nil, err
	}

	cc.optimizer.exprCache.CacheExpression(expression, instructions, constants)
	return instructions, constants, nil
}

// VM Integration methods

// SetMemoryOptimizer sets the memory optimizer for a VM
func (vm *VM) SetMemoryOptimizer(optimizer *MemoryOptimizer) {
	// Replace current allocations with optimized ones
	optimizer.PutOptimizedStack(vm.stack)
	optimizer.PutOptimizedGlobals(vm.globals)

	vm.stack = optimizer.GetOptimizedStack()
	vm.globals = optimizer.GetOptimizedGlobals()
}

// GetOptimizedStringValue returns an optimized string value
func (vm *VM) GetOptimizedStringValue(value string) *types.StringValue {
	return GlobalMemoryOptimizer.stringPool.GetOptimizedString(value)
}
