// data/model_data.go
package models

import (
	"sort"
)

// ModelData describes a RecordSet model
type ModelData struct {
	Name                  string
	SnakeName             string
	ModelsPackageName     string
	QueryPackageName      string
	InterfacesPackageName string
	ModelType             string
	IsModelMixin          bool
	RelModels             []string
	Fields                []FieldData
	Methods               []MethodData
	AllMethods            []MethodData
	ConditionFuncs        []string
	Types                 []FieldType
	Imports               []string // Add Imports field to collect all imports
}

// Sort sorts all slices fields of this ModelData so that the generated code is always the same.
func (m *ModelData) Sort() {
	sort.Slice(m.Fields, func(i, j int) bool {
		return m.Fields[i].Name < m.Fields[j].Name
	})
	sort.Slice(m.Methods, func(i, j int) bool {
		return m.Methods[i].Name < m.Methods[j].Name
	})
	sort.Slice(m.AllMethods, func(i, j int) bool {
		return m.AllMethods[i].Name < m.AllMethods[j].Name
	})
	sort.Strings(m.RelModels)
	sort.Slice(m.Types, func(i, j int) bool {
		return m.Types[i].Type < m.Types[j].Type
	})
}
