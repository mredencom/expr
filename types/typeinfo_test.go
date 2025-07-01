package types

import (
	"testing"
)

// TestTypeInfo tests TypeInfo functionality
func TestTypeInfo(t *testing.T) {
	t.Run("TypeKind constants", func(t *testing.T) {
		// Test that all type kinds are defined
		kinds := []TypeKind{
			KindBool, KindInt, KindInt8, KindInt16, KindInt32, KindInt64,
			KindUint, KindUint8, KindUint16, KindUint32, KindUint64,
			KindFloat32, KindFloat64, KindString, KindArray, KindSlice,
			KindMap, KindStruct, KindInterface, KindPointer, KindFunc, KindNil,
		}

		for i, kind := range kinds {
			if int(kind) != i {
				t.Errorf("Expected kind %d to have value %d, got %d", i, i, int(kind))
			}
		}
	})

	t.Run("Assignable", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}
		interfaceType := TypeInfo{Kind: KindInterface, Name: "interface"}

		// Same types should be assignable
		if !intType.Assignable(intType) {
			t.Error("Expected same types to be assignable")
		}

		// Interface type should accept all types
		if !interfaceType.Assignable(intType) {
			t.Error("Expected interface to accept int")
		}
		if !interfaceType.Assignable(stringType) {
			t.Error("Expected interface to accept string")
		}

		// String should not be assignable to numeric
		if intType.Assignable(stringType) {
			t.Error("Expected int to not accept string")
		}
	})

	t.Run("Compatible", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		floatType := TypeInfo{Kind: KindFloat64, Name: "float"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}

		// Same types should be compatible
		if !intType.Compatible(intType) {
			t.Error("Expected same types to be compatible")
		}

		// Numeric types should be compatible
		if !intType.Compatible(floatType) {
			t.Error("Expected int and float to be compatible")
		}

		// String should not be compatible with numeric
		if stringType.Compatible(intType) {
			t.Error("Expected string and int to be incompatible")
		}
	})

	t.Run("IsNumeric", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		floatType := TypeInfo{Kind: KindFloat64, Name: "float"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}
		boolType := TypeInfo{Kind: KindBool, Name: "bool"}

		if !intType.IsNumeric() {
			t.Error("Expected int to be numeric")
		}
		if !floatType.IsNumeric() {
			t.Error("Expected float to be numeric")
		}
		if stringType.IsNumeric() {
			t.Error("Expected string to not be numeric")
		}
		if boolType.IsNumeric() {
			t.Error("Expected bool to not be numeric")
		}
	})

	t.Run("IsInteger", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		int32Type := TypeInfo{Kind: KindInt, Name: "int32"}
		floatType := TypeInfo{Kind: KindFloat64, Name: "float"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}

		if !intType.IsInteger() {
			t.Error("Expected int64 to be integer")
		}
		if !int32Type.IsInteger() {
			t.Error("Expected int32 to be integer")
		}
		if floatType.IsInteger() {
			t.Error("Expected float to not be integer")
		}
		if stringType.IsInteger() {
			t.Error("Expected string to not be integer")
		}
	})

	t.Run("IsFloat", func(t *testing.T) {
		floatType := TypeInfo{Kind: KindFloat64, Name: "float64"}
		float32Type := TypeInfo{Kind: KindFloat32, Name: "float32"}
		intType := TypeInfo{Kind: KindInt64, Name: "int"}
		stringType := TypeInfo{Kind: KindString, Name: "string"}

		if !floatType.IsFloat() {
			t.Error("Expected float64 to be float")
		}
		if !float32Type.IsFloat() {
			t.Error("Expected float32 to be float")
		}
		if intType.IsFloat() {
			t.Error("Expected int to not be float")
		}
		if stringType.IsFloat() {
			t.Error("Expected string to not be float")
		}
	})

	t.Run("String", func(t *testing.T) {
		intType := TypeInfo{Kind: KindInt64, Name: "int", Size: 8}
		expected := "int"

		if intType.String() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, intType.String())
		}

		// Test with empty name
		noNameType := TypeInfo{Kind: KindInt64, Size: 8}
		expectedKind := "int64"

		if noNameType.String() != expectedKind {
			t.Errorf("Expected '%s', got '%s'", expectedKind, noNameType.String())
		}
	})
}

