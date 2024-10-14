package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

type FieldParser struct{}

func NewFieldParser() *FieldParser {
	return &FieldParser{}
}

// Parse processes the AddFields call and extracts model fields into ModelData
// FieldParser processes the AddFields call and extracts model fields into ModelData
func (p *FieldParser) Parse(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData map[string]*models.ModelData) {
	// Use extractModel to handle all cases of extracting model names
	fNode, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		fmt.Printf("Unexpected node.Fun type %T in FieldParser.Parse\n", node.Fun)
		return
	}
	modelName, err := extractModel(fNode.X) // Use the robust model extraction logic
	if err != nil {
		fmt.Printf("Error extracting model: %v\n", err)
		return
	}

	// Log model name extraction
	fmt.Printf("Extracted model name from fields: %s\n", modelName)

	// Get or create ModelData
	modelData, exists := modelsData[modelName]
	if !exists {
		modelData = &models.ModelData{
			Name:   modelName,
			Fields: []*models.FieldAST{},
		}
		modelsData[modelName] = modelData
	}

	// Parse fields
	if len(node.Args) == 0 {
		fmt.Println("No arguments in AddFields call")
		return
	}
	fieldsArg := node.Args[0]

	switch fieldsArg := fieldsArg.(type) {
	case *ast.CompositeLit:
		fmt.Printf("Parsing composite literal fields for model: %s\n", modelName)
		p.parseCompositeLitFields(fieldsArg, modelData, modInfo)
	case *ast.Ident:
		fmt.Printf("Field defined as identifier: %s for model: %s\n", fieldsArg.Name, modelName)
	default:
		fmt.Printf("Unexpected fields argument type %T for model: %s\n", fieldsArg, modelName)
	}
}

// parseCompositeLitFields handles composite literal fields parsing
func (p *FieldParser) parseCompositeLitFields(fieldsMap *ast.CompositeLit, modelData *models.ModelData, modInfo *models.ModuleInfo) {
	for _, elt := range fieldsMap.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		fieldNameLit, ok := kvExpr.Key.(*ast.BasicLit)
		if !ok {
			continue
		}
		fieldName := strings.Trim(fieldNameLit.Value, "\"`")
		fieldDefExpr := kvExpr.Value
		// Process fieldDefExpr to extract field type and attributes
		fieldAST := p.parseFieldDefinition(fieldDefExpr, fieldName)
		modelData.Fields = append(modelData.Fields, fieldAST)
	}
}

// parseFieldDefinition extracts field details, attributes, and type
// parseFieldDefinition extracts field details, attributes, and type
func (p *FieldParser) parseFieldDefinition(expr ast.Expr, fieldName string) *models.FieldAST {
	// Initialize FieldAST
	fieldAST := &models.FieldAST{
		Name: fieldName,
		Type: models.TypeAST{},
		Attributes: models.FieldAttributes{
			Selection: make(map[string]string),
		},
	}

	// Handle various AST node types for field definitions
	switch field := expr.(type) {
	case *ast.Ident:
		fieldAST.Type.TypeName = CreateTypeIdent(field.Name)
	case *ast.CompositeLit:
		// Extract field type from composite literal
		switch ft := field.Type.(type) {
		case *ast.Ident:
			fieldAST.Type.TypeName = CreateTypeIdent(ft.Name)
		case *ast.SelectorExpr:
			fieldAST.Type.TypeName = CreateTypeIdent(ft.Sel.Name)
		default:
			fmt.Printf("Unexpected field type %T\n", field.Type)
		}
		// Extract field attributes from composite literal elements
		for _, elt := range field.Elts {
			kvExpr, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			keyIdent, ok := kvExpr.Key.(*ast.Ident)
			if !ok {
				continue
			}
			switch keyIdent.Name {
			case "String":
				fieldAST.Attributes.Description = extractStringValue(kvExpr.Value)
			case "Help":
				fieldAST.Attributes.Help = extractStringValue(kvExpr.Value)
			case "Selection":
				fieldAST.Attributes.Selection = extractSelection(kvExpr.Value)
			case "Required":
				fieldAST.Attributes.Required = extractBoolValue(kvExpr.Value)
			case "Default":
				fieldAST.Attributes.Default = extractStringValue(kvExpr.Value)
			case "ReadOnly":
				fieldAST.Attributes.ReadOnly = extractBoolValue(kvExpr.Value)
			case "Index":
				fieldAST.Attributes.Index = extractBoolValue(kvExpr.Value)
			case "Size":
				fieldAST.Attributes.Size = extractIntValue(kvExpr.Value)
			}
		}
	default:
		fmt.Printf("Unexpected field definition type %T\n", expr)
	}

	return fieldAST
}

// Helper functions

func extractStringValue(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		return strings.Trim(v.Value, "\"`")
	case *ast.Ident:
		return v.Name
	default:
		return ""
	}
}
