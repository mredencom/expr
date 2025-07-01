package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

func TestNullSafetyOperators(t *testing.T) {
	fmt.Println("🔒 空值安全操作符测试")
	fmt.Println("========================")

	env := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"profile": map[string]interface{}{
				"bio": "Software Engineer",
			},
		},
		"emptyUser": nil,
		"config": map[string]interface{}{
			"timeout": 30,
		},
		"nullValue":      nil,
		"defaultTimeout": 60,
	}

	tests := []struct {
		name       string
		expression string
		expected   interface{}
		shouldPass bool
	}{
		// Optional chaining tests
		{"Optional chaining - valid path", `user?.name`, "Alice", true},
		{"Optional chaining - nested valid", `user?.profile?.bio`, "Software Engineer", true},
		{"Optional chaining - null object", `emptyUser?.name`, nil, true},
		{"Optional chaining - missing property", `user?.address?.city`, nil, true},

		// Null coalescing tests
		{"Null coalescing - use left value", `config?.timeout ?? 45`, 30, true},
		{"Null coalescing - use default", `nullValue ?? 'default'`, "default", true},
		{"Null coalescing - chain with optional", `emptyUser?.name ?? 'Anonymous'`, "Anonymous", true},
		{"Null coalescing - nested", `user?.profile?.email ?? user?.email ?? 'no-email'`, "no-email", true},

		// Combined operations
		{"Combined - optional + coalescing", `user?.profile?.timeout ?? config?.timeout ?? defaultTimeout`, 30, true},
		{"Combined - complex chain", `emptyUser?.profile?.name ?? user?.name ?? 'Unknown'`, "Alice", true},
	}

	passCount := 0
	for _, test := range tests {
		fmt.Printf("%-40s: ", test.name)

		// 编译测试
		program, err := expr.Compile(test.expression, expr.Env(env))
		if err != nil {
			fmt.Printf("❌ 编译失败: %v\n", err)
			if test.shouldPass {
				t.Errorf("%s 编译失败: %v", test.name, err)
			}
			continue
		}

		// 执行测试
		result, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("❌ 执行失败: %v\n", err)
			if test.shouldPass {
				t.Errorf("%s 执行失败: %v", test.name, err)
			}
			continue
		}

		// 检查结果
		if !compareResults(result, test.expected) {
			fmt.Printf("⚠️ 结果不匹配: 期望 %v, 得到 %v\n", test.expected, result)
			if test.shouldPass {
				t.Errorf("%s 结果不匹配: 期望 %v, 得到 %v", test.name, test.expected, result)
			}
			continue
		}

		fmt.Printf("✅ 成功: %v\n", result)
		passCount++
	}

	fmt.Printf("\n空值安全操作符: %d/%d 通过\n", passCount, len(tests))
	if passCount == len(tests) {
		fmt.Println("✅ 所有空值安全操作符测试通过!")
	}
}

func compareResults(actual, expected interface{}) bool {
	// Handle nil comparisons - if expected is nil, check if actual is nil or NilValue
	if expected == nil {
		if actual == nil {
			return true
		}
		// Check if actual is a string "nil" or <nil>
		if actualStr, ok := actual.(string); ok {
			return actualStr == "nil" || actualStr == "<nil>"
		}
		return false
	}

	if actual == nil {
		return expected == nil
	}

	// Handle string comparison
	if actualStr, ok := actual.(string); ok {
		if expectedStr, ok := expected.(string); ok {
			return actualStr == expectedStr
		}
	}

	// Handle numeric comparison
	if actualNum, ok := actual.(int64); ok {
		if expectedNum, ok := expected.(int); ok {
			return actualNum == int64(expectedNum)
		}
		if expectedNum, ok := expected.(int64); ok {
			return actualNum == expectedNum
		}
	}

	return actual == expected
}
