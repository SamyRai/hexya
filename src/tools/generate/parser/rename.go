package parser

import (
	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
	"log"
	"strings"
)

// GetModelsASTDataForModules extracts models' AST data for the given modules, inflates mixins and embeds if needed.
func GetModelsASTDataForModules(modInfos []*models.ModuleInfo, validate bool) map[string]*models.ModelData {
	modelsData := make(map[string]*models.ModelData)

	for _, modInfo := range modInfos {
		for _, file := range modInfo.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.CallExpr:
					fnctName, err := ExtractFunctionName(node)
					if err != nil {
						return true
					}
					switch fnctName {
					case "addMethod":
						parseAddMethod(node, &modelsData, false)
					case "NewMethod":
						parseAddMethod(node, &modelsData, true)
					case "InheritModel":
						parseMixInModel(node, &modelsData)
					case "AddFields":
						parseAddFields(node, &modelsData)
					case "NewModel", "NewMixinModel", "NewTransientModel":
						parseNewModel(node, &modelsData)
					}
				}
				return true
			})
		}
	}

	if !validate {
		// Skip validation if not required
		return modelsData
	}

	// Validate models and inflate mixins and embedded fields/methods
	for modelName := range modelsData {
		if !modelsData[modelName].Validated {
			delete(modelsData, modelName)
		}
		inflateMixinsForAST(modelName, modelsData)
		inflateEmbedsForAST(modelName, modelsData)
	}

	return modelsData
}

// parseAddFields parses the given node which is an AddFields function
func parseAddFields(node *ast.CallExpr, modelsData *map[string]*models.ModelData) {
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		log.Panic("Unable to extract model while visiting AST", "error", err)
	}

	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = &models.ModelData{Name: modelName}
	}

	var fields *ast.CompositeLit
	switch n := node.Args[0].(type) {
	case *ast.CompositeLit:
		fields = n
	case *ast.Ident:
		fields = n.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.CompositeLit)
	}

	for _, f := range fields.Elts {
		fDef := f.(*ast.KeyValueExpr)
		fieldName := strings.Trim(fDef.Key.(*ast.BasicLit).Value, "\"`")
		var typeStr string

		switch ft := fDef.Value.(*ast.CompositeLit).Type.(type) {
		case *ast.Ident:
			typeStr = strings.TrimSuffix(ft.Name, "Field")
		case *ast.SelectorExpr:
			typeStr = strings.TrimSuffix(ft.Sel.Name, "Field")
		}

		//var fieldParams []ast.Expr
		//switch fd := fDef.Value.(type) {
		//case *ast.Ident:
		//	fieldParams = fd.Obj.Decl.(*ast.CompositeLit).Elts
		//case *ast.CompositeLit:
		//	fieldParams = fd.Elts
		//}

		fType := fieldtype.Type(strings.ToLower(typeStr))
		fData := models.FieldAST{
			Name: fieldName,
			Type: models.TypeAST{
				TypeName:    fType.DefaultGoType().String(),
				ImportPath:  "",
				IsRecordSet: false,
			},
		}

		(*modelsData)[modelName].Fields = append((*modelsData)[modelName].Fields, &fData)
	}
}

// parseAddMethod parses the given node which is an addMethod function.
func parseAddMethod(node *ast.CallExpr, modelsData *map[string]*models.ModelData, toDeclare bool) {
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		log.Panicf("Unable to extract model: %v", err)
	}
	methodName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"`")

	// Check if the second argument is a function literal (*ast.FuncLit)
	var funcType *ast.FuncType
	switch t := node.Args[1].(type) {
	case *ast.FuncLit:
		funcType = t.Type
	case *ast.Ident:
		// Handle case when it's an identifier
		if decl, ok := t.Obj.Decl.(*ast.FuncDecl); ok {
			funcType = decl.Type
		} else {
			log.Panicf("Unhandled identifier type: %T", t.Obj.Decl)
		}
	default:
		log.Panicf("Unexpected argument type for method: %T", t)
	}

	methodData := &models.MethodAST{
		Name:    methodName,
		Params:  extractParams(funcType),     // Custom function to extract method parameters
		Returns: extractReturnType(funcType), // Custom function to extract method returns
	}

	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = &models.ModelData{Name: modelName}
	}
	(*modelsData)[modelName].Methods = append((*modelsData)[modelName].Methods, methodData)
}

// parseNewModel parses the given node which is a NewModel function.
func parseNewModel(node *ast.CallExpr, modelsData *map[string]*models.ModelData) {
	modelName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"`")

	var modelType string
	switch fun := node.Fun.(type) {
	case *ast.Ident:
		modelType = strings.TrimSuffix(strings.TrimPrefix(fun.Name, "New"), "Model")
	case *ast.SelectorExpr:
		modelType = strings.TrimSuffix(strings.TrimPrefix(fun.Sel.Name, "New"), "Model")
	default:
		log.Panicf("Unexpected function type for NewModel: %T", fun)
	}

	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = &models.ModelData{Name: modelName, ModelType: modelType}
	}
}

// parseMixInModel updates the mixin tree with the given node which is an InheritModel function.
func parseMixInModel(node *ast.CallExpr, modelsData *map[string]*models.ModelData) {
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		log.Panic("Unable to extract model", err)
	}
	mixinModel, err := extractModel(node.Args[0])
	if err != nil {
		log.Panic("Unable to extract mixin model", err)
	}
	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = &models.ModelData{Name: modelName}
	}
	(*modelsData)[modelName].Mixins = append((*modelsData)[modelName].Mixins, &models.ModelData{Name: mixinModel})
}

func extractParams(funcType *ast.FuncType) []models.ParamAST {
	var params []models.ParamAST
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			paramAST := models.ParamAST{
				Name: name.Name,
				Type: models.TypeAST{
					TypeName: getTypeString(param.Type),
				},
			}
			params = append(params, paramAST)
		}
	}
	return params
}

func extractReturnType(funcType *ast.FuncType) []models.ReturnAST {
	var returns []models.ReturnAST

	// Check if the function has any return types
	if funcType.Results == nil {
		return returns // Return an empty slice if no return types are defined
	}

	// Iterate over return types if they exist
	for _, ret := range funcType.Results.List {
		returnAST := models.ReturnAST{
			Type: models.TypeAST{
				TypeName: getTypeString(ret.Type),
			},
		}
		returns = append(returns, returnAST)
	}
	return returns
}
