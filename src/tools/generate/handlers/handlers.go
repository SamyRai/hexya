package handlers

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// ProcessSpecificMethod handles specific method logic based on the method's AST structure.
func ProcessSpecificMethod(astData *models.MethodAST, modelData *models.ModelData) {
	// Check method name and apply the corresponding logic
	switch astData.Name {
	case "Search":
		searchMethodHandler(astData, modelData)
	case "SearchByName":
		searchByNameMethodHandler(astData, modelData)
	case "Create":
		createMethodHandler(astData, modelData)
	case "New":
		newMethodHandler(astData, modelData)
	case "Write":
		writeMethodHandler(astData, modelData)
	case "Copy":
		copyMethodHandler(astData, modelData)
	case "CopyData":
		copyDataMethodHandler(astData, modelData)
	case "CartesianProduct":
		cartesianProductMethodHandler(astData, modelData)
	case "Sorted":
		sortedMethodHandler(astData, modelData)
	case "Filtered":
		filteredMethodHandler(astData, modelData)
	case "Aggregates":
		aggregatesMethodHandler(astData, modelData)
	case "First":
		firstMethodHandler(astData, modelData)
	case "All":
		allMethodHandler(astData, modelData)
	case "DefaultGet":
		defaultGetMethodHandler(astData, modelData)
	default:
		fmt.Printf("Unhandled specific method: %s\n", astData.Name)
	}
}

// searchMethodHandler processes the Search method using structured AST data.
func searchMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the Search method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Search",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// searchByNameMethodHandler processes the SearchByName method using structured AST data.
func searchByNameMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return string type for this method
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Append method to model's Methods
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name: "SearchByName",
		Doc: fmt.Sprintf(`// SearchByName searches for %s records that have a display name matching the given
		// "name" pattern when compared with the given "op" operator, while also
		// matching the optional search condition ("additionalCond").
		//
		// This is used for example to provide suggestions based on a partial
		// value for a relational field.`, modelData.Name),
		Params: []models.ParamAST{
			{Name: "name", Type: models.TypeAST{TypeName: "string"}},
			{Name: "op", Type: models.TypeAST{TypeName: "operator.Operator"}},
			{Name: "additionalCond", Type: models.TypeAST{TypeName: fmt.Sprintf("%s.%sCondition", "PoolQueryPackage", modelData.Name)}},
			{Name: "limit", Type: models.TypeAST{TypeName: "int"}},
		},
		Returns: []models.ReturnAST{
			{Type: returnType},
		},
		Dependencies: astData.Dependencies, // Collect dependencies dynamically
	})
}

// createMethodHandler processes the Create method.
func createMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the Create method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Create",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// newMethodHandler processes the New method.
func newMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the New method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "New",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// writeMethodHandler processes the Write method.
func writeMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as boolean
	returnType := models.TypeAST{
		TypeName:    "bool",
		ImportPath:  "",
		IsRecordSet: false,
	}

	// Create and append the Write method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Write",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// copyMethodHandler processes the Copy method.
func copyMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the Copy method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Copy",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// copyDataMethodHandler processes the CopyData method.
func copyDataMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as model data
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sData", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: false,
	}

	// Create and append the CopyData method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "CopyData",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// cartesianProductMethodHandler processes the CartesianProduct method.
func cartesianProductMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as an array of RecordSets
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("[]%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the CartesianProduct method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "CartesianProduct",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// sortedMethodHandler processes the Sorted method.
func sortedMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the Sorted method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Sorted",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// filteredMethodHandler processes the Filtered method.
func filteredMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as a RecordSet
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sSet", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: true,
	}

	// Create and append the Filtered method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Filtered",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// aggregatesMethodHandler processes the Aggregates method.
func aggregatesMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as an array of group aggregate rows
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("[]%sGroupAggregateRow", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: false,
	}

	// Create and append the Aggregates method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "Aggregates",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// firstMethodHandler processes the First method.
func firstMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as model data
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sData", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: false,
	}

	// Create and append the First method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "First",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// allMethodHandler processes the All method.
func allMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as an array of model data
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("[]%sData", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: false,
	}

	// Create and append the All method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "All",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}

// defaultGetMethodHandler processes the DefaultGet method.
func defaultGetMethodHandler(astData *models.MethodAST, modelData *models.ModelData) {
	// Define the return type as model data
	returnType := models.TypeAST{
		TypeName:    fmt.Sprintf("%sData", modelData.Name),
		ImportPath:  "PoolInterfacesPackage", // Adjust to actual package path
		IsRecordSet: false,
	}

	// Create and append the DefaultGet method
	modelData.Methods = append(modelData.Methods, &models.MethodAST{
		Name:         "DefaultGet",
		Params:       astData.Params,                         // Reuse parsed parameters
		Returns:      []models.ReturnAST{{Type: returnType}}, // Structured return type
		Dependencies: astData.Dependencies,                   // Collect dependencies dynamically
	})
}
