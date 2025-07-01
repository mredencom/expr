package types

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

// Value represents any value that can be processed by expressions
type Value interface {
	Type() TypeInfo
	String() string
	Equal(other Value) bool
	Hash() uint64
}

// Basic value implementations

// BoolValue represents a boolean value
type BoolValue struct {
	value bool
}

func NewBool(v bool) *BoolValue {
	return &BoolValue{value: v}
}

func (b *BoolValue) Type() TypeInfo {
	return TypeInfo{Kind: KindBool, Name: "bool", Size: 1}
}

func (b *BoolValue) String() string {
	return strconv.FormatBool(b.value)
}

func (b *BoolValue) Equal(other Value) bool {
	if o, ok := other.(*BoolValue); ok {
		return b.value == o.value
	}
	return false
}

func (b *BoolValue) Hash() uint64 {
	if b.value {
		return 1
	}
	return 0
}

func (b *BoolValue) Value() bool {
	return b.value
}

// IntValue represents an integer value (int64)
type IntValue struct {
	value int64
}

func NewInt(v int64) *IntValue {
	return &IntValue{value: v}
}

func (i *IntValue) Type() TypeInfo {
	return TypeInfo{Kind: KindInt64, Name: "int", Size: 8}
}

func (i *IntValue) String() string {
	return strconv.FormatInt(i.value, 10)
}

func (i *IntValue) Equal(other Value) bool {
	if o, ok := other.(*IntValue); ok {
		return i.value == o.value
	}
	return false
}

func (i *IntValue) Hash() uint64 {
	return uint64(i.value)
}

func (i *IntValue) Value() int64 {
	return i.value
}

// FloatValue represents a floating-point value (float64)
type FloatValue struct {
	value float64
}

func NewFloat(v float64) *FloatValue {
	return &FloatValue{value: v}
}

func (f *FloatValue) Type() TypeInfo {
	return TypeInfo{Kind: KindFloat64, Name: "float", Size: 8}
}

func (f *FloatValue) String() string {
	return strconv.FormatFloat(f.value, 'g', -1, 64)
}

func (f *FloatValue) Equal(other Value) bool {
	if o, ok := other.(*FloatValue); ok {
		return f.value == o.value
	}
	return false
}

func (f *FloatValue) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(f.String()))
	return h.Sum64()
}

func (f *FloatValue) Value() float64 {
	return f.value
}

// StringValue represents a string value
type StringValue struct {
	value string
}

func NewString(v string) *StringValue {
	return &StringValue{value: v}
}

func (s *StringValue) Type() TypeInfo {
	return TypeInfo{Kind: KindString, Name: "string", Size: -1}
}

func (s *StringValue) String() string {
	return s.value
}

func (s *StringValue) Equal(other Value) bool {
	if o, ok := other.(*StringValue); ok {
		return s.value == o.value
	}
	return false
}

func (s *StringValue) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.value))
	return h.Sum64()
}

func (s *StringValue) Value() string {
	return s.value
}

// SliceValue represents a slice/array value
type SliceValue struct {
	values   []Value
	elemType TypeInfo
}

func NewSlice(values []Value, elemType TypeInfo) *SliceValue {
	return &SliceValue{values: values, elemType: elemType}
}

func (s *SliceValue) Type() TypeInfo {
	return TypeInfo{Kind: KindSlice, Name: "[]" + s.elemType.Name, Size: -1}
}

func (s *SliceValue) String() string {
	return fmt.Sprintf("%v", s.values)
}

