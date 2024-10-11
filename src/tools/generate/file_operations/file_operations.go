package file_operations

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/logging"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
)

var log = logging.GetLogger("file_operations")

// ComputeDirs computes the project and pool directory paths.
func ComputeDirs(projectDir string) (string, string, error) {
	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", "", fmt.Errorf("error computing absolute project directory: %w", err)
	}
	poolDir := filepath.Join(absProjectDir, "pool")
	return absProjectDir, poolDir, nil
}

// CleanPoolDir cleans the given pool directory.
func CleanPoolDir(poolDir string) error {
	log.Info(fmt.Sprintf("Cleaning pool directory: %s", poolDir))
	if err := os.RemoveAll(poolDir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", poolDir, err)
	}

	subDirs := []string{"models", "query", "interfaces"}
	for _, subDir := range subDirs {
		if err := os.MkdirAll(filepath.Join(poolDir, subDir), 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", subDir, err)
		}
	}

	log.Info("Pool directory cleaned successfully")
	return nil
}

// DetermineTargetPaths returns the target paths based on whether the test flag is set.
func DetermineTargetPaths(projectDir string, testEnabled bool) []string {
	if testEnabled {
		return []string{projectDir}
	}
	return viper.GetStringSlice("Modules")
}

// CreateEmptyPool creates an empty pool structure in the pool directory.
// This function no longer handles Go module generation.
// Instead, it should just ensure the required folder structure is set up.
func CreateEmptyPool(poolDir string) error {
	log.Info(fmt.Sprintf("Creating empty pool in poolDir: %s", poolDir))
	subDirs := []string{
		filepath.Join(poolDir, config.PoolModelPackage),      // `h`
		filepath.Join(poolDir, config.PoolQueryPackage),      // `q`
		filepath.Join(poolDir, config.PoolInterfacesPackage), // `m`
	}

	for _, subDir := range subDirs {
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", subDir, err)
		}

		// Create `temp.go` to ensure directories aren't removed
		tempFilePath := filepath.Join(subDir, config.TempEmpty)
		if err := createTempGoFile(tempFilePath, filepath.Base(subDir)); err != nil {
			return fmt.Errorf("error creating temp.go file in %s: %w", subDir, err)
		}
	}

	log.Info("Empty pool structure created successfully")
	return nil
}

// createTempGoFile creates a `temp.go` file in the given directory.
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

// CreateMainFile generates a main.go file in the project directory.

// CreateMainFile generates a main.go file in the project directory.
// CreateMainFile generates a main.go file in the project directory.

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

func LoadProgram(targetPaths []string, tests bool) ([]*packages.Package, error) {
	conf := packages.Config{
		Mode: packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: tests,
	}
	packs, err := packages.Load(&conf, targetPaths...)
	return packs, err
}

// createModuleSymlinks create the symlinks of the given module in the
// project directory.
// CreateModuleSymlinks creates symlinks for the given modules in the project directory.
func CreateModuleSymlinks(modules []*models.ModuleInfo, projectDir string) {
	CleanModuleSymlinks(projectDir)

	for _, module := range modules {
		if module.ModType != models.Base {
			continue
		}
		createModuleSymlinks(module, projectDir)
	}
}

// createModuleSymlinks creates the symlinks of the given module in the project directory.
func createModuleSymlinks(mod *models.ModuleInfo, projectDir string) {

	for _, dir := range config.SymlinkDirs {
		mDir := filepath.Dir(mod.GoFiles[0])
		srcPath := filepath.Join(mDir, dir)
		dstPath := filepath.Join(projectDir, config.ResDirRel, dir)
		if _, err := os.Stat(srcPath); err != nil {
			// Subdir doesn't exist, so we don't symlink
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

// cleanModuleSymlinks removes all symlinks in the server symlink directories.
// Note that this function actually removes and recreates the symlink directories.
func CleanModuleSymlinks(projectDir string) {
	for _, dir := range config.SymlinkDirs {
		dirPath := filepath.Join(projectDir, config.ResDirRel, dir)
		os.RemoveAll(dirPath)
		os.Mkdir(dirPath, 0775)
	}
}
