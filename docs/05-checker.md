# Checker模块 - 静态类型检查器

## 概述

Checker模块是表达式引擎的静态类型检查器，负责在编译时验证表达式的类型正确性。它通过遍历AST节点，检查类型兼容性、变量声明、函数调用等，确保表达式在运行时不会出现类型错误。

## 核心功能

### 1. 静态类型检查
- 编译时类型验证
- 类型兼容性检查
- 自动类型推断
- 类型错误定位

### 2. 作用域管理
- 变量作用域跟踪
- 符号表管理
- 闭包变量捕获
- 作用域嵌套处理

### 3. 函数类型检查
- 参数类型验证
- 返回值类型推断
- 函数重载检查
- Lambda表达式类型推断

## 主要类型

### Checker结构体
```go
type Checker struct {
    scopes        []*Scope      // 作用域栈
    currentScope  *Scope        // 当前作用域
    errors        []string      // 类型错误列表
    builtins      map[string]*BuiltinInfo // 内置函数信息
    options       CheckerOptions // 检查选项
}

type CheckerOptions struct {
    AllowUndefinedVars bool  // 允许未定义变量
    StrictTypeCheck    bool  // 严格类型检查
    EnableInference    bool  // 启用类型推断
}
```

### Scope结构体
```go
type Scope struct {
    parent   *Scope                    // 父作用域
    symbols  map[string]*SymbolInfo    // 符号表
    level    int                       // 作用域层级
}

type SymbolInfo struct {
    Name     string              // 符号名称
    Type     types.TypeInfo      // 类型信息
    Kind     SymbolKind          // 符号种类
    Position *lexer.Position     // 定义位置
    Used     bool                // 是否被使用
}

type SymbolKind int
const (
    VarSymbol SymbolKind = iota
    FuncSymbol
    ParamSymbol
    BuiltinSymbol
)
```

## 基本使用

### 1. 创建类型检查器
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/checker"
    "github.com/mredencom/expr/parser"
    "github.com/mredencom/expr/lexer"
)

func main() {
    // 创建检查器
    c := checker.New()
    
    // 添加环境变量类型
    env := map[string]types.TypeInfo{
        "age":  {Kind: types.KindInt64, Name: "int"},
        "name": {Kind: types.KindString, Name: "string"},
    }
    
    for name, typeInfo := range env {
        c.DefineVariable(name, typeInfo)
    }
    
    // 解析表达式
    expression := "age > 18 && name != \"\""
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    // 类型检查
    err := c.Check(program)
    if err != nil {
        fmt.Printf("类型错误: %v\n", err)
        return
    }
    
    fmt.Println("类型检查通过")
}
```

### 2. 函数类型检查
```go
func checkFunctionTypes() {
    c := checker.New()
    
    // 定义内置函数
    c.DefineBuiltin("max", &checker.BuiltinInfo{
        Parameters: []types.TypeInfo{
            {Kind: types.KindInt64, Name: "int"},
            {Kind: types.KindInt64, Name: "int"},
        },
        ReturnType: types.TypeInfo{Kind: types.KindInt64, Name: "int"},
        Variadic:   false,
    })
    
    // 检查函数调用
    expression := "max(10, 20) + 5"
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    err := c.Check(program)
    if err != nil {
        fmt.Printf("类型错误: %v\n", err)
    } else {
        fmt.Println("函数调用类型检查通过")
    }
}
```

### 3. Lambda表达式类型检查
```go
func checkLambdaTypes() {
    c := checker.New()
    
    // 定义数组变量
    c.DefineVariable("numbers", types.TypeInfo{
        Kind: types.KindSlice,
        Name: "[]int",
    })
    
    // 定义filter函数
    c.DefineBuiltin("filter", &checker.BuiltinInfo{
        Parameters: []types.TypeInfo{
            {Kind: types.KindSlice, Name: "[]T"},  // 泛型支持
            {Kind: types.KindFunc, Name: "func(T) bool"},
        },
        ReturnType: types.TypeInfo{Kind: types.KindSlice, Name: "[]T"},
        Generic:    true,
    })
    
    // 检查Lambda表达式
    expression := "filter(numbers, x => x > 10)"
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    err := c.Check(program)
    if err != nil {
        fmt.Printf("Lambda类型错误: %v\n", err)
    } else {
        fmt.Println("Lambda表达式类型检查通过")
    }
}
```

## 高级特性

### 1. 类型推断系统
```go
type TypeInferrer struct {
    checker *Checker
}

