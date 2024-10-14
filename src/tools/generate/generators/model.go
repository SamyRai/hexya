package generators

import (
	"fmt"
	"path/filepath"

	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
)

// ModelGenerator handles code generation for models
type ModelGenerator struct{}

// NewModelGenerator returns a new instance of ModelGenerator
func NewModelGenerator() *ModelGenerator {
	return &ModelGenerator{}
}

// Generate creates all necessary files for the given ModelData
func (g *ModelGenerator) Generate(modelData *models.ModelData, dir string) error {
	modelData.Sort() // Ensure all fields and methods are sorted for consistent output

	// Use the ModelData directly, no need to pass all fields explicitly
	modelViewData := templates.PoolViewObjectConstructor(modelData)

	// Load required templates
	modelTemplate, err := templates.GetTemplate("PoolTemplate")
	if err != nil {
		return fmt.Errorf("failed to load PoolTemplate: %w", err)
	}
	interfaceTemplate, err := templates.GetTemplate("PoolTemplate") // Replace with correct template if needed
	if err != nil {
		return fmt.Errorf("failed to load PoolInterfacesTemplate: %w", err)
	}
	queryTemplate, err := templates.GetTemplate("PoolTemplate") // Replace with correct template if needed
	if err != nil {
		return fmt.Errorf("failed to load PoolQueryTemplate: %w", err)
	}

	// Define file paths
	interfaceFilePath := filepath.Join(dir, config.PoolInterfacesPackage, fmt.Sprintf("%s.go", modelData.SnakeName()))
	modelFilePath := filepath.Join(dir, config.PoolModelPackage, fmt.Sprintf("%s.go", modelData.SnakeName()))
	queryFilePath := filepath.Join(dir, config.PoolQueryPackage, fmt.Sprintf("%s.go", modelData.SnakeName()))

	// Generate all files using the template references
	if err := templates.CreateFileFromTemplate(interfaceFilePath, interfaceTemplate, modelViewData); err != nil {
		return fmt.Errorf("failed to create interface file: %w", err)
	}
	if err := templates.CreateFileFromTemplate(modelFilePath, modelTemplate, modelViewData); err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	if err := templates.CreateFileFromTemplate(queryFilePath, queryTemplate, modelViewData); err != nil {
		return fmt.Errorf("failed to create query file: %w", err)
	}

	return nil
}
