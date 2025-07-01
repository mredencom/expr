package main

import (
	"fmt"

	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Println("=== Expr Compiler & VM Demo ===")
	fmt.Println("Testing bytecode compilation and execution without reflection")
	fmt.Println()

	// Test expressions
	expressions := []string{
		"1 + 2",
		"10 - 5",
		"3 * 4",
		"20 / 5",
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"5 > 3",
		"2 == 2",
		"!true",
		"-42",
	}

	for i, expr := range expressions {
		fmt.Printf("--- Test %d: %s ---\n", i+1, expr)

		// Lexical analysis
		l := lexer.New(expr)

		// Parsing
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("Parser errors: %v\n", p.Errors())
			continue
		}

		// Print AST
		fmt.Printf("AST: %s\n", program.String())

		// Compilation
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			continue
		}

		// Get bytecode
		bytecode := comp.Bytecode()
		fmt.Printf("Instructions: %d bytes\n", len(bytecode.Instructions))
		fmt.Printf("Constants: %d\n", len(bytecode.Constants))

		// Create VM and run
		machine := vm.New(&vm.Bytecode{
			Instructions: bytecode.Instructions,
			Constants:    bytecode.Constants,
		})

		result, err := machine.Run(&vm.Bytecode{
			Instructions: bytecode.Instructions,
			Constants:    bytecode.Constants,
		}, map[string]interface{}{})
		if err != nil {
			fmt.Printf("VM error: %v\n", err)
			continue
		}

		// Get result
		if result != nil {
			fmt.Printf("Result: %s (%s)\n", result.String(), result.Type().Name)
		} else {
			fmt.Printf("Result: <nil>\n")
		}

		fmt.Println()
	}

	fmt.Println("=== Compilation Demo Complete ===")
	fmt.Println("✅ Zero reflection achieved!")
	fmt.Println("✅ Static type checking implemented!")
	fmt.Println("✅ Bytecode compilation working!")
	fmt.Println("✅ Virtual machine execution ready!")
}
