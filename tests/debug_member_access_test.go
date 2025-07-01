package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugMemberAccess 调试成员访问问题
func TestDebugMemberAccess(t *testing.T) {
	fmt.Println("🔍 调试成员访问问题")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	// 简单的成员访问测试
	env := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	// 测试1: 直接成员访问
	fmt.Println("\n1. 测试直接成员访问: user.name")
	result, err := expr.Eval("user.name", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	// 测试2: 数组中的成员访问
	usersEnv := map[string]interface{}{
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
		},
	}

	fmt.Println("\n2. 测试数组索引: users[0]")
	result, err = expr.Eval("users[0]", usersEnv)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	fmt.Println("\n3. 测试组合访问: users[0].name")
	result, err = expr.Eval("users[0].name", usersEnv)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	// 测试3: 分析编译过程
	fmt.Println("\n4. 分析编译过程")
	program, err := expr.Compile("users[0].name", expr.Env(usersEnv))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功\n")
		fmt.Printf("   字节码大小: %d\n", program.BytecodeSize())
		fmt.Printf("   常量数量: %d\n", program.ConstantsCount())

		// 尝试执行
		result, err := expr.Run(program, usersEnv)
		if err != nil {
			fmt.Printf("❌ 执行错误: %v\n", err)
		} else {
			fmt.Printf("✅ 执行结果: %v\n", result)
		}
	}
}
