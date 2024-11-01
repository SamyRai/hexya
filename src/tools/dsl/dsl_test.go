package dsl_test

import (
	"testing"

	"github.com/hexya-erp/hexya/src/tools/dsl/ast"
	"github.com/hexya-erp/hexya/src/tools/dsl/lexer"
	"github.com/hexya-erp/hexya/src/tools/dsl/parser"
	"github.com/stretchr/testify/assert"
)

func TestParserForDSLConstructs(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected *ast.DSLModel
	}{
		{
			name: "Basic NewModel",
			tokens: []lexer.Token{
				{Type: lexer.TokenNewModel, Value: "NewModel"},
				{Type: lexer.TokenLeftParen, Value: "("},
				{Type: lexer.TokenString, Value: "User"},
				{Type: lexer.TokenRightParen, Value: ")"},
				{Type: lexer.TokenEOF, Value: ""},
			},
			expected: &ast.DSLModel{
				Name:    "User",
				Fields:  []*ast.DSLField{},
				Methods: []*ast.DSLMethod{},
			},
		},
		{
			name: "NewModel with Fields",
			tokens: []lexer.Token{
				{Type: lexer.TokenNewModel, Value: "NewModel"},
				{Type: lexer.TokenLeftParen, Value: "("},
				{Type: lexer.TokenString, Value: "User"},
				{Type: lexer.TokenRightParen, Value: ")"},
				{Type: lexer.TokenDot, Value: "."},
				{Type: lexer.TokenAddFields, Value: "AddFields"},
				{Type: lexer.TokenLeftParen, Value: "("},
				{Type: lexer.TokenIdentifier, Value: "Name"},
				{Type: lexer.TokenColon, Value: ":"},
				{Type: lexer.TokenIdentifier, Value: "Char"},
				{Type: lexer.TokenComma, Value: ","},
				{Type: lexer.TokenIdentifier, Value: "Age"},
				{Type: lexer.TokenColon, Value: ":"},
				{Type: lexer.TokenIdentifier, Value: "Integer"},
				{Type: lexer.TokenRightParen, Value: ")"},
				{Type: lexer.TokenEOF, Value: ""},
			},
			expected: &ast.DSLModel{
				Name: "User",
				Fields: []*ast.DSLField{
					{Name: "Name", Type: "Char"},
					{Name: "Age", Type: "Integer"},
				},
				Methods:     []*ast.DSLMethod{},
				Relations:   []*ast.DSLRelation{},
				Validations: []*ast.DSLValidation{},
				Routes:      []*ast.DSLRoute{},
				Mixins:      []*ast.DSLMixin{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parser.NewDSLParser()
			actualAST, err := parser.Parse(tt.tokens)

			assert.Nil(t, err)
			assert.True(t, equalModel(tt.expected, actualAST), "Parsed model does not match expected structure")
		})
	}
}

func equalModel(expected, actual *ast.DSLModel) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	// Compare model names
	if expected.Name != actual.Name {
		return false
	}

	// Check for nil and empty slice equivalence in Fields
	if len(expected.Fields) != len(actual.Fields) {
		return false
	}

	// Check fields
	for i := range expected.Fields {
		if expected.Fields[i].Name != actual.Fields[i].Name || expected.Fields[i].Type != actual.Fields[i].Type {
			return false
		}
	}

	return true
}
