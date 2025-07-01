package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugLiterals 调试数组和对象字面量问题
func TestDebugLiterals(t *testing.T) {
	fmt.Println("🔍 调试数组和对象字面量问题")
	fmt.Println("=" + fmt.Sprintf("%40s", "="))

	// 测试1: 简单数组字面量
	fmt.Println("\n1. 测试简单数组字面量: [1, 2, 3]")
	result, err := expr.Eval("[1, 2, 3]", nil)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	// 测试2: 简单对象字面量
	fmt.Println("\n2. 测试简单对象字面量: {\"name\": \"Alice\"}")
	result, err = expr.Eval(`{"name": "Alice"}`, nil)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	// 测试3: 分析编译过程 - 数组
	fmt.Println("\n3. 分析数组编译过程")
	program, err := expr.Compile("[1, 2, 3]")
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功\n")
		fmt.Printf("   字节码大小: %d\n", program.BytecodeSize())
		fmt.Printf("   常量数量: %d\n", program.ConstantsCount())

		// 尝试执行
		result, err := expr.Run(program, nil)
		if err != nil {
			fmt.Printf("❌ 执行错误: %v\n", err)
		} else {
			fmt.Printf("✅ 执行结果: %v\n", result)
		}
	}

	// 测试4: 分析编译过程 - 对象
	fmt.Println("\n4. 分析对象编译过程")
	program, err = expr.Compile(`{"name": "Alice"}`)
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功\n")
		fmt.Printf("   字节码大小: %d\n", program.BytecodeSize())
		fmt.Printf("   常量数量: %d\n", program.ConstantsCount())

		// 尝试执行
		result, err := expr.Run(program, nil)
		if err != nil {
			fmt.Printf("❌ 执行错误: %v\n", err)
		} else {
			fmt.Printf("✅ 执行结果: %v\n", result)
		}
	}

	// 测试5: 更简单的常量测试
	fmt.Println("\n5. 测试简单常量: 42")
	result, err = expr.Eval("42", nil)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}
}
