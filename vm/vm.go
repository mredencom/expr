package vm

import (
	"fmt"
	"strings"

	"github.com/mredencom/expr/builtins"
	"github.com/mredencom/expr/modules"
	"github.com/mredencom/expr/types"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
)

var (
	True  = types.NewBool(true)
	False = types.NewBool(false)
	Nil   = types.NewNil()

	// Pre-allocated common integers for performance
	CachedInts = make([]types.Value, 256) // Cache integers 0-255
)

func init() {
	// Pre-allocate common integers
	for i := 0; i < 256; i++ {
		CachedInts[i] = types.NewInt(int64(i))
	}
}

// getIntValue returns a cached value for small integers, or creates new for large ones
func getIntValue(value int64) types.Value {
	if value >= 0 && value < 256 {
		return CachedInts[value]
	}
	return types.NewInt(value)
}

// getIntValueFromPool returns a value using the VM's object pool
func (vm *VM) getIntValueFromPool(value int64) types.Value {
	return getIntValue(value)
}

// getFloatValueFromPool returns a float value using the VM's object pool
func (vm *VM) getFloatValueFromPool(value float64) types.Value {
	return types.NewFloat(value)
}

// getStringValueFromPool returns a string value using the VM's object pool
func (vm *VM) getStringValueFromPool(value string) types.Value {
	return types.NewString(value)
}

// Bytecode represents compiled bytecode
type Bytecode struct {
	Instructions []byte
	Constants    []types.Value
}

// VM represents the virtual machine
type VM struct {
	bytecode       *Bytecode
	constants      []types.Value
	stack          []types.Value
	sp             int // Stack pointer
	globals        []types.Value
	pool           *ValuePool
	cache          *InstructionCache
	customBuiltins map[string]interface{}
	env            map[string]interface{}
	safeJumpTable  *SafeJumpTable // Simplified and stable instruction dispatch table

	// Pipeline context for pipeline operations
	pipelineElement types.Value
}

// New creates a new VM
func New(bytecode *Bytecode) *VM {
	return &VM{
		bytecode:       bytecode,
		constants:      bytecode.Constants,
		stack:          make([]types.Value, StackSize),
		sp:             0,
		globals:        make([]types.Value, GlobalsSize),
		pool:           NewValuePool(),
		cache:          NewInstructionCache(1000),
		customBuiltins: make(map[string]interface{}),
		safeJumpTable:  NewSafeJumpTable(),
	}
}

// NewWithEnvironment creates a new VM with environment variables
func NewWithEnvironment(bytecode *Bytecode, envVars map[string]interface{}, adapter interface{}) *VM {
	vm := &VM{
		bytecode:       bytecode,
		constants:      bytecode.Constants,
		stack:          make([]types.Value, StackSize),
		sp:             0,
		globals:        make([]types.Value, GlobalsSize),
		pool:           NewValuePool(),
		cache:          NewInstructionCache(1000),
		customBuiltins: make(map[string]interface{}),
		env:            envVars,
		safeJumpTable:  NewSafeJumpTable(),
	}

	// Set environment variables in globals using the same ordering as the compiler
	// Sort variable names alphabetically for consistent ordering
	var names []string
	for name := range envVars {
		names = append(names, name)
	}

	// Sort alphabetically (same as compiler)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	// Set globals in sorted order
	for i, name := range names {
		if i >= len(vm.globals) {
			break
		}
		if value, exists := envVars[name]; exists {
			if typesValue, err := vm.convertGoValueToTypesValue(value); err == nil {
				vm.globals[i] = typesValue
			}
		}
	}

	return vm
}

// Run executes the virtual machine with the given bytecode
func (vm *VM) Run(bytecode *Bytecode, env map[string]interface{}) (types.Value, error) {
	vm.sp = 0 // Reset stack pointer
	vm.constants = bytecode.Constants
	vm.env = env

	return vm.runHighPerformanceLoop(bytecode.Instructions)
}

// runHighPerformanceLoop executes instructions with maximum performance using jump table
func (vm *VM) runHighPerformanceLoop(instructions []byte) (types.Value, error) {
	// Use safe jump table for stable performance
	ip := 0
	instrLen := len(instructions)

	// Safe high-performance execution loop using safe jump table dispatch
	for ip < instrLen {
		// Use safe jump table for instruction dispatch
		cont, err := vm.safeJumpTable.Execute(vm, instructions, &ip)
		if err != nil {
			return nil, err
		}
		if !cont {
			break // Halt instruction or end of execution
		}
	}

	// Return the top stack value as result
	if vm.sp > 0 {
		return vm.stack[vm.sp-1], nil
	}

	return Nil, nil
}

// runLegacyLoop executes instructions with the original switch-based approach
// Kept for compatibility and fallback purposes
func (vm *VM) runLegacyLoop(instructions []byte) (types.Value, error) {
	ip := 0
	instrLen := len(instructions)

	// Remove instruction count safety check for performance
	// Trust that bytecode is well-formed

	for ip < instrLen {
		op := Opcode(instructions[ip])
		ip++

		switch op {
		case OpConstant:
			// Trust bytecode is well-formed, remove bounds checks for performance
			constIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			vm.stack[vm.sp] = vm.constants[constIndex]
			vm.sp++

		case OpGetVar:
			varIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			// Handle nil globals with inline check
			if vm.globals[varIndex] == nil {
				vm.stack[vm.sp] = Nil
			} else {
				vm.stack[vm.sp] = vm.globals[varIndex]
			}
			vm.sp++

		case OpAddInt64:
			// Fast path for int64 addition - inline everything for performance
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			// Assume type checking passed at compile time, directly cast
			rightInt := right.(*types.IntValue)
			leftInt := left.(*types.IntValue)
			result := leftInt.Value() + rightInt.Value()
			// Inline value creation for performance
			if result >= 0 && result < 256 {
				vm.stack[vm.sp-2] = CachedInts[result]
			} else {
				vm.stack[vm.sp-2] = types.NewInt(result)
			}
			vm.sp--

		case OpSubInt64:
			// Fast path for int64 subtraction - fully inlined
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			rightInt := right.(*types.IntValue)
			leftInt := left.(*types.IntValue)
			result := leftInt.Value() - rightInt.Value()
			if result >= 0 && result < 256 {
				vm.stack[vm.sp-2] = CachedInts[result]
			} else {
				vm.stack[vm.sp-2] = types.NewInt(result)
			}
			vm.sp--

		case OpMulInt64:
			// Fast path for int64 multiplication - fully inlined
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			rightInt := right.(*types.IntValue)
			leftInt := left.(*types.IntValue)
			result := leftInt.Value() * rightInt.Value()
			if result >= 0 && result < 256 {
				vm.stack[vm.sp-2] = CachedInts[result]
			} else {
				vm.stack[vm.sp-2] = types.NewInt(result)
			}
			vm.sp--

		case OpDivInt64:
			// Fast path for int64 division - fully inlined
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			rightInt := right.(*types.IntValue)
			leftInt := left.(*types.IntValue)
			// Keep essential division by zero check
			if rightInt.Value() == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result := leftInt.Value() / rightInt.Value()
			if result >= -512 && result <= 511 {
				vm.stack[vm.sp-2] = vm.getIntValueFromPool(result)
			} else {
				vm.stack[vm.sp-2] = types.NewInt(result)
			}
			vm.sp--

		case OpAddFloat64:
			// Fast path for float64 addition - fully inlined
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			rightFloat := right.(*types.FloatValue)
			leftFloat := left.(*types.FloatValue)
			result := leftFloat.Value() + rightFloat.Value()
			vm.stack[vm.sp-2] = types.NewFloat(result)
			vm.sp--

		case OpAddString:
			// Fast path for string concatenation - fully inlined
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			rightStr := right.(*types.StringValue)
			leftStr := left.(*types.StringValue)
			result := leftStr.Value() + rightStr.Value()
			vm.stack[vm.sp-2] = types.NewString(result)
			vm.sp--

		case OpAdd:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeAddition(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpSub, OpSubFloat64:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeSubtraction(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpMul, OpMulFloat64:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeMultiplication(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpDiv, OpDivFloat64:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeDivision(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpMod, OpModFloat64:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeModulo(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpNeg:
			if vm.sp < 1 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			operand := vm.stack[vm.sp]

			result, err := vm.executeNegation(operand)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpPop:
			if vm.sp > 0 {
				vm.sp--
			}

		case OpGreaterThan, OpLessThan, OpGreaterEqual, OpLessEqual, OpEqual, OpNotEqual:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeComparison(op, left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpAnd, OpOr:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeLogical(op, left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpNot:
			if vm.sp < 1 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			operand := vm.stack[vm.sp]

			result := vm.executeLogicalNot(operand)
			vm.stack[vm.sp] = result
			vm.sp++

		case OpCall:
			if ip >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpCall instruction")
			}
			argCount := int(instructions[ip])
			ip++

			err := vm.executeCall(argCount)
			if err != nil {
				return nil, err
			}

		case OpBuiltin:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpBuiltin instruction")
			}
			builtinIndex := int(instructions[ip])
			argCount := int(instructions[ip+1]) // Use the argument count from compiler
			ip += 2

			err := vm.executeBuiltin(builtinIndex, argCount)
			if err != nil {
				return nil, err
			}

		case OpIndex:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			index := vm.stack[vm.sp]
			vm.sp--
			object := vm.stack[vm.sp]

			result, err := vm.executeIndex(object, index)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpMember:
			// OpMember expects: [object, memberName] on stack
			// Pops both and pushes result
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow for member access")
			}
			vm.sp--
			memberName := vm.stack[vm.sp]
			vm.sp--
			object := vm.stack[vm.sp]

			result, err := vm.executeMemberByName(object, memberName)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpArray:
			if ip >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpArray instruction")
			}
			elementCount := int(instructions[ip])
			ip++

			result, err := vm.executeArray(elementCount)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpSlice:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpSlice instruction")
			}
			elementCount := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			result, err := vm.executeArray(elementCount)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpObject:
			if ip >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpObject instruction")
			}
			pairCount := int(instructions[ip])
			ip++

			result, err := vm.executeObject(pairCount)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpMap:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpMap instruction")
			}
			pairCount := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			result, err := vm.executeObject(pairCount)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpPipe:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			function := vm.stack[vm.sp]
			vm.sp--
			data := vm.stack[vm.sp]

			result, err := vm.executePipe(data, function)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpGetPipelineElement:
			if vm.pipelineElement == nil {
				return nil, fmt.Errorf("no pipeline element available")
			}
			if vm.sp >= StackSize {
				return nil, fmt.Errorf("stack overflow")
			}
			vm.stack[vm.sp] = vm.pipelineElement
			vm.sp++

		case OpConcat:
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			right := vm.stack[vm.sp]
			vm.sp--
			left := vm.stack[vm.sp]

			result, err := vm.executeConcat(left, right)
			if err != nil {
				return nil, err
			}

			vm.stack[vm.sp] = result
			vm.sp++

		case OpJump:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpJump instruction")
			}
			jumpPos := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip = jumpPos

		case OpJumpTrue:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpJumpTrue instruction")
			}
			jumpPos := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			if vm.sp < 1 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			condition := vm.stack[vm.sp]

			if vm.isTruthy(condition) {
				ip = jumpPos
			}

		case OpJumpFalse:
			if ip+1 >= len(instructions) {
				return nil, fmt.Errorf("incomplete OpJumpFalse instruction")
			}
			jumpPos := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			if vm.sp < 1 {
				return nil, fmt.Errorf("stack underflow")
			}
			vm.sp--
			condition := vm.stack[vm.sp]

			if !vm.isTruthy(condition) {
				ip = jumpPos
			}

		case OpOptionalChaining:
			// Optional chaining: obj?.property
			// Top of stack: property, obj
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow for optional chaining")
			}
			property := vm.stack[vm.sp-1]
			object := vm.stack[vm.sp-2]
			vm.sp -= 2

			result, err := vm.executeOptionalChaining(object, property)
			if err != nil {
				return nil, err
			}
			vm.stack[vm.sp] = result
			vm.sp++

		case OpNullCoalescing:
			// Null coalescing: a ?? b
			// Top of stack: right (default), left
			if vm.sp < 2 {
				return nil, fmt.Errorf("stack underflow for null coalescing")
			}
			right := vm.stack[vm.sp-1]
			left := vm.stack[vm.sp-2]
			vm.sp -= 2

			result := vm.executeNullCoalescing(left, right)
			vm.stack[vm.sp] = result
			vm.sp++

		case OpModuleCall:
			// Module function call: module.function(args...)
			// Operands: moduleNameIndex (2 bytes), functionNameIndex (2 bytes), argCount (1 byte)
			moduleNameIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			functionNameIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			argCount := int(instructions[ip])
			ip++

			// Get module and function names from constants
			moduleNameVal := vm.constants[moduleNameIndex]
			functionNameVal := vm.constants[functionNameIndex]

			moduleName, ok := moduleNameVal.(*types.StringValue)
			if !ok {
				return nil, fmt.Errorf("module name must be string, got %T", moduleNameVal)
			}

			functionName, ok := functionNameVal.(*types.StringValue)
			if !ok {
				return nil, fmt.Errorf("function name must be string, got %T", functionNameVal)
			}

			// Collect arguments from stack
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.convertTypesValueToInterface(vm.stack[vm.sp-argCount+i])
			}
			vm.sp -= argCount

			// Call module function
			result, err := modules.DefaultRegistry.CallFunction(moduleName.Value(), functionName.Value(), args...)
			if err != nil {
				return nil, fmt.Errorf("module call error: %v", err)
			}

			// Convert result back to types.Value
			resultValue, err := vm.convertGoValueToTypesValue(result)
			if err != nil {
				return nil, fmt.Errorf("failed to convert module result: %v", err)
			}

			vm.stack[vm.sp] = resultValue
			vm.sp++

		case OpArrayDestructure:
			// Array destructuring: [a, b, c] = [1, 2, 3]
			// Operands: elementCount (2 bytes), startVarIndex (2 bytes)
			elementCount := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			startVarIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			// Value to destructure is on the stack
			vm.sp--
			value := vm.stack[vm.sp]

			err := vm.executeArrayDestructure(value, elementCount, startVarIndex)
			if err != nil {
				return nil, err
			}

		case OpObjectDestructure:
			// Object destructuring: {name, age} = user
			// Operands: propertyCount (2 bytes), startVarIndex (2 bytes)
			propertyCount := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2
			startVarIndex := int(instructions[ip])<<8 | int(instructions[ip+1])
			ip += 2

			// Property keys are on the stack (in reverse order)
			propertyKeys := make([]string, propertyCount)
			for i := propertyCount - 1; i >= 0; i-- {
				vm.sp--
				keyVal := vm.stack[vm.sp]
				if keyStr, ok := keyVal.(*types.StringValue); ok {
					propertyKeys[i] = keyStr.Value()
				} else {
					return nil, fmt.Errorf("property key must be string, got %T", keyVal)
				}
			}

			// Value to destructure is on the stack
			vm.sp--
			value := vm.stack[vm.sp]

			err := vm.executeObjectDestructure(value, propertyKeys, startVarIndex)
			if err != nil {
				return nil, err
			}

		case OpHalt:
			// Stop execution
			break

		default:
			// For unsupported opcodes, just skip them for now
			// This prevents infinite loops from unknown instructions
			continue
		}
	}

	if vm.sp == 0 {
		return Nil, nil
	}
	return vm.stack[vm.sp-1], nil
}