func (ti *TypeInferrer) InferExpressionType(expr ast.Expression) (types.TypeInfo, error) {
    switch e := expr.(type) {
    case *ast.IntegerLiteral:
        return types.TypeInfo{Kind: types.KindInt64, Name: "int"}, nil
        
    case *ast.StringLiteral:
        return types.TypeInfo{Kind: types.KindString, Name: "string"}, nil
        
    case *ast.BooleanLiteral:
        return types.TypeInfo{Kind: types.KindBool, Name: "bool"}, nil
        
    case *ast.Identifier:
        symbol := ti.checker.LookupSymbol(e.Value)
        if symbol == nil {
            return types.TypeInfo{}, fmt.Errorf("undefined variable: %s", e.Value)
        }
        return symbol.Type, nil
        
    case *ast.InfixExpression:
        return ti.inferBinaryOperation(e)
        
    case *ast.CallExpression:
        return ti.inferFunctionCall(e)
        
    case *ast.LambdaExpression:
        return ti.inferLambdaType(e)
        
    default:
        return types.TypeInfo{}, fmt.Errorf("unsupported expression type: %T", expr)
    }
}

func (ti *TypeInferrer) inferBinaryOperation(expr *ast.InfixExpression) (types.TypeInfo, error) {
    leftType, err := ti.InferExpressionType(expr.Left)
    if err != nil {
        return types.TypeInfo{}, err
    }
    
    rightType, err := ti.InferExpressionType(expr.Right)
    if err != nil {
        return types.TypeInfo{}, err
    }
    
    return ti.inferBinaryResultType(leftType, rightType, expr.Operator)
}
```

### 2. 作用域管理器
```go
type ScopeManager struct {
    scopes []*Scope
}

func NewScopeManager() *ScopeManager {
    return &ScopeManager{
        scopes: []*Scope{NewGlobalScope()},
    }
}

func (sm *ScopeManager) PushScope() {
    currentScope := sm.CurrentScope()
    newScope := &Scope{
        parent:  currentScope,
        symbols: make(map[string]*SymbolInfo),
        level:   currentScope.level + 1,
    }
    sm.scopes = append(sm.scopes, newScope)
}

func (sm *ScopeManager) PopScope() *Scope {
    if len(sm.scopes) <= 1 {
        return nil // 不能弹出全局作用域
    }
    
    popped := sm.scopes[len(sm.scopes)-1]
    sm.scopes = sm.scopes[:len(sm.scopes)-1]
    return popped
}

func (sm *ScopeManager) CurrentScope() *Scope {
    return sm.scopes[len(sm.scopes)-1]
}

func (sm *ScopeManager) DefineSymbol(name string, typeInfo types.TypeInfo, kind SymbolKind) error {
    currentScope := sm.CurrentScope()
    
    // 检查当前作用域是否已存在
    if _, exists := currentScope.symbols[name]; exists {
        return fmt.Errorf("symbol '%s' already defined in current scope", name)
    }
    
    currentScope.symbols[name] = &SymbolInfo{
        Name: name,
        Type: typeInfo,
        Kind: kind,
        Used: false,
    }
    
    return nil
}

func (sm *ScopeManager) LookupSymbol(name string) *SymbolInfo {
    // 从当前作用域向上查找
    for i := len(sm.scopes) - 1; i >= 0; i-- {
        scope := sm.scopes[i]
        if symbol, exists := scope.symbols[name]; exists {
            symbol.Used = true
            return symbol
        }
    }
    return nil
}
```

### 3. 错误收集和报告
```go
type ErrorCollector struct {
    errors []CheckError
}

type CheckError struct {
    Message  string
    Position *lexer.Position
    Type     ErrorType
}

type ErrorType int
const (
    TypeError ErrorType = iota
    UndefinedError
    RedefinitionError
    ArgumentError
    ReturnTypeError
)

func (ec *ErrorCollector) AddError(msg string, pos *lexer.Position, errType ErrorType) {
    ec.errors = append(ec.errors, CheckError{
        Message:  msg,
        Position: pos,
        Type:     errType,
    })
}

func (ec *ErrorCollector) HasErrors() bool {
    return len(ec.errors) > 0
}

func (ec *ErrorCollector) GetErrors() []CheckError {
    return ec.errors
}

