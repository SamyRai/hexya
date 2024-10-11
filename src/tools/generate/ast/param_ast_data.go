// generate/param_ast_data.go
package ast

import "github.com/hexya-erp/hexya/src/tools/generate/models"

// ParamASTData describes a parameter's AST data
type ParamASTData struct {
	Name     string
	Type     models.TypeInfo
	Variadic bool
}
