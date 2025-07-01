# 表达式引擎 - 综合使用指南

## 概述

本文档提供表达式引擎的综合使用指南，整合了所有模块的最佳实践。该表达式引擎采用零反射架构，支持25M-39M ops/sec的极高性能，广泛适用于业务规则、动态配置、数据处理等场景。

## 快速开始

### 1. 基础安装和引入
```go
import "github.com/mredencom/expr"

// 最简单的使用
result, err := expr.Eval("2 + 3 * 4", nil)
fmt.Println(result) // 14
```

### 2. 带环境变量的基础使用
```go
// 用户数据
user := map[string]interface{}{
    "name": "Alice",
    "age":  30,
    "vip":  true,
}

// 基础表达式求值
isAdult, _ := expr.Eval("age >= 18", user)
greeting, _ := expr.Eval("'Hello ' + name + '!'", user)
discount, _ := expr.Eval("vip ? 0.2 : 0.1", user)

fmt.Printf("成年人: %v, 问候: %s, 折扣: %v\n", isAdult, greeting, discount)
```

### 3. 预编译性能优化
```go
// 编译一次，多次使用
program, err := expr.Compile("user.age >= minAge && user.active")
if err != nil {
    panic(err)
}

users := []User{...} // 大量用户数据
minAge := 18

for _, user := range users {
    env := map[string]interface{}{
        "user":   user,
        "minAge": minAge,
    }
    
    if result, _ := expr.Run(program, env); result.(bool) {
        fmt.Printf("符合条件的用户: %s\n", user.Name)
    }
}
```

## 核心特性详解

### 1. 支持的数据类型
```go
// 基础类型
numbers := expr.Eval("42", nil)           // int64
decimal := expr.Eval("3.14", nil)         // float64
text := expr.Eval("'hello'", nil)         // string
flag := expr.Eval("true", nil)            // bool
empty := expr.Eval("nil", nil)            // nil

// 集合类型
env := map[string]interface{}{
    "numbers": []int{1, 2, 3, 4, 5},
    "user": map[string]interface{}{
        "name": "Bob",
        "tags": []string{"admin", "user"},
    },
}

// 数组操作
firstNumber, _ := expr.Eval("numbers[0]", env)        // 1
userTags, _ := expr.Eval("user.tags", env)           // ["admin", "user"]
hasAdmin, _ := expr.Eval("'admin' in user.tags", env) // true

// 通配符访问
allUserProps, _ := expr.Eval("user.*", env)          // 获取user的所有属性
wildCardAccess, _ := expr.Eval("*.name", env)        // 通配符对象访问
```

### 2. 运算符支持
```go
env := map[string]interface{}{
    "a": 10, "b": 3, "x": true, "y": false,
    "name": "Alice", "age": 25,
}

// 算术运算
expr.Eval("a + b", env)     // 13 (加法)
expr.Eval("a - b", env)     // 7  (减法)
expr.Eval("a * b", env)     // 30 (乘法)
expr.Eval("a / b", env)     // 3  (除法)
expr.Eval("a % b", env)     // 1  (取模)
expr.Eval("a ** b", env)    // 1000 (幂运算)

// 比较运算
expr.Eval("age > 18", env)      // true
expr.Eval("age <= 30", env)     // true
expr.Eval("name == 'Alice'", env) // true

// 逻辑运算
expr.Eval("x && y", env)        // false
expr.Eval("x || y", env)        // true
expr.Eval("!x", env)           // false

// 字符串运算
expr.Eval("name + ' is ' + string(age)", env) // "Alice is 25"
```

### 3. 高级语法特性
```go
// 三元运算符
result, _ := expr.Eval("age >= 18 ? 'adult' : 'minor'", 
    map[string]interface{}{"age": 25}) // "adult"

// Lambda表达式
numbers := map[string]interface{}{
    "data": []int{1, 2, 3, 4, 5},
}
evens, _ := expr.Eval("filter(data, x => x % 2 == 0)", numbers) // [2, 4]

// 管道占位符语法（推荐）
pipeline1, _ := expr.Eval("data | filter(# % 2 == 0)", numbers) // [2, 4]
pipeline2, _ := expr.Eval("data | map(# * 2)", numbers)         // [2, 4, 6, 8, 10]

// 管道操作
pipeline, _ := expr.Eval(`
    data 
    | filter(x => x > 2) 
    | map(x => x * x) 
    | sum()
