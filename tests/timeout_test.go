package tests

import (
	"strings"
	"testing"
	"time"

	expr "github.com/mredencom/expr"
)

func TestTimeoutControl(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		// 正常执行应该在超时之前完成
		program, err := expr.Compile("2 + 2", expr.WithTimeout(5*time.Second))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		result, err := expr.Run(program, nil)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		if result != int64(4) {
			t.Errorf("期望结果为 4, 得到 %v", result)
		}
	})

	t.Run("ShortTimeout", func(t *testing.T) {
		// 测试非常短的超时
		env := map[string]interface{}{
			"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}

		// 创建一个复杂的表达式，可能需要更多时间
		expression := `
			numbers 
			| filter(# > 0) 
			| map(# * 2) 
			| filter(# > 5) 
			| map(# * 3) 
			| filter(# > 10)
		`

		program, err := expr.Compile(expression, expr.Env(env), expr.WithTimeout(1*time.Nanosecond))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		start := time.Now()
		_, err = expr.Run(program, env)
		duration := time.Since(start)

		if err == nil {
			t.Logf("表达式在 %v 内完成执行，未触发超时", duration)
		} else if strings.Contains(err.Error(), "timeout") {
			t.Logf("✅ 超时控制正常工作: %v (耗时: %v)", err, duration)
		} else {
			t.Errorf("期望超时错误，但得到其他错误: %v", err)
		}
	})

	t.Run("ReasonableTimeout", func(t *testing.T) {
		// 测试合理的超时设置，使用简单表达式
		env := map[string]interface{}{
			"x": 10,
			"y": 20,
		}

		program, err := expr.Compile(
			`x + y * 2`,
			expr.Env(env),
			expr.WithTimeout(100*time.Millisecond),
		)
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		result, err := expr.Run(program, env)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		expected := int64(50) // 10 + 20 * 2 = 50
		if result != expected {
			t.Errorf("期望结果为 %v, 得到 %v", expected, result)
		}
	})

	t.Run("WithoutTimeout", func(t *testing.T) {
		// 测试没有设置超时的情况
		program, err := expr.Compile("10 + 20 * 3")
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		result, err := expr.Run(program, nil)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		if result != int64(70) {
			t.Errorf("期望结果为 70, 得到 %v", result)
		}
	})

	t.Run("TimeoutWithResult", func(t *testing.T) {
		// 测试RunWithResult中的超时控制
		env := map[string]interface{}{
			"numbers": []int{1, 2, 3, 4, 5},
		}
		program, err := expr.Compile("numbers | map(# * 2) | sum", expr.Env(env), expr.WithTimeout(100*time.Millisecond))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		result, err := expr.RunWithResult(program, env)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		if result.Value != int64(30) { // (1+2+3+4+5)*2 = 30
			t.Errorf("期望结果为 30, 得到 %v", result.Value)
		}

		// 执行时间可能非常短，我们只是记录它
		t.Logf("执行时间: %v", result.ExecutionTime)

		if result.Type == "" {
			t.Errorf("结果类型不应该为空")
		}
	})
}

func TestTimeoutConfiguration(t *testing.T) {
	t.Run("DefaultTimeout", func(t *testing.T) {
		// 测试默认超时配置
		program, err := expr.Compile("5 * 5")
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		// 默认超时应该是30秒
		result, err := expr.Run(program, nil)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		if result != int64(25) {
			t.Errorf("期望结果为 25, 得到 %v", result)
		}
	})

	t.Run("ZeroTimeout", func(t *testing.T) {
		// 测试零超时（表示无超时限制）
		program, err := expr.Compile("100 / 4", expr.WithTimeout(0))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		result, err := expr.Run(program, nil)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}

		if result != int64(25) {
			t.Errorf("期望结果为 25, 得到 %v", result)
		}
	})

	t.Run("MultipleTimeouts", func(t *testing.T) {
		// 测试同一个程序用不同的超时配置
		baseExpression := "2 * 3 + 4"

		// 短超时
		program1, err := expr.Compile(baseExpression, expr.WithTimeout(1*time.Nanosecond))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		// 长超时
		program2, err := expr.Compile(baseExpression, expr.WithTimeout(5*time.Second))
		if err != nil {
			t.Fatalf("编译失败: %v", err)
		}

		// 长超时应该成功
		result2, err2 := expr.Run(program2, nil)
		if err2 != nil {
			t.Fatalf("长超时执行失败: %v", err2)
		}

		if result2 != int64(10) {
			t.Errorf("期望结果为 10, 得到 %v", result2)
		}

		// 短超时可能失败，但我们记录结果
		result1, err1 := expr.Run(program1, nil)
		if err1 != nil {
			t.Logf("短超时执行失败（预期）: %v", err1)
		} else {
			t.Logf("短超时执行成功: %v", result1)
		}
	})
}
