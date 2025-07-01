package vm

import (
	"fmt"

	"github.com/mredencom/expr/types"
)

// InstructionHandler represents a function that handles a specific opcode
type InstructionHandler func(vm *VM, instructions []byte, ip *int) (bool, error)

// SafeJumpTable 是一个简化且安全的跳转表实现
// 专注于正确性和稳定性，而非极端性能优化
type SafeJumpTable struct {
	handlers map[Opcode]InstructionHandler
}

// NewSafeJumpTable 创建一个新的安全跳转表
func NewSafeJumpTable() *SafeJumpTable {
	jt := &SafeJumpTable{
		handlers: make(map[Opcode]InstructionHandler),
	}

	jt.initializeSafeHandlers()
	return jt
}

// initializeSafeHandlers 初始化所有安全的指令处理函数
func (jt *SafeJumpTable) initializeSafeHandlers() {
	// 基础栈操作
	jt.handlers[OpConstant] = safeHandleConstant
	jt.handlers[OpPop] = safeHandlePop
	jt.handlers[OpDup] = safeHandleDup
	jt.handlers[OpSwap] = safeHandleSwap

	// 算术运算
	jt.handlers[OpAdd] = safeHandleAdd
	jt.handlers[OpSub] = safeHandleSub
	jt.handlers[OpMul] = safeHandleMul
	jt.handlers[OpDiv] = safeHandleDiv
	jt.handlers[OpMod] = safeHandleMod
	jt.handlers[OpNeg] = safeHandleNeg

	// 类型特化算术运算
	jt.handlers[OpAddInt64] = safeHandleAdd
	jt.handlers[OpSubInt64] = safeHandleSub
	jt.handlers[OpMulInt64] = safeHandleMul
	jt.handlers[OpDivInt64] = safeHandleDiv
	jt.handlers[OpModInt64] = safeHandleMod
	jt.handlers[OpAddFloat64] = safeHandleAdd
	jt.handlers[OpSubFloat64] = safeHandleSub
	jt.handlers[OpMulFloat64] = safeHandleMul
	jt.handlers[OpDivFloat64] = safeHandleDiv
	jt.handlers[OpModFloat64] = safeHandleMod
	jt.handlers[OpAddString] = safeHandleAdd

	// 比较运算
	jt.handlers[OpEqual] = safeHandleEqual
	jt.handlers[OpNotEqual] = safeHandleNotEqual
	jt.handlers[OpGreaterThan] = safeHandleGreaterThan
	jt.handlers[OpGreaterEqual] = safeHandleGreaterEqual
	jt.handlers[OpLessThan] = safeHandleLessThan
	jt.handlers[OpLessEqual] = safeHandleLessEqual

	// 逻辑运算
	jt.handlers[OpAnd] = safeHandleAnd
	jt.handlers[OpOr] = safeHandleOr
	jt.handlers[OpNot] = safeHandleNot

	// 变量操作
	jt.handlers[OpGetVar] = safeHandleGetVar
	jt.handlers[OpSetVar] = safeHandleSetVar

	// 函数和内置函数
	jt.handlers[OpCall] = safeHandleCall
	jt.handlers[OpBuiltin] = safeHandleBuiltin

	// 集合操作
	jt.handlers[OpIndex] = safeHandleIndex
	jt.handlers[OpMember] = safeHandleMember
	jt.handlers[OpArray] = safeHandleArray
	jt.handlers[OpSlice] = safeHandleArray // 使用相同的处理函数
	jt.handlers[OpObject] = safeHandleObject
	jt.handlers[OpMap] = safeHandleObject // 使用相同的处理函数

	// 控制流
	jt.handlers[OpJump] = safeHandleJump
	jt.handlers[OpJumpTrue] = safeHandleJumpTrue
	jt.handlers[OpJumpFalse] = safeHandleJumpFalse
	jt.handlers[OpJumpNil] = safeHandleJumpNil

	// 管道操作
	jt.handlers[OpPipe] = safeHandlePipe
	jt.handlers[OpFilter] = safeHandleFilter
	jt.handlers[OpMapFunc] = safeHandleMapFunc
	jt.handlers[OpReduce] = safeHandleReduce

	// 空值安全操作
	jt.handlers[OpOptionalChaining] = safeHandleOptionalChaining
	jt.handlers[OpNullCoalescing] = safeHandleNullCoalescing

	// 注意：高级操作码暂时注释掉，等待VM中相应方法的实现
	// 位运算操作 (规划中)
	// jt.handlers[OpBitAnd] = safeHandleBitAnd
	// jt.handlers[OpBitOr] = safeHandleBitOr
	// jt.handlers[OpBitXor] = safeHandleBitXor
	// jt.handlers[OpBitNot] = safeHandleBitNot
	// jt.handlers[OpShiftL] = safeHandleShiftL
	// jt.handlers[OpShiftR] = safeHandleShiftR

	// 字符串操作 (规划中)
	// jt.handlers[OpConcat] = safeHandleConcat
	// jt.handlers[OpMatches] = safeHandleMatches
	// jt.handlers[OpContains] = safeHandleContains
	// jt.handlers[OpStartsWith] = safeHandleStartsWith
	// jt.handlers[OpEndsWith] = safeHandleEndsWith

	// 类型转换 (规划中)
	// jt.handlers[OpToString] = safeHandleToString
	// jt.handlers[OpToInt] = safeHandleToInt
	// jt.handlers[OpToFloat] = safeHandleToFloat
	// jt.handlers[OpToBool] = safeHandleToBool

	// 高级算术 (规划中)
	// jt.handlers[OpPow] = safeHandlePow

	// 特殊操作
	jt.handlers[OpHalt] = safeHandleHalt
	jt.handlers[OpNoop] = safeHandleNoop
}

