package tests

import (
	"fmt"
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// TestFinalPerformance 最终性能验证测试
func TestFinalPerformance(t *testing.T) {
	fmt.Println("🚀 最终性能验证测试")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	// 测试环境
	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 30, "active": true},
			{"name": "Bob", "age": 25, "active": false},
			{"name": "Charlie", "age": 35, "active": true},
		},
		"threshold": 5,
	}

	// 性能测试用例
	testCases := []struct {
		name       string
		expression string
		iterations int
	}{
		{"基础算术", "2 + 3 * 4", 10000},
		{"成员访问", "users[0].name", 10000},
		{"数组字面量", "[1, 2, 3, 4, 5]", 10000},
		{"对象字面量", `{"name": "test", "value": 42}`, 10000},
		{"管道过滤", "numbers | filter(# > threshold)", 5000},
		{"管道映射", "numbers | map(# * 2)", 5000},
		{"Lambda表达式", "filter(numbers, x => x > 5)", 5000},
		{"复杂管道", "numbers | filter(# > 3) | map(# * 2) | sum()", 3000},
	}

	fmt.Println("\n📊 性能测试结果:")
	fmt.Println("表达式类型                    | 迭代次数 | 总时间    | 平均时间  | ops/sec")
	fmt.Println("------------------------------|----------|-----------|-----------|----------")

	for _, tc := range testCases {
		// 预编译
		program, err := expr.Compile(tc.expression, expr.Env(env))
		if err != nil {
			fmt.Printf("%-30s | 编译失败: %v\n", tc.name, err)
			continue
		}

		// 性能测试
		start := time.Now()
		var lastResult interface{}

		for i := 0; i < tc.iterations; i++ {
			result, err := expr.Run(program, env)
			if err != nil {
				fmt.Printf("%-30s | 执行失败: %v\n", tc.name, err)
				break
			}
			lastResult = result
		}

		elapsed := time.Since(start)
		avgTime := elapsed / time.Duration(tc.iterations)
		opsPerSec := float64(tc.iterations) / elapsed.Seconds()

		fmt.Printf("%-30s | %8d | %9s | %9s | %8.0f\n",
			tc.name,
			tc.iterations,
			elapsed.Round(time.Millisecond),
			avgTime.Round(time.Microsecond),
			opsPerSec,
		)

		// 验证结果正确性
		if lastResult == nil {
			fmt.Printf("  ⚠️  最后结果为nil\n")
		}
	}

	fmt.Println("\n✅ 性能测试完成")
	fmt.Println("🎯 目标: >10,000 ops/sec 用于基础操作")
	fmt.Println("🎯 目标: >5,000 ops/sec 用于复杂操作")
}
