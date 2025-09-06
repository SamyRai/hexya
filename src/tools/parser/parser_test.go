package parser

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestGetModulePackages(t *testing.T) {
	conf := packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	packs, err := packages.Load(&conf, "github.com/hexya-erp/hexya/src/tests/testmodule")
	if err != nil {
		t.Fatalf("Error loading packages: %v", err)
	}
	if len(packs) == 0 {
		t.Fatal("No packages found")
	}

	mods := GetModulePackages(packs)
	if len(mods) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(mods))
	}
	if mods[0].Name != "testmodule" {
		t.Fatalf("Expected module name 'testmodule', got '%s'", mods[0].Name)
	}
}
