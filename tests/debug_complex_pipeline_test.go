package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugComplexPipeline 调试复杂的pipeline表达式
func TestDebugComplexPipeline(t *testing.T) {
	fmt.Println("🔍 调试复杂pipeline表达式")

	// 测试1：简单的类型方法调用 - 应该工作
	fmt.Println("\n1. 简单类型方法调用:")
	result1, err1 := expr.Eval(`["hello", "world"] | map(#.upper())`, nil)
	if err1 != nil {
		fmt.Printf("❌ 失败: %v\n", err1)
	} else {
		fmt.Printf("✅ 成功: %v\n", result1)
	}

	// 测试2：类型方法在比较表达式中 - 可能有问题
	fmt.Println("\n2. 类型方法在比较表达式中:")
	env := map[string]interface{}{
		"words": []string{"hi", "hello", "world", "a"},
	}

	result2, err2 := expr.Eval(`words | filter(#.length() > 4)`, env)
	if err2 != nil {
		fmt.Printf("❌ 失败: %v\n", err2)
	} else {
		fmt.Printf("✅ 成功: %v\n", result2)
	}

	// 测试3：手动验证length方法
	fmt.Println("\n3. 手动验证length方法:")
	result3, err3 := expr.Eval(`words | map(#.length())`, env)
	if err3 != nil {
		fmt.Printf("❌ 失败: %v\n", err3)
	} else {
		fmt.Printf("✅ 成功: %v\n", result3)
	}

	// 测试4：简化的比较
	fmt.Println("\n4. 简化的比较:")
	result4, err4 := expr.Eval(`[1, 5, 3, 7] | filter(# > 4)`, nil)
	if err4 != nil {
		fmt.Printf("❌ 失败: %v\n", err4)
	} else {
		fmt.Printf("✅ 成功: %v\n", result4)
	}
}
