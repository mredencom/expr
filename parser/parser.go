package parser

import (
	"fmt"
	"strconv"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
)

// Parser implements a recursive descent parser using Pratt parsing for expressions
type Parser struct {
	lexer *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string

	// Parser functions
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.BOOL, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullLiteral)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.SUB, p.parsePrefixExpression)
	p.registerPrefix(lexer.BIT_NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseMapLiteral)
	p.registerPrefix(lexer.WILDCARD, p.parseWildcard)
	p.registerPrefix(lexer.PLACEHOLDER, p.parsePlaceholder)
	// Add builtin functions
	p.registerPrefix(lexer.CONTAINS, p.parseBuiltinFunction)
	p.registerPrefix(lexer.STARTS_WITH, p.parseBuiltinFunction)
	p.registerPrefix(lexer.ENDS_WITH, p.parseBuiltinFunction)
	p.registerPrefix(lexer.MATCHES, p.parseBuiltinFunction)

	// Initialize infix parse functions
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.ADD, p.parseInfixExpression)
	p.registerInfix(lexer.SUB, p.parseInfixExpression)
	p.registerInfix(lexer.MUL, p.parseInfixExpression)
	p.registerInfix(lexer.DIV, p.parseInfixExpression)
	p.registerInfix(lexer.MOD, p.parseInfixExpression)
	p.registerInfix(lexer.POW, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NE, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.LE, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.GE, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.BIT_AND, p.parseInfixExpression)
	p.registerInfix(lexer.BIT_OR, p.parseBitOrOrPipeExpression)
	p.registerInfix(lexer.BIT_XOR, p.parseInfixExpression)
	p.registerInfix(lexer.SHL, p.parseInfixExpression)
	p.registerInfix(lexer.SHR, p.parseInfixExpression)
	p.registerInfix(lexer.IN, p.parseInfixExpression)
	p.registerInfix(lexer.MATCHES, p.parseInfixExpression)
	p.registerInfix(lexer.CONTAINS, p.parseInfixExpression)
	p.registerInfix(lexer.STARTS_WITH, p.parseInfixExpression)
	p.registerInfix(lexer.ENDS_WITH, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseMemberExpression)
	p.registerInfix(lexer.QUESTION, p.parseConditionalExpression)
	p.registerInfix(lexer.ARROW, p.parseLambdaExpression)
	p.registerInfix(lexer.QUESTION_DOT, p.parseOptionalChainingExpression)
	p.registerInfix(lexer.NULL_COALESCING, p.parseNullCoalescingExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// registerPrefix registers a prefix parse function
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram parses the program and returns the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.IMPORT:
		return p.parseImportStatement()
	case lexer.LBRACKET, lexer.LBRACE:
		// For now, skip destructuring assignment and just parse as expressions
		fallthrough
	default:
		// For now, we only support expression statements, import statements and destructuring assignments
		return p.parseExpressionStatement()
	}
}

// parseImportStatement parses an import statement (e.g., import "math" as m)
func (p *Parser) parseImportStatement() ast.Statement {
	pos := p.curToken.Position

	// Expect 'import' keyword
	if !p.expectToken(lexer.IMPORT) {
		return nil
	}

	// Expect module name as string literal
	if p.curToken.Type != lexer.STRING {
		p.errors = append(p.errors, fmt.Sprintf("expected string literal after 'import', got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}

	moduleName := p.curToken.Value
	p.nextToken()

	var alias string
	// Check if there's an 'as' clause
	if p.curToken.Type == lexer.AS {
		p.nextToken()
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("expected identifier after 'as', got %s at %s",
				p.curToken.Type, p.curToken.Position))
			return nil
		}
		alias = p.curToken.Value
		p.nextToken()
	} else {
		// If no alias is provided, use the module name as alias
		alias = moduleName
	}

	return &ast.ImportStatement{
		ModuleName: moduleName,
		Alias:      alias,
		Pos:        pos,
	}
}

// expectToken checks if current token matches expected type and advances
func (p *Parser) expectToken(expectedType lexer.TokenType) bool {
	if p.curToken.Type != expectedType {
		p.errors = append(p.errors, fmt.Sprintf("expected %s, got %s at %s",
			expectedType, p.curToken.Type, p.curToken.Position))
		return false
	}
	p.nextToken()
	return true
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Pos: p.curToken.Position}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// parseExpression parses an expression using Pratt parsing
func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != lexer.EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Value: p.curToken.Value,
		Pos:   p.curToken.Position,
	}
}

