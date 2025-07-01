# Env模块 - 环境适配器

## 概述

Env模块是表达式引擎的环境适配器，负责将Go语言的原生数据结构转换为表达式引擎内部的值类型，实现零反射的高性能数据访问。它支持结构体、映射、切片等复杂数据类型的自动适配。

## 核心功能

### 1. 数据类型适配
- Go原生类型到Value类型转换
- 结构体字段访问适配
- 切片/数组索引访问
- 映射键值访问

### 2. 零反射实现
- 预编译访问器生成
- 静态类型信息缓存
- 高性能字段访问
- 避免运行时反射开销

### 3. 动态环境支持
- 运行时环境变量注入
- 多层嵌套环境
- 变量作用域管理
- 类型安全访问

## 主要类型

### Adapter结构体
```go
type Adapter struct {
    typeAdapters map[reflect.Type]*TypeAdapter
    fieldCache   map[string]*FieldAccessor
    indexCache   map[reflect.Type]*IndexAccessor
}

type TypeAdapter struct {
    Type        reflect.Type
    FieldMap    map[string]*FieldAccessor
    IndexFunc   func(obj interface{}, index int) types.Value
    ConvertFunc func(obj interface{}) types.Value
}

type FieldAccessor struct {
    Name     string
    Type     types.TypeInfo
    GetFunc  func(obj interface{}) types.Value
    Offset   uintptr  // 字段偏移量（优化用）
}
```

## 基本使用

### 1. 简单数据适配
```go
func main() {
    adapter := env.New()
    
    // 基本类型适配
    data := map[string]interface{}{
        "name": "Alice",
        "age":  30,
        "active": true,
    }
    
    // 转换为表达式值
    nameVal := adapter.Convert(data["name"])
    ageVal := adapter.Convert(data["age"])
    
    fmt.Printf("Name: %s, Age: %s\n", nameVal.String(), ageVal.String())
}
```

### 2. 结构体适配
```go
type User struct {
    Name    string `expr:"name"`
    Age     int    `expr:"age"`
    Email   string `expr:"email"`
    Profile *Profile `expr:"profile"`
}

type Profile struct {
    Bio     string `expr:"bio"`
    Website string `expr:"website"`
}

func structExample() {
    adapter := env.New()
    
    user := &User{
        Name:  "Alice",
        Age:   30,
        Email: "alice@example.com",
        Profile: &Profile{
            Bio:     "Software Engineer",
            Website: "https://alice.dev",
        },
    }
    
    // 注册结构体类型
    adapter.RegisterType(reflect.TypeOf(user))
    
    // 访问字段
    nameVal := adapter.GetField(user, "name")
    bioVal := adapter.GetField(user.Profile, "bio")
    
    fmt.Printf("Name: %s, Bio: %s\n", nameVal.String(), bioVal.String())
}
```

### 3. 集合类型适配
```go
func collectionExample() {
    adapter := env.New()
    
    // 切片适配
    numbers := []int{1, 2, 3, 4, 5}
    sliceVal := adapter.Convert(numbers)
    
    // 访问元素
    firstVal := adapter.GetIndex(sliceVal, 0)
    
    // 映射适配
    userMap := map[string]string{
        "name":  "Bob",
        "email": "bob@example.com",
    }
    mapVal := adapter.Convert(userMap)
    
    // 访问键
    emailVal := adapter.GetMapValue(mapVal, "email")
    
    fmt.Printf("First: %s, Email: %s\n", firstVal.String(), emailVal.String())
}
```

## 高级特性

### 1. 自定义类型适配器
```go
func registerCustomAdapter() {
    adapter := env.New()
    
    // 注册时间类型适配器
    adapter.RegisterCustomAdapter(reflect.TypeOf(time.Time{}), &env.TypeAdapter{
        ConvertFunc: func(obj interface{}) types.Value {
            t := obj.(time.Time)
            return types.NewString(t.Format("2006-01-02 15:04:05"))
        },
        FieldMap: map[string]*env.FieldAccessor{
            "year": {
                Name: "year",
                Type: types.TypeInfo{Kind: types.KindInt64, Name: "int"},
                GetFunc: func(obj interface{}) types.Value {
                    t := obj.(time.Time)
                    return types.NewInt(int64(t.Year()))
                },
            },
            "month": {
                Name: "month",
                Type: types.TypeInfo{Kind: types.KindInt64, Name: "int"},
                GetFunc: func(obj interface{}) types.Value {
                    t := obj.(time.Time)
                    return types.NewInt(int64(t.Month()))
                },
            },
        },
    })
}
```

### 2. 性能优化访问器
```go
type FastAccessor struct {
    fieldOffsets map[string]uintptr
    structType   reflect.Type
}

func (fa *FastAccessor) GetFieldDirect(ptr unsafe.Pointer, fieldName string) types.Value {
    offset := fa.fieldOffsets[fieldName]
    fieldPtr := unsafe.Pointer(uintptr(ptr) + offset)
    
    // 直接内存访问（需要类型信息）
    switch fieldName {
    case "name":
        strPtr := (*string)(fieldPtr)
        return types.NewString(*strPtr)
    case "age":
        intPtr := (*int)(fieldPtr)
        return types.NewInt(int64(*intPtr))
    }
    
    return types.NewNil()
}
```

