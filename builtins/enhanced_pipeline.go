package builtins

import (
	"fmt"

	"github.com/mredencom/expr/types"
)

// EnhancedPipelineProcessor provides advanced pipeline processing with Lambda and placeholder support
type EnhancedPipelineProcessor struct {
	// We'll avoid import cycles by using interfaces and simpler implementation
}

// NewEnhancedPipelineProcessor creates a new enhanced pipeline processor
func NewEnhancedPipelineProcessor() *EnhancedPipelineProcessor {
	return &EnhancedPipelineProcessor{}
}

// PredicateType represents the type of predicate in pipeline operations
type PredicateType int

const (
	PredicateUnknown PredicateType = iota
	PredicateLambda
	PredicatePlaceholder
	PredicateString
	PredicateConstant
)

// PipelineExpression represents a compiled pipeline expression
type PipelineExpression struct {
	Type        PredicateType
	Lambda      *types.FuncValue
	Placeholder *types.PlaceholderExprValue
	StringValue string
	Bytecode    []byte
	Constants   []types.Value
}

// DetectPredicateType analyzes a value to determine its predicate type
func (epp *EnhancedPipelineProcessor) DetectPredicateType(predicate types.Value) PredicateType {
	switch predicate.(type) {
	case *types.FuncValue:
		return PredicateLambda
	case *types.StringValue:
		return PredicateString
	case *types.PlaceholderExprValue:
		return PredicatePlaceholder
	case *types.BoolValue, *types.IntValue, *types.FloatValue:
		return PredicateConstant
	default:
		return PredicateUnknown
	}
}

// CompilePipelineExpression compiles different types of pipeline expressions
func (epp *EnhancedPipelineProcessor) CompilePipelineExpression(predicate types.Value) (*PipelineExpression, error) {
	pType := epp.DetectPredicateType(predicate)

	expr := &PipelineExpression{
		Type: pType,
	}

	switch pType {
	case PredicateLambda:
		expr.Lambda = predicate.(*types.FuncValue)

	case PredicatePlaceholder:
		// For placeholder expressions, use the built-in bytecode
		placeholderExpr := predicate.(*types.PlaceholderExprValue)
		expr.Placeholder = placeholderExpr
		expr.Bytecode = placeholderExpr.Instructions()
		expr.Constants = placeholderExpr.Constants()

	case PredicateString:
		// Store string value for simple evaluation
		strVal := predicate.(*types.StringValue).Value()
		expr.StringValue = strVal
	}

	return expr, nil
}

// EvaluatePredicate evaluates a compiled predicate against an item
func (epp *EnhancedPipelineProcessor) EvaluatePredicate(expr *PipelineExpression, item types.Value) (bool, error) {
	switch expr.Type {
	case PredicateLambda:
		// Execute Lambda function
		result, err := epp.executeLambda(expr.Lambda, item)
		if err != nil {
			return false, err
		}
		return isTruthy(result), nil

	case PredicatePlaceholder:
		// Execute placeholder expression with item as context
		if len(expr.Bytecode) > 0 {
			result, err := epp.executePlaceholderBytecode(expr.Bytecode, expr.Constants, item)
			if err != nil {
				return false, err
			}
			return isTruthy(result), nil
		}
		return true, nil

	case PredicateString:
		// Use enhanced string predicate evaluation
		if len(expr.Bytecode) > 0 {
			result, err := epp.executePlaceholderBytecode(expr.Bytecode, expr.Constants, item)
			if err != nil {
				return false, err
			}
			return isTruthy(result), nil
		}
		// Fall back to simple string evaluation
		return evaluateStringPredicate("", item), nil

	case PredicateConstant:
		// Constant predicates are simple
		return true, nil

	default:
		return true, nil
	}
}

