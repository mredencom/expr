package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestFinalComprehensiveTypeMethodPipeline 最终全面测试
func TestFinalComprehensiveTypeMethodPipeline(t *testing.T) {
	fmt.Println("🎯 最终全面测试：TypeMethod与Pipeline的结合")
	fmt.Println("==============================================")

	env := map[string]interface{}{
		"words":   []string{"hi", "hello", "world"},
		"numbers": []int{1, 2, 3, 4, 5},
		"prices":  []float64{10.5, 20.0, 15.75},
		"flags":   []bool{true, false, true},
		"mixed":   []interface{}{"abc", 123, true},
	}

	tests := []struct {
		name       string
		expression string
		expected   interface{}
		category   string
	}{
		// 字符串类型方法
		{"字符串转大写", `words | map(#.upper())`, []interface{}{"HI", "HELLO", "WORLD"}, "string"},
		{"字符串转小写", `words | map(#.lower())`, []interface{}{"hi", "hello", "world"}, "string"},
		{"字符串长度", `words | map(#.length())`, []interface{}{2, 5, 5}, "string"},
		{"字符串包含", `words | filter(#.contains("e"))`, []interface{}{"hello"}, "string"},

		// 复杂表达式过滤
		{"长度过滤>4", `words | filter(#.length() > 4)`, []interface{}{"hello", "world"}, "complex"},
		{"长度过滤>3", `words | filter(#.length() > 3)`, []interface{}{"hello", "world"}, "complex"},
		{"长度过滤>2", `words | filter(#.length() > 2)`, []interface{}{"hello", "world"}, "complex"},
		{"长度过滤==2", `words | filter(#.length() == 2)`, []interface{}{"hi"}, "complex"},

		// 链式操作
		{"链式长度过滤+转大写", `words | filter(#.length() > 3) | map(#.upper())`, []interface{}{"HELLO", "WORLD"}, "chain"},
		{"链式转大写+长度过滤", `words | map(#.upper()) | filter(#.length() > 3)`, []interface{}{"HELLO", "WORLD"}, "chain"},

		// 数值类型方法
		{"数值绝对值", `[-1, 2, -3] | map(#.abs())`, []interface{}{1, 2, 3}, "numeric"},
		{"数值过滤>3", `numbers | filter(# > 3)`, []interface{}{4, 5}, "numeric"},

		// 布尔值过滤
		{"布尔值过滤", `flags | filter(#)`, []interface{}{true, true}, "boolean"},
		{"布尔值反转", `flags | map(!#)`, []interface{}{false, true, false}, "boolean"},

		// 组合操作
		{"字符串替换", `words | map(#.replace("l", "L"))`, []interface{}{"hi", "heLLo", "worLd"}, "replace"},
		{"起始字符判断", `words | filter(#.startsWith("h"))`, []interface{}{"hi", "hello"}, "startsWith"},
		{"结束字符判断", `words | filter(#.endsWith("o"))`, []interface{}{"hello"}, "endsWith"},

		// 高级组合
		{"复杂链式操作", `words | filter(#.length() > 2) | map(#.upper()) | filter(#.startsWith("H"))`, []interface{}{"HELLO"}, "advanced"},
	}

	var passed, failed int
	categoryStats := make(map[string][]bool)

	for _, test := range tests {
		result, err := expr.Eval(test.expression, env)
		success := false

		if err != nil {
			fmt.Printf("❌ %s: 错误 - %v\n", test.name, err)
		} else {
			if fmt.Sprintf("%v", result) == fmt.Sprintf("%v", test.expected) {
				fmt.Printf("✅ %s: %v\n", test.name, result)
				success = true
				passed++
			} else {
				fmt.Printf("❌ %s: 期望 %v, 实际 %v\n", test.name, test.expected, result)
				failed++
			}
		}

		categoryStats[test.category] = append(categoryStats[test.category], success)
	}

	// 分类统计
	fmt.Printf("\n📊 分类统计:\n")
	for category, results := range categoryStats {
		successCount := 0
		for _, success := range results {
			if success {
				successCount++
			}
		}
		percentage := float64(successCount) * 100 / float64(len(results))
		fmt.Printf("   %s: %d/%d (%.1f%%)\n", category, successCount, len(results), percentage)
	}

	// 总体统计
	total := passed + failed
	percentage := float64(passed) * 100 / float64(total)

	fmt.Printf("\n🏆 最终结果:\n")
	fmt.Printf("   通过: %d\n", passed)
	fmt.Printf("   失败: %d\n", failed)
	fmt.Printf("   总计: %d\n", total)
	fmt.Printf("   成功率: %.1f%%\n", percentage)

	if percentage == 100.0 {
		fmt.Printf("\n🎉 恭喜！TypeMethod与Pipeline完美结合，达到100%%通过率！\n")
	} else if percentage >= 90.0 {
		fmt.Printf("\n🎯 优秀！接近完美的实现，成功率超过90%%！\n")
	} else {
		fmt.Printf("\n⚠️  还需要继续完善，目标是100%%通过率\n")
	}

	// 测试一些边界情况
	fmt.Printf("\n🔍 边界情况测试:\n")
	edgeCases := []struct {
		name string
		expr string
	}{
		{"空数组", `[] | filter(#.length() > 0)`},
		{"单元素数组", `["test"] | filter(#.length() > 3)`},
		{"混合类型", `["abc", "defgh"] | filter(#.length() > 3)`},
	}

	for _, edge := range edgeCases {
		result, err := expr.Eval(edge.expr, env)
		if err != nil {
			fmt.Printf("   ❌ %s: %v\n", edge.name, err)
		} else {
			fmt.Printf("   ✅ %s: %v\n", edge.name, result)
		}
	}
}
