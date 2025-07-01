package tests

import (
	"runtime"
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// TestVMPoolingPerformance tests if VM pooling improves performance
func TestVMPoolingPerformance(t *testing.T) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
	}

	// Compile once
	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		t.Fatal(err)
	}

	// Test execution performance
	iterations := 10000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			t.Fatal(err)
		}
	}

	elapsed := time.Since(start)
	opsPerSec := float64(iterations) / elapsed.Seconds()
	avgTimePerOp := float64(elapsed.Nanoseconds()) / float64(iterations) / 1000.0 // μs

	t.Logf("VM池化性能测试:")
	t.Logf("  表达式: %s", expression)
	t.Logf("  性能: %.0f ops/sec (%.2f μs/op)", opsPerSec, avgTimePerOp)

	// Check if performance improved (should be better than previous ~300μs/op)
	if avgTimePerOp < 200 {
		t.Logf("  ✅ 性能改善: 平均执行时间 < 200μs/op")
	} else {
		t.Logf("  ⚠️  性能仍需改善: 平均执行时间 %.2fμs/op", avgTimePerOp)
	}
}

// TestConstantFoldingImprovement tests if constant folding is working
func TestConstantFoldingImprovement(t *testing.T) {
	testCases := []struct {
		name     string
		expr     string
		expected interface{}
	}{
		{"Simple Addition", "1 + 2", int64(3)},
		{"Complex Arithmetic", "10 + 20 * 5", int64(110)},
		{"String Concatenation", `"Hello" + " " + "World"`, "Hello World"},
		{"Comparison", "5 > 3", true},
		{"Logical", "true && false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Compile
			start := time.Now()
			program, err := expr.Compile(tc.expr)
			if err != nil {
				t.Fatal(err)
			}
			compileTime := time.Since(start)

			// Execute
			start = time.Now()
			result, err := expr.Run(program, nil)
			if err != nil {
				t.Fatal(err)
			}
			execTime := time.Since(start)

			t.Logf("常量折叠测试 - %s:", tc.name)
			t.Logf("  表达式: %s", tc.expr)
			t.Logf("  编译时间: %v", compileTime)
			t.Logf("  执行时间: %v", execTime)
			t.Logf("  结果: %v (期望: %v)", result, tc.expected)

			// Check if result is correct
			if result != tc.expected {
				t.Errorf("结果不匹配: 得到 %v, 期望 %v", result, tc.expected)
			}

			// Check if execution is fast (should be very fast for constants)
			if execTime < time.Microsecond*50 {
				t.Logf("  ✅ 执行速度优秀: %v", execTime)
			} else if execTime < time.Microsecond*200 {
				t.Logf("  ⚠️  执行速度一般: %v", execTime)
			} else {
				t.Logf("  ❌ 执行速度较慢: %v", execTime)
			}
		})
	}
}

// TestMemoryUsageImprovement tests if memory usage has improved
func TestMemoryUsageImprovement(t *testing.T) {
	expression := "a + b"
	env := map[string]interface{}{
		"a": 10, "b": 20,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		t.Fatal(err)
	}

	// Record initial memory
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Execute many times
	iterations := 10000
	for i := 0; i < iterations; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Force GC and measure memory
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memoryIncrease := m2.Alloc - m1.Alloc
	t.Logf("内存使用测试:")
	t.Logf("  执行次数: %d", iterations)
	t.Logf("  内存增长: %d bytes", memoryIncrease)
	t.Logf("  平均每次: %.2f bytes", float64(memoryIncrease)/float64(iterations))

	// Check if memory usage is reasonable
	avgMemoryPerOp := float64(memoryIncrease) / float64(iterations)
	if avgMemoryPerOp < 100 {
		t.Logf("  ✅ 内存使用优秀: %.2f bytes/op", avgMemoryPerOp)
	} else if avgMemoryPerOp < 1000 {
		t.Logf("  ⚠️  内存使用一般: %.2f bytes/op", avgMemoryPerOp)
	} else {
		t.Logf("  ❌ 内存使用过高: %.2f bytes/op", avgMemoryPerOp)
	}
}

// BenchmarkImprovedConstantFolding benchmarks constant folding performance
func BenchmarkImprovedConstantFolding(b *testing.B) {
	testCases := []struct {
		name string
		expr string
	}{
		{"Simple Constant", "42"},
		{"Simple Arithmetic", "1 + 2"},
		{"Complex Arithmetic", "10 + 20 * 5 - 15 / 3"},
		{"String Concatenation", `"Hello" + " " + "World"`},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			program, err := expr.Compile(tc.expr)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := expr.Run(program, nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkImprovedVMPooling benchmarks VM pooling performance
func BenchmarkImprovedVMPooling(b *testing.B) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentExecution tests concurrent performance with VM pooling
func BenchmarkConcurrentExecution(b *testing.B) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := expr.Run(program, env)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