// parseNumberLiteral parses a number literal
func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.Literal{Pos: p.curToken.Position}

	// Try to parse as integer first
	if val, err := strconv.ParseInt(p.curToken.Value, 10, 64); err == nil {
		lit.Value = types.NewInt(val)
		return lit
	}

	// Parse as float
	val, err := strconv.ParseFloat(p.curToken.Value, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = types.NewFloat(val)
	return lit
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.Literal{
		Value: types.NewString(p.curToken.Value),
		Pos:   p.curToken.Position,
	}
}

// parseWildcard parses a wildcard expression
func (p *Parser) parseWildcard() ast.Expression {
	return &ast.WildcardExpression{
		TypeInfo: types.TypeInfo{Kind: types.KindInterface, Name: "wildcard"},
		Pos:      p.curToken.Position,
	}
}

// parsePlaceholder parses a placeholder expression
func (p *Parser) parsePlaceholder() ast.Expression {
	return &ast.PlaceholderExpression{
		TypeInfo: types.TypeInfo{Kind: types.KindInterface, Name: "placeholder"},
		Pos:      p.curToken.Position,
	}
}

// parseBooleanLiteral parses a boolean literal
func (p *Parser) parseBooleanLiteral() ast.Expression {
	value := p.curToken.Value == "true"
	return &ast.Literal{
		Value: types.NewBool(value),
		Pos:   p.curToken.Position,
	}
}

// parseNullLiteral parses a null literal
func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.Literal{
		Value: types.NewNil(),
		Pos:   p.curToken.Position,
	}
}

// parsePrefixExpression parses a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.curToken.Value,
		Pos:      p.curToken.Position,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses an infix expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: p.curToken.Value,
		Pos:      p.curToken.Position,
	}

	precedence := p.curPrecedence()

	// Handle right-associative operators
	if IsRightAssociative(p.curToken.Type) {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression parses a grouped expression (parentheses) or lambda parameters
func (p *Parser) parseGroupedExpression() ast.Expression {
	pos := p.curToken.Position
	p.nextToken()

	// Check if this could be lambda parameters: collect identifiers
	var identifiers []string
	firstExpr := p.parseExpression(LOWEST)

	// Check if it's a simple identifier - potential lambda parameter
	if ident, ok := firstExpr.(*ast.Identifier); ok && p.peekToken.Type == lexer.COMMA {
		// This looks like lambda parameters (x, y, ...)
		identifiers = append(identifiers, ident.Value)

		// Parse remaining parameters
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken() // consume comma
			p.nextToken() // move to next identifier

			if p.curToken.Type != lexer.IDENT {
				// Not a valid parameter list, fallback to regular parsing
				break
			}

			identifiers = append(identifiers, p.curToken.Value)
		}

		// If we successfully parsed parameter list and next is ), check for =>
		if p.expectPeek(lexer.RPAREN) && p.peekToken.Type == lexer.ARROW {
			// This is definitely lambda parameters, parse as lambda
			p.nextToken() // consume =>
			p.nextToken() // move to body

			body := p.parseExpression(LOWEST)

			return &ast.LambdaExpression{
				Parameters: identifiers,
				Body:       body,
				Pos:        pos,
			}
		}
	}

	// Not lambda parameters, parse as regular grouped expression
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return firstExpr
}

// parseArrayLiteral parses an array literal
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Pos: p.curToken.Position}
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}

