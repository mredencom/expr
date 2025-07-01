# 调试指南

## 🐛 调试器基础

Expr提供了专业的调试器支持，帮助您深入了解表达式的执行过程。

### 创建调试器

```go
import "github.com/mredencom/expr/debug"

// 创建新的调试器实例
debugger := debug.NewDebugger()
```

## 🔍 断点管理

### 设置断点

```go
// 在特定字节码位置设置断点
debugger.SetBreakpoint(5)
debugger.SetBreakpoint(10)
debugger.SetBreakpoint(15)

// 检查是否设置了断点
hasBreakpoint := debugger.HasBreakpoint(5) // true
```

### 管理断点

```go
// 移除特定断点
debugger.RemoveBreakpoint(5)

// 清除所有断点
debugger.ClearBreakpoints()

// 获取所有断点
breakpoints := debugger.GetBreakpoints()
fmt.Printf("当前断点: %v\n", breakpoints)
```

## 🚶 单步执行

### 基础单步执行

```go
// 编译表达式
program, err := expr.Compile("numbers | filter(# > 5) | map(# * 2)")
if err != nil {
    log.Fatal(err)
}

// 准备数据
env := map[string]interface{}{
    "numbers": []int{1, 6, 3, 8, 2, 9},
}

// 单步执行
result := debugger.StepThrough(program, env)
fmt.Printf("最终结果: %v\n", result)
```

### 获取执行统计

```go
// 执行后获取统计信息
stats := debugger.GetExecutionStats()
fmt.Printf("执行步数: %d\n", stats.Steps)
fmt.Printf("断点命中次数: %d\n", stats.BreakpointHits)
fmt.Printf("执行时间: %v\n", stats.ExecutionTime)
fmt.Printf("访问的变量: %v\n", stats.VariablesAccessed)
```

## 📊 执行回调

### 步骤回调

```go
// 设置执行步骤回调
debugger.SetExecutionCallback(func(step int, opcode string, value interface{}) {
    fmt.Printf("步骤 %d: %s -> %v\n", step, opcode, value)
})

// 执行表达式，查看每个步骤
debugger.StepThrough(program, env)
```

### 断点回调

```go
// 设置断点命中回调
debugger.SetBreakpointCallback(func(step int) {
    fmt.Printf("🔴 断点命中于步骤 %d\n", step)
    
    // 可以在这里检查当前状态
    stack := debugger.GetCurrentStack()
    fmt.Printf("当前栈状态: %v\n", stack)
})
```

## 🔧 高级调试功能

### 变量监控

```go
// 监控特定变量的访问
debugger.WatchVariable("user")
debugger.WatchVariable("settings")

// 设置变量访问回调
debugger.SetVariableAccessCallback(func(name string, value interface{}) {
    fmt.Printf("📍 访问变量 %s: %v\n", name, value)
})
```

### 条件断点

```go
// 设置条件断点（仅在满足条件时暂停）
debugger.SetConditionalBreakpoint(8, func(stack []interface{}) bool {
    // 仅当栈顶值大于10时暂停
    if len(stack) > 0 {
        if val, ok := stack[len(stack)-1].(int); ok {
            return val > 10
        }
    }
    return false
})
```

## 🎯 实际调试示例

### 调试复杂表达式

```go
func debugComplexExpression() {
    debugger := debug.NewDebugger()
    
    // 复杂的业务表达式
    expression := `
        users 
        | filter(u => u.active && u.age >= minAge)
        | map(u => {
            name: u.firstName + " " + u.lastName,
            score: u.baseScore * multiplier + bonus
        })
        | filter(u => u.score > threshold)
        | sort((a, b) => b.score - a.score)
        | take(topN)
    `
    
    program, err := expr.Compile(expression)
    if err != nil {
        log.Fatal("编译失败:", err)
    }
    
    // 设置调试回调
    debugger.SetExecutionCallback(func(step int, opcode string, value interface{}) {
        fmt.Printf("[%03d] %-15s %v\n", step, opcode, value)
    })
    
    // 在关键操作上设置断点
    debugger.SetBreakpoint(20) // filter操作后
    debugger.SetBreakpoint(35) // map操作后
    debugger.SetBreakpoint(50) // sort操作后
    
    // 执行并调试
    env := map[string]interface{}{
        "users": []map[string]interface{}{
            {"firstName": "Alice", "lastName": "Smith", "active": true, "age": 25, "baseScore": 80},
            {"firstName": "Bob", "lastName": "Jones", "active": false, "age": 30, "baseScore": 90},
            {"firstName": "Charlie", "lastName": "Brown", "active": true, "age": 35, "baseScore": 85},
        },
        "minAge":     20,
        "multiplier": 1.2,
        "bonus":      10,
        "threshold":  100,
        "topN":       2,
    }
    
    result := debugger.StepThrough(program, env)
    fmt.Printf("\n最终结果: %+v\n", result)
    
    // 查看执行统计
    stats := debugger.GetExecutionStats()
    fmt.Printf("\n=== 执行统计 ===\n")
    fmt.Printf("总步数: %d\n", stats.Steps)
    fmt.Printf("执行时间: %v\n", stats.ExecutionTime)
    fmt.Printf("断点命中: %d\n", stats.BreakpointHits)
}
```

### 调试Lambda表达式

