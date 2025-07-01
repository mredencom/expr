package checker

import (
	"testing"

	"github.com/mredencom/expr/types"
)

// TestNewScope tests basic scope creation
func TestNewScope(t *testing.T) {
	parent := NewScope(nil)
	child := NewScope(parent)

	if child.parent != parent {
		t.Error("Expected child scope to have correct parent")
	}

	if child.variables == nil {
		t.Error("Expected child scope to have initialized variables map")
	}

	if child.functions == nil {
		t.Error("Expected child scope to have initialized functions map")
	}
}

// TestNewRootScope tests root scope creation with built-in functions
func TestNewRootScope(t *testing.T) {
	rootScope := NewRootScope()

	if rootScope.parent != nil {
		t.Error("Expected root scope to have no parent")
	}

	// Test that built-in functions are registered
	expectedBuiltins := []string{
		"len", "string", "int", "float", "bool",
		"contains", "startsWith", "endsWith", "matches",
		"abs", "max", "min",
	}

	for _, builtin := range expectedBuiltins {
		funcInfo, exists := rootScope.LookupFunction(builtin)
		if !exists {
			t.Errorf("Expected builtin function '%s' to be registered", builtin)
			continue
		}
		if funcInfo == nil {
			t.Errorf("Expected non-nil function info for '%s'", builtin)
			continue
		}
		if !funcInfo.Builtin {
			t.Errorf("Expected '%s' to be marked as builtin", builtin)
		}
		if funcInfo.Name != builtin {
			t.Errorf("Expected function name '%s', got '%s'", builtin, funcInfo.Name)
		}
	}
}

// TestBuiltinFunctionSignatures tests the signatures of built-in functions
func TestBuiltinFunctionSignatures(t *testing.T) {
	rootScope := NewRootScope()

	t.Run("TypeConversionFunctions", func(t *testing.T) {
		conversionFunctions := []string{"string", "int", "float", "bool"}
		for _, funcName := range conversionFunctions {
			funcInfo, exists := rootScope.LookupFunction(funcName)
			if !exists {
				t.Errorf("Expected conversion function '%s' to exist", funcName)
				continue
			}

			// Should accept interface{} parameter
			if len(funcInfo.Params) != 1 {
				t.Errorf("Expected '%s' to have 1 parameter, got %d", funcName, len(funcInfo.Params))
			} else if funcInfo.Params[0].Kind != types.KindInterface {
				t.Errorf("Expected '%s' parameter to be interface{}, got %v", funcName, funcInfo.Params[0].Kind)
			}

			// Should return appropriate type
			if len(funcInfo.Returns) != 1 {
				t.Errorf("Expected '%s' to have 1 return value, got %d", funcName, len(funcInfo.Returns))
			} else {
				var expectedKind types.TypeKind
				switch funcName {
				case "string":
					expectedKind = types.KindString
				case "int":
					expectedKind = types.KindInt64
				case "float":
					expectedKind = types.KindFloat64
				case "bool":
					expectedKind = types.KindBool
				}
				if funcInfo.Returns[0].Kind != expectedKind {
					t.Errorf("Expected '%s' to return %v, got %v", funcName, expectedKind, funcInfo.Returns[0].Kind)
				}
			}
		}
	})

	t.Run("LenFunction", func(t *testing.T) {
		funcInfo, exists := rootScope.LookupFunction("len")
		if !exists {
			t.Fatal("Expected 'len' function to exist")
		}

		if len(funcInfo.Params) != 1 {
			t.Errorf("Expected 'len' to have 1 parameter, got %d", len(funcInfo.Params))
		}
		if len(funcInfo.Returns) != 1 {
			t.Errorf("Expected 'len' to have 1 return value, got %d", len(funcInfo.Returns))
		} else if funcInfo.Returns[0].Kind != types.KindInt64 {
			t.Errorf("Expected 'len' to return int, got %v", funcInfo.Returns[0].Kind)
		}
	})

	t.Run("StringFunctions", func(t *testing.T) {
		stringFunctions := []string{"contains", "startsWith", "endsWith", "matches"}
		for _, funcName := range stringFunctions {
			funcInfo, exists := rootScope.LookupFunction(funcName)
			if !exists {
				t.Errorf("Expected string function '%s' to exist", funcName)
				continue
			}

			// Should accept two string parameters
			if len(funcInfo.Params) != 2 {
				t.Errorf("Expected '%s' to have 2 parameters, got %d", funcName, len(funcInfo.Params))
			} else {
				for i, param := range funcInfo.Params {
					if param.Kind != types.KindString {
						t.Errorf("Expected '%s' parameter %d to be string, got %v", funcName, i, param.Kind)
					}
				}
			}

			// Should return bool
			if len(funcInfo.Returns) != 1 {
				t.Errorf("Expected '%s' to have 1 return value, got %d", funcName, len(funcInfo.Returns))
			} else if funcInfo.Returns[0].Kind != types.KindBool {
				t.Errorf("Expected '%s' to return bool, got %v", funcName, funcInfo.Returns[0].Kind)
			}
		}
	})

	t.Run("MathFunctions", func(t *testing.T) {
		// Test abs function
		absInfo, exists := rootScope.LookupFunction("abs")
		if !exists {
			t.Error("Expected 'abs' function to exist")
		} else {
			if len(absInfo.Params) != 1 {
				t.Errorf("Expected 'abs' to have 1 parameter, got %d", len(absInfo.Params))
			} else if absInfo.Params[0].Kind != types.KindFloat64 {
				t.Errorf("Expected 'abs' parameter to be float64, got %v", absInfo.Params[0].Kind)
			}
			if len(absInfo.Returns) != 1 {
				t.Errorf("Expected 'abs' to have 1 return value, got %d", len(absInfo.Returns))
			} else if absInfo.Returns[0].Kind != types.KindFloat64 {
				t.Errorf("Expected 'abs' to return float64, got %v", absInfo.Returns[0].Kind)
			}
		}

		// Test variadic functions (max, min)
		variadicFunctions := []string{"max", "min"}
		for _, funcName := range variadicFunctions {
			funcInfo, exists := rootScope.LookupFunction(funcName)
			if !exists {
				t.Errorf("Expected variadic function '%s' to exist", funcName)
				continue
			}

			if !funcInfo.Variadic {
				t.Errorf("Expected '%s' to be variadic", funcName)
			}

			if len(funcInfo.Params) != 1 {
				t.Errorf("Expected '%s' to have 1 parameter (variadic), got %d", funcName, len(funcInfo.Params))
			} else if funcInfo.Params[0].Kind != types.KindInterface {
				t.Errorf("Expected '%s' parameter to be interface{}, got %v", funcName, funcInfo.Params[0].Kind)
			}

			if len(funcInfo.Returns) != 1 {
				t.Errorf("Expected '%s' to have 1 return value, got %d", funcName, len(funcInfo.Returns))
			} else if funcInfo.Returns[0].Kind != types.KindInterface {
				t.Errorf("Expected '%s' to return interface{}, got %v", funcName, funcInfo.Returns[0].Kind)
			}
		}
	})
}