// Execute 使用安全的跳转表派发指令
func (jt *SafeJumpTable) Execute(vm *VM, instructions []byte, ip *int) (bool, error) {
	if ip == nil || *ip >= len(instructions) {
		return false, nil
	}

	opcode := Opcode(instructions[*ip])
	*ip++

	handler, exists := jt.handlers[opcode]
	if !exists {
		return jt.fallbackToSwitch(vm, opcode, instructions, ip)
	}

	return handler(vm, instructions, ip)
}

// fallbackToSwitch 回退到原始的switch处理
func (jt *SafeJumpTable) fallbackToSwitch(vm *VM, opcode Opcode, instructions []byte, ip *int) (bool, error) {
	*ip--
	switch opcode {
	case OpReturn:
		return false, nil
	default:
		return false, fmt.Errorf("unsupported opcode: %v", opcode)
	}
}

// 安全的指令处理函数实现
// 这些函数专注于正确性，使用VM现有的方法确保兼容性

func safeHandleConstant(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for constant index")
	}

	constIndex := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	if constIndex >= len(vm.constants) {
		return false, fmt.Errorf("constant index out of bounds")
	}

	if vm.sp >= len(vm.stack) {
		return false, fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = vm.constants[constIndex]
	vm.sp++
	return true, nil
}

func safeHandleGetVar(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for variable index")
	}

	varIndex := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	if varIndex >= len(vm.globals) {
		return false, fmt.Errorf("variable index out of bounds")
	}

	if vm.sp >= len(vm.stack) {
		return false, fmt.Errorf("stack overflow")
	}

	val := vm.globals[varIndex]
	if val == nil {
		vm.stack[vm.sp] = Nil
	} else {
		vm.stack[vm.sp] = val
	}
	vm.sp++
	return true, nil
}

func safeHandlePop(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp <= 0 {
		return false, fmt.Errorf("stack underflow")
	}
	vm.sp--
	return true, nil
}

func safeHandleDup(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp <= 0 {
		return false, fmt.Errorf("stack underflow")
	}
	if vm.sp >= len(vm.stack) {
		return false, fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = vm.stack[vm.sp-1]
	vm.sp++
	return true, nil
}

func safeHandleSwap(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient stack elements")
	}

	vm.stack[vm.sp-1], vm.stack[vm.sp-2] = vm.stack[vm.sp-2], vm.stack[vm.sp-1]
	return true, nil
}

