# API 文档

## 目录

1. [核心API](#核心api)
2. [配置选项](#配置选项)
3. [内置函数](#内置函数)
4. [模块系统](#模块系统)
5. [Lambda表达式](#lambda表达式)
6. [管道操作](#管道操作)
7. [空值安全](#空值安全)
8. [调试支持](#调试支持)
9. [错误处理](#错误处理)

## 核心API

### Eval

**函数签名**: `func Eval(expression string, env interface{}) (interface{}, error)`

一次性编译并执行表达式，适用于简单场景。

**参数**:
- `expression`: 表达式字符串
- `env`: 执行环境，可以是 `map[string]interface{}` 或结构体

**返回值**:
- `interface{}`: 表达式执行结果
- `error`: 执行错误

**示例**:
```go
// 简单算术
result, err := expr.Eval("2 + 3 * 4", nil)
// result: 14

// 使用环境变量
env := map[string]interface{}{"x": 10, "y": 20}
result, err = expr.Eval("x + y", env)
// result: 30

// 复杂表达式
env = map[string]interface{}{
    "users": []map[string]interface{}{
        {"name": "Alice", "age": 25},
        {"name": "Bob", "age": 30},
    },
}
result, err = expr.Eval("users | filter(u => u.age > 20) | map(u => u.name)", env)
// result: ["Alice", "Bob"]
```

### Compile

**函数签名**: `func Compile(expression string, options ...Option) (*Program, error)`

编译表达式为可执行程序，适用于需要多次执行的场景。

**参数**:
- `expression`: 表达式字符串  
- `options`: 编译选项

**返回值**:
- `*Program`: 编译后的程序对象
- `error`: 编译错误

**示例**:
```go
// 基础编译
program, err := expr.Compile("x * 2 + y")
if err != nil {
    log.Fatal(err)
}

// 多次执行
env1 := map[string]interface{}{"x": 5, "y": 3}
result1, _ := expr.Run(program, env1) // 13

env2 := map[string]interface{}{"x": 10, "y": 2}
result2, _ := expr.Run(program, env2) // 22
```

### Run

**函数签名**: `func Run(program *Program, env interface{}) (interface{}, error)`

执行编译后的程序。

**参数**:
- `program`: 编译后的程序对象
- `env`: 执行环境

**返回值**:
- `interface{}`: 执行结果
- `error`: 执行错误

**示例**:
```go
program, _ := expr.Compile("numbers | filter(# > 5) | sum()")

env1 := map[string]interface{}{"numbers": []int{1, 6, 3, 8}}
result1, _ := expr.Run(program, env1) // 14

env2 := map[string]interface{}{"numbers": []int{2, 7, 4, 9, 1}}
result2, _ := expr.Run(program, env2) // 16
```

## 配置选项

### Env

设置默认环境变量。

```go
program, err := expr.Compile("name + ' is ' + toString(age)",
    expr.Env(map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }))
```

### AsInt / AsFloat64 / AsString / AsBool

指定期望的返回类型，启用类型检查。

```go
// 期望整数结果
program, err := expr.Compile("1 + 2", expr.AsInt())

// 期望字符串结果  
program, err := expr.Compile("'hello' + ' world'", expr.AsString())

// 期望浮点数结果
program, err := expr.Compile("3.14 * 2", expr.AsFloat64())

// 期望布尔结果
program, err := expr.Compile("age > 18", expr.AsBool())
```

### WithTimeout

设置执行超时时间。

```go
program, err := expr.Compile("heavyComputation()",
    expr.WithTimeout(5 * time.Second))

result, err := expr.Run(program, env)
if err != nil {
    // 处理超时错误
}
```

### WithMaxIterations

设置最大迭代次数，防止无限循环。

```go
program, err := expr.Compile("numbers | map(# * 2)",
    expr.WithMaxIterations(10000))
```

## 内置函数

### 数学函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `abs(x)` | 绝对值 | `abs(-5)` → `5` |
| `min(a, b, ...)` | 最小值 | `min(3, 1, 4)` → `1` |
| `max(a, b, ...)` | 最大值 | `max(3, 1, 4)` → `4` |
| `sum(array)` | 数组求和 | `sum([1, 2, 3])` → `6` |
| `avg(array)` | 数组平均值 | `avg([1, 2, 3])` → `2` |
| `ceil(x)` | 向上取整 | `ceil(3.2)` → `4` |
| `floor(x)` | 向下取整 | `floor(3.8)` → `3` |
| `round(x)` | 四舍五入 | `round(3.5)` → `4` |

### 字符串函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `length(s)` | 字符串长度 | `length("hello")` → `5` |
| `upper(s)` | 转大写 | `upper("hello")` → `"HELLO"` |
| `lower(s)` | 转小写 | `lower("HELLO")` → `"hello"` |
| `trim(s)` | 去除首尾空格 | `trim(" hello ")` → `"hello"` |
| `contains(s, sub)` | 包含子串 | `contains("hello", "ell")` → `true` |
| `startsWith(s, prefix)` | 前缀匹配 | `startsWith("hello", "he")` → `true` |
| `endsWith(s, suffix)` | 后缀匹配 | `endsWith("hello", "lo")` → `true` |
| `indexOf(s, sub)` | 查找子串位置 | `indexOf("hello", "l")` → `2` |

### 数组函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `filter(array, predicate)` | 过滤数组 | `filter([1,2,3,4], # > 2)` → `[3,4]` |
| `map(array, transform)` | 映射数组 | `map([1,2,3], # * 2)` → `[2,4,6]` |
| `reduce(array, reducer)` | 数组归约 | `reduce([1,2,3,4], # + #)` → `10` |
| `sort(array)` | 数组排序 | `sort([3,1,4,2])` → `[1,2,3,4]` |
| `reverse(array)` | 数组反转 | `reverse([1,2,3])` → `[3,2,1]` |
| `unique(array)` | 数组去重 | `unique([1,2,2,3])` → `[1,2,3]` |
| `take(array, n)` | 取前n个 | `take([1,2,3,4,5], 3)` → `[1,2,3]` |
| `skip(array, n)` | 跳过前n个 | `skip([1,2,3,4,5], 2)` → `[3,4,5]` |

### 类型转换函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `toString(x)` | 转字符串 | `toString(123)` → `"123"` |
| `toNumber(s)` | 转数字 | `toNumber("123")` → `123` |
| `toBool(x)` | 转布尔值 | `toBool(1)` → `true` |
| `type(x)` | 获取类型 | `type(123)` → `"int"` |

### 工具函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `size(x)` | 获取大小 | `size([1,2,3])` → `3` |
| `keys(map)` | 获取map的键 | `keys({"a":1,"b":2})` → `["a","b"]` |
| `values(map)` | 获取map的值 | `values({"a":1,"b":2})` → `[1,2]` |
| `range(n)` | 生成数字范围 | `range(3)` → `[0,1,2]` |
| `first(array)` | 第一个元素 | `first([1,2,3])` → `1` |
| `last(array)` | 最后一个元素 | `last([1,2,3])` → `3` |

## 模块系统

### Math模块 (14个函数)

```go
// 基础数学函数
math.sqrt(16)     // 4
math.pow(2, 3)    // 8
math.abs(-5)      // 5

// 三角函数
math.sin(math.pi / 2)  // 1
math.cos(0)            // 1
math.tan(math.pi / 4)  // 1

// 对数函数
math.log(math.e)   // 1
math.log10(100)    // 2
math.log2(8)       // 3

// 取整函数
math.ceil(3.2)     // 4
math.floor(3.8)    // 3
math.round(3.5)    // 4

// 常数
math.pi      // 3.141592653589793
math.e       // 2.718281828459045
```

### Strings模块 (13个函数)

```go
// 大小写转换
strings.upper("hello")      // "HELLO"
strings.lower("HELLO")      // "hello"
strings.title("hello world") // "Hello World"

// 字符串处理
strings.trim(" hello ")     // "hello"
strings.trimLeft(" hello")  // "hello"
strings.trimRight("hello ") // "hello"

// 字符串操作
strings.replace("hello", "l", "L")  // "heLLo"
strings.split("a,b,c", ",")        // ["a","b","c"]
strings.join(["a","b","c"], ",")   // "a,b,c"

// 字符串查询
strings.contains("hello", "ell")    // true
strings.hasPrefix("hello", "he")    // true
strings.hasSuffix("hello", "lo")    // true
strings.indexOf("hello", "l")       // 2
```

### 自定义模块

```go
// 注册自定义模块
modules.RegisterModule("utils", map[string]interface{}{
    "formatCurrency": func(amount float64) string {
        return fmt.Sprintf("$%.2f", amount)
    },
    "isEven": func(n int) bool {
        return n%2 == 0
    },
})

// 使用自定义模块
result, _ := expr.Eval("utils.formatCurrency(123.456)", nil)
// result: "$123.46"

result, _ = expr.Eval("utils.isEven(4)", nil)
// result: true
```

## Lambda表达式

### 基础语法

```go
// 单参数Lambda
filter(users, user => user.age > 18)
map(numbers, n => n * 2)

// 多参数Lambda
reduce(numbers, (acc, n) => acc + n)
sort(items, (a, b) => a.priority - b.priority)

// 复杂Lambda表达式
map(products, p => {
    name: p.name,
    discountPrice: p.price * 0.8
})
```

### 在管道中使用

```go
users 
| filter(u => u.active && u.age >= 18)
| map(u => {name: u.name, email: u.email})
| sort((a, b) => a.name.localeCompare(b.name))
```

## 管道操作

### 占位符语法

占位符 `#` 代表管道中的当前元素。

```go
// 基础操作
numbers | filter(# > 5)        // 过滤大于5的数
numbers | map(# * 2)           // 每个数乘以2

// 复杂表达式
numbers | filter(# % 2 == 0)   // 过滤偶数
numbers | map(# * # + 1)       // 平方加1

// 条件表达式
numbers | map(# > 10 ? # : # * 2)  // 条件映射

// 链式操作
data 
| filter(# > threshold)
| map(# * multiplier)
| sort()
| take(10)
```

### Lambda与占位符混合

```go
// 在filter中使用Lambda，在map中使用占位符
users 
| filter(u => u.role == "admin")
| map(#.name)

// 复杂的混合使用
products
| filter(p => p.category == "electronics")
| map(# => {name: #.name, salePrice: #.price * 0.9})
| sort((a, b) => a.salePrice - b.salePrice)
```

## 空值安全

### 可选链操作符 `?.`

安全访问可能为null的对象属性。

```go
// 基础用法
user?.profile?.name        // 如果user或profile为null，返回null
user?.addresses?.[0]?.city // 安全访问数组元素

// 在表达式中使用
user?.profile?.name ?? "Unknown User"
order?.items?.length ?? 0
```

### 空值合并操作符 `??`

当左侧为null或undefined时，返回右侧的值。

```go
// 基础用法
user.name ?? "Guest"           // 如果name为null，返回"Guest"
settings.timeout ?? 5000       // 默认超时时间

// 链式使用
user?.profile?.bio ?? user?.description ?? "No description"

// 与其他操作符结合
(user?.age ?? 0) > 18 ? "Adult" : "Minor"
```

### 在管道中使用

```go
users 
| map(#?.profile?.avatar ?? "/default-avatar.png")
| filter(# != "/default-avatar.png")
```

## 调试支持

### 创建调试器

```go
import "github.com/mredencom/expr/debug"

// 创建新的调试器
debugger := debug.NewDebugger()
```

### 断点管理

```go
// 设置断点
debugger.SetBreakpoint(5)        // 在字节码位置5设置断点
debugger.SetBreakpoint(10)       // 在字节码位置10设置断点

// 移除断点
debugger.RemoveBreakpoint(5)

// 清除所有断点
debugger.ClearBreakpoints()

// 检查断点
hasBreakpoint := debugger.HasBreakpoint(5)
```

### 单步执行

```go
// 单步执行程序
result := debugger.StepThrough(program, env)

// 获取执行统计
stats := debugger.GetExecutionStats()
fmt.Printf("执行步数: %d\n", stats.Steps)
fmt.Printf("断点命中: %d\n", stats.BreakpointHits)
fmt.Printf("执行时间: %v\n", stats.ExecutionTime)
```

### 执行回调

```go
// 设置执行回调
debugger.SetExecutionCallback(func(step int, opcode string, value interface{}) {
    fmt.Printf("Step %d: %s -> %v\n", step, opcode, value)
})

// 设置断点回调
debugger.SetBreakpointCallback(func(step int) {
    fmt.Printf("断点命中于步骤 %d\n", step)
})
```

## 错误处理

### 编译错误

```go
program, err := expr.Compile("invalid expression +")
if err != nil {
    if compileErr, ok := err.(*expr.CompileError); ok {
        fmt.Printf("编译错误: %s\n", compileErr.Message)
        fmt.Printf("位置: 行%d 列%d\n", compileErr.Line, compileErr.Column)
    }
}
```

### 运行时错误

```go
result, err := expr.Run(program, env)
if err != nil {
    if runtimeErr, ok := err.(*expr.RuntimeError); ok {
        fmt.Printf("运行时错误: %s\n", runtimeErr.Message)
        fmt.Printf("表达式: %s\n", runtimeErr.Expression)
    }
}
```

### 超时错误

```go
program, _ := expr.Compile("longRunningOperation()",
    expr.WithTimeout(1 * time.Second))

result, err := expr.Run(program, env)
if err != nil {
    if timeoutErr, ok := err.(*expr.TimeoutError); ok {
        fmt.Printf("执行超时: %v\n", timeoutErr.Duration)
    }
}
```

### 类型错误

```go
program, err := expr.Compile("name + age", 
    expr.AsString(),
    expr.Env(map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }))

if err != nil {
    if typeErr, ok := err.(*expr.TypeError); ok {
        fmt.Printf("类型错误: %s\n", typeErr.Message)
        fmt.Printf("期望类型: %s\n", typeErr.Expected)
        fmt.Printf("实际类型: %s\n", typeErr.Actual)
    }
}
```

## 性能优化建议

### 预编译表达式

```go
// ❌ 避免在循环中重复编译
for _, data := range datasets {
    result, _ := expr.Eval("complex expression", data)
}

// ✅ 预编译表达式
program, _ := expr.Compile("complex expression")
for _, data := range datasets {
    result, _ := expr.Run(program, data)
}
```

### 类型提示

```go
// ✅ 提供类型提示以获得更好的性能
program, _ := expr.Compile("x + y", 
    expr.AsFloat64(),
    expr.Env(map[string]interface{}{
        "x": 0.0,
        "y": 0.0,
    }))
```

### 缓存启用

```go
// ✅ 启用内置缓存
program, _ := expr.Compile("expression", expr.EnableCache())
``` 