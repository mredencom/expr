package debug

import (
	"fmt"
	"time"

	"github.com/mredencom/expr/lexer"
)

// Breakpoint represents a debugging breakpoint
type Breakpoint struct {
	PC          int            // Program counter location
	Enabled     bool           // Whether the breakpoint is active
	HitCount    int64          // Number of times this breakpoint has been hit
	Condition   string         // Optional condition for conditional breakpoints
	Source      string         // Original source code line
	Position    lexer.Position // Source code position
	CreatedAt   time.Time      // When the breakpoint was created
	Description string         // Optional description
}

// NewBreakpoint creates a new breakpoint
func NewBreakpoint(pc int) *Breakpoint {
	return &Breakpoint{
		PC:        pc,
		Enabled:   true,
		HitCount:  0,
		CreatedAt: time.Now(),
	}
}

// Enable enables the breakpoint
func (bp *Breakpoint) Enable() {
	bp.Enabled = true
}

// Disable disables the breakpoint
func (bp *Breakpoint) Disable() {
	bp.Enabled = false
}

// SetCondition sets a condition for the breakpoint
func (bp *Breakpoint) SetCondition(condition string) {
	bp.Condition = condition
}

// SetDescription sets a description for the breakpoint
func (bp *Breakpoint) SetDescription(description string) {
	bp.Description = description
}

// ShouldBreak determines if the breakpoint should trigger
func (bp *Breakpoint) ShouldBreak() bool {
	if !bp.Enabled {
		return false
	}

	// For now, always break if enabled
	// TODO: Implement condition evaluation
	return true
}

// Hit records a hit on this breakpoint
func (bp *Breakpoint) Hit() {
	bp.HitCount++
}

// String returns a string representation of the breakpoint
func (bp *Breakpoint) String() string {
	status := "disabled"
	if bp.Enabled {
		status = "enabled"
	}

	result := fmt.Sprintf("Breakpoint %d: %s (hits: %d)", bp.PC, status, bp.HitCount)

	if bp.Description != "" {
		result += fmt.Sprintf(" - %s", bp.Description)
	}

	if bp.Condition != "" {
		result += fmt.Sprintf(" [condition: %s]", bp.Condition)
	}

	return result
}
