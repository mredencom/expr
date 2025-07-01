package tests

import (
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// BenchmarkCompilationOnly measures compilation performance
func BenchmarkCompilationOnly(b *testing.B) {
	expression := "1 + 2"
	env := map[string]interface{}{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.Compile(expression, expr.Env(env))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExecutionOnly measures execution performance with pre-compiled program
func BenchmarkExecutionOnly(b *testing.B) {
	expression := "1 + 2"
	env := map[string]interface{}{}

	// Pre-compile once
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

// TestSeparatePerformance tests compilation vs execution performance separately
func TestSeparatePerformance(t *testing.T) {
	expressions := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Simple Int", "1 + 2", nil},
		{"Simple Float", "1.5 + 2.5", nil},
		{"Simple String", `"a" + "b"`, nil},
		{"Variables", "a + b", map[string]interface{}{"a": 10, "b": 20}},
		{"Complex", "((a + b) * c - d) / e", map[string]interface{}{"a": 10, "b": 20, "c": 5, "d": 15, "e": 3}},
	}

	for _, test := range expressions {
		t.Run(test.name, func(t *testing.T) {
			// Measure compilation time
			var compilationTime time.Duration
			var program *expr.Program
			var err error

			start := time.Now()
			for i := 0; i < 1000; i++ {
				program, err = expr.Compile(test.expr, expr.Env(test.env))
				if err != nil {
					t.Fatal(err)
				}
			}
			compilationTime = time.Since(start) / 1000

			// Measure execution time (with pre-compiled program)
			iterations := 100000
			start = time.Now()
			for i := 0; i < iterations; i++ {
				_, err := expr.Run(program, test.env)
				if err != nil {
					t.Fatal(err)
				}
			}
			executionTime := time.Since(start)
			avgExecutionTime := executionTime / time.Duration(iterations)

			execOpsPerSec := float64(iterations) / executionTime.Seconds()

			t.Logf("%s:", test.name)
			t.Logf("  Compilation: %.2f μs/op", float64(compilationTime.Nanoseconds())/1000.0)
			t.Logf("  Execution: %.2f μs/op (%.0f ops/sec)",
				float64(avgExecutionTime.Nanoseconds())/1000.0, execOpsPerSec)

			// Check if execution meets performance targets
			if execOpsPerSec > 1000000 {
				t.Logf("  🎉 EXCELLENT: >1M ops/sec execution performance!")
			} else if execOpsPerSec > 100000 {
				t.Logf("  ✅ GOOD: >100K ops/sec execution performance")
			} else if execOpsPerSec > 10000 {
				t.Logf("  ⚠️  OK: >10K ops/sec execution performance")
			} else {
				t.Logf("  ❌ POOR: <10K ops/sec execution performance")
			}
		})
	}
}
