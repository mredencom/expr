 # Types模块 - 零反射类型系统

## 概述

Types模块是表达式引擎的零反射类型系统，提供了一套完整的类型表示和转换机制。它摒弃了传统的反射机制，通过预编译的类型信息和专门的值类型实现，达到极高的性能表现。

## 核心设计理念

### 1. 零反射架构
- 预编译类型信息，运行时无反射调用
- 静态类型检查，编译时发现类型错误
- 专用值类型，避免interface{}装箱开销

### 2. 高性能类型系统
- 内联类型转换，减少函数调用开销
- 值缓存池，重用常用值对象
- 快速类型判断，避免类型断言

### 3. 类型安全保证
- 强类型检查，防止运行时类型错误
- 明确的类型转换规则
- 完整的类型兼容性检查

## 核心接口和类型

### Value接口
```go
type Value interface {
    Type() TypeInfo        // 获取类型信息
    String() string        // 字符串表示
    Equal(other Value) bool // 值相等比较
    Hash() uint64          // 计算哈希值
}
```

### TypeInfo结构体
```go
type TypeInfo struct {
    Kind TypeKind  // 类型种类
    Name string    // 类型名称
    Size int       // 类型大小（字节，-1表示动态大小）
}

type TypeKind int

const (
    KindBool TypeKind = iota
    KindInt64
    KindFloat64
    KindString
    KindSlice
    KindMap
    KindFunc
    KindNil
)
```

## 基本值类型

### 1. 布尔值类型
```go
type BoolValue struct {
    value bool
}

// 创建布尔值
func NewBool(v bool) *BoolValue

// 基本使用
boolVal := types.NewBool(true)
fmt.Println(boolVal.String())     // "true"
fmt.Println(boolVal.Type().Name)  // "bool"
fmt.Println(boolVal.Value())      // true
```

### 2. 整数值类型
```go
type IntValue struct {
    value int64
}

// 创建整数值
func NewInt(v int64) *IntValue

// 基本使用
intVal := types.NewInt(42)
fmt.Println(intVal.String())     // "42"
fmt.Println(intVal.Type().Name)  // "int"
fmt.Println(intVal.Value())      // 42
```

### 3. 浮点数值类型
```go
type FloatValue struct {
    value float64
}

// 创建浮点数值
func NewFloat(v float64) *FloatValue

// 基本使用
floatVal := types.NewFloat(3.14)
fmt.Println(floatVal.String())     // "3.14"
fmt.Println(floatVal.Type().Name)  // "float"
fmt.Println(floatVal.Value())      // 3.14
```

### 4. 字符串值类型
```go
type StringValue struct {
    value string
}

// 创建字符串值
func NewString(v string) *StringValue

// 基本使用
strVal := types.NewString("hello")
fmt.Println(strVal.String())     // "hello"
fmt.Println(strVal.Type().Name)  // "string"
fmt.Println(strVal.Value())      // "hello"
```

## 复合值类型

### 1. 切片值类型
```go
type SliceValue struct {
    values   []Value
    elemType TypeInfo
}

// 创建切片值
func NewSlice(values []Value, elemType TypeInfo) *SliceValue

// 基本使用
elements := []types.Value{
    types.NewInt(1),
    types.NewInt(2),
    types.NewInt(3),
}
sliceVal := types.NewSlice(elements, types.TypeInfo{
    Kind: types.KindInt64,
    Name: "int",
    Size: 8,
})

fmt.Println(sliceVal.Len())      // 3
fmt.Println(sliceVal.Get(0))     // 获取第一个元素
```

### 2. 映射值类型
```go
type MapValue struct {
    values           map[string]Value
    keyType, valType TypeInfo
}

// 创建映射值
func NewMap(values map[string]Value, keyType, valType TypeInfo) *MapValue

// 基本使用
mapData := map[string]types.Value{
    "name": types.NewString("Alice"),
    "age":  types.NewInt(30),
}
mapVal := types.NewMap(mapData, 
    types.TypeInfo{Kind: types.KindString, Name: "string"}, 
    types.TypeInfo{Kind: types.KindString, Name: "interface{}"},
)

fmt.Println(mapVal.Get("name"))  // "Alice"
fmt.Println(mapVal.Has("age"))   // true
fmt.Println(mapVal.Len())        // 2
```

### 3. 函数值类型
```go
type FuncValue struct {
    parameters []string         // 参数名列表
    body       interface{}      // 函数体（AST或字节码）
    closure    map[string]Value // 闭包变量
    name       string           // 函数名（可选）
}

// 创建函数值
func NewFunc(parameters []string, body interface{}, closure map[string]Value, name string) *FuncValue

// 基本使用
funcVal := types.NewFunc(
    []string{"x", "y"},  // 参数
    nil,                 // 函数体（由编译器设置）
    nil,                 // 闭包
    "add",              // 函数名
)

fmt.Println(funcVal.Parameters())  // ["x", "y"]
fmt.Println(funcVal.Name())        // "add"
```

