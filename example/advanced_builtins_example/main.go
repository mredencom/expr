package main

import (
	"fmt"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Advanced Built-in Functions Example ===")
	fmt.Println()

	// Test with environment that includes arrays
	testWithArrayEnvironment()
	fmt.Println()

	// Test regex functions
	testRegexFunctions()
	fmt.Println()

	// Test type introspection
	testTypeIntrospection()
}

func testWithArrayEnvironment() {
	fmt.Println("=== Collection Functions with Environment ===")

	// Create environment with arrays
	envVars := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5},
		"words":   []string{"hello", "world", "test"},
		"flags":   []bool{true, false, true},
		"mixed":   []interface{}{1, "hello", true},
		"empty":   []int{},
	}

	adapter := env.New()

	// For now, test expressions that work with the current implementation
	expressions := []string{
		"len(\"hello world\")",
		"contains(\"hello world\", \"world\")",
		"startsWith(\"hello\", \"he\")",
		"endsWith(\"world\", \"ld\")",
		"upper(\"hello\")",
		"lower(\"WORLD\")",
		"trim(\"  test  \")",
		"abs(-42)",
		"max(10, 20, 5)",
		"min(10, 20, 5)",
		"sum(1, 2, 3)",
		"string(123)",
		"int(\"456\")",
		"float(\"3.14\")",
		"bool(1)",
		"type(\"test\")",
		"first(\"hello\")",
		"last(\"world\")",
	}

	for _, expr := range expressions {
		testExpressionWithEnv(expr, envVars, adapter)
	}
}

func testRegexFunctions() {
	fmt.Println("=== Regex Functions ===")

	expressions := []string{
		"matches(\"hello123\", \"[a-z]+\")",
		"matches(\"hello123\", \"\\\\d+\")",
		"matches(\"test@example.com\", \".*@.*\\\\..*\")",
		"matches(\"abc123def\", \"\\\\d+\")",
		"matches(\"nodigits\", \"\\\\d+\")",
	}

	for _, expr := range expressions {
		testExpression(expr)
	}
}

func testTypeIntrospection() {
	fmt.Println("=== Type Introspection ===")

	expressions := []string{
		"type(42)",
		"type(3.14)",
		"type(\"hello\")",
		"type(true)",
		"string(type(42))",
		"string(type(3.14))",
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

func testExpressionWithEnv(expr string, envVars map[string]interface{}, adapter *env.Adapter) {
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

	// Compile with environment
	compiler := compiler.New()
	err := compiler.AddEnvironment(envVars, adapter)
	if err != nil {
		fmt.Printf("  ❌ Environment error: %v\n", err)
		fmt.Println()
		return
	}

	err = compiler.Compile(expression)
	if err != nil {
		fmt.Printf("  ❌ Compilation error: %v\n", err)
		fmt.Println()
		return
	}

	bytecode := compiler.Bytecode()

	// Execute
	machine := vm.NewWithEnvironment(bytecode, envVars, adapter)
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

// Helper function to create test arrays (for future use)
func createTestSlice() *types.SliceValue {
	elements := []types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
	}
	elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8}
	return types.NewSlice(elements, elemType)
}
