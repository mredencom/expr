package lexer

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Literals
	NUMBER // 123, 123.45
	STRING // "abc", 'abc'
	BOOL   // true, false
	NULL   // null

	// Identifiers
	IDENT // variable names, function names

	// Operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %
	POW // **

	// Comparison operators
	EQ // ==
	NE // !=
	LT // <
	LE // <=
	GT // >
	GE // >=

	// Logical operators
	AND // &&
	OR  // ||
	NOT // !

	// Assignment
	ASSIGN // =

	// Bitwise operators
	BIT_AND // &
	BIT_OR  // |
	BIT_XOR // ^
	BIT_NOT // ~
	SHL     // <<
	SHR     // >>

	// Punctuation
	LPAREN      // (
	RPAREN      // )
	LBRACKET    // [
	RBRACKET    // ]
	LBRACE      // {
	RBRACE      // }
	COMMA       // ,
	DOT         // .
	SEMICOLON   // ;
	COLON       // :
	QUESTION    // ?
	ARROW       // => (for lambda expressions)
	PIPE        // | (for pipeline operations)
	WILDCARD    // * (for wildcard operations, different from MUL)
	PLACEHOLDER // # (for pipeline placeholder)

	// Null safety operators
	QUESTION_DOT    // ?. (optional chaining)
	NULL_COALESCING // ?? (null coalescing)

	// Destructuring operators
	SPREAD // ... (spread/rest operator)

	// Keywords
	IF
	ELSE
	IN
	MATCHES
	CONTAINS
	STARTS_WITH
	ENDS_WITH
	IMPORT
	AS
	FROM
)

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

// String returns a string representation of the token
func (t Token) String() string {
	if t.Value != "" {
		return fmt.Sprintf("%s(%s)", t.Type, t.Value)
	}
	return t.Type.String()
}

// IsLiteral returns true if the token is a literal
func (t Token) IsLiteral() bool {
	return t.Type >= NUMBER && t.Type <= NULL
}

// IsOperator returns true if the token is an operator
func (t Token) IsOperator() bool {
	return t.Type >= ADD && t.Type <= SHR
}

// IsComparison returns true if the token is a comparison operator
func (t Token) IsComparison() bool {
	return t.Type >= EQ && t.Type <= GE
}

// IsLogical returns true if the token is a logical operator
func (t Token) IsLogical() bool {
	return t.Type >= AND && t.Type <= NOT
}

// IsKeyword returns true if the token is a keyword
func (t Token) IsKeyword() bool {
	return t.Type >= IF && t.Type <= ENDS_WITH
}

// String returns the string representation of TokenType
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case BOOL:
		return "BOOL"
	case NULL:
		return "NULL"
	case IDENT:
		return "IDENT"
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case MOD:
		return "%"
	case POW:
		return "**"
	case EQ:
		return "=="
	case NE:
		return "!="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case ASSIGN:
		return "="
	case BIT_AND:
		return "&"
	case BIT_OR:
		return "|"
	case BIT_XOR:
		return "^"
	case BIT_NOT:
		return "~"
	case SHL:
		return "<<"
	case SHR:
		return ">>"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case COMMA:
		return ","
	case DOT:
		return "."
	case SEMICOLON:
		return ";"
	case COLON:
		return ":"
	case QUESTION:
		return "?"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case IN:
		return "in"
	case MATCHES:
		return "matches"
	case CONTAINS:
		return "contains"
	case STARTS_WITH:
		return "startsWith"
	case ENDS_WITH:
		return "endsWith"
	case IMPORT:
		return "import"
	case AS:
		return "as"
	case FROM:
		return "from"
	case ARROW:
		return "=>"
	case PIPE:
		return "|"
	case WILDCARD:
		return "*"
	case PLACEHOLDER:
		return "#"
	case QUESTION_DOT:
		return "?."
	case NULL_COALESCING:
		return "??"
	case SPREAD:
		return "..."
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}

// Keywords maps keyword strings to their token types
var Keywords = map[string]TokenType{
	"true":       BOOL,
	"false":      BOOL,
	"null":       NULL,
	"if":         IF,
	"else":       ELSE,
	"in":         IN,
	"matches":    MATCHES,
	"contains":   CONTAINS,
	"startsWith": STARTS_WITH,
	"endsWith":   ENDS_WITH,
	"import":     IMPORT,
	"as":         AS,
	"from":       FROM,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