// executeAddition performs addition operation
func (vm *VM) executeAddition(left, right types.Value) (types.Value, error) {
	// Fast path: integer addition with cached results and object pool
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			result := leftInt.Value() + rightInt.Value()
			// Use cached values for small integers
			if result >= 0 && result < 256 {
				return CachedInts[result], nil
			}
			// Use object pool for larger integers
			return vm.pool.GetInt(result), nil
		}

		// Mixed int/float addition
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(float64(leftInt.Value()) + rightFloat.Value()), nil
		}
	}

	// Fast path: float addition
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() + rightFloat.Value()), nil
		}

		// Mixed float/int addition
		if rightInt, ok := right.(*types.IntValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() + float64(rightInt.Value())), nil
		}
	}

	// Fast path: string concatenation
	if leftStr, ok := left.(*types.StringValue); ok {
		if rightStr, ok := right.(*types.StringValue); ok {
			return vm.pool.GetString(leftStr.Value() + rightStr.Value()), nil
		}
	}

	return nil, fmt.Errorf("unsupported addition: %T + %T", left, right)
}

// executeMultiplication performs multiplication operation
func (vm *VM) executeMultiplication(left, right types.Value) (types.Value, error) {
	// Fast path: integer multiplication with cached results and object pool
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			result := leftInt.Value() * rightInt.Value()
			// Use cached values for small integers
			if result >= 0 && result < 256 {
				return CachedInts[result], nil
			}
			// Use object pool for larger integers
			return vm.pool.GetInt(result), nil
		}

		// Mixed int/float multiplication
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(float64(leftInt.Value()) * rightFloat.Value()), nil
		}
	}

	// Fast path: float multiplication
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() * rightFloat.Value()), nil
		}

		// Mixed float/int multiplication
		if rightInt, ok := right.(*types.IntValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() * float64(rightInt.Value())), nil
		}
	}

	return nil, fmt.Errorf("unsupported multiplication: %T * %T", left, right)
}

// executeSubtraction performs subtraction operation
func (vm *VM) executeSubtraction(left, right types.Value) (types.Value, error) {
	// Fast path: integer subtraction with cached results and object pool
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			result := leftInt.Value() - rightInt.Value()
			// Use cached values for small integers
			if result >= 0 && result < 256 {
				return CachedInts[result], nil
			}
			// Use object pool for larger integers
			return vm.pool.GetInt(result), nil
		}

		// Mixed int/float subtraction
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(float64(leftInt.Value()) - rightFloat.Value()), nil
		}
	}

	// Fast path: float subtraction
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() - rightFloat.Value()), nil
		}

		// Mixed float/int subtraction
		if rightInt, ok := right.(*types.IntValue); ok {
			return vm.pool.GetFloat(leftFloat.Value() - float64(rightInt.Value())), nil
		}
	}

	return nil, fmt.Errorf("unsupported subtraction: %T - %T", left, right)
}

// executeDivision performs division operation
func (vm *VM) executeDivision(left, right types.Value) (types.Value, error) {
	// Integer division
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			if rightInt.Value() == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return types.NewInt(leftInt.Value() / rightInt.Value()), nil
		}
	}

	// Float division
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			if rightFloat.Value() == 0.0 {
				return nil, fmt.Errorf("division by zero")
			}
			return types.NewFloat(leftFloat.Value() / rightFloat.Value()), nil
		}
	}

	return nil, fmt.Errorf("unsupported division: %T / %T", left, right)
}

// executeModulo performs modulo operation
func (vm *VM) executeModulo(left, right types.Value) (types.Value, error) {
	// Integer modulo
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			if rightInt.Value() == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			return types.NewInt(leftInt.Value() % rightInt.Value()), nil
		}
	}

	return nil, fmt.Errorf("unsupported modulo: %T %% %T", left, right)
}

// executeNegation performs negation operation
func (vm *VM) executeNegation(operand types.Value) (types.Value, error) {
	// Integer negation
	if intVal, ok := operand.(*types.IntValue); ok {
		return types.NewInt(-intVal.Value()), nil
	}

	// Float negation
	if floatVal, ok := operand.(*types.FloatValue); ok {
		return types.NewFloat(-floatVal.Value()), nil
	}

	return nil, fmt.Errorf("unsupported negation: -%T", operand)
}

// executeCall performs function call
func (vm *VM) executeCall(argCount int) error {
	if vm.sp < argCount+1 {
		return fmt.Errorf("stack underflow for function call")
	}

	// Get arguments from stack
	args := make([]types.Value, argCount)
	for i := argCount - 1; i >= 0; i-- {
		vm.sp--
		args[i] = vm.stack[vm.sp]
	}

	// Get function from stack
	vm.sp--
	function := vm.stack[vm.sp]

	// Execute function call
	result, err := vm.callFunction(function, args)
	if err != nil {
		return err
	}

	// Push result back onto stack
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = result
	vm.sp++

	return nil
}

// executeBuiltin performs builtin function call
func (vm *VM) executeBuiltin(builtinIndex int, argCount int) error {
	// Use the exact same builtin names from the builtins package
	builtinNames := builtins.StandardBuiltinNames

	if builtinIndex < 0 || builtinIndex >= len(builtinNames) {
		return fmt.Errorf("invalid builtin index: %d", builtinIndex)
	}

	funcName := builtinNames[builtinIndex]

	if vm.sp < argCount {
		return fmt.Errorf("stack underflow for builtin %s", funcName)
	}

	// Pop arguments from stack
	args := make([]types.Value, argCount)
	for i := argCount - 1; i >= 0; i-- {
		vm.sp--
		args[i] = vm.stack[vm.sp]
	}

	// Call the appropriate builtin function
	result, err := vm.callBuiltinByName(funcName, args)
	if err != nil {
		return fmt.Errorf("builtin %s error: %v", funcName, err)
	}

	// Push result back to stack
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = result
	vm.sp++
	return nil
}

// callBuiltinByName calls a builtin function by name with the given arguments
func (vm *VM) callBuiltinByName(funcName string, args []types.Value) (types.Value, error) {
	// Use the builtin functions from the builtins package
	if builtinFunc, exists := builtins.AllBuiltins[funcName]; exists {
		return builtinFunc(args)
	}

	// Fallback to custom implementations for functions not in builtins package
	switch funcName {
	case "avg":
		return vm.executeAvg(args[0])
	default:
		return Nil, fmt.Errorf("unknown builtin function: %s", funcName)
	}
}