// EvaluateTransformer evaluates a compiled transformer against an item
func (epp *EnhancedPipelineProcessor) EvaluateTransformer(expr *PipelineExpression, item types.Value) (types.Value, error) {
	switch expr.Type {
	case PredicateLambda:
		// Execute Lambda function
		return epp.executeLambda(expr.Lambda, item)

	case PredicatePlaceholder:
		// Execute placeholder expression with item as context
		if len(expr.Bytecode) > 0 {
			return epp.executePlaceholderBytecode(expr.Bytecode, expr.Constants, item)
		}
		return item, nil

	case PredicateString:
		// Use enhanced string transformation
		if len(expr.Bytecode) > 0 {
			return epp.executePlaceholderBytecode(expr.Bytecode, expr.Constants, item)
		}
		// Fall back to simple string transformation
		return applyStringTransform("", item), nil

	default:
		return item, nil
	}
}

// EvaluateReducer evaluates a compiled reducer with accumulator and item
func (epp *EnhancedPipelineProcessor) EvaluateReducer(expr *PipelineExpression, acc, item types.Value) (types.Value, error) {
	switch expr.Type {
	case PredicateLambda:
		// Execute Lambda function with two arguments
		return epp.executeLambdaWithTwoArgs(expr.Lambda, acc, item)

	case PredicatePlaceholder:
		// For placeholder reducers, we need special handling
		// This is more complex as placeholders typically take one argument
		return acc, nil

	case PredicateString:
		// Use string-based reducer
		return applyStringReducer("add", acc, item), nil

	default:
		return acc, nil
	}
}

// executeLambda executes a Lambda function with simplified evaluation
func (epp *EnhancedPipelineProcessor) executeLambda(lambda *types.FuncValue, arg types.Value) (types.Value, error) {
	// For now, use simplified Lambda evaluation
	// In a full implementation, this would need proper VM execution

	// Check if it's a simple identity function
	params := lambda.Parameters()
	if len(params) == 1 {
		body := lambda.Body()
		if bodyStr, ok := body.(string); ok && bodyStr == params[0] {
			return arg, nil
		}
	}

	// For complex lambdas, we'd need proper execution
	// This is a placeholder that maintains compatibility
	return arg, nil
}

// executeLambdaWithTwoArgs executes a Lambda function with two arguments
func (epp *EnhancedPipelineProcessor) executeLambdaWithTwoArgs(lambda *types.FuncValue, arg1, arg2 types.Value) (types.Value, error) {
	// Simplified two-argument Lambda evaluation
	// For basic operations like addition, we can handle some cases

	params := lambda.Parameters()
	if len(params) >= 2 {
		// Check for simple addition pattern
		body := lambda.Body()
		if bodyStr, ok := body.(string); ok {
			if bodyStr == params[0]+"+"+params[1] {
				// Simple addition
				if int1, ok1 := arg1.(*types.IntValue); ok1 {
					if int2, ok2 := arg2.(*types.IntValue); ok2 {
						return types.NewInt(int1.Value() + int2.Value()), nil
					}
				}
			}
		}
	}

	return arg1, nil
}

// executePlaceholderBytecode evaluates placeholder expressions
func (epp *EnhancedPipelineProcessor) executePlaceholderBytecode(bytecode []byte, constants []types.Value, item types.Value) (types.Value, error) {
	// For placeholder expressions, we use a simplified evaluation approach
	// This would need proper VM integration in a full implementation

	// For now, return true for simple predicates
	return types.NewBool(true), nil
}

// Enhanced builtin functions that use the new processor

// EnhancedFilterFunc provides improved filter with Lambda and placeholder support
var EnhancedFilterFunc BuiltinFunction = func(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("filter requires exactly 2 arguments")
	}

	collection := args[0]
	predicate := args[1]

	if array, ok := collection.(*types.SliceValue); ok {
		var result []types.Value
		for i := 0; i < array.Len(); i++ {
			item := array.Get(i)
			if shouldIncludeEnhanced(item, predicate) {
				result = append(result, item)
			}
		}
		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("filter can only be applied to arrays")
}

