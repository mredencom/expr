package checker

import (
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
)

func TestNewChecker(t *testing.T) {
	checker := New()
	if checker == nil {
		t.Fatal("Expected non-nil checker")
	}
	if checker.scope == nil {
		t.Fatal("Expected non-nil scope")
	}
	if checker.errors == nil {
		t.Fatal("Expected non-nil errors slice")
	}
}

func TestNewWithScope(t *testing.T) {
	scope := NewRootScope()
	checker := NewWithScope(scope)
	if checker == nil {
		t.Fatal("Expected non-nil checker")
	}
	if checker.scope != scope {
		t.Fatal("Expected checker to use provided scope")
	}
}

func TestCheckLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected types.TypeInfo
	}{
		{"42", types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		{"3.14", types.TypeInfo{Kind: types.KindFloat64, Name: "float"}},
		{"true", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{"false", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{`"hello"`, types.TypeInfo{Kind: types.KindString, Name: "string"}},
		{"null", types.TypeInfo{Kind: types.KindNil, Name: "nil"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			checker := New()

			err := checker.Check(program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			literal := stmt.Expression.(*ast.Literal)

			if literal.Value.Type().Kind != tt.expected.Kind {
				t.Errorf("Expected kind %v, got %v", tt.expected.Kind, literal.Value.Type().Kind)
			}
		})
	}
}

func TestCheckIdentifier(t *testing.T) {
	checker := New()

	// Define a variable
	varType := types.TypeInfo{Kind: types.KindString, Name: "string"}
	checker.scope.DefineVariable("x", varType)

	program := parseProgram(t, "x")

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	ident := stmt.Expression.(*ast.Identifier)

	if ident.TypeInfo.Kind != types.KindString {
		t.Errorf("Expected string type, got %v", ident.TypeInfo.Kind)
	}
}

func TestCheckUndefinedIdentifier(t *testing.T) {
	checker := New()
	program := parseProgram(t, "undefined_var")

	err := checker.Check(program)
	if err == nil {
		t.Fatal("Expected error for undefined variable")
	}

	errors := checker.Errors()
	if len(errors) == 0 {
		t.Fatal("Expected errors for undefined variable")
	}

	if errors[0] != "undefined variable: undefined_var" {
		t.Errorf("Unexpected error message: %s", errors[0])
	}
}

func TestCheckInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected types.TypeInfo
	}{
		{"5 + 3", types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		{"5.0 + 3.0", types.TypeInfo{Kind: types.KindFloat64, Name: "float"}},
		{"5 + 3.0", types.TypeInfo{Kind: types.KindFloat64, Name: "float"}},
		{"5 > 3", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{"5 == 3", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{"true && false", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{`"hello" + " world"`, types.TypeInfo{Kind: types.KindString, Name: "string"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			checker := New()

			err := checker.Check(program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			infix := stmt.Expression.(*ast.InfixExpression)

			if infix.TypeInfo.Kind != tt.expected.Kind {
				t.Errorf("Expected kind %v, got %v", tt.expected.Kind, infix.TypeInfo.Kind)
			}
		})
	}
}

func TestCheckPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected types.TypeInfo
	}{
		{"-5", types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		{"-5.0", types.TypeInfo{Kind: types.KindFloat64, Name: "float"}},
		{"!true", types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
		{"~42", types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			checker := New()

			err := checker.Check(program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			prefix := stmt.Expression.(*ast.PrefixExpression)

			if prefix.TypeInfo.Kind != tt.expected.Kind {
				t.Errorf("Expected kind %v, got %v", tt.expected.Kind, prefix.TypeInfo.Kind)
			}
		})
	}
}

func TestCheckArrayLiteral(t *testing.T) {
	program := parseProgram(t, "[1, 2, 3]")
	checker := New()

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array := stmt.Expression.(*ast.ArrayLiteral)

	if array.TypeInfo.Kind != types.KindSlice {
		t.Errorf("Expected slice type, got %v", array.TypeInfo.Kind)
	}
}

func TestCheckMapLiteral(t *testing.T) {
	program := parseProgram(t, `{"key": "value", "num": 42}`)
	checker := New()

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapLit := stmt.Expression.(*ast.MapLiteral)

	if mapLit.TypeInfo.Kind != types.KindMap {
		t.Errorf("Expected map type, got %v", mapLit.TypeInfo.Kind)
	}
}

func TestCheckBuiltinExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected types.TypeInfo
	}{
		{`len("hello")`, types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		{`string(42)`, types.TypeInfo{Kind: types.KindString, Name: "string"}},
		{`int("42")`, types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		{`float(42)`, types.TypeInfo{Kind: types.KindFloat64, Name: "float"}},
		{`bool(1)`, types.TypeInfo{Kind: types.KindBool, Name: "bool"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			checker := New()

			err := checker.Check(program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			builtin := stmt.Expression.(*ast.BuiltinExpression)

			if builtin.TypeInfo.Kind != tt.expected.Kind {
				t.Errorf("Expected kind %v, got %v", tt.expected.Kind, builtin.TypeInfo.Kind)
			}
		})
	}
}

func TestCheckConditionalExpression(t *testing.T) {
	program := parseProgram(t, "5 > 3 ? 10 : 20")
	checker := New()

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	cond := stmt.Expression.(*ast.ConditionalExpression)

	if cond.TypeInfo.Kind != types.KindInt64 {
		t.Errorf("Expected int type, got %v", cond.TypeInfo.Kind)
	}
}

func TestCheckIndexExpression(t *testing.T) {
	checker := New()

	// Define array variable
	arrayType := types.TypeInfo{
		Kind:     types.KindSlice,
		Name:     "[]int",
		ElemType: &types.TypeInfo{Kind: types.KindInt64, Name: "int"},
	}
	checker.scope.DefineVariable("arr", arrayType)

	program := parseProgram(t, "arr[0]")

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	index := stmt.Expression.(*ast.IndexExpression)

	if index.TypeInfo.Kind != types.KindInt64 {
		t.Errorf("Expected int type, got %v", index.TypeInfo.Kind)
	}
}

func TestCheckMemberExpression(t *testing.T) {
	checker := New()

	// Define a struct type with fields
	structType := types.TypeInfo{
		Kind: types.KindStruct,
		Name: "User",
		Fields: []types.FieldInfo{
			{Name: "name", Type: types.TypeInfo{Kind: types.KindString, Name: "string"}},
			{Name: "age", Type: types.TypeInfo{Kind: types.KindInt64, Name: "int"}},
		},
	}
	checker.scope.DefineVariable("user", structType)

	program := parseProgram(t, "user.name")

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	member := stmt.Expression.(*ast.MemberExpression)

	if member.TypeInfo.Kind != types.KindString {
		t.Errorf("Expected string type, got %v", member.TypeInfo.Kind)
	}
}

func TestWithEnvironment(t *testing.T) {
	checker := New()

	env := map[string]types.TypeInfo{
		"x": {Kind: types.KindInt64, Name: "int"},
		"y": {Kind: types.KindString, Name: "string"},
	}

	checker.WithEnvironment(env)

	program := parseProgram(t, "x + 5")

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestTypeInferenceErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{`"hello" + 5`, "invalid operation"},
		{`5 && true`, "invalid operation"},
		{`!42`, "invalid operation"},
		{`arr[0]`, "undefined variable: arr"},
		{`obj.field`, "undefined variable: obj"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			checker := New()

			err := checker.Check(program)
			if err == nil {
				t.Fatal("Expected error but got none")
			}

			errors := checker.Errors()
			if len(errors) == 0 {
				t.Fatal("Expected errors but got none")
			}
		})
	}
}

func TestCheckExpression(t *testing.T) {
	checker := New()
	program := parseProgram(t, "5 + 3")

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	typeInfo, err := checker.CheckExpression(stmt.Expression)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if typeInfo.Kind != types.KindInt64 {
		t.Errorf("Expected int type, got %v", typeInfo.Kind)
	}
}

func TestScopeAccess(t *testing.T) {
	checker := New()
	scope := checker.Scope()

	if scope == nil {
		t.Fatal("Expected non-nil scope")
	}

	// Test that scope is accessible
	scope.DefineVariable("test", types.TypeInfo{Kind: types.KindString, Name: "string"})

	if _, ok := scope.LookupVariable("test"); !ok {
		t.Error("Expected to find defined variable")
	}
}

func TestComplexExpression(t *testing.T) {
	checker := New()

	// Define variables
	env := map[string]types.TypeInfo{
		"x":    {Kind: types.KindInt64, Name: "int"},
		"y":    {Kind: types.KindFloat64, Name: "float"},
		"name": {Kind: types.KindString, Name: "string"},
	}
	checker.WithEnvironment(env)

	program := parseProgram(t, `x > 5 ? y * 2.0 : float(x)`)

	err := checker.Check(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	cond := stmt.Expression.(*ast.ConditionalExpression)

	if cond.TypeInfo.Kind != types.KindFloat64 {
		t.Errorf("Expected float type, got %v", cond.TypeInfo.Kind)
	}
}

// Helper functions

func parseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	if len(program.Statements) == 0 {
		t.Fatal("No statements parsed")
	}

	return program
}
