package gomod

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/logging"
	"github.com/spf13/viper"
)

var log = logging.GetLogger("generate/gomod")

func init() {
	logging.Initialize()
}

// CreateGoModFiles creates the go.mod files for the pool and project directories.
func CreateGoModFiles(poolDir, projectDir string) error {
	log.Info("Starting to create go.mod files for pool and project directories...")

	// Create go.mod for the pool directory
	if err := createGoModFile(poolDir, true, poolDir); err != nil {
		return fmt.Errorf("failed to create go.mod for pool directory: %w", err)
	}

	// Create go.mod for the project directory
	if err := createGoModFile(projectDir, false, poolDir); err != nil {
		return fmt.Errorf("failed to create go.mod for project directory: %w", err)
	}

	log.Info("go.mod files created successfully for pool and project directories")

	// Run go mod tidy for both directories to ensure proper dependency management
	if err := TidyGoMod(poolDir); err != nil {
		return fmt.Errorf("failed to run go mod tidy for pool directory: %w", err)
	}

	if err := TidyGoMod(projectDir); err != nil {
		return fmt.Errorf("failed to run go mod tidy for project directory: %w", err)
	}

	return nil
}

// createGoModFile creates the go.mod file for either the pool or the project directory.
func createGoModFile(targetDir string, isPool bool, poolDir string) error {
	log.Info(fmt.Sprintf("Creating go.mod for directory: %s", targetDir))

	// Get required modules and replacements
	requiredModules := getRequiredModules(isPool)
	replacements := getReplacements(isPool, poolDir)

	// Prepare template data
	templateData := GoModTemplateData{
		ModuleName:      determineModuleName(targetDir, isPool),
		GoVersion:       "1.23",
		RequiredModules: requiredModules,
		Replaces:        replacements,
	}

	// Create go.mod from template
	return createGoModFileFromTemplate(targetDir, templateData)
}

// createGoModFileFromTemplate creates a go.mod file using the provided template data.
func createGoModFileFromTemplate(dir string, templateData GoModTemplateData) error {
	tmpl, err := templates.GetTemplate("GoModTemplate")
	if err != nil {
		return fmt.Errorf("failed to get go.mod template: %w", err)
	}

	goModPath := filepath.Join(dir, "go.mod")
	if err := templates.CreateFileFromTemplate(goModPath, tmpl, templateData); err != nil {
		return fmt.Errorf("failed to create go.mod file in directory %s: %w", dir, err)
	}

	log.Info(fmt.Sprintf("go.mod file created successfully in directory: %s", dir))
	return nil
}

// runGoModTidy runs `go mod tidy` for the specified directory.
func TidyGoMod(dir string) error {
	log.Info(fmt.Sprintf("Running go mod tidy for directory: %s", dir))
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Error(fmt.Sprintf("go mod tidy failed for directory: %s", dir), "error", err)
		fmt.Printf("stderr: %s\n error: %v\n", stderr.String(), err)
		return fmt.Errorf("failed to run go mod tidy for directory %s: %w", dir, err)
	}

	log.Info(fmt.Sprintf("go mod tidy completed successfully for directory: %s", dir))
	return nil
}

// determineModuleName determines the Go module name for the given directory.
func determineModuleName(dir string, isPool bool) string {
	if isPool {
		return "github.com/hexya-erp/pool"
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to determine absolute path for directory %s: %s", dir, err))
	}
	return generateModuleName(absDir)
}

// generateModuleName generates a suitable Go module name for the provided directory.
func generateModuleName(absDir string) string {
	projectName := filepath.Base(absDir)
	moduleName := strings.ReplaceAll(projectName, " ", "_")
	return fmt.Sprintf("github.com/your-org/%s", moduleName)
}

// getRequiredModules returns the list of required modules, adding a default version if none is specified.
func getRequiredModules(isPool bool) []string {
	if isPool {
		return []string{"github.com/hexya-erp/hexya v1.0.2"}
	}

	modules := viper.GetStringSlice("Modules")
	for i, module := range modules {
		if !strings.Contains(module, "v") {
			modules[i] = module + " v0.1.0"
		}
	}
	return modules
}

// getReplacements returns the replacement directives for the go.mod file.
func getReplacements(isPool bool, poolDir string) []ModuleReplacement {
	var replacements []ModuleReplacement

	if isPool {
		// Pool-specific replacements
		replacements = append(replacements, ModuleReplacement{
			Module: "github.com/hexya-erp/hexya",
			Path:   `"/Users/damirmukimov/projects/Glowing Pixels/hexya-main/hexya"`,
		})
	} else {
		// Project-specific replacements
		replacesFromConfig := viper.GetStringSlice("Replaces")
		for _, replace := range replacesFromConfig {
			parts := strings.Split(replace, " => ")
			if len(parts) == 2 {
				replacements = append(replacements, ModuleReplacement{
					Module: parts[0],
					Path:   filepath.ToSlash(parts[1]),
				})
			}
		}

		// Add specific replacement for the pool
		replacements = append(replacements, ModuleReplacement{
			Module: "github.com/hexya-erp/pool",
			Path:   "./pool",
		})
	}

	return replacements
}

// ModuleReplacement represents a Go module replacement directive.
type ModuleReplacement struct {
	Module string
	Path   string
}

// GoModTemplateData represents the data needed to create a go.mod file.
type GoModTemplateData struct {
	ModuleName      string
	GoVersion       string
	RequiredModules []string
	Replaces        []ModuleReplacement
}
