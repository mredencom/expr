# Debug 模块 - 调试与性能分析

## 概述

Debug模块为表达式引擎提供了强大的调试和性能分析功能。通过断点设置、执行统计、步进调试和热点分析等功能，帮助开发者深入了解表达式的执行过程，优化性能，快速定位问题。

## 🏗️ 核心架构

### 主要组件关系
```
Debugger (调试器)
├── Breakpoint[] (断点集合)
├── ExecutionStats (执行统计)
├── DebugContext (调试上下文) 
└── Event Callbacks (事件回调)
```

## 📊 核心组件

### 1. Debugger - 调试器

调试器是整个调试系统的核心，负责管理断点、统计信息和执行控制。

```go
type Debugger struct {
    breakpoints map[int]*Breakpoint
    stats       *ExecutionStats
    enabled     bool
    stepMode    bool
    
    // 执行状态
    currentPC        int
    currentStack     []types.Value
    instructionCount int64
    
    // 事件回调
    onBreakpoint func(*DebugContext)
    onStep       func(*DebugContext)
    onError      func(error)
}
```

#### 基本操作
```go
// 创建调试器
debugger := debug.New()

// 启用/禁用调试
debugger.Enable()
debugger.Disable()

// 检查状态
isEnabled := debugger.IsEnabled()

// 设置步进模式
debugger.SetStepMode(true)
```

### 2. DebugContext - 调试上下文

提供当前执行状态的详细信息，在断点或步进时传递给回调函数。

```go
type DebugContext struct {
    PC          int                    // 程序计数器
    Instruction []byte                 // 当前指令字节码
    Stack       []types.Value          // 当前栈状态
    Variables   map[string]types.Value // 当前变量值
    Source      string                 // 原始源代码
    Position    lexer.Position         // 源码位置信息
}
```

### 3. ExecutionStats - 执行统计

收集和分析表达式执行过程中的各项性能指标。

```go
type ExecutionStats struct {
    TotalInstructions int64                 // 总指令数
    InstructionCounts map[vm.Opcode]int64   // 各指令执行次数
    ExecutionTime     time.Duration         // 总执行时间
    StartTime         time.Time             // 开始时间
    FunctionCalls     int64                 // 函数调用次数
    MemoryAllocations int64                 // 内存分配次数
    HotSpots          []HotSpot             // 执行热点
}

type HotSpot struct {
    PC         int            // 程序计数器位置
    OpCode     vm.Opcode      // 操作码
    Count      int64          // 执行次数
    Percentage float64        // 占比
    Source     string         // 源码
    Position   lexer.Position // 源码位置
}
```

### 4. Breakpoint - 断点

支持条件断点、计数断点等高级调试功能。

```go
type Breakpoint struct {
    PC          int            // 程序计数器位置
    Enabled     bool           // 是否启用
    HitCount    int64          // 命中次数
    Condition   string         // 条件表达式
    Description string         // 描述信息
    CreatedAt   time.Time      // 创建时间
}
```

## 🔧 调试功能

### 1. 断点管理

#### 基本断点操作
```go
// 设置断点
bp := debugger.SetBreakpoint(10) // 在PC=10处设置断点
bp.SetDescription("主循环入口")

// 移除断点
success := debugger.RemoveBreakpoint(10)

// 获取断点
bp, exists := debugger.GetBreakpoint(10)

// 列出所有断点
breakpoints := debugger.ListBreakpoints()
for _, bp := range breakpoints {
    fmt.Printf("断点: %s\n", bp.String())
}
```

#### 条件断点
```go
// 设置条件断点
bp := debugger.SetBreakpoint(15)
bp.SetCondition("x > 10")
bp.SetDescription("当x大于10时停止")

// 启用/禁用断点
bp.Enable()
bp.Disable()

// 检查断点状态
shouldBreak := bp.ShouldBreak()
```

#### 断点信息
```go
// 断点字符串表示
fmt.Println(bp.String())
// 输出: "Breakpoint 15: enabled (hits: 3) - 当x大于10时停止 [condition: x > 10]"

// 记录断点命中
bp.Hit()
fmt.Printf("命中次数: %d\n", bp.HitCount)
```

### 2. 执行控制

