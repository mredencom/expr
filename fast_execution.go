package expr

import (
	"fmt"

	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// FastExecution provides ultra-fast execution for simple expressions
type FastExecution struct {
	machine *vm.VM
	// Pre-allocated VM to avoid recreation overhead
}

// NewFastExecution creates a reusable fast execution engine
func NewFastExecution() *FastExecution {
	// Create VM with minimal bytecode
	emptyBytecode := &vm.Bytecode{
		Instructions: []byte{},
		Constants:    []types.Value{},
	}

	return &FastExecution{
		machine: vm.New(emptyBytecode),
	}
}

// FastRun executes a pre-compiled program with minimal overhead
func (fe *FastExecution) FastRun(program *Program, environment interface{}) (interface{}, error) {
	// Reset VM state but keep allocated memory
	fe.machine.ResetStack()

	// Set constants directly
	fe.machine.SetConstants(program.bytecode.Constants)

	// Set environment if needed
	if environment != nil {
		if envMap, ok := environment.(map[string]interface{}); ok {
			err := fe.machine.SetEnvironment(envMap, program.variableOrder)
			if err != nil {
				return nil, err
			}
		}
	}

	// Execute instructions directly
	err := fe.machine.RunInstructions(program.bytecode.Instructions)
	if err != nil {
		return nil, err
	}

	// Get result with minimal conversion
	resultValue := fe.machine.StackTop()
	if resultValue == nil {
		return nil, fmt.Errorf("no result value")
	}

	// Fast type conversion
	switch v := resultValue.(type) {
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

// FastRunSimple is optimized for simple constant expressions
func FastRunSimple(program *Program) (interface{}, error) {
	// For constant-folded expressions, just return the first constant
	if len(program.bytecode.Instructions) == 3 && len(program.bytecode.Constants) == 1 {
		// This is likely a constant-folded expression: OpConstant + index
		value := program.bytecode.Constants[0]
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

	// Fall back to normal execution
	fe := NewFastExecution()
	return fe.FastRun(program, nil)
}
