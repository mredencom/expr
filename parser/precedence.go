package parser

import "github.com/mredencom/expr/lexer"

// Precedence represents operator precedence levels
type Precedence int

const (
	LOWEST            Precedence = iota
	LAMBDA                       // => (lambda arrows have lowest precedence)
	PIPE                         // | (pipeline operations)
	NULL_COALESCING              // ?? (null coalescing)
	TERNARY                      // ? :
	LOGICAL_OR                   // ||
	LOGICAL_AND                  // &&
	EQUALS                       // ==, !=
	LESSGREATER                  // > or <
	SUM                          // +, -
	PRODUCT                      // *, /, %
	POWER                        // **
	PREFIX                       // -X, !X
	CALL                         // myFunction(X)
	INDEX                        // array[index], obj.property
	OPTIONAL_CHAINING            // ?. (optional chaining, highest precedence for member access)
)

// precedences maps token types to their precedence levels
var precedences = map[lexer.TokenType]Precedence{
	// Lambda and pipeline
	lexer.ARROW: LAMBDA,
	lexer.PIPE:  PIPE,

	// Ternary operator
	lexer.QUESTION: TERNARY,

	// Logical operators
	lexer.OR:  LOGICAL_OR,
	lexer.AND: LOGICAL_AND,

	// Equality operators
	lexer.EQ: EQUALS,
	lexer.NE: EQUALS,

	// Comparison operators
	lexer.LT: LESSGREATER,
	lexer.LE: LESSGREATER,
	lexer.GT: LESSGREATER,
	lexer.GE: LESSGREATER,

	// String operators
	lexer.IN:          EQUALS,
	lexer.MATCHES:     EQUALS,
	lexer.CONTAINS:    EQUALS,
	lexer.STARTS_WITH: EQUALS,
	lexer.ENDS_WITH:   EQUALS,

	// Arithmetic operators
	lexer.ADD: SUM,
	lexer.SUB: SUM,
	lexer.MUL: PRODUCT,
	lexer.DIV: PRODUCT,
	lexer.MOD: PRODUCT,
	lexer.POW: POWER,

	// Bitwise operators
	lexer.BIT_OR:  SUM,
	lexer.BIT_XOR: SUM,
	lexer.BIT_AND: PRODUCT,
	lexer.SHL:     PRODUCT,
	lexer.SHR:     PRODUCT,

	// Call and index
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: INDEX,
	lexer.DOT:      INDEX,

	// Null safety operators
	lexer.NULL_COALESCING: NULL_COALESCING,
	lexer.QUESTION_DOT:    OPTIONAL_CHAINING,
}

// GetPrecedence returns the precedence of a token type
func GetPrecedence(tokenType lexer.TokenType) Precedence {
	if p, ok := precedences[tokenType]; ok {
		return p
	}
	return LOWEST
}

// IsRightAssociative returns true if the operator is right-associative
func IsRightAssociative(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.POW: // ** is right-associative: 2**3**2 = 2**(3**2) = 2**9 = 512
		return true
	default:
		return false
	}
}

// IsComparison returns true if the token is a comparison operator
func IsComparison(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.EQ, lexer.NE, lexer.LT, lexer.LE, lexer.GT, lexer.GE:
		return true
	case lexer.IN, lexer.MATCHES, lexer.CONTAINS, lexer.STARTS_WITH, lexer.ENDS_WITH:
		return true
	default:
		return false
	}
}

// IsLogical returns true if the token is a logical operator
func IsLogical(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.AND, lexer.OR:
		return true
	default:
		return false
	}
}

// IsArithmetic returns true if the token is an arithmetic operator
func IsArithmetic(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.ADD, lexer.SUB, lexer.MUL, lexer.DIV, lexer.MOD, lexer.POW:
		return true
	default:
		return false
	}
}

// IsBitwise returns true if the token is a bitwise operator
func IsBitwise(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.BIT_AND, lexer.BIT_OR, lexer.BIT_XOR, lexer.BIT_NOT, lexer.SHL, lexer.SHR:
		return true
	default:
		return false
	}
}
