package main

import (
	"fmt"

	"github.com/mredencom/expr/vm"
)

func main() {
	fmt.Printf("OpAnd = %d\n", vm.OpAnd)
	fmt.Printf("OpOr = %d\n", vm.OpOr)
	fmt.Printf("OpNot = %d\n", vm.OpNot)

	// Print all opcodes
	fmt.Println("\nAll opcodes:")
	for i := 0; i < 100; i++ {
		op := vm.Opcode(i)
		fmt.Printf("%d: %s\n", i, op.String())
		if op.String() == fmt.Sprintf("Unknown(%d)", i) {
			break
		}
	}
}
