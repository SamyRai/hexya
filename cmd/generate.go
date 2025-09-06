package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/hexya-erp/hexya/src/tools/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"
)

const (
	// PoolDirRel is the name of the generated pool directory (relative to the current project root)
	PoolDirRel = "pool"
	// ResDirRel is the name of the resources directory (relative to the current project root)
	ResDirRel = "res"
)

var (
	testEnabled bool
)

var GenerateCmd = &cobra.Command{
	Use:   "generate PROJECT_DIR",
	Short: "Generate the source code of the model pool",
	Long: `Generate the source code of the pool package which includes the definition of all the models.
This command must be rerun after each source code modification, including module import.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify the project directory ")
			os.Exit(1)
		}
		viper.Set("LogStdout", true)
		runGenerate(args[0])
	},
}

func runGenerate(projectDir string) {
	fmt.Printf("runGenerate called with projectDir: %s\n", projectDir)
	defer fmt.Println("runGenerate returned")
	poolDir := "pool"
	var targetPaths []string
	targetPaths = []string{projectDir}

	fmt.Println(`Hexya Generate
	--------------`)
	fmt.Println("Modules paths:")
	fmt.Println(" -", strings.Join(targetPaths, "\n - "))

	fmt.Print(`1/2 - Loading program...`)
	packs, err := loadProgram(targetPaths, false)
	if err != nil {
		panic(err)
	}
	mods := parser.GetModulePackages(packs)
	fmt.Println("Ok")

	fmt.Print("2/2 - Generating pool...")
	graph := parser.GetModelsASTData(mods)
	generate.CreatePool(graph, poolDir)
	fmt.Println("Ok")

	fmt.Println("Pool generated successfully")
}

func loadProgram(targetPaths []string, tests bool) ([]*packages.Package, error) {
	conf := packages.Config{
		Mode: packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: tests,
	}
	packs, err := packages.Load(&conf, targetPaths...)
	return packs, err
}

func writeFileFromTemplate(fileName string, tmpl *template.Template, data interface{}) error {
	var buf bytes.Buffer
	tmpl.Execute(&buf, data)
	err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	return err
}

var symlinkDirs = []string{"static", "data", "demo", "resources", "i18n"}

func init() {
	HexyaCmd.AddCommand(GenerateCmd)
}
