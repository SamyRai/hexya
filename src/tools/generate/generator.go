package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/gomod"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"
	"path/filepath"
)

// A Generator is the main object for the generation process
type Generator struct {
	ProjectDir      string
	TestEnabled     bool
	GenerateEmpty   bool
	Replaces        []string
	modulePaths     []string
	packages        []*packages.Package
	moduleInfos     []*models.ModuleInfo
	projectPath     string
	poolPath        string
	skippedModules  map[string]bool
	doneModules     map[string]bool
	allModules      map[string]bool
	modulesToUpdate map[string]bool
}

// NewGenerator returns a new Generator instance
func NewGenerator(projectDir string, testEnabled, generateEmpty bool, replaces []string) *Generator {
	return &Generator{
		ProjectDir:      projectDir,
		TestEnabled:     testEnabled,
		GenerateEmpty:   generateEmpty,
		Replaces:        replaces,
		skippedModules:  make(map[string]bool),
		doneModules:     make(map[string]bool),
		allModules:      make(map[string]bool),
		modulesToUpdate: make(map[string]bool),
	}
}

// Run executes the code generation
func (g *Generator) Run() {
	fmt.Println("Hexya Generate\n--------------")

	// Step 1: Load templates
	if err := g.loadTemplates(); err != nil {
		fmt.Printf("Error in loading templates: %v\n", err)
		return
	}

	// Step 2: Compute directories
	g.computeDirectories()

	// Step 3: Prepare the pool directory
	if err := g.preparePoolDirectory(); err != nil {
		fmt.Printf("Error in preparing pool directory: %v\n", err)
		return
	}

	// If an empty pool is requested, stop here
	if g.GenerateEmpty {
		fmt.Println("Empty pool generated successfully")
		return
	}

	// Step 4: Retrieve the list of modules for the project
	g.getModulePaths()

	// Step 5: Create Go modules files for dependencies
	if err := g.createGoModFiles(); err != nil {
		fmt.Printf("Error in creating go.mod files: %v\n", err)
		return
	}

	// Step 6: Load packages of the modules
	g.loadProgramPackages()

	// Step 7: Extract module information from packages
	var err error
	g.moduleInfos, err = utils.GetModulePackages(g.packages)
	if err != nil {
		fmt.Printf("Error in getting module packages: %v\n", err)
		return
	}

	// Step 8: Create symbolic links for the modules
	if err := g.createModuleSymlinks(); err != nil {
		fmt.Printf("Error in generating symlinks: %v\n", err)
		return
	}

	// Step 9: Generate the pool files for models
	if err := g.generatePoolFiles(); err != nil {
		fmt.Printf("Error in generating pool files: %v\n", err)
		return
	}

	// Step 10: Verify the generated code
	if err := g.verifyGeneratedCode(); err != nil {
		fmt.Printf("Error in verifying generated code: %v\n", err)
		return
	}

	// Step 11: Create or update the main.go file
	if err := g.createMainFile(); err != nil {
		fmt.Printf("Error in creating main.go: %v\n", err)
		return
	}

	// Step 12: Tidy up go.mod files
	if err := g.tidyGoModFiles(); err != nil {
		fmt.Printf("Error in tidying go.mod files: %v\n", err)
		return
	}

	fmt.Println("Pool generated successfully")
}

// loadTemplates loads necessary templates for the generation process.
func (g *Generator) loadTemplates() error {
	fmt.Println("1/9 - Loading templates... ")
	if err := templates.LoadTemplates(); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// computeDirectories calculates the project and pool directory paths.
func (g *Generator) computeDirectories() {
	fmt.Println("2/9 - Computing directories... ")
	var err error
	g.projectPath, g.poolPath, err = file_operations.ComputeDirs(g.ProjectDir)
	if err != nil {
		fmt.Printf("Error computing directories: %v\n", err)
	}
	fmt.Printf("Project directory: %s\nPool directory: %s\n", g.projectPath, g.poolPath)
	fmt.Println("Ok")
}

// preparePoolDirectory cleans and sets up the pool directory.
func (g *Generator) preparePoolDirectory() error {
	fmt.Println("3/9 - Preparing pool directory... ")
	if err := file_operations.CleanPoolDir(g.poolPath); err != nil {
		return err
	}
	if err := file_operations.CreateEmptyPool(g.poolPath); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// getModulePaths retrieves a list of modules for the project.
func (g *Generator) getModulePaths() {
	if g.TestEnabled {
		g.modulePaths = []string{g.ProjectDir}
		return
	}
	g.modulePaths = viper.GetStringSlice("Modules")
}

// createGoModFiles generates go.mod files for dependencies.
func (g *Generator) createGoModFiles() error {
	fmt.Println("4/9 - Creating go.mod files... ")
	replacesToml := viper.GetStringSlice("Replaces")
	if len(g.Replaces) > 0 {
		replacesToml = g.Replaces
	}
	fmt.Printf("Target paths: %v\nCustom replaces: %v\n", g.modulePaths, replacesToml)
	return gomod.CreateGoModFiles(g.poolPath, g.projectPath, replacesToml, g.modulePaths)
}

// loadProgramPackages loads all packages for the modules.
func (g *Generator) loadProgramPackages() {
	fmt.Println("5/9 - Loading all program packages... ")
	fmt.Printf("Modules paths: %+v\n", g.modulePaths)
	var err error
	g.packages, err = file_operations.LoadProgram(g.modulePaths, g.TestEnabled)
	if err != nil {
		fmt.Printf("Error loading program packages: %v\n", err)
	}
	fmt.Println("Ok")
}

// createModuleSymlinks generates symlinks for modules.
func (g *Generator) createModuleSymlinks() error {
	file_operations.CleanModuleSymlinks(g.projectPath)
	for _, moduleInfo := range g.moduleInfos {
		if err := file_operations.CreateModuleSymlinks(moduleInfo, g.projectPath); err != nil {
			return err
		}
	}
	return nil
}

// generatePoolFiles generates pool files for the models.
func (g *Generator) generatePoolFiles() error {
	fmt.Println("7/9 - Generating pool files... ")
	if err := CreatePool(g.moduleInfos, g.poolPath); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// verifyGeneratedCode verifies the generated code for correctness.
func (g *Generator) verifyGeneratedCode() error {
	fmt.Println("8/9 - Verifying generated code... ")
	_, err := file_operations.LoadProgram(g.Replaces, g.TestEnabled)
	if err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// createMainFile creates or updates the main.go file for the project.
func (g *Generator) createMainFile() error {
	fmt.Println("9/9 - Finalizing main.go file... ")
	coreImports, moduleImports := models.GatherCoreAndModuleImports(g.modulePaths)
	return file_operations.CreateMainFile(g.projectPath, coreImports, moduleImports, filepath.Base(g.projectPath))
}

// tidyGoModFiles runs go mod tidy on specified directories.
func (g *Generator) tidyGoModFiles() error {
	fmt.Println("Finalizing go.mod files... ")
	if err := gomod.TidyGoMod(g.projectPath); err != nil {
		return err
	}
	return gomod.TidyGoMod(g.poolPath)
}
