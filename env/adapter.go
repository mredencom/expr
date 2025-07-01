package env

import (
	"fmt"

	"github.com/mredencom/expr/types"
)

// Adapter converts Go types to our type system without using reflection
// This is a zero-reflection implementation that requires explicit type registration
type Adapter struct {
	typeRegistry   map[string]TypeAdapter
	structRegistry map[string]StructAdapter
}

// TypeAdapter handles conversion for a specific Go type
type TypeAdapter interface {
	// Convert converts a Go value to our Value type
	Convert(goValue interface{}) (types.Value, error)
	// TypeInfo returns the type information
	TypeInfo() types.TypeInfo
}

// StructAdapter handles struct type operations
type StructAdapter interface {
	TypeAdapter
	// GetField gets a field value by name
	GetField(obj interface{}, fieldName string) (types.Value, error)
	// HasField checks if a field exists
	HasField(fieldName string) bool
	// ListFields returns all field names
	ListFields() []string
}

// FieldInfo contains information about a struct field
type FieldInfo struct {
	Name     string
	TypeInfo types.TypeInfo
	Getter   func(obj interface{}) (types.Value, error)
}

// New creates a new zero-reflection adapter
func New() *Adapter {
	a := &Adapter{
		typeRegistry:   make(map[string]TypeAdapter),
		structRegistry: make(map[string]StructAdapter),
	}

	// Register built-in types
	a.registerBuiltinTypes()

	return a
}

// registerBuiltinTypes registers basic Go types
func (a *Adapter) registerBuiltinTypes() {
	// Basic types
	a.typeRegistry["bool"] = &BoolAdapter{}
	a.typeRegistry["int"] = &IntAdapter{}
	a.typeRegistry["int64"] = &Int64Adapter{}
	a.typeRegistry["float64"] = &Float64Adapter{}
	a.typeRegistry["string"] = &StringAdapter{}
}

// RegisterType registers a type adapter
func (a *Adapter) RegisterType(typeName string, adapter TypeAdapter) {
	a.typeRegistry[typeName] = adapter
}

// RegisterStruct registers a struct adapter
func (a *Adapter) RegisterStruct(typeName string, adapter StructAdapter) {
	a.structRegistry[typeName] = adapter
	a.typeRegistry[typeName] = adapter
}

// ConvertValue converts a Go value to our Value type
func (a *Adapter) ConvertValue(typeName string, goValue interface{}) (types.Value, error) {
	adapter, exists := a.typeRegistry[typeName]
	if !exists {
		return nil, fmt.Errorf("no adapter registered for type: %s", typeName)
	}

	return adapter.Convert(goValue)
}

// GetTypeInfo returns type information for a registered type
func (a *Adapter) GetTypeInfo(typeName string) (types.TypeInfo, error) {
	adapter, exists := a.typeRegistry[typeName]
	if !exists {
		return types.TypeInfo{}, fmt.Errorf("no adapter registered for type: %s", typeName)
	}

	return adapter.TypeInfo(), nil
}

// GetField gets a field from a struct
func (a *Adapter) GetField(typeName string, obj interface{}, fieldName string) (types.Value, error) {
	structAdapter, exists := a.structRegistry[typeName]
	if !exists {
		return nil, fmt.Errorf("no struct adapter registered for type: %s", typeName)
	}

	return structAdapter.GetField(obj, fieldName)
}

// CreateEnvironment creates an environment from a map of variables
func (a *Adapter) CreateEnvironment(variables map[string]interface{}) (map[string]types.Value, error) {
	env := make(map[string]types.Value)

	for name, value := range variables {
		// For basic types, try to auto-detect
		converted, err := a.autoConvert(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert variable %s: %v", name, err)
		}
		env[name] = converted
	}

	return env, nil
}

