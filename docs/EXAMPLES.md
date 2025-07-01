# 示例代码

## 🚀 快速开始示例

### 基础表达式

```go
package main

import (
    "fmt"
    expr "github.com/mredencom/expr"
)

func main() {
    // 简单算术
    result, _ := expr.Eval("2 + 3 * 4", nil)
    fmt.Println(result) // 14

    // 字符串操作
    result, _ = expr.Eval("upper('hello') + ' WORLD'", nil)
    fmt.Println(result) // "HELLO WORLD"

    // 使用变量
    env := map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }
    result, _ = expr.Eval("name + ' is ' + toString(age) + ' years old'", env)
    fmt.Println(result) // "Alice is 30 years old"
}
```

### 预编译使用

```go
func main() {
    // 预编译表达式
    program, err := expr.Compile("price * (1 - discount) * quantity")
    if err != nil {
        panic(err)
    }

    // 多次执行
    orders := []map[string]interface{}{
        {"price": 100.0, "discount": 0.1, "quantity": 2},
        {"price": 50.0, "discount": 0.2, "quantity": 3},
        {"price": 200.0, "discount": 0.15, "quantity": 1},
    }

    for i, order := range orders {
        total, _ := expr.Run(program, order)
        fmt.Printf("订单%d总价: %.2f\n", i+1, total)
    }
}
```

## 🔧 Lambda表达式示例

### 数据过滤和映射

```go
func main() {
    users := []map[string]interface{}{
        {"name": "Alice", "age": 25, "active": true},
        {"name": "Bob", "age": 16, "active": false},
        {"name": "Charlie", "age": 30, "active": true},
        {"name": "David", "age": 22, "active": true},
    }

    env := map[string]interface{}{"users": users}

    // 过滤成年且活跃的用户
    result, _ := expr.Eval(`
        users 
        | filter(u => u.age >= 18 && u.active) 
        | map(u => u.name)
    `, env)
    fmt.Println("活跃成年用户:", result) // ["Alice", "Charlie", "David"]

    // 计算平均年龄
    result, _ = expr.Eval(`
        users 
        | filter(u => u.active) 
        | map(u => u.age) 
        | avg()
    `, env)
    fmt.Println("活跃用户平均年龄:", result) // 25.666...

    // 复杂的数据转换
    result, _ = expr.Eval(`
        users 
        | filter(u => u.age >= 18)
        | map(u => {
            name: u.name,
            category: u.age >= 25 ? "senior" : "junior",
            status: u.active ? "active" : "inactive"
        })
    `, env)
    fmt.Printf("用户分类: %+v\n", result)
}
```

### 排序和分组

```go
func main() {
    products := []map[string]interface{}{
        {"name": "Laptop", "price": 1200, "category": "electronics"},
        {"name": "Book", "price": 20, "category": "education"},
        {"name": "Phone", "price": 800, "category": "electronics"},
        {"name": "Pen", "price": 5, "category": "office"},
    }

    env := map[string]interface{}{"products": products}

    // 按价格排序
    result, _ := expr.Eval(`
        products 
        | sort((a, b) => a.price - b.price) 
        | map(p => p.name + ": $" + toString(p.price))
    `, env)
    fmt.Println("按价格排序:", result)

    // 获取最贵的电子产品
    result, _ = expr.Eval(`
        products 
        | filter(p => p.category == "electronics")
        | sort((a, b) => b.price - a.price)
        | first()
    `, env)
    fmt.Printf("最贵的电子产品: %+v\n", result)
}
```

## ⚡ 管道和占位符示例

### 占位符操作

```go
func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    env := map[string]interface{}{"numbers": numbers}

    // 链式占位符操作
    result, _ := expr.Eval("numbers | filter(# > 5) | map(# * 2)", env)
    fmt.Println("大于5的数乘2:", result) // [12, 14, 16, 18, 20]

    // 复杂占位符表达式
    result, _ = expr.Eval("numbers | filter(# % 2 == 0) | map(# * # + 1)", env)
    fmt.Println("偶数平方加1:", result) // [5, 17, 37, 65, 101]

    // 条件占位符
    result, _ = expr.Eval("numbers | map(# > 5 ? # * 10 : # * 2)", env)
    fmt.Println("条件映射:", result) // [2, 4, 6, 8, 10, 60, 70, 80, 90, 100]

    // 聚合操作
    result, _ = expr.Eval("numbers | filter(# >= 3 && # <= 7) | sum()", env)
    fmt.Println("3-7的和:", result) // 25
}
```

### 字符串处理管道

