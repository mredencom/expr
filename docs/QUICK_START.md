# 快速开始指南

## 🚀 5分钟上手 Expr

### 安装

```bash
go get github.com/mredencom/expr
```

### 第一个表达式

```go
package main

import (
    "fmt"
    expr "github.com/mredencom/expr"
)

func main() {
    // 简单计算
    result, _ := expr.Eval("2 + 3 * 4", nil)
    fmt.Println(result) // 输出: 14
}
```

### 使用变量

```go
func main() {
    env := map[string]interface{}{
        "name": "世界",
        "count": 42,
    }
    
    result, _ := expr.Eval("'你好, ' + name + '! 答案是 ' + toString(count)", env)
    fmt.Println(result) // 输出: "你好, 世界! 答案是 42"
}
```

### 数据处理

```go
func main() {
    users := []map[string]interface{}{
        {"name": "Alice", "age": 25, "active": true},
        {"name": "Bob", "age": 16, "active": false},
        {"name": "Charlie", "age": 30, "active": true},
    }
    
    env := map[string]interface{}{"users": users}
    
    // 过滤成年活跃用户
    result, _ := expr.Eval("users | filter(u => u.age >= 18 && u.active)", env)
    fmt.Printf("成年活跃用户: %+v\n", result)
}
```

### 占位符语法

```go
func main() {
    numbers := []int{1, 6, 3, 8, 2, 9}
    env := map[string]interface{}{"numbers": numbers}
    
    // 过滤和映射
    result, _ := expr.Eval("numbers | filter(# > 5) | map(# * 2)", env)
    fmt.Println(result) // 输出: [12, 16, 18]
}
```

### 空值安全

```go
func main() {
    data := map[string]interface{}{
        "user": map[string]interface{}{
            "profile": map[string]interface{}{
                "name": "Alice",
            },
        },
        "guest": nil,
    }
    
    // 安全访问
    result, _ := expr.Eval("user?.profile?.name ?? 'Unknown'", data)
    fmt.Println(result) // 输出: "Alice"
    
    result, _ = expr.Eval("guest?.profile?.name ?? 'Guest'", data)
    fmt.Println(result) // 输出: "Guest"
}
```

### 预编译优化

```go
func main() {
    // 预编译表达式（推荐用于重复执行）
    program, err := expr.Compile("price * (1 - discount) * quantity")
    if err != nil {
        panic(err)
    }
    
    // 多次执行
    orders := []map[string]interface{}{
        {"price": 100.0, "discount": 0.1, "quantity": 2},
        {"price": 50.0, "discount": 0.2, "quantity": 3},
    }
    
    for i, order := range orders {
        total, _ := expr.Run(program, order)
        fmt.Printf("订单%d总价: %.2f\n", i+1, total)
    }
}
```

## 📚 下一步

- 阅读 [API文档](API.md) 了解完整功能
- 查看 [示例代码](EXAMPLES.md) 学习更多用法
- 参考 [最佳实践](BEST_PRACTICES.md) 用于生产环境
- 使用 [调试指南](DEBUGGING.md) 解决问题

## 🎯 常用场景

### 业务规则
```go
expr.Eval("customer.level == 'VIP' && order.amount > 1000", data)
```

### 数据转换
```go
expr.Eval("users | map(u => {name: u.firstName + ' ' + u.lastName, adult: u.age >= 18})", data)
```

### 计算公式
```go
expr.Eval("math.sqrt(a * a + b * b)", data)
```

开始您的Expr之旅吧！🚀 