package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	expr "github.com/mredencom/expr"
)

// User 用户结构体
type User struct {
	Name     string
	Age      int
	Email    string
	Active   bool
	Balance  float64
	Tags     []string
	Metadata map[string]interface{}
}

// ToMap implements the StructConverter interface for zero-reflection conversion
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Name":     u.Name,
		"Age":      u.Age,
		"Email":    u.Email,
		"Active":   u.Active,
		"Balance":  u.Balance,
		"Tags":     u.Tags,
		"Metadata": u.Metadata,
	}
}

// Product 产品结构体
type Product struct {
	ID       int
	Name     string
	Price    float64
	Category string
	InStock  bool
	Tags     []string
}

// ToMap implements the StructConverter interface for zero-reflection conversion
func (p Product) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":       p.ID,
		"Name":     p.Name,
		"Price":    p.Price,
		"Category": p.Category,
		"InStock":  p.InStock,
		"Tags":     p.Tags,
	}
}

func main() {
	fmt.Println("🚀 Go Expression Engine - 综合功能演示")
	fmt.Println(strings.Repeat("=", 60))

	// 基础数据准备
	users := []User{
		{Name: "Alice", Age: 28, Email: "alice@example.com", Active: true, Balance: 1500.50, Tags: []string{"vip", "premium"}, Metadata: map[string]interface{}{"level": "gold", "score": 95}},
		{Name: "Bob", Age: 32, Email: "bob@example.com", Active: true, Balance: 2300.75, Tags: []string{"regular"}, Metadata: map[string]interface{}{"level": "silver", "score": 78}},
		{Name: "Charlie", Age: 25, Email: "charlie@example.com", Active: false, Balance: 450.25, Tags: []string{"new"}, Metadata: map[string]interface{}{"level": "bronze", "score": 65}},
		{Name: "Diana", Age: 35, Email: "diana@example.com", Active: true, Balance: 3200.00, Tags: []string{"vip", "enterprise"}, Metadata: map[string]interface{}{"level": "platinum", "score": 98}},
	}

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	prices := []float64{29.99, 199.99, 899.99, 2999.99}

	// 1. 基础表达式演示
	fmt.Println("\n📝 1. 基础表达式演示")
	fmt.Println(strings.Repeat("-", 30))

	basicExpressions := []string{
		"2 + 3 * 4",                   // 算术运算
		"'Hello' + ' ' + 'World'",     // 字符串连接
		"true && (false || true)",     // 布尔逻辑
		"42 > 30 ? 'large' : 'small'", // 三元条件
		"abs(-42)",                    // 内置函数
		"max(1, 5, 3, 9, 2)",          // 多参数函数
	}

	for _, expression := range basicExpressions {
		result, err := expr.Eval(expression, nil)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", expression, err)
			continue
		}
		fmt.Printf("✅ %-30s → %v\n", expression, result)
	}

	// 2. 变量和环境演示
	fmt.Println("\n🔧 2. 变量和环境演示")
	fmt.Println(strings.Repeat("-", 30))

	env := map[string]interface{}{
		"user":       users[0],
		"threshold":  1000.0,
		"multiplier": 2.5,
		"prefix":     "Mr./Ms. ",
	}

	envExpressions := []string{
		"user.Name",
		"user.Age >= 25",
		"user.Balance > threshold",
		"prefix + user.Name",
		"user.Active && user.Balance > threshold",
		"len(user.Tags)",
		"contains(user.Email, '@')",
		"user.Metadata['level']",
	}

	for _, expression := range envExpressions {
		result, err := expr.Eval(expression, env)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", expression, err)
			continue
		}
		fmt.Printf("✅ %-35s → %v\n", expression, result)
	}

	// 3. 数组和集合操作演示
	fmt.Println("\n📊 3. 数组和集合操作演示")
	fmt.Println(strings.Repeat("-", 30))

	arrayEnv := map[string]interface{}{
		"numbers":   numbers,
		"prices":    prices,
		"userCount": len(users),
	}

	arrayExpressions := []string{
		"len(numbers)",
		"sum(numbers)",
		"avg(prices)",
		"max(numbers)",
		"min(prices)",
		"numbers[0]",
		"numbers[len(numbers)-1]",
	}

	for _, expression := range arrayExpressions {
		result, err := expr.Eval(expression, arrayEnv)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", expression, err)
			continue
		}
		fmt.Printf("✅ %-25s → %v\n", expression, result)
	}

	// 4. 管道占位符语法演示 - 核心功能
	fmt.Println("\n🔥 4. 管道占位符语法演示 (核心功能)")
	fmt.Println(strings.Repeat("-", 45))

	pipelineEnv := map[string]interface{}{
		"numbers":    numbers,
		"threshold":  30,
		"userAges":   []int{28, 32, 25, 35},
		"productIds": []int{1, 2, 3, 4},
	}

	pipelineExpressions := []struct {
		expr        string
		description string
	}{
		{"numbers | filter(# > 5)", "过滤大于5的数字"},
		{"numbers | filter(# % 2 == 0)", "过滤偶数"},
		{"numbers | map(# * 2)", "每个数字乘以2"},
		{"numbers | filter(# > 3) | map(# * 2)", "链式操作：过滤后映射"},
		{"numbers | filter(# % 2 == 1) | map(# * #)", "奇数的平方"},
		{"numbers | filter(# > threshold / 10)", "动态阈值过滤"},
	}

	for _, item := range pipelineExpressions {
		result, err := expr.Eval(item.expr, pipelineEnv)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", item.expr, err)
			continue
		}
		fmt.Printf("✅ %-40s → %s\n", item.expr, formatResult(result))
		fmt.Printf("   💡 %s\n", item.description)
		fmt.Println()
	}

	// 5. 复杂表达式演示
	fmt.Println("\n🧠 5. 复杂表达式演示")
	fmt.Println(strings.Repeat("-", 30))

	complexEnv := map[string]interface{}{
		"numbers":    numbers,
		"vipLevel":   "gold",
		"minBalance": 1000.0,
		"discount":   0.1,
		// 使用基本类型而不是结构体数组
		"balances": []float64{1500.50, 2300.75, 450.25, 3200.00},
		"ages":     []int{28, 32, 25, 35},
	}

	complexExpressions := []struct {
		expr        string
		description string
	}{
		{
			"numbers | filter(# > 3) | map(# * 2 + 1) | filter(# % 3 == 0)",
			"多级数值处理管道：过滤>3，转换为2n+1，再过滤3的倍数",
		},
		{
			"numbers | filter(# % 2 == 0 && # > 4) | map(# * # - 1)",
			"偶数且>4的数字，计算平方减1",
		},
		{
			"numbers | map(# > 5 ? # * 10 : # * 2)",
			"条件映射：>5的数字×10，否则×2",
		},
	}

	for _, item := range complexExpressions {
		result, err := expr.Eval(item.expr, complexEnv)
		if err != nil {
			log.Printf("❌ Error evaluating complex expression: %v", err)
			fmt.Printf("   Expression: %s\n", item.expr)
			fmt.Printf("   💡 %s\n", item.description)
			fmt.Println()
			continue
		}
		fmt.Printf("✅ 复杂表达式成功执行\n")
		fmt.Printf("   Expression: %s\n", item.expr)
		fmt.Printf("   Result: %s\n", formatResult(result))
		fmt.Printf("   💡 %s\n", item.description)
		fmt.Println()
	}

	// 6. 字符串处理演示
	fmt.Println("\n📝 6. 字符串处理演示")
	fmt.Println(strings.Repeat("-", 30))

	stringEnv := map[string]interface{}{
		"text":  "Hello, World! This is a test.",
		"email": "user@example.com",
		"csv":   "apple,banana,cherry,date",
		"words": []string{"hello", "world", "test", "demo"},
	}

	stringExpressions := []string{
		"upper(text)",
		"lower(text)",
		"split(csv, ',')",
		"join(words, '-')",
		"contains(email, '@')",
		"startsWith(text, 'Hello')",
		"endsWith(text, 'test.')",
		"len(text)",
	}

	for _, expression := range stringExpressions {
		result, err := expr.Eval(expression, stringEnv)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", expression, err)
			continue
		}
		fmt.Printf("✅ %-35s → %s\n", expression, formatResult(result))
	}

	// 7. 类型转换和验证演示
	fmt.Println("\n🔄 7. 类型转换和验证演示")
	fmt.Println(strings.Repeat("-", 35))

	typeEnv := map[string]interface{}{
		"numbers":   []interface{}{1, 2.5, "3", true},
		"mixedData": []interface{}{"42", 3.14, true, "hello"},
	}

	typeExpressions := []string{
		"type(42)",
		"type('hello')",
		"type(true)",
		"string(123)",
		"int('42')",
		"float('3.14')",
		"bool('true')",
	}

	for _, expression := range typeExpressions {
		result, err := expr.Eval(expression, typeEnv)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", expression, err)
			continue
		}
		fmt.Printf("✅ %-25s → %s\n", expression, formatResult(result))
	}

	// 8. 性能演示
	fmt.Println("\n⚡ 8. 性能演示")
	fmt.Println(strings.Repeat("-", 20))

	performanceTest := func(expression string, env map[string]interface{}, iterations int) {
		// 编译一次
		program, err := expr.Compile(expression)
		if err != nil {
			log.Printf("❌ Compilation error: %v", err)
			return
		}

		start := time.Now()
		for i := 0; i < iterations; i++ {
			_, err := expr.Run(program, env)
			if err != nil {
				log.Printf("❌ Execution error: %v", err)
				return
			}
		}
		duration := time.Since(start)

		opsPerSec := float64(iterations) / duration.Seconds()
		fmt.Printf("✅ %-40s: %d ops in %v (%.0f ops/sec)\n",
			expression, iterations, duration, opsPerSec)
	}

	performanceTest("numbers | filter(# > 5) | map(# * 2)",
		map[string]interface{}{"numbers": numbers}, 10000)
	performanceTest("2 + 3 * 4", nil, 100000)
	performanceTest("'Hello' + ' ' + 'World'", nil, 50000)

	// 9. 错误处理演示
	fmt.Println("\n❌ 9. 错误处理演示")
	fmt.Println(strings.Repeat("-", 25))

	errorExpressions := []string{
		"undefinedVariable",
		"numbers | filter(# > 'invalid')",
		"split('hello', '')",
		"int('not_a_number')",
	}

	for _, expression := range errorExpressions {
		_, err := expr.Eval(expression, map[string]interface{}{
			"numbers": numbers,
		})
		if err != nil {
			fmt.Printf("✅ %-30s → Error caught: %v\n", expression, err)
		} else {
			fmt.Printf("⚠️  %-30s → Unexpected success\n", expression)
		}
	}

	// 10. 高级管道占位符演示
	fmt.Println("\n🚀 10. 高级管道占位符演示")
	fmt.Println(strings.Repeat("-", 35))

	advancedPipelineExpressions := []struct {
		expr        string
		description string
	}{
		{
			"numbers | filter(# > 2 && # < 8) | map(# * 3 - 1)",
			"复合条件过滤 + 复杂映射",
		},
		{
			"numbers | map(# % 3 == 0 ? 'fizz' : string(#))",
			"条件映射：3的倍数显示'fizz'",
		},
		{
			"numbers | filter(# % 2 == 0) | map(# + 10) | filter(# > 15)",
			"三级管道：偶数 → 加10 → 过滤>15",
		},
		{
			"numbers | map((# + 1) * (# - 1))",
			"数学表达式：(n+1)*(n-1)",
		},
		{
			"numbers | filter(# >= 3 && # <= 7) | map(# * 2) | sum",
			"管道链终结于聚合函数",
		},
	}

	for _, item := range advancedPipelineExpressions {
		result, err := expr.Eval(item.expr, map[string]interface{}{"numbers": numbers})
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", item.expr, err)
			continue
		}
		fmt.Printf("✅ Expression: %s\n", item.expr)
		fmt.Printf("   Result: %s\n", formatResult(result))
		fmt.Printf("   💡 %s\n", item.description)
		fmt.Println()
	}

	fmt.Println("\n🎉 演示完成！")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✨ 表达式引擎功能总结:")
	fmt.Println("   • 基础算术、逻辑、字符串操作")
	fmt.Println("   • 变量和环境集成")
	fmt.Println("   • 数组和集合操作")
	fmt.Println("   • 🔥 管道占位符语法 (# 语法) - 核心亮点")
	fmt.Println("   • 复杂表达式链式操作")
	fmt.Println("   • 字符串处理和分割")
	fmt.Println("   • 类型转换和验证")
	fmt.Println("   • 高性能执行")
	fmt.Println("   • 完善的错误处理")
	fmt.Println("   • 多级管道组合")
	fmt.Println("   • 条件映射和复杂逻辑")
}

// formatResult 格式化输出结果
func formatResult(result interface{}) string {
	switch v := result.(type) {
	case []interface{}:
		if len(v) > 5 {
			return fmt.Sprintf("[%v, %v, %v, ... (%d items)]", v[0], v[1], v[2], len(v))
		}
		return fmt.Sprintf("%v", v)
	case string:
		if len(v) > 50 {
			return fmt.Sprintf("'%s...' (%d chars)", v[:47], len(v))
		}
		return fmt.Sprintf("'%s'", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