## 类型转换系统

### 1. 自动类型转换
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/types"
)

func main() {
    // 数值类型转换
    intVal := types.NewInt(42)
    floatVal, err := types.ConvertToFloat(intVal)
    if err != nil {
        panic(err)
    }
    fmt.Println(floatVal.Value()) // 42.0
    
    // 字符串转换
    strVal, err := types.ConvertToString(intVal)
    if err != nil {
        panic(err)
    }
    fmt.Println(strVal.Value()) // "42"
}
```

### 2. 类型兼容性检查
```go
func checkTypeCompatibility() {
    intType := types.TypeInfo{Kind: types.KindInt64, Name: "int"}
    floatType := types.TypeInfo{Kind: types.KindFloat64, Name: "float"}
    stringType := types.TypeInfo{Kind: types.KindString, Name: "string"}
    
    // 检查数值类型兼容性
    fmt.Println(types.IsNumericCompatible(intType, floatType))     // true
    fmt.Println(types.IsNumericCompatible(intType, stringType))    // false
    
    // 检查可比较性
    fmt.Println(types.IsComparable(intType))                       // true
    fmt.Println(types.IsComparable(stringType))                    // true
}
```

### 3. 批量类型转换
```go
func batchConversion() {
    values := []types.Value{
        types.NewInt(1),
        types.NewFloat(2.5),
        types.NewString("3"),
    }
    
    // 转换为统一类型（字符串）
    converted := make([]types.Value, len(values))
    for i, val := range values {
        if strVal, err := types.ConvertToString(val); err == nil {
            converted[i] = strVal
        }
    }
    
    for _, val := range converted {
        fmt.Println(val.String())
    }
}
```

## 高级特性

### 1. 值缓存池
```go
// 常用值的缓存池
var (
    // 缓存小整数值
    IntCache [256]*IntValue
    
    // 缓存常用布尔值
    TrueValue  = NewBool(true)
    FalseValue = NewBool(false)
    
    // 缓存空值
    NilValue = NewNil()
)

// 高效创建整数值
func GetIntValue(n int64) *IntValue {
    if n >= 0 && n < 256 {
        return IntCache[n]  // 零分配
    }
    return NewInt(n)
}

// 高效创建布尔值
func GetBoolValue(b bool) *BoolValue {
    if b {
        return TrueValue
    }
    return FalseValue
}
```

### 2. 类型推断系统
```go
type TypeInferrer struct {
    context map[string]TypeInfo
}

func NewTypeInferrer() *TypeInferrer {
    return &TypeInferrer{
        context: make(map[string]TypeInfo),
    }
}

func (ti *TypeInferrer) InferBinaryExpression(leftType, rightType TypeInfo, operator string) (TypeInfo, error) {
    switch operator {
    case "+", "-", "*", "/":
        // 算术运算类型推断
        if types.IsNumeric(leftType) && types.IsNumeric(rightType) {
            if leftType.Kind == types.KindFloat64 || rightType.Kind == types.KindFloat64 {
                return types.TypeInfo{Kind: types.KindFloat64, Name: "float"}, nil
            }
            return types.TypeInfo{Kind: types.KindInt64, Name: "int"}, nil
        }
        if operator == "+" && (leftType.Kind == types.KindString || rightType.Kind == types.KindString) {
            return types.TypeInfo{Kind: types.KindString, Name: "string"}, nil
        }
        return types.TypeInfo{}, fmt.Errorf("invalid operands for %s", operator)
        
    case "==", "!=", "<", "<=", ">", ">=":
        // 比较运算类型推断
        if types.IsComparable(leftType) && types.IsComparable(rightType) {
            return types.TypeInfo{Kind: types.KindBool, Name: "bool"}, nil
        }
        return types.TypeInfo{}, fmt.Errorf("incomparable types")
        
    case "&&", "||":
        // 逻辑运算类型推断
        if leftType.Kind == types.KindBool && rightType.Kind == types.KindBool {
            return types.TypeInfo{Kind: types.KindBool, Name: "bool"}, nil
        }
        return types.TypeInfo{}, fmt.Errorf("logical operators require boolean operands")
    }
    
    return types.TypeInfo{}, fmt.Errorf("unsupported operator: %s", operator)
}
```

### 3. 动态类型适配器
```go
type TypeAdapter struct {
    converters map[string]func(Value) (Value, error)
}

