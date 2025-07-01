package builtins

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mredencom/expr/types"
)

// BuiltinFunction represents a builtin function
type BuiltinFunction func(args []types.Value) (types.Value, error)

// StandardBuiltinNames defines the standard order of builtin functions for indexing
var StandardBuiltinNames = []string{
	// Core builtins
	"len", "string", "int", "float", "bool", "abs", "max", "min",
	"contains", "startsWith", "endsWith", "upper", "lower", "trim", "type",

	// String functions
	"replace", "substring", "indexOf",

	// Math functions
	"ceil", "floor", "round", "sqrt", "pow",

	// Time functions
	"now",

	// Collection functions
	"flatten", "groupBy",

	// Pipeline functions - Collection processing
	"filter", "map", "reduce", "sort", "reverse", "take", "skip", "unique",

	// Pipeline functions - Aggregation
	"count", "sum", "avg",

	// Pipeline functions - String processing
	"split", "join", "match",

	// Pipeline functions - Utility
	"debug", "pipe",

	// Legacy names for compatibility
	"matches", "all", "any", "first", "last", "keys",
}

// AllBuiltins contains all available builtin functions
var AllBuiltins = map[string]BuiltinFunction{
	// Core builtins
	"len":        lenBuiltin,
	"string":     stringBuiltin,
	"int":        intBuiltin,
	"float":      floatBuiltin,
	"bool":       boolBuiltin,
	"abs":        absBuiltin,
	"max":        maxBuiltin,
	"min":        minBuiltin,
	"contains":   containsBuiltin,
	"startsWith": startsWithBuiltin,
	"endsWith":   endsWithBuiltin,
	"upper":      upperBuiltin,
	"lower":      lowerBuiltin,
	"trim":       trimBuiltin,
	"type":       typeBuiltin,

	// String functions
	"replace":   replaceBuiltin,
	"substring": substringBuiltin,
	"indexOf":   indexOfBuiltin,

	// Math functions
	"ceil":  ceilBuiltin,
	"floor": floorBuiltin,
	"round": roundBuiltin,
	"sqrt":  sqrtBuiltin,
	"pow":   powBuiltin,

	// Time functions
	"now": nowBuiltin,

	// Collection functions
	"flatten": flattenBuiltin,
	"groupBy": groupByBuiltin,

	// Pipeline functions - Collection processing
	"filter":  filterFunc,
	"map":     mapFunc,
	"reduce":  reduceFunc,
	"sort":    sortFunc,
	"reverse": reverseFunc,
	"take":    takeFunc,
	"skip":    skipFunc,
	"unique":  uniqueFunc,

	// Pipeline functions - Aggregation
	"count": countFunc,
	"sum":   sumBuiltin,
	"avg":   avgFunc,

	// Pipeline functions - String processing
	"split": splitFunc,
	"join":  joinFunc,
	"match": matchFunc,

	// Pipeline functions - Type conversion (avoiding duplicates)
	// "string": stringFunc, // Already defined above as stringBuiltin
	// "int":    intFunc,    // Already defined above as intBuiltin
	// "float":  floatFunc,  // Already defined above as floatBuiltin
	// "bool":   boolFunc,   // Already defined above as boolBuiltin

	// Pipeline functions - Utility
	"debug": debugFunc,
	"pipe":  pipeFunc,

	// Collection access
	"first": firstBuiltin,
	"last":  lastBuiltin,

	// String matching
	"matches": matchesBuiltin,
}

// lenBuiltin returns the length of a value
func lenBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.StringValue:
		return types.NewInt(int64(len(v.Value()))), nil
	case *types.SliceValue:
		return types.NewInt(int64(v.Len())), nil
	case *types.MapValue:
		return types.NewInt(int64(v.Len())), nil
	default:
		return nil, fmt.Errorf("object of type %T has no len()", arg)
	}
}

// stringBuiltin converts a value to string
func stringBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string() takes exactly 1 argument, got %d", len(args))
	}

	return types.NewString(args[0].String()), nil
}

// intBuiltin converts a value to integer
func intBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.IntValue:
		return v, nil
	case *types.FloatValue:
		return types.NewInt(int64(v.Value())), nil
	case *types.StringValue:
		if i, err := strconv.ParseInt(v.Value(), 10, 64); err == nil {
			return types.NewInt(i), nil
		}
		return nil, fmt.Errorf("invalid literal for int(): %s", v.Value())
	case *types.BoolValue:
		if v.Value() {
			return types.NewInt(1), nil
		}
		return types.NewInt(0), nil
	default:
		return nil, fmt.Errorf("int() argument must be a string, number or bool, not %T", arg)
	}
}

