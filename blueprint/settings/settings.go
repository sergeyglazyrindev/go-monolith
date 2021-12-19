package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/settings/migrations"
	settingmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/settings/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"mime/multipart"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	settingsAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	settingsAdminPage.PageName = "Settings"
	settingsAdminPage.Slug = "setting"
	settingsAdminPage.BlueprintName = "setting"
	settingsAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(settingsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingcategoriesmodelAdminPage := core.NewGormAdminPage(
		settingsAdminPage,
		&settingmodel.SettingCategory{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"Name", "Icon"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			return form
		},
	)
	settingcategoriesmodelAdminPage.PageName = "Setting categories"
	settingcategoriesmodelAdminPage.Slug = "settingcategory"
	settingcategoriesmodelAdminPage.BlueprintName = "setting"
	settingcategoriesmodelAdminPage.Router = app.GetRouter()
	err = settingsAdminPage.SubPages.AddAdminPage(settingcategoriesmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
	settingmodelAdminPage := core.NewGormAdminPage(
		settingsAdminPage,
		&settingmodel.Setting{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"Name", "DefaultValue", "DataType", "Value", "Help", "Category", "Code"}
			settingModel := modelI.(*settingmodel.Setting)
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/go-monolith/assets/js/settingformhandler.js")
			if settingModel.ID != 0 {
				codeField, _ := form.FieldRegistry.GetByName("Code")
				codeField.FieldConfig.Widget.SetReadonly(true)
			}
			dataTypeField, _ := form.FieldRegistry.GetByName("DataType")
			dataTypeWidget := dataTypeField.FieldConfig.Widget.(*core.SelectWidget)
			dataTypeWidget.OptGroups = make(map[string][]*core.SelectOptGroup)
			dataTypeWidget.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			//case 1:
			//return "string"
			//case 2:
			//return "integer"
			//case 3:
			//return "float"
			//case 4:
			//return "boolean"
			//case 5:
			//return "file"
			//case 6:
			//return "image"
			//case 7:
			//return "datetime"
			//default:
			//return "unknown"
			//}

			dataTypeWidget.Populate = func(w core.IWidget, renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				if renderContext.Context.GetCtx() != nil && renderContext.Context.GetCtx().Query("widgetType") != "" {
					dataType := settingmodel.DataTypeFromString(renderContext.Context.GetCtx().Query("widgetType"))
					return strconv.Itoa(int(dataType))
				}
				m1 := renderContext.Model.(*settingmodel.Setting)
				return strconv.Itoa(int(m1.DataType))
			}
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "String",
				Value:    "1",
				Selected: 1 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "Integer",
				Value:    "2",
				Selected: 2 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "Float",
				Value:    "3",
				Selected: 3 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "Boolean",
				Value:    "4",
				Selected: 4 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "File",
				Value:    "5",
				Selected: 5 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "Image",
				Value:    "6",
				Selected: 6 == settingModel.CategoryID,
			})
			dataTypeWidget.OptGroups[""] = append(dataTypeWidget.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "Datetime",
				Value:    "7",
				Selected: 7 == settingModel.CategoryID,
			})
			dataTypeField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				m1 := modelI.(*settingmodel.Setting)
				dataTypeI, _ := strconv.Atoi(v.(string))
				m1.DataType = settingmodel.DataType(dataTypeI)
				return nil
			}
			defaultValueField, _ := form.FieldRegistry.GetByName("DefaultValue")
			defaultValueWidget := defaultValueField.FieldConfig.Widget.(*core.DynamicWidget)
			defaultValueWidget.GetRealWidget = func(ctx2 *core.FormRenderContext, field *core.Field) core.IWidget {
				m2 := ctx2.Model.(*settingmodel.Setting)
				return m2.GetRealWidget()
			}
			defaultValueWidget.GetRealWidgetForFormProceeding = func(form *multipart.Form, afo core.IAdminFilterObjects) core.IWidget {
				dataTypeS := form.Value["DataType"][0]
				dataTypeI, _ := strconv.Atoi(dataTypeS)
				widgetTypeString := settingmodel.HumanizeDataType(settingmodel.DataType(dataTypeI))
				widgetType := core.GetWidgetByWidgetType(widgetTypeString, nil)
				return widgetType
			}
			valueField, _ := form.FieldRegistry.GetByName("Value")
			valueWidget := valueField.FieldConfig.Widget.(*core.DynamicWidget)
			valueWidget.GetRealWidget = func(ctx2 *core.FormRenderContext, field *core.Field) core.IWidget {
				m2 := ctx2.Model.(*settingmodel.Setting)
				return m2.GetRealWidget()
			}
			valueWidget.GetRealWidgetForFormProceeding = func(form *multipart.Form, afo core.IAdminFilterObjects) core.IWidget {
				dataTypeS := form.Value["DataType"][0]
				dataTypeI, _ := strconv.Atoi(dataTypeS)
				widgetTypeString := settingmodel.HumanizeDataType(settingmodel.DataType(dataTypeI))
				widgetType := core.GetWidgetByWidgetType(widgetTypeString, nil)
				return widgetType
			}
			categoryField, _ := form.FieldRegistry.GetByName("Category")
			initializedwidgetForCategory := categoryField.FieldConfig.Widget
			categoryWidget := &core.ForeignKeyWidget{}
			// configure foreign key widget for category field
			categoryWidget.GenerateModelInterface = func() (interface{}, interface{}) {
				return &settingmodel.SettingCategory{}, &[]*settingmodel.SettingCategory{}
			}
			categoryField.FieldConfig.Widget = categoryWidget
			categoryWidget.RenderForAdmin()
			//categoryWidget.Populate = func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
			//	m3 := renderContext.Model.(*settingmodel.Setting)
			//	return strconv.Itoa(int(m3.CategoryID))
			//}
			categoryWidget.SetRequired()
			categoryWidget.SetValue(settingModel.Category)
			categoryWidget.RenderUsingRenderer(initializedwidgetForCategory.GetRenderer())
			categoryWidget.SetName(initializedwidgetForCategory.GetName())
			categoryWidget.SetFieldDisplayName(initializedwidgetForCategory.GetFieldDisplayName())
			categoryWidget.AddNewLink = settingcategoriesmodelAdminPage.GenerateLinkToAddNewModel(true)
			categoryField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				m1 := modelI.(*settingmodel.Setting)
				categoryID, _ := strconv.Atoi(v.(string))
				var settingCategory settingmodel.SettingCategory
				afo.LoadDataForModelByID(categoryID, &settingCategory)
				m1.Category = settingCategory
				return nil
			}
			return form
		},
	)
	settingmodelAdminPage.PageName = "Settings"
	settingmodelAdminPage.Slug = "setting"
	settingmodelAdminPage.BlueprintName = "setting"
	settingmodelAdminPage.Router = app.GetRouter()
	dataTypeListDisplay, _ := settingmodelAdminPage.ListDisplay.GetFieldByDisplayName("DataType")
	dataTypeListDisplay.Populate = func(m interface{}) string {
		return settingmodel.HumanizeDataType(m.(*settingmodel.Setting).DataType)
	}
	// categoryListDisplay, _ := settingmodelAdminPage.ListDisplay.GetFieldByDisplayName("Category")
	// categoryListDisplay.Field.FieldConfig.Widget.SetReadonly(true)
	//categoryListDisplay.Populate = func(m interface{}) string {
	//	return m.(*settingmodel.Setting).Category.Name
	//}
	settingmodelAdminPage.NoPermissionToAddNew = true
	err = settingsAdminPage.SubPages.AddAdminPage(settingmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing settings blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &settingmodel.Setting{}, &[]*settingmodel.Setting{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) {
		return &settingmodel.SettingCategory{}, &[]*settingmodel.SettingCategory{}
	})
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "settings",
		Description:       "Settings blueprint responsible for wide-project settings",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
