# 无反射类型系统设计

## 概述

在不使用反射的情况下，我们需要构建一个静态类型系统，能够在编译时确定所有类型信息并生成高效的执行代码。

## 核心类型接口

### 基础值接口

```go
// Value 表示任何可以被表达式处理的值
type Value interface {
    Type() TypeInfo
    String() string
    Equal(other Value) bool
    Hash() uint64
}

// TypeInfo 包含类型的静态信息
type TypeInfo struct {
    Kind     TypeKind
    Name     string
    Size     int
    Methods  []MethodInfo
    Fields   []FieldInfo
}

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
)
```

### 具体类型实现

```go
// 基础类型包装器
type BoolValue struct {
    value bool
}

func (b BoolValue) Type() TypeInfo {
    return TypeInfo{Kind: KindBool, Name: "bool", Size: 1}
}

type IntValue struct {
    value int64
}

func (i IntValue) Type() TypeInfo {
    return TypeInfo{Kind: KindInt64, Name: "int", Size: 8}
}

type StringValue struct {
    value string
}

func (s StringValue) Type() TypeInfo {
    return TypeInfo{Kind: KindString, Name: "string", Size: -1}
}

// 复合类型
type SliceValue struct {
    values []Value
    elemType TypeInfo
}

type MapValue struct {
    values map[string]Value
    keyType, valueType TypeInfo
}

type StructValue struct {
    fields map[string]Value
    typeInfo TypeInfo
}
```

## 类型注册系统

### 类型注册器

```go
type TypeRegistry struct {
    types map[string]*RegisteredType
    converters map[TypePair]*Converter
}

type RegisteredType struct {
    Info     TypeInfo
    Factory  func() Value
    Methods  map[string]*Method
    Fields   map[string]*Field
}

type Method struct {
    Name     string
    Params   []TypeInfo
    Returns  []TypeInfo
    Handler  func(receiver Value, args []Value) ([]Value, error)
}

type Field struct {
    Name     string
    Type     TypeInfo
    Getter   func(receiver Value) Value
    Setter   func(receiver Value, value Value) error
}

// 注册新类型
func (r *TypeRegistry) Register(name string, factory func() Value) *RegisteredType {
    rt := &RegisteredType{
        Info:    factory().Type(),
        Factory: factory,
        Methods: make(map[string]*Method),
        Fields:  make(map[string]*Field),
    }
    r.types[name] = rt
    return rt
}
```

### 类型适配器

```go
// 将 Go 原生类型适配到我们的类型系统
type TypeAdapter struct {
    registry *TypeRegistry
}

func (a *TypeAdapter) Adapt(v interface{}) Value {
    switch val := v.(type) {
    case bool:
        return BoolValue{value: val}
    case int:
        return IntValue{value: int64(val)}
    case int64:
        return IntValue{value: val}
    case string:
        return StringValue{value: val}
    case []interface{}:
        values := make([]Value, len(val))
        for i, item := range val {
            values[i] = a.Adapt(item)
        }
        return SliceValue{values: values}
    case map[string]interface{}:
        values := make(map[string]Value)
        for k, v := range val {
            values[k] = a.Adapt(v)
        }
        return MapValue{values: values}
    default:
        // 对于自定义类型，使用注册的适配器
        return a.adaptCustomType(v)
    }
}
```

## 静态类型检查器

### 类型推导引擎

```go
type TypeChecker struct {
    registry *TypeRegistry
    env      *Environment
    errors   []TypeError
}

type TypeError struct {
    Position Position
    Message  string
}

func (tc *TypeChecker) CheckExpression(expr ASTNode) TypeInfo {
    switch e := expr.(type) {
    case *LiteralNode:
        return tc.checkLiteral(e)
    case *IdentifierNode:
        return tc.checkIdentifier(e)
    case *BinaryOpNode:
        return tc.checkBinaryOp(e)
    case *CallNode:
        return tc.checkCall(e)
    case *MemberNode:
        return tc.checkMember(e)
    default:
        tc.addError(expr.Position(), "Unknown expression type")
        return TypeInfo{}
    }
}

func (tc *TypeChecker) checkBinaryOp(node *BinaryOpNode) TypeInfo {
    leftType := tc.CheckExpression(node.Left)
    rightType := tc.CheckExpression(node.Right)
    
    // 检查操作符兼容性
    if !tc.isOperatorCompatible(node.Operator, leftType, rightType) {
        tc.addError(node.Position(), 
            fmt.Sprintf("invalid operation %s (mismatched types %s and %s)",
                node.Operator, leftType.Name, rightType.Name))
        return TypeInfo{}
    }
    
    return tc.getResultType(node.Operator, leftType, rightType)
}
```

