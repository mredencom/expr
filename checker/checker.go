package checker

import (
	"fmt"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
)

// Checker performs static type checking on AST nodes
type Checker struct {
	scope  *Scope
	errors []string
}

// New creates a new type checker
func New() *Checker {
	return &Checker{
		scope:  NewRootScope(),
		errors: []string{},
	}
}

// NewWithScope creates a new type checker with a custom scope
func NewWithScope(scope *Scope) *Checker {
	return &Checker{
		scope:  scope,
		errors: []string{},
	}
}

// Check performs type checking on a program
func (c *Checker) Check(program *ast.Program) error {
	for _, stmt := range program.Statements {
		c.checkStatement(stmt)
	}

	if len(c.errors) > 0 {
		return fmt.Errorf("type checking failed: %v", c.errors)
	}

	return nil
}

// CheckExpression performs type checking on an expression and returns its type
func (c *Checker) CheckExpression(expr ast.Expression) (types.TypeInfo, error) {
	typeInfo := c.checkExpression(expr)

	if len(c.errors) > 0 {
		return types.TypeInfo{}, fmt.Errorf("type checking failed: %v", c.errors)
	}

	return typeInfo, nil
}

// Errors returns the type checking errors
func (c *Checker) Errors() []string {
	return c.errors
}

// Scope returns the current scope
func (c *Checker) Scope() *Scope {
	return c.scope
}

// WithEnvironment adds environment variables to the scope
func (c *Checker) WithEnvironment(env map[string]types.TypeInfo) *Checker {
	for name, typeInfo := range env {
		c.scope.DefineVariable(name, typeInfo)
	}
	return c
}

// checkStatement checks a statement
func (c *Checker) checkStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		c.checkExpression(s.Expression)
	default:
		c.addError(fmt.Sprintf("unknown statement type: %T", stmt))
	}
}

// checkExpression checks an expression and returns its type
func (c *Checker) checkExpression(expr ast.Expression) types.TypeInfo {
	if expr == nil {
		return types.TypeInfo{Kind: types.KindNil, Name: "nil"}
	}

	switch e := expr.(type) {
	case *ast.Literal:
		return c.checkLiteral(e)
	case *ast.Identifier:
		return c.checkIdentifier(e)
	case *ast.InfixExpression:
		return c.checkInfixExpression(e)
	case *ast.PrefixExpression:
		return c.checkPrefixExpression(e)
	case *ast.CallExpression:
		return c.checkCallExpression(e)
	case *ast.IndexExpression:
		return c.checkIndexExpression(e)
	case *ast.MemberExpression:
		return c.checkMemberExpression(e)
	case *ast.ConditionalExpression:
		return c.checkConditionalExpression(e)
	case *ast.ArrayLiteral:
		return c.checkArrayLiteral(e)
	case *ast.MapLiteral:
		return c.checkMapLiteral(e)
	case *ast.BuiltinExpression:
		return c.checkBuiltinExpression(e)
	case *ast.VariableExpression:
		return c.checkVariableExpression(e)
	default:
		c.addError(fmt.Sprintf("unknown expression type: %T", expr))
		return types.TypeInfo{Kind: types.KindNil, Name: "unknown"}
	}
}

// checkLiteral checks a literal expression
func (c *Checker) checkLiteral(lit *ast.Literal) types.TypeInfo {
	if lit.Value == nil {
		return types.NilType
	}
	return lit.Value.Type()
}

// checkIdentifier checks an identifier expression
func (c *Checker) checkIdentifier(ident *ast.Identifier) types.TypeInfo {
	if typeInfo, ok := c.scope.LookupVariable(ident.Value); ok {
		// Update the identifier's type info
		ident.TypeInfo = typeInfo
		return typeInfo
	}

	c.addError(fmt.Sprintf("undefined variable: %s", ident.Value))
	return types.TypeInfo{Kind: types.KindNil, Name: "undefined"}
}

