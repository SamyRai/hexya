package parser

import (
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/dsl/ast"
	"github.com/hexya-erp/hexya/src/tools/dsl/lexer"
)

type DSLParser struct {
	tokens []lexer.Token
	pos    int
}

// NewDSLParser creates a new DSLParser instance.
func NewDSLParser() *DSLParser {
	return &DSLParser{}
}

// Parse takes the tokens from the lexer and parses them into a DSLModel.
func (p *DSLParser) Parse(tokens []lexer.Token) (*ast.DSLModel, error) {
	p.tokens = tokens
	model := &ast.DSLModel{
		Fields:      []*ast.DSLField{},
		Methods:     []*ast.DSLMethod{},
		Relations:   []*ast.DSLRelation{},
		Validations: []*ast.DSLValidation{},
		Routes:      []*ast.DSLRoute{},
		Mixins:      []*ast.DSLMixin{},
	}

	for p.pos < len(p.tokens) {
		token := p.currentToken()

		switch token.Type {
		case lexer.TokenNewModel:
			modelName, err := p.parseModelName()
			if err != nil {
				return nil, err
			}
			model.Name = modelName
		case lexer.TokenAddFields:
			fields, err := p.parseFields()
			if err != nil {
				return nil, err
			}
			model.Fields = append(model.Fields, fields...)
		case lexer.TokenEOF:
			return model, nil
		default:
			return nil, fmt.Errorf("unexpected token: %s", token.Value)
		}
		p.advance()
	}
	return model, nil
}

// parseModelName parses the model's name from the tokens.
func (p *DSLParser) parseModelName() (string, error) {
	p.advance() // Skip 'NewModel'

	if p.currentToken().Value != "(" {
		return "", errors.New("expected '(' after 'NewModel'")
	}
	p.advance() // Move past '('

	if p.currentToken().Type != lexer.TokenString {
		return "", errors.New("expected model name after '('")
	}
	modelName := p.currentToken().Value
	p.advance() // Move past model name

	if p.currentToken().Value != ")" {
		return "", errors.New("expected ')' after model name")
	}
	p.advance() // Move past ')'

	return modelName, nil
}

// parseFields parses the fields for the model.
func (p *DSLParser) parseFields() ([]*ast.DSLField, error) {
	var fields []*ast.DSLField
	p.advance() // Skip 'AddFields'

	if p.currentToken().Type != lexer.TokenLeftParen {
		return nil, fmt.Errorf("expected '(' but got: %s", p.currentToken().Value)
	}
	p.advance() // Move past '('

	for p.currentToken().Type != lexer.TokenRightParen && p.currentToken().Type != lexer.TokenEOF {
		fieldName := p.currentToken().Value
		p.advance() // Move past field name

		if p.currentToken().Type != lexer.TokenColon {
			return nil, fmt.Errorf("expected ':' after field name")
		}
		p.advance() // Move past ':'

		fieldType := p.currentToken().Value
		field := &ast.DSLField{
			Name: fieldName,
			Type: fieldType,
		}
		fields = append(fields, field)

		p.advance() // Move past field type

		if p.currentToken().Type == lexer.TokenComma {
			p.advance() // Move past ','
		}
	}

	if p.currentToken().Type != lexer.TokenRightParen {
		return nil, fmt.Errorf("expected ')' after field definitions")
	}
	p.advance() // Move past ')'

	return fields, nil
}

// Helper functions to manage token stream
func (p *DSLParser) currentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *DSLParser) advance() {
	p.pos++
}
