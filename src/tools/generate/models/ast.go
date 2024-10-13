package models

// ModelAST represents the data structure for a model
type ModelAST struct {
	Name            string
	Type            string           // Type of model (e.g., TransientModel, BaseModel)
	IsModelMixin    bool             // True if this is a model mixin
	Fields          []FieldAST       // List of fields for the model
	Methods         []MethodAST      // List of methods associated with the model
	Mixins          []ModelReference // Models that this model inherits or mixes in
	Dependencies    []DependencyNode // List of dependencies required by this model
	Validated       bool             // Has the model been validated for use?
	EmbedReferences []ModelReference // Embedded models
}

// FieldAST represents the structure of a field inside a model
type FieldAST struct {
	Name          string
	Type          TypeAST         // The type of the field
	RelationModel *ModelReference // Reference to the model that this field relates to (if applicable)
	Attributes    FieldAttributes // Additional attributes such as selection, help, etc.
}

// MethodAST represents the structure of a method within a model
type MethodAST struct {
	Name          string
	Doc           string           // Documentation of the method
	Params        []ParamAST       // Parameters for the method
	Returns       []ReturnAST      // Return types for the method
	RelatedModels []ModelReference // Models this method interacts with
	Dependencies  []DependencyNode // Dependencies required for the method
}

// TypeAST represents type information about fields and method parameters
type TypeAST struct {
	TypeName    string
	ImportPath  string
	IsRecordSet bool
}

// ParamAST represents a parameter of a method
type ParamAST struct {
	Name       string
	Type       TypeAST
	IsVariadic bool
}

// ReturnAST represents return data from a method
type ReturnAST struct {
	Type       TypeAST
	ImportPath string
}

// ModelReference represents a reference to another model
type ModelReference struct {
	ModelName  string
	ImportPath string
}

// DependencyNode represents a dependency for models, methods, or fields
type DependencyNode struct {
	ImportPath     string
	DependencyKind DependencyKind // Field, Method, or Type dependency
}

// DependencyKind defines what kind of dependency a node is (field, method, type, etc.)
type DependencyKind string

const (
	FieldDependency  DependencyKind = "field"
	MethodDependency DependencyKind = "method"
	TypeDependency   DependencyKind = "type"
)
