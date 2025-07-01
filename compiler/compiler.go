package compiler

import (
	"fmt"
	"math"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/builtins"
	"github.com/mredencom/expr/checker"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// Bytecode represents compiled bytecode
type Bytecode struct {
	Instructions []byte
	Constants    []types.Value
	Variables    map[string]int // Variable name to index mapping
	Fields       map[string]int // Field name to index mapping
	Builtins     map[string]int // Builtin function name to index mapping
}

// EmittedInstruction represents an emitted instruction with position
type EmittedInstruction struct {
	Opcode   vm.Opcode
	Position int
}

// CompilationScope represents a compilation scope
type CompilationScope struct {
	instructions        []byte
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// Compiler represents the bytecode compiler
type Compiler struct {
	constants   []types.Value
	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int

	// Jump tracking for control flow
	jumpStack []int

	// Error tracking
	errors []string

	// Pipeline context for placeholder expressions
	inPipelineContext bool

	// Bytecode optimizer
	optimizer *BytecodeOptimizer
}

// New creates a new compiler
func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        []byte{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	// Register all builtin functions using standard names
	for i, name := range builtins.StandardBuiltinNames {
		symbolTable.DefineBuiltin(i, name)
	}

	return &Compiler{
		constants:         []types.Value{},
		symbolTable:       symbolTable,
		scopes:            []CompilationScope{mainScope},
		scopeIndex:        0,
		jumpStack:         []int{},
		errors:            []string{},
		optimizer:         NewBytecodeOptimizer(OptimizationBasic),
		inPipelineContext: false,
	}
}

// NewWithState creates a new compiler with existing state
func NewWithState(s *SymbolTable, constants []types.Value) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

// Compile compiles an AST node to bytecode
func (c *Compiler) Compile(node ast.Node) error {

	switch node := node.(type) {
	case *ast.Program:
		return c.compileProgram(node)

	case *ast.ExpressionStatement:
		return c.compileExpressionStatement(node)

	case *ast.Literal:
		return c.compileLiteral(node)

	case *ast.Identifier:
		return c.compileIdentifier(node)

	case *ast.InfixExpression:
		return c.compileInfixExpression(node)

	case *ast.PrefixExpression:
		return c.compilePrefixExpression(node)

	case *ast.CallExpression:
		return c.compileCallExpression(node)

	case *ast.BuiltinExpression:
		return c.compileBuiltinExpression(node)

	case *ast.IndexExpression:
		return c.compileIndexExpression(node)

	case *ast.MemberExpression:
		return c.compileMemberExpression(node)

	case *ast.ConditionalExpression:
		return c.compileConditionalExpression(node)

	case *ast.ArrayLiteral:
		return c.compileArrayLiteral(node)

	case *ast.MapLiteral:
		return c.compileMapLiteral(node)

	case *ast.LambdaExpression:
		return c.compileLambdaExpression(node)

	case *ast.PipeExpression:
		return c.compilePipeExpression(node)

	case *ast.PlaceholderExpression:
		return c.compilePlaceholderExpression(node)

	case *ast.OptionalChainingExpression:
		return c.compileOptionalChainingExpression(node)

	case *ast.NullCoalescingExpression:
		return c.compileNullCoalescingExpression(node)

	case *ast.ImportStatement:
		return c.compileImportStatement(node)

	case *ast.ModuleCallExpression:
		return c.compileModuleCallExpression(node)

	case *ast.DestructuringAssignment:
		return c.compileDestructuringAssignment(node)

	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

// compileProgram compiles a program node
func (c *Compiler) compileProgram(node *ast.Program) error {
	for _, stmt := range node.Statements {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// compileExpressionStatement compiles an expression statement
func (c *Compiler) compileExpressionStatement(node *ast.ExpressionStatement) error {
	err := c.Compile(node.Expression)
	if err != nil {
		return err
	}
	return nil
}

// compileLiteral compiles a literal expression
func (c *Compiler) compileLiteral(node *ast.Literal) error {
	if node.Value == nil {
		return c.emitError(vm.OpConstant, c.addConstant(types.NewNil()))
	}

	return c.emitError(vm.OpConstant, c.addConstant(node.Value))
}

// compileIdentifier compiles an identifier expression
func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	symbol, ok := c.symbolTable.Resolve(node.Value)
	if !ok {
		return fmt.Errorf("undefined variable %s", node.Value)
	}

	return c.loadSymbol(symbol)
}

// compileInfixExpression compiles an infix expression
func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
	// Check if we're in a pipeline context and this expression contains placeholders
	if c.inPipelineContext && c.hasPlaceholder(node) {
		// Don't compile placeholder expressions immediately in pipeline context
		// Instead, emit them as serialized expressions for later evaluation
		return c.compilePlaceholderInfixExpression(node)
	}

	// Try constant folding first - if successful, emit only the result
	if foldedValue := c.tryConstantFolding(node); foldedValue != nil {
		// Apply constant folding for all supported operations
		return c.emitError(vm.OpConstant, c.addConstant(foldedValue))
	}

	if node.Operator == "<" {
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		err = c.Compile(node.Left)
		if err != nil {
			return err
		}

		return c.emitOptimizedOp(node, vm.OpGreaterThan, vm.OpGreaterThan, vm.OpGreaterThan, vm.OpGreaterThan)
	}

	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "+":
		return c.emitOptimizedOp(node, vm.OpAdd, vm.OpAddInt64, vm.OpAddFloat64, vm.OpAddString)
	case "-":
		return c.emitOptimizedOp(node, vm.OpSub, vm.OpSubInt64, vm.OpSubFloat64, vm.OpSub)
	case "*":
		return c.emitOptimizedOp(node, vm.OpMul, vm.OpMulInt64, vm.OpMulFloat64, vm.OpMul)
	case "/":
		return c.emitOptimizedOp(node, vm.OpDiv, vm.OpDivInt64, vm.OpDivFloat64, vm.OpDiv)
	case "%":
		return c.emitOptimizedOp(node, vm.OpMod, vm.OpModInt64, vm.OpModFloat64, vm.OpMod)
	case "==":
		return c.emitError(vm.OpEqual)
	case "!=":
		return c.emitError(vm.OpNotEqual)
	case ">":
		return c.emitOptimizedOp(node, vm.OpGreaterThan, vm.OpGreaterThan, vm.OpGreaterThan, vm.OpGreaterThan)
	case ">=":
		return c.emitError(vm.OpGreaterEqual)
	case "<=":
		return c.emitError(vm.OpLessEqual)
	case "&&":
		return c.emitError(vm.OpAnd)
	case "||":
		return c.emitError(vm.OpOr)
	case "&":
		return c.emitError(vm.OpBitAnd)
	case "|":
		return c.emitError(vm.OpBitOr)
	case "^":
		return c.emitError(vm.OpBitXor)
	case "<<":
		return c.emitError(vm.OpShiftL)
	case ">>":
		return c.emitError(vm.OpShiftR)
	default:
		return fmt.Errorf("unknown operator %s", node.Operator)
	}
}

// tryConstantFolding attempts to fold constant expressions at compile time
func (c *Compiler) tryConstantFolding(node *ast.InfixExpression) types.Value {
	// Only fold if both operands are literals
	leftLit, leftIsLit := node.Left.(*ast.Literal)
	rightLit, rightIsLit := node.Right.(*ast.Literal)

	if !leftIsLit || !rightIsLit {
		return nil
	}

	left := leftLit.Value
	right := rightLit.Value

	if left == nil || right == nil {
		return nil
	}

	// Enhanced constant folding with more operations
	switch node.Operator {
	case "+":
		return c.foldAddition(left, right)
	case "-":
		return c.foldSubtraction(left, right)
	case "*":
		return c.foldMultiplication(left, right)
	case "/":
		return c.foldDivision(left, right)
	case "%":
		return c.foldModulo(left, right)
	case "==":
		return c.foldComparison(left, right, "==")
	case "!=":
		return c.foldComparison(left, right, "!=")
	case "<":
		return c.foldComparison(left, right, "<")
	case "<=":
		return c.foldComparison(left, right, "<=")
	case ">":
		return c.foldComparison(left, right, ">")
	case ">=":
		return c.foldComparison(left, right, ">=")
	case "&&":
		return c.foldLogical(left, right, "&&")
	case "||":
		return c.foldLogical(left, right, "||")
	case "&":
		return c.foldBitwise(left, right, "&")
	case "|":
		return c.foldBitwise(left, right, "|")
	case "^":
		return c.foldBitwise(left, right, "^")
	case "<<":
		return c.foldBitwise(left, right, "<<")
	case ">>":
		return c.foldBitwise(left, right, ">>")
	case "**":
		return c.foldPower(left, right)
	}

	return nil
}

// foldPower performs constant folding for power operations
func (c *Compiler) foldPower(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok {
			if r.Value() < 0 {
				return nil // Negative powers result in floats
			}
			result := int64(1)
			base := l.Value()
			exp := r.Value()

			// Simple power calculation with overflow protection
			for i := int64(0); i < exp && i < 63; i++ { // Limit to prevent overflow
				if result > 9223372036854775807/base { // Check for overflow
					return nil
				}
				result *= base
			}
			return types.NewInt(result)
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok {
			// Use math.Pow for float calculations
			result := math.Pow(l.Value(), r.Value())
			if math.IsInf(result, 0) || math.IsNaN(result) {
				return nil
			}
			return types.NewFloat(result)
		}
	}
	return nil
}

// foldComparison performs compile-time comparison operations
func (c *Compiler) foldComparison(left, right types.Value, op string) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok {
			lVal, rVal := l.Value(), r.Value()
			switch op {
			case "==":
				return types.NewBool(lVal == rVal)
			case "!=":
				return types.NewBool(lVal != rVal)
			case ">":
				return types.NewBool(lVal > rVal)
			case ">=":
				return types.NewBool(lVal >= rVal)
			case "<":
				return types.NewBool(lVal < rVal)
			case "<=":
				return types.NewBool(lVal <= rVal)
			}
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok {
			lVal, rVal := l.Value(), r.Value()
			switch op {
			case "==":
				return types.NewBool(lVal == rVal)
			case "!=":
				return types.NewBool(lVal != rVal)
			case ">":
				return types.NewBool(lVal > rVal)
			case ">=":
				return types.NewBool(lVal >= rVal)
			case "<":
				return types.NewBool(lVal < rVal)
			case "<=":
				return types.NewBool(lVal <= rVal)
			}
		}
	case *types.StringValue:
		if r, ok := right.(*types.StringValue); ok {
			lVal, rVal := l.Value(), r.Value()
			switch op {
			case "==":
				return types.NewBool(lVal == rVal)
			case "!=":
				return types.NewBool(lVal != rVal)
			case ">":
				return types.NewBool(lVal > rVal)
			case ">=":
				return types.NewBool(lVal >= rVal)
			case "<":
				return types.NewBool(lVal < rVal)
			case "<=":
				return types.NewBool(lVal <= rVal)
			}
		}
	case *types.BoolValue:
		if r, ok := right.(*types.BoolValue); ok {
			lVal, rVal := l.Value(), r.Value()
			switch op {
			case "==":
				return types.NewBool(lVal == rVal)
			case "!=":
				return types.NewBool(lVal != rVal)
			}
		}
	}
	return nil
}

// foldLogical performs compile-time logical operations
func (c *Compiler) foldLogical(left, right types.Value, op string) types.Value {
	// Convert values to boolean for logical operations
	leftBool := c.valueToBool(left)
	rightBool := c.valueToBool(right)

	if leftBool == nil || rightBool == nil {
		return nil
	}

	lVal := leftBool.(*types.BoolValue).Value()
	rVal := rightBool.(*types.BoolValue).Value()

	switch op {
	case "&&":
		return types.NewBool(lVal && rVal)
	case "||":
		return types.NewBool(lVal || rVal)
	}
	return nil
}

// foldBitwise performs compile-time bitwise operations
func (c *Compiler) foldBitwise(left, right types.Value, op string) types.Value {
	// Bitwise operations only work on integers
	leftInt, leftOk := left.(*types.IntValue)
	rightInt, rightOk := right.(*types.IntValue)

	if !leftOk || !rightOk {
		return nil
	}

	lVal := leftInt.Value()
	rVal := rightInt.Value()

	switch op {
	case "&":
		return types.NewInt(lVal & rVal)
	case "|":
		return types.NewInt(lVal | rVal)
	case "^":
		return types.NewInt(lVal ^ rVal)
	case "<<":
		return types.NewInt(lVal << uint(rVal))
	case ">>":
		return types.NewInt(lVal >> uint(rVal))
	}
	return nil
}

// valueToBool converts a value to boolean following expression language rules
func (c *Compiler) valueToBool(val types.Value) types.Value {
	switch v := val.(type) {
	case *types.BoolValue:
		return v
	case *types.IntValue:
		return types.NewBool(v.Value() != 0)
	case *types.FloatValue:
		return types.NewBool(v.Value() != 0.0)
	case *types.StringValue:
		return types.NewBool(v.Value() != "")
	case *types.NilValue:
		return types.NewBool(false)
	}
	return nil
}

// Enhanced foldAddition with mixed type conversion
func (c *Compiler) foldAddition(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok {
			return types.NewInt(l.Value() + r.Value())
		}
		// Int + Float = Float
		if r, ok := right.(*types.FloatValue); ok {
			return types.NewFloat(float64(l.Value()) + r.Value())
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok {
			return types.NewFloat(l.Value() + r.Value())
		}
		// Float + Int = Float
		if r, ok := right.(*types.IntValue); ok {
			return types.NewFloat(l.Value() + float64(r.Value()))
		}
	case *types.StringValue:
		if r, ok := right.(*types.StringValue); ok {
			return types.NewString(l.Value() + r.Value())
		}
	}
	return nil
}

// foldSubtraction performs compile-time subtraction
func (c *Compiler) foldSubtraction(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok {
			return types.NewInt(l.Value() - r.Value())
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok {
			return types.NewFloat(l.Value() - r.Value())
		}
	}
	return nil
}

// foldMultiplication performs compile-time multiplication
func (c *Compiler) foldMultiplication(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok {
			return types.NewInt(l.Value() * r.Value())
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok {
			return types.NewFloat(l.Value() * r.Value())
		}
	}
	return nil
}

// foldDivision performs compile-time division
func (c *Compiler) foldDivision(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok && r.Value() != 0 {
			return types.NewInt(l.Value() / r.Value())
		}
	case *types.FloatValue:
		if r, ok := right.(*types.FloatValue); ok && r.Value() != 0 {
			return types.NewFloat(l.Value() / r.Value())
		}
	}
	return nil
}

// foldModulo performs compile-time modulo
func (c *Compiler) foldModulo(left, right types.Value) types.Value {
	switch l := left.(type) {
	case *types.IntValue:
		if r, ok := right.(*types.IntValue); ok && r.Value() != 0 {
			return types.NewInt(l.Value() % r.Value())
		}
	}
	return nil
}

// emitOptimizedOp emits optimized operations based on operand types
func (c *Compiler) emitOptimizedOp(node *ast.InfixExpression, genericOp, intOp, floatOp, stringOp vm.Opcode) error {
	// Try to detect types of operands at compile time
	leftType := c.detectOperandType(node.Left)
	rightType := c.detectOperandType(node.Right)

	// If both operands are of the same basic type, use specialized instruction
	if leftType == rightType {
		switch leftType {
		case types.KindInt64:
			if intOp != 0 {
				return c.emitError(intOp)
			}
		case types.KindFloat64:
			if floatOp != 0 {
				return c.emitError(floatOp)
			}
		case types.KindString:
			if stringOp != 0 && node.Operator == "+" {
				return c.emitError(stringOp)
			}
		}
	}

	// Enhanced type inference for variables and member access
	if leftType == types.KindUnknown {
		leftType = c.inferTypeFromContext(node.Left)
	}
	if rightType == types.KindUnknown {
		rightType = c.inferTypeFromContext(node.Right)
	}

	// Try again with inferred types
	if leftType == rightType && leftType != types.KindUnknown {
		switch leftType {
		case types.KindInt64:
			if intOp != 0 {
				return c.emitError(intOp)
			}
		case types.KindFloat64:
			if floatOp != 0 {
				return c.emitError(floatOp)
			}
		case types.KindString:
			if stringOp != 0 && node.Operator == "+" {
				return c.emitError(stringOp)
			}
		}
	}

	// Fall back to generic operation
	return c.emitError(genericOp)
}

// detectOperandType detects the type of an operand at compile time
func (c *Compiler) detectOperandType(node ast.Node) types.TypeKind {
	switch n := node.(type) {
	case *ast.Literal:
		if n.Value == nil {
			return types.KindNil
		}
		return n.Value.Type().Kind

	case *ast.Identifier:
		// Try to infer from symbol table or context
		return c.inferIdentifierType(n.Value)

	case *ast.MemberExpression:
		// Try to infer member type
		return c.inferMemberType(n)

	case *ast.InfixExpression:
		// Infer result type based on operation and operands
		return c.inferInfixResultType(n)

	case *ast.CallExpression:
		// Infer result type from function call
		return c.inferCallResultType(n)

	default:
		return types.KindUnknown
	}
}

// inferTypeFromContext tries to infer type from surrounding context
func (c *Compiler) inferTypeFromContext(node ast.Node) types.TypeKind {
	// This is a simplified implementation
	// In a full implementation, this would use type checker information
	return types.KindUnknown
}

// inferIdentifierType tries to infer the type of an identifier
func (c *Compiler) inferIdentifierType(name string) types.TypeKind {
	// This would ideally use type information from the type checker
	// For now, return unknown
	return types.KindUnknown
}

// inferMemberType tries to infer the type of a member access
func (c *Compiler) inferMemberType(node *ast.MemberExpression) types.TypeKind {
	// This would use type information about the object and property
	return types.KindUnknown
}

// inferInfixResultType infers the result type of an infix expression
func (c *Compiler) inferInfixResultType(node *ast.InfixExpression) types.TypeKind {
	leftType := c.detectOperandType(node.Left)
	rightType := c.detectOperandType(node.Right)

	switch node.Operator {
	case "+":
		// String concatenation or numeric addition
		if leftType == types.KindString || rightType == types.KindString {
			return types.KindString
		}
		if leftType == types.KindFloat64 || rightType == types.KindFloat64 {
			return types.KindFloat64
		}
		if leftType == types.KindInt64 && rightType == types.KindInt64 {
			return types.KindInt64
		}
	case "-", "*", "/", "%":
		// Numeric operations
		if leftType == types.KindFloat64 || rightType == types.KindFloat64 {
			return types.KindFloat64
		}
		if leftType == types.KindInt64 && rightType == types.KindInt64 {
			return types.KindInt64
		}
	case "==", "!=", "<", "<=", ">", ">=", "&&", "||":
		// Comparison and logical operations return boolean
		return types.KindBool
	}

	return types.KindUnknown
}

// inferCallResultType infers the result type of a function call
func (c *Compiler) inferCallResultType(node *ast.CallExpression) types.TypeKind {
	// This would use builtin function type information
	return types.KindUnknown
}

// compilePrefixExpression compiles a prefix expression
func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	// Check if this is a placeholder expression in pipeline context
	if c.inPipelineContext && c.hasPlaceholder(node.Right) {
		// This is a prefix expression with placeholder (like !#) in pipeline context
		// Create a PlaceholderExprValue to handle this at runtime
		placeholderExpr := types.NewPlaceholderExpr(
			[]byte{},        // Empty bytecode for now
			[]types.Value{}, // Empty constants
			node.Operator,   // The prefix operator (!, -, ~)
			types.NewNil(),  // Unary operations don't have a right operand
		)

		// Emit the actual placeholder expression directly
		return c.emitError(vm.OpConstant, c.addConstant(placeholderExpr))
	}

	// Skip constant folding for compatibility with existing tests
	// Constant folding can be enabled later for production use

	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "!":
		return c.emitError(vm.OpNot)
	case "-":
		return c.emitError(vm.OpNeg)
	case "~":
		return c.emitError(vm.OpBitNot)
	default:
		return fmt.Errorf("unknown operator %s", node.Operator)
	}
}