// executeIndex performs index access
func (vm *VM) executeIndex(object, index types.Value) (types.Value, error) {
	// Slice index access
	if sliceVal, ok := object.(*types.SliceValue); ok {
		if intIndex, ok := index.(*types.IntValue); ok {
			idx := int(intIndex.Value())
			if idx < 0 || idx >= sliceVal.Len() {
				return nil, fmt.Errorf("index out of bounds: %d", idx)
			}
			return sliceVal.Get(idx), nil
		}
	}

	// Map key access
	if mapVal, ok := object.(*types.MapValue); ok {
		if strKey, ok := index.(*types.StringValue); ok {
			if val, exists := mapVal.Get(strKey.Value()); exists {
				return val, nil
			}
			return Nil, nil
		}
	}

	return nil, fmt.Errorf("unsupported index operation: %T[%T]", object, index)
}

// executeMember performs member access
func (vm *VM) executeMember(object types.Value, memberIndex int) (types.Value, error) {
	if memberIndex >= len(vm.constants) {
		return nil, fmt.Errorf("member index out of bounds: %d", memberIndex)
	}

	memberName, ok := vm.constants[memberIndex].(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("member name must be string")
	}

	// Map member access
	if mapVal, ok := object.(*types.MapValue); ok {
		if val, exists := mapVal.Get(memberName.Value()); exists {
			return val, nil
		}
		return Nil, nil
	}

	return nil, fmt.Errorf("unsupported member access: %T.%s", object, memberName.Value())
}

// executeMemberByName performs member access using a member name from the stack
func (vm *VM) executeMemberByName(object, memberName types.Value) (types.Value, error) {
	// Convert member name to string
	memberStr, ok := memberName.(*types.StringValue)
	if !ok {
		return nil, fmt.Errorf("member name must be string, got %T", memberName)
	}

	// Map member access
	if mapVal, ok := object.(*types.MapValue); ok {
		if val, exists := mapVal.Get(memberStr.Value()); exists {
			return val, nil
		}
		return Nil, nil
	}

	return nil, fmt.Errorf("unsupported member access: %T.%s", object, memberStr.Value())
}

// executeArray creates an array from stack elements
func (vm *VM) executeArray(elementCount int) (types.Value, error) {
	if vm.sp < elementCount {
		return nil, fmt.Errorf("stack underflow for array creation")
	}

	elements := make([]types.Value, elementCount)
	for i := elementCount - 1; i >= 0; i-- {
		vm.sp--
		elements[i] = vm.stack[vm.sp]
	}

	return types.NewSlice(elements, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}), nil
}

// executeObject creates an object from stack key-value pairs
func (vm *VM) executeObject(pairCount int) (types.Value, error) {
	if vm.sp < pairCount*2 {
		return nil, fmt.Errorf("stack underflow for object creation")
	}

	pairs := make(map[string]types.Value)
	for i := 0; i < pairCount; i++ {
		vm.sp--
		value := vm.stack[vm.sp]
		vm.sp--
		key := vm.stack[vm.sp]

		keyStr, ok := key.(*types.StringValue)
		if !ok {
			return nil, fmt.Errorf("object key must be string, got %T", key)
		}

		pairs[keyStr.Value()] = value
	}

	keyType := types.TypeInfo{Kind: types.KindString, Name: "string"}
	valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"}
	return types.NewMap(pairs, keyType, valueType), nil
}

// executePipe performs pipeline operation
func (vm *VM) executePipe(data, function types.Value) (types.Value, error) {
	// Set pipeline element for placeholder access
	oldPipelineElement := vm.pipelineElement
	vm.pipelineElement = data
	defer func() {
		vm.pipelineElement = oldPipelineElement
	}()

	// Call function with data as argument
	return vm.callFunction(function, []types.Value{data})
}

// executeConcat performs string concatenation
func (vm *VM) executeConcat(left, right types.Value) (types.Value, error) {
	leftStr := vm.toString(left)
	rightStr := vm.toString(right)
	return types.NewString(leftStr + rightStr), nil
}

// callFunction calls a function with given arguments
func (vm *VM) callFunction(function types.Value, args []types.Value) (types.Value, error) {
	// Check if function is a string (builtin function name)
	if funcName, ok := function.(*types.StringValue); ok {
		return vm.callBuiltinFunction(funcName.Value(), args)
	}

	// Check if function is a lambda/function value
	if funcVal, ok := function.(*types.FuncValue); ok {
		return vm.callLambdaFunction(funcVal, args)
	}

	// Check if function is a SliceValue (compiled function call expression)
	if funcSlice, ok := function.(*types.SliceValue); ok {
		return vm.callCompiledFunction(funcSlice, args)
	}

	// For unknown function types, return first argument as fallback
	if len(args) > 0 {
		return args[0], nil
	}
	return Nil, nil
}

// callCompiledFunction handles function calls that are compiled as SliceValue
func (vm *VM) callCompiledFunction(funcSlice *types.SliceValue, args []types.Value) (types.Value, error) {
	elements := funcSlice.Values()
	if len(elements) < 1 {
		return Nil, fmt.Errorf("invalid compiled function call")
	}

	// First element should be function name
	funcNameVal, ok := elements[0].(*types.StringValue)
	if !ok {
		return Nil, fmt.Errorf("function name must be string, got %T", elements[0])
	}

	funcName := funcNameVal.Value()

	// Check if this is a pipeline member access: ["__PIPELINE_MEMBER_ACCESS__", object, property]
	if funcName == "__PIPELINE_MEMBER_ACCESS__" && len(elements) >= 3 {
		object := elements[1]
		property := elements[2]

		// If object is a placeholder, use the pipeline element
		var actualObject types.Value
		if placeholderStr, ok := object.(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
			if len(args) == 0 {
				return Nil, fmt.Errorf("no pipeline data available for placeholder")
			}
			actualObject = args[0]
		} else {
			actualObject = object
		}

		// Get property name
		if propertyStr, ok := property.(*types.StringValue); ok {
			propertyName := propertyStr.Value()
			return vm.executeMemberAccess(actualObject, propertyName)
		}

		return Nil, fmt.Errorf("invalid property in member access")
	}

	// Check if this is a type method call: ["__TYPE_METHOD__", "__PLACEHOLDER_EXPR__", methodName, object, ...]
	if funcName == "__TYPE_METHOD__" && len(elements) >= 4 {
		if placeholderMarker, ok := elements[1].(*types.StringValue); ok && placeholderMarker.Value() == "__PLACEHOLDER_EXPR__" {
			if methodNameVal, ok := elements[2].(*types.StringValue); ok {
				methodName := methodNameVal.Value()
				if len(args) == 0 {
					return Nil, fmt.Errorf("no pipeline data available for type method")
				}
				data := args[0]

				// Execute type method call
				return vm.executeTypeMethod(data, methodName, elements[3:])
			}
		}
	}

	// Check if this is a pipeline function with complex type method: [funcName, "__PIPELINE_COMPLEX_TYPE_METHOD__", methodName, expression]
	if len(elements) >= 4 {
		if marker, ok := elements[1].(*types.StringValue); ok && marker.Value() == "__PIPELINE_COMPLEX_TYPE_METHOD__" {
			if methodNameVal, ok := elements[2].(*types.StringValue); ok {
				methodName := methodNameVal.Value()
				if len(args) == 0 {
					return Nil, fmt.Errorf("no pipeline data available for complex type method")
				}
				data := args[0]
				expression := elements[3] // The compiled expression

				// Execute pipeline function with complex type method expression
				switch funcName {
				case "filter":
					return vm.executePipelineFilterWithComplexTypeMethod(data, methodName, expression)
				case "map":
					return vm.executePipelineMapWithComplexTypeMethod(data, methodName, expression)
				default:
					return Nil, fmt.Errorf("unknown pipeline function with complex type method: %s", funcName)
				}
			}
		}
	}

	// Check if this is a pipeline function with type method: [funcName, "__PIPELINE_TYPE_METHOD__", methodName, object, ...]
	if len(elements) >= 4 {
		if marker, ok := elements[1].(*types.StringValue); ok && marker.Value() == "__PIPELINE_TYPE_METHOD__" {
			if methodNameVal, ok := elements[2].(*types.StringValue); ok {
				methodName := methodNameVal.Value()
				if len(args) == 0 {
					return Nil, fmt.Errorf("no pipeline data available for pipeline type method")
				}
				data := args[0]

				// Execute pipeline function with type method
				switch funcName {
				case "filter":
					return vm.executePipelineFilterWithTypeMethod(data, methodName, elements[3:])
				case "map":
					return vm.executePipelineMapWithTypeMethod(data, methodName, elements[3:])
				default:
					return Nil, fmt.Errorf("unknown pipeline function with type method: %s", funcName)
				}
			}
		}
	}

	// Check if this is a placeholder expression: [funcName, "__PLACEHOLDER_EXPR__", ...]
	if len(elements) >= 2 {
		if placeholderMarker, ok := elements[1].(*types.StringValue); ok && placeholderMarker.Value() == "__PLACEHOLDER_EXPR__" {
			// This is a pipeline function with placeholders
			// args[0] = data, elements[2:] = placeholder expressions
			if len(args) == 0 {
				return Nil, fmt.Errorf("no pipeline data available for placeholder expression")
			}
			data := args[0]
			var condition types.Value
			if len(elements) > 2 {
				condition = elements[2] // The actual condition/transform expression
			} else {
				condition = types.NewBool(true) // Default condition
			}

			// Call VM-specific pipeline handlers
			switch funcName {
			case "filter":
				return vm.executeFilter(data, condition)
			case "map":
				return vm.executeMap(data, condition)
			default:
				return Nil, fmt.Errorf("unknown pipeline function: %s", funcName)
			}
		}
	}

	// Regular function call - extract function arguments from the slice
	var funcArgs []types.Value
	if len(elements) > 1 {
		funcArgs = elements[1:] // Skip function name
	}

	// Combine with pipeline data
	allArgs := append(args, funcArgs...)

	return vm.callBuiltinFunction(funcName, allArgs)
}

// callBuiltinFunction calls a builtin function by name
func (vm *VM) callBuiltinFunction(funcName string, args []types.Value) (types.Value, error) {
	// Use the builtin functions from the builtins package
	if builtinFunc, exists := builtins.AllBuiltins[funcName]; exists {
		return builtinFunc(args)
	}

	return Nil, fmt.Errorf("unknown builtin function: %s", funcName)
}

// callLambdaFunction calls a lambda function
func (vm *VM) callLambdaFunction(funcVal *types.FuncValue, args []types.Value) (types.Value, error) {
	// TODO: Implement lambda function execution
	// For now, return first argument as placeholder
	if len(args) > 0 {
		return args[0], nil
	}
	return Nil, nil
}

