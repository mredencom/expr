package builtins

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mredencom/expr/types"
)

// All type method functions are defined below

// String type methods implementation

func stringLengthMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.length() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.length() requires a string argument")
	}
	return types.NewInt(int64(len(str.Value()))), nil
}

func stringCharAtMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.charAt() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.charAt() requires a string argument")
	}
	index, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.charAt() requires an int index")
	}

	s := str.Value()
	idx := int(index.Value())
	if idx < 0 || idx >= len(s) {
		return types.NewString(""), nil
	}
	return types.NewString(string(s[idx])), nil
}

func stringCharCodeAtMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.charCodeAt() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.charCodeAt() requires a string argument")
	}
	index, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.charCodeAt() requires an int index")
	}

	s := str.Value()
	idx := int(index.Value())
	if idx < 0 || idx >= len(s) {
		return types.NewInt(-1), nil
	}
	return types.NewInt(int64(s[idx])), nil
}

func stringSliceMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("string.slice() takes 2 or 3 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.slice() requires a string argument")
	}
	start, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.slice() requires an int start index")
	}

	s := str.Value()
	startIdx := int(start.Value())
	endIdx := len(s)

	if len(args) == 3 {
		end, ok := args[2].(*types.IntValue)
		if !ok {
			return nil, fmt.Errorf("string.slice() requires an int end index")
		}
		endIdx = int(end.Value())
	}

	// Handle negative indices
	if startIdx < 0 {
		startIdx = len(s) + startIdx
	}
	if endIdx < 0 {
		endIdx = len(s) + endIdx
	}

	// Clamp indices
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(s) {
		endIdx = len(s)
	}
	if startIdx > endIdx {
		startIdx = endIdx
	}

	return types.NewString(s[startIdx:endIdx]), nil
}

func stringSplitMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("string.split() takes 2 or 3 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.split() requires a string argument")
	}
	separator, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.split() requires a string separator")
	}

	limit := -1
	if len(args) == 3 {
		limitVal, ok := args[2].(*types.IntValue)
		if !ok {
			return nil, fmt.Errorf("string.split() requires an int limit")
		}
		limit = int(limitVal.Value())
	}

	parts := strings.Split(str.Value(), separator.Value())
	if limit > 0 && len(parts) > limit {
		parts = parts[:limit]
	}

	result := make([]types.Value, len(parts))
	for i, part := range parts {
		result[i] = types.NewString(part)
	}

	elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewSlice(result, elemType), nil
}

func stringJoinMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.join() takes exactly 2 arguments")
	}
	separator, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.join() requires a string separator")
	}
	arr, ok := args[1].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("string.join() requires an array argument")
	}

	parts := make([]string, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		parts[i] = arr.Get(i).String()
	}

	return types.NewString(strings.Join(parts, separator.Value())), nil
}

func stringReplaceMethod(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("string.replace() takes exactly 3 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.replace() requires a string argument")
	}
	old, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.replace() requires a string to replace")
	}
	new, ok := args[2].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.replace() requires a replacement string")
	}

	result := strings.ReplaceAll(str.Value(), old.Value(), new.Value())
	return types.NewString(result), nil
}

func stringTrimMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.trim() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.trim() requires a string argument")
	}
	return types.NewString(strings.TrimSpace(str.Value())), nil
}

func stringTrimLeftMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.trimLeft() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.trimLeft() requires a string argument")
	}
	return types.NewString(strings.TrimLeftFunc(str.Value(), func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})), nil
}

func stringTrimRightMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.trimRight() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.trimRight() requires a string argument")
	}
	return types.NewString(strings.TrimRightFunc(str.Value(), func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})), nil
}

func stringUpperMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.upper() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.upper() requires a string argument")
	}
	return types.NewString(strings.ToUpper(str.Value())), nil
}

func stringLowerMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.lower() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.lower() requires a string argument")
	}
	return types.NewString(strings.ToLower(str.Value())), nil
}

func stringStartsWithMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.startsWith() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.startsWith() requires a string argument")
	}
	prefix, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.startsWith() requires a string prefix")
	}
	return types.NewBool(strings.HasPrefix(str.Value(), prefix.Value())), nil
}

func stringEndsWithMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.endsWith() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.endsWith() requires a string argument")
	}
	suffix, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.endsWith() requires a string suffix")
	}
	return types.NewBool(strings.HasSuffix(str.Value(), suffix.Value())), nil
}

func stringContainsMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.contains() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.contains() requires a string argument")
	}
	substr, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.contains() requires a string to search for")
	}
	return types.NewBool(strings.Contains(str.Value(), substr.Value())), nil
}

func stringIndexOfMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.indexOf() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.indexOf() requires a string argument")
	}
	substr, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.indexOf() requires a string to search for")
	}
	index := strings.Index(str.Value(), substr.Value())
	return types.NewInt(int64(index)), nil
}

func stringLastIndexOfMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.lastIndexOf() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.lastIndexOf() requires a string argument")
	}
	substr, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.lastIndexOf() requires a string to search for")
	}
	index := strings.LastIndex(str.Value(), substr.Value())
	return types.NewInt(int64(index)), nil
}

func stringRepeatMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.repeat() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.repeat() requires a string argument")
	}
	count, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.repeat() requires an int count")
	}
	return types.NewString(strings.Repeat(str.Value(), int(count.Value()))), nil
}

func stringReverseMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.reverse() takes exactly 1 argument")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.reverse() requires a string argument")
	}

	runes := []rune(str.Value())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return types.NewString(string(runes)), nil
}

func stringPadLeftMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("string.padLeft() takes 2 or 3 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.padLeft() requires a string argument")
	}
	length, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.padLeft() requires an int length")
	}

	padChar := " "
	if len(args) == 3 {
		padStr, ok := args[2].(*types.StringValue)
		if !ok {
			return nil, fmt.Errorf("string.padLeft() requires a string pad character")
		}
		if padStr.Value() != "" {
			padChar = string([]rune(padStr.Value())[0])
		}
	}

	s := str.Value()
	targetLen := int(length.Value())
	if len(s) >= targetLen {
		return str, nil
	}

	padding := strings.Repeat(padChar, targetLen-len(s))
	return types.NewString(padding + s), nil
}

func stringPadRightMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("string.padRight() takes 2 or 3 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.padRight() requires a string argument")
	}
	length, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("string.padRight() requires an int length")
	}

	padChar := " "
	if len(args) == 3 {
		padStr, ok := args[2].(*types.StringValue)
		if !ok {
			return nil, fmt.Errorf("string.padRight() requires a string pad character")
		}
		if padStr.Value() != "" {
			padChar = string([]rune(padStr.Value())[0])
		}
	}

	s := str.Value()
	targetLen := int(length.Value())
	if len(s) >= targetLen {
		return str, nil
	}

	padding := strings.Repeat(padChar, targetLen-len(s))
	return types.NewString(s + padding), nil
}

func stringMatchMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.match() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.match() requires a string argument")
	}
	pattern, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.match() requires a string pattern")
	}

	re, err := regexp.Compile(pattern.Value())
	if err != nil {
		return nil, fmt.Errorf("string.match() invalid regex pattern: %v", err)
	}

	matches := re.FindAllString(str.Value(), -1)
	result := make([]types.Value, len(matches))
	for i, match := range matches {
		result[i] = types.NewString(match)
	}

	elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewSlice(result, elemType), nil
}

func stringTestMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string.test() takes exactly 2 arguments")
	}
	str, ok := args[0].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.test() requires a string argument")
	}
	pattern, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("string.test() requires a string pattern")
	}

	matched, err := regexp.MatchString(pattern.Value(), str.Value())
	if err != nil {
		return nil, fmt.Errorf("string.test() invalid regex pattern: %v", err)
	}

	return types.NewBool(matched), nil
}

// Int type methods implementation

func intAbsMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.abs() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.abs() requires an int argument")
	}

	val := intVal.Value()
	if val < 0 {
		val = -val
	}
	return types.NewInt(val), nil
}

func intSignMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.sign() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.sign() requires an int argument")
	}

	val := intVal.Value()
	if val > 0 {
		return types.NewInt(1), nil
	} else if val < 0 {
		return types.NewInt(-1), nil
	}
	return types.NewInt(0), nil
}

func intToStringMethod(args []types.Value) (types.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("int.toString() takes 1 or 2 arguments")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.toString() requires an int argument")
	}

	base := 10
	if len(args) == 2 {
		baseVal, ok := args[1].(*types.IntValue)
		if !ok {
			return nil, fmt.Errorf("int.toString() requires an int base")
		}
		base = int(baseVal.Value())
		if base < 2 || base > 36 {
			return nil, fmt.Errorf("int.toString() base must be between 2 and 36")
		}
	}

	return types.NewString(strconv.FormatInt(intVal.Value(), base)), nil
}

func intToFloatMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.toFloat() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.toFloat() requires an int argument")
	}
	return types.NewFloat(float64(intVal.Value())), nil
}

func intToBoolMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.toBool() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.toBool() requires an int argument")
	}
	return types.NewBool(intVal.Value() != 0), nil
}

func intMinMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("int.min() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.min() requires int arguments")
	}
	b, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.min() requires int arguments")
	}

	if a.Value() < b.Value() {
		return a, nil
	}
	return b, nil
}

func intMaxMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("int.max() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.max() requires int arguments")
	}
	b, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.max() requires int arguments")
	}

	if a.Value() > b.Value() {
		return a, nil
	}
	return b, nil
}

func intClampMethod(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("int.clamp() takes exactly 3 arguments")
	}
	val, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.clamp() requires int arguments")
	}
	min, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.clamp() requires int arguments")
	}
	max, ok := args[2].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.clamp() requires int arguments")
	}

	value := val.Value()
	minVal := min.Value()
	maxVal := max.Value()

	if value < minVal {
		value = minVal
	} else if value > maxVal {
		value = maxVal
	}

	return types.NewInt(value), nil
}

func intIsEvenMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.isEven() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.isEven() requires an int argument")
	}
	return types.NewBool(intVal.Value()%2 == 0), nil
}

func intIsOddMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.isOdd() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.isOdd() requires an int argument")
	}
	return types.NewBool(intVal.Value()%2 != 0), nil
}

func intIsPrimeMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.isPrime() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.isPrime() requires an int argument")
	}

	n := intVal.Value()
	if n < 2 {
		return types.NewBool(false), nil
	}
	if n == 2 {
		return types.NewBool(true), nil
	}
	if n%2 == 0 {
		return types.NewBool(false), nil
	}

	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return types.NewBool(false), nil
		}
	}

	return types.NewBool(true), nil
}

func intFactorialMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int.factorial() takes exactly 1 argument")
	}
	intVal, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.factorial() requires an int argument")
	}

	n := intVal.Value()
	if n < 0 {
		return nil, fmt.Errorf("int.factorial() requires a non-negative integer")
	}
	if n == 0 || n == 1 {
		return types.NewInt(1), nil
	}

	result := int64(1)
	for i := int64(2); i <= n; i++ {
		result *= i
		// Check for overflow
		if result < 0 {
			return nil, fmt.Errorf("int.factorial() result overflow")
		}
	}

	return types.NewInt(result), nil
}

func intGcdMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("int.gcd() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.gcd() requires int arguments")
	}
	b, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.gcd() requires int arguments")
	}

	x, y := a.Value(), b.Value()
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	for y != 0 {
		x, y = y, x%y
	}

	return types.NewInt(x), nil
}

func intLcmMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("int.lcm() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.lcm() requires int arguments")
	}
	b, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("int.lcm() requires int arguments")
	}

	x, y := a.Value(), b.Value()
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	// Calculate GCD first
	gcd := x
	temp := y
	for temp != 0 {
		gcd, temp = temp, gcd%temp
	}

	if gcd == 0 {
		return types.NewInt(0), nil
	}

	return types.NewInt((x / gcd) * y), nil
}

// Float type methods implementation

func floatAbsMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.abs() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.abs() requires a float argument")
	}
	return types.NewFloat(math.Abs(floatVal.Value())), nil
}

func floatSignMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.sign() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.sign() requires a float argument")
	}

	val := floatVal.Value()
	if val > 0 {
		return types.NewFloat(1.0), nil
	} else if val < 0 {
		return types.NewFloat(-1.0), nil
	} else if math.IsNaN(val) {
		return types.NewFloat(math.NaN()), nil
	}
	return types.NewFloat(0.0), nil
}

func floatRoundMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.round() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.round() requires a float argument")
	}
	return types.NewFloat(math.Round(floatVal.Value())), nil
}

func floatFloorMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.floor() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.floor() requires a float argument")
	}
	return types.NewFloat(math.Floor(floatVal.Value())), nil
}

func floatCeilMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.ceil() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.ceil() requires a float argument")
	}
	return types.NewFloat(math.Ceil(floatVal.Value())), nil
}

func floatTruncMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.trunc() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.trunc() requires a float argument")
	}
	return types.NewFloat(math.Trunc(floatVal.Value())), nil
}

func floatToStringMethod(args []types.Value) (types.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("float.toString() takes 1 or 2 arguments")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.toString() requires a float argument")
	}

	precision := -1
	if len(args) == 2 {
		precVal, ok := args[1].(*types.IntValue)
		if !ok {
			return nil, fmt.Errorf("float.toString() requires an int precision")
		}
		precision = int(precVal.Value())
	}

	if precision >= 0 {
		format := fmt.Sprintf("%%.%df", precision)
		return types.NewString(fmt.Sprintf(format, floatVal.Value())), nil
	}
	return types.NewString(fmt.Sprintf("%g", floatVal.Value())), nil
}

func floatToIntMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.toInt() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.toInt() requires a float argument")
	}
	return types.NewInt(int64(floatVal.Value())), nil
}

func floatToBoolMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.toBool() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.toBool() requires a float argument")
	}
	val := floatVal.Value()
	return types.NewBool(val != 0.0 && !math.IsNaN(val)), nil
}

func floatMinMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("float.min() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.min() requires float arguments")
	}
	b, ok := args[1].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.min() requires float arguments")
	}
	return types.NewFloat(math.Min(a.Value(), b.Value())), nil
}

func floatMaxMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("float.max() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.max() requires float arguments")
	}
	b, ok := args[1].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.max() requires float arguments")
	}
	return types.NewFloat(math.Max(a.Value(), b.Value())), nil
}

func floatClampMethod(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("float.clamp() takes exactly 3 arguments")
	}
	val, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.clamp() requires float arguments")
	}
	min, ok := args[1].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.clamp() requires float arguments")
	}
	max, ok := args[2].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.clamp() requires float arguments")
	}

	value := val.Value()
	minVal := min.Value()
	maxVal := max.Value()

	if value < minVal {
		value = minVal
	} else if value > maxVal {
		value = maxVal
	}

	return types.NewFloat(value), nil
}

func floatIsNaNMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.isNaN() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.isNaN() requires a float argument")
	}
	return types.NewBool(math.IsNaN(floatVal.Value())), nil
}

func floatIsInfMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.isInf() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.isInf() requires a float argument")
	}
	return types.NewBool(math.IsInf(floatVal.Value(), 0)), nil
}

func floatIsFiniteMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float.isFinite() takes exactly 1 argument")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.isFinite() requires a float argument")
	}
	val := floatVal.Value()
	return types.NewBool(!math.IsNaN(val) && !math.IsInf(val, 0)), nil
}

func floatPrecisionMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("float.precision() takes exactly 2 arguments")
	}
	floatVal, ok := args[0].(*types.FloatValue)
	if !ok {
		return nil, fmt.Errorf("float.precision() requires a float argument")
	}
	precision, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("float.precision() requires an int precision")
	}

	prec := int(precision.Value())
	if prec < 0 {
		prec = 0
	}

	multiplier := math.Pow(10, float64(prec))
	rounded := math.Round(floatVal.Value()*multiplier) / multiplier
	return types.NewFloat(rounded), nil
}

// Bool type methods implementation

func boolToStringMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool.toString() takes exactly 1 argument")
	}
	boolVal, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.toString() requires a bool argument")
	}
	if boolVal.Value() {
		return types.NewString("true"), nil
	}
	return types.NewString("false"), nil
}

func boolToIntMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool.toInt() takes exactly 1 argument")
	}
	boolVal, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.toInt() requires a bool argument")
	}
	if boolVal.Value() {
		return types.NewInt(1), nil
	}
	return types.NewInt(0), nil
}

func boolToFloatMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool.toFloat() takes exactly 1 argument")
	}
	boolVal, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.toFloat() requires a bool argument")
	}
	if boolVal.Value() {
		return types.NewFloat(1.0), nil
	}
	return types.NewFloat(0.0), nil
}

func boolNotMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool.not() takes exactly 1 argument")
	}
	boolVal, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.not() requires a bool argument")
	}
	return types.NewBool(!boolVal.Value()), nil
}

func boolAndMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("bool.and() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.and() requires bool arguments")
	}
	b, ok := args[1].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.and() requires bool arguments")
	}
	return types.NewBool(a.Value() && b.Value()), nil
}

func boolOrMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("bool.or() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.or() requires bool arguments")
	}
	b, ok := args[1].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.or() requires bool arguments")
	}
	return types.NewBool(a.Value() || b.Value()), nil
}

func boolXorMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("bool.xor() takes exactly 2 arguments")
	}
	a, ok := args[0].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.xor() requires bool arguments")
	}
	b, ok := args[1].(*types.BoolValue)
	if !ok {
		return nil, fmt.Errorf("bool.xor() requires bool arguments")
	}
	return types.NewBool(a.Value() != b.Value()), nil
}

// Slice/Array type methods implementation

func sliceLengthMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.length() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.length() requires a slice argument")
	}
	return types.NewInt(int64(slice.Len())), nil
}

func sliceIsEmptyMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.isEmpty() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.isEmpty() requires a slice argument")
	}
	return types.NewBool(slice.Len() == 0), nil
}

func sliceFirstMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.first() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.first() requires a slice argument")
	}
	if slice.Len() == 0 {
		return types.NewNil(), nil
	}
	return slice.Get(0), nil
}

func sliceLastMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.last() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.last() requires a slice argument")
	}
	if slice.Len() == 0 {
		return types.NewNil(), nil
	}
	return slice.Get(slice.Len() - 1), nil
}

func sliceGetMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.get() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.get() requires a slice argument")
	}
	index, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("slice.get() requires an int index")
	}

	idx := int(index.Value())
	if idx < 0 || idx >= slice.Len() {
		return types.NewNil(), nil
	}
	return slice.Get(idx), nil
}

func sliceContainsMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.contains() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.contains() requires a slice argument")
	}
	target := args[1]

	for i := 0; i < slice.Len(); i++ {
		if compareValues(slice.Get(i), target) == 0 {
			return types.NewBool(true), nil
		}
	}
	return types.NewBool(false), nil
}

func sliceIndexOfMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.indexOf() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.indexOf() requires a slice argument")
	}
	target := args[1]

	for i := 0; i < slice.Len(); i++ {
		if compareValues(slice.Get(i), target) == 0 {
			return types.NewInt(int64(i)), nil
		}
	}
	return types.NewInt(-1), nil
}

func sliceSliceMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("slice.slice() takes 2 or 3 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.slice() requires a slice argument")
	}
	start, ok := args[1].(*types.IntValue)
	if !ok {
		return nil, fmt.Errorf("slice.slice() requires an int start index")
	}

	startIdx := int(start.Value())
	endIdx := slice.Len()

	if len(args) == 3 {
		end, ok := args[2].(*types.IntValue)
		if !ok {
			return nil, fmt.Errorf("slice.slice() requires an int end index")
		}
		endIdx = int(end.Value())
	}

	// Handle negative indices
	if startIdx < 0 {
		startIdx = slice.Len() + startIdx
	}
	if endIdx < 0 {
		endIdx = slice.Len() + endIdx
	}

	// Clamp indices
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > slice.Len() {
		endIdx = slice.Len()
	}
	if startIdx > endIdx {
		startIdx = endIdx
	}

	result := make([]types.Value, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		result[i-startIdx] = slice.Get(i)
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceConcatMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.concat() takes exactly 2 arguments")
	}
	slice1, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.concat() requires slice arguments")
	}
	slice2, ok := args[1].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.concat() requires slice arguments")
	}

	totalLen := slice1.Len() + slice2.Len()
	result := make([]types.Value, totalLen)

	for i := 0; i < slice1.Len(); i++ {
		result[i] = slice1.Get(i)
	}
	for i := 0; i < slice2.Len(); i++ {
		result[slice1.Len()+i] = slice2.Get(i)
	}

	return types.NewSlice(result, slice1.ElementType()), nil
}

func slicePushMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.push() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.push() requires a slice argument")
	}
	element := args[1]

	result := make([]types.Value, slice.Len()+1)
	for i := 0; i < slice.Len(); i++ {
		result[i] = slice.Get(i)
	}
	result[slice.Len()] = element

	return types.NewSlice(result, slice.ElementType()), nil
}

func slicePopMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.pop() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.pop() requires a slice argument")
	}

	if slice.Len() == 0 {
		return types.NewNil(), nil
	}

	return slice.Get(slice.Len() - 1), nil
}

func sliceShiftMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.shift() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.shift() requires a slice argument")
	}

	if slice.Len() == 0 {
		return types.NewNil(), nil
	}

	return slice.Get(0), nil
}

func sliceUnshiftMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.unshift() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.unshift() requires a slice argument")
	}
	element := args[1]

	result := make([]types.Value, slice.Len()+1)
	result[0] = element
	for i := 0; i < slice.Len(); i++ {
		result[i+1] = slice.Get(i)
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceReverseMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.reverse() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.reverse() requires a slice argument")
	}

	result := make([]types.Value, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result[slice.Len()-1-i] = slice.Get(i)
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceSortMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.sort() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.sort() requires a slice argument")
	}

	result := make([]types.Value, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result[i] = slice.Get(i)
	}

	sort.Slice(result, func(i, j int) bool {
		return compareValues(result[i], result[j]) < 0
	})

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceFilterMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.filter() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.filter() requires a slice argument")
	}
	predicate := args[1]

	var result []types.Value
	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if shouldInclude(item, predicate) {
			result = append(result, item)
		}
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceMapMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.map() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.map() requires a slice argument")
	}
	transformer := args[1]

	result := make([]types.Value, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result[i] = applyTransform(slice.Get(i), transformer)
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceReduceMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("slice.reduce() takes 2 or 3 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.reduce() requires a slice argument")
	}
	reducer := args[1]

	if slice.Len() == 0 {
		if len(args) == 3 {
			return args[2], nil
		}
		return types.NewNil(), nil
	}

	var accumulator types.Value
	startIdx := 0

	if len(args) == 3 {
		accumulator = args[2]
	} else {
		accumulator = slice.Get(0)
		startIdx = 1
	}

	for i := startIdx; i < slice.Len(); i++ {
		accumulator = applyReducer(accumulator, slice.Get(i), reducer)
	}

	return accumulator, nil
}

func sliceForEachMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.forEach() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.forEach() requires a slice argument")
	}
	action := args[1]

	for i := 0; i < slice.Len(); i++ {
		applyTransform(slice.Get(i), action)
	}

	return types.NewNil(), nil
}

func sliceFindMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.find() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.find() requires a slice argument")
	}
	predicate := args[1]

	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if shouldInclude(item, predicate) {
			return item, nil
		}
	}

	return types.NewNil(), nil
}

func sliceFindIndexMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.findIndex() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.findIndex() requires a slice argument")
	}
	predicate := args[1]

	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if shouldInclude(item, predicate) {
			return types.NewInt(int64(i)), nil
		}
	}

	return types.NewInt(-1), nil
}

func sliceSomeMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.some() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.some() requires a slice argument")
	}
	predicate := args[1]

	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if shouldInclude(item, predicate) {
			return types.NewBool(true), nil
		}
	}

	return types.NewBool(false), nil
}

func sliceEveryMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.every() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.every() requires a slice argument")
	}
	predicate := args[1]

	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if !shouldInclude(item, predicate) {
			return types.NewBool(false), nil
		}
	}

	return types.NewBool(true), nil
}

func sliceJoinMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice.join() takes exactly 2 arguments")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.join() requires a slice argument")
	}
	separator, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("slice.join() requires a string separator")
	}

	if slice.Len() == 0 {
		return types.NewString(""), nil
	}

	parts := make([]string, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		parts[i] = slice.Get(i).String()
	}

	return types.NewString(strings.Join(parts, separator.Value())), nil
}

func sliceUniqueMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.unique() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.unique() requires a slice argument")
	}

	seen := make(map[string]bool)
	var result []types.Value

	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		key := item.String()
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return types.NewSlice(result, slice.ElementType()), nil
}

func sliceFlattenMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("slice.flatten() takes exactly 1 argument")
	}
	slice, ok := args[0].(*types.SliceValue)
	if !ok {
		return nil, fmt.Errorf("slice.flatten() requires a slice argument")
	}

	var result []types.Value
	for i := 0; i < slice.Len(); i++ {
		item := slice.Get(i)
		if subSlice, ok := item.(*types.SliceValue); ok {
			for j := 0; j < subSlice.Len(); j++ {
				result = append(result, subSlice.Get(j))
			}
		} else {
			result = append(result, item)
		}
	}

	elemType := types.TypeInfo{Kind: types.KindInterface, Name: "any", Size: -1}
	return types.NewSlice(result, elemType), nil
}

// Map/Object type methods implementation

func mapSizeMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.size() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.size() requires a map argument")
	}
	return types.NewInt(int64(mapVal.Len())), nil
}

func mapIsEmptyMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.isEmpty() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.isEmpty() requires a map argument")
	}
	return types.NewBool(mapVal.Len() == 0), nil
}

func mapKeysMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.keys() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.keys() requires a map argument")
	}

	keys := mapVal.Keys()
	result := make([]types.Value, len(keys))
	for i, key := range keys {
		result[i] = types.NewString(key)
	}

	elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewSlice(result, elemType), nil
}

func mapValuesMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.values() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.values() requires a map argument")
	}

	values := make([]types.Value, 0, mapVal.Len())
	for _, key := range mapVal.Keys() {
		if val, exists := mapVal.Get(key); exists {
			values = append(values, val)
		}
	}

	elemType := types.TypeInfo{Kind: types.KindInterface, Name: "any", Size: -1}
	return types.NewSlice(values, elemType), nil
}

func mapEntriesMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.entries() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.entries() requires a map argument")
	}

	entries := make([]types.Value, 0, mapVal.Len())
	for _, key := range mapVal.Keys() {
		if val, exists := mapVal.Get(key); exists {
			entry := []types.Value{types.NewString(key), val}
			entryType := types.TypeInfo{Kind: types.KindInterface, Name: "any", Size: -1}
			entries = append(entries, types.NewSlice(entry, entryType))
		}
	}

	elemType := types.TypeInfo{Kind: types.KindSlice, Name: "slice", Size: -1}
	return types.NewSlice(entries, elemType), nil
}

func mapHasMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.has() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.has() requires a map argument")
	}
	key, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("map.has() requires a string key")
	}

	_, exists := mapVal.Get(key.Value())
	return types.NewBool(exists), nil
}

func mapGetMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.get() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.get() requires a map argument")
	}
	key, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("map.get() requires a string key")
	}

	if val, exists := mapVal.Get(key.Value()); exists {
		return val, nil
	}
	return types.NewNil(), nil
}

func mapSetMethod(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("map.set() takes exactly 3 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.set() requires a map argument")
	}
	key, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("map.set() requires a string key")
	}
	value := args[2]

	// Create a new map with the updated value
	newData := make(map[string]types.Value)
	for _, k := range mapVal.Keys() {
		if v, exists := mapVal.Get(k); exists {
			newData[k] = v
		}
	}
	newData[key.Value()] = value

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(newData, keyType, mapVal.ValueType()), nil
}

func mapDeleteMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.delete() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.delete() requires a map argument")
	}
	key, ok := args[1].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("map.delete() requires a string key")
	}

	// Create a new map without the key
	newData := make(map[string]types.Value)
	for _, k := range mapVal.Keys() {
		if k != key.Value() {
			if v, exists := mapVal.Get(k); exists {
				newData[k] = v
			}
		}
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(newData, keyType, mapVal.ValueType()), nil
}

func mapClearMethod(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map.clear() takes exactly 1 argument")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.clear() requires a map argument")
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(make(map[string]types.Value), keyType, mapVal.ValueType()), nil
}

func mapMergeMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.merge() takes exactly 2 arguments")
	}
	map1, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.merge() requires map arguments")
	}
	map2, ok := args[1].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.merge() requires map arguments")
	}

	newData := make(map[string]types.Value)

	// Copy first map
	for _, k := range map1.Keys() {
		if v, exists := map1.Get(k); exists {
			newData[k] = v
		}
	}

	// Merge second map (overwrites conflicts)
	for _, k := range map2.Keys() {
		if v, exists := map2.Get(k); exists {
			newData[k] = v
		}
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(newData, keyType, map1.ValueType()), nil
}

func mapForEachMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.forEach() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.forEach() requires a map argument")
	}
	action := args[1]

	for _, key := range mapVal.Keys() {
		if val, exists := mapVal.Get(key); exists {
			applyTransform(val, action)
		}
	}

	return types.NewNil(), nil
}

func mapFilterMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.filter() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.filter() requires a map argument")
	}
	predicate := args[1]

	newData := make(map[string]types.Value)
	for _, key := range mapVal.Keys() {
		if val, exists := mapVal.Get(key); exists {
			if shouldInclude(val, predicate) {
				newData[key] = val
			}
		}
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(newData, keyType, mapVal.ValueType()), nil
}

func mapMapMethod(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map.map() takes exactly 2 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.map() requires a map argument")
	}
	transformer := args[1]

	newData := make(map[string]types.Value)
	for _, key := range mapVal.Keys() {
		if val, exists := mapVal.Get(key); exists {
			newData[key] = applyTransform(val, transformer)
		}
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
	return types.NewMap(newData, keyType, mapVal.ValueType()), nil
}

func mapReduceMethod(args []types.Value) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("map.reduce() takes 2 or 3 arguments")
	}
	mapVal, ok := args[0].(*types.MapValue)
	if !ok {
		return nil, fmt.Errorf("map.reduce() requires a map argument")
	}
	reducer := args[1]

	keys := mapVal.Keys()
	if len(keys) == 0 {
		if len(args) == 3 {
			return args[2], nil
		}
		return types.NewNil(), nil
	}

	var accumulator types.Value
	startIdx := 0

	if len(args) == 3 {
		accumulator = args[2]
	} else {
		if val, exists := mapVal.Get(keys[0]); exists {
			accumulator = val
		} else {
			return types.NewNil(), nil
		}
		startIdx = 1
	}

	for i := startIdx; i < len(keys); i++ {
		if val, exists := mapVal.Get(keys[i]); exists {
			accumulator = applyReducer(accumulator, val, reducer)
		}
	}

	return accumulator, nil
}

// Helper functions (using existing ones from pipeline.go)

