# 最佳实践指南

## 🎯 企业级使用模式

### 1. 表达式预编译

**问题**: 在生产环境中重复编译相同表达式会影响性能。

**解决方案**: 使用预编译模式，一次编译多次执行。

```go
// ✅ 推荐：预编译表达式
type UserProcessor struct {
    nameProgram *expr.Program
}

func NewUserProcessor() *UserProcessor {
    program, _ := expr.Compile("user.firstName + ' ' + user.lastName", expr.AsString())
    return &UserProcessor{nameProgram: program}
}

func (p *UserProcessor) ProcessUsers(users []User) []string {
    var results []string
    for _, user := range users {
        env := map[string]interface{}{"user": user}
        result, _ := expr.Run(p.nameProgram, env)
        results = append(results, result.(string))
    }
    return results
}
```

### 2. 类型安全配置

**问题**: 运行时类型错误难以调试。

**解决方案**: 使用类型提示和环境预定义。

```go
// ✅ 推荐：提供完整的类型信息
program, err := expr.Compile("user.age > minAge && user.active",
    expr.AsBool(),
    expr.Env(map[string]interface{}{
        "user": User{}, // 提供类型示例
        "minAge": 0,
    }))
```

### 3. 资源控制

**问题**: 恶意或错误的表达式可能消耗过多资源。

**解决方案**: 设置执行超时和迭代限制。

```go
// ✅ 推荐：设置资源限制
config := expr.Config{
    Timeout:       5 * time.Second,
    MaxIterations: 100000,
}

program, err := expr.CompileWithConfig(expression, config)
```

### 4. 错误处理策略

**问题**: 表达式执行失败影响整个业务流程。

**解决方案**: 实现优雅的错误处理和降级机制。

```go
func EvaluateWithFallback(program *expr.Program, env interface{}, fallback interface{}) interface{} {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("表达式执行panic: %v", r)
        }
    }()

    result, err := expr.Run(program, env)
    if err != nil {
        log.Printf("表达式执行失败: %v", err)
        return fallback
    }
    
    return result
}
```

## 🏗️ 架构设计模式

### 1. 规则引擎模式

```go
type Rule struct {
    Name       string
    Expression string
    Priority   int
    Action     func(interface{}) error
}

type RuleEngine struct {
    rules []*CompiledRule
}

func NewRuleEngine(rules []Rule) (*RuleEngine, error) {
    var compiledRules []*CompiledRule
    
    for _, rule := range rules {
        program, err := expr.Compile(rule.Expression, expr.AsBool())
        if err != nil {
            return nil, fmt.Errorf("编译规则 %s 失败: %w", rule.Name, err)
        }
        
        compiledRules = append(compiledRules, &CompiledRule{
            Rule:    &rule,
            program: program,
        })
    }
    
    return &RuleEngine{rules: compiledRules}, nil
}
```

### 2. 配置驱动模式

```go
type ExpressionConfig struct {
    Name       string            `json:"name"`
    Expression string            `json:"expression"`
    Type       string            `json:"type"`
    Env        map[string]interface{} `json:"env"`
    Timeout    int               `json:"timeout_ms"`
}

type ConfigurableProcessor struct {
    programs map[string]*expr.Program
    configs  map[string]ExpressionConfig
}

func NewConfigurableProcessor(configs []ExpressionConfig) (*ConfigurableProcessor, error) {
    processor := &ConfigurableProcessor{
        programs: make(map[string]*expr.Program),
        configs:  make(map[string]ExpressionConfig),
    }
    
    for _, config := range configs {
        options := []expr.Option{expr.Env(config.Env)}
        
        switch config.Type {
        case "string":
            options = append(options, expr.AsString())
        case "int":
            options = append(options, expr.AsInt())
        case "bool":
            options = append(options, expr.AsBool())
        }
        
        if config.Timeout > 0 {
            options = append(options, 
                expr.WithTimeout(time.Duration(config.Timeout)*time.Millisecond))
        }
        
        program, err := expr.Compile(config.Expression, options...)
        if err != nil {
            return nil, fmt.Errorf("编译配置 %s 失败: %w", config.Name, err)
        }
        
        processor.programs[config.Name] = program
        processor.configs[config.Name] = config
    }
    
    return processor, nil
}

func (cp *ConfigurableProcessor) Execute(name string, data interface{}) (interface{}, error) {
    program, exists := cp.programs[name]
    if !exists {
        return nil, fmt.Errorf("配置 %s 不存在", name)
    }
    
    return expr.Run(program, data)
}
```

### 3. 多租户模式

```go
type TenantExpressionManager struct {
    tenantPrograms map[string]map[string]*expr.Program
    mu             sync.RWMutex
}

func NewTenantExpressionManager() *TenantExpressionManager {
    return &TenantExpressionManager{
        tenantPrograms: make(map[string]map[string]*expr.Program),
    }
}

func (tem *TenantExpressionManager) CompileForTenant(tenantID, name, expression string, options ...expr.Option) error {
    program, err := expr.Compile(expression, options...)
    if err != nil {
        return err
    }
    
    tem.mu.Lock()
    defer tem.mu.Unlock()
    
    if _, exists := tem.tenantPrograms[tenantID]; !exists {
        tem.tenantPrograms[tenantID] = make(map[string]*expr.Program)
    }
    
    tem.tenantPrograms[tenantID][name] = program
    return nil
}

func (tem *TenantExpressionManager) ExecuteForTenant(tenantID, name string, env interface{}) (interface{}, error) {
    tem.mu.RLock()
    defer tem.mu.RUnlock()
    
    tenantPrograms, exists := tem.tenantPrograms[tenantID]
    if !exists {
        return nil, fmt.Errorf("租户 %s 不存在", tenantID)
    }
    
    program, exists := tenantPrograms[name]
    if !exists {
        return nil, fmt.Errorf("租户 %s 的表达式 %s 不存在", tenantID, name)
    }
    
    return expr.Run(program, env)
}
```

