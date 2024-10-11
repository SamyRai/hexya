package templates

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed tmpl/*
var templateFS embed.FS

var templatesMap = map[string]*template.Template{}

var templateFiles = map[string]string{
	"PoolInterfacesTemplate":  "tmpl/poolInterfacesTemplate.tmpl",
	"PoolModelsTemplate":      "tmpl/poolModelsTemplate.tmpl",
	"PoolModelsDirTemplate":   "tmpl/poolModelsDirTemplate.tmpl",
	"PoolQueryTemplate":       "tmpl/poolQueryTemplate.tmpl",
	"PoolModelsQueryTemplate": "tmpl/poolModelsQueryTemplate.tmpl",
	"StartFileTemplateI18n":   "tmpl/startFileTemplateI18n.tmpl",
	"MainFileTemplate":        "tmpl/mainFileTemplate.tmpl",
	"GoModTemplate":           "tmpl/goModTemplate.tmpl", // Consolidated go.mod template
	"EmptyPoolTemplate":       "tmpl/emptyPoolTemplate.tmpl",
}

// LoadTemplates loads all templates from the embedded FS.
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

// GetTemplate returns a template by its name.
func GetTemplate(name string) (*template.Template, error) {
	tmpl, ok := templatesMap[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}

// loadTemplate loads an individual template from the embedded FS.
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

// CreateFileFromTemplate generates a file from the provided template and data.
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