`, numbers) // 50 (3²+4²+5² = 9+16+25)

// 复杂对象访问
user := map[string]interface{}{
    "profile": map[string]interface{}{
        "settings": map[string]interface{}{
            "theme": "dark",
            "notifications": true,
        },
    },
}
theme, _ := expr.Eval("profile.settings.theme", user) // "dark"
```

## 管道占位符语法

### 1. 基础占位符用法
```go
numbers := map[string]interface{}{
    "data": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
}

// 基础过滤：获取大于5的数字
result1, _ := expr.Eval("data | filter(# > 5)", numbers)
// 结果: [6, 7, 8, 9, 10]

// 基础映射：每个数字乘以2
result2, _ := expr.Eval("data | map(# * 2)", numbers)
// 结果: [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

// 复合条件：偶数且大于3
result3, _ := expr.Eval("data | filter(# % 2 == 0 && # > 3)", numbers)
// 结果: [4, 6, 8, 10]
```

### 2. 复杂表达式
```go
// 数学运算组合
expr.Eval("data | filter(# % 2 == 0)", numbers)        // 偶数过滤
expr.Eval("data | map(# * 2 + 1)", numbers)           // 复杂算术变换
expr.Eval("data | filter((# + 1) * 2 > 10)", numbers) // 嵌套运算条件

// 比较运算
expr.Eval("data | filter(# >= 3 && # <= 7)", numbers) // 范围过滤
expr.Eval("data | map(# > 5 ? # * 2 : #)", numbers)   // 条件映射

// 模运算
expr.Eval("data | filter(# % 3 == 0)", numbers)       // 3的倍数
expr.Eval("data | map(# % 2)", numbers)               // 获取奇偶性
```

