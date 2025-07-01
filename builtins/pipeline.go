package builtins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mredencom/expr/types"
)

// PipelineBuiltins contains all pipeline-oriented builtin functions
var PipelineBuiltins = map[string]BuiltinFunction{
	// Collection processing - Enhanced versions with Lambda and placeholder support
	"filter":  EnhancedFilterFunc,
	"map":     EnhancedMapFunc,
	"reduce":  EnhancedReduceFunc,
	"sort":    sortFunc,
	"reverse": reverseFunc,
	"take":    takeFunc,
	"skip":    skipFunc,
	"unique":  uniqueFunc,

	// Aggregation functions
	"count": countFunc,
	"sum":   sumFunc,
	"avg":   avgFunc,
	"min":   minFunc,
	"max":   maxFunc,

	// String processing
	"split": splitFunc,
	"join":  joinFunc,
	"trim":  trimFunc,
	"upper": upperFunc,
	"lower": lowerFunc,
	"match": matchFunc,

	// Type checking and conversion
	"type":   typeFunc,
	"string": stringFunc,
	"int":    intFunc,
	"float":  floatFunc,
	"bool":   boolFunc,

	// Utility functions
	"debug": debugFunc,
	"pipe":  pipeFunc,
}

// Collection processing functions

func filterFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("filter requires exactly 2 arguments")
	}

	collection := args[0]
	predicate := args[1]

	if array, ok := collection.(*types.SliceValue); ok {
		var result []types.Value
		for i := 0; i < array.Len(); i++ {
			item := array.Get(i)
			// Apply predicate (for now, simplified)
			if shouldInclude(item, predicate) {
				result = append(result, item)
			}
		}
		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("filter can only be applied to arrays")
}

func mapFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map requires exactly 2 arguments")
	}

	collection := args[0]
	transformer := args[1]

	if array, ok := collection.(*types.SliceValue); ok {
		var result []types.Value
		for i := 0; i < array.Len(); i++ {
			item := array.Get(i)
			// Apply transformer (for now, simplified)
			transformed := applyTransform(item, transformer)
			result = append(result, transformed)
		}
		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("map can only be applied to arrays")
}

func reduceFunc(args []types.Value) (types.Value, error) {
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
			acc = applyReducer(acc, array.Get(i), reducer)
		}

		return acc, nil
	}

	return nil, fmt.Errorf("reduce can only be applied to arrays")
}

func sortFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sort requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		// Simple numeric/string sorting for now
		result := make([]types.Value, array.Len())
		for i := 0; i < array.Len(); i++ {
			result[i] = array.Get(i)
		}

		// Basic bubble sort implementation for demonstration
		for i := 0; i < len(result); i++ {
			for j := 0; j < len(result)-i-1; j++ {
				if compareValues(result[j], result[j+1]) > 0 {
					result[j], result[j+1] = result[j+1], result[j]
				}
			}
		}

		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("sort can only be applied to arrays")
}

func reverseFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("reverse requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		result := make([]types.Value, array.Len())
		for i := 0; i < array.Len(); i++ {
			result[array.Len()-1-i] = array.Get(i)
		}
		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("reverse can only be applied to arrays")
}

func takeFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("take requires exactly 2 arguments")
	}

	array, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("take can only be applied to arrays")
	}

	count, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("take count must be an integer")
	}

	n := int(count.Value())
	if n < 0 {
		n = 0
	}
	if n > array.Len() {
		n = array.Len()
	}

	result := make([]types.Value, n)
	for i := 0; i < n; i++ {
		result[i] = array.Get(i)
	}

	return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
}

func skipFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("skip requires exactly 2 arguments")
	}

	array, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("skip can only be applied to arrays")
	}

	count, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("skip count must be an integer")
	}

	n := int(count.Value())
	if n < 0 {
		n = 0
	}
	if n >= array.Len() {
		return types.NewSlice([]types.Value{}, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	result := make([]types.Value, array.Len()-n)
	for i := n; i < array.Len(); i++ {
		result[i-n] = array.Get(i)
	}

	return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
}

func uniqueFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unique requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		seen := make(map[string]bool)
		var result []types.Value

		for i := 0; i < array.Len(); i++ {
			elem := array.Get(i)
			key := elem.String()
			if !seen[key] {
				seen[key] = true
				result = append(result, elem)
			}
		}

		return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
	}

	return nil, fmt.Errorf("unique can only be applied to arrays")
}

