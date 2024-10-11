// data/field_type.go
package models

// FieldType holds the name and valid operators on a field type
type FieldType struct {
	Type      string
	SanType   string
	IsRS      bool
	Operators []OperatorDef
}
