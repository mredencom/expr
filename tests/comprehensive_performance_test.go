package tests

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// 性能基准测试 - 常量折叠优化
func BenchmarkConstantFolding(b *testing.B) {
	testCases := []struct {
		name string
		expr string
	}{
		{"Simple Addition", "1 + 2"},
		{"Simple Multiplication", "3 * 4"},
		{"Complex Arithmetic", "10 + 20 * 5 - 15 / 3"},
		{"Comparison", "5 > 3 && 10 < 20"},
		{"String Concatenation", `"Hello" + " " + "World"`},
		{"Mixed Types", "10 + 5.5 * 2"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			program, err := expr.Compile(tc.expr, expr.Env(nil))
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

// 性能基准测试 - 快速路径优化
func BenchmarkFastPathOptimizations(b *testing.B) {
	testCases := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Int64 Addition", "a + b", map[string]interface{}{"a": int64(10), "b": int64(20)}},
		{"Float64 Addition", "a + b", map[string]interface{}{"a": 10.5, "b": 20.3}},
		{"String Addition", "a + b", map[string]interface{}{"a": "Hello", "b": "World"}},
		{"Int64 Comparison", "a > b", map[string]interface{}{"a": int64(15), "b": int64(10)}},
		{"Float64 Comparison", "a < b", map[string]interface{}{"a": 5.5, "b": 10.2}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			program, err := expr.Compile(tc.expr, expr.Env(tc.env))
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := expr.Run(program, tc.env)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// 性能基准测试 - 对象池优化
func BenchmarkObjectPooling(b *testing.B) {
	// 测试大量小对象的创建和销毁
	expressions := []string{
		"1", "2", "3", "4", "5",
		"1.1", "2.2", "3.3", "4.4", "5.5",
		`"a"`, `"b"`, `"c"`, `"d"`, `"e"`,
	}

	b.Run("Sequential Operations", func(b *testing.B) {
		programs := make([]*expr.Program, len(expressions))
		for i, exprStr := range expressions {
			program, err := expr.Compile(exprStr, expr.Env(nil))
			if err != nil {
				b.Fatal(err)
			}
			programs[i] = program
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, program := range programs {
				_, err := expr.Run(program, nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// 性能基准测试 - 字节码缓存
func BenchmarkBytecodeCaching(b *testing.B) {
	expressions := []string{
		"a + b * c",
		"x > y && z < w",
		`"Hello" + name`,
		"price * quantity * (1 - discount)",
	}

	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
		"x": 15, "y": 10, "z": 5, "w": 20,
		"name":  "World",
		"price": 100.0, "quantity": 5, "discount": 0.1,
	}

	b.Run("Repeated Compilation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, exprStr := range expressions {
				program, err := expr.Compile(exprStr, expr.Env(env))
				if err != nil {
					b.Fatal(err)
				}
				_, err = expr.Run(program, env)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("Cached Compilation", func(b *testing.B) {
		programs := make([]*expr.Program, len(expressions))
		for i, exprStr := range expressions {
			program, err := expr.Compile(exprStr, expr.Env(env))
			if err != nil {
				b.Fatal(err)
			}
			programs[i] = program
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, program := range programs {
				_, err := expr.Run(program, env)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// 性能基准测试 - 栈优化
func BenchmarkStackOptimizations(b *testing.B) {
	// 测试深度嵌套表达式
	deepExpression := "1"
	for i := 0; i < 100; i++ {
		deepExpression = "(" + deepExpression + " + 1)"
	}

	b.Run("Deep Nesting", func(b *testing.B) {
		program, err := expr.Compile(deepExpression, expr.Env(nil))
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

	// 测试大量变量访问
	largeEnv := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		largeEnv[fmt.Sprintf("var%d", i)] = i
	}

	largeExpression := "var0"
	for i := 1; i < 50; i++ {
		largeExpression += fmt.Sprintf(" + var%d", i)
	}

	b.Run("Many Variables", func(b *testing.B) {
		program, err := expr.Compile(largeExpression, expr.Env(largeEnv))
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := expr.Run(program, largeEnv)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 性能基准测试 - 内存管理优化
func BenchmarkMemoryManagement(b *testing.B) {
	// 测试大量字符串操作
	b.Run("String Operations", func(b *testing.B) {
		expression := `"a" + "b" + "c" + "d" + "e"`
		program, err := expr.Compile(expression, expr.Env(nil))
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

	// 测试数组操作
	b.Run("Array Operations", func(b *testing.B) {
		expression := "[1, 2, 3, 4, 5]"
		program, err := expr.Compile(expression, expr.Env(nil))
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

// 性能基准测试 - 并发安全
func BenchmarkConcurrency(b *testing.B) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Concurrent Execution", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := expr.Run(program, env)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// 综合性能测试
func TestComprehensivePerformance(t *testing.T) {
	testCases := []struct {
		name        string
		expression  string
		env         map[string]interface{}
		targetOps   float64 // 目标每秒操作数
		description string
	}{
		{
			name:        "Simple Constant",
			expression:  "42",
			env:         nil,
			targetOps:   100000,
			description: "常量折叠优化测试",
		},
		{
			name:        "Constant Arithmetic",
			expression:  "10 + 20 * 5",
			env:         nil,
			targetOps:   50000,
			description: "常量算术运算优化测试",
		},
		{
			name:        "Variable Access",
			expression:  "a + b",
			env:         map[string]interface{}{"a": 10, "b": 20},
			targetOps:   30000,
			description: "变量访问优化测试",
		},
		{
			name:        "Complex Expression",
			expression:  "(a + b) * c - d / e",
			env:         map[string]interface{}{"a": 10, "b": 20, "c": 5, "d": 100, "e": 2},
			targetOps:   20000,
			description: "复杂表达式优化测试",
		},
		{
			name:        "String Operations",
			expression:  `"Hello" + " " + "World"`,
			env:         nil,
			targetOps:   25000,
			description: "字符串操作优化测试",
		},
		{
			name:        "Comparison Operations",
			expression:  "a > b && c < d",
			env:         map[string]interface{}{"a": 15, "b": 10, "c": 5, "d": 20},
			targetOps:   25000,
			description: "比较操作优化测试",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			program, err := expr.Compile(tc.expression, expr.Env(tc.env))
			if err != nil {
				t.Fatalf("编译失败: %v", err)
			}

			// 预热
			for i := 0; i < 1000; i++ {
				_, err := expr.Run(program, tc.env)
				if err != nil {
					t.Fatalf("预热失败: %v", err)
				}
			}

			// 性能测试
			iterations := 100000
			start := time.Now()
			for i := 0; i < iterations; i++ {
				_, err := expr.Run(program, tc.env)
				if err != nil {
					t.Fatalf("执行失败: %v", err)
				}
			}
			elapsed := time.Since(start)

			opsPerSec := float64(iterations) / elapsed.Seconds()
			avgTimePerOp := float64(elapsed.Nanoseconds()) / float64(iterations) / 1000.0 // μs

			t.Logf("%s:", tc.description)
			t.Logf("  表达式: %s", tc.expression)
			t.Logf("  性能: %.0f ops/sec (%.2f μs/op)", opsPerSec, avgTimePerOp)
			t.Logf("  目标: %.0f ops/sec", tc.targetOps)

			// 性能检查
			if opsPerSec < tc.targetOps {
				t.Logf("  ⚠️  性能低于目标: %.0f < %.0f ops/sec", opsPerSec, tc.targetOps)
			} else {
				t.Logf("  ✅ 性能达到目标: %.0f >= %.0f ops/sec", opsPerSec, tc.targetOps)
			}

			// 内存使用分析
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			t.Logf("  内存使用: %d KB", m.Alloc/1024)
		})
	}
}

// 压力测试
func TestStressTest(t *testing.T) {
	// 生成大量随机表达式进行压力测试
	rand.Seed(time.Now().UnixNano())

	expressions := make([]string, 100)
	envs := make([]map[string]interface{}, 100)

	for i := 0; i < 100; i++ {
		// 生成随机表达式
		expr := generateRandomExpression(rand.Intn(10) + 1)
		expressions[i] = expr

		// 生成随机环境
		env := make(map[string]interface{})
		for j := 0; j < 5; j++ {
			key := fmt.Sprintf("var%d", j)
			switch rand.Intn(3) {
			case 0:
				env[key] = rand.Int63()
			case 1:
				env[key] = rand.Float64() * 1000
			case 2:
				env[key] = fmt.Sprintf("str%d", rand.Intn(1000))
			}
		}
		envs[i] = env
	}

	// 编译所有表达式
	programs := make([]*expr.Program, len(expressions))
	start := time.Now()
	for i, exprStr := range expressions {
		program, err := expr.Compile(exprStr, expr.Env(envs[i]))
		if err != nil {
			t.Logf("表达式编译失败: %s, 错误: %v", exprStr, err)
			continue
		}
		programs[i] = program
	}
	compileTime := time.Since(start)

	// 执行所有表达式
	execStart := time.Now()
	successCount := 0
	for i, program := range programs {
		if program == nil {
			continue
		}
		_, err := expr.Run(program, envs[i])
		if err == nil {
			successCount++
		}
	}
	execTime := time.Since(execStart)

	t.Logf("压力测试结果:")
	t.Logf("  编译时间: %v", compileTime)
	t.Logf("  执行时间: %v", execTime)
	t.Logf("  成功执行: %d/%d", successCount, len(expressions))
	t.Logf("  平均编译时间: %v", compileTime/time.Duration(len(expressions)))
	t.Logf("  平均执行时间: %v", execTime/time.Duration(successCount))
}

// 生成随机表达式
func generateRandomExpression(depth int) string {
	if depth <= 0 {
		// 生成叶子节点
		switch rand.Intn(4) {
		case 0:
			return fmt.Sprintf("%d", rand.Intn(100))
		case 1:
			return fmt.Sprintf("%.2f", rand.Float64()*100)
		case 2:
			return fmt.Sprintf(`"str%d"`, rand.Intn(100))
		case 3:
			return fmt.Sprintf("var%d", rand.Intn(5))
		}
	}

	// 生成操作符
	operators := []string{"+", "-", "*", "/", ">", "<", "==", "!=", "&&", "||"}
	op := operators[rand.Intn(len(operators))]

	// 递归生成子表达式
	left := generateRandomExpression(depth - 1)
	right := generateRandomExpression(depth - 1)

	return fmt.Sprintf("(%s %s %s)", left, op, right)
}

// 内存泄漏测试
func TestMemoryLeak(t *testing.T) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10, "b": 20, "c": 5,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		t.Fatal(err)
	}

	// 记录初始内存使用
	var initialMem runtime.MemStats
	runtime.ReadMemStats(&initialMem)

	// 执行大量操作
	for i := 0; i < 100000; i++ {
		_, err := expr.Run(program, env)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 强制垃圾回收
	runtime.GC()

	// 记录最终内存使用
	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)

	memoryIncrease := finalMem.Alloc - initialMem.Alloc
	t.Logf("内存使用变化: %d bytes", memoryIncrease)

	// 检查内存泄漏（允许一定的内存增长）
	if memoryIncrease > 1024*1024 { // 1MB
		t.Errorf("可能存在内存泄漏，内存增长: %d bytes", memoryIncrease)
	}
}