## 🚀 性能优化技巧

### 1. 批量处理

```go
func ProcessBatch(program *expr.Program, dataList []interface{}) []interface{} {
    results := make([]interface{}, len(dataList))
    
    var wg sync.WaitGroup
    for i, data := range dataList {
        wg.Add(1)
        go func(idx int, d interface{}) {
            defer wg.Done()
            result, _ := expr.Run(program, d)
            results[idx] = result
        }(i, data)
    }
    
    wg.Wait()
    return results
}
```

### 2. 环境对象重用

```go
type EnvironmentPool struct {
    pool sync.Pool
}

func NewEnvironmentPool() *EnvironmentPool {
    return &EnvironmentPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make(map[string]interface{})
            },
        },
    }
}

func (ep *EnvironmentPool) ExecuteWithPool(program *expr.Program, data interface{}) interface{} {
    env := ep.pool.Get().(map[string]interface{})
    defer func() {
        for k := range env {
            delete(env, k)
        }
        ep.pool.Put(env)
    }()
    
    env["data"] = data
    result, _ := expr.Run(program, env)
    return result
}
```

## 🛡️ 安全性考虑

### 1. 输入验证

```go
func ValidateExpression(expression string) error {
    if len(expression) > 10000 {
        return errors.New("表达式过长")
    }
    
    dangerousPatterns := []string{"os.", "exec.", "syscall."}
    for _, pattern := range dangerousPatterns {
        if strings.Contains(expression, pattern) {
            return fmt.Errorf("表达式包含危险关键字: %s", pattern)
        }
    }
    
    _, err := expr.Compile(expression)
    return err
}
```

### 2. 权限控制

```go
type PermissionChecker struct {
    allowedFunctions map[string]bool
    allowedModules   map[string]bool
}

func (pc *PermissionChecker) ValidateExpression(expression string, userRole string) error {
    // 根据用户角色检查权限
    if userRole != "admin" {
        // 普通用户不能使用系统函数
        restrictedFunctions := []string{"exec", "system", "file"}
        for _, fn := range restrictedFunctions {
            if strings.Contains(expression, fn) {
                return fmt.Errorf("用户 %s 无权使用函数 %s", userRole, fn)
            }
        }
    }
    
    return nil
}
```

## 📊 监控和观测

### 1. 执行统计

```go
type ExecutionMetrics struct {
    TotalExecutions int64
    TotalDuration   time.Duration
    ErrorCount      int64
    mu             sync.RWMutex
}

func (em *ExecutionMetrics) RecordExecution(duration time.Duration, err error) {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    em.TotalExecutions++
    em.TotalDuration += duration
    if err != nil {
        em.ErrorCount++
    }
}

func MonitoredRun(program *expr.Program, env interface{}, metrics *ExecutionMetrics) (interface{}, error) {
    start := time.Now()
    result, err := expr.Run(program, env)
    duration := time.Since(start)
    
    metrics.RecordExecution(duration, err)
    
    if duration > 100*time.Millisecond {
        log.Printf("慢表达式执行: 耗时 %v", duration)
    }
    
    return result, err
}
```

### 2. 性能监控

```go
func MonitoredRun(program *expr.Program, env interface{}, metrics *ExecutionMetrics) (interface{}, error) {
    start := time.Now()
    result, err := expr.Run(program, env)
    duration := time.Since(start)
    
    // 记录指标
    metrics.RecordExecution(duration, err)
    
    // 慢查询告警
    if duration > 100*time.Millisecond {
        log.Printf("慢表达式执行: 耗时 %v", duration)
    }
    
    return result, err
}
```

## 🔧 故障排除

### 1. 表达式诊断

```go
func DiagnoseExpression(expression string, env interface{}) {
    fmt.Printf("=== 表达式诊断 ===\n")
    fmt.Printf("表达式: %s\n", expression)
    
    program, err := expr.Compile(expression)
    if err != nil {
        fmt.Printf("❌ 编译失败: %v\n", err)
        return
    }
    fmt.Printf("✅ 编译成功\n")
    
    result, err := expr.Run(program, env)
    if err != nil {
        fmt.Printf("❌ 执行失败: %v\n", err)
        return
    }
    
    fmt.Printf("✅ 执行成功\n")
    fmt.Printf("结果: %v (类型: %T)\n", result, result)
}
```

### 2. 调试模式

```go
func RunWithDebug(program *expr.Program, env interface{}) (interface{}, error) {
    debugger := debug.NewDebugger()
    
    // 设置详细日志
    debugger.SetExecutionCallback(func(step int, opcode string, value interface{}) {
        fmt.Printf("Step %d: %s -> %v\n", step, opcode, value)
    })
    
    return debugger.StepThrough(program, env), nil
}
```

这些最佳实践将帮助您在企业环境中安全、高效地使用Expr表达式引擎。 