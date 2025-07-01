package main

import (
	"fmt"
	"time"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Performance Test Example ===")
	fmt.Println()

	// Test basic performance
	testBasicPerformance()
	fmt.Println()

	// Test environment performance
	testEnvironmentPerformance()
	fmt.Println()

	// Test complex expressions
	testComplexExpressions()
}

func testBasicPerformance() {
	fmt.Println("=== Basic Performance Test ===")

	expressions := []string{
		"1 + 2 * 3",
		"(10 + 5) * 2 - 3",
		"100 / 4 + 50 - 25",
		"true && false || true",
		"\"hello\" + \" \" + \"world\"",
		"abs(-42) + max(1, 2, 3)",
	}

	iterations := 10000

	for _, expr := range expressions {
		fmt.Printf("Testing: %s\n", expr)

		// Compile once
		start := time.Now()
		bytecode, err := compileExpression(expr)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			continue
		}
		compileTime := time.Since(start)

		// Execute many times
		start = time.Now()
		var result interface{}
		for i := 0; i < iterations; i++ {
			machine := vm.New(bytecode)
			err := machine.RunInstructions(bytecode.Instructions)
			if err != nil {
				fmt.Printf("Execution error: %v\n", err)
				break
			}
			result = machine.StackTop()
		}
		execTime := time.Since(start)

		fmt.Printf("  Compile time: %v\n", compileTime)
		fmt.Printf("  Execute time: %v (%d iterations)\n", execTime, iterations)
		fmt.Printf("  Avg per execution: %v\n", execTime/time.Duration(iterations))
		fmt.Printf("  Result: %v\n", result)
		fmt.Println()
	}
}

func testEnvironmentPerformance() {
	fmt.Println("=== Environment Performance Test ===")

	// Create environment
	envVars := map[string]interface{}{
		"x":      42,
		"y":      3.14,
		"name":   "test",
		"active": true,
	}

	adapter := env.New()

	expressions := []string{
		"x + 10",
		"y * 2.0",
		"x > 40 && active",
		"name == \"test\" ? x * 2 : x / 2",
	}

	iterations := 5000

	for _, expr := range expressions {
		fmt.Printf("Testing: %s\n", expr)

		// Compile with environment
		start := time.Now()
		bytecode, err := compileExpressionWithEnv(expr, envVars, adapter)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			continue
		}
		compileTime := time.Since(start)

		// Execute many times
		start = time.Now()
		var result interface{}
		for i := 0; i < iterations; i++ {
			machine := vm.NewWithEnvironment(bytecode, envVars, adapter)
			err := machine.RunInstructions(bytecode.Instructions)
			if err != nil {
				fmt.Printf("Execution error: %v\n", err)
				break
			}
			result = machine.StackTop()
		}
		execTime := time.Since(start)

		fmt.Printf("  Compile time: %v\n", compileTime)
		fmt.Printf("  Execute time: %v (%d iterations)\n", execTime, iterations)
		fmt.Printf("  Avg per execution: %v\n", execTime/time.Duration(iterations))
		fmt.Printf("  Result: %v\n", result)
		fmt.Println()
	}
}

func testComplexExpressions() {
	fmt.Println("=== Complex Expression Performance Test ===")

	complexExpressions := []string{
		"(1 + 2) * (3 + 4) - (5 + 6) / (7 + 8)",
		"abs(-10) + max(1, 2, 3) + min(4, 5, 6)",
		"true && (false || true) && !(false && true)",
		"\"prefix\" + \"_\" + \"middle\" + \"_\" + \"suffix\"",
	}

	iterations := 1000

	for _, expr := range complexExpressions {
		fmt.Printf("Testing complex: %s\n", expr)

		start := time.Now()
		bytecode, err := compileExpression(expr)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			continue
		}
		compileTime := time.Since(start)

		start = time.Now()
		var result interface{}
		for i := 0; i < iterations; i++ {
			machine := vm.New(bytecode)
			err := machine.RunInstructions(bytecode.Instructions)
			if err != nil {
				fmt.Printf("Execution error: %v\n", err)
				break
			}
			result = machine.StackTop()
		}
		execTime := time.Since(start)

		fmt.Printf("  Compile time: %v\n", compileTime)
		fmt.Printf("  Execute time: %v (%d iterations)\n", execTime, iterations)
		fmt.Printf("  Avg per execution: %v\n", execTime/time.Duration(iterations))
		fmt.Printf("  Ops per second: %.0f\n", float64(iterations)/execTime.Seconds())
		fmt.Printf("  Result: %v\n", result)
		fmt.Println()
	}

	fmt.Println("=== Performance Summary ===")
	fmt.Println("✓ Zero reflection implementation")
	fmt.Println("✓ Optimized bytecode compilation")
	fmt.Println("✓ Fast stack-based virtual machine")
	fmt.Println("✓ Environment variable support")
	fmt.Println("✓ Built-in function library")
	fmt.Println("✓ Performance monitoring ready")
}

func compileExpression(expr string) (*vm.Bytecode, error) {
	// Parse
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("no statements found")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		return nil, fmt.Errorf("expected expression statement")
	}

	// Compile
	compiler := compiler.New()
	err := compiler.Compile(stmt.Expression)
	if err != nil {
		return nil, err
	}

	return compiler.Bytecode(), nil
}

func compileExpressionWithEnv(expr string, envVars map[string]interface{}, adapter *env.Adapter) (*vm.Bytecode, error) {
	// Parse
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("no statements found")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		return nil, fmt.Errorf("expected expression statement")
	}

	// Compile with environment
	compiler := compiler.New()
	err := compiler.AddEnvironment(envVars, adapter)
	if err != nil {
		return nil, err
	}

	err = compiler.Compile(stmt.Expression)
	if err != nil {
		return nil, err
	}

	return compiler.Bytecode(), nil
}
