// generate/method_ast_data.go
package ast

import "github.com/hexya-erp/hexya/src/tools/generate/data"

// MethodASTData describes a method's AST data
type MethodASTData struct {
	Name      string
	Doc       string
	PkgPath   string
	Params    []data.ParamData
	Returns   []data.TypeData
	ToDeclare bool
}
