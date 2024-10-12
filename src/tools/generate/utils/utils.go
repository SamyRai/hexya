package utils

import (
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

// TrimRecordSetSuffix Helper function to trim "Set" suffix from a type
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
	var nonEmptyStrs []string
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

// FormatDocString formats the given string by stripping whitespaces at the
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
// GetModulePackages returns a slice of PackageInfo for packages that are hexya modules, that is:
// - A package that declares a "MODULE_NAME" constant
// - A package that is in a subdirectory of a package
// Also returns the 'hexya/models' package since all models are initialized there
// GetModulePackages returns a list of module information along with an error, if any occurs.
// This function gathers packages with 'MODULE_NAME', handles core hexya models, and ensures syntax is loaded.
func GetModulePackages(packs []*packages.Package) ([]*models.ModuleInfo, error) {
	// Create a map to store found modules
	modules := make(map[string]*models.ModuleInfo)
	fmt.Println("[INFO] Starting to collect Hexya modules")

	// Visit each package to process and identify modules
	packages.Visit(packs, func(pack *packages.Package) bool {

		// Step 2: Ensure package contains syntax to be processed
		if len(pack.Syntax) == 0 {
			fmt.Printf("[WARNING] Skipping package %s: No syntax available.\n", pack.PkgPath)
			return true
		}

		// Step 3: Check for 'MODULE_NAME' constant which identifies a Hexya module
		obj := pack.Types.Scope().Lookup("MODULE_NAME")
		if obj != nil {
			// This package is a Hexya module
			fmt.Printf("[INFO] Found MODULE_NAME in package: %s\n", pack.PkgPath)
			modules[pack.Types.Path()] = NewModuleInfo(pack, models.Base, pack.Fset)
			return true
		}

		// Step 4: Check if this package matches Hexya core models path
		if pack.PkgPath == config.ModelsPath {
			fmt.Printf("[INFO] Found Hexya core models package: %s\n", pack.PkgPath)
			modules[pack.Types.Path()] = NewModuleInfo(pack, models.Models, pack.Fset)
			return true
		}

		// Step 5: Handle other addon or model packages
		if strings.Contains(pack.PkgPath, "/addons/") || strings.HasSuffix(pack.PkgPath, "src/models") {
			fmt.Printf("[INFO] Adding package by path match: %s\n", pack.PkgPath)
			modules[pack.Types.Path()] = NewModuleInfo(pack, models.Base, pack.Fset)
			return true
		}

		return true
	}, nil)

	// Return an error if no valid modules are found
	if len(modules) == 0 {
		return nil, fmt.Errorf("no valid modules found with MODULE_NAME or matching models path")
	}

	// Convert the map to a slice for return
	modSlice := make([]*models.ModuleInfo, 0, len(modules))
	for _, mod := range modules {
		modSlice = append(modSlice, mod)
	}

	fmt.Printf("[DEBUG] Collected Hexya modules: %v\n", modSlice)

	return modSlice, nil
}

// NewModuleInfo returns a pointer to a new moduleInfo instance
func NewModuleInfo(pack *packages.Package, modType models.PackageType, fSet *token.FileSet) *models.ModuleInfo {
	return &models.ModuleInfo{
		Package: *pack,
		ModType: modType,
		FSet:    fSet,
		Syntax:  pack.Syntax,
	}
}
