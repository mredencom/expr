# Compiler模块 - 字节码编译器

## 概述

Compiler模块是表达式引擎的字节码编译器，负责将经过类型检查的AST转换为高效的字节码指令序列。它实现了多种编译优化技术，包括常量折叠、死代码消除、指令合并等，为虚拟机提供高性能的执行代码。

## 核心功能

### 1. AST到字节码编译
- 递归遍历AST节点
- 生成优化的指令序列
- 常量池管理
- 符号表构建

### 2. 编译优化
- 常量折叠
- 死代码消除
- 指令合并
- 跳转优化

### 3. 字节码生成
- 紧凑的指令格式
- 高效的操作码设计
- 内存友好的数据布局
- 快速执行路径

## 主要类型

### Compiler结构体
```go
type Compiler struct {
    constants       []types.Value          // 常量池
    symbolTable     *SymbolTable          // 符号表
    scopes          []*CompilationScope   // 编译作用域栈
    instructions    []byte                // 指令序列
    lastInstruction *EmittedInstruction   // 最后发出的指令
    builtins        map[string]int        // 内置函数索引
}

type CompilationScope struct {
    instructions        []byte
    lastInstruction     *EmittedInstruction
    previousInstruction *EmittedInstruction
}

type EmittedInstruction struct {
    Opcode   vm.Opcode
    Position int
}
```

### 字节码结构
```go
type Bytecode struct {
    Instructions []byte        // 指令序列
    Constants    []types.Value // 常量池
}
```

## 基本使用

### 1. 创建编译器
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/compiler"
    "github.com/mredencom/expr/parser"
    "github.com/mredencom/expr/lexer"
)

func main() {
    // 创建编译器
    comp := compiler.New()
    
    // 解析表达式
    expression := "2 + 3 * 4"
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        panic(p.Errors())
    }
    
    // 编译
    err := comp.Compile(program.Statements[0])
    if err != nil {
        panic(err)
    }
    
    // 获取字节码
    bytecode := comp.Bytecode()
    fmt.Printf("字节码: %v\n", bytecode.Instructions)
    fmt.Printf("常量池: %v\n", bytecode.Constants)
}
```

### 2. 环境变量编译
```go
func compileWithEnvironment() {
    comp := compiler.New()
    
    // 添加环境变量
    env := map[string]interface{}{
        "x": 10,
        "y": 20,
    }
    
    adapter := env.New()
    err := comp.AddEnvironment(env, adapter)
    if err != nil {
        panic(err)
    }
    
    // 编译包含变量的表达式
    expression := "x + y * 2"
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    err = comp.Compile(program.Statements[0])
    if err != nil {
        panic(err)
    }
    
    bytecode := comp.Bytecode()
    fmt.Printf("编译完成，指令数: %d\n", len(bytecode.Instructions))
}
```

### 3. 函数调用编译
```go
func compileFunctionCalls() {
    comp := compiler.New()
    
    // 定义内置函数
    comp.DefineBuiltin("max")
    comp.DefineBuiltin("min")
    
    // 编译函数调用
    expression := "max(10, min(5, 15))"
    l := lexer.New(expression)
    p := parser.New(l)
    program := p.ParseProgram()
    
    err := comp.Compile(program.Statements[0])
    if err != nil {
        panic(err)
    }
    
    bytecode := comp.Bytecode()
    fmt.Printf("函数调用编译完成\n")
    fmt.Printf("指令序列长度: %d\n", len(bytecode.Instructions))
}
```

## 编译优化技术

### 1. 常量折叠优化
```go
type ConstantFolder struct {
    compiler *Compiler
}

func (cf *ConstantFolder) OptimizeExpression(expr ast.Expression) ast.Expression {
    switch e := expr.(type) {
    case *ast.InfixExpression:
        return cf.optimizeInfixExpression(e)
    case *ast.PrefixExpression:
        return cf.optimizePrefixExpression(e)
    default:
        return expr
    }
}

