package vm

import (
	"testing"
)

// TestLookup tests the Lookup function for opcode definitions
func TestLookup(t *testing.T) {
	t.Run("ValidOpcodes", func(t *testing.T) {
		testCases := []struct {
			opcode   Opcode
			expected string
		}{
			{OpConstant, "OpConstant"},
			{OpPop, "OpPop"},
			{OpAdd, "OpAdd"},
			{OpSub, "OpSub"},
			{OpMul, "OpMul"},
			{OpDiv, "OpDiv"},
			{OpEqual, "OpEqual"},
			{OpNotEqual, "OpNotEqual"},
			{OpGreaterThan, "OpGreaterThan"},
			{OpGetVar, "OpGetVar"},
			{OpSetVar, "OpSetVar"},
			{OpCall, "OpCall"},
			{OpBuiltin, "OpBuiltin"},
		}

		for _, test := range testCases {
			def, err := Lookup(test.opcode)
			if err != nil {
				t.Errorf("Expected no error for opcode %v, got %v", test.opcode, err)
				continue
			}
			if def == nil {
				t.Errorf("Expected definition for opcode %v, got nil", test.opcode)
				continue
			}
			if def.Name != test.expected {
				t.Errorf("Expected name %s for opcode %v, got %s", test.expected, test.opcode, def.Name)
			}
		}
	})

	t.Run("InvalidOpcode", func(t *testing.T) {
		invalidOpcode := Opcode(255) // Should be invalid
		def, err := Lookup(invalidOpcode)
		if err == nil {
			t.Error("Expected error for invalid opcode")
		}
		if def != nil {
			t.Error("Expected nil definition for invalid opcode")
		}
	})

	t.Run("OpcodeDefinitionStructure", func(t *testing.T) {
		// Test that OpConstant has correct operand width
		def, err := Lookup(OpConstant)
		if err != nil {
			t.Fatalf("Expected no error for OpConstant, got %v", err)
		}
		if len(def.OperandWidth) != 1 {
			t.Errorf("Expected OpConstant to have 1 operand, got %d", len(def.OperandWidth))
		}
		if def.OperandWidth[0] != 2 {
			t.Errorf("Expected OpConstant operand width to be 2, got %d", def.OperandWidth[0])
		}

		// Test that OpGetVar and OpSetVar have correct operand width
		getVarDef, err := Lookup(OpGetVar)
		if err != nil {
			t.Fatalf("Expected no error for OpGetVar, got %v", err)
		}
		if len(getVarDef.OperandWidth) != 1 {
			t.Errorf("Expected OpGetVar to have 1 operand, got %d", len(getVarDef.OperandWidth))
		}
		if getVarDef.OperandWidth[0] != 2 {
			t.Errorf("Expected OpGetVar operand width to be 2, got %d", getVarDef.OperandWidth[0])
		}

		setVarDef, err := Lookup(OpSetVar)
		if err != nil {
			t.Fatalf("Expected no error for OpSetVar, got %v", err)
		}
		if len(setVarDef.OperandWidth) != 1 {
			t.Errorf("Expected OpSetVar to have 1 operand, got %d", len(setVarDef.OperandWidth))
		}
		if setVarDef.OperandWidth[0] != 2 {
			t.Errorf("Expected OpSetVar operand width to be 2, got %d", setVarDef.OperandWidth[0])
		}
	})
}

// TestMake tests the Make function for creating bytecode instructions
func TestMake(t *testing.T) {
	t.Run("OpConstant", func(t *testing.T) {
		instruction := Make(OpConstant, 65534)
		expected := []byte{byte(OpConstant), 255, 254} // 65534 = 0xFFFE

		if len(instruction) != len(expected) {
			t.Errorf("Expected instruction length %d, got %d", len(expected), len(instruction))
		}

		for i, b := range expected {
			if instruction[i] != b {
				t.Errorf("Expected byte[%d] = %d, got %d", i, b, instruction[i])
			}
		}
	})

	t.Run("OpPop", func(t *testing.T) {
		instruction := Make(OpPop)
		expected := []byte{byte(OpPop)}

		if len(instruction) != len(expected) {
			t.Errorf("Expected instruction length %d, got %d", len(expected), len(instruction))
		}

		if instruction[0] != expected[0] {
			t.Errorf("Expected opcode %d, got %d", expected[0], instruction[0])
		}
	})

	t.Run("OpGetVar", func(t *testing.T) {
		instruction := Make(OpGetVar, 255)
		expected := []byte{byte(OpGetVar), 0, 255}

		if len(instruction) != len(expected) {
			t.Errorf("Expected instruction length %d, got %d", len(expected), len(instruction))
		}

		for i, b := range expected {
			if instruction[i] != b {
				t.Errorf("Expected byte[%d] = %d, got %d", i, b, instruction[i])
			}
		}
	})

	t.Run("MultipleOperands", func(t *testing.T) {
		// Test instruction with multiple operands - OpBuiltin has 2 operands
		instruction := Make(OpBuiltin, 1, 2)
		if len(instruction) < 1 {
			t.Error("Expected instruction to have at least opcode byte")
		}
		if instruction[0] != byte(OpBuiltin) {
			t.Errorf("Expected opcode %d, got %d", byte(OpBuiltin), instruction[0])
		}
	})

	t.Run("EdgeCaseValues", func(t *testing.T) {
		// Test with 0 value
		instruction := Make(OpConstant, 0)
		expected := []byte{byte(OpConstant), 0, 0}

		if len(instruction) != len(expected) {
			t.Errorf("Expected instruction length %d, got %d", len(expected), len(instruction))
		}

		for i, b := range expected {
			if instruction[i] != b {
				t.Errorf("Expected byte[%d] = %d, got %d", i, b, instruction[i])
			}
		}

		// Test with maximum 2-byte value
		instruction2 := Make(OpConstant, 65535)
		expected2 := []byte{byte(OpConstant), 255, 255}

		for i, b := range expected2 {
			if instruction2[i] != b {
				t.Errorf("Expected byte[%d] = %d, got %d", i, b, instruction2[i])
			}
		}
	})
}

