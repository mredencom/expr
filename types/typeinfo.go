package types

// TypeKind represents the kind of a type
type TypeKind uint8

const (
	KindBool TypeKind = iota
	KindInt
	KindInt8
	KindInt16
	KindInt32
	KindInt64
	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64
	KindFloat32
	KindFloat64
	KindString
	KindArray
	KindSlice
	KindMap
	KindStruct
	KindInterface
	KindPointer
	KindFunc
	KindNil
	KindUnknown // Added for type inference when type cannot be determined
)

// String returns the string representation of TypeKind
func (k TypeKind) String() string {
	switch k {
	case KindBool:
		return "bool"
	case KindInt:
		return "int"
	case KindInt8:
		return "int8"
	case KindInt16:
		return "int16"
	case KindInt32:
		return "int32"
	case KindInt64:
		return "int64"
	case KindUint:
		return "uint"
	case KindUint8:
		return "uint8"
	case KindUint16:
		return "uint16"
	case KindUint32:
		return "uint32"
	case KindUint64:
		return "uint64"
	case KindFloat32:
		return "float32"
	case KindFloat64:
		return "float64"
	case KindString:
		return "string"
	case KindArray:
		return "array"
	case KindSlice:
		return "slice"
	case KindMap:
		return "map"
	case KindStruct:
		return "struct"
	case KindInterface:
		return "interface"
	case KindPointer:
		return "pointer"
	case KindFunc:
		return "func"
	case KindNil:
		return "nil"
	case KindUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// TypeInfo contains static information about a type
type TypeInfo struct {
	Kind     TypeKind
	Name     string
	Size     int // -1 for variable size
	Methods  []MethodInfo
	Fields   []FieldInfo
	ElemType *TypeInfo // for arrays, slices, pointers
	KeyType  *TypeInfo // for maps
	ValType  *TypeInfo // for maps
}

// MethodInfo contains information about a method
type MethodInfo struct {
	Name     string
	Params   []TypeInfo
	Returns  []TypeInfo
	Variadic bool
}

// FieldInfo contains information about a field
type FieldInfo struct {
	Name   string
	Type   TypeInfo
	Index  int
	Offset int
}

// IsNumeric returns true if the type is numeric
func (t TypeInfo) IsNumeric() bool {
	return t.Kind >= KindInt && t.Kind <= KindFloat64
}

// IsInteger returns true if the type is an integer type
func (t TypeInfo) IsInteger() bool {
	return t.Kind >= KindInt && t.Kind <= KindUint64
}

// IsFloat returns true if the type is a floating-point type
func (t TypeInfo) IsFloat() bool {
	return t.Kind == KindFloat32 || t.Kind == KindFloat64
}

// IsComparable returns true if values of this type can be compared
func (t TypeInfo) IsComparable() bool {
	switch t.Kind {
	case KindBool, KindString:
		return true
	case KindInt, KindInt8, KindInt16, KindInt32, KindInt64:
		return true
	case KindUint, KindUint8, KindUint16, KindUint32, KindUint64:
		return true
	case KindFloat32, KindFloat64:
		return true
	case KindArray:
		return t.ElemType != nil && t.ElemType.IsComparable()
	case KindStruct:
		for _, field := range t.Fields {
			if !field.Type.IsComparable() {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// IsOrdered returns true if values of this type can be ordered (support <, >, etc.)
func (t TypeInfo) IsOrdered() bool {
	switch t.Kind {
	case KindString:
		return true
	case KindInt, KindInt8, KindInt16, KindInt32, KindInt64:
		return true
	case KindUint, KindUint8, KindUint16, KindUint32, KindUint64:
		return true
	case KindFloat32, KindFloat64:
		return true
	default:
		return false
	}
}

// Compatible returns true if this type is compatible with another type
func (t TypeInfo) Compatible(other TypeInfo) bool {
	if t.Kind == other.Kind {
		return true
	}

	// Numeric type conversions
	if t.IsNumeric() && other.IsNumeric() {
		return true
	}

	// Interface compatibility
	if t.Kind == KindInterface || other.Kind == KindInterface {
		return true
	}

	return false
}

// Assignable returns true if a value of 'other' type can be assigned to this type
func (t TypeInfo) Assignable(other TypeInfo) bool {
	if t.Kind == other.Kind {
		return true
	}

	// Interface can accept any type
	if t.Kind == KindInterface {
		return true
	}

	// Nil can be assigned to pointers, slices, maps, interfaces
	if other.Kind == KindNil {
		switch t.Kind {
		case KindPointer, KindSlice, KindMap, KindInterface:
			return true
		}
	}

	return false
}

// String returns the string representation of TypeInfo
func (t TypeInfo) String() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Kind.String()
}

// Predefined type info instances
var (
	BoolType = TypeInfo{
		Kind: KindBool,
		Name: "bool",
		Size: 1,
	}

	IntType = TypeInfo{
		Kind: KindInt64,
		Name: "int",
		Size: 8,
	}

	FloatType = TypeInfo{
		Kind: KindFloat64,
		Name: "float",
		Size: 8,
	}

	StringType = TypeInfo{
		Kind: KindString,
		Name: "string",
		Size: -1,
	}

	NilType = TypeInfo{
		Kind: KindNil,
		Name: "nil",
		Size: 0,
	}
)