func (ec *ErrorCollector) FormatErrors() string {
    var sb strings.Builder
    
    for _, err := range ec.errors {
        sb.WriteString(fmt.Sprintf("[%s] %s", 
            ec.errorTypeString(err.Type), err.Message))
        
        if err.Position != nil {
            sb.WriteString(fmt.Sprintf(" at line %d, column %d", 
                err.Position.Line, err.Position.Column))
        }
        sb.WriteString("\n")
    }
    
    return sb.String()
}

func (ec *ErrorCollector) errorTypeString(errType ErrorType) string {
    switch errType {
    case TypeError:
        return "TYPE ERROR"
    case UndefinedError:
        return "UNDEFINED ERROR"
    case RedefinitionError:
        return "REDEFINITION ERROR"
    case ArgumentError:
        return "ARGUMENT ERROR"
    case ReturnTypeError:
        return "RETURN TYPE ERROR"
    default:
        return "UNKNOWN ERROR"
    }
}
```

## 类型检查规则

### 1. 算术运算检查
```go
func (c *Checker) checkArithmeticOperation(expr *ast.InfixExpression) error {
    leftType, err := c.inferType(expr.Left)
    if err != nil {
        return err
    }
    
    rightType, err := c.inferType(expr.Right)
    if err != nil {
        return err
    }
    
    switch expr.Operator {
    case "+":
        // 数值加法或字符串连接
        if (types.IsNumeric(leftType) && types.IsNumeric(rightType)) ||
           (leftType.Kind == types.KindString && rightType.Kind == types.KindString) {
            return nil
        }
        return fmt.Errorf("invalid operands for +: %s and %s", leftType.Name, rightType.Name)
        
    case "-", "*", "/", "%":
        // 仅数值运算
        if types.IsNumeric(leftType) && types.IsNumeric(rightType) {
            return nil
        }
        return fmt.Errorf("invalid operands for %s: %s and %s", 
            expr.Operator, leftType.Name, rightType.Name)
        
    case "**":
        // 幂运算
        if types.IsNumeric(leftType) && types.IsNumeric(rightType) {
            return nil
        }
        return fmt.Errorf("power operator requires numeric operands")
    }
    
    return fmt.Errorf("unknown arithmetic operator: %s", expr.Operator)
}
```

### 2. 比较运算检查
```go
func (c *Checker) checkComparisonOperation(expr *ast.InfixExpression) error {
    leftType, err := c.inferType(expr.Left)
    if err != nil {
        return err
    }
    
    rightType, err := c.inferType(expr.Right)
    if err != nil {
        return err
    }
    
    switch expr.Operator {
    case "==", "!=":
        // 相等比较：类型必须兼容
        if c.areTypesCompatible(leftType, rightType) {
            return nil
        }
        return fmt.Errorf("incompatible types for equality: %s and %s", 
            leftType.Name, rightType.Name)
            
    case "<", "<=", ">", ">=":
        // 大小比较：必须是可排序的类型
        if c.areTypesComparable(leftType, rightType) {
            return nil
        }
        return fmt.Errorf("incomparable types: %s and %s", 
            leftType.Name, rightType.Name)
    }
    
    return fmt.Errorf("unknown comparison operator: %s", expr.Operator)
}

func (c *Checker) areTypesCompatible(left, right types.TypeInfo) bool {
    // 相同类型
    if left.Kind == right.Kind {
        return true
    }
    
    // 数值类型之间兼容
    if types.IsNumeric(left) && types.IsNumeric(right) {
        return true
    }
    
    // nil与任何类型兼容
    if left.Kind == types.KindNil || right.Kind == types.KindNil {
        return true
    }
    
    return false
}

