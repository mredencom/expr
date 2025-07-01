package expr

import (
	"fmt"
	"time"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/env"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// Program represents a compiled expression program
type Program struct {
	bytecode      *vm.Bytecode
	envAdapter    *env.Adapter
	config        *Config
	variableOrder []string

	// Performance metrics
	compileTime time.Duration
	source      string
}

// Config holds configuration options for the expression engine
type Config struct {
	env                     interface{}
	allowUndefinedVariables bool
	disableAllBuiltins      bool
	builtins                map[string]interface{}
	operators               map[string]int

	// Type checking options
	expectedType       AsKind
	enableTypeChecking bool

	// Performance options
	enableCache        bool
	enableOptimization bool
	maxExecutionTime   time.Duration

	// Debug options
	enableDebug     bool
	enableProfiling bool
}

// Option represents a configuration option function
type Option func(*Config)

// Result represents the result of expression evaluation
type Result struct {
	Value interface{}
	Type  string

	// Performance metrics
	ExecutionTime time.Duration
	MemoryUsed    int64
}

// Statistics holds performance statistics
type Statistics struct {
	TotalCompilations  int64
	TotalExecutions    int64
	AverageCompileTime time.Duration
	AverageExecTime    time.Duration
	CacheHitRate       float64
	MemoryUsage        int64
}

// Global statistics
var globalStats = &Statistics{}

// Main API Functions

// Compile compiles an expression string into a Program
func Compile(expression string, options ...Option) (*Program, error) {
	start := time.Now()

	// Apply configuration options
	config := &Config{
		enableCache:        true,
		enableOptimization: true,
		maxExecutionTime:   time.Second * 30,
		builtins:           make(map[string]interface{}),
		operators:          make(map[string]int),
	}

	for _, option := range options {
		option(config)
	}

	// Parse the expression
	l := lexer.New(expression)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.Errors())
	}

	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("no statements found in expression")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		return nil, fmt.Errorf("expected expression statement")
	}

	// Compile to bytecode
	comp := compiler.New()

	// Add custom built-in functions
	for name := range config.builtins {
		comp.DefineBuiltin(name)
	}

	// Add environment if provided
	if config.env != nil {
		adapter := env.New()
		if envMap, ok := config.env.(map[string]interface{}); ok {
			err := comp.AddEnvironment(envMap, adapter)
			if err != nil {
				return nil, fmt.Errorf("environment error: %v", err)
			}
		}
	}

	err := comp.Compile(stmt.Expression)
	if err != nil {
		return nil, fmt.Errorf("compilation error: %v", err)
	}

	// Perform type checking if enabled
	if config.enableTypeChecking {
		err = validateExpectedType(stmt.Expression, config.expectedType)
		if err != nil {
			return nil, fmt.Errorf("type validation error: %v", err)
		}
	}

	bytecode := comp.Bytecode()
	variableOrder := comp.GetVariableOrder()
	compileTime := time.Since(start)

	// Update global statistics
	globalStats.TotalCompilations++
	globalStats.AverageCompileTime = updateAverage(
		globalStats.AverageCompileTime,
		compileTime,
		globalStats.TotalCompilations,
	)

	return &Program{
		bytecode:      bytecode,
		envAdapter:    env.New(),
		config:        config,
		variableOrder: variableOrder,
		compileTime:   compileTime,
		source:        expression,
	}, nil
}

