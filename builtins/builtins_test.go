package builtins

import (
	"testing"

	"github.com/mredencom/expr/types"
)

// Test type conversion builtins

func TestLenBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected int64
		hasError bool
	}{
		{
			name:     "string length",
			args:     []types.Value{types.NewString("hello")},
			expected: 5,
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: 0,
		},
		{
			name: "slice length",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: 3,
		},
		{
			name: "map length",
			args: []types.Value{types.NewMap(
				map[string]types.Value{"a": types.NewInt(1), "b": types.NewInt(2)},
				types.TypeInfo{Kind: types.KindString, Name: "string"},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: 2,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
		{
			name:     "unsupported type",
			args:     []types.Value{types.NewInt(42)},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lenBuiltin(tt.args)
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
}

func TestStringBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected string
		hasError bool
	}{
		{
			name:     "int to string",
			args:     []types.Value{types.NewInt(42)},
			expected: "42",
		},
		{
			name:     "float to string",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: "3.14",
		},
		{
			name:     "bool to string",
			args:     []types.Value{types.NewBool(true)},
			expected: "true",
		},
		{
			name:     "string to string",
			args:     []types.Value{types.NewString("test")},
			expected: "test",
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := stringBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestIntBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected int64
		hasError bool
	}{
		{
			name:     "int to int",
			args:     []types.Value{types.NewInt(42)},
			expected: 42,
		},
		{
			name:     "float to int",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: 3,
		},
		{
			name:     "string to int",
			args:     []types.Value{types.NewString("123")},
			expected: 123,
		},
		{
			name:     "bool true to int",
			args:     []types.Value{types.NewBool(true)},
			expected: 1,
		},
		{
			name:     "bool false to int",
			args:     []types.Value{types.NewBool(false)},
			expected: 0,
		},
		{
			name:     "invalid string",
			args:     []types.Value{types.NewString("abc")},
			hasError: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := intBuiltin(tt.args)
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
}

func TestFloatBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected float64
		hasError bool
	}{
		{
			name:     "float to float",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: 3.14,
		},
		{
			name:     "int to float",
			args:     []types.Value{types.NewInt(42)},
			expected: 42.0,
		},
		{
			name:     "string to float",
			args:     []types.Value{types.NewString("3.14")},
			expected: 3.14,
		},
		{
			name:     "invalid string",
			args:     []types.Value{types.NewString("abc")},
			hasError: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := floatBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestBoolBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name:     "bool to bool",
			args:     []types.Value{types.NewBool(true)},
			expected: true,
		},
		{
			name:     "non-zero int to bool",
			args:     []types.Value{types.NewInt(42)},
			expected: true,
		},
		{
			name:     "zero int to bool",
			args:     []types.Value{types.NewInt(0)},
			expected: false,
		},
		{
			name:     "non-zero float to bool",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: true,
		},
		{
			name:     "zero float to bool",
			args:     []types.Value{types.NewFloat(0.0)},
			expected: false,
		},
		{
			name:     "non-empty string to bool",
			args:     []types.Value{types.NewString("test")},
			expected: true,
		},
		{
			name:     "empty string to bool",
			args:     []types.Value{types.NewString("")},
			expected: false,
		},
		{
			name:     "nil to bool",
			args:     []types.Value{types.NewNil()},
			expected: false,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := boolBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			boolVal, ok := result.(*types.BoolValue)
			if !ok {
				t.Fatalf("Expected BoolValue, got %T", result)
			}
			if boolVal.Value() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolVal.Value())
			}
		})
	}
}

// Test math builtins

func TestAbsBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name:     "positive int",
			args:     []types.Value{types.NewInt(42)},
			expected: types.NewInt(42),
		},
		{
			name:     "negative int",
			args:     []types.Value{types.NewInt(-42)},
			expected: types.NewInt(42),
		},
		{
			name:     "positive float",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: types.NewFloat(3.14),
		},
		{
			name:     "negative float",
			args:     []types.Value{types.NewFloat(-3.14)},
			expected: types.NewFloat(3.14),
		},
		{
			name:     "zero",
			args:     []types.Value{types.NewInt(0)},
			expected: types.NewInt(0),
		},
		{
			name:     "unsupported type",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := absBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMaxBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name:     "integers",
			args:     []types.Value{types.NewInt(1), types.NewInt(5), types.NewInt(3)},
			expected: types.NewInt(5),
		},
		{
			name:     "floats",
			args:     []types.Value{types.NewFloat(1.1), types.NewFloat(5.5), types.NewFloat(3.3)},
			expected: types.NewFloat(5.5),
		},
		{
			name:     "strings",
			args:     []types.Value{types.NewString("apple"), types.NewString("zebra"), types.NewString("banana")},
			expected: types.NewString("zebra"),
		},
		{
			name:     "single value",
			args:     []types.Value{types.NewInt(42)},
			expected: types.NewInt(42),
		},
		{
			name:     "no args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := maxBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMinBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name:     "integers",
			args:     []types.Value{types.NewInt(5), types.NewInt(1), types.NewInt(3)},
			expected: types.NewInt(1),
		},
		{
			name:     "floats",
			args:     []types.Value{types.NewFloat(5.5), types.NewFloat(1.1), types.NewFloat(3.3)},
			expected: types.NewFloat(1.1),
		},
		{
			name:     "strings",
			args:     []types.Value{types.NewString("zebra"), types.NewString("apple"), types.NewString("banana")},
			expected: types.NewString("apple"),
		},
		{
			name:     "single value",
			args:     []types.Value{types.NewInt(42)},
			expected: types.NewInt(42),
		},
		{
			name:     "no args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := minBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test string builtins

func TestContainsBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name:     "contains substring",
			args:     []types.Value{types.NewString("hello world"), types.NewString("world")},
			expected: true,
		},
		{
			name:     "does not contain substring",
			args:     []types.Value{types.NewString("hello world"), types.NewString("foo")},
			expected: false,
		},
		{
			name:     "empty substring",
			args:     []types.Value{types.NewString("hello"), types.NewString("")},
			expected: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
		{
			name:     "first arg not string",
			args:     []types.Value{types.NewInt(42), types.NewString("test")},
			hasError: true,
		},
		{
			name:     "second arg not string",
			args:     []types.Value{types.NewString("test"), types.NewInt(42)},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := containsBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			boolVal, ok := result.(*types.BoolValue)
			if !ok {
				t.Fatalf("Expected BoolValue, got %T", result)
			}
			if boolVal.Value() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolVal.Value())
			}
		})
	}
}

func TestStartsWithBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name:     "starts with prefix",
			args:     []types.Value{types.NewString("hello world"), types.NewString("hello")},
			expected: true,
		},
		{
			name:     "does not start with prefix",
			args:     []types.Value{types.NewString("hello world"), types.NewString("world")},
			expected: false,
		},
		{
			name:     "empty prefix",
			args:     []types.Value{types.NewString("hello"), types.NewString("")},
			expected: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := startsWithBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			boolVal, ok := result.(*types.BoolValue)
			if !ok {
				t.Fatalf("Expected BoolValue, got %T", result)
			}
			if boolVal.Value() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolVal.Value())
			}
		})
	}
}

func TestEndsWithBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name:     "ends with suffix",
			args:     []types.Value{types.NewString("hello world"), types.NewString("world")},
			expected: true,
		},
		{
			name:     "does not end with suffix",
			args:     []types.Value{types.NewString("hello world"), types.NewString("hello")},
			expected: false,
		},
		{
			name:     "empty suffix",
			args:     []types.Value{types.NewString("hello"), types.NewString("")},
			expected: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := endsWithBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			boolVal, ok := result.(*types.BoolValue)
			if !ok {
				t.Fatalf("Expected BoolValue, got %T", result)
			}
			if boolVal.Value() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolVal.Value())
			}
		})
	}
}

func TestUpperBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected string
		hasError bool
	}{
		{
			name:     "lowercase to uppercase",
			args:     []types.Value{types.NewString("hello")},
			expected: "HELLO",
		},
		{
			name:     "mixed case to uppercase",
			args:     []types.Value{types.NewString("Hello World")},
			expected: "HELLO WORLD",
		},
		{
			name:     "already uppercase",
			args:     []types.Value{types.NewString("HELLO")},
			expected: "HELLO",
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: "",
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
		{
			name:     "not a string",
			args:     []types.Value{types.NewInt(42)},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := upperBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestLowerBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected string
		hasError bool
	}{
		{
			name:     "uppercase to lowercase",
			args:     []types.Value{types.NewString("HELLO")},
			expected: "hello",
		},
		{
			name:     "mixed case to lowercase",
			args:     []types.Value{types.NewString("Hello World")},
			expected: "hello world",
		},
		{
			name:     "already lowercase",
			args:     []types.Value{types.NewString("hello")},
			expected: "hello",
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: "",
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lowerBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestTrimBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected string
		hasError bool
	}{
		{
			name:     "trim spaces",
			args:     []types.Value{types.NewString("  hello  ")},
			expected: "hello",
		},
		{
			name:     "trim tabs and newlines",
			args:     []types.Value{types.NewString("\t\nhello\n\t")},
			expected: "hello",
		},
		{
			name:     "no whitespace to trim",
			args:     []types.Value{types.NewString("hello")},
			expected: "hello",
		},
		{
			name:     "only whitespace",
			args:     []types.Value{types.NewString("   ")},
			expected: "",
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: "",
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := trimBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestTypeBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected string
		hasError bool
	}{
		{
			name:     "int type",
			args:     []types.Value{types.NewInt(42)},
			expected: "int",
		},
		{
			name:     "float type",
			args:     []types.Value{types.NewFloat(3.14)},
			expected: "float",
		},
		{
			name:     "string type",
			args:     []types.Value{types.NewString("hello")},
			expected: "string",
		},
		{
			name:     "bool type",
			args:     []types.Value{types.NewBool(true)},
			expected: "bool",
		},
		{
			name:     "nil type",
			args:     []types.Value{types.NewNil()},
			expected: "nil",
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := typeBuiltin(tt.args)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
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
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a, b     types.Value
		expected int
	}{
		{
			name:     "int equal",
			a:        types.NewInt(5),
			b:        types.NewInt(5),
			expected: 0,
		},
		{
			name:     "int a less than b",
			a:        types.NewInt(3),
			b:        types.NewInt(5),
			expected: -1,
		},
		{
			name:     "int a greater than b",
			a:        types.NewInt(7),
			b:        types.NewInt(5),
			expected: 1,
		},
		{
			name:     "float equal",
			a:        types.NewFloat(3.14),
			b:        types.NewFloat(3.14),
			expected: 0,
		},
		{
			name:     "string equal",
			a:        types.NewString("apple"),
			b:        types.NewString("apple"),
			expected: 0,
		},
		{
			name:     "string a less than b",
			a:        types.NewString("apple"),
			b:        types.NewString("banana"),
			expected: -1,
		},
		{
			name:     "different types",
			a:        types.NewInt(5),
			b:        types.NewString("test"),
			expected: 0, // Different types return 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
