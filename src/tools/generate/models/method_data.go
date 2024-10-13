package models

// MethodData describes a method in a RecordSet
type MethodData struct {
	Name             string
	Doc              string
	Params           string
	ParamsWithType   string
	IParamsWithTypes string
	ParamsTypes      string
	ReturnAsserts    string
	Returns          string
	ReturnString     string
	IReturnString    string
	Call             string
	ToDeclare        bool
	ImportPaths      []string
}

// MethodASTData describes a method's AST data
type MethodASTData struct {
	Name      string
	Doc       string
	PkgPath   string
	Params    []ParamData
	Returns   []TypeData
	ToDeclare bool
}
