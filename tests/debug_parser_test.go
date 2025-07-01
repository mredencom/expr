package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
)

// TestDebugParser 测试解析器如何解析 #.upper()
func TestDebugParser(t *testing.T) {
	fmt.Println("🔍 调试解析器处理 #.upper()")
	fmt.Println("=" + fmt.Sprintf("%40s", "="))

	expressions := []string{
		"#",
		"#.upper",
		"#.upper()",
		"# > 5",
		"#.length() > 5",
		"words | map(#.upper())",
	}

	for _, expr := range expressions {
		fmt.Printf("\n📝 解析表达式: %s\n", expr)

		l := lexer.New(expr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("❌ 解析错误: %v\n", p.Errors())
			continue
		}

		if len(program.Statements) == 0 {
			fmt.Printf("❌ 没有解析到语句\n")
			continue
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			fmt.Printf("❌ 不是表达式语句\n")
			continue
		}

		fmt.Printf("✅ AST 节点类型: %T\n", stmt.Expression)
		fmt.Printf("✅ AST 字符串: %s\n", stmt.Expression.String())

		// 详细分析特定类型的表达式
		switch e := stmt.Expression.(type) {
		case *ast.CallExpression:
			fmt.Printf("   📞 函数调用: %T\n", e.Function)
			if member, ok := e.Function.(*ast.MemberExpression); ok {
				fmt.Printf("   👉 对象: %T - %s\n", member.Object, member.Object.String())
				fmt.Printf("   👉 属性: %T - %s\n", member.Property, member.Property.String())
			}
		case *ast.MemberExpression:
			fmt.Printf("   👉 对象: %T - %s\n", e.Object, e.Object.String())
			fmt.Printf("   👉 属性: %T - %s\n", e.Property, e.Property.String())
		case *ast.PipeExpression:
			fmt.Printf("   🔄 左侧: %T - %s\n", e.Left, e.Left.String())
			fmt.Printf("   🔄 右侧: %T - %s\n", e.Right, e.Right.String())

			// 如果右侧是函数调用，进一步分析
			if call, ok := e.Right.(*ast.CallExpression); ok {
				fmt.Printf("   📞 右侧函数: %T\n", call.Function)
				if member, ok := call.Function.(*ast.MemberExpression); ok {
					fmt.Printf("   👉 函数对象: %T - %s\n", member.Object, member.Object.String())
					fmt.Printf("   👉 函数属性: %T - %s\n", member.Property, member.Property.String())
				}
			}
		}
	}
}
