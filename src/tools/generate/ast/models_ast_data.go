// generate/models_ast_data.go
package ast

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/strutils"
	"go/printer"
	"go/types"
	"log"
	"strings"

	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"go/ast"
)

// GetModelsASTDataForModules returns the MethodASTData for all methods in given modules.
// If validate is true, then only models that have been explicitly declared will appear in
// the result. Mixins and embeddings will be inflated too. Use this if you want validate the
// whole application.
func GetModelsASTData(modules []*models.ModuleInfo) map[string]ModelASTData {
	return GetModelsASTDataForModules(modules, true)
}

// GetModelsASTDataForModules returns the MethodASTData for all methods in given modules.
// If validate is true, then only models that have been explicitly declared will appear in
// the result. Mixins and embeddings will be inflated too. Use this if you want validate the
// whole application.
func GetModelsASTDataForModules(modInfos []*models.ModuleInfo, validate bool) map[string]ModelASTData {
	fmt.Println("[INFO] Starting to parse modules for AST data")
	modelsData := make(map[string]ModelASTData)

	for _, modInfo := range modInfos {
		fmt.Printf("[INFO] Processing module: %s\n", modInfo.PkgPath)
		for _, file := range modInfo.Syntax {
			fmt.Printf("[DEBUG] Inspecting file: %s\n", modInfo.FSet.Position(file.Pos()))
			ast.Inspect(file, func(n ast.Node) bool {
				if n == nil {
					return false
				}

				switch node := n.(type) {
				case *ast.CallExpr:
					fnctName, err := ExtractFunctionName(node)
					if err != nil {
						fmt.Printf("[ERROR] Failed to extract function name: %v\n", err)
						return true
					}
					fmt.Printf("[INFO] Found function call: %s\n", fnctName)

					switch {
					case fnctName == "addMethod":
						fmt.Printf("[INFO] Parsing addMethod for module: %s\n", modInfo.PkgPath)
						parseAddMethod(node, modInfo, &modelsData, false)
					case fnctName == "NewMethod":
						fmt.Printf("[INFO] Parsing NewMethod for module: %s\n", modInfo.PkgPath)
						parseAddMethod(node, modInfo, &modelsData, true)
					case fnctName == "InheritModel":
						fmt.Printf("[INFO] Parsing InheritModel for module: %s\n", modInfo.PkgPath)
						parseMixInModel(node, modInfo, &modelsData)
					case fnctName == "AddFields":
						fmt.Printf("[INFO] Parsing AddFields for module: %s\n", modInfo.PkgPath)
						parseAddFields(node, modInfo, &modelsData)
					case strutils.StartsAndEndsWith(fnctName, "New", "Model"):
						fmt.Printf("[INFO] Parsing NewModel for function: %s\n", fnctName)
						parseNewModel(node, &modelsData)
					default:
						fmt.Printf("[DEBUG] Unhandled function call: %s\n", fnctName)
					}
				}
				return true
			})
		}
	}

	if !validate {
		fmt.Println("[INFO] Validation disabled, returning models data")
		return modelsData
	}

	fmt.Println("[INFO] Starting validation and inflating mixins and embeds")
	for modelName, md := range modelsData {
		if !md.Validated {
			fmt.Printf("[INFO] Deleting unvalidated model: %s\n", modelName)
			delete(modelsData, modelName)
			continue
		}
		fmt.Printf("[INFO] Inflating mixins and embeds for model: %s\n", modelName)
		inflateMixins(modelName, &modelsData)
		inflateEmbeds(modelName, &modelsData)
	}

	fmt.Println("[INFO] Finished parsing modules for AST data")
	return modelsData
}

