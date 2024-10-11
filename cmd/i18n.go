// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package cmd

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

const startFileNameI18n = "i18nUpdate.go"

var i18nCmd = &cobra.Command{
	Use:   "i18n",
	Short: "Internationalization utilities",
	Long:  `Internationalization utilities for Hexya`,
}

var i18nUpdate = &cobra.Command{
	Use:   "update [dir]",
	Short: "Create or update PO files",
	Long: `Create or update PO files of the module specified by 'dir'.
PO files will be generated for each loaded language (--language flag)
in the i18n directory of the module.`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleDir, _ := filepath.Abs(".")
		if len(args) > 0 {
			moduleDir = args[0]
		}
		langs, err := cmd.Flags().GetStringSlice("languages")
		if err != nil {
			log.Panic("Unable to read languages from the command line")
		}

		// Load the appropriate template using GetTemplate function
		i18nTemplate, err := templates.GetTemplate("StartFileTemplateI18n")
		if err != nil {
			log.Panic("Failed to retrieve template StartFileTemplateI18n: %v", err)
		}

		// Generate and update PO files
		generateAndUpdatePOFiles(moduleDir, langs, i18nTemplate)
	},
}

// generateAndUpdatePOFile creates the startup file of the translation update and runs it.
func generateAndUpdatePOFiles(moduleDir string, langs []string, tmpl *template.Template) {
	testEnabled = true
	runGenerate(moduleDir)
	fmt.Println("Please wait, Hexya is starting ...")
	moduleDir, _ = filepath.Abs(moduleDir)
	loadConf := packages.Config{
		Mode: packages.NeedName,
	}
	packs, err := packages.Load(&loadConf, moduleDir)
	if err != nil {
		panic(err)
	}
	modulePath := packs[0].PkgPath
	conf := make(map[string]interface{})
	conf["moduleDir"] = moduleDir
	conf["modulePath"] = modulePath
	conf["langs"] = langs
	tmplData := struct {
		Imports []string
		Config  string
	}{
		Imports: []string{modulePath},
		Config:  fmt.Sprintf("%#v", conf),
	}
	fileName := filepath.Join(os.TempDir(), startFileNameI18n)

	// Create the startup file from the provided template
	if err := templates.CreateFileFromTemplate(fileName, tmpl, tmplData); err != nil {
		log.Panic("Failed to create temporary file for PO update: %v", err)
	}

	cmd := exec.Command("go", "run", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.Remove(fileName)
}

func init() {
	i18nUpdate.PersistentFlags().StringSliceP("languages", "l", []string{}, "Comma separated list of languages codes to load (ex: fr,de,es).")
	HexyaCmd.AddCommand(i18nCmd)
	i18nCmd.AddCommand(i18nUpdate)
}
