package main

import (
	"fmt"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Built-in Functions Library Example ===")
	fmt.Println()

	// Test type conversion functions
	testTypeConversions()
	fmt.Println()

	// Test math functions
	testMathFunctions()
	fmt.Println()

	// Test string functions
	testStringFunctions()
	fmt.Println()

	// Test new functions
	testNewFunctions()
}

func testTypeConversions() {
	fmt.Println("=== Type Conversion Functions ===")

	expressions := []string{
		"string(42)",
		"string(3.14)",
		"string(true)",
		"int(\"123\")",
		"int(3.14)",
		"int(true)",
		"float(\"3.14\")",
		"float(42)",
		"bool(1)",
		"bool(0)",
		"bool(\"hello\")",
		"bool(\"\")",
	}

	for _, expr := range expressions {
		testExpression(expr)
	}
}

func testMathFunctions() {
	fmt.Println("=== Math Functions ===")

	expressions := []string{
		"abs(-42)",
		"abs(3.14)",
		"max(1, 2, 3)",
		"max(3.14, 2.71, 1.41)",
		"min(5, 3, 8)",
		"min(3.14, 2.71, 1.41)",
		"sum(1, 2, 3, 4, 5)",
		"sum(1.1, 2.2, 3.3)",
	}

	for _, expr := range expressions {
		testExpression(expr)
	}
}

func testStringFunctions() {
	fmt.Println("=== String Functions ===")

	expressions := []string{
		"len(\"hello\")",
		"contains(\"hello world\", \"world\")",
		"startsWith(\"hello\", \"he\")",
		"endsWith(\"hello\", \"lo\")",
		"upper(\"hello\")",
		"lower(\"HELLO\")",
		"trim(\"  hello  \")",
		"type(\"hello\")",
		"type(42)",
		"type(true)",
	}

	for _, expr := range expressions {
		testExpression(expr)
	}
}

func testNewFunctions() {
	fmt.Println("=== New Functions (String manipulation) ===")

	expressions := []string{
		"first(\"hello\")",
		"last(\"hello\")",
		"matches(\"hello\", \"h.*o\")",
		"matches(\"test123\", \"\\\\d+\")",
	}

	for _, expr := range expressions {
		testExpression(expr)
	}
}

func testExpression(expr string) {
	fmt.Printf("Testing: %s\n", expr)

	// Parse
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("  ❌ Parser errors: %v\n", p.Errors())
		fmt.Println()
		return
	}

	if len(program.Statements) == 0 {
		fmt.Printf("  ❌ No statements found\n")
		fmt.Println()
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		fmt.Printf("  ❌ Expected expression statement, got %T\n", program.Statements[0])
		fmt.Println()
		return
	}

	expression := stmt.Expression

	// Compile
	compiler := compiler.New()
	err := compiler.Compile(expression)
	if err != nil {
		fmt.Printf("  ❌ Compilation error: %v\n", err)
		fmt.Println()
		return
	}

	bytecode := compiler.Bytecode()

	// Execute
	machine := vm.New(bytecode)
	err = machine.RunInstructions(bytecode.Instructions)
	if err != nil {
		fmt.Printf("  ❌ VM error: %v\n", err)
		fmt.Println()
		return
	}

	// Get result
	result := machine.StackTop()
	if result != nil {
		fmt.Printf("  ✅ Result: %s (type: %s)\n", result.String(), result.Type().Name)
	} else {
		fmt.Printf("  ❌ Result: <nil>\n")
	}

	fmt.Println()
}
