package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
	"github.com/mredencom/expr/types"
)

// TestDebugPlaceholder 调试占位符处理
func TestDebugPlaceholder(t *testing.T) {
	fmt.Println("🔍 调试占位符处理")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	// 创建一个测试用的VM方法来查看占位符内容
	testData := types.NewSlice([]types.Value{
		types.NewInt(1), types.NewInt(6), types.NewInt(3),
	}, types.TypeInfo{Kind: types.KindInt64, Name: "int"})

	fmt.Printf("测试数据: %v\n", testData)

	// 测试1: 编译并查看常量
	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	fmt.Println("\n1. 编译 'numbers | filter(# > 5)' 并查看常量")
	program, err := expr.Compile("numbers | filter(# > 5)", expr.Env(env))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 编译成功，常量数量: %d\n", program.ConstantsCount())

	// 测试2: 编译简单占位符
	fmt.Println("\n2. 编译 '# > 5' 单独的占位符表达式")
	program2, err := expr.Compile("# > 5")
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功，常量数量: %d\n", program2.ConstantsCount())
	}

	// 测试3: 编译 filter 函数调用
	fmt.Println("\n3. 编译 'filter(numbers, # > 5)' 函数调用形式")
	program3, err := expr.Compile("filter(numbers, # > 5)", expr.Env(env))
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功，常量数量: %d\n", program3.ConstantsCount())

		// 尝试执行
		result, err := expr.Run(program3, env)
		if err != nil {
			fmt.Printf("❌ 执行错误: %v\n", err)
		} else {
			fmt.Printf("✅ 执行结果: %v\n", result)
		}
	}

	// 测试4: 直接测试占位符
	fmt.Println("\n4. 测试占位符 '#' 的编译")
	program4, err := expr.Compile("#")
	if err != nil {
		fmt.Printf("❌ 编译错误: %v\n", err)
	} else {
		fmt.Printf("✅ 编译成功，常量数量: %d\n", program4.ConstantsCount())
	}
}
