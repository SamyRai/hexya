package templates

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"os"
	"path/filepath"
	"strings"
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

// ImportsForPoolInterfacesTemplate generates a list of imports required for the PoolInterfaces template.
func ImportsForPoolInterfacesTemplate(mData *models.ModelData) []string {
	var imports []string

	// Always include the pool query package import
	imports = append(imports, fmt.Sprintf("\"github.com/hexya-erp/pool/%s\"", mData.QueryPackageName))

	return removeDuplicateImports(imports)
}

// ImportsForPoolQueryTemplate generates a list of imports required for the PoolQuery template.
func ImportsForPoolQueryTemplate(mData *models.ModelData) []string {
	imports := []string{
		`"github.com/hexya-erp/hexya/src/models"`,
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s/%s"`, mData.QueryPackageName, mData.SnakeName),
	}

	return removeDuplicateImports(imports)
}

// ImportsForPoolModelsTemplate generates a list of imports required for the PoolModels template.
func ImportsForPoolModelsTemplate(mData *models.ModelData) []string {
	imports := []string{
		`"github.com/hexya-erp/hexya/src/models"`,
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s"`, mData.QueryPackageName),
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s/%s"`, mData.ModelsPackageName, mData.SnakeName),
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s"`, mData.InterfacesPackageName),
	}

	return removeDuplicateImports(imports)
}

// ImportsForPoolModelsDirTemplate generates a list of imports required for the PoolModelsDir template.
func ImportsForPoolModelsDirTemplate(mData *models.ModelData) []string {
	imports := []string{
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s"`, mData.QueryPackageName),
		fmt.Sprintf(`"github.com/hexya-erp/pool/%s"`, mData.InterfacesPackageName),
	}

	for _, f := range mData.Fields {
		if f.ImportPath != "" {
			imports = append(imports, fmt.Sprintf(`"%s"`, getPackagePathFromAST(f.ImportPath)))
		}
	}

	for _, d := range mData.Deps {
		if d != "" {
			imports = append(imports, fmt.Sprintf(`"%s"`, getPackagePathFromAST(d)))
		}
	}

	for _, r := range mData.Methods {
		if r.ParamsTypes != "" {
			fmt.Printf("\n\n\n!!!Method: %v\n", r.ParamsTypes)
		}
	}

	return removeDuplicateImports(imports)
}

// ImportsForPoolModelsQueryTemplate generates a list of imports required for the PoolModelsQuery template.
func ImportsForPoolModelsQueryTemplate(mData *models.ModelData) []string {
	imports := []string{
		`"github.com/hexya-erp/hexya/src/models/operator"`,
		`"github.com/hexya-erp/hexya/src/models"`,
	}

	for _, f := range mData.Fields {
		if f.ImportPath != "" {
			imports = append(imports, fmt.Sprintf(`"%s"`, getPackagePathFromAST(f.ImportPath)))
		}
	}

	return removeDuplicateImports(imports)
}

// removeDuplicateImports ensures that the imports list contains only unique entries.
func removeDuplicateImports(imports []string) []string {
	seen := make(map[string]bool)
	uniqueImports := []string{}

	for _, imp := range imports {
		if !seen[imp] {
			seen[imp] = true
			uniqueImports = append(uniqueImports, imp)
		}
	}
	return uniqueImports
}

func getPackagePathFromAST(importPath string) string {
	// Split by `/` to get each part of the path and avoid confusion with dots
	parts := strings.Split(importPath, "/")

	// Check the last part for any type references (contains ".")
	lastPart := parts[len(parts)-1]
	if strings.Contains(lastPart, ".") {
		// It's a type reference, so we remove everything after the last dot
		parts[len(parts)-1] = lastPart[:strings.LastIndex(lastPart, ".")]
	}

	// Rejoin the parts to form the package path
	packagePath := strings.Join(parts, "/")

	// Return the final package path
	return packagePath
}
