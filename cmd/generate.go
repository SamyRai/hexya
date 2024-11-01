package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/gomod"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"
	"path/filepath"
)

var (
	generateEmptyPool bool
	testEnabled       bool
	replaces          []string
)

// generateCmd defines the command to generate the pool and related files.
var generateCmd = &cobra.Command{
	Use:   "generate PROJECT_DIR",
	Short: "Generate the source code of the model pool",
	Long: `Generate the source code of the pool package which includes the definition of all the models.
This command also:
- Creates the resource directory by symlinking all modules resources into the project directory.
- Creates or updates the main.go of the project.
This command must be rerun after each source code modification, including module import.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must specify the project directory")
			return
		}
		runGenerate(args[0])
	},
}

func init() {
	HexyaCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&testEnabled, "test", "t", false, "Generate pool for testing a module.")
	generateCmd.Flags().BoolVar(&generateEmptyPool, "empty", false, "Generate an empty pool package and return.")
	generateCmd.Flags().StringSliceVar(&replaces, "replace", []string{}, "Custom replace directives for go.mod files")
}

func runGenerate(projectDir string) {
	fmt.Println("Hexya Generate\n--------------")

	// Step 1: Load templates
	if err := loadTemplates(); err != nil {
		fmt.Printf("Error in loading templates: %v\n", err)
		return
	}

	// Step 2: Compute directories
	projectPath, poolPath := computeDirectories(projectDir)

	// Step 3: Prepare the pool directory
	if err := preparePoolDirectory(poolPath); err != nil {
		fmt.Printf("Error in preparing pool directory: %v\n", err)
		return
	}

	// If an empty pool is requested, stop here
	if generateEmptyPool {
		fmt.Println("Empty pool generated successfully")
		return
	}

	// Step 4: Retrieve the list of modules for the project
	modulePaths := getModulePaths(projectPath)

	// Step 5: Create Go modules files for dependencies
	if err := createGoModFiles(poolPath, projectPath, modulePaths); err != nil {
		fmt.Printf("Error in creating go.mod files: %v\n", err)
		return
	}

	// Step 6: Load packages of the modules
	packages := loadProgramPackages(modulePaths)

	// Step 7: Extract module information from packages
	moduleInfos, err := utils.GetModulePackages(packages)
	if err != nil {
		fmt.Printf("Error in getting module packages: %v\n", err)
		return
	}

	// Step 8: Create symbolic links for the modules
	if err := createModuleSymlinks(moduleInfos, projectPath); err != nil {
		fmt.Printf("Error in generating symlinks: %v\n", err)
		return
	}

	// Step 9: Generate the pool files for models
	if err := generatePoolFiles(moduleInfos, poolPath); err != nil {
		fmt.Printf("Error in generating pool files: %v\n", err)
		return
	}

	// Step 10: Verify the generated code
	if err := verifyGeneratedCode(replaces, testEnabled); err != nil {
		fmt.Printf("Error in verifying generated code: %v\n", err)
		return
	}

	// Step 11: Create or update the main.go file
	if err := createMainFile(projectPath, modulePaths); err != nil {
		fmt.Printf("Error in creating main.go: %v\n", err)
		return
	}

	// Step 12: Tidy up go.mod files
	if err := tidyGoModFiles(projectPath, poolPath); err != nil {
		fmt.Printf("Error in tidying go.mod files: %v\n", err)
		return
	}

	fmt.Println("Pool generated successfully")
}

// loadTemplates loads necessary templates for the generation process.
func loadTemplates() error {
	fmt.Println("1/9 - Loading templates... ")
	if err := templates.LoadTemplates(); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// computeDirectories calculates the project and pool directory paths.
func computeDirectories(projectDir string) (string, string) {
	fmt.Println("2/9 - Computing directories... ")
	projectPath, poolPath, err := file_operations.ComputeDirs(projectDir)
	if err != nil {
		fmt.Printf("Error computing directories: %v\n", err)
	}
	fmt.Printf("Project directory: %s\nPool directory: %s\n", projectPath, poolPath)
	fmt.Println("Ok")
	return projectPath, poolPath
}

// preparePoolDirectory cleans and sets up the pool directory.
func preparePoolDirectory(poolPath string) error {
	fmt.Println("3/9 - Preparing pool directory... ")
	if err := file_operations.CleanPoolDir(poolPath); err != nil {
		return err
	}
	if err := file_operations.CreateEmptyPool(poolPath); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// getModulePaths retrieves a list of modules for the project.
func getModulePaths(projectDir string) []string {
	if testEnabled {
		return []string{projectDir}
	}
	return viper.GetStringSlice("Modules")
}

// createGoModFiles generates go.mod files for dependencies.
func createGoModFiles(poolPath, projectPath string, modulePaths []string) error {
	fmt.Println("4/9 - Creating go.mod files... ")
	replacesToml := viper.GetStringSlice("Replaces")
	if len(replaces) > 0 {
		replacesToml = replaces
	}
	fmt.Printf("Target paths: %v\nCustom replaces: %v\n", modulePaths, replacesToml)
	return gomod.CreateGoModFiles(poolPath, projectPath, replacesToml, modulePaths)
}

// loadProgramPackages loads all packages for the modules.
func loadProgramPackages(modulePaths []string) []*packages.Package {
	fmt.Println("5/9 - Loading all program packages... ")
	fmt.Printf("Modules paths: %+v\n", modulePaths)
	packs, err := file_operations.LoadProgram(modulePaths, testEnabled)
	if err != nil {
		fmt.Printf("Error loading program packages: %v\n", err)
	}
	fmt.Println("Ok")
	return packs
}

// createModuleSymlinks generates symlinks for modules.
func createModuleSymlinks(moduleInfos []*models.ModuleInfo, projectDir string) error {
	file_operations.CleanModuleSymlinks(projectDir)
	for _, moduleInfo := range moduleInfos {
		if err := file_operations.CreateModuleSymlinks(moduleInfo, projectDir); err != nil {
			return err
		}
	}
	return nil
}

// generatePoolFiles generates pool files for the models.
func generatePoolFiles(moduleInfos []*models.ModuleInfo, poolPath string) error {
	fmt.Println("7/9 - Generating pool files... ")
	if err := generate.CreatePool(moduleInfos, poolPath); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// verifyGeneratedCode verifies the generated code for correctness.
func verifyGeneratedCode(replaces []string, testEnabled bool) error {
	fmt.Println("8/9 - Verifying generated code... ")
	_, err := file_operations.LoadProgram(replaces, testEnabled)
	if err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

// createMainFile creates or updates the main.go file for the project.
func createMainFile(projectPath string, modulePaths []string) error {
	fmt.Println("9/9 - Finalizing main.go file... ")
	coreImports, moduleImports := models.GatherCoreAndModuleImports(modulePaths)
	return file_operations.CreateMainFile(projectPath, coreImports, moduleImports, filepath.Base(projectPath))
}

// tidyGoModFiles runs go mod tidy on specified directories.
func tidyGoModFiles(projectPath, poolPath string) error {
	fmt.Println("Finalizing go.mod files... ")
	if err := gomod.TidyGoMod(projectPath); err != nil {
		return err
	}
	return gomod.TidyGoMod(poolPath)
}
