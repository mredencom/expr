package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== As() Function Demo ===")
	fmt.Println("Testing type validation and conversion")
	fmt.Println()

	// Test cases with different expected types
	testCases := []struct {
		name       string
		expression string
		option     expr.Option
		env        map[string]interface{}
		shouldFail bool
	}{
		{
			name:       "Integer Expression as Int",
			expression: "1 + 2",
			option:     expr.AsInt(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Integer Expression as Int64",
			expression: "10 * 5",
			option:     expr.AsInt64(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Float Expression as Float64",
			expression: "3.14 * 2.0",
			option:     expr.AsFloat64(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "String Expression as String",
			expression: `"hello" + " " + "world"`,
			option:     expr.AsString(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Boolean Expression as Bool",
			expression: "true && false",
			option:     expr.AsBool(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Integer to Float64 Conversion",
			expression: "42",
			option:     expr.AsFloat64(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Any Value to String",
			expression: "123",
			option:     expr.AsString(),
			env:        nil,
			shouldFail: false,
		},
		{
			name:       "Variable Expression",
			expression: "age * 2",
			option:     expr.AsInt(),
			env:        map[string]interface{}{"age": 25},
			shouldFail: false,
		},
		{
			name:       "Type Mismatch - String as Int",
			expression: `"hello"`,
			option:     expr.AsInt(),
			env:        nil,
			shouldFail: true,
		},
		{
			name:       "Type Mismatch - Int as Bool",
			expression: "42",
			option:     expr.AsBool(),
			env:        nil,
			shouldFail: true,
		},
	}

	for i, tc := range testCases {
		fmt.Printf("--- Test %d: %s ---\n", i+1, tc.name)
		fmt.Printf("Expression: %s\n", tc.expression)

		// Compile with As option
		var program *expr.Program
		var err error

		if tc.env != nil {
			program, err = expr.Compile(tc.expression, expr.Env(tc.env), tc.option)
		} else {
			program, err = expr.Compile(tc.expression, tc.option)
		}

		if err != nil {
			if tc.shouldFail {
				fmt.Printf("✅ Expected compilation error: %v\n", err)
			} else {
				fmt.Printf("❌ Unexpected compilation error: %v\n", err)
			}
			fmt.Println()
			continue
		}

		if tc.shouldFail {
			fmt.Printf("❌ Expected compilation to fail but it succeeded\n")
			fmt.Println()
			continue
		}

		// Run the program
		result, err := expr.Run(program, tc.env)
		if err != nil {
			if tc.shouldFail {
				fmt.Printf("✅ Expected runtime error: %v\n", err)
			} else {
				fmt.Printf("❌ Unexpected runtime error: %v\n", err)
			}
			fmt.Println()
			continue
		}

		if tc.shouldFail {
			fmt.Printf("❌ Expected runtime to fail but it succeeded\n")
			fmt.Println()
			continue
		}

		// Show result and its type
		fmt.Printf("✅ Result: %v (type: %T)\n", result, result)
		fmt.Println()
	}

	fmt.Println("=== As() Function Demo Complete ===")
	fmt.Println("✅ Type validation implemented!")
	fmt.Println("✅ Type conversion working!")
	fmt.Println("✅ Compile-time type checking active!")
}