// TestTypeKindString tests the TypeKind.String method for all enum values
func TestTypeKindString(t *testing.T) {
	testCases := []struct {
		kind     TypeKind
		expected string
	}{
		{KindBool, "bool"},
		{KindInt, "int"},
		{KindInt8, "int8"},
		{KindInt16, "int16"},
		{KindInt32, "int32"},
		{KindInt64, "int64"},
		{KindUint, "uint"},
		{KindUint8, "uint8"},
		{KindUint16, "uint16"},
		{KindUint32, "uint32"},
		{KindUint64, "uint64"},
		{KindFloat32, "float32"},
		{KindFloat64, "float64"},
		{KindString, "string"},
		{KindArray, "array"},
		{KindSlice, "slice"},
		{KindMap, "map"},
		{KindStruct, "struct"},
		{KindInterface, "interface"},
		{KindPointer, "pointer"},
		{KindFunc, "func"},
		{KindNil, "nil"},
	}

	for _, test := range testCases {
		result := test.kind.String()
		if result != test.expected {
			t.Errorf("Expected %s.String() = %s, got %s", test.kind, test.expected, result)
		}
	}

	// Test unknown kind
	unknownKind := TypeKind(255)
	result := unknownKind.String()
	if result != "unknown" {
		t.Errorf("Expected unknown kind to return 'unknown', got %s", result)
	}
}

// TestIsComparable tests the IsComparable method
func TestIsComparable(t *testing.T) {
	t.Run("BasicComparableTypes", func(t *testing.T) {
		comparableTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindString, Name: "string"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindInt8, Name: "int8"},
			{Kind: KindInt16, Name: "int16"},
			{Kind: KindInt32, Name: "int32"},
			{Kind: KindInt64, Name: "int64"},
			{Kind: KindUint, Name: "uint"},
			{Kind: KindUint8, Name: "uint8"},
			{Kind: KindUint16, Name: "uint16"},
			{Kind: KindUint32, Name: "uint32"},
			{Kind: KindUint64, Name: "uint64"},
			{Kind: KindFloat32, Name: "float32"},
			{Kind: KindFloat64, Name: "float64"},
		}

		for _, typeInfo := range comparableTypes {
			if !typeInfo.IsComparable() {
				t.Errorf("Expected %s to be comparable", typeInfo.Name)
			}
		}
	})

	t.Run("NonComparableTypes", func(t *testing.T) {
		nonComparableTypes := []TypeInfo{
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindFunc, Name: "func"},
			{Kind: KindInterface, Name: "interface"},
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindNil, Name: "nil"},
		}

		for _, typeInfo := range nonComparableTypes {
			if typeInfo.IsComparable() {
				t.Errorf("Expected %s to not be comparable", typeInfo.Name)
			}
		}
	})

	t.Run("ArrayComparability", func(t *testing.T) {
		// Array with comparable element type
		comparableArray := TypeInfo{
			Kind:     KindArray,
			Name:     "array",
			ElemType: &TypeInfo{Kind: KindInt, Name: "int"},
		}
		if !comparableArray.IsComparable() {
			t.Error("Expected array with comparable element type to be comparable")
		}

		// Array with non-comparable element type
		nonComparableArray := TypeInfo{
			Kind:     KindArray,
			Name:     "array",
			ElemType: &TypeInfo{Kind: KindSlice, Name: "slice"},
		}
		if nonComparableArray.IsComparable() {
			t.Error("Expected array with non-comparable element type to not be comparable")
		}

		// Array with nil element type
		nilElemArray := TypeInfo{
			Kind:     KindArray,
			Name:     "array",
			ElemType: nil,
		}
		if nilElemArray.IsComparable() {
			t.Error("Expected array with nil element type to not be comparable")
		}
	})

	t.Run("StructComparability", func(t *testing.T) {
		// Struct with all comparable fields
		comparableStruct := TypeInfo{
			Kind: KindStruct,
			Name: "struct",
			Fields: []FieldInfo{
				{Name: "field1", Type: TypeInfo{Kind: KindInt, Name: "int"}},
				{Name: "field2", Type: TypeInfo{Kind: KindString, Name: "string"}},
			},
		}
		if !comparableStruct.IsComparable() {
			t.Error("Expected struct with all comparable fields to be comparable")
		}

		// Struct with non-comparable field
		nonComparableStruct := TypeInfo{
			Kind: KindStruct,
			Name: "struct",
			Fields: []FieldInfo{
				{Name: "field1", Type: TypeInfo{Kind: KindInt, Name: "int"}},
				{Name: "field2", Type: TypeInfo{Kind: KindSlice, Name: "slice"}},
			},
		}
		if nonComparableStruct.IsComparable() {
			t.Error("Expected struct with non-comparable field to not be comparable")
		}

		// Struct with no fields
		emptyStruct := TypeInfo{
			Kind:   KindStruct,
			Name:   "struct",
			Fields: []FieldInfo{},
		}
		if !emptyStruct.IsComparable() {
			t.Error("Expected empty struct to be comparable")
		}
	})
}

