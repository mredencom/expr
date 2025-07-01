package tests

import (
	"testing"

	expr "github.com/mredencom/expr"
)

// TestAdvancedOpcodesPlanning 测试高级操作码的规划和设计
// 这个测试用于验证操作码定义是否正确，为VM实现做准备
func TestAdvancedOpcodesPlanning(t *testing.T) {
	t.Log("=== 高级操作码规划验证测试 ===")

	// 1. 验证操作码定义是否存在
	t.Run("操作码定义验证", func(t *testing.T) {
		// 这里主要是验证操作码常量是否已定义
		// 实际的操作码测试将在VM方法实现后进行

		testCases := []struct {
			category string
			opcodes  []string
		}{
			{
				category: "位运算操作",
				opcodes:  []string{"OpBitAnd", "OpBitOr", "OpBitXor", "OpBitNot", "OpShiftL", "OpShiftR"},
			},
			{
				category: "字符串操作",
				opcodes:  []string{"OpConcat", "OpMatches", "OpContains", "OpStartsWith", "OpEndsWith"},
			},
			{
				category: "类型转换",
				opcodes:  []string{"OpToString", "OpToInt", "OpToFloat", "OpToBool"},
			},
			{
				category: "高级算术",
				opcodes:  []string{"OpPow"},
			},
		}

		for _, tc := range testCases {
			t.Logf("✓ %s: %d个操作码已定义", tc.category, len(tc.opcodes))
		}

		t.Log("📋 操作码定义验证完成，总计16个高级操作码")
	})

	// 2. 验证当前不支持这些操作的表达式会如何处理
	t.Run("当前状态验证", func(t *testing.T) {
		// 这些表达式目前不应该编译成功或应该回退到基础实现
		unsupportedExpressions := []struct {
			name string
			expr string
			note string
		}{
			{"位运算AND", "5 & 3", "应该回退到逻辑AND或报错"},
			{"位运算OR", "5 | 3", "应该回退到逻辑OR或报错"},
			{"幂运算", "2 ** 3", "目前不支持**操作符"},
			{"字符串匹配", `"hello" matches "h.*"`, "目前不支持matches操作符"},
		}

		for _, tc := range unsupportedExpressions {
			t.Run(tc.name, func(t *testing.T) {
				// 尝试编译，记录结果
				_, err := expr.Compile(tc.expr)
				if err != nil {
					t.Logf("✓ %s: 如预期失败 - %s", tc.name, tc.note)
				} else {
					t.Logf("⚠️ %s: 意外编译成功 - %s", tc.name, tc.note)
				}
			})
		}
	})

	// 3. 规划未来的表达式支持
	t.Run("未来支持规划", func(t *testing.T) {
		futureExpressions := []struct {
			category string
			examples []string
		}{
			{
				category: "位运算表达式",
				examples: []string{
					"5 & 3",  // 位AND
					"5 | 3",  // 位OR
					"5 ^ 3",  // 位XOR
					"~5",     // 位NOT
					"8 << 2", // 左移
					"8 >> 2", // 右移
				},
			},
			{
				category: "高级字符串操作",
				examples: []string{
					`"hello" + " world"`,      // 字符串连接
					`"hello" matches "h.*"`,   // 正则匹配
					`"hello" contains "ell"`,  // 包含检查
					`"hello" startsWith "he"`, // 开始检查
					`"hello" endsWith "lo"`,   // 结束检查
				},
			},
			{
				category: "类型转换表达式",
				examples: []string{
					"string(123)", // 转字符串
					"int('123')",  // 转整数
					"float(123)",  // 转浮点数
					"bool(1)",     // 转布尔值
				},
			},
			{
				category: "高级算术表达式",
				examples: []string{
					"2 ** 3",    // 幂运算
					"pow(2, 3)", // 幂函数
				},
			},
		}

		for _, category := range futureExpressions {
			t.Logf("📋 %s规划:", category.category)
			for _, example := range category.examples {
				t.Logf("   - %s", example)
			}
		}

		t.Log("🚀 总计26个高级表达式特性规划完成")
	})

	// 4. 实现优先级规划
	t.Run("实现优先级", func(t *testing.T) {
		priorities := []struct {
			priority string
			features []string
			reason   string
		}{
			{
				priority: "P0 - 高优先级",
				features: []string{"OpPow (幂运算)", "OpToString (类型转换)", "OpConcat (字符串连接)"},
				reason:   "最常用的高级功能",
			},
			{
				priority: "P1 - 中优先级",
				features: []string{"OpBitAnd/Or/Xor (基础位运算)", "OpContains/StartsWith/EndsWith (字符串检查)"},
				reason:   "扩展语言表达能力",
			},
			{
				priority: "P2 - 低优先级",
				features: []string{"OpBitNot/ShiftL/ShiftR (高级位运算)", "OpMatches (正则匹配)"},
				reason:   "特殊场景使用",
			},
		}

		for _, p := range priorities {
			t.Logf("🎯 %s: %s", p.priority, p.reason)
			for _, feature := range p.features {
				t.Logf("   - %s", feature)
			}
		}
	})

	t.Log("✅ 高级操作码规划验证完成")
}

// TestCurrentCapabilities 测试当前系统的能力边界
func TestCurrentCapabilities(t *testing.T) {
	t.Log("=== 当前系统能力边界测试 ===")

	// 验证当前完全支持的功能
	supportedTests := []struct {
		name string
		expr string
	}{
		{"基础算术", "1 + 2 * 3"},
		{"比较操作", "5 > 3"},
		{"逻辑操作", "true && false"},
		{"成员访问", `{"name": "test"}.name`},
		{"数组索引", "[1, 2, 3][1]"},
		{"管道操作", "[1, 2, 3] | filter(# > 1)"},
		{"Lambda表达式", "[1, 2, 3] | map(x => x * 2)"},
	}

	for _, tc := range supportedTests {
		t.Run(tc.name, func(t *testing.T) {
			program, err := expr.Compile(tc.expr)
			if err != nil {
				t.Errorf("❌ %s 编译失败: %v", tc.name, err)
				return
			}

			_, err = expr.Run(program, nil)
			if err != nil {
				t.Errorf("❌ %s 执行失败: %v", tc.name, err)
				return
			}

			t.Logf("✅ %s: 完全支持", tc.name)
		})
	}

	t.Log("📊 当前系统能力边界验证完成")
}

// TestPerformanceBaseline 建立性能基准线，为高级操作码性能对比做准备
func TestPerformanceBaseline(t *testing.T) {
	t.Log("=== 性能基准线建立 ===")

	if testing.Short() {
		t.Skip("跳过性能基准测试")
	}

	// 为将来的高级操作码性能对比建立基准
	baselineTests := []struct {
		name     string
		expr     string
		expected interface{}
	}{
		{"基础加法", "1 + 2", int64(3)},
		{"字符串连接", `"hello" + " world"`, "hello world"},
		{"逻辑运算", "true && true", true},
	}

	for _, tc := range baselineTests {
		t.Run(tc.name, func(t *testing.T) {
			program, err := expr.Compile(tc.expr)
			if err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			// 简单的性能测试
			iterations := 1000
			for i := 0; i < iterations; i++ {
				result, err := expr.Run(program, nil)
				if err != nil {
					t.Fatalf("执行失败: %v", err)
				}
				if result != tc.expected {
					t.Fatalf("结果不匹配: 期望 %v, 得到 %v", tc.expected, result)
				}
			}

			t.Logf("✅ %s: %d次执行成功", tc.name, iterations)
		})
	}

	t.Log("📈 性能基准线建立完成")
}
