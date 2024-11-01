package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/builders"
	"github.com/hexya-erp/hexya/src/tools/generate/generators"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"
)

// Pool processes and generates model data based on parsed model information.
func Pool(modelDataMap map[string]*models.ModelData, outputDir string) error {
	fmt.Println("Starting generation phase...")

	// Inflate models by adding mixins and embedded models
	fmt.Println("Inflating models (mixins and embedded models)...")
	models.InflateModels(modelDataMap)

	// Process predefined mixins if applicable
	fmt.Println("Processing predefined mixins...")
	if err := generators.ProcessPredefinedMixins(modelDataMap, outputDir); err != nil {
		return err
	}

	// Process each model in the model data map
	for modelName, modelData := range modelDataMap {
		fmt.Printf("Processing model: %s\n", modelName)

		// Initialize a map to store dependencies for the current model
		dependencyMap := make(map[string]bool)

		// Populate fields and methods for the model data
		builders.AddFieldsToModelData(modelData)
		builders.AddMethodsToModelData(modelData, modelData, &dependencyMap)

		// Generate code for the processed model
		if err := generators.ProcessModel(modelName, modelData, outputDir); err != nil {
			return err
		}
	}

	fmt.Println("Generation phase completed successfully.")
	return nil
}

// CreatePool orchestrates the parsing and generation phases to create the pool.
func CreatePool(moduleInfoList []*models.ModuleInfo, outputDir string) error {

	// Phase 1: Retrieve initial AST data with mixins/embedded models
	initialASTData := parser.GetModelsASTDataForModules(moduleInfoList, true)

	// Phase 2: Convert AST data into fully populated model data structures
	modelDataMap := parser.ParseModels(initialASTData)

	// Phase 3: Generate code based on the parsed model data
	if err := Pool(modelDataMap, outputDir); err != nil {
		return err
	}

	fmt.Println("Pool creation completed successfully.")
	return nil
}