// TypeMethodBuiltins contains type-specific method implementations
var TypeMethodBuiltins = map[string]BuiltinFunction{
	// String type methods
	"string.length":      stringLengthMethod,
	"string.charAt":      stringCharAtMethod,
	"string.charCodeAt":  stringCharCodeAtMethod,
	"string.slice":       stringSliceMethod,
	"string.split":       stringSplitMethod,
	"string.join":        stringJoinMethod,
	"string.replace":     stringReplaceMethod,
	"string.trim":        stringTrimMethod,
	"string.trimLeft":    stringTrimLeftMethod,
	"string.trimRight":   stringTrimRightMethod,
	"string.upper":       stringUpperMethod,
	"string.lower":       stringLowerMethod,
	"string.startsWith":  stringStartsWithMethod,
	"string.endsWith":    stringEndsWithMethod,
	"string.contains":    stringContainsMethod,
	"string.indexOf":     stringIndexOfMethod,
	"string.lastIndexOf": stringLastIndexOfMethod,
	"string.repeat":      stringRepeatMethod,
	"string.reverse":     stringReverseMethod,
	"string.padLeft":     stringPadLeftMethod,
	"string.padRight":    stringPadRightMethod,
	"string.match":       stringMatchMethod,
	"string.test":        stringTestMethod,

	// Int type methods
	"int.abs":       intAbsMethod,
	"int.sign":      intSignMethod,
	"int.toString":  intToStringMethod,
	"int.toFloat":   intToFloatMethod,
	"int.toBool":    intToBoolMethod,
	"int.min":       intMinMethod,
	"int.max":       intMaxMethod,
	"int.clamp":     intClampMethod,
	"int.isEven":    intIsEvenMethod,
	"int.isOdd":     intIsOddMethod,
	"int.isPrime":   intIsPrimeMethod,
	"int.factorial": intFactorialMethod,
	"int.gcd":       intGcdMethod,
	"int.lcm":       intLcmMethod,

	// Float type methods
	"float.abs":       floatAbsMethod,
	"float.sign":      floatSignMethod,
	"float.round":     floatRoundMethod,
	"float.floor":     floatFloorMethod,
	"float.ceil":      floatCeilMethod,
	"float.trunc":     floatTruncMethod,
	"float.toString":  floatToStringMethod,
	"float.toInt":     floatToIntMethod,
	"float.toBool":    floatToBoolMethod,
	"float.min":       floatMinMethod,
	"float.max":       floatMaxMethod,
	"float.clamp":     floatClampMethod,
	"float.isNaN":     floatIsNaNMethod,
	"float.isInf":     floatIsInfMethod,
	"float.isFinite":  floatIsFiniteMethod,
	"float.precision": floatPrecisionMethod,

	// Bool type methods
	"bool.toString": boolToStringMethod,
	"bool.toInt":    boolToIntMethod,
	"bool.toFloat":  boolToFloatMethod,
	"bool.not":      boolNotMethod,
	"bool.and":      boolAndMethod,
	"bool.or":       boolOrMethod,
	"bool.xor":      boolXorMethod,

	// Slice/Array type methods
	"slice.length":    sliceLengthMethod,
	"slice.isEmpty":   sliceIsEmptyMethod,
	"slice.first":     sliceFirstMethod,
	"slice.last":      sliceLastMethod,
	"slice.get":       sliceGetMethod,
	"slice.contains":  sliceContainsMethod,
	"slice.indexOf":   sliceIndexOfMethod,
	"slice.slice":     sliceSliceMethod,
	"slice.concat":    sliceConcatMethod,
	"slice.push":      slicePushMethod,
	"slice.pop":       slicePopMethod,
	"slice.shift":     sliceShiftMethod,
	"slice.unshift":   sliceUnshiftMethod,
	"slice.reverse":   sliceReverseMethod,
	"slice.sort":      sliceSortMethod,
	"slice.filter":    sliceFilterMethod,
	"slice.map":       sliceMapMethod,
	"slice.reduce":    sliceReduceMethod,
	"slice.forEach":   sliceForEachMethod,
	"slice.find":      sliceFindMethod,
	"slice.findIndex": sliceFindIndexMethod,
	"slice.some":      sliceSomeMethod,
	"slice.every":     sliceEveryMethod,
	"slice.join":      sliceJoinMethod,
	"slice.unique":    sliceUniqueMethod,
	"slice.flatten":   sliceFlattenMethod,

	// Map/Object type methods
	"map.size":    mapSizeMethod,
	"map.isEmpty": mapIsEmptyMethod,
	"map.keys":    mapKeysMethod,
	"map.values":  mapValuesMethod,
	"map.entries": mapEntriesMethod,
	"map.has":     mapHasMethod,
	"map.get":     mapGetMethod,
	"map.set":     mapSetMethod,
	"map.delete":  mapDeleteMethod,
	"map.clear":   mapClearMethod,
	"map.merge":   mapMergeMethod,
	"map.forEach": mapForEachMethod,
	"map.filter":  mapFilterMethod,
	"map.map":     mapMapMethod,
	"map.reduce":  mapReduceMethod,
}