func (cf *ConstantFolder) optimizeInfixExpression(expr *ast.InfixExpression) ast.Expression {
    left := cf.OptimizeExpression(expr.Left)
    right := cf.OptimizeExpression(expr.Right)
    
    // 检查是否都是常量
    leftLit, leftOk := left.(*ast.IntegerLiteral)
    rightLit, rightOk := right.(*ast.IntegerLiteral)
    
    if leftOk && rightOk {
        // 在编译时计算结果
        switch expr.Operator {
        case "+":
            return &ast.IntegerLiteral{
                Token: expr.Token,
                Value: leftLit.Value + rightLit.Value,
            }
        case "-":
            return &ast.IntegerLiteral{
                Token: expr.Token,
                Value: leftLit.Value - rightLit.Value,
            }
        case "*":
            return &ast.IntegerLiteral{
                Token: expr.Token,
                Value: leftLit.Value * rightLit.Value,
            }
        case "/":
            if rightLit.Value != 0 {
                return &ast.IntegerLiteral{
                    Token: expr.Token,
                    Value: leftLit.Value / rightLit.Value,
                }
            }
        }
    }
    
    // 无法优化，返回原表达式
    return &ast.InfixExpression{
        Token:    expr.Token,
        Left:     left,
        Operator: expr.Operator,
        Right:    right,
    }
}
```

### 2. 指令合并优化
```go
type InstructionOptimizer struct {
    instructions []byte
}

func (io *InstructionOptimizer) Optimize() []byte {
    optimized := make([]byte, 0, len(io.instructions))
    
    for i := 0; i < len(io.instructions); {
        instruction := io.instructions[i]
        
        // 查找可优化的指令模式
        if optimized := io.tryOptimizePattern(i); optimized != nil {
            optimized = append(optimized, optimized...)
            i += len(optimized)
        } else {
            optimized = append(optimized, instruction)
            i++
        }
    }
    
    return optimized
}

func (io *InstructionOptimizer) tryOptimizePattern(pos int) []byte {
    if pos+2 >= len(io.instructions) {
        return nil
    }
    
    // 优化模式: OpConstant + OpConstant + OpAdd => OpConstantAdd
    if io.instructions[pos] == byte(vm.OpConstant) &&
       io.instructions[pos+3] == byte(vm.OpConstant) &&
       io.instructions[pos+6] == byte(vm.OpAdd) {
        
        // 检查是否是小整数常量
        const1 := int(io.instructions[pos+1])<<8 | int(io.instructions[pos+2])
        const2 := int(io.instructions[pos+4])<<8 | int(io.instructions[pos+5])
        
        if const1 < 256 && const2 < 256 {
            // 生成优化指令
            return []byte{
                byte(vm.OpConstantAdd),
                byte(const1),
                byte(const2),
            }
        }
    }
    
    return nil
}
```

### 3. 死代码消除
```go
type DeadCodeEliminator struct {
    reachable map[int]bool
}

func (dce *DeadCodeEliminator) EliminateDeadCode(instructions []byte) []byte {
    dce.reachable = make(map[int]bool)
    
    // 标记可达代码
    dce.markReachable(instructions, 0)
    
    // 移除不可达代码
    return dce.removeUnreachable(instructions)
}

func (dce *DeadCodeEliminator) markReachable(instructions []byte, pos int) {
    for pos < len(instructions) && !dce.reachable[pos] {
        dce.reachable[pos] = true
        
        opcode := vm.Opcode(instructions[pos])
        switch opcode {
        case vm.OpJump:
            // 跳转指令：标记目标地址
            target := int(instructions[pos+1])<<8 | int(instructions[pos+2])
            dce.markReachable(instructions, target)
            return // 无条件跳转后的代码不可达
            
        case vm.OpJumpNotTruthy:
            // 条件跳转：标记两个分支
            target := int(instructions[pos+1])<<8 | int(instructions[pos+2])
            dce.markReachable(instructions, target)
            pos += 3
            
        case vm.OpReturn:
            return // 返回后的代码不可达
            
        default:
            pos += dce.getInstructionWidth(opcode)
        }
    }
}

