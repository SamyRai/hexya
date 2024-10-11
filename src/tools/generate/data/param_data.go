package data

// A ParamData holds the name and type of a method parameter
type ParamData struct {
	Name     string
	Variadic bool
	Type     TypeData
}