// Aggregation functions

func countFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("count requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		return types.NewInt(int64(array.Len())), nil
	}

	if str, ok := args[0].(*types.StringValue); ok {
		return types.NewInt(int64(len(str.Value()))), nil
	}

	return nil, fmt.Errorf("count can only be applied to arrays or strings")
}

func sumFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sum requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		var intSum int64
		var floatSum float64
		hasFloat := false

		for i := 0; i < array.Len(); i++ {
			elem := array.Get(i)
			switch v := elem.(type) {
			case *types.IntValue:
				intSum += v.Value()
			case *types.FloatValue:
				floatSum += v.Value()
				hasFloat = true
			default:
				return nil, fmt.Errorf("sum can only be applied to numeric arrays")
			}
		}

		if hasFloat {
			return types.NewFloat(floatSum + float64(intSum)), nil
		}
		return types.NewInt(intSum), nil
	}

	return nil, fmt.Errorf("sum can only be applied to arrays")
}

func avgFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("avg requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		if array.Len() == 0 {
			return nil, fmt.Errorf("cannot calculate average of empty array")
		}

		sum, err := sumFunc(args)
		if err != nil {
			return nil, err
		}

		count := float64(array.Len())

		if intSum, ok := sum.(*types.IntValue); ok {
			return types.NewFloat(float64(intSum.Value()) / count), nil
		}
		if floatSum, ok := sum.(*types.FloatValue); ok {
			return types.NewFloat(floatSum.Value() / count), nil
		}
	}

	return nil, fmt.Errorf("avg can only be applied to numeric arrays")
}

func minFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("min requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		if array.Len() == 0 {
			return nil, fmt.Errorf("cannot find minimum of empty array")
		}

		min := array.Get(0)
		for i := 1; i < array.Len(); i++ {
			elem := array.Get(i)
			if compareValues(elem, min) < 0 {
				min = elem
			}
		}
		return min, nil
	}

	return nil, fmt.Errorf("min can only be applied to arrays")
}

func maxFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("max requires exactly 1 argument")
	}

	if array, ok := args[0].(*types.SliceValue); ok {
		if array.Len() == 0 {
			return nil, fmt.Errorf("cannot find maximum of empty array")
		}

		max := array.Get(0)
		for i := 1; i < array.Len(); i++ {
			elem := array.Get(i)
			if compareValues(elem, max) > 0 {
				max = elem
			}
		}
		return max, nil
	}

	return nil, fmt.Errorf("max can only be applied to arrays")
}

// String processing functions

func splitFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split requires exactly 2 arguments")
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("split first argument must be a string")
	}

	separator, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("split separator must be a string")
	}

	parts := strings.Split(str.Value(), separator.Value())
	result := make([]types.Value, len(parts))
	for i, part := range parts {
		result[i] = types.NewString(part)
	}

	return types.NewSlice(result, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
}

func joinFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join requires exactly 2 arguments")
	}

	array, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("join first argument must be an array")
	}

	separator, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("join separator must be a string")
	}

	var parts []string
	for i := 0; i < array.Len(); i++ {
		parts = append(parts, array.Get(i).String())
	}

	return types.NewString(strings.Join(parts, separator.Value())), nil
}

func trimFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim requires exactly 1 argument")
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("trim can only be applied to strings")
	}

	return types.NewString(strings.TrimSpace(str.Value())), nil
}

func upperFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper requires exactly 1 argument")
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("upper can only be applied to strings")
	}

	return types.NewString(strings.ToUpper(str.Value())), nil
}

func lowerFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower requires exactly 1 argument")
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("lower can only be applied to strings")
	}

	return types.NewString(strings.ToLower(str.Value())), nil
}

func matchFunc(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("match requires exactly 2 arguments")
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("match first argument must be a string")
	}

	pattern, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("match pattern must be a string")
	}

	matched, err := regexp.MatchString(pattern.Value(), str.Value())
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	return types.NewBool(matched), nil
}

// Type checking and conversion functions

func typeFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type requires exactly 1 argument")
	}

	return types.NewString(args[0].Type().Name), nil
}

func stringFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string requires exactly 1 argument")
	}

	return types.NewString(args[0].String()), nil
}

func intFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int requires exactly 1 argument")
	}

	switch v := args[0].(type) {
	case *types.IntValue:
		return v, nil
	case *types.FloatValue:
		return types.NewInt(int64(v.Value())), nil
	case *types.StringValue:
		// Try to parse the string as an integer
		return nil, fmt.Errorf("string to int conversion not implemented")
	case *types.BoolValue:
		if v.Value() {
			return types.NewInt(1), nil
		}
		return types.NewInt(0), nil
	default:
		return nil, fmt.Errorf("cannot convert to int")
	}
}

func floatFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float requires exactly 1 argument")
	}

	switch v := args[0].(type) {
	case *types.FloatValue:
		return v, nil
	case *types.IntValue:
		return types.NewFloat(float64(v.Value())), nil
	case *types.StringValue:
		return nil, fmt.Errorf("string to float conversion not implemented")
	case *types.BoolValue:
		if v.Value() {
			return types.NewFloat(1.0), nil
		}
		return types.NewFloat(0.0), nil
	default:
		return nil, fmt.Errorf("cannot convert to float")
	}
}

func boolFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool requires exactly 1 argument")
	}

	switch v := args[0].(type) {
	case *types.BoolValue:
		return v, nil
	case *types.IntValue:
		return types.NewBool(v.Value() != 0), nil
	case *types.FloatValue:
		return types.NewBool(v.Value() != 0.0), nil
	case *types.StringValue:
		return types.NewBool(v.Value() != ""), nil
	case *types.NilValue:
		return types.NewBool(false), nil
	default:
		return types.NewBool(true), nil
	}
}

// Utility functions

func debugFunc(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("debug requires exactly 1 argument")
	}

	fmt.Printf("DEBUG: %s (type: %s)\n", args[0].String(), args[0].Type().Name)
	return args[0], nil
}

func pipeFunc(args []types.Value) (types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("pipe requires at least 2 arguments")
	}

	// Simple pipe implementation - apply functions in sequence
	result := args[0]
	for i := 1; i < len(args); i++ {
		if fn, ok := args[i].(*types.FuncValue); ok {
			result = applyFunction(result, fn)
		} else {
			return nil, fmt.Errorf("pipe arguments must be functions")
		}
	}

	return result, nil
}

// Helper functions

func shouldInclude(item types.Value, predicate types.Value) bool {
	// Apply predicate function to item
	switch pred := predicate.(type) {
	case *types.FuncValue:
		// Execute lambda function with item as parameter
		result := evaluateLambda(pred, item)
		return isTruthy(result)
	case *types.StringValue:
		// String predicates for simple filtering
		return evaluateStringPredicate(pred.Value(), item)
	case *types.BoolValue:
		// Boolean predicate - constant filter
		return pred.Value()
	default:
		// Default: include all items
		return true
	}
}

func applyTransform(item types.Value, transformer types.Value) types.Value {
	// Apply transformer function to item
	switch trans := transformer.(type) {
	case *types.FuncValue:
		// Execute lambda function with item as parameter
		return evaluateLambda(trans, item)
	case *types.StringValue:
		// String transformers for simple operations
		return applyStringTransform(trans.Value(), item)
	default:
		// Default: return item unchanged
		return item
	}
}

func applyReducer(acc types.Value, item types.Value, reducer types.Value) types.Value {
	// Apply reducer function to accumulator and item
	switch red := reducer.(type) {
	case *types.FuncValue:
		// Execute lambda function with acc and item as parameters
		return evaluateLambdaWithTwoArgs(red, acc, item)
	case *types.StringValue:
		// String reducers for simple operations
		return applyStringReducer(red.Value(), acc, item)
	default:
		// Default: return accumulator
		return acc
	}
}