func NewTypeAdapter() *TypeAdapter {
    ta := &TypeAdapter{
        converters: make(map[string]func(Value) (Value, error)),
    }
    ta.registerBuiltinConverters()
    return ta
}

func (ta *TypeAdapter) registerBuiltinConverters() {
    // 注册内置类型转换器
    ta.converters["int->string"] = func(v Value) (Value, error) {
        if intVal, ok := v.(*IntValue); ok {
            return NewString(strconv.FormatInt(intVal.Value(), 10)), nil
        }
        return nil, fmt.Errorf("expected int value")
    }
    
    ta.converters["string->int"] = func(v Value) (Value, error) {
        if strVal, ok := v.(*StringValue); ok {
            if n, err := strconv.ParseInt(strVal.Value(), 10, 64); err == nil {
                return NewInt(n), nil
            }
            return nil, fmt.Errorf("invalid number format")
        }
        return nil, fmt.Errorf("expected string value")
    }
}

func (ta *TypeAdapter) Convert(value Value, targetType TypeInfo) (Value, error) {
    sourceType := value.Type()
    key := fmt.Sprintf("%s->%s", sourceType.Name, targetType.Name)
    
    if converter, exists := ta.converters[key]; exists {
        return converter(value)
    }
    
    return nil, fmt.Errorf("no converter from %s to %s", sourceType.Name, targetType.Name)
}
```

## 性能优化技术

### 1. 内存池管理
```go
type ValuePool struct {
    intPool    sync.Pool
    floatPool  sync.Pool
    stringPool sync.Pool
    slicePool  sync.Pool
    mapPool    sync.Pool
}

func NewValuePool() *ValuePool {
    return &ValuePool{
        intPool: sync.Pool{
            New: func() interface{} {
                return &IntValue{}
            },
        },
        floatPool: sync.Pool{
            New: func() interface{} {
                return &FloatValue{}
            },
        },
        stringPool: sync.Pool{
            New: func() interface{} {
                return &StringValue{}
            },
        },
    }
}

func (vp *ValuePool) GetInt(value int64) *IntValue {
    intVal := vp.intPool.Get().(*IntValue)
    intVal.value = value
    return intVal
}

func (vp *ValuePool) PutInt(intVal *IntValue) {
    intVal.value = 0
    vp.intPool.Put(intVal)
}
```

### 2. 快速类型判断
```go
// 使用位掩码进行快速类型判断
const (
    NumericTypeMask = (1 << types.KindInt64) | (1 << types.KindFloat64)
    ComparableTypeMask = NumericTypeMask | (1 << types.KindString) | (1 << types.KindBool)
    IterableTypeMask = (1 << types.KindSlice) | (1 << types.KindMap) | (1 << types.KindString)
)

func IsNumeric(typeInfo TypeInfo) bool {
    return (1 << typeInfo.Kind) & NumericTypeMask != 0
}

func IsComparable(typeInfo TypeInfo) bool {
    return (1 << typeInfo.Kind) & ComparableTypeMask != 0
}

func IsIterable(typeInfo TypeInfo) bool {
    return (1 << typeInfo.Kind) & IterableTypeMask != 0
}
```

### 3. 内联类型转换
```go
// 针对热点路径的内联优化
func FastIntToFloat(intVal *IntValue) *FloatValue {
    // 直接内联，避免函数调用开销
    return &FloatValue{value: float64(intVal.value)}
}

func FastBoolToString(boolVal *BoolValue) *StringValue {
    if boolVal.value {
        return PrebuiltTrueString  // 预构建的常量
    }
    return PrebuiltFalseString
}

// 预构建的常用字符串值
var (
    PrebuiltTrueString  = &StringValue{value: "true"}
    PrebuiltFalseString = &StringValue{value: "false"}
    PrebuiltEmptyString = &StringValue{value: ""}
)
```

## 实际应用示例

### 1. 表达式值计算器
```go
type ValueCalculator struct {
    typeAdapter *TypeAdapter
}

func NewValueCalculator() *ValueCalculator {
    return &ValueCalculator{
        typeAdapter: NewTypeAdapter(),
    }
}

func (vc *ValueCalculator) Add(left, right Value) (Value, error) {
    leftType := left.Type()
    rightType := right.Type()
    
    // 数值加法
    if IsNumeric(leftType) && IsNumeric(rightType) {
        if leftType.Kind == KindFloat64 || rightType.Kind == KindFloat64 {
            leftFloat, _ := vc.ensureFloat(left)
            rightFloat, _ := vc.ensureFloat(right)
            return NewFloat(leftFloat.Value() + rightFloat.Value()), nil
        }
        
        leftInt := left.(*IntValue)
        rightInt := right.(*IntValue)
        return NewInt(leftInt.Value() + rightInt.Value()), nil
    }
    
    // 字符串连接
    if leftType.Kind == KindString || rightType.Kind == KindString {
        leftStr, _ := vc.ensureString(left)
        rightStr, _ := vc.ensureString(right)
        return NewString(leftStr.Value() + rightStr.Value()), nil
    }
    
    return nil, fmt.Errorf("unsupported addition between %s and %s", leftType.Name, rightType.Name)
}