// inflateMixins populates the given model with fields and methods defined in its mixins
// inflateMixins populates the given model with fields and methods defined in its mixins
func inflateMixins(modelName string, modelsData *map[string]ModelASTData) {
	fmt.Printf("[INFO] Inflating mixins for model: %s\n", modelName)

	modelData, exists := (*modelsData)[modelName]
	if !exists {
		fmt.Printf("[WARNING] Model '%s' does not exist in modelsData. Initializing it.\n", modelName)
		modelData = NewModelASTData(modelName)
		(*modelsData)[modelName] = modelData
	}

	for mixin := range modelData.Mixins {
		if _, isPredefinedMixin := config.ModelMixins[mixin]; isPredefinedMixin {
			if _, exists := (*modelsData)[mixin]; !exists {
				fmt.Printf("[INFO] Mixin '%s' is predefined but missing. Initializing.\n", mixin)
				(*modelsData)[mixin] = NewModelASTData(mixin)
			}
		}

		inflateMixins(mixin, modelsData)

		mixinData, mixinExists := (*modelsData)[mixin]
		if !mixinExists {
			fmt.Printf("[ERROR] Mixin '%s' still does not exist in modelsData.\n", mixin)
			continue
		}

		// Copy fields and methods from mixin
		for fieldName, field := range mixinData.Fields {
			if fieldName == "ID" || modelData.Fields[fieldName].MixinField {
				continue // Skip ID field and already inherited fields
			}
			field.MixinField = true
			modelData.Fields[fieldName] = field
			fmt.Printf("[INFO] Added field '%s' from mixin '%s' to model '%s'\n", fieldName, mixin, modelName)
		}

		for methodName, method := range mixinData.Methods {
			if _, exists := modelData.Methods[methodName]; exists {
				continue // Skip existing methods
			}
			method.ToDeclare = true
			modelData.Methods[methodName] = method
			fmt.Printf("[INFO] Added method '%s' from mixin '%s' to model '%s'\n", methodName, mixin, modelName)
		}
	}

	modelData.Validated = true
	(*modelsData)[modelName] = modelData
}

// inflateEmbeds populates the given model with fields from the embedded type
// inflateEmbeds populates the given model with fields from the embedded type
func inflateEmbeds(modelName string, modelsData *map[string]ModelASTData) {
	fmt.Printf("[INFO] Inflating embeds for model: %s\n", modelName)

	modelData, exists := (*modelsData)[modelName]
	if !exists {
		fmt.Printf("[WARNING] Model '%s' does not exist in modelsData. Initializing it.\n", modelName)
		modelData = NewModelASTData(modelName)
		(*modelsData)[modelName] = modelData
	}

	for emb := range modelData.Embeds {
		relModel := modelData.Fields[emb].RelModel

		// Ensure the embedded model exists
		if _, exists := (*modelsData)[relModel]; !exists {
			fmt.Printf("[INFO] Embedded model '%s' does not exist. Initializing it.\n", relModel)
			(*modelsData)[relModel] = NewModelASTData(relModel)
		}

		inflateEmbeds(relModel, modelsData)

		// Copy fields from embedded model
		for fieldName, field := range (*modelsData)[relModel].Fields {
			if _, exists := modelData.Fields[fieldName]; exists {
				continue // Skip existing fields
			}
			field.EmbedField = true
			modelData.Fields[fieldName] = field
			fmt.Printf("[INFO] Added embedded field '%s' from model '%s' to model '%s'\n", fieldName, relModel, modelName)
		}
	}

	(*modelsData)[modelName] = modelData
}

