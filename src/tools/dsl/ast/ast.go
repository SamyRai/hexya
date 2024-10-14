package ast

// DSLModel represents the structure of a model in Hexya DSL.
type DSLModel struct {
	Name        string           // Name of the model (e.g., "User")
	Fields      []*DSLField      // List of fields associated with the model
	Methods     []*DSLMethod     // List of methods associated with the model
	Relations   []*DSLRelation   // List of relationships (e.g., Many2One, One2Many)
	Validations []*DSLValidation // List of validation functions (e.g., ValidateAge)
	Routes      []*DSLRoute      // List of custom routes
	Mixins      []*DSLMixin      // Mixins that extend this model
}

// DSLField represents a field definition in a model (e.g., Name: Char)
type DSLField struct {
	Name       string // Field name
	Type       string // Field type (e.g., Char, Integer, Date)
	IsRequired bool   // Whether this field is required
}

// DSLMethod represents a method definition for a model (e.g., GreetUser: return fmt.Sprintf)
type DSLMethod struct {
	Name   string   // Method name
	Body   string   // Method body as a string
	Params []string // List of parameters for the method
}

// DSLRelation defines the relationship between models (e.g., Many2One, One2Many)
type DSLRelation struct {
	FieldName    string // Name of the field (e.g., "Customer")
	RelatedModel string // Name of the related model (e.g., "CustomerModel")
	RelationType string // Type of relation (Many2One, One2Many, etc.)
}

// DSLValidation represents validation logic for a model field or condition
type DSLValidation struct {
	Name string // Validation name (e.g., ValidateAge)
	Body string // Validation logic (e.g., if Age < 18 { return false })
}

// DSLRoute defines a route for the model (e.g., /greet => GreetUser)
type DSLRoute struct {
	Route   string // The URL path (e.g., "/greet")
	Handler string // The method to handle the route (e.g., GreetUser)
}

// DSLMixin defines a mixin (inheritance or behavior reuse) in a model
type DSLMixin struct {
	MixinName string // Name of the mixin (e.g., BaseModel)
}
