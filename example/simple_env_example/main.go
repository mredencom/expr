package main

import (
	"fmt"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Simple Environment Example ===")
	fmt.Println()

	// Create simple environment with basic variables
	envVars := map[string]interface{}{
		"count":  42,
		"price":  19.99,
		"name":   "Alice",
		"active": true,
	}

	// Create environment adapter
	adapter := env.New()

	// Test expressions with simple variables
	expressions := []string{
		"count",
		"price",
		"name",
		"active",
		"count + 10",
		"price * 2.0",
		"count > 40",
		"name == \"Alice\"",
		"active && count > 30",
		"count > 50 ? \"high\" : \"low\"",
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

	fmt.Println("=== Environment Variables Summary ===")
	for name, value := range envVars {
		fmt.Printf("%s = %v (%T)\n", name, value, value)
	}
}
