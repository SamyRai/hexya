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

// AddFieldsToModelData extracts data from modelASTData to populate fields in modelData
func AddFieldsToModelData(modelASTData ast.ModelASTData, modelData *data.ModelData, depsMap *map[string]bool) {
	relModels := make(map[string]bool)
	for fieldName, fieldASTData := range modelASTData.Fields {
		typStr := fieldASTData.Type.Type
		iTypStr := utils.TrimInterfacePackagePrefix(typStr)
		if fieldASTData.RelModel != "" {
			relModels[fieldASTData.RelModel] = true
			typStr = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, fieldASTData.RelModel)
			iTypStr = fmt.Sprintf("%sSet", fieldASTData.RelModel)
		}
		jsonName := strutils.GetDefaultString(fieldASTData.JSON, models.SnakeCaseFieldName(fieldName, fieldASTData.FType))
		modelData.Fields = append(modelData.Fields, data.FieldData{
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
		})
		(*depsMap)[fieldASTData.Type.ImportPath] = true
	}
	for rm := range relModels {
		modelData.RelModels = append(modelData.RelModels, rm)
	}
}
