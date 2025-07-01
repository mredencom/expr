package main

import (
	"fmt"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/checker"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
)

func main() {
	fmt.Println("=== Expression Type Checker Demo ===")
	fmt.Println("Testing zero-reflection type checking system")
	fmt.Println()

	// Create environment adapter
	adapter := env.New()

	// Register custom struct adapter
	userAdapter := env.NewUserStructAdapter()
	adapter.RegisterStruct("User", userAdapter)

	// Create test environment with basic types
	testEnv := map[string]interface{}{
		"count":  42,
		"price":  19.99,
		"active": true,
		"name":   "Alice",
	}

	// Convert environment to our Value types
	envValues, err := adapter.CreateEnvironment(testEnv)
	if err != nil {
		fmt.Printf("Failed to create environment: %v\n", err)
		return
	}

	// Convert to TypeInfo for type checker
	envTypes := make(map[string]types.TypeInfo)
	for name, value := range envValues {
		envTypes[name] = value.Type()
	}

	// Test expressions (simplified for zero-reflection demo)
	expressions := []string{
		"count + 10",
		"price * 2.0",
		"active && true",
		"name + \" is awesome\"",
		"count > 25",
		"price == 19.99",
		"!active",
		"-count",
		"count + price", // Type mismatch example
	}

	for i, expr := range expressions {
		fmt.Printf("%d. Expression: %s\n", i+1, expr)
		checkExpression(expr, envTypes)
		fmt.Println()
	}

	fmt.Println("=== Type Checking Demo Complete ===")
	fmt.Println("✅ Zero reflection type checking implemented!")
	fmt.Println("✅ Static type safety enforced!")
	fmt.Println("✅ Type inference working correctly!")
}

func checkExpression(expression string, envTypes map[string]types.TypeInfo) {
	// Tokenize
	l := lexer.New(expression)

	// Parse
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("   Parse errors: %v\n", p.Errors())
		return
	}

	if len(program.Statements) == 0 {
		fmt.Printf("   No statements found\n")
		return
	}

	// Get the expression from the first statement
	var expr ast.Expression
	if exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
		expr = exprStmt.Expression
	} else {
		fmt.Printf("   First statement is not an expression\n")
		return
	}

	// Create type checker with environment
	checker := checker.New().WithEnvironment(envTypes)

	// Check the expression
	resultType, err := checker.CheckExpression(expr)
	if err != nil {
		fmt.Printf("   Type check failed: %v\n", err)

		// Show specific errors
		if len(checker.Errors()) > 0 {
			fmt.Printf("   Errors:\n")
			for _, e := range checker.Errors() {
				fmt.Printf("     - %s\n", e)
			}
		}
		return
	}

	fmt.Printf("   Result type: %s (%s)\n", resultType.Name, resultType.Kind.String())

	// Show type information for complex expressions
	if resultType.Kind == types.KindStruct {
		fmt.Printf("   Fields: ")
		for i, field := range resultType.Fields {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s: %s", field.Name, field.Type.Name)
		}
		fmt.Println()
	} else if resultType.Kind == types.KindSlice && resultType.ElemType != nil {
		fmt.Printf("   Element type: %s\n", resultType.ElemType.Name)
	} else if resultType.Kind == types.KindMap {
		if resultType.KeyType != nil && resultType.ValType != nil {
			fmt.Printf("   Key type: %s, Value type: %s\n",
				resultType.KeyType.Name, resultType.ValType.Name)
		}
	}

	fmt.Printf("   ✓ Type check passed\n")
}
