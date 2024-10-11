package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/utils"
	"path/filepath"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/gomod"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	generateEmptyPool bool
	testEnabled       bool
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
			log.Panic("You must specify the project directory")
		}
		runGenerate(args[0])
	},
}

func init() {
	HexyaCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&testEnabled, "test", "t", false, "Generate pool for testing a module. When set projectDir must be the source directory of the module.")
	generateCmd.Flags().BoolVar(&generateEmptyPool, "empty", false, "Generate an empty pool package and return. When set, resource dir and main.go are untouched.")
}

func runGenerate(projectDir string) {
	fmt.Println("Hexya Generate\n--------------")

	// Step 1: Load all templates
	fmt.Print("1/9 - Loading templates... ")
	if err := templates.LoadTemplates(); err != nil {
		log.Panic(fmt.Sprintf("Error loading templates: %v", err))
	}
	fmt.Println("Ok")

	// Step 2: Compute project and pool directories
	fmt.Print("2/9 - Computing directories... ")
	projectDir, poolDir, err := file_operations.ComputeDirs(projectDir)
	if err != nil {
		log.Panic(fmt.Sprintf("Error computing directories: %v", err))
	}
	fmt.Printf("Project directory: %s\nPool directory: %s\n", projectDir, poolDir)
	fmt.Println("Ok")

	// Step 3: Clean and prepare Pool Directory
	fmt.Print("3/9 - Preparing pool directory... ")
	if err := file_operations.CleanPoolDir(poolDir); err != nil {
		log.Panic(fmt.Sprintf("Error cleaning pool directory: %v", err))
	}
	if err := file_operations.CreateEmptyPool(poolDir); err != nil {
		log.Panic(fmt.Sprintf("Error creating empty pool structure: %v", err))
	}
	fmt.Println("Ok")

	// Early exit if empty pool is requested
	if generateEmptyPool {
		fmt.Println("Empty pool generated successfully")
		return
	}

	// Step 4: Generate Go Mod Files for Project and Pool
	fmt.Print("4/9 - Creating go.mod files... ")
	if err := gomod.CreateGoModFiles(poolDir, projectDir); err != nil {
		log.Panic(fmt.Sprintf("Error creating go.mod files: %v", err))
	}
	fmt.Println("Ok")

	err = gomod.TidyGoMod(poolDir)
	if err != nil {
		log.Panic(fmt.Sprintf("Error running go mod tidy for pool directory: %v", err))
	}

	// Step 5: Load Program Packages using GetModulePackages
	fmt.Print("5/9 - Loading program packages... ")
	var targetPaths []string
	if testEnabled {
		targetPaths = []string{projectDir}
	} else {
		targetPaths = viper.GetStringSlice("Modules")
	}

	// Using GetModulePackages to load the packages
	packs, err := file_operations.LoadProgram(targetPaths, testEnabled)
	if err != nil {
		log.Panic(fmt.Sprintf("Error loading program packages: %v", err))
	}

	modules, err := utils.GetModulePackages(packs) // Updated part
	if err != nil {
		log.Panic(fmt.Sprintf("Error getting module packages: %v", err))
	}

	fmt.Println("Modules loaded successfully: ", len(modules))
	fmt.Printf("Modules: %v\n", modules)

	// Step 6: Generate Symlinks for Resources
	fmt.Print("6/9 - Generating symlinks for resources... ")
	fmt.Println("Modules paths:")
	fmt.Println(" -", strings.Join(targetPaths, "\n - "))

	file_operations.CleanModuleSymlinks(projectDir)
	file_operations.CreateModuleSymlinks(modules, projectDir)
	fmt.Println("Ok")

	// Step 7: Generate Pool Files
	fmt.Print("7/9 - Generating pool files... ")

	// Log model data to ensure correctness before generation
	fmt.Println("[INFO] Validating models before generating pool files:")
	for _, mod := range modules {
		if len(mod.Syntax) == 0 {
			fmt.Printf("[WARNING] Module %s has no syntax information loaded. Skipping...\n", mod.PkgPath)
			continue
		}
		fmt.Printf("[INFO] Module: %s\n", mod.PkgPath)
		for _, file := range mod.Syntax {
			fmt.Printf("[INFO] Loaded file: %s\n", mod.FSet.Position(file.Pos()))
		}
	}

	// Ensure models have the correct data before passing them for pool generation
	modelData := ast.GetModelsASTData(modules)
	if len(modelData) == 0 {
		log.Panic("[ERROR] No valid models found for pool generation.")
	}

	// Proceed with pool file generation
	if err := generate.CreatePool(modelData, poolDir); err != nil {
		log.Panic(fmt.Sprintf("Error generating pool: %v", err))
	}
	fmt.Println("Ok")

	// Step 8: Verify Generated Code by Reloading
	fmt.Print("8/9 - Verifying generated code... ")
	_, err = file_operations.LoadProgram(targetPaths, testEnabled)
	if err != nil {
		log.Panic(fmt.Sprintf("Generated code verification failed: %v", err))
	}
	fmt.Println("Ok")

	// Step 9: Generate `main.go` File
	fmt.Print("9/9 - Generating main.go file... ")
	coreImports, moduleImports := models.GatherCoreAndModuleImports(modules)
	if err := file_operations.CreateMainFile(projectDir, coreImports, moduleImports, filepath.Base(projectDir)); err != nil {
		log.Panic(fmt.Sprintf("Error generating main.go: %v", err))
	}
	fmt.Println("Ok")

	// Final Step: Run `go mod tidy` for Dependency Cleanup
	fmt.Print("Finalizing go.mod files... ")
	if err := gomod.TidyGoMod(projectDir); err != nil {
		log.Panic(fmt.Sprintf("Error running go mod tidy for project directory: %v", err))
	}
	if err := gomod.TidyGoMod(poolDir); err != nil {
		log.Panic(fmt.Sprintf("Error running go mod tidy for pool directory: %v", err))
	}
	fmt.Println("Ok")

	fmt.Println("Pool generated successfully")
}
