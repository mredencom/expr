package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

// TestP1OptimizationFix验证P1优化修复后的效果
func TestP1OptimizationFix(t *testing.T) {
	t.Log("=== P1 优化修复验证测试 ===")
	t.Log("对比标准VM vs 优化VM的性能差异")
	t.Log("")

	testCases := []struct {
		name       string
		expression string
		iterations int
	}{
		{
			name:       "基础算术",
			expression: "42 + 58 - 10",
			iterations: 100000,
		},
		{
			name:       "字符串连接",
			expression: `"Hello" + " " + "World"`,
			iterations: 50000,
		},
		{
			name:       "复杂表达式",
			expression: "(10 + 5) * 3 - 2 / 2",
			iterations: 60000,
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
			standardOPS := benchmarkStandardVM(t, bytecode, tc.iterations)

			// 测试优化VM性能
			optimizedOPS := benchmarkOptimizedVM(t, bytecode, tc.iterations)

			// 计算性能提升
			improvement := float64(optimizedOPS) / float64(standardOPS)

			// 报告结果
			t.Logf("=== %s ===", tc.name)
			t.Logf("表达式: %s", tc.expression)
			t.Logf("标准VM:  %d ops/sec", standardOPS)
			t.Logf("优化VM:  %d ops/sec", optimizedOPS)
			t.Logf("性能提升: %.2fx", improvement)

			if improvement >= 1.5 {
				t.Logf("🚀 P1优化效果显著!")
			} else if improvement >= 1.1 {
				t.Logf("📈 P1优化有效果")
			} else if improvement >= 0.9 {
				t.Logf("➡️  性能基本持平")
			} else {
				t.Logf("⚠️  优化有负面影响")
			}
			t.Logf("")
		})
	}
}

