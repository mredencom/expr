package main

import (
	"fmt"
	"os"

	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"expression\"")
		fmt.Println("Example: go run main.go \"1 + 2 * 3\"")
		os.Exit(1)
	}

	input := os.Args[1]
	fmt.Printf("Input: %s\n", input)

	// Tokenize
	fmt.Println("\n=== Lexical Analysis ===")
	l := lexer.New(input)

	for {
		tok := l.NextToken()
		fmt.Printf("Token: %s\n", tok)
		if tok.Type == lexer.EOF {
			break
		}
	}

	// Parse
	fmt.Println("\n=== Parsing ===")
	l = lexer.New(input) // Reset lexer
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("AST: %s\n", program.String())

	if len(program.Statements) > 0 {
		fmt.Printf("First statement: %s\n", program.Statements[0].String())
	}

	fmt.Println("\n=== Successfully parsed! ===")
}