func (c *Checker) areTypesComparable(left, right types.TypeInfo) bool {
    // 数值类型可比较
    if types.IsNumeric(left) && types.IsNumeric(right) {
        return true
    }
    
    // 字符串可比较
    if left.Kind == types.KindString && right.Kind == types.KindString {
        return true
    }
    
    return false
}
```

### 3. 函数调用检查
```go
func (c *Checker) checkFunctionCall(expr *ast.CallExpression) error {
    // 检查函数是否存在
    funcName := ""
    if ident, ok := expr.Function.(*ast.Identifier); ok {
        funcName = ident.Value
    } else {
        return fmt.Errorf("complex function expressions not yet supported")
    }
    
    builtin := c.builtins[funcName]
    if builtin == nil {
        return fmt.Errorf("undefined function: %s", funcName)
    }
    
    // 检查参数数量
    expectedArgs := len(builtin.Parameters)
    actualArgs := len(expr.Arguments)
    
    if !builtin.Variadic && actualArgs != expectedArgs {
        return fmt.Errorf("function %s expects %d arguments, got %d", 
            funcName, expectedArgs, actualArgs)
    }
    
    if builtin.Variadic && actualArgs < expectedArgs {
        return fmt.Errorf("function %s expects at least %d arguments, got %d", 
            funcName, expectedArgs, actualArgs)
    }
    
    // 检查参数类型
    for i, arg := range expr.Arguments {
        argType, err := c.inferType(arg)
        if err != nil {
            return fmt.Errorf("error in argument %d: %v", i+1, err)
        }
        
        var expectedType types.TypeInfo
        if i < len(builtin.Parameters) {
            expectedType = builtin.Parameters[i]
        } else if builtin.Variadic {
            // 可变参数使用最后一个参数类型
            expectedType = builtin.Parameters[len(builtin.Parameters)-1]
        }
        
        if !c.isAssignable(argType, expectedType) {
            return fmt.Errorf("argument %d: expected %s, got %s", 
                i+1, expectedType.Name, argType.Name)
        }
    }
    
    return nil
}

func (c *Checker) isAssignable(from, to types.TypeInfo) bool {
    // 相同类型
    if from.Kind == to.Kind {
        return true
    }
    
    // 数值类型的隐式转换
    if types.IsNumeric(from) && types.IsNumeric(to) {
        return true
    }
    
    // 任何类型都可以转换为interface{}
    if to.Name == "interface{}" {
        return true
    }
    
    return false
}
```

## 性能优化

### 1. 增量类型检查
```go
type IncrementalChecker struct {
    baseChecker *Checker
    cache       map[string]*CheckResult
    version     int
}

type CheckResult struct {
    Version   int
    Errors    []CheckError
    TypeInfo  types.TypeInfo
    Symbols   map[string]*SymbolInfo
}

func NewIncrementalChecker() *IncrementalChecker {
    return &IncrementalChecker{
        baseChecker: New(),
        cache:       make(map[string]*CheckResult),
        version:     0,
    }
}

func (ic *IncrementalChecker) CheckExpression(expr string, env map[string]types.TypeInfo) (*CheckResult, error) {
    // 计算缓存键
    key := ic.computeCacheKey(expr, env)
    
    // 检查缓存
    if cached, exists := ic.cache[key]; exists && cached.Version == ic.version {
        return cached, nil
    }
    
    // 执行类型检查
    result, err := ic.performCheck(expr, env)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果
    result.Version = ic.version
    ic.cache[key] = result
    
    return result, nil
}

func (ic *IncrementalChecker) InvalidateCache() {
    ic.version++
    // 可选择性清理旧版本缓存
    for key, result := range ic.cache {
        if result.Version < ic.version-10 { // 保留最近10个版本
            delete(ic.cache, key)
        }
    }
}
```

### 2. 并行类型检查
```go
type ParallelChecker struct {
    workerCount int
    taskQueue   chan CheckTask
    resultQueue chan CheckResult
}

type CheckTask struct {
    ID         string
    Expression ast.Expression
    Context    *CheckContext
}

func NewParallelChecker(workerCount int) *ParallelChecker {
    pc := &ParallelChecker{
        workerCount: workerCount,
        taskQueue:   make(chan CheckTask, workerCount*2),
        resultQueue: make(chan CheckResult, workerCount*2),
    }
    
    // 启动工作者
    for i := 0; i < workerCount; i++ {
        go pc.worker()
    }
    
    return pc
}

func (pc *ParallelChecker) worker() {
    for task := range pc.taskQueue {
        checker := New()
        // 设置检查上下文
        checker.SetContext(task.Context)
        
        // 执行检查
        err := checker.CheckExpression(task.Expression)
        
        result := CheckResult{
            ID:    task.ID,
            Error: err,
        }
        
        pc.resultQueue <- result
    }
}

func (pc *ParallelChecker) CheckBatch(expressions []ast.Expression) map[string]error {
    results := make(map[string]error)
    
    // 提交任务
    for i, expr := range expressions {
        task := CheckTask{
            ID:         fmt.Sprintf("task_%d", i),
            Expression: expr,
            Context:    NewCheckContext(),
        }
        pc.taskQueue <- task
    }
    
    // 收集结果
    for i := 0; i < len(expressions); i++ {
        result := <-pc.resultQueue
        results[result.ID] = result.Error
    }
    
    return results
}
```

## 实际应用示例

### 1. 表达式验证器
```go
type ExpressionValidator struct {
    checker *Checker
    rules   []ValidationRule
}

