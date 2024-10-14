// generators/mixin.go
package generators

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// ProcessPredefinedMixins handles predefined mixins from the configuration.
func ProcessPredefinedMixins(modelsData map[string]*models.ModelData, dir string) error {
	for mixinName := range config.ModelMixins {
		modelData, exists := modelsData[mixinName]
		if !exists {
			modelData = InitializePredefinedMixin(mixinName)
			modelsData[mixinName] = modelData
		}
		if err := ProcessModel(mixinName, modelData, dir); err != nil {
			return fmt.Errorf("failed to process predefined mixin %s: %w", mixinName, err)
		}
	}
	return nil
}

// InitializePredefinedMixin initializes an empty predefined mixin.
func InitializePredefinedMixin(mixinName string) *models.ModelData {
	return &models.ModelData{
		Name:         mixinName,
		IsModelMixin: true,
		Fields:       []*models.FieldAST{},
		Methods:      []*models.MethodAST{},
		Validated:    true,
	}
}
