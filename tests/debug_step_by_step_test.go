package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/builtins"
	"github.com/mredencom/expr/types"
)

// TestDebugStepByStep 逐步调试表达式求值过程
func TestDebugStepByStep(t *testing.T) {
	fmt.Println("🔍 逐步调试表达式求值过程")

	// 测试字符串：应该长度为5，大于4
	testString := "hello"
	testStringValue := types.NewString(testString)

	fmt.Printf("测试字符串: \"%s\"\n", testString)
	fmt.Printf("测试字符串长度: %d\n", len(testString))
	fmt.Printf("预期 length() > 4: %t\n", len(testString) > 4)

	// 1. 测试string.length方法
	fmt.Printf("\n1. 测试 string.length 方法:\n")
	if lengthMethod, exists := builtins.TypeMethodBuiltins["string.length"]; exists {
		lengthResult, err := lengthMethod([]types.Value{testStringValue})
		if err != nil {
			fmt.Printf("   ❌ 错误: %v\n", err)
		} else {
			fmt.Printf("   ✅ 结果: %v (%T)\n", lengthResult, lengthResult)
			if intVal, ok := lengthResult.(*types.IntValue); ok {
				fmt.Printf("   ✅ 整数值: %d\n", intVal.Value())
			}
		}
	} else {
		fmt.Printf("   ❌ string.length 方法不存在\n")
	}

	// 2. 手动模拟比较过程
	fmt.Printf("\n2. 手动模拟比较过程:\n")

	// 创建长度值
	lengthValue := types.NewInt(int64(len(testString)))
	fmt.Printf("   左值 (长度): %v (%T)\n", lengthValue, lengthValue)

	// 创建常量4
	constantFour := types.NewInt(4)
	fmt.Printf("   右值 (常量4): %v (%T)\n", constantFour, constantFour)

	// 手动执行比较
	leftInt := lengthValue.Value()
	rightInt := constantFour.Value()

	fmt.Printf("   比较: %d > %d = %t\n", leftInt, rightInt, leftInt > rightInt)

	// 3. 使用VM的类型转换函数测试
	fmt.Printf("\n3. 测试类型转换函数:\n")

	// 模拟VM的tryConvertToInt函数
	testTryConvertToInt := func(value types.Value) (int64, bool) {
		switch v := value.(type) {
		case *types.IntValue:
			return v.Value(), true
		case *types.FloatValue:
			return int64(v.Value()), true
		default:
			return 0, false
		}
	}

	leftConverted, leftOk := testTryConvertToInt(lengthValue)
	rightConverted, rightOk := testTryConvertToInt(constantFour)

	fmt.Printf("   左值转换: %d, 成功=%t\n", leftConverted, leftOk)
	fmt.Printf("   右值转换: %d, 成功=%t\n", rightConverted, rightOk)

	if leftOk && rightOk {
		result := leftConverted > rightConverted
		fmt.Printf("   比较结果: %d > %d = %t\n", leftConverted, rightConverted, result)
	}

	// 4. 测试实际的字符串数组
	fmt.Printf("\n4. 测试字符串数组:\n")
	words := []string{"hi", "hello", "world"}
	for i, word := range words {
		length := len(word)
		shouldInclude := length > 4
		fmt.Printf("   words[%d] = \"%s\", 长度=%d, >4=%t\n", i, word, length, shouldInclude)
	}

	// 5. 验证string.length对每个字符串的结果
	fmt.Printf("\n5. 验证string.length对每个字符串的结果:\n")
	if lengthMethod, exists := builtins.TypeMethodBuiltins["string.length"]; exists {
		for i, word := range words {
			wordValue := types.NewString(word)
			lengthResult, err := lengthMethod([]types.Value{wordValue})
			if err != nil {
				fmt.Printf("   words[%d] (\"%s\") length() 错误: %v\n", i, word, err)
			} else {
				if intVal, ok := lengthResult.(*types.IntValue); ok {
					actualLength := intVal.Value()
					greaterThan4 := actualLength > 4
					fmt.Printf("   words[%d] (\"%s\") length()=%d, >4=%t\n", i, word, actualLength, greaterThan4)
				}
			}
		}
	}
}
