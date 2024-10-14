package generators

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"path/filepath"

	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
)

// MethodGenerator handles code generation for methods
type MethodGenerator struct{}

// NewMethodGenerator returns a new instance of MethodGenerator
func NewMethodGenerator() *MethodGenerator {
	return &MethodGenerator{}
}

// Generate creates all necessary files for the given MethodData
func (g *MethodGenerator) Generate(modelData *models.ModelData, dir string) error {
	modelData.Sort() // Ensure methods are sorted for consistent output

	// Define the file path for methods in the model's directory
	methodFilePath := filepath.Join(dir, config.PoolModelPackage, fmt.Sprintf("%s_methods.go", modelData.SnakeName()))

	// Use PoolViewObjectConstructor for creating view objects for templates
	viewData := templates.PoolViewObjectConstructor(modelData)
	tmpl, err := templates.GetTemplate("PoolTemplate")
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}
	// Generate the method file using the common helper
	if err := templates.CreateFileFromTemplate(methodFilePath, tmpl, viewData); err != nil {
		return fmt.Errorf("failed to create method file: %w", err)
	}

	return nil
}