### 类型兼容性检查

```go
func (tc *TypeChecker) isOperatorCompatible(op Operator, left, right TypeInfo) bool {
    switch op {
    case OpAdd:
        return (isNumeric(left) && isNumeric(right)) || 
               (left.Kind == KindString && right.Kind == KindString)
    case OpSub, OpMul, OpDiv:
        return isNumeric(left) && isNumeric(right)
    case OpEq, OpNe:
        return tc.isComparable(left, right)
    case OpLt, OpLe, OpGt, OpGe:
        return tc.isOrdered(left, right)
    case OpAnd, OpOr:
        return left.Kind == KindBool && right.Kind == KindBool
    default:
        return false
    }
}

func isNumeric(t TypeInfo) bool {
    return t.Kind >= KindInt && t.Kind <= KindFloat64
}
```

## 预编译环境生成

### 环境结构体生成器

```go
// 根据用户提供的环境类型，生成优化的访问代码
type EnvGenerator struct {
    registry *TypeRegistry
}

func (eg *EnvGenerator) GenerateEnvironment(envType interface{}) *CompiledEnvironment {
    // 分析环境类型结构
    envInfo := eg.analyzeType(envType)
    
    // 生成字段访问器
    accessors := make(map[string]FieldAccessor)
    for _, field := range envInfo.Fields {
        accessors[field.Name] = eg.generateAccessor(field)
    }
    
    // 生成方法调用器
    methods := make(map[string]MethodCaller)
    for _, method := range envInfo.Methods {
        methods[method.Name] = eg.generateMethodCaller(method)
    }
    
    return &CompiledEnvironment{
        TypeInfo:  envInfo,
        Accessors: accessors,
        Methods:   methods,
    }
}

type CompiledEnvironment struct {
    TypeInfo  TypeInfo
    Accessors map[string]FieldAccessor
    Methods   map[string]MethodCaller
}

type FieldAccessor func(env interface{}) Value
type MethodCaller func(env interface{}, args []Value) (Value, error)
```

### 代码生成

```go
func (eg *EnvGenerator) generateAccessor(field FieldInfo) FieldAccessor {
    switch field.Type.Kind {
    case KindString:
        return func(env interface{}) Value {
            // 直接类型断言，避免反射
            if e, ok := env.(*UserEnvType); ok {
                return StringValue{value: e.StringField}
            }
            panic("type assertion failed")
        }
    case KindInt:
        return func(env interface{}) Value {
            if e, ok := env.(*UserEnvType); ok {
                return IntValue{value: int64(e.IntField)}
            }
            panic("type assertion failed")
        }
    // ... 其他类型
    }
}
```

## 性能优化

### 类型缓存

```go
type TypeCache struct {
    typeInfos map[reflect.Type]TypeInfo
    mutex     sync.RWMutex
}

func (tc *TypeCache) GetTypeInfo(t reflect.Type) TypeInfo {
    tc.mutex.RLock()
    if info, exists := tc.typeInfos[t]; exists {
        tc.mutex.RUnlock()
        return info
    }
    tc.mutex.RUnlock()
    
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    // 双重检查
    if info, exists := tc.typeInfos[t]; exists {
        return info
    }
    
    info := tc.analyzeType(t)
    tc.typeInfos[t] = info
    return info
}
```

### 内联优化

```go
// 对于简单类型操作，生成内联代码
func (compiler *Compiler) compileInlineOperation(op Operator, left, right TypeInfo) []Instruction {
    if left.Kind == KindInt64 && right.Kind == KindInt64 {
        switch op {
        case OpAdd:
            return []Instruction{
                {OpCode: OpAddInt64},
            }
        case OpSub:
            return []Instruction{
                {OpCode: OpSubInt64},
            }
        }
    }
    
    // 回退到通用操作
    return []Instruction{
        {OpCode: OpGenericBinaryOp, Operand: int(op)},
    }
}
```

## 内存管理

### 值对象池

```go
type ValuePool struct {
    boolPool   sync.Pool
    intPool    sync.Pool
    stringPool sync.Pool
}

func (vp *ValuePool) GetBool() *BoolValue {
    if v := vp.boolPool.Get(); v != nil {
        return v.(*BoolValue)
    }
    return &BoolValue{}
}

func (vp *ValuePool) PutBool(v *BoolValue) {
    v.value = false
    vp.boolPool.Put(v)
}
```

这个类型系统设计完全避免了反射的使用，通过预定义的类型接口、静态类型检查和编译时代码生成来实现高性能的表达式求值。 