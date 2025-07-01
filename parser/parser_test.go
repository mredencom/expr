package parser

import (
	"strings"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
)

func TestParseSimpleExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5"},
		{"true", "true"},
		{"false", "false"},
		{"foobar", "foobar"},
		{"\"hello world\"", "hello world"},
		{"null", "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			if stmt.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, stmt.String())
			}
		})
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Literal)
	if !ok {
		t.Fatalf("exp not *ast.Literal. got=%T", stmt.Expression)
	}

	intVal, ok := literal.Value.(*types.IntValue)
	if !ok {
		t.Fatalf("literal.Value not *types.IntValue. got=%T", literal.Value)
	}

	if intVal.Value() != 5 {
		t.Errorf("intVal.Value not %d. got=%d", 5, intVal.Value())
	}
}

func TestParseFloatLiteral(t *testing.T) {
	input := "3.14"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Literal)
	if !ok {
		t.Fatalf("exp not *ast.Literal. got=%T", stmt.Expression)
	}

	floatVal, ok := literal.Value.(*types.FloatValue)
	if !ok {
		t.Fatalf("literal.Value not *types.FloatValue. got=%T", literal.Value)
	}

	if floatVal.Value() != 3.14 {
		t.Errorf("floatVal.Value not %f. got=%f", 3.14, floatVal.Value())
	}
}

func TestParseBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			literal, ok := stmt.Expression.(*ast.Literal)
			if !ok {
				t.Fatalf("exp not *ast.Literal. got=%T", stmt.Expression)
			}

			boolVal, ok := literal.Value.(*types.BoolValue)
			if !ok {
				t.Fatalf("literal.Value not *types.BoolValue. got=%T", literal.Value)
			}

			if boolVal.Value() != tt.expected {
				t.Errorf("boolVal.Value not %t. got=%t", tt.expected, boolVal.Value())
			}
		})
	}
}

func TestParseStringLiteral(t *testing.T) {
	input := `'hello world'`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Literal)
	if !ok {
		t.Fatalf("exp not *ast.Literal. got=%T", stmt.Expression)
	}

	strVal, ok := literal.Value.(*types.StringValue)
	if !ok {
		t.Fatalf("literal.Value not *types.StringValue. got=%T", literal.Value)
	}

	if strVal.Value() != "hello world" {
		t.Errorf("strVal.Value not %q. got=%q", "hello world", strVal.Value())
	}
}

func TestParseIdentifier(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"~42", "~", 42},
	}

	for _, tt := range prefixTests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
			}

			if exp.Operator != tt.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
				return
			}
		})
	}
}

func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 % 5", 5, "%", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 && 5", 5, "&&", 5},
		{"5 || 5", 5, "||", 5},
	}

	for _, tt := range infixTests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},

		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) == 0 {
				t.Fatalf("no statements parsed")
			}

			actual := program.Statements[0].String()
			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		})
	}
}

func TestParseCallExpression(t *testing.T) {
	input := "myFunc(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.BuiltinExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.BuiltinExpression. got=%T", stmt.Expression)
	}

	if exp.Name != "myFunc" {
		t.Errorf("exp.Name not 'myFunc'. got=%s", exp.Name)
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParseArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParseMapLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	mapLit, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(mapLit.Pairs) != 3 {
		t.Fatalf("map literal pairs wrong. got=%d", len(mapLit.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for _, pair := range mapLit.Pairs {
		literal, ok := pair.Key.(*ast.Literal)
		if !ok {
			t.Errorf("key is not ast.Literal. got=%T", pair.Key)
			continue
		}

		strVal, ok := literal.Value.(*types.StringValue)
		if !ok {
			t.Errorf("key is not string. got=%T", literal.Value)
			continue
		}

		expectedValue := expected[strVal.Value()]
		testIntegerLiteral(t, pair.Value, expectedValue)
	}
}

func TestParseIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParseMemberExpression(t *testing.T) {
	input := "obj.property"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	memberExp, ok := stmt.Expression.(*ast.MemberExpression)
	if !ok {
		t.Fatalf("exp not *ast.MemberExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, memberExp.Object, "obj") {
		return
	}

	if !testIdentifier(t, memberExp.Property, "property") {
		return
	}
}

func TestParseConditionalExpression(t *testing.T) {
	input := "x > 5 ? 10 : 20"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	condExp, ok := stmt.Expression.(*ast.ConditionalExpression)
	if !ok {
		t.Fatalf("exp not *ast.ConditionalExpression. got=%T", stmt.Expression)
	}

	testInfixExpression(t, condExp.Test, "x", ">", 5)
	testIntegerLiteral(t, condExp.Consequent, 10)
	testIntegerLiteral(t, condExp.Alternative, 20)
}

func TestParseBuiltinExpression(t *testing.T) {
	input := `contains("hello", "lo")`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	builtinExp, ok := stmt.Expression.(*ast.BuiltinExpression)
	if !ok {
		t.Fatalf("exp not *ast.BuiltinExpression. got=%T", stmt.Expression)
	}

	if builtinExp.Name != "contains" {
		t.Errorf("builtin.Name not 'contains'. got=%s", builtinExp.Name)
	}

	if len(builtinExp.Arguments) != 2 {
		t.Fatalf("wrong number of arguments. got=%d", len(builtinExp.Arguments))
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"5 +", "no prefix parse function for EOF found"},
		{"!true == false", ""}, // This should parse correctly
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			p.ParseProgram()

			errors := p.Errors()
			if tt.expectedError == "" {
				if len(errors) != 0 {
					t.Errorf("Expected no errors, got %v", errors)
				}
			} else {
				if len(errors) == 0 {
					t.Errorf("Expected error containing %q, got no errors", tt.expectedError)
				} else {
					found := false
					for _, err := range errors {
						if strings.Contains(err, tt.expectedError) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error containing %q, got %v", tt.expectedError, errors)
					}
				}
			}
		})
	}
}

// Helper functions

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	literal, ok := il.(*ast.Literal)
	if !ok {
		t.Errorf("il not *ast.Literal. got=%T", il)
		return false
	}

	intVal, ok := literal.Value.(*types.IntValue)
	if !ok {
		t.Errorf("literal.Value not *types.IntValue. got=%T", literal.Value)
		return false
	}

	if intVal.Value() != value {
		t.Errorf("intVal.Value not %d. got=%d", value, intVal.Value())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	literal, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("exp not *ast.Literal. got=%T", exp)
		return false
	}

	boolVal, ok := literal.Value.(*types.BoolValue)
	if !ok {
		t.Errorf("literal.Value not *types.BoolValue. got=%T", literal.Value)
		return false
	}

	if boolVal.Value() != value {
		t.Errorf("boolVal.Value not %t. got=%t", value, boolVal.Value())
		return false
	}

	return true
}
