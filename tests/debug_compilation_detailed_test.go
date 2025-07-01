package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

// TestDebugCompilationDetailed 详细分析编译过程
func TestDebugCompilationDetailed(t *testing.T) {
	fmt.Println("🔍 详细分析 #.length() 的编译过程")

	// 测试1：解析 #.length()
	fmt.Printf("\n1. 解析 #.length():\n")
	expr1 := `#.length()`
	l1 := lexer.New(expr1)
	p1 := parser.New(l1)
	program1 := p1.ParseProgram()

	if len(p1.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p1.Errors())
	} else {
		stmt1 := program1.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt1.Expression)

		if call, ok := stmt1.Expression.(*ast.CallExpression); ok {
			fmt.Printf("   是CallExpression: %s\n", call.String())
			fmt.Printf("   Function类型: %T\n", call.Function)

			if member, ok := call.Function.(*ast.MemberExpression); ok {
				fmt.Printf("   MemberExpression.Object: %T (%s)\n", member.Object, member.Object.String())
				fmt.Printf("   MemberExpression.Property: %T (%s)\n", member.Property, member.Property.String())
			}
		}
	}

	// 测试2：解析 #.length() > 4
	fmt.Printf("\n2. 解析 #.length() > 4:\n")
	expr2 := `#.length() > 4`
	l2 := lexer.New(expr2)
	p2 := parser.New(l2)
	program2 := p2.ParseProgram()

	if len(p2.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p2.Errors())
	} else {
		stmt2 := program2.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt2.Expression)

		if infix, ok := stmt2.Expression.(*ast.InfixExpression); ok {
			fmt.Printf("   是InfixExpression: %s\n", infix.String())
			fmt.Printf("   操作符: %s\n", infix.Operator)
			fmt.Printf("   左操作数类型: %T (%s)\n", infix.Left, infix.Left.String())
			fmt.Printf("   右操作数类型: %T (%s)\n", infix.Right, infix.Right.String())

			// 检查左操作数是否是CallExpression
			if call, ok := infix.Left.(*ast.CallExpression); ok {
				fmt.Printf("   左操作数是CallExpression!\n")
				fmt.Printf("   Function类型: %T\n", call.Function)

				if member, ok := call.Function.(*ast.MemberExpression); ok {
					fmt.Printf("   MemberExpression.Object: %T (%s)\n", member.Object, member.Object.String())
					fmt.Printf("   MemberExpression.Property: %T (%s)\n", member.Property, member.Property.String())

					// 检查是否是占位符
					if placeholder, ok := member.Object.(*ast.Identifier); ok && placeholder.Value == "#" {
						fmt.Printf("   ✅ 找到占位符类型方法调用!\n")
					}
				}
			}
		}
	}

	// 测试3：解析 filter(#.length() > 4)
	fmt.Printf("\n3. 解析 filter(#.length() > 4):\n")
	expr3 := `filter(#.length() > 4)`
	l3 := lexer.New(expr3)
	p3 := parser.New(l3)
	program3 := p3.ParseProgram()

	if len(p3.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p3.Errors())
	} else {
		stmt3 := program3.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt3.Expression)

		if call, ok := stmt3.Expression.(*ast.CallExpression); ok {
			fmt.Printf("   是CallExpression (filter调用): %s\n", call.String())
			fmt.Printf("   参数数量: %d\n", len(call.Arguments))

			if len(call.Arguments) > 0 {
				arg := call.Arguments[0]
				fmt.Printf("   第一个参数类型: %T (%s)\n", arg, arg.String())

				if infix, ok := arg.(*ast.InfixExpression); ok {
					fmt.Printf("   参数是InfixExpression: %s\n", infix.String())
					fmt.Printf("   左操作数: %T (%s)\n", infix.Left, infix.Left.String())

					if leftCall, ok := infix.Left.(*ast.CallExpression); ok {
						fmt.Printf("   左操作数是CallExpression!\n")
						if member, ok := leftCall.Function.(*ast.MemberExpression); ok {
							if placeholder, ok := member.Object.(*ast.Identifier); ok && placeholder.Value == "#" {
								fmt.Printf("   ✅ 找到嵌套的占位符类型方法调用!\n")
								fmt.Printf("   方法名: %s\n", member.Property.String())
							}
						}
					}
				}
			}
		}
	}
}
