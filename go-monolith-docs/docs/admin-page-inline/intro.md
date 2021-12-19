---
sidebar_position: 1
---

# Admin page inlines

Admin page inline is a way to show data related to the object on the model page.  
This is a structure with following fields:
```go
type AdminPageInline struct {
	Ordering          int
	GenerateModelI    func(m interface{}) (interface{}, interface{})
	GetQueryset       func(afo IAdminFilterObjects, model interface{}, rp *AdminRequestParams) IAdminFilterObjects
	Actions           *AdminModelActionRegistry
	EmptyValueDisplay string
	ExcludeFields     IFieldRegistry
	FieldsToShow      IFieldRegistry
	ShowAllFields     bool
	Validators        *ValidatorRegistry
	Classes           []string
	Extra             int
	MaxNum            int
	MinNum            int
	VerboseName       string
	VerboseNamePlural string
	ShowChangeLink    bool
	Template          string
	ContentType       *ContentType
	Permission        CustomPermission
	InlineType        InlineType
	Prefix            string
	ListDisplay       *ListDisplayRegistry `json:"-"`
}
```
The easiest way to create new admin page inline for Gorm is:
```go
abTestValueInline := core.NewAdminPageInline(
	"AB Test Values",
	// it could be be tabular or stacked inline.
	core.TabularInline, func(m interface{}) (interface{}, interface{}) {
		// here you return interface for one entity and for multiple
		if m != nil {
			mO := m.(*abtestmodel.ABTest)
			return &abtestmodel.ABTestValue{ABTestID: mO.ID}, &[]*abtestmodel.ABTestValue{}
		}
		return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
	}, func(afo core.IAdminFilterObjects, model interface{}, rp *core.AdminRequestParams) core.IAdminFilterObjects {
		// here you load data for your inline
		abTest := model.(*abtestmodel.ABTest)
		var db *core.Database
		if afo == nil {
			db = core.NewDatabaseInstance()
		} else {
			db = afo.(*core.GormAdminFilterObjects).Database
		}
		return &core.GormAdminFilterObjects{
			GormQuerySet:   core.NewGormPersistenceStorage(db.Db.Model(&abtestmodel.ABTestValue{}).Where(&abtestmodel.ABTestValue{ABTestID: abTest.ID})),
			Model:          &abtestmodel.ABTestValue{},
			Database: db,
			GenerateModelI: func() (interface{}, interface{}) {
				return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{}
			},
		}
	},
)
abTestValueInline.VerboseName = "AB Test Value"
// add custom fields to abTestValue inline
abTestValueInline.ListDisplay.AddField(&core.ListDisplay{
	DisplayName: "Click through rate",
	MethodName:  "ClickThroughRate",
})
abTestValueInline.ListDisplay.AddField(&core.ListDisplay{
	DisplayName: "Preview",
	MethodName:  "PreviewFormList",
})
abtestmodelAdminPage.InlineRegistry.Add(abTestValueInline)
err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
```
