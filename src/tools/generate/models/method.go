// models/methods.go
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

// MethodData is the processed representation of a model method.
type MethodData struct {
	Name        string   // Method name
	Params      string   // Method parameters as a string
	Returns     string   // Return types as a string
	ImportPaths []string // Import paths for dependencies
	Source      string   // Source of the method
}

// ParamAST defines the structure of a method parameter.
type ParamAST struct {
	Name       string  // Parameter name
	Type       TypeAST // Parameter type information
	IsVariadic bool    // Indicates if the parameter is variadic
}

// ReturnAST defines the structure of a method return type.
type ReturnAST struct {
	Type       TypeAST // Type information for the return type
	ImportPath string  // Import path for the return type, if any
}

// GetImportPaths returns a list of import paths from method dependencies.
func (m *MethodAST) GetImportPaths() []string {
	var importPaths []string
	for _, dep := range m.Dependencies {
		importPaths = append(importPaths, dep.ImportPath)
	}
	return importPaths
}
