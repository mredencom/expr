package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mredencom/expr/types"
)

// registerStringsModule registers the strings module with string manipulation functions
func (r *Registry) registerStringsModule() {
	functions := map[string]*ModuleFunction{
		"upper": {
			Name:        "upper",
			Description: "Converts string to uppercase",
			Handler:     stringsUpper,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"lower": {
			Name:        "lower",
			Description: "Converts string to lowercase",
			Handler:     stringsLower,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"trim": {
			Name:        "trim",
			Description: "Removes leading and trailing whitespace",
			Handler:     stringsTrim,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"length": {
			Name:        "length",
			Description: "Returns the length of the string",
			Handler:     stringsLength,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindInt, Name: "int"},
			Variadic:   false,
		},
		"contains": {
			Name:        "contains",
			Description: "Checks if string contains substring",
			Handler:     stringsContains,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			Variadic:   false,
		},
		"startsWith": {
			Name:        "startsWith",
			Description: "Checks if string starts with prefix",
			Handler:     stringsStartsWith,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			Variadic:   false,
		},
		"endsWith": {
			Name:        "endsWith",
			Description: "Checks if string ends with suffix",
			Handler:     stringsEndsWith,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindBool, Name: "bool"},
			Variadic:   false,
		},
		"replace": {
			Name:        "replace",
			Description: "Replaces all occurrences of old with new",
			Handler:     stringsReplace,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"split": {
			Name:        "split",
			Description: "Splits string by separator",
			Handler:     stringsSplit,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindSlice, Name: "[]string"},
			Variadic:   false,
		},
		"join": {
			Name:        "join",
			Description: "Joins string slice with separator",
			Handler:     stringsJoin,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindSlice, Name: "[]string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"substring": {
			Name:        "substring",
			Description: "Extracts substring from start to end",
			Handler:     stringsSubstring,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindInt, Name: "int"},
				{Kind: types.KindInt, Name: "int"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
		"indexOf": {
			Name:        "indexOf",
			Description: "Returns the index of the first occurrence of substring",
			Handler:     stringsIndexOf,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindInt, Name: "int"},
			Variadic:   false,
		},
		"repeat": {
			Name:        "repeat",
			Description: "Repeats string n times",
			Handler:     stringsRepeat,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
				{Kind: types.KindInt, Name: "int"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
	}

	r.RegisterModule("strings", "String manipulation functions", functions)
}

// String function implementations

func stringsUpper(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper expects 1 argument, got %d", len(args))
	}
	s := toString(args[0])
	return strings.ToUpper(s), nil
}

func stringsLower(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower expects 1 argument, got %d", len(args))
	}
	s := toString(args[0])
	return strings.ToLower(s), nil
}

func stringsTrim(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim expects 1 argument, got %d", len(args))
	}
	s := toString(args[0])
	return strings.TrimSpace(s), nil
}

func stringsLength(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("length expects 1 argument, got %d", len(args))
	}
	s := toString(args[0])
	return len(s), nil
}

func stringsContains(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	substr := toString(args[1])
	return strings.Contains(s, substr), nil
}

func stringsStartsWith(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("startsWith expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	prefix := toString(args[1])
	return strings.HasPrefix(s, prefix), nil
}

func stringsEndsWith(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("endsWith expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	suffix := toString(args[1])
	return strings.HasSuffix(s, suffix), nil
}

func stringsReplace(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace expects 3 arguments, got %d", len(args))
	}
	s := toString(args[0])
	old := toString(args[1])
	new := toString(args[2])
	return strings.ReplaceAll(s, old, new), nil
}

func stringsSplit(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	sep := toString(args[1])
	return strings.Split(s, sep), nil
}

func stringsJoin(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join expects 2 arguments, got %d", len(args))
	}

	// Convert first argument to string slice
	arr := args[0]
	var strs []string
	switch v := arr.(type) {
	case []string:
		strs = v
	case []interface{}:
		strs = make([]string, len(v))
		for i, item := range v {
			strs[i] = toString(item)
		}
	default:
		return nil, fmt.Errorf("join expects first argument to be string array, got %T", arr)
	}

	sep := toString(args[1])
	return strings.Join(strs, sep), nil
}

func stringsSubstring(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("substring expects 3 arguments, got %d", len(args))
	}
	s := toString(args[0])
	start := toInt(args[1])
	end := toInt(args[2])

	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		start = end
	}

	return s[start:end], nil
}

func stringsIndexOf(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("indexOf expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	substr := toString(args[1])
	return strings.Index(s, substr), nil
}

func stringsRepeat(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("repeat expects 2 arguments, got %d", len(args))
	}
	s := toString(args[0])
	count := toInt(args[1])
	if count < 0 {
		count = 0
	}
	return strings.Repeat(s, count), nil
}

// toString converts an interface{} to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// toInt converts an interface{} to int
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int8:
		return int(val)
	case int16:
		return int(val)
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint8:
		return int(val)
	case uint16:
		return int(val)
	case uint32:
		return int(val)
	case uint64:
		return int(val)
	case float32:
		return int(val)
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}
