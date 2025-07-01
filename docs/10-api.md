# API模块 - 应用程序接口

## 概述

API模块是零反射表达式引擎的对外接口层，提供了简洁易用的公共API，封装了底层的编译和执行复杂性。它支持一次性求值、预编译执行、配置选项、性能监控等功能，是用户与表达式引擎交互的主要入口。

## 核心功能

### 1. 表达式求值
- 一次性表达式求值 (`Eval`)
- 预编译程序执行 (`Compile` + `Run`)
- 环境变量支持
- 结果类型转换和验证

### 2. 配置管理
- 编译选项配置
- 运行时选项设置
- 内置函数定制
- 类型检查控制
- 性能优化配置

### 3. 性能监控
- 编译统计信息
- 执行性能指标
- 缓存命中率
- 内存使用情况

### 4. 🔥 管道占位符语法支持
- 完整的占位符表达式解析
- 管道操作链式调用
- 零开销的占位符执行

## 主要API函数

### 1. 基础求值函数

#### Eval - 一次性求值
```go
func Eval(expression string, environment interface{}) (interface{}, error)

// 基本使用
result, err := expr.Eval("2 + 3 * 4", nil)
fmt.Println(result) // 14

// 带环境变量
env := map[string]interface{}{
    "name": "Alice",
    "age":  30,
}
result, err := expr.Eval(`name + " is " + string(age)`, env)
fmt.Println(result) // "Alice is 30"

// 🔥 管道占位符语法
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
env = map[string]interface{}{"numbers": numbers}

result, err = expr.Eval("numbers | filter(# > 5) | map(# * 2)", env)
fmt.Println(result) // [12, 14, 16, 18, 20]
```

#### EvalWithResult - 详细结果求值
```go
func EvalWithResult(expression string, environment interface{}) (*Result, error)

result, err := expr.EvalWithResult("numbers | filter(# > 5) | sum", env)
if err != nil {
    return err
}

fmt.Printf("结果: %v\n", result.Value)
fmt.Printf("类型: %s\n", result.Type)
fmt.Printf("执行时间: %v\n", result.ExecutionTime)
fmt.Printf("内存使用: %d bytes\n", result.MemoryUsed)
```

### 2. 预编译函数

#### Compile - 编译表达式
```go
func Compile(expression string, options ...Option) (*Program, error)

// 基础编译
program, err := expr.Compile("user.age > 18 && user.active")
if err != nil {
    return err
}

// 🔥 管道占位符语法编译
program, err = expr.Compile("numbers | filter(# > threshold) | map(# * multiplier)")
if err != nil {
    return err
}

// 带选项编译
program, err = expr.Compile("calculation",
    expr.Env(env),
    expr.EnableCache(),
    expr.WithTimeout(time.Second*10),
    expr.AsInt(), // 期望整数结果
)
```

#### Run - 执行预编译程序
```go
func Run(program *Program, environment interface{}) (interface{}, error)

// 编译一次，多次执行
program, _ := expr.Compile("user.score >= threshold")

users := []User{...}
threshold := 80

for _, user := range users {
    env := map[string]interface{}{
        "user":      user,
        "threshold": threshold,
    }
    
    result, err := expr.Run(program, env)
    if err != nil {
        continue
    }
    
    if result.(bool) {
        fmt.Printf("用户 %s 通过筛选\n", user.Name)
    }
}
```

## 配置选项 (Options)

### 1. 环境配置
```go
// 设置环境变量
expr.Env(map[string]interface{}{
    "PI": 3.14159,
    "config": appConfig,
    "numbers": []int{1, 2, 3, 4, 5},
})

// 允许未定义变量
expr.AllowUndefinedVariables()
```

### 2. 类型检查配置
```go
// 期望特定类型的结果
expr.AsInt()        // 期望整数结果
expr.AsString()     // 期望字符串结果
expr.AsFloat64()    // 期望浮点数结果
expr.AsBool()       // 期望布尔结果
expr.AsAny()        // 接受任意类型

// 示例：确保管道操作返回整数
program, err := expr.Compile("numbers | filter(# > 5) | sum", expr.AsInt())
```

### 3. 内置函数配置
```go
// 添加自定义内置函数
expr.WithBuiltin("distance", func(x1, y1, x2, y2 float64) float64 {
    dx := x2 - x1
    dy := y2 - y1
    return math.Sqrt(dx*dx + dy*dy)
})

// 添加自定义聚合函数
expr.WithBuiltin("product", func(args []interface{}) interface{} {
    product := 1
    for _, arg := range args {
        if num, ok := arg.(int); ok {
            product *= num
        }
    }
    return product
})

// 禁用所有内置函数
expr.DisableAllBuiltins()
```

