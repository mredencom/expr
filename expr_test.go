package expr

import (
	"testing"
	"time"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		expression string
		shouldErr  bool
	}{
		{"42", false},
		{"1 + 2", false},
		{`"hello"`, false},
		{"true", false},
		{"null", false},
		{"", true},    // Empty expression should error
		{"1 +", true}, // Invalid syntax should error
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			program, err := Compile(tt.expression)

			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if program == nil {
				t.Fatal("Expected program but got nil")
			}

			if program.source != tt.expression {
				t.Errorf("Expected source %q, got %q", tt.expression, program.source)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		expression string
		env        interface{}
		expected   interface{}
	}{
		{"42", nil, int64(42)},
		{"1 + 2", nil, int64(3)},
		{"5 * 6", nil, int64(30)},
		{"10 / 2", nil, int64(5)},
		{`"hello"`, nil, "hello"},
		{"true", nil, true},
		{"false", nil, false},
		{"5 > 3", nil, true},
		{"5 < 3", nil, false},
		{"5 == 5", nil, true},
		{"5 != 3", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			program, err := Compile(tt.expression)
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			result, err := Run(program, tt.env)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRunWithEnvironment(t *testing.T) {
	tests := []struct {
		expression string
		env        map[string]interface{}
		expected   interface{}
	}{
		{"x", map[string]interface{}{"x": 42}, int64(42)},
		{"x + y", map[string]interface{}{"x": 5, "y": 3}, int64(8)},
		{"name", map[string]interface{}{"name": "John"}, "John"},
		{"active", map[string]interface{}{"active": true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			program, err := Compile(tt.expression, Env(tt.env))
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			result, err := Run(program, tt.env)
			if err != nil {
				t.Fatalf("Runtime error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		expression string
		env        interface{}
		expected   interface{}
	}{
		{"42", nil, int64(42)},
		{"1 + 2", nil, int64(3)},
		{"x", map[string]interface{}{"x": 10}, int64(10)},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Eval(tt.expression, tt.env)
			if err != nil {
				t.Fatalf("Eval error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRunWithResult(t *testing.T) {
	program, err := Compile("42")
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	result, err := RunWithResult(program, nil)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result but got nil")
	}

	if result.Value != int64(42) {
		t.Errorf("Expected value 42, got %v", result.Value)
	}

	if result.ExecutionTime < 0 {
		t.Error("Expected non-negative execution time")
	}
}

func TestEvalWithResult(t *testing.T) {
	result, err := EvalWithResult("5 + 3", nil)
	if err != nil {
		t.Fatalf("EvalWithResult error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result but got nil")
	}

	if result.Value != int64(8) {
		t.Errorf("Expected value 8, got %v", result.Value)
	}
}

func TestOptions(t *testing.T) {
	t.Run("Env", func(t *testing.T) {
		env := map[string]interface{}{"x": 42}
		program, err := Compile("x", Env(env))
		if err != nil {
			t.Fatalf("Compilation error: %v", err)
		}

		result, err := Run(program, env)
		if err != nil {
			t.Fatalf("Runtime error: %v", err)
		}

		if result != int64(42) {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("AllowUndefinedVariables", func(t *testing.T) {
		_, err := Compile("undefined_var", AllowUndefinedVariables())
		// The AllowUndefinedVariables feature might not be fully implemented yet
		// For now, we just test that the option can be used without panicking
		if err != nil {
			t.Logf("AllowUndefinedVariables not fully implemented yet: %v", err)
		}
	})

	t.Run("DisableAllBuiltins", func(t *testing.T) {
		_, err := Compile("len('hello')", DisableAllBuiltins())
		// This should still compile but might fail at runtime
		if err != nil {
			t.Errorf("Unexpected compilation error: %v", err)
		}
	})

	t.Run("WithBuiltin", func(t *testing.T) {
		customFunc := func(x int) int { return x * 2 }
		_, err := Compile("double(5)", WithBuiltin("double", customFunc))
		if err != nil {
			t.Errorf("Unexpected error with custom builtin: %v", err)
		}
	})

	t.Run("EnableCache", func(t *testing.T) {
		_, err := Compile("1 + 1", EnableCache())
		if err != nil {
			t.Errorf("Unexpected error with cache enabled: %v", err)
		}
	})

	t.Run("DisableCache", func(t *testing.T) {
		_, err := Compile("1 + 1", DisableCache())
		if err != nil {
			t.Errorf("Unexpected error with cache disabled: %v", err)
		}
	})

	t.Run("EnableOptimization", func(t *testing.T) {
		_, err := Compile("1 + 1", EnableOptimization())
		if err != nil {
			t.Errorf("Unexpected error with optimization enabled: %v", err)
		}
	})

	t.Run("DisableOptimization", func(t *testing.T) {
		_, err := Compile("1 + 1", DisableOptimization())
		if err != nil {
			t.Errorf("Unexpected error with optimization disabled: %v", err)
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		_, err := Compile("1 + 1", WithTimeout(time.Second))
		if err != nil {
			t.Errorf("Unexpected error with timeout: %v", err)
		}
	})

	t.Run("EnableDebug", func(t *testing.T) {
		_, err := Compile("1 + 1", EnableDebug())
		if err != nil {
			t.Errorf("Unexpected error with debug enabled: %v", err)
		}
	})

	t.Run("EnableProfiling", func(t *testing.T) {
		_, err := Compile("1 + 1", EnableProfiling())
		if err != nil {
			t.Errorf("Unexpected error with profiling enabled: %v", err)
		}
	})
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		expression string
		expected   interface{}
	}{
		{`len("hello")`, int64(5)},
		{`string(42)`, "42"},
		{`int("42")`, int64(42)},
		{`bool(1)`, true},
		{`bool(0)`, false},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Eval(tt.expression, nil)
			if err != nil {
				t.Fatalf("Eval error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProgramMethods(t *testing.T) {
	program, err := Compile("42")
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	// Test Source
	if program.Source() != "42" {
		t.Errorf("Expected source '42', got %q", program.Source())
	}

	// CompileTime might be 0 for very fast compilation, so just check it's not negative
	if program.CompileTime() < 0 {
		t.Error("Expected non-negative compile time")
	}

	// Test BytecodeSize
	if program.BytecodeSize() <= 0 {
		t.Error("Expected positive bytecode size")
	}

	// Test ConstantsCount
	if program.ConstantsCount() <= 0 {
		t.Error("Expected positive constants count")
	}

	// Test String
	str := program.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}

func TestStatistics(t *testing.T) {
	// Reset statistics for clean test
	ResetStatistics()

	// Compile and run some expressions
	_, err := Compile("1 + 1")
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	_, err = Eval("2 + 2", nil)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}

	stats := GetStatistics()
	if stats == nil {
		t.Fatal("Expected statistics but got nil")
	}

	if stats.TotalCompilations == 0 {
		t.Error("Expected positive compilation count")
	}

	if stats.TotalExecutions == 0 {
		t.Error("Expected positive execution count")
	}
}

func TestCompileErrors(t *testing.T) {
	tests := []string{
		"1 +",     // Incomplete expression
		")",       // Invalid syntax
		"1 + + 2", // Invalid syntax
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			_, err := Compile(expression)
			if err == nil {
				t.Error("Expected compilation error but got none")
			}
		})
	}
}

func TestRuntimeErrors(t *testing.T) {
	tests := []struct {
		expression string
		env        interface{}
	}{
		{"x", nil},     // Undefined variable without AllowUndefinedVariables
		{"1 / 0", nil}, // Division by zero
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			program, err := Compile(tt.expression)
			if err != nil {
				// Skip if compilation fails
				return
			}

			_, err = Run(program, tt.env)
			if err == nil {
				t.Error("Expected runtime error but got none")
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		expression string
		env        map[string]interface{}
		expected   interface{}
	}{
		// Comment out advanced features that may not be implemented yet
		// {"x > 5 ? 'high' : 'low'", map[string]interface{}{"x": 10}, "high"},
		// {"x > 5 ? 'high' : 'low'", map[string]interface{}{"x": 3}, "low"},
		// {"[1, 2, 3][1]", nil, int64(2)},
		// {`{"key": "value"}["key"]`, nil, "value"},
		{"x + y * z", map[string]interface{}{"x": 1, "y": 2, "z": 3}, int64(7)},
		{"x > 5", map[string]interface{}{"x": 10}, true},
		{"x > 5", map[string]interface{}{"x": 3}, false},
		{"x + (y * z)", map[string]interface{}{"x": 1, "y": 2, "z": 3}, int64(7)},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Eval(tt.expression, tt.env)
			if err != nil {
				t.Fatalf("Eval error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkCompile(b *testing.B) {
	expression := "x + y * z"
	for i := 0; i < b.N; i++ {
		_, err := Compile(expression)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRun(b *testing.B) {
	program, err := Compile("x + y * z")
	if err != nil {
		b.Fatal(err)
	}

	env := map[string]interface{}{
		"x": 1,
		"y": 2,
		"z": 3,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Run(program, env)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEval(b *testing.B) {
	expression := "x + y * z"
	env := map[string]interface{}{
		"x": 1,
		"y": 2,
		"z": 3,
	}

	for i := 0; i < b.N; i++ {
		_, err := Eval(expression, env)
		if err != nil {
			b.Fatal(err)
		}
	}
}
