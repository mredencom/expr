package compiler

import (
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

func TestNew(t *testing.T) {
	compiler := New()
	if compiler == nil {
		t.Fatal("Expected non-nil compiler")
	}
	if compiler.constants == nil {
		t.Fatal("Expected non-nil constants")
	}
	if compiler.symbolTable == nil {
		t.Fatal("Expected non-nil symbol table")
	}
	if len(compiler.scopes) != 1 {
		t.Errorf("Expected 1 scope, got %d", len(compiler.scopes))
	}
}

func TestNewWithState(t *testing.T) {
	symbolTable := NewSymbolTable()
	constants := []types.Value{types.NewInt(42)}

	compiler := NewWithState(symbolTable, constants)
	if compiler == nil {
		t.Fatal("Expected non-nil compiler")
	}
	if compiler.symbolTable != symbolTable {
		t.Error("Expected compiler to use provided symbol table")
	}
	if len(compiler.constants) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(compiler.constants))
	}
}

func TestCompileLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"42", []byte{byte(vm.OpConstant), 0, 0}},
		{"3.14", []byte{byte(vm.OpConstant), 0, 0}},
		{"true", []byte{byte(vm.OpConstant), 0, 0}},
		{"false", []byte{byte(vm.OpConstant), 0, 0}},
		{`"hello"`, []byte{byte(vm.OpConstant), 0, 0}},
		{"null", []byte{byte(vm.OpConstant), 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			compiler := New()

			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			bytecode := compiler.Bytecode()
			if len(bytecode.Instructions) < 3 {
				t.Errorf("Expected at least 3 instructions, got %d", len(bytecode.Instructions))
			}

			// Check that the opcode is correct
			if bytecode.Instructions[0] != byte(vm.OpConstant) {
				t.Errorf("Expected OpConstant, got %d", bytecode.Instructions[0])
			}

			// Check that we have one constant
			if len(bytecode.Constants) != 1 {
				t.Errorf("Expected 1 constant, got %d", len(bytecode.Constants))
			}
		})
	}
}

func TestCompileInfixExpressions(t *testing.T) {
	tests := []struct {
		input             string
		expectedConstants int
		expectedOps       []vm.Opcode
	}{
		// Constant folding optimization: arithmetic operations on literals are computed at compile time
		{"5 + 3", 1, []vm.Opcode{vm.OpConstant}}, // Optimized to constant 8
		{"5 - 3", 1, []vm.Opcode{vm.OpConstant}}, // Optimized to constant 2
		{"5 * 3", 1, []vm.Opcode{vm.OpConstant}}, // Optimized to constant 15
		{"5 / 3", 1, []vm.Opcode{vm.OpConstant}}, // Optimized to constant 1
		// Comparison operations cannot be constant folded as easily
		{"5 > 3", 2, []vm.Opcode{vm.OpConstant, vm.OpConstant, vm.OpGreaterThan}},
		{"5 < 3", 2, []vm.Opcode{vm.OpConstant, vm.OpConstant, vm.OpGreaterThan}}, // < becomes > with swapped operands
		{"5 == 3", 2, []vm.Opcode{vm.OpConstant, vm.OpConstant, vm.OpEqual}},
		{"5 != 3", 2, []vm.Opcode{vm.OpConstant, vm.OpConstant, vm.OpNotEqual}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			compiler := New()

			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			bytecode := compiler.Bytecode()
			if len(bytecode.Constants) != tt.expectedConstants {
				t.Errorf("Expected %d constants, got %d", tt.expectedConstants, len(bytecode.Constants))
			}

			// Check that we have the expected operations in the bytecode
			ops := extractOpcodes(bytecode.Instructions)
			if len(ops) < len(tt.expectedOps) {
				t.Errorf("Expected at least %d opcodes, got %d", len(tt.expectedOps), len(ops))
			}

			// Check the main operation is present
			lastOp := tt.expectedOps[len(tt.expectedOps)-1]
			found := false
			for _, op := range ops {
				if op == lastOp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find opcode %v in %v", lastOp, ops)
			}
		})
	}
}

func TestCompilePrefixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		expectedOp vm.Opcode
		constants  int
	}{
		{"-5", vm.OpNeg, 1},
		{"!true", vm.OpNot, 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			compiler := New()

			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			bytecode := compiler.Bytecode()
			if len(bytecode.Constants) != tt.constants {
				t.Errorf("Expected %d constants, got %d", tt.constants, len(bytecode.Constants))
			}

			ops := extractOpcodes(bytecode.Instructions)
			found := false
			for _, op := range ops {
				if op == tt.expectedOp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find opcode %v in %v", tt.expectedOp, ops)
			}
		})
	}
}

