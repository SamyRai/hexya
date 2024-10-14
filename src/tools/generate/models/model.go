package models

import (
	"fmt"
	"strings"
)

// ModelData represents the full structure of a model, including fields, methods, mixins, and dependencies.
type ModelData struct {
	Name           string           // Model name.
	ModelType      string           // Type of model (e.g., Base, Transient).
	IsModelMixin   bool             // Is this model a mixin?
	Fields         []*FieldAST      // List of fields within the model.
	Methods        []*MethodAST     // List of methods within the model.
	EmbeddedModels []*ModelData     // Models embedded in this one.
	Mixins         []*ModelData     // Mixins this model inherits.
	Dependencies   []DependencyNode // List of dependencies for this model.
	Validated      bool             // Validation status of the model.

	ProcessedFields  []FieldData
	ProcessedMethods []MethodData
	ProcessedImports []string
}

// ModelReference defines a reference to another model (for fields or relations).
type ModelReference struct {
	ModelName  string
	ImportPath string
}

// HasField checks if a field already exists in the model.
func (m *ModelData) HasField(fieldName string) bool {
	fieldName = strings.ToLower(fieldName) // Normalize case
	for _, field := range m.Fields {
		if strings.ToLower(field.Name) == fieldName {
			return true
		}
	}
	return false
}

// HasMethod checks if a method already exists in the model.
func (m *ModelData) HasMethod(methodName string) bool {
	methodName = strings.ToLower(methodName) // Normalize case
	for _, method := range m.Methods {
		if strings.ToLower(method.Name) == methodName {
			return true
		}
	}
	return false
}

// InflateMixins inflates mixin fields and methods into the current model.
func (m *ModelData) InflateMixins(visited map[string]bool) {
	fmt.Printf("Inflating mixins for model: %s\n", m.Name)

	if visited[m.Name] {
		return // Prevent infinite recursion
	}
	visited[m.Name] = true
	for _, mixin := range m.Mixins {
		for _, mixinField := range mixin.Fields {
			if !m.HasField(mixinField.Name) {
				m.Fields = append(m.Fields, mixinField)
			}
		}
		for _, mixinMethod := range mixin.Methods {
			if !m.HasMethod(mixinMethod.Name) {
				m.Methods = append(m.Methods, mixinMethod)
			}
		}
	}
}

// InflateEmbeds inflates embedded model fields and methods into the current model.
func (m *ModelData) InflateEmbeds() {
	for _, embed := range m.EmbeddedModels {
		for _, embedField := range embed.Fields {
			if !m.HasField(embedField.Name) {
				m.Fields = append(m.Fields, embedField)
			}
		}
		for _, embedMethod := range embed.Methods {
			if !m.HasMethod(embedMethod.Name) {
				m.Methods = append(m.Methods, embedMethod)
			}
		}
	}
}

func (m *ModelData) collectImports(imports map[string]bool) {
	// Collect imports from fields
	for _, field := range m.Fields {
		if field.Type.ImportPath != "" {
			imports[field.Type.ImportPath] = true
		}
	}

	// Collect imports from methods
	for _, method := range m.Methods {
		for _, imp := range method.GetImportPaths() {
			imports[imp] = true
		}
	}

	// Collect imports from mixins
	for _, mixin := range m.Mixins {
		mixin.collectImports(imports)
	}
}

func (m *ModelData) GetImports() []string {
	if m.ProcessedImports != nil {
		return m.ProcessedImports
	}

	imports := make(map[string]bool)
	m.collectImports(imports)

	m.ProcessedImports = removeDuplicateImports(imports)
	return m.ProcessedImports
}

// GetFields processes and returns all fields in FieldData format.
func (m *ModelData) GetFields() []FieldData {
	if m.ProcessedFields != nil {
		return m.ProcessedFields // Return cached result if already processed.
	}
	var result []FieldData
	for _, field := range m.Fields {
		result = append(result, field.ToFieldData())
	}
	m.ProcessedFields = result // Cache the result.
	return result
}

// GetMethods processes and returns all methods in MethodData format.
func (m *ModelData) GetMethods() []MethodData {
	if m.ProcessedMethods != nil {
		return m.ProcessedMethods // Return cached result if already processed.
	}
	var result []MethodData
	for _, method := range m.Methods {
		result = append(result, method.ToMethodData())
	}
	m.ProcessedMethods = result // Cache the result.
	return result
}
