package language

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/language/migrations"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"mime/multipart"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	languageAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	languageAdminPage.PageName = "Languages"
	languageAdminPage.Slug = "language"
	languageAdminPage.BlueprintName = "language"
	languageAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(languageAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
	languagemodelAdminPage := core.NewGormAdminPage(
		languageAdminPage,
		&core.Language{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"EnglishName", "Name", "Flag", "Code", "RTL", "Default", "Active", "AvailableInGui"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			defaultField, _ := form.FieldRegistry.GetByName("Default")
			// add custom validator for defaultField
			defaultField.Validators.AddValidator("only_one_default_language", func(i interface{}, o interface{}) error {
				isDefault := i.(bool)
				if !isDefault {
					return nil
				}
				d := o.(*multipart.Form)
				ID := d.Value["ID"][0]
				database := core.NewDatabaseInstance()
				lang := &core.Language{}
				database.Db.Where(&core.Language{Default: true}).First(lang)
				if lang.ID != 0 && ID != strconv.Itoa(int(lang.ID)) {
					return core.NewHTTPErrorResponse("only_one_default_language_allowed", "only one default language could be configured")
				}
				return nil
			})
			return form
		},
	)
	languagemodelAdminPage.PageName = "Languages"
	languagemodelAdminPage.Slug = "language"
	languagemodelAdminPage.BlueprintName = "language"
	languagemodelAdminPage.Router = app.GetRouter()
	languagemodelAdminPage.NoPermissionToAddNew = true
	err = languageAdminPage.SubPages.AddAdminPage(languagemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.Language{}, &[]*core.Language{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
