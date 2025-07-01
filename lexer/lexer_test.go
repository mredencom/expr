package lexer

import (
	"testing"
)

// TestNewLexer tests lexer creation
func TestNewLexer(t *testing.T) {
	input := "1 + 2"
	lexer := New(input)

	if lexer == nil {
		t.Fatal("Expected lexer to be created")
	}

	if lexer.input != input {
		t.Errorf("Expected input '%s', got '%s'", input, lexer.input)
	}

	if lexer.line != 1 {
		t.Errorf("Expected line 1, got %d", lexer.line)
	}

	if lexer.column != 2 {
		t.Errorf("Expected column 2, got %d", lexer.column)
	}
}

// TestNextToken tests basic tokenization
func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			tokenType TokenType
			value     string
		}
	}{
		{
			"1 + 2",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{NUMBER, "1"},
				{ADD, "+"},
				{NUMBER, "2"},
				{EOF, ""},
			},
		},
		{
			"hello world",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{IDENT, "hello"},
				{IDENT, "world"},
				{EOF, ""},
			},
		},
		{
			"\"hello world\"",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{STRING, "hello world"},
				{EOF, ""},
			},
		},
		{
			"== != < > <= >=",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{EQ, "=="},
				{NE, "!="},
				{LT, "<"},
				{GT, ">"},
				{LE, "<="},
				{GE, ">="},
				{EOF, ""},
			},
		},
		{
			"&& || !",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{AND, "&&"},
				{OR, "||"},
				{NOT, "!"},
				{EOF, ""},
			},
		},
		{
			"( ) [ ] { }",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{LPAREN, "("},
				{RPAREN, ")"},
				{LBRACKET, "["},
				{RBRACKET, "]"},
				{LBRACE, "{"},
				{RBRACE, "}"},
				{EOF, ""},
			},
		},
		{
			", . ; :",
			[]struct {
				tokenType TokenType
				value     string
			}{
				{COMMA, ","},
				{DOT, "."},
				{SEMICOLON, ";"},
				{COLON, ":"},
				{EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		lexer := New(tt.input)

		for i, expected := range tt.expected {
			token := lexer.NextToken()

			if token.Type != expected.tokenType {
				t.Errorf("test[%s] token[%d] - tokentype wrong. expected=%q, got=%q",
					tt.input, i, expected.tokenType, token.Type)
			}

			if token.Value != expected.value {
				t.Errorf("test[%s] token[%d] - value wrong. expected=%q, got=%q",
					tt.input, i, expected.value, token.Value)
			}
		}
	}
}

// TestNumbers tests number tokenization
func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"0", "0"},
		{"123.456", "123.456"},
		{"0.5", "0.5"},
	}

	for _, tt := range tests {
		lexer := New(tt.input)
		token := lexer.NextToken()

		if token.Type != NUMBER {
			t.Errorf("Expected NUMBER token, got %q", token.Type)
		}

		if token.Value != tt.expected {
			t.Errorf("Expected value '%s', got '%s'", tt.expected, token.Value)
		}
	}
}

// TestStrings tests string tokenization
func TestStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\"", "hello"},
		{"\"hello world\"", "hello world"},
		{"\"\"", ""},
		{"\"with\\\"quotes\"", "with\"quotes"},
		{"\"with\\nlines\"", "with\nlines"},
		{"\"with\\ttabs\"", "with\ttabs"},
		{"\"with\\\\backslash\"", "with\\backslash"},
	}

	for _, tt := range tests {
		lexer := New(tt.input)
		token := lexer.NextToken()

		if token.Type != STRING {
			t.Errorf("Expected STRING token, got %q", token.Type)
		}

		if token.Value != tt.expected {
			t.Errorf("Expected value '%s', got '%s'", tt.expected, token.Value)
		}
	}
}

// TestIdentifiers tests identifier tokenization
func TestIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"hello", IDENT},
		{"_hello", IDENT},
		{"hello123", IDENT},
		{"_", IDENT},
		{"if", IF},
		{"else", ELSE},
		{"in", IN},
	}

	for _, tt := range tests {
		lexer := New(tt.input)
		token := lexer.NextToken()

		if token.Type != tt.expected {
			t.Errorf("Expected token type %q, got %q", tt.expected, token.Type)
		}
	}
}

// TestWhitespace tests whitespace handling
func TestWhitespace(t *testing.T) {
	input := "  \t\n\r  1  \t\n  2  "
	lexer := New(input)

	// Should skip whitespace
	token1 := lexer.NextToken()
	if token1.Type != NUMBER || token1.Value != "1" {
		t.Errorf("Expected NUMBER '1', got %q '%s'", token1.Type, token1.Value)
	}

	token2 := lexer.NextToken()
	if token2.Type != NUMBER || token2.Value != "2" {
		t.Errorf("Expected NUMBER '2', got %q '%s'", token2.Type, token2.Value)
	}

	token3 := lexer.NextToken()
	if token3.Type != EOF {
		t.Errorf("Expected EOF, got %q", token3.Type)
	}
}

// TestPosition tests position tracking
func TestPosition(t *testing.T) {
	input := "a\nb"
	lexer := New(input)

	token1 := lexer.NextToken()
	if token1.Position.Line != 1 || token1.Position.Column != 2 {
		t.Errorf("Expected position (1,2), got (%d,%d)",
			token1.Position.Line, token1.Position.Column)
	}

	token2 := lexer.NextToken()
	if token2.Position.Line != 2 || token2.Position.Column != 2 {
		t.Errorf("Expected position (2,2), got (%d,%d)",
			token2.Position.Line, token2.Position.Column)
	}
}