// parseMapLiteral parses a map literal
func (p *Parser) parseMapLiteral() ast.Expression {
	hash := &ast.MapLiteral{Pos: p.curToken.Position}
	hash.Pairs = []ast.MapPair{}

	for p.peekToken.Type != lexer.RBRACE {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs = append(hash.Pairs, ast.MapPair{Key: key, Value: value})

		if p.peekToken.Type != lexer.RBRACE && !p.expectPeek(lexer.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return hash
}

// parseCallExpression parses a function call expression
func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	// Check if this is a function call on an identifier
	if ident, ok := fn.(*ast.Identifier); ok {
		// Parse as builtin function (compiler will decide if it's actually builtin)
		args := p.parseExpressionList(lexer.RPAREN)
		return &ast.BuiltinExpression{
			Name:      ident.Value,
			Arguments: args,
			Pos:       ident.Pos,
		}
	}

	// Parse as regular function call for complex expressions
	exp := &ast.CallExpression{Function: fn, Pos: p.curToken.Position}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

// isBuiltinFunction checks if a name is a builtin function
func (p *Parser) isBuiltinFunction(name string) bool {
	builtins := map[string]bool{
		// Type conversion functions
		"string": true,
		"int":    true,
		"float":  true,
		"bool":   true,

		// Math functions
		"abs": true,
		"max": true,
		"min": true,
		"sum": true,

		// String functions
		"len":        true,
		"contains":   true,
		"startsWith": true,
		"endsWith":   true,
		"matches":    true,
		"upper":      true,
		"lower":      true,
		"trim":       true,

		// Collection functions
		"all":    true,
		"any":    true,
		"filter": true,
		"map":    true,
		"count":  true,
		"first":  true,
		"last":   true,

		// Utility functions
		"type": true,
		"keys": true,
	}
	return builtins[name]
}

// isValidPropertyToken checks if a token type can be used as a property name
func (p *Parser) isValidPropertyToken(tokenType lexer.TokenType) bool {
	// Allow certain keywords to be used as property names in member access
	validProperties := map[lexer.TokenType]bool{
		lexer.CONTAINS:    true,
		lexer.STARTS_WITH: true,
		lexer.ENDS_WITH:   true,
		lexer.MATCHES:     true,
	}
	return validProperties[tokenType]
}

// parseIndexExpression parses an index expression
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Left: left, Pos: p.curToken.Position}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

// parseMemberExpression parses a member access expression
func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{Object: left, Pos: p.curToken.Position}

	// Check for enhanced array access syntax: user.[0].name
	if p.peekToken.Type == lexer.LBRACKET {
		// This is actually an index expression, not a member expression
		// Convert to index expression
		p.nextToken() // consume the '.'
		return p.parseIndexExpression(left)
	}

	// Check for wildcard member access: user.*
	if p.peekToken.Type == lexer.WILDCARD {
		p.nextToken() // consume the '.'
		exp.Property = &ast.WildcardExpression{
			TypeInfo: types.TypeInfo{Kind: types.KindInterface, Name: "wildcard"},
			Pos:      p.curToken.Position,
		}
		return exp
	}

	// Standard member access: user.name
	// Allow certain keywords as property names
	if p.peekToken.Type == lexer.IDENT {
		p.nextToken()
		exp.Property = &ast.Identifier{
			Value: p.curToken.Value,
			Pos:   p.curToken.Position,
		}
	} else if p.isValidPropertyToken(p.peekToken.Type) {
		// Allow certain keywords as property names (like contains, matches, etc.)
		p.nextToken()
		exp.Property = &ast.Identifier{
			Value: p.curToken.Value,
			Pos:   p.curToken.Position,
		}
	} else {
		p.peekError(lexer.IDENT)
		return nil
	}

	return exp
}

// parseConditionalExpression parses a ternary conditional expression
func (p *Parser) parseConditionalExpression(condition ast.Expression) ast.Expression {
	exp := &ast.ConditionalExpression{Test: condition, Pos: p.curToken.Position}

	p.nextToken()
	exp.Consequent = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	exp.Alternative = p.parseExpression(LOWEST)

	return exp
}

