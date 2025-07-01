package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/ast"
	"github.com/mredencom/expr/compiler"
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/parser"
	"github.com/mredencom/expr/types"
)

// TestDebugExpressionStructure 调试表达式结构
func TestDebugExpressionStructure(t *testing.T) {
	fmt.Println("🔍 调试表达式结构和执行流程")

	// 测试表达式
	expr := `["hi", "hello", "world"] | filter(#.length() > 4)`

	fmt.Printf("测试表达式: %s\n", expr)

	// 1. 解析
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

	// 2. 分析字节码常量
	fmt.Println("\n📋 字节码常量详细分析:")
	for i, constant := range bytecode.Constants {
		fmt.Printf("  [%d] %T: %v\n", i, constant, constant)

		// 如果是切片，深入分析其结构
		if slice, ok := constant.(*types.SliceValue); ok {
			fmt.Printf("      切片详细内容:\n")
			analyzeSliceStructure(slice, "        ")
		}
	}

	// 3. 手动分析第4个常量（应该是复杂表达式）
	if len(bytecode.Constants) > 3 {
		fmt.Println("\n🔍 第4个常量详细分析（应该是完整的filter调用）:")
		if slice, ok := bytecode.Constants[3].(*types.SliceValue); ok {
			analyzeFilterCall(slice)
		}
	}
}

// analyzeSliceStructure 递归分析切片结构
func analyzeSliceStructure(slice *types.SliceValue, indent string) {
	elements := slice.Values()
	for j, elem := range elements {
		fmt.Printf("%s[%d] %T: %v\n", indent, j, elem, elem)

		// 如果元素也是切片，递归分析
		if nestedSlice, ok := elem.(*types.SliceValue); ok {
			fmt.Printf("%s    嵌套切片:\n", indent)
			analyzeSliceStructure(nestedSlice, indent+"      ")
		}
	}
}

// analyzeFilterCall 分析filter调用的结构
func analyzeFilterCall(slice *types.SliceValue) {
	elements := slice.Values()
	fmt.Printf("Filter调用包含 %d 个元素:\n", len(elements))

	for i, elem := range elements {
		fmt.Printf("  [%d] %T: %v\n", i, elem, elem)

		if strVal, ok := elem.(*types.StringValue); ok {
			switch strVal.Value() {
			case "filter":
				fmt.Printf("      → 这是filter函数名\n")
			case "__PIPELINE_COMPLEX_TYPE_METHOD__":
				fmt.Printf("      → 这是复杂类型方法标记\n")
			case "length":
				fmt.Printf("      → 这是方法名\n")
			case ">":
				fmt.Printf("      → 这是比较操作符\n")
			case "__PIPELINE_MEMBER_ACCESS__":
				fmt.Printf("      → 这是成员访问标记\n")
			case "__PLACEHOLDER__":
				fmt.Printf("      → 这是占位符\n")
			default:
				fmt.Printf("      → 字符串值: %s\n", strVal.Value())
			}
		} else if intVal, ok := elem.(*types.IntValue); ok {
			fmt.Printf("      → 整数值: %d\n", intVal.Value())
		} else if slice, ok := elem.(*types.SliceValue); ok {
			fmt.Printf("      → 嵌套表达式，包含 %d 个元素:\n", len(slice.Values()))
			analyzeSliceStructure(slice, "        ")
		}
	}
}
