package parser

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// A generalMixinError is returned if the mixin is
// a general mixin set in NewXXXXModel function.
type generalMixinError struct{}

// Error method for generalMixinError
func (gme generalMixinError) Error() string {
	return "General Mixin Error"
}

var _ error = generalMixinError{}

func extractModel(ident ast.Expr) (string, error) {
	switch idt := ident.(type) {
	case *ast.Ident:
		// Check if the identifier declaration exists before accessing it.
		if idt.Obj == nil || idt.Obj.Decl == nil {
			return "", fmt.Errorf("identifier %s is not properly declared or initialized", idt.Name)
		}

		fmt.Printf("Ident: %s\n", idt.Name)

		// Method is called on an identifier without selector such as user.addMethod.
		// In this case, we try to find out the model from the identifier declaration.
		switch decl := idt.Obj.Decl.(type) {
		case *ast.AssignStmt:
			// The declaration is also an assignment
			if len(decl.Rhs) == 0 {
				return "", fmt.Errorf("no right-hand side in assignment for %s", idt.Name)
			}

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
					return extractModelNameFromFunc(rd)
				}
			case *ast.Ident:
				// The assignment is another identifier, we go to the declaration of this new ident.
				return extractModel(rd)
			default:
				return "", fmt.Errorf("unmanaged type %T at %s", rd, idt.Name)
			}
		}
	case *ast.CallExpr:
		return extractModelNameFromFunc(idt)
	default:
		return "", fmt.Errorf("unmanaged call. ident: %s (%T)", idt, idt)
	}
	return "", fmt.Errorf("unmanaged situation")
}

func extractModelNameFromFunc(ce *ast.CallExpr) (string, error) {
	switch ft := ce.Fun.(type) {
	case *ast.Ident:
		// func is called without selector, then it is not from pool
		return "", fmt.Errorf("function call without selector")
	case *ast.SelectorExpr:
		switch ftt := ft.X.(type) {
		case *ast.Ident:
			if ftt.Name != config.PoolModelPackage && ftt.Name != "Registry" {
				return extractModel(ftt)
			}
			return ft.Sel.Name, nil
		case *ast.CallExpr:
			return extractModel(ftt)
		default:
			return "", fmt.Errorf("selector is of not managed type: %T", ftt)
		}
	}
	return "", fmt.Errorf("unparsable function call")
}

// Helper function to extract the type of parameter or return value as a string.
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

// Helper function to extract the import path from an expression.
func getImportPathFromExpr(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return getImportPathFromExpr(x.X) + "." + x.Sel.Name
	}
	return ""
}

// Helper function to determine if the type is a RecordSet.
func isRecordSetType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return strings.HasSuffix(sel.Sel.Name, "Set")
	}
	return false
}

// Helper function to extract selection from a CompositeLit.
func extractSelection(expr ast.Expr) map[string]string {
	selection := make(map[string]string)
	switch e := expr.(type) {
	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			kv := elt.(*ast.KeyValueExpr)
			key := strings.Trim(kv.Key.(*ast.BasicLit).Value, "\"`")
			value := strings.Trim(kv.Value.(*ast.BasicLit).Value, "\"`")
			selection[key] = value
		}
	}
	return selection
}

// extractBoolValue handles boolean extraction from an AST expression
func extractBoolValue(expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name == "true"
	}
	return false
}

// extractIntValue handles integer extraction from an AST expression
func extractIntValue(expr ast.Expr) int {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.INT {
			intValue, err := strconv.Atoi(v.Value)
			if err == nil {
				return intValue
			}
		}
	}
	return 0
}

// CreateTypeIdent creates a string from the given type that can be used inside an identifier.
func CreateTypeIdent(typStr string) string {
	res := strings.Replace(typStr, ".", "", -1)
	res = strings.Replace(res, "[", "Slice", -1)
	res = strings.Replace(res, "map[", "Map", -1)
	res = strings.Replace(res, "]", "", -1)
	res = CapitalizeFirst(res)
	return res
}

// CapitalizeFirst capitalizes the first letter of the given string.
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
