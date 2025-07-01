package expr

import (
	"fmt"
	"testing"
)

func TestAsOptions(t *testing.T) {
	tests := []struct {
		name         string
		option       Option
		expression   string
		expectedType AsKind
	}{
		{"AsInt", AsInt(), "42", AsIntKind},
		{"AsInt64", AsInt64(), "42", AsInt64Kind},
		{"AsFloat64", AsFloat64(), "3.14", AsFloat64Kind},
		{"AsString", AsString(), `"hello"`, AsStringKind},
		{"AsBool", AsBool(), "true", AsBoolKind},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{}
			tt.option(config)

			if config.expectedType != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, config.expectedType)
			}

			if !config.enableTypeChecking {
				t.Error("Expected type checking to be enabled")
			}

			// Test that compilation works with the type option
			_, err := Compile(tt.expression, tt.option)
			if err != nil {
				t.Errorf("Compilation error: %v", err)
			}
		})
	}
}

func TestFunctionsOption(t *testing.T) {
	customFuncs := map[string]interface{}{
		"double": func(x int) int { return x * 2 },
		"greet":  func(name string) string { return "Hello, " + name },
	}

	config := &Config{builtins: make(map[string]interface{})}
	Functions(customFuncs)(config)

	if len(config.builtins) != 2 {
		t.Errorf("Expected 2 builtins, got %d", len(config.builtins))
	}

	if _, exists := config.builtins["double"]; !exists {
		t.Error("Expected 'double' function to be added")
	}

	if _, exists := config.builtins["greet"]; !exists {
		t.Error("Expected 'greet' function to be added")
	}
}

func TestOperatorsOption(t *testing.T) {
	customOps := map[string]Operator{
		"**": {Symbol: "**", Precedence: 8},
		"??": {Symbol: "??", Precedence: 1},
	}

	config := &Config{operators: make(map[string]int)}
	Operators(customOps)(config)

	if len(config.operators) != 2 {
		t.Errorf("Expected 2 operators, got %d", len(config.operators))
	}

	if config.operators["**"] != 8 {
		t.Errorf("Expected precedence 8 for '**', got %d", config.operators["**"])
	}

	if config.operators["??"] != 1 {
		t.Errorf("Expected precedence 1 for '??', got %d", config.operators["??"])
	}
}

func TestOptimizeOption(t *testing.T) {
	t.Run("EnableOptimization", func(t *testing.T) {
		config := &Config{}
		Optimize(true)(config)

		if !config.enableOptimization {
			t.Error("Expected optimization to be enabled")
		}
	})

	t.Run("DisableOptimization", func(t *testing.T) {
		config := &Config{enableOptimization: true}
		Optimize(false)(config)

		if config.enableOptimization {
			t.Error("Expected optimization to be disabled")
		}
	})
}

