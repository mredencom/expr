package main

import (
	"fmt"
	"time"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== Week 10: Comprehensive Test Suite ===")

	// Run all test categories
	fmt.Println("\n🧪 Running Comprehensive Tests...")

	results := TestResults{}

	// Basic functionality tests
	fmt.Println("\n1. Basic Expression Tests:")
	results.Add(runBasicTests())

	// Type system tests
	fmt.Println("\n2. Type System Tests:")
	results.Add(runTypeTests())

	// Built-in function tests
	fmt.Println("\n3. Built-in Function Tests:")
	results.Add(runBuiltinTests())

	// Error handling tests
	fmt.Println("\n4. Error Handling Tests:")
	results.Add(runErrorTests())

	// Performance tests
	fmt.Println("\n5. Performance Tests:")
	results.Add(runPerformanceTests())

	// Print summary
	fmt.Println("\n📊 Test Summary:")
	results.PrintSummary()

	// Run benchmarks
	fmt.Println("\n⚡ Performance Benchmarks:")
	runBenchmarks()
}

type TestResults struct {
	Total  int
	Passed int
	Failed int
}

func (r *TestResults) Add(other TestResults) {
	r.Total += other.Total
	r.Passed += other.Passed
	r.Failed += other.Failed
}

func (r *TestResults) PrintSummary() {
	fmt.Printf("  Total Tests: %d\n", r.Total)
	fmt.Printf("  Passed: %d (%.1f%%)\n", r.Passed, float64(r.Passed)/float64(r.Total)*100)
	fmt.Printf("  Failed: %d (%.1f%%)\n", r.Failed, float64(r.Failed)/float64(r.Total)*100)

	if r.Failed == 0 {
		fmt.Println("  🎉 All tests passed!")
	} else {
		fmt.Printf("  ⚠️  %d tests failed\n", r.Failed)
	}
}

func runBasicTests() TestResults {
	tests := []TestCase{
		// Arithmetic
		{"1 + 2", nil, 3},
		{"10 - 5", nil, 5},
		{"3 * 4", nil, 12},
		{"15 / 3", nil, 5},
		{"2 + 3 * 4", nil, 14},   // Precedence
		{"(2 + 3) * 4", nil, 20}, // Parentheses

		// String operations
		{"\"hello\" + \" world\"", nil, "hello world"},
		{"\"test\" == \"test\"", nil, true},
		{"\"a\" < \"b\"", nil, true},

		// Boolean operations
		{"true && false", nil, false},
		{"true || false", nil, true},
		{"!true", nil, false},
		{"!false", nil, true},

		// Comparisons
		{"5 > 3", nil, true},
		{"2 < 1", nil, false},
		{"5 >= 5", nil, true},
		{"3 <= 2", nil, false},
		{"10 == 10", nil, true},
		{"7 != 8", nil, true},

		// Variables
		{"x", map[string]interface{}{"x": 42}, 42},
		{"name", map[string]interface{}{"name": "Alice"}, "Alice"},
		{"active", map[string]interface{}{"active": true}, true},

		// Mixed expressions
		{"x + y", map[string]interface{}{"x": 10, "y": 20}, 30},
		{"age > 18", map[string]interface{}{"age": 25}, true},
		{"name + \" is \" + string(age)", map[string]interface{}{"name": "Bob", "age": 30}, "Bob is 30"},
	}

	return runTests("Basic", tests)
}

func runTypeTests() TestResults {
	tests := []TestCase{
		// Type conversions
		{"string(42)", nil, "42"},
		{"int(\"123\")", nil, 123},
		{"float(\"3.14\")", nil, 3.14},
		{"bool(1)", nil, true},
		{"bool(0)", nil, false},

		// Type checking
		{"type(42)", nil, "int"},
		{"type(\"hello\")", nil, "string"},
		{"type(true)", nil, "bool"},
		{"type(3.14)", nil, "float"},

		// Mixed type operations
		{"1 + 2.5", nil, 3.5},
		{"\"Value: \" + string(42)", nil, "Value: 42"},
	}

	return runTests("Type", tests)
}

func runBuiltinTests() TestResults {
	tests := []TestCase{
		// Math functions
		{"abs(-42)", nil, 42},
		{"abs(42)", nil, 42},
		{"max(1, 5, 3)", nil, 5},
		{"min(1, 5, 3)", nil, 1},

		// String functions
		{"len(\"hello\")", nil, 5},
		{"contains(\"hello world\", \"world\")", nil, true},
		{"contains(\"hello\", \"xyz\")", nil, false},
		{"startsWith(\"hello\", \"he\")", nil, true},
		{"endsWith(\"world\", \"ld\")", nil, true},
		{"upper(\"hello\")", nil, "HELLO"},
		{"lower(\"WORLD\")", nil, "world"},
		{"trim(\"  test  \")", nil, "test"},

		// Utility functions
		{"count(\"hello\")", nil, 5},
		{"first(\"hello\")", nil, "h"},
		{"last(\"hello\")", nil, "o"},
	}

	return runTests("Builtin", tests)
}

