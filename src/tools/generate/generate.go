package generate

import (
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/builders"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/file_operations"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"github.com/hexya-erp/hexya/src/tools/strutils"
	"golang.org/x/tools/go/packages"
)

// CreatePool generates the pool package by parsing the source code AST of the given program.
func CreatePool(modelsASTData map[string]ast.ModelASTData, dir string) error {
	fmt.Println("Starting createPoolFilesFromASTData...")

	// First, process predefined mixins from config.ModelMixins
	for mixinName := range config.ModelMixins {
		if modelASTData, exists := modelsASTData[mixinName]; exists {
			fmt.Printf("Processing mixin model: %s\n", mixinName)
			if err := processModel(mixinName, modelASTData, modelsASTData, dir); err != nil {
				return fmt.Errorf("failed to process mixin %s: %w", mixinName, err)
			}
		} else {
			// Ensure predefined mixins are initialized even if they don't exist in AST data
			fmt.Printf("Initializing predefined mixin: %s\n", mixinName)
			modelsASTData[mixinName] = ast.ModelASTData{
				Name:         mixinName,
				IsModelMixin: true,
				Fields:       map[string]ast.FieldASTData{},
				Methods:      map[string]ast.MethodASTData{},
				Mixins:       map[string]bool{},
				Embeds:       map[string]bool{},
				Validated:    true,
			}
			if err := processModel(mixinName, modelsASTData[mixinName], modelsASTData, dir); err != nil {
				return fmt.Errorf("failed to process predefined mixin %s: %w", mixinName, err)
			}
		}
	}

	// Process remaining models after mixins
	for modelName, modelASTData := range modelsASTData {
		if _, isMixin := config.ModelMixins[modelName]; isMixin {
			// Skip models that are already processed as mixins
			continue
		}

		fmt.Printf("Processing regular model: %s\n", modelName)
		if !modelASTData.Validated {
			fmt.Printf("Skipping unvalidated model: %s\n", modelName)
			continue
		}

		if err := processModel(modelName, modelASTData, modelsASTData, dir); err != nil {
			return fmt.Errorf("failed to process model %s: %w", modelName, err)
		}
	}

	fmt.Println("createPoolFilesFromASTData completed successfully.")
	return nil
}

// processModel processes each model, inflating mixins and embeddings, and creates pool files
func processModel(modelName string, modelASTData ast.ModelASTData, modelsASTData map[string]ast.ModelASTData, dir string) error {

	// Prepare model data
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
	builders.AddFieldsToModelData(modelASTData, &mData, &depsMap)

	// Add field types
	builders.AddFieldTypesToModelData(&mData)

	// Add methods (including those added by mixins and MethodsToAdd)
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
	if err := file_operations.CreatePoolFiles(dir, &mData); err != nil {
		return fmt.Errorf("failed to create pool files for model %s: %w", modelName, err)
	}

	return nil
}

// LoadModelFiles uses the already loaded packages to gather model information
func LoadModelFiles(packs []*packages.Package) ([]*models.ModuleInfo, error) {
	var moduleFiles []*models.ModuleInfo

	for _, pack := range packs {
		if len(pack.Errors) > 0 {
			for _, err := range pack.Errors {
				fmt.Printf("Error in package %s: %v\n", pack.PkgPath, err)
			}
			return nil, fmt.Errorf("errors encountered in package %s", pack.PkgPath)
		}

		// Skip packages without syntax information
		if len(pack.Syntax) == 0 {
			fmt.Printf("Skipping package %s: No syntax available.\n", pack.PkgPath)
			continue
		}

		// Add the package's syntax and other relevant data to the module info
		moduleFiles = append(moduleFiles, &models.ModuleInfo{
			Syntax: pack.Syntax,
			FSet:   pack.Fset,
		})
	}

	// Return all gathered module files
	return moduleFiles, nil
}