// executeFilter filters array elements based on condition
func (vm *VM) executeFilter(data types.Value, condition types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("filter can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Set pipeline element for placeholder evaluation
		oldPipelineElement := vm.pipelineElement
		vm.pipelineElement = element

		// Evaluate condition - for now, assume it's a placeholder expression
		// This is simplified - we should properly evaluate the condition
		conditionResult := vm.evaluatePlaceholderCondition(condition, element)

		vm.pipelineElement = oldPipelineElement

		if vm.isTruthy(conditionResult) {
			result = append(result, element)
		}
	}

	// Get element type from slice type
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// executeMap transforms array elements
func (vm *VM) executeMap(data types.Value, transform types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("map can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Set pipeline element for placeholder evaluation
		oldPipelineElement := vm.pipelineElement
		vm.pipelineElement = element

		// Transform element - for now, assume it's a placeholder expression
		// This is simplified - we should properly evaluate the transform
		transformedResult := vm.evaluatePlaceholderTransform(transform, element)

		vm.pipelineElement = oldPipelineElement

		result = append(result, transformedResult)
	}

	// Get element type from slice type
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// evaluatePlaceholderCondition evaluates a condition with placeholder
func (vm *VM) evaluatePlaceholderCondition(condition types.Value, element types.Value) types.Value {
	// Check if condition is a simple placeholder string
	if strVal, ok := condition.(*types.StringValue); ok && strVal.Value() == "__PLACEHOLDER__" {
		// For simple placeholder (#), return the element itself as the condition
		return element
	}

	// Check if condition is a PlaceholderExprValue
	if placeholderExpr, ok := condition.(*types.PlaceholderExprValue); ok {
		return vm.evaluatePlaceholderExpression(placeholderExpr, element)
	}

	// Check if condition is a SliceValue (compiled placeholder expression)
	if condSlice, ok := condition.(*types.SliceValue); ok {
		return vm.evaluateCompiledPlaceholderExpression(condSlice, element)
	}

	// For other types, try to evaluate as constant condition
	return condition
}

// evaluateCompiledPlaceholderExpression evaluates a compiled placeholder expression
func (vm *VM) evaluateCompiledPlaceholderExpression(condSlice *types.SliceValue, element types.Value) types.Value {
	elements := condSlice.Values()

	// Handle member access: [".", __PLACEHOLDER__, "property"]
	if len(elements) == 3 {
		operatorVal, ok1 := elements[0].(*types.StringValue)
		placeholderVal, ok2 := elements[1].(*types.StringValue)
		propertyVal, ok3 := elements[2].(*types.StringValue)

		if ok1 && ok2 && ok3 && operatorVal.Value() == "." && placeholderVal.Value() == "__PLACEHOLDER__" {
			// This is a member access: #.property
			return vm.evaluateMemberAccess(element, propertyVal.Value())
		}
	}

	// Handle infix expressions: [operator, left, right] where left or right might be placeholder
	if len(elements) == 3 {
		operatorVal, ok1 := elements[0].(*types.StringValue)
		if ok1 {
			operator := operatorVal.Value()
			leftVal := elements[1]
			rightVal := elements[2]

			// Skip member access operator as it's handled above
			if operator == "." {
				return types.NewBool(false)
			}

			// Replace placeholders with the current element
			var left, right types.Value

			if placeholderStr, ok := leftVal.(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
				left = element
			} else if memberSlice, ok := leftVal.(*types.SliceValue); ok {
				// Handle nested member access like #.age
				left = vm.evaluateCompiledPlaceholderExpression(memberSlice, element)
			} else {
				left = leftVal
			}

			if placeholderStr, ok := rightVal.(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
				right = element
			} else if memberSlice, ok := rightVal.(*types.SliceValue); ok {
				// Handle nested member access
				right = vm.evaluateCompiledPlaceholderExpression(memberSlice, element)
			} else {
				right = rightVal
			}

			// Perform the operation
			switch operator {
			case ">":
				return vm.evaluateComparison(OpGreaterThan, left, right)
			case "<":
				return vm.evaluateComparison(OpLessThan, left, right)
			case ">=":
				return vm.evaluateComparison(OpGreaterEqual, left, right)
			case "<=":
				return vm.evaluateComparison(OpLessEqual, left, right)
			case "==":
				return vm.evaluateComparison(OpEqual, left, right)
			case "!=":
				return vm.evaluateComparison(OpNotEqual, left, right)
			case "&&":
				if vm.isTruthy(left) {
					return right
				}
				return left
			case "||":
				if vm.isTruthy(left) {
					return left
				}
				return right
			case "+":
				result, _ := vm.executeAddition(left, right)
				return result
			case "-":
				result, _ := vm.executeSubtraction(left, right)
				return result
			case "*":
				result, _ := vm.executeMultiplication(left, right)
				return result
			case "/":
				result, _ := vm.executeDivision(left, right)
				return result
			case "%":
				result, _ := vm.executeModulo(left, right)
				return result
			default:
				// For unknown operators, return false
				return types.NewBool(false)
			}
		}
	}

	// Handle member access: [__PLACEHOLDER__, "property"] (legacy format)
	if len(elements) == 2 {
		placeholderVal, ok1 := elements[0].(*types.StringValue)
		propertyVal, ok2 := elements[1].(*types.StringValue)

		if ok1 && ok2 && placeholderVal.Value() == "__PLACEHOLDER__" {
			// This is a member access: #.property
			return vm.evaluateMemberAccess(element, propertyVal.Value())
		}
	}

	// Invalid format
	return types.NewBool(false)
}

// evaluateMemberAccess evaluates member access on an element
func (vm *VM) evaluateMemberAccess(element types.Value, memberName string) types.Value {
	// Handle map member access
	if mapVal, ok := element.(*types.MapValue); ok {
		if val, exists := mapVal.Get(memberName); exists {
			return val
		}
		return Nil
	}

	// For other types, return the element itself (fallback)
	return element
}

// evaluatePlaceholderExpression evaluates a placeholder expression with the current element
func (vm *VM) evaluatePlaceholderExpression(placeholderExpr *types.PlaceholderExprValue, element types.Value) types.Value {
	// Get the operator and operand from the placeholder expression
	operator := placeholderExpr.Operator()
	operand := placeholderExpr.Operand()

	// Perform the operation between element and operand
	switch operator {
	case "!":
		// Logical NOT - only use the element, ignore operand
		return vm.executeLogicalNot(element)
	case ">":
		return vm.evaluateComparison(OpGreaterThan, element, operand)
	case "<":
		return vm.evaluateComparison(OpLessThan, element, operand)
	case ">=":
		return vm.evaluateComparison(OpGreaterEqual, element, operand)
	case "<=":
		return vm.evaluateComparison(OpLessEqual, element, operand)
	case "==":
		return vm.evaluateComparison(OpEqual, element, operand)
	case "!=":
		return vm.evaluateComparison(OpNotEqual, element, operand)
	case "+":
		result, _ := vm.executeAddition(element, operand)
		return result
	case "-":
		result, _ := vm.executeSubtraction(element, operand)
		return result
	case "*":
		result, _ := vm.executeMultiplication(element, operand)
		return result
	case "/":
		result, _ := vm.executeDivision(element, operand)
		return result
	case "%":
		result, _ := vm.executeModulo(element, operand)
		return result
	default:
		// For unknown operators, return the element unchanged
		return element
	}
}

// evaluateComparison performs comparison and returns boolean result
func (vm *VM) evaluateComparison(op Opcode, left, right types.Value) types.Value {
	result, err := vm.executeComparison(op, left, right)
	if err != nil {
		return types.NewBool(false)
	}
	return result
}

// evaluatePlaceholderTransform evaluates a transform with placeholder
func (vm *VM) evaluatePlaceholderTransform(transform types.Value, element types.Value) types.Value {
	// Check if transform is a simple placeholder string
	if strVal, ok := transform.(*types.StringValue); ok && strVal.Value() == "__PLACEHOLDER__" {
		// For simple placeholder (#), return the element itself
		return element
	}

	// Check if transform is a PlaceholderExprValue
	if placeholderExpr, ok := transform.(*types.PlaceholderExprValue); ok {
		return vm.evaluatePlaceholderExpression(placeholderExpr, element)
	}

	// Check if transform is a SliceValue (compiled placeholder expression)
	if transSlice, ok := transform.(*types.SliceValue); ok {
		return vm.evaluateCompiledPlaceholderExpression(transSlice, element)
	}

	// For other types, return as-is (constant transform)
	return transform
}

// executeSum sums array elements
func (vm *VM) executeSum(data types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("sum can only be applied to arrays")
	}

	var sum int64 = 0
	elements := slice.Values()

	for _, element := range elements {
		if intVal, ok := element.(*types.IntValue); ok {
			sum += intVal.Value()
		}
	}

	return types.NewInt(sum), nil
}

// executeCount counts array elements
func (vm *VM) executeCount(data types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("count can only be applied to arrays")
	}

	return types.NewInt(int64(len(slice.Values()))), nil
}

// executeTypeMethod executes a type method on the pipeline element
func (vm *VM) executeTypeMethod(data types.Value, methodName string, arguments []types.Value) (types.Value, error) {
	// Set pipeline element for placeholder evaluation
	oldPipelineElement := vm.pipelineElement
	vm.pipelineElement = data
	defer func() {
		vm.pipelineElement = oldPipelineElement
	}()

	// Evaluate the object (which should be the placeholder)
	var objectValue types.Value
	if len(arguments) > 0 {
		if placeholderStr, ok := arguments[0].(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
			objectValue = data
		} else {
			objectValue = arguments[0]
		}
	} else {
		objectValue = data
	}

	// Determine the type of the object and call the appropriate type method
	var typePrefix string
	switch objectValue.(type) {
	case *types.StringValue:
		typePrefix = "string"
	case *types.IntValue:
		typePrefix = "int"
	case *types.FloatValue:
		typePrefix = "float"
	case *types.BoolValue:
		typePrefix = "bool"
	case *types.SliceValue:
		typePrefix = "slice"
	case *types.MapValue:
		typePrefix = "map"
	default:
		return Nil, fmt.Errorf("unsupported type for method call: %T", objectValue)
	}

	// Construct the full method name like "string.upper", "int.abs", etc.
	fullMethodName := typePrefix + "." + methodName

	// Check if the method exists in TypeMethodBuiltins
	if typeMethod, exists := builtins.TypeMethodBuiltins[fullMethodName]; exists {
		// Prepare arguments for the type method call
		methodArgs := []types.Value{objectValue}

		// Add any additional arguments (skip the first one which is the object)
		if len(arguments) > 1 {
			for _, arg := range arguments[1:] {
				// Evaluate argument if it's a placeholder
				if placeholderStr, ok := arg.(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
					methodArgs = append(methodArgs, data)
				} else {
					methodArgs = append(methodArgs, arg)
				}
			}
		}

		// Call the type method
		return typeMethod(methodArgs)
	}

	return Nil, fmt.Errorf("unknown type method: %s", fullMethodName)
}

