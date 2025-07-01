# Parser模块 - 语法分析器

## 概述

Parser模块是表达式引擎的语法分析层，负责将词法分析器产生的Token序列转换为抽象语法树（AST）。它实现了递归下降解析算法，支持运算符优先级处理、错误恢复和详细的语法错误报告。

## 核心功能

### 1. 语法解析
- 将Token序列转换为AST
- 支持递归下降解析
- 实现运算符优先级处理
- 支持左递归消除

### 2. 错误处理
- 详细的语法错误报告
- 错误位置定位
- 错误恢复机制
- 多错误收集

### 3. 优先级管理
- 可配置的运算符优先级
- 支持自定义操作符
- 左结合和右结合支持
- 前缀、中缀、后缀表达式

## 主要类型

### Parser结构体
```go
type Parser struct {
    l *lexer.Lexer  // 词法分析器
    
    curToken  lexer.Token  // 当前Token
    peekToken lexer.Token  // 下一个Token
    
    errors []string  // 错误列表
    
    // 解析函数映射
    prefixParseFns map[lexer.TokenType]prefixParseFn
    infixParseFns  map[lexer.TokenType]infixParseFn
}
```

### 优先级常量
```go
const (
    _ int = iota
    LOWEST      // 最低优先级
    PIPE        // | 管道操作
    TERNARY     // ? : 三元运算符
    OR          // ||
    AND         // &&
    EQUALS      // == !=
    LESSGREATER // > < >= <=
    SUM         // + -
    PRODUCT     // * / %
    POWER       // **
    PREFIX      // -x !x
    CALL        // myFunction()
    INDEX       // array[index]
    HIGHEST     // 最高优先级
)
```

## 基本使用

### 1. 创建语法分析器
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/lexer"
    "github.com/mredencom/expr/parser"
)

func main() {
    input := "age > 18 && name == 'Alice'"
    
    l := lexer.New(input)
    p := parser.New(l)
    
    program := p.ParseProgram()
    
    // 检查解析错误
    if len(p.Errors()) > 0 {
        fmt.Println("解析错误:")
        for _, err := range p.Errors() {
            fmt.Printf("  %s\n", err)
        }
        return
    }
    
    fmt.Printf("AST: %s\n", program.String())
}
```

### 2. 解析不同类型的表达式
```go
func parseExpressionTypes() {
    expressions := []string{
        "42",                          // 整数字面量
        "'hello world'",               // 字符串字面量
        "x + y * z",                   // 算术表达式
        "user.name",                   // 成员访问
        "items[0]",                    // 索引访问
        "user.*",                      // 通配符访问
        "filter(x => x > 10)",         // Lambda表达式
        "data | map(transform)",       // 管道操作
        "age > 18 ? 'adult' : 'minor'", // 三元运算符
    }
    
    for _, expr := range expressions {
        l := lexer.New(expr)
        p := parser.New(l)
        
        program := p.ParseProgram()
        if len(p.Errors()) == 0 {
            fmt.Printf("表达式: %s\n", expr)
            fmt.Printf("AST: %s\n\n", program.String())
        }
    }
}
```

### 3. 错误处理示例
```go
func handleParseErrors() {
    // 包含语法错误的表达式
    invalidExpressions := []string{
        "age > ",              // 不完整的比较表达式
        "user.",               // 不完整的成员访问
        "items[",              // 不完整的索引访问
        "( x + y",             // 不匹配的括号
    }
    
    for _, expr := range invalidExpressions {
        l := lexer.New(expr)
        p := parser.New(l)
        
        program := p.ParseProgram()
        
        if len(p.Errors()) > 0 {
            fmt.Printf("表达式: %s\n", expr)
            fmt.Println("错误:")
            for _, err := range p.Errors() {
                fmt.Printf("  %s\n", err)
            }
            fmt.Println()
        }
    }
}
```

## 高级特性

### 1. 自定义运算符优先级
```go
type CustomParser struct {
    *parser.Parser
    customPrecedences map[lexer.TokenType]int
}

