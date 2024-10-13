// data/field_data.go
package models

// FieldData describes a field in a RecordSet
type FieldData struct {
	Name        string
	JSON        string
	RelModel    string
	Type        string
	IType       string
	TypeWrapper string
	SanType     string
	ImportPath  string // Holds the import path for the field type
	IsRS        bool
	MixinField  bool
	EmbedField  bool
}