// floatBuiltin converts a value to float
func floatBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.FloatValue:
		return v, nil
	case *types.IntValue:
		return types.NewFloat(float64(v.Value())), nil
	case *types.StringValue:
		if f, err := strconv.ParseFloat(v.Value(), 64); err == nil {
			return types.NewFloat(f), nil
		}
		return nil, fmt.Errorf("invalid literal for float(): %s", v.Value())
	default:
		return nil, fmt.Errorf("float() argument must be a string or number, not %T", arg)
	}
}

// boolBuiltin converts a value to boolean
func boolBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
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

// absBuiltin returns absolute value
func absBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.IntValue:
		val := v.Value()
		if val < 0 {
			return types.NewInt(-val), nil
		}
		return v, nil
	case *types.FloatValue:
		val := v.Value()
		if val < 0 {
			return types.NewFloat(-val), nil
		}
		return v, nil
	default:
		return nil, fmt.Errorf("abs() argument must be a number, not %T", arg)
	}
}

// maxBuiltin returns maximum value
func maxBuiltin(args []types.Value) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("max() expected at least 1 argument, got 0")
	}

	// If single argument and it's an array, find max in array
	if len(args) == 1 {
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
	}

	// Multiple arguments - find max among them
	max := args[0]
	for i := 1; i < len(args); i++ {
		if compareValues(args[i], max) > 0 {
			max = args[i]
		}
	}
	return max, nil
}

// minBuiltin returns minimum value
func minBuiltin(args []types.Value) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() expected at least 1 argument, got 0")
	}

	// If single argument and it's an array, find min in array
	if len(args) == 1 {
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
	}

	// Multiple arguments - find min among them
	min := args[0]
	for i := 1; i < len(args); i++ {
		if compareValues(args[i], min) < 0 {
			min = args[i]
		}
	}
	return min, nil
}

// sumBuiltin returns sum of values
func sumBuiltin(args []types.Value) (types.Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("sum() expected at least 1 argument, got 0")
	}

	// If single argument and it's an array, sum the array elements
	if len(args) == 1 {
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
	}

	// Multiple arguments - sum them directly
	var intSum int64
	var floatSum float64
	hasFloat := false

	for _, arg := range args {
		switch v := arg.(type) {
		case *types.IntValue:
			intSum += v.Value()
		case *types.FloatValue:
			floatSum += v.Value()
			hasFloat = true
		default:
			return nil, fmt.Errorf("sum() arguments must be numbers, got %T", arg)
		}
	}

	if hasFloat {
		return types.NewFloat(floatSum + float64(intSum)), nil
	}
	return types.NewInt(intSum), nil
}

// containsBuiltin checks if string contains substring
func containsBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains() takes exactly 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("contains() first argument must be string, not %T", args[0])
	}

	substr, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("contains() second argument must be string, not %T", args[1])
	}

	result := strings.Contains(str.Value(), substr.Value())
	return types.NewBool(result), nil
}

// startsWithBuiltin checks if string starts with prefix
func startsWithBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("startsWith() takes exactly 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("startsWith() first argument must be string, not %T", args[0])
	}

	prefix, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("startsWith() second argument must be string, not %T", args[1])
	}

	result := strings.HasPrefix(str.Value(), prefix.Value())
	return types.NewBool(result), nil
}

// endsWithBuiltin checks if string ends with suffix
func endsWithBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("endsWith() takes exactly 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("endsWith() first argument must be string, not %T", args[0])
	}

	suffix, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("endsWith() second argument must be string, not %T", args[1])
	}

	result := strings.HasSuffix(str.Value(), suffix.Value())
	return types.NewBool(result), nil
}

// upperBuiltin converts string to uppercase
func upperBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper() takes exactly 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("upper() argument must be string, not %T", args[0])
	}

	return types.NewString(strings.ToUpper(str.Value())), nil
}

// lowerBuiltin converts string to lowercase
func lowerBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower() takes exactly 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("lower() argument must be string, not %T", args[0])
	}

	return types.NewString(strings.ToLower(str.Value())), nil
}

