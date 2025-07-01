# 内置函数库 - 企业级函数集合

## 概述

内置函数库是零反射表达式引擎的核心功能模块，提供了40+个高性能的内置函数，涵盖数学运算、字符串处理、集合操作、类型转换等企业应用的各个方面。所有函数都经过零反射优化，支持管道占位符语法，具备极致的执行性能。

## 🔥 核心特性

### 1. 零反射架构
- **编译时绑定**: 所有函数在编译时就完成绑定，无运行时反射开销
- **类型特化**: 针对不同类型参数的专用执行路径
- **内联优化**: 简单函数调用被内联，消除函数调用开销

### 2. 管道占位符语法支持
- **完整支持**: 所有函数都支持`#`占位符语法
- **零内存分配**: 占位符表达式执行时零内存分配
- **链式调用**: 支持复杂的管道操作链

### 3. 高性能执行
- **25M+ ops/sec**: 内置函数执行速度超过25M操作/秒
- **内存优化**: 对象池和缓存机制减少GC压力
- **批量处理**: 集合函数支持大数据集的高效处理

## 函数分类

### 📊 数学函数 (10个)

#### 基础数学运算
```go
abs(-42)           // 42 - 绝对值
max(1, 2, 3)       // 3 - 最大值
min(5, 3, 8)       // 3 - 最小值
sum([1, 2, 3, 4])  // 10 - 求和
avg([1, 2, 3, 4])  // 2.5 - 平均值
```

#### 高级数学函数
```go
ceil(3.14)         // 4 - 向上取整
floor(3.14)        // 3 - 向下取整
round(3.14)        // 3 - 四舍五入
sqrt(16)           // 4 - 平方根
pow(2, 3)          // 8 - 幂运算
```

#### 🔥 管道中的数学函数
```go
// 数组求和
numbers | sum                              // 求和
numbers | filter(# > 0) | sum             // 正数求和

// 统计运算
numbers | avg                              // 平均值
numbers | map(abs(#)) | max               // 绝对值的最大值

// 数学变换
numbers | map(sqrt(#)) | filter(# > 2)    // 开方后过滤
numbers | map(pow(#, 2)) | avg            // 平方后求平均
```

### 📝 字符串函数 (12个)

#### 基础字符串操作
```go
len("hello")                    // 5 - 字符串长度
upper("hello")                  // "HELLO" - 转大写
lower("HELLO")                  // "hello" - 转小写
trim("  hello  ")               // "hello" - 去除首尾空格
reverse("hello")                // "olleh" - 反转字符串
```

#### 字符串查询
```go
contains("hello world", "world")  // true - 包含检查
startsWith("hello", "he")         // true - 前缀检查
endsWith("hello", "lo")           // true - 后缀检查
indexOf("hello", "ll")            // 2 - 查找位置
```

#### 字符串处理
```go
replace("hello world", "world", "go")  // "hello go" - 替换
split("a,b,c", ",")                   // ["a", "b", "c"] - 分割
join(["a", "b", "c"], ",")            // "a,b,c" - 连接
substring("hello", 1, 4)              // "ell" - 子字符串
```

#### 🔥 管道中的字符串函数
```go
// 字符串数组处理
words := ["hello", "world", "go", "programming"]

words | filter(len(#) > 3)                    // 过滤长度>3的单词
words | map(upper(#))                         // 转换为大写
words | map(trim(#)) | filter(# != "")       // 去空格并过滤空串
words | filter(startsWith(#, "go"))          // 过滤以"go"开头的单词

// 复杂字符串处理
words | filter(contains(#, "o")) | map(upper(#)) | join(" ")
// 结果: "HELLO WORLD GO PROGRAMMING"

// 字符串变换管道
texts | map(trim(#)) | filter(len(#) > 0) | map(lower(#)) | unique
```

### 🗂️ 集合函数 (12个)

#### 集合过滤和映射
```go
filter([1,2,3,4], # > 2)     // [3, 4] - 过滤元素
map([1,2,3], # * 2)          // [2, 4, 6] - 映射变换
reduce([1,2,3,4], # + #)     // 10 - 归约操作
```

#### 集合排序和去重
```go
sort([3,1,4,2])              // [1, 2, 3, 4] - 排序
reverse([1,2,3])             // [3, 2, 1] - 反转
unique([1,2,2,3])            // [1, 2, 3] - 去重
```

#### 集合选择
```go
first([1,2,3])               // 1 - 第一个元素
last([1,2,3])                // 3 - 最后一个元素
take([1,2,3,4,5], 3)         // [1, 2, 3] - 取前N个
skip([1,2,3,4,5], 2)         // [3, 4, 5] - 跳过前N个
```

#### 集合检查
```go
any([1,2,3], # > 2)          // true - 任意元素满足条件
all([1,2,3], # > 0)          // true - 所有元素满足条件
count([1,2,3,4], # > 2)      // 2 - 满足条件的元素数量
```