// EnhancedMapFunc provides improved map with Lambda and placeholder support
var EnhancedMapFunc BuiltinFunction = func(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map requires exactly 2 arguments")
	}

	collection := args[0]
	transformer := args[1]

	if array, ok := collection.(*types.SliceValue); ok {
		var result []types.Value
		for i := 0; i < array.Len(); i++ {
			item := array.Get(i)
			transformed := applyTransformEnhanced(item, transformer)
			result = append(result, transformed)
		}
		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("map can only be applied to arrays")
}

// EnhancedReduceFunc provides improved reduce with Lambda and placeholder support
var EnhancedReduceFunc BuiltinFunction = func(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("reduce requires 2 or 3 arguments")
	}

	collection := args[0]
	reducer := args[1]

	if array, ok := collection.(*types.SliceValue); ok {
		if array.Len() == 0 {
			if len(args) == 3 {
				return args[2], nil // Return initial value
			}
			return types.NewNil(), nil
		}

		var acc types.Value
		var start int

		if len(args) == 3 {
			acc = args[2]
			start = 0
		} else {
			acc = array.Get(0)
			start = 1
		}

		for i := start; i < array.Len(); i++ {
			acc = applyReducerEnhanced(acc, array.Get(i), reducer)
		}

		return acc, nil
	}

	return nil, fmt.Errorf("reduce can only be applied to arrays")
}

// shouldIncludeEnhanced evaluates predicates with enhanced Lambda and placeholder support
func shouldIncludeEnhanced(item types.Value, predicate types.Value) bool {
	switch pred := predicate.(type) {
	case *types.FuncValue:
		// Enhanced Lambda evaluation
		return evaluateLambdaPredicate(pred, item)
	case *types.PlaceholderExprValue:
		// Enhanced placeholder evaluation
		return evaluatePlaceholderPredicate(pred, item)
	case *types.StringValue:
		// Enhanced string predicates
		return evaluateEnhancedStringPredicate(pred.Value(), item)
	case *types.BoolValue:
		return pred.Value()
	default:
		return true
	}
}

// applyTransformEnhanced applies transformations with enhanced Lambda and placeholder support
func applyTransformEnhanced(item types.Value, transformer types.Value) types.Value {
	switch trans := transformer.(type) {
	case *types.FuncValue:
		// Enhanced Lambda transformation
		return evaluateLambdaTransform(trans, item)
	case *types.PlaceholderExprValue:
		// Enhanced placeholder transformation
		return evaluatePlaceholderTransform(trans, item)
	case *types.StringValue:
		// Enhanced string transformations
		return applyEnhancedStringTransform(trans.Value(), item)
	default:
		return item
	}
}

// applyReducerEnhanced applies reducers with enhanced Lambda and placeholder support
func applyReducerEnhanced(acc types.Value, item types.Value, reducer types.Value) types.Value {
	switch red := reducer.(type) {
	case *types.FuncValue:
		// Enhanced Lambda reduction
		return evaluateLambdaReducer(red, acc, item)
	case *types.StringValue:
		// Enhanced string reducers
		return applyEnhancedStringReducer(red.Value(), acc, item)
	default:
		return acc
	}
}

// evaluateLambdaPredicate evaluates Lambda predicates
func evaluateLambdaPredicate(lambda *types.FuncValue, item types.Value) bool {
	// Enhanced Lambda evaluation for predicates
	params := lambda.Parameters()
	if len(params) == 0 {
		return true
	}

	// For simple Lambda expressions, we can do basic pattern matching
	body := lambda.Body()
	if bodyStr, ok := body.(string); ok {
		// Handle common Lambda patterns
		paramName := params[0]

		// Simple identity check: x => x
		if bodyStr == paramName {
			return isTruthy(item)
		}

		// Simple comparison: x => x > 5
		if len(bodyStr) > len(paramName)+3 && bodyStr[:len(paramName)] == paramName {
			operator := ""
			operandStr := ""
			rest := bodyStr[len(paramName):]

			// Extract operator and operand
			for i, op := range []string{" >= ", " <= ", " > ", " < ", " == ", " != "} {
				if len(rest) > len(op) && rest[:len(op)] == op {
					operator = []string{">=", "<=", ">", "<", "==", "!="}[i]
					operandStr = rest[len(op):]
					break
				}
			}

			if operator != "" {
				return evaluateComparison(item, operator, operandStr)
			}
		}
	}

	// Fallback to truthy evaluation
	return isTruthy(item)
}

