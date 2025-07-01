package tests

import (
	"fmt"
	"testing"

	"github.com/mredencom/expr/modules"
	"github.com/mredencom/expr/types"
)

func TestModuleSystemBasics(t *testing.T) {
	fmt.Println("\n🔧 模块系统基础测试")
	fmt.Println("========================")

	// Test 1: Check if built-in modules are registered
	fmt.Print("  检查内置模块注册      : ")
	mathModule, err := modules.DefaultRegistry.GetModule("math")
	if err != nil {
		t.Errorf("Math module not found: %v", err)
		return
	}
	fmt.Printf("✅ math模块已注册 (函数数量: %d)\n", len(mathModule.Functions))

	stringsModule, err := modules.DefaultRegistry.GetModule("strings")
	if err != nil {
		t.Errorf("Strings module not found: %v", err)
		return
	}
	fmt.Printf("  检查strings模块注册   : ✅ strings模块已注册 (函数数量: %d)\n", len(stringsModule.Functions))

	// Test 2: Test module function calls directly
	fmt.Print("  测试math.sqrt直接调用 : ")
	result, err := modules.DefaultRegistry.CallFunction("math", "sqrt", 16.0)
	if err != nil {
		t.Errorf("Math sqrt call failed: %v", err)
		return
	}
	if result != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
		return
	}
	fmt.Printf("✅ sqrt(16) = %v\n", result)

	fmt.Print("  测试strings.upper调用: ")
	result, err = modules.DefaultRegistry.CallFunction("strings", "upper", "hello")
	if err != nil {
		t.Errorf("Strings upper call failed: %v", err)
		return
	}
	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got %v", result)
		return
	}
	fmt.Printf("✅ upper('hello') = '%v'\n", result)

	fmt.Println("模块系统基础测试: 全部通过 ✅")
}

func TestModuleRegistration(t *testing.T) {
	fmt.Println("\n🔧 自定义模块注册测试")
	fmt.Println("========================")

	// Create a custom module registry for testing
	registry := modules.NewRegistry()

	// Define custom functions
	customFunctions := map[string]*modules.ModuleFunction{
		"double": {
			Name:        "double",
			Description: "Doubles a number",
			Handler: func(args ...interface{}) (interface{}, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("double expects 1 argument, got %d", len(args))
				}
				switch v := args[0].(type) {
				case int:
					return v * 2, nil
				case int64:
					return v * 2, nil
				case float64:
					return v * 2, nil
				default:
					return nil, fmt.Errorf("double expects numeric argument, got %T", v)
				}
			},
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindFloat64, Name: "float64"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindFloat64, Name: "float64"},
			Variadic:   false,
		},
		"greet": {
			Name:        "greet",
			Description: "Greets a person",
			Handler: func(args ...interface{}) (interface{}, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("greet expects 1 argument, got %d", len(args))
				}
				name := fmt.Sprintf("%v", args[0])
				return fmt.Sprintf("Hello, %s!", name), nil
			},
			ParamTypes: []types.TypeInfo{
				{Kind: types.KindString, Name: "string"},
			},
			ReturnType: types.TypeInfo{Kind: types.KindString, Name: "string"},
			Variadic:   false,
		},
	}

	// Register custom module
	fmt.Print("  注册自定义模块        : ")
	err := registry.RegisterModule("custom", "Custom functions for testing", customFunctions)
	if err != nil {
		t.Errorf("Failed to register custom module: %v", err)
		return
	}
	fmt.Println("✅ 自定义模块注册成功")

	// Test custom module functions
	fmt.Print("  测试custom.double调用 : ")
	result, err := registry.CallFunction("custom", "double", 5.0)
	if err != nil {
		t.Errorf("Custom double call failed: %v", err)
		return
	}
	if result != 10.0 {
		t.Errorf("Expected 10.0, got %v", result)
		return
	}
	fmt.Printf("✅ double(5.0) = %v\n", result)

	fmt.Print("  测试custom.greet调用  : ")
	result, err = registry.CallFunction("custom", "greet", "World")
	if err != nil {
		t.Errorf("Custom greet call failed: %v", err)
		return
	}
	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expected '%s', got %v", expected, result)
		return
	}
	fmt.Printf("✅ greet('World') = '%v'\n", result)

	// Test module listing
	fmt.Print("  列出所有模块          : ")
	moduleNames := registry.ListModules()
	if len(moduleNames) < 3 { // math, strings, custom
		t.Errorf("Expected at least 3 modules, got %d", len(moduleNames))
		return
	}
	fmt.Printf("✅ 找到 %d 个模块: %v\n", len(moduleNames), moduleNames)

	fmt.Println("自定义模块注册测试: 全部通过 ✅")
}

func TestModuleFunctionValidation(t *testing.T) {
	fmt.Println("\n🔧 模块函数验证测试")
	fmt.Println("====================")

	// Test error cases
	fmt.Print("  测试不存在的模块      : ")
	_, err := modules.DefaultRegistry.CallFunction("nonexistent", "function", 1, 2, 3)
	if err == nil {
		t.Errorf("Expected error for nonexistent module")
		return
	}
	fmt.Printf("✅ 正确返回错误: %v\n", err)

	fmt.Print("  测试不存在的函数      : ")
	_, err = modules.DefaultRegistry.CallFunction("math", "nonexistent", 1, 2, 3)
	if err == nil {
		t.Errorf("Expected error for nonexistent function")
		return
	}
	fmt.Printf("✅ 正确返回错误: %v\n", err)

	fmt.Print("  测试参数数量错误      : ")
	_, err = modules.DefaultRegistry.CallFunction("math", "sqrt") // sqrt requires 1 argument
	if err == nil {
		t.Errorf("Expected error for wrong argument count")
		return
	}
	fmt.Printf("✅ 正确返回错误: %v\n", err)

	fmt.Println("模块函数验证测试: 全部通过 ✅")
}