// tryPrefixConstantFolding attempts to fold prefix constant expressions
func (c *Compiler) tryPrefixConstantFolding(node *ast.PrefixExpression) types.Value {
	// Check if operand is a literal
	lit, ok := node.Right.(*ast.Literal)
	if !ok || lit.Value == nil {
		return nil
	}

	switch node.Operator {
	case "-":
		return c.foldNegation(lit.Value)
	case "!":
		return c.foldLogicalNot(lit.Value)
	case "~":
		return c.foldBitwiseNot(lit.Value)
	}

	return nil
}

// foldNegation performs compile-time negation
func (c *Compiler) foldNegation(val types.Value) types.Value {
	switch v := val.(type) {
	case *types.IntValue:
		return types.NewInt(-v.Value())
	case *types.FloatValue:
		return types.NewFloat(-v.Value())
	}
	return nil
}

// foldLogicalNot performs compile-time logical NOT
func (c *Compiler) foldLogicalNot(val types.Value) types.Value {
	boolVal := c.valueToBool(val)
	if boolVal == nil {
		return nil
	}

	return types.NewBool(!boolVal.(*types.BoolValue).Value())
}

// foldBitwiseNot performs compile-time bitwise NOT
func (c *Compiler) foldBitwiseNot(val types.Value) types.Value {
	// Bitwise NOT only works on integers
	intVal, ok := val.(*types.IntValue)
	if !ok {
		return nil
	}

	return types.NewInt(^intVal.Value())
}