func (dce *DeadCodeEliminator) removeUnreachable(instructions []byte) []byte {
    var result []byte
    
    for i := 0; i < len(instructions); {
        if dce.reachable[i] {
            opcode := vm.Opcode(instructions[i])
            width := dce.getInstructionWidth(opcode)
            result = append(result, instructions[i:i+width]...)
        }
        i += dce.getInstructionWidth(vm.Opcode(instructions[i]))
    }
    
    return result
}
```

## 高级编译技术

### 1. Lambda表达式编译
```go
func (c *Compiler) compileLambdaExpression(node *ast.LambdaExpression) error {
    // 创建新的编译作用域
    c.enterScope()
    
    // 定义参数
    for _, param := range node.Parameters {
        c.symbolTable.Define(param.Value)
    }
    
    // 编译函数体
    err := c.Compile(node.Body)
    if err != nil {
        return err
    }
    
    // 添加返回指令
    c.emit(vm.OpReturnValue)
    
    // 获取编译后的指令
    instructions := c.leaveScope()
    
    // 创建函数对象
    numParams := len(node.Parameters)
    compiledFn := &vm.CompiledFunction{
        Instructions:  instructions,
        NumParameters: numParams,
        NumLocals:     c.symbolTable.numDefinitions,
    }
    
    // 添加到常量池并生成加载指令
    constIndex := c.addConstant(compiledFn)
    c.emit(vm.OpConstant, constIndex)
    
    return nil
}
```

### 2. 管道操作编译
```go
func (c *Compiler) compilePipeExpression(node *ast.PipeExpression) error {
    // 编译左操作数
    err := c.Compile(node.Left)
    if err != nil {
        return err
    }
    
    // 检查右操作数是否包含占位符
    if c.containsPlaceholder(node.Right) {
        return c.compilePipelineFunction(node.Right)
    }
    
    // 编译右操作数（通常是函数调用）
    if call, ok := node.Right.(*ast.CallExpression); ok {
        // 将左操作数作为第一个参数
        err = c.Compile(call.Function)
        if err != nil {
            return err
        }
        
        // 交换栈顶两个元素（函数和数据）
        c.emit(vm.OpSwap)
        
        // 编译其他参数
        for _, arg := range call.Arguments {
            err = c.Compile(arg)
            if err != nil {
                return err
            }
        }
        
        // 调用函数
        c.emit(vm.OpCall, len(call.Arguments)+1) // +1 for piped value
    } else {
        return fmt.Errorf("pipe right operand must be a function call")
    }
    
    return nil
}

// 占位符表达式编译
func (c *Compiler) compilePlaceholderExpression(node *ast.PlaceholderExpression) error {
    // 在管道上下文中，占位符被编译为特殊常量
    if c.inPipelineContext {
        constIndex := c.addConstant("__PLACEHOLDER__")
        c.emit(vm.OpConstant, constIndex)
        return nil
    }
    
    // 在非管道上下文中，占位符无效
    return fmt.Errorf("placeholder # can only be used in pipeline context")
}

// 检查表达式是否包含占位符
func (c *Compiler) containsPlaceholder(expr ast.Expression) bool {
    switch e := expr.(type) {
    case *ast.PlaceholderExpression:
        return true
    case *ast.CallExpression:
        for _, arg := range e.Arguments {
            if c.containsPlaceholder(arg) {
                return true
            }
        }
    case *ast.InfixExpression:
        return c.containsPlaceholder(e.Left) || c.containsPlaceholder(e.Right)
    case *ast.PrefixExpression:
        return c.containsPlaceholder(e.Right)
    }
    return false
}

