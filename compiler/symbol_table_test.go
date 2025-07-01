package compiler

import "testing"

func TestNewSymbolTable(t *testing.T) {
	st := NewSymbolTable()
	if st == nil {
		t.Fatal("Expected non-nil symbol table")
	}
	if st.store == nil {
		t.Fatal("Expected non-nil store")
	}
	if st.numDefinitions != 0 {
		t.Errorf("Expected 0 definitions, got %d", st.numDefinitions)
	}
}

func TestNewEnclosedSymbolTable(t *testing.T) {
	outer := NewSymbolTable()
	inner := NewEnclosedSymbolTable(outer)

	if inner == nil {
		t.Fatal("Expected non-nil symbol table")
	}
	if inner.Outer != outer {
		t.Error("Expected inner table to reference outer table")
	}
}

func TestDefine(t *testing.T) {
	st := NewSymbolTable()

	symbol := st.Define("x")
	if symbol.Name != "x" {
		t.Errorf("Expected name 'x', got %s", symbol.Name)
	}
	if symbol.Index != 0 {
		t.Errorf("Expected index 0, got %d", symbol.Index)
	}
	if symbol.Scope != GlobalScope {
		t.Errorf("Expected GlobalScope, got %s", symbol.Scope)
	}

	// Test that the symbol is stored
	resolved, ok := st.Resolve("x")
	if !ok {
		t.Fatal("Expected to resolve symbol 'x'")
	}
	if resolved != symbol {
		t.Error("Expected resolved symbol to match defined symbol")
	}
}

func TestDefineLocal(t *testing.T) {
	global := NewSymbolTable()
	local := NewEnclosedSymbolTable(global)

	symbol := local.Define("x")
	if symbol.Scope != LocalScope {
		t.Errorf("Expected LocalScope, got %s", symbol.Scope)
	}
}

func TestDefineBuiltin(t *testing.T) {
	st := NewSymbolTable()

	symbol := st.DefineBuiltin(0, "len")
	if symbol.Name != "len" {
		t.Errorf("Expected name 'len', got %s", symbol.Name)
	}
	if symbol.Index != 0 {
		t.Errorf("Expected index 0, got %d", symbol.Index)
	}
	if symbol.Scope != BuiltinScope {
		t.Errorf("Expected BuiltinScope, got %s", symbol.Scope)
	}

	// Test that the builtin is stored
	resolved, ok := st.Resolve("len")
	if !ok {
		t.Fatal("Expected to resolve builtin 'len'")
	}
	if resolved != symbol {
		t.Error("Expected resolved builtin to match defined builtin")
	}
}

func TestDefineFunctionName(t *testing.T) {
	st := NewSymbolTable()

	symbol := st.DefineFunctionName("myFunc")
	if symbol.Name != "myFunc" {
		t.Errorf("Expected name 'myFunc', got %s", symbol.Name)
	}
	if symbol.Index != 0 {
		t.Errorf("Expected index 0, got %d", symbol.Index)
	}
	if symbol.Scope != FunctionScope {
		t.Errorf("Expected FunctionScope, got %s", symbol.Scope)
	}
}