// compileCallExpression compiles a function call expression
func (c *Compiler) compileCallExpression(node *ast.CallExpression) error {
	err := c.Compile(node.Function)
	if err != nil {
		return err
	}

	for _, arg := range node.Arguments {
		err := c.Compile(arg)
		if err != nil {
			return err
		}
	}

	return c.emitError(vm.OpCall, len(node.Arguments))
}

// compileBuiltinExpression compiles a builtin expression
func (c *Compiler) compileBuiltinExpression(node *ast.BuiltinExpression) error {
	// Check if any argument contains a placeholder - if so, treat this as a pipeline function
	hasPlaceholder := c.containsPlaceholder(node.Arguments)

	if hasPlaceholder {
		// This is a pipeline function with placeholders, use the pipeline compilation logic
		return c.compilePipelineFunction(node.Name, node.Arguments)
	}

	// Regular builtin function compilation (no placeholders)
	for _, arg := range node.Arguments {
		err := c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// Find the index in StandardBuiltinNames
	builtinIndex := -1
	for i, name := range builtins.StandardBuiltinNames {
		if name == node.Name {
			builtinIndex = i
			break
		}
	}

	if builtinIndex == -1 {
		return fmt.Errorf("undefined builtin function %s", node.Name)
	}

	return c.emitError(vm.OpBuiltin, builtinIndex, len(node.Arguments))
}

// compileIndexExpression compiles an index expression
func (c *Compiler) compileIndexExpression(node *ast.IndexExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	err = c.Compile(node.Index)
	if err != nil {
		return err
	}

	return c.emitError(vm.OpIndex)
}

// compileMemberExpression compiles a member expression
func (c *Compiler) compileMemberExpression(node *ast.MemberExpression) error {
	// Check if we're in pipeline context and this involves a placeholder
	if c.inPipelineContext && c.hasPlaceholder(node) {
		// For placeholder member access like #.name, serialize it for runtime evaluation
		return c.compilePlaceholderMemberExpression(node)
	}

	// Compile the object
	err := c.Compile(node.Object)
	if err != nil {
		return err
	}

	// Handle different property types
	switch prop := node.Property.(type) {
	case *ast.Identifier:
		// Standard member access: obj.property
		propertyObj := types.NewString(prop.Value)
		err = c.emitError(vm.OpConstant, c.addConstant(propertyObj))
		if err != nil {
			return err
		}
		// Emit member access operation
		return c.emitError(vm.OpMember)
	case *ast.WildcardExpression:
		// Wildcard access: obj.*
		wildcardObj := types.NewString("*")
		err = c.emitError(vm.OpConstant, c.addConstant(wildcardObj))
		if err != nil {
			return err
		}
		// Emit member access operation for wildcard
		return c.emitError(vm.OpMember)
	default:
		// Other property types - compile the expression
		err = c.Compile(node.Property)
		if err != nil {
			return err
		}
		// Emit member access operation
		return c.emitError(vm.OpMember)
	}
}

// compilePlaceholderMemberExpression compiles a member expression containing placeholders
// in pipeline context by serializing it for later evaluation
func (c *Compiler) compilePlaceholderMemberExpression(node *ast.MemberExpression) error {
	// Emit a special marker to indicate this is a pipeline member access
	memberAccessMarker := types.NewString("__PIPELINE_MEMBER_ACCESS__")
	err := c.emitError(vm.OpConstant, c.addConstant(memberAccessMarker))
	if err != nil {
		return err
	}

	// Compile the object (usually a placeholder)
	err = c.Compile(node.Object)
	if err != nil {
		return err
	}

	// Handle the property
	switch prop := node.Property.(type) {
	case *ast.Identifier:
		// Standard member access: #.property
		propertyObj := types.NewString(prop.Value)
		err = c.emitError(vm.OpConstant, c.addConstant(propertyObj))
		if err != nil {
			return err
		}
	case *ast.WildcardExpression:
		// Wildcard access: #.*
		wildcardObj := types.NewString("*")
		err = c.emitError(vm.OpConstant, c.addConstant(wildcardObj))
		if err != nil {
			return err
		}
	default:
		// Other property types
		err = c.Compile(node.Property)
		if err != nil {
			return err
		}
	}

	// Create an array: ["__PIPELINE_MEMBER_ACCESS__", object, property]
	return c.emitError(vm.OpSlice, 3)
}

// compileConditionalExpression compiles a conditional (ternary) expression
func (c *Compiler) compileConditionalExpression(node *ast.ConditionalExpression) error {
	// Compile condition
	err := c.Compile(node.Test)
	if err != nil {
		return err
	}

	// Jump to false branch if condition is false (use 0 as placeholder)
	jumpNotTruthy := c.emit(vm.OpJumpFalse, 0)

	// Compile true branch
	err = c.Compile(node.Consequent)
	if err != nil {
		return err
	}

	// Jump over false branch (use 0 as placeholder)
	jumpPos := c.emit(vm.OpJump, 0)

	// Update jump position for false branch
	jumpNotTruthyPos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthy, jumpNotTruthyPos)

	// Compile false branch
	err = c.Compile(node.Alternative)
	if err != nil {
		return err
	}

	// Update jump position to skip false branch
	afterConditionalPos := len(c.currentInstructions())
	c.changeOperand(jumpPos, afterConditionalPos)

	return nil
}