// checkVariableExpression checks a variable expression
func (c *Checker) checkVariableExpression(varExpr *ast.VariableExpression) types.TypeInfo {
	if typeInfo, ok := c.scope.LookupVariable(varExpr.Name); ok {
		varExpr.TypeInfo = typeInfo
		return typeInfo
	}

	c.addError(fmt.Sprintf("undefined variable: %s", varExpr.Name))
	return types.TypeInfo{Kind: types.KindNil, Name: "undefined"}
}

// checkInfixExpression checks an infix expression
func (c *Checker) checkInfixExpression(infix *ast.InfixExpression) types.TypeInfo {
	leftType := c.checkExpression(infix.Left)
	rightType := c.checkExpression(infix.Right)

	resultType := c.inferInfixType(infix.Operator, leftType, rightType, infix.Pos)
	infix.TypeInfo = resultType
	return resultType
}

// checkPrefixExpression checks a prefix expression
func (c *Checker) checkPrefixExpression(prefix *ast.PrefixExpression) types.TypeInfo {
	rightType := c.checkExpression(prefix.Right)

	resultType := c.inferPrefixType(prefix.Operator, rightType, prefix.Pos)
	prefix.TypeInfo = resultType
	return resultType
}

// checkCallExpression checks a function call expression
func (c *Checker) checkCallExpression(call *ast.CallExpression) types.TypeInfo {
	// Check if it's a function identifier
	if ident, ok := call.Function.(*ast.Identifier); ok {
		if funcInfo, exists := c.scope.LookupFunction(ident.Value); exists {
			return c.checkFunctionCall(funcInfo, call.Arguments, call.Pos)
		}
		c.addError(fmt.Sprintf("undefined function: %s", ident.Value))
		return types.TypeInfo{Kind: types.KindNil, Name: "undefined"}
	}

	// For other function expressions, we need to check the function type
	funcType := c.checkExpression(call.Function)
	if funcType.Kind != types.KindFunc {
		c.addError(fmt.Sprintf("cannot call non-function type: %s", funcType.Name))
		return types.TypeInfo{Kind: types.KindNil, Name: "error"}
	}

	// For now, assume it returns interface{}
	return types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
}

// checkIndexExpression checks an index expression
func (c *Checker) checkIndexExpression(index *ast.IndexExpression) types.TypeInfo {
	leftType := c.checkExpression(index.Left)
	indexType := c.checkExpression(index.Index)

	switch leftType.Kind {
	case types.KindSlice, types.KindArray:
		if !indexType.IsInteger() {
			c.addError(fmt.Sprintf("array/slice index must be integer, got %s", indexType.Name))
		}
		if leftType.ElemType != nil {
			index.TypeInfo = *leftType.ElemType
			return *leftType.ElemType
		}
		return types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}

	case types.KindMap:
		if leftType.KeyType != nil && !types.CanConvert(indexType, *leftType.KeyType) {
			c.addError(fmt.Sprintf("map key type mismatch: expected %s, got %s",
				leftType.KeyType.Name, indexType.Name))
		}
		if leftType.ValType != nil {
			index.TypeInfo = *leftType.ValType
			return *leftType.ValType
		}
		return types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}

	case types.KindString:
		if !indexType.IsInteger() {
			c.addError(fmt.Sprintf("string index must be integer, got %s", indexType.Name))
		}
		result := types.StringType
		index.TypeInfo = result
		return result

	default:
		c.addError(fmt.Sprintf("cannot index type %s", leftType.Name))
		return types.TypeInfo{Kind: types.KindNil, Name: "error"}
	}
}