// 编译管道函数（包含占位符）
func (c *Compiler) compilePipelineFunction(expr ast.Expression) error {
    c.inPipelineContext = true
    defer func() { c.inPipelineContext = false }()
    
    // 序列化表达式为特殊格式
    if infixExpr, ok := expr.(*ast.InfixExpression); ok {
        return c.compilePlaceholderInfixExpression(infixExpr)
    }
    
    // 编译包含占位符的函数调用
    if call, ok := expr.(*ast.CallExpression); ok {
        // 生成特殊的占位符函数标记
        constIndex := c.addConstant("__PLACEHOLDER_EXPR__")
        c.emit(vm.OpConstant, constIndex)
        
        // 编译函数名
        err := c.Compile(call.Function)
        if err != nil {
            return err
        }
        
        // 编译参数
        for _, arg := range call.Arguments {
            err = c.Compile(arg)
            if err != nil {
                return err
            }
        }
        
        // 发出管道操作指令
        c.emit(vm.OpPipeOperation, len(call.Arguments))
        return nil
    }
    
    return fmt.Errorf("unsupported placeholder expression type")
}
```

### 3. 条件表达式编译
```go
func (c *Compiler) compileConditionalExpression(node *ast.ConditionalExpression) error {
    // 编译条件
    err := c.Compile(node.Condition)
    if err != nil {
        return err
    }
    
    // 条件跳转：如果为假，跳转到false分支
    jumpNotTruthyPos := c.emit(vm.OpJumpNotTruthy, 9999)
    
    // 编译true分支
    err = c.Compile(node.TrueExpr)
    if err != nil {
        return err
    }
    
    // 跳过false分支
    jumpPos := c.emit(vm.OpJump, 9999)
    
    // 修正第一个跳转的目标地址
    jumpNotTruthyAddr := len(c.instructions)
    c.changeOperand(jumpNotTruthyPos, jumpNotTruthyAddr)
    
    // 编译false分支
    err = c.Compile(node.FalseExpr)
    if err != nil {
        return err
    }
    
    // 修正第二个跳转的目标地址
    jumpAddr := len(c.instructions)
    c.changeOperand(jumpPos, jumpAddr)
    
    return nil
}
```

## 性能优化

### 1. 指令缓存
```go
type InstructionCache struct {
    cache map[string][]byte
    mutex sync.RWMutex
}

func NewInstructionCache() *InstructionCache {
    return &InstructionCache{
        cache: make(map[string][]byte),
    }
}

func (ic *InstructionCache) Get(key string) ([]byte, bool) {
    ic.mutex.RLock()
    defer ic.mutex.RUnlock()
    
    instructions, exists := ic.cache[key]
    if !exists {
        return nil, false
    }
    
    // 返回副本以避免修改
    result := make([]byte, len(instructions))
    copy(result, instructions)
    return result, true
}

func (ic *InstructionCache) Put(key string, instructions []byte) {
    ic.mutex.Lock()
    defer ic.mutex.Unlock()
    
    // 存储副本
    cached := make([]byte, len(instructions))
    copy(cached, instructions)
    ic.cache[key] = cached
}
```

### 2. 快速路径编译
```go
func (c *Compiler) tryFastPath(expr ast.Expression) bool {
    switch e := expr.(type) {
    case *ast.IntegerLiteral:
        // 小整数快速路径
        if e.Value >= 0 && e.Value <= 255 {
            c.emit(vm.OpConstantFast, int(e.Value))
            return true
        }
        
    case *ast.InfixExpression:
        // 简单算术快速路径
        if c.isSimpleArithmetic(e) {
            return c.compileSimpleArithmetic(e)
        }
    }
    
    return false
}

func (c *Compiler) isSimpleArithmetic(expr *ast.InfixExpression) bool {
    // 检查是否是简单的整数算术
    leftLit, leftOk := expr.Left.(*ast.IntegerLiteral)
    rightLit, rightOk := expr.Right.(*ast.IntegerLiteral)
    
    if !leftOk || !rightOk {
        return false
    }
    
    // 检查操作符
    switch expr.Operator {
    case "+", "-", "*":
        return leftLit.Value >= 0 && leftLit.Value <= 255 &&
               rightLit.Value >= 0 && rightLit.Value <= 255
    }
    
    return false
}

