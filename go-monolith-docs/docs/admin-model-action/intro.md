---
sidebar_position: 1
---

# Admin model actions

Admin model actions are to help you to manipulate entities in your database. For example, right now we have one predefined model action: Remove it from database.  
It builds removal tree with all dependencies. And let you remove entity from database safely.  
An example of how it could be used:
```go
adminModelAction := core.NewAdminModelAction(
	"TurnSuperusersToNormalUsers", &core.AdminActionPlacement{},
)
adminModelAction.Handler = func(ap *core.AdminPage, afo core.IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
	tx := afo.GetFullQuerySet().Update("IsSuperUser", false).Commit()
	return tx.(*core.GormPersistenceStorage).Db.Error == nil, tx.(*core.GormPersistenceStorage).Db.RowsAffected
}
adminUserPage.ModelActionsRegistry.AddModelAction(adminModelAction)
```
You have to add your action to GlobalActionRegistry, if you want to make your action globally available, for every model in your admin panel, for example:
```go
removalModelAction := NewAdminModelAction(
	"Delete permanently", &AdminActionPlacement{
		ShowOnTheListPage: true,
	},
)
removalModelAction.RequiresExtraSteps = true
removalModelAction.Description = "Delete users permanently"
removalModelAction.Handler = func(ap *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
	removalPlan := make([]RemovalTreeList, 0)
	removalConfirmed := ctx.PostForm("removal_confirmed")
	afo.GetDatabase().Db.Transaction(func(tx *gorm.DB) error {
		database := &Database{Db: tx, Adapter: afo.GetDatabase().Adapter}
		for modelIterated := range afo.IterateThroughWholeQuerySet() {
			removalTreeNode := BuildRemovalTree(database, modelIterated.Model)
			if removalConfirmed == "" {
				deletionStringified := removalTreeNode.BuildDeletionTreeStringified(database)
				removalPlan = append(removalPlan, deletionStringified)
			} else {
				err := removalTreeNode.RemoveFromDatabase(database)
				if err != nil {
					return err
				}
			}
		}
		if removalConfirmed != "" {
			truncateLastPartOfPath := regexp.MustCompile("/[^/]+/?$")
			newPath := truncateLastPartOfPath.ReplaceAll([]byte(ctx.Request.URL.RawPath), []byte(""))
			clonedURL := CloneNetURL(ctx.Request.URL)
			clonedURL.RawPath = string(newPath)
			clonedURL.Path = string(newPath)
			query := clonedURL.Query()
			query.Set("message", "Objects were removed succesfully")
			clonedURL.RawQuery = query.Encode()
			ctx.Redirect(http.StatusFound, clonedURL.String())
			return nil
		}
		type Context struct {
			AdminContext
			RemovalPlan []RemovalTreeList
			AdminPage   *AdminPage
			ObjectIds   string
		}
		c := &Context{}
		adminRequestParams := NewAdminRequestParams()
		c.RemovalPlan = removalPlan
		c.AdminPage = ap
		c.ObjectIds = ctx.PostForm("object_ids")
		PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

		tr := NewTemplateRenderer(fmt.Sprintf("Remove %s ?", ap.ModelName))
		tr.Render(ctx, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("remove_objects"), c, FuncMap)
		return nil
	})
	return true, 1
}
core.GlobalModelActionRegistry.AddModelAction(removalModelAction)
```
Later on we will migrate it to interface as well, so it could be used easily for any type of list filter functionality.