// TestScopeVariableOperations tests variable management in scopes
func TestScopeVariableOperations(t *testing.T) {
	rootScope := NewScope(nil)
	childScope := NewScope(rootScope)

	// Test defining variables
	intType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
	stringType := types.TypeInfo{Kind: types.KindString, Name: "string"}

	rootScope.DefineVariable("rootVar", intType)
	childScope.DefineVariable("childVar", stringType)

	// Test looking up variables in current scope
	if !childScope.HasVariable("childVar") {
		t.Error("Expected child scope to have childVar")
	}

	if rootScope.HasVariable("childVar") {
		t.Error("Expected root scope to not have childVar")
	}

	// Test looking up variables in parent scope
	foundType, exists := childScope.LookupVariable("rootVar")
	if !exists {
		t.Error("Expected to find rootVar in parent scope")
	} else if foundType.Kind != intType.Kind {
		t.Errorf("Expected rootVar to be int type, got %v", foundType.Kind)
	}

	// Test variable shadowing
	childScope.DefineVariable("rootVar", stringType) // Shadow parent variable
	foundType, exists = childScope.LookupVariable("rootVar")
	if !exists {
		t.Error("Expected to find shadowed rootVar")
	} else if foundType.Kind != stringType.Kind {
		t.Errorf("Expected shadowed rootVar to be string type, got %v", foundType.Kind)
	}

	// Test that parent still has original variable
	foundType, exists = rootScope.LookupVariable("rootVar")
	if !exists {
		t.Error("Expected to find original rootVar in parent")
	} else if foundType.Kind != intType.Kind {
		t.Errorf("Expected original rootVar to be int type, got %v", foundType.Kind)
	}
}