// parseExpressionList parses a list of expressions separated by commas
func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// parseBuiltinFunction parses a builtin function call
func (p *Parser) parseBuiltinFunction() ast.Expression {
	name := p.curToken.Value
	pos := p.curToken.Position

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	args := p.parseExpressionList(lexer.RPAREN)

	return &ast.BuiltinExpression{
		Name:      name,
		Arguments: args,
		Pos:       pos,
	}
}

// Helper methods

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() Precedence {
	return GetPrecedence(p.peekToken.Type)
}

// curPrecedence returns the precedence of the current token
func (p *Parser) curPrecedence() Precedence {
	return GetPrecedence(p.curToken.Type)
}

// expectPeek checks the peek token type and advances if it matches
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Error handling

// Errors returns the parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError adds a peek error
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// noPrefixParseFnError adds a no prefix parse function error
func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseLambdaExpression parses lambda expressions (e.g., x => x * 2)
func (p *Parser) parseLambdaExpression(left ast.Expression) ast.Expression {
	// The left side should be parameter(s)
	var parameters []string

	// Handle different parameter formats
	switch leftExpr := left.(type) {
	case *ast.Identifier:
		// Single parameter: x => ...
		parameters = []string{leftExpr.Value}
	case *ast.ArrayLiteral:
		// Multiple parameters: [x, y] => ...  (we'll parse this as array for now)
		// This is a simplified approach - in practice you might want special syntax
		for _, elem := range leftExpr.Elements {
			if ident, ok := elem.(*ast.Identifier); ok {
				parameters = append(parameters, ident.Value)
			} else {
				p.errors = append(p.errors, "lambda parameters must be identifiers")
				return nil
			}
		}
	default:
		p.errors = append(p.errors, "invalid lambda parameter format")
		return nil
	}

	expression := &ast.LambdaExpression{
		Parameters: parameters,
		Pos:        p.curToken.Position,
	}

	// Parse the body expression after =>
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Body = p.parseExpression(precedence)

	return expression
}

// parseBitOrOrPipeExpression parses BIT_OR as either bitwise OR or pipe operation based on context
func (p *Parser) parseBitOrOrPipeExpression(left ast.Expression) ast.Expression {
	// Determine if this should be a pipe operation or bitwise OR
	// Heuristic: if the right side looks like a function call or identifier that could be a pipeline function,
	// treat it as a pipe operation
	if p.isPipeOperation(left) {
		// Parse as pipe expression
		expression := &ast.PipeExpression{
			Left: left,
			Pos:  p.curToken.Position,
		}

		precedence := p.curPrecedence()
		p.nextToken()
		expression.Right = p.parseExpression(precedence)

		return expression
	} else {
		// Parse as regular bitwise OR
		return p.parseInfixExpression(left)
	}
}

// isPipeOperation determines if a BIT_OR should be treated as a pipe operation
func (p *Parser) isPipeOperation(left ast.Expression) bool {
	// Simple heuristic: if the next token is an identifier that looks like a function name,
	// or if the left side is a complex expression (not just a number), treat as pipe

	// Check if the next token is an identifier (potential function name)
	if p.peekToken.Type == lexer.IDENT {
		// Common pipeline function names
		funcName := p.peekToken.Value
		pipelineFuncs := map[string]bool{
			"filter": true, "map": true, "reduce": true, "sort": true, "reverse": true,
			"take": true, "skip": true, "join": true, "split": true, "match": true,
			"sum": true, "avg": true, "count": true, "len": true, "unique": true,
			"first": true, "last": true, "max": true, "min": true,
		}
		if pipelineFuncs[funcName] {
			return true
		}
	}

	// If left side is not a simple literal, it's more likely to be a pipeline
	switch left.(type) {
	case *ast.Literal:
		// Only numbers would be used in bitwise operations
		if lit, ok := left.(*ast.Literal); ok {
			if lit.Value != nil && lit.Value.Type().IsInteger() {
				return false // Likely bitwise operation
			}
		}
		return true // Other literals (strings, etc.) more likely to be piped
	case *ast.Identifier, *ast.MemberExpression, *ast.IndexExpression, *ast.CallExpression:
		return true // Complex expressions are more likely to be piped
	default:
		return false
	}
}

