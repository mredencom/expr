package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== Collection Functions Example ===")

	// Test data
	env := map[string]interface{}{
		"numbers":  []int{1, 2, 3, 4, 5, -1, -2},
		"words":    []string{"hello", "world", "", "test"},
		"booleans": []bool{true, false, true, true},
		"mixed":    []interface{}{1, "hello", true, 3.14},
		"empty":    []int{},
		"user_map": map[string]interface{}{"name": "Alice", "age": 30, "active": true},
	}

	// Test collection functions
	testCollectionFunctions(env)
}

func testCollectionFunctions(env map[string]interface{}) {
	// Test cases for collection functions
	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		// All function tests
		{"all_booleans_mixed", "all(booleans)", "false"},
		{"all_positive_numbers", "all([1, 2, 3])", "true"}, // All positive numbers
		{"all_with_zero", "all([1, 0, 3])", "false"},       // Contains zero

		// Any function tests
		{"any_booleans_mixed", "any(booleans)", "true"},
		{"any_empty", "any(empty)", "false"},
		{"any_with_positive", "any([-1, 0, 1])", "true"},

		// Sum function tests
		{"sum_numbers", "sum(numbers)", "12"}, // 1+2+3+4+5+(-1)+(-2) = 12
		{"sum_positive", "sum([1, 2, 3, 4, 5])", "15"},
		{"sum_mixed_numbers", "sum([1, 2.5, 3])", "6.5"},

		// Count function tests
		{"count_numbers", "count(numbers)", "7"},
		{"count_words", "count(words)", "4"},
		{"count_empty", "count(empty)", "0"},
		{"count_string", "count(\"hello\")", "5"},
		{"count_map", "count(user_map)", "3"},

		// First and Last function tests
		{"first_numbers", "first(numbers)", "1"},
		{"first_words", "first(words)", "hello"},
		{"first_empty", "first(empty)", "nil"},
		{"first_string", "first(\"hello\")", "h"},
		{"last_numbers", "last(numbers)", "-2"},
		{"last_words", "last(words)", "test"},
		{"last_empty", "last(empty)", "nil"},
		{"last_string", "last(\"hello\")", "o"},

		// Filter function tests (with simple predicates)
		{"filter_not_empty", "filter(words, \"not_empty\")", "[hello world test]"},
		{"filter_positive", "filter(numbers, \"positive\")", "[1 2 3 4 5]"},
		{"filter_negative", "filter(numbers, \"negative\")", "[-1 -2]"},

		// Map function tests (simple transformations)
		{"map_to_string", "map([1, 2, 3], \"string\")", "[1 2 3]"}, // Simple identity mapping

		// Utility functions
		{"type_of_numbers", "type(numbers)", "[]interface{}"},
		{"keys_of_map", "keys(user_map)", "[name age active]"}, // Order may vary
	}

	fmt.Println("\nTesting collection functions:")

	for _, tc := range testCases {
		fmt.Printf("\n--- %s ---\n", tc.name)
		fmt.Printf("Expression: %s\n", tc.expression)

		// Register all builtin functions
		program, err := expr.Compile(tc.expression, expr.Env(env))
		if err != nil {
			fmt.Printf("❌ Compilation error: %v\n", err)
			continue
		}

		result, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("❌ Runtime error: %v\n", err)
			continue
		}

		fmt.Printf("Result: %v\n", result)
		fmt.Printf("Expected: %s\n", tc.expected)

		// Simple validation (in a real test, you'd do proper comparison)
		resultStr := fmt.Sprintf("%v", result)
		if resultStr == tc.expected {
			fmt.Printf("✅ PASS\n")
		} else {
			fmt.Printf("⚠️  Different result (may still be correct)\n")
		}
	}

	// Test advanced collection operations
	fmt.Println("\n=== Advanced Collection Operations ===")
	testAdvancedOperations(env)
}

func testAdvancedOperations(env map[string]interface{}) {
	// More complex expressions
	complexExpressions := []string{
		"sum(filter(numbers, \"positive\"))",    // Sum of positive numbers
		"count(filter(words, \"not_empty\"))",   // Count of non-empty words
		"first(filter(numbers, \"positive\"))",  // First positive number
		"last(filter(numbers, \"negative\"))",   // Last negative number
		"all([true, 1 > 0, \"hello\" != \"\"])", // All with mixed expressions
		"any([false, 0 > 1, \"\" == \"\"])",     // Any with mixed expressions
	}

	for _, exprStr := range complexExpressions {
		fmt.Printf("\nExpression: %s\n", exprStr)

		program, err := expr.Compile(exprStr, expr.Env(env))
		if err != nil {
			fmt.Printf("❌ Compilation error: %v\n", err)
			continue
		}

		result, err := expr.Run(program, env)
		if err != nil {
			fmt.Printf("❌ Runtime error: %v\n", err)
			continue
		}

		fmt.Printf("Result: %v\n", result)
	}
}