func runErrorTests() TestResults {
	errorTests := []ErrorTestCase{
		// Syntax errors
		{"1 +", nil, "parse error"},
		{"(1 + 2", nil, "parse error"},
		{"1 + + 2", nil, "parse error"},

		// Type errors
		{"1 + \"hello\"", nil, "type error"},
		{"\"hello\" - 1", nil, "type error"},

		// Function errors
		{"unknownFunction()", nil, "unknown function"},
		{"abs()", nil, "wrong number of arguments"},
		{"abs(1, 2)", nil, "wrong number of arguments"},
	}

	results := TestResults{}
	for _, test := range errorTests {
		results.Total++
		fmt.Printf("    %-30s: ", test.Expression)

		_, err := expr.Eval(test.Expression, test.Env)
		if err != nil {
			fmt.Printf("✅ Error caught: %v\n", err)
			results.Passed++
		} else {
			fmt.Printf("❌ Expected error but got result\n")
			results.Failed++
		}
	}

	return results
}

func runPerformanceTests() TestResults {
	tests := []PerformanceTestCase{
		{"Simple arithmetic", "1 + 2 * 3", nil, 1000},
		{"String operations", "\"hello\" + \" \" + \"world\"", nil, 1000},
		{"Variable access", "x + y + z", map[string]interface{}{"x": 1, "y": 2, "z": 3}, 1000},
		{"Function calls", "abs(max(1, 2, 3))", nil, 1000},
		{"Complex expression", "len(name) > 3 && age >= 18", map[string]interface{}{"name": "Alice", "age": 25}, 1000},
	}

	results := TestResults{}
	for _, test := range tests {
		results.Total++
		fmt.Printf("    %-25s: ", test.Name)

		// Compile once
		program, err := expr.Compile(test.Expression, expr.Env(test.Env))
		if err != nil {
			fmt.Printf("❌ Compile error: %v\n", err)
			results.Failed++
			continue
		}

		// Time multiple executions
		start := time.Now()
		for i := 0; i < test.Iterations; i++ {
			_, err := expr.Run(program, test.Env)
			if err != nil {
				fmt.Printf("❌ Runtime error: %v\n", err)
				results.Failed++
				break
			}
		}

		if err == nil {
			duration := time.Since(start)
			avgTime := duration / time.Duration(test.Iterations)
			fmt.Printf("✅ %d ops in %v (avg: %v)\n", test.Iterations, duration, avgTime)
			results.Passed++
		}
	}

	return results
}

func runBenchmarks() {
	benchmarks := []BenchmarkCase{
		{"Compilation Speed", func() {
			expr.Compile("1 + 2 * 3", nil)
		}},
		{"Simple Execution", func() {
			program, _ := expr.Compile("1 + 2 * 3", nil)
			expr.Run(program, nil)
		}},
		{"Variable Access", func() {
			env := map[string]interface{}{"x": 10, "y": 20}
			program, _ := expr.Compile("x + y", expr.Env(env))
			expr.Run(program, env)
		}},
		{"String Operations", func() {
			program, _ := expr.Compile("\"hello\" + \" \" + \"world\"", nil)
			expr.Run(program, nil)
		}},
		{"Function Calls", func() {
			program, _ := expr.Compile("abs(max(1, 2, 3))", nil)
			expr.Run(program, nil)
		}},
	}

	for _, bench := range benchmarks {
		fmt.Printf("  %-20s: ", bench.Name)

		// Warmup
		for i := 0; i < 100; i++ {
			bench.Func()
		}

		// Benchmark
		iterations := 10000
		start := time.Now()
		for i := 0; i < iterations; i++ {
			bench.Func()
		}
		duration := time.Since(start)

		avgTime := duration / time.Duration(iterations)
		opsPerSec := float64(iterations) / duration.Seconds()

		fmt.Printf("%v per op, %.0f ops/sec\n", avgTime, opsPerSec)
	}
}

// Test case types
type TestCase struct {
	Expression string
	Env        map[string]interface{}
	Expected   interface{}
}

type ErrorTestCase struct {
	Expression string
	Env        map[string]interface{}
	ErrorType  string
}

type PerformanceTestCase struct {
	Name       string
	Expression string
	Env        map[string]interface{}
	Iterations int
}

type BenchmarkCase struct {
	Name string
	Func func()
}

func runTests(category string, tests []TestCase) TestResults {
	results := TestResults{}

	for _, test := range tests {
		results.Total++
		fmt.Printf("    %-30s: ", test.Expression)

		result, err := expr.Eval(test.Expression, test.Env)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			results.Failed++
			continue
		}

		if fmt.Sprintf("%v", result) == fmt.Sprintf("%v", test.Expected) {
			fmt.Printf("✅ %v\n", result)
			results.Passed++
		} else {
			fmt.Printf("❌ Expected %v, got %v\n", test.Expected, result)
			results.Failed++
		}
	}

	return results
}
