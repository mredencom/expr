package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugLengthComparison 调试长度比较问题
func TestDebugLengthComparison(t *testing.T) {
	fmt.Println("🔍 调试长度比较问题")

	env := map[string]interface{}{
		"words": []string{"hi", "hello", "world"},
	}

	// 测试1：直接获取长度
	fmt.Println("\n1. 测试直接获取长度:")
	result1, err1 := expr.Eval(`words | map(#.length())`, env)
	if err1 != nil {
		fmt.Printf("❌ 失败: %v\n", err1)
	} else {
		fmt.Printf("✅ 成功: %v\n", result1)
	}

	// 测试2：简单的数值比较
	fmt.Println("\n2. 测试简单数值比较:")
	result2, err2 := expr.Eval(`[2, 5, 5] | filter(# > 4)`, env)
	if err2 != nil {
		fmt.Printf("❌ 失败: %v\n", err2)
	} else {
		fmt.Printf("✅ 成功: %v\n", result2)
	}

	// 测试3：复杂的长度比较
	fmt.Println("\n3. 测试复杂长度比较:")
	result3, err3 := expr.Eval(`words | filter(#.length() > 4)`, env)
	if err3 != nil {
		fmt.Printf("❌ 失败: %v\n", err3)
	} else {
		fmt.Printf("✅ 成功: %v\n", result3)
	}

	// 测试4: 手动验证预期结果
	fmt.Println("\n4. 手动验证字符串长度:")
	words := []string{"hi", "hello", "world"}
	for i, word := range words {
		length := len(word)
		shouldInclude := length > 4
		fmt.Printf("   words[%d] = \"%s\", length = %d, > 4 = %t\n", i, word, length, shouldInclude)
	}
}
