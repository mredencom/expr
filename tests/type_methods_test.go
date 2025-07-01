package tests

import (
	"math"
	"testing"

	"github.com/mredencom/expr/builtins"
	"github.com/mredencom/expr/types"
)

// TestStringTypeMethods tests all string type methods
func TestStringTypeMethods(t *testing.T) {
	t.Run("string.length", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected int64
			hasError bool
		}{
			{
				name:     "normal string",
				args:     []types.Value{types.NewString("hello")},
				expected: 5,
			},
			{
				name:     "empty string",
				args:     []types.Value{types.NewString("")},
				expected: 0,
			},
			{
				name:     "wrong number of args",
				args:     []types.Value{},
				hasError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["string.length"](tt.args)
				if tt.hasError {
					if err == nil {
						t.Error("Expected error but got none")
					}
					return
				}
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				intVal, ok := result.(*types.IntValue)
				if !ok {
					t.Fatalf("Expected IntValue, got %T", result)
				}
				if intVal.Value() != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, intVal.Value())
				}
			})
		}
	})

	t.Run("string.charAt", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.charAt"]([]types.Value{
			types.NewString("hello"), types.NewInt(1),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "e" {
			t.Errorf("Expected 'e', got %s", strVal.Value())
		}
	})

	t.Run("string.upper", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.upper"]([]types.Value{
			types.NewString("hello"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "HELLO" {
			t.Errorf("Expected 'HELLO', got %s", strVal.Value())
		}
	})

	t.Run("string.lower", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.lower"]([]types.Value{
			types.NewString("HELLO"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "hello" {
			t.Errorf("Expected 'hello', got %s", strVal.Value())
		}
	})

	t.Run("string.replace", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.replace"]([]types.Value{
			types.NewString("hello world"),
			types.NewString("world"),
			types.NewString("universe"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "hello universe" {
			t.Errorf("Expected 'hello universe', got %s", strVal.Value())
		}
	})

	t.Run("string.split", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.split"]([]types.Value{
			types.NewString("a,b,c"),
			types.NewString(","),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		sliceVal, ok := result.(*types.SliceValue)
		if !ok {
			t.Fatalf("Expected SliceValue, got %T", result)
		}
		if sliceVal.Len() != 3 {
			t.Errorf("Expected 3 elements, got %d", sliceVal.Len())
		}
	})

	t.Run("string.contains", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.contains"]([]types.Value{
			types.NewString("hello world"),
			types.NewString("world"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		boolVal, ok := result.(*types.BoolValue)
		if !ok {
			t.Fatalf("Expected BoolValue, got %T", result)
		}
		if !boolVal.Value() {
			t.Error("Expected true, got false")
		}
	})

	t.Run("string.indexOf", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.indexOf"]([]types.Value{
			types.NewString("hello world"),
			types.NewString("world"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 6 {
			t.Errorf("Expected 6, got %d", intVal.Value())
		}
	})

	t.Run("string.startsWith", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.startsWith"]([]types.Value{
			types.NewString("hello world"),
			types.NewString("hello"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		boolVal, ok := result.(*types.BoolValue)
		if !ok {
			t.Fatalf("Expected BoolValue, got %T", result)
		}
		if !boolVal.Value() {
			t.Error("Expected true, got false")
		}
	})

	t.Run("string.endsWith", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["string.endsWith"]([]types.Value{
			types.NewString("hello world"),
			types.NewString("world"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		boolVal, ok := result.(*types.BoolValue)
		if !ok {
			t.Fatalf("Expected BoolValue, got %T", result)
		}
		if !boolVal.Value() {
			t.Error("Expected true, got false")
		}
	})
}

// TestIntTypeMethods tests all int type methods
func TestIntTypeMethods(t *testing.T) {
	t.Run("int.abs", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected int64
		}{
			{
				name:     "positive number",
				args:     []types.Value{types.NewInt(5)},
				expected: 5,
			},
			{
				name:     "negative number",
				args:     []types.Value{types.NewInt(-5)},
				expected: 5,
			},
			{
				name:     "zero",
				args:     []types.Value{types.NewInt(0)},
				expected: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.abs"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				intVal, ok := result.(*types.IntValue)
				if !ok {
					t.Fatalf("Expected IntValue, got %T", result)
				}
				if intVal.Value() != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, intVal.Value())
				}
			})
		}
	})

	t.Run("int.sign", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected int64
		}{
			{
				name:     "positive number",
				args:     []types.Value{types.NewInt(5)},
				expected: 1,
			},
			{
				name:     "negative number",
				args:     []types.Value{types.NewInt(-5)},
				expected: -1,
			},
			{
				name:     "zero",
				args:     []types.Value{types.NewInt(0)},
				expected: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.sign"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				intVal, ok := result.(*types.IntValue)
				if !ok {
					t.Fatalf("Expected IntValue, got %T", result)
				}
				if intVal.Value() != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, intVal.Value())
				}
			})
		}
	})

	t.Run("int.toString", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["int.toString"]([]types.Value{
			types.NewInt(42),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "42" {
			t.Errorf("Expected '42', got %s", strVal.Value())
		}
	})

	t.Run("int.isEven", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "even number",
				args:     []types.Value{types.NewInt(4)},
				expected: true,
			},
			{
				name:     "odd number",
				args:     []types.Value{types.NewInt(5)},
				expected: false,
			},
			{
				name:     "zero",
				args:     []types.Value{types.NewInt(0)},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.isEven"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("int.isPrime", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "prime number 2",
				args:     []types.Value{types.NewInt(2)},
				expected: true,
			},
			{
				name:     "prime number 7",
				args:     []types.Value{types.NewInt(7)},
				expected: true,
			},
			{
				name:     "non-prime number 4",
				args:     []types.Value{types.NewInt(4)},
				expected: false,
			},
			{
				name:     "non-prime number 1",
				args:     []types.Value{types.NewInt(1)},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.isPrime"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("int.factorial", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected int64
			hasError bool
		}{
			{
				name:     "factorial 0",
				args:     []types.Value{types.NewInt(0)},
				expected: 1,
			},
			{
				name:     "factorial 5",
				args:     []types.Value{types.NewInt(5)},
				expected: 120,
			},
			{
				name:     "negative number",
				args:     []types.Value{types.NewInt(-1)},
				hasError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.factorial"](tt.args)
				if tt.hasError {
					if err == nil {
						t.Error("Expected error but got none")
					}
					return
				}
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				intVal, ok := result.(*types.IntValue)
				if !ok {
					t.Fatalf("Expected IntValue, got %T", result)
				}
				if intVal.Value() != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, intVal.Value())
				}
			})
		}
	})

	t.Run("int.toFloat", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["int.toFloat"]([]types.Value{
			types.NewInt(42),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		floatVal, ok := result.(*types.FloatValue)
		if !ok {
			t.Fatalf("Expected FloatValue, got %T", result)
		}
		if floatVal.Value() != 42.0 {
			t.Errorf("Expected 42.0, got %f", floatVal.Value())
		}
	})

	t.Run("int.toBool", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "non-zero to true",
				args:     []types.Value{types.NewInt(42)},
				expected: true,
			},
			{
				name:     "zero to false",
				args:     []types.Value{types.NewInt(0)},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["int.toBool"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("int.min", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["int.min"]([]types.Value{
			types.NewInt(5), types.NewInt(3),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 3 {
			t.Errorf("Expected 3, got %d", intVal.Value())
		}
	})

	t.Run("int.max", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["int.max"]([]types.Value{
			types.NewInt(5), types.NewInt(3),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 5 {
			t.Errorf("Expected 5, got %d", intVal.Value())
		}
	})

	t.Run("int.clamp", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["int.clamp"]([]types.Value{
			types.NewInt(15), types.NewInt(10), types.NewInt(20),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 15 {
			t.Errorf("Expected 15, got %d", intVal.Value())
		}
	})
}

// TestFloatTypeMethods tests all float type methods
func TestFloatTypeMethods(t *testing.T) {
	t.Run("float.abs", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected float64
		}{
			{
				name:     "positive number",
				args:     []types.Value{types.NewFloat(3.14)},
				expected: 3.14,
			},
			{
				name:     "negative number",
				args:     []types.Value{types.NewFloat(-3.14)},
				expected: 3.14,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["float.abs"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				floatVal, ok := result.(*types.FloatValue)
				if !ok {
					t.Fatalf("Expected FloatValue, got %T", result)
				}
				if floatVal.Value() != tt.expected {
					t.Errorf("Expected %f, got %f", tt.expected, floatVal.Value())
				}
			})
		}
	})

	t.Run("float.round", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["float.round"]([]types.Value{
			types.NewFloat(3.14),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		floatVal, ok := result.(*types.FloatValue)
		if !ok {
			t.Fatalf("Expected FloatValue, got %T", result)
		}
		if floatVal.Value() != 3.0 {
			t.Errorf("Expected 3.0, got %f", floatVal.Value())
		}
	})

	t.Run("float.ceil", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["float.ceil"]([]types.Value{
			types.NewFloat(3.14),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		floatVal, ok := result.(*types.FloatValue)
		if !ok {
			t.Fatalf("Expected FloatValue, got %T", result)
		}
		if floatVal.Value() != 4.0 {
			t.Errorf("Expected 4.0, got %f", floatVal.Value())
		}
	})

	t.Run("float.floor", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["float.floor"]([]types.Value{
			types.NewFloat(3.14),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		floatVal, ok := result.(*types.FloatValue)
		if !ok {
			t.Fatalf("Expected FloatValue, got %T", result)
		}
		if floatVal.Value() != 3.0 {
			t.Errorf("Expected 3.0, got %f", floatVal.Value())
		}
	})

	t.Run("float.isNaN", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "normal number",
				args:     []types.Value{types.NewFloat(3.14)},
				expected: false,
			},
			{
				name:     "NaN",
				args:     []types.Value{types.NewFloat(math.NaN())},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["float.isNaN"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("float.isInf", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "normal number",
				args:     []types.Value{types.NewFloat(3.14)},
				expected: false,
			},
			{
				name:     "positive infinity",
				args:     []types.Value{types.NewFloat(math.Inf(1))},
				expected: true,
			},
			{
				name:     "negative infinity",
				args:     []types.Value{types.NewFloat(math.Inf(-1))},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["float.isInf"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("float.toString", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["float.toString"]([]types.Value{
			types.NewFloat(3.14),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "3.14" {
			t.Errorf("Expected '3.14', got %s", strVal.Value())
		}
	})

	t.Run("float.toInt", func(t *testing.T) {
		result, err := builtins.TypeMethodBuiltins["float.toInt"]([]types.Value{
			types.NewFloat(3.14),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 3 {
			t.Errorf("Expected 3, got %d", intVal.Value())
		}
	})

	t.Run("float.sign", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected float64
		}{
			{
				name:     "positive",
				args:     []types.Value{types.NewFloat(3.14)},
				expected: 1.0,
			},
			{
				name:     "negative",
				args:     []types.Value{types.NewFloat(-3.14)},
				expected: -1.0,
			},
			{
				name:     "zero",
				args:     []types.Value{types.NewFloat(0.0)},
				expected: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["float.sign"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				floatVal, ok := result.(*types.FloatValue)
				if !ok {
					t.Fatalf("Expected FloatValue, got %T", result)
				}
				if floatVal.Value() != tt.expected {
					t.Errorf("Expected %f, got %f", tt.expected, floatVal.Value())
				}
			})
		}
	})
}

// TestBoolTypeMethods tests all bool type methods
func TestBoolTypeMethods(t *testing.T) {
	t.Run("bool.toString", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected string
		}{
			{
				name:     "true",
				args:     []types.Value{types.NewBool(true)},
				expected: "true",
			},
			{
				name:     "false",
				args:     []types.Value{types.NewBool(false)},
				expected: "false",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["bool.toString"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				strVal, ok := result.(*types.StringValue)
				if !ok {
					t.Fatalf("Expected StringValue, got %T", result)
				}
				if strVal.Value() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, strVal.Value())
				}
			})
		}
	})

	t.Run("bool.not", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "not true",
				args:     []types.Value{types.NewBool(true)},
				expected: false,
			},
			{
				name:     "not false",
				args:     []types.Value{types.NewBool(false)},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["bool.not"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("bool.and", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "true and true",
				args:     []types.Value{types.NewBool(true), types.NewBool(true)},
				expected: true,
			},
			{
				name:     "true and false",
				args:     []types.Value{types.NewBool(true), types.NewBool(false)},
				expected: false,
			},
			{
				name:     "false and false",
				args:     []types.Value{types.NewBool(false), types.NewBool(false)},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["bool.and"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("bool.or", func(t *testing.T) {
		tests := []struct {
			name     string
			args     []types.Value
			expected bool
		}{
			{
				name:     "true or false",
				args:     []types.Value{types.NewBool(true), types.NewBool(false)},
				expected: true,
			},
			{
				name:     "false or false",
				args:     []types.Value{types.NewBool(false), types.NewBool(false)},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["bool.or"](tt.args)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})
}

// TestSliceTypeMethods tests all slice type methods
func TestSliceTypeMethods(t *testing.T) {
	t.Run("slice.length", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.length"]([]types.Value{slice})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 3 {
			t.Errorf("Expected 3, got %d", intVal.Value())
		}
	})

	t.Run("slice.isEmpty", func(t *testing.T) {
		tests := []struct {
			name     string
			slice    *types.SliceValue
			expected bool
		}{
			{
				name: "non-empty slice",
				slice: types.NewSlice([]types.Value{types.NewInt(1)},
					types.TypeInfo{Kind: types.KindInt64, Name: "int"}),
				expected: false,
			},
			{
				name: "empty slice",
				slice: types.NewSlice([]types.Value{},
					types.TypeInfo{Kind: types.KindInt64, Name: "int"}),
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["slice.isEmpty"]([]types.Value{tt.slice})
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("slice.first", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.first"]([]types.Value{slice})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 1 {
			t.Errorf("Expected 1, got %d", intVal.Value())
		}
	})

	t.Run("slice.last", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.last"]([]types.Value{slice})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 3 {
			t.Errorf("Expected 3, got %d", intVal.Value())
		}
	})

	t.Run("slice.get", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.get"]([]types.Value{
			slice, types.NewInt(1),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 2 {
			t.Errorf("Expected 2, got %d", intVal.Value())
		}
	})

	t.Run("slice.contains", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.contains"]([]types.Value{
			slice, types.NewInt(2),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		boolVal, ok := result.(*types.BoolValue)
		if !ok {
			t.Fatalf("Expected BoolValue, got %T", result)
		}
		if !boolVal.Value() {
			t.Error("Expected true, got false")
		}
	})

	t.Run("slice.reverse", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{
			types.NewInt(1), types.NewInt(2), types.NewInt(3),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.reverse"]([]types.Value{slice})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		resultSlice, ok := result.(*types.SliceValue)
		if !ok {
			t.Fatalf("Expected SliceValue, got %T", result)
		}
		if resultSlice.Len() != 3 {
			t.Errorf("Expected length 3, got %d", resultSlice.Len())
		}
		firstVal, ok := resultSlice.Get(0).(*types.IntValue)
		if !ok || firstVal.Value() != 3 {
			t.Errorf("Expected first element to be 3, got %v", resultSlice.Get(0))
		}
	})

	t.Run("slice.join", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		slice := types.NewSlice([]types.Value{
			types.NewString("a"), types.NewString("b"), types.NewString("c"),
		}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.join"]([]types.Value{
			slice, types.NewString(","),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		strVal, ok := result.(*types.StringValue)
		if !ok {
			t.Fatalf("Expected StringValue, got %T", result)
		}
		if strVal.Value() != "a,b,c" {
			t.Errorf("Expected 'a,b,c', got %s", strVal.Value())
		}
	})
}

// TestMapTypeMethods tests all map type methods
func TestMapTypeMethods(t *testing.T) {
	t.Run("map.size", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
			"b": types.NewInt(2),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.size"]([]types.Value{mapVal})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 2 {
			t.Errorf("Expected 2, got %d", intVal.Value())
		}
	})

	t.Run("map.isEmpty", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}

		tests := []struct {
			name     string
			mapVal   *types.MapValue
			expected bool
		}{
			{
				name: "non-empty map",
				mapVal: types.NewMap(map[string]types.Value{
					"a": types.NewInt(1),
				}, keyType, valType),
				expected: false,
			},
			{
				name:     "empty map",
				mapVal:   types.NewMap(map[string]types.Value{}, keyType, valType),
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := builtins.TypeMethodBuiltins["map.isEmpty"]([]types.Value{tt.mapVal})
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				boolVal, ok := result.(*types.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", result)
				}
				if boolVal.Value() != tt.expected {
					t.Errorf("Expected %t, got %t", tt.expected, boolVal.Value())
				}
			})
		}
	})

	t.Run("map.has", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.has"]([]types.Value{
			mapVal, types.NewString("a"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		boolVal, ok := result.(*types.BoolValue)
		if !ok {
			t.Fatalf("Expected BoolValue, got %T", result)
		}
		if !boolVal.Value() {
			t.Error("Expected true, got false")
		}
	})

	t.Run("map.get", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.get"]([]types.Value{
			mapVal, types.NewString("a"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		intVal, ok := result.(*types.IntValue)
		if !ok {
			t.Fatalf("Expected IntValue, got %T", result)
		}
		if intVal.Value() != 1 {
			t.Errorf("Expected 1, got %d", intVal.Value())
		}
	})

	t.Run("map.keys", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
			"b": types.NewInt(2),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.keys"]([]types.Value{mapVal})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		sliceVal, ok := result.(*types.SliceValue)
		if !ok {
			t.Fatalf("Expected SliceValue, got %T", result)
		}
		if sliceVal.Len() != 2 {
			t.Errorf("Expected 2 keys, got %d", sliceVal.Len())
		}
	})

	t.Run("map.values", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
			"b": types.NewInt(2),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.values"]([]types.Value{mapVal})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		sliceVal, ok := result.(*types.SliceValue)
		if !ok {
			t.Fatalf("Expected SliceValue, got %T", result)
		}
		if sliceVal.Len() != 2 {
			t.Errorf("Expected 2 values, got %d", sliceVal.Len())
		}
	})
}

// TestTypeMethodErrorCases tests error handling in type methods
func TestTypeMethodErrorCases(t *testing.T) {
	t.Run("wrong argument count", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["string.length"]([]types.Value{})
		if err == nil {
			t.Error("Expected error for wrong argument count")
		}
	})

	t.Run("wrong argument type", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["string.length"]([]types.Value{types.NewInt(123)})
		if err == nil {
			t.Error("Expected error for wrong argument type")
		}
	})

	t.Run("negative factorial", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["int.factorial"]([]types.Value{types.NewInt(-5)})
		if err == nil {
			t.Error("Expected error for negative factorial")
		}
	})

	t.Run("invalid base for toString", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["int.toString"]([]types.Value{
			types.NewInt(42), types.NewInt(1), // base 1 is invalid
		})
		if err == nil {
			t.Error("Expected error for invalid base")
		}
	})

	t.Run("empty slice first", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		emptySlice := types.NewSlice([]types.Value{}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.first"]([]types.Value{emptySlice})
		if err != nil {
			// Expected error - this is correct behavior
			return
		}
		// If no error, check if it returns nil/empty value
		if result == nil {
			t.Error("Expected result not to be nil")
		}
	})

	t.Run("empty slice last", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		emptySlice := types.NewSlice([]types.Value{}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.last"]([]types.Value{emptySlice})
		if err != nil {
			// Expected error - this is correct behavior
			return
		}
		// If no error, check if it returns nil/empty value
		if result == nil {
			t.Error("Expected result not to be nil")
		}
	})

	t.Run("slice index out of bounds", func(t *testing.T) {
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		slice := types.NewSlice([]types.Value{types.NewInt(1)}, elemType)

		result, err := builtins.TypeMethodBuiltins["slice.get"]([]types.Value{slice, types.NewInt(10)})
		if err != nil {
			// Expected error - this is correct behavior
			return
		}
		// If no error, check if it returns nil/empty value
		if result == nil {
			t.Error("Expected result not to be nil")
		}
	})

	t.Run("large factorial", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["int.factorial"]([]types.Value{types.NewInt(25)})
		if err == nil {
			t.Error("Expected error for large factorial")
		}
	})

	t.Run("map get non-existent key", func(t *testing.T) {
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
		valType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
		mapVal := types.NewMap(map[string]types.Value{
			"a": types.NewInt(1),
		}, keyType, valType)

		result, err := builtins.TypeMethodBuiltins["map.get"]([]types.Value{
			mapVal, types.NewString("nonexistent"),
		})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Should return nil/empty value for non-existent key
		if result == nil {
			t.Error("Expected non-nil result for non-existent key")
		}
	})

	t.Run("wrong number of arguments for binary operations", func(t *testing.T) {
		_, err := builtins.TypeMethodBuiltins["bool.and"]([]types.Value{types.NewBool(true)})
		if err == nil {
			t.Error("Expected error for missing second argument")
		}
	})
}
