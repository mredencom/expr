package lexer

import "fmt"

// Position represents a position in the source code
type Position struct {
	Line   int // line number (1-based)
	Column int // column number (1-based)
	Offset int // byte offset (0-based)
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}

// Valid returns true if the position is valid
func (p Position) Valid() bool {
	return p.Line > 0 && p.Column > 0
}

// Before returns true if this position is before another position
func (p Position) Before(other Position) bool {
	return p.Offset < other.Offset
}

// After returns true if this position is after another position
func (p Position) After(other Position) bool {
	return p.Offset > other.Offset
}

// NoPos represents an invalid position
var NoPos = Position{}

// Pos creates a new position
func Pos(line, column, offset int) Position {
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}