// compileArrayLiteral compiles an array literal
func (c *Compiler) compileArrayLiteral(node *ast.ArrayLiteral) error {
	for _, el := range node.Elements {
		err := c.Compile(el)
		if err != nil {
			return err
		}
	}

	return c.emitError(vm.OpSlice, len(node.Elements))
}

// compileMapLiteral compiles a map literal
func (c *Compiler) compileMapLiteral(node *ast.MapLiteral) error {
	for _, pair := range node.Pairs {
		err := c.Compile(pair.Key)
		if err != nil {
			return err
		}
		err = c.Compile(pair.Value)
		if err != nil {
			return err
		}
	}

	return c.emitError(vm.OpMap, len(node.Pairs))
}

// compileLambdaExpression compiles a lambda expression
func (c *Compiler) compileLambdaExpression(node *ast.LambdaExpression) error {
	// Create a function value with the lambda parameters and body
	funcValue := types.NewFunc(node.Parameters, node.Body, nil, "")

	// Add the function as a constant
	return c.emitError(vm.OpConstant, c.addConstant(funcValue))
}

// compilePlaceholderExpression compiles a placeholder expression
func (c *Compiler) compilePlaceholderExpression(node *ast.PlaceholderExpression) error {
	// Instead of emitting OpGetPipelineElement immediately,
	// we emit a special placeholder constant that can be processed later
	placeholderValue := types.NewString("__PLACEHOLDER__")
	return c.emitError(vm.OpConstant, c.addConstant(placeholderValue))
}

