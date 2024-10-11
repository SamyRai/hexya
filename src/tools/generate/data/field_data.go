package data

// FieldData describes a field in a RecordSet
type FieldData struct {
	Name        string
	JSON        string
	RelModel    string
	Type        string
	IType       string
	TypeWrapper string
	SanType     string
	ImportPath  string
	IsRS        bool
	MixinField  bool
	EmbedField  bool
}
