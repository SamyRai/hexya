package models

// MethodAST defines a method within a model, along with parameters, returns, and dependencies.
type MethodAST struct {
	Name         string           // Method name.
	Doc          string           // Documentation for the method.
	Params       []ParamAST       // Parameters for the method.
	Returns      []ReturnAST      // Return types for the method.
	Dependencies []DependencyNode // Dependencies for the method (e.g., import paths).
	ToDeclare    bool             // Whether the method should be declared in the model.
}

// GetImportPaths returns import paths from method dependencies.
func (m *MethodAST) GetImportPaths() []string {
	var importPaths []string
	for _, dep := range m.Dependencies {
		importPaths = append(importPaths, dep.ImportPath)
	}
	return importPaths
}

// ToMethodData converts MethodAST into MethodData for processing or generation.
func (m *MethodAST) ToMethodData() MethodData {
	return MethodData{
		Name:        m.Name,
		Params:      formatParams(m.Params),
		Returns:     formatReturns(m.Returns),
		ImportPaths: getImportPathsFromDependencies(m.Dependencies),
	}
}

// MethodData is the processed representation of a model method.
type MethodData struct {
	Name        string
	Params      string
	Returns     string
	ImportPaths []string
	Source      string
}

// ParamAST defines the structure of a method parameter.
type ParamAST struct {
	Name       string
	Type       TypeAST
	IsVariadic bool
}

// ReturnAST defines the structure of a method return type.
type ReturnAST struct {
	Type       TypeAST
	ImportPath string
}