// parseMixInModel updates the mixin tree with the given node which is a InheritModel function
func parseMixInModel(node *ast.CallExpr, modInfo *models.ModuleInfo, modelsData *map[string]ModelASTData) {
	fmt.Printf("[INFO] Parsing mixin model for node: %v", node)
	fNode := node.Fun.(*ast.SelectorExpr)
	modelName, err := extractModel(fNode.X, modInfo)
	if err != nil {
		if _, ok := err.(generalMixinError); ok {
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
func ExtractFunctionName(node *ast.CallExpr) (string, error) {
	switch nf := node.Fun.(type) {
	case *ast.SelectorExpr:
		return nf.Sel.Name, nil
	case *ast.Ident:
		return nf.Name, nil
	default:
		return "", errors.New("unexpected node type")
	}
}

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
	fmt.Printf("[INFO] Parsing add fields for model: %s", modInfo.PkgPath)
	fNode := node.Fun.(*ast.SelectorExpr)
	modelName, err := extractModel(fNode.X, modInfo)
	if err != nil {
		log.Panic("Unable to extract model while visiting AST", "error", err)
	}
	if _, exists := (*modelsData)[modelName]; !exists {
		(*modelsData)[modelName] = NewModelASTData(modelName)
	}
	fields := extractFieldCompositeLit(node.Args[0])
	for _, f := range fields.Elts {
		fDef := f.(*ast.KeyValueExpr)
		fieldName := strings.Trim(fDef.Key.(*ast.BasicLit).Value, "\"")
		fData := parseFieldDefinition(fDef, modInfo)
		if fData.embed {
			(*modelsData)[modelName].Embeds[fieldName] = true
		}
		(*modelsData)[modelName].Fields[fieldName] = fData
	}
}

// parseFieldDefinition parses the field definition for the given KeyValueExpr
func parseFieldDefinition(fDef *ast.KeyValueExpr, modInfo *models.ModuleInfo) FieldASTData {
	var typeStr string
	switch ft := fDef.Value.(*ast.CompositeLit).Type.(type) {
	case *ast.Ident:
		typeStr = strings.TrimSuffix(ft.Name, "Field")
	case *ast.SelectorExpr:
		typeStr = strings.TrimSuffix(ft.Sel.Name, "Field")
	}
	fType := fieldtype.Type(strings.ToLower(typeStr))
	return FieldASTData{
		Name:  strings.Trim(fDef.Key.(*ast.BasicLit).Value, "\""),
		FType: fType,
		Type: data.TypeData{
			Type:       fType.DefaultGoType().String(),
			ImportPath: getFieldImportPath(typeStr),
		},
	}
}

// extractFieldCompositeLit extracts the *ast.CompositeLit from the given argument
func extractFieldCompositeLit(arg ast.Expr) *ast.CompositeLit {
	switch n := arg.(type) {
	case *ast.CompositeLit:
		return n
	case *ast.Ident:
		return n.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.CompositeLit)
	}
	return nil
}

// getFieldImportPath returns the appropriate import path for the given field type
func getFieldImportPath(typeStr string) string {
	if typeStr == "Date" || typeStr == "DateTime" {
		return config.DatesPath
	}
	return ""
}

// parseNewModel parses the given node which is a NewXXXModel function

// setModelData adds a model with the given name and type to the given modelsData

// ExtractFunctionName returns the name of the called function
// in the given call expression.

// parseAddMethod parses the given node which is an addMethod function

// A generalMixinError is returned if the mixin is
// a general mixin set in NewXXXXModel function.
type generalMixinError struct{}

// Error method for generalMixinError
func (gme generalMixinError) Error() string {
	return "General Mixin Error"
}

// parseAddFields parses the given node which is an AddFields function

// parseStringValue returns the value of a string expr which can be a literal
// or an identifier for a string.
func parseStringValue(expr ast.Expr) string {
	var str string
	switch v := expr.(type) {
	case *ast.BasicLit:
		str = v.Value
	case *ast.Ident:
		str = parseStringValue(v.Obj.Decl.(*ast.ValueSpec).Values[0])
	}
	return strings.Trim(str, "\"")
}

// extractSelection returns a map with the keys and values of the Selection
// specified by expr.
func extractSelection(expr ast.Expr) map[string]string {
	res := make(map[string]string)
	switch e := expr.(type) {
	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			elem := elt.(*ast.KeyValueExpr)
			key := elem.Key.(*ast.BasicLit).Value
			value := strings.Trim(elem.Value.(*ast.BasicLit).Value, "\"")
			res[key] = value
		}
	}
	return res
}

