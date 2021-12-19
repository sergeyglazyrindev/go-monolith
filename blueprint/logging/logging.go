package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/logging/migrations"
	logmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/logging/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	logAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	logAdminPage.PageName = "Logs"
	logAdminPage.Slug = "log"
	logAdminPage.BlueprintName = "logging"
	logAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(logAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
	logmodelAdminPage := core.NewGormAdminPage(
		logAdminPage,
		&logmodel.Log{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"Username", "Action", "Activity", "CreatedAt", "ContentType", "ModelPK"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			// load extra css/js for logging admin page
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/go-monolith/assets/highlight.js/highlight.pack.js")
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/go-monolith/assets/js/initialize.highlight.js")
			form.ExtraStatic.ExtraCSS = append(form.ExtraStatic.ExtraCSS, "/static-inbuilt/go-monolith/assets/highlight.js/styles/default.css")
			actionField, _ := form.FieldRegistry.GetByName("Action")
			actionField.FieldConfig.Widget.SetPopulate(func(w core.IWidget, renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				a := renderContext.Model.(*logmodel.Log).Action
				return logmodel.HumanizeAction(a)
			})
			actionField.Ordering = 0
			usernameField, _ := form.FieldRegistry.GetByName("Username")
			usernameField.Ordering = 0
			activityField, _ := form.FieldRegistry.GetByName("Activity")
			// customize template for widget
			activityField.FieldConfig.Widget.SetTemplateName("admin/widgets/textareajson")
			activityField.Ordering = 9
			createdAtField, _ := form.FieldRegistry.GetByName("CreatedAt")
			// set CreatedAt field as readonly
			createdAtField.FieldConfig.Widget.SetReadonly(true)
			createdAtField.Ordering = 10
			return form
		},
	)
	logmodelAdminPage.PageName = "Logs"
	logmodelAdminPage.Slug = "log"
	logmodelAdminPage.BlueprintName = "logging"
	logmodelAdminPage.Router = app.GetRouter()
	contentTypeListDisplay, _ := logmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*logmodel.Log).ContentType.String()
	}
	actionListDisplay, _ := logmodelAdminPage.ListDisplay.GetFieldByDisplayName("Action")
	actionListDisplay.Populate = func(m interface{}) string {
		return logmodel.HumanizeAction(m.(*logmodel.Log).Action)
	}
	// this is internal data, so no need to add it as new.
	logmodelAdminPage.NoPermissionToAddNew = true
	err = logAdminPage.SubPages.AddAdminPage(logmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing log blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &logmodel.Log{}, &[]*logmodel.Log{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "logging",
		Description:       "Logging blueprint is responsible to store all actions made through admin panel",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
