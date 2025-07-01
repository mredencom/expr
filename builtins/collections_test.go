package builtins

import (
	"testing"

	"github.com/mredencom/expr/types"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name: "all true booleans",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewBool(true), types.NewBool(true), types.NewBool(true)},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: true,
		},
		{
			name: "one false boolean",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewBool(true), types.NewBool(false), types.NewBool(true)},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: false,
		},
		{
			name: "all non-zero integers",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: true,
		},
		{
			name: "one zero integer",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(0), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: false,
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
		{
			name:     "not a collection",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := All(tt.args)
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

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected bool
		hasError bool
	}{
		{
			name: "one true boolean",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewBool(false), types.NewBool(true), types.NewBool(false)},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: true,
		},
		{
			name: "all false booleans",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewBool(false), types.NewBool(false), types.NewBool(false)},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: false,
		},
		{
			name: "one non-zero integer",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(0), types.NewInt(1), types.NewInt(0)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: true,
		},
		{
			name: "all zero integers",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(0), types.NewInt(0), types.NewInt(0)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: false,
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			)},
			expected: false,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
		{
			name:     "not a collection",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Any(tt.args)
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

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name: "sum integers",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewInt(6),
		},
		{
			name: "sum floats",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewFloat(1.1), types.NewFloat(2.2), types.NewFloat(3.3)},
				types.TypeInfo{Kind: types.KindFloat64, Name: "float"},
			)},
			expected: types.NewFloat(6.6),
		},
		{
			name: "sum mixed int and float",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewFloat(2.5), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"},
			)},
			expected: types.NewFloat(6.5),
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewInt(0),
		},
		{
			name: "non-numeric values",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewString("test"), types.NewInt(1)},
				types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"},
			)},
			hasError: true,
		},
		{
			name:     "wrong number of args",
			args:     []types.Value{},
			hasError: true,
		},
		{
			name:     "not a collection",
			args:     []types.Value{types.NewString("test")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Sum(tt.args)
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

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected int64
		hasError bool
	}{
		{
			name: "count slice elements",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: 3,
		},
		{
			name: "count map elements",
			args: []types.Value{types.NewMap(
				map[string]types.Value{"a": types.NewInt(1), "b": types.NewInt(2)},
				types.TypeInfo{Kind: types.KindString, Name: "string"},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: 2,
		},
		{
			name:     "count string characters",
			args:     []types.Value{types.NewString("hello")},
			expected: 5,
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: 0,
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
		{
			name:     "unsupported type",
			args:     []types.Value{types.NewInt(42)},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Count(tt.args)
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

func TestFirst(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name: "first of slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewInt(1),
		},
		{
			name:     "first character of string",
			args:     []types.Value{types.NewString("hello")},
			expected: types.NewString("h"),
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewNil(),
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: types.NewString(""),
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
			result, err := First(tt.args)
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

func TestLast(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		hasError bool
	}{
		{
			name: "last of slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewInt(3),
		},
		{
			name:     "last character of string",
			args:     []types.Value{types.NewString("hello")},
			expected: types.NewString("o"),
		},
		{
			name: "empty slice",
			args: []types.Value{types.NewSlice(
				[]types.Value{},
				types.TypeInfo{Kind: types.KindInt64, Name: "int"},
			)},
			expected: types.NewNil(),
		},
		{
			name:     "empty string",
			args:     []types.Value{types.NewString("")},
			expected: types.NewString(""),
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
			result, err := Last(tt.args)
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

func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		value    types.Value
		expected bool
	}{
		{
			name:     "zero int",
			value:    types.NewInt(0),
			expected: true,
		},
		{
			name:     "non-zero int",
			value:    types.NewInt(42),
			expected: false,
		},
		{
			name:     "zero float",
			value:    types.NewFloat(0.0),
			expected: true,
		},
		{
			name:     "non-zero float",
			value:    types.NewFloat(3.14),
			expected: false,
		},
		{
			name:     "empty string",
			value:    types.NewString(""),
			expected: true,
		},
		{
			name:     "non-empty string",
			value:    types.NewString("hello"),
			expected: false,
		},
		{
			name:     "false bool",
			value:    types.NewBool(false),
			expected: true,
		},
		{
			name:     "true bool",
			value:    types.NewBool(true),
			expected: false,
		},
		{
			name:     "nil value",
			value:    types.NewNil(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isZeroValue(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		value    types.Value
		expected bool
	}{
		{
			name:     "non-zero int is truthy",
			value:    types.NewInt(42),
			expected: true,
		},
		{
			name:     "zero int is falsy",
			value:    types.NewInt(0),
			expected: false,
		},
		{
			name:     "non-zero float is truthy",
			value:    types.NewFloat(3.14),
			expected: true,
		},
		{
			name:     "zero float is falsy",
			value:    types.NewFloat(0.0),
			expected: false,
		},
		{
			name:     "non-empty string is truthy",
			value:    types.NewString("hello"),
			expected: true,
		},
		{
			name:     "empty string is falsy",
			value:    types.NewString(""),
			expected: false,
		},
		{
			name:     "true bool is truthy",
			value:    types.NewBool(true),
			expected: true,
		},
		{
			name:     "false bool is falsy",
			value:    types.NewBool(false),
			expected: false,
		},
		{
			name:     "nil is falsy",
			value:    types.NewNil(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTruthy(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
