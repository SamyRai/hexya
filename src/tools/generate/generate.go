package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"
	"github.com/hexya-erp/hexya/src/tools/generate/parser/builders"
	"github.com/hexya-erp/hexya/src/tools/strutils"
)

// CreatePool generates the pool package from parsed AST data.
func CreatePool(modelsASTData map[string]parser.ModelASTData, dir string) error {
	fmt.Println("Starting pool creation... with models:", len(modelsASTData))

	for _, modelASTData := range modelsASTData {
		fmt.Println("Model:", modelASTData.Name)
	}

	// Process predefined mixins first.
	if err := processPredefinedMixins(modelsASTData, dir); err != nil {
		return fmt.Errorf("failed to process predefined mixins: %w", err)
	}

	// Process the remaining regular models.
	for modelName, modelASTData := range modelsASTData {
		if isMixin(modelName) {
			continue // Skip already processed mixins.
		}

		if !modelASTData.Validated {
			fmt.Printf("Skipping unvalidated model: %s\n", modelName)
			continue
		}

		if err := processModel(modelName, modelASTData, modelsASTData, dir); err != nil {
			return fmt.Errorf("failed to process model %s: %w", modelName, err)
		}
	}

	fmt.Println("Pool creation completed successfully.")
	return nil
}

// processPredefinedMixins processes the predefined mixins in the configuration.
func processPredefinedMixins(modelsASTData map[string]parser.ModelASTData, dir string) error {
	for mixinName := range config.ModelMixins {
		modelASTData, exists := modelsASTData[mixinName]
		if !exists {
			// Initialize empty predefined mixins if they do not exist in the AST data.
			modelASTData = initializePredefinedMixin(mixinName)
		}

		if err := processModel(mixinName, modelASTData, modelsASTData, dir); err != nil {
			return fmt.Errorf("failed to process predefined mixin %s: %w", mixinName, err)
		}
	}
	return nil
}

// initializePredefinedMixin initializes an empty mixin if not present in the AST data.
func initializePredefinedMixin(mixinName string) parser.ModelASTData {
	fmt.Printf("Initializing predefined mixin: %s\n", mixinName)
	return parser.ModelASTData{
		Name:         mixinName,
		IsModelMixin: true,
		Fields:       map[string]models.FieldASTData{},
		Methods:      map[string]models.MethodASTData{},
		Mixins:       map[string]bool{},
		Embeds:       map[string]bool{},
		Validated:    true,
	}
}

// isMixin checks if the given model name is a predefined mixin.
func isMixin(modelName string) bool {
	_, exists := config.ModelMixins[modelName]
	return exists
}

// processModel processes a single model, adding fields, methods, and validating dependencies.
func processModel(modelName string, modelASTData parser.ModelASTData, modelsASTData map[string]parser.ModelASTData, dir string) error {
	mData := prepareModelData(modelName, modelASTData)

	// Process and validate the model fields, types, and methods.
	processModelFields(modelASTData, &mData)
	processModelTypes(&mData)
	processModelMethods(modelsASTData, &mData)

	// Write model data to pool files.
	if err := file_operations.CreatePoolFiles(dir, &mData); err != nil {
		return fmt.Errorf("failed to create pool files for model %s: %w", modelName, err)
	}
	return nil
}

// prepareModelData initializes and prepares the model data structure.
func prepareModelData(modelName string, modelASTData parser.ModelASTData) models.ModelData {
	return models.ModelData{
		Name:                  modelName,
		SnakeName:             strutils.SnakeCase(modelName),
		ModelsPackageName:     config.PoolModelPackage,
		QueryPackageName:      config.PoolQueryPackage,
		InterfacesPackageName: config.PoolInterfacesPackage,
		ModelType:             modelASTData.ModelType,
		IsModelMixin:          modelASTData.IsModelMixin,
		ConditionFuncs:        config.ConditionFuncs,
		Deps:                  []string{},
	}
}

// processModelFields processes fields of the model and updates dependencies.
func processModelFields(modelASTData parser.ModelASTData, mData *models.ModelData) {
	builders.AddFieldsToModelData(modelASTData, mData)
}

// processModelTypes processes field types and updates dependencies.
func processModelTypes(mData *models.ModelData) {
	builders.AddFieldTypesToModelData(mData)
}

// processModelMethods processes methods of the model and updates dependencies.
func processModelMethods(modelsASTData map[string]parser.ModelASTData, mData *models.ModelData) {
	builders.AddMethodsToModelData(modelsASTData, mData)
}
