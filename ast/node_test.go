package ast

import (
	"testing"

	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
)

// TestProgramNode tests Program node
func TestProgramNode(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &Identifier{Value: "test"},
			},
		},
	}

	// Test Position method
	pos := program.Position()
	if pos.Line != 0 || pos.Column != 0 {
		t.Errorf("Expected position (0,0), got (%d,%d)", pos.Line, pos.Column)
	}

	// Test String method
	str := program.String()
	if str != "program" {
		t.Errorf("Expected 'program', got '%s'", str)
	}

	// Test with no statements
	emptyProgram := &Program{Statements: []Statement{}}
	if emptyProgram.String() != "program" {
		t.Error("Empty program should return 'program'")
	}
}

// TestIdentifierNode tests Identifier node
func TestIdentifierNode(t *testing.T) {
	pos := lexer.Position{Line: 1, Column: 5}
	identifier := &Identifier{
		Value: "myVar",
		Pos:   pos,
	}

	// Test Position method
	nodePos := identifier.Position()
	if nodePos.Line != 1 || nodePos.Column != 5 {
		t.Errorf("Expected position (1,5), got (%d,%d)", nodePos.Line, nodePos.Column)
	}

	// Test String method
	str := identifier.String()
	if str != "myVar" {
		t.Errorf("Expected 'myVar', got %s", str)
	}
}

// TestLiteralExpressions tests various literal expressions
func TestLiteralExpressions(t *testing.T) {
	// Test Integer Literal
	intLit := &Literal{
		Value: types.NewInt(42),
	}

	if intLit.String() != "42" {
		t.Errorf("Expected '42', got %s", intLit.String())
	}

	// Test Float Literal
	floatLit := &Literal{
		Value: types.NewFloat(3.14),
	}

	if floatLit.String() != "3.14" {
		t.Errorf("Expected '3.14', got %s", floatLit.String())
	}

	// Test String Literal
	stringLit := &Literal{
		Value: types.NewString("hello"),
	}

	if stringLit.String() != "hello" {
		t.Errorf("Expected 'hello', got %s", stringLit.String())
	}

	// Test Boolean Literal
	boolLit := &Literal{
		Value: types.NewBool(true),
	}

	if boolLit.String() != "true" {
		t.Errorf("Expected 'true', got %s", boolLit.String())
	}
}

// TestInfixExpression tests infix expressions
func TestInfixExpression(t *testing.T) {
	left := &Literal{
		Value: types.NewInt(5),
	}

	right := &Literal{
		Value: types.NewInt(3),
	}

	infix := &InfixExpression{
		Left:     left,
		Operator: "+",
		Right:    right,
	}

	// Test String method
	expected := "(5 + 3)"
	if infix.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, infix.String())
	}

	// Test with different operators
	operators := []string{"-", "*", "/", "==", "!=", "<", ">", "<=", ">="}
	for _, op := range operators {
		infix.Operator = op
		str := infix.String()
		if str == "" {
			t.Errorf("Infix expression with operator '%s' should not be empty", op)
		}
	}
}

// TestPrefixExpression tests prefix expressions
func TestPrefixExpression(t *testing.T) {
	operand := &Literal{
		Value: types.NewInt(5),
	}

	prefix := &PrefixExpression{
		Operator: "-",
		Right:    operand,
	}

	// Test String method
	expected := "(-5)"
	if prefix.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, prefix.String())
	}

	// Test with NOT operator
	prefix.Operator = "!"
	expected = "(!5)"
	if prefix.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, prefix.String())
	}
}

// TestCallExpression tests function call expressions
func TestCallExpression(t *testing.T) {
	function := &Identifier{
		Value: "add",
	}

	args := []Expression{
		&Literal{Value: types.NewInt(1)},
		&Literal{Value: types.NewInt(2)},
	}

	call := &CallExpression{
		Function:  function,
		Arguments: args,
	}

	// Test String method
	str := call.String()
	expected := "add(1, 2)"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}

	// Test with no arguments
	call.Arguments = []Expression{}
	str = call.String()
	expected = "add()"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

// TestIndexExpression tests index expressions
func TestIndexExpression(t *testing.T) {
	left := &Identifier{
		Value: "arr",
	}

	index := &Literal{
		Value: types.NewInt(0),
	}

	indexExpr := &IndexExpression{
		Left:  left,
		Index: index,
	}

	// Test String method
	expected := "(arr[0])"
	if indexExpr.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, indexExpr.String())
	}
}

// TestMemberExpression tests member access expressions
func TestMemberExpression(t *testing.T) {
	object := &Identifier{
		Value: "user",
	}

	property := &Identifier{
		Value: "name",
	}

	member := &MemberExpression{
		Object:   object,
		Property: property,
	}

	// Test String method
	expected := "user.name"
	if member.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, member.String())
	}
}