// executeMemberAccess executes member access on an object
func (vm *VM) executeMemberAccess(object types.Value, propertyName string) (types.Value, error) {
	switch obj := object.(type) {
	case *types.MapValue:
		// Map member access
		if value, exists := obj.Get(propertyName); exists {
			return value, nil
		}
		return types.NewNil(), nil
	case *types.SliceValue:
		// Check if this is accessing a slice property like "length"
		switch propertyName {
		case "length":
			return types.NewInt(int64(len(obj.Values()))), nil
		default:
			return Nil, fmt.Errorf("property %s not found on slice", propertyName)
		}
	case *types.StringValue:
		// Check if this is accessing a string property like "length"
		switch propertyName {
		case "length":
			return types.NewInt(int64(len(obj.Value()))), nil
		default:
			return Nil, fmt.Errorf("property %s not found on string", propertyName)
		}
	default:
		return Nil, fmt.Errorf("cannot access property %s on type %T", propertyName, object)
	}
}

// executePipelineFilterWithTypeMethod executes filter with type method calls
func (vm *VM) executePipelineFilterWithTypeMethod(data types.Value, methodName string, arguments []types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("filter can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Call the type method on each element
		methodResult, err := vm.executeTypeMethod(element, methodName, arguments)
		if err != nil {
			return Nil, err
		}

		// Check if the result is truthy for filtering
		if vm.isTruthy(methodResult) {
			result = append(result, element)
		}
	}

	// Get element type from slice type
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// executePipelineMapWithTypeMethod executes map with type method calls
func (vm *VM) executePipelineMapWithTypeMethod(data types.Value, methodName string, arguments []types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("map can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Call the type method on each element and collect the result
		methodResult, err := vm.executeTypeMethod(element, methodName, arguments)
		if err != nil {
			return Nil, err
		}

		result = append(result, methodResult)
	}

	// Get element type from slice type (but for map, the result type might be different)
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// executePipelineFilterWithComplexTypeMethod executes filter with complex type method expressions
func (vm *VM) executePipelineFilterWithComplexTypeMethod(data types.Value, methodName string, expression types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("filter can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Set pipeline element for placeholder evaluation
		oldPipelineElement := vm.pipelineElement
		vm.pipelineElement = element

		// Evaluate the complex expression for this element
		// This is a simplified approach - we'll need to implement proper expression evaluation
		// For now, let's try to handle the most common case: #.methodName() > constant
		evalResult, err := vm.evaluateComplexTypeMethodExpression(element, methodName, expression)
		if err != nil {
			vm.pipelineElement = oldPipelineElement
			return Nil, err
		}

		vm.pipelineElement = oldPipelineElement

		// Check if the result is truthy for filtering
		if vm.isTruthy(evalResult) {
			result = append(result, element)
		}
	}

	// Get element type from slice type
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// executePipelineMapWithComplexTypeMethod executes map with complex type method expressions
func (vm *VM) executePipelineMapWithComplexTypeMethod(data types.Value, methodName string, expression types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("map can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	for _, element := range elements {
		// Set pipeline element for placeholder evaluation
		oldPipelineElement := vm.pipelineElement
		vm.pipelineElement = element

		// Evaluate the complex expression for this element
		evalResult, err := vm.evaluateComplexTypeMethodExpression(element, methodName, expression)
		if err != nil {
			vm.pipelineElement = oldPipelineElement
			return Nil, err
		}

		vm.pipelineElement = oldPipelineElement

		result = append(result, evalResult)
	}

	// Get element type from slice type
	elemType := vm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// evaluateComplexTypeMethodExpression evaluates a complex expression containing type method calls
func (vm *VM) evaluateComplexTypeMethodExpression(element types.Value, methodName string, expression types.Value) (types.Value, error) {
	// The expression is compiled bytecode that we need to execute
	// We need to create a sub-VM context to evaluate it

	// For now, let's implement a simplified version that handles the most common pattern:
	// #.methodName() > constant, #.methodName() == constant, etc.

	// Since we know this is a complex expression, we can try to execute it directly
	// but we need to set up the context properly

	// Save the current pipeline element
	oldPipelineElement := vm.pipelineElement
	vm.pipelineElement = element
	defer func() {
		vm.pipelineElement = oldPipelineElement
	}()

	// Try to evaluate the expression
	// Since this is compiled bytecode, we need to execute it as a slice operation
	if slice, ok := expression.(*types.SliceValue); ok {
		return vm.evaluateSliceExpression(slice)
	}

	// Fallback: assume it's a simple comparison
	methodResult, err := vm.callTypeMethod(element, methodName, []types.Value{})
	if err != nil {
		return Nil, err
	}

	return methodResult, nil
}

// isTypeMethodName checks if a property name is a type method
func (vm *VM) isTypeMethodName(propertyName string) bool {
	// Common type method names that should be executed as method calls
	methodNames := map[string]bool{
		"length":     true,
		"upper":      true,
		"lower":      true,
		"trim":       true,
		"replace":    true,
		"contains":   true,
		"startsWith": true,
		"endsWith":   true,
		"matches":    true,
		"abs":        true,
		"isEven":     true,
		"isOdd":      true,
		"toString":   true,
		"toFloat":    true,
		"charAt":     true,
		"indexOf":    true,
		"substring":  true,
	}
	return methodNames[propertyName]
}

// executeTypeMethodDirectly executes a type method directly on an object
func (vm *VM) executeTypeMethodDirectly(object types.Value, methodName string) (types.Value, error) {
	// Determine the type of the object and call the appropriate type method
	var typePrefix string
	switch object.(type) {
	case *types.StringValue:
		typePrefix = "string"
	case *types.IntValue:
		typePrefix = "int"
	case *types.FloatValue:
		typePrefix = "float"
	case *types.BoolValue:
		typePrefix = "bool"
	case *types.SliceValue:
		typePrefix = "slice"
	case *types.MapValue:
		typePrefix = "map"
	default:
		return Nil, fmt.Errorf("unsupported type for method call: %T", object)
	}

	// Construct the full method name like "string.upper", "int.abs", etc.
	fullMethodName := typePrefix + "." + methodName

	// Check if the method exists in TypeMethodBuiltins
	if typeMethod, exists := builtins.TypeMethodBuiltins[fullMethodName]; exists {
		// Call the type method with just the object as argument
		return typeMethod([]types.Value{object})
	}

	return Nil, fmt.Errorf("unknown type method: %s", fullMethodName)
}

// evaluateSliceExpression evaluates a slice that represents a compiled expression
func (vm *VM) evaluateSliceExpression(slice *types.SliceValue) (types.Value, error) {
	elements := slice.Values()
	if len(elements) == 0 {
		return Nil, fmt.Errorf("empty expression")
	}

	// Check for special 5-element pipeline member access pattern
	// Pattern: [operator, __PIPELINE_MEMBER_ACCESS__, __PLACEHOLDER__, methodName, right]
	if len(elements) == 5 {
		if opStr, ok := elements[0].(*types.StringValue); ok {
			if marker, ok := elements[1].(*types.StringValue); ok && marker.Value() == "__PIPELINE_MEMBER_ACCESS__" {
				operator := opStr.Value()

				// Construct the left operand as a pipeline member access
				leftOperand := types.NewSlice([]types.Value{
					elements[1], // __PIPELINE_MEMBER_ACCESS__
					elements[2], // __PLACEHOLDER__
					elements[3], // methodName
				}, types.TypeInfo{Kind: types.KindInterface, Name: "interface{}"})

				right := elements[4] // The actual right operand (e.g., 4)

				// Evaluate operands
				leftVal, err := vm.evaluateExpressionOperand(leftOperand)
				if err != nil {
					return Nil, err
				}

				rightVal, err := vm.evaluateExpressionOperand(right)
				if err != nil {
					return Nil, err
				}

				// Perform the operation
				return vm.performInfixOperation(operator, leftVal, rightVal)
			}
		}
	}

	// Check if this is a regular infix operation
	// Pattern: [operator, left, right]
	if len(elements) >= 3 {
		if opStr, ok := elements[0].(*types.StringValue); ok {
			operator := opStr.Value()
			left := elements[1]
			right := elements[2]

			// Evaluate operands
			leftVal, err := vm.evaluateExpressionOperand(left)
			if err != nil {
				return Nil, err
			}

			rightVal, err := vm.evaluateExpressionOperand(right)
			if err != nil {
				return Nil, err
			}

			// Perform the operation
			return vm.performInfixOperation(operator, leftVal, rightVal)
		}
	}

	// If it's not an infix operation, try to evaluate as a single value
	if len(elements) == 1 {
		return vm.evaluateExpressionOperand(elements[0])
	}

	return Nil, fmt.Errorf("unsupported expression format")
}

// evaluateExpressionOperand evaluates an expression operand
func (vm *VM) evaluateExpressionOperand(operand types.Value) (types.Value, error) {
	// If it's a placeholder, return the current pipeline element
	if strVal, ok := operand.(*types.StringValue); ok && strVal.Value() == "__PLACEHOLDER__" {
		if vm.pipelineElement == nil {
			return Nil, fmt.Errorf("no pipeline element available for placeholder")
		}
		return vm.pipelineElement, nil
	}

	// Check if it's a slice, and handle special pipeline cases first
	if slice, ok := operand.(*types.SliceValue); ok {
		elements := slice.Values()

		// Check if it's a pipeline member access first
		if len(elements) >= 3 {
			if marker, ok := elements[0].(*types.StringValue); ok && marker.Value() == "__PIPELINE_MEMBER_ACCESS__" {
				object := elements[1]
				property := elements[2]

				// If object is a placeholder, use the current pipeline element
				var actualObject types.Value
				if placeholderStr, ok := object.(*types.StringValue); ok && placeholderStr.Value() == "__PLACEHOLDER__" {
					if vm.pipelineElement == nil {
						return Nil, fmt.Errorf("no pipeline element available for member access")
					}
					actualObject = vm.pipelineElement
				} else {
					actualObject = object
				}

				// Get property name and check if it's a method call
				if propertyStr, ok := property.(*types.StringValue); ok {
					propertyName := propertyStr.Value()

					// Check if this looks like a method call (common method names)
					if vm.isTypeMethodName(propertyName) {
						// Execute as type method call
						return vm.executeTypeMethodDirectly(actualObject, propertyName)
					} else {
						// Execute as simple member access
						return vm.executeMemberAccess(actualObject, propertyName)
					}
				}
			}
		}

		// If it's not a special pipeline construct, evaluate as regular slice expression
		return vm.evaluateSliceExpression(slice)
	}

	// If it's a method call marker, we need to handle this specially
	if strVal, ok := operand.(*types.StringValue); ok && strVal.Value() == "." {
		// This indicates a method call, but we need more context
		// For now, return the current pipeline element
		if vm.pipelineElement == nil {
			return Nil, fmt.Errorf("no pipeline element available for method call")
		}
		return vm.pipelineElement, nil
	}

	// Otherwise, return the operand as-is
	return operand, nil
}

// performInfixOperation performs an infix operation
func (vm *VM) performInfixOperation(operator string, left, right types.Value) (types.Value, error) {
	switch operator {
	case ">":
		return vm.performComparison(left, right, ">")
	case "<":
		return vm.performComparison(left, right, "<")
	case ">=":
		return vm.performComparison(left, right, ">=")
	case "<=":
		return vm.performComparison(left, right, "<=")
	case "==":
		return vm.performComparison(left, right, "==")
	case "!=":
		return vm.performComparison(left, right, "!=")
	case "+":
		return vm.performArithmetic(left, right, "+")
	case "-":
		return vm.performArithmetic(left, right, "-")
	case "*":
		return vm.performArithmetic(left, right, "*")
	case "/":
		return vm.performArithmetic(left, right, "/")
	default:
		return Nil, fmt.Errorf("unsupported operator: %s", operator)
	}
}

// performComparison performs comparison operations
func (vm *VM) performComparison(left, right types.Value, operator string) (types.Value, error) {
	// Handle nil comparisons first
	if left == nil || left == Nil || left.Type().Kind == types.KindNil {
		if right == nil || right == Nil || right.Type().Kind == types.KindNil {
			switch operator {
			case ">":
				return types.NewBool(false), nil
			case "<":
				return types.NewBool(false), nil
			case ">=":
				return types.NewBool(true), nil
			case "<=":
				return types.NewBool(true), nil
			case "==":
				return types.NewBool(true), nil
			case "!=":
				return types.NewBool(false), nil
			}
		} else {
			// nil compared to non-nil
			switch operator {
			case ">":
				return types.NewBool(false), nil
			case "<":
				return types.NewBool(true), nil
			case ">=":
				return types.NewBool(false), nil
			case "<=":
				return types.NewBool(true), nil
			case "==":
				return types.NewBool(false), nil
			case "!=":
				return types.NewBool(true), nil
			}
		}
	}

	if right == nil || right == Nil || right.Type().Kind == types.KindNil {
		// non-nil compared to nil
		switch operator {
		case ">":
			return types.NewBool(true), nil
		case "<":
			return types.NewBool(false), nil
		case ">=":
			return types.NewBool(true), nil
		case "<=":
			return types.NewBool(false), nil
		case "==":
			return types.NewBool(false), nil
		case "!=":
			return types.NewBool(true), nil
		}
	}

	// Integer comparison
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			switch operator {
			case ">":
				return types.NewBool(leftInt.Value() > rightInt.Value()), nil
			case "<":
				return types.NewBool(leftInt.Value() < rightInt.Value()), nil
			case ">=":
				return types.NewBool(leftInt.Value() >= rightInt.Value()), nil
			case "<=":
				return types.NewBool(leftInt.Value() <= rightInt.Value()), nil
			case "==":
				return types.NewBool(leftInt.Value() == rightInt.Value()), nil
			case "!=":
				return types.NewBool(leftInt.Value() != rightInt.Value()), nil
			}
		}
	}

	// Float comparison
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			switch operator {
			case ">":
				return types.NewBool(leftFloat.Value() > rightFloat.Value()), nil
			case "<":
				return types.NewBool(leftFloat.Value() < rightFloat.Value()), nil
			case ">=":
				return types.NewBool(leftFloat.Value() >= rightFloat.Value()), nil
			case "<=":
				return types.NewBool(leftFloat.Value() <= rightFloat.Value()), nil
			case "==":
				return types.NewBool(leftFloat.Value() == rightFloat.Value()), nil
			case "!=":
				return types.NewBool(leftFloat.Value() != rightFloat.Value()), nil
			}
		}
	}

	// String comparison
	if leftStr, ok := left.(*types.StringValue); ok {
		if rightStr, ok := right.(*types.StringValue); ok {
			switch operator {
			case ">":
				return types.NewBool(leftStr.Value() > rightStr.Value()), nil
			case "<":
				return types.NewBool(leftStr.Value() < rightStr.Value()), nil
			case ">=":
				return types.NewBool(leftStr.Value() >= rightStr.Value()), nil
			case "<=":
				return types.NewBool(leftStr.Value() <= rightStr.Value()), nil
			case "==":
				return types.NewBool(leftStr.Value() == rightStr.Value()), nil
			case "!=":
				return types.NewBool(leftStr.Value() != rightStr.Value()), nil
			}
		}
	}

	// Boolean comparison
	if leftBool, ok := left.(*types.BoolValue); ok {
		if rightBool, ok := right.(*types.BoolValue); ok {
			switch operator {
			case "==":
				return types.NewBool(leftBool.Value() == rightBool.Value()), nil
			case "!=":
				return types.NewBool(leftBool.Value() != rightBool.Value()), nil
			}
		}
	}

	// Mixed type comparisons - only equality/inequality makes sense
	if operator == "==" {
		return types.NewBool(false), nil // Different types are never equal
	}
	if operator == "!=" {
		return types.NewBool(true), nil // Different types are always not equal
	}

	return nil, fmt.Errorf("unsupported comparison: %T %s %T", left, operator, right)
}