// parseOptionalChainingExpression parses optional chaining expressions (e.g., obj?.property)
func (p *Parser) parseOptionalChainingExpression(left ast.Expression) ast.Expression {
	expression := &ast.OptionalChainingExpression{
		Object: left,
		Pos:    p.curToken.Position,
	}

	if !p.expectPeek(lexer.IDENT) {
		// Try to parse as computed property access: obj?.[expr]
		if p.peekToken.Type == lexer.LBRACKET {
			p.nextToken() // consume [
			p.nextToken() // move to expression
			expression.Property = p.parseExpression(LOWEST)
			if !p.expectPeek(lexer.RBRACKET) {
				return nil
			}
		} else {
			return nil
		}
	} else {
		// Simple property access: obj?.property
		expression.Property = &ast.Identifier{
			Value: p.curToken.Value,
			Pos:   p.curToken.Position,
		}
	}

	return expression
}

// parseNullCoalescingExpression parses null coalescing expressions (e.g., a ?? b)
// isDestructuringAssignment checks if current position is a destructuring assignment
// For now, this is disabled to avoid conflicts with expression parsing
func (p *Parser) isDestructuringAssignment() bool {
	return false
}

// parseDestructuringAssignment parses a destructuring assignment statement
func (p *Parser) parseDestructuringAssignment() ast.Statement {
	pos := p.curToken.Position

	// Parse the left side (destructuring pattern)
	left := p.parseDestructuringPattern()
	if left == nil {
		return nil
	}

	// Expect assignment operator
	if p.peekToken.Type != lexer.ASSIGN {
		p.errors = append(p.errors, fmt.Sprintf("expected '=' in destructuring assignment, got %s at %s",
			p.peekToken.Type, p.peekToken.Position))
		return nil
	}
	p.nextToken() // consume '='
	p.nextToken() // move to right side

	// Parse the right side (value expression)
	right := p.parseExpression(LOWEST)
	if right == nil {
		return nil
	}

	return &ast.DestructuringAssignment{
		Left:  left,
		Right: right,
		Pos:   pos,
	}
}

// parseDestructuringPattern parses a destructuring pattern (array or object)
func (p *Parser) parseDestructuringPattern() ast.DestructuringPattern {
	switch p.curToken.Type {
	case lexer.LBRACKET:
		return p.parseArrayDestructuringPattern()
	case lexer.LBRACE:
		return p.parseObjectDestructuringPattern()
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected '[' or '{' for destructuring pattern, got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}
}

// parseArrayDestructuringPattern parses array destructuring pattern like [a, b, c]
func (p *Parser) parseArrayDestructuringPattern() ast.DestructuringPattern {
	pos := p.curToken.Position
	elements := []ast.DestructuringElement{}

	// Current token should already be LBRACKET
	if p.curToken.Type != lexer.LBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("expected '[' at start of array destructuring, got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}

	// Move to first element or ]
	p.nextToken()

	if p.curToken.Type == lexer.RBRACKET {
		// Empty array pattern []
		return &ast.ArrayDestructuringPattern{
			Elements: elements,
			Pos:      pos,
		}
	}

	for {
		if p.curToken.Type == lexer.RBRACKET {
			break
		}

		element := p.parseDestructuringElement()
		if element != nil {
			elements = append(elements, element)
		}

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken() // consume ','
		p.nextToken() // move to next element
	}

	// Expect closing bracket
	if p.peekToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("expected ']' at end of array destructuring, got %s at %s",
			p.peekToken.Type, p.peekToken.Position))
		return nil
	}
	p.nextToken() // consume ']'

	return &ast.ArrayDestructuringPattern{
		Elements: elements,
		Pos:      pos,
	}
}

