# 性能优化策略

## 概述

本文档详细说明无反射 expr 实现的各种性能优化策略，目标是达到或超越原版 expr 的性能表现。

## 编译时优化

### 1. 常量折叠 (Constant Folding)

在编译阶段计算常量表达式，减少运行时计算开销。

```go
type ConstantFolder struct {
    changed bool
}

func (cf *ConstantFolder) OptimizeAST(node ASTNode) ASTNode {
    switch n := node.(type) {
    case *BinaryOpNode:
        left := cf.OptimizeAST(n.left)
        right := cf.OptimizeAST(n.right)
        
        // 如果两个操作数都是常量，直接计算结果
        if leftLit, leftOk := left.(*LiteralNode); leftOk {
            if rightLit, rightOk := right.(*LiteralNode); rightOk {
                return cf.foldBinaryOp(n.operator, leftLit, rightLit)
            }
        }
        
        return &BinaryOpNode{
            pos: n.pos,
            left: left,
            operator: n.operator,
            right: right,
            typ: n.typ,
        }
    }
}

func (cf *ConstantFolder) foldBinaryOp(op Operator, left, right *LiteralNode) *LiteralNode {
    switch op {
    case OpAdd:
        if leftInt, ok := left.value.(IntValue); ok {
            if rightInt, ok := right.value.(IntValue); ok {
                return &LiteralNode{
                    value: IntValue{value: leftInt.value + rightInt.value},
                    typ: TypeInfo{Kind: KindInt64},
                }
            }
        }
    // ... 其他操作符
    }
    return nil
}
```

### 2. 死代码消除 (Dead Code Elimination)

移除永远不会执行的代码分支。

```go
func (optimizer *Optimizer) eliminateDeadCode(node ASTNode) ASTNode {
    switch n := node.(type) {
    case *ConditionalNode:
        condition := optimizer.eliminateDeadCode(n.condition)
        
        // 如果条件是常量，直接返回对应分支
        if condLit, ok := condition.(*LiteralNode); ok {
            if condBool, ok := condLit.value.(BoolValue); ok {
                if condBool.value {
                    return optimizer.eliminateDeadCode(n.trueExpr)
                } else {
                    return optimizer.eliminateDeadCode(n.falseExpr)
                }
            }
        }
    }
    return node
}
```

### 3. 字节码级别优化

#### 指令合并
```go
// 将连续的 Push + Pop 操作优化掉
func (optimizer *BytecodeOptimizer) mergeInstructions(instructions []Instruction) []Instruction {
    result := make([]Instruction, 0, len(instructions))
    
    for i := 0; i < len(instructions); i++ {
        current := instructions[i]
        
        // 检查 Push + Pop 模式
        if current.OpCode == OpPush && i+1 < len(instructions) {
            if instructions[i+1].OpCode == OpPop {
                i++ // 跳过这两条指令
                continue
            }
        }
        
        result = append(result, current)
    }
    
    return result
}
```

## 运行时优化

### 1. 内联执行 (Inline Execution)

对于简单操作，生成专门的内联指令。

```go
// 为基本类型操作生成专门的指令
func (compiler *Compiler) emitBinaryOp(op Operator, leftType, rightType TypeInfo) {
    if leftType.Kind == KindInt64 && rightType.Kind == KindInt64 {
        switch op {
        case OpAdd:
            compiler.emit(OpAddInt64)
            return
        case OpSub:
            compiler.emit(OpSubInt64) 
            return
        case OpMul:
            compiler.emit(OpMulInt64)
            return
        }
    }
    
    if leftType.Kind == KindString && rightType.Kind == KindString {
        switch op {
        case OpAdd:
            compiler.emit(OpAddString)
            return
        }
    }
    
    // 回退到通用操作
    compiler.emit(OpGenericBinaryOp, int(op))
}
```

### 2. 快速路径优化 (Fast Path)

为常见情况提供优化的执行路径。

```go
func (vm *VM) executeAddInt64() error {
    // 直接操作栈顶两个元素，避免类型检查
    if vm.sp < 2 {
        return fmt.Errorf("stack underflow")
    }
    
    // 快速类型断言
    right := vm.stack[vm.sp-1].(IntValue)
    left := vm.stack[vm.sp-2].(IntValue)
    
    // 直接计算并更新栈
    vm.stack[vm.sp-2] = IntValue{value: left.value + right.value}
    vm.sp--
    
    return nil
}
```

### 3. 栈优化

