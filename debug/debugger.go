package debug

import (
	"fmt"
	"time"

	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
	"github.com/mredencom/expr/vm"
)

// Debugger provides debugging capabilities for expression execution
type Debugger struct {
	breakpoints map[int]*Breakpoint
	stats       *ExecutionStats
	enabled     bool
	stepMode    bool

	// Current execution state
	currentPC        int
	currentStack     []types.Value
	instructionCount int64

	// Event callbacks
	onBreakpoint func(*DebugContext)
	onStep       func(*DebugContext)
	onError      func(error)
}

// DebugContext provides context information during debugging
type DebugContext struct {
	PC          int                    // Program counter
	Instruction []byte                 // Current instruction bytes
	Stack       []types.Value          // Current stack state
	Variables   map[string]types.Value // Current variable values
	Source      string                 // Original source code
	Position    lexer.Position         // Source position
}

// ExecutionStats tracks execution statistics
type ExecutionStats struct {
	TotalInstructions int64
	InstructionCounts map[vm.Opcode]int64
	ExecutionTime     time.Duration
	StartTime         time.Time
	FunctionCalls     int64
	MemoryAllocations int64

	// Hot spots - most executed instructions
	HotSpots []HotSpot
}

// HotSpot represents a frequently executed instruction location
type HotSpot struct {
	PC         int
	OpCode     vm.Opcode
	Count      int64
	Percentage float64
	Source     string
	Position   lexer.Position
}

// New creates a new debugger instance
func New() *Debugger {
	return &Debugger{
		breakpoints: make(map[int]*Breakpoint),
		stats: &ExecutionStats{
			InstructionCounts: make(map[vm.Opcode]int64),
			StartTime:         time.Now(),
		},
		enabled: false,
	}
}

// Enable enables the debugger
func (d *Debugger) Enable() {
	d.enabled = true
	d.stats.StartTime = time.Now()
}

// Disable disables the debugger
func (d *Debugger) Disable() {
	d.enabled = false
}

// IsEnabled returns whether the debugger is enabled
func (d *Debugger) IsEnabled() bool {
	return d.enabled
}

// SetStepMode enables or disables step-by-step execution
func (d *Debugger) SetStepMode(enabled bool) {
	d.stepMode = enabled
}

// SetBreakpoint sets a breakpoint at the specified program counter
func (d *Debugger) SetBreakpoint(pc int) *Breakpoint {
	bp := &Breakpoint{
		PC:       pc,
		Enabled:  true,
		HitCount: 0,
	}
	d.breakpoints[pc] = bp
	return bp
}

// RemoveBreakpoint removes a breakpoint at the specified program counter
func (d *Debugger) RemoveBreakpoint(pc int) bool {
	_, exists := d.breakpoints[pc]
	if exists {
		delete(d.breakpoints, pc)
	}
	return exists
}

// GetBreakpoint returns the breakpoint at the specified PC, if any
func (d *Debugger) GetBreakpoint(pc int) (*Breakpoint, bool) {
	bp, exists := d.breakpoints[pc]
	return bp, exists
}

// ListBreakpoints returns all active breakpoints
func (d *Debugger) ListBreakpoints() []*Breakpoint {
	var bps []*Breakpoint
	for _, bp := range d.breakpoints {
		bps = append(bps, bp)
	}
	return bps
}

// OnBreakpoint sets a callback function to be called when a breakpoint is hit
func (d *Debugger) OnBreakpoint(callback func(*DebugContext)) {
	d.onBreakpoint = callback
}

// OnStep sets a callback function to be called on each step
func (d *Debugger) OnStep(callback func(*DebugContext)) {
	d.onStep = callback
}

// OnError sets a callback function to be called when an error occurs
func (d *Debugger) OnError(callback func(error)) {
	d.onError = callback
}

// ShouldBreak determines if execution should break at the current PC
func (d *Debugger) ShouldBreak(pc int) bool {
	if !d.enabled {
		return false
	}

	// Check for breakpoint
	if bp, exists := d.breakpoints[pc]; exists && bp.Enabled {
		bp.HitCount++
		return true
	}

	// Check for step mode
	return d.stepMode
}