func (vc *ValueCalculator) ensureFloat(value Value) (*FloatValue, error) {
    if floatVal, ok := value.(*FloatValue); ok {
        return floatVal, nil
    }
    if intVal, ok := value.(*IntValue); ok {
        return NewFloat(float64(intVal.Value())), nil
    }
    return nil, fmt.Errorf("cannot convert to float")
}

func (vc *ValueCalculator) ensureString(value Value) (*StringValue, error) {
    if strVal, ok := value.(*StringValue); ok {
        return strVal, nil
    }
    return NewString(value.String()), nil
}
```

### 2. 集合操作处理器
```go
type CollectionProcessor struct{}

func (cp *CollectionProcessor) Filter(slice *SliceValue, predicate func(Value) bool) *SliceValue {
    var filtered []Value
    
    for i := 0; i < slice.Len(); i++ {
        element := slice.Get(i)
        if predicate(element) {
            filtered = append(filtered, element)
        }
    }
    
    return NewSlice(filtered, slice.elemType)
}

func (cp *CollectionProcessor) Map(slice *SliceValue, transform func(Value) Value) *SliceValue {
    mapped := make([]Value, slice.Len())
    
    for i := 0; i < slice.Len(); i++ {
        element := slice.Get(i)
        mapped[i] = transform(element)
    }
    
    // 推断新的元素类型
    var newElemType TypeInfo
    if len(mapped) > 0 {
        newElemType = mapped[0].Type()
    } else {
        newElemType = slice.elemType
    }
    
    return NewSlice(mapped, newElemType)
}

func (cp *CollectionProcessor) Reduce(slice *SliceValue, reducer func(Value, Value) Value, initial Value) Value {
    result := initial
    
    for i := 0; i < slice.Len(); i++ {
        element := slice.Get(i)
        result = reducer(result, element)
    }
    
    return result
}
```

## 最佳实践

### 1. 类型安全编程
```go
// 好的做法：使用类型检查
func safeOperation(value Value) error {
    switch v := value.(type) {
    case *IntValue:
        // 处理整数
        fmt.Printf("Integer: %d\n", v.Value())
    case *StringValue:
        // 处理字符串
        fmt.Printf("String: %s\n", v.Value())
    default:
        return fmt.Errorf("unsupported type: %s", value.Type().Name)
    }
    return nil
}

// 避免的做法：盲目类型断言
func unsafeOperation(value Value) {
    intVal := value.(*IntValue)  // 可能panic
    fmt.Println(intVal.Value())
}
```

### 2. 性能优化
```go
// 使用值缓存池
func optimizedIntCreation(n int64) *IntValue {
    if n >= 0 && n < 256 {
        return IntCache[n]  // 重用缓存值
    }
    return NewInt(n)  // 创建新值
}

// 批量操作优化
func batchProcess(values []Value) []Value {
    result := make([]Value, 0, len(values))  // 预分配容量
    
    for _, value := range values {
        // 处理逻辑
        processed := processValue(value)
        result = append(result, processed)
    }
    
    return result
}
```

### 3. 错误处理
```go
func robustTypeConversion(value Value, targetType TypeInfo) (Value, error) {
    // 检查输入有效性
    if value == nil {
        return NewNil(), nil
    }
    
    sourceType := value.Type()
    
    // 相同类型直接返回
    if sourceType.Kind == targetType.Kind {
        return value, nil
    }
    
    // 尝试转换
    converter := getConverter(sourceType, targetType)
    if converter == nil {
        return nil, fmt.Errorf("no conversion from %s to %s", sourceType.Name, targetType.Name)
    }
    
    return converter(value)
}
```

## 与其他模块的集成

Types模块在整个架构中的作用：

```
表达式输入 → Lexer → Parser → AST → Checker(使用Types) → Compiler → VM(使用Types) → 结果
```

Types模块是整个表达式引擎的类型基础，它：
1. 为Checker模块提供类型检查支持
2. 为Compiler模块提供类型信息
3. 为VM模块提供运行时值表示
4. 为Builtins模块提供函数参数和返回值类型
5. 为API模块提供类型转换功能

通过零反射的设计，Types模块实现了高性能的类型系统，是整个表达式引擎性能优势的核心基础。