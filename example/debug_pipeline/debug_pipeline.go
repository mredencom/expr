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
	fmt.Println("=== Pipeline Operation Debug ===")

	// Test expressions
	expressions := []string{
		"numbers | len",
		"len(numbers)",
	}

	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5},
	}

	for _, expr := range expressions {
		fmt.Printf("\n--- Debugging: %s ---\n", expr)
		debugPipelineExpression(expr, env)
	}
}

func debugPipelineExpression(expr string, env map[string]interface{}) {
	// Parse
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Parse errors: %v\n", p.Errors())
		return
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	expression := stmt.Expression

	// Print AST structure
	fmt.Printf("AST Type: %T\n", expression)
	fmt.Printf("AST String: %s\n", expression.String())

	// If it's a pipe expression, examine the parts
	if pipe, ok := expression.(*ast.PipeExpression); ok {
		fmt.Printf("  Left: %T = %s\n", pipe.Left, pipe.Left.String())
		fmt.Printf("  Right: %T = %s\n", pipe.Right, pipe.Right.String())

		// Check if right side is builtin or identifier
		if builtin, ok := pipe.Right.(*ast.BuiltinExpression); ok {
			fmt.Printf("    Builtin name: %s\n", builtin.Name)
			fmt.Printf("    Builtin args: %d\n", len(builtin.Arguments))
		}
		if ident, ok := pipe.Right.(*ast.Identifier); ok {
			fmt.Printf("    Identifier: %s\n", ident.Value)
		}
	}

	// Compile
	comp := compiler.New()

	// Add environment variables as custom builtins
	for name := range env {
		comp.DefineBuiltin(name) // Add custom builtin
	}

	err := comp.Compile(expression)
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		return
	}

	bytecode := comp.Bytecode()
	fmt.Printf("Constants: %d\n", len(bytecode.Constants))
	for i, c := range bytecode.Constants {
		fmt.Printf("  [%d] %T: %s\n", i, c, c.String())
	}

	fmt.Printf("Instructions: %d bytes\n", len(bytecode.Instructions))

	// Decode instructions
	instructions := bytecode.Instructions
	fmt.Printf("Bytecode:\n")
	for i := 0; i < len(instructions); {
		op := vm.Opcode(instructions[i])
		fmt.Printf("  %04d: %s", i, op.String())

		// Decode operands based on opcode
		switch op {
		case vm.OpConstant:
			if i+2 < len(instructions) {
				constIndex := int(instructions[i+1])<<8 | int(instructions[i+2])
				fmt.Printf(" %d", constIndex)
				i += 3
			} else {
				i++
			}
		case vm.OpGetVar:
			if i+2 < len(instructions) {
				varIndex := int(instructions[i+1])<<8 | int(instructions[i+2])
				fmt.Printf(" %d", varIndex)
				i += 3
			} else {
				i++
			}
		case vm.OpBuiltin:
			if i+2 < len(instructions) {
				builtinIndex := int(instructions[i+1])
				argCount := int(instructions[i+2])
				fmt.Printf(" %d %d", builtinIndex, argCount)
				i += 3
			} else {
				i++
			}
		default:
			i++
		}
		fmt.Println()
	}

	fmt.Println("--- End Debug ---")
}
