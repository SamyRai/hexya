package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/generators"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"
	"golang.org/x/tools/go/packages"
)

// Pool handles the generation of models based on parsed data.
func Pool(modelsData map[string]*models.ModelData, dir string) error {
	fmt.Println("Starting generation phase...")

	// Inflate models (add mixins, embedded models)
	fmt.Println("Inflating models (mixins and embedded models)...")
	models.InflateModels(modelsData)

	// Process predefined mixins
	fmt.Println("Processing predefined mixins...")
	if err := generators.ProcessPredefinedMixins(modelsData, dir); err != nil {
		return err
	}

	// Process each validated model
	for modelName, modelData := range modelsData {
		if !modelData.Validated {
			fmt.Printf("Skipping unvalidated model: %s\n", modelName)
			continue
		}
		fmt.Printf("Processing validated model: %s\n", modelName)
		if err := generators.ProcessModel(modelName, modelData, dir); err != nil {
			return err
		}
	}

	fmt.Println("Generation phase completed successfully.")
	return nil
}

// CreatePool orchestrates the entire pool creation process by first parsing and then generating the pool.
func CreatePool(packs []*packages.Package, dir string) error {
	// Convert packages to ModuleInfo objects
	modules := models.ConvertPackagesToModules(packs)

	// Phase 1: Parsing
	modelsData := parser.GetModelsASTDataForModules(modules, true) // You can choose whether to validate or not

	// Parse models
	modelsData = parser.ParseModels(packs, true)

	// Phase 2: Generation
	if err := Pool(modelsData, dir); err != nil {
		return err
	}

	fmt.Println("Pool creation completed successfully.")
	return nil
}