### 3. 环境变量管理
```go
type Environment struct {
    variables map[string]types.Value
    parent    *Environment
    adapter   *Adapter
}

func (e *Environment) Define(name string, value interface{}) {
    e.variables[name] = e.adapter.Convert(value)
}

func (e *Environment) Get(name string) (types.Value, bool) {
    // 当前环境查找
    if value, exists := e.variables[name]; exists {
        return value, true
    }
    
    // 父环境查找
    if e.parent != nil {
        return e.parent.Get(name)
    }
    
    return types.NewNil(), false
}

func (e *Environment) Child() *Environment {
    return &Environment{
        variables: make(map[string]types.Value),
        parent:    e,
        adapter:   e.adapter,
    }
}
```

## 类型转换策略

### 1. 自动类型推断
```go
func (a *Adapter) inferAndConvert(value interface{}) types.Value {
    if value == nil {
        return types.NewNil()
    }
    
    switch v := value.(type) {
    case bool:
        return types.NewBool(v)
    case int, int8, int16, int32, int64:
        return types.NewInt(convertToInt64(v))
    case uint, uint8, uint16, uint32, uint64:
        return types.NewInt(int64(convertToUint64(v)))
    case float32, float64:
        return types.NewFloat(convertToFloat64(v))
    case string:
        return types.NewString(v)
    case []interface{}:
        return a.convertSlice(v)
    case map[string]interface{}:
        return a.convertMap(v)
    default:
        return a.convertStruct(v)
    }
}
```

### 2. 批量转换优化
```go
func (a *Adapter) ConvertBatch(values []interface{}) []types.Value {
    result := make([]types.Value, len(values))
    
    // 并行转换（适用于大批量数据）
    if len(values) > 1000 {
        return a.convertParallel(values)
    }
    
    // 串行转换
    for i, value := range values {
        result[i] = a.Convert(value)
    }
    
    return result
}

func (a *Adapter) convertParallel(values []interface{}) []types.Value {
    result := make([]types.Value, len(values))
    
    const batchSize = 100
    var wg sync.WaitGroup
    
    for i := 0; i < len(values); i += batchSize {
        end := i + batchSize
        if end > len(values) {
            end = len(values)
        }
        
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            for j := start; j < end; j++ {
                result[j] = a.Convert(values[j])
            }
        }(i, end)
    }
    
    wg.Wait()
    return result
}
```

## 性能优化

### 1. 字段访问缓存
```go
type FieldCache struct {
    cache map[string]*CachedField
    mutex sync.RWMutex
}

type CachedField struct {
    Accessor *FieldAccessor
    HitCount int64
    LastUsed time.Time
}

func (fc *FieldCache) GetAccessor(typeName, fieldName string) *FieldAccessor {
    key := typeName + "." + fieldName
    
    fc.mutex.RLock()
    cached, exists := fc.cache[key]
    fc.mutex.RUnlock()
    
    if exists {
        atomic.AddInt64(&cached.HitCount, 1)
        cached.LastUsed = time.Now()
        return cached.Accessor
    }
    
    // 创建新的访问器
    accessor := fc.createAccessor(typeName, fieldName)
    
    fc.mutex.Lock()
    fc.cache[key] = &CachedField{
        Accessor: accessor,
        HitCount: 1,
        LastUsed: time.Now(),
    }
    fc.mutex.Unlock()
    
    return accessor
}
```

### 2. 内存池管理
```go
type ValuePool struct {
    intPool    sync.Pool
    stringPool sync.Pool
    slicePool  sync.Pool
}

func (vp *ValuePool) GetInt(value int64) *types.IntValue {
    intVal := vp.intPool.Get().(*types.IntValue)
    intVal.SetValue(value)
    return intVal
}

func (vp *ValuePool) PutInt(intVal *types.IntValue) {
    intVal.Reset()
    vp.intPool.Put(intVal)
}
```

## 实际应用示例

### 1. 配置文件处理
```go
type Config struct {
    Database struct {
        Host     string `expr:"host"`
        Port     int    `expr:"port"`
        Username string `expr:"username"`
    } `expr:"database"`
    
    Features struct {
        EnableCache bool `expr:"enable_cache"`
        MaxUsers    int  `expr:"max_users"`
    } `expr:"features"`
}

func processConfig() {
    adapter := env.New()
    adapter.RegisterType(reflect.TypeOf(&Config{}))
    
    config := loadConfig()
    
    // 表达式: database.port > 0 && features.enable_cache
    // 可以直接访问嵌套字段
}
```

### 2. API响应处理
```go
type APIResponse struct {
    Data   interface{} `expr:"data"`
    Status string      `expr:"status"`
    Meta   struct {
        Total int `expr:"total"`
        Page  int `expr:"page"`
    } `expr:"meta"`
}

func handleAPIResponse() {
    adapter := env.New()
    
    response := &APIResponse{
        Status: "success",
        Meta: struct {
            Total int `expr:"total"`
            Page  int `expr:"page"`
        }{Total: 100, Page: 1},
    }
    
    // 表达式: status == "success" && meta.total > 0
    env := map[string]interface{}{
        "response": response,
    }
    
    result := evaluateExpression(`response.status == "success"`, env)
    fmt.Printf("Valid response: %v\n", result)
}
```

## 最佳实践

1. **类型注册**: 提前注册常用结构体类型
2. **缓存利用**: 重用适配器实例和字段访问器
3. **批量处理**: 大量数据时使用批量转换
4. **内存管理**: 使用对象池减少GC压力
5. **性能监控**: 监控字段访问热点进行优化

## 与其他模块的集成

Env模块为整个表达式引擎提供数据访问支持：
```
Go数据 → Env适配器 → types.Value → 表达式引擎
```

通过零反射的设计，Env模块确保了数据访问的高性能，是整个表达式引擎性能优势的重要基础。 