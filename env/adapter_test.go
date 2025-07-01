package env

import (
	"testing"

	"github.com/mredencom/expr/types"
)

func TestNewAdapter(t *testing.T) {
	adapter := New()
	if adapter == nil {
		t.Fatal("Expected non-nil adapter")
	}
	if adapter.typeRegistry == nil {
		t.Fatal("Expected initialized type registry")
	}
	if adapter.structRegistry == nil {
		t.Fatal("Expected initialized struct registry")
	}
}

func TestBoolAdapter(t *testing.T) {
	adapter := &BoolAdapter{}

	// Test successful conversion
	result, err := adapter.Convert(true)
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

	// Test type error
	_, err = adapter.Convert("not a bool")
	if err == nil {
		t.Error("Expected error for non-bool input")
	}

	// Test TypeInfo
	typeInfo := adapter.TypeInfo()
	if typeInfo.Kind != types.KindBool {
		t.Errorf("Expected KindBool, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "bool" {
		t.Errorf("Expected name 'bool', got %s", typeInfo.Name)
	}
}

func TestIntAdapter(t *testing.T) {
	adapter := &IntAdapter{}

	// Test successful conversion
	result, err := adapter.Convert(42)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	intVal, ok := result.(*types.IntValue)
	if !ok {
		t.Fatalf("Expected IntValue, got %T", result)
	}

	if intVal.Value() != 42 {
		t.Errorf("Expected 42, got %d", intVal.Value())
	}

	// Test type error
	_, err = adapter.Convert("not an int")
	if err == nil {
		t.Error("Expected error for non-int input")
	}

	// Test TypeInfo
	typeInfo := adapter.TypeInfo()
	if typeInfo.Kind != types.KindInt64 {
		t.Errorf("Expected KindInt64, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "int" {
		t.Errorf("Expected name 'int', got %s", typeInfo.Name)
	}
}

func TestFloat64Adapter(t *testing.T) {
	adapter := &Float64Adapter{}

	// Test successful conversion
	result, err := adapter.Convert(3.14)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	floatVal, ok := result.(*types.FloatValue)
	if !ok {
		t.Fatalf("Expected FloatValue, got %T", result)
	}

	if floatVal.Value() != 3.14 {
		t.Errorf("Expected 3.14, got %f", floatVal.Value())
	}

	// Test type error
	_, err = adapter.Convert("not a float")
	if err == nil {
		t.Error("Expected error for non-float input")
	}
}

func TestStringAdapter(t *testing.T) {
	adapter := &StringAdapter{}

	// Test successful conversion
	result, err := adapter.Convert("hello")
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

	// Test type error
	_, err = adapter.Convert(42)
	if err == nil {
		t.Error("Expected error for non-string input")
	}
}

func TestRegisterType(t *testing.T) {
	adapter := New()

	// Create a custom adapter
	customAdapter := &StringAdapter{}

	// Register it
	adapter.RegisterType("custom", customAdapter)

	// Test conversion using registered type
	result, err := adapter.ConvertValue("custom", "test")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	strVal, ok := result.(*types.StringValue)
	if !ok {
		t.Fatalf("Expected StringValue, got %T", result)
	}

	if strVal.Value() != "test" {
		t.Errorf("Expected 'test', got %s", strVal.Value())
	}
}

func TestAutoConvert(t *testing.T) {
	adapter := New()

	tests := []struct {
		name     string
		input    interface{}
		expected types.Value
	}{
		{
			name:     "bool true",
			input:    true,
			expected: types.NewBool(true),
		},
		{
			name:     "int",
			input:    42,
			expected: types.NewInt(42),
		},
		{
			name:     "int64",
			input:    int64(123),
			expected: types.NewInt(123),
		},
		{
			name:     "float64",
			input:    3.14,
			expected: types.NewFloat(3.14),
		},
		{
			name:     "string",
			input:    "hello",
			expected: types.NewString("hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.autoConvert(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

	// Test unsupported type
	_, err := adapter.autoConvert(struct{}{})
	if err == nil {
		t.Error("Expected error for unsupported type")
	}
}

func TestCreateEnvironment(t *testing.T) {
	adapter := New()

	variables := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
		"score":  95.5,
		"userId": int64(12345),
	}

	env, err := adapter.CreateEnvironment(variables)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(env) != 5 {
		t.Errorf("Expected 5 variables, got %d", len(env))
	}

	// Test individual variables
	nameVal, exists := env["name"]
	if !exists {
		t.Error("Expected name variable to exist")
	}
	if nameVal.(*types.StringValue).Value() != "John" {
		t.Error("Expected name to be 'John'")
	}

	ageVal, exists := env["age"]
	if !exists {
		t.Error("Expected age variable to exist")
	}
	if ageVal.(*types.IntValue).Value() != 30 {
		t.Error("Expected age to be 30")
	}
}

func TestUserStructAdapter(t *testing.T) {
	adapter := NewUserStructAdapter()

	user := User{
		Name:   "Alice",
		Age:    25,
		Active: true,
	}

	// Test that struct conversion is not implemented
	_, err := adapter.Convert(&user)
	if err == nil {
		t.Error("Expected error for struct conversion")
	}

	// Test field access instead
	nameVal, err := adapter.GetField(user, "Name")
	if err != nil {
		t.Fatalf("Unexpected error getting Name field: %v", err)
	}
	if nameVal.(*types.StringValue).Value() != "Alice" {
		t.Error("Expected name to be 'Alice'")
	}

	ageVal, err := adapter.GetField(user, "Age")
	if err != nil {
		t.Fatalf("Unexpected error getting Age field: %v", err)
	}
	if ageVal.(*types.IntValue).Value() != 25 {
		t.Error("Expected age to be 25")
	}

	activeVal, err := adapter.GetField(user, "Active")
	if err != nil {
		t.Fatalf("Unexpected error getting Active field: %v", err)
	}
	if !activeVal.(*types.BoolValue).Value() {
		t.Error("Expected active to be true")
	}
}

func TestUserStructAdapterGetField(t *testing.T) {
	adapter := NewUserStructAdapter()

	user := User{
		Name:   "Bob",
		Age:    30,
		Active: false,
	}

	// Test getting existing field
	nameVal, err := adapter.GetField(user, "Name")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if nameVal.(*types.StringValue).Value() != "Bob" {
		t.Error("Expected name to be 'Bob'")
	}

	// Test getting non-existent field
	_, err = adapter.GetField(user, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent field")
	}

	// Test with wrong type
	_, err = adapter.GetField("not a user", "Name")
	if err == nil {
		t.Error("Expected error for wrong type")
	}
}

func TestUserStructAdapterHasField(t *testing.T) {
	adapter := NewUserStructAdapter()

	if !adapter.HasField("Name") {
		t.Error("Expected Name field to exist")
	}

	if !adapter.HasField("Age") {
		t.Error("Expected Age field to exist")
	}

	if !adapter.HasField("Active") {
		t.Error("Expected Active field to exist")
	}

	if adapter.HasField("nonexistent") {
		t.Error("Expected nonexistent field to not exist")
	}
}

func TestUserStructAdapterListFields(t *testing.T) {
	adapter := NewUserStructAdapter()

	fields := adapter.ListFields()

	expectedFields := []string{"Name", "Age", "Active"}
	if len(fields) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(fields))
	}

	// Check all expected fields are present
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	for _, expectedField := range expectedFields {
		if !fieldMap[expectedField] {
			t.Errorf("Expected field %s to be present", expectedField)
		}
	}
}

func TestUserStructAdapterTypeInfo(t *testing.T) {
	adapter := NewUserStructAdapter()

	typeInfo := adapter.TypeInfo()

	if typeInfo.Kind != types.KindStruct {
		t.Errorf("Expected KindStruct, got %v", typeInfo.Kind)
	}

	if typeInfo.Name != "User" {
		t.Errorf("Expected name 'User', got %s", typeInfo.Name)
	}

	if len(typeInfo.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(typeInfo.Fields))
	}
}
