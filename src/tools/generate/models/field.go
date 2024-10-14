package models

// FieldAttributes holds metadata for fields (e.g., mixin, embedded, selection).
type FieldAttributes struct {
	MixinField  bool              // Marks if the field is from a mixin.
	EmbedField  bool              // Marks if the field is embedded from another model.
	Help        string            // Help or documentation for the field.
	Description string            // Description of the field.
	Selection   map[string]string // Selection options (for enums, etc.).
	Required    bool
	Default     string
	ReadOnly    bool
	Index       bool
	Size        int
}

// TypeAST represents the type of fields and method parameters.
type TypeAST struct {
	TypeName    string // Go type (e.g., string, int64).
	ImportPath  string // Import path, if needed.
	IsRecordSet bool   // Marks if it's a RecordSet field.
}

// FieldAST defines a field within a model, along with its type and attributes.
type FieldAST struct {
	Name          string
	Type          TypeAST         // Type of the field.
	RelationModel *ModelReference // Reference to the related model (for relationships).
	Attributes    FieldAttributes // Additional metadata for the field.
}

// ToFieldData converts FieldAST into FieldData for processing or code generation.
func (f *FieldAST) ToFieldData() FieldData {
	return FieldData{
		Name:       f.Name,
		Type:       f.Type.TypeName,
		ImportPath: f.Type.ImportPath,
		IsRS:       f.RelationModel != nil,
		RelModel:   f.RelationModel.ModelName,
		MixinField: f.Attributes.MixinField,
		EmbedField: f.Attributes.EmbedField,
	}
}

// FieldData is the processed representation of a model field.
type FieldData struct {
	Name       string
	Type       string
	ImportPath string
	IsRS       bool
	RelModel   string
	MixinField bool
	EmbedField bool
}
