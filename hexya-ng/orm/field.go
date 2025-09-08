package orm

// FieldInfo holds the definition of a model field.
type FieldInfo struct {
	Name     string
	Type     string // e.g., "Char", "Integer", "Many2One"
	String   string
	Help     string
	Required bool
	// Add other field attributes here as needed
}
