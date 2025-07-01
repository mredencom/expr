package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// TestDebugComplexExecution 调试复杂表达式的执行过程
func TestDebugComplexExecution(t *testing.T) {
	fmt.Println("🔍 调试复杂表达式执行过程")

	// 测试表达式
	expr := `["hi", "hello", "world"] | filter(#.length() > 4)`

	fmt.Printf("测试表达式: %s\n", expr)

	// 1. 解析和编译
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("❌ 解析错误: %v\n", p.Errors())
		return
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	comp := compiler.New()
	err := comp.Compile(stmt.Expression)
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
		return
	}

	bytecode := comp.Bytecode()
	fmt.Printf("✅ 编译成功\n")

	// 打印详细的常量信息
	fmt.Println("\n📋 编译后的常量:")
	for i, constant := range bytecode.Constants {
		fmt.Printf("  [%d] %T: %v\n", i, constant, constant)

		// 如果是切片，打印其内容
		if slice, ok := constant.(*types.SliceValue); ok {
			fmt.Printf("      切片内容:\n")
			for j, elem := range slice.Values() {
				fmt.Printf("        [%d] %T: %v\n", j, elem, elem)
			}
		}
	}

	// 2. 尝试执行
	fmt.Println("\n🚀 执行阶段:")
	machine := vm.New(bytecode)

	// 环境
	env := map[string]interface{}{}

	result, err := machine.Run(bytecode, env)
	if err != nil {
		fmt.Printf("❌ 执行错误: %v\n", err)
		fmt.Println("分析错误原因...")
	} else {
		fmt.Printf("✅ 执行成功: %v\n", result)
	}
}
