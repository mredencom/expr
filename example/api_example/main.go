package main

import (
	"fmt"
	"log"
	"time"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== expr API Example ===")

	// Test basic evaluation
	fmt.Println("\n1. Basic Evaluation:")
	testBasicEval()

	// Test with environment
	fmt.Println("\n2. Environment Variables:")
	testEnvironment()

	// Test compilation and reuse
	fmt.Println("\n3. Compilation and Reuse:")
	testCompilationReuse()

	// Test configuration options
	fmt.Println("\n4. Configuration Options:")
	testConfigOptions()

	// Test result details
	fmt.Println("\n5. Detailed Results:")
	testDetailedResults()

	// Test statistics
	fmt.Println("\n6. Performance Statistics:")
	testStatistics()

	// Test compatibility API
	fmt.Println("\n7. Compatibility API:")
	testCompatibilityAPI()

	// Test new generic type checking
	fmt.Println("\n8. Generic Type Checking:")
	testGenericTypeChecking()
}

func testBasicEval() {
	expressions := []string{
		"1 + 2 * 3",
		"10 > 5",
		"\"hello\" + \" \" + \"world\"",
		"abs(-42)",
		"max(1, 2, 3, 4, 5)",
	}

	for _, exprStr := range expressions {
		result, err := expr.Eval(exprStr, nil)
		if err != nil {
			fmt.Printf("  %s = ERROR: %v\n", exprStr, err)
		} else {
			fmt.Printf("  %s = %v\n", exprStr, result)
		}
	}
}

func testEnvironment() {
	env := map[string]interface{}{
		"name":   "Alice",
		"age":    30,
		"active": true,
		"score":  95.5,
	}

	expressions := []string{
		"name",
		"age > 25",
		"active && age >= 18",
		"score * 0.1",
		"\"Hello, \" + name + \"!\"",
	}

	for _, exprStr := range expressions {
		result, err := expr.Eval(exprStr, env)
		if err != nil {
			fmt.Printf("  %s = ERROR: %v\n", exprStr, err)
		} else {
			fmt.Printf("  %s = %v\n", exprStr, result)
		}
	}
}

func testCompilationReuse() {
	// Compile once, run multiple times
	program, err := expr.Compile("age * factor", expr.Env(map[string]interface{}{
		"age":    0,
		"factor": 0,
	}))
	if err != nil {
		log.Printf("Compilation error: %v", err)
		return
	}

	fmt.Printf("  Compiled program: %s\n", program.String())
	fmt.Printf("  Compile time: %v\n", program.CompileTime())

	// Run with different environments
	environments := []map[string]interface{}{
		{"age": 25, "factor": 2.0},
		{"age": 30, "factor": 1.5},
		{"age": 35, "factor": 1.2},
	}

	for _, env := range environments {
		result, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("  env %v = ERROR: %v\n", env, err)
		} else {
			fmt.Printf("  env %v = %v\n", env, result)
		}
	}
}

func testConfigOptions() {
	// Test with custom built-in function
	customFn := func(x float64) float64 {
		return x * x
	}

	program, err := expr.Compile("square(5.0)",
		expr.WithBuiltin("square", customFn),
		expr.EnableOptimization(),
		expr.EnableCache(),
	)
	if err != nil {
		log.Printf("Compilation error: %v", err)
		return
	}

	result, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("  square(5.0) = ERROR: %v\n", err)
	} else {
		fmt.Printf("  square(5.0) = %v\n", result)
	}

	// Test with timeout
	program2, err := expr.Compile("1 + 2", expr.WithTimeout(time.Millisecond))
	if err != nil {
		log.Printf("Compilation error: %v", err)
		return
	}

	result2, err := expr.Run(program2, nil)
	if err != nil {
		fmt.Printf("  1 + 2 (with timeout) = ERROR: %v\n", err)
	} else {
		fmt.Printf("  1 + 2 (with timeout) = %v\n", result2)
	}
}