// trimBuiltin trims whitespace from string
func trimBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim() takes exactly 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("trim() argument must be string, not %T", args[0])
	}

	return types.NewString(strings.TrimSpace(str.Value())), nil
}

// typeBuiltin returns the type name of a value
func typeBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type() takes exactly 1 argument, got %d", len(args))
	}

	return types.NewString(args[0].Type().Name), nil
}

// ListBuiltinNames returns a sorted list of all builtin function names
func ListBuiltinNames() []string {
	names := make([]string, 0, len(AllBuiltins))
	for name := range AllBuiltins {
		names = append(names, name)
	}
	return names
}

// compareValues compares two values, returns -1, 0, or 1
func compareValues(a, b types.Value) int {
	switch va := a.(type) {
	case *types.IntValue:
		if vb, ok := b.(*types.IntValue); ok {
			if va.Value() < vb.Value() {
				return -1
			} else if va.Value() > vb.Value() {
				return 1
			}
			return 0
		}
	case *types.FloatValue:
		if vb, ok := b.(*types.FloatValue); ok {
			if va.Value() < vb.Value() {
				return -1
			} else if va.Value() > vb.Value() {
				return 1
			}
			return 0
		}
	case *types.StringValue:
		if vb, ok := b.(*types.StringValue); ok {
			if va.Value() < vb.Value() {
				return -1
			} else if va.Value() > vb.Value() {
				return 1
			}
			return 0
		}
	}
	return 0
}

// firstBuiltin returns the first element of a collection
func firstBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.SliceValue:
		if v.Len() == 0 {
			return types.NewNil(), nil
		}
		return v.Get(0), nil
	case *types.StringValue:
		str := v.Value()
		if len(str) == 0 {
			return types.NewString(""), nil
		}
		return types.NewString(string(str[0])), nil
	default:
		return nil, fmt.Errorf("first() can only be applied to arrays or strings, got %T", arg)
	}
}

// lastBuiltin returns the last element of a collection
func lastBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("last() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.SliceValue:
		if v.Len() == 0 {
			return types.NewNil(), nil
		}
		return v.Get(v.Len() - 1), nil
	case *types.StringValue:
		str := v.Value()
		if len(str) == 0 {
			return types.NewString(""), nil
		}
		return types.NewString(string(str[len(str)-1])), nil
	default:
		return nil, fmt.Errorf("last() can only be applied to arrays or strings, got %T", arg)
	}
}

// matchesBuiltin checks if a string matches a regular expression pattern
func matchesBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("matches() takes exactly 2 arguments, got %d", len(args))
	}

	strVal, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("matches() first argument must be a string, got %T", args[0])
	}

	patternVal, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("matches() second argument must be a string, got %T", args[1])
	}

	// Use Go's regexp package for real regex matching
	matched, err := regexp.MatchString(patternVal.Value(), strVal.Value())
	if err != nil {
		return nil, fmt.Errorf("matches() invalid regex pattern: %v", err)
	}
	return types.NewBool(matched), nil
}

// replaceBuiltin replaces all occurrences of a substring in a string
func replaceBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace() takes exactly 3 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() first argument must be string, not %T", args[0])
	}

	old, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() second argument must be string, not %T", args[1])
	}

	new, ok := args[2].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() third argument must be string, not %T", args[2])
	}

	result := strings.ReplaceAll(str.Value(), old.Value(), new.Value())
	return types.NewString(result), nil
}

// substringBuiltin extracts a substring from a string
func substringBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("substring() takes exactly 3 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("substring() first argument must be string, not %T", args[0])
	}

	start, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("substring() second argument must be int, not %T", args[1])
	}

	end, ok := args[2].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("substring() third argument must be int, not %T", args[2])
	}

	s := str.Value()
	startIdx := int(start.Value())
	endIdx := int(end.Value())

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(s) {
		endIdx = len(s)
	}
	if startIdx > endIdx {
		startIdx = endIdx
	}

	result := s[startIdx:endIdx]
	return types.NewString(result), nil
}

// indexOfBuiltin finds the index of a substring in a string
func indexOfBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("indexOf() takes exactly 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("indexOf() first argument must be string, not %T", args[0])
	}

	substr, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("indexOf() second argument must be string, not %T", args[1])
	}

	index := strings.Index(str.Value(), substr.Value())
	return types.NewInt(int64(index)), nil
}

