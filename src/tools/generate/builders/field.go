package builders

import (
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// AddFieldsToModelData populates fields from the AST into ModelData and handles their dependencies.
func AddFieldsToModelData(modelData *models.ModelData) {
	// Process and append each field
	for _, field := range modelData.Fields {
		modelData.ProcessedFields = append(modelData.ProcessedFields, buildFieldData(field))
	}
}

// buildFieldData creates a FieldData structure from the FieldAST.
func buildFieldData(fieldAST *models.FieldAST) models.FieldData {
	return models.FieldData{
		Name:       fieldAST.Name,
		Type:       fieldAST.Type.TypeName,
		IsRS:       fieldAST.RelationModel != nil,
		RelModel:   getRelModelName(fieldAST),
		ImportPath: fieldAST.Type.ImportPath,
		MixinField: fieldAST.Attributes.MixinField,
		EmbedField: fieldAST.Attributes.EmbedField,
	}
}

// getRelModelName safely extracts the related model's name if available.
func getRelModelName(fieldAST *models.FieldAST) string {
	if fieldAST.RelationModel != nil {
		return fieldAST.RelationModel.ModelName
	}
	return ""
}
