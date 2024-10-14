package dsl_test

import (
	"testing"

	"github.com/hexya-erp/hexya/src/tools/dsl"
	"github.com/hexya-erp/hexya/src/tools/dsl/ast"
	"github.com/hexya-erp/hexya/src/tools/dsl/lexer"
	"github.com/hexya-erp/hexya/src/tools/dsl/parser"
	"github.com/hexya-erp/hexya/src/tools/dsl/validator"
	"github.com/stretchr/testify/assert"
)

// Test the Lexer for various DSL constructs
func TestLexerForDSLConstructs(t *testing.T) {
	tests := []struct {
		input    string
		expected []lexer.Token
	}{
		{
			input: `NewModel("User")`,
			expected: []lexer.Token{
				{Type: lexer.TokenModel, Value: dsl.KeywordNewModel},
				{Type: lexer.TokenUnknown, Value: "User"}, // Without quotes
				{Type: lexer.TokenEOF},
			},
		},
		{
			input: `AddFields(Name:Char, Age:Integer)`,
			expected: []lexer.Token{
				{Type: lexer.TokenField, Value: dsl.KeywordAddFields},
				{Type: lexer.TokenUnknown, Value: `Name`},
				{Type: lexer.TokenUnknown, Value: `Char`},
				{Type: lexer.TokenUnknown, Value: `Age`},
				{Type: lexer.TokenUnknown, Value: `Integer`},
				{Type: lexer.TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		lexer := lexer.NewLexer(tt.input)
		tokens := lexer.Lex()

		assert.Equal(t, tt.expected, tokens)
	}
}

// Test the Parser for correct AST generation
// Test the Parser for correct AST generation
func TestParserForDSLConstructs(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.DSLModel
	}{
		{
			input: `NewModel("User")`,
			expected: &ast.DSLModel{
				Name:    "User",
				Fields:  []*ast.DSLField{},
				Methods: []*ast.DSLMethod{},
			},
		},
		{
			input: `NewModel("User").AddFields(Name:Char, Age:Integer)`,
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
		lexer := lexer.NewLexer(tt.input)
		tokens := lexer.Lex()
		parser := parser.NewDSLParser()
		actualAST, err := parser.Parse(tokens)

		assert.Nil(t, err)
		assert.True(t, equalModel(tt.expected, actualAST), "Parsed model does not match expected structure")
	}
}

// Test the Validator for correct model validation
func TestValidatorForDSLModels(t *testing.T) {
	tests := []struct {
		model    *ast.DSLModel
		expected error
	}{
		{
			model: &ast.DSLModel{
				Name: "User",
				Fields: []*ast.DSLField{
					{Name: "Name", Type: "Char"},
					{Name: "Age", Type: "Integer"},
				},
				Methods: []*ast.DSLMethod{},
			},
			expected: nil,
		},
		{
			model: &ast.DSLModel{
				Name: "",
				Fields: []*ast.DSLField{
					{Name: "Name", Type: "Char"},
				},
			},
			expected: validator.ErrModelNameEmpty,
		},
	}

	for _, tt := range tests {
		err := validator.ValidateModel(tt.model)
		assert.Equal(t, tt.expected, err)
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
	if expected.Fields != nil && actual.Fields == nil {
		return false
	}
	if expected.Fields == nil && actual.Fields != nil {
		return false
	}
	// Similarly for methods, relations, etc.
	if len(expected.Methods) != len(actual.Methods) {
		return false
	}

	// Additional slice comparison logic for Relations, Validations, Routes, and Mixins
	// ...

	return true
}