// checkMemberExpression checks a member expression
func (c *Checker) checkMemberExpression(member *ast.MemberExpression) types.TypeInfo {
	objectType := c.checkExpression(member.Object)

	// Handle wildcard property
	if wildcard, ok := member.Property.(*ast.WildcardExpression); ok {
		// Wildcard access - return interface{} type
		result := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
		wildcard.TypeInfo = result
		member.TypeInfo = result
		return result
	}

	// Handle identifier property
	if ident, ok := member.Property.(*ast.Identifier); ok {
		if objectType.Kind == types.KindStruct {
			// Look for the field in the struct
			for _, field := range objectType.Fields {
				if field.Name == ident.Value {
					member.TypeInfo = field.Type
					ident.TypeInfo = field.Type
					return field.Type
				}
			}
			c.addError(fmt.Sprintf("field %s not found in struct %s",
				ident.Value, objectType.Name))
		} else {
			c.addError(fmt.Sprintf("cannot access member %s of non-struct type %s",
				ident.Value, objectType.Name))
		}
	} else {
		// Handle other property types (should not happen in normal cases)
		c.addError("unsupported member property type")
	}

	return types.TypeInfo{Kind: types.KindNil, Name: "error"}
}

// checkConditionalExpression checks a conditional expression
func (c *Checker) checkConditionalExpression(cond *ast.ConditionalExpression) types.TypeInfo {
	testType := c.checkExpression(cond.Test)
	consequentType := c.checkExpression(cond.Consequent)
	alternativeType := c.checkExpression(cond.Alternative)

	// Test must be boolean
	if testType.Kind != types.KindBool {
		c.addError(fmt.Sprintf("conditional test must be boolean, got %s", testType.Name))
	}

	// Both branches should have compatible types
	if consequentType.Compatible(alternativeType) {
		cond.TypeInfo = consequentType
		return consequentType
	} else if alternativeType.Compatible(consequentType) {
		cond.TypeInfo = alternativeType
		return alternativeType
	} else {
		// Return interface{} for incompatible types
		result := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
		cond.TypeInfo = result
		return result
	}
}

// checkArrayLiteral checks an array literal
func (c *Checker) checkArrayLiteral(array *ast.ArrayLiteral) types.TypeInfo {
	if len(array.Elements) == 0 {
		// Empty array, assume interface{} elements
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
		result := types.TypeInfo{
			Kind:     types.KindSlice,
			Name:     "[]interface{}",
			ElemType: &elemType,
		}
		array.TypeInfo = result
		return result
	}

	// Check first element to determine array type
	firstType := c.checkExpression(array.Elements[0])

	// Check all elements are compatible
	for i, elem := range array.Elements[1:] {
		elemType := c.checkExpression(elem)
		if !firstType.Compatible(elemType) {
			c.addError(fmt.Sprintf("array element %d type mismatch: expected %s, got %s",
				i+1, firstType.Name, elemType.Name))
		}
	}

	result := types.TypeInfo{
		Kind:     types.KindSlice,
		Name:     "[]" + firstType.Name,
		ElemType: &firstType,
	}
	array.TypeInfo = result
	return result
}

// checkMapLiteral checks a map literal
func (c *Checker) checkMapLiteral(mapLit *ast.MapLiteral) types.TypeInfo {
	if len(mapLit.Pairs) == 0 {
		// Empty map, assume string keys and interface{} values
		keyType := types.StringType
		valType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
		result := types.TypeInfo{
			Kind:    types.KindMap,
			Name:    "map[string]interface{}",
			KeyType: &keyType,
			ValType: &valType,
		}
		mapLit.TypeInfo = result
		return result
	}

	// Check first pair to determine initial map types
	firstPair := mapLit.Pairs[0]
	keyType := c.checkExpression(firstPair.Key)
	valType := c.checkExpression(firstPair.Value)

	// Check if all values are compatible, if not use interface{}
	allValuesCompatible := true
	for _, pair := range mapLit.Pairs[1:] {
		pairKeyType := c.checkExpression(pair.Key)
		pairValType := c.checkExpression(pair.Value)

		if !keyType.Compatible(pairKeyType) {
			c.addError(fmt.Sprintf("map key type mismatch: expected %s, got %s",
				keyType.Name, pairKeyType.Name))
		}
		if !valType.Compatible(pairValType) {
			allValuesCompatible = false
		}
	}

	// If values are not all compatible, use interface{} as value type
	if !allValuesCompatible {
		valType = types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
	}

	result := types.TypeInfo{
		Kind:    types.KindMap,
		Name:    "map[" + keyType.Name + "]" + valType.Name,
		KeyType: &keyType,
		ValType: &valType,
	}
	mapLit.TypeInfo = result
	return result
}