// autoConvert attempts to automatically convert basic Go types
func (a *Adapter) autoConvert(value interface{}) (types.Value, error) {
	switch v := value.(type) {
	case bool:
		return types.NewBool(v), nil
	case int:
		return types.NewInt(int64(v)), nil
	case int64:
		return types.NewInt(v), nil
	case float64:
		return types.NewFloat(v), nil
	case string:
		return types.NewString(v), nil
	case []int:
		// Convert []int to []types.Value
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewInt(int64(item))
		}
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8}
		return types.NewSlice(values, elemType), nil
	case []float64:
		// Convert []float64 to []types.Value
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewFloat(item)
		}
		elemType := types.TypeInfo{Kind: types.KindFloat64, Name: "float64", Size: 8}
		return types.NewSlice(values, elemType), nil
	case []string:
		// Convert []string to []types.Value
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewString(item)
		}
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		return types.NewSlice(values, elemType), nil
	case []interface{}:
		// Convert []interface{} to []types.Value
		values := make([]types.Value, len(v))
		for i, item := range v {
			converted, err := a.autoConvert(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert slice element %d: %v", i, err)
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewSlice(values, elemType), nil
	case []map[string]interface{}:
		// Convert []map[string]interface{} to []types.Value
		values := make([]types.Value, len(v))
		for i, item := range v {
			converted, err := a.autoConvert(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert slice element %d: %v", i, err)
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindMap, Name: "map[string]interface{}", Size: -1}
		return types.NewSlice(values, elemType), nil
	case map[string]interface{}:
		// Convert map[string]interface{} to MapValue
		values := make(map[string]types.Value)
		for key, val := range v {
			converted, err := a.autoConvert(val)
			if err != nil {
				return nil, fmt.Errorf("failed to convert map value for key %s: %v", key, err)
			}
			values[key] = converted
		}
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewMap(values, keyType, valueType), nil
	default:
		// Try to handle custom struct types using reflection-like approach
		return a.convertCustomType(value)
	}
}

// convertCustomType handles custom struct types
func (a *Adapter) convertCustomType(value interface{}) (types.Value, error) {
	// Handle slice of custom structs first
	if slice, ok := a.tryConvertSlice(value); ok {
		return slice, nil
	}

	// Handle individual custom struct
	if structValue, ok := a.tryConvertStruct(value); ok {
		return structValue, nil
	}

	return nil, fmt.Errorf("unsupported type: %T", value)
}

// tryConvertSlice attempts to convert slice types
func (a *Adapter) tryConvertSlice(value interface{}) (types.Value, bool) {
	switch v := value.(type) {
	case []User:
		values := make([]types.Value, len(v))
		for i, user := range v {
			converted, err := a.convertUserStruct(user)
			if err != nil {
				return nil, false
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindStruct, Name: "User", Size: -1}
		return types.NewSlice(values, elemType), true
	case []Product:
		values := make([]types.Value, len(v))
		for i, product := range v {
			converted, err := a.convertProductStruct(product)
			if err != nil {
				return nil, false
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindStruct, Name: "Product", Size: -1}
		return types.NewSlice(values, elemType), true
	}
	return nil, false
}

// tryConvertStruct attempts to convert struct types
func (a *Adapter) tryConvertStruct(value interface{}) (types.Value, bool) {
	switch v := value.(type) {
	case User:
		converted, err := a.convertUserStruct(v)
		if err != nil {
			return nil, false
		}
		return converted, true
	case Product:
		converted, err := a.convertProductStruct(v)
		if err != nil {
			return nil, false
		}
		return converted, true
	}
	return nil, false
}

// convertUserStruct converts User struct to MapValue for field access
func (a *Adapter) convertUserStruct(user User) (types.Value, error) {
	fields := make(map[string]types.Value)

	fields["Name"] = types.NewString(user.Name)
	fields["Age"] = types.NewInt(int64(user.Age))
	fields["Email"] = types.NewString(user.Email)
	fields["Active"] = types.NewBool(user.Active)
	fields["Balance"] = types.NewFloat(user.Balance)

	// Convert Tags slice
	if user.Tags != nil {
		tagValues := make([]types.Value, len(user.Tags))
		for i, tag := range user.Tags {
			tagValues[i] = types.NewString(tag)
		}
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		fields["Tags"] = types.NewSlice(tagValues, elemType)
	}

	// Convert Metadata map
	if user.Metadata != nil {
		metadataValues := make(map[string]types.Value)
		for key, val := range user.Metadata {
			converted, err := a.autoConvert(val)
			if err != nil {
				return nil, fmt.Errorf("failed to convert metadata field %s: %v", key, err)
			}
			metadataValues[key] = converted
		}
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		fields["Metadata"] = types.NewMap(metadataValues, keyType, valueType)
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
	return types.NewMap(fields, keyType, valueType), nil
}

// convertProductStruct converts Product struct to MapValue for field access
func (a *Adapter) convertProductStruct(product Product) (types.Value, error) {
	fields := make(map[string]types.Value)

	fields["ID"] = types.NewInt(int64(product.ID))
	fields["Name"] = types.NewString(product.Name)
	fields["Price"] = types.NewFloat(product.Price)
	fields["Category"] = types.NewString(product.Category)
	fields["InStock"] = types.NewBool(product.InStock)

	// Convert Tags slice
	if product.Tags != nil {
		tagValues := make([]types.Value, len(product.Tags))
		for i, tag := range product.Tags {
			tagValues[i] = types.NewString(tag)
		}
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		fields["Tags"] = types.NewSlice(tagValues, elemType)
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
	return types.NewMap(fields, keyType, valueType), nil
}

// User struct definition for type compatibility
type User struct {
	Name     string
	Age      int
	Email    string
	Active   bool
	Balance  float64
	Tags     []string
	Metadata map[string]interface{}
}

// Product struct definition for type compatibility
type Product struct {
	ID       int
	Name     string
	Price    float64
	Category string
	InStock  bool
	Tags     []string
}

// Built-in type adapters

// BoolAdapter handles bool type
type BoolAdapter struct{}

func (a *BoolAdapter) Convert(goValue interface{}) (types.Value, error) {
	if v, ok := goValue.(bool); ok {
		return types.NewBool(v), nil
	}
	return nil, fmt.Errorf("expected bool, got %T", goValue)
}

func (a *BoolAdapter) TypeInfo() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindBool, Name: "bool", Size: 1}
}

// IntAdapter handles int type
type IntAdapter struct{}

func (a *IntAdapter) Convert(goValue interface{}) (types.Value, error) {
	if v, ok := goValue.(int); ok {
		return types.NewInt(int64(v)), nil
	}
	return nil, fmt.Errorf("expected int, got %T", goValue)
}

func (a *IntAdapter) TypeInfo() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8}
}

// Int64Adapter handles int64 type
type Int64Adapter struct{}

func (a *Int64Adapter) Convert(goValue interface{}) (types.Value, error) {
	if v, ok := goValue.(int64); ok {
		return types.NewInt(v), nil
	}
	return nil, fmt.Errorf("expected int64, got %T", goValue)
}

func (a *Int64Adapter) TypeInfo() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindInt64, Name: "int64", Size: 8}
}

