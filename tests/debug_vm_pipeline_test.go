package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// TestDebugVMPipeline 调试VM管道处理
func TestDebugVMPipeline(t *testing.T) {
	fmt.Println("🔍 调试VM管道处理")
	fmt.Println("=" + fmt.Sprintf("%30s", "="))

	// 创建自定义VM来添加调试输出
	originalVM := &DebugVM{}

	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	// 测试管道操作
	fmt.Println("\n测试: numbers | filter(# > 5)")
	result, err := expr.Eval("numbers | filter(# > 5)", env)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v\n", result)
	}

	_ = originalVM
}

// DebugVM 包装VM以添加调试功能
type DebugVM struct {
	*vm.VM
}

// 重写callBuiltinFunction以添加调试输出
func (dvm *DebugVM) callBuiltinFunction(funcName string, args []types.Value) (types.Value, error) {
	fmt.Printf("🔍 调用内置函数: %s\n", funcName)
	fmt.Printf("   参数数量: %d\n", len(args))

	for i, arg := range args {
		fmt.Printf("   参数[%d]: 类型=%T, 值=%v\n", i, arg, arg)
	}

	if len(args) == 0 {
		return types.NewNil(), fmt.Errorf("builtin function %s requires at least one argument", funcName)
	}

	data := args[0] // First argument is the data being piped

	switch funcName {
	case "filter":
		if len(args) < 2 {
			return types.NewNil(), fmt.Errorf("filter requires a condition")
		}
		fmt.Printf("🔍 Filter条件: 类型=%T, 值=%v\n", args[1], args[1])
		return dvm.executeFilter(data, args[1])
	default:
		return types.NewNil(), fmt.Errorf("unknown builtin function: %s", funcName)
	}
}

// 重写executeFilter以添加调试输出
func (dvm *DebugVM) executeFilter(data types.Value, condition types.Value) (types.Value, error) {
	fmt.Printf("🔍 执行Filter:\n")
	fmt.Printf("   数据类型: %T\n", data)
	fmt.Printf("   条件类型: %T\n", condition)
	fmt.Printf("   条件值: %v\n", condition)

	slice, ok := data.(*types.SliceValue)
	if !ok {
		return types.NewNil(), fmt.Errorf("filter can only be applied to arrays")
	}

	var result []types.Value
	elements := slice.Values()

	fmt.Printf("   处理 %d 个元素\n", len(elements))

	for i, element := range elements {
		fmt.Printf("   元素[%d]: %v\n", i, element)

		// 这里是关键：看看condition到底是什么
		conditionResult := dvm.evaluatePlaceholderCondition(condition, element)
		fmt.Printf("   条件结果: %v\n", conditionResult)

		if dvm.isTruthy(conditionResult) {
			result = append(result, element)
		}
	}

	elemType := dvm.getSliceElementType(slice)
	return types.NewSlice(result, elemType), nil
}

// 辅助方法
func (dvm *DebugVM) evaluatePlaceholderCondition(condition types.Value, element types.Value) types.Value {
	fmt.Printf("     🔍 评估占位符条件:\n")
	fmt.Printf("       条件类型: %T\n", condition)
	fmt.Printf("       元素类型: %T, 值: %v\n", element, element)

	// 检查是否是PlaceholderExprValue
	if placeholderExpr, ok := condition.(*types.PlaceholderExprValue); ok {
		fmt.Printf("       ✅ 是PlaceholderExprValue\n")
		fmt.Printf("       操作符: %s\n", placeholderExpr.Operator())
		fmt.Printf("       操作数: %v\n", placeholderExpr.Operand())
		return types.NewBool(true) // 简化返回
	}

	// 检查是否是字符串
	if condStr, ok := condition.(*types.StringValue); ok {
		fmt.Printf("       ❌ 是StringValue: %s\n", condStr.Value())
		return types.NewBool(true) // 简化返回
	}

	fmt.Printf("       ❓ 未知类型\n")
	return types.NewBool(true)
}

func (dvm *DebugVM) isTruthy(value types.Value) bool {
	if boolVal, ok := value.(*types.BoolValue); ok {
		return boolVal.Value()
	}
	return true
}

func (dvm *DebugVM) getSliceElementType(slice *types.SliceValue) types.TypeInfo {
	if len(slice.Values()) > 0 {
		firstElement := slice.Values()[0]
		return firstElement.Type()
	}
	return types.TypeInfo{Kind: types.KindInterface, Name: "interface{}", Size: -1}
}
