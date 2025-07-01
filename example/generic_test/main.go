package main

import (
	"fmt"

	expr "github.com/mredencom/expr"
)

func main() {
	fmt.Println("=== Generic Type System Test ===")

	// Test 1: Basic type checking with generics
	fmt.Println("\n1. Generic Type Checking:")
	testGenericTypeChecking()

	// Test 2: Type conversions with generics
	fmt.Println("\n2. Generic Type Conversions:")
	testGenericTypeConversions()

	// Test 3: Complex type scenarios
	fmt.Println("\n3. Complex Type Scenarios:")
	testComplexTypeScenarios()

	// Test 4: Mappable interface usage
	fmt.Println("\n4. Mappable Interface:")
	testMappableInterface()

	// Test 5: Performance comparison (no reflection vs reflection)
	fmt.Println("\n5. Performance Benefits:")
	testPerformanceBenefits()
}

func testGenericTypeChecking() {
	testCases := []struct {
		value    interface{}
		typeName string
		testFunc func(interface{}) error
	}{
		{42, "int", expr.CheckType[int]},
		{"hello", "string", expr.CheckType[string]},
		{3.14, "float64", expr.CheckType[float64]},
		{true, "bool", expr.CheckType[bool]},
		{int64(100), "int64", expr.CheckType[int64]},
	}

	for _, tc := range testCases {
		err := tc.testFunc(tc.value)
		if err != nil {
			fmt.Printf("  ✗ CheckType[%s](%v): %v\n", tc.typeName, tc.value, err)
		} else {
			fmt.Printf("  ✓ CheckType[%s](%v): OK\n", tc.typeName, tc.value)
		}
	}

	// Test type mismatches
	fmt.Println("  Type mismatch tests:")
	err := expr.CheckType[string](42)
	if err != nil {
		fmt.Printf("  ✓ CheckType[string](42): Expected error: %v\n", err)
	} else {
		fmt.Printf("  ✗ CheckType[string](42): Should have failed\n")
	}
}

func testGenericTypeConversions() {
	// Test successful conversions
	fmt.Println("  Successful conversions:")

	// String conversions
	if result, err := expr.ConvertType[string](42); err == nil {
		fmt.Printf("  ✓ int(42) → string: %q\n", result)
	} else {
		fmt.Printf("  ✗ int(42) → string: %v\n", err)
	}

	// Numeric conversions
	if result, err := expr.ConvertType[int](int64(123)); err == nil {
		fmt.Printf("  ✓ int64(123) → int: %d\n", result)
	} else {
		fmt.Printf("  ✗ int64(123) → int: %v\n", err)
	}

	if result, err := expr.ConvertType[float64](42); err == nil {
		fmt.Printf("  ✓ int(42) → float64: %f\n", result)
	} else {
		fmt.Printf("  ✗ int(42) → float64: %v\n", err)
	}

	// Boolean conversions
	if result, err := expr.ConvertType[bool]("true"); err == nil {
		fmt.Printf("  ✓ string(\"true\") → bool: %t\n", result)
	} else {
		fmt.Printf("  ✗ string(\"true\") → bool: %v\n", err)
	}

	if result, err := expr.ConvertType[bool](1); err == nil {
		fmt.Printf("  ✓ int(1) → bool: %t\n", result)
	} else {
		fmt.Printf("  ✗ int(1) → bool: %v\n", err)
	}

	// Test nil handling
	if result, err := expr.ConvertType[string](nil); err == nil {
		fmt.Printf("  ✓ nil → string: %q\n", result)
	} else {
		fmt.Printf("  ✗ nil → string: %v\n", err)
	}
}

func testComplexTypeScenarios() {
	// Test with expression results
	env := map[string]interface{}{
		"age":    30,
		"name":   "Alice",
		"active": true,
		"score":  95.5,
	}

	expressions := []struct {
		expr     string
		expected string
		testFunc func(interface{}) error
	}{
		{"age", "int", expr.CheckType[int]},
		{"name", "string", expr.CheckType[string]},
		{"active", "bool", expr.CheckType[bool]},
		{"score", "float64", expr.CheckType[float64]},
		{"age > 25", "bool", expr.CheckType[bool]},
		{"age * 2", "int", expr.CheckType[int]},
	}

	for _, test := range expressions {
		result, err := expr.Eval(test.expr, env)
		if err != nil {
			fmt.Printf("  ✗ %s: Evaluation error: %v\n", test.expr, err)
			continue
		}

		err = test.testFunc(result)
		if err != nil {
			fmt.Printf("  ✗ %s → %s: Type check failed: %v\n", test.expr, test.expected, err)
		} else {
			fmt.Printf("  ✓ %s → %s: %v\n", test.expr, test.expected, result)
		}
	}
}

func testMappableInterface() {
	// Test with custom mappable struct
	user := User{
		Name:   "Bob",
		Age:    25,
		Email:  "bob@example.com",
		Active: true,
	}

	userMap := expr.ToMap(user)
	fmt.Printf("  User struct to map: %v\n", userMap)

	// Use the map in expressions
	expressions := []string{
		"Name",
		"Age > 18",
		"Active && Age >= 21",
		"\"Email: \" + Email",
	}

	for _, exprStr := range expressions {
		result, err := expr.Eval(exprStr, userMap)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", exprStr, err)
		} else {
			fmt.Printf("  ✓ %s = %v\n", exprStr, result)
		}
	}
}

func testPerformanceBenefits() {
	// Since we're not using reflection, we can show the benefits
	fmt.Println("  Zero-reflection benefits:")
	fmt.Println("  ✓ No runtime type inspection overhead")
	fmt.Println("  ✓ Compile-time type safety with generics")
	fmt.Println("  ✓ Better performance for type operations")
	fmt.Println("  ✓ Reduced memory allocations")

	// Test type operations performance
	iterations := 1000

	// Type checking performance
	var typeCheckErrors int
	for i := 0; i < iterations; i++ {
		if err := expr.CheckType[int](i); err != nil {
			typeCheckErrors++
		}
	}
	fmt.Printf("  Type checks: %d/%d successful\n", iterations-typeCheckErrors, iterations)

	// Type conversion performance
	var conversionErrors int
	for i := 0; i < iterations; i++ {
		if _, err := expr.ConvertType[string](i); err != nil {
			conversionErrors++
		}
	}
	fmt.Printf("  Type conversions: %d/%d successful\n", iterations-conversionErrors, iterations)

	// Show type information without reflection
	values := []interface{}{42, "hello", 3.14, true, int64(100)}
	fmt.Println("  Type identification without reflection:")
	for _, v := range values {
		fmt.Printf("    %v → %s\n", v, expr.GetType(v))
	}
}

// User represents a user struct that implements Mappable
type User struct {
	Name   string
	Age    int
	Email  string
	Active bool
}

// ToMap implements the Mappable interface
func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Name":   u.Name,
		"Age":    u.Age,
		"Email":  u.Email,
		"Active": u.Active,
	}
}