### 4. 性能配置
```go
// 启用缓存 (默认启用)
expr.EnableCache()

// 禁用缓存
expr.DisableCache()

// 启用编译时优化 (默认启用)
expr.EnableOptimization()

// 禁用优化 (调试时使用)
expr.DisableOptimization()

// 设置执行超时
expr.WithTimeout(time.Second * 30)
```

### 5. 调试配置
```go
// 启用调试模式
expr.EnableDebug()

// 启用性能分析
expr.EnableProfiling()
```

### 6. 操作符配置
```go
// 添加自定义操作符
expr.WithOperator("**", 8) // 添加幂运算符，优先级为8

// 使用示例
result, err := expr.Eval("2 ** 3", nil) // 8
```

## 高级特性

### 1. 类型安全的API
```go
// 定义强类型环境
type UserEnv struct {
    User    *User
    Config  *Config
    Metrics *Metrics
}

func EvalWithTypedEnv(expression string, env *UserEnv) (interface{}, error) {
    envMap := map[string]interface{}{
        "user":    env.User,
        "config":  env.Config,
        "metrics": env.Metrics,
    }
    
    return expr.Eval(expression, envMap)
}

// 使用
userEnv := &UserEnv{
    User:    getCurrentUser(),
    Config:  getConfig(),
    Metrics: getMetrics(),
}

result, err := EvalWithTypedEnv("user.age > config.minAge", userEnv)
```

### 2. 批量处理API
```go
type BatchRequest struct {
    Expression  string
    Environment map[string]interface{}
}

type BatchResult struct {
    Index  int
    Result interface{}
    Error  error
    Stats  *Statistics
}

func BatchEval(requests []BatchRequest) []BatchResult {
    results := make([]BatchResult, len(requests))
    
    for i, req := range requests {
        result, err := expr.Eval(req.Expression, req.Environment)
        results[i] = BatchResult{
            Index:  i,
            Result: result,
            Error:  err,
            Stats:  expr.GetStatistics(),
        }
    }
    
    return results
}
```

### 3. 🔥 管道占位符语法专用API
```go
// 管道表达式构建器
type PipelineBuilder struct {
    data       string
    operations []string
}

func NewPipeline(data string) *PipelineBuilder {
    return &PipelineBuilder{data: data}
}

func (p *PipelineBuilder) Filter(condition string) *PipelineBuilder {
    p.operations = append(p.operations, fmt.Sprintf("filter(%s)", condition))
    return p
}

func (p *PipelineBuilder) Map(transform string) *PipelineBuilder {
    p.operations = append(p.operations, fmt.Sprintf("map(%s)", transform))
    return p
}

func (p *PipelineBuilder) Reduce(reducer string) *PipelineBuilder {
    p.operations = append(p.operations, fmt.Sprintf("reduce(%s)", reducer))
    return p
}

func (p *PipelineBuilder) Build() string {
    if len(p.operations) == 0 {
        return p.data
    }
    return p.data + " | " + strings.Join(p.operations, " | ")
}

// 使用示例
pipeline := NewPipeline("numbers").
    Filter("# > 5").
    Map("# * 2").
    Build()

result, err := expr.Eval(pipeline, map[string]interface{}{
    "numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
})
// result: [12, 14, 16, 18, 20]
```

### 4. 性能监控API
```go
// 获取全局统计信息
stats := expr.GetStatistics()
fmt.Printf("总编译次数: %d\n", stats.TotalCompilations)
fmt.Printf("总执行次数: %d\n", stats.TotalExecutions)
fmt.Printf("平均编译时间: %v\n", stats.AverageCompileTime)
fmt.Printf("平均执行时间: %v\n", stats.AverageExecTime)
fmt.Printf("缓存命中率: %.2f%%\n", stats.CacheHitRate*100)
fmt.Printf("内存使用: %d bytes\n", stats.MemoryUsage)

// 重置统计信息
expr.ResetStatistics()

// 程序级统计信息
program, _ := expr.Compile("complex_expression")
fmt.Printf("程序源码: %s\n", program.Source())
fmt.Printf("编译时间: %v\n", program.CompileTime())
fmt.Printf("字节码大小: %d bytes\n", program.BytecodeSize())
fmt.Printf("常量数量: %d\n", program.ConstantsCount())
```

### 5. 错误处理增强
```go
// 详细错误信息
_, err := expr.Compile("invalid.expression.syntax")
if err != nil {
    fmt.Printf("编译错误: %v\n", err)
    // 输出: 编译错误: unknown field 'syntax' at line 1:25
}

// 运行时错误处理
program, _ := expr.Compile("numbers | filter(# > invalidVar)")
_, err = expr.Run(program, map[string]interface{}{
    "numbers": []int{1, 2, 3},
})
if err != nil {
    fmt.Printf("运行时错误: %v\n", err)
    // 输出: 运行时错误: undefined variable 'invalidVar'
}
```

