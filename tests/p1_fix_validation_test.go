package tests

import (
	"testing"
	"time"

	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

// TestP1FixValidation验证P1优化修复后的性能
func TestP1FixValidation(t *testing.T) {
	t.Log("=== P1 修复验证测试 ===")
	t.Log("测试智能清理策略的性能效果")
	t.Log("")

	testCases := []struct {
		name       string
		expression string
		iterations int
	}{
		{
			name:       "基础算术（修复后）",
			expression: "42 + 58 - 10",
			iterations: 100000,
		},
		{
			name:       "字符串连接（修复后）",
			expression: `"Hello" + " " + "World"`,
			iterations: 50000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 编译表达式
			l := lexer.New(tc.expression)
			p := parser.New(l)
			ast := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("解析失败: %v", p.Errors())
			}

			c := compiler.New()
			if err := c.Compile(ast); err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			bytecode := c.Bytecode()

			// 测试标准VM性能
			standardOPS := benchmarkStandardVMFixed(t, bytecode, tc.iterations)

			// 测试修复后的优化VM性能
			optimizedOPS := benchmarkOptimizedVMFixed(t, bytecode, tc.iterations)

			// 计算性能提升
			improvement := float64(optimizedOPS) / float64(standardOPS)

			// 报告结果
			t.Logf("=== %s ===", tc.name)
			t.Logf("表达式: %s", tc.expression)
			t.Logf("标准VM:   %d ops/sec", standardOPS)
			t.Logf("修复优化VM: %d ops/sec", optimizedOPS)
			t.Logf("性能提升: %.2fx", improvement)

			if improvement >= 1.3 {
				t.Logf("🚀 P1修复效果显著!")
			} else if improvement >= 1.1 {
				t.Logf("📈 P1修复有效果")
			} else if improvement >= 0.95 {
				t.Logf("➡️  性能基本持平")
			} else {
				t.Logf("⚠️  仍有性能问题")
			}
			t.Logf("")
		})
	}
}

// benchmarkStandardVMFixed使用正确的资源管理测试标准VM
func benchmarkStandardVMFixed(t *testing.T, bytecode *vm.Bytecode, iterations int) int {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		vmInstance := vm.New(bytecode)
		_, err := vmInstance.Run(bytecode, nil)
		if err != nil {
			t.Fatalf("标准VM执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}

// benchmarkOptimizedVMFixed使用正确的资源管理测试优化VM
func benchmarkOptimizedVMFixed(t *testing.T, bytecode *vm.Bytecode, iterations int) int {
	start := time.Now()

	// 创建工厂
	factory := vm.DefaultOptimizedFactory()

	for i := 0; i < iterations; i++ {
		vmInstance := factory.CreateVM(bytecode)
		_, err := vmInstance.Run(bytecode, nil)
		if err != nil {
			t.Fatalf("优化VM执行失败: %v", err)
		}
		// 显式释放资源
		factory.ReleaseVM(vmInstance)
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}

// TestSmartClearingValidation验证智能清理策略
func TestSmartClearingValidation(t *testing.T) {
	t.Log("=== 智能清理策略验证 ===")

	// 准备测试数据
	expr := "1 + 2 + 3"
	l := lexer.New(expr)
	p := parser.New(l)
	ast := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("解析失败: %v", p.Errors())
	}

	c := compiler.New()
	if err := c.Compile(ast); err != nil {
		t.Fatalf("编译失败: %v", err)
	}

	bytecode := c.Bytecode()

	t.Run("清理开销对比", func(t *testing.T) {
		// 测试1: 全量清理（原来的方式）- 模拟
		start1 := time.Now()
		for i := 0; i < 10000; i++ {
			// 模拟全量清理开销
			stack := make([]interface{}, 2048)
			globals := make([]interface{}, 65536)
			for j := range stack {
				stack[j] = nil
			}
			for j := range globals {
				globals[j] = nil
			}
		}
		fullClearTime := time.Since(start1)

		// 测试2: 智能清理（新方式）
		start2 := time.Now()
		for i := 0; i < 10000; i++ {
			// 模拟智能清理开销
			stack := make([]interface{}, 2048)
			globals := make([]interface{}, 65536)
			// 仅清理前256个栈位置
			for j := 0; j < 256 && j < len(stack); j++ {
				stack[j] = nil
			}
			// 仅清理前64个全局变量
			for j := 0; j < 64 && j < len(globals); j++ {
				globals[j] = nil
			}
		}
		smartClearTime := time.Since(start2)

		improvement := float64(fullClearTime.Nanoseconds()) / float64(smartClearTime.Nanoseconds())

		t.Logf("全量清理时间: %v", fullClearTime)
		t.Logf("智能清理时间: %v", smartClearTime)
		t.Logf("清理效率提升: %.1fx", improvement)

		if improvement > 10 {
			t.Logf("✅ 智能清理策略显著降低开销!")
		} else {
			t.Logf("⚠️  清理策略改进有限")
		}
	})

	t.Run("内存池统计", func(t *testing.T) {
		// 重置统计
		vm.GlobalMemoryOptimizer.ResetStats()

		factory := vm.DefaultOptimizedFactory()

		// 创建和释放多个VM
		for i := 0; i < 100; i++ {
			vmInstance := factory.CreateVM(bytecode)
			vmInstance.Run(bytecode, nil)
			factory.ReleaseVM(vmInstance)
		}

		stats := vm.GlobalMemoryOptimizer.GetOptimizationStats()
		t.Logf("内存池命中: %d", stats.PoolHits)
		t.Logf("内存池未命中: %d", stats.PoolMisses)
		t.Logf("命中率: %.1f%%", stats.HitRatio*100)

		if stats.PoolHits > 100 {
			t.Logf("✅ 内存池高效运作!")
		}
	})
}

// TestP1ProgressAfterFix测试修复后的P0目标进度
func TestP1ProgressAfterFix(t *testing.T) {
	t.Log("=== P0 目标进度（修复后）===")

	targets := map[string]int{
		"基础算术":  50000, // P0目标
		"字符串操作": 25000, // P0目标
	}

	expressions := map[string]string{
		"基础算术":  "42 + 58 - 10",
		"字符串操作": `"Hello" + " " + "World"`,
	}

	for name, expr := range expressions {
		t.Run(name, func(t *testing.T) {
			l := lexer.New(expr)
			p := parser.New(l)
			ast := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("解析失败: %v", p.Errors())
			}

			c := compiler.New()
			if err := c.Compile(ast); err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			bytecode := c.Bytecode()

			// 测试修复后的优化VM
			optimizedOPS := benchmarkOptimizedVMFixed(t, bytecode, 80000)
			target := targets[name]
			achievement := float64(optimizedOPS) / float64(target)

			t.Logf("表达式: %s", expr)
			t.Logf("修复后性能: %d ops/sec", optimizedOPS)
			t.Logf("P0目标: %d ops/sec", target)
			t.Logf("目标达成率: %.1f%%", achievement*100)

			if achievement >= 1.0 {
				t.Logf("🎉 已达到P0目标!")
			} else if achievement >= 0.7 {
				t.Logf("📈 接近P0目标")
			} else if achievement >= 0.3 {
				t.Logf("⚠️  距离P0目标还有距离")
			} else {
				t.Logf("❌ 距离P0目标很远")
			}
			t.Logf("")
		})
	}
}
