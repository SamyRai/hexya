package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
	"go/printer"
	"go/types"
	"strings"
)

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

// extractModel extracts the model name from the given AST node.
func extractModel(ident ast.Expr, modInfo *models.ModuleInfo) (string, error) {
	switch idt := ident.(type) {
	case *ast.Ident:
		// Method is called on an identifier without selector such as
		// user.addMethod. In this case, we try to find out the model from
		// the identifier declaration.
		switch decl := idt.Obj.Decl.(type) {
		case *ast.AssignStmt:
			// The declaration is also an assignment
			switch rd := decl.Rhs[0].(type) {
			case *ast.CallExpr:
				// The assignment is a call to a function
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
					return strings.Trim(rd.Args[0].(*ast.BasicLit).Value, "\"`"), nil
				case "CreateModel", "getOrCreateModel":
					// This is a call from inside a NewXXXXModel function
					return "", generalMixinError{}
				default:
					return extractModelNameFromFunc(rd, modInfo)
				}
			case *ast.Ident:
				// The assignment is another identifier, we go to the declaration of this new ident.
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

// getTypeData returns a data.TypeData instance representing the typ AST Expression
func getTypeData(typ ast.Expr, modInfo *models.ModuleInfo) models.TypeData {
	typStr := types.TypeString(modInfo.TypesInfo.TypeOf(typ), (*types.Package).Name)
	if strings.Contains(typStr, "invalid type") {
		var byts bytes.Buffer
		err := printer.Fprint(&byts, modInfo.FSet, typ)
		if err != nil {
			log.Panic("Unable to print type", "error", err)
		}
		typStr = byts.String()
	}
	importPath := computeExportPath(modInfo.TypesInfo.TypeOf(typ))
	return models.TypeData{
		Type:       typStr,
		ImportPath: importPath,
	}
}

// computeExportPath returns the import path of the given type.
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

// extractReturnType returns the return type of the first returned value
// of the given FuncType as a string and an import path if needed.
func extractReturnType(ft *ast.FuncType, modInfo *models.ModuleInfo) []models.TypeData {
	var res []models.TypeData
	if ft.Results != nil {
		for _, l := range ft.Results.List {
			res = append(res, getTypeData(l.Type, modInfo))
		}
	}
	return res
}

// defaultFields returns the map of default fields for the model with the given name
func defaultFields(name string) map[string]models.FieldASTData {
	res := make(map[string]models.FieldASTData)
	idField := models.FieldASTData{
		Name: "ID",
		JSON: "id",
		Type: models.TypeData{
			Type: "int64",
		},
		FType: fieldtype.Integer,
	}
	res["ID"] = idField
	switch name {
	case "BaseMixin":
		res["CreateDate"] = models.FieldASTData{
			Name:        "CreateDate",
			JSON:        "create_date",
			Description: "Created On",
			Type: models.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["CreateUID"] = models.FieldASTData{
			Name:        "CreateUID",
			JSON:        "create_uid",
			Description: "Created By",
			Type:        models.TypeData{Type: "int64"},
			FType:       fieldtype.Integer,
		}
		res["WriteDate"] = models.FieldASTData{
			Name:        "WriteDate",
			JSON:        "write_date",
			Description: "Updated On",
			Type: models.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["WriteUID"] = models.FieldASTData{
			Name:        "WriteUID",
			JSON:        "write_uid",
			Description: "Updated By",
			Type:        models.TypeData{Type: "int64"},
			FType:       fieldtype.Integer,
		}
		res["LastUpdate"] = models.FieldASTData{
			Name:        "LastUpdate",
			JSON:        "__last_update",
			Description: "Last Updated On",
			Type: models.TypeData{
				Type:       "dates.DateTime",
				ImportPath: config.DatesPath,
			},
			FType: fieldtype.DateTime,
		}
		res["DisplayName"] = models.FieldASTData{
			Name:        "DisplayName",
			JSON:        "display_name",
			Description: "Display Name",
			Type:        models.TypeData{Type: "string"},
			FType:       fieldtype.Char,
		}
	case "ModelMixin":
		res["HexyaExternalID"] = models.FieldASTData{
			Name:        "HexyaExternalID",
			JSON:        "hexya_external_id",
			Description: "External ID",
			Type:        models.TypeData{Type: "string"},
			FType:       fieldtype.Char,
		}
		res["HexyaVersion"] = models.FieldASTData{
			Name:        "HexyaVersion",
			JSON:        "hexya_version",
			Description: "External Version",
			Type:        models.TypeData{Type: "int"},
			FType:       fieldtype.Integer,
		}
	}
	return res
}

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
func parseFieldAttribute(fElem *ast.KeyValueExpr, fData models.FieldASTData, modInfo *models.ModuleInfo) models.FieldASTData {
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
			fData.Embed = true
		}
	}
	return fData
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
func extractParams(ft *ast.FuncType, modInfo *models.ModuleInfo) []models.ParamData {
	var params []models.ParamData
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
			params = append(params, models.ParamData{
				Name:     nn.Name,
				Variadic: variadic,
				Type:     getTypeData(typ, modInfo)})
		}
	}
	return params
}
