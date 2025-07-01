package tests

import (
	"fmt"
	"testing"

	expr "github.com/mredencom/expr"
)

// TestComprehensiveFunctionality 全面测试所有计划功能
func TestComprehensiveFunctionality(t *testing.T) {
	fmt.Println("🔍 全面功能测试 - 验证plan vs 实现")
	fmt.Println("=" + fmt.Sprintf("%50s", "="))

	// 测试环境
	env := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"users": []map[string]interface{}{
			{"name": "Alice", "age": 30, "active": true, "score": 95.5},
			{"name": "Bob", "age": 25, "active": false, "score": 87.2},
			{"name": "Charlie", "age": 35, "active": true, "score": 92.8},
		},
		"text":      "Hello World",
		"threshold": 5,
	}

	// 1. 核心基础设施测试
	testCoreInfrastructure(t, env)

	// 2. 高级语言特性测试
	testAdvancedLanguageFeatures(t, env)

	// 3. 管道占位符语法测试
	testPipelinePlaceholderSyntax(t, env)

	// 4. 内置函数库测试
	testBuiltinFunctions(t, env)

	// 5. Lambda表达式测试
	testLambdaExpressions(t, env)

	// 6. 缺失功能检测
	testMissingFeatures(t, env)

	fmt.Println("\n📊 测试总结完成")
}

func testCoreInfrastructure(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n✅ 1. 核心基础设施测试")
	fmt.Println("------------------------")

	tests := []struct {
		name       string
		expression string
		expected   interface{}
	}{
		{"基础算术", "2 + 3 * 4", 14},
		{"字符串连接", "\"hello\" + \" \" + \"world\"", "hello world"},
		{"布尔逻辑", "true && false || true", true},
		{"三元运算", "5 > 3 ? \"yes\" : \"no\"", "yes"},
		{"变量访问", "threshold", 5},
		{"成员访问", "users[0].name", "Alice"},
		{"索引访问", "numbers[0]", 1},
	}

	passCount := 0
	for _, test := range tests {
		result, err := expr.Eval(test.expression, env)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", test.name, err)
		} else if result == test.expected {
			fmt.Printf("  ✅ %-12s: %v\n", test.name, result)
			passCount++
		} else {
			fmt.Printf("  ⚠️  %-12s: 期望 %v, 得到 %v\n", test.name, test.expected, result)
		}
	}

	fmt.Printf("核心基础设施: %d/%d 通过\n", passCount, len(tests))
}

func testAdvancedLanguageFeatures(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n✅ 2. 高级语言特性测试")
	fmt.Println("------------------------")

	tests := []struct {
		name       string
		expression string
		shouldPass bool
	}{
		{"数组字面量", "[1, 2, 3]", true},
		{"对象字面量", "{\"key\": \"value\"}", true},
		{"复杂成员访问", "users[0].active", true},
		{"嵌套索引", "users[0].name", true},
	}

	passCount := 0
	for _, test := range tests {
		_, err := expr.Eval(test.expression, env)
		if err != nil && test.shouldPass {
			fmt.Printf("  ❌ %-15s: %v\n", test.name, err)
		} else if err == nil && test.shouldPass {
			fmt.Printf("  ✅ %-15s: 编译和执行成功\n", test.name)
			passCount++
		} else if err != nil && !test.shouldPass {
			fmt.Printf("  ✅ %-15s: 正确拒绝\n", test.name)
			passCount++
		}
	}

	fmt.Printf("高级语言特性: %d/%d 通过\n", passCount, len(tests))
}

func testPipelinePlaceholderSyntax(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n🔥 3. 管道占位符语法测试 (核心创新)")
	fmt.Println("------------------------------------------")

	tests := []struct {
		name       string
		expression string
		shouldPass bool
	}{
		{"基础过滤", "numbers | filter(# > 5)", true},
		{"基础映射", "numbers | map(# * 2)", true},
		{"复杂条件", "numbers | filter(# % 2 == 0 && # > 3)", true},
		{"链式操作", "numbers | filter(# > 3) | map(# * 2)", true},
		{"对象属性", "users | filter(#.age > 25)", true},
		{"复杂表达式", "numbers | map((# + 1) * (# - 1))", true},
		{"聚合终结", "numbers | filter(# > 5) | sum", true},
	}

	passCount := 0
	for _, test := range tests {
		result, err := expr.Eval(test.expression, env)
		if err != nil {
			fmt.Printf("  ❌ %-15s: %v\n", test.name, err)
		} else {
			fmt.Printf("  ✅ %-15s: %v\n", test.name, result)
			passCount++
		}
	}

	fmt.Printf("管道占位符语法: %d/%d 通过\n", passCount, len(tests))
}