func (c *Compiler) compileSimpleArithmetic(expr *ast.InfixExpression) bool {
    leftLit := expr.Left.(*ast.IntegerLiteral)
    rightLit := expr.Right.(*ast.IntegerLiteral)
    
    switch expr.Operator {
    case "+":
        c.emit(vm.OpAddFast, int(leftLit.Value), int(rightLit.Value))
        return true
    case "-":
        c.emit(vm.OpSubFast, int(leftLit.Value), int(rightLit.Value))
        return true
    case "*":
        c.emit(vm.OpMulFast, int(leftLit.Value), int(rightLit.Value))
        return true
    }
    
    return false
}
```

## 调试和诊断

### 1. 字节码反汇编
```go
type Disassembler struct {
    constants []types.Value
}

func NewDisassembler(constants []types.Value) *Disassembler {
    return &Disassembler{constants: constants}
}

func (d *Disassembler) Disassemble(instructions []byte) string {
    var output strings.Builder
    
    for i := 0; i < len(instructions); {
        opcode := vm.Opcode(instructions[i])
        
        output.WriteString(fmt.Sprintf("%04d %s", i, opcode.String()))
        
        switch opcode {
        case vm.OpConstant:
            constIndex := int(instructions[i+1])<<8 | int(instructions[i+2])
            output.WriteString(fmt.Sprintf(" %d", constIndex))
            if constIndex < len(d.constants) {
                output.WriteString(fmt.Sprintf(" (%s)", d.constants[constIndex].String()))
            }
            i += 3
            
        case vm.OpJump, vm.OpJumpNotTruthy:
            target := int(instructions[i+1])<<8 | int(instructions[i+2])
            output.WriteString(fmt.Sprintf(" %d", target))
            i += 3
            
        case vm.OpCall:
            numArgs := int(instructions[i+1])
            output.WriteString(fmt.Sprintf(" %d", numArgs))
            i += 2
            
        default:
            i++
        }
        
        output.WriteString("\n")
    }
    
    return output.String()
}
```

### 2. 编译统计
```go
type CompilerStats struct {
    ExpressionsCompiled int
    InstructionsEmitted int
    ConstantsCreated    int
    OptimizationsApplied int
    CompileTime         time.Duration
}

func (c *Compiler) GetStats() CompilerStats {
    return CompilerStats{
        ExpressionsCompiled: c.stats.expressionsCompiled,
        InstructionsEmitted: len(c.instructions),
        ConstantsCreated:    len(c.constants),
        OptimizationsApplied: c.stats.optimizationsApplied,
        CompileTime:         c.stats.compileTime,
    }
}
```

## 最佳实践

### 1. 编译错误处理
```go
func robustCompile(expr ast.Expression) (*vm.Bytecode, error) {
    comp := compiler.New()
    
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Compiler panic: %v", r)
        }
    }()
    
    err := comp.Compile(expr)
    if err != nil {
        return nil, fmt.Errorf("compilation failed: %v", err)
    }
    
    bytecode := comp.Bytecode()
    
    // 验证字节码完整性
    if len(bytecode.Instructions) == 0 {
        return nil, fmt.Errorf("empty bytecode generated")
    }
    
    return bytecode, nil
}
```

### 2. 内存管理
```go
type CompilerPool struct {
    pool sync.Pool
}

func NewCompilerPool() *CompilerPool {
    return &CompilerPool{
        pool: sync.Pool{
            New: func() interface{} {
                return New()
            },
        },
    }
}

func (cp *CompilerPool) Get() *Compiler {
    return cp.pool.Get().(*Compiler)
}

func (cp *CompilerPool) Put(comp *Compiler) {
    comp.Reset() // 重置状态
    cp.pool.Put(comp)
}
```

## 与其他模块的集成

Compiler模块在编译流水线中的位置：
```
Parser → AST → Checker → Compiler → Bytecode → VM → 结果
```

Compiler模块：
1. 接收类型检查后的AST
2. 生成优化的字节码指令
3. 管理常量池和符号表
4. 为VM提供可执行代码
5. 支持运行时的高效执行

通过多种编译优化技术，Compiler模块确保生成的字节码具有最佳的执行性能。 