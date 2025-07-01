# 管道占位符完整指南

## 概述

管道占位符语法是表达式引擎的高级特性，使用 `#` 符号作为管道操作中当前元素的占位符。这种语法提供了比传统Lambda表达式更简洁、更直观的数据处理方式。

## 核心概念

### 1. 占位符符号 `#`
- `#` 代表管道中当前正在处理的元素
- 只能在管道操作的右侧使用
- 支持复杂的表达式和嵌套运算

### 2. 管道操作符 `|`
- 将左侧数据传递给右侧函数
- 支持链式操作
- 自动处理数据流转换

## 基础语法

### 1. 简单过滤
```go
// 基础数字过滤
"[1, 2, 3, 4, 5] | filter(# > 3)"        // [4, 5]
"[1, 2, 3, 4, 5] | filter(# % 2 == 0)"   // [2, 4]

// 字符串过滤
"['apple', 'banana', 'cherry'] | filter(len(#) > 5)"  // ['banana', 'cherry']
```

### 2. 简单映射
```go
// 数值变换
"[1, 2, 3, 4, 5] | map(# * 2)"           // [2, 4, 6, 8, 10]
"[1, 2, 3, 4, 5] | map(# * # + 1)"       // [2, 5, 10, 17, 26]

// 字符串变换
"['hello', 'world'] | map(upper(#))"     // ['HELLO', 'WORLD']
```

### 3. 链式操作
```go
// 多步骤处理
"[1, 2, 3, 4, 5, 6, 7, 8, 9, 10] | filter(# > 5) | map(# * 2)"
// 结果: [12, 14, 16, 18, 20]

// 复杂数据流
"data | filter(# % 2 == 0) | map(# * # + 1) | filter(# > 10)"
```

## 高级用法

### 1. 复杂表达式
```go
// 算术表达式
"numbers | filter(# * 2 + 1 > 10)"       // 复合算术条件
"numbers | map((# + 1) * (# - 1))"       // 嵌套括号运算

// 比较表达式
"numbers | filter(# >= 3 && # <= 7)"     // 范围过滤
"numbers | map(# > 5 ? # * 2 : #)"       // 条件映射

// 模运算
"numbers | filter(# % 3 == 0 || # % 5 == 0)"  // 3或5的倍数
```

### 2. 对象属性访问
```go
// 对象数组处理
users := []map[string]interface{}{
    {"name": "Alice", "age": 30, "salary": 75000},
    {"name": "Bob", "age": 25, "salary": 65000},
    {"name": "Charlie", "age": 35, "salary": 85000},
}

// 属性过滤
"users | filter(#.age > 30)"             // 年龄大于30的用户
"users | filter(#.salary >= 70000)"      // 高薪用户

// 属性提取
"users | map(#.name)"                    // 提取所有用户名
"users | map(#.salary)"                  // 提取所有薪资

// 复合条件
"users | filter(#.age >= 30 && #.salary > 70000) | map(#.name)"
```

### 3. 嵌套数据处理
```go
// 多层对象
departments := []map[string]interface{}{
    {
        "name": "Engineering",
        "employees": []map[string]interface{}{
            {"name": "Alice", "level": "Senior"},
            {"name": "Bob", "level": "Junior"},
        },
    },
}

// 嵌套访问
"departments | map(#.employees) | reduce((acc, emp) => acc + emp, [])"
"departments | filter(#.name == 'Engineering') | map(#.employees) | map(#.name)"
```

## 实际应用示例

