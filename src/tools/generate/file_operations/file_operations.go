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

// CleanPoolDir removes all contents of the pool directory.
func CleanPoolDir(poolDir string) error {
	log.Info(fmt.Sprintf("Cleaning pool directory: %s", poolDir))
	if err := os.RemoveAll(poolDir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", poolDir, err)
	}
	log.Info("Pool directory cleaned successfully.")
	return nil
}

// CreateEmptyPool sets up an empty pool structure (folders and temp files) in the pool directory.
func CreateEmptyPool(poolDir string) error {
	if err := os.MkdirAll(poolDir, 0755); err != nil {
		return fmt.Errorf("error creating pool directory %s: %w", poolDir, err)
	}

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
	conf := &packages.Config{
		Mode: packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedTypes | packages.NeedTypesSizes | packages.NeedImports |
			packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: tests,
	}

	packs, err := packages.Load(conf, targetPaths...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	return packs, nil
}

// CreateModuleSymlinks sets up symlinks for the provided module in the project directory.
// CreateModuleSymlinks sets up symlinks for the provided module in the project directory.
func CreateModuleSymlinks(module *models.ModuleInfo, projectDir string) error {
	resDir := filepath.Join(projectDir, config.ResDirRel)

	for _, dir := range config.SymlinkDirs {
		// Ensure module has Go files
		if len(module.Package.GoFiles) == 0 {
			continue // No Go files in module, skip
		}
		// Get the file system path to the module
		mDir := filepath.Dir(module.Package.GoFiles[0])
		srcPath := filepath.Join(mDir, dir)
		dstPath := filepath.Join(resDir, dir)

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			// Source directory doesn't exist, skip
			continue
		}

		if err := os.MkdirAll(dstPath, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dstPath, err)
		}

		linkName := filepath.Join(dstPath, module.Package.Name)

		// Check if the destination exists and whether it's a directory or symlink
		if info, err := os.Lstat(linkName); err == nil {
			// If it's a symlink, we can safely remove it
			if info.Mode()&os.ModeSymlink != 0 {
				if err := os.Remove(linkName); err != nil {
					return fmt.Errorf("error removing existing symlink %s: %v", linkName, err)
				}
			} else {
				// If it's a directory and not empty, log a warning and skip
				fmt.Printf("Warning: %s is a directory and not empty, skipping symlink creation\n", linkName)
				continue
			}
		}

		// Create the symlink
		if err := os.Symlink(srcPath, linkName); err != nil {
			return fmt.Errorf("error creating symlink %s -> %s: %v", linkName, srcPath, err)
		}
	}
	return nil
}

// CleanModuleSymlinks removes all existing symlinks in the project symlink directories.
func CleanModuleSymlinks(projectDir string) {
	resDir := filepath.Join(projectDir, config.ResDirRel)

	for _, dir := range config.SymlinkDirs {
		dirPath := filepath.Join(resDir, dir)
		if err := os.RemoveAll(dirPath); err != nil {
			fmt.Printf("Error removing directory %s: %v\n", dirPath, err)
		}
		if err := os.Mkdir(dirPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
		}
	}
}

// CreateMainFile generates the main.go file with core and module imports in the specified directory.
func CreateMainFile(dir string, coreImports, moduleImports []string, executableName string) error {
	// Use MainFileViewObjectConstructor to build the view object for main.go
	viewData := templates.MainFileViewObjectConstructor(executableName, coreImports, moduleImports)

	tmpl, err := templates.GetTemplate("MainFileTemplate")
	if err != nil {
		return fmt.Errorf("failed to get main.go template: %w\n", err)
	}

	mainFilePath := filepath.Join(dir, "main.go")
	if err := templates.CreateFileFromTemplate(mainFilePath, tmpl, viewData); err != nil {
		return fmt.Errorf("failed to create main.go file: %w\n", err)
	}

	log.Info(fmt.Sprintf("main.go file created successfully in directory: %s\n", dir))
	return nil
}