func (s *SliceValue) Equal(other Value) bool {
	if o, ok := other.(*SliceValue); ok {
		if len(s.values) != len(o.values) {
			return false
		}
		for i, v := range s.values {
			if !v.Equal(o.values[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (s *SliceValue) Hash() uint64 {
	h := fnv.New64a()
	for _, v := range s.values {
		h.Write([]byte(strconv.FormatUint(v.Hash(), 10)))
	}
	return h.Sum64()
}

func (s *SliceValue) Values() []Value {
	return s.values
}

func (s *SliceValue) Len() int {
	return len(s.values)
}

func (s *SliceValue) Get(index int) Value {
	if index < 0 || index >= len(s.values) {
		return nil
	}
	return s.values[index]
}

// ElementType returns the element type of the slice
func (s *SliceValue) ElementType() TypeInfo {
	return s.elemType
}

// MapValue represents a map value
type MapValue struct {
	values           map[string]Value
	keyType, valType TypeInfo
}

func NewMap(values map[string]Value, keyType, valType TypeInfo) *MapValue {
	return &MapValue{values: values, keyType: keyType, valType: valType}
}

func (m *MapValue) Type() TypeInfo {
	return TypeInfo{
		Kind: KindMap,
		Name: "map[" + m.keyType.Name + "]" + m.valType.Name,
		Size: -1,
	}
}

func (m *MapValue) String() string {
	return fmt.Sprintf("%v", m.values)
}

func (m *MapValue) Equal(other Value) bool {
	if o, ok := other.(*MapValue); ok {
		if len(m.values) != len(o.values) {
			return false
		}
		for k, v := range m.values {
			if ov, exists := o.values[k]; !exists || !v.Equal(ov) {
				return false
			}
		}
		return true
	}
	return false
}

func (m *MapValue) Hash() uint64 {
	h := fnv.New64a()
	for k, v := range m.values {
		h.Write([]byte(k))
		h.Write([]byte(strconv.FormatUint(v.Hash(), 10)))
	}
	return h.Sum64()
}

func (m *MapValue) Values() map[string]Value {
	return m.values
}

func (m *MapValue) Len() int {
	return len(m.values)
}

// Get returns the value for a key and whether the key exists
func (m *MapValue) Get(key string) (Value, bool) {
	val, exists := m.values[key]
	return val, exists
}

// Keys returns all keys in the map
func (m *MapValue) Keys() []string {
	keys := make([]string, 0, len(m.values))
	for k := range m.values {
		keys = append(keys, k)
	}
	return keys
}

// ValueType returns the value type of the map
func (m *MapValue) ValueType() TypeInfo {
	return m.valType
}

func (m *MapValue) Has(key string) bool {
	_, exists := m.values[key]
	return exists
}

// NilValue represents a nil/null value
type NilValue struct{}

func NewNil() *NilValue {
	return &NilValue{}
}

func (n *NilValue) Type() TypeInfo {
	return TypeInfo{Kind: KindNil, Name: "nil", Size: 0}
}

func (n *NilValue) String() string {
	return "nil"
}

func (n *NilValue) Equal(other Value) bool {
	_, ok := other.(*NilValue)
	return ok
}

func (n *NilValue) Hash() uint64 {
	return 0
}

// FuncValue represents a function/lambda value
type FuncValue struct {
	parameters []string         // Parameter names
	body       interface{}      // AST node for function body or compiled bytecode
	closure    map[string]Value // Captured variables from outer scope
	name       string           // Optional function name
}

func NewFunc(parameters []string, body interface{}, closure map[string]Value, name string) *FuncValue {
	if closure == nil {
		closure = make(map[string]Value)
	}
	return &FuncValue{
		parameters: parameters,
		body:       body,
		closure:    closure,
		name:       name,
	}
}

func (f *FuncValue) Type() TypeInfo {
	return TypeInfo{Kind: KindFunc, Name: "function", Size: -1}
}

func (f *FuncValue) String() string {
	if f.name != "" {
		return f.name
	}
	params := ""
	for i, param := range f.parameters {
		if i > 0 {
			params += ", "
		}
		params += param
	}
	if len(f.parameters) == 1 {
		return params + " => <function>"
	}
	return "(" + params + ") => <function>"
}

func (f *FuncValue) Equal(other Value) bool {
	if o, ok := other.(*FuncValue); ok {
		// Functions are equal if they have same parameters and same body
		// This is a simplified comparison - in practice might want pointer equality
		if len(f.parameters) != len(o.parameters) {
			return false
		}
		for i, param := range f.parameters {
			if param != o.parameters[i] {
				return false
			}
		}
		return f.body == o.body
	}
	return false
}

func (f *FuncValue) Hash() uint64 {
	h := fnv.New64a()
	for _, param := range f.parameters {
		h.Write([]byte(param))
	}
	// Simple hash based on parameters
	return h.Sum64()
}

func (f *FuncValue) Parameters() []string {
	return f.parameters
}

func (f *FuncValue) Body() interface{} {
	return f.body
}

func (f *FuncValue) Closure() map[string]Value {
	return f.closure
}

func (f *FuncValue) Name() string {
	return f.name
}

// PlaceholderExprValue represents a compiled placeholder expression that can be evaluated later
type PlaceholderExprValue struct {
	instructions []byte  // Compiled bytecode for the expression
	constants    []Value // Constants used in the expression
	operator     string  // The operation type (e.g., ">", "*", etc.)
	operand      Value   // The operand (e.g., 3 in "# > 3")
}

func NewPlaceholderExpr(instructions []byte, constants []Value, operator string, operand Value) *PlaceholderExprValue {
	return &PlaceholderExprValue{
		instructions: instructions,
		constants:    constants,
		operator:     operator,
		operand:      operand,
	}
}

func (p *PlaceholderExprValue) Type() TypeInfo {
	return TypeInfo{Kind: KindInterface, Name: "placeholder_expr", Size: -1}
}

func (p *PlaceholderExprValue) String() string {
	return fmt.Sprintf("# %s %v", p.operator, p.operand)
}

func (p *PlaceholderExprValue) Equal(other Value) bool {
	if o, ok := other.(*PlaceholderExprValue); ok {
		return p.operator == o.operator && p.operand.Equal(o.operand)
	}
	return false
}

func (p *PlaceholderExprValue) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(p.operator))
	h.Write([]byte(strconv.FormatUint(p.operand.Hash(), 10)))
	return h.Sum64()
}

func (p *PlaceholderExprValue) Instructions() []byte {
	return p.instructions
}

func (p *PlaceholderExprValue) Constants() []Value {
	return p.constants
}

func (p *PlaceholderExprValue) Operator() string {
	return p.operator
}

func (p *PlaceholderExprValue) Operand() Value {
	return p.operand
}