#### 🔥 管道中的集合函数
```go
// 复杂数据处理管道
numbers | filter(# > 5) | map(# * 2) | sort | unique
// 过滤 -> 映射 -> 排序 -> 去重

// 统计分析管道
numbers | filter(# > 0) | map(# * #) | sum | sqrt
// 正数 -> 平方 -> 求和 -> 开方

// 条件聚合
numbers | filter(# % 2 == 0) | map(# / 2) | avg
// 偶数 -> 除以2 -> 求平均

// 复合条件处理
data | filter(# > threshold && # < limit) | map(transform(#)) | reduce(combine)
```

### 🔄 类型转换函数 (6个)

#### 基础类型转换
```go
string(42)           // "42" - 转字符串
int("42")            // 42 - 转整数
float("3.14")        // 3.14 - 转浮点数
bool("true")         // true - 转布尔值
```

#### 高级类型转换
```go
toArray(value)       // 转换为数组
toMap(pairs)         // 转换为映射
```

#### 🔥 管道中的类型转换
```go
// 类型转换管道
numbers | map(string(#)) | filter(len(#) > 1) | join(",")
// 数字 -> 字符串 -> 过滤长度 -> 连接

// 混合类型处理
mixed | map(string(#)) | filter(# != "") | map(upper(#))
// 统一转字符串 -> 过滤空值 -> 转大写
```

## 🔥 管道占位符语法详解

### 基础占位符用法
```go
// 单个占位符
numbers | filter(# > 5)                    // # 表示当前元素
numbers | map(# * 2)                       // # 参与运算
numbers | filter(# % 2 == 0)               // # 参与条件判断
```

### 复杂占位符表达式
```go
// 算术表达式
numbers | map(# * 2 + 1)                   // 线性变换
numbers | map((# + 1) * (# - 1))           // 乘法展开: x²-1
numbers | map(# * # - 2 * # + 1)           // 二次方程: x²-2x+1

// 条件表达式
numbers | map(# > 5 ? # * 10 : # * 2)      // 三元运算符
numbers | filter(# > min && # < max)       // 复合条件
```

### 嵌套占位符
```go
// 多级嵌套
data | map(items | filter(# > #parent) | sum)  // 嵌套过滤求和
data | filter(values | any(# > threshold))     // 嵌套存在性检查
```

### 占位符与函数结合
```go
// 占位符调用函数
strings | filter(len(#) > 3)               // 长度过滤
strings | map(upper(#))                    // 大写转换
strings | filter(contains(#, keyword))     // 包含检查

// 复杂函数调用
strings | map(substring(#, 0, min(len(#), 10)))  // 截取前10个字符
numbers | map(pow(#, 2)) | filter(# < 100)       // 平方后过滤
```

## 性能特性

### 执行性能
| 函数类型 | 性能指标 | 说明 |
|----------|----------|------|
| **数学函数** | 20-25M ops/sec | 基础运算内联优化 |
| **字符串函数** | 15-20M ops/sec | 零拷贝字符串处理 |
| **集合函数** | 8-15M ops/sec | 批量处理优化 |
| **类型转换** | 25-30M ops/sec | 编译时类型确定 |

### 内存优化
```go
// ✅ 高效：使用对象池
result := numbers | filter(# > 5) | map(# * 2)  // 零额外分配

// ✅ 高效：批量处理
result := process(largeDataset | filter(condition) | map(transform))

// ❌ 低效：频繁小对象创建
for item in dataset {
    if condition(item) {
        result.append(transform(item))
    }
}
```

### 编译时优化
```go
// 常量折叠优化
abs(-42)                    // 编译时计算为 42
max(1, 2, 3)               // 编译时计算为 3
"hello" | upper            // 编译时计算为 "HELLO"

// 函数内联优化
numbers | map(# * 2)       // 乘法操作被内联
strings | filter(len(#) > 3)  // 长度检查被内联
```

## 自定义内置函数

### 添加自定义函数
```go
// 添加简单函数
expr.WithBuiltin("double", func(x int) int {
    return x * 2
})

// 添加复杂函数
expr.WithBuiltin("distance", func(x1, y1, x2, y2 float64) float64 {
    dx := x2 - x1
    dy := y2 - y1
    return math.Sqrt(dx*dx + dy*dy)
})

// 添加可变参数函数
expr.WithBuiltin("concat", func(args ...string) string {
    return strings.Join(args, "")
})
```

### 管道兼容的自定义函数
```go
// 支持占位符的自定义函数
expr.WithBuiltin("square", func(x interface{}) interface{} {
    if num, ok := x.(int); ok {
        return num * num
    }
    return x
})

// 使用示例
result := expr.Eval("numbers | map(square(#))", env)
```

