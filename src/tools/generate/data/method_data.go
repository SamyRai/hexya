package data

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
}
