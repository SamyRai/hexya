package orm

import "fmt"

// ModelInfo holds the definition of a model.
type ModelInfo struct {
	Name   string
	Fields map[string]*FieldInfo
}

// Registry is the registry of all models.
var Registry = make(map[string]*ModelInfo)

// NewModel creates a new model and adds it to the registry.
func NewModel(name string) *ModelInfo {
	if _, ok := Registry[name]; ok {
		panic(fmt.Sprintf("model %s already exists", name))
	}
	model := &ModelInfo{
		Name:   name,
		Fields: make(map[string]*FieldInfo),
	}
	Registry[name] = model
	return model
}

// AddFields adds the given fields to the model.
func (m *ModelInfo) AddFields(fields map[string]*FieldInfo) {
	for name, field := range fields {
		m.Fields[name] = field
	}
}

// GetModel returns the model with the given name from the registry.
func GetModel(name string) *ModelInfo {
	model, ok := Registry[name]
	if !ok {
		panic(fmt.Sprintf("model %s not found", name))
	}
	return model
}
