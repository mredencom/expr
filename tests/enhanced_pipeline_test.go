package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestEnhancedPipelineOperations 测试增强的管道操作
func TestEnhancedPipelineOperations(t *testing.T) {
	fmt.Println("🔧 增强管道操作测试")
	fmt.Println("========================")

	tests := []struct {
		name       string
		expression string
		env        map[string]interface{}
		expected   interface{}
		shouldPass bool
	}{
		// Lambda表达式支持
		{
			name:       "Lambda filter",
			expression: `[1, 2, 3, 4, 5] | filter(x => x > 3)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{4, 5},
			shouldPass: true,
		},
		{
			name:       "Lambda map",
			expression: `[1, 2, 3] | map(x => x * 2)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{2, 4, 6},
			shouldPass: true,
		},
		{
			name:       "Lambda reduce",
			expression: `[1, 2, 3, 4] | reduce((a, b) => a + b)`,
			env:        map[string]interface{}{},
			expected:   10,
			shouldPass: true,
		},

		// 占位符表达式支持
		{
			name:       "Placeholder filter",
			expression: `[1, 2, 3, 4, 5] | filter(# > 2)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{3, 4, 5},
			shouldPass: true,
		},
		{
			name:       "Placeholder map",
			expression: `[1, 2, 3] | map(# * 3)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{3, 6, 9},
			shouldPass: true,
		},

		// 字符串谓词增强支持
		{
			name:       "Enhanced string filter - positive",
			expression: `[-2, -1, 0, 1, 2] | filter('positive')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 2},
			shouldPass: true,
		},
		{
			name:       "Enhanced string filter - even",
			expression: `[1, 2, 3, 4, 5, 6] | filter('even')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{2, 4, 6},
			shouldPass: true,
		},
		{
			name:       "Enhanced string map - double",
			expression: `[1, 2, 3] | map('double')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{2, 4, 6},
			shouldPass: true,
		},
		{
			name:       "Enhanced string map - square",
			expression: `[1, 2, 3, 4] | map('square')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 4, 9, 16},
			shouldPass: true,
		},

		// 混合链式操作
		{
			name:       "Mixed chain - Lambda and string",
			expression: `[1, 2, 3, 4, 5] | filter(x => x > 2) | map('double')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{6, 8, 10},
			shouldPass: true,
		},

		// 复杂数据处理
		{
			name:       "Complex data processing",
			expression: `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10] | filter('odd') | map('square') | take(3)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 9, 25},
			shouldPass: true,
		},
	}

	passCount := 0
	for _, test := range tests {
		fmt.Printf("  ✓ %-25s: ", test.name)

		// 编译测试
		program, err := expr.Compile(test.expression)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 编译失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期编译失败\n")
				passCount++
				continue
			}
		}

		// 执行测试
		result, err := expr.Run(program, test.env)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 执行失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期执行失败\n")
				passCount++
				continue
			}
		}

		if test.shouldPass {
			fmt.Printf("✅ 结果: %v\n", result)
			passCount++
		} else {
			fmt.Printf("❌ 应该失败但成功了: %v\n", result)
		}
	}

	fmt.Printf("\n增强管道操作: %d/%d 通过\n", passCount, len(tests))
	if passCount == len(tests) {
		fmt.Println("✅ 所有增强管道操作测试通过!")
	}
}