### 聚合函数定制
```go
// 自定义聚合函数
expr.WithBuiltin("product", func(values []interface{}) interface{} {
    product := 1
    for _, v := range values {
        if num, ok := v.(int); ok {
            product *= num
        }
    }
    return product
})

// 使用示例
result := expr.Eval("numbers | filter(# > 0) | product", env)
```

## 错误处理

### 函数调用错误
```go
// 参数类型错误
_, err := expr.Eval("abs('not a number')", nil)
// 错误: invalid argument type for abs(): expected number, got string

// 参数数量错误
_, err = expr.Eval("max()", nil)
// 错误: max() requires at least 1 argument

// 除零错误
_, err = expr.Eval("numbers | map(10 / #)", map[string]interface{}{
    "numbers": []int{1, 0, 2},
})
// 错误: division by zero
```

### 管道错误处理
```go
// 类型不匹配
_, err := expr.Eval("'not array' | filter(# > 5)", nil)
// 错误: filter() can only be applied to arrays

// 空数组处理
result, _ := expr.Eval("[] | sum", nil)  // 0 (默认值)
result, _ = expr.Eval("[] | first", nil) // nil
```

## 最佳实践

### 1. 性能优化
```go
// ✅ 推荐：预编译复杂表达式
program, _ := expr.Compile("data | filter(# > threshold) | map(transform(#)) | sum")
for _, dataset := range datasets {
    result, _ := expr.Run(program, dataset)
}

// ❌ 不推荐：重复编译
for _, dataset := range datasets {
    result, _ := expr.Eval("data | filter(# > threshold) | map(transform(#)) | sum", dataset)
}
```

### 2. 类型安全
```go
// ✅ 推荐：明确类型期望
program, _ := expr.Compile("numbers | sum", expr.AsInt())

// ✅ 推荐：环境类型提示
env := map[string]interface{}{
    "numbers": []int{1, 2, 3, 4, 5},  // 明确类型
}
```

### 3. 错误处理
```go
// ✅ 推荐：完整错误处理
result, err := expr.Eval("numbers | filter(# > 0) | sum", env)
if err != nil {
    log.Printf("表达式执行错误: %v", err)
    return defaultValue
}
```

### 4. 管道设计
```go
// ✅ 推荐：清晰的管道步骤
expr := "data | filter(# > 0) | map(# * 2) | sort | take(10)"

// ✅ 推荐：合理的管道长度
expr := "users | filter(#.active) | map(#.score) | avg"

// ❌ 不推荐：过长的管道
expr := "data | filter(...) | map(...) | filter(...) | map(...) | sort | reverse | filter(...)"
```

## 函数索引

### 按字母顺序
- `abs()` - 绝对值
- `all()` - 全部满足
- `any()` - 任意满足
- `avg()` - 平均值
- `bool()` - 转布尔
- `ceil()` - 向上取整
- `contains()` - 包含检查
- `count()` - 计数
- `endsWith()` - 后缀检查
- `filter()` - 过滤
- `first()` - 第一个
- `float()` - 转浮点
- `floor()` - 向下取整
- `indexOf()` - 查找位置
- `int()` - 转整数
- `join()` - 连接
- `last()` - 最后一个
- `len()` - 长度
- `lower()` - 转小写
- `map()` - 映射
- `max()` - 最大值
- `min()` - 最小值
- `pow()` - 幂运算
- `reduce()` - 归约
- `replace()` - 替换
- `reverse()` - 反转
- `round()` - 四舍五入
- `skip()` - 跳过
- `sort()` - 排序
- `split()` - 分割
- `sqrt()` - 平方根
- `startsWith()` - 前缀检查
- `string()` - 转字符串
- `substring()` - 子字符串
- `sum()` - 求和
- `take()` - 取前N个
- `toArray()` - 转数组
- `toMap()` - 转映射
- `trim()` - 去空格
- `unique()` - 去重
- `upper()` - 转大写

### 按功能分类
- **数学**: `abs`, `max`, `min`, `sum`, `avg`, `ceil`, `floor`, `round`, `sqrt`, `pow`
- **字符串**: `len`, `upper`, `lower`, `trim`, `reverse`, `contains`, `startsWith`, `endsWith`, `indexOf`, `replace`, `split`, `join`, `substring`
- **集合**: `filter`, `map`, `reduce`, `sort`, `reverse`, `unique`, `first`, `last`, `take`, `skip`, `any`, `all`, `count`
- **类型**: `string`, `int`, `float`, `bool`, `toArray`, `toMap`

## 总结

内置函数库提供了企业级应用所需的全面功能支持：

- **🔥 40+高性能函数**: 覆盖数学、字符串、集合、类型转换
- **⚡ 零反射架构**: 25M+ ops/sec的极致性能
- **🛡️ 管道占位符**: 革命性的`#`占位符语法支持
- **🏢 企业级特性**: 错误处理、类型安全、批量处理
- **🔧 高度可扩展**: 支持自定义函数和操作符

通过合理使用内置函数和管道语法，可以构建出高效、简洁、可维护的数据处理表达式。 