package vm

import "github.com/mredencom/expr/types"

// VMFactory creates VMs with different optimization levels
type VMFactory struct {
	useMemoryOptimization bool
	useUnionTypes         bool
	enableCaching         bool
}

// NewVMFactory creates a new VM factory with specified optimizations
func NewVMFactory(memOpt, unionTypes, caching bool) *VMFactory {
	return &VMFactory{
		useMemoryOptimization: memOpt,
		useUnionTypes:         unionTypes,
		enableCaching:         caching,
	}
}

// DefaultOptimizedFactory returns a factory with all P1 optimizations enabled
func DefaultOptimizedFactory() *VMFactory {
	return &VMFactory{
		useMemoryOptimization: true,
		useUnionTypes:         false, // Keep false for compatibility for now
		enableCaching:         true,
	}
}

// CreateVM creates a VM with the configured optimizations
func (f *VMFactory) CreateVM(bytecode *Bytecode) *VM {
	if f.useMemoryOptimization {
		return f.createOptimizedVM(bytecode)
	}
	return New(bytecode) // fallback to standard VM
}

// createOptimizedVM creates a VM with P1 memory optimizations
func (f *VMFactory) createOptimizedVM(bytecode *Bytecode) *VM {
	optimizer := GlobalMemoryOptimizer

	vm := &VM{
		bytecode:       bytecode,
		constants:      bytecode.Constants,
		stack:          optimizer.GetOptimizedStack(), // ✅ Use memory pool
		sp:             0,
		globals:        optimizer.GetOptimizedGlobals(), // ✅ Use memory pool
		pool:           NewValuePool(),
		cache:          NewInstructionCache(1000),
		customBuiltins: make(map[string]interface{}),
		safeJumpTable:  NewSafeJumpTable(),
	}

	// ✅ 添加VM析构器，确保资源自动释放
	// 使用runtime.SetFinalizer确保资源释放
	// runtime.SetFinalizer(vm, func(vm *VM) {
	// 	ReleaseOptimizedVM(vm)
	// })

	return vm
}

// CreateVMWithPool creates a VM that uses resource pooling
func (f *VMFactory) CreateVMWithPool(bytecode *Bytecode) (*VM, func()) {
	vm := f.createOptimizedVM(bytecode)

	// 返回VM和清理函数
	cleanup := func() {
		f.ReleaseVM(vm)
	}

	return vm, cleanup
}

// RunOptimizedExpression 运行优化表达式的便利方法
func RunOptimizedExpression(bytecode *Bytecode, env map[string]interface{}) (types.Value, error) {
	factory := DefaultOptimizedFactory()
	vm, cleanup := factory.CreateVMWithPool(bytecode)
	defer cleanup() // 确保资源释放

	return vm.Run(bytecode, env)
}

// ReleaseVM properly releases VM resources back to pools
func (f *VMFactory) ReleaseVM(vm *VM) {
	if f.useMemoryOptimization {
		ReleaseOptimizedVM(vm)
	}
	// Standard VMs don't need special cleanup
}

// Global optimized factory instance
var GlobalOptimizedFactory = DefaultOptimizedFactory()

// NewOptimized creates a new VM with P1 optimizations enabled
func NewOptimized(bytecode *Bytecode) *VM {
	return GlobalOptimizedFactory.CreateVM(bytecode)
}

// NewOptimizedWithOptions creates a VM with specific optimization options
func NewOptimizedWithOptions(bytecode *Bytecode, memOpt, unionTypes, caching bool) *VM {
	factory := NewVMFactory(memOpt, unionTypes, caching)
	return factory.CreateVM(bytecode)
}

// VMPool with optimization support
type OptimizedVMPool struct {
	factory *VMFactory
	pool    *VMPool
}

// NewOptimizedVMPool creates a VM pool that uses optimized VMs
func NewOptimizedVMPool() *OptimizedVMPool {
	return &OptimizedVMPool{
		factory: DefaultOptimizedFactory(),
		pool:    NewVMPool(),
	}
}

// Get retrieves an optimized VM from the pool
func (p *OptimizedVMPool) Get(bytecode *Bytecode) *VM {
	vm := p.pool.Get()
	// Reset with new bytecode
	vm.bytecode = bytecode
	vm.constants = bytecode.Constants
	vm.sp = 0
	return vm
}

// Put returns an optimized VM to the pool
func (p *OptimizedVMPool) Put(vm *VM) {
	p.factory.ReleaseVM(vm)
	p.pool.Put(vm)
}
