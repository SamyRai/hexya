package parser

import "github.com/hexya-erp/hexya/src/tools/generate/models"

// ModelASTData holds the AST data of a model
type ModelASTData struct {
	Name         string
	ModelType    string
	IsModelMixin bool
	Fields       map[string]models.FieldASTData
	Methods      map[string]models.MethodASTData
	Mixins       map[string]bool
	Embeds       map[string]bool
	Validated    bool
}
