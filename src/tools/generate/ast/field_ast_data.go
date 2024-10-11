// generate/field_ast_data.go
package ast

import (
	"github.com/hexya-erp/hexya/src/models/fieldtype"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
)

type FieldASTData struct {
	Name        string
	JSON        string
	Help        string
	Description string
	Selection   map[string]string
	RelModel    string
	Type        data.TypeData
	FType       fieldtype.Type
	IsRS        bool
	MixinField  bool
	EmbedField  bool
	embed       bool
}
