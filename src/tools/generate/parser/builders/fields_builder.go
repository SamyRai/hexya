package builders

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	models2 "github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"github.com/hexya-erp/hexya/src/tools/strutils"
)

// AddFieldsToModelData extracts data from modelASTData and populates it into modelData.
// It also ensures that dependencies are updated in depsMap.

// AddFieldsToModelData extracts data from modelASTData and populates it into modelData.
func AddFieldsToModelData(modelASTData parser.ModelASTData, modelData *models.ModelData) {
	relatedModels := collectRelatedModels(modelASTData, modelData)
	fmt.Printf("\nRelated models: %v\n", relatedModels)

	for relModel := range relatedModels {
		modelData.RelModels = append(modelData.RelModels, relModel)
	}

	// Populate ImportPath in Fields
	for _, fieldAST := range modelASTData.Fields {
		field := buildFieldData(fieldAST, modelData)
		modelData.Fields = append(modelData.Fields, field)
	}
}

// collectRelatedModels identifies related models from the field definitions.
func collectRelatedModels(modelASTData parser.ModelASTData, modelData *models.ModelData) map[string]bool {
	relModels := make(map[string]bool)

	for _, fieldAST := range modelASTData.Fields {
		if fieldAST.RelModel != "" {
			relModels[fieldAST.RelModel] = true
		}
	}

	return relModels
}

// getFieldTypeStrings returns the type strings for a given field, including interface type.
func getFieldTypeStrings(fieldASTData models2.FieldASTData) (typStr, iTypStr string) {
	typStr = fieldASTData.Type.Type
	iTypStr = utils.TrimInterfacePackagePrefix(typStr)

	if fieldASTData.RelModel != "" {
		typStr = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, fieldASTData.RelModel)
		iTypStr = fmt.Sprintf("%sSet", fieldASTData.RelModel)
	}
	return
}

// buildFieldData creates and returns a FieldData object from the AST field models2.
func buildFieldData(fieldName string, fieldASTData models2.FieldASTData, typStr, iTypStr string) models2.FieldData {
	jsonName := getJSONName(fieldName, fieldASTData)

	return models2.FieldData{
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
func getJSONName(fieldName string, fieldASTData models2.FieldASTData) string {
	return strutils.GetDefaultString(fieldASTData.JSON, models.SnakeCaseFieldName(fieldName, fieldASTData.FType))
}
