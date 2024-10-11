// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/spf13/cobra"
)

var symlinkDirs = []string{"static", "data", "demo", "resources", "i18n"}

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Module development utilities",
	Long:  `Hexya utilities for module development.`,
}

var moduleInitCmd = &cobra.Command{
	Use:   "init MODULE_PATH",
	Short: "Initialize a module",
	Long: `Initialize and scaffold a new module in the current directory with the given path (e.g. github.com/myuser/my-hexya-module).
Use this command if you plan to distribute your module.
Note that you will need to commit your module to its remote repository before consuming it in a project.
Alternatively, you can manually set the replace directive in your project go.mod to point to this directory.

For local only modules (i.e. modules tied to a project), use 'hexya module new' from the project directory instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the go.mod file
		if len(args) == 0 {
			fmt.Println("You must specify a module path.")
			os.Exit(1)
		}
		modulePath := args[0]
		if err := runCommand("go", "mod", "init", modulePath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Load the template
		hexyaGoTemplate, err := templates.GetTemplate("HexyaGoTemplate")
		if err != nil {
			fmt.Printf("Error retrieving hexya go template: %v\n", err)
			os.Exit(1)
		}

		// Create the 000hexya.go file
		data := struct {
			ModuleName string
		}{
			ModuleName: path.Base(modulePath),
		}
		if err := templates.CreateFileFromTemplate("000hexya.go", hexyaGoTemplate, data); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create standard directories
		for _, dir := range symlinkDirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println(err)
			}
		}
		runCommand("go", "mod", "tidy")
	},
}

var moduleNewCmd = &cobra.Command{
	Use:   "new MODULE_NAME",
	Short: "Initialize a new local module in current project",
	Long: `Initialize and scaffold a new local module in the current project. 
The current directory must be a Hexya project directory.

If you plan to make a module and distribute it on its own, you should create a new directory and run 'hexya module init' inside instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify a module name.")
			os.Exit(1)
		}
		moduleName := args[0]

		// Check we are in a project dir (at least a dir with go.mod)
		c := exec.Command("go", "list", "-f", "'{{ .Name }}")
		if res, err := c.Output(); err != nil || string(res) != "main" {
			fmt.Println("You must call hexya module new from a project directory.")
			fmt.Println(err)
			os.Exit(1)
		}

		// Get this project path
		c = exec.Command("go", "list", "-m")

		// Create hexya module subdir
		os.MkdirAll(moduleName, 0755)

		// Load the template
		hexyaGoTemplate, err := templates.GetTemplate("HexyaGoTemplate")
		if err != nil {
			fmt.Printf("Error retrieving hexya go template: %v\n", err)
			os.Exit(1)
		}

		// Create the 000hexya.go file
		data := struct {
			ModuleName string
		}{
			ModuleName: moduleName,
		}
		if err := templates.CreateFileFromTemplate(filepath.Join(moduleName, "000hexya.go"), hexyaGoTemplate, data); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create standard directories
		for _, dir := range symlinkDirs {
			if err := os.MkdirAll(filepath.Join(moduleName, dir), 0755); err != nil {
				fmt.Println(err)
			}
		}
		runCommand("go", "mod", "tidy")
	},
}

var moduleCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the module directory",
	Long: `Clean the current directory from all generated and test artifacts.
You should use this command before committing your work.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCommand("go", "mod", "edit", "-dropreplace", "github.com/hexya-erp/pool@v1.0.2")
		if err := removeProjectDir(config.PoolDirRel); err != nil {
			fmt.Println(err)
		}
		if err := removeProjectDir(config.ResDirRel); err != nil {
			fmt.Println(err)
		}
		runCommand("go", "mod", "tidy")
	},
}

func init() {
	HexyaCmd.AddCommand(moduleCmd)
	moduleCmd.AddCommand(moduleInitCmd)
	moduleCmd.AddCommand(moduleNewCmd)
	moduleCmd.AddCommand(moduleCleanCmd)
}
