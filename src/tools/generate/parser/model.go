package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

type ModelParser struct{}

func NewModelParser() *ModelParser {
	return &ModelParser{}
}

// Parse processes the model creation function and updates the model data.
// ModelParser processes model creation and updates ModelData
func (p *ModelParser) Parse(node *ast.CallExpr, modelsData map[string]*models.ModelData) {
	// Use the more robust extractModel function
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil || modelName == "" {
		fmt.Printf("Error extracting model name: %v\n", err)
		return
	}
	fmt.Printf("Parsing model: %s\n", modelName)

	// Get or create ModelData
	modelData, exists := modelsData[modelName]
	if !exists {
		modelData = &models.ModelData{
			Name:    modelName,
			Fields:  []*models.FieldAST{},
			Methods: []*models.MethodAST{},
		}
		modelsData[modelName] = modelData
	}
}

// Helper function to extract model name from CallExpr
func extractModelNameFromCallExpr(node *ast.CallExpr) (string, error) {
	if len(node.Args) == 0 {
		return "", fmt.Errorf("no arguments in NewModel call")
	}

	arg := node.Args[0]
	switch lit := arg.(type) {
	case *ast.BasicLit:
		modelName := strings.Trim(lit.Value, "\"`")
		fmt.Printf("Extracted model name from model: %s\n", modelName) // Log extracted model name
		return modelName, nil
	case *ast.Ident:
		fmt.Printf("Extracted model name from Ident: %s\n", lit.Name) // Log Ident type model name
		return lit.Name, nil
	default:
		// Log unexpected argument type
		fmt.Printf("Unexpected argument type %T in NewModel call\n", arg)
		return "", fmt.Errorf("unexpected argument type %T in NewModel call", arg)
	}
}
