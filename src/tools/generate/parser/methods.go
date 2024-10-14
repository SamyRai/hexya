package parser

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/handlers"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
	"strings"
)

// MethodParser parses methods from the AST nodes.
type MethodParser struct{}

// NewMethodParser returns a new MethodParser instance.
func NewMethodParser() *MethodParser {
	return &MethodParser{}
}

// Parse extracts methods from AST nodes and updates the corresponding ModelData.
func (p *MethodParser) Parse(node *ast.CallExpr, modelsData map[string]*models.ModelData, modInfo *models.ModuleInfo) {
	// Extract the model from the AST node using a robust model extraction logic.
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		fmt.Printf("Unable to extract model: %v\n", err)
		return
	}

	// Retrieve or initialize the ModelData.
	modelData, exists := modelsData[modelName]
	if !exists {
		modelData = &models.ModelData{
			Name:    modelName,
			Methods: []*models.MethodAST{},
		}
		modelsData[modelName] = modelData
	}

	// Extract the method name.
	methodName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"`")

	// Extract the method's parameters, return types, and documentation.
	methodAST := p.extractMethodData(node)
	methodAST.Name = methodName

	// Process the method with additional handling from handlers package.
	handlers.ProcessSpecificMethod(&methodAST, modelData)

	// Append the method to the model's Methods slice.
	modelData.Methods = append(modelData.Methods, &methodAST)
}

// extractMethodData parses the method's parameters, return types, and documentation.
func (p *MethodParser) extractMethodData(node *ast.CallExpr) models.MethodAST {
	var methodAST models.MethodAST

	// Parse the method signature (parameters and return types).
	switch fd := node.Args[1].(type) {
	case *ast.Ident:
		funcDecl := fd.Obj.Decl.(*ast.FuncDecl)
		methodAST.Doc = funcDecl.Doc.Text()
		methodAST.Params = p.extractParams(funcDecl.Type)
		methodAST.Returns = p.extractReturnTypes(funcDecl.Type)
	case *ast.FuncLit:
		methodAST.Params = p.extractParams(fd.Type)
		methodAST.Returns = p.extractReturnTypes(fd.Type)
	}

	return methodAST
}

// extractParams parses the parameters of a function.
func (p *MethodParser) extractParams(ft *ast.FuncType) []models.ParamAST {
	var params []models.ParamAST
	for i, pl := range ft.Params.List {
		if i == 0 {
			// Skip the first argument (usually "rs").
			continue
		}
		for _, param := range pl.Names {
			params = append(params, models.ParamAST{
				Name: param.Name,
				Type: models.TypeAST{
					TypeName:    getTypeString(pl.Type),
					ImportPath:  getImportPathFromExpr(pl.Type),
					IsRecordSet: isRecordSetType(pl.Type),
				},
			})
		}
	}
	return params
}

// extractReturnTypes parses the return types of a function.
func (p *MethodParser) extractReturnTypes(ft *ast.FuncType) []models.ReturnAST {
	var returns []models.ReturnAST
	if ft.Results != nil {
		for _, res := range ft.Results.List {
			returns = append(returns, models.ReturnAST{
				Type: models.TypeAST{
					TypeName:   getTypeString(res.Type),
					ImportPath: getImportPathFromExpr(res.Type),
				},
			})
		}
	}
	return returns
}