// Run executes a compiled program with the given environment
func Run(program *Program, environment interface{}) (interface{}, error) {
	result, err := RunWithResult(program, environment)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

// RunWithResult executes a program and returns detailed result information
func RunWithResult(program *Program, environment interface{}) (*Result, error) {
	start := time.Now()

	// Get VM from pool instead of creating new one
	machine := vm.GlobalVMPool.Get()
	defer vm.GlobalVMPool.Put(machine) // Return to pool when done

	// Set up the VM with program data
	machine.SetConstants(program.bytecode.Constants)

	if environment != nil {
		if envMap, ok := environment.(map[string]interface{}); ok {
			err := machine.SetEnvironment(envMap, program.variableOrder)
			if err != nil {
				return nil, fmt.Errorf("environment setup error: %v", err)
			}
		}
	}

	// Set custom builtins if any
	if len(program.config.builtins) > 0 {
		for name, fn := range program.config.builtins {
			machine.SetCustomBuiltin(name, fn)
		}
	}

	// Execute with timeout if configured
	var result types.Value
	var execErr error

	if program.config.maxExecutionTime > 0 {
		done := make(chan struct{})
		go func() {
			defer close(done)
			result, execErr = machine.RunInstructionsWithResult(program.bytecode.Instructions)
		}()

		select {
		case <-done:
			// Execution completed
		case <-time.After(program.config.maxExecutionTime):
			return nil, fmt.Errorf("execution timeout after %v", program.config.maxExecutionTime)
		}
	} else {
		result, execErr = machine.RunInstructionsWithResult(program.bytecode.Instructions)
	}

	if execErr != nil {
		return nil, fmt.Errorf("execution error: %v", execErr)
	}

	execTime := time.Since(start)

	// Update global statistics
	globalStats.TotalExecutions++
	globalStats.AverageExecTime = updateAverage(
		globalStats.AverageExecTime,
		execTime,
		globalStats.TotalExecutions,
	)

	// Convert result to Go value
	var goValue interface{}
	if result != nil {
		goValue = convertTypesValueToGoValue(result)
	}

	return &Result{
		Value:         goValue,
		Type:          inferResultType(result),
		ExecutionTime: execTime,
		MemoryUsed:    0, // TODO: implement memory tracking
	}, nil
}

// Eval is a convenience function that compiles and runs an expression in one call
func Eval(expression string, environment interface{}) (interface{}, error) {
	program, err := Compile(expression, Env(environment))
	if err != nil {
		return nil, err
	}

	return Run(program, environment)
}

// EvalWithResult is like Eval but returns detailed result information
func EvalWithResult(expression string, environment interface{}) (*Result, error) {
	program, err := Compile(expression, Env(environment))
	if err != nil {
		return &Result{
			Value: nil,
			Type:  "error",
		}, err
	}

	return RunWithResult(program, environment)
}

// Configuration Options

// Env sets the environment for variable resolution
func Env(env interface{}) Option {
	return func(c *Config) {
		c.env = env
	}
}

// AllowUndefinedVariables allows undefined variables to be used
func AllowUndefinedVariables() Option {
	return func(c *Config) {
		c.allowUndefinedVariables = true
	}
}

// DisableAllBuiltins disables all built-in functions
func DisableAllBuiltins() Option {
	return func(c *Config) {
		c.disableAllBuiltins = true
	}
}

// WithBuiltin adds a custom built-in function
func WithBuiltin(name string, fn interface{}) Option {
	return func(c *Config) {
		c.builtins[name] = fn
	}
}

// WithOperator adds a custom operator with precedence
func WithOperator(op string, precedence int) Option {
	return func(c *Config) {
		c.operators[op] = precedence
	}
}

// EnableCache enables instruction caching
func EnableCache() Option {
	return func(c *Config) {
		c.enableCache = true
	}
}

// DisableCache disables instruction caching
func DisableCache() Option {
	return func(c *Config) {
		c.enableCache = false
	}
}

// EnableOptimization enables bytecode optimization
func EnableOptimization() Option {
	return func(c *Config) {
		c.enableOptimization = true
	}
}

// DisableOptimization disables bytecode optimization
func DisableOptimization() Option {
	return func(c *Config) {
		c.enableOptimization = false
	}
}

// WithTimeout sets the maximum execution time
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.maxExecutionTime = timeout
	}
}

// EnableDebug enables debug mode
func EnableDebug() Option {
	return func(c *Config) {
		c.enableDebug = true
	}
}

// EnableProfiling enables performance profiling
func EnableProfiling() Option {
	return func(c *Config) {
		c.enableProfiling = true
	}
}

// Utility Functions

// GetStatistics returns global performance statistics
func GetStatistics() *Statistics {
	return &Statistics{
		TotalCompilations:  globalStats.TotalCompilations,
		TotalExecutions:    globalStats.TotalExecutions,
		AverageCompileTime: globalStats.AverageCompileTime,
		AverageExecTime:    globalStats.AverageExecTime,
		CacheHitRate:       globalStats.CacheHitRate,
		MemoryUsage:        globalStats.MemoryUsage,
	}
}

// ResetStatistics resets global performance statistics
func ResetStatistics() {
	globalStats = &Statistics{}
}

// Program Methods

// Source returns the original expression source
func (p *Program) Source() string {
	return p.source
}

// CompileTime returns the compilation time
func (p *Program) CompileTime() time.Duration {
	return p.compileTime
}

// BytecodeSize returns the size of the compiled bytecode
func (p *Program) BytecodeSize() int {
	return len(p.bytecode.Instructions)
}

// ConstantsCount returns the number of constants in the program
func (p *Program) ConstantsCount() int {
	return len(p.bytecode.Constants)
}

// String returns a string representation of the program
func (p *Program) String() string {
	return fmt.Sprintf("Program{source: %q, bytecode: %d bytes, constants: %d}",
		p.source, p.BytecodeSize(), p.ConstantsCount())
}

// Helper functions

func updateAverage(currentAvg time.Duration, newValue time.Duration, count int64) time.Duration {
	if count == 1 {
		return newValue
	}
	return time.Duration((int64(currentAvg)*(count-1) + int64(newValue)) / count)
}

