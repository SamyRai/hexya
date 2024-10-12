// Package ast generate/models_ast_data.go
package ast

import (
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/logging"
	"github.com/hexya-erp/hexya/src/tools/strutils"
	"strings"

	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"go/ast"
)

var log = logging.GetLogger("ast")

// GetModelsASTData GetModelsASTDataForModules returns the MethodASTData for all methods in given modules.
// If validate is true, then only models that have been explicitly declared will appear in
// the result. Mixins and embeddings will be inflated too. Use this if you want to validate the
// whole application.
func GetModelsASTData(modInfos []*models.ModuleInfo, validate bool) map[string]ModelASTData {
	fmt.Printf("\n\n\n[INFO] Getting models AST data for modules: %d\n", len(modInfos))

	modelsData := make(map[string]ModelASTData)
	for _, modInfo := range modInfos {
		fmt.Printf("\n\t[INFO] Processing module: %s\n", modInfo.Name)
		for _, file := range modInfo.Syntax {
			fmt.Printf("\n\t\t[INFO] Processing file: %s\n", modInfo.FSet.Position(file.Pos()))
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.CallExpr:
					fnctName, err := ExtractFunctionName(node)
					if err != nil {
						fmt.Printf("[ERROR] Unable to extract function name while visiting AST: %v\n", err)
						return true
					}
					switch {
					case fnctName == "addMethod":
						parseAddMethod(node, modInfo, &modelsData, false)
					case fnctName == "NewMethod":
						parseAddMethod(node, modInfo, &modelsData, true)
					case fnctName == "InheritModel":
						parseMixInModel(node, modInfo, &modelsData)
					case fnctName == "AddFields":
						parseAddFields(node, modInfo, &modelsData)
					case strutils.StartsAndEndsWith(fnctName, "New", "Model"):
						parseNewModel(node, &modelsData)
					}
				}
				return true
			})
		}
	}
	if !validate {
		fmt.Printf("[INFO] Skipping validation of models AST data\n")
		return modelsData
	}
	for modelName, md := range modelsData {
		// Delete models that have not been declared explicitly
		// Because it means we have a typing error
		if !md.Validated {
			delete(modelsData, modelName)
		}
		inflateMixins(modelName, &modelsData)
		inflateEmbeds(modelName, &modelsData)
	}
	return modelsData
}

// inflateMixins populates the given model with fields and methods defined in its mixins
func inflateMixins(modelName string, modelsData *map[string]ModelASTData) {
	for mixin := range (*modelsData)[modelName].Mixins {
		inflateMixins(mixin, modelsData) // Recursively inflate mixins

		// Adding mixin fields
		for fieldName, field := range (*modelsData)[mixin].Fields {
			if fieldName == "ID" {
				continue
			}
			field.MixinField = true
			(*modelsData)[modelName].Fields[fieldName] = field
		}

		// Adding mixin methods
		for methodName, method := range (*modelsData)[mixin].Methods {
			method.ToDeclare = true
			(*modelsData)[modelName].Methods[methodName] = method
		}
	}
}

// inflateEmbeds populates the given model with fields from the embedded type
func inflateEmbeds(modelName string, modelsData *map[string]ModelASTData) {
	for emb := range (*modelsData)[modelName].Embeds {
		relModel := (*modelsData)[modelName].Fields[emb].RelModel
		inflateEmbeds(relModel, modelsData)
		for fieldName, field := range (*modelsData)[relModel].Fields {
			if _, exists := (*modelsData)[modelName].Fields[fieldName]; exists {
				continue
			}
			embeddedField := field
			embeddedField.EmbedField = true
			(*modelsData)[modelName].Fields[fieldName] = embeddedField
		}
	}
}

// parseMixInModel updates the mixin tree with the given node which is a InheritModel function
func parseMixInModel(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData *map[string]ModelASTData) {
	fmt.Printf("[INFO] Parsing mixin model for node: %v", node)
	fNode := node.Fun.(*ast.SelectorExpr)
	modelName, err := extractModel(fNode.X, modInfo)
	if err != nil {
		var generalMixinError generalMixinError
		if errors.As(err, &generalMixinError) {
			return
		}
		log.Panic("Unable to extract model while visiting AST", "error", err, "node", modInfo.FSet.Position(node.Pos()))
	}
	mixinModel, err := extractModel(node.Args[0], modInfo)
	if err != nil {
		log.Panic("Unable to extract mixin model while visiting AST", "error", err)
	}
	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = NewModelASTData(modelName)
	}
	(*modelsData)[modelName].Mixins[mixinModel] = true
}

// parseNewModel parses the given node which is a NewXXXModel function
func parseNewModel(node *ast.CallExpr, modelsData *map[string]ModelASTData) {
	fmt.Printf("[INFO] Parsing new model for node: %+v", node.Args)
	fName, _ := ExtractFunctionName(node)
	modelName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"")
	modelType := strings.TrimSuffix(strings.TrimPrefix(fName, "New"), "Model")
	fmt.Printf("[INFO] Setting model data for model: %s, type: %s", modelName, modelType)
	setModelData(modelsData, modelName, modelType)
}