// TestReadOperands tests the ReadOperands function
func TestReadOperands(t *testing.T) {
	t.Run("OpConstant", func(t *testing.T) {
		def, err := Lookup(OpConstant)
		if err != nil {
			t.Fatalf("Expected no error looking up OpConstant, got %v", err)
		}

		instruction := []byte{255, 254} // 65534 in big-endian
		operands, bytesRead := ReadOperands(def, instruction)

		if len(operands) != 1 {
			t.Errorf("Expected 1 operand, got %d", len(operands))
		}
		if operands[0] != 65534 {
			t.Errorf("Expected operand 65534, got %d", operands[0])
		}
		if bytesRead != 2 {
			t.Errorf("Expected 2 bytes read, got %d", bytesRead)
		}
	})

	t.Run("OpGetVar", func(t *testing.T) {
		def, err := Lookup(OpGetVar)
		if err != nil {
			t.Fatalf("Expected no error looking up OpGetVar, got %v", err)
		}

		instruction := []byte{0, 255} // 255 in big-endian
		operands, bytesRead := ReadOperands(def, instruction)

		if len(operands) != 1 {
			t.Errorf("Expected 1 operand, got %d", len(operands))
		}
		if operands[0] != 255 {
			t.Errorf("Expected operand 255, got %d", operands[0])
		}
		if bytesRead != 2 {
			t.Errorf("Expected 2 bytes read, got %d", bytesRead)
		}
	})

	t.Run("EmptyOperands", func(t *testing.T) {
		def, err := Lookup(OpPop)
		if err != nil {
			t.Fatalf("Expected no error looking up OpPop, got %v", err)
		}

		instruction := []byte{}
		operands, bytesRead := ReadOperands(def, instruction)

		if len(operands) != 0 {
			t.Errorf("Expected 0 operands, got %d", len(operands))
		}
		if bytesRead != 0 {
			t.Errorf("Expected 0 bytes read, got %d", bytesRead)
		}
	})

	t.Run("MultipleOperands", func(t *testing.T) {
		def, err := Lookup(OpBuiltin)
		if err != nil {
			t.Fatalf("Expected no error looking up OpBuiltin, got %v", err)
		}

		// OpBuiltin has 2 operands, test with multiple values
		if len(def.OperandWidth) > 0 {
			// Create instruction bytes based on operand widths
			var instruction []byte
			expectedOperands := []int{1, 2}

			for i, width := range def.OperandWidth {
				if i >= len(expectedOperands) {
					break
				}

				operand := expectedOperands[i]
				for j := width - 1; j >= 0; j-- {
					instruction = append(instruction, byte(operand>>(j*8)))
				}
			}

			operands, bytesRead := ReadOperands(def, instruction)

			if bytesRead != len(instruction) {
				t.Errorf("Expected %d bytes read, got %d", len(instruction), bytesRead)
			}

			// Check that we got the expected number of operands
			expectedCount := len(def.OperandWidth)
			if len(operands) != expectedCount {
				t.Errorf("Expected %d operands, got %d", expectedCount, len(operands))
			}
		}
	})

	t.Run("InsufficientBytes", func(t *testing.T) {
		def, err := Lookup(OpConstant)
		if err != nil {
			t.Fatalf("Expected no error looking up OpConstant, got %v", err)
		}

		// Only provide 1 byte when 2 are needed
		instruction := []byte{255}
		operands, bytesRead := ReadOperands(def, instruction)

		// Should handle gracefully (might return partial data or empty)
		if bytesRead > len(instruction) {
			t.Errorf("Expected bytes read <= %d, got %d", len(instruction), bytesRead)
		}

		// The function should not panic
		if len(operands) > 1 {
			t.Errorf("Expected at most 1 operand with insufficient bytes, got %d", len(operands))
		}
	})
}

