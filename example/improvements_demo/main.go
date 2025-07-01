package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

func main() {
	fmt.Println("🚀 表达式引擎改进功能演示")
	fmt.Println("================================")

	demonstrateSingleQuotes()
	demonstrateWildcards()
	demonstrateEnhancedArrayAccess()
	demonstratePipelineOperations()
	demonstrateBitwiseOperations()
	demonstrateRealWorldUsage()
}

func demonstrateSingleQuotes() {
	fmt.Println("\n📝 1. 单引号字符串支持 (避免转义)")
	fmt.Println("--------------------------------")

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"基础单引号",
			"'Hello, World!'",
			"使用单引号定义字符串",
		},
		{
			"避免双引号转义",
			"'He said \"Hello!\" to me'",
			"单引号内可以直接使用双引号",
		},
		{
			"单引号转义",
			"'It\\'s a beautiful day'",
			"只需要转义单引号本身",
		},
		{
			"混合使用",
			"'JSON: {\"name\": \"Alice\", \"age\": 30}'",
			"处理JSON字符串非常方便",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, nil)
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-15s: %s\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstrateWildcards() {
	fmt.Println("\n🔍 2. 通配符支持")
	fmt.Println("----------------")

	examples := []string{
		"user.*",
		"*.field",
		"data.*.name",
		"config.*.settings.*",
	}

	fmt.Println("  通配符语法解析测试:")
	for _, expr := range examples {
		l := lexer.New(expr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("  ❌ %s: 解析错误\n", expr)
		} else {
			fmt.Printf("  ✅ %-20s: 解析成功\n", expr)
			if len(program.Statements) > 0 {
				stmt := program.Statements[0].(*ast.ExpressionStatement)
				fmt.Printf("     AST: %s\n", stmt.Expression.String())
			}
		}
	}

	fmt.Println("\n  通配符应用场景:")
	scenarios := []struct {
		scenario string
		example  string
		use_case string
	}{
		{
			"对象属性提取",
			"user.*",
			"获取用户对象的所有属性",
		},
		{
			"动态字段访问",
			"*.name",
			"访问任意对象的name字段",
		},
		{
			"嵌套通配符",
			"data.*.config.*",
			"多级通配符访问",
		},
		{
			"管道中使用",
			"users | map(u => u.*)",
			"在管道操作中使用通配符",
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("  📋 %-15s: %s\n", scenario.scenario, scenario.example)
		fmt.Printf("     用途: %s\n", scenario.use_case)
	}
}

func demonstrateEnhancedArrayAccess() {
	fmt.Println("\n🔢 3. 增强的数组访问")
	fmt.Println("--------------------")

	env := map[string]interface{}{
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 30, "city": "NYC"},
			{"name": "Bob", "age": 25, "city": "LA"},
			{"name": "Charlie", "age": 35, "city": "Chicago"},
		},
		"matrix": [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		"config": map[string]interface{}{
			"servers": []string{"web1", "web2", "web3"},
			"ports":   []int{8080, 8081, 8082},
		},
	}

	examples := []struct {
		name        string
		expression  string
		description string
	}{
		{
			"基础数组访问",
			"users[0].name",
			"访问第一个用户的姓名",
		},
		{
			"嵌套数组访问",
			"matrix[1][2]",
			"访问二维数组元素",
		},
		{
			"配置数组访问",
			"config.servers[0]",
			"访问配置中的第一个服务器",
		},
		{
			"数组长度",
			"len(users)",
			"获取用户数组长度",
		},
		{
			"最后一个元素",
			"users[len(users)-1].name",
			"访问最后一个用户的姓名",
		},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, env)
		if err != nil {
			fmt.Printf("  ❌ %-15s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-15s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s\n", example.expression)
			fmt.Printf("     说明: %s\n", example.description)
			fmt.Println()
		}
	}
}

