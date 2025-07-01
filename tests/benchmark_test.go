package tests

import (
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

// BenchmarkIntegerArithmetic tests the performance of integer arithmetic operations
func BenchmarkIntegerArithmetic(b *testing.B) {
	expression := "10 + 20 * 5"
	env := map[string]interface{}{}

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

// BenchmarkFloatArithmetic tests the performance of float arithmetic operations
func BenchmarkFloatArithmetic(b *testing.B) {
	expression := "10.5 + 20.3 * 5.7"
	env := map[string]interface{}{}

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

// BenchmarkStringConcatenation tests the performance of string concatenation
func BenchmarkStringConcatenation(b *testing.B) {
	expression := `"Hello" + " " + "World"`
	env := map[string]interface{}{}

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

// BenchmarkVariableAccess tests the performance of variable operations
func BenchmarkVariableAccess(b *testing.B) {
	expression := "a + b * c"
	env := map[string]interface{}{
		"a": 10,
		"b": 20,
		"c": 5,
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

// BenchmarkComplexExpression tests the performance of complex expressions
func BenchmarkComplexExpression(b *testing.B) {
	expression := "((a + b) * c - d) / e + f % g"
	env := map[string]interface{}{
		"a": 10,
		"b": 20,
		"c": 5,
		"d": 15,
		"e": 3,
		"f": 100,
		"g": 7,
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

// Performance comparison test
func TestPerformanceComparison(t *testing.T) {
	expressions := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Simple Int", "1 + 2", nil},
		{"Simple Float", "1.5 + 2.5", nil},
		{"Simple String", `"a" + "b"`, nil},
		{"Variables", "a + b", map[string]interface{}{"a": 10, "b": 20}},
	}

	for _, test := range expressions {
		t.Run(test.name, func(t *testing.T) {
			program, err := expr.Compile(test.expr, expr.Env(test.env))
			if err != nil {
				t.Fatal(err)
			}

			// Warmup
			for i := 0; i < 1000; i++ {
				_, err := expr.Run(program, test.env)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Measure performance
			iterations := 100000
			start := time.Now()
			for i := 0; i < iterations; i++ {
				_, err := expr.Run(program, test.env)
				if err != nil {
					t.Fatal(err)
				}
			}
			elapsed := time.Since(start)

			opsPerSec := float64(iterations) / elapsed.Seconds()
			t.Logf("%s: %.0f ops/sec (%.2f μs/op)",
				test.name, opsPerSec, float64(elapsed.Nanoseconds())/float64(iterations)/1000.0)

			// Check if we're meeting performance targets
			if opsPerSec < 10000 {
				t.Logf("WARNING: Performance below 10K ops/sec target: %.0f", opsPerSec)
			}
		})
	}
}