// TestIsOrdered tests the IsOrdered method
func TestIsOrdered(t *testing.T) {
	t.Run("OrderedTypes", func(t *testing.T) {
		orderedTypes := []TypeInfo{
			{Kind: KindString, Name: "string"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindInt8, Name: "int8"},
			{Kind: KindInt16, Name: "int16"},
			{Kind: KindInt32, Name: "int32"},
			{Kind: KindInt64, Name: "int64"},
			{Kind: KindUint, Name: "uint"},
			{Kind: KindUint8, Name: "uint8"},
			{Kind: KindUint16, Name: "uint16"},
			{Kind: KindUint32, Name: "uint32"},
			{Kind: KindUint64, Name: "uint64"},
			{Kind: KindFloat32, Name: "float32"},
			{Kind: KindFloat64, Name: "float64"},
		}

		for _, typeInfo := range orderedTypes {
			if !typeInfo.IsOrdered() {
				t.Errorf("Expected %s to be ordered", typeInfo.Name)
			}
		}
	})

	t.Run("NonOrderedTypes", func(t *testing.T) {
		nonOrderedTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindStruct, Name: "struct"},
			{Kind: KindInterface, Name: "interface"},
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindFunc, Name: "func"},
			{Kind: KindNil, Name: "nil"},
		}

		for _, typeInfo := range nonOrderedTypes {
			if typeInfo.IsOrdered() {
				t.Errorf("Expected %s to not be ordered", typeInfo.Name)
			}
		}
	})
}

// TestExtendedAssignable tests more edge cases for the Assignable method
func TestExtendedAssignable(t *testing.T) {
	t.Run("NilAssignability", func(t *testing.T) {
		nilType := TypeInfo{Kind: KindNil, Name: "nil"}

		// Nil should be assignable to pointers, slices, maps, interfaces
		assignableToNil := []TypeInfo{
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindInterface, Name: "interface"},
		}

		for _, targetType := range assignableToNil {
			if !targetType.Assignable(nilType) {
				t.Errorf("Expected %s to accept nil", targetType.Name)
			}
		}

		// Nil should not be assignable to basic types
		nonAssignableToNil := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindFloat64, Name: "float"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindStruct, Name: "struct"},
		}

		for _, targetType := range nonAssignableToNil {
			if targetType.Assignable(nilType) {
				t.Errorf("Expected %s to not accept nil", targetType.Name)
			}
		}
	})

	t.Run("SameTypeAssignability", func(t *testing.T) {
		types := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindFloat64, Name: "float"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindStruct, Name: "struct"},
		}

		for _, typeInfo := range types {
			if !typeInfo.Assignable(typeInfo) {
				t.Errorf("Expected %s to be assignable to itself", typeInfo.Name)
			}
		}
	})

	t.Run("InterfaceAcceptsAll", func(t *testing.T) {
		interfaceType := TypeInfo{Kind: KindInterface, Name: "interface"}

		allTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindFloat64, Name: "float"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindStruct, Name: "struct"},
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindFunc, Name: "func"},
			{Kind: KindNil, Name: "nil"},
		}

		for _, sourceType := range allTypes {
			if !interfaceType.Assignable(sourceType) {
				t.Errorf("Expected interface to accept %s", sourceType.Name)
			}
		}
	})
}

