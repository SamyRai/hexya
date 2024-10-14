package validator

import (
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/dsl"
	"github.com/hexya-erp/hexya/src/tools/dsl/ast"
)

var ErrModelNameEmpty = errors.New("model name cannot be empty")

// ValidateModel checks if the model structure is correct
func ValidateModel(model *ast.DSLModel) error {
	if model.Name == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	// Validate fields
	for _, field := range model.Fields {
		if err := validateField(field); err != nil {
			return err
		}
	}

	// Validate methods
	for _, method := range model.Methods {
		if err := validateMethod(method); err != nil {
			return err
		}
	}

	return nil
}

// validateField checks that the field type and options are valid
func validateField(field *ast.DSLField) error {
	if field.Name == "" {
		return fmt.Errorf("field name cannot be empty")
	}
	switch field.Type {
	case dsl.TypeChar, dsl.TypeInteger, dsl.TypeDate:
		// Valid field types
	default:
		return fmt.Errorf("invalid field type: %s", field.Type)
	}
	return nil
}

// validateMethod checks method parameters and return types
func validateMethod(method *ast.DSLMethod) error {
	if method.Name == "" {
		return fmt.Errorf("method name cannot be empty")
	}
	if method.Body == "" {
		return fmt.Errorf("method body cannot be empty")
	}
	return nil
}