// evaluateLambdaTransform evaluates Lambda transformations
func evaluateLambdaTransform(lambda *types.FuncValue, item types.Value) types.Value {
	// Enhanced Lambda evaluation for transformations
	params := lambda.Parameters()
	if len(params) == 0 {
		return item
	}

	// For simple Lambda expressions, we can do basic pattern matching
	body := lambda.Body()
	if bodyStr, ok := body.(string); ok {
		paramName := params[0]

		// Simple identity: x => x
		if bodyStr == paramName {
			return item
		}

		// Simple arithmetic: x => x * 2
		if len(bodyStr) > len(paramName)+3 && bodyStr[:len(paramName)] == paramName {
			operator := ""
			operandStr := ""
			rest := bodyStr[len(paramName):]

			// Extract operator and operand
			for i, op := range []string{" * ", " + ", " - ", " / "} {
				if len(rest) > len(op) && rest[:len(op)] == op {
					operator = []string{"*", "+", "-", "/"}[i]
					operandStr = rest[len(op):]
					break
				}
			}

			if operator != "" {
				return evaluateArithmetic(item, operator, operandStr)
			}
		}
	}

	return item
}

// evaluateLambdaReducer evaluates Lambda reducers
func evaluateLambdaReducer(lambda *types.FuncValue, acc, item types.Value) types.Value {
	// Enhanced Lambda evaluation for reducers
	params := lambda.Parameters()
	if len(params) < 2 {
		return acc
	}

	// For simple Lambda expressions, we can handle basic operations
	body := lambda.Body()
	if bodyStr, ok := body.(string); ok {
		accParam := params[0]
		itemParam := params[1]

		// Simple addition: (a, b) => a + b
		if bodyStr == accParam+" + "+itemParam {
			return addValues(acc, item)
		}

		// Simple multiplication: (a, b) => a * b
		if bodyStr == accParam+" * "+itemParam {
			return multiplyValues(acc, item)
		}
	}

	return acc
}

// evaluatePlaceholderPredicate evaluates placeholder predicates
func evaluatePlaceholderPredicate(placeholder *types.PlaceholderExprValue, item types.Value) bool {
	// Use the placeholder's operator and operand
	operator := placeholder.Operator()
	operand := placeholder.Operand()

	return evaluateComparison(item, operator, operand.String())
}

// evaluatePlaceholderTransform evaluates placeholder transformations
func evaluatePlaceholderTransform(placeholder *types.PlaceholderExprValue, item types.Value) types.Value {
	// For placeholder transformations, we need to apply the operation
	operator := placeholder.Operator()
	operand := placeholder.Operand()

	return evaluateArithmetic(item, operator, operand.String())
}

// Helper functions for enhanced string operations

func evaluateEnhancedStringPredicate(predicate string, item types.Value) bool {
	// Enhanced string predicate evaluation with more patterns
	switch predicate {
	case "not_empty", "non_empty", "truthy":
		return !isZeroValue(item)
	case "empty", "falsy":
		return isZeroValue(item)
	case "positive":
		return isPositive(item)
	case "negative":
		return isNegative(item)
	case "even":
		return isEven(item)
	case "odd":
		return isOdd(item)
	case "numeric":
		return isNumeric(item)
	case "string":
		return isString(item)
	default:
		return evaluateStringPredicate(predicate, item)
	}
}

func applyEnhancedStringTransform(transform string, item types.Value) types.Value {
	// Enhanced string transformation with more operations
	switch transform {
	case "double", "times2", "*2":
		return multiplyByTwo(item)
	case "square", "^2":
		return square(item)
	case "abs", "absolute":
		return absolute(item)
	case "upper", "uppercase":
		return toUpper(item)
	case "lower", "lowercase":
		return toLower(item)
	case "length", "len":
		return getLength(item)
	case "string", "toString":
		return types.NewString(item.String())
	default:
		return applyStringTransform(transform, item)
	}
}

