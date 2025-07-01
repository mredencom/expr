package checker

import (
	"github.com/mredencom/expr/types"
)

// Scope represents a lexical scope for variables and functions
type Scope struct {
	parent    *Scope
	variables map[string]types.TypeInfo
	functions map[string]*FunctionInfo
}

// FunctionInfo contains information about a function
type FunctionInfo struct {
	Name     string
	Params   []types.TypeInfo
	Returns  []types.TypeInfo
	Variadic bool
	Builtin  bool
}

// NewScope creates a new scope
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]types.TypeInfo),
		functions: make(map[string]*FunctionInfo),
	}
}

// NewRootScope creates a new root scope with built-in functions
func NewRootScope() *Scope {
	scope := NewScope(nil)

	// Add built-in functions
	scope.functions["len"] = &FunctionInfo{
		Name:    "len",
		Params:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns: []types.TypeInfo{types.IntType},
		Builtin: true,
	}

	scope.functions["string"] = &FunctionInfo{
		Name:    "string",
		Params:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns: []types.TypeInfo{types.StringType},
		Builtin: true,
	}

	scope.functions["int"] = &FunctionInfo{
		Name:    "int",
		Params:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns: []types.TypeInfo{types.IntType},
		Builtin: true,
	}

	scope.functions["float"] = &FunctionInfo{
		Name:    "float",
		Params:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns: []types.TypeInfo{types.FloatType},
		Builtin: true,
	}

	scope.functions["bool"] = &FunctionInfo{
		Name:    "bool",
		Params:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns: []types.TypeInfo{types.BoolType},
		Builtin: true,
	}

	// String functions
	scope.functions["contains"] = &FunctionInfo{
		Name:    "contains",
		Params:  []types.TypeInfo{types.StringType, types.StringType},
		Returns: []types.TypeInfo{types.BoolType},
		Builtin: true,
	}

	scope.functions["startsWith"] = &FunctionInfo{
		Name:    "startsWith",
		Params:  []types.TypeInfo{types.StringType, types.StringType},
		Returns: []types.TypeInfo{types.BoolType},
		Builtin: true,
	}

	scope.functions["endsWith"] = &FunctionInfo{
		Name:    "endsWith",
		Params:  []types.TypeInfo{types.StringType, types.StringType},
		Returns: []types.TypeInfo{types.BoolType},
		Builtin: true,
	}

	scope.functions["matches"] = &FunctionInfo{
		Name:    "matches",
		Params:  []types.TypeInfo{types.StringType, types.StringType},
		Returns: []types.TypeInfo{types.BoolType},
		Builtin: true,
	}

	// Math functions
	scope.functions["abs"] = &FunctionInfo{
		Name:    "abs",
		Params:  []types.TypeInfo{{Kind: types.KindFloat64, Name: "float"}},
		Returns: []types.TypeInfo{types.FloatType},
		Builtin: true,
	}

	scope.functions["max"] = &FunctionInfo{
		Name:     "max",
		Params:   []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Variadic: true,
		Builtin:  true,
	}

	scope.functions["min"] = &FunctionInfo{
		Name:     "min",
		Params:   []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Returns:  []types.TypeInfo{{Kind: types.KindInterface, Name: "interface{}"}},
		Variadic: true,
		Builtin:  true,
	}

	return scope
}

// DefineVariable defines a variable in the current scope
func (s *Scope) DefineVariable(name string, typeInfo types.TypeInfo) {
	s.variables[name] = typeInfo
}

// DefineFunction defines a function in the current scope
func (s *Scope) DefineFunction(name string, funcInfo *FunctionInfo) {
	s.functions[name] = funcInfo
}

// LookupVariable looks up a variable in the scope chain
func (s *Scope) LookupVariable(name string) (types.TypeInfo, bool) {
	if typeInfo, ok := s.variables[name]; ok {
		return typeInfo, true
	}

	if s.parent != nil {
		return s.parent.LookupVariable(name)
	}

	return types.TypeInfo{}, false
}

// LookupFunction looks up a function in the scope chain
func (s *Scope) LookupFunction(name string) (*FunctionInfo, bool) {
	if funcInfo, ok := s.functions[name]; ok {
		return funcInfo, true
	}

	if s.parent != nil {
		return s.parent.LookupFunction(name)
	}

	return nil, false
}

// HasVariable checks if a variable exists in the current scope (not parent scopes)
func (s *Scope) HasVariable(name string) bool {
	_, ok := s.variables[name]
	return ok
}

// HasFunction checks if a function exists in the current scope (not parent scopes)
func (s *Scope) HasFunction(name string) bool {
	_, ok := s.functions[name]
	return ok
}

// Parent returns the parent scope
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Variables returns all variables in the current scope
func (s *Scope) Variables() map[string]types.TypeInfo {
	result := make(map[string]types.TypeInfo)
	for name, typeInfo := range s.variables {
		result[name] = typeInfo
	}
	return result
}

// Functions returns all functions in the current scope
func (s *Scope) Functions() map[string]*FunctionInfo {
	result := make(map[string]*FunctionInfo)
	for name, funcInfo := range s.functions {
		result[name] = funcInfo
	}
	return result
}
