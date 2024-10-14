package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType represents the type of token.
type TokenType string

// Token types.
const (
	// Keywords
	TokenModel      TokenType = "MODEL"
	TokenField      TokenType = "FIELD"
	TokenMethod     TokenType = "METHOD"
	TokenFunc       TokenType = "FUNC"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenNumber     TokenType = "NUMBER"

	// Symbols
	TokenLeftParen    TokenType = "LEFT_PAREN"
	TokenRightParen   TokenType = "RIGHT_PAREN"
	TokenLeftCurly    TokenType = "LEFT_CURLY_BRACE"
	TokenRightCurly   TokenType = "RIGHT_CURLY_BRACE"
	TokenLeftBracket  TokenType = "LEFT_BRACKET"
	TokenRightBracket TokenType = "RIGHT_BRACKET"
	TokenComma        TokenType = "COMMA"
	TokenColon        TokenType = "COLON"
	TokenDot          TokenType = "DOT"
	TokenSemicolon    TokenType = "SEMICOLON"
	TokenEllipsis     TokenType = "ELLIPSIS"
	TokenOperator     TokenType = "OPERATOR"
	TokenSpecial      TokenType = "SPECIAL"

	// Others
	TokenEOF   TokenType = "EOF"
	TokenError TokenType = "ERROR"
)

// DSL Keywords
const (
	NewModelKeyword      = "NewModel"
	AddFieldsKeyword     = "AddFields"
	AddMethodKeyword     = "AddMethod"
	SetMixinKeyword      = "SetMixin"
	AddValidationKeyword = "AddValidation"
	SetRouteKeyword      = "SetRoute"

	// Method Keywords
	SearchKeyword           = "Search"
	SearchByNameKeyword     = "SearchByName"
	CreateKeyword           = "Create"
	NewKeyword              = "New"
	WriteKeyword            = "Write"
	CopyKeyword             = "Copy"
	CopyDataKeyword         = "CopyData"
	CartesianProductKeyword = "CartesianProduct"
	SortedKeyword           = "Sorted"
	FilteredKeyword         = "Filtered"
	AggregatesKeyword       = "Aggregates"
	FirstKeyword            = "First"
	AllKeyword              = "All"
	DefaultGetKeyword       = "DefaultGet"
)

// Token represents a lexed token with type and value.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer for tokenizing input.
type Lexer struct {
	input        string
	position     int  // current position in input (in bytes)
	readPosition int  // current reading position in input (in bytes)
	ch           rune // current character under examination
}

// NewLexer initializes a new lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character and advances the positions.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
	} else {
		var width int
		l.ch, width = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += width
	}
}

// Lex tokenizes the input and returns a slice of tokens.
func (l *Lexer) Lex() []Token {
	var tokens []Token

	for l.ch != 0 {
		l.skipWhitespace()

		var tok Token

		switch {
		case l.ch == '"', l.ch == '\'':
			tok = l.readString()
		case l.ch == '(':
			tok = Token{Type: TokenLeftParen, Value: string(l.ch)}
			l.readChar()
		case l.ch == ')':
			tok = Token{Type: TokenRightParen, Value: string(l.ch)}
			l.readChar()
		case l.ch == '{':
			tok = Token{Type: TokenLeftCurly, Value: string(l.ch)}
			l.readChar()
		case l.ch == '}':
			tok = Token{Type: TokenRightCurly, Value: string(l.ch)}
			l.readChar()
		case l.ch == '[':
			tok = Token{Type: TokenLeftBracket, Value: string(l.ch)}
			l.readChar()
		case l.ch == ']':
			tok = Token{Type: TokenRightBracket, Value: string(l.ch)}
			l.readChar()
		case l.ch == ',':
			tok = Token{Type: TokenComma, Value: string(l.ch)}
			l.readChar()
		case l.ch == ':':
			tok = Token{Type: TokenColon, Value: string(l.ch)}
			l.readChar()
		case l.ch == ';':
			tok = Token{Type: TokenSemicolon, Value: string(l.ch)}
			l.readChar()
		case l.ch == '.':
			if l.peekChar() == '.' && l.peekSecondChar() == '.' {
				tok = Token{Type: TokenEllipsis, Value: "..."}
				l.readChar() // Read first dot
				l.readChar() // Read second dot
				l.readChar() // Read third dot
			} else {
				tok = Token{Type: TokenDot, Value: string(l.ch)}
				l.readChar()
			}
		case isOperator(l.ch):
			tok = Token{Type: TokenOperator, Value: l.readOperator()}
		case isLetter(l.ch):
			ident := l.readIdentifier()
			tokType := lookupIdent(ident)
			tok = Token{Type: tokType, Value: ident}
		case isDigit(l.ch):
			tok = l.readNumber()
		default:
			tok = Token{Type: TokenError, Value: string(l.ch)}
			l.readChar()
		}

		tokens = append(tokens, tok)
	}

	tokens = append(tokens, Token{Type: TokenEOF})
	return tokens
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch != 0 && isWhitespace(l.ch) {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword.
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isIdentifierChar(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber reads a numeric literal.
func (l *Lexer) readNumber() Token {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar() // Consume dot
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return Token{Type: TokenNumber, Value: l.input[start:l.position]}
}

// readString reads a string literal.
func (l *Lexer) readString() Token {
	quote := l.ch
	l.readChar()
	start := l.position
	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // Skip escape character
		}
		l.readChar()
	}
	str := l.input[start:l.position]
	l.readChar() // Skip the closing quote
	return Token{Type: TokenString, Value: str}
}

// readOperator reads an operator token.
func (l *Lexer) readOperator() string {
	start := l.position
	for isOperator(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// peekChar peeks at the next character without consuming it.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return ch
}

// peekSecondChar peeks at the character after the next one.
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

// Helper functions
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isIdentifierChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch)
}

func isOperator(ch rune) bool {
	return strings.ContainsRune("+-*/%=&|<>!", ch)
}

// Keywords map for lookup.
var keywords = map[string]TokenType{
	NewModelKeyword:         TokenModel,
	AddFieldsKeyword:        TokenField,
	AddMethodKeyword:        TokenMethod,
	SetMixinKeyword:         TokenMethod,
	AddValidationKeyword:    TokenMethod,
	SetRouteKeyword:         TokenMethod,
	SearchKeyword:           TokenMethod,
	SearchByNameKeyword:     TokenMethod,
	CreateKeyword:           TokenMethod,
	NewKeyword:              TokenMethod,
	WriteKeyword:            TokenMethod,
	CopyKeyword:             TokenMethod,
	CopyDataKeyword:         TokenMethod,
	CartesianProductKeyword: TokenMethod,
	SortedKeyword:           TokenMethod,
	FilteredKeyword:         TokenMethod,
	AggregatesKeyword:       TokenMethod,
	FirstKeyword:            TokenMethod,
	AllKeyword:              TokenMethod,
	DefaultGetKeyword:       TokenMethod,
	"func":                  TokenFunc,   // Added func keyword
	"NewMethod":             TokenMethod, // Added NewMethod
}

// lookupIdent checks if the identifier is a keyword.
func lookupIdent(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return TokenIdentifier
}