func testBuiltinFunctions(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n📚 4. 内置函数库测试")
	fmt.Println("---------------------")

	// 已实现的函数
	implementedTests := []struct {
		name       string
		expression string
		shouldPass bool
	}{
		// 基础函数
		{"len函数", "len(numbers)", true},
		{"string转换", "string(42)", true},
		{"int转换", "int(\"42\")", true},
		{"bool转换", "bool(1)", true},

		// 字符串函数
		{"contains", "contains(text, \"World\")", true},
		{"upper", "upper(text)", true},
		{"lower", "lower(text)", true},

		// 集合函数
		{"filter", "filter(numbers, x => x > 5)", true},
		{"map", "map(numbers, x => x * 2)", true},
		{"sum", "sum(numbers)", true},
		{"count", "count(numbers)", true},
		{"max", "max(numbers)", true},
		{"min", "min(numbers)", true},
	}

	passCount := 0
	for _, test := range implementedTests {
		_, err := expr.Eval(test.expression, env)
		if err != nil {
			fmt.Printf("  ❌ %-12s: %v\n", test.name, err)
		} else {
			fmt.Printf("  ✅ %-12s: 执行成功\n", test.name)
			passCount++
		}
	}

	fmt.Printf("已实现函数: %d/%d 通过\n", passCount, len(implementedTests))

	// 缺失的函数测试
	fmt.Println("\n❌ 缺失的内置函数:")
	missingFunctions := []string{
		"replace(text, \"World\", \"Universe\")", // 字符串替换
		"substring(text, 0, 5)",                  // 字符串截取
		"indexOf(text, \"World\")",               // 查找位置
		"ceil(3.14)",                             // 向上取整
		"floor(3.14)",                            // 向下取整
		"round(3.14)",                            // 四舍五入
		"sqrt(16)",                               // 平方根
		"pow(2, 3)",                              // 幂运算
		"flatten([[1, 2], [3, 4]])",              // 数组扁平化
		"groupBy(users, u => u.age > 30)",        // 分组
		"now()",                                  // 当前时间
	}

	for _, fn := range missingFunctions {
		_, err := expr.Eval(fn, env)
		if err != nil {
			fmt.Printf("  ❌ %s\n", fn)
		} else {
			fmt.Printf("  ✅ %s (意外成功)\n", fn)
		}
	}
}

func testLambdaExpressions(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n🔧 5. Lambda表达式测试")
	fmt.Println("----------------------")

	tests := []struct {
		name       string
		expression string
		shouldPass bool
	}{
		{"单参数Lambda", "x => x * 2", true},
		{"多参数Lambda", "(x, y) => x + y", true},
		{"Lambda在filter中", "filter(numbers, x => x > 5)", true},
		{"Lambda在map中", "map(numbers, x => x * 2)", true},
		{"复杂Lambda", "filter(users, u => u.age > 25 && u.active)", true},
	}

	passCount := 0
	for _, test := range tests {
		// 测试编译
		_, err := expr.Compile(test.expression, expr.Env(env))
		if err != nil {
			fmt.Printf("  ❌ %-15s: 编译失败 - %v\n", test.name, err)
		} else {
			fmt.Printf("  ✅ %-15s: 编译成功\n", test.name)
			passCount++

			// 如果可能，测试执行
			if test.name == "单参数Lambda" || test.name == "多参数Lambda" {
				// 这些是纯Lambda，无法直接执行
				continue
			}

			_, execErr := expr.Eval(test.expression, env)
			if execErr != nil {
				fmt.Printf("      ⚠️  执行失败: %v\n", execErr)
			} else {
				fmt.Printf("      ✅ 执行成功\n")
			}
		}
	}

	fmt.Printf("Lambda表达式: %d/%d 编译通过\n", passCount, len(tests))
}

func testMissingFeatures(t *testing.T, env map[string]interface{}) {
	fmt.Println("\n❌ 6. 缺失功能检测")
	fmt.Println("-------------------")

	fmt.Println("模块系统:")
	moduleTests := []string{
		"import \"math\" as m",
		"m.sqrt(16)",
	}

	for _, test := range moduleTests {
		_, err := expr.Eval(test, env)
		if err != nil {
			fmt.Printf("  ❌ %s\n", test)
		} else {
			fmt.Printf("  ✅ %s (意外成功)\n", test)
		}
	}

	fmt.Println("\n错误处理增强:")
	fmt.Println("  ❌ 详细错误位置信息")
	fmt.Println("  ❌ 错误恢复机制")
	fmt.Println("  ❌ 错误建议功能")

	fmt.Println("\n调试工具:")
	fmt.Println("  ❌ 表达式调试器")
	fmt.Println("  ❌ 性能分析器")
	fmt.Println("  ❌ 字节码可视化")

	fmt.Println("\n高级优化:")
	fmt.Println("  ❌ JIT编译")
	fmt.Println("  ❌ SIMD指令")
	fmt.Println("  ❌ 分支预测优化")
}
