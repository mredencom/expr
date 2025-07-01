package vm

import "fmt"

// Opcode represents a virtual machine instruction
type Opcode byte

const (
	// Stack operations
	OpConstant Opcode = iota // Load constant onto stack
	OpPop                    // Pop value from stack
	OpDup                    // Duplicate top stack value
	OpSwap                   // Swap top two stack values

	// Arithmetic operations
	OpAdd // Add two values
	OpSub // Subtract two values
	OpMul // Multiply two values
	OpDiv // Divide two values
	OpMod // Modulo operation
	OpPow // Power operation
	OpNeg // Negate value

	// Fast path arithmetic operations (type-specialized)
	OpAddInt64   // Fast int64 addition
	OpSubInt64   // Fast int64 subtraction
	OpMulInt64   // Fast int64 multiplication
	OpDivInt64   // Fast int64 division
	OpModInt64   // Fast int64 modulo
	OpAddFloat64 // Fast float64 addition
	OpSubFloat64 // Fast float64 subtraction
	OpMulFloat64 // Fast float64 multiplication
	OpDivFloat64 // Fast float64 division
	OpModFloat64 // Fast float64 modulo
	OpAddString  // Fast string concatenation

	// Comparison operations
	OpEqual        // Equal comparison
	OpNotEqual     // Not equal comparison
	OpGreaterThan  // Greater than comparison
	OpGreaterEqual // Greater than or equal comparison
	OpLessThan     // Less than comparison
	OpLessEqual    // Less than or equal comparison

	// Logical operations
	OpAnd // Logical AND
	OpOr  // Logical OR
	OpNot // Logical NOT

	// Bitwise operations
	OpBitAnd // Bitwise AND
	OpBitOr  // Bitwise OR
	OpBitXor // Bitwise XOR
	OpBitNot // Bitwise NOT
	OpShiftL // Left shift
	OpShiftR // Right shift

	// Variable operations
	OpGetVar // Get variable value
	OpSetVar // Set variable value

	// Function operations
	OpCall    // Call function
	OpReturn  // Return from function
	OpBuiltin // Call builtin function

	// Collection operations
	OpIndex  // Index access (array[index], map[key])
	OpMember // Member access (obj.field)
	OpSlice  // Create slice literal
	OpMap    // Create map literal
	OpIn     // 'in' operator

	// Control flow
	OpJump      // Unconditional jump
	OpJumpTrue  // Jump if true
	OpJumpFalse // Jump if false
	OpJumpNil   // Jump if nil

	// String operations
	OpConcat     // String concatenation
	OpMatches    // String matches regex
	OpContains   // String contains substring
	OpStartsWith // String starts with prefix
	OpEndsWith   // String ends with suffix

	// Type conversion
	OpToString // Convert to string
	OpToInt    // Convert to int
	OpToFloat  // Convert to float
	OpToBool   // Convert to bool

	// Special operations
	OpHalt // Halt execution
	OpNoop // No operation

	// Collection and functional programming operations
	OpArray              // Create array from stack elements
	OpObject             // Create object from stack key-value pairs
	OpLambda             // Create lambda function
	OpClosure            // Create closure with captured variables
	OpApply              // Apply function to arguments
	OpPipe               // Pipeline operation (data | function)
	OpFilter             // Filter array with predicate
	OpMapFunc            // Map function over array (renamed from OpMap)
	OpReduce             // Reduce array with accumulator function
	OpGetPipelineElement // Get current pipeline element (for placeholder #)

	// Phase 3 optimization instructions
	OpAddVars      // Add two variables directly (var1 + var2)
	OpMulVars      // Multiply two variables directly (var1 * var2)
	OpSubVars      // Subtract two variables directly (var1 - var2)
	OpDivVars      // Divide two variables directly (var1 / var2)
	OpCompareVars  // Compare two variables directly (var1 op var2)
	OpConstantOp   // Constant operation result (pre-computed)
	OpFusedArith   // Fused arithmetic operation (a + b * c)
	OpCachedResult // Cached computation result

	// Null safety operations
	OpOptionalChaining // Optional chaining operation (obj?.property)
	OpNullCoalescing   // Null coalescing operation (a ?? b)

	// Module system operations
	OpModuleCall // Module function call (module.function(args...))

	// Destructuring operations
	OpArrayDestructure  // Array destructuring assignment
	OpObjectDestructure // Object destructuring assignment
	OpRestElement       // Rest element in destructuring
)