// TestScopeFunctionOperations tests function management in scopes
func TestScopeFunctionOperations(t *testing.T) {
	rootScope := NewScope(nil)
	childScope := NewScope(rootScope)

	// Test defining functions
	customFunc := &FunctionInfo{
		Name:    "customFunc",
		Params:  []types.TypeInfo{{Kind: types.KindInt64, Name: "int"}},
		Returns: []types.TypeInfo{{Kind: types.KindString, Name: "string"}},
		Builtin: false,
	}

	rootScope.DefineFunction("customFunc", customFunc)

	// Test looking up functions in current scope
	if !rootScope.HasFunction("customFunc") {
		t.Error("Expected root scope to have customFunc")
	}

	if childScope.HasFunction("customFunc") {
		t.Error("Expected child scope to not have customFunc in local scope")
	}

	// Test looking up functions in parent scope
	foundFunc, exists := childScope.LookupFunction("customFunc")
	if !exists {
		t.Error("Expected to find customFunc in parent scope")
	} else {
		if foundFunc.Name != "customFunc" {
			t.Errorf("Expected function name 'customFunc', got '%s'", foundFunc.Name)
		}
		if foundFunc.Builtin {
			t.Error("Expected customFunc to not be builtin")
		}
	}

	// Test function shadowing
	shadowFunc := &FunctionInfo{
		Name:    "customFunc",
		Params:  []types.TypeInfo{{Kind: types.KindString, Name: "string"}},
		Returns: []types.TypeInfo{{Kind: types.KindBool, Name: "bool"}},
		Builtin: false,
	}

	childScope.DefineFunction("customFunc", shadowFunc)
	foundFunc, exists = childScope.LookupFunction("customFunc")
	if !exists {
		t.Error("Expected to find shadowed customFunc")
	} else {
		if len(foundFunc.Params) != 1 || foundFunc.Params[0].Kind != types.KindString {
			t.Error("Expected shadowed function to have string parameter")
		}
		if len(foundFunc.Returns) != 1 || foundFunc.Returns[0].Kind != types.KindBool {
			t.Error("Expected shadowed function to return bool")
		}
	}
}

// TestScopeParentAccess tests parent scope access
func TestScopeParentAccess(t *testing.T) {
	rootScope := NewScope(nil)
	childScope := NewScope(rootScope)
	grandChildScope := NewScope(childScope)

	if rootScope.Parent() != nil {
		t.Error("Expected root scope to have no parent")
	}

	if childScope.Parent() != rootScope {
		t.Error("Expected child scope parent to be root scope")
	}

	if grandChildScope.Parent() != childScope {
		t.Error("Expected grandchild scope parent to be child scope")
	}
}

// TestScopeVariablesAndFunctionsListing tests Variables() and Functions() methods
func TestScopeVariablesAndFunctionsListing(t *testing.T) {
	scope := NewScope(nil)

	// Initially empty (except for root scope which has builtins)
	variables := scope.Variables()
	functions := scope.Functions()

	if len(variables) != 0 {
		t.Errorf("Expected empty variables map, got %d items", len(variables))
	}

	if len(functions) != 0 {
		t.Errorf("Expected empty functions map, got %d items", len(functions))
	}

	// Add some variables and functions
	intType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
	stringType := types.TypeInfo{Kind: types.KindString, Name: "string"}

	scope.DefineVariable("var1", intType)
	scope.DefineVariable("var2", stringType)

	testFunc := &FunctionInfo{
		Name:    "testFunc",
		Params:  []types.TypeInfo{intType},
		Returns: []types.TypeInfo{stringType},
		Builtin: false,
	}

	scope.DefineFunction("testFunc", testFunc)

	// Check that Variables() returns correct data
	variables = scope.Variables()
	if len(variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(variables))
	}

	if varType, exists := variables["var1"]; !exists {
		t.Error("Expected var1 to exist in variables map")
	} else if varType.Kind != intType.Kind {
		t.Errorf("Expected var1 to be int type, got %v", varType.Kind)
	}

	if varType, exists := variables["var2"]; !exists {
		t.Error("Expected var2 to exist in variables map")
	} else if varType.Kind != stringType.Kind {
		t.Errorf("Expected var2 to be string type, got %v", varType.Kind)
	}

	// Check that Functions() returns correct data
	functions = scope.Functions()
	if len(functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(functions))
	}

	if funcInfo, exists := functions["testFunc"]; !exists {
		t.Error("Expected testFunc to exist in functions map")
	} else if funcInfo.Name != "testFunc" {
		t.Errorf("Expected function name 'testFunc', got '%s'", funcInfo.Name)
	}
}

