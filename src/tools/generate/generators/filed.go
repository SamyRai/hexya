package generators

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
)

// FieldGenerator handles code generation for fields
type FieldGenerator struct{}

// NewFieldGenerator returns a new instance of FieldGenerator
func NewFieldGenerator() *FieldGenerator {
	return &FieldGenerator{}
}

// Generate creates all necessary files for the given ModelData (specifically fields)
func (g *FieldGenerator) Generate(modelData *models.ModelData, dir string) error {
	modelData.Sort() // Ensure fields are sorted for consistent output

	// Define the file path for fields in the model's directory
	fieldFilePath := filepath.Join(dir, "m", fmt.Sprintf("%s_fields.go", modelData.SnakeName()))

	// Create the fields file in the "m" directory
	if err := g.createFileFromTemplate(fieldFilePath, "FieldTemplate", modelData); err != nil {
		return fmt.Errorf("failed to create field file: %w", err)
	}

	return nil
}

// createFileFromTemplate generates a new file from the embedded template and data
func (g *FieldGenerator) createFileFromTemplate(fileName, templateName string, data *models.ModelData) error {
	// Get the template for fields
	tmpl, err := templates.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return fmt.Errorf("failed to create directory for file %s: %w", fileName, err)
	}

	// Execute the template and format the source code
	var srcBuffer bytes.Buffer
	if err := tmpl.Execute(&srcBuffer, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formattedSource, err := format.Source(srcBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format source code: %w", err)
	}

	// Write the formatted code to the file
	if err := os.WriteFile(fileName, formattedSource, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
