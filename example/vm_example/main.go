package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Virtual Machine Example ===")
	fmt.Println()

	// Test expressions
	expressions := []string{
		"1 + 2",
		"10 - 5",
		"3 * 4",
		"15 / 3",
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"10 > 5",
		"3 == 3",
		"5 != 3",
		"true && false",
		"true || false",
		"!true",
		"-42",
		`"hello" + " " + "world"`,
		"1.5 + 2.5",
		"abs(-10)",
		"max(1, 2, 3)",
		"min(5, 3, 8)",
		"len(\"hello\")",
	}

	for i, expr := range expressions {
		fmt.Printf("Test %d: %s\n", i+1, expr)
		fmt.Println(strings.Repeat("-", 40))

		// Lexical analysis
		l := lexer.New(expr)
		fmt.Printf("Tokens: ")
		var tokens []lexer.Token
		for {
			tok := l.NextToken()
			tokens = append(tokens, tok)
			fmt.Printf("%s ", tok.Type)
			if tok.Type == lexer.EOF {
				break
			}
		}
		fmt.Println()

		// Parse
		l = lexer.New(expr) // Reset lexer
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
		err := compiler.Compile(expression)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			fmt.Println()
			continue
		}

		bytecode := compiler.Bytecode()
		fmt.Printf("Bytecode: %d instructions, %d constants\n",
			len(bytecode.Instructions), len(bytecode.Constants))

		// Create VM and run
		machine := vm.New(bytecode)
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

	// Test complex expression with variables
	fmt.Println("=== Complex Expression Test ===")
	testComplexExpression()
}

func testComplexExpression() {
	// Test: 2 * (3 + 4) > 10
	expr := "2 * (3 + 4) > 10"
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
	err := compiler.Compile(expression)
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
	machine := vm.New(bytecode)
	err = machine.RunInstructions(bytecode.Instructions)
	if err != nil {
		log.Fatalf("VM error: %v", err)
	}

	result := machine.StackTop()
	fmt.Printf("Result: %s\n", result.String())

	// Verify the result: 2 * (3 + 4) = 2 * 7 = 14, 14 > 10 = true
	if boolResult, ok := result.(*types.BoolValue); ok {
		if boolResult.Value() {
			fmt.Println("✓ Correct result: true")
		} else {
			fmt.Println("✗ Incorrect result: expected true")
		}
	} else {
		fmt.Printf("✗ Expected boolean result, got %T\n", result)
	}
}
