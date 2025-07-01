package types

import (
	"testing"
)

// TestConvertValue tests the main ConvertValue function
func TestConvertValue(t *testing.T) {
	t.Run("SameType", func(t *testing.T) {
		val := NewInt(42)
		result, err := ConvertValue(val, TypeInfo{Kind: KindInt64, Name: "int"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result != val {
			t.Error("Expected same value for same type conversion")
		}
	})

	t.Run("ToString", func(t *testing.T) {
		val := NewInt(42)
		result, err := ConvertValue(val, TypeInfo{Kind: KindString, Name: "string"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		strVal, ok := result.(*StringValue)
		if !ok {
			t.Fatal("Expected StringValue")
		}
		if strVal.Value() != "42" {
			t.Errorf("Expected '42', got %s", strVal.Value())
		}
	})

	t.Run("FromString", func(t *testing.T) {
		val := NewString("123")
		result, err := ConvertValue(val, TypeInfo{Kind: KindInt64, Name: "int"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatal("Expected IntValue")
		}
		if intVal.Value() != 123 {
			t.Errorf("Expected 123, got %d", intVal.Value())
		}
	})

	t.Run("NumericConversion", func(t *testing.T) {
		val := NewInt(42)
		result, err := ConvertValue(val, TypeInfo{Kind: KindFloat64, Name: "float"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		floatVal, ok := result.(*FloatValue)
		if !ok {
			t.Fatal("Expected FloatValue")
		}
		if floatVal.Value() != 42.0 {
			t.Errorf("Expected 42.0, got %f", floatVal.Value())
		}
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		val := NewString("hello")
		_, err := ConvertValue(val, TypeInfo{Kind: KindSlice, Name: "slice"})
		if err == nil {
			t.Error("Expected error for invalid conversion")
		}
	})
}

// TestConvertFromString tests string to other type conversions
func TestConvertFromString(t *testing.T) {
	t.Run("ToBool", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
			hasError bool
		}{
			{"true", true, false},
			{"false", false, false},
			{"1", true, false},
			{"0", false, false},
			{"invalid", false, true},
		}

		for _, test := range tests {
			result, err := convertFromString(test.input, TypeInfo{Kind: KindBool, Name: "bool"})
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s", test.input)
				}
				continue
			}
			if err != nil {
				t.Fatalf("Unexpected error for input %s: %v", test.input, err)
			}
			boolVal, ok := result.(*BoolValue)
			if !ok {
				t.Fatal("Expected BoolValue")
			}
			if boolVal.Value() != test.expected {
				t.Errorf("For input %s, expected %v, got %v", test.input, test.expected, boolVal.Value())
			}
		}
	})

	t.Run("ToInt", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
			hasError bool
		}{
			{"123", 123, false},
			{"-456", -456, false},
			{"0", 0, false},
			{"abc", 0, true},
			{"123.45", 0, true},
		}

		for _, test := range tests {
			result, err := convertFromString(test.input, TypeInfo{Kind: KindInt64, Name: "int"})
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s", test.input)
				}
				continue
			}
			if err != nil {
				t.Fatalf("Unexpected error for input %s: %v", test.input, err)
			}
			intVal, ok := result.(*IntValue)
			if !ok {
				t.Fatal("Expected IntValue")
			}
			if intVal.Value() != test.expected {
				t.Errorf("For input %s, expected %d, got %d", test.input, test.expected, intVal.Value())
			}
		}
	})

	t.Run("ToFloat", func(t *testing.T) {
		tests := []struct {
			input    string
			expected float64
			hasError bool
		}{
			{"123.45", 123.45, false},
			{"-456.78", -456.78, false},
			{"0", 0.0, false},
			{"123", 123.0, false},
			{"abc", 0.0, true},
		}

		for _, test := range tests {
			result, err := convertFromString(test.input, TypeInfo{Kind: KindFloat64, Name: "float"})
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s", test.input)
				}
				continue
			}
			if err != nil {
				t.Fatalf("Unexpected error for input %s: %v", test.input, err)
			}
			floatVal, ok := result.(*FloatValue)
			if !ok {
				t.Fatal("Expected FloatValue")
			}
			if floatVal.Value() != test.expected {
				t.Errorf("For input %s, expected %f, got %f", test.input, test.expected, floatVal.Value())
			}
		}
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		_, err := convertFromString("test", TypeInfo{Kind: KindSlice, Name: "slice"})
		if err == nil {
			t.Error("Expected error for unsupported type conversion")
		}
	})
}