func demonstratePipelineOperations() {
	fmt.Println("\n🔄 4. 管道操作 (智能|符号处理)")
	fmt.Println("-----------------------------")

	// 展示管道操作的解析
	pipelineExamples := []string{
		"data | filter(x => x > 5)",
		"numbers | map(n => n * 2)",
		"users | filter(u => u.age > 18) | map(u => u.name)",
		"data | filter(# > 5)",                     // 占位符语法
		"numbers | map(# * 2)",                     // 占位符语法
		"users | filter(#.age > 18) | map(#.name)", // 占位符语法
	}

	fmt.Println("  管道表达式解析测试:")
	for _, expr := range pipelineExamples {
		l := lexer.New(expr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("  ❌ %s: 解析错误 - %v\n", expr, p.Errors()[0])
		} else {
			fmt.Printf("  ✅ %-40s: 解析成功\n", expr)
			if len(program.Statements) > 0 {
				stmt := program.Statements[0].(*ast.ExpressionStatement)
				fmt.Printf("     AST: %s\n", stmt.Expression.String())
			}
		}
	}

	// 实际执行示例
	fmt.Println("\n  管道操作执行示例:")

	data := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 17},
			{"name": "Charlie", "age": 25},
		},
	}

	executionExamples := []struct {
		name       string
		expression string
		syntax     string
	}{
		{
			"基础过滤",
			"numbers | filter(# > 5)",
			"占位符语法",
		},
		{
			"数值映射",
			"numbers | map(# * 2)",
			"占位符语法",
		},
		{
			"复合条件",
			"numbers | filter(# % 2 == 0 && # > 3)",
			"占位符语法",
		},
		{
			"链式操作",
			"numbers | filter(# > 3) | map(# * 2)",
			"占位符语法",
		},
		{
			"对象过滤",
			"users | filter(#.age >= 18) | map(#.name)",
			"占位符语法",
		},
	}

	for _, example := range executionExamples {
		result, err := expr.Eval(example.expression, data)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", example.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: %v\n", example.name, result)
			fmt.Printf("     表达式: %s (%s)\n", example.expression, example.syntax)
		}
		fmt.Println()
	}
}

func demonstrateBitwiseOperations() {
	fmt.Println("\n⚡ 5. 位运算操作 (保持兼容性)")
	fmt.Println("---------------------------")

	examples := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"位或运算", "5 | 3", int64(7)},
		{"位与运算", "5 & 3", int64(1)},
		{"位异或运算", "5 ^ 3", int64(6)},
		{"左移运算", "5 << 1", int64(10)},
		{"右移运算", "10 >> 1", int64(5)},
		{"位非运算", "~5", int64(-6)},
	}

	for _, example := range examples {
		result, err := expr.Eval(example.expression, nil)
		if err != nil {
			fmt.Printf("  ❌ %-10s: %v\n", example.name, err)
		} else {
			status := "✅"
			if result != example.expected {
				status = "❌"
			}
			fmt.Printf("  %s %-10s: %s = %v (期望: %v)\n",
				status, example.name, example.expression, result, example.expected)
		}
	}
}

func demonstrateRealWorldUsage() {
	fmt.Println("\n🌍 6. 实际应用场景")
	fmt.Println("------------------")

	// 配置管理场景
	fmt.Println("  📋 配置管理:")
	configEnv := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "MyApp",
			"version": "1.2.3",
			"debug":   true,
		},
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"ssl":  false,
		},
	}

	configExpressions := []string{
		"app.name + ' v' + app.version",
		"database.host + ':' + string(database.port)",
		"app.debug ? 'Development' : 'Production'",
	}

	for _, exprStr := range configExpressions {
		result, err := expr.Eval(exprStr, configEnv)
		if err != nil {
			fmt.Printf("    ❌ %s\n", err)
		} else {
			fmt.Printf("    ✅ %s → %v\n", exprStr, result)
		}
	}

	// 业务规则场景
	fmt.Println("\n  💼 业务规则:")
	businessEnv := map[string]interface{}{
		"user": map[string]interface{}{
			"age":        28,
			"membership": "premium",
			"totalSpent": 1500.0,
			"country":    "US",
		},
		"order": map[string]interface{}{
			"amount":      250.0,
			"items":       3,
			"destination": "domestic",
		},
	}

	businessRules := []struct {
		rule        string
		expression  string
		description string
	}{
		{
			"VIP用户检查",
			"user.membership == 'premium' && user.totalSpent > 1000",
			"检查是否为VIP用户",
		},
		{
			"免费配送",
			"order.amount > 100 || (user.membership == 'premium' && order.destination == 'domestic')",
			"确定是否符合免费配送条件",
		},
		{
			"折扣计算",
			"user.age > 25 ? (user.membership == 'premium' ? 0.15 : 0.10) : 0.05",
			"根据年龄和会员级别计算折扣",
		},
	}

	for _, rule := range businessRules {
		result, err := expr.Eval(rule.expression, businessEnv)
		if err != nil {
			fmt.Printf("    ❌ %s: %v\n", rule.rule, err)
		} else {
			fmt.Printf("    ✅ %-12s: %v\n", rule.rule, result)
			fmt.Printf("       表达式: %s\n", rule.expression)
			fmt.Printf("       说明: %s\n", rule.description)
			fmt.Println()
		}
	}

	fmt.Println("🎉 改进功能演示完成!")
	fmt.Println("\n主要改进总结:")
	fmt.Println("  • 单引号字符串 - 减少转义，提高可读性")
	fmt.Println("  • 通配符支持 - 灵活的对象属性访问")
	fmt.Println("  • 增强数组访问 - 更自然的语法")
	fmt.Println("  • 智能管道操作 - 上下文敏感的|符号处理")
	fmt.Println("  • 位运算兼容 - 保持原有功能不变")
}
