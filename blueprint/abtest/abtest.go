package abtest

/*
	Blueprint abtest is designed to do A/B tests for the project. Currently not usable after migration from previous
	GoMonolith implementation.
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/abtest/migrations"
	abtestmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/abtest/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	// add new gorm admin page for blueprint
	abTestAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	abTestAdminPage.PageName = "AB Tests"
	abTestAdminPage.Slug = "abtest"
	abTestAdminPage.BlueprintName = "abtest"
	abTestAdminPage.Router = app.GetRouter()

	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(abTestAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
	// add abtest gorm page
	abtestmodelAdminPage := core.NewGormAdminPage(
		abTestAdminPage,
		&abtestmodel.ABTest{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"ContentType", "Type", "Name", "Field", "PrimaryKey", "Active", "Group", "StaticPath"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/go-monolith/assets/js/abtestformhandler.js")
			typeField, _ := form.FieldRegistry.GetByName("Type")
			w := typeField.FieldConfig.Widget.(*core.SelectWidget)
			w.OptGroups = make(map[string][]*core.SelectOptGroup)
			w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "unknown",
				Value:    "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "static",
				Value:    "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "model",
				Value:    "2",
			})
			typeField.FieldConfig.Widget.SetPopulate(func(w core.IWidget, renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				a := renderContext.Model.(*abtestmodel.ABTest).Type
				return strconv.Itoa(int(a))
			})
			typeField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				abTestM := m.(*abtestmodel.ABTest)
				vI, _ := strconv.Atoi(v.(string))
				abTestM.Type = abtestmodel.ABTestType(vI)
				return nil
			}
			contentTypeField, _ := form.FieldRegistry.GetByName("ContentType")
			w1 := contentTypeField.FieldConfig.Widget.(*core.ContentTypeSelectorWidget)
			w1.LoadFieldsOfAllModels = true
			fieldField, _ := form.FieldRegistry.GetByName("Field")
			w2 := fieldField.FieldConfig.Widget.(*core.SelectWidget)
			w2.SetAttr("data-initialized", "false")
			w2.DontValidateForExistence = true
			return form
		},
	)
	abtestmodelAdminPage.PageName = "AB Tests"
	abtestmodelAdminPage.Slug = "abtest"
	abtestmodelAdminPage.BlueprintName = "abtest"
	abtestmodelAdminPage.Router = app.GetRouter()
	// customize list display for abtest admin page
	typeListDisplay, _ := abtestmodelAdminPage.ListDisplay.GetFieldByDisplayName("Type")
	// custom populate method for type list display
	typeListDisplay.Populate = func(m interface{}) string {
		return abtestmodel.HumanizeAbTestType(m.(*abtestmodel.ABTest).Type)
	}
	contentTypeListDisplay, _ := abtestmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	// custom populate method for contentType list display
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*abtestmodel.ABTest).ContentType.String()
	}
	// initialize inline for abtest, it shows all abtest values that belong to the current abtest object
	abTestValueInline := core.NewAdminPageInline(
		"AB Test Values",
		core.TabularInline, func(m interface{}) interface{} {
			if m != nil {
				mO := m.(*abtestmodel.ABTest)
				return &abtestmodel.ABTestValue{ABTestID: mO.ID}
			}
			return &abtestmodel.ABTestValue{}
		}, func(adminContext core.IAdminContext, afo core.IAdminFilterObjects, model interface{}) core.IAdminFilterObjects {
			abTest := model.(*abtestmodel.ABTest)
			var db *core.ProjectDatabase
			if afo == nil {
				db = core.NewDatabaseInstance()
			} else {
				db = afo.(*core.GormAdminFilterObjects).DatabaseInstance
			}
			return &core.GormAdminFilterObjects{
				GormQuerySet:     core.NewGormPersistenceStorage(db.Db.Model(&abtestmodel.ABTestValue{}).Where(&abtestmodel.ABTestValue{ABTestID: abTest.ID})),
				Model:            &abtestmodel.ABTestValue{},
				DatabaseInstance: db,
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
	if err != nil {
		panic(fmt.Errorf("error initializing abtest blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	// add models to ProjectModels, so we can determine admin pages for these models, etc
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &abtestmodel.ABTestValue{}, &[]*abtestmodel.ABTestValue{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &abtestmodel.ABTest{}, &[]*abtestmodel.ABTest{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "abtest",
		Description:       "ABTest blueprint is responsible for ab tests",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