func NewCustomParser(l *lexer.Lexer) *CustomParser {
    p := &CustomParser{
        Parser: parser.New(l),
        customPrecedences: make(map[lexer.TokenType]int),
    }
    
    // 自定义操作符优先级
    p.customPrecedences[lexer.CUSTOM_OP] = parser.PRODUCT
    
    return p
}

func (cp *CustomParser) getPrecedence(tokenType lexer.TokenType) int {
    if prec, ok := cp.customPrecedences[tokenType]; ok {
        return prec
    }
    return cp.Parser.GetPrecedence(tokenType)
}
```

### 2. 表达式验证器
```go
type ExpressionValidator struct {
    parser   *parser.Parser
    maxDepth int
    maxNodes int
}

func NewValidator(maxDepth, maxNodes int) *ExpressionValidator {
    return &ExpressionValidator{
        maxDepth: maxDepth,
        maxNodes: maxNodes,
    }
}

func (ev *ExpressionValidator) Validate(expression string) error {
    l := lexer.New(expression)
    p := parser.New(l)
    
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        return fmt.Errorf("语法错误: %v", p.Errors())
    }
    
    // 检查表达式复杂度
    depth := ev.calculateDepth(program)
    if depth > ev.maxDepth {
        return fmt.Errorf("表达式嵌套层级过深: %d > %d", depth, ev.maxDepth)
    }
    
    nodes := ev.countNodes(program)
    if nodes > ev.maxNodes {
        return fmt.Errorf("表达式节点数量过多: %d > %d", nodes, ev.maxNodes)
    }
    
    return nil
}

func (ev *ExpressionValidator) calculateDepth(node ast.Node) int {
    // 实现深度计算逻辑
    switch n := node.(type) {
    case *ast.InfixExpression:
        leftDepth := ev.calculateDepth(n.Left)
        rightDepth := ev.calculateDepth(n.Right)
        return 1 + max(leftDepth, rightDepth)
    case *ast.PrefixExpression:
        return 1 + ev.calculateDepth(n.Right)
    case *ast.CallExpression:
        maxDepth := ev.calculateDepth(n.Function)
        for _, arg := range n.Arguments {
            argDepth := ev.calculateDepth(arg)
            if argDepth > maxDepth {
                maxDepth = argDepth
            }
        }
        return 1 + maxDepth
    default:
        return 1
    }
}

func (ev *ExpressionValidator) countNodes(node ast.Node) int {
    count := 1
    
    switch n := node.(type) {
    case *ast.InfixExpression:
        count += ev.countNodes(n.Left)
        count += ev.countNodes(n.Right)
    case *ast.PrefixExpression:
        count += ev.countNodes(n.Right)
    case *ast.CallExpression:
        count += ev.countNodes(n.Function)
        for _, arg := range n.Arguments {
            count += ev.countNodes(arg)
        }
    // ... 处理其他节点类型
    }
    
    return count
}
```

### 3. 增量解析支持
```go
type IncrementalParser struct {
    baseParser *parser.Parser
    cache      map[string]*ast.Program
    stats      ParseStats
}

type ParseStats struct {
    TotalParses int
    CacheHits   int
    CacheMisses int
}

func NewIncrementalParser() *IncrementalParser {
    return &IncrementalParser{
        cache: make(map[string]*ast.Program),
    }
}

func (ip *IncrementalParser) Parse(expression string) (*ast.Program, error) {
    ip.stats.TotalParses++
    
    // 检查缓存
    if cached, exists := ip.cache[expression]; exists {
        ip.stats.CacheHits++
        return cached, nil
    }
    
    ip.stats.CacheMisses++
    
    // 解析新表达式
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        return nil, fmt.Errorf("解析错误: %v", p.Errors())
    }
    
    // 缓存结果
    ip.cache[expression] = program
    
    return program, nil
}

func (ip *IncrementalParser) GetStats() ParseStats {
    return ip.stats
}

func (ip *IncrementalParser) ClearCache() {
    ip.cache = make(map[string]*ast.Program)
}
```

## 语法支持

### 1. 基础表达式
```go
// 字面量
"42"          // 整数
"3.14"        // 浮点数
"\"hello\""   // 字符串
"true"        // 布尔值
"nil"         // 空值

