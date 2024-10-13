package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/parser"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"path/filepath"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/gomod"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	generateEmptyPool bool
	testEnabled       bool
	replaces          []string
)

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
			log.Error("You must specify the project directory")
		}
		runGenerate(args[0])
	},
}

func init() {
	HexyaCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&testEnabled, "test", "t", false, "Generate pool for testing a module. When set, projectDir must be the source directory of the module.")
	generateCmd.Flags().BoolVar(&generateEmptyPool, "empty", false, "Generate an empty pool package and return. When set, resource dir and main.go are untouched.")
	generateCmd.Flags().StringSliceVar(&replaces, "replace", []string{}, "Custom replace directives for go.mod files")
}

func runGenerate(projectDir string) {
	fmt.Println("Hexya Generate\n--------------")

	// Step 1: Load all templates
	fmt.Print("1/9 - Loading templates... ")
	if err := templates.LoadTemplates(); err != nil {
		log.Error(fmt.Sprintf("Error loading templates: %v", err))
	}
	fmt.Println("Ok")

	// Step 2: Compute project and pool directories
	fmt.Print("2/9 - Computing directories... ")
	projectDir, poolDir, err := file_operations.ComputeDirs(projectDir)
	if err != nil {
		log.Error(fmt.Sprintf("Error computing directories: %v", err))
	}
	fmt.Printf("Project directory: %s\nPool directory: %s\n", projectDir, poolDir)
	fmt.Println("Ok")

	// Step 3: Clean and prepare Pool Directory
	fmt.Print("3/9 - Preparing pool directory... ")
	if err := file_operations.CleanPoolDir(poolDir); err != nil {
		log.Error(fmt.Sprintf("Error cleaning pool directory: %v", err))
	}
	if err := file_operations.CreateEmptyPool(poolDir); err != nil {
		log.Error(fmt.Sprintf("Error creating empty pool structure: %v", err))
	}
	fmt.Println("Ok")

	// Early exit if empty pool is requested
	if generateEmptyPool {
		fmt.Println("Empty pool generated successfully")
		return
	}

	// Step 4: Create Go Mod Files
	fmt.Print("4/9 - Creating go.mod files... ")
	replacesToml := viper.GetStringSlice("Replaces")
	if len(replaces) > 0 {
		replacesToml = replaces
	}
	var modulesToml []string
	if testEnabled {
		modulesToml = []string{projectDir}
	} else {
		modulesToml = viper.GetStringSlice("Modules")
	}
	fmt.Printf("Target paths: %v\n", modulesToml)
	fmt.Print("Custom replaces: ", replacesToml)
	// Copy modules into addons list
	addonsList := make([]string, len(modulesToml))
	copy(addonsList, modulesToml)
	if err := gomod.CreateGoModFiles(poolDir, projectDir, replacesToml, addonsList); err != nil {
		log.Error(fmt.Sprintf("Error creating go.mod files: %v", err))
	}
	fmt.Println("Ok")

	// Step 5: Load all Program Packages (core and addon together)
	fmt.Print("5/9 - Loading all program packages... ")
	fmt.Printf("Modules paths: %+v\n", modulesToml)
	packs, err := file_operations.LoadProgram(modulesToml, testEnabled)
	if err != nil {
		log.Error(fmt.Sprintf("Error loading program packages: %v", err))
	}

	// Step 6: Process all modules (core and addons)
	fmt.Println("Program packages loaded successfully. Identifying modules...")
	fmt.Printf("Packages loaded: %+v\n", packs)
	modules, err := utils.GetModulePackages(packs)
	if err != nil {
		log.Error(fmt.Sprintf("Error identifying modules: %v", err))
	}

	fmt.Printf("Modules identified: %d\n", len(modules))

	// Step 7: Generate Symlinks for Resources
	fmt.Print("7/9 - Generating symlinks for resources... ")
	fmt.Println("Modules paths:")
	fmt.Println(" -", strings.Join(modulesToml, "\n - "))

	file_operations.CleanModuleSymlinks(projectDir)
	file_operations.CreateModuleSymlinks(modules, projectDir)
	fmt.Println("Ok")

	// Step 8: Generate Pool Files (for all modules)
	fmt.Print("8/9 - Generating pool files... ")

	// Gather AST data for all models (core and addons)
	modelsASTData := parser.GetModelsASTData(modules, true)
	if len(modelsASTData) == 0 {
		log.Error("[ERROR] No valid models found for pool generation.")
	}

	// Generate pool for all modules at once
	if err := generate.CreatePool(modelsASTData, poolDir); err != nil {
		log.Error(fmt.Sprintf("Error generating pool: %v", err))
	}
	fmt.Println("Ok")

	// Step 9: Verify Generated Code by Reloading
	fmt.Print("9/9 - Verifying generated code... ")
	_, err = file_operations.LoadProgram(replaces, testEnabled)
	if err != nil {
		log.Error(fmt.Sprintf("Generated code verification failed: %v", err))
	}
	fmt.Println("Ok")

	// Final Step: Generate `main.go` File
	fmt.Print("Finalizing main.go file... ")
	coreImports, moduleImports := models.GatherCoreAndModuleImports(modulesToml)
	if err := file_operations.CreateMainFile(projectDir, coreImports, moduleImports, filepath.Base(projectDir)); err != nil {
		log.Error(fmt.Sprintf("Error generating main.go: %v", err))
	}
	fmt.Println("Ok")

	// Final Cleanup Step: Run `go mod tidy`
	fmt.Print("Finalizing go.mod files... ")
	if err := gomod.TidyGoMod(projectDir); err != nil {
		log.Error(fmt.Sprintf("Error running go mod tidy for project directory: %v", err))
	}
	if err := gomod.TidyGoMod(poolDir); err != nil {
		log.Error(fmt.Sprintf("Error running go mod tidy for pool directory: %v", err))
	}
	fmt.Println("Ok")

	fmt.Println("Pool generated successfully")
}
