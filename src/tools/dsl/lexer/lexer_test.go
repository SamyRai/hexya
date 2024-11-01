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
				{Type: TokenNewModel, Value: "NewModel", Line: 1, Column: 1},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 9},
				{Type: TokenString, Value: "User", Line: 1, Column: 10},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 16},
				{Type: TokenEOF, Value: "", Line: 1, Column: 16},
			},
		},
		{
			name:  "Basic AddFields",
			input: `AddFields(Name:Char)`,
			expected: []Token{
				{Type: TokenAddFields, Value: "AddFields", Line: 1, Column: 1},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 10},
				{Type: TokenIdentifier, Value: "Name", Line: 1, Column: 11},
				{Type: TokenColon, Value: ":", Line: 1, Column: 15},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 16},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 20},
				{Type: TokenEOF, Value: "", Line: 1, Column: 20},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, _ := lexer.Lex()
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
				{Type: TokenNewModel, Value: "NewModel", Line: 1, Column: 1},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 9},
				{Type: TokenString, Value: "User", Line: 1, Column: 10},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 16},
				{Type: TokenDot, Value: ".", Line: 1, Column: 17},
				{Type: TokenAddFields, Value: "AddFields", Line: 1, Column: 18},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 27},
				{Type: TokenIdentifier, Value: "Name", Line: 1, Column: 28},
				{Type: TokenColon, Value: ":", Line: 1, Column: 32},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 33},
				{Type: TokenComma, Value: ",", Line: 1, Column: 37},
				{Type: TokenIdentifier, Value: "Email", Line: 1, Column: 39},
				{Type: TokenColon, Value: ":", Line: 1, Column: 44},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 45},
				{Type: TokenComma, Value: ",", Line: 1, Column: 49},
				{Type: TokenIdentifier, Value: "Age", Line: 1, Column: 51},
				{Type: TokenColon, Value: ":", Line: 1, Column: 54},
				{Type: TokenIdentifier, Value: "Integer", Line: 1, Column: 55},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 62},
				{Type: TokenEOF, Value: "", Line: 1, Column: 62},
			},
		},
		{
			name:  "Model with Relations",
			input: `NewModel("Product").AddFields(Name:Char, Price:Float, Available:Boolean).SetRelation(Category:Many2One)`,
			expected: []Token{
				{Type: TokenNewModel, Value: "NewModel", Line: 1, Column: 1},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 9},
				{Type: TokenString, Value: "Product", Line: 1, Column: 10},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 18},
				{Type: TokenDot, Value: ".", Line: 1, Column: 19},
				{Type: TokenAddFields, Value: "AddFields", Line: 1, Column: 20},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 29},
				{Type: TokenIdentifier, Value: "Name", Line: 1, Column: 30},
				{Type: TokenColon, Value: ":", Line: 1, Column: 34},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 35},
				{Type: TokenComma, Value: ",", Line: 1, Column: 39},
				{Type: TokenIdentifier, Value: "Price", Line: 1, Column: 41},
				{Type: TokenColon, Value: ":", Line: 1, Column: 46},
				{Type: TokenIdentifier, Value: "Float", Line: 1, Column: 47},
				{Type: TokenComma, Value: ",", Line: 1, Column: 52},
				{Type: TokenIdentifier, Value: "Available", Line: 1, Column: 54},
				{Type: TokenColon, Value: ":", Line: 1, Column: 63},
				{Type: TokenIdentifier, Value: "Boolean", Line: 1, Column: 64},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 71},
				{Type: TokenDot, Value: ".", Line: 1, Column: 72},
				{Type: TokenIdentifier, Value: "SetRelation", Line: 1, Column: 73},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 84},
				{Type: TokenIdentifier, Value: "Category", Line: 1, Column: 85},
				{Type: TokenColon, Value: ":", Line: 1, Column: 93},
				{Type: TokenIdentifier, Value: "Many2One", Line: 1, Column: 94},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 102},
				{Type: TokenEOF, Value: "", Line: 1, Column: 102},
			},
		},
		{
			name:  "Field with Compute Method",
			input: `DaysUntilDue:Integer{String:"Days Until Due", Compute:h.TodoItem().NewMethod("ComputeDaysUntilDue", func(rs m.TodoItemSet) m.TodoItemData {...})}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "DaysUntilDue", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 13},
				{Type: TokenIdentifier, Value: "Integer", Line: 1, Column: 14},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 21},
				{Type: TokenIdentifier, Value: "String", Line: 1, Column: 22},
				{Type: TokenColon, Value: ":", Line: 1, Column: 28},
				{Type: TokenString, Value: "Days Until Due", Line: 1, Column: 29},
				{Type: TokenComma, Value: ",", Line: 1, Column: 44},
				{Type: TokenIdentifier, Value: "Compute", Line: 1, Column: 46},
				{Type: TokenColon, Value: ":", Line: 1, Column: 53},
				{Type: TokenIdentifier, Value: "h", Line: 1, Column: 54},
				{Type: TokenDot, Value: ".", Line: 1, Column: 55},
				{Type: TokenIdentifier, Value: "TodoItem", Line: 1, Column: 56},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 64},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 65},
				{Type: TokenDot, Value: ".", Line: 1, Column: 66},
				{Type: TokenNewMethod, Value: "NewMethod", Line: 1, Column: 67},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 76},
				{Type: TokenString, Value: "ComputeDaysUntilDue", Line: 1, Column: 77},
				{Type: TokenComma, Value: ",", Line: 1, Column: 97},
				{Type: TokenFunc, Value: "func", Line: 1, Column: 99},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 103},
				{Type: TokenIdentifier, Value: "rs", Line: 1, Column: 104},
				{Type: TokenIdentifier, Value: "m", Line: 1, Column: 107},
				{Type: TokenDot, Value: ".", Line: 1, Column: 108},
				{Type: TokenIdentifier, Value: "TodoItemSet", Line: 1, Column: 109},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 120},
				{Type: TokenIdentifier, Value: "m", Line: 1, Column: 122},
				{Type: TokenDot, Value: ".", Line: 1, Column: 123},
				{Type: TokenIdentifier, Value: "TodoItemData", Line: 1, Column: 124},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 137},
				{Type: TokenEllipsis, Value: "...", Line: 1, Column: 138},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 141},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 142},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 143},
				{Type: TokenEOF, Value: "", Line: 1, Column: 143},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, _ := lexer.Lex()
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
				{Type: TokenIdentifier, Value: "OnChange", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 9},
				{Type: TokenIdentifier, Value: "h", Line: 1, Column: 10},
				{Type: TokenDot, Value: ".", Line: 1, Column: 11},
				{Type: TokenIdentifier, Value: "TodoItem", Line: 1, Column: 12},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 20},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 21},
				{Type: TokenDot, Value: ".", Line: 1, Column: 22},
				{Type: TokenNewMethod, Value: "NewMethod", Line: 1, Column: 23},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 32},
				{Type: TokenString, Value: "OnDueDateChange", Line: 1, Column: 33},
				{Type: TokenComma, Value: ",", Line: 1, Column: 50},
				{Type: TokenFunc, Value: "func", Line: 1, Column: 52},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 56},
				{Type: TokenIdentifier, Value: "rs", Line: 1, Column: 57},
				{Type: TokenIdentifier, Value: "m", Line: 1, Column: 60},
				{Type: TokenDot, Value: ".", Line: 1, Column: 61},
				{Type: TokenIdentifier, Value: "TodoItemSet", Line: 1, Column: 62},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 73},
				{Type: TokenIdentifier, Value: "m", Line: 1, Column: 75},
				{Type: TokenDot, Value: ".", Line: 1, Column: 76},
				{Type: TokenIdentifier, Value: "TodoItemData", Line: 1, Column: 77},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 90},
				{Type: TokenEllipsis, Value: "...", Line: 1, Column: 91},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 94},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 95},
				{Type: TokenEOF, Value: "", Line: 1, Column: 95},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, _ := lexer.Lex()
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
				{Type: TokenIdentifier, Value: "field", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 6},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 7},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 11},
				{Type: TokenIdentifier, Value: "String", Line: 1, Column: 12},
				{Type: TokenColon, Value: ":", Line: 1, Column: 18},
				{Type: TokenString, Value: `Name with \"quotes\" and \n newlines`, Line: 1, Column: 19},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 54},
				{Type: TokenEOF, Value: "", Line: 1, Column: 54},
			},
		},
		{
			name:  "Integer Default Value",
			input: `Count:Integer{Default:42}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Count", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 6},
				{Type: TokenIdentifier, Value: "Integer", Line: 1, Column: 7},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 14},
				{Type: TokenIdentifier, Value: "Default", Line: 1, Column: 15},
				{Type: TokenColon, Value: ":", Line: 1, Column: 22},
				{Type: TokenNumber, Value: "42", Line: 1, Column: 23},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 25},
				{Type: TokenEOF, Value: "", Line: 1, Column: 25},
			},
		},
		{
			name:  "Float Default Value",
			input: `Price:Float{Default:99.99}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Price", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 6},
				{Type: TokenIdentifier, Value: "Float", Line: 1, Column: 7},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 12},
				{Type: TokenIdentifier, Value: "Default", Line: 1, Column: 13},
				{Type: TokenColon, Value: ":", Line: 1, Column: 20},
				{Type: TokenFloat, Value: "99.99", Line: 1, Column: 21},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 26},
				{Type: TokenEOF, Value: "", Line: 1, Column: 26},
			},
		},
		{
			name:  "Function with Operators",
			input: `ComputeTotal:func(a int, b int) int { return a + b }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "ComputeTotal", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 13},
				{Type: TokenFunc, Value: "func", Line: 1, Column: 14},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 18},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 19},
				{Type: TokenIdentifier, Value: "int", Line: 1, Column: 21},
				{Type: TokenComma, Value: ",", Line: 1, Column: 24},
				{Type: TokenIdentifier, Value: "b", Line: 1, Column: 26},
				{Type: TokenIdentifier, Value: "int", Line: 1, Column: 28},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 31},
				{Type: TokenIdentifier, Value: "int", Line: 1, Column: 33},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 37},
				{Type: TokenReturn, Value: "return", Line: 1, Column: 39},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 46},
				{Type: TokenOperator, Value: "+", Line: 1, Column: 48},
				{Type: TokenIdentifier, Value: "b", Line: 1, Column: 50},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 52},
				{Type: TokenEOF, Value: "", Line: 1, Column: 52},
			},
		},
		{
			name:  "Unicode Identifiers",
			input: `用户:Char{String:"用户名称"}`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "用户", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 4},
				{Type: TokenIdentifier, Value: "Char", Line: 1, Column: 5},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 9},
				{Type: TokenIdentifier, Value: "String", Line: 1, Column: 10},
				{Type: TokenColon, Value: ":", Line: 1, Column: 16},
				{Type: TokenString, Value: "用户名称", Line: 1, Column: 17},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 22},
				{Type: TokenEOF, Value: "", Line: 1, Column: 22},
			},
		},
		{
			name:  "Nested Function Calls",
			input: `Calculate:func(a int) int { return square(multiply(a, a)) }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Calculate", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 10},
				{Type: TokenFunc, Value: "func", Line: 1, Column: 11},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 15},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 16},
				{Type: TokenIdentifier, Value: "int", Line: 1, Column: 18},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 21},
				{Type: TokenIdentifier, Value: "int", Line: 1, Column: 23},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 27},
				{Type: TokenReturn, Value: "return", Line: 1, Column: 29},
				{Type: TokenIdentifier, Value: "square", Line: 1, Column: 36},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 42},
				{Type: TokenIdentifier, Value: "multiply", Line: 1, Column: 43},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 51},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 52},
				{Type: TokenComma, Value: ",", Line: 1, Column: 53},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 55},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 56},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 57},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 58},
				{Type: TokenEOF, Value: "", Line: 1, Column: 58},
			},
		},
		{
			name:  "Complex Expression with Operators",
			input: `Expression:func() bool { return (a > b) && (c != d) }`,
			expected: []Token{
				{Type: TokenIdentifier, Value: "Expression", Line: 1, Column: 1},
				{Type: TokenColon, Value: ":", Line: 1, Column: 11},
				{Type: TokenFunc, Value: "func", Line: 1, Column: 12},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 16},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 17},
				{Type: TokenIdentifier, Value: "bool", Line: 1, Column: 19},
				{Type: TokenLeftCurly, Value: "{", Line: 1, Column: 24},
				{Type: TokenReturn, Value: "return", Line: 1, Column: 26},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 33},
				{Type: TokenIdentifier, Value: "a", Line: 1, Column: 34},
				{Type: TokenOperator, Value: ">", Line: 1, Column: 36},
				{Type: TokenIdentifier, Value: "b", Line: 1, Column: 38},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 39},
				{Type: TokenOperator, Value: "&&", Line: 1, Column: 41},
				{Type: TokenLeftParen, Value: "(", Line: 1, Column: 44},
				{Type: TokenIdentifier, Value: "c", Line: 1, Column: 45},
				{Type: TokenOperator, Value: "!=", Line: 1, Column: 47},
				{Type: TokenIdentifier, Value: "d", Line: 1, Column: 50},
				{Type: TokenRightParen, Value: ")", Line: 1, Column: 51},
				{Type: TokenRightCurly, Value: "}", Line: 1, Column: 53},
				{Type: TokenEOF, Value: "", Line: 1, Column: 53},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, _ := lexer.Lex()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}