func TestResolveGlobal(t *testing.T) {
	st := NewSymbolTable()
	st.Define("x")
	st.Define("y")

	// Test resolving existing symbols
	symbol, ok := st.Resolve("x")
	if !ok {
		t.Fatal("Expected to resolve 'x'")
	}
	if symbol.Name != "x" || symbol.Index != 0 || symbol.Scope != GlobalScope {
		t.Errorf("Unexpected symbol: %+v", symbol)
	}

	symbol, ok = st.Resolve("y")
	if !ok {
		t.Fatal("Expected to resolve 'y'")
	}
	if symbol.Name != "y" || symbol.Index != 1 || symbol.Scope != GlobalScope {
		t.Errorf("Unexpected symbol: %+v", symbol)
	}

	// Test resolving non-existent symbol
	_, ok = st.Resolve("nonexistent")
	if ok {
		t.Error("Expected not to resolve non-existent symbol")
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	// Test resolving local symbols
	symbol, ok := local.Resolve("c")
	if !ok {
		t.Fatal("Expected to resolve local 'c'")
	}
	if symbol.Scope != LocalScope {
		t.Errorf("Expected LocalScope, got %s", symbol.Scope)
	}

	// Test resolving global symbols from local scope
	symbol, ok = local.Resolve("a")
	if !ok {
		t.Fatal("Expected to resolve global 'a' from local scope")
	}
	if symbol.Scope != GlobalScope {
		t.Errorf("Expected GlobalScope, got %s", symbol.Scope)
	}
}

func TestResolveBuiltin(t *testing.T) {
	global := NewSymbolTable()
	global.DefineBuiltin(0, "len")
	global.DefineBuiltin(1, "string")

	local := NewEnclosedSymbolTable(global)

	// Test resolving builtins from local scope
	symbol, ok := local.Resolve("len")
	if !ok {
		t.Fatal("Expected to resolve builtin 'len' from local scope")
	}
	if symbol.Scope != BuiltinScope {
		t.Errorf("Expected BuiltinScope, got %s", symbol.Scope)
	}

	symbol, ok = local.Resolve("string")
	if !ok {
		t.Fatal("Expected to resolve builtin 'string' from local scope")
	}
	if symbol.Scope != BuiltinScope {
		t.Errorf("Expected BuiltinScope, got %s", symbol.Scope)
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	// Test resolving from second local scope
	// Should create free variables for first local scope variables
	symbol, ok := secondLocal.Resolve("c")
	if !ok {
		t.Fatal("Expected to resolve 'c' as free variable")
	}
	if symbol.Scope != FreeScope {
		t.Errorf("Expected FreeScope, got %s", symbol.Scope)
	}

	// Test that free symbols are tracked
	if len(secondLocal.FreeSymbols) != 1 {
		t.Errorf("Expected 1 free symbol, got %d", len(secondLocal.FreeSymbols))
	}

	if secondLocal.FreeSymbols[0].Name != "c" {
		t.Errorf("Expected free symbol 'c', got %s", secondLocal.FreeSymbols[0].Name)
	}
}

func TestDefineFree(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("a")

	secondLocal := NewEnclosedSymbolTable(firstLocal)

	// Manually call defineFree
	originalSymbol := Symbol{Name: "a", Scope: LocalScope, Index: 0}
	freeSymbol := secondLocal.defineFree(originalSymbol)

	if freeSymbol.Name != "a" {
		t.Errorf("Expected name 'a', got %s", freeSymbol.Name)
	}
	if freeSymbol.Scope != FreeScope {
		t.Errorf("Expected FreeScope, got %s", freeSymbol.Scope)
	}
	if freeSymbol.Index != 0 {
		t.Errorf("Expected index 0, got %d", freeSymbol.Index)
	}

	// Check that it's added to FreeSymbols
	if len(secondLocal.FreeSymbols) != 1 {
		t.Errorf("Expected 1 free symbol, got %d", len(secondLocal.FreeSymbols))
	}
}

func TestSymbolScopes(t *testing.T) {
	// Test all scope types
	scopes := []SymbolScope{
		LocalScope,
		GlobalScope,
		BuiltinScope,
		FreeScope,
		FunctionScope,
	}

	expectedStrings := []string{
		"LOCAL",
		"GLOBAL",
		"BUILTIN",
		"FREE",
		"FUNCTION",
	}

	for i, scope := range scopes {
		if string(scope) != expectedStrings[i] {
			t.Errorf("Expected scope %s, got %s", expectedStrings[i], string(scope))
		}
	}
}

func TestMultipleDefinitions(t *testing.T) {
	st := NewSymbolTable()

	vars := []string{"x", "y", "z"}
	for i, name := range vars {
		symbol := st.Define(name)
		if symbol.Index != i {
			t.Errorf("Expected index %d for %s, got %d", i, name, symbol.Index)
		}
	}

	if st.numDefinitions != 3 {
		t.Errorf("Expected 3 definitions, got %d", st.numDefinitions)
	}
}
