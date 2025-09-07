package parser

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
	"log"
	"strings"
)

// GetModelsASTDataForModules extracts AST data for given modules, with mixin and embed inflation if needed.
func GetModelsASTDataForModules(moduleInfos []*models.ModuleInfo, validate bool) map[string]*models.ModelData {
	modelASTDataMap := make(map[string]*models.ModelData)

	for _, moduleInfo := range moduleInfos {
		for _, file := range moduleInfo.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CallExpr:
					funcName, err := ExtractFunctionName(n)
					if err != nil {
						fmt.Printf("Failed to extract function name: %v\n", err)
						return true
					}
					switch funcName {
					case "addMethod":
						parseAddMethod(n, &modelASTDataMap)
					case "NewMethod":
						parseAddMethod(n, &modelASTDataMap)
					case "InheritModel":
						parseMixinModel(n, moduleInfo, &modelASTDataMap)
					case "AddFields":
						parseAddFields(n, &modelASTDataMap)
					case "NewModel", "NewMixinModel", "NewTransientModel":
						parseNewModel(n, &modelASTDataMap)
					}
				}
				return true
			})
		}
	}

	if !validate {
		return modelASTDataMap
	}

	// Validate models and inflate mixins and embedded fields/methods
	for modelName := range modelASTDataMap {
		if !modelASTDataMap[modelName].Validated {
			delete(modelASTDataMap, modelName)
		}
		inflateMixinsForAST(modelName, modelASTDataMap)
		inflateEmbedsForAST(modelName, modelASTDataMap)
	}

	return modelASTDataMap
}

// parseAddFields processes AddFields functions in the AST node.
func parseAddFields(node *ast.CallExpr, modelASTDataMap *map[string]*models.ModelData) {
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		log.Panic("Error extracting model name", err)
	}

	if _, exists := (*modelASTDataMap)[modelName]; !exists {
		(*modelASTDataMap)[modelName] = &models.ModelData{Name: modelName}
	}

	var fields *ast.CompositeLit
	switch n := node.Args[0].(type) {
	case *ast.CompositeLit:
		fields = n
	case *ast.Ident:
		fields = n.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.CompositeLit)
	}

	for _, f := range fields.Elts {
		fieldExpr := f.(*ast.KeyValueExpr)
		fieldName := strings.Trim(fieldExpr.Key.(*ast.BasicLit).Value, "\"`")
		fieldTypeStr := extractFieldType(fieldExpr.Value.(*ast.CompositeLit).Type)

		fieldType := fieldtype.Type(strings.ToLower(fieldTypeStr))
		fieldData := models.FieldAST{
			Name: fieldName,
			Type: models.TypeAST{
				TypeName:    fieldType.DefaultGoType().String(),
				ImportPath:  "",
				IsRecordSet: false,
			},
		}

		(*modelASTDataMap)[modelName].Fields = append((*modelASTDataMap)[modelName].Fields, &fieldData)
	}
}

// parseAddMethod parses AST nodes for addMethod function calls.
func parseAddMethod(node *ast.CallExpr, modelASTDataMap *map[string]*models.ModelData) {
	modelName, err := extractModel(node.Fun.(*ast.SelectorExpr).X)
	if err != nil {
		log.Panicf("Error extracting model: %v", err)
	}
	methodName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"`")

	var funcType *ast.FuncType
	switch t := node.Args[1].(type) {
	case *ast.FuncLit:
		funcType = t.Type
	case *ast.Ident:
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
		Params:  extractParams(funcType),
		Returns: extractReturnType(funcType),
	}

	if _, exists := (*modelASTDataMap)[modelName]; !exists {
		(*modelASTDataMap)[modelName] = &models.ModelData{Name: modelName}
	}
	(*modelASTDataMap)[modelName].Methods = append((*modelASTDataMap)[modelName].Methods, methodData)
}

// parseNewModel processes AST nodes for NewModel functions.
func parseNewModel(node *ast.CallExpr, modelASTDataMap *map[string]*models.ModelData) {
	modelName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"`")
	modelType := getModelTypeFromFunc(node.Fun)

	if _, exists := (*modelASTDataMap)[modelName]; !exists {
		(*modelASTDataMap)[modelName] = &models.ModelData{Name: modelName, ModelType: modelType}
	}
}

// parseMixinModel handles mixin relationships in the AST data.
// parseMixinModel updates the mixin tree with the given node which is an InheritModel function.
func parseMixinModel(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData *map[string]*models.ModelData) {
	fNode := node.Fun.(*ast.SelectorExpr)
	modelName, err := extractModel(fNode.X)
	if err != nil {
		// Check if it's a General Mixin Error, skip if true
		if _, ok := err.(generalMixinError); ok {
			fmt.Printf("Skipping general mixin error for model: %s\n", modelName)
			return
		}
		// Log and continue for other errors without stopping execution
		fmt.Printf("Error extracting model: %v\n", err)
		return
	}

	mixinModelName, err := extractModel(node.Args[0])
	if err != nil {
		fmt.Printf("Unable to extract mixin model: %v\n", err)
		return
	}

	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = &models.ModelData{Name: modelName}
	}
	(*modelsData)[modelName].Mixins = append((*modelsData)[modelName].Mixins, &models.ModelData{Name: mixinModelName})
}

// extractFieldType determines the field type from a given AST expression node.
func extractFieldType(expr ast.Expr) string {
	switch fieldType := expr.(type) {
	case *ast.Ident:
		return strings.TrimSuffix(fieldType.Name, "Field")
	case *ast.SelectorExpr:
		return strings.TrimSuffix(fieldType.Sel.Name, "Field")
	default:
		return ""
	}
}

// Additional helper functions for extracting and processing AST nodes.

func extractParams(funcType *ast.FuncType) []models.ParamAST {
	var params []models.ParamAST
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			params = append(params, models.ParamAST{
				Name: name.Name,
				Type: models.TypeAST{
					TypeName: getTypeString(param.Type),
				},
			})
		}
	}
	return params
}

func extractReturnType(funcType *ast.FuncType) []models.ReturnAST {
	var returns []models.ReturnAST
	if funcType.Results == nil {
		return returns
	}
	for _, ret := range funcType.Results.List {
		returns = append(returns, models.ReturnAST{
			Type: models.TypeAST{
				TypeName: getTypeString(ret.Type),
			},
		})
	}
	return returns
}

func getModelTypeFromFunc(funNode ast.Expr) string {
	switch fun := funNode.(type) {
	case *ast.Ident:
		return strings.TrimSuffix(strings.TrimPrefix(fun.Name, "New"), "Model")
	case *ast.SelectorExpr:
		return strings.TrimSuffix(strings.TrimPrefix(fun.Sel.Name, "New"), "Model")
	default:
		log.Panicf("Unexpected function type for NewModel: %T", fun)
		return ""
	}
}