// 算术运算处理函数 - 使用VM现有的方法确保兼容性
func safeHandleAdd(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeAddition(left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleSub(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeSubtraction(left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleMul(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeMultiplication(left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleDiv(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeDivision(left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleMod(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeModulo(left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleNeg(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}

	operand := vm.stack[vm.sp-1]
	result, err := vm.executeNegation(operand)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-1] = result
	return true, nil
}

// 比较运算处理函数
func safeHandleEqual(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpEqual, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleNotEqual(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpNotEqual, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleGreaterThan(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpGreaterThan, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleGreaterEqual(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpGreaterEqual, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleLessThan(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpLessThan, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleLessEqual(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeComparison(OpLessEqual, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

// 逻辑运算处理函数
func safeHandleAnd(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeLogical(OpAnd, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleOr(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeLogical(OpOr, left, right)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleNot(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}

	operand := vm.stack[vm.sp-1]
	result := vm.executeLogicalNot(operand)
	vm.stack[vm.sp-1] = result
	return true, nil
}

// 其余处理函数的简化实现
func safeHandleSetVar(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) || vm.sp < 1 {
		return false, fmt.Errorf("invalid setvar instruction")
	}

	varIndex := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	if varIndex >= len(vm.globals) {
		return false, fmt.Errorf("variable index out of bounds")
	}

	vm.globals[varIndex] = vm.stack[vm.sp-1]
	vm.sp--
	return true, nil
}

func safeHandleCall(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for call")
	}

	argCount := int(instructions[*ip])
	*ip++

	return true, vm.executeCall(argCount)
}

func safeHandleBuiltin(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("incomplete OpBuiltin instruction")
	}

	builtinIndex := int(instructions[*ip])
	argCount := int(instructions[*ip+1])
	*ip += 2

	err := vm.executeBuiltin(builtinIndex, argCount)
	if err != nil {
		return false, err
	}

	return true, nil
}

func safeHandleIndex(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	index := vm.stack[vm.sp-1]
	object := vm.stack[vm.sp-2]
	result, err := vm.executeIndex(object, index)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleMember(vm *VM, instructions []byte, ip *int) (bool, error) {
	// OpMember expects: [object, memberName] on stack
	// Pops both and pushes result
	if vm.sp < 2 {
		return false, fmt.Errorf("stack underflow for member access")
	}

	// Pop memberName from stack
	vm.sp--
	memberName := vm.stack[vm.sp]

	// Pop object from stack
	vm.sp--
	object := vm.stack[vm.sp]

	// Execute member access
	result, err := vm.executeMemberByName(object, memberName)
	if err != nil {
		return false, err
	}

	// Push result to stack
	vm.stack[vm.sp] = result
	vm.sp++

	return true, nil
}

func safeHandleArray(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for array")
	}

	elementCount := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	result, err := vm.executeArray(elementCount)
	if err != nil {
		return false, err
	}

	if vm.sp >= len(vm.stack) {
		return false, fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = result
	vm.sp++
	return true, nil
}

func safeHandleObject(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for object")
	}

	pairCount := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	result, err := vm.executeObject(pairCount)
	if err != nil {
		return false, err
	}

	if vm.sp >= len(vm.stack) {
		return false, fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = result
	vm.sp++
	return true, nil
}

// 控制流处理函数
func safeHandleJump(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) {
		return false, fmt.Errorf("insufficient bytes for jump")
	}

	offset := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip = offset
	return true, nil
}

func safeHandleJumpTrue(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) || vm.sp < 1 {
		return false, fmt.Errorf("invalid jump true instruction")
	}

	offset := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	condition := vm.stack[vm.sp-1]
	vm.sp--

	if vm.isTruthy(condition) {
		*ip = offset
	}
	return true, nil
}

func safeHandleJumpFalse(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) || vm.sp < 1 {
		return false, fmt.Errorf("invalid jump false instruction")
	}

	offset := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	condition := vm.stack[vm.sp-1]
	vm.sp--

	if !vm.isTruthy(condition) {
		*ip = offset
	}
	return true, nil
}

func safeHandleJumpNil(vm *VM, instructions []byte, ip *int) (bool, error) {
	if *ip+1 >= len(instructions) || vm.sp < 1 {
		return false, fmt.Errorf("invalid jump nil instruction")
	}

	offset := int(instructions[*ip])<<8 | int(instructions[*ip+1])
	*ip += 2

	value := vm.stack[vm.sp-1]
	vm.sp--

	if value == nil || value.Type().Kind == types.KindNil {
		*ip = offset
	}
	return true, nil
}

// 管道操作处理函数
func safeHandlePipe(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	function := vm.stack[vm.sp-1]
	data := vm.stack[vm.sp-2]
	result, err := vm.executePipe(data, function)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleFilter(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	condition := vm.stack[vm.sp-1]
	data := vm.stack[vm.sp-2]
	result, err := vm.executeFilter(data, condition)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleMapFunc(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	transform := vm.stack[vm.sp-1]
	data := vm.stack[vm.sp-2]
	result, err := vm.executeMap(data, transform)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

// 空值安全操作处理函数
func safeHandleOptionalChaining(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	property := vm.stack[vm.sp-1]
	object := vm.stack[vm.sp-2]
	result, err := vm.executeOptionalChaining(object, property)
	if err != nil {
		return false, err
	}

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleNullCoalescing(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}

	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result := vm.executeNullCoalescing(left, right)

	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

// 特殊操作处理函数
func safeHandleReduce(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 3 {
		return false, fmt.Errorf("insufficient operands for reduce operation")
	}

	// 基本的reduce操作，调用VM的现有方法
	// 这里简化处理，实际的reduce逻辑在VM中实现
	return false, fmt.Errorf("reduce operation not yet fully implemented in safe jump table")
}

func safeHandleHalt(vm *VM, instructions []byte, ip *int) (bool, error) {
	return false, nil
}

func safeHandleNoop(vm *VM, instructions []byte, ip *int) (bool, error) {
	return true, nil
}

// 新增的高级操作处理函数
// 注意：这些函数暂时注释掉，等待VM中相应方法的实现

/*
// 位运算操作处理函数
func safeHandleBitAnd(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeBitwiseOperation("&", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleBitOr(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeBitwiseOperation("|", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleBitXor(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeBitwiseOperation("^", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleBitNot(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}
	operand := vm.stack[vm.sp-1]
	result, err := vm.executeBitwiseOperation("~", operand, nil)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-1] = result
	return true, nil
}

func safeHandleShiftL(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeBitwiseOperation("<<", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleShiftR(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeBitwiseOperation(">>", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

// 字符串操作处理函数
func safeHandleConcat(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]
	result, err := vm.executeStringOperation("concat", left, right)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleMatches(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	pattern := vm.stack[vm.sp-1]
	str := vm.stack[vm.sp-2]
	result, err := vm.executeStringOperation("matches", str, pattern)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleContains(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	substr := vm.stack[vm.sp-1]
	str := vm.stack[vm.sp-2]
	result, err := vm.executeStringOperation("contains", str, substr)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleStartsWith(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	prefix := vm.stack[vm.sp-1]
	str := vm.stack[vm.sp-2]
	result, err := vm.executeStringOperation("startsWith", str, prefix)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

func safeHandleEndsWith(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	suffix := vm.stack[vm.sp-1]
	str := vm.stack[vm.sp-2]
	result, err := vm.executeStringOperation("endsWith", str, suffix)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}

// 类型转换处理函数
func safeHandleToString(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}
	operand := vm.stack[vm.sp-1]
	result, err := vm.executeTypeConversion("string", operand)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-1] = result
	return true, nil
}

func safeHandleToInt(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}
	operand := vm.stack[vm.sp-1]
	result, err := vm.executeTypeConversion("int", operand)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-1] = result
	return true, nil
}

func safeHandleToFloat(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}
	operand := vm.stack[vm.sp-1]
	result, err := vm.executeTypeConversion("float", operand)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-1] = result
	return true, nil
}

func safeHandleToBool(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 1 {
		return false, fmt.Errorf("insufficient operands")
	}
	operand := vm.stack[vm.sp-1]
	result, err := vm.executeTypeConversion("bool", operand)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-1] = result
	return true, nil
}

// 幂运算处理函数
func safeHandlePow(vm *VM, instructions []byte, ip *int) (bool, error) {
	if vm.sp < 2 {
		return false, fmt.Errorf("insufficient operands")
	}
	exponent := vm.stack[vm.sp-1]
	base := vm.stack[vm.sp-2]
	result, err := vm.executePowerOperation(base, exponent)
	if err != nil {
		return false, err
	}
	vm.stack[vm.sp-2] = result
	vm.sp--
	return true, nil
}
*/