func TestCompileArrayLiteral(t *testing.T) {
	program := parseProgram(t, "[1, 2, 3]")
	compiler := New()

	err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := compiler.Bytecode()
	if len(bytecode.Constants) != 3 {
		t.Errorf("Expected 3 constants, got %d", len(bytecode.Constants))
	}

	ops := extractOpcodes(bytecode.Instructions)
	found := false
	for _, op := range ops {
		if op == vm.OpSlice {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find OpSlice in instructions")
	}
}

func TestCompileMapLiteral(t *testing.T) {
	program := parseProgram(t, `{"key": "value", "num": 42}`)
	compiler := New()

	err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := compiler.Bytecode()
	if len(bytecode.Constants) != 4 { // "key", "value", "num", 42
		t.Errorf("Expected 4 constants, got %d", len(bytecode.Constants))
	}

	ops := extractOpcodes(bytecode.Instructions)
	found := false
	for _, op := range ops {
		if op == vm.OpMap {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find OpMap in instructions")
	}
}

func TestCompileBuiltinExpression(t *testing.T) {
	tests := []string{
		`len("hello")`,
		`string(42)`,
		`int("42")`,
		`bool(1)`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			program := parseProgram(t, input)
			compiler := New()

			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("Compilation error: %v", err)
			}

			bytecode := compiler.Bytecode()
			ops := extractOpcodes(bytecode.Instructions)

			// Should contain a call to a builtin function
			found := false
			for _, op := range ops {
				if op == vm.OpBuiltin || op == vm.OpCall {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find OpBuiltin or OpCall in instructions for builtin, got: %v", ops)
			}
		})
	}
}

func TestCompileConditionalExpression(t *testing.T) {
	program := parseProgram(t, "5 > 3 ? 10 : 20")
	compiler := New()

	err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := compiler.Bytecode()
	if len(bytecode.Constants) < 3 { // 5, 3, 10, 20 (at least 3)
		t.Errorf("Expected at least 3 constants, got %d", len(bytecode.Constants))
	}

	ops := extractOpcodes(bytecode.Instructions)
	// Should contain jump instructions for conditional logic
	hasJump := false
	for _, op := range ops {
		if op == vm.OpJumpFalse || op == vm.OpJump {
			hasJump = true
			break
		}
	}
	if !hasJump {
		t.Error("Expected conditional to generate jump instructions")
	}
}

func TestBytecode(t *testing.T) {
	compiler := New()
	compiler.constants = []types.Value{types.NewInt(42)}

	bytecode := compiler.Bytecode()
	if bytecode == nil {
		t.Fatal("Expected non-nil bytecode")
	}
	if len(bytecode.Constants) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(bytecode.Constants))
	}
}

func TestAddConstant(t *testing.T) {
	compiler := New()
	val := types.NewInt(42)

	index := compiler.addConstant(val)
	if index != 0 {
		t.Errorf("Expected index 0, got %d", index)
	}

	if len(compiler.constants) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(compiler.constants))
	}

	if !compiler.constants[0].Equal(val) {
		t.Error("Expected constant to be added correctly")
	}
}

func TestErrors(t *testing.T) {
	compiler := New()
	compiler.errors = []string{"test error"}

	errors := compiler.Errors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
	if errors[0] != "test error" {
		t.Errorf("Expected 'test error', got %s", errors[0])
	}
}

func TestCompileUndefinedVariable(t *testing.T) {
	program := parseProgram(t, "undefined_var")
	compiler := New()

	err := compiler.Compile(program)
	if err == nil {
		t.Fatal("Expected compilation error for undefined variable")
	}
}

func TestGetVariableOrder(t *testing.T) {
	compiler := New()
	compiler.symbolTable.Define("x")
	compiler.symbolTable.Define("y")

	order := compiler.GetVariableOrder()
	if len(order) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(order))
	}
}

// Helper functions

func parseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	if len(program.Statements) == 0 {
		t.Fatal("No statements parsed")
	}

	return program
}

func extractOpcodes(instructions []byte) []vm.Opcode {
	var ops []vm.Opcode
	i := 0
	for i < len(instructions) {
		op := vm.Opcode(instructions[i])
		ops = append(ops, op)

		// Skip operands based on opcode
		switch op {
		case vm.OpConstant, vm.OpGetVar, vm.OpSetVar, vm.OpSlice, vm.OpMap, vm.OpJump, vm.OpJumpFalse:
			i += 3 // 1 byte opcode + 2 byte operand
		case vm.OpCall, vm.OpBuiltin:
			i += 2 // 1 byte opcode + 1 byte operand
		default:
			i += 1 // Just the opcode
		}
	}
	return ops
}