```go
func main() {
    words := []string{"hello", "world", "go", "programming", "is", "fun"}
    env := map[string]interface{}{"words": words}

    // 字符串处理管道
    result, _ := expr.Eval(`
        words 
        | filter(length(#) > 3) 
        | map(upper(#)) 
        | sort()
    `, env)
    fmt.Println("长单词大写排序:", result) // ["HELLO", "PROGRAMMING", "WORLD"]

    // 使用字符串方法
    result, _ = expr.Eval(`
        words 
        | map(#.upper().length()) 
        | filter(# > 4)
    `, env)
    fmt.Println("长单词的长度:", result) // [5, 11]
}
```

## 🛡️ 空值安全示例

### 安全属性访问

```go
func main() {
    data := map[string]interface{}{
        "user": map[string]interface{}{
            "profile": map[string]interface{}{
                "name": "Alice",
                "bio":  "Software Developer",
            },
            "settings": map[string]interface{}{
                "theme": "dark",
            },
        },
        "admin": map[string]interface{}{
            "profile": nil,
        },
        "guest": nil,
    }

    // 安全访问存在的属性
    result, _ := expr.Eval("user?.profile?.name ?? 'Unknown'", data)
    fmt.Println("用户名:", result) // "Alice"

    // 安全访问不存在的属性
    result, _ = expr.Eval("admin?.profile?.name ?? 'No Name'", data)
    fmt.Println("管理员名:", result) // "No Name"

    // 安全访问null对象
    result, _ = expr.Eval("guest?.profile?.name ?? 'Guest User'", data)
    fmt.Println("访客名:", result) // "Guest User"

    // 复杂的空值处理
    result, _ = expr.Eval(`
        user?.profile?.bio ?? user?.profile?.description ?? "No description available"
    `, data)
    fmt.Println("用户简介:", result) // "Software Developer"
}
```

### 数组安全访问

```go
func main() {
    data := map[string]interface{}{
        "users": []map[string]interface{}{
            {"name": "Alice", "emails": []string{"alice@example.com"}},
            {"name": "Bob", "emails": nil},
            {"name": "Charlie"},
        },
        "emptyList": []interface{}{},
        "nullList":  nil,
    }

    // 安全访问数组元素
    result, _ := expr.Eval("users?.[0]?.name ?? 'No user'", data)
    fmt.Println("第一个用户:", result) // "Alice"

    // 安全访问嵌套数组
    result, _ = expr.Eval("users?.[0]?.emails?.[0] ?? 'No email'", data)
    fmt.Println("第一个用户邮箱:", result) // "alice@example.com"

    result, _ = expr.Eval("users?.[1]?.emails?.[0] ?? 'No email'", data)
    fmt.Println("第二个用户邮箱:", result) // "No email"

    // 空列表安全访问
    result, _ = expr.Eval("emptyList?.[0] ?? 'Empty'", data)
    fmt.Println("空列表访问:", result) // "Empty"

    result, _ = expr.Eval("nullList?.[0] ?? 'Null list'", data)
    fmt.Println("空引用访问:", result) // "Null list"
}
```

## 📦 模块系统示例

### Math模块使用

```go
func main() {
    // 基础数学计算
    result, _ := expr.Eval("math.sqrt(16) + math.pow(2, 3)", nil)
    fmt.Println("数学计算:", result) // 12

    // 三角函数
    result, _ = expr.Eval("math.sin(math.pi / 2)", nil)
    fmt.Println("sin(π/2):", result) // 1

    // 在表达式中使用
    data := map[string]interface{}{
        "radius": 5.0,
    }
    result, _ = expr.Eval("math.pi * math.pow(radius, 2)", data)
    fmt.Printf("圆面积: %.2f\n", result) // 78.54

    // 数组中的数学运算
    numbers := []float64{1.2, 2.7, 3.1, 4.9}
    env := map[string]interface{}{"numbers": numbers}
    result, _ = expr.Eval("numbers | map(math.ceil(#))", env)
    fmt.Println("向上取整:", result) // [2, 3, 4, 5]
}
```

### Strings模块使用

```go
func main() {
    // 字符串处理
    result, _ := expr.Eval(`strings.upper("hello") + " " + strings.lower("WORLD")`, nil)
    fmt.Println("字符串操作:", result) // "HELLO world"

    // 字符串分割和连接
    result, _ = expr.Eval(`strings.split("a,b,c", ",") | map(strings.trim(#))`, nil)
    fmt.Println("分割处理:", result) // ["a", "b", "c"]

    // 在数据处理中使用
    names := []string{" Alice ", " Bob ", " Charlie "}
    env := map[string]interface{}{"names": names}
    result, _ = expr.Eval(`
        names 
        | map(strings.trim(#)) 
        | map(strings.upper(#)) 
        | filter(strings.hasPrefix(#, "A"))
    `, env)
    fmt.Println("以A开头的名字:", result) // ["ALICE"]
}
```

### 自定义模块

```go
func main() {
    // 注册自定义模块
    customFunctions := map[string]interface{}{
        "formatPrice": func(price float64) string {
            return fmt.Sprintf("$%.2f", price)
        },
        "isWeekend": func(day string) bool {
            return day == "Saturday" || day == "Sunday"
        },
        "calculateTax": func(amount float64, rate float64) float64 {
            return amount * rate
        },
    }
    
    // 这里需要模块注册的实际API
    // modules.RegisterModule("custom", customFunctions)

    // 使用自定义模块
    data := map[string]interface{}{
        "price": 99.99,
        "taxRate": 0.08,
        "today": "Saturday",
    }
    
    // result, _ := expr.Eval("custom.formatPrice(price + custom.calculateTax(price, taxRate))", data)
    // fmt.Println("含税价格:", result)
    
    // result, _ = expr.Eval("custom.isWeekend(today) ? 'Weekend!' : 'Weekday'", data)
    // fmt.Println("今天:", result)
}
```

## 🏢 企业级应用示例

### 业务规则引擎

```go
type BusinessRule struct {
    Name       string
    Expression string
    Priority   int
}

func main() {
    rules := []BusinessRule{
        {
            Name:       "VIP客户优惠",
            Expression: "customer.vipLevel >= 3 && order.amount > 1000",
            Priority:   1,
        },
        {
            Name:       "新客户优惠",
            Expression: "customer.isNew && order.amount > 100",
            Priority:   2,
        },
        {
            Name:       "批量订单优惠",
            Expression: "order.items | length() > 10",
            Priority:   3,
        },
    }

    // 编译规则
    compiledRules := make([]*expr.Program, len(rules))
    for i, rule := range rules {
        program, err := expr.Compile(rule.Expression, expr.AsBool())
        if err != nil {
            fmt.Printf("规则 %s 编译失败: %v\n", rule.Name, err)
            continue
        }
        compiledRules[i] = program
    }

    // 测试数据
    testData := map[string]interface{}{
        "customer": map[string]interface{}{
            "vipLevel": 4,
            "isNew":    false,
        },
        "order": map[string]interface{}{
            "amount": 1500,
            "items":  make([]interface{}, 15),
        },
    }

    // 执行规则
    for i, program := range compiledRules {
        if program == nil {
            continue
        }
        
        result, err := expr.Run(program, testData)
        if err != nil {
            fmt.Printf("规则 %s 执行失败: %v\n", rules[i].Name, err)
            continue
        }
        
        if result.(bool) {
            fmt.Printf("✅ 触发规则: %s\n", rules[i].Name)
        } else {
            fmt.Printf("❌ 未触发: %s\n", rules[i].Name)
        }
    }
}
```

### 配置驱动的数据处理

```go
type ProcessingConfig struct {
    Name       string `json:"name"`
    Expression string `json:"expression"`
    OutputType string `json:"output_type"`
}

func main() {
    configs := []ProcessingConfig{
        {
            Name:       "用户全名",
            Expression: "user.firstName + ' ' + user.lastName",
            OutputType: "string",
        },
        {
            Name:       "年龄分组",
            Expression: "user.age >= 18 ? 'adult' : 'minor'",
            OutputType: "string",
        },
        {
            Name:       "活跃度评分",
            Expression: "user.loginDays * 2 + user.posts * 5",
            OutputType: "int",
        },
    }

    // 编译所有配置
    programs := make(map[string]*expr.Program)
    for _, config := range configs {
        var options []expr.Option
        switch config.OutputType {
        case "string":
            options = append(options, expr.AsString())
        case "int":
            options = append(options, expr.AsInt())
        case "bool":
            options = append(options, expr.AsBool())
        }

        program, err := expr.Compile(config.Expression, options...)
        if err != nil {
            fmt.Printf("配置 %s 编译失败: %v\n", config.Name, err)
            continue
        }
        programs[config.Name] = program
    }

    // 处理用户数据
    users := []map[string]interface{}{
        {
            "firstName": "Alice",
            "lastName":  "Smith",
            "age":       25,
            "loginDays": 30,
            "posts":     12,
        },
        {
            "firstName": "Bob",
            "lastName":  "Jones",
            "age":       17,
            "loginDays": 15,
            "posts":     8,
        },
    }

    for i, user := range users {
        fmt.Printf("=== 用户 %d ===\n", i+1)
        env := map[string]interface{}{"user": user}

        for _, config := range configs {
            program := programs[config.Name]
            if program == nil {
                continue
            }

            result, err := expr.Run(program, env)
            if err != nil {
                fmt.Printf("%s: 错误 - %v\n", config.Name, err)
            } else {
                fmt.Printf("%s: %v\n", config.Name, result)
            }
        }
        fmt.Println()
    }
}
```

这些示例展示了Expr表达式引擎在各种场景下的强大功能和灵活性。 