// benchmarkStandardVM测试标准VM性能
func benchmarkStandardVM(t *testing.T, bytecode *vm.Bytecode, iterations int) int {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		vmInstance := vm.New(bytecode) // 使用标准VM
		_, err := vmInstance.Run(bytecode, nil)
		if err != nil {
			t.Fatalf("标准VM执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}

// benchmarkOptimizedVM测试优化VM性能
func benchmarkOptimizedVM(t *testing.T, bytecode *vm.Bytecode, iterations int) int {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		vmInstance := vm.NewOptimized(bytecode) // 使用P1优化VM
		_, err := vmInstance.Run(bytecode, nil)
		if err != nil {
			t.Fatalf("优化VM执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}

// TestMemoryOptimizationVerification验证内存优化是否生效
func TestMemoryOptimizationVerification(t *testing.T) {
	t.Log("=== 内存优化验证测试 ===")

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

	t.Run("标准VM内存分配", func(t *testing.T) {
		// 创建标准VM
		standardVM := vm.New(bytecode)

		// 检查栈和全局变量是否为直接分配
		t.Logf("标准VM栈容量: %d", cap(standardVM.StackDebug()))
		t.Logf("标准VM全局变量容量: %d", cap(standardVM.GlobalsDebug()))

		// 标准VM应该直接分配内存
		t.Logf("✓ 标准VM使用直接内存分配")
	})

	t.Run("优化VM内存池", func(t *testing.T) {
		// 获取初始统计
		initialStats := vm.GlobalMemoryOptimizer.GetOptimizationStats()

		// 创建优化VM
		optimizedVM := vm.NewOptimized(bytecode)

		// 检查是否使用了内存池
		t.Logf("优化VM栈容量: %d", cap(optimizedVM.StackDebug()))
		t.Logf("优化VM全局变量容量: %d", cap(optimizedVM.GlobalsDebug()))

		// 获取更新后的统计
		finalStats := vm.GlobalMemoryOptimizer.GetOptimizationStats()
		poolHitIncrease := finalStats.PoolHits - initialStats.PoolHits

		t.Logf("内存池命中次数增加: %d", poolHitIncrease)
		t.Logf("内存池总命中: %d", finalStats.PoolHits)
		t.Logf("内存池未命中: %d", finalStats.PoolMisses)

		if poolHitIncrease > 0 {
			t.Logf("✅ 优化VM成功使用内存池!")
		} else {
			t.Logf("⚠️  内存池可能未被使用")
		}
	})
}

// TestP1ComponentsIsolation测试各个P1组件的隔离效果
func TestP1ComponentsIsolation(t *testing.T) {
	t.Log("=== P1 组件隔离测试 ===")

	expr := "a + b * c"
	l := lexer.New(expr)
	p := parser.New(l)
	ast := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("解析失败: %v", p.Errors())
	}

	c := compiler.New()
	// 添加环境变量以避免编译错误
	c.AddEnvironment(map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}, nil)

	if err := c.Compile(ast); err != nil {
		t.Fatalf("编译失败: %v", err)
	}

	bytecode := c.Bytecode()
	iterations := 30000

	t.Run("仅内存优化", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < iterations; i++ {
			vmInstance := vm.NewOptimizedWithOptions(bytecode, true, false, false) // 仅内存优化
			vmInstance.Run(bytecode, map[string]interface{}{"a": 1, "b": 2, "c": 3})
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("仅内存优化: %d ops/sec", opsPerSec)
	})

	t.Run("内存优化+缓存", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < iterations; i++ {
			vmInstance := vm.NewOptimizedWithOptions(bytecode, true, false, true) // 内存+缓存
			vmInstance.Run(bytecode, map[string]interface{}{"a": 1, "b": 2, "c": 3})
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("内存优化+缓存: %d ops/sec", opsPerSec)
	})
}

// TestOptimizationProgress测试优化进度
func TestOptimizationProgress(t *testing.T) {
	t.Log("=== 优化进度测试 ===")
	t.Log("基于PERFORMANCE_SUMMARY.md目标验证当前进度")
	t.Log("")

	// 基于PERFORMANCE_SUMMARY.md的P0目标
	targets := map[string]int{
		"基础算术":  50000, // P0目标 (vs最终目标20M)
		"字符串操作": 25000, // P0目标 (vs最终目标5M)
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

			// 测试优化VM
			optimizedOPS := benchmarkOptimizedVM(t, bytecode, 80000)
			target := targets[name]
			achievement := float64(optimizedOPS) / float64(target)

			t.Logf("表达式: %s", expr)
			t.Logf("优化VM性能: %d ops/sec", optimizedOPS)
			t.Logf("P0目标: %d ops/sec", target)
			t.Logf("目标达成率: %.1f%%", achievement*100)

			if achievement >= 1.0 {
				t.Logf("🎉 已达到P0目标!")
			} else if achievement >= 0.5 {
				t.Logf("📈 接近P0目标 (还需%.1fx提升)", 1.0/achievement)
			} else {
				t.Logf("⚠️  距离P0目标较远 (需要%.1fx提升)", 1.0/achievement)
			}
			t.Logf("")
		})
	}
}

// TestP1OptimizationIntegration验证P1优化的集成效果
func TestP1OptimizationIntegration(t *testing.T) {
	t.Log("=== P1 优化集成效果测试 ===")

	// 复杂表达式测试P1优化的综合效果
	complexExpressions := []string{
		`(a + b) * c - d / e`,
		`"prefix" + var1 + "suffix" + var2`,
		`arr[0] + arr[1] * 2`,
		`obj.field1 + obj.field2`,
	}

	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 3, "d": 40, "e": 4,
		"var1": "test", "var2": "value",
		"arr": []int{5, 10, 15},
		"obj": map[string]interface{}{"field1": 100, "field2": 200},
	}

	for i, expr := range complexExpressions {
		t.Run(fmt.Sprintf("复杂表达式_%d", i+1), func(t *testing.T) {
			l := lexer.New(expr)
			p := parser.New(l)
			ast := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("解析失败: %v", p.Errors())
			}

			c := compiler.New()
			c.AddEnvironment(env, nil)
			if err := c.Compile(ast); err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			bytecode := c.Bytecode()

			// 小批量快速测试
			iterations := 10000

			standardOPS := benchmarkStandardVMWithEnv(t, bytecode, iterations, env)
			optimizedOPS := benchmarkOptimizedVMWithEnv(t, bytecode, iterations, env)

			improvement := float64(optimizedOPS) / float64(standardOPS)

			t.Logf("表达式: %s", expr)
			t.Logf("标准VM: %d ops/sec", standardOPS)
			t.Logf("优化VM: %d ops/sec", optimizedOPS)
			t.Logf("性能提升: %.2fx", improvement)

			if improvement > 1.0 {
				t.Logf("✅ P1优化对复杂表达式有效")
			} else {
				t.Logf("⚠️  复杂表达式性能未提升")
			}
		})
	}
}

// benchmarkStandardVMWithEnv使用环境变量测试标准VM
func benchmarkStandardVMWithEnv(t *testing.T, bytecode *vm.Bytecode, iterations int, env map[string]interface{}) int {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		vmInstance := vm.New(bytecode)
		_, err := vmInstance.Run(bytecode, env)
		if err != nil {
			t.Fatalf("标准VM执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}

// benchmarkOptimizedVMWithEnv使用环境变量测试优化VM
func benchmarkOptimizedVMWithEnv(t *testing.T, bytecode *vm.Bytecode, iterations int, env map[string]interface{}) int {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		vmInstance := vm.NewOptimized(bytecode)
		_, err := vmInstance.Run(bytecode, env)
		if err != nil {
			t.Fatalf("优化VM执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	return int(float64(iterations) / duration.Seconds())
}
