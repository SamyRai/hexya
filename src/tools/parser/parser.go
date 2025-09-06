package parser

import (
	"go/token"

	"golang.org/x/tools/go/packages"
)

// A PackageType describes a type of module
type PackageType int8

const (
	// Base is the PackageType for the base package of a module
	Base PackageType = iota
	// Models is the PackageType for the hexya/models package
	Models
)

// A ModuleInfo is a wrapper around packages.Package with additional data to
// describe a module.
type ModuleInfo struct {
	packages.Package
	ModType PackageType
	FSet    *token.FileSet
}

// NewModuleInfo returns a pointer to a new moduleInfo instance
func NewModuleInfo(pack *packages.Package, modType PackageType, fSet *token.FileSet) *ModuleInfo {
	return &ModuleInfo{
		Package: *pack,
		ModType: modType,
		FSet:    fSet,
	}
}

// GetModulePackages returns a slice of PackageInfo for packages that are hexya modules, that is:
// - A package that declares a "MODULE_NAME" constant
// - A package that is in a subdirectory of a package
// Also returns the 'hexya/models' package since all models are initialized there
func GetModulePackages(packs []*packages.Package) []*ModuleInfo {
	modules := make(map[string]*ModuleInfo)
	// We add to the modulePaths all packages which define a MODULE_NAME constant
	// and we check for 'hexya/models' package
	packages.Visit(packs, func(pack *packages.Package) bool {
		obj := pack.Types.Scope().Lookup("MODULE_NAME")
		if obj != nil {
			modules[pack.Types.Path()] = NewModuleInfo(pack, Base, pack.Fset)
			return true
		}
		if pack.PkgPath == "github.com/hexya-erp/hexya/src/models" {
			modules[pack.Types.Path()] = NewModuleInfo(pack, Models, pack.Fset)
		}
		return true
	}, func(pack *packages.Package) {})

	// Finally, we build up our result slice from modules map
	modSlice := make([]*ModuleInfo, len(modules))
	var i int
	for _, mod := range modules {
		modSlice[i] = mod
		i++
	}
	return modSlice
}