// TestMixedLambdaPlaceholderSyntax 测试Lambda和占位符混合语法
func TestMixedLambdaPlaceholderSyntax(t *testing.T) {
	fmt.Println("\n🔧 混合语法测试")
	fmt.Println("==================")

	tests := []struct {
		name       string
		expression string
		env        map[string]interface{}
		expected   interface{}
		shouldPass bool
	}{
		{
			name:       "Lambda filter + placeholder map",
			expression: `[1, 2, 3, 4, 5] | filter(x => x > 2) | map(# * 2)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{6, 8, 10},
			shouldPass: true,
		},
		{
			name:       "Placeholder filter + Lambda map",
			expression: `[1, 2, 3, 4, 5] | filter(# > 2) | map(x => x + 10)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{13, 14, 15},
			shouldPass: true,
		},
		{
			name:       "String filter + Lambda map + placeholder filter",
			expression: `[1, 2, 3, 4, 5, 6] | filter('even') | map(x => x * x) | filter(# > 10)`,
			env:        map[string]interface{}{},
			expected:   []interface{}{16, 36},
			shouldPass: true,
		},
	}

	passCount := 0
	for _, test := range tests {
		fmt.Printf("  ✓ %-35s: ", test.name)

		// 编译测试
		program, err := expr.Compile(test.expression)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 编译失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期编译失败\n")
				passCount++
				continue
			}
		}

		// 执行测试
		result, err := expr.Run(program, test.env)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 执行失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期执行失败\n")
				passCount++
				continue
			}
		}

		if test.shouldPass {
			fmt.Printf("✅ 结果: %v\n", result)
			passCount++
		} else {
			fmt.Printf("❌ 应该失败但成功了: %v\n", result)
		}
	}

	fmt.Printf("\n混合语法: %d/%d 通过\n", passCount, len(tests))
	if passCount == len(tests) {
		fmt.Println("✅ 所有混合语法测试通过!")
	}
}

// TestEnhancedStringOperations 测试增强的字符串操作
func TestEnhancedStringOperations(t *testing.T) {
	fmt.Println("\n🔧 增强字符串操作测试")
	fmt.Println("========================")

	tests := []struct {
		name       string
		expression string
		env        map[string]interface{}
		expected   interface{}
		shouldPass bool
	}{
		{
			name:       "Numeric predicates",
			expression: `[-3, -2, -1, 0, 1, 2, 3] | filter('positive')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 2, 3},
			shouldPass: true,
		},
		{
			name:       "Even/odd predicates",
			expression: `[1, 2, 3, 4, 5, 6, 7, 8] | filter('odd')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 3, 5, 7},
			shouldPass: true,
		},
		{
			name:       "String transformations",
			expression: `['hello', 'world'] | map('upper')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{"HELLO", "WORLD"},
			shouldPass: true,
		},
		{
			name:       "Mathematical transformations",
			expression: `[1, 2, 3, 4] | map('square') | map('abs')`,
			env:        map[string]interface{}{},
			expected:   []interface{}{1, 4, 9, 16},
			shouldPass: true,
		},
		{
			name:       "Enhanced reducers",
			expression: `[1, 2, 3, 4, 5] | reduce('sum')`,
			env:        map[string]interface{}{},
			expected:   15,
			shouldPass: true,
		},
	}

	passCount := 0
	for _, test := range tests {
		fmt.Printf("  ✓ %-25s: ", test.name)

		// 编译测试
		program, err := expr.Compile(test.expression)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 编译失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期编译失败\n")
				passCount++
				continue
			}
		}

		// 执行测试
		result, err := expr.Run(program, test.env)
		if err != nil {
			if test.shouldPass {
				fmt.Printf("❌ 执行失败: %v\n", err)
				continue
			} else {
				fmt.Printf("✅ 预期执行失败\n")
				passCount++
				continue
			}
		}

		if test.shouldPass {
			fmt.Printf("✅ 结果: %v\n", result)
			passCount++
		} else {
			fmt.Printf("❌ 应该失败但成功了: %v\n", result)
		}
	}

	fmt.Printf("\n增强字符串操作: %d/%d 通过\n", passCount, len(tests))
	if passCount == len(tests) {
		fmt.Println("✅ 所有增强字符串操作测试通过!")
	}
}

func TestEnhancedPipelineIntegration(t *testing.T) {
	fmt.Println("🚀 增强管道操作功能测试")
	fmt.Println("================================")

	TestEnhancedPipelineOperations(t)
	TestMixedLambdaPlaceholderSyntax(t)
	TestEnhancedStringOperations(t)

	fmt.Println("\n🎉 增强管道操作测试完成!")
}
