package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

// TestDebugPlaceholderReplacement 调试占位符替换问题
func TestDebugPlaceholderReplacement(t *testing.T) {
	fmt.Println("🔍 调试占位符替换问题")
	fmt.Println("==========================================")

	// 测试1：简单的占位符 filter(#)
	fmt.Printf("\n1. 解析 filter(#):\n")
	expr1 := `filter(#)`
	l1 := lexer.New(expr1)
	p1 := parser.New(l1)
	program1 := p1.ParseProgram()

	var stmt1 *ast.ExpressionStatement
	if len(p1.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p1.Errors())
		return
	} else {
		stmt1 = program1.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt1.Expression)

		if builtin, ok := stmt1.Expression.(*ast.BuiltinExpression); ok {
			fmt.Printf("   函数名: %s\n", builtin.Name)
			fmt.Printf("   参数数量: %d\n", len(builtin.Arguments))

			if len(builtin.Arguments) > 0 {
				arg := builtin.Arguments[0]
				fmt.Printf("   第一个参数类型: %T (%s)\n", arg, arg.String())

				if placeholder, ok := arg.(*ast.PlaceholderExpression); ok {
					fmt.Printf("   ✅ 是占位符表达式: %s\n", placeholder.String())
				}
			}
		}
	}

	// 编译这个表达式
	fmt.Printf("\n2. 编译 filter(#):\n")
	comp1 := compiler.New()
	err1 := comp1.Compile(stmt1.Expression)
	if err1 != nil {
		fmt.Printf("   ❌ 编译错误: %v\n", err1)
	} else {
		fmt.Printf("   ✅ 编译成功\n")
		bytecode1 := comp1.Bytecode()
		fmt.Printf("   常量数量: %d\n", len(bytecode1.Constants))

		for i, constant := range bytecode1.Constants {
			fmt.Printf("   常量[%d]: %v (%T)\n", i, constant, constant)
		}
	}

	// 测试2：复杂的占位符 filter(#.length() > 4)
	fmt.Printf("\n3. 解析 filter(#.length() > 4):\n")
	expr2 := `filter(#.length() > 4)`
	l2 := lexer.New(expr2)
	p2 := parser.New(l2)
	program2 := p2.ParseProgram()

	var stmt2 *ast.ExpressionStatement
	if len(p2.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p2.Errors())
		return
	} else {
		stmt2 = program2.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt2.Expression)

		if builtin, ok := stmt2.Expression.(*ast.BuiltinExpression); ok {
			fmt.Printf("   函数名: %s\n", builtin.Name)
			fmt.Printf("   参数数量: %d\n", len(builtin.Arguments))

			if len(builtin.Arguments) > 0 {
				arg := builtin.Arguments[0]
				fmt.Printf("   第一个参数类型: %T (%s)\n", arg, arg.String())

				if infix, ok := arg.(*ast.InfixExpression); ok {
					fmt.Printf("   中缀表达式操作符: %s\n", infix.Operator)
					fmt.Printf("   左操作数: %T (%s)\n", infix.Left, infix.Left.String())
					fmt.Printf("   右操作数: %T (%s)\n", infix.Right, infix.Right.String())
				}
			}
		}
	}

	// 编译这个表达式
	fmt.Printf("\n4. 编译 filter(#.length() > 4):\n")
	comp2 := compiler.New()
	err2 := comp2.Compile(stmt2.Expression)
	if err2 != nil {
		fmt.Printf("   ❌ 编译错误: %v\n", err2)
	} else {
		fmt.Printf("   ✅ 编译成功\n")
		bytecode2 := comp2.Bytecode()
		fmt.Printf("   常量数量: %d\n", len(bytecode2.Constants))

		for i, constant := range bytecode2.Constants {
			fmt.Printf("   常量[%d]: %v (%T)\n", i, constant, constant)
		}
	}

	// 测试3：pipeline形式
	fmt.Printf("\n5. 解析 [1,2,3] | filter(#):\n")
	expr3 := `[1,2,3] | filter(#)`
	l3 := lexer.New(expr3)
	p3 := parser.New(l3)
	program3 := p3.ParseProgram()

	var stmt3 *ast.ExpressionStatement
	if len(p3.Errors()) > 0 {
		fmt.Printf("   ❌ 解析错误: %v\n", p3.Errors())
		return
	} else {
		stmt3 = program3.Statements[0].(*ast.ExpressionStatement)
		fmt.Printf("   ✅ 解析成功: %T\n", stmt3.Expression)

		if pipe, ok := stmt3.Expression.(*ast.PipeExpression); ok {
			fmt.Printf("   左操作数: %T (%s)\n", pipe.Left, pipe.Left.String())
			fmt.Printf("   右操作数: %T (%s)\n", pipe.Right, pipe.Right.String())
		}
	}

	// 编译这个表达式
	fmt.Printf("\n6. 编译 [1,2,3] | filter(#):\n")
	comp3 := compiler.New()
	err3 := comp3.Compile(stmt3.Expression)
	if err3 != nil {
		fmt.Printf("   ❌ 编译错误: %v\n", err3)
	} else {
		fmt.Printf("   ✅ 编译成功\n")
		bytecode3 := comp3.Bytecode()
		fmt.Printf("   常量数量: %d\n", len(bytecode3.Constants))

		for i, constant := range bytecode3.Constants {
			fmt.Printf("   常量[%d]: %v (%T)\n", i, constant, constant)
		}
	}
}