// TestConditionalExpression tests conditional (ternary) expressions
func TestConditionalExpression(t *testing.T) {
	condition := &Literal{
		Value: types.NewBool(true),
	}

	trueExpr := &Literal{
		Value: types.NewInt(1),
	}

	falseExpr := &Literal{
		Value: types.NewInt(2),
	}

	conditional := &ConditionalExpression{
		Test:        condition,
		Consequent:  trueExpr,
		Alternative: falseExpr,
	}

	// Test String method
	str := conditional.String()
	expected := "(true ? 1 : 2)"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

// TestArrayLiteral tests array literal expressions
func TestArrayLiteral(t *testing.T) {
	elements := []Expression{
		&Literal{Value: types.NewInt(1)},
		&Literal{Value: types.NewInt(2)},
		&Literal{Value: types.NewInt(3)},
	}

	array := &ArrayLiteral{
		Elements: elements,
	}

	// Test String method
	str := array.String()
	expected := "[1, 2, 3]"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}

	// Test empty array
	array.Elements = []Expression{}
	str = array.String()
	if str != "[]" {
		t.Errorf("Expected '[]', got '%s'", str)
	}
}

// TestMapLiteral tests map literal expressions
func TestMapLiteral(t *testing.T) {
	pairs := []MapPair{
		{
			Key:   &Literal{Value: types.NewString("name")},
			Value: &Literal{Value: types.NewString("John")},
		},
		{
			Key:   &Literal{Value: types.NewString("age")},
			Value: &Literal{Value: types.NewInt(30)},
		},
	}

	mapLit := &MapLiteral{
		Pairs: pairs,
	}

	// Test String method
	str := mapLit.String()
	if str == "" {
		t.Error("Map literal string should not be empty")
	}

	// Test empty map
	mapLit.Pairs = []MapPair{}
	str = mapLit.String()
	if str != "{}" {
		t.Errorf("Expected '{}', got '%s'", str)
	}
}

// TestBuiltinExpression tests builtin function expressions
func TestBuiltinExpression(t *testing.T) {
	args := []Expression{
		&Literal{Value: types.NewString("hello")},
	}

	builtin := &BuiltinExpression{
		Name:      "len",
		Arguments: args,
	}

	// Test String method
	str := builtin.String()
	expected := "len(hello)"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}

	// Test with no arguments
	builtin.Arguments = []Expression{}
	str = builtin.String()
	expected = "len()"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

// TestVariableExpression tests variable expressions
func TestVariableExpression(t *testing.T) {
	varExpr := &VariableExpression{
		Name: "x",
	}

	// Test String method
	if varExpr.String() != "x" {
		t.Errorf("Expected 'x', got '%s'", varExpr.String())
	}
}

// TestExpressionStatement tests expression statements
func TestExpressionStatement(t *testing.T) {
	expr := &Identifier{Value: "test"}
	stmt := &ExpressionStatement{
		Expression: expr,
	}

	// Test String method
	if stmt.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", stmt.String())
	}

	// Test Type method - should return the expression's type
	typeInfo := stmt.Type()
	// Since Identifier doesn't have TypeInfo set, it returns zero value
	if typeInfo.Name != "" {
		t.Errorf("Expected empty type name, got %v", typeInfo.Name)
	}
}

// TestNodeTypeAssertions tests that nodes implement correct interfaces
func TestNodeTypeAssertions(t *testing.T) {
	// Test that all expressions implement Expression interface
	expressions := []Expression{
		&Identifier{},
		&Literal{},
		&InfixExpression{},
		&PrefixExpression{},
		&CallExpression{},
		&IndexExpression{},
		&MemberExpression{},
		&ConditionalExpression{},
		&ArrayLiteral{},
		&MapLiteral{},
		&BuiltinExpression{},
		&VariableExpression{},
	}

	for i, expr := range expressions {
		if expr == nil {
			t.Errorf("Expression %d is nil", i)
		}
		// Test that it implements Node interface
		if _, ok := expr.(Node); !ok {
			t.Errorf("Expression %d does not implement Node interface", i)
		}
	}

	// Test that statements implement Statement interface
	statements := []Statement{
		&Program{},
		&ExpressionStatement{},
	}

	for i, stmt := range statements {
		if stmt == nil {
			t.Errorf("Statement %d is nil", i)
		}
		// Test that it implements Node interface
		if _, ok := stmt.(Node); !ok {
			t.Errorf("Statement %d does not implement Node interface", i)
		}
	}
}

// TestComplexExpression tests a complex nested expression
func TestComplexExpression(t *testing.T) {
	// Build: (a + b) * (c - d)
	a := &Identifier{Value: "a"}
	b := &Identifier{Value: "b"}
	c := &Identifier{Value: "c"}
	d := &Identifier{Value: "d"}

	leftInfix := &InfixExpression{
		Left:     a,
		Operator: "+",
		Right:    b,
	}

	rightInfix := &InfixExpression{
		Left:     c,
		Operator: "-",
		Right:    d,
	}

	complex := &InfixExpression{
		Left:     leftInfix,
		Operator: "*",
		Right:    rightInfix,
	}

	// Test String method
	str := complex.String()
	expected := "((a + b) * (c - d))"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}