#### 栈空间预分配
```go
func NewVM(bytecode *Bytecode) *VM {
    // 根据字节码分析预估栈深度
    maxDepth := analyzeStackDepth(bytecode)
    
    return &VM{
        stack:    make([]Value, maxDepth),
        sp:       0,
        maxStack: maxDepth,
        bytecode: bytecode,
    }
}

func analyzeStackDepth(bytecode *Bytecode) int {
    maxDepth := 0
    currentDepth := 0
    
    for _, inst := range bytecode.Instructions {
        switch inst.OpCode {
        case OpPush:
            currentDepth++
        case OpPop:
            currentDepth--
        case OpAddInt64, OpSubInt64, OpMulInt64, OpDivInt64:
            currentDepth-- // 两个操作数变成一个结果
        }
        
        if currentDepth > maxDepth {
            maxDepth = currentDepth
        }
    }
    
    return maxDepth + 10 // 添加一些缓冲
}
```

### 4. 内存管理优化

#### 对象池 (Object Pooling)
```go
type ValuePool struct {
    intPool    sync.Pool
    floatPool  sync.Pool
    stringPool sync.Pool
}

func (vp *ValuePool) GetInt() *IntValue {
    if v := vp.intPool.Get(); v != nil {
        return v.(*IntValue)
    }
    return &IntValue{}
}

func (vp *ValuePool) PutInt(v *IntValue) {
    v.value = 0
    vp.intPool.Put(v)
}
```

## 特化优化 (Specialization)

### 1. 类型特化

为特定类型组合生成专门的函数。

```go
// 代码生成器生成的特化函数
func executeAddIntInt(vm *VM) {
    right := vm.stack[vm.sp-1].(*IntValue)
    left := vm.stack[vm.sp-2].(*IntValue)
    
    // 溢出检查（可选）
    if (left.value > 0 && right.value > math.MaxInt64-left.value) ||
       (left.value < 0 && right.value < math.MinInt64-left.value) {
        panic("integer overflow")
    }
    
    vm.stack[vm.sp-2] = &IntValue{value: left.value + right.value}
    vm.sp--
}
```

### 2. 环境访问特化

为不同的环境类型生成特化的访问函数。

```go
// 为用户环境类型生成的特化访问器
type UserEnvAccessor struct{}

func (uea *UserEnvAccessor) GetField(env interface{}, fieldName string) Value {
    userEnv := env.(*UserEnvironment) // 编译时确定的类型
    
    switch fieldName {
    case "name":
        return StringValue{value: userEnv.Name}
    case "age":
        return IntValue{value: int64(userEnv.Age)}
    case "active":
        return BoolValue{value: userEnv.Active}
    default:
        panic("unknown field: " + fieldName)
    }
}
```

## 缓存策略

### 1. 字节码缓存

```go
type BytecodeCache struct {
    cache map[string]*CachedBytecode
    mu    sync.RWMutex
}

type CachedBytecode struct {
    bytecode   *Bytecode
    lastAccess time.Time
    accessCount int64
}

func (bc *BytecodeCache) Get(expression string, envType reflect.Type) *Bytecode {
    key := fmt.Sprintf("%s|%s", expression, envType.String())
    
    bc.mu.RLock()
    if cached, exists := bc.cache[key]; exists {
        cached.lastAccess = time.Now()
        atomic.AddInt64(&cached.accessCount, 1)
        bc.mu.RUnlock()
        return cached.bytecode
    }
    bc.mu.RUnlock()
    
    return nil
}
```

### 2. 类型信息缓存

```go
type TypeInfoCache struct {
    cache map[reflect.Type]*TypeInfo
    mu    sync.RWMutex
}

func (tic *TypeInfoCache) GetTypeInfo(t reflect.Type) *TypeInfo {
    tic.mu.RLock()
    if info, exists := tic.cache[t]; exists {
        tic.mu.RUnlock()
        return info
    }
    tic.mu.RUnlock()
    
    tic.mu.Lock()
    defer tic.mu.Unlock()
    
    info := tic.analyzeType(t)
    tic.cache[t] = info
    return info
}
```

## 并发优化

### 1. 并发安全的值类型

```go
// 不可变值类型，天然并发安全
type ImmutableIntValue struct {
    value int64
}

func (iiv ImmutableIntValue) Add(other ImmutableIntValue) ImmutableIntValue {
    return ImmutableIntValue{value: iiv.value + other.value}
}

// 写时复制的复合类型
type COWSliceValue struct {
    data   []Value
    shared bool
}

func (csv *COWSliceValue) Set(index int, value Value) *COWSliceValue {
    if csv.shared {
        // 执行写时复制
        newData := make([]Value, len(csv.data))
        copy(newData, csv.data)
        csv = &COWSliceValue{data: newData, shared: false}
    }
    csv.data[index] = value
    return csv
}
```

通过这些全面的性能优化策略，无反射的 expr 实现可以在编译时间、执行速度和内存使用方面都超越原版实现。 