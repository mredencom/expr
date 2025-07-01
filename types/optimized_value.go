package types

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"unsafe"
)

// ValueType represents the type of a value in the union
type ValueType uint8

const (
	TypeBool ValueType = iota
	TypeInt64
	TypeFloat64
	TypeString
	TypeSlice
	TypeMap
	TypeNil
	TypeFunc
	TypePlaceholder
)

// OptimizedValue is a union type that eliminates interface overhead
// This is the P0 optimization suggested in PERFORMANCE_SUMMARY.md
type OptimizedValue struct {
	Type      ValueType
	IntVal    int64
	FloatVal  float64
	StringVal string
	BoolVal   bool

	// For complex types, we use a pointer to avoid copying large data
	ComplexVal unsafe.Pointer
}

// Performance optimized constructors
// These are inlined for maximum performance

// NewOptimizedBool creates an optimized boolean value
func NewOptimizedBool(v bool) OptimizedValue {
	return OptimizedValue{
		Type:    TypeBool,
		BoolVal: v,
	}
}

// NewOptimizedInt creates an optimized integer value
func NewOptimizedInt(v int64) OptimizedValue {
	return OptimizedValue{
		Type:   TypeInt64,
		IntVal: v,
	}
}

// NewOptimizedFloat creates an optimized float value
func NewOptimizedFloat(v float64) OptimizedValue {
	return OptimizedValue{
		Type:     TypeFloat64,
		FloatVal: v,
	}
}

// NewOptimizedString creates an optimized string value
func NewOptimizedString(v string) OptimizedValue {
	return OptimizedValue{
		Type:      TypeString,
		StringVal: v,
	}
}

// NewOptimizedNil creates an optimized nil value
func NewOptimizedNil() OptimizedValue {
	return OptimizedValue{
		Type: TypeNil,
	}
}

// Type information accessors

// GetType returns the TypeInfo for this value
func (v *OptimizedValue) GetType() TypeInfo {
	switch v.Type {
	case TypeBool:
		return TypeInfo{Kind: KindBool, Name: "bool", Size: 1}
	case TypeInt64:
		return TypeInfo{Kind: KindInt64, Name: "int", Size: 8}
	case TypeFloat64:
		return TypeInfo{Kind: KindFloat64, Name: "float", Size: 8}
	case TypeString:
		return TypeInfo{Kind: KindString, Name: "string", Size: -1}
	case TypeNil:
		return TypeInfo{Kind: KindNil, Name: "nil", Size: 0}
	default:
		return TypeInfo{Kind: KindInterface, Name: "unknown", Size: -1}
	}
}

// String returns string representation
func (v *OptimizedValue) String() string {
	switch v.Type {
	case TypeBool:
		return strconv.FormatBool(v.BoolVal)
	case TypeInt64:
		return strconv.FormatInt(v.IntVal, 10)
	case TypeFloat64:
		return strconv.FormatFloat(v.FloatVal, 'g', -1, 64)
	case TypeString:
		return v.StringVal
	case TypeNil:
		return "nil"
	default:
		return "unknown"
	}
}

// Equal compares two optimized values for equality
func (v *OptimizedValue) Equal(other *OptimizedValue) bool {
	if v.Type != other.Type {
		return false
	}

	switch v.Type {
	case TypeBool:
		return v.BoolVal == other.BoolVal
	case TypeInt64:
		return v.IntVal == other.IntVal
	case TypeFloat64:
		return v.FloatVal == other.FloatVal
	case TypeString:
		return v.StringVal == other.StringVal
	case TypeNil:
		return true
	default:
		return false
	}
}

// Hash returns a hash value for this value
func (v *OptimizedValue) Hash() uint64 {
	switch v.Type {
	case TypeBool:
		if v.BoolVal {
			return 1
		}
		return 0
	case TypeInt64:
		return uint64(v.IntVal)
	case TypeFloat64:
		h := fnv.New64a()
		h.Write([]byte(strconv.FormatFloat(v.FloatVal, 'g', -1, 64)))
		return h.Sum64()
	case TypeString:
		h := fnv.New64a()
		h.Write([]byte(v.StringVal))
		return h.Sum64()
	case TypeNil:
		return 0
	default:
		return 0
	}
}

// Accessors for specific types

// IsBool returns true if this is a boolean value
func (v *OptimizedValue) IsBool() bool {
	return v.Type == TypeBool
}

// IsInt returns true if this is an integer value
func (v *OptimizedValue) IsInt() bool {
	return v.Type == TypeInt64
}

// IsFloat returns true if this is a float value
func (v *OptimizedValue) IsFloat() bool {
	return v.Type == TypeFloat64
}

// IsString returns true if this is a string value
func (v *OptimizedValue) IsString() bool {
	return v.Type == TypeString
}

// IsNil returns true if this is a nil value
func (v *OptimizedValue) IsNil() bool {
	return v.Type == TypeNil
}

// GetBool returns the boolean value (assumes IsBool() is true)
func (v *OptimizedValue) GetBool() bool {
	return v.BoolVal
}

// GetInt returns the integer value (assumes IsInt() is true)
func (v *OptimizedValue) GetInt() int64 {
	return v.IntVal
}

// GetFloat returns the float value (assumes IsFloat() is true)
func (v *OptimizedValue) GetFloat() float64 {
	return v.FloatVal
}

// GetString returns the string value (assumes IsString() is true)
func (v *OptimizedValue) GetString() string {
	return v.StringVal
}

// Performance-critical operations optimized for the union type

