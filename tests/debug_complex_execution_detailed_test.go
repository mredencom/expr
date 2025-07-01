package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugComplexExecutionDetailed 详细调试复杂表达式执行
func TestDebugComplexExecutionDetailed(t *testing.T) {
	fmt.Println("🔍 详细调试复杂表达式执行过程")

	exprStr := `words | filter(#.length() > 4)`
	env := map[string]interface{}{
		"words": []string{"hi", "hello", "world"},
	}

	fmt.Printf("表达式: %s\n", exprStr)
	fmt.Printf("环境: %v\n", env)

	// 测试各个组成部分
	fmt.Printf("\n🧪 分步测试:\n")

	// 1. 测试基础数据
	fmt.Printf("1. 基础数据访问:\n")
	result1, err1 := expr.Eval(`words`, env)
	if err1 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err1)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result1)
	}

	// 2. 测试长度方法
	fmt.Printf("2. 长度方法调用:\n")
	result2, err2 := expr.Eval(`words | map(#.length())`, env)
	if err2 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err2)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result2)
	}

	// 3. 测试数值比较
	fmt.Printf("3. 数值比较:\n")
	result3, err3 := expr.Eval(`[2, 5, 5] | filter(# > 4)`, env)
	if err3 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err3)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result3)
	}

	// 4. 测试简化的复杂表达式
	fmt.Printf("4. 简化的字符串长度过滤:\n")
	result4, err4 := expr.Eval(`["hi", "hello", "world"] | filter(#.length() > 4)`, env)
	if err4 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err4)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result4)
	}

	// 5. 测试目标表达式
	fmt.Printf("5. 目标表达式:\n")
	result5, err5 := expr.Eval(exprStr, env)
	if err5 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err5)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result5)
	}

	// 6. 手动验证预期
	fmt.Printf("\n📝 手动验证预期结果:\n")
	words := []string{"hi", "hello", "world"}
	var expected []string
	for _, word := range words {
		if len(word) > 4 {
			expected = append(expected, word)
		}
	}
	fmt.Printf("   预期结果: %v\n", expected)

	// 7. 测试其他已工作的表达式
	fmt.Printf("\n✅ 测试已知工作的表达式:\n")
	result6, err6 := expr.Eval(`words | filter(#.contains("e"))`, env)
	if err6 != nil {
		fmt.Printf("   contains 过滤 ❌: %v\n", err6)
	} else {
		fmt.Printf("   contains 过滤 ✅: %v\n", result6)
	}
}
