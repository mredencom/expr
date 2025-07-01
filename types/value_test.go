package types

import (
	"testing"
)

// TestBoolValue tests BoolValue implementation
func TestBoolValue(t *testing.T) {
	t.Run("NewBool", func(t *testing.T) {
		trueVal := NewBool(true)
		falseVal := NewBool(false)

		if trueVal.Value() != true {
			t.Errorf("Expected true, got %v", trueVal.Value())
		}
		if falseVal.Value() != false {
			t.Errorf("Expected false, got %v", falseVal.Value())
		}
	})

	t.Run("Type", func(t *testing.T) {
		val := NewBool(true)
		typeInfo := val.Type()

		if typeInfo.Kind != KindBool {
			t.Errorf("Expected KindBool, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "bool" {
			t.Errorf("Expected 'bool', got %s", typeInfo.Name)
		}
		if typeInfo.Size != 1 {
			t.Errorf("Expected size 1, got %d", typeInfo.Size)
		}
	})

	t.Run("String", func(t *testing.T) {
		trueVal := NewBool(true)
		falseVal := NewBool(false)

		if trueVal.String() != "true" {
			t.Errorf("Expected 'true', got %s", trueVal.String())
		}
		if falseVal.String() != "false" {
			t.Errorf("Expected 'false', got %s", falseVal.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		val1 := NewBool(true)
		val2 := NewBool(true)
		val3 := NewBool(false)
		val4 := NewString("true")

		if !val1.Equal(val2) {
			t.Error("Expected equal bool values to be equal")
		}
		if val1.Equal(val3) {
			t.Error("Expected different bool values to be unequal")
		}
		if val1.Equal(val4) {
			t.Error("Expected bool and string to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		trueVal := NewBool(true)
		falseVal := NewBool(false)

		if trueVal.Hash() != 1 {
			t.Errorf("Expected hash 1 for true, got %d", trueVal.Hash())
		}
		if falseVal.Hash() != 0 {
			t.Errorf("Expected hash 0 for false, got %d", falseVal.Hash())
		}
	})
}

// TestIntValue tests IntValue implementation
func TestIntValue(t *testing.T) {
	t.Run("NewInt", func(t *testing.T) {
		val := NewInt(42)
		if val.Value() != 42 {
			t.Errorf("Expected 42, got %v", val.Value())
		}
	})

	t.Run("Type", func(t *testing.T) {
		val := NewInt(42)
		typeInfo := val.Type()

		if typeInfo.Kind != KindInt64 {
			t.Errorf("Expected KindInt64, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "int" {
			t.Errorf("Expected 'int', got %s", typeInfo.Name)
		}
		if typeInfo.Size != 8 {
			t.Errorf("Expected size 8, got %d", typeInfo.Size)
		}
	})

	t.Run("String", func(t *testing.T) {
		val := NewInt(42)
		if val.String() != "42" {
			t.Errorf("Expected '42', got %s", val.String())
		}

		negVal := NewInt(-123)
		if negVal.String() != "-123" {
			t.Errorf("Expected '-123', got %s", negVal.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		val1 := NewInt(42)
		val2 := NewInt(42)
		val3 := NewInt(24)
		val4 := NewString("42")

		if !val1.Equal(val2) {
			t.Error("Expected equal int values to be equal")
		}
		if val1.Equal(val3) {
			t.Error("Expected different int values to be unequal")
		}
		if val1.Equal(val4) {
			t.Error("Expected int and string to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		val := NewInt(42)
		if val.Hash() != uint64(42) {
			t.Errorf("Expected hash %d, got %d", uint64(42), val.Hash())
		}
	})
}

// TestFloatValue tests FloatValue implementation
func TestFloatValue(t *testing.T) {
	t.Run("NewFloat", func(t *testing.T) {
		val := NewFloat(3.14)
		if val.Value() != 3.14 {
			t.Errorf("Expected 3.14, got %v", val.Value())
		}
	})

	t.Run("Type", func(t *testing.T) {
		val := NewFloat(3.14)
		typeInfo := val.Type()

		if typeInfo.Kind != KindFloat64 {
			t.Errorf("Expected KindFloat64, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "float" {
			t.Errorf("Expected 'float', got %s", typeInfo.Name)
		}
		if typeInfo.Size != 8 {
			t.Errorf("Expected size 8, got %d", typeInfo.Size)
		}
	})

	t.Run("String", func(t *testing.T) {
		val := NewFloat(3.14)
		expected := "3.14"
		if val.String() != expected {
			t.Errorf("Expected '%s', got %s", expected, val.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		val1 := NewFloat(3.14)
		val2 := NewFloat(3.14)
		val3 := NewFloat(2.71)
		val4 := NewString("3.14")

		if !val1.Equal(val2) {
			t.Error("Expected equal float values to be equal")
		}
		if val1.Equal(val3) {
			t.Error("Expected different float values to be unequal")
		}
		if val1.Equal(val4) {
			t.Error("Expected float and string to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		val1 := NewFloat(3.14)
		val2 := NewFloat(3.14)

		// Same values should have same hash
		if val1.Hash() != val2.Hash() {
			t.Error("Expected same float values to have same hash")
		}

		val3 := NewFloat(2.71)
		// Different values should have different hash (usually)
		if val1.Hash() == val3.Hash() {
			t.Log("Warning: Different float values have same hash (possible but unlikely)")
		}
	})
}

// TestStringValue tests StringValue implementation
func TestStringValue(t *testing.T) {
	t.Run("NewString", func(t *testing.T) {
		val := NewString("hello")
		if val.Value() != "hello" {
			t.Errorf("Expected 'hello', got %v", val.Value())
		}
	})

	t.Run("Type", func(t *testing.T) {
		val := NewString("hello")
		typeInfo := val.Type()

		if typeInfo.Kind != KindString {
			t.Errorf("Expected KindString, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "string" {
			t.Errorf("Expected 'string', got %s", typeInfo.Name)
		}
		if typeInfo.Size != -1 {
			t.Errorf("Expected size -1, got %d", typeInfo.Size)
		}
	})

	t.Run("String", func(t *testing.T) {
		val := NewString("hello world")
		if val.String() != "hello world" {
			t.Errorf("Expected 'hello world', got %s", val.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		val1 := NewString("hello")
		val2 := NewString("hello")
		val3 := NewString("world")
		val4 := NewInt(42)

		if !val1.Equal(val2) {
			t.Error("Expected equal string values to be equal")
		}
		if val1.Equal(val3) {
			t.Error("Expected different string values to be unequal")
		}
		if val1.Equal(val4) {
			t.Error("Expected string and int to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		val1 := NewString("hello")
		val2 := NewString("hello")
		val3 := NewString("world")

		// Same values should have same hash
		if val1.Hash() != val2.Hash() {
			t.Error("Expected same string values to have same hash")
		}

		// Different values should have different hash (usually)
		if val1.Hash() == val3.Hash() {
			t.Log("Warning: Different string values have same hash (possible but unlikely)")
		}
	})
}

// TestSliceValue tests SliceValue implementation
func TestSliceValue(t *testing.T) {
	t.Run("NewSlice", func(t *testing.T) {
		values := []Value{NewInt(1), NewInt(2), NewInt(3)}
		elemType := TypeInfo{Kind: KindInt64, Name: "int"}
		slice := NewSlice(values, elemType)

		if slice.Len() != 3 {
			t.Errorf("Expected length 3, got %d", slice.Len())
		}

		vals := slice.Values()
		if len(vals) != 3 {
			t.Errorf("Expected 3 values, got %d", len(vals))
		}
	})

	t.Run("Type", func(t *testing.T) {
		values := []Value{NewInt(1), NewInt(2)}
		elemType := TypeInfo{Kind: KindInt64, Name: "int"}
		slice := NewSlice(values, elemType)

		typeInfo := slice.Type()
		if typeInfo.Kind != KindSlice {
			t.Errorf("Expected KindSlice, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "[]int" {
			t.Errorf("Expected '[]int', got %s", typeInfo.Name)
		}
		if typeInfo.Size != -1 {
			t.Errorf("Expected size -1, got %d", typeInfo.Size)
		}
	})

	t.Run("Get", func(t *testing.T) {
		values := []Value{NewInt(1), NewInt(2), NewInt(3)}
		elemType := TypeInfo{Kind: KindInt64, Name: "int"}
		slice := NewSlice(values, elemType)

		val := slice.Get(1)
		if intVal, ok := val.(*IntValue); ok {
			if intVal.Value() != 2 {
				t.Errorf("Expected 2, got %v", intVal.Value())
			}
		} else {
			t.Error("Expected IntValue")
		}

		// Test out of bounds
		nilVal := slice.Get(10)
		if nilVal != nil {
			t.Error("Expected nil for out of bounds access")
		}
	})

	t.Run("Equal", func(t *testing.T) {
		values1 := []Value{NewInt(1), NewInt(2)}
		values2 := []Value{NewInt(1), NewInt(2)}
		values3 := []Value{NewInt(1), NewInt(3)}
		elemType := TypeInfo{Kind: KindInt64, Name: "int"}

		slice1 := NewSlice(values1, elemType)
		slice2 := NewSlice(values2, elemType)
		slice3 := NewSlice(values3, elemType)
		str := NewString("test")

		if !slice1.Equal(slice2) {
			t.Error("Expected equal slices to be equal")
		}
		if slice1.Equal(slice3) {
			t.Error("Expected different slices to be unequal")
		}
		if slice1.Equal(str) {
			t.Error("Expected slice and string to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		values1 := []Value{NewInt(1), NewInt(2)}
		values2 := []Value{NewInt(1), NewInt(2)}
		elemType := TypeInfo{Kind: KindInt64, Name: "int"}

		slice1 := NewSlice(values1, elemType)
		slice2 := NewSlice(values2, elemType)

		// Same slices should have same hash
		if slice1.Hash() != slice2.Hash() {
			t.Error("Expected same slices to have same hash")
		}
	})
}

// TestMapValue tests MapValue implementation
func TestMapValue(t *testing.T) {
	t.Run("NewMap", func(t *testing.T) {
		values := map[string]Value{
			"key1": NewInt(1),
			"key2": NewString("value"),
		}
		keyType := TypeInfo{Kind: KindString, Name: "string"}
		valType := TypeInfo{Kind: KindInterface, Name: "interface"}
		mapVal := NewMap(values, keyType, valType)

		vals := mapVal.Values()
		if len(vals) != 2 {
			t.Errorf("Expected 2 values, got %d", len(vals))
		}
	})

	t.Run("Type", func(t *testing.T) {
		values := map[string]Value{"key": NewInt(1)}
		keyType := TypeInfo{Kind: KindString, Name: "string"}
		valType := TypeInfo{Kind: KindInt64, Name: "int"}
		mapVal := NewMap(values, keyType, valType)

		typeInfo := mapVal.Type()
		if typeInfo.Kind != KindMap {
			t.Errorf("Expected KindMap, got %v", typeInfo.Kind)
		}
		expectedName := "map[string]int"
		if typeInfo.Name != expectedName {
			t.Errorf("Expected '%s', got %s", expectedName, typeInfo.Name)
		}
	})

	t.Run("Get and Has", func(t *testing.T) {
		values := map[string]Value{
			"existing": NewInt(42),
		}
		keyType := TypeInfo{Kind: KindString, Name: "string"}
		valType := TypeInfo{Kind: KindInt64, Name: "int"}
		mapVal := NewMap(values, keyType, valType)

		// Test existing key
		if !mapVal.Has("existing") {
			t.Error("Expected map to have 'existing' key")
		}

		val, exists := mapVal.Get("existing")
		if !exists {
			t.Error("Expected key to exist")
		}
		if intVal, ok := val.(*IntValue); ok {
			if intVal.Value() != 42 {
				t.Errorf("Expected 42, got %v", intVal.Value())
			}
		} else {
			t.Error("Expected IntValue")
		}

		// Test non-existing key
		if mapVal.Has("nonexistent") {
			t.Error("Expected map to not have 'nonexistent' key")
		}

		nilVal, exists := mapVal.Get("nonexistent")
		if exists {
			t.Error("Expected key to not exist")
		}
		if nilVal != nil {
			t.Error("Expected nil for non-existent key")
		}
	})

	t.Run("Equal", func(t *testing.T) {
		values1 := map[string]Value{"key": NewInt(1)}
		values2 := map[string]Value{"key": NewInt(1)}
		values3 := map[string]Value{"key": NewInt(2)}
		keyType := TypeInfo{Kind: KindString, Name: "string"}
		valType := TypeInfo{Kind: KindInt64, Name: "int"}

		map1 := NewMap(values1, keyType, valType)
		map2 := NewMap(values2, keyType, valType)
		map3 := NewMap(values3, keyType, valType)
		str := NewString("test")

		if !map1.Equal(map2) {
			t.Error("Expected equal maps to be equal")
		}
		if map1.Equal(map3) {
			t.Error("Expected different maps to be unequal")
		}
		if map1.Equal(str) {
			t.Error("Expected map and string to be unequal")
		}
	})
}

// TestNilValue tests NilValue implementation
func TestNilValue(t *testing.T) {
	t.Run("NewNil", func(t *testing.T) {
		nilVal := NewNil()
		if nilVal == nil {
			t.Error("Expected non-nil NilValue")
		}
	})

	t.Run("Type", func(t *testing.T) {
		nilVal := NewNil()
		typeInfo := nilVal.Type()

		if typeInfo.Kind != KindNil {
			t.Errorf("Expected KindNil, got %v", typeInfo.Kind)
		}
		if typeInfo.Name != "nil" {
			t.Errorf("Expected 'nil', got %s", typeInfo.Name)
		}
		if typeInfo.Size != 0 {
			t.Errorf("Expected size 0, got %d", typeInfo.Size)
		}
	})

	t.Run("String", func(t *testing.T) {
		nilVal := NewNil()
		if nilVal.String() != "nil" {
			t.Errorf("Expected 'nil', got %s", nilVal.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		nil1 := NewNil()
		nil2 := NewNil()
		str := NewString("nil")

		if !nil1.Equal(nil2) {
			t.Error("Expected nil values to be equal")
		}
		if nil1.Equal(str) {
			t.Error("Expected nil and string to be unequal")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		nilVal := NewNil()
		if nilVal.Hash() != 0 {
			t.Errorf("Expected hash 0 for nil, got %d", nilVal.Hash())
		}
	})
}