// compilePipeExpression compiles a pipe expression
func (c *Compiler) compilePipeExpression(node *ast.PipeExpression) error {
	// Compile the left side of the pipe (data)
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	// Handle the right side differently based on its type
	switch right := node.Right.(type) {
	case *ast.Identifier:
		// For identifiers in pipe operations, emit them as string constants
		// so the VM can look them up in the builtin function registry
		funcName := types.NewString(right.Value)
		err = c.emitError(vm.OpConstant, c.addConstant(funcName))
		if err != nil {
			return err
		}
	case *ast.BuiltinExpression:
		// For builtin expressions with no arguments in pipe context,
		// emit as string constant
		if len(right.Arguments) == 0 {
			funcName := types.NewString(right.Name)
			err = c.emitError(vm.OpConstant, c.addConstant(funcName))
			if err != nil {
				return err
			}
		} else {
			// Check if any argument contains a placeholder
			hasPlaceholder := c.containsPlaceholder(right.Arguments)

			if hasPlaceholder {
				// For expressions with placeholders, emit special pipeline function bytecode
				err = c.compilePipelineFunction(right.Name, right.Arguments)
				if err != nil {
					return err
				}
			} else {
				// For builtin expressions with arguments but no placeholders, create an array:
				// [functionName, arg1, arg2, ...]

				// First emit the function name
				funcName := types.NewString(right.Name)
				err = c.emitError(vm.OpConstant, c.addConstant(funcName))
				if err != nil {
					return err
				}

				// Then emit each argument
				for _, arg := range right.Arguments {
					err = c.Compile(arg)
					if err != nil {
						return err
					}
				}

				// Create an array with the function name and arguments
				arraySize := 1 + len(right.Arguments)
				err = c.emitError(vm.OpSlice, arraySize)
				if err != nil {
					return err
				}
			}
		}
	case *ast.CallExpression:
		// Handle function calls that might contain placeholders
		if ident, ok := right.Function.(*ast.Identifier); ok {
			hasPlaceholder := c.containsPlaceholder(right.Arguments)

			if hasPlaceholder {
				// For expressions with placeholders, emit special pipeline function bytecode
				err = c.compilePipelineFunction(ident.Value, right.Arguments)
				if err != nil {
					return err
				}
			} else {
				// Regular function call compilation
				err = c.Compile(node.Right)
				if err != nil {
					return err
				}
			}
		} else if member, ok := right.Function.(*ast.MemberExpression); ok {
			// Handle type method calls like #.upper(), #.length(), etc.
			if c.hasPlaceholder(member) {
				// This is a type method call with placeholder: #.upper()
				err = c.compileTypeMethodCall(member, right.Arguments)
				if err != nil {
					return err
				}
			} else {
				// Regular member expression call
				err = c.Compile(node.Right)
				if err != nil {
					return err
				}
			}
		} else {
			// For other expressions (lambdas, etc.), compile normally
			err = c.Compile(node.Right)
			if err != nil {
				return err
			}
		}
	default:
		// For other expressions (lambdas, etc.), compile normally
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
	}

	// Emit the pipe operation
	return c.emitError(vm.OpPipe)
}

// emitError is a helper that wraps emit and converts the int result to error
func (c *Compiler) emitError(op vm.Opcode, operands ...int) error {
	c.emit(op, operands...)
	return nil
}

