package lexer

import (
	"strings"
	"unicode"
)

// Helper functions for character classification.
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isOperatorStart(ch rune) bool {
	return strings.ContainsRune("+-*/%=&|<>!^", ch)
}
