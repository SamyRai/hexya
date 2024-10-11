// generate.go - Updated to look for models in folders and generate accordingly
package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/builders"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/generate/templates"
	"github.com/hexya-erp/hexya/src/tools/strutils"
	goast "go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CreatePool generates the pool package by parsing the source code AST of the given program.
func CreatePool(modelsASTData map[string]ast.ModelASTData, dir string) error {
	fmt.Println("Starting createPoolFilesFromASTData...") // Debugging

	wg := sync.WaitGroup{}
	errChan := make(chan error, len(modelsASTData))
	wg.Add(len(modelsASTData))

	for mName, mASTData := range modelsASTData {
		fmt.Printf("Processing model: %s\n", mName) // Debugging
		go func(modelName string, modelASTData ast.ModelASTData) {
			defer wg.Done()
			depsMap := map[string]bool{config.ModelsPath: true}
			mData := data.ModelData{
				Name:                  modelName,
				SnakeName:             strutils.SnakeCase(modelName),
				ModelsPackageName:     config.PoolModelPackage,
				QueryPackageName:      config.PoolQueryPackage,
				InterfacesPackageName: config.PoolInterfacesPackage,
				ModelType:             modelASTData.ModelType,
				IsModelMixin:          modelASTData.IsModelMixin,
				ConditionFuncs:        config.ConditionFuncs,
			}

			// Add fields
			fmt.Printf("Adding fields for model: %s\n", modelName) // Debugging
			builders.AddFieldsToModelData(modelASTData, &mData, &depsMap)

			// Add field types
			fmt.Printf("Adding field types for model: %s\n", modelName) // Debugging
			builders.AddFieldTypesToModelData(&mData)

			// Add methods
			fmt.Printf("Adding methods for model: %s\n", modelName) // Debugging
			builders.AddMethodsToModelData(modelsASTData, &mData, &depsMap)

			// Setting imports
			var deps []string
			for dep := range depsMap {
				if dep == "" {
					continue
				}
				deps = append(deps, dep)
			}
			mData.Deps = deps

			// Writing to file
			fmt.Printf("Creating pool files for model: %s\n", modelName) // Debugging
			if err := createPoolFiles(dir, &mData); err != nil {
				errChan <- fmt.Errorf("failed to create pool files for model %s: %w", modelName, err)
				return
			}
			errChan <- nil
		}(mName, mASTData)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	fmt.Println("createPoolFilesFromASTData completed successfully.") // Debugging
	return nil
}

// loadModelFiles loads Go source files from the provided directories and parses them into *ast.File objects.
func loadModelFiles(dirs ...string) ([]*models.ModuleInfo, error) {
	var moduleFiles []*models.ModuleInfo
	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
				continue
			}
			filePath := filepath.Join(dir, file.Name())
			fset := token.NewFileSet()
			parsedFile, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
			if err != nil {
				return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
			}
			moduleFiles = append(moduleFiles, &models.ModuleInfo{
				Syntax: []*goast.File{parsedFile},
				FSet:   fset,
			})
		}
	}
	return moduleFiles, nil
}

// inflateMixins populates the given model with fields and methods defined in its mixins
func inflateMixins(modelName string, modelsData *map[string]ast.ModelASTData) {
	for mixin := range (*modelsData)[modelName].Mixins {
		inflateMixins(mixin, modelsData)
		for fieldName, field := range (*modelsData)[mixin].Fields {
			if fieldName == "ID" {
				continue
			}
			field.MixinField = true
			(*modelsData)[modelName].Fields[fieldName] = field
		}
		for methodName, method := range (*modelsData)[mixin].Methods {
			method.ToDeclare = true
			(*modelsData)[modelName].Methods[methodName] = method
		}
	}
}

// inflateEmbeds populates the given model with fields from the embedded type
func inflateEmbeds(modelName string, modelsData *map[string]ast.ModelASTData) {
	for emb := range (*modelsData)[modelName].Embeds {
		relModel := (*modelsData)[modelName].Fields[emb].RelModel
		inflateEmbeds(relModel, modelsData)
		for fieldName, field := range (*modelsData)[relModel].Fields {
			if _, exists := (*modelsData)[modelName].Fields[fieldName]; exists {
				continue
			}
			embeddedField := field
			embeddedField.EmbedField = true
			(*modelsData)[modelName].Fields[fieldName] = embeddedField
		}
	}
}

// createPoolFiles creates all pool files for the given model data
func createPoolFiles(dir string, mData *data.ModelData) error {
	fmt.Printf("Sorting model data for model: %s\n", mData.Name) // Debugging
	mData.Sort()

	// Helper function to create files from a template
	createFile := func(path, templateName string, mData *data.ModelData) error {
		fmt.Printf("Creating file: %s using template: %s\n", path, templateName) // Debugging
		tmpl, err := templates.GetTemplate(templateName)
		if err != nil {
			return fmt.Errorf("failed to get template %s: %w", templateName, err)
		}
		if err := templates.CreateFileFromTemplate(path, tmpl, mData); err != nil {
			return fmt.Errorf("failed to create file from template %s: %w", templateName, err)
		}
		return nil
	}

	// Create the model's interface file in interface directory
	interfaceFile := filepath.Join(dir, config.PoolInterfacesPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	fmt.Printf("Creating interface file for model: %s\n", mData.Name) // Debugging
	if err := createFile(interfaceFile, "PoolInterfacesTemplate", mData); err != nil {
		return err
	}

	// Create the models directory if not exists
	modelDir := filepath.Join(dir, config.PoolModelPackage, mData.SnakeName)
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		fmt.Printf("Creating models directory: %s\n", modelDir) // Debugging
		if err = os.MkdirAll(modelDir, 0755); err != nil {
			return fmt.Errorf("failed to create models directory for %s: %w", mData.Name, err)
		}
	}

	// Create the model's file in models directory
	modelFile := filepath.Join(dir, config.PoolModelPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	fmt.Printf("Creating model file for model: %s\n", mData.Name) // Debugging
	if err := createFile(modelFile, "PoolModelsTemplate", mData); err != nil {
		return err
	}

	// Create the model's file in model's subdirectory
	modelSubFile := filepath.Join(dir, config.PoolModelPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName))
	fmt.Printf("Creating model subdirectory file for model: %s\n", mData.Name) // Debugging
	if err := createFile(modelSubFile, "PoolModelsDirTemplate", mData); err != nil {
		return err
	}

	// Create the model's query directory if not exists
	queryDir := filepath.Join(dir, config.PoolQueryPackage, mData.SnakeName)
	if _, err := os.Stat(queryDir); os.IsNotExist(err) {
		fmt.Printf("Creating query directory: %s\n", queryDir) // Debugging
		if err = os.MkdirAll(queryDir, 0755); err != nil {
			return fmt.Errorf("failed to create query directory for %s: %w", mData.Name, err)
		}
	}

	// Create the model's query file in query directory
	queryFile := filepath.Join(dir, config.PoolQueryPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	fmt.Printf("Creating query file for model: %s\n", mData.Name) // Debugging
	if err := createFile(queryFile, "PoolQueryTemplate", mData); err != nil {
		return err
	}

	// Create the model's query file in model's query subdirectory
	querySubFile := filepath.Join(dir, config.PoolQueryPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName))
	fmt.Printf("Creating model's query subdirectory file for model: %s\n", mData.Name) // Debugging
	if err := createFile(querySubFile, "PoolModelsQueryTemplate", mData); err != nil {
		return err
	}

	fmt.Printf("Completed creating pool files for model: %s\n", mData.Name) // Debugging
	return nil
}
