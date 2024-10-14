package gomod

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/logging"
)

var log = logging.GetLogger("generate/gomod")

func init() {
	logging.Initialize()
}

// CreateGoModFiles creates the go.mod files for the pool and project directories.
func CreateGoModFiles(poolDir, projectDir string, customReplaces, addons []string) error {
	log.Info("Starting to create go.mod files for pool and project directories...")

	// Create empty go.mod files
	if err := createEmptyGoModFile(poolDir); err != nil {
		return fmt.Errorf("failed to create go.mod for pool directory: %w", err)
	}
	if err := createEmptyGoModFile(projectDir); err != nil {
		return fmt.Errorf("failed to create go.mod for project directory: %w", err)
	}

	// Add custom replacements
	if err := addReplacements(poolDir, customReplaces); err != nil {
		return fmt.Errorf("failed to add replacements for pool directory: %w", err)
	}
	if err := addReplacements(projectDir, customReplaces); err != nil {
		return fmt.Errorf("failed to add replacements for project directory: %w", err)
	}

	// Add dependencies to go.mod files, not needed for pool directory
	if err := addDependencies(projectDir, addons); err != nil {
		return fmt.Errorf("failed to add dependencies for project directory: %w", err)
	}

	// Tidy up go.mod files
	if err := TidyGoMod(poolDir); err != nil {
		return fmt.Errorf("failed to run go mod tidy for pool directory: %w", err)
	}
	if err := TidyGoMod(projectDir); err != nil {
		return fmt.Errorf("failed to run go mod tidy for project directory: %w", err)
	}

	return nil
}

// createEmptyGoModFile creates an empty go.mod file in the specified directory.
func createEmptyGoModFile(dir string) error {
	goModPath := filepath.Join(dir, "go.mod")

	// Check if go.mod file exists
	if _, err := os.Stat(goModPath); err == nil {
		// File exists, remove it
		log.Info(fmt.Sprintf("go.mod file already exists in directory: %s. Overwriting it.", dir))
		if err := os.Remove(goModPath); err != nil {
			return fmt.Errorf("failed to remove existing go.mod file in directory %s: %w", dir, err)
		}
	} else if !os.IsNotExist(err) {
		// Error checking file existence
		return fmt.Errorf("error checking go.mod existence in directory %s: %w", dir, err)
	}

	// Create a new go.mod file
	log.Info(fmt.Sprintf("Creating new go.mod file for directory: %s", dir))
	return runGoCommand(dir, "go", "mod", "init", determineModuleName(dir))
}

// addDependencies adds the provided dependencies to the go.mod file using the go get command.
func addDependencies(dir string, addons []string) error {
	for _, addon := range addons {
		fmt.Printf("Adding dependency %s to directory: %s\n", addon, dir)

		// Check if version is specified, if not, fetch the latest version
		if !strings.Contains(addon, "@") {
			latestVersion, err := getLatestVersion(addon)
			if err != nil {
				return fmt.Errorf("failed to fetch latest version for %s: %w", addon, err)
			}
			addon = fmt.Sprintf("%s@%s", addon, latestVersion)
		}

		// Use go mod edit to add the dependency
		if err := runGoCommand(dir, "go", "mod", "edit", "-require", addon); err != nil {
			return fmt.Errorf("failed to add dependency %s in directory %s: %w", addon, dir, err)
		}
	}
	return nil
}

// getLatestVersion fetches the latest available version of a module
func getLatestVersion(module string) (string, error) {
	// Use `go list -m -versions` to fetch the available versions of the module
	cmd := exec.Command("go", "list", "-m", "-versions", module)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "v0.0.0", nil
	}

	// Parse the versions output, split by space
	output := strings.TrimSpace(out.String())
	versions := strings.Split(output, " ")

	if len(versions) < 2 {
		return "v0.0.0", nil
	}

	// Return the last version (which is usually the latest)
	return versions[len(versions)-1], nil
}

// addReplacements adds custom replacements to the go.mod file.
func addReplacements(dir string, customReplaces []string) error {
	for _, replace := range customReplaces {
		parts := strings.Split(replace, " => ")
		if len(parts) == 2 {
			// Remove additional quotes around the path
			module := parts[0]
			path := strings.Trim(parts[1], `"`) // Clean the path from extra quotes

			fmt.Printf(fmt.Sprintf("\n\nAdding replacement %s => %s to directory: %s", module, path, dir))
			if err := runGoCommand(dir, "go", "mod", "edit", "-replace", fmt.Sprintf("%s=%s", module, path)); err != nil {
				return fmt.Errorf("failed to add replacement %s => %s in directory %s: %w", module, path, dir, err)
			}
		} else {
			return fmt.Errorf("invalid replacement format: %s", replace)
		}
	}
	return nil
}

// TidyGoMod runs `go mod tidy` for the specified directory.
func TidyGoMod(dir string) error {
	log.Info(fmt.Sprintf("Running go mod tidy for directory: %s", dir))
	return runGoCommand(dir, "go", "mod", "download")
}

// runGoCommand is a helper function to execute Go-related commands and capture the output.
func runGoCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		fmt.Printf(fmt.Sprintf("Command failed in directory: %s\nOutput: %s", dir, output.String()))
		return fmt.Errorf("command %s %s failed in directory %s: %w", name, strings.Join(args, " "), dir, err)
	}

	fmt.Printf(fmt.Sprintf("Command %s %s executed successfully in directory: %s", name, strings.Join(args, " "), dir))
	return nil
}

// determineModuleName generates a Go module name based on the directory structure.
func determineModuleName(dir string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to determine absolute path for directory %s: %s", dir, err))
	}
	projectName := filepath.Base(absDir)
	return fmt.Sprintf("github.com/your-org/%s", strings.ReplaceAll(projectName, " ", "_"))
}
