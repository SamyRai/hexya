// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/tools/strutils"
)

// A fieldData describes a field in a RecordSet
type fieldData struct {
	Name        string
	JSON        string
	RelModel    string
	Type        string
	IType       string
	TypeWrapper string
	SanType     string
	ImportPath  string
	IsRS        bool
	MixinField  bool
	EmbedField  bool
}

// A methodData describes a method in a RecordSet
type methodData struct {
	Name             string
	Doc              string
	Params           string
	ParamsWithType   string
	IParamsWithTypes string
	ParamsTypes      string
	ReturnAsserts    string
	Returns          string
	ReturnString     string
	IReturnString    string
	Call             string
	ToDeclare        bool
}

// an operatorDef defines an operator func
type operatorDef struct {
	Name  string
	Multi bool
}

// A fieldType holds the name and valid operators on a field type
type fieldType struct {
	Type      string
	SanType   string
	IsRS      bool
	Operators []operatorDef
}

// A modelData describes a RecordSet model
type modelData struct {
	Name                  string
	SnakeName             string
	ModelsPackageName     string
	QueryPackageName      string
	InterfacesPackageName string
	ModelType             string
	IsModelMixin          bool
	Deps                  []string
	RelModels             []string
	Fields                []fieldData
	Methods               []methodData
	AllMethods            []methodData
	ConditionFuncs        []string
	Types                 []fieldType
	TypesDeps             []string
}

// sort sorts all slices fields of this modelData so that the generated code is always the same.
func (m *modelData) sort() {
	sort.Strings(m.Deps)
	sort.Slice(m.Fields, func(i, j int) bool {
		return m.Fields[i].Name < m.Fields[j].Name
	})
	sort.Slice(m.Methods, func(i, j int) bool {
		return m.Methods[i].Name < m.Methods[j].Name
	})
	sort.Slice(m.AllMethods, func(i, j int) bool {
		return m.AllMethods[i].Name < m.AllMethods[j].Name
	})
	sort.Strings(m.RelModels)
	sort.Slice(m.Types, func(i, j int) bool {
		return m.Types[i].Type < m.Types[j].Type
	})
}

// createTypeIdent creates a string from the given type that
// can be used inside an identifier.
func createTypeIdent(typStr string) string {
	res := strings.Replace(typStr, ".", "", -1)
	res = strings.Replace(res, "[", "Slice", -1)
	res = strings.Replace(res, "map[", "Map", -1)
	res = strings.Replace(res, "]", "", -1)
	res = strings.Title(res)
	return res
}

// trimInterfacePackagePrefix removes the 'm.' prefix from types
func trimInterfacePackagePrefix(typ string) string {
	toks := strings.Split(typ, "]")
	lastTok := strings.TrimPrefix(toks[len(toks)-1], PoolInterfacesPackage+".")
	toks = append(toks[:len(toks)-1], lastTok)
	return strings.Join(toks, "]")
}

