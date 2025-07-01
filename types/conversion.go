package types

import (
	"fmt"
	"strconv"
)

// ConvertValue converts a value from one type to another
func ConvertValue(value Value, targetType TypeInfo) (Value, error) {
	sourceType := value.Type()

	// Same type, no conversion needed
	if sourceType.Kind == targetType.Kind {
		return value, nil
	}

	// Convert to string
	if targetType.Kind == KindString {
		return NewString(value.String()), nil
	}

	// Convert from string
	if sourceType.Kind == KindString {
		str := value.(*StringValue).Value()
		return convertFromString(str, targetType)
	}

	// Numeric conversions
	if sourceType.IsNumeric() && targetType.IsNumeric() {
		return convertNumeric(value, targetType)
	}

	return nil, fmt.Errorf("cannot convert %s to %s", sourceType.Name, targetType.Name)
}

// convertFromString converts a string to the target type
func convertFromString(str string, targetType TypeInfo) (Value, error) {
	switch targetType.Kind {
	case KindBool:
		v, err := strconv.ParseBool(str)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to bool: %v", str, err)
		}
		return NewBool(v), nil

	case KindInt, KindInt64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to int: %v", str, err)
		}
		return NewInt(v), nil

	case KindFloat64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to float: %v", str, err)
		}
		return NewFloat(v), nil

	default:
		return nil, fmt.Errorf("cannot convert string to %s", targetType.Name)
	}
}

// convertNumeric converts between numeric types
func convertNumeric(value Value, targetType TypeInfo) (Value, error) {
	switch v := value.(type) {
	case *IntValue:
		return convertFromInt(v.Value(), targetType)
	case *FloatValue:
		return convertFromFloat(v.Value(), targetType)
	default:
		return nil, fmt.Errorf("not a numeric value: %T", value)
	}
}

// convertFromInt converts an int64 to the target numeric type
func convertFromInt(val int64, targetType TypeInfo) (Value, error) {
	switch targetType.Kind {
	case KindInt, KindInt64:
		return NewInt(val), nil
	case KindFloat64:
		return NewFloat(float64(val)), nil
	default:
		return nil, fmt.Errorf("cannot convert int to %s", targetType.Name)
	}
}

// convertFromFloat converts a float64 to the target numeric type
func convertFromFloat(val float64, targetType TypeInfo) (Value, error) {
	switch targetType.Kind {
	case KindInt, KindInt64:
		return NewInt(int64(val)), nil
	case KindFloat64:
		return NewFloat(val), nil
	default:
		return nil, fmt.Errorf("cannot convert float to %s", targetType.Name)
	}
}

// CanConvert returns true if a value can be converted to the target type
func CanConvert(sourceType, targetType TypeInfo) bool {
	// Same type
	if sourceType.Kind == targetType.Kind {
		return true
	}

	// To string
	if targetType.Kind == KindString {
		return true
	}

	// From string to basic types
	if sourceType.Kind == KindString {
		switch targetType.Kind {
		case KindBool, KindInt, KindInt64, KindFloat64:
			return true
		}
	}

	// Numeric conversions
	if sourceType.IsNumeric() && targetType.IsNumeric() {
		return true
	}

	return false
}

// ConvertToGo converts a Value to a native Go value
func ConvertToGo(value Value) interface{} {
	switch v := value.(type) {
	case *BoolValue:
		return v.Value()
	case *IntValue:
		return v.Value()
	case *FloatValue:
		return v.Value()
	case *StringValue:
		return v.Value()
	case *SliceValue:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = ConvertToGo(v.Get(i))
		}
		return result
	case *MapValue:
		result := make(map[string]interface{})
		for k, val := range v.Values() {
			result[k] = ConvertToGo(val)
		}
		return result
	case *NilValue:
		return nil
	default:
		return nil
	}
}

// ConvertFromGo converts a native Go value to a Value
func ConvertFromGo(v interface{}) Value {
	switch val := v.(type) {
	case bool:
		return NewBool(val)
	case int:
		return NewInt(int64(val))
	case int8:
		return NewInt(int64(val))
	case int16:
		return NewInt(int64(val))
	case int32:
		return NewInt(int64(val))
	case int64:
		return NewInt(val)
	case uint:
		return NewInt(int64(val))
	case uint8:
		return NewInt(int64(val))
	case uint16:
		return NewInt(int64(val))
	case uint32:
		return NewInt(int64(val))
	case uint64:
		return NewInt(int64(val))
	case float32:
		return NewFloat(float64(val))
	case float64:
		return NewFloat(val)
	case string:
		return NewString(val)
	case []interface{}:
		values := make([]Value, len(val))
		for i, item := range val {
			values[i] = ConvertFromGo(item)
		}
		return NewSlice(values, TypeInfo{Kind: KindInterface, Name: "interface{}"})
	case map[string]interface{}:
		values := make(map[string]Value)
		for k, item := range val {
			values[k] = ConvertFromGo(item)
		}
		return NewMap(values, StringType, TypeInfo{Kind: KindInterface, Name: "interface{}"})
	case nil:
		return NewNil()
	default:
		// For unknown types, convert to string
		return NewString(fmt.Sprintf("%v", val))
	}
}
