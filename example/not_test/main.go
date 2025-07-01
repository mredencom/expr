package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== NOT Operation Test ===")

	tests := []struct {
		expression string
		expected   bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!(1 > 2)", true},
		{"!(2 > 1)", false},
	}

	for _, test := range tests {
		fmt.Printf("%-12s: ", test.expression)

		result, err := expr.Eval(test.expression, nil)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		if result == test.expected {
			fmt.Printf("✅ %v\n", result)
		} else {
			fmt.Printf("❌ Expected %v, got %v\n", test.expected, result)
		}
	}
}
