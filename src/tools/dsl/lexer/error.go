package lexer

// ErrorType represents the type of a lexical error.
type ErrorType string

// Error represents a lexical error.
type Error struct {
	Type    ErrorType
	Message string
	Line    int
	Column  int
}

// Error types.
const (
	ErrorUnexpectedChar        ErrorType = "UnexpectedChar"
	ErrorInvalidToken          ErrorType = "InvalidToken"
	ErrorUnterminatedString    ErrorType = "UnterminatedString"
	ErrorInvalidEscapeSequence ErrorType = "InvalidEscapeSequence"
	ErrorUnterminatedComment   ErrorType = "UnterminatedComment"
	ErrorNestingLevelReached   ErrorType = "NestingLevelReached"
	ErrorMissingClosingParen   ErrorType = "MissingClosingParen"
)
