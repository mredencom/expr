# Lexer模块 - 词法分析器

## 概述

Lexer模块是表达式引擎的词法分析层，负责将输入的表达式字符串分解为标记（Token）序列。它是整个编译流水线的第一个阶段，为后续的语法分析提供规范化的输入。

## 核心功能

### 1. 标记化（Tokenization）
- 将输入字符串转换为标记序列
- 支持UTF-8编码的Unicode字符
- 准确的位置跟踪（行号、列号、偏移量）

### 2. 支持的标记类型
- **算术操作符**: `+`, `-`, `*`, `/`, `%`, `**`（幂运算）
- **比较操作符**: `==`, `!=`, `<`, `<=`, `>`, `>=`
- **逻辑操作符**: `&&`, `||`, `!`
- **位运算操作符**: `&`, `|`, `^`, `~`, `<<`, `>>`
- **分隔符**: `(`, `)`, `[`, `]`, `{`, `}`, `,`, `.`, `;`, `:`
- **特殊标记**: `=>` (Lambda箭头), `|` (管道), `?` (三元运算符), `#` (管道占位符)
- **字面量**: 整数、浮点数、字符串、布尔值
- **标识符**: 变量名、函数名
- **关键字**: `true`, `false`, `nil`, `if`, `else`, `in`

## 主要类型

### Lexer结构体
```go
type Lexer struct {
    input     string  // 输入字符串
    position  int     // 当前位置
    readPos   int     // 读取位置
    char      rune    // 当前字符
    line      int     // 行号（1-based）
    column    int     // 列号（1-based）
    lineStart int     // 当前行起始位置
}
```

### Token结构体
```go
type Token struct {
    Type     TokenType  // 标记类型
    Value    string     // 标记值
    Position Position   // 位置信息
}
```

### Position结构体
```go
type Position struct {
    Line   int  // 行号
    Column int  // 列号
    Offset int  // 字符偏移量
}
```

## 基本使用

### 1. 创建词法分析器
```go
package main

import (
    "fmt"
    "github.com/mredencom/expr/lexer"
)

func main() {
    input := "age > 18 && name == 'Alice'"
    l := lexer.New(input)
    
    // 逐个获取标记
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        fmt.Printf("Type: %s, Value: %s, Position: %d:%d\n", 
            token.Type, token.Value, token.Position.Line, token.Position.Column)
    }
}
```

### 2. 标记序列解析
```go
func tokenizeExpression(expression string) []lexer.Token {
    l := lexer.New(expression)
    var tokens []lexer.Token
    
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        tokens = append(tokens, token)
    }
    
    return tokens
}
```

## 高级特性

### 1. 字符串处理
```go
// 推荐使用单引号字符串（避免转义）
input1 := `name == 'Alice'`
input2 := `message == 'It\'s a great day!'`

// 双引号字符串（需要转义）
input3 := `name == "Bob"`
input4 := `message == "He said \"Hello\""`

// 混合使用
input5 := `text == 'No "escaping" needed here'`
```

### 2. 数字处理
```go
// 整数
input1 := "age == 25"

// 浮点数
input2 := "price == 19.99"
input3 := "ratio == .5"  // 等同于 0.5
input4 := "large == 1e6" // 科学记数法
```

### 3. 通配符支持
```go
// 通配符在成员访问中的使用
input1 := "user.*"           // 访问user的所有属性
input2 := "*.field"          // 通配符对象的field属性
input3 := "data.*.name"      // 嵌套通配符访问

// 通配符在管道操作中的应用
input4 := "users | map(u => u.*)"  // 提取所有用户的所有属性
```

### 4. 标识符和关键字
```go
// 普通标识符
input1 := "userName > minLength"

// 包含Unicode字符的标识符
input2 := "用户名 == '张三'"

// 关键字识别
input3 := "isValid == true && status != nil"
```

### 5. 复杂表达式
```go
// Lambda表达式
input1 := "users | filter(x => x.age > 18)"

// 管道操作
input2 := "data | map(transform) | filter(condition)"

// 管道占位符
input2_5 := "numbers | filter(# > 3) | map(# * 2)"

// 三元运算符
input3 := "score >= 60 ? 'Pass' : 'Fail'"
```

