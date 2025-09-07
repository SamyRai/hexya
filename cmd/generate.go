package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/spf13/cobra"
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
		generator := generate.NewGenerator(args[0], testEnabled, generateEmptyPool, replaces)
		generator.Run()
	},
}

func init() {
	HexyaCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&testEnabled, "test", "t", false, "Generate pool for testing a module.")
	generateCmd.Flags().BoolVar(&generateEmptyPool, "empty", false, "Generate an empty pool package and return.")
	generateCmd.Flags().StringSliceVar(&replaces, "replace", []string{}, "Custom replace directives for go.mod files")
}