### 1. 数据分析
```go
package main

import (
    "fmt"
    expr "github.com/mredencom/expr"
)

func analyzeData() {
    sales := map[string]interface{}{
        "transactions": []map[string]interface{}{
            {"amount": 120.50, "category": "electronics", "date": "2024-01-15"},
            {"amount": 89.99, "category": "books", "date": "2024-01-16"},
            {"amount": 250.00, "category": "electronics", "date": "2024-01-17"},
            {"amount": 45.00, "category": "books", "date": "2024-01-18"},
        },
    }

    // 高价值交易
    highValue, _ := expr.Eval(`
        transactions 
        | filter(#.amount > 100) 
        | map(#.amount)
    `, sales)
    fmt.Printf("高价值交易: %v\n", highValue)

    // 电子产品销售总额
    electronicsTotal, _ := expr.Eval(`
        transactions 
        | filter(#.category == "electronics") 
        | map(#.amount) 
        | sum()
    `, sales)
    fmt.Printf("电子产品总销售额: %v\n", electronicsTotal)

    // 平均交易金额
    avgAmount, _ := expr.Eval(`
        transactions 
        | map(#.amount) 
        | avg()
    `, sales)
    fmt.Printf("平均交易金额: %v\n", avgAmount)
}
```

### 2. 业务规则引擎
```go
func businessRules() {
    customers := map[string]interface{}{
        "customers": []map[string]interface{}{
            {"name": "Alice", "age": 30, "orders": 15, "totalSpent": 2500},
            {"name": "Bob", "age": 25, "orders": 8, "totalSpent": 1200},
            {"name": "Charlie", "age": 35, "orders": 25, "totalSpent": 4500},
        },
    }

    // VIP客户识别
    vipCustomers, _ := expr.Eval(`
        customers 
        | filter(#.totalSpent > 2000 && #.orders > 10) 
        | map(#.name)
    `, customers)
    fmt.Printf("VIP客户: %v\n", vipCustomers)

    // 潜在流失客户
    riskCustomers, _ := expr.Eval(`
        customers 
        | filter(#.orders < 10 && #.age > 30) 
        | map({name: #.name, risk: "high"})
    `, customers)
    fmt.Printf("风险客户: %v\n", riskCustomers)
}
```

### 3. 配置驱动处理
```go
func configDriven() {
    config := map[string]interface{}{
        "rules": []map[string]interface{}{
            {"field": "age", "operator": ">", "value": 18, "action": "allow"},
            {"field": "score", "operator": ">=", "value": 80, "action": "promote"},
            {"field": "status", "operator": "==", "value": "active", "action": "process"},
        },
        "users": []map[string]interface{}{
            {"name": "Alice", "age": 25, "score": 85, "status": "active"},
            {"name": "Bob", "age": 17, "score": 90, "status": "inactive"},
            {"name": "Charlie", "age": 30, "score": 75, "status": "active"},
        },
    }

    // 应用年龄规则
    adults, _ := expr.Eval("users | filter(#.age > 18)", config)
    fmt.Printf("成年用户: %v\n", adults)

    // 应用综合规则
    qualified, _ := expr.Eval(`
        users 
        | filter(#.age > 18 && #.score >= 80 && #.status == "active")
        | map(#.name)
    `, config)
    fmt.Printf("合格用户: %v\n", qualified)
}
```

## 性能优化

### 1. 预编译优化
```go
// 预编译复用
func optimizedProcessing() {
    // 编译一次，多次使用
    program, err := expr.Compile("data | filter(# > threshold) | map(# * multiplier)")
    if err != nil {
        panic(err)
    }

    // 多次执行不同数据
    datasets := []map[string]interface{}{
        {"data": []int{1, 2, 3, 4, 5}, "threshold": 2, "multiplier": 3},
        {"data": []int{6, 7, 8, 9, 10}, "threshold": 7, "multiplier": 2},
    }

    for _, env := range datasets {
        result, _ := expr.Run(program, env)
        fmt.Printf("结果: %v\n", result)
    }
}
```

### 2. 批量处理
```go
// 批量数据处理
func batchProcessing() {
    largeDataset := make([]int, 100000)
    for i := range largeDataset {
        largeDataset[i] = i + 1
    }

    env := map[string]interface{}{"data": largeDataset}

    // 高效的管道处理
    result, _ := expr.Eval(`
        data 
        | filter(# % 1000 == 0)  // 每1000个取一个
        | map(# / 1000)          // 转换为千位数
        | filter(# <= 50)        // 限制范围
    `, env)

    fmt.Printf("批量处理结果: %v\n", result)
}
```

## 错误处理和调试

### 1. 常见错误
```go
// ❌ 错误：在非管道上下文中使用占位符
// expr.Eval("# > 5", nil)  // 错误

// ✅ 正确：在管道上下文中使用
// expr.Eval("data | filter(# > 5)", env)

// ❌ 错误：类型不匹配
// expr.Eval("'hello' | filter(# > 5)", nil)  // 字符串不能用于数值比较

// ✅ 正确：类型匹配
// expr.Eval("'hello' | map(len(#))", nil)  // 获取字符串长度
```

### 2. 调试技巧
```go
func debugging() {
    env := map[string]interface{}{
        "data": []int{1, 2, 3, 4, 5},
    }

    // 分步调试
    step1, _ := expr.Eval("data | filter(# > 2)", env)
    fmt.Printf("步骤1: %v\n", step1)

    step2, _ := expr.Eval("data | filter(# > 2) | map(# * 2)", env)
    fmt.Printf("步骤2: %v\n", step2)

    // 添加调试输出
    debug, _ := expr.Eval(`
        data 
        | filter(# > 2) 
        | map(# * 2) 
        | filter(# < 10)
    `, env)
    fmt.Printf("最终结果: %v\n", debug)
}
```

## 最佳实践

### 1. 代码可读性
```go
// ✅ 推荐：简洁清晰
"users | filter(#.active) | map(#.email)"

// ❌ 不推荐：过度复杂
"users | filter(#.active && #.age > 18 && #.score > 80 && #.premium) | map(#.email)"

// ✅ 推荐：分步处理
`users 
 | filter(#.active) 
 | filter(#.age > 18) 
 | filter(#.score > 80) 
 | filter(#.premium) 
 | map(#.email)`
```

### 2. 性能考虑
```go
// ✅ 推荐：先过滤后映射
"data | filter(# > 100) | map(# * 2)"

// ❌ 不推荐：先映射后过滤（处理更多数据）
"data | map(# * 2) | filter(# > 200)"

// ✅ 推荐：预编译重用
program, _ := expr.Compile("data | filter(# > threshold)")
```

### 3. 类型安全
```go
// ✅ 推荐：明确类型期望
"numbers | filter(# > 0) | map(# * 1.5)"  // 数值处理

// ✅ 推荐：对象属性访问
"users | filter(#.age > 18) | map(#.name)"  // 对象处理

// ❌ 避免：混合类型处理
"mixed | filter(# > 0)"  // 如果mixed包含字符串会出错
```

## 与其他语法的对比

### 1. vs Lambda表达式
```go
// Lambda语法（传统）
"filter(numbers, x => x > 5)"
"map(filter(numbers, x => x % 2 == 0), y => y * 2)"

// 占位符语法（推荐）
"numbers | filter(# > 5)"
"numbers | filter(# % 2 == 0) | map(# * 2)"
```

### 2. vs 函数调用
```go
// 函数调用语法
"map(filter(sort(numbers), x => x > 5), y => y * 2)"

// 管道语法（更清晰）
"numbers | sort() | filter(# > 5) | map(# * 2)"
```

## 总结

管道占位符语法提供了：

1. **简洁性**: 减少冗余的参数声明
2. **可读性**: 清晰的数据流向
3. **可组合性**: 易于构建复杂的数据处理管道
4. **性能**: 优化的执行路径
5. **类型安全**: 编译时类型检查

这种语法特别适合于数据密集型应用、业务规则引擎和实时数据处理场景。 