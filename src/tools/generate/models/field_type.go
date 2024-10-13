// data/field_type.go
package models

import (
	"github.com/hexya-erp/hexya/src/models/fieldtype"
)

// FieldType holds the name and valid operators on a field type
type FieldType struct {
	Type      string
	SanType   string
	IsRS      bool
	Operators []OperatorDef
}

type FieldASTData struct {
	Name        string
	JSON        string
	Help        string
	Description string
	Selection   map[string]string
	RelModel    string
	Type        TypeData
	FType       fieldtype.Type
	IsRS        bool
	MixinField  bool
	EmbedField  bool
	Embed       bool
}
