package file_operations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/logging"
	"golang.org/x/tools/go/packages"
)

var log = logging.GetLogger("file_operations")

// ComputeDirs returns the absolute project and pool directory paths.
func ComputeDirs(projectDir string) (string, string, error) {
	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", "", fmt.Errorf("error computing absolute project directory: %w", err)
	}
	poolDir := filepath.Join(absProjectDir, "pool")
	return absProjectDir, poolDir, nil
}

// CleanPoolDir removes all contents of the pool directory and recreates the necessary subdirectories.
func CleanPoolDir(poolDir string) error {
	log.Info(fmt.Sprintf("Cleaning pool directory: %s", poolDir))
	if err := os.RemoveAll(poolDir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", poolDir, err)
	}

	log.Info("Pool directory cleaned and recreated successfully.")
	return nil
}

// CreateEmptyPool sets up an empty pool structure (folders and temp files) in the pool directory.
func CreateEmptyPool(poolDir string) error {
	log.Info(fmt.Sprintf("Creating empty pool structure in poolDir: %s", poolDir))

	subDirs := []string{
		filepath.Join(poolDir, config.PoolModelPackage),      // `h`
		filepath.Join(poolDir, config.PoolQueryPackage),      // `q`
		filepath.Join(poolDir, config.PoolInterfacesPackage), // `m`
	}

	for _, subDir := range subDirs {
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", subDir, err)
		}
		if err := createTempGoFile(filepath.Join(subDir, config.TempEmpty), filepath.Base(subDir)); err != nil {
			return fmt.Errorf("error creating temp.go file in %s: %w", subDir, err)
		}
	}

	log.Info("Empty pool structure created successfully.")
	return nil
}

// createTempGoFile creates a `temp.go` file to prevent directory removal.
func createTempGoFile(filePath, packageName string) error {
	templateData := struct {
		PackageName string
	}{
		PackageName: packageName,
	}

	tmpl, err := templates.GetTemplate("EmptyPoolTemplate")
	if err != nil {
		return fmt.Errorf("failed to get temp.go template: %w", err)
	}

	return templates.CreateFileFromTemplate(filePath, tmpl, templateData)
}

// LoadProgram loads Go packages for the provided target paths.
func LoadProgram(targetPaths []string, tests bool) ([]*packages.Package, error) {
	fmt.Printf("Loading packages for target paths: %v\n", targetPaths)
	conf := packages.Config{
		Mode:  packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: tests,
	}

	packs, err := packages.Load(&conf, targetPaths...)
	return packs, err
}

// CreateModuleSymlinks sets up symlinks for the provided modules in the project directory.
func CreateModuleSymlinks(modules []*models.ModuleInfo, projectDir string) {
	CleanModuleSymlinks(projectDir)

	for _, module := range modules {
		if module.ModType != models.Base {
			continue
		}
		createModuleSymlinks(module, projectDir)
	}
}

// createModuleSymlinks creates symlinks for individual modules.
func createModuleSymlinks(mod *models.ModuleInfo, projectDir string) {
	for _, dir := range config.SymlinkDirs {
		mDir := filepath.Dir(mod.GoFiles[0])
		srcPath := filepath.Join(mDir, dir)
		dstPath := filepath.Join(projectDir, config.ResDirRel, dir)

		if _, err := os.Stat(srcPath); err != nil {
			continue
		}

		if err := os.MkdirAll(dstPath, 0755); err != nil {
			panic(err)
		}

		if err := os.Symlink(srcPath, filepath.Join(dstPath, mod.Name)); err != nil {
			panic(err)
		}
	}
}

// CleanModuleSymlinks removes all existing symlinks in the project symlink directories.
func CleanModuleSymlinks(projectDir string) {
	for _, dir := range config.SymlinkDirs {
		dirPath := filepath.Join(projectDir, config.ResDirRel, dir)
		os.RemoveAll(dirPath)
		os.Mkdir(dirPath, 0775)
	}
}