type ValidationRule func(ast.Expression, *Checker) error

func NewExpressionValidator() *ExpressionValidator {
    ev := &ExpressionValidator{
        checker: New(),
        rules:   make([]ValidationRule, 0),
    }
    
    // 添加默认规则
    ev.AddRule(ev.checkDepthLimit)
    ev.AddRule(ev.checkComplexityLimit)
    ev.AddRule(ev.checkSecurityConstraints)
    
    return ev
}

func (ev *ExpressionValidator) AddRule(rule ValidationRule) {
    ev.rules = append(ev.rules, rule)
}

func (ev *ExpressionValidator) Validate(expr ast.Expression) []error {
    var errors []error
    
    // 基本类型检查
    if err := ev.checker.CheckExpression(expr); err != nil {
        errors = append(errors, err)
    }
    
    // 应用自定义规则
    for _, rule := range ev.rules {
        if err := rule(expr, ev.checker); err != nil {
            errors = append(errors, err)
        }
    }
    
    return errors
}

func (ev *ExpressionValidator) checkDepthLimit(expr ast.Expression, checker *Checker) error {
    const maxDepth = 20
    depth := calculateExpressionDepth(expr)
    if depth > maxDepth {
        return fmt.Errorf("expression depth %d exceeds limit %d", depth, maxDepth)
    }
    return nil
}

func (ev *ExpressionValidator) checkComplexityLimit(expr ast.Expression, checker *Checker) error {
    const maxNodes = 100
    nodes := countNodes(expr)
    if nodes > maxNodes {
        return fmt.Errorf("expression complexity %d exceeds limit %d", nodes, maxNodes)
    }
    return nil
}

func (ev *ExpressionValidator) checkSecurityConstraints(expr ast.Expression, checker *Checker) error {
    // 检查危险函数调用
    dangerousFunctions := []string{"eval", "exec", "system"}
    
    return ast.Walk(expr, func(node ast.Node) error {
        if call, ok := node.(*ast.CallExpression); ok {
            if ident, ok := call.Function.(*ast.Identifier); ok {
                for _, dangerous := range dangerousFunctions {
                    if ident.Value == dangerous {
                        return fmt.Errorf("dangerous function call: %s", dangerous)
                    }
                }
            }
        }
        return nil
    })
}
```

## 最佳实践

### 1. 错误处理策略
```go
func robustTypeCheck(expression string, env map[string]types.TypeInfo) error {
    defer func() {
        if r := recover(); r != nil {
            // 记录panic信息用于调试
            log.Printf("Type checker panic: %v", r)
        }
    }()
    
    // 预验证
    if strings.TrimSpace(expression) == "" {
        return fmt.Errorf("empty expression")
    }
    
    // 解析
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        return fmt.Errorf("parse errors: %v", p.Errors())
    }
    
    // 类型检查
    checker := New()
    for name, typeInfo := range env {
        checker.DefineVariable(name, typeInfo)
    }
    
    return checker.Check(program)
}
```

### 2. 性能监控
```go
type CheckerMetrics struct {
    CheckCount     int64
    SuccessCount   int64
    ErrorCount     int64
    TotalTime      time.Duration
    AverageTime    time.Duration
}

func (c *Checker) CheckWithMetrics(expr ast.Expression) (error, *CheckerMetrics) {
    start := time.Now()
    
    err := c.CheckExpression(expr)
    
    duration := time.Since(start)
    
    metrics := &CheckerMetrics{
        CheckCount:  1,
        TotalTime:   duration,
        AverageTime: duration,
    }
    
    if err != nil {
        metrics.ErrorCount = 1
    } else {
        metrics.SuccessCount = 1
    }
    
    return err, metrics
}
```

## 与其他模块的集成

Checker模块在编译流水线中的作用：

```
Parser → AST → Checker → 类型化AST → Compiler → 优化字节码
```

Checker模块：
1. 接收Parser产生的AST
2. 执行静态类型检查和语义分析
3. 产生类型化的AST或报告错误
4. 为Compiler提供类型信息用于优化
5. 确保运行时类型安全

通过静态类型检查，Checker模块显著提高了表达式的可靠性和执行效率。 