package builders

import (
	"fmt"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"github.com/hexya-erp/hexya/src/tools/strutils"
)

// AddFieldsToModelData extracts data from modelASTData and populates it into modelData.
// It also ensures that dependencies are updated in depsMap.
func AddFieldsToModelData(modelASTData ast.ModelASTData, modelData *data.ModelData, depsMap *map[string]bool) {
	// Track related models for later
	relatedModels := collectRelatedModels(modelASTData, modelData, depsMap)

	// Append related models to modelData
	for relModel := range relatedModels {
		modelData.RelModels = append(modelData.RelModels, relModel)
	}
}

// collectRelatedModels processes each field in modelASTData and adds it to modelData.
// It returns a set of related models found in the fields.
func collectRelatedModels(modelASTData ast.ModelASTData, modelData *data.ModelData, depsMap *map[string]bool) map[string]bool {
	relModels := make(map[string]bool)

	// Process fields and populate the modelData
	for fieldName, fieldASTData := range modelASTData.Fields {
		// Get the field types and sanitized versions
		typStr, iTypStr := getFieldTypeStrings(fieldASTData)

		// Add related models to the map if present
		if fieldASTData.RelModel != "" {
			relModels[fieldASTData.RelModel] = true
		}

		// Build the FieldData and append it to modelData
		modelData.Fields = append(modelData.Fields, buildFieldData(fieldName, fieldASTData, typStr, iTypStr))

		// Add field import path to dependency map
		//addDependency(depsMap, fieldASTData.Type.ImportPath)
	}

	return relModels
}

// getFieldTypeStrings returns the type strings for a given field, including interface type.
func getFieldTypeStrings(fieldASTData ast.FieldASTData) (typStr, iTypStr string) {
	typStr = fieldASTData.Type.Type
	iTypStr = utils.TrimInterfacePackagePrefix(typStr)

	if fieldASTData.RelModel != "" {
		typStr = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, fieldASTData.RelModel)
		iTypStr = fmt.Sprintf("%sSet", fieldASTData.RelModel)
	}
	return
}

// buildFieldData creates and returns a FieldData object from the AST field data.
func buildFieldData(fieldName string, fieldASTData ast.FieldASTData, typStr, iTypStr string) data.FieldData {
	jsonName := getJSONName(fieldName, fieldASTData)

	return data.FieldData{
		Name:       fieldName,
		JSON:       jsonName,
		Type:       typStr,
		IType:      iTypStr,
		IsRS:       fieldASTData.IsRS,
		RelModel:   fieldASTData.RelModel,
		SanType:    utils.CreateTypeIdent(typStr),
		MixinField: fieldASTData.MixinField,
		EmbedField: fieldASTData.EmbedField,
		ImportPath: fieldASTData.Type.ImportPath,
	}
}

// getJSONName returns the JSON field name, defaulting to snake case if not provided.
func getJSONName(fieldName string, fieldASTData ast.FieldASTData) string {
	return strutils.GetDefaultString(fieldASTData.JSON, models.SnakeCaseFieldName(fieldName, fieldASTData.FType))
}

// addDependency adds a valid import path to the dependency map.
func addDependency(depsMap *map[string]bool, importPath string) {
	if importPath != "" {
		(*depsMap)[importPath] = true
	}
}