// checkBuiltinExpression checks a builtin expression
func (c *Checker) checkBuiltinExpression(builtin *ast.BuiltinExpression) types.TypeInfo {
	if funcInfo, ok := c.scope.LookupFunction(builtin.Name); ok {
		result := c.checkFunctionCall(funcInfo, builtin.Arguments, builtin.Pos)
		builtin.TypeInfo = result
		return result
	}

	c.addError(fmt.Sprintf("undefined builtin function: %s", builtin.Name))
	return types.TypeInfo{Kind: types.KindNil, Name: "undefined"}
}

// Helper methods

// addError adds an error to the error list
func (c *Checker) addError(msg string) {
	c.errors = append(c.errors, msg)
}

// addErrorAt adds an error with position information
func (c *Checker) addErrorAt(pos lexer.Position, msg string) {
	c.errors = append(c.errors, fmt.Sprintf("%s: %s", pos.String(), msg))
}

// checkFunctionCall checks a function call
func (c *Checker) checkFunctionCall(funcInfo *FunctionInfo, args []ast.Expression, pos lexer.Position) types.TypeInfo {
	// Check argument count
	expectedArgs := len(funcInfo.Params)
	actualArgs := len(args)

	if funcInfo.Variadic {
		if actualArgs < expectedArgs-1 {
			c.addErrorAt(pos, fmt.Sprintf("function %s expects at least %d arguments, got %d",
				funcInfo.Name, expectedArgs-1, actualArgs))
		}
	} else {
		if actualArgs != expectedArgs {
			c.addErrorAt(pos, fmt.Sprintf("function %s expects %d arguments, got %d",
				funcInfo.Name, expectedArgs, actualArgs))
		}
	}

	// Check argument types
	for i, arg := range args {
		argType := c.checkExpression(arg)

		var expectedType types.TypeInfo
		if i < len(funcInfo.Params) {
			expectedType = funcInfo.Params[i]
		} else if funcInfo.Variadic {
			// Use the last parameter type for variadic arguments
			expectedType = funcInfo.Params[len(funcInfo.Params)-1]
		} else {
			continue // Too many arguments, already reported
		}

		if !expectedType.Assignable(argType) && expectedType.Kind != types.KindInterface {
			c.addErrorAt(pos, fmt.Sprintf("function %s argument %d type mismatch: expected %s, got %s",
				funcInfo.Name, i+1, expectedType.Name, argType.Name))
		}
	}

	// Return the function's return type
	if len(funcInfo.Returns) > 0 {
		return funcInfo.Returns[0]
	}

	return types.TypeInfo{Kind: types.KindNil, Name: "void"}
}

