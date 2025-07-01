package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== Simple Collection Functions Test ===")

	// Simple test cases
	testCases := []struct {
		name       string
		expression string
		env        map[string]interface{}
	}{
		{
			name:       "sum_simple",
			expression: "sum([1, 2, 3, 4, 5])",
			env:        map[string]interface{}{},
		},
		{
			name:       "count_numbers",
			expression: "count(numbers)",
			env:        map[string]interface{}{"numbers": []int{1, 2, 3, 4, 5}},
		},
		{
			name:       "first_numbers",
			expression: "first(numbers)",
			env:        map[string]interface{}{"numbers": []int{10, 20, 30}},
		},
		{
			name:       "last_numbers",
			expression: "last(numbers)",
			env:        map[string]interface{}{"numbers": []int{10, 20, 30}},
		},
		{
			name:       "all_booleans",
			expression: "all(bools)",
			env:        map[string]interface{}{"bools": []bool{true, true, true}},
		},
		{
			name:       "any_booleans",
			expression: "any(bools)",
			env:        map[string]interface{}{"bools": []bool{false, true, false}},
		},
		{
			name:       "count_string",
			expression: "count(\"hello\")",
			env:        map[string]interface{}{},
		},
		{
			name:       "type_of_number",
			expression: "type(42)",
			env:        map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n--- %s ---\n", tc.name)
		fmt.Printf("Expression: %s\n", tc.expression)

		program, err := expr.Compile(tc.expression, expr.Env(tc.env))
		if err != nil {
			fmt.Printf("❌ Compilation error: %v\n", err)
			continue
		}

		result, err := expr.Run(program, tc.env)
		if err != nil {
			fmt.Printf("❌ Runtime error: %v\n", err)
			continue
		}

		fmt.Printf("✅ Result: %v\n", result)
	}

	fmt.Println("\n=== Test Complete ===")
}
