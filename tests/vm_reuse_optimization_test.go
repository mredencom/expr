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

// TestVMReuseOptimization测试VM重用模式的性能效果
func TestVMReuseOptimization(t *testing.T) {
	t.Log("=== VM重用优化测试 ===")
	t.Log("对比VM重用 vs 每次创建的性能差异")
	t.Log("")

	// 准备测试表达式
	expr := "42 + 58 - 10"
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
	iterations := 100000

	t.Run("标准VM每次创建", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < iterations; i++ {
			vmInstance := vm.New(bytecode) // 每次创建新VM
			_, err := vmInstance.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("标准VM每次创建: %d ops/sec", opsPerSec)
	})

	t.Run("标准VM重用", func(t *testing.T) {
		start := time.Now()

		vmInstance := vm.New(bytecode) // 创建一次，重复使用
		for i := 0; i < iterations; i++ {
			vmInstance.ResetStack() // 仅重置栈
			_, err := vmInstance.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("标准VM重用: %d ops/sec", opsPerSec)
	})

	t.Run("优化VM每次创建", func(t *testing.T) {
		start := time.Now()

		factory := vm.DefaultOptimizedFactory()
		for i := 0; i < iterations; i++ {
			vmInstance := factory.CreateVM(bytecode) // 每次创建新VM
			_, err := vmInstance.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
			factory.ReleaseVM(vmInstance) // 释放到池
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("优化VM每次创建: %d ops/sec", opsPerSec)
	})

	t.Run("优化VM重用", func(t *testing.T) {
		start := time.Now()

		factory := vm.DefaultOptimizedFactory()
		vmInstance := factory.CreateVM(bytecode) // 创建一次，重复使用
		defer factory.ReleaseVM(vmInstance)

		for i := 0; i < iterations; i++ {
			vmInstance.ResetStack() // 仅重置栈
			_, err := vmInstance.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
		}

		duration := time.Since(start)
		opsPerSec := int(float64(iterations) / duration.Seconds())
		t.Logf("优化VM重用: %d ops/sec", opsPerSec)
	})
}

// TestOptimalVMUsagePattern测试最优VM使用模式
func TestOptimalVMUsagePattern(t *testing.T) {
	t.Log("=== 最优VM使用模式测试 ===")

	expressions := []string{
		"42 + 58 - 10",
		`"Hello" + " " + "World"`,
		"(10 + 5) * 3 - 2 / 2",
	}

	for i, expr := range expressions {
		t.Run(fmt.Sprintf("表达式_%d", i+1), func(t *testing.T) {
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
			iterations := 50000

			// 测试最优模式：优化VM + 重用
			start := time.Now()

			factory := vm.DefaultOptimizedFactory()
			vmInstance := factory.CreateVM(bytecode)
			defer factory.ReleaseVM(vmInstance)

			for j := 0; j < iterations; j++ {
				vmInstance.ResetStack()
				_, err := vmInstance.Run(bytecode, nil)
				if err != nil {
					t.Fatalf("执行失败: %v", err)
				}
			}

			duration := time.Since(start)
			opsPerSec := int(float64(iterations) / duration.Seconds())

			t.Logf("表达式: %s", expr)
			t.Logf("最优模式性能: %d ops/sec", opsPerSec)

			// 与P0目标对比
			var target int
			if i == 0 { // 基础算术
				target = 50000
			} else if i == 1 { // 字符串操作
				target = 25000
			} else { // 复杂表达式
				target = 35000
			}

			achievement := float64(opsPerSec) / float64(target)
			t.Logf("P0目标: %d ops/sec", target)
			t.Logf("目标达成率: %.1f%%", achievement*100)

			if achievement >= 1.0 {
				t.Logf("🎉 已达到P0目标!")
			} else if achievement >= 0.7 {
				t.Logf("📈 接近P0目标")
			} else {
				t.Logf("⚠️  仍需进一步优化")
			}
			t.Logf("")
		})
	}
}

// TestMemoryPoolEfficiency测试内存池效率
func TestMemoryPoolEfficiency(t *testing.T) {
	t.Log("=== 内存池效率测试 ===")

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

	// 重置统计
	vm.GlobalMemoryOptimizer.ResetStats()

	factory := vm.DefaultOptimizedFactory()

	// 创建和释放多个VM来测试池效率
	numVMs := 1000

	t.Run("内存池循环测试", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < numVMs; i++ {
			vmInstance := factory.CreateVM(bytecode)
			_, err := vmInstance.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}
			factory.ReleaseVM(vmInstance)
		}

		duration := time.Since(start)

		stats := vm.GlobalMemoryOptimizer.GetOptimizationStats()

		t.Logf("总共创建VM: %d", numVMs)
		t.Logf("执行时间: %v", duration)
		t.Logf("平均每个VM: %v", duration/time.Duration(numVMs))
		t.Logf("内存池命中: %d", stats.PoolHits)
		t.Logf("内存池未命中: %d", stats.PoolMisses)
		t.Logf("命中率: %.1f%%", stats.HitRatio*100)

		expectedHits := int64(numVMs * 2) // 每个VM需要栈和全局变量
		if stats.PoolHits >= expectedHits {
			t.Logf("✅ 内存池高效运作!")
		} else {
			t.Logf("⚠️  内存池效率可能需要改进")
		}
	})
}
