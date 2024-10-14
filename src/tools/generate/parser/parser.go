package parser

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

// inflateMixinsForAST populates the given model with fields and methods defined in its mixins for AST data.
func inflateMixinsForAST(modelName string, modelsData map[string]*models.ModelData) {
	model, exists := (modelsData)[modelName]
	if !exists || model == nil {
		fmt.Printf("Model %s not found or nil in modelsData", modelName)
		return
	}

	// Ensure mixins are initialized
	if len(model.Mixins) == 0 {
		fmt.Printf("Model %s has no mixins to inflate\n", modelName)
		return
	}
	for _, mixin := range (modelsData)[modelName].Mixins {
		// Recursively inflate mixins
		inflateMixinsForAST(mixin.Name, modelsData)

		// Add fields from mixins to the model
		for _, field := range mixin.Fields {
			// Skip adding the "ID" field
			if field.Name == "ID" {
				continue
			}
			// Set MixinField attribute if the field comes from a mixin
			field.Attributes.MixinField = true
			model := (modelsData)[modelName]
			if !model.HasField(field.Name) {
				model.Fields = append(model.Fields, field)
			}
		}

		// Add methods from mixins to the model
		for _, method := range mixin.Methods {
			method.ToDeclare = true //
			model := (modelsData)[modelName]

			for _, method := range mixin.Methods {
				method.ToDeclare = true
				if !model.HasMethod(method.Name) {
					model.Methods = append(model.Methods, method)
				}
			}

			// Reassign the modified model back to the map
			(modelsData)[modelName] = model // Mark method for declaration
		}
	}
}

// inflateEmbedsForAST populates the given model with fields from the embedded model types for AST data.
func inflateEmbedsForAST(modelName string, modelsData map[string]*models.ModelData) {
	model, exists := (modelsData)[modelName]
	if !exists || model == nil {
		fmt.Printf("Model %s not found or nil in modelsData", modelName)
		return
	}
	for _, embed := range (modelsData)[modelName].EmbeddedModels {
		// Recursively inflate embedded models
		inflateEmbedsForAST(embed.Name, modelsData)

		// Add fields from embedded models to the current model
		for _, field := range embed.Fields {
			field.Attributes.EmbedField = true // Mark field as an embedded field
			model := (modelsData)[modelName]
			if !model.HasField(field.Name) {
				model.Fields = append(model.Fields, field)
			}
		}
	}
}

// ParseModels parses the models, fields, and methods from the given packages, inflates mixins, and embeds fields/methods.
func ParseModels(packs []*packages.Package, validate bool) map[string]*models.ModelData {
	fmt.Println("Starting parsing phase...")

	// Initialize the parsers
	modelParser := NewModelParser()   // Assumes a ModelParser is defined elsewhere
	fieldParser := NewFieldParser()   // Assumes a FieldParser is defined elsewhere
	methodParser := NewMethodParser() // Assumes a MethodParser is defined elsewhere

	modelsData := make(map[string]*models.ModelData)

	for _, pack := range packs {
		for _, file := range pack.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.CallExpr:
					funcName, err := ExtractFunctionName(node)
					if err != nil {
						fmt.Printf("Failed to extract function name: %v\n", err)
						return true
					}

					// Delegate parsing based on the function name
					switch funcName {
					case "NewModel", "NewMixinModel", "NewTransientModel":
						modelParser.Parse(node, modelsData)
					case "AddFields":
						fieldParser.Parse(node, &models.ModuleInfo{Package: *pack}, modelsData)
					case "AddMethod":
						methodParser.Parse(node, modelsData, &models.ModuleInfo{Package: *pack})
					default:
						fmt.Printf("Unhandled function: %s\n", funcName)
					}
				default:
					fmt.Printf("Ignoring node type: %T\n", node)
				}
				return true
			})
		}
	}

	for modelName := range modelsData {
		inflateMixins(modelName, &modelsData)
		inflateEmbeds(modelName, &modelsData)
	}

	if validate {
		for modelName, modelData := range modelsData {
			if !modelData.Validated {
				delete(modelsData, modelName)
			}
		}
	}

	fmt.Printf("Total models parsed: %d\n", len(modelsData))
	return modelsData
}

// inflateMixins populates the given model with fields and methods defined in its mixins.
func inflateMixins(modelName string, modelsData *map[string]*models.ModelData) {
	for _, mixin := range (*modelsData)[modelName].Mixins {
		inflateMixins(mixin.Name, modelsData)

		for _, field := range mixin.Fields {
			if field.Name == "ID" {
				continue
			}
			field.Attributes.MixinField = true
			if !(*modelsData)[modelName].HasField(field.Name) {
				(*modelsData)[modelName].Fields = append((*modelsData)[modelName].Fields, field)
			}
		}

		for _, method := range mixin.Methods {
			method.ToDeclare = true
			if !(*modelsData)[modelName].HasMethod(method.Name) {
				(*modelsData)[modelName].Methods = append((*modelsData)[modelName].Methods, method)
			}
		}
	}
}

// inflateEmbeds populates the given model with fields from the embedded model types.
func inflateEmbeds(modelName string, modelsData *map[string]*models.ModelData) {
	for _, embed := range (*modelsData)[modelName].EmbeddedModels {
		inflateEmbeds(embed.Name, modelsData)

		for _, field := range embed.Fields {
			field.Attributes.EmbedField = true
			if !(*modelsData)[modelName].HasField(field.Name) {
				(*modelsData)[modelName].Fields = append((*modelsData)[modelName].Fields, field)
			}
		}
	}
}

// Utility function to extract function name from a CallExpr.
func ExtractFunctionName(node *ast.CallExpr) (string, error) {
	switch fun := node.Fun.(type) {
	case *ast.Ident:
		return fun.Name, nil
	case *ast.SelectorExpr:
		return fun.Sel.Name, nil
	default:
		return "", fmt.Errorf("unsupported function type: %T", fun)
	}
}