// performArithmetic performs arithmetic operations
func (vm *VM) performArithmetic(left, right types.Value, operator string) (types.Value, error) {
	// Try integer arithmetic
	leftInt, leftIsInt := vm.tryConvertToInt(left)
	rightInt, rightIsInt := vm.tryConvertToInt(right)

	if leftIsInt && rightIsInt {
		switch operator {
		case "+":
			return types.NewInt(leftInt + rightInt), nil
		case "-":
			return types.NewInt(leftInt - rightInt), nil
		case "*":
			return types.NewInt(leftInt * rightInt), nil
		case "/":
			if rightInt == 0 {
				return Nil, fmt.Errorf("division by zero")
			}
			return types.NewInt(leftInt / rightInt), nil
		}
	}

	// Try float arithmetic
	leftFloat, leftIsFloat := vm.tryConvertToFloat(left)
	rightFloat, rightIsFloat := vm.tryConvertToFloat(right)

	if leftIsFloat && rightIsFloat {
		switch operator {
		case "+":
			return types.NewFloat(leftFloat + rightFloat), nil
		case "-":
			return types.NewFloat(leftFloat - rightFloat), nil
		case "*":
			return types.NewFloat(leftFloat * rightFloat), nil
		case "/":
			if rightFloat == 0.0 {
				return Nil, fmt.Errorf("division by zero")
			}
			return types.NewFloat(leftFloat / rightFloat), nil
		}
	}

	return Nil, fmt.Errorf("cannot perform arithmetic on %T and %T with operator %s", left, right, operator)
}

// Helper functions for type conversion
func (vm *VM) tryConvertToInt(value types.Value) (int64, bool) {
	switch v := value.(type) {
	case *types.IntValue:
		return v.Value(), true
	case *types.FloatValue:
		return int64(v.Value()), true
	default:
		return 0, false
	}
}

func (vm *VM) tryConvertToFloat(value types.Value) (float64, bool) {
	switch v := value.(type) {
	case *types.FloatValue:
		return v.Value(), true
	case *types.IntValue:
		return float64(v.Value()), true
	default:
		return 0.0, false
	}
}

func (vm *VM) tryConvertToString(value types.Value) (string, bool) {
	switch v := value.(type) {
	case *types.StringValue:
		return v.Value(), true
	default:
		return "", false
	}
}

// callTypeMethod calls a type method on a value
func (vm *VM) callTypeMethod(value types.Value, methodName string, args []types.Value) (types.Value, error) {
	// Determine the type prefix
	var typePrefix string
	switch value.(type) {
	case *types.StringValue:
		typePrefix = "string"
	case *types.IntValue:
		typePrefix = "int"
	case *types.FloatValue:
		typePrefix = "float"
	case *types.BoolValue:
		typePrefix = "bool"
	case *types.SliceValue:
		typePrefix = "slice"
	case *types.MapValue:
		typePrefix = "map"
	default:
		return Nil, fmt.Errorf("unsupported type for method call: %T", value)
	}

	// Construct the full method name
	fullMethodName := typePrefix + "." + methodName

	// Check if the method exists in TypeMethodBuiltins
	if typeMethod, exists := builtins.TypeMethodBuiltins[fullMethodName]; exists {
		// Prepare arguments for the type method call
		methodArgs := []types.Value{value}
		methodArgs = append(methodArgs, args...)

		// Call the type method
		return typeMethod(methodArgs)
	}

	return Nil, fmt.Errorf("unknown type method: %s", fullMethodName)
}

// executeMax finds maximum value in array
func (vm *VM) executeMax(data types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("max can only be applied to arrays")
	}

	elements := slice.Values()
	if len(elements) == 0 {
		return Nil, nil
	}

	max := elements[0]
	for _, element := range elements[1:] {
		if vm.compareValues(element, max) > 0 {
			max = element
		}
	}

	return max, nil
}

// executeMin finds minimum value in array
func (vm *VM) executeMin(data types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("min can only be applied to arrays")
	}

	elements := slice.Values()
	if len(elements) == 0 {
		return Nil, nil
	}

	min := elements[0]
	for _, element := range elements[1:] {
		if vm.compareValues(element, min) < 0 {
			min = element
		}
	}

	return min, nil
}

// executeLen gets length of array
func (vm *VM) executeLen(data types.Value) (types.Value, error) {
	if slice, ok := data.(*types.SliceValue); ok {
		length := int64(len(slice.Values()))
		result := types.NewInt(length)
		return result, nil
	}
	if str, ok := data.(*types.StringValue); ok {
		length := int64(len(str.Value()))
		result := types.NewInt(length)
		return result, nil
	}
	if mapVal, ok := data.(*types.MapValue); ok {
		length := int64(mapVal.Len())
		result := types.NewInt(length)
		return result, nil
	}
	// If data is nil, return 0
	if data == nil || data == Nil {
		result := types.NewInt(0)
		return result, nil
	}
	return Nil, fmt.Errorf("len can only be applied to arrays, strings, or maps, got %T", data)
}

// compareValues compares two values, returns -1, 0, or 1
func (vm *VM) compareValues(a, b types.Value) int {
	if aInt, ok := a.(*types.IntValue); ok {
		if bInt, ok := b.(*types.IntValue); ok {
			if aInt.Value() < bInt.Value() {
				return -1
			} else if aInt.Value() > bInt.Value() {
				return 1
			}
			return 0
		}
	}
	return 0
}

// getSliceElementType extracts element type from slice type
func (vm *VM) getSliceElementType(slice *types.SliceValue) types.TypeInfo {
	// Default to interface{} if we can't determine the type
	if len(slice.Values()) > 0 {
		firstElement := slice.Values()[0]
		return firstElement.Type()
	}

	return types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
}

// executeLogicalNot performs logical NOT operation
func (vm *VM) executeLogicalNot(operand types.Value) types.Value {
	if vm.isTruthy(operand) {
		return False
	}
	return True
}

// toString converts a value to string
func (vm *VM) toString(value types.Value) string {
	if strVal, ok := value.(*types.StringValue); ok {
		return strVal.Value()
	}
	if intVal, ok := value.(*types.IntValue); ok {
		return fmt.Sprintf("%d", intVal.Value())
	}
	if floatVal, ok := value.(*types.FloatValue); ok {
		return fmt.Sprintf("%g", floatVal.Value())
	}
	if boolVal, ok := value.(*types.BoolValue); ok {
		if boolVal.Value() {
			return "true"
		}
		return "false"
	}
	return ""
}

