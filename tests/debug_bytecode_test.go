package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugBytecode 调试字节码生成
func TestDebugBytecode(t *testing.T) {
	fmt.Println("🔍 调试字节码生成")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	// 测试1: 简单管道操作
	fmt.Println("\n1. 分析: numbers | filter(# > 5)")
	program, err := expr.Compile("numbers | filter(# > 5)", expr.Env(env))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 编译成功\n")
	fmt.Printf("   字节码大小: %d\n", program.BytecodeSize())
	fmt.Printf("   常量数量: %d\n", program.ConstantsCount())

	// 简化：只显示基本信息，不深入字节码细节

	// 测试执行
	fmt.Println("\n2. 执行测试")
	result, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("❌ 执行错误: %v\n", err)
	} else {
		fmt.Printf("✅ 执行结果: %v\n", result)
	}

	// 测试2: 简单变量访问
	fmt.Println("\n3. 对比: numbers (变量访问)")
	program2, err := expr.Compile("numbers", expr.Env(env))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 编译成功\n")
	fmt.Printf("   字节码大小: %d\n", program2.BytecodeSize())
	fmt.Printf("   常量数量: %d\n", program2.ConstantsCount())

	result2, err := expr.Run(program2, env)
	if err != nil {
		fmt.Printf("❌ 执行错误: %v\n", err)
	} else {
		fmt.Printf("✅ 执行结果: %v\n", result2)
	}

	// 检查是否是同一个结果（说明管道没有工作）
	if fmt.Sprintf("%v", result) == fmt.Sprintf("%v", result2) {
		fmt.Printf("\n🚨 问题确认: 管道操作和变量访问返回相同结果！\n")
		fmt.Printf("   这说明管道操作根本没有执行\n")
	}
}
