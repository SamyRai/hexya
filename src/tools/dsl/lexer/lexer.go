package lexer

import (
	"fmt"
)

// Lexer represents the state of the lexer.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char

	line   int // current line number
	column int // current column number

	errors       []Error
	nestingLevel int
}

// NewLexer creates a new instance of Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// Lex processes the entire input string and returns tokens and errors.
func (l *Lexer) Lex() ([]Token, []Error) {
	var tokens []Token

	for l.ch != 0 {
		l.skipWhitespace()

		var token Token
		switch {
		case l.ch == '/' && l.peekChar() == '/':
			l.skipSingleLineComment()
		case l.ch == '/' && l.peekChar() == '*':
			l.skipMultiLineComment()
		case l.ch == '"', l.ch == '\'':
			token = l.readString()
			tokens = append(tokens, token)
		case l.ch == '`':
			token = l.readRawString()
			tokens = append(tokens, token)
		case l.ch == '{':
			token = l.makeToken(TokenLeftCurly, "{")
			tokens = append(tokens, token)
			l.readChar()
			l.updateNesting(1)
		case l.ch == '}':
			token = l.makeToken(TokenRightCurly, "}")
			tokens = append(tokens, token)
			l.updateNesting(-1)
			l.readChar()
		case l.ch == '(' && l.peekChar() == '*':
			methodTokens := l.handleMethodReceiver()
			tokens = append(tokens, methodTokens...)
		case l.ch == '(':
			token = l.makeToken(TokenLeftParen, "(")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == ')':
			token = l.makeToken(TokenRightParen, ")")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == '[':
			token = l.makeToken(TokenLeftBracket, "[")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == ']':
			token = l.makeToken(TokenRightBracket, "]")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == ':':
			token = l.handleColon()
			tokens = append(tokens, token)
		case l.ch == ',':
			token = l.makeToken(TokenComma, ",")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == ';':
			token = l.makeToken(TokenSemicolon, ";")
			tokens = append(tokens, token)
			l.readChar()
		case l.ch == '.':
			token = l.handleDot()
			tokens = append(tokens, token)
		case isOperatorStart(l.ch):
			token = l.readOperator()
			tokens = append(tokens, token)
		case isLetter(l.ch):
			token = l.readIdentifier()
			tokens = append(tokens, token)
		case isDigit(l.ch):
			token = l.readNumber()
			tokens = append(tokens, token)
		default:
			l.appendError(ErrorUnexpectedChar, fmt.Sprintf("Unexpected character: %q (%U)", l.ch, l.ch))
			l.readChar()
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Line: l.line, Column: l.column})
	return tokens, l.errors
}

// updateNesting safely updates the nesting level.
func (l *Lexer) updateNesting(change int) {
	l.nestingLevel += change
	if l.nestingLevel < 0 {
		l.appendError(ErrorNestingLevelReached, "Nesting level went below zero")
		l.nestingLevel = 0
	}
}

// makeToken creates a new token with the current line and column.
func (l *Lexer) makeToken(tokenType TokenType, value string) Token {
	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column,
	}
}

// appendError adds an error to the lexer's error list.
func (l *Lexer) appendError(errorType ErrorType, message string) {
	lexError := Error{
		Type:    errorType,
		Message: message,
		Line:    l.line,
		Column:  l.column,
	}
	l.errors = append(l.errors, lexError)
}
