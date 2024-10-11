package utils

import (
	"errors"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/generate/config"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
	"go/token"
	"golang.org/x/tools/go/packages"
	"strings"
)

// CreateTypeIdent creates a string from the given type that can be used inside an identifier.
func CreateTypeIdent(typStr string) string {
	res := strings.Replace(typStr, ".", "", -1)
	res = strings.Replace(res, "[", "Slice", -1)
	res = strings.Replace(res, "map[", "Map", -1)
	res = strings.Replace(res, "]", "", -1)
	res = CapitalizeFirst(res)
	return res
}

// TrimInterfacePackagePrefix removes the 'm.' prefix from types
func TrimInterfacePackagePrefix(typ string) string {
	toks := strings.Split(typ, "]")
	lastTok := strings.TrimPrefix(toks[len(toks)-1], config.PoolInterfacesPackage+".") // Updated reference
	toks = append(toks[:len(toks)-1], lastTok)
	return strings.Join(toks, "]")
}

// CapitalizeFirst capitalizes the first letter of the given string.
func CapitalizeFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

// Helper function to trim "Set" suffix from a type
func TrimRecordSetSuffix(typ string) string {
	return strings.TrimSuffix(typ, "Set")
}

// AddImport adds an import path to the dependency map if it's not empty.
func AddImport(depsMap *map[string]bool, importPath string) {
	if importPath != "" {
		(*depsMap)[importPath] = true
	}
}

// JoinStrings joins multiple strings with a given separator, ignoring empty strings.
func JoinStrings(strs []string, sep string) string {
	nonEmptyStrs := []string{}
	for _, str := range strs {
		if str != "" {
			nonEmptyStrs = append(nonEmptyStrs, str)
		}
	}
	return strings.Join(nonEmptyStrs, sep)
}

// TrimTrailingNewline trims any trailing newline character from the given string.
func TrimTrailingNewline(str string) string {
	return strings.TrimRight(str, "\n")
}

// TrimTrailingComma trims any trailing comma from the given string.
func TrimTrailingComma(str string) string {
	return strings.TrimRight(str, ",")
}

// formatDocString formats the given string by stripping whitespaces at the
// beginning of each line and prepend "// ". It also strips empty lines at
// the beginning.
func FormatDocString(doc string) string {
	var res string
	var dataStarted bool
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if line == "" && !dataStarted {
			continue
		}
		dataStarted = true
		res += fmt.Sprintf("// %s\n", line)
	}
	return strings.TrimRight(res, "/ \n")
}

// GetModulePackages returns a slice of PackageInfo for packages that are hexya modules, that is:
// - A package that declares a "MODULE_NAME" constant
// - A package that is in a subdirectory of a package
// Also returns the 'hexya/models' package since all models are initialized there

// GetModulePackages returns a list of module information along with an error, if any occurs.
func GetModulePackages(packs []*packages.Package) ([]*models.ModuleInfo, error) {
	modules := make(map[string]*models.ModuleInfo)

	// Iterate over packages to find modules and handle potential errors
	packages.Visit(packs, func(pack *packages.Package) bool {
		fmt.Printf("Checking package: %s\n", pack.PkgPath)

		// Logging all errors encountered by the package
		if len(pack.Errors) > 0 {
			fmt.Printf("Errors encountered in package %s:\n", pack.PkgPath)
			for _, err := range pack.Errors {
				fmt.Printf("  Error: %v\n", err)
			}
		}

		obj := pack.Types.Scope().Lookup("MODULE_NAME")
		if obj != nil {
			fmt.Printf("Found MODULE_NAME in package: %s\n", pack.PkgPath)
			modules[pack.PkgPath] = NewModuleInfo(pack, config.Base, pack.Fset)
			return true
		} else {
			fmt.Printf("MODULE_NAME not found in package: %s\n", pack.PkgPath)
		}

		return true
	}, nil)

	// Check if any modules were found
	if len(modules) == 0 {
		return nil, errors.New("no modules found with MODULE_NAME constant or matching models path")
	}

	// Build up the result slice from the modules map
	modSlice := make([]*models.ModuleInfo, len(modules))
	var i int
	for _, mod := range modules {
		modSlice[i] = mod
		i++
	}

	return modSlice, nil
}

// NewModuleInfo returns a pointer to a new moduleInfo instance
func NewModuleInfo(pack *packages.Package, modType models.PackageType, fSet *token.FileSet) *models.ModuleInfo {
	return &models.ModuleInfo{
		Package: *pack,
		ModType: modType,
		FSet:    fSet,
	}
}
