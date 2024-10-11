package models

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
)

// A ModuleInfo is a wrapper around packages.Package with additional data to
// describe a module.
type ModuleInfo struct {
	packages.Package
	ModType PackageType
	FSet    *token.FileSet
	Syntax  []*ast.File // Added to handle the loaded syntax files (ASTs)
}

// A PackageType describes a type of module
type PackageType int8

const (
	// Base is the PackageType for the base package of a module
	Base PackageType = iota
	// Models is the PackageType for the hexya/models package
	Models
)

// ConvertPackagesToModules converts a list of *packages.Package to a list of *models.ModuleInfo.
func ConvertPackagesToModules(packages []*packages.Package) []*ModuleInfo {
	modules := make([]*ModuleInfo, len(packages))
	for i, pkg := range packages {
		modules[i] = &ModuleInfo{
			Package: *pkg,
			FSet:    pkg.Fset,
			Syntax:  pkg.Syntax, // Adding the AST files from the loaded package
		}
	}
	return modules
}

// GatherImportsFromModules gathers the necessary import paths from the provided modules.

// GatherCoreAndModuleImports splits core imports from module imports.
func GatherCoreAndModuleImports(modules []*ModuleInfo) (coreImports []string, moduleImports []string) {
	// Core imports
	coreImports = []string{
		"github.com/hexya-erp/hexya/cmd",
		"github.com/spf13/cobra",
	}

	// Module imports
	for _, module := range modules {
		if module.PkgPath != "" {
			moduleImports = append(moduleImports, module.PkgPath)
		}
	}

	return coreImports, moduleImports
}