#### 步进调试
```go
// 启用步进模式
debugger.SetStepMode(true)

// 设置步进回调
debugger.OnStep(func(ctx *DebugContext) {
    fmt.Printf("步进执行: PC=%d, 指令=%v\n", 
        ctx.PC, ctx.Instruction)
    
    // 检查栈状态
    fmt.Printf("栈深度: %d\n", len(ctx.Stack))
    for i, val := range ctx.Stack {
        fmt.Printf("  [%d]: %v\n", i, val)
    }
})
```

#### 断点回调
```go
// 设置断点命中回调
debugger.OnBreakpoint(func(ctx *DebugContext) {
    fmt.Printf("断点命中: PC=%d\n", ctx.PC)
    fmt.Printf("当前指令: %v\n", ctx.Instruction)
    fmt.Printf("栈状态: %v\n", ctx.Stack)
    
    // 检查变量
    for name, value := range ctx.Variables {
        fmt.Printf("变量 %s = %v\n", name, value)
    }
    
    // 显示源码位置
    if ctx.Position.Line > 0 {
        fmt.Printf("位置: 第%d行，第%d列\n", 
            ctx.Position.Line, ctx.Position.Column)
    }
})
```

#### 错误处理
```go
// 设置错误回调
debugger.OnError(func(err error) {
    fmt.Printf("执行错误: %v\n", err)
    
    // 获取当前状态
    stats := debugger.GetStats()
    fmt.Printf("错误发生时已执行指令: %d\n", stats.TotalInstructions)
})
```

### 3. 性能分析

#### 执行统计
```go
// 获取详细统计
stats := debugger.GetStats()
fmt.Printf("执行统计:\n")
fmt.Printf("  总指令数: %d\n", stats.TotalInstructions)
fmt.Printf("  执行时间: %v\n", stats.ExecutionTime)
fmt.Printf("  函数调用: %d\n", stats.FunctionCalls)
fmt.Printf("  内存分配: %d\n", stats.MemoryAllocations)

// 格式化输出统计信息
fmt.Println(debugger.FormatStats())
```

#### 指令分析
```go
// 指令计数统计
fmt.Println(debugger.FormatInstructionCounts())

// 手动获取指令统计
stats := debugger.GetStats()
for opcode, count := range stats.InstructionCounts {
    percentage := float64(count) / float64(stats.TotalInstructions) * 100
    fmt.Printf("%s: %d次 (%.2f%%)\n", opcode, count, percentage)
}
```

#### 热点分析
```go
// 获取执行热点
stats := debugger.GetStats()
fmt.Println("执行热点:")
for i, hotspot := range stats.HotSpots {
    fmt.Printf("%d. PC=%d, 操作=%v, 次数=%d, 占比=%.2f%%\n",
        i+1, hotspot.PC, hotspot.OpCode, 
        hotspot.Count, hotspot.Percentage)
    
    if hotspot.Source != "" {
        fmt.Printf("   源码: %s\n", hotspot.Source)
    }
}
```

#### 统计重置
```go
// 重置统计信息
debugger.ResetStats()

// 重新开始计时
stats := debugger.GetStats()
stats.StartTime = time.Now()
```

## 🚀 使用场景

### 1. 开发调试

#### 表达式执行跟踪
```go
func debugExpression(expr string, env map[string]interface{}) {
    debugger := debug.New()
    debugger.Enable()
    debugger.SetStepMode(true)
    
    // 记录所有执行步骤
    var steps []string
    debugger.OnStep(func(ctx *DebugContext) {
        step := fmt.Sprintf("PC=%d, Stack=%v", ctx.PC, ctx.Stack)
        steps = append(steps, step)
    })
    
    // 执行表达式 (需要VM集成)
    // result := executeWithDebugger(expr, env, debugger)
    
    // 输出执行轨迹
    fmt.Println("执行轨迹:")
    for i, step := range steps {
        fmt.Printf("%d: %s\n", i+1, step)
    }
}
```

#### 条件调试
```go
func conditionalDebug(expr string) {
    debugger := debug.New()
    debugger.Enable()
    
    // 在特定条件下断点
    bp := debugger.SetBreakpoint(20)
    bp.SetCondition("result > 100")
    bp.SetDescription("结果超过阈值时停止")
    
    debugger.OnBreakpoint(func(ctx *DebugContext) {
        fmt.Println("警告：结果超过预期阈值！")
        // 进行详细分析...
    })
}
```

### 2. 性能优化