// validateExpectedType performs compile-time type validation
func validateExpectedType(expr ast.Expression, expectedType AsKind) error {
	if expectedType == AsAny {
		return nil // No validation needed for any type
	}

	// Infer type from expression
	inferredType := inferExpressionType(expr)

	// Validate based on expected type
	switch expectedType {
	case AsIntKind:
		if inferredType != "int" && inferredType != "int64" && inferredType != "numeric" {
			return fmt.Errorf("expected integer type, got %s", inferredType)
		}
	case AsInt64Kind:
		if inferredType != "int64" && inferredType != "int" && inferredType != "numeric" {
			return fmt.Errorf("expected int64 type, got %s", inferredType)
		}
	case AsFloat64Kind:
		if inferredType != "float64" && inferredType != "numeric" && inferredType != "int" && inferredType != "int64" {
			return fmt.Errorf("expected numeric type, got %s", inferredType)
		}
	case AsStringKind:
		// Any type can be converted to string, so we're permissive here
		if inferredType != "string" && inferredType != "unknown" && inferredType != "int" &&
			inferredType != "float64" && inferredType != "bool" && inferredType != "numeric" {
			return fmt.Errorf("expected string type, got %s", inferredType)
		}
	case AsBoolKind:
		if inferredType != "bool" {
			return fmt.Errorf("expected boolean type, got %s", inferredType)
		}
	default:
		return fmt.Errorf("unknown expected type: %v", expectedType)
	}

	return nil
}

// inferExpressionType infers the type of an expression
func inferExpressionType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Literal:
		if e.Value == nil {
			return "nil"
		}
		switch e.Value.(type) {
		case *types.IntValue:
			return "int"
		case *types.FloatValue:
			return "float64"
		case *types.StringValue:
			return "string"
		case *types.BoolValue:
			return "bool"
		default:
			return "unknown"
		}
	case *ast.InfixExpression:
		// Infer based on operator and operands
		switch e.Operator {
		case "+", "-", "*", "/", "%":
			leftType := inferExpressionType(e.Left)
			rightType := inferExpressionType(e.Right)
			if leftType == "float64" || rightType == "float64" {
				return "float64"
			}
			if leftType == "string" || rightType == "string" {
				return "string" // String concatenation
			}
			return "numeric" // Could be int or int64
		case "==", "!=", "<", "<=", ">", ">=", "&&", "||":
			return "bool"
		default:
			return "unknown"
		}
	case *ast.PrefixExpression:
		switch e.Operator {
		case "!":
			return "bool"
		case "-":
			rightType := inferExpressionType(e.Right)
			if rightType == "float64" {
				return "float64"
			}
			return "numeric"
		default:
			return "unknown"
		}
	case *ast.Identifier:
		// For identifiers, we can't easily infer the type without environment
		// Return a generic type that will pass most validations
		return "unknown"
	case *ast.CallExpression, *ast.BuiltinExpression:
		// Function calls can return any type, we'll be permissive
		return "unknown"
	default:
		return "unknown"
	}
}

// convertToExpectedType converts a runtime value to the expected type
func convertToExpectedType(value types.Value, expectedType AsKind) (interface{}, error) {
	if expectedType == AsAny {
		// Return the native Go value
		switch v := value.(type) {
		case *types.IntValue:
			return v.Value(), nil
		case *types.FloatValue:
			return v.Value(), nil
		case *types.StringValue:
			return v.Value(), nil
		case *types.BoolValue:
			return v.Value(), nil
		default:
			return v.String(), nil
		}
	}

	switch expectedType {
	case AsIntKind:
		switch v := value.(type) {
		case *types.IntValue:
			return int(v.Value()), nil
		case *types.FloatValue:
			return int(v.Value()), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int", value)
		}
	case AsInt64Kind:
		switch v := value.(type) {
		case *types.IntValue:
			return v.Value(), nil
		case *types.FloatValue:
			return int64(v.Value()), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int64", value)
		}
	case AsFloat64Kind:
		switch v := value.(type) {
		case *types.IntValue:
			return float64(v.Value()), nil
		case *types.FloatValue:
			return v.Value(), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float64", value)
		}
	case AsStringKind:
		switch v := value.(type) {
		case *types.StringValue:
			return v.Value(), nil
		default:
			return v.String(), nil // All types can be converted to string
		}
	case AsBoolKind:
		switch v := value.(type) {
		case *types.BoolValue:
			return v.Value(), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", value)
		}
	default:
		return nil, fmt.Errorf("unknown expected type: %v", expectedType)
	}
}

// convertTypesValueToGoValue converts a types.Value to a Go value
func convertTypesValueToGoValue(value types.Value) interface{} {
	switch v := value.(type) {
	case *types.IntValue:
		return v.Value()
	case *types.FloatValue:
		return v.Value()
	case *types.StringValue:
		return v.Value()
	case *types.BoolValue:
		return v.Value()
	default:
		return v.String()
	}
}

// inferResultType infers the type of a result value
func inferResultType(value types.Value) string {
	switch value.(type) {
	case *types.IntValue:
		return "int"
	case *types.FloatValue:
		return "float64"
	case *types.StringValue:
		return "string"
	case *types.BoolValue:
		return "bool"
	default:
		return "unknown"
	}
}