// TestScopeChainLookup tests variable and function lookup through scope chain
func TestScopeChainLookup(t *testing.T) {
	// Create a scope chain: root -> child -> grandchild
	rootScope := NewScope(nil)
	childScope := NewScope(rootScope)
	grandChildScope := NewScope(childScope)

	// Define variables at different levels
	intType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
	stringType := types.TypeInfo{Kind: types.KindString, Name: "string"}
	boolType := types.TypeInfo{Kind: types.KindBool, Name: "bool"}

	rootScope.DefineVariable("rootVar", intType)
	childScope.DefineVariable("childVar", stringType)
	grandChildScope.DefineVariable("grandChildVar", boolType)

	// Define functions at different levels
	rootFunc := &FunctionInfo{Name: "rootFunc", Params: []types.TypeInfo{}, Returns: []types.TypeInfo{intType}, Builtin: false}
	childFunc := &FunctionInfo{Name: "childFunc", Params: []types.TypeInfo{}, Returns: []types.TypeInfo{stringType}, Builtin: false}
	grandChildFunc := &FunctionInfo{Name: "grandChildFunc", Params: []types.TypeInfo{}, Returns: []types.TypeInfo{boolType}, Builtin: false}

	rootScope.DefineFunction("rootFunc", rootFunc)
	childScope.DefineFunction("childFunc", childFunc)
	grandChildScope.DefineFunction("grandChildFunc", grandChildFunc)

	// Test variable lookup from grandchild scope
	tests := []struct {
		varName      string
		shouldExist  bool
		expectedType types.TypeKind
	}{
		{"grandChildVar", true, types.KindBool},
		{"childVar", true, types.KindString},
		{"rootVar", true, types.KindInt64},
		{"nonexistentVar", false, types.KindBool}, // Kind doesn't matter for non-existent
	}

	for _, test := range tests {
		foundType, exists := grandChildScope.LookupVariable(test.varName)
		if exists != test.shouldExist {
			t.Errorf("Variable '%s': expected exists=%v, got %v", test.varName, test.shouldExist, exists)
		}
		if exists && foundType.Kind != test.expectedType {
			t.Errorf("Variable '%s': expected type %v, got %v", test.varName, test.expectedType, foundType.Kind)
		}
	}

	// Test function lookup from grandchild scope
	funcTests := []struct {
		funcName       string
		shouldExist    bool
		expectedReturn types.TypeKind
	}{
		{"grandChildFunc", true, types.KindBool},
		{"childFunc", true, types.KindString},
		{"rootFunc", true, types.KindInt64},
		{"nonexistentFunc", false, types.KindBool}, // Kind doesn't matter for non-existent
	}

	for _, test := range funcTests {
		foundFunc, exists := grandChildScope.LookupFunction(test.funcName)
		if exists != test.shouldExist {
			t.Errorf("Function '%s': expected exists=%v, got %v", test.funcName, test.shouldExist, exists)
		}
		if exists && len(foundFunc.Returns) > 0 && foundFunc.Returns[0].Kind != test.expectedReturn {
			t.Errorf("Function '%s': expected return type %v, got %v", test.funcName, test.expectedReturn, foundFunc.Returns[0].Kind)
		}
	}
}

// TestFunctionInfo tests FunctionInfo structure
func TestFunctionInfo(t *testing.T) {
	funcInfo := &FunctionInfo{
		Name: "testFunction",
		Params: []types.TypeInfo{
			{Kind: types.KindInt64, Name: "int"},
			{Kind: types.KindString, Name: "string"},
		},
		Returns: []types.TypeInfo{
			{Kind: types.KindBool, Name: "bool"},
		},
		Variadic: true,
		Builtin:  false,
	}

	if funcInfo.Name != "testFunction" {
		t.Errorf("Expected function name 'testFunction', got '%s'", funcInfo.Name)
	}

	if len(funcInfo.Params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(funcInfo.Params))
	}

	if len(funcInfo.Returns) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(funcInfo.Returns))
	}

	if !funcInfo.Variadic {
		t.Error("Expected function to be variadic")
	}

	if funcInfo.Builtin {
		t.Error("Expected function to not be builtin")
	}
}

// TestRootScopeBuiltinCompletion tests that all expected builtins are present
func TestRootScopeBuiltinCompletion(t *testing.T) {
	rootScope := NewRootScope()

	// Verify specific builtin function properties
	testCases := []struct {
		name        string
		paramCount  int
		returnCount int
		variadic    bool
	}{
		{"len", 1, 1, false},
		{"string", 1, 1, false},
		{"int", 1, 1, false},
		{"float", 1, 1, false},
		{"bool", 1, 1, false},
		{"contains", 2, 1, false},
		{"startsWith", 2, 1, false},
		{"endsWith", 2, 1, false},
		{"matches", 2, 1, false},
		{"abs", 1, 1, false},
		{"max", 1, 1, true},
		{"min", 1, 1, true},
	}

	for _, test := range testCases {
		funcInfo, exists := rootScope.LookupFunction(test.name)
		if !exists {
			t.Errorf("Expected builtin function '%s' to exist", test.name)
			continue
		}

		if len(funcInfo.Params) != test.paramCount {
			t.Errorf("Function '%s': expected %d parameters, got %d", test.name, test.paramCount, len(funcInfo.Params))
		}

		if len(funcInfo.Returns) != test.returnCount {
			t.Errorf("Function '%s': expected %d return values, got %d", test.name, test.returnCount, len(funcInfo.Returns))
		}

		if funcInfo.Variadic != test.variadic {
			t.Errorf("Function '%s': expected variadic=%v, got %v", test.name, test.variadic, funcInfo.Variadic)
		}

		if !funcInfo.Builtin {
			t.Errorf("Function '%s': expected to be marked as builtin", test.name)
		}
	}
}
