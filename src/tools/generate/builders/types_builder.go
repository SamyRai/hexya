package builders

import (
	"github.com/hexya-erp/hexya/src/tools/generate/data"
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// AddFieldTypesToModelData extracts field types from mData.Fields and adds them to mData.Types
func AddFieldTypesToModelData(mData *data.ModelData) {
	fTypes := make(map[string]bool)
	tDeps := make(map[string]bool)
	for _, f := range mData.Fields {
		if fTypes[f.IType] {
			continue
		}
		fTypes[f.IType] = true
		tDeps[f.ImportPath] = true
		mData.Types = append(mData.Types, models.FieldType{
			Type:    f.IType,
			SanType: f.SanType,
			IsRS:    f.IsRS,
			Operators: []models.OperatorDef{
				{Name: "Equals"}, {Name: "NotEquals"}, {Name: "Greater"}, {Name: "GreaterOrEqual"},
				{Name: "Lower"}, {Name: "LowerOrEqual"}, {Name: "Like"}, {Name: "Contains"},
				{Name: "NotContains"}, {Name: "IContains"}, {Name: "NotIContains"}, {Name: "ILike"},
				{Name: "In", Multi: true}, {Name: "NotIn", Multi: true}, {Name: "ChildOf"},
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
