package builders

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"strings"
)

// AddMethodsToModelData extracts data from modelsASTData to populate methods in modelData
func AddMethodsToModelData(modelsASTData map[string]ast.ModelASTData, modelData *data.ModelData, depsMap *map[string]bool) {
	modelASTData := modelsASTData[modelData.Name]
	specificMethodsHandlers := getSpecificMethodsHandlers()

	// Iterate over all methods in the model
	for methodName, methodASTData := range modelASTData.Methods {
		// Check for specific handlers, if applicable
		if handler, exists := specificMethodsHandlers[methodName]; exists {
			handler(&methodASTData, modelData, depsMap)
			continue
		}

		// Process parameters and return types
		params, paramsWithType, iParamsWithType, paramsType := processParams(methodASTData, modelData, modelsASTData, depsMap)
		call, returns, returnAsserts, returnString, iReturnString := processReturns(methodASTData, modelData, modelsASTData, depsMap)

		// Add processed method data to modelData
		addToModelData(methodName, methodASTData, modelData, params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString)
	}
}

// Splitting parameter processing logic
func processParams(methodASTData ast.MethodASTData, modelData *data.ModelData, modelsASTData map[string]ast.ModelASTData, depsMap *map[string]bool) (params, paramsWithType, iParamsWithType, paramsType string) {
	for _, astParam := range methodASTData.Params {
		paramType := astParam.Type.Type
		iParamType := utils.TrimInterfacePackagePrefix(paramType)
		paramName := fmt.Sprintf("%s,", astParam.Name)

		// Check if the parameter is a RecordSet type and adjust accordingly
		if isRS, _ := IsRecordSetType(paramType, modelsASTData); isRS {
			// Adjust RecordSet type references
			iParamType = fmt.Sprintf("%sSet", modelData.Name)
			paramType = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
		}

		// Handle variadic parameters correctly
		if astParam.Variadic {
			iParamType = fmt.Sprintf("...%s", iParamType)
			paramType = fmt.Sprintf("...%s", paramType)
		}

		// Append parameter names and types
		params += paramName
		paramsWithType += fmt.Sprintf("%s %s,", astParam.Name, paramType)
		iParamsWithType += fmt.Sprintf("%s %s,", astParam.Name, iParamType)
		paramsType += fmt.Sprintf("%s,", paramType)

		// Ensure the dependency map is updated for the import path
		addPackageToDepsMap(depsMap, astParam.Type.ImportPath)
	}
	return
}

// Splitting return processing logic
func processReturns(methodASTData ast.MethodASTData, modelData *data.ModelData, modelsASTData map[string]ast.ModelASTData, depsMap *map[string]bool) (call, returns, returnAsserts, returnString, iReturnString string) {
	// Single return case
	if len(methodASTData.Returns) == 1 {
		call = "Call"
		returnType := methodASTData.Returns[0].Type
		iReturnType := utils.TrimInterfacePackagePrefix(returnType)

		// Handle RecordSet types properly
		if isRS, _ := IsRecordSetType(returnType, modelsASTData); isRS {
			returnType = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
			iReturnType = fmt.Sprintf("%sSet", modelData.Name)
			returnAsserts = fmt.Sprintf("resTyped := res.(models.RecordSet).Collection().Wrap(\"%s\").(%s)", modelData.Name, returnType)
		} else {
			returnAsserts = fmt.Sprintf("resTyped, _ := res.(%s)", returnType)
		}

		// Finalize return strings
		returns = "resTyped"
		returnString = returnType
		iReturnString = iReturnType
		addPackageToDepsMap(depsMap, methodASTData.Returns[0].ImportPath)
	}

	// Multiple returns case
	if len(methodASTData.Returns) > 1 {
		call = "CallMulti"
		for i, ret := range methodASTData.Returns {
			returnType := ret.Type
			iReturnType := utils.TrimInterfacePackagePrefix(returnType)

			// Handle RecordSet types for multiple returns
			if isRS, _ := IsRecordSetType(returnType, modelsASTData); isRS {
				returnType = fmt.Sprintf("%s.%sSet", config.PoolInterfacesPackage, modelData.Name)
				iReturnType = fmt.Sprintf("%sSet", modelData.Name)
				returnAsserts += fmt.Sprintf("resTyped%d := res[%d].(models.RecordSet).Collection().Wrap(\"%s\").(%s)\n", i, i, modelData.Name, returnType)
			} else {
				returnAsserts += fmt.Sprintf("resTyped%d, _ := res[%d].(%s)\n", i, i, returnType)
			}

			// Accumulate return strings
			returnString += fmt.Sprintf("%s,", returnType)
			iReturnString += fmt.Sprintf("%s,", iReturnType)
			returns += fmt.Sprintf("resTyped%d,", i)

			// Add necessary imports
			addPackageToDepsMap(depsMap, ret.ImportPath)
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

// Adds only the package to depsMap, ensuring no duplicate imports and no type-specific paths
func addPackageToDepsMap(depsMap *map[string]bool, fullPath string) {
	packagePath := getPackagePath(fullPath)
	if packagePath != "" {
		(*depsMap)[packagePath] = true
	}
}

// Extracts just the package path from a full import path
func getPackagePath(fullPath string) string {
	// Check if the full path contains a type reference after a dot (.)
	if idx := strings.LastIndex(fullPath, "."); idx != -1 {
		// Get the substring up to the last '/'
		return fullPath[:idx]
	}
	return fullPath
}

// isRecordSetType returns true if the given typ is a RecordSet according
// to the AST data stored in models.
// The second returned value is true if typ is models.RecordCollection or models.RecordSet
// and false if it is a specific RecordSet type
func IsRecordSetType(typ string, models map[string]ast.ModelASTData) (bool, bool) {
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
