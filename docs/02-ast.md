# AST模块 - 抽象语法树

## 概述

AST（Abstract Syntax Tree，抽象语法树）模块定义了表达式引擎的语法结构表示。它提供了一套完整的节点类型，用于表示各种表达式、语句和声明的语法结构，是连接词法分析器和语义分析器的重要桥梁。

## 核心功能

### 1. 语法结构表示
- 定义完整的表达式语法树节点
- 支持复杂的嵌套表达式结构
- 提供类型安全的节点访问接口

### 2. 节点访问模式
- 实现访问者模式（Visitor Pattern）
- 支持树遍历和转换操作
- 提供节点类型检查和转换

### 3. 调试支持
- 节点位置信息跟踪
- 语法树可视化支持
- 结构化输出和调试

## 核心接口

### Node接口
```go
type Node interface {
    String() string
    TokenLiteral() string
    GetPosition() *lexer.Position
}
```

### Statement接口
```go
type Statement interface {
    Node
    statementNode()
}
```

### Expression接口
```go
type Expression interface {
    Node
    expressionNode()
}
```

## 主要节点类型

### 1. 基础表达式节点

#### 字面量节点
```go
// 整数字面量
type IntegerLiteral struct {
    Token *lexer.Token
    Value int64
}

// 浮点数字面量
type FloatLiteral struct {
    Token *lexer.Token
    Value float64
}

// 字符串字面量
type StringLiteral struct {
    Token *lexer.Token
    Value string
}

// 布尔字面量
type BooleanLiteral struct {
    Token *lexer.Token
    Value bool
}
```

#### 标识符节点
```go
type Identifier struct {
    Token *lexer.Token
    Value string
}
```

### 2. 运算表达式节点

#### 中缀表达式（二元运算）
```go
type InfixExpression struct {
    Token    *lexer.Token // 操作符标记
    Left     Expression   // 左操作数
    Operator string       // 操作符
    Right    Expression   // 右操作数
}
```

#### 前缀表达式（一元运算）
```go
type PrefixExpression struct {
    Token    *lexer.Token // 操作符标记
    Operator string       // 操作符 (!、-等)
    Right    Expression   // 操作数
}
```

### 3. 复杂表达式节点

#### 数组字面量
```go
type ArrayLiteral struct {
    Token    *lexer.Token // [
    Elements []Expression // 数组元素
}
```

#### 索引表达式
```go
type IndexExpression struct {
    Token *lexer.Token // [
    Left  Expression   // 被索引的表达式
    Index Expression   // 索引表达式
}
```

#### 函数调用表达式
```go
type CallExpression struct {
    Token     *lexer.Token  // (
    Function  Expression    // 函数表达式
    Arguments []Expression  // 参数列表
}
```

#### Lambda表达式
```go
type LambdaExpression struct {
    Token      *lexer.Token    // => 或 lambda关键字
    Parameters []*Identifier  // 参数列表
    Body       Expression     // 函数体表达式
}
```

#### 管道表达式
```go
type PipeExpression struct {
    Token *lexer.Token // |
    Left  Expression   // 管道左侧表达式
    Right Expression   // 管道右侧表达式
}
```

#### 占位符表达式
```go
type PlaceholderExpression struct {
    Token *lexer.Token // #
}
```

#### 三元运算表达式
```go
type ConditionalExpression struct {
    Token     *lexer.Token // ?
    Condition Expression   // 条件表达式
    TrueExpr  Expression   // 真值表达式
    FalseExpr Expression   // 假值表达式
}
```

### 4. 访问表达式节点

#### 成员访问
```go
type MemberExpression struct {
    Token    *lexer.Token // .
    Object   Expression   // 被访问的对象
    Property *Identifier  // 属性名
}
```

#### 计算属性访问
```go
type ComputedMemberExpression struct {
    Token    *lexer.Token // [
    Object   Expression   // 被访问的对象
    Property Expression   // 动态属性表达式
}
```

## 基本使用

### 1. 创建简单表达式节点
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/ast"
    "github.com/mredencom/expr/lexer"
)

func main() {
    // 创建整数字面量节点
    intToken := &lexer.Token{Type: lexer.INT, Value: "42"}
    intLiteral := &ast.IntegerLiteral{
        Token: intToken,
        Value: 42,
    }
    
    // 创建标识符节点
    identToken := &lexer.Token{Type: lexer.IDENT, Value: "age"}
    identifier := &ast.Identifier{
        Token: identToken,
        Value: "age",
    }
    
    fmt.Printf("整数: %s\n", intLiteral.String())
    fmt.Printf("标识符: %s\n", identifier.String())
}
```

### 2. 构建复杂表达式
```go
func buildComplexExpression() ast.Expression {
    // age > 18
    left := &ast.Identifier{
        Token: &lexer.Token{Type: lexer.IDENT, Value: "age"},
        Value: "age",
    }
    
    right := &ast.IntegerLiteral{
        Token: &lexer.Token{Type: lexer.INT, Value: "18"},
        Value: 18,
    }
    
    return &ast.InfixExpression{
        Token:    &lexer.Token{Type: lexer.GT, Value: ">"},
        Left:     left,
        Operator: ">",
        Right:    right,
    }
}
```

### 3. 构建Lambda表达式
```go
func buildLambdaExpression() ast.Expression {
    // x => x * 2
    param := &ast.Identifier{
        Token: &lexer.Token{Type: lexer.IDENT, Value: "x"},
        Value: "x",
    }
    
    body := &ast.InfixExpression{
        Token: &lexer.Token{Type: lexer.MUL, Value: "*"},
        Left: &ast.Identifier{
            Token: &lexer.Token{Type: lexer.IDENT, Value: "x"},
            Value: "x",
        },
        Operator: "*",
        Right: &ast.IntegerLiteral{
            Token: &lexer.Token{Type: lexer.INT, Value: "2"},
            Value: 2,
        },
    }
    
    return &ast.LambdaExpression{
        Token:      &lexer.Token{Type: lexer.ARROW, Value: "=>"},
        Parameters: []*ast.Identifier{param},
        Body:       body,
    }
}
```

## 最佳实践

1. **类型安全**: 使用类型断言时进行充分的检查
2. **内存管理**: 对于频繁创建的节点，考虑使用对象池
3. **位置跟踪**: 保持准确的位置信息，便于错误报告
4. **不可变性**: 尽量保持AST节点的不可变性，需要修改时创建新节点
5. **访问者模式**: 对于复杂的AST操作，使用访问者模式保持代码清晰

## 与其他模块的集成

AST模块在编译流水线中的位置：
```
Lexer → Token序列 → Parser → AST → Checker → 类型化AST → Compiler → 字节码
```

AST模块为上层模块提供了结构化的语法表示，是语法分析和语义分析的基础。通过合理使用AST模块，可以实现强大的表达式分析、转换和优化功能。 