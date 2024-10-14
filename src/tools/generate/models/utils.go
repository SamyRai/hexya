package models

import (
	"sort"
	"strings"
	"unicode"
)

// Helper Functions

// removeDuplicateImports removes duplicate import paths.
func removeDuplicateImports(imports map[string]bool) []string {
	var uniqueImports []string
	for imp := range imports {
		uniqueImports = append(uniqueImports, imp)
	}
	sort.Strings(uniqueImports)
	return uniqueImports
}

// formatParams formats the parameters of a method into a string.
func formatParams(params []ParamAST) string {
	var paramStrs []string
	for _, param := range params {
		paramStr := param.Name + " " + param.Type.TypeName
		if param.IsVariadic {
			paramStr = "..." + paramStr
		}
		paramStrs = append(paramStrs, paramStr)
	}
	return strings.Join(paramStrs, ", ")
}

// formatReturns formats the return types of a method into a string.
func formatReturns(returns []ReturnAST) string {
	var returnStrs []string
	for _, ret := range returns {
		returnStrs = append(returnStrs, ret.Type.TypeName)
	}
	return strings.Join(returnStrs, ", ")
}

// getImportPathsFromDependencies retrieves import paths from method or field dependencies.
func getImportPathsFromDependencies(dependencies []DependencyNode) []string {
	var importPaths []string
	for _, dep := range dependencies {
		importPaths = append(importPaths, dep.ImportPath)
	}
	return importPaths
}

// SnakeName converts the ModelData's name into snake_case.
func (m *ModelData) SnakeName() string {
	return toSnakeCase(m.Name)
}

// toSnakeCase converts a CamelCase string into snake_case.
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// Sort organizes fields and methods for consistent code generation output
func (m *ModelData) Sort() {
	// Sort fields by their name
	sort.Slice(m.Fields, func(i, j int) bool {
		return m.Fields[i].Name < m.Fields[j].Name
	})

	// Sort methods by their name
	sort.Slice(m.Methods, func(i, j int) bool {
		return m.Methods[i].Name < m.Methods[j].Name
	})
}

// InflateModels inflates mixins and embedded models into their respective models.
func InflateModels(modelsData map[string]*ModelData) {
	for _, modelData := range modelsData {
		InflateFieldsAndMethods(modelData, modelData.Mixins)
		InflateFieldsAndMethods(modelData, modelData.EmbeddedModels)
	}
}

// InflateFieldsAndMethods inflates the fields and methods of source models into the target model.
func InflateFieldsAndMethods(targetModel *ModelData, sourceModels []*ModelData) {
	for _, sourceModel := range sourceModels {
		for _, field := range sourceModel.Fields {
			if !targetModel.HasField(field.Name) {
				targetModel.Fields = append(targetModel.Fields, field)
			}
		}
		for _, method := range sourceModel.Methods {
			if !targetModel.HasMethod(method.Name) {
				targetModel.Methods = append(targetModel.Methods, method)
			}
		}
	}
}
