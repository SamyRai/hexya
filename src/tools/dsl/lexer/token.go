package lexer

// TokenType represents the type of a token.
type TokenType string

// Token represents a lexical token.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Token types.
const (
	TokenEOF          TokenType = "EOF"
	TokenIdentifier   TokenType = "IDENTIFIER"
	TokenNumber       TokenType = "NUMBER"
	TokenFloat        TokenType = "FLOAT"
	TokenString       TokenType = "STRING"
	TokenOperator     TokenType = "OPERATOR"
	TokenLeftParen    TokenType = "LEFT_PAREN"
	TokenRightParen   TokenType = "RIGHT_PAREN"
	TokenLeftCurly    TokenType = "LEFT_CURLY"
	TokenRightCurly   TokenType = "RIGHT_CURLY"
	TokenLeftBracket  TokenType = "LEFT_BRACKET"
	TokenRightBracket TokenType = "RIGHT_BRACKET"
	TokenComma        TokenType = "COMMA"
	TokenSemicolon    TokenType = "SEMICOLON"
	TokenColon        TokenType = "COLON"
	TokenDot          TokenType = "DOT"
	TokenEllipsis     TokenType = "ELLIPSIS"
	TokenDefine       TokenType = "DEFINE"
	TokenNewline      TokenType = "NEWLINE"
	TokenUnknown      TokenType = "UNKNOWN"
	// Go token types
	TokenIllegal     TokenType = "ILLEGAL"
	TokenComment     TokenType = "COMMENT"
	TokenImaginary   TokenType = "IMAGINARY"
	TokenBoolLiteral TokenType = "BOOL_LITERAL"
	TokenNil         TokenType = "NIL"
	TokenFunc        TokenType = "FUNC"
	TokenStruct      TokenType = "STRUCT"
	TokenPackage     TokenType = "PACKAGE"
	TokenImport      TokenType = "IMPORT"
	TokenReturn      TokenType = "RETURN"
	TokenVar         TokenType = "VAR"
	TokenConst       TokenType = "CONST"
	TokenTypeKeyword TokenType = "TYPE"
	TokenInterface   TokenType = "INTERFACE"
	TokenMap         TokenType = "MAP"
	TokenChan        TokenType = "CHAN"
	TokenGo          TokenType = "GO"
	TokenSelect      TokenType = "SELECT"
	TokenCase        TokenType = "CASE"
	TokenDefault     TokenType = "DEFAULT"
	TokenIf          TokenType = "IF"
	TokenElse        TokenType = "ELSE"
	TokenSwitch      TokenType = "SWITCH"
	TokenFallthrough TokenType = "FALLTHROUGH"
	TokenFor         TokenType = "FOR"
	TokenRange       TokenType = "RANGE"
	TokenBreak       TokenType = "BREAK"
	TokenContinue    TokenType = "CONTINUE"
	TokenGoto        TokenType = "GOTO"
	TokenDefer       TokenType = "DEFER"
	// Operators and delimiters
	TokenAsterisk TokenType = "ASTERISK"
	// Hexya-specific tokens
	TokenNewModel  TokenType = "NEWMODEL"
	TokenAddFields TokenType = "ADDFIELDS"
	TokenNewMethod TokenType = "NEWMETHOD"
)