// 标识符
"userName"
"user_name"
"_private"
```

### 2. 运算表达式
```go
// 算术运算
"a + b"       // 加法
"a - b"       // 减法
"a * b"       // 乘法
"a / b"       // 除法
"a % b"       // 取模
"a ** b"      // 幂运算

// 比较运算
"a == b"      // 等于
"a != b"      // 不等于
"a < b"       // 小于
"a <= b"      // 小于等于
"a > b"       // 大于
"a >= b"      // 大于等于

// 逻辑运算
"a && b"      // 逻辑与
"a || b"      // 逻辑或
"!a"          // 逻辑非

// 位运算
"a & b"       // 按位与
"a | b"       // 按位或
"a ^ b"       // 按位异或
"~a"          // 按位取反
"a << b"      // 左移
"a >> b"      // 右移
```

### 3. 复杂表达式
```go
// 成员访问
"user.name"
"user.profile.email"

// 索引访问
"items[0]"
"matrix[i][j]"
"map[\"key\"]"

// 函数调用
"func()"
"func(a, b, c)"
"math.max(a, b)"

// 数组字面量
"[1, 2, 3]"
"[\"a\", \"b\", \"c\"]"
"[]"

// 对象字面量（映射）
"{\"name\": \"Alice\", \"age\": 30}"
"{key: value, \"other\": 42}"

// Lambda表达式
"x => x * 2"
"(x, y) => x + y"
"item => item.name"

// 管道操作
"data | filter(condition)"
"numbers | map(x => x * 2) | sum()"

// 三元运算符
"age >= 18 ? \"adult\" : \"minor\""
"value != nil ? value : \"default\""
```

## 错误处理和调试

### 1. 详细错误报告
```go
func parseWithDetailedErrors(expression string) {
    l := lexer.New(expression)
    p := parser.New(l)
    
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        fmt.Printf("表达式: %s\n", expression)
        fmt.Println("语法错误:")
        
        for _, err := range p.Errors() {
            fmt.Printf("  错误: %s\n", err)
            
            // 可以通过错误消息提取位置信息
            if pos := extractPosition(err); pos != nil {
                fmt.Printf("  位置: 第%d行, 第%d列\n", pos.Line, pos.Column)
            }
        }
    }
}

func extractPosition(errorMsg string) *lexer.Position {
    // 从错误消息中提取位置信息的实现
    // 这取决于具体的错误消息格式
    return nil
}
```

### 2. 解析过程跟踪
```go
type TracingParser struct {
    *parser.Parser
    traceLevel int
    output     io.Writer
}

func NewTracingParser(l *lexer.Lexer, level int, output io.Writer) *TracingParser {
    return &TracingParser{
        Parser:     parser.New(l),
        traceLevel: level,
        output:     output,
    }
}

func (tp *TracingParser) trace(msg string) {
    if tp.traceLevel > 0 {
        fmt.Fprintf(tp.output, "TRACE: %s\n", msg)
    }
}

func (tp *TracingParser) ParseExpression() ast.Expression {
    tp.trace("开始解析表达式")
    expr := tp.Parser.ParseExpression()
    tp.trace(fmt.Sprintf("完成解析表达式: %s", expr.String()))
    return expr
}
```

### 3. 性能分析
```go
type PerformanceParser struct {
    *parser.Parser
    metrics ParseMetrics
}

type ParseMetrics struct {
    TotalTime     time.Duration
    TokensRead    int
    NodesCreated  int
    ErrorsFound   int
    StartTime     time.Time
}

func NewPerformanceParser(l *lexer.Lexer) *PerformanceParser {
    return &PerformanceParser{
        Parser: parser.New(l),
    }
}

func (pp *PerformanceParser) ParseProgram() *ast.Program {
    pp.metrics.StartTime = time.Now()
    
    program := pp.Parser.ParseProgram()
    
    pp.metrics.TotalTime = time.Since(pp.metrics.StartTime)
    pp.metrics.ErrorsFound = len(pp.Parser.Errors())
    pp.metrics.NodesCreated = pp.countNodes(program)
    
    return program
}