func applyEnhancedStringReducer(reducer string, acc, item types.Value) types.Value {
	// Enhanced string reducer operations
	switch reducer {
	case "add", "sum", "+":
		return addValues(acc, item)
	case "multiply", "mul", "*":
		return multiplyValues(acc, item)
	case "max", "maximum":
		return maxValue(acc, item)
	case "min", "minimum":
		return minValue(acc, item)
	case "concat", "join":
		return types.NewString(acc.String() + item.String())
	default:
		return applyStringReducer(reducer, acc, item)
	}
}

// Helper utility functions

func evaluateComparison(left types.Value, operator, rightStr string) bool {
	// Convert right operand based on left type
	var right types.Value

	switch left.(type) {
	case *types.IntValue:
		if val, err := parseIntValue(rightStr); err == nil {
			right = types.NewInt(val)
		} else {
			return false
		}
	case *types.FloatValue:
		if val, err := parseFloatValue(rightStr); err == nil {
			right = types.NewFloat(val)
		} else {
			return false
		}
	case *types.StringValue:
		right = types.NewString(rightStr)
	default:
		return false
	}

	return compareValuesWithOperator(left, right, operator)
}

func evaluateArithmetic(left types.Value, operator, rightStr string) types.Value {
	// Convert right operand based on left type
	switch leftVal := left.(type) {
	case *types.IntValue:
		if val, err := parseIntValue(rightStr); err == nil {
			return applyArithmeticInt(leftVal.Value(), int64(val), operator)
		}
	case *types.FloatValue:
		if val, err := parseFloatValue(rightStr); err == nil {
			return applyArithmeticFloat(leftVal.Value(), val, operator)
		}
	}

	return left
}

// Additional utility functions will be implemented as needed...

// parseIntValue parses string to int64
func parseIntValue(s string) (int64, error) {
	// Simple integer parsing
	var result int64
	var sign int64 = 1

	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}

	i := 0
	if s[0] == '-' {
		sign = -1
		i = 1
	} else if s[0] == '+' {
		i = 1
	}

	for ; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int64(s[i]-'0')
		} else {
			return 0, fmt.Errorf("invalid integer")
		}
	}

	return result * sign, nil
}

// parseFloatValue parses string to float64
func parseFloatValue(s string) (float64, error) {
	// Simplified float parsing - in production would use strconv.ParseFloat
	if intVal, err := parseIntValue(s); err == nil {
		return float64(intVal), nil
	}
	return 0.0, fmt.Errorf("invalid float")
}

// Additional helper functions for enhanced operations...
func addValues(a, b types.Value) types.Value {
	switch av := a.(type) {
	case *types.IntValue:
		if bv, ok := b.(*types.IntValue); ok {
			return types.NewInt(av.Value() + bv.Value())
		}
	case *types.FloatValue:
		if bv, ok := b.(*types.FloatValue); ok {
			return types.NewFloat(av.Value() + bv.Value())
		}
	}
	return a
}

func multiplyValues(a, b types.Value) types.Value {
	switch av := a.(type) {
	case *types.IntValue:
		if bv, ok := b.(*types.IntValue); ok {
			return types.NewInt(av.Value() * bv.Value())
		}
	case *types.FloatValue:
		if bv, ok := b.(*types.FloatValue); ok {
			return types.NewFloat(av.Value() * bv.Value())
		}
	}
	return a
}

func compareValuesWithOperator(left, right types.Value, operator string) bool {
	switch operator {
	case ">":
		return compareValues(left, right) > 0
	case "<":
		return compareValues(left, right) < 0
	case ">=":
		return compareValues(left, right) >= 0
	case "<=":
		return compareValues(left, right) <= 0
	case "==":
		return left.Equal(right)
	case "!=":
		return !left.Equal(right)
	}
	return false
}

