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

	if err := loadTemplates(); err != nil {
		fmt.Printf("Error in loading templates: %v\n", err)
		return
	}

	projectDir, poolDir := computeDirectories(projectDir)

	if err := preparePoolDirectory(poolDir); err != nil {
		fmt.Printf("Error in preparing pool directory: %v\n", err)
		return
	}

	if generateEmptyPool {
		fmt.Println("Empty pool generated successfully")
		return
	}

	modulesToml := getModulesList(projectDir)

	if err := createGoModFiles(poolDir, projectDir, modulesToml); err != nil {
		fmt.Printf("Error in creating go.mod files: %v\n", err)
		return
	}

	packs := loadProgramPackages(modulesToml)

	mods, err := utils.GetModulePackages(packs)
	if err != nil {
		fmt.Printf("Error in getting module packages: %v\n", err)
		return
	}

	if err := createSymlinks(mods, projectDir); err != nil {
		fmt.Printf("Error in generating symlinks: %v\n", err)
		return
	}

	if err := generatePoolFiles(packs, poolDir); err != nil {
		fmt.Printf("Error in generating pool files: %v\n", err)
		return
	}

	if err := verifyGeneratedCode(replaces, testEnabled); err != nil {
		fmt.Printf("Error in verifying generated code: %v\n", err)
		return
	}

	if err := createMainFile(projectDir, modulesToml); err != nil {
		fmt.Printf("Error in creating main.go: %v\n", err)
		return
	}

	if err := tidyGoModFiles(projectDir, poolDir); err != nil {
		fmt.Printf("Error in tidying go.mod files: %v\n", err)
		return
	}

	fmt.Println("Pool generated successfully")
}

func loadTemplates() error {
	fmt.Print("1/9 - Loading templates... ")
	if err := templates.LoadTemplates(); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

func computeDirectories(projectDir string) (string, string) {
	fmt.Print("2/9 - Computing directories... ")
	projectDir, poolDir, err := file_operations.ComputeDirs(projectDir)
	if err != nil {
		fmt.Printf("Error computing directories: %v\n", err)
	}
	fmt.Printf("Project directory: %s\nPool directory: %s\n", projectDir, poolDir)
	fmt.Println("Ok")
	return projectDir, poolDir
}

func preparePoolDirectory(poolDir string) error {
	fmt.Print("3/9 - Preparing pool directory... ")
	if err := file_operations.CleanPoolDir(poolDir); err != nil {
		return err
	}
	if err := file_operations.CreateEmptyPool(poolDir); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

func getModulesList(projectDir string) []string {
	if testEnabled {
		return []string{projectDir}
	}
	return viper.GetStringSlice("Modules")
}

func createGoModFiles(poolDir, projectDir string, modulesToml []string) error {
	fmt.Print("4/9 - Creating go.mod files... ")
	replacesToml := viper.GetStringSlice("Replaces")
	if len(replaces) > 0 {
		replacesToml = replaces
	}
	fmt.Printf("Target paths: %v\nCustom replaces: %v\n", modulesToml, replacesToml)
	return gomod.CreateGoModFiles(poolDir, projectDir, replacesToml, modulesToml)
}

func loadProgramPackages(modulesToml []string) []*packages.Package {
	fmt.Print("5/9 - Loading all program packages... ")
	fmt.Printf("Modules paths: %+v\n", modulesToml)
	packs, err := file_operations.LoadProgram(modulesToml, testEnabled)
	if err != nil {
		fmt.Printf("Error loading program packages: %v\n", err)
	}
	fmt.Println("Ok")
	return packs
}

func createSymlinks(modules []*models.ModuleInfo, projectDir string) error {
	file_operations.CleanModuleSymlinks(projectDir)
	for _, m := range modules {
		// Create symlinks for all modules
		if err := file_operations.CreateModuleSymlinks(m, projectDir); err != nil {
			return err
		}
	}
	return nil
}

func generatePoolFiles(packs []*packages.Package, poolDir string) error {
	fmt.Print("7/9 - Generating pool files... ")
	if err := generate.CreatePool(packs, poolDir); err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

func verifyGeneratedCode(replaces []string, testEnabled bool) error {
	fmt.Print("8/9 - Verifying generated code... ")
	_, err := file_operations.LoadProgram(replaces, testEnabled)
	if err != nil {
		return err
	}
	fmt.Println("Ok")
	return nil
}

func createMainFile(projectDir string, modulesToml []string) error {
	fmt.Print("9/9 - Finalizing main.go file... ")
	coreImports, moduleImports := models.GatherCoreAndModuleImports(modulesToml)
	return file_operations.CreateMainFile(projectDir, coreImports, moduleImports, filepath.Base(projectDir))
}

func tidyGoModFiles(projectDir, poolDir string) error {
	fmt.Print("Finalizing go.mod files... ")
	if err := gomod.TidyGoMod(projectDir); err != nil {
		return err
	}
	return gomod.TidyGoMod(poolDir)
}
