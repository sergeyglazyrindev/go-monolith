package approval

/*
	Blueprint approval is designed to provide Approval functionality for administrators. Currently not usable after
	migration from previous UAdmin implementation.
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/approval/migrations"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/approval/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	approvalAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	approvalAdminPage.PageName = "Approvals"
	approvalAdminPage.Slug = "approval"
	approvalAdminPage.BlueprintName = "approval"
	approvalAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(approvalAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
	approvalmodelAdminPage := core.NewGormAdminPage(
		approvalAdminPage,
		&models.Approval{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"ContentType", "ModelPK", "ColumnName", "OldValue", "NewValue", "NewValueDescription", "ChangedBy", "ChangeDate", "ApprovalAction", "ApprovalBy", "ApprovalDate"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			approvalField, _ := form.FieldRegistry.GetByName("ApprovalAction")
			w := approvalField.FieldConfig.Widget.(*core.SelectWidget)
			// initialize select options for ApprovalAction field
			w.OptGroups = make(map[string][]*core.SelectOptGroup)
			w.OptGroups[""] = make([]*core.SelectOptGroup, 0)
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "unknown",
				Value:    "0",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "approved",
				Value:    "1",
			})
			w.OptGroups[""] = append(w.OptGroups[""], &core.SelectOptGroup{
				OptLabel: "rejected",
				Value:    "2",
			})
			// populate ApprovalAction field value
			approvalField.FieldConfig.Widget.SetPopulate(func(w core.IWidget, renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
				a := renderContext.Model.(*models.Approval).ApprovalAction
				return strconv.Itoa(int(a))
			})
			// setup ApprovalAction field during saving
			approvalField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				approvalM := m.(*models.Approval)
				vI, _ := strconv.Atoi(v.(string))
				approvalM.ApprovalAction = models.ApprovalAction(vI)
				return nil
			}
			return form
		},
	)
	approvalmodelAdminPage.PageName = "Approval"
	approvalmodelAdminPage.Slug = "approval"
	approvalmodelAdminPage.BlueprintName = "approval"
	approvalmodelAdminPage.Router = app.GetRouter()
	approvalActionListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ApprovalAction")
	approvalActionListDisplay.Populate = func(m interface{}) string {
		return models.HumanizeApprovalAction(m.(*models.Approval).ApprovalAction)
	}
	approvalDateListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ApprovalDate")
	approvalDateListDisplay.Populate = func(m interface{}) string {
		aD := m.(*models.Approval).ApprovalDate
		if aD == nil {
			return ""
		}
		return aD.Format(core.CurrentConfig.D.GoMonolith.DateTimeFormat)
	}
	contentTypeListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ContentType")
	contentTypeListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ContentType.String()
	}
	changeDateListDisplay, _ := approvalmodelAdminPage.ListDisplay.GetFieldByDisplayName("ChangeDate")
	changeDateListDisplay.Populate = func(m interface{}) string {
		return m.(*models.Approval).ChangeDate.Format(core.CurrentConfig.D.GoMonolith.DateTimeFormat)
	}
	err = approvalAdminPage.SubPages.AddAdminPage(approvalmodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing approval blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &models.Approval{}, &[]*models.Approval{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "approval",
		Description:       "Approval blueprint is responsible for approving things in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
