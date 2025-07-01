package tests

import (
	"testing"

	"github.com/mredencom/expr/debug"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

func TestDebugger(t *testing.T) {
	t.Run("CreateDebugger", func(t *testing.T) {
		debugger := debug.New()
		if debugger == nil {
			t.Fatal("调试器创建失败")
		}

		if debugger.IsEnabled() {
			t.Error("调试器默认应该是禁用状态")
		}
	})

	t.Run("EnableDisableDebugger", func(t *testing.T) {
		debugger := debug.New()

		debugger.Enable()
		if !debugger.IsEnabled() {
			t.Error("调试器启用失败")
		}

		debugger.Disable()
		if debugger.IsEnabled() {
			t.Error("调试器禁用失败")
		}
	})

	t.Run("BreakpointManagement", func(t *testing.T) {
		debugger := debug.New()

		// 设置断点
		bp := debugger.SetBreakpoint(10)
		if bp == nil {
			t.Fatal("断点创建失败")
		}

		if bp.PC != 10 {
			t.Errorf("期望断点PC为10，得到%d", bp.PC)
		}

		if !bp.Enabled {
			t.Error("断点应该默认启用")
		}

		// 获取断点
		retrieved, exists := debugger.GetBreakpoint(10)
		if !exists {
			t.Error("应该能找到设置的断点")
		}

		if retrieved.PC != 10 {
			t.Errorf("期望断点PC为10，得到%d", retrieved.PC)
		}

		// 列出断点
		breakpoints := debugger.ListBreakpoints()
		if len(breakpoints) != 1 {
			t.Errorf("期望1个断点，得到%d个", len(breakpoints))
		}

		// 删除断点
		removed := debugger.RemoveBreakpoint(10)
		if !removed {
			t.Error("断点删除失败")
		}

		_, exists = debugger.GetBreakpoint(10)
		if exists {
			t.Error("断点删除后不应该存在")
		}
	})

	t.Run("ShouldBreakLogic", func(t *testing.T) {
		debugger := debug.New()

		// 禁用状态下不应该中断
		if debugger.ShouldBreak(10) {
			t.Error("禁用状态下不应该中断")
		}

		// 启用但没有断点，不应该中断
		debugger.Enable()
		if debugger.ShouldBreak(10) {
			t.Error("没有断点时不应该中断")
		}

		// 设置断点后应该中断
		debugger.SetBreakpoint(10)
		if !debugger.ShouldBreak(10) {
			t.Error("有断点时应该中断")
		}

		// 步进模式下应该中断
		debugger.SetStepMode(true)
		if !debugger.ShouldBreak(20) {
			t.Error("步进模式下应该中断")
		}
	})

	t.Run("InstructionStatistics", func(t *testing.T) {
		debugger := debug.New()
		debugger.Enable()

		// 模拟一些指令执行
		instructions := [][]byte{
			{byte(vm.OpConstant), 0, 0},
			{byte(vm.OpConstant), 0, 1},
			{byte(vm.OpAdd)},
		}

		stack := []types.Value{}

		for i, inst := range instructions {
			debugger.OnInstruction(i, inst, stack)
		}

		stats := debugger.GetStats()
		if stats.TotalInstructions != 3 {
			t.Errorf("期望总指令数为3，得到%d", stats.TotalInstructions)
		}

		if stats.InstructionCounts[vm.OpConstant] != 2 {
			t.Errorf("期望OpConstant指令2次，得到%d", stats.InstructionCounts[vm.OpConstant])
		}

		if stats.InstructionCounts[vm.OpAdd] != 1 {
			t.Errorf("期望OpAdd指令1次，得到%d", stats.InstructionCounts[vm.OpAdd])
		}
	})

	t.Run("StatsFormatting", func(t *testing.T) {
		debugger := debug.New()
		debugger.Enable()

		// 执行一些指令
		for i := 0; i < 10; i++ {
			debugger.OnInstruction(i, []byte{byte(vm.OpAdd)}, []types.Value{})
		}

		statsText := debugger.FormatStats()
		if statsText == "" {
			t.Error("统计信息格式化不应该为空")
		}

		t.Logf("统计信息:\n%s", statsText)

		instructionStats := debugger.FormatInstructionCounts()
		if instructionStats == "" {
			t.Error("指令计数格式化不应该为空")
		}

		t.Logf("指令计数:\n%s", instructionStats)
	})

	t.Run("CallbackExecution", func(t *testing.T) {
		debugger := debug.New()
		debugger.Enable()

		var breakpointCalled bool
		var stepCalled bool
		var breakpointPC int
		var stepPC int

		debugger.OnBreakpoint(func(ctx *debug.DebugContext) {
			breakpointCalled = true
			breakpointPC = ctx.PC
		})

		debugger.OnStep(func(ctx *debug.DebugContext) {
			stepCalled = true
			stepPC = ctx.PC
		})

		// 设置断点并触发
		debugger.SetBreakpoint(5)
		debugger.OnInstruction(5, []byte{byte(vm.OpAdd)}, []types.Value{})

		if !breakpointCalled {
			t.Error("断点回调应该被调用")
		}

		if breakpointPC != 5 {
			t.Errorf("期望断点PC为5，得到%d", breakpointPC)
		}

		// 启用步进模式
		debugger.SetStepMode(true)
		debugger.OnInstruction(6, []byte{byte(vm.OpSub)}, []types.Value{})

		if !stepCalled {
			t.Error("步进回调应该被调用")
		}

		if stepPC != 6 {
			t.Errorf("期望步进PC为6，得到%d", stepPC)
		}
	})

	t.Run("ResetStats", func(t *testing.T) {
		debugger := debug.New()
		debugger.Enable()

		// 执行一些指令
		debugger.OnInstruction(0, []byte{byte(vm.OpAdd)}, []types.Value{})
		debugger.OnInstruction(1, []byte{byte(vm.OpSub)}, []types.Value{})

		stats := debugger.GetStats()
		if stats.TotalInstructions != 2 {
			t.Errorf("期望2条指令，得到%d", stats.TotalInstructions)
		}

		// 重置统计
		debugger.ResetStats()
		stats = debugger.GetStats()
		if stats.TotalInstructions != 0 {
			t.Errorf("重置后期望0条指令，得到%d", stats.TotalInstructions)
		}
	})

	t.Run("Trace", func(t *testing.T) {
		debugger := debug.New()
		debugger.Enable()

		// 执行一些指令来构建跟踪状态
		stack := []types.Value{types.NewInt(42)}
		debugger.OnInstruction(10, []byte{byte(vm.OpConstant)}, stack)

		trace := debugger.Trace()
		if trace == "" {
			t.Error("跟踪信息不应该为空")
		}

		t.Logf("跟踪信息:\n%s", trace)
	})
}

func TestBreakpoint(t *testing.T) {
	t.Run("CreateBreakpoint", func(t *testing.T) {
		bp := debug.NewBreakpoint(15)
		if bp == nil {
			t.Fatal("断点创建失败")
		}

		if bp.PC != 15 {
			t.Errorf("期望PC为15，得到%d", bp.PC)
		}

		if !bp.Enabled {
			t.Error("断点应该默认启用")
		}

		if bp.HitCount != 0 {
			t.Errorf("期望命中次数为0，得到%d", bp.HitCount)
		}
	})

	t.Run("BreakpointOperations", func(t *testing.T) {
		bp := debug.NewBreakpoint(20)

		// 测试禁用/启用
		bp.Disable()
		if bp.Enabled {
			t.Error("断点禁用失败")
		}

		if bp.ShouldBreak() {
			t.Error("禁用的断点不应该中断")
		}

		bp.Enable()
		if !bp.Enabled {
			t.Error("断点启用失败")
		}

		if !bp.ShouldBreak() {
			t.Error("启用的断点应该中断")
		}

		// 测试命中计数
		bp.Hit()
		if bp.HitCount != 1 {
			t.Errorf("期望命中次数为1，得到%d", bp.HitCount)
		}

		bp.Hit()
		if bp.HitCount != 2 {
			t.Errorf("期望命中次数为2，得到%d", bp.HitCount)
		}
	})

	t.Run("BreakpointDescription", func(t *testing.T) {
		bp := debug.NewBreakpoint(25)

		bp.SetDescription("测试断点")
		bp.SetCondition("x > 10")

		str := bp.String()
		if str == "" {
			t.Error("断点字符串表示不应该为空")
		}

		t.Logf("断点信息: %s", str)
	})
}
