# 模块系统 - 功能扩展架构

## 概述

模块系统为表达式引擎提供了可扩展的功能架构，通过模块注册和管理机制，支持自定义函数和功能的动态加载。当前实现了数学模块和字符串模块，为表达式提供丰富的内置功能。

> 💡 **调试支持**: 表达式引擎还提供了强大的调试功能，详见 [Debug 模块文档](14-debug.md)。

## 🏗️ 核心架构

### ModuleRegistry - 模块注册器
```go
type Registry struct {
    modules map[string]*Module
    mu      sync.RWMutex
}

// 模块定义
type Module struct {
    Name        string
    Description string
    Functions   map[string]*ModuleFunction
}

// 模块函数定义
type ModuleFunction struct {
    Name        string
    Description string
    Handler     func(args ...interface{}) (interface{}, error)
    ParamTypes  []types.TypeInfo
    ReturnType  types.TypeInfo
    Variadic    bool
}

// 全局模块注册器
var DefaultRegistry = NewRegistry()
```

## 📊 内置模块

### 1. Math 模块 - 数学运算

#### 基础数学函数
```go
// 基本运算
math.abs(x)          // 绝对值
math.sqrt(x)         // 平方根  
math.pow(x, y)       // 幂运算 x^y
math.max(a, b)       // 最大值
math.min(a, b)       // 最小值

// 取整函数
math.floor(x)        // 向下取整
math.ceil(x)         // 向上取整
math.round(x)        // 四舍五入

// 三角函数
math.sin(x)          // 正弦值 (弧度)
math.cos(x)          // 余弦值 (弧度)  
math.tan(x)          // 正切值 (弧度)

// 对数函数
math.log(x)          // 自然对数
math.log10(x)        // 以10为底的对数
math.exp(x)          // e^x
```

#### 使用示例
```go
// 计算圆的面积
expr := "math.pi * math.pow(radius, 2)"
result, _ := Run(expr, map[string]interface{}{
    "radius": 5.0,
})
// 结果: 78.54

// 计算三角形斜边
expr := "math.sqrt(math.pow(a, 2) + math.pow(b, 2))"
result, _ := Run(expr, map[string]interface{}{
    "a": 3.0,
    "b": 4.0,
})
// 结果: 5.0
```

### 2. Strings 模块 - 字符串处理

#### 字符串操作函数
```go
// 大小写转换
strings.upper(s)           // 转大写
strings.lower(s)           // 转小写

// 空格处理
strings.trim(s)            // 去除前后空格
strings.length(s)          // 字符串长度

// 查找和匹配
strings.contains(s, substr)    // 包含检查
strings.startsWith(s, prefix)  // 前缀检查
strings.endsWith(s, suffix)    // 后缀检查
strings.indexOf(s, substr)     // 查找位置

// 字符串操作
strings.replace(s, old, new)   // 替换内容
strings.split(s, sep)          // 分割字符串
strings.join(arr, sep)         // 连接字符串
strings.substring(s, start, end) // 截取子串
strings.repeat(s, n)           // 重复字符串
```

#### 使用示例
```go
// 处理用户姓名
expr := "strings.upper(strings.trim(firstName))"
result, _ := Run(expr, map[string]interface{}{
    "firstName": " john ",
})
// 结果: "JOHN"

// 构建完整姓名
expr := "strings.join([lastName, firstName], ', ')"
result, _ := Run(expr, map[string]interface{}{
    "firstName": "John",
    "lastName": "Doe",
})
// 结果: "Doe, John"

// 检查邮箱格式
expr := "strings.contains(email, '@') && strings.contains(email, '.')"
result, _ := Run(expr, map[string]interface{}{
    "email": "user@example.com",
})
// 结果: true
```



## 🔧 模块使用

### 基本使用
```go
import (
    "github.com/mredencom/expr"
    "github.com/mredencom/expr/modules"
)

func main() {
    // 使用模块函数
    result, err := api.Run("math.sqrt(16) + strings.length('hello')", nil)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // 9 (4 + 5)
}
```

### 自定义模块注册
```go
// 创建自定义函数
customFunctions := map[string]*modules.ModuleFunction{
    "greet": {
        Name:        "greet",
        Description: "问候函数",
        Handler: func(args ...interface{}) (interface{}, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("greet expects 1 argument")
            }
            name := args[0].(string)
            return fmt.Sprintf("Hello, %s!", name), nil
        },
        ParamTypes: []types.TypeInfo{
            {Kind: types.KindString, Name: "string"},
        },
        ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
    },
}

// 注册自定义模块
err := modules.DefaultRegistry.RegisterModule("custom", "自定义模块", customFunctions)
if err != nil {
    panic(err)
}

// 使用自定义模块
result, _ := api.Run("custom.greet('World')", nil)
fmt.Println(result) // "Hello, World!"
```

### 模块信息查询
```go
// 获取所有模块
modules := modules.DefaultRegistry.ListModules()
fmt.Println("可用模块:", modules)

// 获取模块信息
mathModule, _ := modules.DefaultRegistry.GetModuleInfo("math")
fmt.Printf("模块: %s - %s\n", mathModule.Name, mathModule.Description)

// 遍历模块函数
for name, function := range mathModule.Functions {
    fmt.Printf("  函数: %s - %s\n", name, function.Description)
}

// 调用模块函数
result, _ := modules.DefaultRegistry.CallFunction("math", "sqrt", 25.0)
fmt.Println("sqrt(25) =", result) // 5.0
```

## 📈 性能优化

### 模块函数缓存
- 模块函数在注册时被缓存
- 避免在运行时频繁查找函数
- 使用全局注册器提高访问效率

模块系统为表达式引擎提供了强大的扩展能力，通过模块化设计实现了功能的高内聚低耦合，便于维护和扩展。
