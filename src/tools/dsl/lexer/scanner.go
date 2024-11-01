package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// skipWhitespace advances over any whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch != 0 && isWhitespace(l.ch) {
		l.readChar()
	}
}

// skipSingleLineComment skips over single-line comments.
func (l *Lexer) skipSingleLineComment() {
	for l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
}

// skipMultiLineComment skips over multi-line comments.
func (l *Lexer) skipMultiLineComment() {
	l.readChar() // Consume '*'
	for {
		if l.ch == 0 {
			l.appendError(ErrorUnterminatedComment, "Unterminated multi-line comment")
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		l.readChar()
	}
}

// peekChar looks at the next character without advancing the lexer
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return ch
}

// peekSecondChar looks at the character after the next character
func (l *Lexer) peekSecondChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	_, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
	nextPos := l.readPosition + width
	if nextPos >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[nextPos:])
	return ch
}

// readNumber reads a number and returns a token.
func (l *Lexer) readNumber() Token {
	startPosition := l.position
	startLine := l.line
	startColumn := l.column
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		isFloat = true
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	value := l.input[startPosition:l.position]
	if isFloat {
		return Token{
			Type:   TokenFloat,
			Value:  value,
			Line:   startLine,
			Column: startColumn,
		}
	}
	return Token{
		Type:   TokenNumber,
		Value:  value,
		Line:   startLine,
		Column: startColumn,
	}
}

// readString reads a quoted string and returns a token.
func (l *Lexer) readString() Token {
	startLine, startColumn := l.line, l.column
	quote := l.ch
	l.readChar() // Consume the opening quote
	var value strings.Builder

	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				value.WriteRune('\n')
			case 't':
				value.WriteRune('\t')
			case 'r':
				value.WriteRune('\r')
			case '\\':
				value.WriteRune('\\')
			case '"':
				value.WriteRune('"')
			case '\'':
				value.WriteRune('\'')
			case 'u':
				unicodeChar := l.readUnicodeEscape()
				value.WriteRune(unicodeChar)
			case 'x':
				hexChar := l.readHexEscape(2)
				value.WriteRune(hexChar)
			default:
				// Unrecognized escape sequence
				l.appendError(ErrorInvalidEscapeSequence, fmt.Sprintf("Unrecognized escape sequence: \\%c", l.ch))
				value.WriteRune(l.ch)
			}
		} else {
			value.WriteRune(l.ch)
		}
		l.readChar()
	}

	if l.ch != quote {
		l.appendError(ErrorUnterminatedString, "Unterminated string literal")
	} else {
		l.readChar() // Consume the closing quote
	}

	return Token{
		Type:   TokenString,
		Value:  value.String(),
		Line:   startLine,
		Column: startColumn,
	}
}

// readRawString reads a raw string (backtick-quoted) and returns a token.
func (l *Lexer) readRawString() Token {
	startLine, startColumn := l.line, l.column
	l.readChar() // Consume starting backtick
	var value strings.Builder

	for l.ch != '`' && l.ch != 0 {
		value.WriteRune(l.ch)
		l.readChar()
	}

	if l.ch == '`' {
		l.readChar() // Consume ending backtick
	} else {
		l.appendError(ErrorUnterminatedString, "Unterminated raw string literal")
	}

	return Token{
		Type:   TokenString,
		Value:  value.String(),
		Line:   startLine,
		Column: startColumn,
	}
}

// handleColon handles colon tokens, including ':='.
func (l *Lexer) handleColon() Token {
	startLine, startColumn := l.line, l.column
	if l.peekChar() == '=' {
		l.readChar()
		l.readChar()
		return Token{
			Type:   TokenDefine,
			Value:  ":=",
			Line:   startLine,
			Column: startColumn,
		}
	}
	l.readChar()
	return Token{
		Type:   TokenColon,
		Value:  ":",
		Line:   startLine,
		Column: startColumn,
	}
}

// handleDot handles dot tokens, including ellipsis '...'.
func (l *Lexer) handleDot() Token {
	startLine, startColumn := l.line, l.column
	if l.peekChar() == '.' && l.peekSecondChar() == '.' {
		l.readChar() // Consume first dot
		l.readChar() // Consume second dot
		l.readChar() // Consume third dot
		return Token{
			Type:   TokenEllipsis,
			Value:  "...",
			Line:   startLine,
			Column: startColumn,
		}
	}
	l.readChar()
	return Token{
		Type:   TokenDot,
		Value:  ".",
		Line:   startLine,
		Column: startColumn,
	}
}

// handleMethodReceiver handles method receiver syntax (e.g., (h *User)).
func (l *Lexer) handleMethodReceiver() []Token {
	// Implement this function based on your language's syntax.
	// For now, we'll skip it.
	return nil
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		var width int
		l.ch, width = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += width

		if l.ch == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column += 1 // Increment by rune width, not just 1
		}
	}
}