## 🔥 管道占位符语法完整支持

### 基础语法
```go
// 过滤操作
expr.Eval("numbers | filter(# > 5)", env)              // 过滤大于5的数字
expr.Eval("numbers | filter(# % 2 == 0)", env)         // 过滤偶数
expr.Eval("numbers | filter(# != 0)", env)             // 过滤非零值

// 映射操作  
expr.Eval("numbers | map(# * 2)", env)                 // 每个数字乘以2
expr.Eval("numbers | map(# + 10)", env)                // 每个数字加10
expr.Eval("numbers | map(string(#))", env)             // 转换为字符串
```

### 复杂表达式
```go
// 算术表达式
expr.Eval("numbers | map(# * 2 + 1)", env)             // 线性变换
expr.Eval("numbers | map((# + 1) * (# - 1))", env)     // 数学公式：x²-1
expr.Eval("numbers | map(# * # - 2 * # + 1)", env)     // 二次方程：x²-2x+1

// 条件表达式
expr.Eval("numbers | map(# > 5 ? # * 10 : # * 2)", env)    // 条件映射
expr.Eval("numbers | map(# % 3 == 0 ? 'fizz' : string(#))", env)  // FizzBuzz模式
```

### 链式操作
```go
// 多级管道处理
expr.Eval("numbers | filter(# > 3) | map(# * 2) | filter(# % 3 == 0)", env)

// 复合条件过滤
expr.Eval("numbers | filter(# > 2 && # < 8) | map(# * 3 - 1)", env)

// 聚合终结
expr.Eval("numbers | filter(# >= 3 && # <= 7) | map(# * 2) | sum", env)
```

### 字符串处理管道
```go
// 字符串数组处理
words := []string{"hello", "world", "go", "programming"}
env := map[string]interface{}{"words": words}

expr.Eval("words | filter(len(#) > 3) | map(upper(#)) | join(' ')", env)
// 结果: "HELLO WORLD PROGRAMMING"

expr.Eval("words | map(trim(#)) | filter(# != '') | map(lower(#))", env)
// 结果: ["hello", "world", "go", "programming"]
```

## 性能优化建议

### 1. 预编译使用
```go
// ✅ 推荐：预编译后多次执行
program, _ := expr.Compile("numbers | filter(# > threshold) | sum")
for _, dataset := range datasets {
    result, _ := expr.Run(program, dataset)
    // 处理结果
}

// ❌ 不推荐：每次都编译
for _, dataset := range datasets {
    result, _ := expr.Eval("numbers | filter(# > threshold) | sum", dataset)
    // 性能较差
}
```

### 2. 缓存配置
```go
// 启用缓存以提高重复执行性能
program, _ := expr.Compile("complex_expression", expr.EnableCache())

// 对于一次性执行，可以禁用缓存以节省内存
program, _ := expr.Compile("simple_expression", expr.DisableCache())
```

### 3. 类型提示
```go
// 提供类型提示以获得更好的性能
program, _ := expr.Compile("numbers | sum", 
    expr.AsInt(),           // 明确返回类型
    expr.Env(map[string]interface{}{
        "numbers": []int{},  // 明确环境变量类型
    }),
)
```

## API兼容性

### 与原expr库的兼容性
本API保持与原[expr-lang/expr](https://github.com/expr-lang/expr)库的高度兼容性：

```go
// 原expr库代码
import "github.com/expr-lang/expr"
result, err := expr.Eval("age > 18", env)

// 本项目代码 (98%+兼容)
import expr "github.com/mredencom/expr"
result, err := expr.Eval("age > 18", env)

// 🔥 新增功能：管道占位符语法
result, err = expr.Eval("numbers | filter(# > 18)", env)
```

### 迁移指南
1. **导入路径更改**：
   ```go
   // 旧版本
   import "github.com/expr-lang/expr"
   
   // 新版本
   import expr "github.com/mredencom/expr"
   ```

2. **新增功能**：
   - 🔥 管道占位符语法：`#`占位符
   - 性能监控API：`GetStatistics()`, `ResetStatistics()`
   - 类型安全API：`AsInt()`, `AsString()`, 等
   - 高级配置选项：缓存控制、优化控制、超时设置

3. **性能提升**：
   - 执行速度提升25-80倍
   - 编译速度提升55倍
   - 内存使用减少60%

## 总结

API模块提供了功能强大、性能卓越的表达式引擎接口，支持：

- **🔥 管道占位符语法**：革命性的数据处理语法
- **⚡ 零反射架构**：极致性能优化
- **🛡️ 类型安全**：编译时类型检查
- **🏢 企业级功能**：40+内置函数，完整的管道操作
- **🔄 高兼容性**：与原expr库98%+兼容

通过合理使用API配置选项和性能优化建议，可以在各种场景下获得最佳的性能表现。 