// TestFormatInstruction tests the FormatInstruction function
func TestFormatInstruction(t *testing.T) {
	t.Run("OpConstant", func(t *testing.T) {
		def, err := Lookup(OpConstant)
		if err != nil {
			t.Fatalf("Expected no error looking up OpConstant, got %v", err)
		}

		operands := []int{65534}
		formatted := FormatInstruction(def, operands)

		if formatted == "" {
			t.Error("Expected non-empty formatted instruction")
		}

		// Should contain the opcode name
		if !contains(formatted, "OpConstant") {
			t.Errorf("Expected formatted instruction to contain 'OpConstant', got %s", formatted)
		}

		// Should contain the operand value
		if !contains(formatted, "65534") {
			t.Errorf("Expected formatted instruction to contain '65534', got %s", formatted)
		}
	})

	t.Run("OpPop", func(t *testing.T) {
		def, err := Lookup(OpPop)
		if err != nil {
			t.Fatalf("Expected no error looking up OpPop, got %v", err)
		}

		operands := []int{} // No operands
		formatted := FormatInstruction(def, operands)

		if formatted == "" {
			t.Error("Expected non-empty formatted instruction")
		}

		// Should contain the opcode name
		if !contains(formatted, "OpPop") {
			t.Errorf("Expected formatted instruction to contain 'OpPop', got %s", formatted)
		}
	})

	t.Run("MultipleOperands", func(t *testing.T) {
		def, err := Lookup(OpBuiltin)
		if err != nil {
			t.Fatalf("Expected no error looking up OpBuiltin, got %v", err)
		}

		operands := []int{1, 2}
		formatted := FormatInstruction(def, operands)

		if formatted == "" {
			t.Error("Expected non-empty formatted instruction")
		}

		// Should contain the opcode name
		if !contains(formatted, "OpBuiltin") {
			t.Errorf("Expected formatted instruction to contain 'OpBuiltin', got %s", formatted)
		}
	})

	t.Run("EmptyOperands", func(t *testing.T) {
		def, err := Lookup(OpAdd)
		if err != nil {
			t.Fatalf("Expected no error looking up OpAdd, got %v", err)
		}

		operands := []int{}
		formatted := FormatInstruction(def, operands)

		if formatted == "" {
			t.Error("Expected non-empty formatted instruction")
		}

		// Should contain the opcode name
		if !contains(formatted, "OpAdd") {
			t.Errorf("Expected formatted instruction to contain 'OpAdd', got %s", formatted)
		}
	})
}

// TestMakeReadOperandsRoundTrip tests that Make and ReadOperands are consistent
func TestMakeReadOperandsRoundTrip(t *testing.T) {
	testCases := []struct {
		opcode   Opcode
		operands []int
	}{
		{OpConstant, []int{65534}},
		{OpConstant, []int{0}},
		{OpConstant, []int{65535}},
		{OpGetVar, []int{255}},
		{OpGetVar, []int{0}},
		{OpSetVar, []int{100}},
	}

	for _, test := range testCases {
		def, err := Lookup(test.opcode)
		if err != nil {
			t.Fatalf("Expected no error looking up %v, got %v", test.opcode, err)
		}

		// Make instruction
		instruction := Make(test.opcode, test.operands...)

		// Skip the opcode byte and read operands
		if len(instruction) > 1 {
			operands, bytesRead := ReadOperands(def, instruction[1:])

			// Check that we read the expected number of bytes
			expectedBytes := len(instruction) - 1
			if bytesRead != expectedBytes {
				t.Errorf("Expected %d bytes read, got %d", expectedBytes, bytesRead)
			}

			// Check that operands match
			if len(operands) != len(test.operands) {
				t.Errorf("Expected %d operands, got %d", len(test.operands), len(operands))
				continue
			}

			for i, expected := range test.operands {
				if i < len(operands) && operands[i] != expected {
					t.Errorf("Expected operand[%d] = %d, got %d", i, expected, operands[i])
				}
			}
		}
	}
}

// TestOpcodeDefinitionConsistency tests that all opcodes have proper definitions
func TestOpcodeDefinitionConsistency(t *testing.T) {
	// Test that commonly used opcodes have definitions
	commonOpcodes := []Opcode{
		OpConstant, OpPop, OpAdd, OpSub, OpMul, OpDiv,
		OpEqual, OpNotEqual, OpGreaterThan, OpLessThan,
		OpAnd, OpOr, OpNot, OpGetVar, OpSetVar,
		OpCall, OpReturn, OpBuiltin,
	}

	for _, opcode := range commonOpcodes {
		def, err := Lookup(opcode)
		if err != nil {
			t.Errorf("Expected definition for common opcode %v, got error: %v", opcode, err)
			continue
		}
		if def == nil {
			t.Errorf("Expected non-nil definition for opcode %v", opcode)
			continue
		}
		if def.Name == "" {
			t.Errorf("Expected non-empty name for opcode %v", opcode)
		}

		// OperandWidth should be non-nil (but can be empty)
		if def.OperandWidth == nil {
			t.Errorf("Expected non-nil OperandWidth for opcode %v", opcode)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && s[:len(substr)] == substr ||
		findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
