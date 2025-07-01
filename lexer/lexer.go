package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Lexer performs lexical analysis of input text
type Lexer struct {
	input     string
	position  int  // current position in input (points to current char)
	readPos   int  // current reading position in input (after current char)
	char      rune // current char under examination
	line      int  // current line number (1-based)
	column    int  // current column number (1-based)
	lineStart int  // position where current line starts
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:     input,
		line:      1,
		column:    1,
		lineStart: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0 // ASCII NUL character represents EOF
	} else {
		var size int
		l.char, size = utf8.DecodeRuneInString(l.input[l.readPos:])
		if l.char == utf8.RuneError && size == 1 {
			l.char = rune(l.input[l.readPos])
		}
	}

	l.position = l.readPos
	l.readPos += utf8.RuneLen(l.char)

	if l.char == '\n' {
		l.line++
		l.lineStart = l.readPos
		l.column = 1
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}

	char, _ := utf8.DecodeRuneInString(l.input[l.readPos:])
	return char
}

// peekCharN returns the character n positions ahead without advancing the position
func (l *Lexer) peekCharN(n int) rune {
	pos := l.readPos
	for i := 0; i < n && pos < len(l.input); i++ {
		_, size := utf8.DecodeRuneInString(l.input[pos:])
		pos += size
	}

	if pos >= len(l.input) {
		return 0
	}

	char, _ := utf8.DecodeRuneInString(l.input[pos:])
	return char
}

// currentPosition returns the current position
func (l *Lexer) currentPosition() Position {
	return Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}
}