// CreatePool generates the pool package by parsing the source code AST
// of the given program.
// The generated package will be put in the given dir.
func CreatePool(modules []*ModuleInfo, dir string) error {
	modelsASTData := GetModelsASTData(modules)
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(modelsASTData))
	wg.Add(len(modelsASTData))
	for mName, mASTData := range modelsASTData {
		for methToADD := range methodsToAdd {
			mASTData.Methods[methToADD] = MethodASTData{}
		}
		go func(modelName string, modelASTData ModelASTData) {
			defer wg.Done()
			depsMap := map[string]bool{ModelsPath: true}
			mData := modelData{
				Name:                  modelName,
				SnakeName:             strutils.SnakeCase(modelName),
				ModelsPackageName:     PoolModelPackage,
				QueryPackageName:      PoolQueryPackage,
				InterfacesPackageName: PoolInterfacesPackage,
				ModelType:             modelASTData.ModelType,
				IsModelMixin:          modelASTData.IsModelMixin,
				ConditionFuncs:        []string{"And", "AndNot", "Or", "OrNot"},
			}
			// Add fields
			addFieldsToModelData(modelASTData, &mData, &depsMap)
			// Add field types
			addFieldTypesToModelData(&mData)
			// Add methods
			addMethodsToModelData(modelsASTData, &mData, &depsMap)
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
	return nil
}

// addMethodsToModelData extracts data from modelsASTData to populate methods in modelData
func addMethodsToModelData(modelsASTData map[string]ModelASTData, modelData *modelData, depsMap *map[string]bool) {
	modelASTData := modelsASTData[modelData.Name]
	for methodName, methodASTData := range modelASTData.Methods {
		if handler, exists := specificMethodsHandlers[methodName]; exists {
			handler(&methodASTData, modelData, depsMap)
			continue
		}
		var params, paramsWithType, iParamsWithType, paramsType, call, returns, returnAsserts, returnString, iReturnString string
		for _, astParam := range methodASTData.Params {
			paramType := astParam.Type.Type
			iParamType := trimInterfacePackagePrefix(paramType)
			p := fmt.Sprintf("%s,", astParam.Name)
			if isRS, _ := isRecordSetType(astParam.Type.Type, modelsASTData); isRS {
				iParamType = fmt.Sprintf("%sSet", modelData.Name)
				paramType = fmt.Sprintf("%s.%sSet", PoolInterfacesPackage, modelData.Name)
			}
			if astParam.Variadic {
				iParamType = fmt.Sprintf("...%s", iParamType)
				paramType = fmt.Sprintf("...%s", paramType)
			}
			params += p
			paramsWithType += fmt.Sprintf("%s %s,", astParam.Name, paramType)
			iParamsWithType += fmt.Sprintf("%s %s,", astParam.Name, iParamType)
			paramsType += fmt.Sprintf("%s,", paramType)
			(*depsMap)[astParam.Type.ImportPath] = true
		}
		if len(methodASTData.Returns) == 1 {
			call = "Call"
			(*depsMap)[methodASTData.Returns[0].ImportPath] = true
			typ := methodASTData.Returns[0].Type
			iTyp := trimInterfacePackagePrefix(typ)
			returnAsserts = fmt.Sprintf("resTyped, _ := res.(%s)", typ)
			returns = "resTyped"
			if isRS, _ := isRecordSetType(typ, modelsASTData); isRS {
				typ = fmt.Sprintf("%s.%sSet", PoolInterfacesPackage, modelData.Name)
				iTyp = fmt.Sprintf("%sSet", modelData.Name)
				returnAsserts = fmt.Sprintf("resTyped := res.(models.RecordSet).Collection().Wrap(\"%s\").(%s)", modelData.Name, typ)
			}
			returnString = typ
			iReturnString = iTyp
		} else if len(methodASTData.Returns) > 1 {
			for i, ret := range methodASTData.Returns {
				typ := ret.Type
				iTyp := trimInterfacePackagePrefix(typ)
				call = "CallMulti"
				(*depsMap)[ret.ImportPath] = true
				if isRS, _ := isRecordSetType(ret.Type, modelsASTData); isRS {
					typ = fmt.Sprintf("%s.%sSet", PoolInterfacesPackage, modelData.Name)
					iTyp = fmt.Sprintf("%sSet", modelData.Name)
					returnAsserts += fmt.Sprintf("resTyped%d := res[%d].(models.RecordSet).Collection().Wrap(\"%s\").(%s)\n", i, i, modelData.Name, typ)
				} else {
					returnAsserts += fmt.Sprintf("resTyped%d, _ := res[%d].(%s)\n", i, i, typ)
				}
				returnString += fmt.Sprintf("%s,", typ)
				iReturnString += fmt.Sprintf("%s,", iTyp)
				returns += fmt.Sprintf("resTyped%d,", i)
			}
		}
		modelData.AllMethods = append(modelData.AllMethods, methodData{
			Name:             methodName,
			Doc:              methodASTData.Doc,
			ToDeclare:        methodASTData.ToDeclare,
			ParamsTypes:      strings.TrimRight(paramsType, ","),
			IParamsWithTypes: strings.TrimRight(iParamsWithType, ","),
			ReturnString:     strings.TrimSuffix(returnString, ","),
			IReturnString:    strings.TrimSuffix(iReturnString, ","),
		})
		modelData.Methods = append(modelData.Methods, methodData{
			Name:           methodName,
			Doc:            methodASTData.Doc,
			ToDeclare:      methodASTData.ToDeclare,
			Params:         strings.TrimRight(params, ","),
			ParamsWithType: strings.TrimRight(paramsWithType, ","),
			ReturnAsserts:  strings.TrimSuffix(returnAsserts, "\n"),
			Returns:        strings.TrimSuffix(returns, ","),
			ReturnString:   strings.TrimSuffix(returnString, ","),
			Call:           call,
		})
	}
}

// addFieldsToModelData extracts data from modelASTData to populate fields in modelData
func addFieldsToModelData(modelASTData ModelASTData, modelData *modelData, depsMap *map[string]bool) {
	relModels := make(map[string]bool)
	for fieldName, fieldASTData := range modelASTData.Fields {
		typStr := fieldASTData.Type.Type
		iTypStr := trimInterfacePackagePrefix(typStr)
		if fieldASTData.RelModel != "" {
			relModels[fieldASTData.RelModel] = true
			typStr = fmt.Sprintf("%s.%sSet", PoolInterfacesPackage, fieldASTData.RelModel)
			iTypStr = fmt.Sprintf("%sSet", fieldASTData.RelModel)
		}
		jsonName := strutils.GetDefaultString(fieldASTData.JSON, models.SnakeCaseFieldName(fieldName, fieldASTData.FType))
		modelData.Fields = append(modelData.Fields, fieldData{
			Name:       fieldName,
			JSON:       jsonName,
			Type:       typStr,
			IType:      iTypStr,
			IsRS:       fieldASTData.IsRS,
			RelModel:   fieldASTData.RelModel,
			SanType:    createTypeIdent(typStr),
			MixinField: fieldASTData.MixinField,
			EmbedField: fieldASTData.EmbedField,
			ImportPath: fieldASTData.Type.ImportPath,
		})
		(*depsMap)[fieldASTData.Type.ImportPath] = true
	}
	for rm := range relModels {
		modelData.RelModels = append(modelData.RelModels, rm)
	}
}

// addFieldTypesToModelData extracts field types from mData.Fields
// and add them to mData.Types
func addFieldTypesToModelData(mData *modelData) {
	fTypes := make(map[string]bool)
	tDeps := make(map[string]bool)
	for _, f := range mData.Fields {
		if fTypes[f.IType] {
			continue
		}
		fTypes[f.IType] = true
		tDeps[f.ImportPath] = true
		mData.Types = append(mData.Types, fieldType{
			Type:    f.IType,
			SanType: f.SanType,
			IsRS:    f.IsRS,
			Operators: []operatorDef{
				{Name: "Equals"}, {Name: "NotEquals"}, {Name: "Greater"}, {Name: "GreaterOrEqual"}, {Name: "Lower"},
				{Name: "LowerOrEqual"}, {Name: "Like"}, {Name: "Contains"}, {Name: "NotContains"}, {Name: "IContains"},
				{Name: "NotIContains"}, {Name: "ILike"}, {Name: "In", Multi: true}, {Name: "NotIn", Multi: true},
				{Name: "ChildOf"},
			},
		})
	}
	for dep := range tDeps {
		if dep == "" {
			continue
		}
		mData.TypesDeps = append(mData.TypesDeps, dep)
	}
}

// createPoolFiles creates all pool files for the given model data
func createPoolFiles(dir string, mData *modelData) error {
	mData.sort()
	// create the model's interface file in interface directory (m)
	fileName := filepath.Join(dir, PoolInterfacesPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	if err := CreateFileFromTemplate(fileName, poolInterfacesTemplate, mData); err != nil {
		return fmt.Errorf("failed to create interface file for %s: %w", mData.Name, err)
	}

	// create the models directories (h)
	if _, err := os.Stat(filepath.Join(dir, PoolModelPackage, mData.SnakeName)); err != nil {
		if err = os.MkdirAll(filepath.Join(dir, PoolModelPackage, mData.SnakeName), 0755); err != nil {
			return fmt.Errorf("failed to create models directory for %s: %w", mData.Name, err)
		}
	}
	// create the model's file in models directory (h)
	fileName = filepath.Join(dir, PoolModelPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	if err := CreateFileFromTemplate(fileName, poolModelsTemplate, mData); err != nil {
		return fmt.Errorf("failed to create model file for %s: %w", mData.Name, err)
	}
	// create the model's file in model's dir (q/model)
	fileName = filepath.Join(dir, PoolModelPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName))
	if err := CreateFileFromTemplate(fileName, poolModelsDirTemplate, mData); err != nil {
		return fmt.Errorf("failed to create model directory file for %s: %w", mData.Name, err)
	}

	// create the model's query directory (q)
	if _, err := os.Stat(filepath.Join(dir, PoolQueryPackage, mData.SnakeName)); err != nil {
		if err = os.MkdirAll(filepath.Join(dir, PoolQueryPackage, mData.SnakeName), 0755); err != nil {
			return fmt.Errorf("failed to create query directory for %s: %w", mData.Name, err)
		}
	}
	// create the model's query file in query dir (q)
	fileName = filepath.Join(dir, PoolQueryPackage, fmt.Sprintf("%s.go", mData.SnakeName))
	if err := CreateFileFromTemplate(fileName, poolQueryTemplate, mData); err != nil {
		return fmt.Errorf("failed to create query file for %s: %w", mData.Name, err)
	}
	// create the model's query file in model's query dir (q/model)
	fileName = filepath.Join(dir, PoolQueryPackage, mData.SnakeName, fmt.Sprintf("%s.go", mData.SnakeName))
	if err := CreateFileFromTemplate(fileName, poolModelsQueryTemplate, mData); err != nil {
		return fmt.Errorf("failed to create model query file for %s: %w", mData.Name, err)
	}

	return nil
}

// CreateFileFromTemplate generates a new file from the given template and data
func CreateFileFromTemplate(fileName string, template *template.Template, data interface{}) error {
	var srcBuffer bytes.Buffer
	err := template.Execute(&srcBuffer, data)
	if err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", fileName, err)
	}
	srcData, err := format.Source(srcBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format source for %s: %w", fileName, err)
	}
	// Write to file
	err = ioutil.WriteFile(fileName, srcData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileName, err)
	}
	return nil
}

// isRecordSetType returns true if the given typ is a RecordSet according
// to the AST data stored in models.
// The second returned value is true if typ is models.RecordCollection or models.RecordSet
// and false if it is a specific RecordSet type
func isRecordSetType(typ string, models map[string]ModelASTData) (bool, bool) {
	if typ == "*RecordCollection" || typ == "*models.RecordCollection" {
		return true, true
	}
	if typ == "RecordSet" || typ == "models.RecordSet" {
		return true, true
	}
	if _, exists := models[strings.TrimSuffix(typ, "Set")]; exists {
		return true, false
	}
	return false, false
}
