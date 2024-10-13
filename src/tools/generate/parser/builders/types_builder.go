package builders

import (
	"github.com/hexya-erp/hexya/src/tools/generate/models"
)

// AddFieldTypesToModelData extracts field types from mData.Fields and adds them to mData.Types.
// It now uses the DepsManager for managing dependencies.
func AddFieldTypesToModelData(mData *models.ModelData) {
	fTypes := make(map[string]bool)

	// Iterate over the fields and gather types and dependencies
	for _, f := range mData.Fields {
		if fTypes[f.IType] {
			continue
		}
		fTypes[f.IType] = true

		// Add the field type to the model data
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

		// Add the field's import path to dependencies using DepsManager

	}
}
