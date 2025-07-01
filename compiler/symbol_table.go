package compiler

// SymbolScope represents the scope of a symbol
type SymbolScope string

const (
	LocalScope    SymbolScope = "LOCAL"
	GlobalScope   SymbolScope = "GLOBAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

// Symbol represents a symbol in the symbol table
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable manages symbols and their scopes
type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free}
}

// NewEnclosedSymbolTable creates a new symbol table with an outer scope
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define defines a new symbol
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// DefineBuiltin defines a builtin symbol
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltinScope, Index: index}
	s.store[name] = symbol
	return symbol
}

// DefineFunctionName defines a function name symbol
func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Scope: FunctionScope, Index: 0}
	s.store[name] = symbol
	return symbol
}

// Resolve resolves a symbol by name
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}

	return obj, ok
}

// defineFree defines a free symbol
func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	s.store[original.Name] = symbol
	return symbol
}

// GetStore returns the symbol store for debugging
func (s *SymbolTable) GetStore() map[string]Symbol {
	return s.store
}