// emit emits an instruction with operands
func (c *Compiler) emit(op vm.Opcode, operands ...int) int {
	ins := vm.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

// addInstruction adds an instruction to the current scope
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

// setLastInstruction sets the last instruction info
func (c *Compiler) setLastInstruction(op vm.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// lastInstructionIs checks if the last instruction is of the given opcode
func (c *Compiler) lastInstructionIs(op vm.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// removeLastPop removes the last POP instruction if it exists
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

// replaceInstruction replaces an instruction at the given position
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

// changeOperand changes the operand of an instruction at the given position
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := vm.Opcode(c.currentInstructions()[opPos])
	newInstruction := vm.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

// currentInstructions returns the current instructions
func (c *Compiler) currentInstructions() []byte {
	return c.scopes[c.scopeIndex].instructions
}

// addConstant adds a constant to the constants pool
func (c *Compiler) addConstant(obj types.Value) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// loadSymbol loads a symbol based on its scope
func (c *Compiler) loadSymbol(s Symbol) error {
	switch s.Scope {
	case GlobalScope:
		return c.emitError(vm.OpGetVar, s.Index)
	case LocalScope:
		return c.emitError(vm.OpGetVar, s.Index)
	case BuiltinScope:
		return c.emitError(vm.OpBuiltin, s.Index, 0)
	case FreeScope:
		return c.emitError(vm.OpGetVar, s.Index)
	default:
		return fmt.Errorf("unknown symbol scope: %s", s.Scope)
	}
}

// Bytecode returns the compiled bytecode with optimizations applied
func (c *Compiler) Bytecode() *vm.Bytecode {
	instructions := c.currentInstructions()

	// Apply bytecode optimizations
	if c.optimizer != nil {
		instructions = c.optimizer.OptimizeInstructions(instructions)
	}

	return &vm.Bytecode{
		Instructions: instructions,
		Constants:    c.constants,
	}
}

// Errors returns the compilation errors
func (c *Compiler) Errors() []string {
	return c.errors
}

// GetVariableOrder returns the order of variables as they were defined
func (c *Compiler) GetVariableOrder() []string {
	// Find the maximum index first
	maxIndex := -1
	for _, symbol := range c.symbolTable.store {
		if symbol.Scope == GlobalScope && symbol.Index > maxIndex {
			maxIndex = symbol.Index
		}
	}

	if maxIndex < 0 {
		return []string{}
	}

	// Create order slice with the correct size
	order := make([]string, maxIndex+1)

	// Fill in the variable names at their correct indices
	for name, symbol := range c.symbolTable.store {
		if symbol.Scope == GlobalScope {
			order[symbol.Index] = name
		}
	}

	return order
}

// GetSymbolTable returns the symbol table for debugging
func (c *Compiler) GetSymbolTable() *SymbolTable {
	return c.symbolTable
}

// DefineBuiltin defines a custom builtin function
func (c *Compiler) DefineBuiltin(name string) {
	// Find the next available builtin index
	index := len(c.symbolTable.store)
	for i := 25; i < 100; i++ { // Start from 25 (after core builtins)
		found := false
		for _, symbol := range c.symbolTable.store {
			if symbol.Scope == BuiltinScope && symbol.Index == i {
				found = true
				break
			}
		}
		if !found {
			index = i
			break
		}
	}
	c.symbolTable.DefineBuiltin(index, name)
}

// CompileWithChecker compiles an AST node with type checking
func CompileWithChecker(node ast.Node, env interface{}) (*vm.Bytecode, error) {
	// For now, we'll skip type checking if it's not a Program
	// In a complete implementation, we'd need to handle individual expressions
	if program, ok := node.(*ast.Program); ok {
		// Create type checker
		checker := checker.New()

		// Type check the AST
		err := checker.Check(program)
		if err != nil {
			return nil, fmt.Errorf("type check failed: %w", err)
		}
	}

	// Compile to bytecode (with existing constant folding optimizations)
	compiler := New()
	err := compiler.Compile(node)
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	return compiler.Bytecode(), nil
}

// AddEnvironment adds environment variables to the compiler
func (c *Compiler) AddEnvironment(envVars map[string]interface{}, adapter *env.Adapter) error {
	// Sort variable names for consistent ordering
	var names []string
	for name := range envVars {
		names = append(names, name)
	}

	// Sort alphabetically for consistent ordering
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	// Define symbols in sorted order and store the mapping
	for i, name := range names {
		// Manually set the index to match the sorted order
		symbol := Symbol{Name: name, Index: i, Scope: GlobalScope}
		c.symbolTable.store[name] = symbol
	}

	// Update numDefinitions to reflect the total count
	c.symbolTable.numDefinitions = len(names)

	return nil
}

// containsPlaceholder checks if any expression in the slice contains a placeholder
func (c *Compiler) containsPlaceholder(expressions []ast.Expression) bool {
	for _, expr := range expressions {
		if c.hasPlaceholder(expr) {
			return true
		}
	}
	return false
}

// hasPlaceholder recursively checks if an expression contains a placeholder
func (c *Compiler) hasPlaceholder(expr ast.Expression) bool {
	switch node := expr.(type) {
	case *ast.PlaceholderExpression:
		return true
	case *ast.InfixExpression:
		return c.hasPlaceholder(node.Left) || c.hasPlaceholder(node.Right)
	case *ast.PrefixExpression:
		return c.hasPlaceholder(node.Right)
	case *ast.CallExpression:
		if c.hasPlaceholder(node.Function) {
			return true
		}
		for _, arg := range node.Arguments {
			if c.hasPlaceholder(arg) {
				return true
			}
		}
		return false
	case *ast.IndexExpression:
		return c.hasPlaceholder(node.Left) || c.hasPlaceholder(node.Index)
	case *ast.MemberExpression:
		return c.hasPlaceholder(node.Object) || c.hasPlaceholder(node.Property)
	case *ast.ConditionalExpression:
		return c.hasPlaceholder(node.Test) || c.hasPlaceholder(node.Consequent) || c.hasPlaceholder(node.Alternative)
	default:
		return false
	}
}

// compilePipelineFunction compiles a pipeline function with placeholders
func (c *Compiler) compilePipelineFunction(functionName string, arguments []ast.Expression) error {
	// Check each argument for type method calls
	for _, arg := range arguments {
		// Check if this is a direct type method call (e.g., #.upper())
		if call, ok := arg.(*ast.CallExpression); ok {
			if member, ok := call.Function.(*ast.MemberExpression); ok {
				if c.hasPlaceholder(member) {
					// This is a simple type method call: functionName(#.methodName())
					return c.compilePipelineWithTypeMethod(functionName, member, call.Arguments)
				}
			}
		}

		// Check if this is a complex expression containing type method calls (e.g., #.length() > 4)
		if typeMethodCall := c.findTypeMethodCallInExpression(arg); typeMethodCall != nil {
			// This argument contains a type method call within a complex expression
			member := typeMethodCall.Function.(*ast.MemberExpression)
			return c.compilePipelineWithTypeMethodInExpression(functionName, arg, member, typeMethodCall.Arguments)
		}
	}

	// Emit the function name as a constant
	funcName := types.NewString(functionName)
	err := c.emitError(vm.OpConstant, c.addConstant(funcName))
	if err != nil {
		return err
	}

	// For pipeline functions with placeholders, we need to emit special bytecode
	// that tells the VM this is a placeholder expression that needs to be evaluated
	// for each element in the pipeline

	// Emit a special marker to indicate this is a placeholder expression
	err = c.emitError(vm.OpConstant, c.addConstant(types.NewString("__PLACEHOLDER_EXPR__")))
	if err != nil {
		return err
	}

	// Set pipeline context to prevent immediate compilation of placeholder expressions
	oldContext := c.inPipelineContext
	c.inPipelineContext = true
	defer func() {
		c.inPipelineContext = oldContext
	}()

	// Compile the expression with placeholders
	// The VM will later replace placeholders with actual values
	for _, arg := range arguments {
		err = c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// Create an array: [functionName, "__PLACEHOLDER_EXPR__", arg1, arg2, ...]
	arraySize := 2 + len(arguments)
	return c.emitError(vm.OpSlice, arraySize)
}

// findTypeMethodCall recursively searches for type method calls in complex expressions (not simple calls)
func (c *Compiler) findTypeMethodCall(expr ast.Expression) *ast.CallExpression {
	switch e := expr.(type) {
	case *ast.CallExpression:
		// Don't return top-level call expressions - they are handled separately
		// Only check arguments recursively
		for _, arg := range e.Arguments {
			if result := c.findTypeMethodCall(arg); result != nil {
				return result
			}
		}
	case *ast.InfixExpression:
		// Check both sides of infix expression
		if result := c.findTypeMethodCall(e.Left); result != nil {
			return result
		}
		if result := c.findTypeMethodCall(e.Right); result != nil {
			return result
		}
	case *ast.PrefixExpression:
		// Check operand
		if result := c.findTypeMethodCall(e.Right); result != nil {
			return result
		}
	}
	return nil
}

// findTypeMethodCallInExpression searches for type method calls including top-level calls
func (c *Compiler) findTypeMethodCallInExpression(expr ast.Expression) *ast.CallExpression {
	switch e := expr.(type) {
	case *ast.CallExpression:
		// Check if this is a type method call
		if member, ok := e.Function.(*ast.MemberExpression); ok {
			if c.hasPlaceholder(member) {
				return e
			}
		}
		// Check arguments recursively
		for _, arg := range e.Arguments {
			if result := c.findTypeMethodCallInExpression(arg); result != nil {
				return result
			}
		}
	case *ast.InfixExpression:
		// Check both sides of infix expression
		if result := c.findTypeMethodCallInExpression(e.Left); result != nil {
			return result
		}
		if result := c.findTypeMethodCallInExpression(e.Right); result != nil {
			return result
		}
	case *ast.PrefixExpression:
		// Check operand
		if result := c.findTypeMethodCallInExpression(e.Right); result != nil {
			return result
		}
	}
	return nil
}

// compilePipelineWithTypeMethod compiles a pipeline function that contains a type method call
func (c *Compiler) compilePipelineWithTypeMethod(functionName string, memberExpr *ast.MemberExpression, methodArgs []ast.Expression) error {
	// Extract method name from the member expression
	var methodName string
	if ident, ok := memberExpr.Property.(*ast.Identifier); ok {
		methodName = ident.Value
	} else {
		return fmt.Errorf("unsupported property type in method call: %T", memberExpr.Property)
	}

	// Emit function name constant
	funcNameConst := types.NewString(functionName)
	err := c.emitError(vm.OpConstant, c.addConstant(funcNameConst))
	if err != nil {
		return err
	}

	// Emit a special marker to indicate this is a pipeline function with type method
	err = c.emitError(vm.OpConstant, c.addConstant(types.NewString("__PIPELINE_TYPE_METHOD__")))
	if err != nil {
		return err
	}

	// Emit the method name
	methodNameConst := types.NewString(methodName)
	err = c.emitError(vm.OpConstant, c.addConstant(methodNameConst))
	if err != nil {
		return err
	}

	// Set pipeline context and compile the object (placeholder)
	oldContext := c.inPipelineContext
	c.inPipelineContext = true
	defer func() {
		c.inPipelineContext = oldContext
	}()

	err = c.Compile(memberExpr.Object)
	if err != nil {
		return err
	}

	// Compile method arguments
	for _, arg := range methodArgs {
		err = c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// Create an array: [functionName, "__PIPELINE_TYPE_METHOD__", methodName, object, arg1, arg2, ...]
	arraySize := 3 + 1 + len(methodArgs) // 3 constants + object + method arguments
	return c.emitError(vm.OpSlice, arraySize)
}

// compilePipelineWithTypeMethodInExpression compiles a pipeline function that contains a type method call within a complex expression
func (c *Compiler) compilePipelineWithTypeMethodInExpression(functionName string, fullExpr ast.Expression, memberExpr *ast.MemberExpression, methodArgs []ast.Expression) error {
	// Extract method name from the member expression
	var methodName string
	if ident, ok := memberExpr.Property.(*ast.Identifier); ok {
		methodName = ident.Value
	} else {
		return fmt.Errorf("unsupported property type in method call: %T", memberExpr.Property)
	}

	// Create a nested compilation scope to capture the full expression as a single constant
	nestedCompiler := New()
	nestedCompiler.inPipelineContext = true

	// Compile the full expression in the nested compiler
	err := nestedCompiler.Compile(fullExpr)
	if err != nil {
		return err
	}

	// Get the compiled expression as a slice of constants
	nestedBytecode := nestedCompiler.Bytecode()
	expressionConstants := nestedBytecode.Constants

	// Create a slice value containing all the expression constants
	var expressionSlice []types.Value
	for _, constant := range expressionConstants {
		expressionSlice = append(expressionSlice, constant)
	}

	compiledExpression := types.NewSlice(expressionSlice, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"})

	// Now emit the pipeline function structure
	// Emit function name constant
	funcNameConst := types.NewString(functionName)
	err = c.emitError(vm.OpConstant, c.addConstant(funcNameConst))
	if err != nil {
		return err
	}

	// Emit a special marker to indicate this is a pipeline function with complex type method expression
	err = c.emitError(vm.OpConstant, c.addConstant(types.NewString("__PIPELINE_COMPLEX_TYPE_METHOD__")))
	if err != nil {
		return err
	}

	// Emit the method name
	methodNameConst := types.NewString(methodName)
	err = c.emitError(vm.OpConstant, c.addConstant(methodNameConst))
	if err != nil {
		return err
	}

	// Emit the compiled expression as a single constant
	err = c.emitError(vm.OpConstant, c.addConstant(compiledExpression))
	if err != nil {
		return err
	}

	// Create an array: [functionName, "__PIPELINE_COMPLEX_TYPE_METHOD__", methodName, compiledExpression]
	arraySize := 4 // 3 constants + compiled expression
	return c.emitError(vm.OpSlice, arraySize)
}

// compilePlaceholderInfixExpression compiles an infix expression containing placeholders
// in pipeline context by serializing it for later evaluation
func (c *Compiler) compilePlaceholderInfixExpression(node *ast.InfixExpression) error {
	// For now, we'll emit a serialized representation of the expression
	// that the VM can interpret at runtime

	// Emit the operator as a string constant
	operatorStr := types.NewString(node.Operator)
	err := c.emitError(vm.OpConstant, c.addConstant(operatorStr))
	if err != nil {
		return err
	}

	// Compile the left operand (might be a placeholder)
	err = c.Compile(node.Left)
	if err != nil {
		return err
	}

	// Compile the right operand (might be a placeholder)
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	// Create an array: [operator, left, right]
	return c.emitError(vm.OpSlice, 3)
}

// compileTypeMethodCall compiles a type method call with placeholders
func (c *Compiler) compileTypeMethodCall(node *ast.MemberExpression, arguments []ast.Expression) error {
	// Extract method name from the property
	var methodName string
	if ident, ok := node.Property.(*ast.Identifier); ok {
		methodName = ident.Value
	} else {
		return fmt.Errorf("unsupported property type in method call: %T", node.Property)
	}

	// Emit function name constant to indicate this is a type method call
	funcNameConst := types.NewString("__TYPE_METHOD__")
	err := c.emitError(vm.OpConstant, c.addConstant(funcNameConst))
	if err != nil {
		return err
	}

	// Emit a special marker to indicate this is a placeholder expression
	err = c.emitError(vm.OpConstant, c.addConstant(types.NewString("__PLACEHOLDER_EXPR__")))
	if err != nil {
		return err
	}

	// Emit the method name as a string constant
	methodNameConst := types.NewString(methodName)
	err = c.emitError(vm.OpConstant, c.addConstant(methodNameConst))
	if err != nil {
		return err
	}

	// Set pipeline context and compile the object (usually a placeholder)
	oldContext := c.inPipelineContext
	c.inPipelineContext = true
	defer func() {
		c.inPipelineContext = oldContext
	}()

	err = c.Compile(node.Object)
	if err != nil {
		return err
	}

	// Compile the arguments
	for _, arg := range arguments {
		err = c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// Create an array: ["__TYPE_METHOD__", "__PLACEHOLDER_EXPR__", methodName, object, arg1, arg2, ...]
	arraySize := 3 + 1 + len(arguments) // 3 special markers + object + arguments
	return c.emitError(vm.OpSlice, arraySize)
}

// Phase 3: Simplified compilation optimizations - focus on what works
func (c *Compiler) optimizeExpression(node ast.Expression) ast.Expression {
	// For now, just return the node as-is to avoid type complexity
	// Focus on runtime optimizations instead of compile-time
	return node
}

// compileOptionalChainingExpression compiles optional chaining expressions (obj?.property)
func (c *Compiler) compileOptionalChainingExpression(node *ast.OptionalChainingExpression) error {
	// Compile the object expression
	err := c.Compile(node.Object)
	if err != nil {
		return err
	}

	// Handle property based on its type
	switch prop := node.Property.(type) {
	case *ast.Identifier:
		// Simple property access: obj?.property
		// Emit property name as a string constant
		err = c.emitError(vm.OpConstant, c.addConstant(types.NewString(prop.Value)))
		if err != nil {
			return err
		}
	default:
		// Computed property access: obj?.[expr]
		err = c.Compile(node.Property)
		if err != nil {
			return err
		}
	}

	// Emit the optional chaining instruction
	return c.emitError(vm.OpOptionalChaining)
}

// compileNullCoalescingExpression compiles null coalescing expressions (a ?? b)
func (c *Compiler) compileNullCoalescingExpression(node *ast.NullCoalescingExpression) error {
	// Compile the left operand
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	// Compile the right operand (default value)
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	// Emit the null coalescing instruction
	return c.emitError(vm.OpNullCoalescing)
}

// compileImportStatement compiles an import statement
func (c *Compiler) compileImportStatement(node *ast.ImportStatement) error {
	// Import statements are handled at compile time, so we emit nothing
	// The module registry will be populated during compilation
	return nil
}

// compileModuleCallExpression compiles a module function call
func (c *Compiler) compileModuleCallExpression(node *ast.ModuleCallExpression) error {
	// Compile arguments
	for _, arg := range node.Arguments {
		err := c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// Emit module call instruction
	// We need to store the module name and function name as constants
	moduleNameIndex := c.addConstant(types.NewString(node.Module))
	functionNameIndex := c.addConstant(types.NewString(node.Function))

	return c.emitError(vm.OpModuleCall, moduleNameIndex, functionNameIndex, len(node.Arguments))
}

// compileDestructuringAssignment compiles a destructuring assignment statement
func (c *Compiler) compileDestructuringAssignment(node *ast.DestructuringAssignment) error {
	// First, compile the right side (the value to destructure)
	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	// Then, compile the left side (the destructuring pattern)
	return c.compileDestructuringPattern(node.Left)
}

// compileDestructuringPattern compiles a destructuring pattern
func (c *Compiler) compileDestructuringPattern(pattern ast.DestructuringPattern) error {
	switch p := pattern.(type) {
	case *ast.ArrayDestructuringPattern:
		return c.compileArrayDestructuringPattern(p)
	case *ast.ObjectDestructuringPattern:
		return c.compileObjectDestructuringPattern(p)
	default:
		return fmt.Errorf("unknown destructuring pattern type: %T", pattern)
	}
}

// compileArrayDestructuringPattern compiles array destructuring like [a, b, c] = [1, 2, 3]
func (c *Compiler) compileArrayDestructuringPattern(pattern *ast.ArrayDestructuringPattern) error {
	// Collect variable assignments
	varAssignments := make([]string, 0)
	hasRestElement := false

	// Process each element in the destructuring pattern
	for _, element := range pattern.Elements {
		switch elem := element.(type) {
		case *ast.IdentifierElement:
			varAssignments = append(varAssignments, elem.Name)
		case *ast.RestElement:
			hasRestElement = true
			varAssignments = append(varAssignments, elem.Name)
		default:
			return fmt.Errorf("unsupported array destructuring element type: %T", element)
		}
	}

	// Define symbols for variables that don't exist
	for _, name := range varAssignments {
		if symbol, ok := c.symbolTable.Resolve(name); !ok {
			// Define new variable
			c.symbolTable.Define(name)
		} else if symbol.Scope != LocalScope {
			// Redefine in local scope
			c.symbolTable.Define(name)
		}
	}

	// Get the starting variable index
	startVarIndex := 0
	if len(varAssignments) > 0 {
		if symbol, ok := c.symbolTable.Resolve(varAssignments[0]); ok {
			startVarIndex = symbol.Index
		}
	}

	// Emit array destructuring instruction
	elementCount := len(pattern.Elements)
	if hasRestElement {
		// For rest elements, we need special handling
		return c.emitError(vm.OpArrayDestructure, elementCount-1, startVarIndex) // -1 because rest element consumes remaining
	}

	return c.emitError(vm.OpArrayDestructure, elementCount, startVarIndex)
}

// compileObjectDestructuringPattern compiles object destructuring like {name, age} = user
func (c *Compiler) compileObjectDestructuringPattern(pattern *ast.ObjectDestructuringPattern) error {
	// Collect variable assignments
	varAssignments := make([]string, 0)
	propertyKeys := make([]string, 0)

	// Process each property in the destructuring pattern
	for _, prop := range pattern.Properties {
		varAssignments = append(varAssignments, prop.Value) // Variable name
		propertyKeys = append(propertyKeys, prop.Key)       // Property key
	}

	// Define symbols for variables that don't exist
	for _, name := range varAssignments {
		if symbol, ok := c.symbolTable.Resolve(name); !ok {
			// Define new variable
			c.symbolTable.Define(name)
		} else if symbol.Scope != LocalScope {
			// Redefine in local scope
			c.symbolTable.Define(name)
		}
	}

	// Get the starting variable index
	startVarIndex := 0
	if len(varAssignments) > 0 {
		if symbol, ok := c.symbolTable.Resolve(varAssignments[0]); ok {
			startVarIndex = symbol.Index
		}
	}

	// Emit property keys as constants
	for _, key := range propertyKeys {
		c.emitError(vm.OpConstant, c.addConstant(types.NewString(key)))
	}

	// Emit object destructuring instruction
	propertyCount := len(pattern.Properties)
	return c.emitError(vm.OpObjectDestructure, propertyCount, startVarIndex)
}

// Phase 3: Focus on runtime optimizations instead of complex compile-time optimizations
// This simplifies the codebase and avoids type complexity issues