func testDetailedResults() {
	expressions := []string{
		"42",
		"3.14159",
		"\"Hello, World!\"",
		"true",
		"1 + 2 * 3",
	}

	for _, exprStr := range expressions {
		result, err := expr.EvalWithResult(exprStr, nil)
		if err != nil {
			fmt.Printf("  %s = ERROR: %v\n", exprStr, err)
		} else {
			fmt.Printf("  %s = %v (type: %s, time: %v)\n",
				exprStr, result.Value, result.Type, result.ExecutionTime)
		}
	}
}

func testStatistics() {
	// Reset statistics
	expr.ResetStatistics()

	// Run some expressions to generate statistics
	for i := 0; i < 5; i++ {
		_, _ = expr.Eval("1 + 2", nil)
		_, _ = expr.Eval("\"hello\" + \" world\"", nil)
	}

	stats := expr.GetStatistics()
	fmt.Printf("  Total Compilations: %d\n", stats.TotalCompilations)
	fmt.Printf("  Total Executions: %d\n", stats.TotalExecutions)
	fmt.Printf("  Average Compile Time: %v\n", stats.AverageCompileTime)
	fmt.Printf("  Average Execution Time: %v\n", stats.AverageExecTime)
}

func testCompatibilityAPI() {
	// Test deprecated functions
	env := expr.NewEnv()
	env["x"] = 10
	env["y"] = 20

	program, err := expr.CompileWithEnv("x + y", env)
	if err != nil {
		log.Printf("Compilation error: %v", err)
		return
	}

	result, err := expr.RunWithEnv(program, env)
	if err != nil {
		fmt.Printf("  x + y = ERROR: %v\n", err)
	} else {
		fmt.Printf("  x + y = %v (using compatibility API)\n", result)
	}

	// Test type functions (now without reflection)
	fmt.Printf("  Type of 42: %s\n", expr.GetType(42))
	fmt.Printf("  Type of \"hello\": %s\n", expr.GetType("hello"))
	fmt.Printf("  IsNil(nil): %v\n", expr.IsNil(nil))
	fmt.Printf("  IsNil(42): %v\n", expr.IsNil(42))

	// Test struct to map conversion with Mappable interface
	type Person struct {
		Name string
		Age  int
	}

	// Implement ToMap method to make Person mappable
	person := struct {
		Name string
		Age  int
	}{Name: "Bob", Age: 25}

	// Create a simple mappable version
	mappablePerson := SimpleMappable{
		"Name": person.Name,
		"Age":  person.Age,
	}

	personMap := expr.ToMap(mappablePerson)
	fmt.Printf("  Struct to map: %v\n", personMap)
}

func testGenericTypeChecking() {
	// Test generic type checking
	var value interface{} = 42

	// Check if value is an int
	err := expr.CheckType[int](value)
	if err != nil {
		fmt.Printf("  CheckType[int](42): ERROR: %v\n", err)
	} else {
		fmt.Printf("  CheckType[int](42): OK\n")
	}

	// Check if value is a string (should fail)
	err = expr.CheckType[string](value)
	if err != nil {
		fmt.Printf("  CheckType[string](42): Expected error: %v\n", err)
	} else {
		fmt.Printf("  CheckType[string](42): Unexpected success\n")
	}

	// Test type conversion
	converted, err := expr.ConvertType[string](42)
	if err != nil {
		fmt.Printf("  ConvertType[string](42): ERROR: %v\n", err)
	} else {
		fmt.Printf("  ConvertType[string](42): %s\n", converted)
	}

	// Test int64 to int conversion
	var int64Value interface{} = int64(123)
	convertedInt, err := expr.ConvertType[int](int64Value)
	if err != nil {
		fmt.Printf("  ConvertType[int](int64(123)): ERROR: %v\n", err)
	} else {
		fmt.Printf("  ConvertType[int](int64(123)): %d\n", convertedInt)
	}
}

// SimpleMappable is a helper type that implements the Mappable interface
type SimpleMappable map[string]interface{}

func (s SimpleMappable) ToMap() map[string]interface{} {
	return map[string]interface{}(s)
}
