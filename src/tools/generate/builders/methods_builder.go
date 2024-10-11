package builders

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
)

// AddMethodsToModelData extracts data from modelsASTData to populate methods in modelData
func AddMethodsToModelData(modelsASTData map[string]ast.ModelASTData, modelData *data.ModelData, depsMap *map[string]bool) {
	modelASTData := modelsASTData[modelData.Name]
	specificMethodsHandlers := getSpecificMethodsHandlers()

	for methodName, methodASTData := range modelASTData.Methods {
		if handler, exists := specificMethodsHandlers[methodName]; exists {
			handler(&methodASTData, modelData, depsMap)
			continue
		}

		params, paramsWithType, iParamsWithType, paramsType := processParams(methodASTData, modelData, modelsASTData, depsMap)
		call, returns, returnAsserts, returnString, iReturnString := processReturns(methodASTData, modelData, modelsASTData, depsMap)

		addToModelData(methodName, methodASTData, modelData, params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString)
	}
}

// Splitting parameter processing logic
func processParams(methodASTData ast.MethodASTData, modelData *data.ModelData, modelsASTData map[string]ast.ModelASTData, depsMap *map[string]bool) (params, paramsWithType, iParamsWithType, paramsType string) {
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

		utils.AddImport(depsMap, astParam.Type.ImportPath)
	}
	return
}

// Splitting return processing logic
func processReturns(methodASTData ast.MethodASTData, modelData *data.ModelData, modelsASTData map[string]ast.ModelASTData, depsMap *map[string]bool) (call, returns, returnAsserts, returnString, iReturnString string) {
	if len(methodASTData.Returns) == 1 {
		call = "Call"
		utils.AddImport(depsMap, methodASTData.Returns[0].ImportPath)
		typ := methodASTData.Returns[0].Type
		iTyp := utils.TrimInterfacePackagePrefix(typ)
		returnAsserts = fmt.Sprintf("resTyped, _ := res.(%s)", typ)
		returns = "resTyped"

		if isRS, _ := IsRecordSetType(typ, modelsASTData); isRS {
			typ = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
			iTyp = fmt.Sprintf("%sSet", modelData.Name)
			returnAsserts = fmt.Sprintf("resTyped := res.(models.RecordSet).Collection().Wrap(\"%s\").(%s)", modelData.Name, typ)
		}

		returnString = typ
		iReturnString = iTyp
	} else if len(methodASTData.Returns) > 1 {
		call = "CallMulti"
		for i, ret := range methodASTData.Returns {
			typ := ret.Type
			iTyp := utils.TrimInterfacePackagePrefix(typ)
			utils.AddImport(depsMap, ret.ImportPath)

			if isRS, _ := IsRecordSetType(ret.Type, modelsASTData); isRS {
				typ = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
				iTyp = fmt.Sprintf("%sSet", modelData.Name)
				returnAsserts += fmt.Sprintf("resTyped%d := res[%d].(models.RecordSet).Collection().Wrap(\"%s\").(%s)\n", i, i, modelData.Name, typ)
			} else {
				returnAsserts += fmt.Sprintf("resTyped%d, _ := res[%d].(%s)\n", i, i, typ)
			}

			returnString += fmt.Sprintf("%s,", typ)
			iReturnString += fmt.Sprintf("%s,", iTyp)
			returns += fmt.Sprintf("resTyped%d,", i)
		}
	}
	return
}

// Function to add processed data to modelData
func addToModelData(methodName string, methodASTData ast.MethodASTData, modelData *data.ModelData, params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString string) {
	modelData.AllMethods = append(modelData.AllMethods, data.MethodData{
		Name:             methodName,
		Doc:              methodASTData.Doc,
		ToDeclare:        methodASTData.ToDeclare,
		ParamsTypes:      utils.JoinStrings([]string{paramsType}, ","),
		IParamsWithTypes: utils.JoinStrings([]string{iParamsWithType}, ","),
		ReturnString:     utils.JoinStrings([]string{returnString}, ","),
		IReturnString:    utils.JoinStrings([]string{iReturnString}, ","),
	})

	modelData.Methods = append(modelData.Methods, data.MethodData{
		Name:           methodName,
		Doc:            methodASTData.Doc,
		ToDeclare:      methodASTData.ToDeclare,
		Params:         utils.JoinStrings([]string{params}, ","),
		ParamsWithType: utils.JoinStrings([]string{paramsWithType}, ","),
		ReturnAsserts:  utils.TrimTrailingNewline(returnAsserts),
		Returns:        utils.TrimTrailingComma(returns),
		ReturnString:   utils.TrimTrailingComma(returnString),
		Call:           call,
	})
}

// IsRecordSetType returns true if the given type is a RecordSet according to modelsASTData.
func IsRecordSetType(typ string, modelsASTData map[string]ast.ModelASTData) (bool, bool) {
	if typ == "*RecordCollection" || typ == "*models.RecordCollection" {
		return true, true
	}
	if typ == "RecordSet" || typ == "models.RecordSet" {
		return true, true
	}
	if _, exists := modelsASTData[utils.TrimRecordSetSuffix(typ)]; exists {
		return true, false
	}
	return false, false
}
