package vm

import (
	"fmt"

	"github.com/mredencom/expr/types"
)

// OptimizedVM is a high-performance VM that uses union types instead of interfaces
// This implements the P0 optimization from PERFORMANCE_SUMMARY.md
type OptimizedVM struct {
	// Stack using union types (no interface overhead)
	stack []types.OptimizedValue
	sp    int

	// Globals using union types
	globals []types.OptimizedValue

	// Instructions and constants remain the same
	instructions []byte
	constants    []types.Value // Converting these to OptimizedValue would be next step
	ip           int

	// Jump table for fast instruction dispatch
	jumpTable *OptimizedJumpTable
}

// OptimizedJumpTable for ultra-fast instruction dispatch
type OptimizedJumpTable struct {
	handlers [256]func(*OptimizedVM) error
}

// NewOptimizedVM creates a new optimized VM instance
func NewOptimizedVM(instructions []byte, constants []types.Value) *OptimizedVM {
	vm := &OptimizedVM{
		stack:        make([]types.OptimizedValue, StackSize),
		sp:           0,
		globals:      make([]types.OptimizedValue, GlobalsSize),
		instructions: instructions,
		constants:    constants,
		ip:           0,
		jumpTable:    NewOptimizedJumpTable(),
	}
	return vm
}

// NewOptimizedJumpTable creates the optimized jump table
func NewOptimizedJumpTable() *OptimizedJumpTable {
	jt := &OptimizedJumpTable{}

	// Initialize handlers for maximum performance
	jt.handlers[OpConstant] = (*OptimizedVM).handleConstant
	jt.handlers[OpAdd] = (*OptimizedVM).handleAdd
	jt.handlers[OpSub] = (*OptimizedVM).handleSub
	jt.handlers[OpMul] = (*OptimizedVM).handleMul
	jt.handlers[OpDiv] = (*OptimizedVM).handleDiv
	jt.handlers[OpEqual] = (*OptimizedVM).handleEqual
	jt.handlers[OpNotEqual] = (*OptimizedVM).handleNotEqual
	jt.handlers[OpGreaterThan] = (*OptimizedVM).handleGreaterThan
	jt.handlers[OpGreaterEqual] = (*OptimizedVM).handleGreaterEqual
	jt.handlers[OpLessThan] = (*OptimizedVM).handleLessThan
	jt.handlers[OpLessEqual] = (*OptimizedVM).handleLessEqual
	jt.handlers[OpAnd] = (*OptimizedVM).handleAnd
	jt.handlers[OpOr] = (*OptimizedVM).handleOr
	jt.handlers[OpNot] = (*OptimizedVM).handleNot
	jt.handlers[OpPop] = (*OptimizedVM).handlePop
	jt.handlers[OpHalt] = (*OptimizedVM).handleHalt

	return jt
}

// Run executes the optimized VM with maximum performance
func (vm *OptimizedVM) Run() (types.OptimizedValue, error) {
	for vm.ip < len(vm.instructions) {
		opcode := vm.instructions[vm.ip]
		vm.ip++

		handler := vm.jumpTable.handlers[opcode]
		if handler == nil {
			return types.NewOptimizedNil(), fmt.Errorf("unknown opcode: %d", opcode)
		}

		if err := handler(vm); err != nil {
			return types.NewOptimizedNil(), err
		}
	}

	if vm.sp > 0 {
		return vm.stack[vm.sp-1], nil
	}

	return types.NewOptimizedNil(), nil
}

// High-performance instruction handlers using union types

// handleConstant loads a constant (optimized)
func (vm *OptimizedVM) handleConstant() error {
	if vm.ip+1 >= len(vm.instructions) {
		return fmt.Errorf("insufficient bytes for constant")
	}

	constIndex := int(vm.instructions[vm.ip])<<8 | int(vm.instructions[vm.ip+1])
	vm.ip += 2

	if constIndex >= len(vm.constants) {
		return fmt.Errorf("constant index out of bounds")
	}

	// Convert interface Value to OptimizedValue
	constant := vm.constants[constIndex]
	optimizedVal := vm.convertToOptimized(constant)

	vm.stack[vm.sp] = optimizedVal
	vm.sp++
	return nil
}

