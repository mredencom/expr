package expr

import (
	"fmt"
)

// Compatibility API for expr-lang/expr library

// AsKind represents the different types that can be used in As() function
type AsKind int

const (
	AsAny AsKind = iota
	AsIntKind
	AsInt64Kind
	AsFloat64Kind
	AsStringKind
	AsBoolKind
)

// As configures the expected return type for the expression
func As(kind AsKind) Option {
	return func(c *Config) {
		c.expectedType = kind
		c.enableTypeChecking = true
	}
}

// AsInt configures the expression to return an integer
func AsInt() Option {
	return As(AsIntKind)
}

// AsInt64 configures the expression to return an int64
func AsInt64() Option {
	return As(AsInt64Kind)
}

// AsFloat64 configures the expression to return a float64
func AsFloat64() Option {
	return As(AsFloat64Kind)
}

// AsString configures the expression to return a string
func AsString() Option {
	return As(AsStringKind)
}

// AsBool configures the expression to return a boolean
func AsBool() Option {
	return As(AsBoolKind)
}

// Function represents a custom function that can be used in expressions
type Function struct {
	Name     string
	Func     interface{}
	Variadic bool
}

// Functions option for adding multiple custom functions
func Functions(funcs map[string]interface{}) Option {
	return func(c *Config) {
		for name, fn := range funcs {
			c.builtins[name] = fn
		}
	}
}

// Operator represents a custom operator
type Operator struct {
	Symbol     string
	Precedence int
	Func       interface{}
}

// Operators option for adding multiple custom operators
func Operators(ops map[string]Operator) Option {
	return func(c *Config) {
		for symbol, op := range ops {
			c.operators[symbol] = op.Precedence
		}
	}
}

// Patch represents a patch to apply to the AST
type Patch struct {
	// Implementation would go here for AST patching
}

// Patch option for applying AST patches
func Patches(patches ...Patch) Option {
	return func(c *Config) {
		// AST patching would be implemented here
	}
}

// Tag represents a struct tag configuration
type Tag struct {
	Name string
}

// Tags option for configuring struct tag handling
func Tags(tag Tag) Option {
	return func(c *Config) {
		// Struct tag handling would be implemented here
	}
}

// Optimize enables various optimizations
func Optimize(enabled bool) Option {
	return func(c *Config) {
		c.enableOptimization = enabled
	}
}

// ConstExpr marks an expression as constant for optimization
func ConstExpr(name string) Option {
	return func(c *Config) {
		// Constant expression optimization would be implemented here
	}
}

// Deprecated compatibility functions for older versions

// NewEnv creates a new environment (deprecated, use Env() instead)
func NewEnv() map[string]interface{} {
	return make(map[string]interface{})
}

// CompileWithEnv compiles an expression with an environment (deprecated)
func CompileWithEnv(expression string, env interface{}) (*Program, error) {
	return Compile(expression, Env(env))
}

// RunWithEnv runs a program with an environment (deprecated)
func RunWithEnv(program *Program, env interface{}) (interface{}, error) {
	return Run(program, env)
}

// EvalWithEnv evaluates an expression with an environment (deprecated)
func EvalWithEnv(expression string, env interface{}) (interface{}, error) {
	return Eval(expression, env)
}

// Type checking helpers for compatibility (using generics instead of reflection)

// CheckType validates that a value matches the expected type using generics
func CheckType[T any](value interface{}) error {
	var zero T

	// Direct type match
	if _, ok := value.(T); ok {
		return nil
	}

	// Handle numeric type compatibility
	switch any(zero).(type) {
	case int:
		// Accept int64 as int for compatibility
		if _, ok := value.(int64); ok {
			return nil
		}
	case int64:
		// Accept int as int64 for compatibility
		if _, ok := value.(int); ok {
			return nil
		}
	case float64:
		// Accept numeric types as float64
		switch value.(type) {
		case int, int64, float32:
			return nil
		}
	}

	return fmt.Errorf("type mismatch: expected %T, got %T", zero, value)
}