// parseFieldAttribute parses the given KeyValueExpr of a field definition
func parseFieldAttribute(fElem *ast.KeyValueExpr, fData FieldASTData, modInfo *models.ModuleInfo) FieldASTData {
	switch fElem.Key.(*ast.Ident).Name {
	case "JSON":
		fData.JSON = parseStringValue(fElem.Value)
	case "Help":
		fData.Help = parseStringValue(fElem.Value)
	case "String":
		fData.Description = parseStringValue(fElem.Value)
	case "Selection":
		fData.Selection = extractSelection(fElem.Value)
	case "RelationModel":
		modName, err := extractModel(fElem.Value, modInfo)
		if err != nil {
			log.Panic("Unable to parse RelationModel", "field", fData.Name, "error", err)
		}
		fData.RelModel = modName
		fData.IsRS = true
	case "GoType":
		fData.Type = getTypeData(fElem.Value.(*ast.CallExpr).Args[0], modInfo)
	case "Embed":
		if fElem.Value.(*ast.Ident).Name == "true" {
			fData.embed = true
		}
	}
	return fData
}

// extractModel returns the string name of the model of the given ident variable
// ident must point to the expr which represents a model
// Returns an error if it cannot determine the model
func extractModel(ident ast.Expr, modInfo *models.ModuleInfo) (string, error) {
	switch idt := ident.(type) {
	case *ast.Ident:
		if decl, ok := idt.Obj.Decl.(*ast.AssignStmt); ok {
			switch rd := decl.Rhs[0].(type) {
			case *ast.CallExpr:
				var fnIdent *ast.Ident
				switch ft := rd.Fun.(type) {
				case *ast.Ident:
					fnIdent = ft
				case *ast.SelectorExpr:
					fnIdent = ft.Sel
				default:
					return "", fmt.Errorf("unexpected function identifier: %v (%T)", rd.Fun, rd.Fun)
				}
				switch fnIdent.Name {
				case "Get", "MustGet", "NewModel", "NewMixinModel", "NewTransientModel", "NewManualModel":
					return strings.Trim(rd.Args[0].(*ast.BasicLit).Value, "\""), nil
				case "CreateModel", "getOrCreateModel":
					return "", generalMixinError{}
				default:
					return extractModelNameFromFunc(rd, modInfo)
				}
			case *ast.Ident:
				return extractModel(rd, modInfo)
			default:
				return "", fmt.Errorf("unmanaged type %T at %s for %s", rd, modInfo.FSet.Position(rd.Pos()), idt.Name)
			}
		}
	case *ast.CallExpr:
		return extractModelNameFromFunc(idt, modInfo)
	default:
		return "", fmt.Errorf("unmanaged call. ident: %s (%T)", idt, idt)
	}
	return "", errors.New("unmanaged situation")
}

// extractModelNameFromFunc extracts the model name from a h.ModelName()
// expression or an error if this is not a pool function.
func extractModelNameFromFunc(ce *ast.CallExpr, modInfo *models.ModuleInfo) (string, error) {
	switch ft := ce.Fun.(type) {
	case *ast.Ident:
		// func is called without selector, then it is not from pool
		return "", errors.New("function call without selector")
	case *ast.SelectorExpr:
		switch ftt := ft.X.(type) {
		case *ast.Ident:
			if ftt.Name != config.PoolModelPackage && ftt.Name != "Registry" {
				return extractModel(ftt, modInfo)
			}
			return ft.Sel.Name, nil
		case *ast.CallExpr:
			return extractModel(ftt, modInfo)
		default:
			return "", fmt.Errorf("selector is of not managed type: %T", ftt)
		}
	}
	return "", errors.New("unparsable function call")
}

// extractParams extracts the parameters of the given FuncType
func extractParams(ft *ast.FuncType, modInfo *models.ModuleInfo) []data.ParamData {
	var params []data.ParamData
	for i, pl := range ft.Params.List {
		if i == 0 {
			// pass the first argument (rs)
			continue
		}
		for _, nn := range pl.Names {
			var variadic bool
			typ := pl.Type
			if el, ok := typ.(*ast.Ellipsis); ok {
				typ = el.Elt
				variadic = true
			}
			params = append(params, data.ParamData{
				Name:     nn.Name,
				Variadic: variadic,
				Type:     getTypeData(typ, modInfo)})
		}
	}
	return params
}

