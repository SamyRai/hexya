package models

// A ParamData holds the name and type of a method parameter
type ParamData struct {
	Name     string
	Variadic bool
	Type     TypeData
}

type ParamASTData struct {
	Name     string
	Type     TypeInfo
	Variadic bool
}
