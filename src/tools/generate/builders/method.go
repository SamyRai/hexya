package builders

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"strings"
)

// AddMethodsToModelData processes methods from ModelASTData and populates ModelData.
func AddMethodsToModelData(modelASTData *models.ModelData, modelData *models.ModelData, depsMap *map[string]bool) {
	for _, methodAST := range modelASTData.Methods {
		processMethod(methodAST, modelData, depsMap)
	}
}

// processMethod processes each MethodAST into MethodData and appends it to ModelData.
func processMethod(methodAST *models.MethodAST, modelData *models.ModelData, depsMap *map[string]bool) {
	var params, paramsWithType, returnAsserts, returnString string

	// Generate parameters for method
	for _, param := range methodAST.Params {
		paramType := param.Type.TypeName
		if param.IsVariadic {
			paramType = "..." + paramType
		}
		params += param.Name + ","
		paramsWithType += fmt.Sprintf("%s %s,", param.Name, paramType)
		(*depsMap)[param.Type.ImportPath] = true
	}

	// Process returns with logic for single and multiple return values
	if len(methodAST.Returns) == 1 {
		returnType := methodAST.Returns[0].Type.TypeName
		returnAsserts = fmt.Sprintf("resTyped, _ := res.(%s)", returnType)
		returnString = returnType
		(*depsMap)[methodAST.Returns[0].Type.ImportPath] = true
	} else if len(methodAST.Returns) > 1 {
		for i, ret := range methodAST.Returns {
			retType := ret.Type.TypeName
			returnAsserts += fmt.Sprintf("resTyped%d, _ := res[%d].(%s)\n", i, i, retType)
			returnString += retType + ","
			(*depsMap)[ret.ImportPath] = true
		}
		returnString = strings.TrimRight(returnString, ",")
	}

	// Append method data to modelData.ProcessedMethods
	modelData.ProcessedMethods = append(modelData.ProcessedMethods, models.MethodData{
		Name:        methodAST.Name,
		Params:      strings.TrimRight(params, ","),
		Returns:     returnString,
		ImportPaths: removeDuplicateImportPaths(*depsMap),
		Source:      "local",
	})
}

// Helper function to remove duplicate import paths
func removeDuplicateImportPaths(importPaths map[string]bool) []string {
	uniquePaths := make([]string, 0, len(importPaths))
	for path := range importPaths {
		if path != "" {
			uniquePaths = append(uniquePaths, path)
		}
	}
	return uniquePaths
}