// TestExtendedCompatible tests more edge cases for the Compatible method
func TestExtendedCompatible(t *testing.T) {
	t.Run("NumericTypeCompatibility", func(t *testing.T) {
		numericTypes := []TypeInfo{
			{Kind: KindInt, Name: "int"},
			{Kind: KindInt8, Name: "int8"},
			{Kind: KindInt16, Name: "int16"},
			{Kind: KindInt32, Name: "int32"},
			{Kind: KindInt64, Name: "int64"},
			{Kind: KindUint, Name: "uint"},
			{Kind: KindUint8, Name: "uint8"},
			{Kind: KindUint16, Name: "uint16"},
			{Kind: KindUint32, Name: "uint32"},
			{Kind: KindUint64, Name: "uint64"},
			{Kind: KindFloat32, Name: "float32"},
			{Kind: KindFloat64, Name: "float64"},
		}

		// All numeric types should be compatible with each other
		for i, type1 := range numericTypes {
			for j, type2 := range numericTypes {
				if !type1.Compatible(type2) {
					t.Errorf("Expected %s to be compatible with %s", type1.Name, type2.Name)
				}
				if i != j && !type1.Compatible(type2) {
					t.Errorf("Expected %s to be compatible with %s", type1.Name, type2.Name)
				}
			}
		}
	})

	t.Run("InterfaceCompatibility", func(t *testing.T) {
		interfaceType := TypeInfo{Kind: KindInterface, Name: "interface"}
		otherTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindInt, Name: "int"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
		}

		for _, otherType := range otherTypes {
			if !interfaceType.Compatible(otherType) {
				t.Errorf("Expected interface to be compatible with %s", otherType.Name)
			}
			if !otherType.Compatible(interfaceType) {
				t.Errorf("Expected %s to be compatible with interface", otherType.Name)
			}
		}
	})

	t.Run("IncompatibleTypes", func(t *testing.T) {
		incompatiblePairs := []struct {
			type1, type2 TypeInfo
		}{
			{
				type1: TypeInfo{Kind: KindString, Name: "string"},
				type2: TypeInfo{Kind: KindBool, Name: "bool"},
			},
			{
				type1: TypeInfo{Kind: KindString, Name: "string"},
				type2: TypeInfo{Kind: KindArray, Name: "array"},
			},
			{
				type1: TypeInfo{Kind: KindBool, Name: "bool"},
				type2: TypeInfo{Kind: KindArray, Name: "array"},
			},
			{
				type1: TypeInfo{Kind: KindArray, Name: "array"},
				type2: TypeInfo{Kind: KindSlice, Name: "slice"},
			},
		}

		for _, pair := range incompatiblePairs {
			if pair.type1.Compatible(pair.type2) {
				t.Errorf("Expected %s to be incompatible with %s", pair.type1.Name, pair.type2.Name)
			}
		}
	})
}

// TestIsNumericExtended tests extended numeric type checking
func TestIsNumericExtended(t *testing.T) {
	t.Run("AllNumericTypes", func(t *testing.T) {
		numericTypes := []TypeInfo{
			{Kind: KindInt, Name: "int"},
			{Kind: KindInt8, Name: "int8"},
			{Kind: KindInt16, Name: "int16"},
			{Kind: KindInt32, Name: "int32"},
			{Kind: KindInt64, Name: "int64"},
			{Kind: KindUint, Name: "uint"},
			{Kind: KindUint8, Name: "uint8"},
			{Kind: KindUint16, Name: "uint16"},
			{Kind: KindUint32, Name: "uint32"},
			{Kind: KindUint64, Name: "uint64"},
			{Kind: KindFloat32, Name: "float32"},
			{Kind: KindFloat64, Name: "float64"},
		}

		for _, typeInfo := range numericTypes {
			if !typeInfo.IsNumeric() {
				t.Errorf("Expected %s to be numeric", typeInfo.Name)
			}
		}
	})

	t.Run("AllNonNumericTypes", func(t *testing.T) {
		nonNumericTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindStruct, Name: "struct"},
			{Kind: KindInterface, Name: "interface"},
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindFunc, Name: "func"},
			{Kind: KindNil, Name: "nil"},
		}

		for _, typeInfo := range nonNumericTypes {
			if typeInfo.IsNumeric() {
				t.Errorf("Expected %s to not be numeric", typeInfo.Name)
			}
		}
	})
}

