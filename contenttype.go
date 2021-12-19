package gomonolith

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"os"
)

type ContentTypeCommand struct {
}

func (c ContentTypeCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("sync", &SyncContentTypes{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.IsRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	commandRegistry.RunAction(subaction, "", args)
	return nil
}

func (c ContentTypeCommand) GetHelpText() string {
	return "Content type for GoMonolith project"
}

type SyncContentTypes struct {
}

func (command SyncContentTypes) Proceed(subaction string, args []string) error {
	appInstance.InitializeRouter()
	database := core.NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	var contentType core.ContentType
	var permission core.Permission
	for blueprintRootAdminPage := range core.CurrentDashboardAdminPanel.AdminPages.GetAll() {
		core.Trail(core.INFO, "Sync content types for blueprint %s", blueprintRootAdminPage.BlueprintName)
		for modelPage := range blueprintRootAdminPage.SubPages.GetAll() {
			if modelPage.Model == nil {
				continue
			}
			core.Trail(core.INFO, "Add content type for model %s - %s", modelPage.BlueprintName, modelPage.ModelName)
			db.Model(&core.ContentType{}).Where(
				&core.ContentType{BlueprintName: modelPage.BlueprintName, ModelName: modelPage.ModelName},
			).First(&contentType)
			if contentType.ID == 0 {
				contentType = core.ContentType{BlueprintName: modelPage.BlueprintName, ModelName: modelPage.ModelName}
				db.Create(&contentType)
				core.Trail(core.INFO, "Created content type for blueprint %s model %s", modelPage.BlueprintName, modelPage.ModelName)
			}
			for permDescribed := range core.ProjectPermRegistry.GetAllPermissions() {
				db.Model(&core.Permission{}).Where(
					&core.Permission{ContentTypeID: contentType.ID, PermissionBits: permDescribed.Bit},
				).First(&permission)
				if permission.ID == 0 {
					permission = core.Permission{ContentTypeID: contentType.ID, PermissionBits: permDescribed.Bit}
					db.Create(&permission)
					core.Trail(core.INFO, "Created permission %s for blueprint %s model %s", permDescribed.Name, modelPage.BlueprintName, modelPage.ModelName)
					permission = core.Permission{}
				}
				permission = core.Permission{}
			}
			contentType = core.ContentType{}
		}
	}
	return nil
}

func (command SyncContentTypes) GetHelpText() string {
	return "Sync your content types"
}