// TestConvertNumeric tests numeric type conversions
func TestConvertNumeric(t *testing.T) {
	t.Run("IntToFloat", func(t *testing.T) {
		val := NewInt(42)
		result, err := convertNumeric(val, TypeInfo{Kind: KindFloat64, Name: "float"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		floatVal, ok := result.(*FloatValue)
		if !ok {
			t.Fatal("Expected FloatValue")
		}
		if floatVal.Value() != 42.0 {
			t.Errorf("Expected 42.0, got %f", floatVal.Value())
		}
	})

	t.Run("FloatToInt", func(t *testing.T) {
		val := NewFloat(42.7)
		result, err := convertNumeric(val, TypeInfo{Kind: KindInt64, Name: "int"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatal("Expected IntValue")
		}
		if intVal.Value() != 42 {
			t.Errorf("Expected 42, got %d", intVal.Value())
		}
	})

	t.Run("NonNumericValue", func(t *testing.T) {
		val := NewString("not a number")
		_, err := convertNumeric(val, TypeInfo{Kind: KindInt64, Name: "int"})
		if err == nil {
			t.Error("Expected error for non-numeric value")
		}
	})
}

// TestConvertFromInt tests integer conversion functions
func TestConvertFromInt(t *testing.T) {
	t.Run("ToInt", func(t *testing.T) {
		result, err := convertFromInt(42, TypeInfo{Kind: KindInt64, Name: "int"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatal("Expected IntValue")
		}
		if intVal.Value() != 42 {
			t.Errorf("Expected 42, got %d", intVal.Value())
		}
	})

	t.Run("ToFloat", func(t *testing.T) {
		result, err := convertFromInt(42, TypeInfo{Kind: KindFloat64, Name: "float"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		floatVal, ok := result.(*FloatValue)
		if !ok {
			t.Fatal("Expected FloatValue")
		}
		if floatVal.Value() != 42.0 {
			t.Errorf("Expected 42.0, got %f", floatVal.Value())
		}
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		_, err := convertFromInt(42, TypeInfo{Kind: KindString, Name: "string"})
		if err == nil {
			t.Error("Expected error for unsupported type conversion")
		}
	})
}

// TestConvertFromFloat tests float conversion functions
func TestConvertFromFloat(t *testing.T) {
	t.Run("ToFloat", func(t *testing.T) {
		result, err := convertFromFloat(42.5, TypeInfo{Kind: KindFloat64, Name: "float"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		floatVal, ok := result.(*FloatValue)
		if !ok {
			t.Fatal("Expected FloatValue")
		}
		if floatVal.Value() != 42.5 {
			t.Errorf("Expected 42.5, got %f", floatVal.Value())
		}
	})

	t.Run("ToInt", func(t *testing.T) {
		result, err := convertFromFloat(42.7, TypeInfo{Kind: KindInt64, Name: "int"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatal("Expected IntValue")
		}
		if intVal.Value() != 42 {
			t.Errorf("Expected 42, got %d", intVal.Value())
		}
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		_, err := convertFromFloat(42.5, TypeInfo{Kind: KindString, Name: "string"})
		if err == nil {
			t.Error("Expected error for unsupported type conversion")
		}
	})
}

// TestCanConvert tests the CanConvert function
func TestCanConvert(t *testing.T) {
	t.Run("SameType", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		if !CanConvert(intType, intType) {
			t.Error("Expected same types to be convertible")
		}
	})

	t.Run("ToString", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}
		if !CanConvert(intType, stringType) {
			t.Error("Expected any type to be convertible to string")
		}
	})

	t.Run("FromStringToBasicTypes", func(t *testing.T) {
		stringType := TypeInfo{Kind: KindString, Name: "string"}
		basicTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindInt64, Name: "int"},
			{Kind: KindFloat64, Name: "float"},
		}

		for _, targetType := range basicTypes {
			if !CanConvert(stringType, targetType) {
				t.Errorf("Expected string to be convertible to %s", targetType.Name)
			}
		}
	})

	t.Run("NumericConversions", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		floatType := TypeInfo{Kind: KindFloat64, Name: "float"}

		if !CanConvert(intType, floatType) {
			t.Error("Expected int to be convertible to float")
		}
		if !CanConvert(floatType, intType) {
			t.Error("Expected float to be convertible to int")
		}
	})

	t.Run("UnsupportedConversions", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		sliceType := TypeInfo{Kind: KindSlice, Name: "slice"}

		if CanConvert(intType, sliceType) {
			t.Error("Expected int to slice conversion to be unsupported")
		}
	})
}

// TestConvertToGo tests conversion from Value to native Go types
func TestConvertToGo(t *testing.T) {
	t.Run("BoolValue", func(t *testing.T) {
		val := NewBool(true)
		result := ConvertToGo(val)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("IntValue", func(t *testing.T) {
		val := NewInt(42)
		result := ConvertToGo(val)
		if result != int64(42) {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("FloatValue", func(t *testing.T) {
		val := NewFloat(3.14)
		result := ConvertToGo(val)
		if result != 3.14 {
			t.Errorf("Expected 3.14, got %v", result)
		}
	})

	t.Run("StringValue", func(t *testing.T) {
		val := NewString("hello")
		result := ConvertToGo(val)
		if result != "hello" {
			t.Errorf("Expected 'hello', got %v", result)
		}
	})

	t.Run("SliceValue", func(t *testing.T) {
		values := []Value{NewInt(1), NewInt(2), NewInt(3)}
		val := NewSlice(values, IntType)
		result := ConvertToGo(val)

		slice, ok := result.([]interface{})
		if !ok {
			t.Fatal("Expected []interface{}")
		}
		if len(slice) != 3 {
			t.Errorf("Expected length 3, got %d", len(slice))
		}
		if slice[0] != int64(1) {
			t.Errorf("Expected 1, got %v", slice[0])
		}
	})

	t.Run("MapValue", func(t *testing.T) {
		values := map[string]Value{
			"a": NewInt(1),
			"b": NewInt(2),
		}
		val := NewMap(values, StringType, IntType)
		result := ConvertToGo(val)

		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatal("Expected map[string]interface{}")
		}
		if len(m) != 2 {
			t.Errorf("Expected length 2, got %d", len(m))
		}
		if m["a"] != int64(1) {
			t.Errorf("Expected 1, got %v", m["a"])
		}
	})

	t.Run("NilValue", func(t *testing.T) {
		val := NewNil()
		result := ConvertToGo(val)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("UnknownValue", func(t *testing.T) {
		// This tests the default case
		var unknownVal Value = nil
		result := ConvertToGo(unknownVal)
		if result != nil {
			t.Errorf("Expected nil for unknown value, got %v", result)
		}
	})
}

// TestConvertFromGo tests conversion from native Go types to Value
func TestConvertFromGo(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		result := ConvertFromGo(true)
		boolVal, ok := result.(*BoolValue)
		if !ok {
			t.Fatal("Expected BoolValue")
		}
		if boolVal.Value() != true {
			t.Error("Expected true")
		}
	})

	t.Run("IntTypes", func(t *testing.T) {
		tests := []interface{}{
			int(42), int8(42), int16(42), int32(42), int64(42),
			uint(42), uint8(42), uint16(42), uint32(42), uint64(42),
		}

		for _, test := range tests {
			result := ConvertFromGo(test)
			intVal, ok := result.(*IntValue)
			if !ok {
				t.Fatalf("Expected IntValue for %T", test)
			}
			if intVal.Value() != 42 {
				t.Errorf("Expected 42, got %d for type %T", intVal.Value(), test)
			}
		}
	})

	t.Run("FloatTypes", func(t *testing.T) {
		tests := []interface{}{
			float32(3.14), float64(3.14),
		}

		for _, test := range tests {
			result := ConvertFromGo(test)
			floatVal, ok := result.(*FloatValue)
			if !ok {
				t.Fatalf("Expected FloatValue for %T", test)
			}
			expected := 3.14
			if f32, ok := test.(float32); ok {
				expected = float64(f32)
			}
			if floatVal.Value() != expected {
				t.Errorf("Expected %f, got %f for type %T", expected, floatVal.Value(), test)
			}
		}
	})

	t.Run("String", func(t *testing.T) {
		result := ConvertFromGo("hello")
		stringVal, ok := result.(*StringValue)
		if !ok {
			t.Fatal("Expected StringValue")
		}
		if stringVal.Value() != "hello" {
			t.Error("Expected 'hello'")
		}
	})

	t.Run("Slice", func(t *testing.T) {
		slice := []interface{}{1, 2, 3}
		result := ConvertFromGo(slice)
		sliceVal, ok := result.(*SliceValue)
		if !ok {
			t.Fatal("Expected SliceValue")
		}
		if sliceVal.Len() != 3 {
			t.Errorf("Expected length 3, got %d", sliceVal.Len())
		}
	})

	t.Run("Map", func(t *testing.T) {
		m := map[string]interface{}{
			"a": 1,
			"b": 2,
		}
		result := ConvertFromGo(m)
		mapVal, ok := result.(*MapValue)
		if !ok {
			t.Fatal("Expected MapValue")
		}
		if mapVal.Len() != 2 {
			t.Errorf("Expected length 2, got %d", mapVal.Len())
		}
	})

	t.Run("Nil", func(t *testing.T) {
		result := ConvertFromGo(nil)
		_, ok := result.(*NilValue)
		if !ok {
			t.Fatal("Expected NilValue")
		}
	})

	t.Run("UnknownType", func(t *testing.T) {
		type CustomType struct{}
		custom := CustomType{}
		result := ConvertFromGo(custom)
		stringVal, ok := result.(*StringValue)
		if !ok {
			t.Fatal("Expected StringValue for unknown type")
		}
		if stringVal.Value() != "{}" {
			t.Errorf("Expected '{}', got %s", stringVal.Value())
		}
	})
}