// OnInstruction should be called before each instruction execution
func (d *Debugger) OnInstruction(pc int, instruction []byte, stack []types.Value) {
	if !d.enabled {
		return
	}

	d.currentPC = pc
	d.currentStack = stack
	d.instructionCount++

	// Update statistics
	if len(instruction) > 0 {
		opcode := vm.Opcode(instruction[0])
		d.stats.InstructionCounts[opcode]++
		d.stats.TotalInstructions++
	}

	// Create debug context
	ctx := &DebugContext{
		PC:          pc,
		Instruction: instruction,
		Stack:       stack,
		Variables:   make(map[string]types.Value), // TODO: populate from VM
	}

	// Check if we should break (but don't call ShouldBreak as it modifies state)
	shouldBreak := false
	if bp, exists := d.breakpoints[pc]; exists && bp.Enabled {
		bp.HitCount++
		shouldBreak = true
	} else if d.stepMode {
		shouldBreak = true
	}

	if shouldBreak {
		if d.onBreakpoint != nil {
			d.onBreakpoint(ctx)
		}
	}

	if d.stepMode && d.onStep != nil {
		d.onStep(ctx)
	}
}

// GetStats returns current execution statistics
func (d *Debugger) GetStats() *ExecutionStats {
	d.stats.ExecutionTime = time.Since(d.stats.StartTime)
	d.updateHotSpots()
	return d.stats
}

// ResetStats resets execution statistics
func (d *Debugger) ResetStats() {
	d.stats = &ExecutionStats{
		InstructionCounts: make(map[vm.Opcode]int64),
		StartTime:         time.Now(),
	}
	d.instructionCount = 0
}

// updateHotSpots calculates the most frequently executed instructions
func (d *Debugger) updateHotSpots() {
	d.stats.HotSpots = nil

	for opcode, count := range d.stats.InstructionCounts {
		if count > 0 {
			percentage := float64(count) / float64(d.stats.TotalInstructions) * 100
			hotspot := HotSpot{
				OpCode:     opcode,
				Count:      count,
				Percentage: percentage,
			}
			d.stats.HotSpots = append(d.stats.HotSpots, hotspot)
		}
	}

	// Sort by count (descending)
	for i := 0; i < len(d.stats.HotSpots)-1; i++ {
		for j := i + 1; j < len(d.stats.HotSpots); j++ {
			if d.stats.HotSpots[j].Count > d.stats.HotSpots[i].Count {
				d.stats.HotSpots[i], d.stats.HotSpots[j] = d.stats.HotSpots[j], d.stats.HotSpots[i]
			}
		}
	}

	// Keep only top 10
	if len(d.stats.HotSpots) > 10 {
		d.stats.HotSpots = d.stats.HotSpots[:10]
	}
}

// FormatStats formats execution statistics for display
func (d *Debugger) FormatStats() string {
	stats := d.GetStats()

	result := fmt.Sprintf("Execution Statistics:\n")
	result += fmt.Sprintf("  Total Instructions: %d\n", stats.TotalInstructions)
	result += fmt.Sprintf("  Execution Time: %v\n", stats.ExecutionTime)
	result += fmt.Sprintf("  Function Calls: %d\n", stats.FunctionCalls)
	result += fmt.Sprintf("  Memory Allocations: %d\n", stats.MemoryAllocations)

	if len(stats.HotSpots) > 0 {
		result += fmt.Sprintf("\nHot Spots (Top Instructions):\n")
		for i, hotspot := range stats.HotSpots {
			if i >= 5 { // Show top 5
				break
			}
			result += fmt.Sprintf("  %d. %s: %d times (%.1f%%)\n",
				i+1, hotspot.OpCode.String(), hotspot.Count, hotspot.Percentage)
		}
	}

	return result
}

// FormatInstructionCounts formats instruction counts for display
func (d *Debugger) FormatInstructionCounts() string {
	stats := d.GetStats()

	result := "Instruction Counts:\n"
	for opcode, count := range stats.InstructionCounts {
		if count > 0 {
			percentage := float64(count) / float64(stats.TotalInstructions) * 100
			result += fmt.Sprintf("  %-20s: %6d (%.1f%%)\n",
				opcode.String(), count, percentage)
		}
	}

	return result
}

// Trace returns a formatted trace of current execution state
func (d *Debugger) Trace() string {
	result := fmt.Sprintf("Debug Trace:\n")
	result += fmt.Sprintf("  PC: %d\n", d.currentPC)
	result += fmt.Sprintf("  Instructions Executed: %d\n", d.instructionCount)
	result += fmt.Sprintf("  Stack Size: %d\n", len(d.currentStack))

	if len(d.currentStack) > 0 {
		result += fmt.Sprintf("  Stack Top: %v\n", d.currentStack[len(d.currentStack)-1])
	}

	result += fmt.Sprintf("  Breakpoints: %d\n", len(d.breakpoints))

	return result
}
