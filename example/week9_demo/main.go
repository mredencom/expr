package main

import (
	"fmt"
	"log"

	expr "github.com/mredencom/expr"
	"github.com/mredencom/expr/builtins"
)

func main() {
	fmt.Println("=== Week 9: Extended Built-in Functions Demo ===")

	// Demonstrate what we've achieved so far
	fmt.Println("\n🎯 Current Implementation Status:")
	fmt.Println("✅ Core expression evaluation")
	fmt.Println("✅ Static type checking")
	fmt.Println("✅ Bytecode compilation")
	fmt.Println("✅ Virtual machine execution")
	fmt.Println("✅ Zero-reflection API")
	fmt.Println("✅ Extended built-in functions")

	// Show available built-in functions
	fmt.Println("\n📚 Available Built-in Functions:")
	builtinNames := builtins.ListBuiltinNames()
	for i, name := range builtinNames {
		if i%5 == 0 && i > 0 {
			fmt.Println()
		}
		fmt.Printf("%-12s", name)
	}
	fmt.Println()

	// Test basic functionality that we know works
	fmt.Println("\n🧪 Testing Core Functionality:")
	testBasicFunctions()

	// Test string functions
	fmt.Println("\n📝 Testing String Functions:")
	testStringFunctions()

	// Test math functions
	fmt.Println("\n🔢 Testing Math Functions:")
	testMathFunctions()

	// Test type functions
	fmt.Println("\n🏷️  Testing Type Functions:")
	testTypeFunctions()

	// Show what's coming next
	fmt.Println("\n🚀 Next Steps (Week 10):")
	fmt.Println("- Complete array literal support")
	fmt.Println("- Improve environment type conversion")
	fmt.Println("- Add comprehensive test suite")
	fmt.Println("- Performance benchmarking")
	fmt.Println("- Documentation and examples")
}

func testBasicFunctions() {
	tests := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Simple arithmetic", "1 + 2 * 3", nil},
		{"String concatenation", "\"Hello\" + \" \" + \"World\"", nil},
		{"Boolean logic", "true && false || true", nil},
		{"Comparison", "10 > 5 && 3 < 7", nil},
		{"Variable access", "name", map[string]interface{}{"name": "Alice"}},
		{"Mixed expression", "age > 18 && active", map[string]interface{}{"age": 25, "active": true}},
	}

	for _, test := range tests {
		runTest(test.name, test.expr, test.env)
	}
}

func testStringFunctions() {
	tests := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"String length", "len(\"hello world\")", nil},
		{"String contains", "contains(\"hello world\", \"world\")", nil},
		{"String starts with", "startsWith(\"hello\", \"he\")", nil},
		{"String ends with", "endsWith(\"world\", \"ld\")", nil},
		{"String upper case", "upper(\"hello\")", nil},
		{"String lower case", "lower(\"HELLO\")", nil},
		{"String trim", "trim(\"  hello  \")", nil},
		{"Variable string length", "len(message)", map[string]interface{}{"message": "Hello, World!"}},
	}

	for _, test := range tests {
		runTest(test.name, test.expr, test.env)
	}
}

func testMathFunctions() {
	tests := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Absolute value", "abs(-42)", nil},
		{"Maximum of numbers", "max(1, 5, 3, 9, 2)", nil},
		{"Minimum of numbers", "min(1, 5, 3, 9, 2)", nil},
		{"Variable math", "abs(value)", map[string]interface{}{"value": -15}},
		{"Complex expression", "max(abs(a), abs(b))", map[string]interface{}{"a": -10, "b": 7}},
	}

	for _, test := range tests {
		runTest(test.name, test.expr, test.env)
	}
}

func testTypeFunctions() {
	tests := []struct {
		name string
		expr string
		env  map[string]interface{}
	}{
		{"Type of integer", "type(42)", nil},
		{"Type of string", "type(\"hello\")", nil},
		{"Type of boolean", "type(true)", nil},
		{"Type of float", "type(3.14)", nil},
		{"Type conversion int", "int(\"123\")", nil},
		{"Type conversion string", "string(456)", nil},
		{"Type conversion float", "float(\"3.14\")", nil},
		{"Type conversion bool", "bool(1)", nil},
	}

	for _, test := range tests {
		runTest(test.name, test.expr, test.env)
	}
}

func runTest(name, expression string, env map[string]interface{}) {
	fmt.Printf("  %-25s: ", name)

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		fmt.Printf("❌ Compile Error: %v\n", err)
		return
	}

	result, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("❌ Runtime Error: %v\n", err)
		return
	}

	fmt.Printf("✅ %v\n", result)
}

func init() {
	// Ensure we're using the latest built-in functions
	log.SetFlags(0)
}