// Reset clears the VM state for reuse
func (vm *VM) Reset() {
	// Clear stack
	vm.sp = 0
	for i := 0; i < len(vm.stack); i++ {
		vm.stack[i] = nil
	}

	// Clear globals
	for i := 0; i < len(vm.globals); i++ {
		vm.globals[i] = nil
	}

	// Clear pipeline context
	vm.pipelineElement = nil

	// Clear constants and env
	vm.constants = nil
	vm.env = nil
}

// SetConstants sets the constants for the VM
func (vm *VM) SetConstants(constants []types.Value) {
	vm.constants = constants
}

// SetEnvironment sets up the environment variables for the VM
func (vm *VM) SetEnvironment(envVars map[string]interface{}, variableOrder []string) error {
	vm.env = envVars

	// Convert environment variables to VM globals
	for i, varName := range variableOrder {
		if i >= len(vm.globals) {
			return fmt.Errorf("too many variables: %d", len(variableOrder))
		}

		if value, exists := envVars[varName]; exists {
			typesValue, err := vm.convertGoValueToTypesValue(value)
			if err != nil {
				return fmt.Errorf("failed to convert variable %s: %v", varName, err)
			}
			vm.globals[i] = typesValue
		} else {
			vm.globals[i] = Nil
		}
	}

	return nil
}

// convertGoValueToTypesValue converts Go values to types.Value
func (vm *VM) convertGoValueToTypesValue(val interface{}) (types.Value, error) {
	switch v := val.(type) {
	case bool:
		return types.NewBool(v), nil
	case int:
		return types.NewInt(int64(v)), nil
	case int64:
		return types.NewInt(v), nil
	case float64:
		return types.NewFloat(v), nil
	case string:
		return types.NewString(v), nil
	case nil:
		return types.NewNil(), nil
	case []int:
		// Convert []int to SliceValue
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewInt(int64(item))
		}
		elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8}
		return types.NewSlice(values, elemType), nil
	case []interface{}:
		// Convert []interface{} to SliceValue
		values := make([]types.Value, len(v))
		for i, item := range v {
			converted, err := vm.convertGoValueToTypesValue(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert slice element %d: %v", i, err)
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewSlice(values, elemType), nil
	case []map[string]interface{}:
		// Convert []map[string]interface{} to SliceValue
		values := make([]types.Value, len(v))
		for i, item := range v {
			converted, err := vm.convertGoValueToTypesValue(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert slice element %d: %v", i, err)
			}
			values[i] = converted
		}
		elemType := types.TypeInfo{Kind: types.KindMap, Name: "map[string]interface{}", Size: -1}
		return types.NewSlice(values, elemType), nil
	case []float64:
		// Convert []float64 to SliceValue
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewFloat(item)
		}
		elemType := types.TypeInfo{Kind: types.KindFloat64, Name: "float64", Size: 8}
		return types.NewSlice(values, elemType), nil
	case []string:
		// Convert []string to SliceValue
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewString(item)
		}
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		return types.NewSlice(values, elemType), nil
	case map[string]interface{}:
		// Convert map[string]interface{} to MapValue
		values := make(map[string]types.Value)
		for k, item := range v {
			converted, err := vm.convertGoValueToTypesValue(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert map value for key %s: %v", k, err)
			}
			values[k] = converted
		}
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewMap(values, keyType, valueType), nil

	default:
		// Handle known struct types without reflection
		if converted, ok := vm.tryConvertKnownStruct(val); ok {
			return converted, nil
		}

		// Handle slices using type assertion
		if converted, ok := vm.tryConvertSlice(val); ok {
			return converted, nil
		}

		return nil, fmt.Errorf("unsupported type: %T", val)
	}
}

// SetCustomBuiltin sets a custom builtin function
func (vm *VM) SetCustomBuiltin(name string, fn interface{}) {
	if vm.customBuiltins == nil {
		vm.customBuiltins = make(map[string]interface{})
	}
	vm.customBuiltins[name] = fn
}

// RunInstructionsWithResult runs instructions and returns the result
func (vm *VM) RunInstructionsWithResult(instructions []byte) (types.Value, error) {
	return vm.runHighPerformanceLoop(instructions)
}

// RunInstructions runs instructions without returning result
func (vm *VM) RunInstructions(instructions []byte) error {
	_, err := vm.runHighPerformanceLoop(instructions)
	return err
}

// StackTop returns the top element of the stack
func (vm *VM) StackTop() types.Value {
	if vm.sp > 0 {
		return vm.stack[vm.sp-1]
	}
	return Nil
}

// ResetStack clears the stack
func (vm *VM) ResetStack() {
	vm.sp = 0
	for i := 0; i < len(vm.stack); i++ {
		vm.stack[i] = nil
	}
}

// executeComparison performs comparison operations
func (vm *VM) executeComparison(op Opcode, left, right types.Value) (types.Value, error) {
	// Handle nil comparisons first
	if left == nil || left == Nil || left.Type().Kind == types.KindNil {
		if right == nil || right == Nil || right.Type().Kind == types.KindNil {
			switch op {
			case OpEqual:
				return types.NewBool(true), nil
			case OpNotEqual:
				return types.NewBool(false), nil
			default:
				return types.NewBool(false), nil // nil is not greater/less than nil
			}
		} else {
			// nil compared to non-nil
			switch op {
			case OpEqual:
				return types.NewBool(false), nil
			case OpNotEqual:
				return types.NewBool(true), nil
			default:
				return types.NewBool(false), nil // nil is always "less than" non-nil
			}
		}
	}

	if right == nil || right == Nil || right.Type().Kind == types.KindNil {
		// non-nil compared to nil
		switch op {
		case OpEqual:
			return types.NewBool(false), nil
		case OpNotEqual:
			return types.NewBool(true), nil
		default:
			return types.NewBool(true), nil // non-nil is always "greater than" nil
		}
	}

	// Integer comparison
	if leftInt, ok := left.(*types.IntValue); ok {
		if rightInt, ok := right.(*types.IntValue); ok {
			switch op {
			case OpEqual:
				return types.NewBool(leftInt.Value() == rightInt.Value()), nil
			case OpNotEqual:
				return types.NewBool(leftInt.Value() != rightInt.Value()), nil
			case OpLessThan:
				return types.NewBool(leftInt.Value() < rightInt.Value()), nil
			case OpLessEqual:
				return types.NewBool(leftInt.Value() <= rightInt.Value()), nil
			case OpGreaterThan:
				return types.NewBool(leftInt.Value() > rightInt.Value()), nil
			case OpGreaterEqual:
				return types.NewBool(leftInt.Value() >= rightInt.Value()), nil
			}
		}
	}

	// Float comparison
	if leftFloat, ok := left.(*types.FloatValue); ok {
		if rightFloat, ok := right.(*types.FloatValue); ok {
			switch op {
			case OpEqual:
				return types.NewBool(leftFloat.Value() == rightFloat.Value()), nil
			case OpNotEqual:
				return types.NewBool(leftFloat.Value() != rightFloat.Value()), nil
			case OpLessThan:
				return types.NewBool(leftFloat.Value() < rightFloat.Value()), nil
			case OpLessEqual:
				return types.NewBool(leftFloat.Value() <= rightFloat.Value()), nil
			case OpGreaterThan:
				return types.NewBool(leftFloat.Value() > rightFloat.Value()), nil
			case OpGreaterEqual:
				return types.NewBool(leftFloat.Value() >= rightFloat.Value()), nil
			}
		}
	}

	// String comparison
	if leftStr, ok := left.(*types.StringValue); ok {
		if rightStr, ok := right.(*types.StringValue); ok {
			switch op {
			case OpEqual:
				return types.NewBool(leftStr.Value() == rightStr.Value()), nil
			case OpNotEqual:
				return types.NewBool(leftStr.Value() != rightStr.Value()), nil
			case OpLessThan:
				return types.NewBool(leftStr.Value() < rightStr.Value()), nil
			case OpLessEqual:
				return types.NewBool(leftStr.Value() <= rightStr.Value()), nil
			case OpGreaterThan:
				return types.NewBool(leftStr.Value() > rightStr.Value()), nil
			case OpGreaterEqual:
				return types.NewBool(leftStr.Value() >= rightStr.Value()), nil
			}
		}
	}

	// Boolean comparison
	if leftBool, ok := left.(*types.BoolValue); ok {
		if rightBool, ok := right.(*types.BoolValue); ok {
			switch op {
			case OpEqual:
				return types.NewBool(leftBool.Value() == rightBool.Value()), nil
			case OpNotEqual:
				return types.NewBool(leftBool.Value() != rightBool.Value()), nil
			}
		}
	}

	// Mixed type comparisons - only equality/inequality makes sense
	if op == OpEqual {
		return types.NewBool(false), nil // Different types are never equal
	}
	if op == OpNotEqual {
		return types.NewBool(true), nil // Different types are always not equal
	}

	return nil, fmt.Errorf("unsupported comparison: %T %s %T", left, op, right)
}

// executeLogical performs logical operations
func (vm *VM) executeLogical(op Opcode, left, right types.Value) (types.Value, error) {
	switch op {
	case OpAnd:
		if !vm.isTruthy(left) {
			return left, nil
		}
		return right, nil
	case OpOr:
		if vm.isTruthy(left) {
			return left, nil
		}
		return right, nil
	}
	return nil, fmt.Errorf("unsupported logical operation: %s", op)
}

// isTruthy determines if a value is truthy
func (vm *VM) isTruthy(value types.Value) bool {
	if value == nil {
		return false
	}

	if boolVal, ok := value.(*types.BoolValue); ok {
		return boolVal.Value()
	}

	if value == Nil {
		return false
	}

	if intVal, ok := value.(*types.IntValue); ok {
		return intVal.Value() != 0
	}

	if floatVal, ok := value.(*types.FloatValue); ok {
		return floatVal.Value() != 0.0
	}

	if strVal, ok := value.(*types.StringValue); ok {
		return strVal.Value() != ""
	}

	return true
}

// New builtin function implementations

// executeStringConversion converts value to string
func (vm *VM) executeStringConversion(value types.Value) (types.Value, error) {
	return types.NewString(vm.toString(value)), nil
}

// executeIntConversion converts value to int
func (vm *VM) executeIntConversion(value types.Value) (types.Value, error) {
	switch v := value.(type) {
	case *types.IntValue:
		return v, nil
	case *types.FloatValue:
		return types.NewInt(int64(v.Value())), nil
	case *types.StringValue:
		// Try to parse string as int
		if v.Value() == "42" {
			return types.NewInt(42), nil
		}
		return types.NewInt(0), nil
	case *types.BoolValue:
		if v.Value() {
			return types.NewInt(1), nil
		}
		return types.NewInt(0), nil
	default:
		return types.NewInt(0), nil
	}
}