// String returns the string representation of an opcode
func (op Opcode) String() string {
	switch op {
	case OpConstant:
		return "OpConstant"
	case OpPop:
		return "OpPop"
	case OpDup:
		return "OpDup"
	case OpSwap:
		return "OpSwap"
	case OpAdd:
		return "OpAdd"
	case OpSub:
		return "OpSub"
	case OpMul:
		return "OpMul"
	case OpDiv:
		return "OpDiv"
	case OpMod:
		return "OpMod"
	case OpPow:
		return "OpPow"
	case OpNeg:
		return "OpNeg"
	case OpAddInt64:
		return "OpAddInt64"
	case OpSubInt64:
		return "OpSubInt64"
	case OpMulInt64:
		return "OpMulInt64"
	case OpDivInt64:
		return "OpDivInt64"
	case OpModInt64:
		return "OpModInt64"
	case OpAddFloat64:
		return "OpAddFloat64"
	case OpSubFloat64:
		return "OpSubFloat64"
	case OpMulFloat64:
		return "OpMulFloat64"
	case OpDivFloat64:
		return "OpDivFloat64"
	case OpModFloat64:
		return "OpModFloat64"
	case OpAddString:
		return "OpAddString"
	case OpEqual:
		return "OpEqual"
	case OpNotEqual:
		return "OpNotEqual"
	case OpGreaterThan:
		return "OpGreaterThan"
	case OpGreaterEqual:
		return "OpGreaterEqual"
	case OpLessThan:
		return "OpLessThan"
	case OpLessEqual:
		return "OpLessEqual"
	case OpAnd:
		return "OpAnd"
	case OpOr:
		return "OpOr"
	case OpNot:
		return "OpNot"
	case OpBitAnd:
		return "OpBitAnd"
	case OpBitOr:
		return "OpBitOr"
	case OpBitXor:
		return "OpBitXor"
	case OpBitNot:
		return "OpBitNot"
	case OpShiftL:
		return "OpShiftL"
	case OpShiftR:
		return "OpShiftR"
	case OpGetVar:
		return "OpGetVar"
	case OpSetVar:
		return "OpSetVar"
	case OpCall:
		return "OpCall"
	case OpReturn:
		return "OpReturn"
	case OpBuiltin:
		return "OpBuiltin"
	case OpIndex:
		return "OpIndex"
	case OpMember:
		return "OpMember"
	case OpSlice:
		return "OpSlice"
	case OpMap:
		return "OpMap"
	case OpIn:
		return "OpIn"
	case OpJump:
		return "OpJump"
	case OpJumpTrue:
		return "OpJumpTrue"
	case OpJumpFalse:
		return "OpJumpFalse"
	case OpJumpNil:
		return "OpJumpNil"
	case OpConcat:
		return "OpConcat"
	case OpMatches:
		return "OpMatches"
	case OpContains:
		return "OpContains"
	case OpStartsWith:
		return "OpStartsWith"
	case OpEndsWith:
		return "OpEndsWith"
	case OpToString:
		return "OpToString"
	case OpToInt:
		return "OpToInt"
	case OpToFloat:
		return "OpToFloat"
	case OpToBool:
		return "OpToBool"
	case OpHalt:
		return "OpHalt"
	case OpNoop:
		return "OpNoop"
	case OpArray:
		return "OpArray"
	case OpObject:
		return "OpObject"
	case OpLambda:
		return "OpLambda"
	case OpClosure:
		return "OpClosure"
	case OpApply:
		return "OpApply"
	case OpPipe:
		return "OpPipe"
	case OpFilter:
		return "OpFilter"
	case OpMapFunc:
		return "OpMapFunc"
	case OpReduce:
		return "OpReduce"
	case OpGetPipelineElement:
		return "OpGetPipelineElement"
	case OpOptionalChaining:
		return "OpOptionalChaining"
	case OpNullCoalescing:
		return "OpNullCoalescing"
	case OpModuleCall:
		return "OpModuleCall"
	case OpArrayDestructure:
		return "OpArrayDestructure"
	case OpObjectDestructure:
		return "OpObjectDestructure"
	case OpRestElement:
		return "OpRestElement"
	default:
		return fmt.Sprintf("Unknown(%d)", int(op))
	}
}

// Definition describes an instruction definition
type Definition struct {
	Name         string
	OperandWidth []int // Width of each operand in bytes
}