// TestIsIntegerExtended tests extended integer type checking
func TestIsIntegerExtended(t *testing.T) {
	t.Run("AllIntegerTypes", func(t *testing.T) {
		integerTypes := []TypeInfo{
			{Kind: KindInt, Name: "int"},
			{Kind: KindInt8, Name: "int8"},
			{Kind: KindInt16, Name: "int16"},
			{Kind: KindInt32, Name: "int32"},
			{Kind: KindInt64, Name: "int64"},
			{Kind: KindUint, Name: "uint"},
			{Kind: KindUint8, Name: "uint8"},
			{Kind: KindUint16, Name: "uint16"},
			{Kind: KindUint32, Name: "uint32"},
			{Kind: KindUint64, Name: "uint64"},
		}

		for _, typeInfo := range integerTypes {
			if !typeInfo.IsInteger() {
				t.Errorf("Expected %s to be integer", typeInfo.Name)
			}
		}
	})

	t.Run("AllNonIntegerTypes", func(t *testing.T) {
		nonIntegerTypes := []TypeInfo{
			{Kind: KindBool, Name: "bool"},
			{Kind: KindFloat32, Name: "float32"},
			{Kind: KindFloat64, Name: "float64"},
			{Kind: KindString, Name: "string"},
			{Kind: KindArray, Name: "array"},
			{Kind: KindSlice, Name: "slice"},
			{Kind: KindMap, Name: "map"},
			{Kind: KindStruct, Name: "struct"},
			{Kind: KindInterface, Name: "interface"},
			{Kind: KindPointer, Name: "pointer"},
			{Kind: KindFunc, Name: "func"},
			{Kind: KindNil, Name: "nil"},
		}

		for _, typeInfo := range nonIntegerTypes {
			if typeInfo.IsInteger() {
				t.Errorf("Expected %s to not be integer", typeInfo.Name)
			}
		}
	})
}

// TestPredefinedTypes tests the predefined type instances
func TestPredefinedTypes(t *testing.T) {
	t.Run("BoolType", func(t *testing.T) {
		if BoolType.Kind != KindBool {
			t.Errorf("Expected BoolType.Kind = KindBool, got %v", BoolType.Kind)
		}
		if BoolType.Name != "bool" {
			t.Errorf("Expected BoolType.Name = 'bool', got %s", BoolType.Name)
		}
		if BoolType.Size != 1 {
			t.Errorf("Expected BoolType.Size = 1, got %d", BoolType.Size)
		}
	})

	t.Run("IntType", func(t *testing.T) {
		if IntType.Kind != KindInt64 {
			t.Errorf("Expected IntType.Kind = KindInt64, got %v", IntType.Kind)
		}
		if IntType.Name != "int" {
			t.Errorf("Expected IntType.Name = 'int', got %s", IntType.Name)
		}
		if IntType.Size != 8 {
			t.Errorf("Expected IntType.Size = 8, got %d", IntType.Size)
		}
	})

	t.Run("FloatType", func(t *testing.T) {
		if FloatType.Kind != KindFloat64 {
			t.Errorf("Expected FloatType.Kind = KindFloat64, got %v", FloatType.Kind)
		}
		if FloatType.Name != "float" {
			t.Errorf("Expected FloatType.Name = 'float', got %s", FloatType.Name)
		}
		if FloatType.Size != 8 {
			t.Errorf("Expected FloatType.Size = 8, got %d", FloatType.Size)
		}
	})

	t.Run("StringType", func(t *testing.T) {
		if StringType.Kind != KindString {
			t.Errorf("Expected StringType.Kind = KindString, got %v", StringType.Kind)
		}
		if StringType.Name != "string" {
			t.Errorf("Expected StringType.Name = 'string', got %s", StringType.Name)
		}
		if StringType.Size != -1 {
			t.Errorf("Expected StringType.Size = -1, got %d", StringType.Size)
		}
	})

	t.Run("NilType", func(t *testing.T) {
		if NilType.Kind != KindNil {
			t.Errorf("Expected NilType.Kind = KindNil, got %v", NilType.Kind)
		}
		if NilType.Name != "nil" {
			t.Errorf("Expected NilType.Name = 'nil', got %s", NilType.Name)
		}
		if NilType.Size != 0 {
			t.Errorf("Expected NilType.Size = 0, got %d", NilType.Size)
		}
	})
}
