package builtins

import (
	"fmt"

	"github.com/mredencom/expr/types"
)

// Collection operation built-in functions

// All returns true if all elements in the collection satisfy the predicate
func All(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("all() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		for _, item := range items {
			if boolVal, ok := item.(*types.BoolValue); ok {
				if !boolVal.Value() {
					return types.NewBool(false), nil
				}
			} else {
				// Non-boolean values are considered truthy if not zero/empty
				if isZeroValue(item) {
					return types.NewBool(false), nil
				}
			}
		}
		return types.NewBool(true), nil
	default:
		return nil, fmt.Errorf("all() requires a collection, got %s", collection.Type().Name)
	}
}

// Any returns true if any element in the collection satisfies the predicate
func Any(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("any() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		for _, item := range items {
			if boolVal, ok := item.(*types.BoolValue); ok {
				if boolVal.Value() {
					return types.NewBool(true), nil
				}
			} else {
				// Non-boolean values are considered truthy if not zero/empty
				if !isZeroValue(item) {
					return types.NewBool(true), nil
				}
			}
		}
		return types.NewBool(false), nil
	default:
		return nil, fmt.Errorf("any() requires a collection, got %s", collection.Type().Name)
	}
}

// Filter returns a new collection with elements that satisfy the predicate
func Filter(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("filter() takes exactly 2 arguments (%d given)", len(args))
	}

	collection := args[0]
	predicate := args[1]

	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		var filtered []types.Value

		for _, item := range items {
			// Apply predicate to each item
			result, err := applyPredicate(predicate, item)
			if err != nil {
				return nil, fmt.Errorf("filter predicate error: %v", err)
			}

			if isTruthy(result) {
				filtered = append(filtered, item)
			}
		}

		// Create a new slice with interface{} element type
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewSlice(filtered, elemType), nil
	default:
		return nil, fmt.Errorf("filter() requires a collection, got %s", collection.Type().Name)
	}
}

// Map applies a function to each element and returns a new collection
func Map(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map() takes exactly 2 arguments (%d given)", len(args))
	}

	collection := args[0]
	mapper := args[1]

	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		var mapped []types.Value

		for _, item := range items {
			// Apply mapper to each item
			result, err := applyPredicate(mapper, item)
			if err != nil {
				return nil, fmt.Errorf("map function error: %v", err)
			}

			mapped = append(mapped, result)
		}

		// Create a new slice with interface{} element type
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewSlice(mapped, elemType), nil
	default:
		return nil, fmt.Errorf("map() requires a collection, got %s", collection.Type().Name)
	}
}

// Sum calculates the sum of numeric values in a collection
func Sum(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sum() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		var intSum int64 = 0
		var floatSum float64 = 0
		hasFloat := false

		for _, item := range items {
			switch val := item.(type) {
			case *types.IntValue:
				if hasFloat {
					floatSum += float64(val.Value())
				} else {
					intSum += val.Value()
				}
			case *types.FloatValue:
				if !hasFloat {
					floatSum = float64(intSum) + val.Value()
					hasFloat = true
				} else {
					floatSum += val.Value()
				}
			default:
				return nil, fmt.Errorf("sum() requires numeric values, got %s", item.Type().Name)
			}
		}

		if hasFloat {
			return types.NewFloat(floatSum), nil
		}
		return types.NewInt(intSum), nil
	default:
		return nil, fmt.Errorf("sum() requires a collection, got %s", collection.Type().Name)
	}
}

// Count returns the number of elements in a collection
func Count(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("count() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		return types.NewInt(int64(len(coll.Values()))), nil
	case *types.MapValue:
		return types.NewInt(int64(len(coll.Values()))), nil
	case *types.StringValue:
		return types.NewInt(int64(len(coll.Value()))), nil
	default:
		return nil, fmt.Errorf("count() requires a collection, got %s", collection.Type().Name)
	}
}

// First returns the first element of a collection
func First(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		if len(items) == 0 {
			return types.NewNil(), nil
		}
		return items[0], nil
	case *types.StringValue:
		str := coll.Value()
		if len(str) == 0 {
			return types.NewString(""), nil
		}
		return types.NewString(string(str[0])), nil
	default:
		return nil, fmt.Errorf("first() requires a collection, got %s", collection.Type().Name)
	}
}

// Last returns the last element of a collection
func Last(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("last() takes exactly 1 argument (%d given)", len(args))
	}

	collection := args[0]
	switch coll := collection.(type) {
	case *types.SliceValue:
		items := coll.Values()
		if len(items) == 0 {
			return types.NewNil(), nil
		}
		return items[len(items)-1], nil
	case *types.StringValue:
		str := coll.Value()
		if len(str) == 0 {
			return types.NewString(""), nil
		}
		return types.NewString(string(str[len(str)-1])), nil
	default:
		return nil, fmt.Errorf("last() requires a collection, got %s", collection.Type().Name)
	}
}

// Helper functions

// isZeroValue checks if a value is considered "zero" or "empty"
func isZeroValue(value types.Value) bool {
	switch val := value.(type) {
	case *types.BoolValue:
		return !val.Value()
	case *types.IntValue:
		return val.Value() == 0
	case *types.FloatValue:
		return val.Value() == 0.0
	case *types.StringValue:
		return val.Value() == ""
	case *types.SliceValue:
		return len(val.Values()) == 0
	case *types.MapValue:
		return len(val.Values()) == 0
	case *types.NilValue:
		return true
	default:
		return false
	}
}

// isTruthy checks if a value is considered "truthy"
func isTruthy(value types.Value) bool {
	return !isZeroValue(value)
}

// applyPredicate applies a predicate function to a value
func applyPredicate(predicate types.Value, arg types.Value) (types.Value, error) {
	// For now, we'll treat the predicate as a simple expression
	// In a full implementation, this would involve calling a compiled function

	// Simple implementation: if predicate is a function name, we could look it up
	// For now, we'll just return the argument as-is for basic filtering
	if stringPred, ok := predicate.(*types.StringValue); ok {
		switch stringPred.Value() {
		case "not_empty":
			return types.NewBool(!isZeroValue(arg)), nil
		case "positive":
			if intVal, ok := arg.(*types.IntValue); ok {
				return types.NewBool(intVal.Value() > 0), nil
			}
			if floatVal, ok := arg.(*types.FloatValue); ok {
				return types.NewBool(floatVal.Value() > 0), nil
			}
			return types.NewBool(false), nil
		case "negative":
			if intVal, ok := arg.(*types.IntValue); ok {
				return types.NewBool(intVal.Value() < 0), nil
			}
			if floatVal, ok := arg.(*types.FloatValue); ok {
				return types.NewBool(floatVal.Value() < 0), nil
			}
			return types.NewBool(false), nil
		default:
			return types.NewBool(true), nil // Default to true for unknown predicates
		}
	}

	// For other predicate types, return the argument
	return arg, nil
}

// Collection operation registry
var CollectionBuiltins = map[string]func([]types.Value) (types.Value, error){
	"all":    All,
	"any":    Any,
	"filter": Filter,
	"map":    Map,
	"sum":    Sum,
	"count":  Count,
	"first":  First,
	"last":   Last,
}
