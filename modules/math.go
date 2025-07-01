package modules

import (
	"fmt"
	"math"

	"github.com/mredencom/expr/types"
)

// registerMathModule registers the math module with mathematical functions
func (r *Registry) registerMathModule() {
	functions := map[string]*ModuleFunction{
		"sqrt": {
			Name:        "sqrt",
			Description: "Returns the square root of x",
			Handler:     mathSqrt,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"pow": {
			Name:        "pow",
			Description: "Returns x raised to the power of y",
			Handler:     mathPow,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"abs": {
			Name:        "abs",
			Description: "Returns the absolute value of x",
			Handler:     mathAbs,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"floor": {
			Name:        "floor",
			Description: "Returns the largest integer less than or equal to x",
			Handler:     mathFloor,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"ceil": {
			Name:        "ceil",
			Description: "Returns the smallest integer greater than or equal to x",
			Handler:     mathCeil,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"round": {
			Name:        "round",
			Description: "Returns the nearest integer to x",
			Handler:     mathRound,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"max": {
			Name:        "max",
			Description: "Returns the maximum of two numbers",
			Handler:     mathMax,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"min": {
			Name:        "min",
			Description: "Returns the minimum of two numbers",
			Handler:     mathMin,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"sin": {
			Name:        "sin",
			Description: "Returns the sine of x (in radians)",
			Handler:     mathSin,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"cos": {
			Name:        "cos",
			Description: "Returns the cosine of x (in radians)",
			Handler:     mathCos,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"tan": {
			Name:        "tan",
			Description: "Returns the tangent of x (in radians)",
			Handler:     mathTan,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"log": {
			Name:        "log",
			Description: "Returns the natural logarithm of x",
			Handler:     mathLog,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"log10": {
			Name:        "log10",
			Description: "Returns the base-10 logarithm of x",
			Handler:     mathLog10,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"exp": {
			Name:        "exp",
			Description: "Returns e raised to the power of x",
			Handler:     mathExp,
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
	}

	r.RegisterModule("math", "Mathematical functions", functions)
}

// Math function implementations

func mathSqrt(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Sqrt(x), nil
}

func mathPow(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow expects 2 arguments, got %d", len(args))
	}
	x := toFloat64(args[0])
	y := toFloat64(args[1])
	return math.Pow(x, y), nil
}

func mathAbs(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Abs(x), nil
}

func mathFloor(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Floor(x), nil
}

func mathCeil(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Ceil(x), nil
}

func mathRound(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("round expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Round(x), nil
}

func mathMax(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("max expects 2 arguments, got %d", len(args))
	}
	x := toFloat64(args[0])
	y := toFloat64(args[1])
	return math.Max(x, y), nil
}

func mathMin(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("min expects 2 arguments, got %d", len(args))
	}
	x := toFloat64(args[0])
	y := toFloat64(args[1])
	return math.Min(x, y), nil
}

func mathSin(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Sin(x), nil
}

func mathCos(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Cos(x), nil
}

func mathTan(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Tan(x), nil
}

func mathLog(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Log(x), nil
}

func mathLog10(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log10 expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Log10(x), nil
}

func mathExp(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("exp expects 1 argument, got %d", len(args))
	}
	x := toFloat64(args[0])
	return math.Exp(x), nil
}

// toFloat64 converts an interface{} to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	default:
		return 0.0
	}
}
