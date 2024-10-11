package builders

import (
	"github.com/hexya-erp/hexya/src/tools/generate/ast"
	"github.com/hexya-erp/hexya/src/tools/generate/data"
)

func getSpecificMethodsHandlers() map[string]func(*ast.MethodASTData, *data.ModelData, *map[string]bool) {
	return map[string]func(*ast.MethodASTData, *data.ModelData, *map[string]bool){
		// Add specific handlers here if required
	}
}
