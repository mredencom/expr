package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestDebugMethodVsProperty 测试方法调用和属性访问的差异
func TestDebugMethodVsProperty(t *testing.T) {
	fmt.Println("🔍 测试方法调用 vs 属性访问的编译差异")

	env := map[string]interface{}{
		"words": []string{"hi", "hello", "world"},
	}

	fmt.Printf("环境: %v\n", env)

	// 测试1：属性访问 (如果存在)
	fmt.Printf("\n1. 测试属性访问 #.length (无括号):\n")
	result1, err1 := expr.Eval(`words | map(#.length)`, env)
	if err1 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err1)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result1)
	}

	// 测试2：方法调用
	fmt.Printf("\n2. 测试方法调用 #.length() (有括号):\n")
	result2, err2 := expr.Eval(`words | map(#.length())`, env)
	if err2 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err2)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result2)
	}

	// 测试3：复杂表达式中的属性访问
	fmt.Printf("\n3. 测试复杂表达式中的属性访问 #.length > 4 (无括号):\n")
	result3, err3 := expr.Eval(`words | filter(#.length > 4)`, env)
	if err3 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err3)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result3)
	}

	// 测试4：复杂表达式中的方法调用
	fmt.Printf("\n4. 测试复杂表达式中的方法调用 #.length() > 4 (有括号):\n")
	result4, err4 := expr.Eval(`words | filter(#.length() > 4)`, env)
	if err4 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err4)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result4)
	}

	// 测试5：直接数组测试方法调用
	fmt.Printf("\n5. 测试直接数组中的方法调用:\n")
	result5, err5 := expr.Eval(`[\"hi\", \"hello\", \"world\"] | filter(#.length() > 4)`, env)
	if err5 != nil {
		fmt.Printf("   ❌ 失败: %v\n", err5)
	} else {
		fmt.Printf("   ✅ 成功: %v\n", result5)
	}

	// 测试6：另一个字符串方法的对比
	fmt.Printf("\n6. 测试其他字符串方法:\n")

	// 属性访问方式
	result6a, err6a := expr.Eval(`words | filter(#.contains("e"))`, env)
	if err6a != nil {
		fmt.Printf("   #.contains(\"e\") ❌: %v\n", err6a)
	} else {
		fmt.Printf("   #.contains(\"e\") ✅: %v\n", result6a)
	}

	// 测试7：验证我们的假设
	fmt.Printf("\n7. 验证假设 - length是属性还是方法？:\n")

	// 测试对象是否有length属性
	result7, err7 := expr.Eval(`\"hello\".length`, env)
	if err7 != nil {
		fmt.Printf("   \"hello\".length (属性访问) ❌: %v\n", err7)
	} else {
		fmt.Printf("   \"hello\".length (属性访问) ✅: %v\n", result7)
	}

	// 对比方法调用
	result8, err8 := expr.Eval(`\"hello\".length()`, env)
	if err8 != nil {
		fmt.Printf("   \"hello\".length() (方法调用) ❌: %v\n", err8)
	} else {
		fmt.Printf("   \"hello\".length() (方法调用) ✅: %v\n", result8)
	}
}
