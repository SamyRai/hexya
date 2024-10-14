package parser

import (
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/dsl"
	"github.com/hexya-erp/hexya/src/tools/dsl/ast"
	"github.com/hexya-erp/hexya/src/tools/dsl/lexer"
	"strings"
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
// Parse takes the tokens from the lexer and parses them into a DSLModel.
func (p *DSLParser) Parse(tokens []lexer.Token) (*ast.DSLModel, error) {
	p.tokens = tokens
	model := &ast.DSLModel{
		Fields:      []*ast.DSLField{},      // Initialize as empty slice
		Methods:     []*ast.DSLMethod{},     // Initialize as empty slice
		Relations:   []*ast.DSLRelation{},   // Initialize as empty slice
		Validations: []*ast.DSLValidation{}, // Initialize as empty slice
		Routes:      []*ast.DSLRoute{},      // Initialize as empty slice
		Mixins:      []*ast.DSLMixin{},      // Initialize as empty slice
	}

	for p.pos < len(p.tokens) {
		token := p.currentToken()

		fmt.Printf("Processing token: %+v\n", token) // Debug info

		switch token.Type {
		case lexer.TokenModel:
			modelName, err := p.parseModelName()
			if err != nil {
				return nil, err
			}
			model.Name = modelName
		case lexer.TokenField:
			fields, err := p.parseFields()
			if err != nil {
				return nil, err
			}
			model.Fields = append(model.Fields, fields...)
		case lexer.TokenMethod:
			methods, err := p.parseMethods()
			if err != nil {
				return nil, err
			}
			model.Methods = append(model.Methods, methods...)
		case lexer.TokenEOF:
			// Stop parsing when EOF is reached
			fmt.Println("Reached EOF. Stopping parsing.")
			return model, nil
		default:
			return nil, fmt.Errorf("unexpected token: %s", token.Value)
		}
		p.advance()
	}
	return model, nil
}

// parseModelName parses the model's name from the tokens.
// parseModelName parses the model's name from the tokens.
// parseModelName parses the model's name from the tokens.
// parseModelName parses the model's name from the tokens.
func (p *DSLParser) parseModelName() (string, error) {
	p.advance() // Skip 'NewModel'

	// Expect the opening parenthesis
	if p.currentToken().Value != "(" {
		return "", errors.New("expected '(' after 'NewModel'")
	}
	p.advance() // Move past '('

	// Expect the actual model name (e.g., "User")
	if p.currentToken().Type != lexer.TokenUnknown {
		return "", errors.New("expected model name after '('")
	}
	modelName := p.currentToken().Value
	p.advance() // Move past model name

	// Expect the closing parenthesis
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

	// Ensure the current token is valid for field definitions
	if p.currentToken().Type != lexer.TokenUnknown {
		return nil, fmt.Errorf("expected field definition but got: %s", p.currentToken().Value)
	}

	fieldStr := p.currentToken().Value
	fieldTokens := strings.Split(fieldStr, dsl.SymbolComma) // Split by comma

	for _, fieldToken := range fieldTokens {
		parts := strings.Split(fieldToken, dsl.SymbolColon) // Split by colon
		fmt.Printf("FieldToken: %s, Parts: %+v\n", fieldToken, parts)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field format: %s", fieldToken)
		}

		field := &ast.DSLField{
			Name: strings.TrimSpace(parts[0]), // Trim spaces around the name
			Type: strings.TrimSpace(parts[1]), // Trim spaces around the type
		}

		// Validate field format (optional check)
		if field.Name == "" || field.Type == "" {
			return nil, fmt.Errorf("invalid field format: %s", fieldToken)
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// parseMethods parses methods for the model.
func (p *DSLParser) parseMethods() ([]*ast.DSLMethod, error) {
	var methods []*ast.DSLMethod
	p.advance() // Skip 'AddMethod'

	methodStr := p.currentToken().Value
	methodTokens := strings.Split(methodStr, dsl.SymbolSemicolon)

	for _, methodToken := range methodTokens {
		parts := strings.Split(methodToken, dsl.SymbolColon)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid method format: %s", methodToken)
		}
		method := &ast.DSLMethod{
			Name: parts[0],
			Body: parts[1],
		}
		methods = append(methods, method)
	}

	return methods, nil
}

// Helper functions to manage token stream
func (p *DSLParser) currentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *DSLParser) advance() {
	fmt.Printf("Advancing from pos: %d, token: %+v\n", p.pos, p.currentToken())
	p.pos++
}
