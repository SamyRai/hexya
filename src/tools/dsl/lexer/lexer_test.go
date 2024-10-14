package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_BasicKeywordsAndIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Basic NewModel",
			input: `NewModel("User")`,
			expected: []Token{
				{Type: TokenModel, Value: "NewModel"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenString, Value: "User"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Basic AddFields",
			input: `AddFields(Name:Char)`,
			expected: []Token{
				{Type: TokenField, Value: "AddFields"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "Name"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Lex()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_ComplexInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Model with AddFields",
			input: `NewModel("User").AddFields(Name:Char, Email:Char, Age:Integer)`,
			expected: []Token{
				{Type: TokenModel, Value: "NewModel"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenString, Value: "User"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenDot, Value: "."},
				{Type: TokenField, Value: "AddFields"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "Name"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "Email"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "Age"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Integer"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Model with Relations",
			input: `NewModel("Product").AddFields(Name:Char, Price:Float, Available:Boolean).SetRelation(Category:Many2One)`,
			expected: []Token{
				{Type: TokenModel, Value: "NewModel"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenString, Value: "Product"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenDot, Value: "."},
				{Type: TokenField, Value: "AddFields"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "Name"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "Price"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Float"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "Available"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Boolean"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "SetRelation"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "Category"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Many2One"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Field with Compute Method",
			input: `DaysUntilDue:Integer{String:"Days Until Due", Compute:h.TodoItem().NewMethod("ComputeDaysUntilDue", func(rs m.TodoItemSet) m.TodoItemData {...})}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "DaysUntilDue"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Integer"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "String"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: "Days Until Due"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "Compute"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "h"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItem"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenDot, Value: "."},
				{Type: TokenMethod, Value: "NewMethod"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenString, Value: "ComputeDaysUntilDue"},
				{Type: TokenComma, Value: ","},
				{Type: TokenFunc, Value: "func"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "rs"},
				{Type: TokenIdentifier, Value: "m"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItemSet"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenIdentifier, Value: "m"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItemData"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenEllipsis, Value: "..."},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Lex()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_ComplexHexyaModels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "OnChange Method",
			input: `OnChange:h.TodoItem().NewMethod("OnDueDateChange", func(rs m.TodoItemSet) m.TodoItemData {...})`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "OnChange"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "h"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItem"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenDot, Value: "."},
				{Type: TokenMethod, Value: "NewMethod"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenString, Value: "OnDueDateChange"},
				{Type: TokenComma, Value: ","},
				{Type: TokenFunc, Value: "func"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "rs"},
				{Type: TokenIdentifier, Value: "m"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItemSet"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenIdentifier, Value: "m"},
				{Type: TokenDot, Value: "."},
				{Type: TokenIdentifier, Value: "TodoItemData"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenEllipsis, Value: "..."},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Lex()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "String with Escaped Characters",
			input: `field:Char{String:"Name with \"quotes\" and \n newlines"}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "field"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "String"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: `Name with \"quotes\" and \n newlines`},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Integer Default Value",
			input: `Count:Integer{Default:42}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Count"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Integer"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "Default"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenNumber, Value: "42"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Float Default Value",
			input: `Price:Float{Default:99.99}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Price"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Float"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "Default"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenNumber, Value: "99.99"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Function with Operators",
			input: `ComputeTotal:func(a int, b int) int { return a + b }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "ComputeTotal"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenFunc, Value: "func"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenIdentifier, Value: "int"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "b"},
				{Type: TokenIdentifier, Value: "int"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenIdentifier, Value: "int"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "return"},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenIdentifier, Value: "b"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Unicode Identifiers",
			input: `用户:Char{String:"用户名称"}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "用户"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenIdentifier, Value: "Char"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "String"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: "用户名称"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Nested Function Calls",
			input: `Calculate:func(a int) int { return square(multiply(a, a)) }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Calculate"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenFunc, Value: "func"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenIdentifier, Value: "int"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenIdentifier, Value: "int"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "return"},
				{Type: TokenIdentifier, Value: "square"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "multiply"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Complex Expression with Operators",
			input: `Expression:func() bool { return (a > b) && (c != d) }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Expression"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenFunc, Value: "func"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenIdentifier, Value: "bool"},
				{Type: TokenLeftCurly, Value: "{"},
				{Type: TokenIdentifier, Value: "return"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "a"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenIdentifier, Value: "b"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenOperator, Value: "&&"},
				{Type: TokenLeftParen, Value: "("},
				{Type: TokenIdentifier, Value: "c"},
				{Type: TokenOperator, Value: "!="},
				{Type: TokenIdentifier, Value: "d"},
				{Type: TokenRightParen, Value: ")"},
				{Type: TokenRightCurly, Value: "}"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Lex()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}