### 3. 管道链式操作
```go
// 多级管道处理
complexPipeline, _ := expr.Eval(`
    data 
    | filter(# > 3)           // 筛选大于3的数字
    | map(# * 2)              // 每个数字乘以2  
    | filter(# % 4 == 0)      // 筛选4的倍数
    | sum()                   // 求和
`, numbers)
// 结果: 60 (8 + 12 + 16 + 20 + 4)

// 数据转换链
transformChain, _ := expr.Eval(`
    data
    | filter(# % 2 == 1)      // 奇数
    | map(# * # + 1)          // 平方加1
    | filter(# > 10)          // 大于10
`, numbers)
// 结果: [26, 50, 82] (5²+1, 7²+1, 9²+1)
```

### 4. 对象数组处理
```go
users := map[string]interface{}{
    "people": []map[string]interface{}{
        {"name": "Alice", "age": 30, "salary": 75000, "department": "Engineering"},
        {"name": "Bob", "age": 25, "salary": 65000, "department": "Sales"},
        {"name": "Charlie", "age": 35, "salary": 85000, "department": "Engineering"},
        {"name": "Diana", "age": 28, "salary": 70000, "department": "Marketing"},
    },
}

// 筛选高薪工程师
engineers, _ := expr.Eval(`
    people 
    | filter(#.department == "Engineering" && #.salary > 70000)
    | map(#.name)
`, users)
// 结果: ["Alice", "Charlie"]

// 计算部门平均工资
avgSalary, _ := expr.Eval(`
    people
    | filter(#.department == "Engineering")
    | map(#.salary)
    | avg()
`, users)
// 结果: 80000

// 复杂业务逻辑
seniorHighEarners, _ := expr.Eval(`
    people
    | filter(#.age >= 30 && #.salary >= 75000)
    | map({
        name: #.name,
        level: #.age > 32 ? "Senior" : "Mid",
        bonus: #.salary * 0.1
    })
`, users)
```

### 5. 嵌套数据处理
```go
departments := map[string]interface{}{
    "orgs": []map[string]interface{}{
        {
            "name": "Tech",
            "teams": []map[string]interface{}{
                {"name": "Backend", "size": 8, "budget": 800000},
                {"name": "Frontend", "size": 6, "budget": 600000},
            },
        },
        {
            "name": "Business",
            "teams": []map[string]interface{}{
                {"name": "Sales", "size": 12, "budget": 500000},
                {"name": "Marketing", "size": 5, "budget": 300000},
            },
        },
    },
}

// 扁平化所有团队
allTeams, _ := expr.Eval(`
    orgs 
    | map(#.teams) 
    | reduce((acc, teams) => acc + teams, [])
    | map(#.name)
`, departments)
// 结果: ["Backend", "Frontend", "Sales", "Marketing"]

// 计算总预算
totalBudget, _ := expr.Eval(`
    orgs
    | map(#.teams)
    | reduce((acc, teams) => acc + teams, [])
    | map(#.budget)
    | sum()
`, departments)
// 结果: 2200000
```

### 6. 性能优化建议
```go
// ✅ 推荐：使用占位符语法（简洁高效）
"numbers | filter(# > 5) | map(# * 2)"

// ❌ 不推荐：传统Lambda语法（冗长）
"filter(numbers, x => x > 5) | map(y => y * 2)"

// ✅ 推荐：链式管道（清晰的数据流）
`data 
 | filter(# % 2 == 0)
 | map(# * # + 1)
 | take(5)`

// ✅ 推荐：复杂条件拆分
`users
 | filter(#.active == true)
 | filter(#.age >= 18)
 | map(#.email)`
```

## 内置函数使用

### 1. 集合操作函数
```go
data := map[string]interface{}{
    "numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
    "words": []string{"apple", "banana", "cherry", "date"},
    "users": []map[string]interface{}{
        {"name": "Alice", "age": 30, "city": "NYC"},
        {"name": "Bob", "age": 25, "city": "LA"},
        {"name": "Charlie", "age": 35, "city": "NYC"},
    },
}

// 过滤
evens, _ := expr.Eval("filter(numbers, n => n % 2 == 0)", data)
nycUsers, _ := expr.Eval(`filter(users, u => u.city == "NYC")`, data)

// 映射转换
squares, _ := expr.Eval("map(numbers, n => n * n)", data)
names, _ := expr.Eval("map(users, u => u.name)", data)

// 聚合计算
total, _ := expr.Eval("sum(numbers)", data)       // 55
average, _ := expr.Eval("avg(numbers)", data)     // 5.5
maxAge, _ := expr.Eval("max(map(users, u => u.age))", data) // 35

// 排序和选择
sorted, _ := expr.Eval("sort(numbers)", data)
top3, _ := expr.Eval("take(sort(numbers), 3)", data)
```

### 2. 字符串处理函数
```go
textData := map[string]interface{}{
    "sentence": "Hello, World! Welcome to Go programming.",
    "words": []string{"go", "rust", "python", "javascript"},
}

// 字符串分割和连接
words, _ := expr.Eval(`split(sentence, " ")`, textData)
joined, _ := expr.Eval(`join(words, " | ")`, textData)

// 字符串检查
hasWorld, _ := expr.Eval(`contains(sentence, "World")`, textData)
startsHello, _ := expr.Eval(`startsWith(sentence, "Hello")`, textData)

// 字符串格式化
upper, _ := expr.Eval("upper(sentence)", textData)
lower, _ := expr.Eval("lower(sentence)", textData)
trimmed, _ := expr.Eval(`trim("  spaced  ")`, nil)
```

### 3. 类型转换函数
```go
mixedData := map[string]interface{}{
    "numStr": "123",
    "floatStr": "45.67",
    "boolStr": "true",
    "number": 42,
}

// 类型转换
number, _ := expr.Eval("int(numStr)", mixedData)     // 123
decimal, _ := expr.Eval("float(floatStr)", mixedData) // 45.67
flag, _ := expr.Eval("bool(boolStr)", mixedData)     // true
text, _ := expr.Eval("string(number)", mixedData)    // "42"

// 类型检查
numType, _ := expr.Eval("type(number)", mixedData)   // "int"
```

## 配置选项和优化

### 1. 编译选项配置
```go
// 基础配置
program, err := expr.Compile("expression",
    expr.Env(environment),              // 设置环境
    expr.EnableOptimization(),          // 启用优化
    expr.EnableCache(),                 // 启用缓存
    expr.WithTimeout(time.Second*30),   // 设置超时
)

// 调试配置
debugProgram, err := expr.Compile("complex_expression",
    expr.EnableDebug(),      // 启用调试
    expr.EnableProfiling(),  // 启用性能分析
)

// 自定义函数
customProgram, err := expr.Compile("distance(p1, p2)",
    expr.WithBuiltin("distance", func(p1, p2 map[string]float64) float64 {
        dx := p1["x"] - p2["x"]
        dy := p1["y"] - p2["y"]
        return math.Sqrt(dx*dx + dy*dy)
    }),
)
```

### 2. 性能优化最佳实践
```go
// 1. 预编译复用
type ExpressionCache struct {
    programs map[string]*expr.Program
    mutex    sync.RWMutex
}

func (ec *ExpressionCache) GetProgram(expression string) (*expr.Program, error) {
    ec.mutex.RLock()
    program, exists := ec.programs[expression]
    ec.mutex.RUnlock()
    
    if exists {
        return program, nil
    }
    
    // 编译新表达式
    newProgram, err := expr.Compile(expression, expr.EnableOptimization())
    if err != nil {
        return nil, err
    }
    
    ec.mutex.Lock()
    ec.programs[expression] = newProgram
    ec.mutex.Unlock()
    
    return newProgram, nil
}

// 2. 环境变量池化
var envPool = sync.Pool{
    New: func() interface{} {
        return make(map[string]interface{})
    },
}

func evalWithPool(program *expr.Program, data interface{}) (interface{}, error) {
    env := envPool.Get().(map[string]interface{})
    defer func() {
        // 清空并回收
        for k := range env {
            delete(env, k)
        }
        envPool.Put(env)
    }()
    
    env["data"] = data
    return expr.Run(program, env)
}
```

## 实际应用场景

### 1. 业务规则引擎
```go
type BusinessRuleEngine struct {
    rules map[string]*expr.Program
}

func NewBusinessRuleEngine() *BusinessRuleEngine {
    engine := &BusinessRuleEngine{
        rules: make(map[string]*expr.Program),
    }
    
    // 预定义业务规则
    ruleDefinitions := map[string]string{
        "vip_customer": `
            customer.totalSpent > 10000 && 
            customer.membershipYears >= 2
        `,
        "discount_eligible": `
            customer.age < 25 || 
            customer.isStudent || 
            customer.firstOrder
        `,
        "free_shipping": `
            order.amount >= 50 || 
            customer.isPremium ||
            order.destination == "local"
        `,
        "credit_approved": `
            applicant.creditScore >= 650 &&
            applicant.income >= 50000 &&
            applicant.employment == "stable"
        `,
    }
    
    for name, rule := range ruleDefinitions {
        program, err := expr.Compile(rule, expr.EnableOptimization())
        if err != nil {
            panic(fmt.Sprintf("规则 %s 编译失败: %v", name, err))
        }
        engine.rules[name] = program
    }
    
    return engine
}

func (bre *BusinessRuleEngine) Evaluate(ruleName string, context interface{}) (bool, error) {
    program, exists := bre.rules[ruleName]
    if !exists {
        return false, fmt.Errorf("规则 %s 不存在", ruleName)
    }
    
    result, err := expr.Run(program, context)
    if err != nil {
        return false, err
    }
    
    return result.(bool), nil
}

// 使用示例
func main() {
    engine := NewBusinessRuleEngine()
    
    // 客户数据
    customer := map[string]interface{}{
        "totalSpent":      15000,
        "membershipYears": 3,
        "age":            28,
        "isPremium":       true,
    }
    
    order := map[string]interface{}{
        "amount":      75,
        "destination": "remote",
    }
    
    context := map[string]interface{}{
        "customer": customer,
        "order":    order,
    }
    
    // 评估规则
    isVIP, _ := engine.Evaluate("vip_customer", context)
    freeShipping, _ := engine.Evaluate("free_shipping", context)
    
    fmt.Printf("VIP客户: %v, 免费配送: %v\n", isVIP, freeShipping)
}
```

### 2. 动态配置管理
```go
type ConfigurationManager struct {
    configs map[string]*expr.Program
    mutex   sync.RWMutex
}

func NewConfigurationManager() *ConfigurationManager {
    return &ConfigurationManager{
        configs: make(map[string]*expr.Program),
    }
}

func (cm *ConfigurationManager) SetConfig(key, expression string) error {
    program, err := expr.Compile(expression,
        expr.AllowUndefinedVariables(),
        expr.EnableOptimization(),
    )
    if err != nil {
        return fmt.Errorf("配置 %s 编译失败: %v", key, err)
    }
    
    cm.mutex.Lock()
    cm.configs[key] = program
    cm.mutex.Unlock()
    
    return nil
}

func (cm *ConfigurationManager) GetValue(key string, context map[string]interface{}) (interface{}, error) {
    cm.mutex.RLock()
    program, exists := cm.configs[key]
    cm.mutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("配置 %s 不存在", key)
    }
    
    return expr.Run(program, context)
}

// 使用示例
func main() {
    config := NewConfigurationManager()
    
    // 设置动态配置表达式
    config.SetConfig("database.pool_size", `
        env == "production" ? 
            min(50, max(10, serverCount * 5)) : 
            5
    `)
    
    config.SetConfig("cache.ttl_seconds", `
        dataType == "user_profile" ? 3600 :
        dataType == "product_info" ? 1800 :
        dataType == "real_time" ? 60 :
        300
    `)
    
    config.SetConfig("feature.enabled", `
        region in ["US", "EU"] && 
        userCount > 1000 && 
        experimentGroup == "test"
    `)
    
    // 运行时获取配置值
    context := map[string]interface{}{
        "env":             "production",
        "serverCount":     8,
        "dataType":        "user_profile",
        "region":          "US",
        "userCount":       1500,
        "experimentGroup": "test",
    }
    
    poolSize, _ := config.GetValue("database.pool_size", context)
    cacheTTL, _ := config.GetValue("cache.ttl_seconds", context)
    featureEnabled, _ := config.GetValue("feature.enabled", context)
    
    fmt.Printf("数据库连接池: %v\n", poolSize)      // 40
    fmt.Printf("缓存TTL: %v秒\n", cacheTTL)         // 3600
    fmt.Printf("功能开关: %v\n", featureEnabled)    // true
}
```

### 3. 数据处理和分析
```go
type DataProcessor struct {
    transformers map[string]*expr.Program
}

func NewDataProcessor() *DataProcessor {
    dp := &DataProcessor{
        transformers: make(map[string]*expr.Program),
    }
    
    // 注册数据处理表达式
    transformations := map[string]string{
        "clean_sales_data": `
            data 
            | filter(record => record.amount > 0 && record.date != nil)
            | map(record => {
                amount: record.amount,
                date: record.date,
                region: record.region ?? "unknown",
                normalized_amount: record.amount / exchangeRate
            })
        `,
        "monthly_summary": `
            salesData
            | groupBy(record => substring(record.date, 0, 7))
            | map(group => {
                month: group.key,
                total_sales: sum(map(group.values, r => r.amount)),
                transaction_count: count(group.values),
                avg_transaction: avg(map(group.values, r => r.amount))
            })
        `,
        "top_regions": `
            salesData
            | groupBy(record => record.region)
            | map(group => {
                region: group.key,
                total: sum(map(group.values, r => r.amount))
            })
            | sort((a, b) => b.total - a.total)
            | take(5)
        `,
    }
    
    for name, expr := range transformations {
        program, err := expr.Compile(expr, expr.EnableOptimization())
        if err != nil {
            panic(fmt.Sprintf("转换器 %s 编译失败: %v", name, err))
        }
        dp.transformers[name] = program
    }
    
    return dp
}

func (dp *DataProcessor) Process(transformerName string, data interface{}, context map[string]interface{}) (interface{}, error) {
    program, exists := dp.transformers[transformerName]
    if !exists {
        return nil, fmt.Errorf("转换器 %s 不存在", transformerName)
    }
    
    env := make(map[string]interface{})
    env["data"] = data
    env["salesData"] = data
    for k, v := range context {
        env[k] = v
    }
    
    return expr.Run(program, env)
}

// 使用示例
func main() {
    processor := NewDataProcessor()
    
    // 原始销售数据
    rawSalesData := []map[string]interface{}{
        {"amount": 1000, "date": "2023-01-15", "region": "North"},
        {"amount": 1500, "date": "2023-01-16", "region": "South"},
        {"amount": 0, "date": "2023-01-17", "region": "North"},      // 无效数据
        {"amount": 2000, "date": "2023-02-01", "region": "West"},
        {"amount": 800, "date": "2023-02-02", "region": "North"},
    }
    
    context := map[string]interface{}{
        "exchangeRate": 1.0,
    }
    
    // 清理数据
    cleanedData, _ := processor.Process("clean_sales_data", rawSalesData, context)
    
    // 月度汇总
    monthlySummary, _ := processor.Process("monthly_summary", cleanedData, context)
    
    // 顶级区域
    topRegions, _ := processor.Process("top_regions", cleanedData, context)
    
    fmt.Printf("清理后数据: %v\n", cleanedData)
    fmt.Printf("月度汇总: %v\n", monthlySummary)
    fmt.Printf("顶级区域: %v\n", topRegions)
}
```

## 错误处理和调试

### 1. 错误处理最佳实践
```go
func safeEvaluate(expression string, env interface{}) (result interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("表达式执行panic: %v", r)
        }
    }()
    
    // 预验证
    if strings.TrimSpace(expression) == "" {
        return nil, fmt.Errorf("表达式不能为空")
    }
    
    // 尝试编译
    program, err := expr.Compile(expression)
    if err != nil {
        return nil, fmt.Errorf("编译错误: %v", err)
    }
    
    // 执行
    result, err = expr.Run(program, env)
    if err != nil {
        return nil, fmt.Errorf("执行错误: %v", err)
    }
    
    return result, nil
}
```

### 2. 调试和监控
```go
// 性能监控
func evaluateWithMetrics(expression string, env interface{}) {
    start := time.Now()
    
    result, err := expr.Eval(expression, env)
    
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("表达式执行失败: %s, 错误: %v, 耗时: %v", 
            expression, err, duration)
    } else {
        log.Printf("表达式执行成功: %s, 结果: %v, 耗时: %v", 
            expression, result, duration)
    }
}

// 获取全局统计
func printStats() {
    stats := expr.GetStatistics()
    fmt.Printf("==== 表达式引擎统计 ====\n")
    fmt.Printf("编译次数: %d\n", stats.TotalCompilations)
    fmt.Printf("执行次数: %d\n", stats.TotalExecutions)
    fmt.Printf("平均编译时间: %v\n", stats.AverageCompileTime)
    fmt.Printf("平均执行时间: %v\n", stats.AverageExecTime)
    fmt.Printf("缓存命中率: %.2f%%\n", stats.CacheHitRate*100)
    fmt.Printf("内存使用: %d bytes\n", stats.MemoryUsage)
}
```

## 总结

表达式引擎提供了强大且高性能的表达式处理能力，通过：

1. **零反射架构** - 确保极高的执行性能
2. **丰富的语法支持** - 覆盖各种业务场景需求
3. **内置函数库** - 提供40+个常用函数
4. **灵活的配置选项** - 支持各种定制需求
5. **完善的错误处理** - 提供详细的错误信息和调试支持

适用于业务规则引擎、动态配置、数据处理、API过滤等多种场景，是构建灵活系统的理想选择。

### 关键建议

1. **性能优化**: 使用预编译来提高重复执行的性能
2. **错误处理**: 始终进行错误检查和异常捕获
3. **资源管理**: 合理使用对象池来减少GC压力
4. **监控诊断**: 利用统计信息来优化系统性能
5. **模块化使用**: 根据具体需求选择合适的模块和功能

通过本指南，您应该能够充分利用表达式引擎的强大功能，构建高性能、灵活可配置的应用系统。 