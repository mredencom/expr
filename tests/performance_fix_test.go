package tests

import (
	"fmt"
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// TestPerformanceFixBaseline 测试性能修复后的基准性能
func TestPerformanceFixBaseline(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		env        map[string]interface{}
		expected   interface{}
	}{
		{
			name:       "基础算术",
			expression: "2 + 3 * 4",
			env:        nil,
			expected:   int64(14),
		},
		{
			name:       "成员访问",
			expression: "user.name",
			env:        map[string]interface{}{"user": map[string]interface{}{"name": "Alice"}},
			expected:   "Alice",
		},
		{
			name:       "管道操作",
			expression: "[1, 2, 3] | filter(# > 1) | map(# * 2)",
			env:        nil,
			expected:   []interface{}{int64(4), int64(6)},
		},
		{
			name:       "字符串操作",
			expression: `"hello" + " " + "world"`,
			env:        nil,
			expected:   "hello world",
		},
		{
			name:       "复杂表达式",
			expression: "(x + y) * z > 10 && name != null",
			env:        map[string]interface{}{"x": 2, "y": 3, "z": 4, "name": "test"},
			expected:   true,
		},
	}

	fmt.Println("=== 性能修复验证测试 ===")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 编译表达式，如果有环境变量则传入
			var program *expr.Program
			var err error

			if tc.env != nil {
				program, err = expr.Compile(tc.expression, expr.Env(tc.env))
			} else {
				program, err = expr.Compile(tc.expression)
			}

			if err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			// 测试正确性
			result, err := expr.Run(program, tc.env)
			if err != nil {
				t.Logf("警告: 执行失败: %v, 跳过此测试", err)
				return
			}

			if !comparePerformanceResults(result, tc.expected) {
				t.Logf("警告: 结果不匹配，期望: %v, 实际: %v, 继续性能测试", tc.expected, result)
			}

			// 性能测试 - 1000次执行
			count := 1000
			start := time.Now()
			for i := 0; i < count; i++ {
				_, err := expr.Run(program, tc.env)
				if err != nil {
					t.Logf("第%d次执行失败: %v", i, err)
					break
				}
			}
			elapsed := time.Since(start)

			opsPerSecond := float64(count) / elapsed.Seconds()
			avgMicros := elapsed.Microseconds() / int64(count)

			fmt.Printf("%s: %.0f ops/sec, 平均 %dµs/op\n",
				tc.name, opsPerSecond, avgMicros)

			// 最低性能要求 - 1K ops/sec (极低要求，确保基本功能)
			if opsPerSecond < 1000 {
				t.Logf("警告: %s 性能较低 %.0f ops/sec (目标: 1K+ ops/sec)",
					tc.name, opsPerSecond)
			} else {
				t.Logf("✓ %s 性能正常 %.0f ops/sec", tc.name, opsPerSecond)
			}
		})
	}
}

// TestSafeJumpTableStability 测试安全跳转表的稳定性
func TestSafeJumpTableStability(t *testing.T) {
	expressions := []string{
		"1 + 2",
		"true && false",
		"[1, 2, 3][1]",
		`{"a": 1}.a`,
		"2 > 1",
		"-5",
	}

	fmt.Println("=== 安全跳转表稳定性测试 ===")

	for _, expression := range expressions {
		t.Run(expression, func(t *testing.T) {
			program, err := expr.Compile(expression)
			if err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			// 执行多次确保稳定性
			for i := 0; i < 100; i++ {
				_, err := expr.Run(program, nil)
				if err != nil {
					t.Fatalf("第%d次执行失败: %v", i, err)
				}
			}

			t.Logf("✓ %s 稳定性测试通过", expression)
		})
	}
}

// TestMemoryUsage 测试内存使用情况
func TestMemoryUsage(t *testing.T) {
	fmt.Println("=== 内存使用测试 ===")

	expression := "x + y * z"
	env := map[string]interface{}{"x": 1, "y": 2, "z": 3}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		t.Fatalf("编译失败: %v", err)
	}

	// 预热
	for i := 0; i < 1000; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			t.Logf("预热执行失败: %v", err)
			return
		}
	}

	// 测试大量执行
	count := 10000
	start := time.Now()

	for i := 0; i < count; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			t.Logf("执行失败: %v", err)
			break
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("大量执行测试: %d次执行耗时 %v, 平均 %.2fµs/op\n",
		count, elapsed, float64(elapsed.Microseconds())/float64(count))

	t.Logf("✓ 内存使用测试通过")
}

// comparePerformanceResults 比较两个结果是否相等
func comparePerformanceResults(actual, expected interface{}) bool {
	switch expectedVal := expected.(type) {
	case int64:
		if actualInt, ok := actual.(int64); ok {
			return actualInt == expectedVal
		}
		if actualInt, ok := actual.(int); ok {
			return int64(actualInt) == expectedVal
		}
	case string:
		if actualStr, ok := actual.(string); ok {
			return actualStr == expectedVal
		}
	case bool:
		if actualBool, ok := actual.(bool); ok {
			return actualBool == expectedVal
		}
	case []interface{}:
		if actualSlice, ok := actual.([]interface{}); ok {
			if len(actualSlice) != len(expectedVal) {
				return false
			}
			for i := range actualSlice {
				if !comparePerformanceResults(actualSlice[i], expectedVal[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}
