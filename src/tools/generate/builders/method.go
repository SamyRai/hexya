package builders

import (
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"strings"
)

// AddMethodsToModelData populates the methods from the model's AST into ModelData.
func AddMethodsToModelData(modelData *models.ModelData) {
	for _, methodAST := range modelData.Methods {
		processMethod(methodAST, modelData)
	}
}

// processMethod processes each MethodAST into MethodData and adds it to the ModelData.
func processMethod(methodAST *models.MethodAST, modelData *models.ModelData) {
	// Gather import paths for method parameters and return types.
	importPaths := gatherImportPathsFromMethod(methodAST)

	// Construct MethodData from MethodAST.
	methodData := models.MethodData{
		Name:        methodAST.Name,
		Params:      formatParams(methodAST.Params),
		Returns:     formatReturns(methodAST.Returns),
		ImportPaths: removeDuplicateImportPaths(importPaths),
		Source:      "local", // Assuming all methods are local unless specified otherwise.
	}

	// Append the constructed MethodData to the model's method list.
	modelData.ProcessedMethods = append(modelData.ProcessedMethods, methodData)
}

// gatherImportPathsFromMethod collects import paths from method parameters and return types.
func gatherImportPathsFromMethod(methodAST *models.MethodAST) []string {
	var importPaths []string
	for _, param := range methodAST.Params {
		if param.Type.ImportPath != "" {
			importPaths = append(importPaths, param.Type.ImportPath)
		}
	}
	for _, ret := range methodAST.Returns {
		if ret.Type.ImportPath != "" {
			importPaths = append(importPaths, ret.Type.ImportPath)
		}
	}
	return removeDuplicateImportPaths(importPaths)
}

// formatParams formats method parameters into a string for use in MethodData.
func formatParams(params []models.ParamAST) string {
	var formattedParams []string
	for _, param := range params {
		paramStr := param.Name + " " + param.Type.TypeName
		if param.IsVariadic {
			paramStr = "..." + paramStr
		}
		formattedParams = append(formattedParams, paramStr)
	}
	return strings.Join(formattedParams, ", ")
}

// formatReturns formats method return types into a string for use in MethodData.
func formatReturns(returns []models.ReturnAST) string {
	var formattedReturns []string
	for _, ret := range returns {
		formattedReturns = append(formattedReturns, ret.Type.TypeName)
	}
	return strings.Join(formattedReturns, ", ")
}

// removeDuplicateImportPaths removes duplicate import paths from the list.
func removeDuplicateImportPaths(importPaths []string) []string {
	seen := make(map[string]bool)
	var uniquePaths []string
	for _, path := range importPaths {
		if !seen[path] {
			seen[path] = true
			uniquePaths = append(uniquePaths, path)
		}
	}
	return uniquePaths
}
