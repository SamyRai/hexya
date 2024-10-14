package generators

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// ProcessModel handles the conversion of ModelData and generates the corresponding pool files.
func ProcessModel(modelName string, modelData *models.ModelData, dir string) error {
	// Build and validate model fields and methods
	if err := BuildModelData(modelData); err != nil {
		return fmt.Errorf("failed to build model data for %s: %w", modelName, err)
	}

	// Generate pool files using the generator
	modelGenerator := NewModelGenerator()
	if err := modelGenerator.Generate(modelData, dir); err != nil {
		return fmt.Errorf("failed to create pool files for model %s: %w", modelName, err)
	}

	return nil
}

// BuildModelData processes and validates the model data.
func BuildModelData(modelData *models.ModelData) error {
	modelData.Validated = true
	modelData.Sort()
	return nil
}