// evaluateLambda executes a lambda function with a single argument
func evaluateLambda(lambda *types.FuncValue, arg types.Value) types.Value {
	// This is a simplified lambda evaluation
	// In a full implementation, we would create a new VM context

	params := lambda.Parameters()
	if len(params) == 0 {
		return arg
	}

	// For basic lambda expressions, try to evaluate them directly
	body := lambda.Body()

	// Handle simple variable reference
	if bodyStr, ok := body.(string); ok {
		if bodyStr == params[0] {
			// Simple identity function: x => x
			return arg
		}
	}

	// For complex expressions, return a placeholder result
	// In a real implementation, this would execute the lambda body
	return arg
}

// evaluateLambdaWithTwoArgs executes a lambda function with two arguments
func evaluateLambdaWithTwoArgs(lambda *types.FuncValue, arg1, arg2 types.Value) types.Value {
	// Simplified implementation for two-argument lambdas
	params := lambda.Parameters()
	if len(params) < 2 {
		return arg1
	}

	// For now, return first argument as default
	return arg1
}

// evaluateStringPredicate evaluates simple string-based predicates
func evaluateStringPredicate(predicate string, item types.Value) bool {
	switch predicate {
	case "not_empty", "non_empty":
		return !isZeroValue(item)
	case "positive":
		if intVal, ok := item.(*types.IntValue); ok {
			return intVal.Value() > 0
		}
		if floatVal, ok := item.(*types.FloatValue); ok {
			return floatVal.Value() > 0
		}
		return false
	case "negative":
		if intVal, ok := item.(*types.IntValue); ok {
			return intVal.Value() < 0
		}
		if floatVal, ok := item.(*types.FloatValue); ok {
			return floatVal.Value() < 0
		}
		return false
	case "even":
		if intVal, ok := item.(*types.IntValue); ok {
			return intVal.Value()%2 == 0
		}
		return false
	case "odd":
		if intVal, ok := item.(*types.IntValue); ok {
			return intVal.Value()%2 != 0
		}
		return false
	default:
		return true
	}
}

// applyStringTransform applies simple string-based transformations
func applyStringTransform(transform string, item types.Value) types.Value {
	switch transform {
	case "double", "times2":
		if intVal, ok := item.(*types.IntValue); ok {
			return types.NewInt(intVal.Value() * 2)
		}
		if floatVal, ok := item.(*types.FloatValue); ok {
			return types.NewFloat(floatVal.Value() * 2.0)
		}
		return item
	case "square":
		if intVal, ok := item.(*types.IntValue); ok {
			val := intVal.Value()
			return types.NewInt(val * val)
		}
		if floatVal, ok := item.(*types.FloatValue); ok {
			val := floatVal.Value()
			return types.NewFloat(val * val)
		}
		return item
	case "upper":
		if strVal, ok := item.(*types.StringValue); ok {
			return types.NewString(strings.ToUpper(strVal.Value()))
		}
		return item
	case "lower":
		if strVal, ok := item.(*types.StringValue); ok {
			return types.NewString(strings.ToLower(strVal.Value()))
		}
		return item
	default:
		return item
	}
}

// applyStringReducer applies simple string-based reducers
func applyStringReducer(reducer string, acc, item types.Value) types.Value {
	switch reducer {
	case "add", "sum":
		if accInt, ok := acc.(*types.IntValue); ok {
			if itemInt, ok := item.(*types.IntValue); ok {
				return types.NewInt(accInt.Value() + itemInt.Value())
			}
		}
		if accFloat, ok := acc.(*types.FloatValue); ok {
			if itemFloat, ok := item.(*types.FloatValue); ok {
				return types.NewFloat(accFloat.Value() + itemFloat.Value())
			}
		}
		return acc
	case "multiply", "mul":
		if accInt, ok := acc.(*types.IntValue); ok {
			if itemInt, ok := item.(*types.IntValue); ok {
				return types.NewInt(accInt.Value() * itemInt.Value())
			}
		}
		if accFloat, ok := acc.(*types.FloatValue); ok {
			if itemFloat, ok := item.(*types.FloatValue); ok {
				return types.NewFloat(accFloat.Value() * itemFloat.Value())
			}
		}
		return acc
	case "concat":
		return types.NewString(acc.String() + item.String())
	default:
		return acc
	}
}

func applyFunction(value types.Value, fn *types.FuncValue) types.Value {
	// Apply function to value - improved implementation
	return evaluateLambda(fn, value)
}