// ceilBuiltin returns the ceiling of a number
func ceilBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.IntValue:
		return v, nil // Integers are already whole numbers
	case *types.FloatValue:
		return types.NewFloat(math.Ceil(v.Value())), nil
	default:
		return nil, fmt.Errorf("ceil() argument must be a number, not %T", arg)
	}
}

// floorBuiltin returns the floor of a number
func floorBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.IntValue:
		return v, nil // Integers are already whole numbers
	case *types.FloatValue:
		return types.NewFloat(math.Floor(v.Value())), nil
	default:
		return nil, fmt.Errorf("floor() argument must be a number, not %T", arg)
	}
}

// roundBuiltin rounds a number to the nearest integer
func roundBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("round() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	switch v := arg.(type) {
	case *types.IntValue:
		return v, nil // Integers are already whole numbers
	case *types.FloatValue:
		return types.NewFloat(math.Round(v.Value())), nil
	default:
		return nil, fmt.Errorf("round() argument must be a number, not %T", arg)
	}
}

// sqrtBuiltin returns the square root of a number
func sqrtBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt() takes exactly 1 argument, got %d", len(args))
	}

	arg := args[0]
	var value float64
	switch v := arg.(type) {
	case *types.IntValue:
		value = float64(v.Value())
	case *types.FloatValue:
		value = v.Value()
	default:
		return nil, fmt.Errorf("sqrt() argument must be a number, not %T", arg)
	}

	if value < 0 {
		return nil, fmt.Errorf("sqrt() argument must be non-negative")
	}

	return types.NewFloat(math.Sqrt(value)), nil
}

// powBuiltin returns the power of a number
func powBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow() takes exactly 2 arguments, got %d", len(args))
	}

	var base, exp float64

	switch v := args[0].(type) {
	case *types.IntValue:
		base = float64(v.Value())
	case *types.FloatValue:
		base = v.Value()
	default:
		return nil, fmt.Errorf("pow() first argument must be a number, not %T", args[0])
	}

	switch v := args[1].(type) {
	case *types.IntValue:
		exp = float64(v.Value())
	case *types.FloatValue:
		exp = v.Value()
	default:
		return nil, fmt.Errorf("pow() second argument must be a number, not %T", args[1])
	}

	result := math.Pow(base, exp)
	return types.NewFloat(result), nil
}

// nowBuiltin returns the current time as Unix timestamp
func nowBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("now() takes no arguments, got %d", len(args))
	}

	return types.NewInt(time.Now().Unix()), nil
}

// flattenBuiltin flattens a nested array
func flattenBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("flatten() takes exactly 1 argument, got %d", len(args))
	}

	arg, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("flatten() argument must be an array, not %T", args[0])
	}

	var result []types.Value
	for i := 0; i < arg.Len(); i++ {
		element := arg.Get(i)
		if slice, ok := element.(*types.SliceValue); ok {
			// Recursively flatten nested arrays
			for j := 0; j < slice.Len(); j++ {
				result = append(result, slice.Get(j))
			}
		} else {
			result = append(result, element)
		}
	}

	// Use a generic type for the flattened result
	elemType := types.TypeInfo{Kind: types.KindUnknown, Name: "any", Size: -1}
	return types.NewSlice(result, elemType), nil
}

// groupByBuiltin groups elements by a given condition (simplified version)
func groupByBuiltin(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("groupBy() takes exactly 2 arguments, got %d", len(args))
	}

	data, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("groupBy() first argument must be an array, not %T", args[0])
	}

	// For now, simplified implementation that groups by boolean condition
	// This is a basic implementation - a full version would need lambda support
	trueGroup := make([]types.Value, 0)
	falseGroup := make([]types.Value, 0)

	for i := 0; i < data.Len(); i++ {
		element := data.Get(i)
		// For simplicity, group by whether the element is truthy
		if isTruthy(element) {
			trueGroup = append(trueGroup, element)
		} else {
			falseGroup = append(falseGroup, element)
		}
	}

	// Return a map with "true" and "false" keys
	elemType := types.TypeInfo{Kind: types.KindUnknown, Name: "any", Size: -1}
	result := map[string]types.Value{
		"true":  types.NewSlice(trueGroup, elemType),
		"false": types.NewSlice(falseGroup, elemType),
	}
	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	valType := types.TypeInfo{Kind: types.KindSlice, Name: "slice", Size: -1}
	return types.NewMap(result, keyType, valType), nil
}

// isTruthy function is defined in collections.go, removing duplicate