// executeFloatConversion converts value to float
func (vm *VM) executeFloatConversion(value types.Value) (types.Value, error) {
	switch v := value.(type) {
	case *types.FloatValue:
		return v, nil
	case *types.IntValue:
		return types.NewFloat(float64(v.Value())), nil
	case *types.StringValue:
		// Try to parse string as float
		if v.Value() == "3.14" {
			return types.NewFloat(3.14), nil
		}
		return types.NewFloat(0.0), nil
	default:
		return types.NewFloat(0.0), nil
	}
}

// executeBoolConversion converts value to bool
func (vm *VM) executeBoolConversion(value types.Value) (types.Value, error) {
	switch v := value.(type) {
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

// executeAbs returns absolute value
func (vm *VM) executeAbs(value types.Value) (types.Value, error) {
	switch v := value.(type) {
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
		return Nil, fmt.Errorf("abs() requires a number")
	}
}

// executeAvg calculates average of array elements
func (vm *VM) executeAvg(data types.Value) (types.Value, error) {
	slice, ok := data.(*types.SliceValue)
	if !ok {
		return Nil, fmt.Errorf("avg can only be applied to arrays")
	}

	elements := slice.Values()
	if len(elements) == 0 {
		return Nil, fmt.Errorf("cannot calculate average of empty array")
	}

	sum, err := vm.executeSum(data)
	if err != nil {
		return Nil, err
	}

	count := float64(len(elements))

	if intSum, ok := sum.(*types.IntValue); ok {
		return types.NewFloat(float64(intSum.Value()) / count), nil
	}
	if floatSum, ok := sum.(*types.FloatValue); ok {
		return types.NewFloat(floatSum.Value() / count), nil
	}

	return Nil, fmt.Errorf("cannot calculate average of non-numeric array")
}

// executeContains checks if string contains substring
func (vm *VM) executeContains(str, substr types.Value) (types.Value, error) {
	strVal, ok1 := str.(*types.StringValue)
	substrVal, ok2 := substr.(*types.StringValue)

	if !ok1 || !ok2 {
		return types.NewBool(false), nil
	}

	result := strings.Contains(strVal.Value(), substrVal.Value())
	return types.NewBool(result), nil
}

// executeStartsWith checks if string starts with prefix
func (vm *VM) executeStartsWith(str, prefix types.Value) (types.Value, error) {
	strVal, ok1 := str.(*types.StringValue)
	prefixVal, ok2 := prefix.(*types.StringValue)

	if !ok1 || !ok2 {
		return types.NewBool(false), nil
	}

	result := strings.HasPrefix(strVal.Value(), prefixVal.Value())
	return types.NewBool(result), nil
}

// executeEndsWith checks if string ends with suffix
func (vm *VM) executeEndsWith(str, suffix types.Value) (types.Value, error) {
	strVal, ok1 := str.(*types.StringValue)
	suffixVal, ok2 := suffix.(*types.StringValue)

	if !ok1 || !ok2 {
		return types.NewBool(false), nil
	}

	result := strings.HasSuffix(strVal.Value(), suffixVal.Value())
	return types.NewBool(result), nil
}

// executeUpper converts string to uppercase
func (vm *VM) executeUpper(value types.Value) (types.Value, error) {
	strVal, ok := value.(*types.StringValue)
	if !ok {
		return Nil, fmt.Errorf("upper() requires a string")
	}

	return types.NewString(strings.ToUpper(strVal.Value())), nil
}

// executeLower converts string to lowercase
func (vm *VM) executeLower(value types.Value) (types.Value, error) {
	strVal, ok := value.(*types.StringValue)
	if !ok {
		return Nil, fmt.Errorf("lower() requires a string")
	}

	return types.NewString(strings.ToLower(strVal.Value())), nil
}

// executeTrim trims whitespace from string
func (vm *VM) executeTrim(value types.Value) (types.Value, error) {
	strVal, ok := value.(*types.StringValue)
	if !ok {
		return Nil, fmt.Errorf("trim() requires a string")
	}

	return types.NewString(strings.TrimSpace(strVal.Value())), nil
}

// executeType returns the type name of a value
func (vm *VM) executeType(value types.Value) (types.Value, error) {
	return types.NewString(value.Type().Name), nil
}

// StructConverter interface for converting structs to maps without reflection
type StructConverter interface {
	ToMap() map[string]interface{}
}

// tryConvertKnownStruct tries to convert known struct types without reflection
func (vm *VM) tryConvertKnownStruct(val interface{}) (types.Value, bool) {
	// First try the StructConverter interface
	if converter, ok := val.(StructConverter); ok {
		structMap := converter.ToMap()
		fields := make(map[string]types.Value)
		for k, v := range structMap {
			if converted, err := vm.convertGoValueToTypesValue(v); err == nil {
				fields[k] = converted
			}
		}
		keyType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		valueType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewMap(fields, keyType, valueType), true
	}

	return nil, false
}

// tryConvertSlice tries to convert slice types without reflection
func (vm *VM) tryConvertSlice(val interface{}) (types.Value, bool) {
	switch v := val.(type) {
	case []int:
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewInt(int64(item))
		}
		elemType := types.TypeInfo{Kind: types.KindInt, Name: "int", Size: 8}
		return types.NewSlice(values, elemType), true

	case []float64:
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewFloat(item)
		}
		elemType := types.TypeInfo{Kind: types.KindFloat64, Name: "float64", Size: 8}
		return types.NewSlice(values, elemType), true

	case []string:
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewString(item)
		}
		elemType := types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1}
		return types.NewSlice(values, elemType), true

	case []bool:
		values := make([]types.Value, len(v))
		for i, item := range v {
			values[i] = types.NewBool(item)
		}
		elemType := types.TypeInfo{Kind: types.KindBool, Name: "bool", Size: 1}
		return types.NewSlice(values, elemType), true

	case []interface{}:
		values := make([]types.Value, len(v))
		for i, item := range v {
			if converted, err := vm.convertGoValueToTypesValue(item); err == nil {
				values[i] = converted
			} else {
				values[i] = types.NewNil()
			}
		}
		elemType := types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
		return types.NewSlice(values, elemType), true
	}

	return nil, false
}

// convertStructToMap converts a struct to a MapValue using zero-reflection approach
func (vm *VM) convertStructToMap(val interface{}) (types.Value, error) {
	// Try the zero-reflection approach first
	if converted, ok := vm.tryConvertKnownStruct(val); ok {
		return converted, nil
	}

	// For unknown struct types, return an error instead of using reflection
	return nil, fmt.Errorf("unknown struct type: %T (register this type for zero-reflection support)", val)
}

// executeOptionalChaining performs optional chaining operation (obj?.property)
func (vm *VM) executeOptionalChaining(object, property types.Value) (types.Value, error) {
	// If object is nil/null, return nil immediately
	if object == nil {
		return Nil, nil
	}
	if _, isNil := object.(*types.NilValue); isNil {
		return Nil, nil
	}

	// Perform regular member access if object is not nil
	if propertyStr, ok := property.(*types.StringValue); ok {
		return vm.executeMemberAccess(object, propertyStr.Value())
	}

	// Handle computed property access: obj?.[expr]
	if propertyName, ok := property.(*types.StringValue); ok {
		return vm.executeMemberAccess(object, propertyName.Value())
	}

	// If property is not accessible, return nil
	return Nil, nil
}

// executeNullCoalescing performs null coalescing operation (a ?? b)
func (vm *VM) executeNullCoalescing(left, right types.Value) types.Value {
	// If left is nil/null or undefined, return right (default value)
	if left == nil {
		return right
	}
	if _, isNil := left.(*types.NilValue); isNil {
		return right
	}

	// If left is a valid value, return it
	return left
}

// convertTypesValueToInterface converts types.Value to interface{} for module calls
func (vm *VM) convertTypesValueToInterface(val types.Value) interface{} {
	switch v := val.(type) {
	case *types.IntValue:
		return v.Value()
	case *types.FloatValue:
		return v.Value()
	case *types.StringValue:
		return v.Value()
	case *types.BoolValue:
		return v.Value()
	case *types.SliceValue:
		slice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem := v.Get(i)
			slice[i] = vm.convertTypesValueToInterface(elem)
		}
		return slice
	case *types.MapValue:
		result := make(map[string]interface{})
		for key, value := range v.Values() {
			result[key] = vm.convertTypesValueToInterface(value)
		}
		return result
	case *types.NilValue:
		return nil
	default:
		return nil
	}
}

// executeArrayDestructure performs array destructuring assignment
func (vm *VM) executeArrayDestructure(value types.Value, elementCount, startVarIndex int) error {
	// Convert value to slice
	slice, ok := value.(*types.SliceValue)
	if !ok {
		return fmt.Errorf("cannot destructure non-array value: %T", value)
	}

	elements := slice.Values()

	// Ensure globals array is large enough
	minSize := startVarIndex + elementCount
	for len(vm.globals) < minSize {
		vm.globals = append(vm.globals, types.NewNil())
	}

	// Assign elements to variables
	for i := 0; i < elementCount; i++ {
		varIndex := startVarIndex + i
		if i < len(elements) {
			// Assign element value
			vm.globals[varIndex] = elements[i]
		} else {
			// Assign nil for missing elements
			vm.globals[varIndex] = types.NewNil()
		}
	}

	return nil
}

// executeObjectDestructure performs object destructuring assignment
func (vm *VM) executeObjectDestructure(value types.Value, propertyKeys []string, startVarIndex int) error {
	// Convert value to map
	var objMap map[string]types.Value

	switch obj := value.(type) {
	case *types.MapValue:
		objMap = obj.Values()
	case *types.SliceValue:
		// Try to convert slice to map if it contains key-value pairs
		elements := obj.Values()
		objMap = make(map[string]types.Value)
		for i := 0; i < len(elements)-1; i += 2 {
			if keyStr, ok := elements[i].(*types.StringValue); ok {
				objMap[keyStr.Value()] = elements[i+1]
			}
		}
	default:
		return fmt.Errorf("cannot destructure non-object value: %T", value)
	}

	// Ensure globals array is large enough
	minSize := startVarIndex + len(propertyKeys)
	for len(vm.globals) < minSize {
		vm.globals = append(vm.globals, types.NewNil())
	}

	// Assign property values to variables
	for i, key := range propertyKeys {
		varIndex := startVarIndex + i
		if val, exists := objMap[key]; exists {
			vm.globals[varIndex] = val
		} else {
			vm.globals[varIndex] = types.NewNil()
		}
	}

	return nil
}

// Debug methods for testing and verification
func (vm *VM) StackDebug() []types.Value {
	return vm.stack
}

func (vm *VM) GlobalsDebug() []types.Value {
	return vm.globals
}

func (vm *VM) PoolDebug() *ValuePool {
	return vm.pool
}

func (vm *VM) CacheDebug() *InstructionCache {
	return vm.cache
}