// setModelData adds a model with the given name and type to the given modelsData
func setModelData(modelsData *map[string]ModelASTData, modelName string, modelType string) {
	fmt.Printf("[INFO] Setting model data for model: %s, type: %s\n", modelName, modelType)
	model, exists := (*modelsData)[modelName]
	if !exists {
		model = NewModelASTData(modelName)
	}

	// Ensure common mixins are included
	if modelName != "CommonMixin" {
		model.Mixins["CommonMixin"] = true
	}

	// Set model type and mixins
	switch modelType {
	case "":
		model.Mixins["BaseMixin"] = true
		model.Mixins["ModelMixin"] = true
	case "Transient":
		model.Mixins["BaseMixin"] = true
	}

	model.ModelType = modelType
	model.Validated = true // Mark the model as validated
	(*modelsData)[modelName] = model
}

// ExtractFunctionName returns the name of the called function in the given call expression.

// parseAddMethod parses the given node which is an addMethod function
func parseAddMethod(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData *map[string]ModelASTData, toDeclare bool) {
	fNode, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		fmt.Printf("[ERROR] Unexpected function node type: %T\n", node.Fun)
		return
	}

	modelName, err := extractModel(fNode.X, modInfo)
	if err != nil {
		fmt.Printf("[ERROR] Unable to extract model while visiting AST: %v\n", err)
		return
	}
	methodName := strings.Trim(node.Args[0].(*ast.BasicLit).Value, "\"")

	var funcType *ast.FuncType
	var doc string

	switch fd := node.Args[1].(type) {
	case *ast.Ident:
		if funcDecl, ok := fd.Obj.Decl.(*ast.FuncDecl); ok {
			funcType = funcDecl.Type
			doc = funcDecl.Doc.Text()
		} else {
			fmt.Printf("[ERROR] Unexpected function declaration type: %T\n", fd.Obj.Decl)
			return
		}
	case *ast.FuncLit:
		funcType = fd.Type
	default:
		fmt.Printf("[ERROR] Unsupported argument type for function: %T\n", fd)
		return
	}

	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = NewModelASTData(modelName)
	}

	methData := MethodASTData{
		Name:      methodName,
		Doc:       utils.FormatDocString(doc),
		PkgPath:   modInfo.PkgPath,
		Params:    extractParams(funcType, modInfo),
		Returns:   extractReturnType(funcType, modInfo),
		ToDeclare: toDeclare,
	}
	fmt.Printf("[INFO] Adding method '%s' to model '%s'\n", methodName, modelName)
	(*modelsData)[modelName].Methods[methodName] = methData
}

// parseAddFields parses the given node which is an AddFields function
func parseAddFields(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData *map[string]ModelASTData) {
	fNode := node.Fun.(*ast.SelectorExpr)
	modelName, err := extractModel(fNode.X, modInfo)
	if err != nil {
		log.Panic("Unable to extract model while visiting AST", "error", err)
	}
	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = NewModelASTData(modelName)
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
		var importPath string
		if typeStr == "Date" || typeStr == "DateTime" {
			importPath = config.DatesPath
		}

		var fieldParams []ast.Expr
		switch fd := fDef.Value.(type) {
		case *ast.Ident:
			fieldParams = fd.Obj.Decl.(*ast.CompositeLit).Elts
		case *ast.CompositeLit:
			fieldParams = fd.Elts
		}
		fType := fieldtype.Type(strings.ToLower(typeStr))
		fData := FieldASTData{
			Name:  fieldName,
			FType: fType,
			Type: data.TypeData{
				Type:       fType.DefaultGoType().String(),
				ImportPath: importPath,
			},
		}
		for _, elem := range fieldParams {
			fElem := elem.(*ast.KeyValueExpr)
			fData = parseFieldAttribute(fElem, fData, modInfo)
			if fData.embed {
				(*modelsData)[modelName].Embeds[fieldName] = true
			}
		}
		(*modelsData)[modelName].Fields[fieldName] = fData
	}
}

// A generalMixinError is returned if the mixin is
// a general mixin set in NewXXXXModel function.
type generalMixinError struct{}

// Error method for generalMixinError
func (gme generalMixinError) Error() string {
	return "General Mixin Error"
}

// NewModelASTData computeExportPath returns the import path of the given type
func NewModelASTData(name string) ModelASTData {
	return ModelASTData{
		Name:         name,
		Fields:       defaultFields(name),
		IsModelMixin: config.ModelMixins[name],
		Methods:      make(map[string]MethodASTData),
		Mixins:       make(map[string]bool),
		Embeds:       make(map[string]bool),
		ModelType:    "",
	}
}
