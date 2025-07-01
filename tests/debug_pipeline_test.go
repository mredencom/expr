package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugPipeline 调试管道操作问题
func TestDebugPipeline(t *testing.T) {
	fmt.Println("🔍 调试管道操作问题")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	// 测试1: 简单过滤
	fmt.Println("\n1. 测试简单过滤: numbers | filter(# > 5)")
	result, err := expr.Eval("numbers | filter(# > 5)", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
		fmt.Printf("   期望: [6, 7, 8, 9, 10]\n")
	}

	// 测试2: 简单映射
	fmt.Println("\n2. 测试简单映射: numbers | map(# * 2)")
	result, err = expr.Eval("numbers | map(# * 2)", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
		fmt.Printf("   期望: [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]\n")
	}

	// 测试3: 链式操作
	fmt.Println("\n3. 测试链式操作: numbers | filter(# > 3) | map(# * 2)")
	result, err = expr.Eval("numbers | filter(# > 3) | map(# * 2)", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
		fmt.Printf("   期望: [8, 10, 12, 14, 16, 18, 20]\n")
	}

	// 测试4: 分析内置函数
	fmt.Println("\n4. 测试内置函数: filter(numbers, # > 5)")
	result, err = expr.Eval("filter(numbers, # > 5)", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
		fmt.Printf("   期望: [6, 7, 8, 9, 10]\n")
	}

	// 测试5: 分析编译过程
	fmt.Println("\n5. 分析编译过程")
	program, err := expr.Compile("numbers | filter(# > 5)", expr.Env(env))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功\n")
		fmt.Printf("   字节码大小: %d\n", program.BytecodeSize())
		fmt.Printf("   常量数量: %d\n", program.ConstantsCount())
	}
}