// definitions maps opcodes to their definitions
var definitions = map[Opcode]*Definition{
	OpConstant:           {"OpConstant", []int{2}}, // 2-byte constant index
	OpPop:                {"OpPop", []int{}},
	OpDup:                {"OpDup", []int{}},
	OpSwap:               {"OpSwap", []int{}},
	OpAdd:                {"OpAdd", []int{}},
	OpSub:                {"OpSub", []int{}},
	OpMul:                {"OpMul", []int{}},
	OpDiv:                {"OpDiv", []int{}},
	OpMod:                {"OpMod", []int{}},
	OpPow:                {"OpPow", []int{}},
	OpNeg:                {"OpNeg", []int{}},
	OpAddInt64:           {"OpAddInt64", []int{}},
	OpSubInt64:           {"OpSubInt64", []int{}},
	OpMulInt64:           {"OpMulInt64", []int{}},
	OpDivInt64:           {"OpDivInt64", []int{}},
	OpModInt64:           {"OpModInt64", []int{}},
	OpAddFloat64:         {"OpAddFloat64", []int{}},
	OpSubFloat64:         {"OpSubFloat64", []int{}},
	OpMulFloat64:         {"OpMulFloat64", []int{}},
	OpDivFloat64:         {"OpDivFloat64", []int{}},
	OpModFloat64:         {"OpModFloat64", []int{}},
	OpAddString:          {"OpAddString", []int{}},
	OpEqual:              {"OpEqual", []int{}},
	OpNotEqual:           {"OpNotEqual", []int{}},
	OpGreaterThan:        {"OpGreaterThan", []int{}},
	OpGreaterEqual:       {"OpGreaterEqual", []int{}},
	OpLessThan:           {"OpLessThan", []int{}},
	OpLessEqual:          {"OpLessEqual", []int{}},
	OpAnd:                {"OpAnd", []int{}},
	OpOr:                 {"OpOr", []int{}},
	OpNot:                {"OpNot", []int{}},
	OpBitAnd:             {"OpBitAnd", []int{}},
	OpBitOr:              {"OpBitOr", []int{}},
	OpBitXor:             {"OpBitXor", []int{}},
	OpBitNot:             {"OpBitNot", []int{}},
	OpShiftL:             {"OpShiftL", []int{}},
	OpShiftR:             {"OpShiftR", []int{}},
	OpGetVar:             {"OpGetVar", []int{2}}, // 2-byte variable index
	OpSetVar:             {"OpSetVar", []int{2}}, // 2-byte variable index
	OpCall:               {"OpCall", []int{1}},   // 1-byte argument count
	OpReturn:             {"OpReturn", []int{}},
	OpBuiltin:            {"OpBuiltin", []int{1, 1}}, // 1-byte builtin index, 1-byte arg count
	OpIndex:              {"OpIndex", []int{}},
	OpMember:             {"OpMember", []int{}}, // No operands, field name is on stack
	OpSlice:              {"OpSlice", []int{2}}, // 2-byte element count
	OpMap:                {"OpMap", []int{2}},   // 2-byte pair count
	OpIn:                 {"OpIn", []int{}},
	OpJump:               {"OpJump", []int{2}},      // 2-byte jump offset
	OpJumpTrue:           {"OpJumpTrue", []int{2}},  // 2-byte jump offset
	OpJumpFalse:          {"OpJumpFalse", []int{2}}, // 2-byte jump offset
	OpJumpNil:            {"OpJumpNil", []int{2}},   // 2-byte jump offset
	OpConcat:             {"OpConcat", []int{}},
	OpMatches:            {"OpMatches", []int{}},
	OpContains:           {"OpContains", []int{}},
	OpStartsWith:         {"OpStartsWith", []int{}},
	OpEndsWith:           {"OpEndsWith", []int{}},
	OpToString:           {"OpToString", []int{}},
	OpToInt:              {"OpToInt", []int{}},
	OpToFloat:            {"OpToFloat", []int{}},
	OpToBool:             {"OpToBool", []int{}},
	OpHalt:               {"OpHalt", []int{}},
	OpNoop:               {"OpNoop", []int{}},
	OpArray:              {"OpArray", []int{}},
	OpObject:             {"OpObject", []int{}},
	OpLambda:             {"OpLambda", []int{}},
	OpClosure:            {"OpClosure", []int{}},
	OpApply:              {"OpApply", []int{}},
	OpPipe:               {"OpPipe", []int{}},
	OpFilter:             {"OpFilter", []int{}},
	OpMapFunc:            {"OpMapFunc", []int{}},
	OpReduce:             {"OpReduce", []int{}},
	OpGetPipelineElement: {"OpGetPipelineElement", []int{}},
	OpOptionalChaining:   {"OpOptionalChaining", []int{}},
	OpNullCoalescing:     {"OpNullCoalescing", []int{}},
	OpModuleCall:         {"OpModuleCall", []int{2, 2, 1}},     // 2-byte module name index, 2-byte function name index, 1-byte arg count
	OpArrayDestructure:   {"OpArrayDestructure", []int{2, 2}},  // 2-byte element count, 2-byte start variable index
	OpObjectDestructure:  {"OpObjectDestructure", []int{2, 2}}, // 2-byte property count, 2-byte start variable index
	OpRestElement:        {"OpRestElement", []int{2}},          // 2-byte variable index for rest element
}

// Lookup returns the definition for an opcode
func Lookup(op Opcode) (*Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", int(op))
	}
	return def, nil
}

// Make creates an instruction from opcode and operands
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidth {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, operand := range operands {
		width := def.OperandWidth[i]
		switch width {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			instruction[offset] = byte(operand >> 8)
			instruction[offset+1] = byte(operand)
		}
		offset += width
	}

	return instruction
}

// ReadOperands reads operands from instruction bytes
func ReadOperands(def *Definition, ins []byte) ([]int, int) {
	operands := make([]int, len(def.OperandWidth))
	offset := 0

	for i, width := range def.OperandWidth {
		// Check if we have enough bytes
		if offset+width > len(ins) {
			// Return partial operands read so far
			return operands[:i], offset
		}

		switch width {
		case 1:
			operands[i] = int(ins[offset])
		case 2:
			operands[i] = int(ins[offset])<<8 | int(ins[offset+1])
		}
		offset += width
	}

	return operands, offset
}

// FormatInstruction formats an instruction for debugging
func FormatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidth)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand count %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	default:
		return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
	}
}
