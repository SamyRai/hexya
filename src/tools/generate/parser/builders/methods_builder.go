package builders

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"

	"github.com/hexya-erp/hexya/src/tools/generate/utils"
)

// AddMethodsToModelData extracts methods from AST data and populates ModelData with MethodData.
func AddMethodsToModelData(modelsASTData map[string]parser.ModelASTData, modelData *models.ModelData) {
	modelASTData := modelsASTData[modelData.Name]

	// Iterate over all methods in the model
	for methodName, methodASTData := range modelASTData.Methods {
		processMethod(methodName, methodASTData, modelData, modelsASTData)
	}
}

func processMethod(methodName string, methodASTData models.MethodASTData, modelData *models.ModelData, modelsASTData map[string]parser.ModelASTData) {
	// Process parameters and return types
	params, paramsWithType, iParamsWithType, paramsType, importPaths := processParams(methodASTData, modelData, modelsASTData)
	call, returns, returnAsserts, returnString, iReturnString, returnImportPaths := processReturns(methodASTData, modelData, modelsASTData)

	// Combine import paths from parameters and returns
	importPaths = append(importPaths, returnImportPaths...)

	// Add processed method data to ModelData
	addMethodToModelData(methodName, methodASTData, modelData, params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString, importPaths)
}

// processParams handles parameters processing and collects import paths.
func processParams(methodASTData models.MethodASTData, modelData *models.ModelData, modelsASTData map[string]parser.ModelASTData) (params, paramsWithType, iParamsWithType, paramsType string, importPaths []string) {
	for _, astParam := range methodASTData.Params {
		paramType := astParam.Type.Type
		iParamType := utils.TrimInterfacePackagePrefix(paramType)
		paramName := fmt.Sprintf("%s,", astParam.Name)

		if isRS, _ := IsRecordSetType(paramType, modelsASTData); isRS {
			iParamType = fmt.Sprintf("%sSet", modelData.Name)
			paramType = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
		}

		if astParam.Variadic {
			iParamType = fmt.Sprintf("...%s", iParamType)
			paramType = fmt.Sprintf("...%s", paramType)
		}

		params += paramName
		paramsWithType += fmt.Sprintf("%s %s,", astParam.Name, paramType)
		iParamsWithType += fmt.Sprintf("%s %s,", astParam.Name, iParamType)
		paramsType += fmt.Sprintf("%s,", paramType)

		// Collect the import path directly in MethodData
		if astParam.Type.ImportPath != "" {
			importPaths = append(importPaths, astParam.Type.ImportPath)
		}
	}
	return
}

// processReturns handles return values processing and collects import paths.
func processReturns(methodASTData models.MethodASTData, modelData *models.ModelData, modelsASTData map[string]parser.ModelASTData) (call, returns, returnAsserts, returnString, iReturnString string, importPaths []string) {
	if len(methodASTData.Returns) == 1 {
		call = "Call"
		returnType := methodASTData.Returns[0].Type
		iReturnType := utils.TrimInterfacePackagePrefix(returnType)

		if isRS, _ := IsRecordSetType(returnType, modelsASTData); isRS {
			returnType = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
			iReturnType = fmt.Sprintf("%sSet", modelData.Name)
			returnAsserts = fmt.Sprintf("resTyped := res.(models.RecordSet).Collection().Wrap(\"%s\").(%s)", modelData.Name, returnType)
		} else {
			returnAsserts = fmt.Sprintf("resTyped, _ := res.(%s)", returnType)
		}

		returns = "resTyped"
		returnString = returnType
		iReturnString = iReturnType

		// Collect the import path directly in MethodData
		if methodASTData.Returns[0].ImportPath != "" {
			importPaths = append(importPaths, methodASTData.Returns[0].ImportPath)
		}
	}
	return
}

// addMethodToModelData adds the processed MethodData to ModelData.
func addMethodToModelData(methodName string, methodASTData models.MethodASTData, modelData *models.ModelData, params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString string, importPaths []string) {
	methodData := models.MethodData{
		Name:           methodName,
		Params:         params,
		ParamsWithType: paramsWithType,
		ReturnAsserts:  returnAsserts,
		Returns:        returns,
		ReturnString:   returnString,
		IReturnString:  iReturnString,
		Call:           call,
		ToDeclare:      methodASTData.ToDeclare,
		ImportPaths:    importPaths,
	}
	modelData.Methods = append(modelData.Methods, methodData)
}

// isRecordSetType returns true if the given typ is a RecordSet according to the AST data stored in models.
// The second returned value is true if typ is models.RecordCollection or models.RecordSet and false if it is a specific RecordSet type
func IsRecordSetType(typ string, models map[string]parser.ModelASTData) (bool, bool) {
	if typ == "*RecordCollection" || typ == "*models.RecordCollection" {
		return true, true
	}
	if typ == "RecordSet" || typ == "models.RecordSet" {
		return true, true
	}
	if _, exists := models[strings.TrimSuffix(typ, "Set")]; exists {
		return true, false
	}
	return false, false
}
