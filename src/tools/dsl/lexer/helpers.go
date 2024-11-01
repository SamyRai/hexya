package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Keywords map for lookup.
var keywords = map[string]TokenType{
	// Go keywords
	"break":       TokenBreak,
	"case":        TokenCase,
	"chan":        TokenChan,
	"const":       TokenConst,
	"continue":    TokenContinue,
	"default":     TokenDefault,
	"defer":       TokenDefer,
	"else":        TokenElse,
	"fallthrough": TokenFallthrough,
	"for":         TokenFor,
	"func":        TokenFunc,
	"go":          TokenGo,
	"goto":        TokenGoto,
	"if":          TokenIf,
	"import":      TokenImport,
	"interface":   TokenInterface,
	"map":         TokenMap,
	"package":     TokenPackage,
	"range":       TokenRange,
	"return":      TokenReturn,
	"select":      TokenSelect,
	"struct":      TokenStruct,
	"switch":      TokenSwitch,
	"type":        TokenTypeKeyword,
	"var":         TokenVar,
	// Boolean literals
	"true":  TokenBoolLiteral,
	"false": TokenBoolLiteral,
	// Nil literal
	"nil": TokenNil,
	// Hexya-specific keywords
	"NewModel":  TokenNewModel,
	"AddFields": TokenAddFields,
	"NewMethod": TokenNewMethod,
	// Add other Hexya-specific keywords as needed
}

func (l *Lexer) readStringToken() Token {
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
				// Handle hexadecimal escapes like \x7F
				hexChar := l.readHexEscape(2) // 2 digits for hex escape
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

	return l.makeToken(TokenString, value.String())
}

func (l *Lexer) readRawStringToken() Token {
	l.readChar() // Consume starting backtick
	startPos := l.position

	for l.ch != '`' && l.ch != 0 {
		l.readChar()
	}

	value := l.input[startPos:l.position]

	if l.ch == '`' {
		l.readChar() // Consume ending backtick
	} else {
		l.appendError(ErrorUnterminatedString, "Unterminated raw string literal")
	}

	return l.makeToken(TokenString, value)
}

func IsOperatorStart(ch rune) bool {
	return strings.ContainsRune("+-*/%=&|<>!^", ch)
}

// readOperator reads an operator token, including multi-character operators.
func (l *Lexer) readOperator() Token {
	startLine, startColumn := l.line, l.column
	ch := l.ch
	var operator string

	switch ch {
	case '+':
		if l.peekChar() == '+' {
			operator = "++"
			l.readChar()
		} else if l.peekChar() == '=' {
			operator = "+="
			l.readChar()
		} else {
			operator = "+"
		}
	case '-':
		if l.peekChar() == '-' {
			operator = "--"
			l.readChar()
		} else if l.peekChar() == '=' {
			operator = "-="
			l.readChar()
		} else {
			operator = "-"
		}
	case '*':
		if l.peekChar() == '=' {
			operator = "*="
			l.readChar()
		} else {
			operator = "*"
		}
	case '/':
		if l.peekChar() == '=' {
			operator = "/="
			l.readChar()
		} else {
			operator = "/"
		}
	case '%':
		if l.peekChar() == '=' {
			operator = "%="
			l.readChar()
		} else {
			operator = "%"
		}
	case '=':
		if l.peekChar() == '=' {
			operator = "=="
			l.readChar()
		} else {
			operator = "="
		}
	case '!':
		if l.peekChar() == '=' {
			operator = "!="
			l.readChar()
		} else {
			operator = "!"
		}
	case '<':
		if l.peekChar() == '=' {
			operator = "<="
			l.readChar()
		} else if l.peekChar() == '<' {
			operator = "<<"
			l.readChar()
			if l.peekChar() == '=' {
				operator += "="
				l.readChar()
			}
		} else {
			operator = "<"
		}
	case '>':
		if l.peekChar() == '=' {
			operator = ">="
			l.readChar()
		} else if l.peekChar() == '>' {
			operator = ">>"
			l.readChar()
			if l.peekChar() == '=' {
				operator += "="
				l.readChar()
			}
		} else {
			operator = ">"
		}
	case '&':
		if l.peekChar() == '&' {
			operator = "&&"
			l.readChar()
		} else if l.peekChar() == '=' {
			operator = "&="
			l.readChar()
		} else {
			operator = "&"
		}
	case '|':
		if l.peekChar() == '|' {
			operator = "||"
			l.readChar()
		} else if l.peekChar() == '=' {
			operator = "|="
			l.readChar()
		} else {
			operator = "|"
		}
	case '^':
		if l.peekChar() == '=' {
			operator = "^="
			l.readChar()
		} else {
			operator = "^"
		}
	default:
		operator = string(ch)
	}

	l.readChar()
	return Token{
		Type:   TokenOperator,
		Value:  operator,
		Line:   startLine,
		Column: startColumn,
	}
}
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func (l *Lexer) readIdentifier() Token {
	startPosition := l.position
	startLine := l.line
	startColumn := l.column

	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	ident := l.input[startPosition:l.position]
	tokenType := LookupIdent(ident)
	return Token{
		Type:   tokenType,
		Value:  ident,
		Line:   startLine,
		Column: startColumn,
	}
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return TokenIdentifier
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func (l *Lexer) readUnicodeEscape() rune {
	var hexDigits []rune
	for i := 0; i < 4; i++ {
		l.readChar()
		if isHexDigit(l.ch) {
			hexDigits = append(hexDigits, l.ch)
		} else {
			l.appendError(ErrorInvalidEscapeSequence, fmt.Sprintf("Invalid Unicode escape sequence: \\u%s", string(hexDigits)))
			return 0
		}
	}
	hexValue := string(hexDigits)
	codePoint, err := parseHexToRune(hexValue)
	if err != nil {
		l.appendError(ErrorInvalidEscapeSequence, fmt.Sprintf("Invalid Unicode code point: %s", hexValue))
		return 0
	}
	return codePoint
}

func (l *Lexer) readHexEscape(digits int) rune {
	var hexDigits []rune
	for i := 0; i < digits; i++ {
		l.readChar()
		if isHexDigit(l.ch) {
			hexDigits = append(hexDigits, l.ch)
		} else {
			l.appendError(ErrorInvalidEscapeSequence, fmt.Sprintf("Invalid hexadecimal escape sequence: \\x%s", string(hexDigits)))
			return 0
		}
	}
	hexValue := string(hexDigits)
	codePoint, err := parseHexToRune(hexValue)
	if err != nil {
		l.appendError(ErrorInvalidEscapeSequence, fmt.Sprintf("Invalid hexadecimal code point: %s", hexValue))
		return 0
	}
	return codePoint
}

func parseHexToRune(hexStr string) (rune, error) {
	var codePoint rune
	_, err := fmt.Sscanf(hexStr, "%04x", &codePoint)
	if err != nil {
		return 0, err
	}
	return codePoint, nil
}

func isHexDigit(ch rune) bool {
	return unicode.IsDigit(ch) || ('a' <= unicode.ToLower(ch) && unicode.ToLower(ch) <= 'f')
}
