package parser

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
)

// inflateMixinsForAST populates fields and methods from mixins in the given model.
func inflateMixinsForAST(modelName string, modelsData map[string]*models.ModelData) {
	model, exists := modelsData[modelName]
	if !exists || model == nil {
		fmt.Printf("Model %s not found or nil in modelsData\n", modelName)
		return
	}

	if len(model.Mixins) == 0 {
		fmt.Printf("Model %s has no mixins to inflate\n", modelName)
		return
	}

	for _, mixin := range model.Mixins {
		inflateMixinsForAST(mixin.Name, modelsData) // Recursively inflate mixins

		for _, field := range mixin.Fields {
			if field.Name == "ID" {
				continue
			}
			field.Attributes.MixinField = true
			if !model.HasField(field.Name) {
				model.Fields = append(model.Fields, field)
			}
		}

		for _, method := range mixin.Methods {
			method.ToDeclare = true
			if !model.HasMethod(method.Name) {
				model.Methods = append(model.Methods, method)
			}
		}
	}
}

// inflateEmbedsForAST adds fields from embedded models to the current model.
func inflateEmbedsForAST(modelName string, modelsData map[string]*models.ModelData) {
	model, exists := modelsData[modelName]
	if !exists || model == nil {
		fmt.Printf("Model %s not found or nil in modelsData\n", modelName)
		return
	}

	for _, embed := range model.EmbeddedModels {
		inflateEmbedsForAST(embed.Name, modelsData) // Recursively inflate embedded models

		for _, field := range embed.Fields {
			field.Attributes.EmbedField = true
			if !model.HasField(field.Name) {
				model.Fields = append(model.Fields, field)
			}
		}
	}
}

// ParseModels extracts models, fields, and methods from modules' AST, then inflates mixins and embeds.

// ParseModels transforms modelASTData into a format suitable for code generation.
func ParseModels(astData map[string]*models.ModelData) map[string]*models.ModelData {
	modelDataMap := make(map[string]*models.ModelData)

	// Convert AST data to ModelData, including fields, methods, dependencies
	for modelName, astModel := range astData {
		modelData := &models.ModelData{
			Name:         modelName,
			ModelType:    astModel.ModelType,
			IsModelMixin: astModel.IsModelMixin,
		}

		// Transfer fields
		for _, fieldAST := range astModel.Fields {
			modelData.Fields = append(modelData.Fields, fieldAST)
		}

		// Transfer methods
		for _, methodAST := range astModel.Methods {
			modelData.Methods = append(modelData.Methods, methodAST)
		}

		modelDataMap[modelName] = modelData
	}

	return modelDataMap
}

// ExtractFunctionName retrieves the function name from a CallExpr node.
func ExtractFunctionName(node *ast.CallExpr) (string, error) {
	switch fun := node.Fun.(type) {
	case *ast.Ident:
		return fun.Name, nil
	case *ast.SelectorExpr:
		return fun.Sel.Name, nil
	default:
		return "", nil
	}
}