// CreatePoolFiles generates the pool files (interface, model, query) for the given model data.
// It now uses DepsManager for handling dependencies.
// CreatePoolFiles generates the pool files (interface, model, query) for the given model data.
// CreatePoolFiles generates the pool files (interface, model, query) for the given model data.
// It now uses DepsManager for handling dependencies.
func CreatePoolFiles(dir string, mData *models.ModelData) error {
	fmt.Printf("Sorting model data for model: %s\n", mData.Name)
	fmt.Printf("Dependencies: %+v\n", mData.Deps)
	mData.Sort()

	// Create all required files using templates
	files := []struct {
		Path         string
		TemplateName string
	}{
		{filepath.Join(dir, config.PoolInterfacesPackage, fmt.Sprintf("%s.go", mData.SnakeName)), "PoolInterfacesTemplate"},
		{filepath.Join(dir, config.PoolModelPackage, fmt.Sprintf("%s.go", mData.SnakeName)), "PoolModelsTemplate"},
		{filepath.Join(dir, config.PoolModelPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName)), "PoolModelsDirTemplate"},
		{filepath.Join(dir, config.PoolQueryPackage, fmt.Sprintf("%s.go", mData.SnakeName)), "PoolQueryTemplate"},
		{filepath.Join(dir, config.PoolQueryPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName)), "PoolModelsQueryTemplate"},
	}

	for _, file := range files {
		fmt.Printf("Creating file: %s using template: %s\n", file.Path, file.TemplateName)

		// Gather imports dynamically based on the template being used
		var imports []string
		switch file.TemplateName {
		case "PoolInterfacesTemplate":
			imports = templates.ImportsForPoolInterfacesTemplate(mData)
		case "PoolModelsTemplate":
			imports = templates.ImportsForPoolModelsTemplate(mData)
		case "PoolModelsQueryTemplate":
			imports = templates.ImportsForPoolModelsQueryTemplate(mData)
		case "PoolQueryTemplate":
			imports = templates.ImportsForPoolQueryTemplate(mData)
		case "PoolModelsDirTemplate":
			imports = templates.ImportsForPoolModelsDirTemplate(mData)
		}

		if err := createFile(file.Path, file.TemplateName, mData, imports); err != nil {
			return fmt.Errorf("failed to create file: %s: %w", file.Path, err)
		}
	}

	fmt.Printf("Completed creating pool files for model: %s\n", mData.Name)
	return nil
}

// createFile is a helper function to create files using a specified template.
// createFile is a helper function to create files using a specified template.
func createFile(path, templateName string, mData *models.ModelData, imports []string) error {
	tmpl, err := templates.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template %s: %w", templateName, err)
	}

	// Combine imports into the model data
	templateData := struct {
		*models.ModelData
		Imports []string
	}{
		ModelData: mData,
		Imports:   imports,
	}

	if err := templates.CreateFileFromTemplate(path, tmpl, templateData); err != nil {
		fmt.Printf("Error creating file from template %s: %v\n", templateName, err)
		return fmt.Errorf("failed to create file from template %s: %w", templateName, err)
	}
	return nil
}

// CreateMainFile generates the main.go file with core and module imports in the specified directory.
func CreateMainFile(dir string, coreImports, moduleImports []string, executableName string) error {
	templateData := struct {
		Executable  string
		CoreImports []string
		Modules     []string
	}{
		Executable:  executableName,
		CoreImports: coreImports,
		Modules:     moduleImports,
	}

	tmpl, err := templates.GetTemplate("MainFileTemplate")
	if err != nil {
		return fmt.Errorf("failed to get main.go template: %w", err)
	}

	mainFilePath := filepath.Join(dir, "main.go")
	if err := templates.CreateFileFromTemplate(mainFilePath, tmpl, templateData); err != nil {
		return fmt.Errorf("failed to create main.go file: %w", err)
	}

	log.Info(fmt.Sprintf("main.go file created successfully in directory: %s", dir))
	return nil
}