#### 性能瓶颈识别
```go
func analyzePerformance(expr string) {
    debugger := debug.New()
    debugger.Enable()
    
    // 执行表达式
    // result := executeWithDebugger(expr, nil, debugger)
    
    // 分析性能
    stats := debugger.GetStats()
    
    // 查找最耗时的操作
    fmt.Println("性能热点:")
    for _, hotspot := range stats.HotSpots {
        if hotspot.Percentage > 10.0 { // 占比超过10%
            fmt.Printf("热点: %v 占用 %.2f%% 执行时间\n",
                hotspot.OpCode, hotspot.Percentage)
        }
    }
    
    // 分析执行效率
    avgInstructionTime := stats.ExecutionTime.Nanoseconds() / stats.TotalInstructions
    fmt.Printf("平均指令执行时间: %d纳秒\n", avgInstructionTime)
}
```

#### 内存使用分析
```go
func analyzeMemoryUsage(expr string) {
    debugger := debug.New()
    debugger.Enable()
    
    var maxStackDepth int
    debugger.OnStep(func(ctx *DebugContext) {
        if len(ctx.Stack) > maxStackDepth {
            maxStackDepth = len(ctx.Stack)
        }
    })
    
    // 执行表达式
    // result := executeWithDebugger(expr, nil, debugger)
    
    stats := debugger.GetStats()
    fmt.Printf("最大栈深度: %d\n", maxStackDepth)
    fmt.Printf("内存分配次数: %d\n", stats.MemoryAllocations)
}
```

### 3. 测试验证

#### 执行路径验证
```go
func verifyExecutionPath(expr string, expectedPath []int) {
    debugger := debug.New()
    debugger.Enable()
    debugger.SetStepMode(true)
    
    var actualPath []int
    debugger.OnStep(func(ctx *DebugContext) {
        actualPath = append(actualPath, ctx.PC)
    })
    
    // 执行表达式
    // result := executeWithDebugger(expr, nil, debugger)
    
    // 验证执行路径
    if len(actualPath) != len(expectedPath) {
        fmt.Printf("路径长度不匹配: 期望%d，实际%d\n", 
            len(expectedPath), len(actualPath))
        return
    }
    
    for i, expected := range expectedPath {
        if actualPath[i] != expected {
            fmt.Printf("路径不匹配: 位置%d，期望PC=%d，实际PC=%d\n",
                i, expected, actualPath[i])
            return
        }
    }
    
    fmt.Println("执行路径验证通过")
}
```

## ⚙️ 最佳实践

### 1. 性能考虑
```go
// 生产环境禁用调试
func createProductionDebugger() *debug.Debugger {
    debugger := debug.New()
    // 生产环境不启用调试器
    debugger.Disable()
    return debugger
}

// 开发环境启用完整调试
func createDevelopmentDebugger() *debug.Debugger {
    debugger := debug.New()
    debugger.Enable()
    
    // 只在需要时启用步进模式
    if needStepDebugging() {
        debugger.SetStepMode(true)
    }
    
    return debugger
}
```

### 2. 内存管理
```go
// 定期清理统计数据
func periodicStatsCleanup(debugger *debug.Debugger) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := debugger.GetStats()
        if stats.TotalInstructions > 1000000 { // 超过100万指令
            debugger.ResetStats()
            fmt.Println("调试统计已重置")
        }
    }
}
```

### 3. 错误处理
```go
// 健壮的调试器配置
func setupRobustDebugger() *debug.Debugger {
    debugger := debug.New()
    debugger.Enable()
    
    // 设置错误恢复
    debugger.OnError(func(err error) {
        log.Printf("调试器错误: %v", err)
        
        // 尝试恢复
        debugger.ResetStats()
        
        // 记录错误
        // errorReporter.Report(err)
    })
    
    return debugger
}
```

## 📈 性能影响

### 调试开销
- **断点检查**: 每条指令约增加 10-20ns 开销
- **统计收集**: 每条指令约增加 5-10ns 开销  
- **步进模式**: 每条指令约增加 50-100ns 开销
- **内存使用**: 约增加 1-2MB 内存占用

### 建议
1. **生产环境**: 完全禁用调试器
2. **测试环境**: 启用基本统计，按需设置断点
3. **开发环境**: 启用完整调试功能
4. **性能测试**: 使用纯统计模式，避免断点和步进

Debug模块为表达式引擎提供了全面的调试和性能分析能力，是开发、测试和优化过程中的重要工具。合理使用调试功能能够显著提高开发效率和代码质量。 