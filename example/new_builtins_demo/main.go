package main

import (
	"fmt"
	"log"

	"github.com/mredencom/expr/builtins"
	"github.com/mredencom/expr/types"
)

func main() {
	fmt.Println("=== 新增内置函数演示 ===\n")

	// 字符串函数
	fmt.Println("1. 字符串函数:")

	// replace 函数
	result, err := builtins.AllBuiltins["replace"]([]types.Value{
		types.NewString("Hello World, Hello Universe"),
		types.NewString("Hello"),
		types.NewString("Hi"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   replace(\"Hello World, Hello Universe\", \"Hello\", \"Hi\") = %v\n", result)

	// substring 函数
	result, err = builtins.AllBuiltins["substring"]([]types.Value{
		types.NewString("Hello World"),
		types.NewInt(0),
		types.NewInt(5),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   substring(\"Hello World\", 0, 5) = %v\n", result)

	// indexOf 函数
	result, err = builtins.AllBuiltins["indexOf"]([]types.Value{
		types.NewString("Hello World"),
		types.NewString("World"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   indexOf(\"Hello World\", \"World\") = %v\n", result)

	fmt.Println()

	// 数学函数
	fmt.Println("2. 数学函数:")

	// ceil 函数
	result, err = builtins.AllBuiltins["ceil"]([]types.Value{
		types.NewFloat(3.14159),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ceil(3.14159) = %v\n", result)

	// floor 函数
	result, err = builtins.AllBuiltins["floor"]([]types.Value{
		types.NewFloat(3.14159),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   floor(3.14159) = %v\n", result)

	// round 函数
	result, err = builtins.AllBuiltins["round"]([]types.Value{
		types.NewFloat(3.14159),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   round(3.14159) = %v\n", result)

	// sqrt 函数
	result, err = builtins.AllBuiltins["sqrt"]([]types.Value{
		types.NewInt(16),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   sqrt(16) = %v\n", result)

	// pow 函数
	result, err = builtins.AllBuiltins["pow"]([]types.Value{
		types.NewInt(2),
		types.NewInt(3),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   pow(2, 3) = %v\n", result)

	fmt.Println()

	// 时间函数
	fmt.Println("3. 时间函数:")

	result, err = builtins.AllBuiltins["now"]([]types.Value{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   now() = %v (Unix timestamp)\n", result)

	fmt.Println()

	// 集合函数
	fmt.Println("4. 集合函数:")

	// 创建嵌套数组 [[1, 2], [3, 4]]
	elemType := types.TypeInfo{Kind: types.KindInt64, Name: "int", Size: 8}
	inner1 := types.NewSlice([]types.Value{types.NewInt(1), types.NewInt(2)}, elemType)
	inner2 := types.NewSlice([]types.Value{types.NewInt(3), types.NewInt(4)}, elemType)
	nested := types.NewSlice([]types.Value{inner1, inner2}, elemType)

	// flatten 函数
	result, err = builtins.AllBuiltins["flatten"]([]types.Value{nested})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   flatten([[1, 2], [3, 4]]) = %v\n", result)

	// groupBy 函数 (简化版本)
	users := types.NewSlice([]types.Value{
		types.NewString("Alice"),
		types.NewString("Bob"),
		types.NewString("Charlie"),
	}, types.TypeInfo{Kind: types.KindString, Name: "string", Size: -1})

	result, err = builtins.AllBuiltins["groupBy"]([]types.Value{users, types.NewBool(true)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   groupBy(users, condition) = %v\n", result)

	fmt.Println()

	// 显示所有可用的内置函数
	fmt.Println("5. 所有可用的内置函数:")
	names := builtins.ListBuiltinNames()
	for i, name := range names {
		if i > 0 && i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("%-12s", name)
	}
	fmt.Println("\n\n=== 演示完成 ===")
}