// getTypeData returns a data.TypeData instance representing the typ AST Expression
func getTypeData(typ ast.Expr, modInfo *models.ModuleInfo) data.TypeData {
	typStr := types.TypeString(modInfo.TypesInfo.TypeOf(typ), (*types.Package).Name)
	if strings.Contains(typStr, "invalid type") {
		// Maybe this is a pool type that is not yet defined
		byts := bytes.Buffer{}
		printer.Fprint(&byts, modInfo.FSet, typ)
		typStr = byts.String()
	}
	importPath := computeExportPath(modInfo.TypesInfo.TypeOf(typ))
	if strings.Contains(importPath, config.PoolPath) {
		importPath = ""
	}

	importPathTokens := strings.Split(importPath, ".")
	if len(importPathTokens) > 0 {
		importPath = strings.Join(importPathTokens[:len(importPathTokens)-1], ".")
	}
	return data.TypeData{
		Type:       typStr,
		ImportPath: importPath,
	}
}

// extractReturnType returns the return type of the first returned value
// of the given FuncType as a string and an import path if needed.
func extractReturnType(ft *ast.FuncType, modInfo *models.ModuleInfo) []data.TypeData {
	var res []data.TypeData
	if ft.Results != nil {
		for _, l := range ft.Results.List {
			res = append(res, getTypeData(l.Type, modInfo))
		}
	}
	return res
}

// computeExportPath returns the import path of the given type
func computeExportPath(typ types.Type) string {
	var res string
	switch typTyped := typ.(type) {
	case *types.Struct, *types.Named:
		res = types.TypeString(typTyped, (*types.Package).Path)
	case *types.Pointer:
		res = computeExportPath(typTyped.Elem())
	case *types.Slice:
		res = computeExportPath(typTyped.Elem())
	}
	return res
}

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

// defaultFields returns the map of default fields for the model with the given name
func defaultFields(name string) map[string]FieldASTData {
	res := make(map[string]FieldASTData)
	idField := FieldASTData{
		Name: "ID",
		JSON: "id",
		Type: data.TypeData{
			Type: "int64",
		},
		FType: fieldtype.Integer,
	}
	res["ID"] = idField
	switch name {
	case "BaseMixin":
		res["CreateDate"] = FieldASTData{
			Name:        "CreateDate",
			JSON:        "create_date",
			Description: "Created On",
			Type: data.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["CreateUID"] = FieldASTData{
			Name:        "CreateUID",
			JSON:        "create_uid",
			Description: "Created By",
			Type:        data.TypeData{Type: "int64"},
			FType:       fieldtype.Integer,
		}
		res["WriteDate"] = FieldASTData{
			Name:        "WriteDate",
			JSON:        "write_date",
			Description: "Updated On",
			Type: data.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["WriteUID"] = FieldASTData{
			Name:        "WriteUID",
			JSON:        "write_uid",
			Description: "Updated By",
			Type:        data.TypeData{Type: "int64"},
			FType:       fieldtype.Integer,
		}
		res["LastUpdate"] = FieldASTData{
			Name:        "LastUpdate",
			JSON:        "__last_update",
			Description: "Last Updated On",
			Type: data.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["DisplayName"] = FieldASTData{
			Name:        "DisplayName",
			JSON:        "display_name",
			Description: "Display Name",
			Type:        data.TypeData{Type: "string"},
			FType:       fieldtype.Char,
		}
	case "ModelMixin":
		res["HexyaExternalID"] = FieldASTData{
			Name:        "HexyaExternalID",
			JSON:        "hexya_external_id",
			Description: "External ID",
			Type:        data.TypeData{Type: "string"},
			FType:       fieldtype.Char,
		}
		res["HexyaVersion"] = FieldASTData{
			Name:        "HexyaVersion",
			JSON:        "hexya_version",
			Description: "External Version",
			Type:        data.TypeData{Type: "int"},
			FType:       fieldtype.Integer,
		}
	}
	return res
}