// TestComplexExpression tests complex expression tokenization
func TestComplexExpression(t *testing.T) {
	input := `user.age >= 18 && user.name != "" || admin == true`

	expected := []TokenType{
		IDENT, DOT, IDENT, GE, NUMBER, AND, IDENT, DOT, IDENT, NE, STRING, OR, IDENT, EQ, BOOL, EOF,
	}

	lexer := New(input)

	for i, expectedType := range expected {
		token := lexer.NextToken()
		if token.Type != expectedType {
			t.Errorf("Token[%d] - expected %q, got %q (value: %s)", i, expectedType, token.Type, token.Value)
		}
	}
}

// TestPeekChar tests peek functionality
func TestPeekChar(t *testing.T) {
	input := "abc"
	lexer := New(input)

	// Initial state
	if lexer.char != 'a' {
		t.Errorf("Expected current char 'a', got %c", lexer.char)
	}

	// Peek should return next char without advancing
	if lexer.peekChar() != 'b' {
		t.Errorf("Expected peek char 'b', got %c", lexer.peekChar())
	}

	// Current char should still be 'a'
	if lexer.char != 'a' {
		t.Errorf("Expected current char still 'a', got %c", lexer.char)
	}

	// Advance and check
	lexer.readChar()
	if lexer.char != 'b' {
		t.Errorf("Expected current char 'b', got %c", lexer.char)
	}
}

// TestArithmeticOperators tests arithmetic operators
func TestArithmeticOperators(t *testing.T) {
	input := "+ - * / % **"
	expected := []TokenType{ADD, SUB, MUL, DIV, MOD, POW, EOF}

	lexer := New(input)

	for i, expectedType := range expected {
		token := lexer.NextToken()
		if token.Type != expectedType {
			t.Errorf("Token[%d] - expected %q, got %q", i, expectedType, token.Type)
		}
	}
}

// TestBitwiseOperators tests bitwise operators
func TestBitwiseOperators(t *testing.T) {
	input := "& | ^ ~ << >>"
	expected := []TokenType{BIT_AND, BIT_OR, BIT_XOR, BIT_NOT, SHL, SHR, EOF}

	lexer := New(input)

	for i, expectedType := range expected {
		token := lexer.NextToken()
		if token.Type != expectedType {
			t.Errorf("Token[%d] - expected %q, got %q", i, expectedType, token.Type)
		}
	}
}

// TestTokenString tests Token.String() method
func TestTokenString(t *testing.T) {
	token := Token{Type: NUMBER, Value: "123"}
	expected := "NUMBER(123)"

	if token.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, token.String())
	}

	token2 := Token{Type: ADD, Value: "+"}
	expected2 := "+(+)"

	if token2.String() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, token2.String())
	}
}

// TestTokenMethods tests Token helper methods
func TestTokenMethods(t *testing.T) {
	// Test IsLiteral
	numberToken := Token{Type: NUMBER}
	if !numberToken.IsLiteral() {
		t.Error("Expected NUMBER token to be literal")
	}

	addToken := Token{Type: ADD}
	if addToken.IsLiteral() {
		t.Error("Expected ADD token to not be literal")
	}

	// Test IsOperator
	if !addToken.IsOperator() {
		t.Error("Expected ADD token to be operator")
	}

	if numberToken.IsOperator() {
		t.Error("Expected NUMBER token to not be operator")
	}

	// Test IsComparison
	eqToken := Token{Type: EQ}
	if !eqToken.IsComparison() {
		t.Error("Expected EQ token to be comparison")
	}

	if addToken.IsComparison() {
		t.Error("Expected ADD token to not be comparison")
	}

	// Test IsLogical
	andToken := Token{Type: AND}
	if !andToken.IsLogical() {
		t.Error("Expected AND token to be logical")
	}

	if numberToken.IsLogical() {
		t.Error("Expected NUMBER token to not be logical")
	}

	// Test IsKeyword
	ifToken := Token{Type: IF}
	if !ifToken.IsKeyword() {
		t.Error("Expected IF token to be keyword")
	}

	if numberToken.IsKeyword() {
		t.Error("Expected NUMBER token to not be keyword")
	}
}

// TestLookupIdent tests identifier lookup
func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident    string
		expected TokenType
	}{
		{"if", IF},
		{"else", ELSE},
		{"in", IN},
		{"hello", IDENT},
		{"variable", IDENT},
	}

	for _, tt := range tests {
		result := LookupIdent(tt.ident)
		if result != tt.expected {
			t.Errorf("LookupIdent(%s) - expected %q, got %q", tt.ident, tt.expected, result)
		}
	}
}

// TestReset tests lexer reset functionality
func TestReset(t *testing.T) {
	lexer := New("old input")

	// Consume some tokens
	lexer.NextToken()
	lexer.NextToken()

	// Reset with new input
	newInput := "new input"
	lexer.Reset(newInput)

	if lexer.input != newInput {
		t.Errorf("Expected input '%s', got '%s'", newInput, lexer.input)
	}

	if lexer.line != 1 {
		t.Errorf("Expected line 1 after reset, got %d", lexer.line)
	}

	if lexer.column != 2 {
		t.Errorf("Expected column 2 after reset, got %d", lexer.column)
	}
}

// TestError tests error token creation
func TestError(t *testing.T) {
	lexer := New("test")
	errorMsg := "test error"

	errorToken := lexer.Error(errorMsg)

	if errorToken.Type != ILLEGAL {
		t.Errorf("Expected ILLEGAL token type, got %q", errorToken.Type)
	}

	expectedValue := "lexer error: " + errorMsg
	if errorToken.Value != expectedValue {
		t.Errorf("Expected error value '%s', got '%s'", expectedValue, errorToken.Value)
	}
}