// NextToken scans the input and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Position = l.currentPosition()

	switch l.char {
	case '+':
		tok = Token{Type: ADD, Value: "+", Position: tok.Position}
	case '-':
		tok = Token{Type: SUB, Value: "-", Position: tok.Position}
	case '*':
		if l.peekChar() == '*' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: POW, Value: "**", Position: pos}
		} else {
			// Check if this could be a wildcard in member access context
			// Look for patterns like "user.*" or "data.*.field"
			if l.isWildcardContext() {
				tok = Token{Type: WILDCARD, Value: "*", Position: tok.Position}
			} else {
				tok = Token{Type: MUL, Value: "*", Position: tok.Position}
			}
		}
	case '/':
		tok = Token{Type: DIV, Value: "/", Position: tok.Position}
	case '%':
		tok = Token{Type: MOD, Value: "%", Position: tok.Position}
	case '=':
		if l.peekChar() == '=' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: EQ, Value: "==", Position: pos}
		} else if l.peekChar() == '>' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: ARROW, Value: "=>", Position: pos}
		} else {
			tok = Token{Type: ASSIGN, Value: "=", Position: tok.Position}
		}
	case '!':
		if l.peekChar() == '=' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: NE, Value: "!=", Position: pos}
		} else {
			tok = Token{Type: NOT, Value: "!", Position: tok.Position}
		}
	case '<':
		if l.peekChar() == '=' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: LE, Value: "<=", Position: pos}
		} else if l.peekChar() == '<' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: SHL, Value: "<<", Position: pos}
		} else {
			tok = Token{Type: LT, Value: "<", Position: tok.Position}
		}
	case '>':
		if l.peekChar() == '=' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: GE, Value: ">=", Position: pos}
		} else if l.peekChar() == '>' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: SHR, Value: ">>", Position: pos}
		} else {
			tok = Token{Type: GT, Value: ">", Position: tok.Position}
		}
	case '&':
		if l.peekChar() == '&' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: AND, Value: "&&", Position: pos}
		} else {
			tok = Token{Type: BIT_AND, Value: "&", Position: tok.Position}
		}
	case '|':
		if l.peekChar() == '|' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: OR, Value: "||", Position: pos}
		} else {
			tok = Token{Type: BIT_OR, Value: "|", Position: tok.Position}
		}
	case '^':
		tok = Token{Type: BIT_XOR, Value: "^", Position: tok.Position}
	case '~':
		tok = Token{Type: BIT_NOT, Value: "~", Position: tok.Position}
	case '(':
		tok = Token{Type: LPAREN, Value: "(", Position: tok.Position}
	case ')':
		tok = Token{Type: RPAREN, Value: ")", Position: tok.Position}
	case '[':
		tok = Token{Type: LBRACKET, Value: "[", Position: tok.Position}
	case ']':
		tok = Token{Type: RBRACKET, Value: "]", Position: tok.Position}
	case '{':
		tok = Token{Type: LBRACE, Value: "{", Position: tok.Position}
	case '}':
		tok = Token{Type: RBRACE, Value: "}", Position: tok.Position}
	case ',':
		tok = Token{Type: COMMA, Value: ",", Position: tok.Position}
	case '.':
		if l.peekChar() == '.' && l.peekCharN(2) == '.' {
			pos := tok.Position
			l.readChar() // consume second '.'
			l.readChar() // consume third '.'
			tok = Token{Type: SPREAD, Value: "...", Position: pos}
		} else {
			tok = Token{Type: DOT, Value: ".", Position: tok.Position}
		}
	case ';':
		tok = Token{Type: SEMICOLON, Value: ";", Position: tok.Position}
	case ':':
		tok = Token{Type: COLON, Value: ":", Position: tok.Position}
	case '?':
		if l.peekChar() == '.' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: QUESTION_DOT, Value: "?.", Position: pos}
		} else if l.peekChar() == '?' {
			pos := tok.Position
			l.readChar()
			tok = Token{Type: NULL_COALESCING, Value: "??", Position: pos}
		} else {
			tok = Token{Type: QUESTION, Value: "?", Position: tok.Position}
		}
	case '#':
		tok = Token{Type: PLACEHOLDER, Value: "#", Position: tok.Position}
	case '"':
		tok.Type = STRING
		tok.Value = l.readString('"')
		return tok // Don't advance char, readString already did
	case '\'':
		tok.Type = STRING
		tok.Value = l.readString('\'')
		return tok // Don't advance char, readString already did
	case 0:
		tok = Token{Type: EOF, Value: "", Position: tok.Position}
	default:
		if isLetter(l.char) {
			tok.Value = l.readIdentifier()
			tok.Type = LookupIdent(tok.Value)
			return tok // Don't advance char, readIdentifier already did
		} else if isDigit(l.char) {
			tok.Type = NUMBER
			tok.Value = l.readNumber()
			return tok // Don't advance char, readNumber already did
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.char), Position: tok.Position}
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) || isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() string {
	position := l.position

	// Read integer part
	for isDigit(l.char) {
		l.readChar()
	}

	// Check for decimal point
	if l.char == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.char) {
			l.readChar()
		}
	}

	// Check for scientific notation
	if l.char == 'e' || l.char == 'E' {
		l.readChar()
		if l.char == '+' || l.char == '-' {
			l.readChar()
		}
		if !isDigit(l.char) {
			// Invalid scientific notation, backtrack
			return l.input[position : l.position-1]
		}
		for isDigit(l.char) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString(delimiter rune) string {
	position := l.position + utf8.RuneLen(delimiter)
	l.readChar() // skip opening delimiter

	for l.char != delimiter && l.char != 0 {
		if l.char == '\\' {
			l.readChar() // skip escape character
			if l.char != 0 {
				l.readChar() // skip escaped character
			}
		} else {
			l.readChar()
		}
	}

	result := l.input[position:l.position]

	// Skip the closing delimiter
	if l.char == delimiter {
		l.readChar()
	}

	return l.unescapeString(result)
}

// unescapeString handles escape sequences in strings
func (l *Lexer) unescapeString(s string) string {
	if len(s) == 0 {
		return s
	}

	result := make([]rune, 0, len(s))
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) {
			switch runes[i+1] {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			case '\'':
				result = append(result, '\'')
			default:
				result = append(result, runes[i+1])
			}
			i++ // skip next character
		} else {
			result = append(result, runes[i])
		}
	}

	return string(result)
}

// isLetter checks if a character is a letter or underscore
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isDigit checks if a character is a digit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// Reset resets the lexer with new input
func (l *Lexer) Reset(input string) {
	l.input = input
	l.position = 0
	l.readPos = 0
	l.line = 1
	l.column = 1
	l.lineStart = 0
	l.readChar()
}

// isWildcardContext determines if * should be treated as wildcard or multiplication
func (l *Lexer) isWildcardContext() bool {
	// Look backward to see if we're in a member access context
	// This is a simple heuristic - in practice, we might need more sophisticated parsing
	if l.position > 0 {
		// Check if preceded by a dot (e.g., "user.*")
		prevPos := l.position - 1
		for prevPos >= 0 && unicode.IsSpace(rune(l.input[prevPos])) {
			prevPos--
		}
		if prevPos >= 0 && l.input[prevPos] == '.' {
			return true
		}
	}

	// Look forward to see if followed by a dot or identifier (e.g., "*.field")
	nextChar := l.peekChar()
	if nextChar == '.' || isLetter(nextChar) {
		return true
	}

	return false
}

// Error creates an error token
func (l *Lexer) Error(msg string) Token {
	return Token{
		Type:     ILLEGAL,
		Value:    fmt.Sprintf("lexer error: %s", msg),
		Position: l.currentPosition(),
	}
}
