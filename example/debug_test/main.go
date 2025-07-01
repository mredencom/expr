package main

import (
	"fmt"

	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

func main() {
	// Test different string expressions
	expressions := []string{
		`config["debug"]`,
		`user.Age > 18 ? "adult" : "minor"`,
		`contains(user.Name, "Ali")`,
	}

	for i, expr := range expressions {
		fmt.Printf("%d. Testing: %s\n", i+1, expr)

		// Tokenize
		l := lexer.New(expr)
		fmt.Printf("   Tokens: ")
		for {
			tok := l.NextToken()
			fmt.Printf("%s ", tok.Type.String())
			if tok.Type == lexer.EOF {
				break
			}
		}
		fmt.Println()

		// Parse
		l2 := lexer.New(expr)
		p := parser.New(l2)
		_ = p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("   Parse errors: %v\n", p.Errors())
		} else {
			fmt.Printf("   ✓ Parsed successfully\n")
		}
		fmt.Println()
	}
}
