package tests

import (
	"testing"
	"sort"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/fields"
	"github.com/hexya-erp/hexya/src/models/security"
	. "github.com/smartystreets/goconvey/convey"
)

type TestModelData struct {
	*models.ModelData
	Name string
	Age  int
}

func TestGenericTypes(t *testing.T) {
	Convey("Testing Generic Types", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			testModel := models.GenericModel[TestModelData]{
				Model: models.Registry.MustGet("TestModel"),
			}
			testData := testModel.NewData(models.FieldMap{
				"Name": "Test 1",
				"Age":  42,
			})
			record1 := testModel.Create(env, testData)

			testData2 := testModel.NewData(models.FieldMap{
				"Name": "Test 2",
				"Age":  24,
			})
			record2 := testModel.Create(env, testData2)

			Convey("Testing GenericModel", func() {
				So(testModel.IsValid(), ShouldBeTrue)
				So(testModel.NewSet(env).IsValid(), ShouldBeTrue)
				So(testModel.NewData().IsValid(), ShouldBeTrue)
				So(testModel.Browse(env, []int64{record1.Ids()[0], record2.Ids()[0]}).Len(), ShouldEqual, 2)
				So(testModel.BrowseOne(env, record1.Ids()[0]).Len(), ShouldEqual, 1)
			})
			Convey("Testing GenericRecordSet", func() {
				testSet := testModel.Search(env, testModel.Condition())
				So(testSet.IsValid(), ShouldBeTrue)
				So(testSet.Len(), ShouldEqual, 2)

				Convey("First", func() {
					first := testSet.First()
					So(first.Name, ShouldBeIn, []string{"Test 1", "Test 2"})
				})
				Convey("All", func() {
					all := testSet.All()
					So(len(all), ShouldEqual, 2)
				})
				Convey("Records", func() {
					records := testSet.Records()
					So(len(records), ShouldEqual, 2)
				})
				Convey("Sorted", func() {
					sorted := testSet.Sorted(func(a, b TestModelData) bool {
						return a.Name < b.Name
					})
					So(sorted.Records()[0].First().Name, ShouldEqual, "Test 1")
					So(sorted.Records()[1].First().Name, ShouldEqual, "Test 2")
				})
				Convey("Filtered", func() {
					filtered := testSet.Filtered(func(r TestModelData) bool {
						return r.Name == "Test 1"
					})
					So(filtered.Len(), ShouldEqual, 1)
					So(filtered.First().Name, ShouldEqual, "Test 1")
				})
			})
			Convey("Testing GenericRecordData", func() {
				testData := testModel.NewData()
				So(testData.IsValid(), ShouldBeTrue)
			})
		}), ShouldBeNil)
	})
}

func init() {
	models.NewModel("TestModel").AddFields(map[string]models.FieldDefinition{
		"Name": fields.Char{},
		"Age":  fields.Integer{},
	})
}