```go
func debugLambdaExpression() {
    debugger := debug.NewDebugger()
    
    // Lambda表达式调试
    expression := "numbers | filter(n => n > threshold) | map(n => n * multiplier)"
    program, _ := expr.Compile(expression)
    
    // 监控Lambda变量
    debugger.WatchVariable("n")
    debugger.WatchVariable("threshold")
    debugger.WatchVariable("multiplier")
    
    debugger.SetVariableAccessCallback(func(name string, value interface{}) {
        fmt.Printf("🔍 Lambda变量 %s = %v\n", name, value)
    })
    
    env := map[string]interface{}{
        "numbers":    []int{1, 5, 3, 8, 2, 9},
        "threshold":  4,
        "multiplier": 3,
    }
    
    result := debugger.StepThrough(program, env)
    fmt.Printf("Lambda结果: %v\n", result)
}
```

## 🔍 错误诊断

### 表达式错误诊断

```go
func diagnoseExpression(expression string, env interface{}) {
    fmt.Printf("=== 表达式诊断 ===\n")
    fmt.Printf("表达式: %s\n", expression)
    
    // 1. 编译检查
    program, err := expr.Compile(expression)
    if err != nil {
        fmt.Printf("❌ 编译失败: %v\n", err)
        
        // 详细错误分析
        if compileErr, ok := err.(*expr.CompileError); ok {
            fmt.Printf("错误位置: 行%d 列%d\n", compileErr.Line, compileErr.Column)
            fmt.Printf("错误类型: %s\n", compileErr.Type)
        }
        return
    }
    fmt.Printf("✅ 编译成功\n")
    
    // 2. 执行检查
    debugger := debug.NewDebugger()
    
    // 捕获执行错误
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("❌ 执行panic: %v\n", r)
        }
    }()
    
    result := debugger.StepThrough(program, env)
    
    if result == nil {
        fmt.Printf("⚠️ 执行返回nil\n")
    } else {
        fmt.Printf("✅ 执行成功\n")
        fmt.Printf("结果: %v (类型: %T)\n", result, result)
    }
    
    // 3. 性能分析
    stats := debugger.GetExecutionStats()
    fmt.Printf("\n=== 性能分析 ===\n")
    fmt.Printf("执行步数: %d\n", stats.Steps)
    fmt.Printf("执行时间: %v\n", stats.ExecutionTime)
    
    if stats.ExecutionTime > 10*time.Millisecond {
        fmt.Printf("⚠️ 执行时间较长，可能需要优化\n")
    }
}
```

### 常见问题诊断

```go
func commonIssuesDiagnosis() {
    fmt.Println("=== 常见问题诊断 ===")
    
    // 1. 类型错误
    fmt.Println("\n1. 类型错误检查:")
    diagnoseExpression("name + age", map[string]interface{}{
        "name": "Alice",
        "age":  30,
    })
    
    // 2. 变量不存在
    fmt.Println("\n2. 变量不存在检查:")
    diagnoseExpression("unknownVar + 10", map[string]interface{}{
        "knownVar": 5,
    })
    
    // 3. 函数调用错误
    fmt.Println("\n3. 函数调用错误:")
    diagnoseExpression("unknownFunction(42)", nil)
    
    // 4. 数组越界
    fmt.Println("\n4. 数组访问检查:")
    diagnoseExpression("arr[10]", map[string]interface{}{
        "arr": []int{1, 2, 3},
    })
    
    // 5. 空值访问
    fmt.Println("\n5. 空值访问检查:")
    diagnoseExpression("user.profile.name", map[string]interface{}{
        "user": map[string]interface{}{
            "profile": nil,
        },
    })
}
```

## 📝 调试日志

### 启用详细日志

```go
// 创建带详细日志的调试器
debugger := debug.NewDebugger()
debugger.EnableVerboseLogging(true)

// 设置日志输出
debugger.SetLogOutput(os.Stdout)

// 执行时会输出详细的调试信息
result := debugger.StepThrough(program, env)
```

### 自定义日志格式

```go
// 自定义日志记录器
type CustomLogger struct {
    file *os.File
}

func (cl *CustomLogger) Log(level string, message string, args ...interface{}) {
    timestamp := time.Now().Format("2006-01-02 15:04:05.000")
    fmt.Fprintf(cl.file, "[%s] %s: %s\n", timestamp, level, fmt.Sprintf(message, args...))
}

// 使用自定义日志记录器
logger := &CustomLogger{file: logFile}
debugger.SetLogger(logger)
```

## 🎯 调试最佳实践

### 1. 分步骤调试

```go
// 复杂表达式分解调试
expressions := []string{
    "users | filter(u => u.active)",
    "users | filter(u => u.active) | map(u => u.score)",
    "users | filter(u => u.active) | map(u => u.score) | sort()",
}

for i, expr := range expressions {
    fmt.Printf("=== 步骤 %d ===\n", i+1)
    diagnoseExpression(expr, env)
}
```

### 2. 使用测试数据

```go
// 创建简化的测试数据
testEnv := map[string]interface{}{
    "users": []map[string]interface{}{
        {"name": "Test1", "active": true, "score": 100},
        {"name": "Test2", "active": false, "score": 200},
    },
    "threshold": 150,
}

debugger.StepThrough(program, testEnv)
```

### 3. 性能调试

```go
// 比较不同实现的性能
expressions := []string{
    "users | filter(u => u.score > 100)",  // Lambda版本
    "users | filter(#.score > 100)",       // 占位符版本
}

for _, expr := range expressions {
    debugger := debug.NewDebugger()
    program, _ := expr.Compile(expr)
    
    start := time.Now()
    debugger.StepThrough(program, env)
    duration := time.Since(start)
    
    stats := debugger.GetExecutionStats()
    fmt.Printf("表达式: %s\n", expr)
    fmt.Printf("执行时间: %v\n", duration)
    fmt.Printf("执行步数: %d\n", stats.Steps)
    fmt.Println()
}
```

这些调试功能和技巧将帮助您快速定位和解决表达式执行中的问题。 