func (pp *PerformanceParser) GetMetrics() ParseMetrics {
    return pp.metrics
}

func (pp *PerformanceParser) countNodes(node ast.Node) int {
    // 实现节点计数逻辑
    count := 1
    // ... 递归计算子节点
    return count
}
```

## 最佳实践

### 1. 错误处理
```go
func robustParse(expression string) (*ast.Program, error) {
    // 预检查：验证输入
    if strings.TrimSpace(expression) == "" {
        return nil, fmt.Errorf("表达式不能为空")
    }
    
    l := lexer.New(expression)
    p := parser.New(l)
    
    // 解析
    program := p.ParseProgram()
    
    // 检查错误
    if errors := p.Errors(); len(errors) > 0 {
        return nil, fmt.Errorf("语法错误: %s", strings.Join(errors, "; "))
    }
    
    // 验证AST完整性
    if len(program.Statements) == 0 {
        return nil, fmt.Errorf("表达式为空或无效")
    }
    
    return program, nil
}
```

### 2. 性能优化
```go
type OptimizedParser struct {
    pool   sync.Pool
    cache  *lru.Cache
    stats  ParseStats
}

func NewOptimizedParser(cacheSize int) *OptimizedParser {
    cache, _ := lru.New(cacheSize)
    
    return &OptimizedParser{
        pool: sync.Pool{
            New: func() interface{} {
                return &parser.Parser{}
            },
        },
        cache: cache,
    }
}

func (op *OptimizedParser) Parse(expression string) (*ast.Program, error) {
    // 检查缓存
    if cached, found := op.cache.Get(expression); found {
        return cached.(*ast.Program), nil
    }
    
    // 从池中获取解析器
    p := op.pool.Get().(*parser.Parser)
    defer op.pool.Put(p)
    
    // 重置解析器状态
    l := lexer.New(expression)
    p.Reset(l)
    
    // 解析
    program := p.ParseProgram()
    if len(p.Errors()) > 0 {
        return nil, fmt.Errorf("解析错误: %v", p.Errors())
    }
    
    // 缓存结果
    op.cache.Add(expression, program)
    
    return program, nil
}
```

### 3. 扩展语法支持
```go
type ExtendedParser struct {
    *parser.Parser
    extensions map[string]func(*parser.Parser, lexer.Token) ast.Expression
}

func NewExtendedParser(l *lexer.Lexer) *ExtendedParser {
    ep := &ExtendedParser{
        Parser:     parser.New(l),
        extensions: make(map[string]func(*parser.Parser, lexer.Token) ast.Expression),
    }
    
    // 注册扩展语法
    ep.registerExtensions()
    
    return ep
}

func (ep *ExtendedParser) registerExtensions() {
    // 注册自定义操作符解析
    ep.extensions["MATCH"] = ep.parseMatchExpression
    ep.extensions["REGEX"] = ep.parseRegexExpression
}

func (ep *ExtendedParser) parseMatchExpression(p *parser.Parser, token lexer.Token) ast.Expression {
    // 实现 match 表达式解析
    // 例如: match value { pattern1 => result1, pattern2 => result2 }
    return nil
}

func (ep *ExtendedParser) parseRegexExpression(p *parser.Parser, token lexer.Token) ast.Expression {
    // 实现正则表达式解析
    // 例如: /pattern/flags
    return nil
}
```

## 与其他模块的集成

Parser模块在编译流水线中的位置：
```
Lexer → Token序列 → Parser → AST → Checker → 类型检查 → Compiler → 字节码
```

Parser模块是语法分析的核心，它：
1. 接收Lexer产生的Token序列
2. 根据语法规则构建AST
3. 处理运算符优先级和结合性
4. 提供详细的错误报告
5. 为后续的语义分析阶段提供结构化输入

通过合理使用Parser模块，可以构建强大且灵活的表达式解析系统。 