// parseObjectDestructuringPattern parses object destructuring pattern like {name, age}
func (p *Parser) parseObjectDestructuringPattern() ast.DestructuringPattern {
	pos := p.curToken.Position
	properties := []ast.ObjectDestructuringProperty{}

	// Current token should already be LBRACE
	if p.curToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, fmt.Sprintf("expected '{' at start of object destructuring, got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}

	// Move to first property or }
	p.nextToken()

	if p.curToken.Type == lexer.RBRACE {
		// Empty object pattern {}
		return &ast.ObjectDestructuringPattern{
			Properties: properties,
			Pos:        pos,
		}
	}

	for {
		if p.curToken.Type == lexer.RBRACE {
			break
		}

		property := p.parseObjectDestructuringProperty()
		if property != nil {
			properties = append(properties, *property)
		}

		if p.peekToken.Type != lexer.COMMA {
			break
		}
		p.nextToken() // consume ','
		p.nextToken() // move to next property
	}

	// Expect closing brace
	if p.peekToken.Type != lexer.RBRACE {
		p.errors = append(p.errors, fmt.Sprintf("expected '}' at end of object destructuring, got %s at %s",
			p.peekToken.Type, p.peekToken.Position))
		return nil
	}
	p.nextToken() // consume '}'

	return &ast.ObjectDestructuringPattern{
		Properties: properties,
		Pos:        pos,
	}
}

// parseDestructuringElement parses a destructuring element (identifier or rest element)
func (p *Parser) parseDestructuringElement() ast.DestructuringElement {
	switch p.curToken.Type {
	case lexer.SPREAD:
		return p.parseRestElement()
	case lexer.IDENT:
		return p.parseIdentifierElement()
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected identifier or '...' in destructuring element, got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}
}

// parseIdentifierElement parses an identifier element in destructuring
func (p *Parser) parseIdentifierElement() ast.DestructuringElement {
	pos := p.curToken.Position
	name := p.curToken.Value

	var defaultValue ast.Expression
	// Check for default value
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // consume '='
		p.nextToken() // move to default value
		defaultValue = p.parseExpression(LOWEST)
	}

	return &ast.IdentifierElement{
		Name:    name,
		Default: defaultValue,
		Pos:     pos,
	}
}

// parseRestElement parses a rest element like ...rest
func (p *Parser) parseRestElement() ast.DestructuringElement {
	pos := p.curToken.Position

	if !p.expectPeek(lexer.SPREAD) {
		return nil
	}

	if p.peekToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier after '...', got %s at %s",
			p.peekToken.Type, p.peekToken.Position))
		return nil
	}

	p.nextToken() // move to identifier
	name := p.curToken.Value

	return &ast.RestElement{
		Name: name,
		Pos:  pos,
	}
}

// parseObjectDestructuringProperty parses an object destructuring property
func (p *Parser) parseObjectDestructuringProperty() *ast.ObjectDestructuringProperty {
	pos := p.curToken.Position

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier in object destructuring, got %s at %s",
			p.curToken.Type, p.curToken.Position))
		return nil
	}

	key := p.curToken.Value
	value := key // default to same as key

	// Check for key: value syntax
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // consume ':'
		if p.peekToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("expected identifier after ':', got %s at %s",
				p.peekToken.Type, p.peekToken.Position))
			return nil
		}
		p.nextToken() // move to value identifier
		value = p.curToken.Value
	}

	var defaultValue ast.Expression
	// Check for default value
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // consume '='
		p.nextToken() // move to default value
		defaultValue = p.parseExpression(LOWEST)
	}

	return &ast.ObjectDestructuringProperty{
		Key:     key,
		Value:   value,
		Default: defaultValue,
		Pos:     pos,
	}
}

func (p *Parser) parseNullCoalescingExpression(left ast.Expression) ast.Expression {
	expression := &ast.NullCoalescingExpression{
		Left: left,
		Pos:  p.curToken.Position,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
