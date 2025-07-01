package main

import (
	"fmt"
	"log"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// User represents a user struct for testing
type User struct {
	Name   string
	Age    int
	Active bool
	Score  float64
}

// ToMap implements the StructConverter interface for zero-reflection conversion
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Name":   u.Name,
		"Age":    u.Age,
		"Active": u.Active,
		"Score":  u.Score,
	}
}

func main() {
	fmt.Println("=== Environment Integration Example ===")
	fmt.Println()

	// Create test environment
	user := User{
		Name:   "Alice",
		Age:    30,
		Active: true,
		Score:  95.5,
	}

	count := 42
	price := 19.99
	tags := []string{"admin", "user", "premium"}
	config := map[string]interface{}{
		"debug":   true,
		"timeout": 30,
		"host":    "localhost",
	}

	// Create environment adapter
	adapter := env.New()

	// Register struct type
	userAdapter := env.NewUserStructAdapter()
	adapter.RegisterStruct("User", userAdapter)

	// Add variables to environment
	envVars := map[string]interface{}{
		"user":   user,
		"count":  count,
		"price":  price,
		"tags":   tags,
		"config": config,
	}

	// Test expressions with variables
	expressions := []string{
		"user.Name",
		"user.Age",
		"user.Active",
		"user.Score",
		"user.Age > 25",
		"user.Name == \"Alice\"",
		"count + 10",
		"price * 2.0",
		"user.Active && count > 40",
		"user.Score > 90.0 ? \"excellent\" : \"good\"",
		// "len(tags)",  // Will implement array support later
		// "tags[0]",    // Will implement indexing later
		// "config[\"debug\"]", // Will implement map access later
	}

	for i, expr := range expressions {
		fmt.Printf("Test %d: %s\n", i+1, expr)
		fmt.Println("----------------------------------------")

		// Parse
		l := lexer.New(expr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("Parser errors: %v\n", p.Errors())
			fmt.Println()
			continue
		}

		// Get the expression from the program
		if len(program.Statements) == 0 {
			fmt.Println("No statements found")
			fmt.Println()
			continue
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			fmt.Printf("Expected expression statement, got %T\n", program.Statements[0])
			fmt.Println()
			continue
		}

		expression := stmt.Expression
		fmt.Printf("AST: %s\n", expression.String())

		// Compile to bytecode
		compiler := compiler.New()

		// Add environment variables to compiler
		err := compiler.AddEnvironment(envVars, adapter)
		if err != nil {
			fmt.Printf("Environment error: %v\n", err)
			fmt.Println()
			continue
		}

		err = compiler.Compile(expression)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			fmt.Println()
			continue
		}

		bytecode := compiler.Bytecode()
		fmt.Printf("Bytecode: %d instructions, %d constants\n",
			len(bytecode.Instructions), len(bytecode.Constants))

		// Create VM with environment
		machine := vm.NewWithEnvironment(bytecode, envVars, adapter)
		err = machine.RunInstructions(bytecode.Instructions)
		if err != nil {
			fmt.Printf("VM error: %v\n", err)
			fmt.Println()
			continue
		}

		// Get result
		result := machine.StackTop()
		if result != nil {
			fmt.Printf("Result: %s (type: %s)\n", result.String(), result.Type().Name)
		} else {
			fmt.Println("Result: <nil>")
		}

		fmt.Println()
	}

	// Test complex expression
	fmt.Println("=== Complex Environment Expression ===")
	testComplexEnvironmentExpression(envVars, adapter)
}

func testComplexEnvironmentExpression(envVars map[string]interface{}, adapter *env.Adapter) {
	// Test: user.Age > 25 && user.Active ? user.Score * 1.1 : user.Score * 0.9
	expr := "user.Age > 25 && user.Active ? user.Score * 1.1 : user.Score * 0.9"
	fmt.Printf("Expression: %s\n", expr)

	// Parse
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		log.Fatalf("Parser errors: %v", p.Errors())
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	expression := stmt.Expression

	// Compile
	compiler := compiler.New()
	err := compiler.AddEnvironment(envVars, adapter)
	if err != nil {
		log.Fatalf("Environment error: %v", err)
	}

	err = compiler.Compile(expression)
	if err != nil {
		log.Fatalf("Compilation error: %v", err)
	}

	bytecode := compiler.Bytecode()

	// Show bytecode details
	fmt.Printf("Constants pool:\n")
	for i, constant := range bytecode.Constants {
		fmt.Printf("  [%d] %s (%s)\n", i, constant.String(), constant.Type().Name)
	}

	fmt.Printf("Instructions: %d bytes\n", len(bytecode.Instructions))

	// Execute
	machine := vm.NewWithEnvironment(bytecode, envVars, adapter)
	err = machine.RunInstructions(bytecode.Instructions)
	if err != nil {
		log.Fatalf("VM error: %v", err)
	}

	result := machine.StackTop()
	fmt.Printf("Result: %s\n", result.String())

	// Verify the result: Alice is 30 (>25) and Active, so score should be 95.5 * 1.1 = 105.05
	if floatResult, ok := result.(*types.FloatValue); ok {
		expected := 95.5 * 1.1
		if floatResult.Value() == expected {
			fmt.Printf("✓ Correct result: %.2f\n", expected)
		} else {
			fmt.Printf("✗ Incorrect result: expected %.2f, got %.2f\n", expected, floatResult.Value())
		}
	} else {
		fmt.Printf("✗ Expected float result, got %T\n", result)
	}
}