func applyArithmeticInt(left int64, right int64, operator string) types.Value {
	switch operator {
	case "+":
		return types.NewInt(left + right)
	case "-":
		return types.NewInt(left - right)
	case "*":
		return types.NewInt(left * right)
	case "/":
		if right != 0 {
			return types.NewInt(left / right)
		}
	}
	return types.NewInt(left)
}

func applyArithmeticFloat(left float64, right float64, operator string) types.Value {
	switch operator {
	case "+":
		return types.NewFloat(left + right)
	case "-":
		return types.NewFloat(left - right)
	case "*":
		return types.NewFloat(left * right)
	case "/":
		if right != 0 {
			return types.NewFloat(left / right)
		}
	}
	return types.NewFloat(left)
}

// Additional helper functions for type checking and transformations
func isPositive(v types.Value) bool {
	switch val := v.(type) {
	case *types.IntValue:
		return val.Value() > 0
	case *types.FloatValue:
		return val.Value() > 0
	}
	return false
}

func isNegative(v types.Value) bool {
	switch val := v.(type) {
	case *types.IntValue:
		return val.Value() < 0
	case *types.FloatValue:
		return val.Value() < 0
	}
	return false
}

func isEven(v types.Value) bool {
	if intVal, ok := v.(*types.IntValue); ok {
		return intVal.Value()%2 == 0
	}
	return false
}

func isOdd(v types.Value) bool {
	if intVal, ok := v.(*types.IntValue); ok {
		return intVal.Value()%2 != 0
	}
	return false
}

func isNumeric(v types.Value) bool {
	switch v.(type) {
	case *types.IntValue, *types.FloatValue:
		return true
	}
	return false
}

func isString(v types.Value) bool {
	_, ok := v.(*types.StringValue)
	return ok
}

func multiplyByTwo(v types.Value) types.Value {
	switch val := v.(type) {
	case *types.IntValue:
		return types.NewInt(val.Value() * 2)
	case *types.FloatValue:
		return types.NewFloat(val.Value() * 2.0)
	}
	return v
}

func square(v types.Value) types.Value {
	switch val := v.(type) {
	case *types.IntValue:
		x := val.Value()
		return types.NewInt(x * x)
	case *types.FloatValue:
		x := val.Value()
		return types.NewFloat(x * x)
	}
	return v
}

func absolute(v types.Value) types.Value {
	switch val := v.(type) {
	case *types.IntValue:
		x := val.Value()
		if x < 0 {
			return types.NewInt(-x)
		}
		return types.NewInt(x)
	case *types.FloatValue:
		x := val.Value()
		if x < 0 {
			return types.NewFloat(-x)
		}
		return types.NewFloat(x)
	}
	return v
}

func toUpper(v types.Value) types.Value {
	if strVal, ok := v.(*types.StringValue); ok {
		s := strVal.Value()
		upper := ""
		for _, r := range s {
			if r >= 'a' && r <= 'z' {
				upper += string(r - 32)
			} else {
				upper += string(r)
			}
		}
		return types.NewString(upper)
	}
	return v
}

func toLower(v types.Value) types.Value {
	if strVal, ok := v.(*types.StringValue); ok {
		s := strVal.Value()
		lower := ""
		for _, r := range s {
			if r >= 'A' && r <= 'Z' {
				lower += string(r + 32)
			} else {
				lower += string(r)
			}
		}
		return types.NewString(lower)
	}
	return v
}

func getLength(v types.Value) types.Value {
	switch val := v.(type) {
	case *types.StringValue:
		return types.NewInt(int64(len(val.Value())))
	case *types.SliceValue:
		return types.NewInt(int64(val.Len()))
	}
	return types.NewInt(0)
}

func maxValue(a, b types.Value) types.Value {
	if compareValues(a, b) > 0 {
		return a
	}
	return b
}

func minValue(a, b types.Value) types.Value {
	if compareValues(a, b) < 0 {
		return a
	}
	return b
}
