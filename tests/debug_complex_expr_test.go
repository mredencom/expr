package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

// TestDebugComplexExpression 调试复杂表达式的解析和编译
func TestDebugComplexExpression(t *testing.T) {
	fmt.Println("🔍 调试复杂表达式: words | filter(#.length() > 4)")

	// 测试表达式 - 使用字面量数组避免变量问题
	expr := `["hi", "hello", "world"] | filter(#.length() > 4)`

	fmt.Printf("测试表达式: %s\n", expr)

	// 1. 解析
	fmt.Println("\n1. 解析阶段:")
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("❌ 解析错误: %v\n", p.Errors())
		return
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	pipeExpr := stmt.Expression.(*ast.PipeExpression)

	fmt.Printf("✅ 左侧: %T - %s\n", pipeExpr.Left, pipeExpr.Left.String())
	fmt.Printf("✅ 右侧: %T - %s\n", pipeExpr.Right, pipeExpr.Right.String())

	// 分析右侧的BuiltinExpression
	if builtin, ok := pipeExpr.Right.(*ast.BuiltinExpression); ok {
		fmt.Printf("   函数名: %s\n", builtin.Name)
		fmt.Printf("   参数数量: %d\n", len(builtin.Arguments))

		for i, arg := range builtin.Arguments {
			fmt.Printf("   参数 %d: %T - %s\n", i, arg, arg.String())

			// 深入分析第一个参数（应该是InfixExpression）
			if infix, ok := arg.(*ast.InfixExpression); ok {
				fmt.Printf("     操作符: %s\n", infix.Operator)
				fmt.Printf("     左侧: %T - %s\n", infix.Left, infix.Left.String())
				fmt.Printf("     右侧: %T - %s\n", infix.Right, infix.Right.String())

				// 分析左侧（应该是#.length()）
				if call, ok := infix.Left.(*ast.CallExpression); ok {
					fmt.Printf("       调用函数: %T - %s\n", call.Function, call.Function.String())
					if member, ok := call.Function.(*ast.MemberExpression); ok {
						fmt.Printf("       对象: %T - %s\n", member.Object, member.Object.String())
						fmt.Printf("       属性: %T - %s\n", member.Property, member.Property.String())
					}
				}
			}
		}
	}

	// 2. 编译
	fmt.Println("\n2. 编译阶段:")
	comp := compiler.New()

	err := comp.Compile(stmt.Expression)
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
		return
	}

	bytecode := comp.Bytecode()
	fmt.Printf("✅ 编译成功\n")
	fmt.Printf("   常量数量: %d\n", len(bytecode.Constants))
	fmt.Printf("   指令长度: %d\n", len(bytecode.Instructions))

	// 打印常量
	for i, constant := range bytecode.Constants {
		fmt.Printf("   常量 %d: %T - %v\n", i, constant, constant)
	}
}
