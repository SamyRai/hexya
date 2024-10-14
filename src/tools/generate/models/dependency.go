package models

// DependencyNode represents a dependency for models, methods, or fields (used for import paths).
type DependencyNode struct {
	ImportPath     string
	DependencyKind DependencyKind // Field, Method, or Type dependency.
}

// DependencyKind is an enum defining the type of dependency (field, method, etc.).
type DependencyKind string

const (
	FieldDependency  DependencyKind = "field"
	MethodDependency DependencyKind = "method"
	TypeDependency   DependencyKind = "type"
)