// AddOptimized performs optimized addition for union types
func AddOptimized(left, right *OptimizedValue) (OptimizedValue, error) {
	// Integer + Integer (most common case)
	if left.Type == TypeInt64 && right.Type == TypeInt64 {
		return NewOptimizedInt(left.IntVal + right.IntVal), nil
	}

	// Float + Float
	if left.Type == TypeFloat64 && right.Type == TypeFloat64 {
		return NewOptimizedFloat(left.FloatVal + right.FloatVal), nil
	}

	// Int + Float
	if left.Type == TypeInt64 && right.Type == TypeFloat64 {
		return NewOptimizedFloat(float64(left.IntVal) + right.FloatVal), nil
	}

	// Float + Int
	if left.Type == TypeFloat64 && right.Type == TypeInt64 {
		return NewOptimizedFloat(left.FloatVal + float64(right.IntVal)), nil
	}

	// String + String
	if left.Type == TypeString && right.Type == TypeString {
		return NewOptimizedString(left.StringVal + right.StringVal), nil
	}

	return NewOptimizedNil(), fmt.Errorf("unsupported addition: %v + %v", left.Type, right.Type)
}

// CompareOptimized performs optimized comparison for union types
func CompareOptimized(left, right *OptimizedValue, op string) (OptimizedValue, error) {
	// Same type comparisons (most common)
	if left.Type == right.Type {
		switch left.Type {
		case TypeInt64:
			switch op {
			case "==":
				return NewOptimizedBool(left.IntVal == right.IntVal), nil
			case "!=":
				return NewOptimizedBool(left.IntVal != right.IntVal), nil
			case ">":
				return NewOptimizedBool(left.IntVal > right.IntVal), nil
			case ">=":
				return NewOptimizedBool(left.IntVal >= right.IntVal), nil
			case "<":
				return NewOptimizedBool(left.IntVal < right.IntVal), nil
			case "<=":
				return NewOptimizedBool(left.IntVal <= right.IntVal), nil
			}
		case TypeFloat64:
			switch op {
			case "==":
				return NewOptimizedBool(left.FloatVal == right.FloatVal), nil
			case "!=":
				return NewOptimizedBool(left.FloatVal != right.FloatVal), nil
			case ">":
				return NewOptimizedBool(left.FloatVal > right.FloatVal), nil
			case ">=":
				return NewOptimizedBool(left.FloatVal >= right.FloatVal), nil
			case "<":
				return NewOptimizedBool(left.FloatVal < right.FloatVal), nil
			case "<=":
				return NewOptimizedBool(left.FloatVal <= right.FloatVal), nil
			}
		case TypeString:
			switch op {
			case "==":
				return NewOptimizedBool(left.StringVal == right.StringVal), nil
			case "!=":
				return NewOptimizedBool(left.StringVal != right.StringVal), nil
			}
		case TypeBool:
			switch op {
			case "==":
				return NewOptimizedBool(left.BoolVal == right.BoolVal), nil
			case "!=":
				return NewOptimizedBool(left.BoolVal != right.BoolVal), nil
			}
		case TypeNil:
			switch op {
			case "==":
				return NewOptimizedBool(true), nil
			case "!=":
				return NewOptimizedBool(false), nil
			}
		}
	}

	// Cross-type comparisons
	if (left.Type == TypeInt64 && right.Type == TypeFloat64) ||
		(left.Type == TypeFloat64 && right.Type == TypeInt64) {
		var leftFloat, rightFloat float64
		if left.Type == TypeInt64 {
			leftFloat = float64(left.IntVal)
			rightFloat = right.FloatVal
		} else {
			leftFloat = left.FloatVal
			rightFloat = float64(right.IntVal)
		}

		switch op {
		case "==":
			return NewOptimizedBool(leftFloat == rightFloat), nil
		case "!=":
			return NewOptimizedBool(leftFloat != rightFloat), nil
		case ">":
			return NewOptimizedBool(leftFloat > rightFloat), nil
		case ">=":
			return NewOptimizedBool(leftFloat >= rightFloat), nil
		case "<":
			return NewOptimizedBool(leftFloat < rightFloat), nil
		case "<=":
			return NewOptimizedBool(leftFloat <= rightFloat), nil
		}
	}

	// Handle nil comparisons
	if left.Type == TypeNil || right.Type == TypeNil {
		switch op {
		case "==":
			return NewOptimizedBool(left.Type == TypeNil && right.Type == TypeNil), nil
		case "!=":
			return NewOptimizedBool(!(left.Type == TypeNil && right.Type == TypeNil)), nil
		default:
			return NewOptimizedBool(false), nil
		}
	}

	return NewOptimizedNil(), fmt.Errorf("unsupported comparison: %v %s %v", left.Type, op, right.Type)
}

// ToBool converts an optimized value to boolean for truthiness testing
func (v *OptimizedValue) ToBool() bool {
	switch v.Type {
	case TypeBool:
		return v.BoolVal
	case TypeInt64:
		return v.IntVal != 0
	case TypeFloat64:
		return v.FloatVal != 0.0
	case TypeString:
		return v.StringVal != ""
	case TypeNil:
		return false
	default:
		return false
	}
}

// IsNumeric returns true if this value is numeric (int or float)
func (v *OptimizedValue) IsNumeric() bool {
	return v.Type == TypeInt64 || v.Type == TypeFloat64
}

// ToFloat64 converts numeric values to float64
func (v *OptimizedValue) ToFloat64() (float64, bool) {
	switch v.Type {
	case TypeInt64:
		return float64(v.IntVal), true
	case TypeFloat64:
		return v.FloatVal, true
	default:
		return 0, false
	}
}

// ToInt64 converts numeric values to int64
func (v *OptimizedValue) ToInt64() (int64, bool) {
	switch v.Type {
	case TypeInt64:
		return v.IntVal, true
	case TypeFloat64:
		return int64(v.FloatVal), true
	default:
		return 0, false
	}
}
