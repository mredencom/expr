package modules

import (
	"fmt"
	"sync"

	"github.com/mredencom/expr/types"
)

// ModuleFunction represents a function in a module
type ModuleFunction struct {
	Name        string
	Description string
	Handler     func(args ...interface{}) (interface{}, error)
	ParamTypes  []types.TypeInfo
	ReturnType  types.TypeInfo
	Variadic    bool
}

// Module represents a module with its functions
type Module struct {
	Name        string
	Description string
	Functions   map[string]*ModuleFunction
}

// Registry manages all registered modules
type Registry struct {
	modules map[string]*Module
	mu      sync.RWMutex
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	registry := &Registry{
		modules: make(map[string]*Module),
	}

	// Register built-in modules
	registry.registerBuiltinModules()

	return registry
}

// RegisterModule registers a new module
func (r *Registry) RegisterModule(name string, description string, functions map[string]*ModuleFunction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module '%s' already registered", name)
	}

	module := &Module{
		Name:        name,
		Description: description,
		Functions:   functions,
	}

	r.modules[name] = module
	return nil
}

// GetModule retrieves a module by name
func (r *Registry) GetModule(name string) (*Module, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	module, exists := r.modules[name]
	if !exists {
		return nil, fmt.Errorf("module '%s' not found", name)
	}

	return module, nil
}

// HasModule checks if a module exists
func (r *Registry) HasModule(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.modules[name]
	return exists
}

// GetFunction retrieves a function from a module
func (r *Registry) GetFunction(moduleName, functionName string) (*ModuleFunction, error) {
	module, err := r.GetModule(moduleName)
	if err != nil {
		return nil, err
	}

	function, exists := module.Functions[functionName]
	if !exists {
		return nil, fmt.Errorf("function '%s' not found in module '%s'", functionName, moduleName)
	}

	return function, nil
}

// CallFunction calls a function from a module
func (r *Registry) CallFunction(moduleName, functionName string, args ...interface{}) (interface{}, error) {
	function, err := r.GetFunction(moduleName, functionName)
	if err != nil {
		return nil, err
	}

	return function.Handler(args...)
}

// ListModules returns all registered module names
func (r *Registry) ListModules() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}

	return names
}

// GetModuleInfo returns information about a module
func (r *Registry) GetModuleInfo(name string) (*Module, error) {
	return r.GetModule(name)
}

// registerBuiltinModules registers all built-in modules
func (r *Registry) registerBuiltinModules() {
	// Register math module
	r.registerMathModule()

	// Register strings module
	r.registerStringsModule()
}

// Global registry instance
var DefaultRegistry = NewRegistry()