// ConvertType attempts to convert a value to the expected type using generics
func ConvertType[T any](value interface{}) (T, error) {
	var zero T

	if value == nil {
		return zero, nil
	}

	// Direct type match
	if converted, ok := value.(T); ok {
		return converted, nil
	}

	// Type-specific conversions using type switches
	switch any(zero).(type) {
	case string:
		result := fmt.Sprintf("%v", value)
		if converted, ok := any(result).(T); ok {
			return converted, nil
		}
	case int:
		switch v := value.(type) {
		case int64:
			if converted, ok := any(int(v)).(T); ok {
				return converted, nil
			}
		case float64:
			if converted, ok := any(int(v)).(T); ok {
				return converted, nil
			}
		}
	case int64:
		switch v := value.(type) {
		case int:
			if converted, ok := any(int64(v)).(T); ok {
				return converted, nil
			}
		case float64:
			if converted, ok := any(int64(v)).(T); ok {
				return converted, nil
			}
		}
	case float64:
		switch v := value.(type) {
		case int:
			if converted, ok := any(float64(v)).(T); ok {
				return converted, nil
			}
		case int64:
			if converted, ok := any(float64(v)).(T); ok {
				return converted, nil
			}
		}
	case bool:
		switch v := value.(type) {
		case string:
			result := v == "true"
			if converted, ok := any(result).(T); ok {
				return converted, nil
			}
		case int:
			result := v != 0
			if converted, ok := any(result).(T); ok {
				return converted, nil
			}
		case int64:
			result := v != 0
			if converted, ok := any(result).(T); ok {
				return converted, nil
			}
		case float64:
			result := v != 0.0
			if converted, ok := any(result).(T); ok {
				return converted, nil
			}
		}
	}

	return zero, fmt.Errorf("cannot convert %T to %T", value, zero)
}

// Error types for compatibility

// CompileError represents a compilation error
type CompileError struct {
	Message  string
	Position int
	Line     int
	Column   int
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("compile error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// RuntimeError represents a runtime error
type RuntimeError struct {
	Message string
	Cause   error
}

func (e *RuntimeError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("runtime error: %s (caused by: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("runtime error: %s", e.Message)
}

// Unwrap returns the underlying cause
func (e *RuntimeError) Unwrap() error {
	return e.Cause
}

// Helper functions for creating errors

// NewCompileError creates a new compile error
func NewCompileError(message string, line, column int) *CompileError {
	return &CompileError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}

// NewRuntimeError creates a new runtime error
func NewRuntimeError(message string, cause error) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Cause:   cause,
	}
}

// Utility functions for compatibility (using type assertions instead of reflection)

// GetType returns the type name of a value using type assertions
func GetType(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch value.(type) {
	case bool:
		return "bool"
	case int:
		return "int"
	case int8:
		return "int8"
	case int16:
		return "int16"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case uint:
		return "uint"
	case uint8:
		return "uint8"
	case uint16:
		return "uint16"
	case uint32:
		return "uint32"
	case uint64:
		return "uint64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case string:
		return "string"
	case []interface{}:
		return "[]interface{}"
	case map[string]interface{}:
		return "map[string]interface{}"
	default:
		return fmt.Sprintf("%T", value)
	}
}

// IsNil checks if a value is nil using type assertions
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}

	// Check for typed nil values
	switch v := value.(type) {
	case *int:
		return v == nil
	case *string:
		return v == nil
	case *bool:
		return v == nil
	case *float64:
		return v == nil
	case []interface{}:
		return v == nil
	case map[string]interface{}:
		return v == nil
	case func():
		return v == nil
	}

	return false
}

// ToMap converts a struct to a map for environment usage using type assertions
func ToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if v == nil {
		return result
	}

	if m, ok := v.(map[string]interface{}); ok {
		return m
	}

	// For struct conversion, we need to use a different approach without reflection
	// We can provide a simple interface-based solution
	if mapper, ok := v.(interface{ ToMap() map[string]interface{} }); ok {
		return mapper.ToMap()
	}

	// For basic types, create a simple map
	switch val := v.(type) {
	case string:
		result["value"] = val
	case int, int64, float64, bool:
		result["value"] = val
	default:
		result["value"] = fmt.Sprintf("%v", val)
	}

	return result
}

// Mappable interface for types that can convert themselves to maps
type Mappable interface {
	ToMap() map[string]interface{}
}

// Example implementation helper for structs
func StructToMap(fields map[string]interface{}) map[string]interface{} {
	return fields
}
