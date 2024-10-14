package templates

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

//go:embed tmpl/*
var templateFS embed.FS

var templatesMap = map[string]*template.Template{}

// templateFiles holds paths to the individual template files
var templateFiles = map[string]string{
	"PoolTemplate":          "tmpl/poolTemplate.tmpl",
	"MainFileTemplate":      "tmpl/mainFileTemplate.tmpl",
	"EmptyPoolTemplate":     "tmpl/emptyPoolTemplate.tmpl",
	"StartFileTemplateI18n": "tmpl/startFileTemplateI18n.tmpl",
}

// LoadTemplates loads all templates from the embedded FS into templatesMap
func LoadTemplates() error {
	for name, file := range templateFiles {
		tmpl, err := loadTemplate(file)
		if err != nil {
			return fmt.Errorf("failed to load template %s: %w", name, err)
		}
		templatesMap[name] = tmpl
	}
	return nil
}

// GetTemplate returns a compiled template by its name
func GetTemplate(name string) (*template.Template, error) {
	tmpl, ok := templatesMap[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}

// loadTemplate loads and parses a template from the embedded FS
func loadTemplate(path string) (*template.Template, error) {
	tmplData, err := templateFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", path, err)
	}
	tmpl, err := template.New(filepath.Base(path)).Parse(string(tmplData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
	}
	return tmpl, nil
}

// CreateFileFromTemplate generates a file from the provided template and data
func CreateFileFromTemplate(fileName string, tmpl *template.Template, data interface{}) error {
	var buffer bytes.Buffer

	// Execute the template with the provided data
	if err := tmpl.Execute(&buffer, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write the output to the specified file
	if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return fmt.Errorf("failed to create directory for file %s: %w", fileName, err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer file.Close()

	if _, err := file.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", fileName, err)
	}

	return nil
}
func PoolViewObjectConstructor(
	mData *models.ModelData,
) map[string]interface{} {
	// Extract the necessary dynamic imports
	dynamicImports := mData.GetImports()

	// Add hardcoded imports for models package
	hardcodedImports := []string{
		`"github.com/hexya-erp/hexya/src/models"`,
	}

	// Combine dynamic and hardcoded imports
	imports := append(hardcodedImports, dynamicImports...)

	return map[string]interface{}{
		"Name":         mData.Name,
		"SnakeName":    mData.SnakeName(),
		"Fields":       mData.GetFields(),
		"Methods":      mData.GetMethods(),
		"Imports":      imports,
		"ModelType":    mData.ModelType,    // Type of model: Base, Transient, Mixin, etc.
		"IsModelMixin": mData.IsModelMixin, // True if it's a Mixin model
	}
}
func InterfaceViewObjectConstructor(
	mData *models.ModelData,
) map[string]interface{} {
	// Gather the dynamic imports required for interfaces
	dynamicImports := mData.GetImports()

	// Add hardcoded imports for models package
	hardcodedImports := []string{
		`"github.com/hexya-erp/hexya/src/models"`,
	}

	// Combine dynamic and hardcoded imports
	imports := append(hardcodedImports, dynamicImports...)

	return map[string]interface{}{
		"Name":      mData.Name,
		"SnakeName": mData.SnakeName(),
		"Fields":    mData.GetFields(),
		"Methods":   mData.GetMethods(),
		"Imports":   imports,
	}
}

func QueryViewObjectConstructor(
	mData *models.ModelData,
) map[string]interface{} {
	// Collect dynamic imports from the model data
	dynamicImports := mData.GetImports()

	// Add hardcoded imports for query processing
	hardcodedImports := []string{
		`"github.com/hexya-erp/hexya/src/models/operator"`,
		`"github.com/hexya-erp/hexya/src/models"`,
	}

	// Combine dynamic and hardcoded imports
	imports := append(hardcodedImports, dynamicImports...)

	return map[string]interface{}{
		"Name":           mData.Name,
		"SnakeName":      mData.SnakeName(),
		"Fields":         mData.GetFields(),
		"ConditionFuncs": []string{"And", "Or", "Not"}, // Hardcoded condition functions
		"Imports":        imports,
	}
}

func MainFileViewObjectConstructor(
	executableName string,
	coreImports []string,
	modules []string,
) map[string]interface{} {
	// Hardcode core imports (e.g., logging, ORM, etc.)
	hardcodedImports := []string{
		`"log"`,
		`"github.com/hexya-erp/hexya/src/models"`,
	}

	// Combine core imports with module imports
	allImports := append(hardcodedImports, coreImports...)
	allImports = append(allImports, modules...)

	return map[string]interface{}{
		"Executable":  executableName, // e.g., "hexya"
		"CoreImports": allImports,     // Combine core and module imports
	}
}