// Float64Adapter handles float64 type
type Float64Adapter struct{}

func (a *Float64Adapter) Convert(goValue interface{}) (types.Value, error) {
	if v, ok := goValue.(float64); ok {
		return types.NewFloat(v), nil
	}
	return nil, fmt.Errorf("expected float64, got %T", goValue)
}

func (a *Float64Adapter) TypeInfo() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindFloat64, Name: "float64", Size: 8}
}

// StringAdapter handles string type
type StringAdapter struct{}

func (a *StringAdapter) Convert(goValue interface{}) (types.Value, error) {
	if v, ok := goValue.(string); ok {
		return types.NewString(v), nil
	}
	return nil, fmt.Errorf("expected string, got %T", goValue)
}

func (a *StringAdapter) TypeInfo() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
}

// Example struct adapter implementation
// Users can implement similar adapters for their custom types

// UserStructAdapter is an example of how to implement a struct adapter
type UserStructAdapter struct {
	fields map[string]FieldInfo
}

// NewUserStructAdapter creates a new user struct adapter
func NewUserStructAdapter() *UserStructAdapter {
	return &UserStructAdapter{
		fields: map[string]FieldInfo{
			"Name": {
				Name:     "Name",
				TypeInfo: types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1},
				Getter: func(obj interface{}) (types.Value, error) {
					if user, ok := obj.(User); ok {
						return types.NewString(user.Name), nil
					}
					return nil, fmt.Errorf("expected User, got %T", obj)
				},
			},
			"Age": {
				Name:     "Age",
				TypeInfo: types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8},
				Getter: func(obj interface{}) (types.Value, error) {
					if user, ok := obj.(User); ok {
						return types.NewInt(int64(user.Age)), nil
					}
					return nil, fmt.Errorf("expected User, got %T", obj)
				},
			},
			"Active": {
				Name:     "Active",
				TypeInfo: types.TypeInfo{Kind: types.KindBool, Name: "bool", Size: 1},
				Getter: func(obj interface{}) (types.Value, error) {
					if user, ok := obj.(User); ok {
						return types.NewBool(user.Active), nil
					}
					return nil, fmt.Errorf("expected User, got %T", obj)
				},
			},
		},
	}
}

func (a *UserStructAdapter) Convert(goValue interface{}) (types.Value, error) {
	// For structs, we might return a special struct value
	// For now, just return nil as structs are accessed via fields
	return nil, fmt.Errorf("struct conversion not implemented")
}

func (a *UserStructAdapter) TypeInfo() types.TypeInfo {
	fieldInfos := make([]types.FieldInfo, 0, len(a.fields))
	for _, field := range a.fields {
		fieldInfos = append(fieldInfos, types.FieldInfo{
			Name: field.Name,
			Type: field.TypeInfo,
		})
	}

	return types.TypeInfo{
		Kind:   types.KindStruct,
		Name:   "User",
		Size:   -1,
		Fields: fieldInfos,
	}
}

func (a *UserStructAdapter) GetField(obj interface{}, fieldName string) (types.Value, error) {
	field, exists := a.fields[fieldName]
	if !exists {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return field.Getter(obj)
}

func (a *UserStructAdapter) HasField(fieldName string) bool {
	_, exists := a.fields[fieldName]
	return exists
}

func (a *UserStructAdapter) ListFields() []string {
	fields := make([]string, 0, len(a.fields))
	for name := range a.fields {
		fields = append(fields, name)
	}
	return fields
}