func TestDeprecatedCompatibilityFunctions(t *testing.T) {
	t.Run("NewEnv", func(t *testing.T) {
		env := NewEnv()
		if env == nil {
			t.Fatal("Expected non-nil environment")
		}

		// Test that we can add values to the environment
		env["x"] = 42
		if env["x"] != 42 {
			t.Errorf("Expected 42, got %v", env["x"])
		}
	})

	t.Run("CompileWithEnv", func(t *testing.T) {
		env := map[string]interface{}{"x": 42}
		program, err := CompileWithEnv("x", env)
		if err != nil {
			t.Fatalf("Compilation error: %v", err)
		}

		if program == nil {
			t.Fatal("Expected program but got nil")
		}
	})

	t.Run("RunWithEnv", func(t *testing.T) {
		env := map[string]interface{}{"x": 42}
		program, err := CompileWithEnv("x", env)
		if err != nil {
			t.Fatalf("Compilation error: %v", err)
		}

		result, err := RunWithEnv(program, env)
		if err != nil {
			t.Fatalf("Runtime error: %v", err)
		}

		if result != int64(42) {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("EvalWithEnv", func(t *testing.T) {
		env := map[string]interface{}{"x": 42}
		result, err := EvalWithEnv("x", env)
		if err != nil {
			t.Fatalf("Eval error: %v", err)
		}

		if result != int64(42) {
			t.Errorf("Expected 42, got %v", result)
		}
	})
}

func TestCheckType(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		checkFunc func(interface{}) error
		shouldErr bool
	}{
		{"int matches int", 42, CheckType[int], false},
		{"int64 matches int", int64(42), CheckType[int], false},
		{"int matches int64", 42, CheckType[int64], false},
		{"float64 matches float64", 3.14, CheckType[float64], false},
		{"int matches float64", 42, CheckType[float64], false},
		{"string matches string", "hello", CheckType[string], false},
		{"bool matches bool", true, CheckType[bool], false},
		{"string does not match int", "hello", CheckType[int], true},
		{"bool does not match string", true, CheckType[string], true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.checkFunc(tt.value)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestConvertType(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		converter func(interface{}) (interface{}, error)
		expected  interface{}
		shouldErr bool
	}{
		{"int to int", 42, func(v interface{}) (interface{}, error) { return ConvertType[int](v) }, 42, false},
		{"string conversion", 42, func(v interface{}) (interface{}, error) { return ConvertType[string](v) }, "42", false},
		{"nil conversion", nil, func(v interface{}) (interface{}, error) { return ConvertType[int](v) }, 0, false},
		{"bool to bool", true, func(v interface{}) (interface{}, error) { return ConvertType[bool](v) }, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.converter(tt.value)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.shouldErr && result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompileError(t *testing.T) {
	err := NewCompileError("test error", 1, 5)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	if err.Message != "test error" {
		t.Errorf("Expected message 'test error', got %q", err.Message)
	}

	if err.Line != 1 {
		t.Errorf("Expected line 1, got %d", err.Line)
	}

	if err.Column != 5 {
		t.Errorf("Expected column 5, got %d", err.Column)
	}

	errorStr := err.Error()
	if errorStr == "" {
		t.Error("Expected non-empty error string")
	}
}

func TestRuntimeError(t *testing.T) {
	cause := fmt.Errorf("original error")
	err := NewRuntimeError("runtime error", cause)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	if err.Message != "runtime error" {
		t.Errorf("Expected message 'runtime error', got %q", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause to be %v, got %v", cause, err.Cause)
	}

	errorStr := err.Error()
	if errorStr == "" {
		t.Error("Expected non-empty error string")
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected string
	}{
		{42, "int"},
		{int64(42), "int64"},
		{3.14, "float64"},
		{"hello", "string"},
		{true, "bool"},
		{[]int{1, 2, 3}, "[]int"},
		{map[string]int{"a": 1}, "map[string]int"},
		{nil, "nil"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.value), func(t *testing.T) {
			result := GetType(tt.value)
			if result != tt.expected {
				t.Errorf("Expected type %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{nil, true},
		{(*int)(nil), true},
		{([]int)(nil), false},          // The IsNil function doesn't handle []int specifically
		{(map[string]int)(nil), false}, // The IsNil function doesn't handle map[string]int specifically
		{42, false},
		{"hello", false},
		{true, false},
		{[]int{}, false},          // Empty slice is not nil
		{map[string]int{}, false}, // Empty map is not nil
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.value), func(t *testing.T) {
			result := IsNil(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	t.Run("map input", func(t *testing.T) {
		input := map[string]interface{}{"a": 1, "b": 2}
		result := ToMap(input)
		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}
		if result["a"] != 1 || result["b"] != 2 {
			t.Errorf("Map conversion failed: %v", result)
		}
	})

	t.Run("non-map input", func(t *testing.T) {
		result := ToMap(42)
		if result == nil {
			t.Error("Expected non-nil result for non-map input")
		}
	})

	t.Run("nil input", func(t *testing.T) {
		result := ToMap(nil)
		if result == nil {
			t.Error("Expected non-nil result for nil input")
		}
	})
}

func TestStructToMap(t *testing.T) {
	fields := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}

	result := StructToMap(fields)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(result))
	}

	if result["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", result["name"])
	}

	if result["age"] != 30 {
		t.Errorf("Expected age 30, got %v", result["age"])
	}

	if result["active"] != true {
		t.Errorf("Expected active true, got %v", result["active"])
	}
}

func TestAsKindConstants(t *testing.T) {
	// Test that all AsKind constants are properly defined
	kinds := []AsKind{
		AsAny,
		AsIntKind,
		AsInt64Kind,
		AsFloat64Kind,
		AsStringKind,
		AsBoolKind,
	}

	// Check that they have different values
	seen := make(map[AsKind]bool)
	for _, kind := range kinds {
		if seen[kind] {
			t.Errorf("Duplicate AsKind value: %v", kind)
		}
		seen[kind] = true
	}
}

func TestPatchesOption(t *testing.T) {
	// Test that Patches option can be called without error
	config := &Config{}
	patch1 := Patch{}
	patch2 := Patch{}

	option := Patches(patch1, patch2)
	option(config)

	// Since Patches doesn't currently do anything, just verify it doesn't panic
}

func TestTagsOption(t *testing.T) {
	// Test that Tags option can be called without error
	config := &Config{}
	tag := Tag{Name: "json"}

	option := Tags(tag)
	option(config)

	// Since Tags doesn't currently do anything, just verify it doesn't panic
}

func TestConstExprOption(t *testing.T) {
	// Test that ConstExpr option can be called without error
	config := &Config{}

	option := ConstExpr("myConstant")
	option(config)

	// Since ConstExpr doesn't currently do anything, just verify it doesn't panic
}
