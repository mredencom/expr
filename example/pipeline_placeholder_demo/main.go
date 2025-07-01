package main

import (
	"fmt"
	"time"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("🔥 管道占位符功能演示")
	fmt.Println("===================")

	demonstrateBasicPlaceholders()
	demonstrateComplexExpressions()
	demonstrateChainedPipelines()
	demonstrateObjectProcessing()
	demonstrateRealWorldExamples()
	demonstratePerformanceComparison()
}

func demonstrateBasicPlaceholders() {
	fmt.Println("\n📝 1. 基础占位符用法")
	fmt.Println("-------------------")

	numbers := map[string]interface{}{
		"data": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"基础过滤",
			"data | filter(# > 5)",
			"获取大于5的数字",
		},
		{
			"基础映射",
			"data | map(# * 2)",
			"每个数字乘以2",
		},
		{
			"偶数过滤",
			"data | filter(# % 2 == 0)",
			"筛选偶数",
		},
		{
			"复合条件",
			"data | filter(# % 2 == 0 && # > 3)",
			"偶数且大于3",
		},
		{
			"平方映射",
			"data | map(# * #)",
			"计算每个数字的平方",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, numbers)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstrateComplexExpressions() {
	fmt.Println("\n🧮 2. 复杂表达式")
	fmt.Println("---------------")

	numbers := map[string]interface{}{
		"data": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"模运算",
			"data | filter(# % 3 == 0)",
			"3的倍数",
		},
		{
			"复杂算术",
			"data | map(# * 2 + 1)",
			"乘以2再加1",
		},
		{
			"嵌套运算",
			"data | filter((# + 1) * 2 > 10)",
			"(x+1)*2 > 10的数字",
		},
		{
			"范围过滤",
			"data | filter(# >= 3 && # <= 7)",
			"3到7之间的数字",
		},
		{
			"平方加一",
			"data | filter(# % 2 == 1) | map(# * # + 1)",
			"奇数的平方加1",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, numbers)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstrateChainedPipelines() {
	fmt.Println("\n⛓️  3. 链式管道操作")
	fmt.Println("------------------")

	numbers := map[string]interface{}{
		"data": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"多级过滤",
			"data | filter(# > 3) | filter(# % 2 == 0)",
			"大于3的偶数",
		},
		{
			"过滤映射",
			"data | filter(# > 5) | map(# * 2)",
			"大于5的数字乘以2",
		},
		{
			"复合变换",
			"data | filter(# % 2 == 1) | map(# * # + 1) | filter(# > 10)",
			"奇数平方加1后大于10",
		},
		{
			"三级管道",
			"data | filter(# > 3) | map(# * 2) | filter(# % 4 == 0)",
			"大于3，乘以2，再筛选4的倍数",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, numbers)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstrateObjectProcessing() {
	fmt.Println("\n👥 4. 对象数组处理")
	fmt.Println("-----------------")

	users := map[string]interface{}{
		"people": []map[string]interface{}{
			{"name": "Alice", "age": 30, "salary": 75000, "department": "Engineering"},
			{"name": "Bob", "age": 25, "salary": 65000, "department": "Sales"},
			{"name": "Charlie", "age": 35, "salary": 85000, "department": "Engineering"},
			{"name": "Diana", "age": 28, "salary": 70000, "department": "Marketing"},
			{"name": "Eve", "age": 32, "salary": 90000, "department": "Engineering"},
		},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"年龄过滤",
			"people | filter(#.age >= 30) | map(#.name)",
			"30岁及以上员工姓名",
		},
		{
			"部门筛选",
			"people | filter(#.department == 'Engineering') | map(#.name)",
			"工程部员工姓名",
		},
		{
			"高薪员工",
			"people | filter(#.salary > 70000) | map(#.name)",
			"薪资超过7万的员工",
		},
		{
			"复合条件",
			"people | filter(#.age >= 30 && #.salary > 75000) | map(#.name)",
			"30岁以上且高薪员工",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, users)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstrateRealWorldExamples() {
	fmt.Println("\n🌍 5. 实际应用场景")
	fmt.Println("-----------------")

	// 电商订单数据
	orders := map[string]interface{}{
		"orders": []map[string]interface{}{
			{"id": "001", "amount": 120.50, "status": "completed", "items": 3, "customer": "premium"},
			{"id": "002", "amount": 89.99, "status": "pending", "items": 2, "customer": "regular"},
			{"id": "003", "amount": 250.00, "status": "completed", "items": 5, "customer": "premium"},
			{"id": "004", "amount": 45.00, "status": "cancelled", "items": 1, "customer": "regular"},
			{"id": "005", "amount": 180.75, "status": "completed", "items": 4, "customer": "premium"},
		},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"大额订单",
			"orders | filter(#.amount > 100) | map(#.id)",
			"金额超过100的订单ID",
		},
		{
			"已完成订单",
			"orders | filter(#.status == 'completed') | map(#.amount)",
			"已完成订单的金额",
		},
		{
			"高价值客户",
			"orders | filter(#.customer == 'premium' && #.amount > 150) | map(#.id)",
			"高价值客户的大额订单",
		},
		{
			"多商品订单",
			"orders | filter(#.items >= 3) | map({id: #.id, value: #.amount})",
			"多商品订单的ID和金额",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, orders)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstratePerformanceComparison() {
	fmt.Println("\n⚡ 6. 性能对比")
	fmt.Println("-------------")

	// 生成大数据集
	largeData := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		largeData[i] = i + 1
	}

	env := map[string]interface{}{
		"data": largeData,
	}

	// 测试表达式
	placeholderExpr := "data | filter(# % 2 == 0) | map(# * 2) | filter(# > 1000)"

	// 预编译
	program, err := expr.Compile(placeholderExpr)
	if err != nil {
		fmt.Printf("编译失败: %v\n", err)
		return
	}

	// 性能测试
	iterations := 100

	// 测试解释执行
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := expr.Eval(placeholderExpr, env)
		if err != nil {
			fmt.Printf("执行失败: %v\n", err)
			return
		}
	}
	interpretTime := time.Since(start)

	// 测试编译执行
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("执行失败: %v\n", err)
			return
		}
	}
	compiledTime := time.Since(start)

	fmt.Printf("  📊 数据规模: %d 个元素\n", len(largeData))
	fmt.Printf("  🔄 迭代次数: %d 次\n", iterations)
	fmt.Printf("  ⏱️  解释执行: %v\n", interpretTime)
	fmt.Printf("  ⚡ 编译执行: %v\n", compiledTime)
	fmt.Printf("  📈 性能提升: %.2fx\n", float64(interpretTime.Nanoseconds())/float64(compiledTime.Nanoseconds()))
	fmt.Printf("  🎯 表达式: %s\n", placeholderExpr)

	// 验证结果
	result, _ := expr.Run(program, env)
	if arr, ok := result.([]interface{}); ok {
		fmt.Printf("  ✅ 结果数量: %d 个元素\n", len(arr))
		if len(arr) > 0 {
			fmt.Printf("  📋 前5个结果: %v\n", arr[:min(5, len(arr))])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