// inferInfixType infers the result type of an infix operation
func (c *Checker) inferInfixType(op string, left, right types.TypeInfo, pos lexer.Position) types.TypeInfo {
	switch op {
	case "==", "!=":
		if !left.IsComparable() || !right.IsComparable() {
			c.addErrorAt(pos, fmt.Sprintf("cannot compare types %s and %s", left.Name, right.Name))
		}
		return types.BoolType

	case "<", "<=", ">", ">=":
		if !left.IsOrdered() || !right.IsOrdered() {
			c.addErrorAt(pos, fmt.Sprintf("cannot compare types %s and %s", left.Name, right.Name))
		}
		if !left.Compatible(right) {
			c.addErrorAt(pos, fmt.Sprintf("cannot compare incompatible types %s and %s", left.Name, right.Name))
		}
		return types.BoolType

	case "&&", "||":
		if left.Kind != types.KindBool {
			c.addErrorAt(pos, fmt.Sprintf("left operand of %s must be boolean, got %s", op, left.Name))
		}
		if right.Kind != types.KindBool {
			c.addErrorAt(pos, fmt.Sprintf("right operand of %s must be boolean, got %s", op, right.Name))
		}
		return types.BoolType

	case "+":
		if left.Kind == types.KindString && right.Kind == types.KindString {
			return types.StringType
		}
		if left.Kind == types.KindString || right.Kind == types.KindString {
			c.addErrorAt(pos, fmt.Sprintf("invalid operation: cannot mix string and non-string types in addition, got %s and %s",
				left.Name, right.Name))
			return types.TypeInfo{Kind: types.KindNil, Name: "error"}
		}
		fallthrough
	case "-", "*", "/", "%", "**":
		if !left.IsNumeric() || !right.IsNumeric() {
			c.addErrorAt(pos, fmt.Sprintf("operator %s requires numeric operands, got %s and %s",
				op, left.Name, right.Name))
			return types.TypeInfo{Kind: types.KindNil, Name: "error"}
		}

		// If either operand is float, result is float
		if left.IsFloat() || right.IsFloat() {
			return types.FloatType
		}
		return types.IntType

	case "&", "|", "^", "<<", ">>":
		if !left.IsInteger() || !right.IsInteger() {
			c.addErrorAt(pos, fmt.Sprintf("bitwise operator %s requires integer operands, got %s and %s",
				op, left.Name, right.Name))
			return types.TypeInfo{Kind: types.KindNil, Name: "error"}
		}
		return types.IntType

	case "in":
		// Check if right is a collection type
		switch right.Kind {
		case types.KindSlice, types.KindArray:
			if right.ElemType != nil && !left.Compatible(*right.ElemType) {
				c.addErrorAt(pos, fmt.Sprintf("cannot check if %s is in %s", left.Name, right.Name))
			}
		case types.KindMap:
			if right.KeyType != nil && !left.Compatible(*right.KeyType) {
				c.addErrorAt(pos, fmt.Sprintf("cannot check if %s is in map with key type %s",
					left.Name, right.KeyType.Name))
			}
		case types.KindString:
			if left.Kind != types.KindString {
				c.addErrorAt(pos, fmt.Sprintf("cannot check if %s is in string", left.Name))
			}
		default:
			c.addErrorAt(pos, fmt.Sprintf("cannot use 'in' with type %s", right.Name))
		}
		return types.BoolType

	case "matches", "contains", "startsWith", "endsWith":
		if left.Kind != types.KindString || right.Kind != types.KindString {
			c.addErrorAt(pos, fmt.Sprintf("string operator %s requires string operands, got %s and %s",
				op, left.Name, right.Name))
		}
		return types.BoolType

	default:
		c.addErrorAt(pos, fmt.Sprintf("unknown operator: %s", op))
		return types.TypeInfo{Kind: types.KindNil, Name: "error"}
	}
}

// inferPrefixType infers the result type of a prefix operation
func (c *Checker) inferPrefixType(op string, right types.TypeInfo, pos lexer.Position) types.TypeInfo {
	switch op {
	case "!":
		if right.Kind != types.KindBool {
			c.addErrorAt(pos, fmt.Sprintf("logical NOT requires boolean operand, got %s", right.Name))
		}
		return types.BoolType

	case "-":
		if !right.IsNumeric() {
			c.addErrorAt(pos, fmt.Sprintf("unary minus requires numeric operand, got %s", right.Name))
			return types.TypeInfo{Kind: types.KindNil, Name: "error"}
		}
		return right

	case "~":
		if !right.IsInteger() {
			c.addErrorAt(pos, fmt.Sprintf("bitwise NOT requires integer operand, got %s", right.Name))
			return types.TypeInfo{Kind: types.KindNil, Name: "error"}
		}
		return right

	default:
		c.addErrorAt(pos, fmt.Sprintf("unknown prefix operator: %s", op))
		return types.TypeInfo{Kind: types.KindNil, Name: "error"}
	}
}