## 错误处理

### 1. 位置追踪
```go
func analyzeWithPositions(expression string) {
    l := lexer.New(expression)
    
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        
        if token.Type == lexer.ILLEGAL {
            fmt.Printf("错误标记: %s, 位置: %d:%d\n", 
                token.Value, token.Position.Line, token.Position.Column)
        }
    }
}
```

### 2. 重置和重用
```go
func reuseAnalyzer() {
    l := lexer.New("")
    
    expressions := []string{
        "x + y",
        "a * b",
        "name == \"test\"",
    }
    
    for _, expr := range expressions {
        l.Reset(expr)  // 重置到新的输入
        
        for {
            token := l.NextToken()
            if token.Type == lexer.EOF {
                break
            }
            fmt.Printf("%s ", token.Value)
        }
        fmt.Println()
    }
}
```

## 性能优化

### 1. 内存效率
- 使用UTF-8解码，支持多语言字符
- 零拷贝字符串处理
- 位置信息的高效计算

### 2. 处理速度
```go
// 批量标记化的高效方式
func tokenizeBatch(expressions []string) [][]lexer.Token {
    l := lexer.New("")
    results := make([][]lexer.Token, len(expressions))
    
    for i, expr := range expressions {
        l.Reset(expr)
        var tokens []lexer.Token
        
        for {
            token := l.NextToken()
            if token.Type == lexer.EOF {
                break
            }
            tokens = append(tokens, token)
        }
        
        results[i] = tokens
    }
    
    return results
}
```

## 调试和诊断

### 1. 标记流可视化
```go
func visualizeTokens(expression string) {
    l := lexer.New(expression)
    fmt.Printf("表达式: %s\n", expression)
    fmt.Println("标记流:")
    
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        
        fmt.Printf("  [%s] '%s' @%d:%d\n",
            token.Type, token.Value, 
            token.Position.Line, token.Position.Column)
    }
}
```

### 2. 性能分析
```go
func benchmarkTokenization(expression string, iterations int) {
    start := time.Now()
    
    for i := 0; i < iterations; i++ {
        l := lexer.New(expression)
        for {
            token := l.NextToken()
            if token.Type == lexer.EOF {
                break
            }
        }
    }
    
    duration := time.Since(start)
    fmt.Printf("处理 %d 次，总时间: %v，平均: %v\n", 
        iterations, duration, duration/time.Duration(iterations))
}
```

## 常见使用模式

### 1. 预检查模式
```go
func validateExpression(expression string) error {
    l := lexer.New(expression)
    
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        
        if token.Type == lexer.ILLEGAL {
            return fmt.Errorf("非法字符 '%s' 在位置 %d:%d", 
                token.Value, token.Position.Line, token.Position.Column)
        }
    }
    
    return nil
}
```

### 2. 统计分析模式
```go
func analyzeTokenStats(expression string) map[lexer.TokenType]int {
    l := lexer.New(expression)
    stats := make(map[lexer.TokenType]int)
    
    for {
        token := l.NextToken()
        if token.Type == lexer.EOF {
            break
        }
        stats[token.Type]++
    }
    
    return stats
}
```

## 最佳实践

1. **错误处理**: 始终检查ILLEGAL标记类型
2. **性能考虑**: 对于重复解析，考虑重用Lexer实例
3. **内存管理**: 及时处理完毕的标记，避免累积
4. **调试支持**: 利用位置信息提供精确的错误定位
5. **编码支持**: 确保输入字符串是有效的UTF-8编码

## 与其他模块的集成

Lexer模块是整个编译流水线的起点：
```
输入字符串 → Lexer → Token序列 → Parser → AST → Compiler → 字节码 → VM → 结果
```

通常情况下，用户不需要直接使用Lexer模块，而是通过更高层的API（如parser或expr包）来间接使用。但在需要精细控制或调试的场景下，直接使用Lexer模块可以提供更好的灵活性。 