// handleAdd performs optimized addition
func (vm *OptimizedVM) handleAdd() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.AddOptimized(left, right)
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleSub performs optimized subtraction
func (vm *OptimizedVM) handleSub() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	// Optimized subtraction using union types
	result, err := vm.subtractOptimized(left, right)
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleMul performs optimized multiplication
func (vm *OptimizedVM) handleMul() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := vm.multiplyOptimized(left, right)
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleDiv performs optimized division
func (vm *OptimizedVM) handleDiv() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := vm.divideOptimized(left, right)
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleEqual performs optimized equality comparison
func (vm *OptimizedVM) handleEqual() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, "==")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleNotEqual performs optimized inequality comparison
func (vm *OptimizedVM) handleNotEqual() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, "!=")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleGreaterThan performs optimized greater than comparison
func (vm *OptimizedVM) handleGreaterThan() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, ">")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleGreaterEqual performs optimized greater than or equal comparison
func (vm *OptimizedVM) handleGreaterEqual() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, ">=")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleLessThan performs optimized less than comparison
func (vm *OptimizedVM) handleLessThan() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, "<")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleLessEqual performs optimized less than or equal comparison
func (vm *OptimizedVM) handleLessEqual() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result, err := types.CompareOptimized(left, right, "<=")
	if err != nil {
		return err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleAnd performs optimized logical AND
func (vm *OptimizedVM) handleAnd() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result := types.NewOptimizedBool(left.ToBool() && right.ToBool())

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleOr performs optimized logical OR
func (vm *OptimizedVM) handleOr() error {
	if vm.sp < 2 {
		return fmt.Errorf("insufficient operands")
	}

	right := &vm.stack[vm.sp-1]
	left := &vm.stack[vm.sp-2]

	result := types.NewOptimizedBool(left.ToBool() || right.ToBool())

	vm.stack[vm.sp-2] = result
	vm.sp--
	return nil
}

// handleNot performs optimized logical NOT
func (vm *OptimizedVM) handleNot() error {
	if vm.sp < 1 {
		return fmt.Errorf("insufficient operands")
	}

	operand := &vm.stack[vm.sp-1]
	result := types.NewOptimizedBool(!operand.ToBool())

	vm.stack[vm.sp-1] = result
	return nil
}

// handlePop removes the top stack element
func (vm *OptimizedVM) handlePop() error {
	if vm.sp <= 0 {
		return fmt.Errorf("stack underflow")
	}
	vm.sp--
	return nil
}

// handleHalt stops execution
func (vm *OptimizedVM) handleHalt() error {
	return nil
}

// Helper methods for optimized arithmetic operations

// subtractOptimized performs optimized subtraction
func (vm *OptimizedVM) subtractOptimized(left, right *types.OptimizedValue) (types.OptimizedValue, error) {
	// Integer - Integer (most common case)
	if left.Type == types.TypeInt64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedInt(left.IntVal - right.IntVal), nil
	}

	// Float - Float
	if left.Type == types.TypeFloat64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(left.FloatVal - right.FloatVal), nil
	}

	// Int - Float
	if left.Type == types.TypeInt64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(float64(left.IntVal) - right.FloatVal), nil
	}

	// Float - Int
	if left.Type == types.TypeFloat64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedFloat(left.FloatVal - float64(right.IntVal)), nil
	}

	return types.NewOptimizedNil(), fmt.Errorf("unsupported subtraction: %v - %v", left.Type, right.Type)
}

// multiplyOptimized performs optimized multiplication
func (vm *OptimizedVM) multiplyOptimized(left, right *types.OptimizedValue) (types.OptimizedValue, error) {
	// Integer * Integer (most common case)
	if left.Type == types.TypeInt64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedInt(left.IntVal * right.IntVal), nil
	}

	// Float * Float
	if left.Type == types.TypeFloat64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(left.FloatVal * right.FloatVal), nil
	}

	// Int * Float
	if left.Type == types.TypeInt64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(float64(left.IntVal) * right.FloatVal), nil
	}

	// Float * Int
	if left.Type == types.TypeFloat64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedFloat(left.FloatVal * float64(right.IntVal)), nil
	}

	return types.NewOptimizedNil(), fmt.Errorf("unsupported multiplication: %v * %v", left.Type, right.Type)
}

// divideOptimized performs optimized division
func (vm *OptimizedVM) divideOptimized(left, right *types.OptimizedValue) (types.OptimizedValue, error) {
	// Check for division by zero
	if (right.Type == types.TypeInt64 && right.IntVal == 0) ||
		(right.Type == types.TypeFloat64 && right.FloatVal == 0.0) {
		return types.NewOptimizedNil(), fmt.Errorf("division by zero")
	}

	// Integer / Integer (return float for precision)
	if left.Type == types.TypeInt64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedFloat(float64(left.IntVal) / float64(right.IntVal)), nil
	}

	// Float / Float
	if left.Type == types.TypeFloat64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(left.FloatVal / right.FloatVal), nil
	}

	// Int / Float
	if left.Type == types.TypeInt64 && right.Type == types.TypeFloat64 {
		return types.NewOptimizedFloat(float64(left.IntVal) / right.FloatVal), nil
	}

	// Float / Int
	if left.Type == types.TypeFloat64 && right.Type == types.TypeInt64 {
		return types.NewOptimizedFloat(left.FloatVal / float64(right.IntVal)), nil
	}

	return types.NewOptimizedNil(), fmt.Errorf("unsupported division: %v / %v", left.Type, right.Type)
}

// convertToOptimized converts interface Value to OptimizedValue
func (vm *OptimizedVM) convertToOptimized(val types.Value) types.OptimizedValue {
	switch v := val.(type) {
	case *types.BoolValue:
		return types.NewOptimizedBool(v.Value())
	case *types.IntValue:
		return types.NewOptimizedInt(v.Value())
	case *types.FloatValue:
		return types.NewOptimizedFloat(v.Value())
	case *types.StringValue:
		return types.NewOptimizedString(v.Value())
	case *types.NilValue:
		return types.NewOptimizedNil()
	default:
		// For complex types, we'd need more sophisticated conversion
		return types.NewOptimizedNil()
	}
}

// GetLastValue returns the last value on the stack for testing
func (vm *OptimizedVM) GetLastValue() types.OptimizedValue {
	if vm.sp > 0 {
		return vm.stack[vm.sp-1]
	}
	return types.NewOptimizedNil()
}
