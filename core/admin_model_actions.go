package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

var GlobalModelActionRegistry *AdminModelActionRegistry

type RemovalTreeList []*RemovalTreeNodeStringified

func init() {
	GlobalModelActionRegistry = NewAdminModelActionRegistry()
	removalModelAction := NewAdminModelAction(
		"Delete permanently", &AdminActionPlacement{
			ShowOnTheListPage: true,
		},
	)
	removalModelAction.RequiresExtraSteps = true
	removalModelAction.PermName = "delete"
	removalModelAction.Description = "Delete permanently"
	removalModelAction.Handler = func(ap *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
		type Context struct {
			AdminContext
		}
		c := &Context{}
		adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
		PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)
		user := c.GetUserObject()
		if !ap.DoesUserHavePermission(user, "delete") {
			ctx.AbortWithStatus(409)
			return false, 0
		}
		removalPlan := make([]RemovalTreeList, 0)
		removalConfirmed := ctx.PostForm("removal_confirmed")
		removalError := afo.WithTransaction(func(afo1 IAdminFilterObjects) error {
			for modelIterated := range afo.IterateThroughWholeQuerySet() {
				removalTreeNode := BuildRemovalTree(afo1.GetDatabase(), modelIterated.Model)
				if removalConfirmed == "" {
					deletionStringified := removalTreeNode.BuildDeletionTreeStringified(afo1.GetDatabase())
					removalPlan = append(removalPlan, deletionStringified)
				} else {
					err := removalTreeNode.RemoveFromDatabase(afo1.GetDatabase())
					if err != nil {
						return err
					}
				}
			}
			if removalConfirmed != "" {
				query := ctx.Request.URL.Query()
				query.Set("message", Tf(c.Language.Code, "Objects were removed succesfully"))
				ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/?%s", CurrentConfig.D.GoMonolith.RootAdminURL, ap.ParentPage.Slug, ap.ModelName, query.Encode()))
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
			tr.Render(ctx, CurrentConfig.GetPathToTemplate("remove_objects"), c, FuncMap)
			return nil
		})
		if removalError != nil {
			return false, 0
		}
		return true, int64(len(ctx.PostForm("object_ids")))
	}
	GlobalModelActionRegistry.AddModelAction(removalModelAction)
}

type AdminActionPlacement struct {
	DisplayOnEditPage bool
	//DisplayToTheTop bool
	//DisplayToTheBottom bool
	//DisplayToTheRight bool
	//DisplayToTheLeft bool
	ShowOnTheListPage bool
}

type IAdminModelActionInterface interface {
}

type AdminModelAction struct {
	ActionName              string
	Description             string
	ShowFutureChanges       bool
	RedirectToRootModelPage bool
	Placement               *AdminActionPlacement
	PermName                CustomPermission
	Handler                 func(adminPage *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64)
	IsDisabled              func(afo IAdminFilterObjects, ctx *gin.Context) bool
	SlugifiedActionName     string
	RequestMethod           string
	RequiresExtraSteps      bool
}

func prepareAdminModelActionName(adminModelAction string) string {
	slugifiedAdminModelAction := ASCIIRegex.ReplaceAllLiteralString(adminModelAction, "")
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), " ", "_", -1)
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), ".", "_", -1)
	return slugifiedAdminModelAction
}

func NewAdminModelAction(actionName string, placement *AdminActionPlacement) *AdminModelAction {
	return &AdminModelAction{
		RedirectToRootModelPage: true,
		ActionName:              actionName,
		Placement:               placement,
		SlugifiedActionName:     prepareAdminModelActionName(actionName),
		RequestMethod:           "POST",
	}
}

type AdminModelActionRegistry struct {
	AdminModelActions map[string]*AdminModelAction
}

func (amar *AdminModelActionRegistry) AddModelAction(ma *AdminModelAction) {
	amar.AdminModelActions[ma.SlugifiedActionName] = ma
}

func (amar *AdminModelActionRegistry) IsThereAnyActions() bool {
	return len(amar.AdminModelActions) > 0
}

func (amar *AdminModelActionRegistry) IsThereAnyActionsToShowOnEditPage(user IUser, ap *AdminPage) bool {
	for adminModelAction := range amar.GetAllModelActionsForUser(user, ap) {
		if adminModelAction.Placement.DisplayOnEditPage {
			return true
		}
	}
	return false
}

func (amar *AdminModelActionRegistry) GetAllModelActions() <-chan *AdminModelAction {
	chnl := make(chan *AdminModelAction)
	go func() {
		defer close(chnl)
		mActions := make([]*AdminModelAction, 0)
		for _, mAction := range amar.AdminModelActions {
			mActions = append(mActions, mAction)
		}
		sort.Slice(mActions, func(i, j int) bool {
			return mActions[i].ActionName < mActions[j].ActionName
		})
		for _, mAction := range mActions {
			chnl <- mAction
		}
	}()
	return chnl
}

func (amar *AdminModelActionRegistry) GetAllModelActionsForUser(user IUser, adminPage *AdminPage) <-chan *AdminModelAction {
	chnl := make(chan *AdminModelAction)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Trail(CRITICAL, "Recovering from panic in GetAllModelActionsForUser error is: %v \n", r)
			}
			close(chnl)
		}()
		permissionRegistry := user.BuildPermissionRegistry()
		mActions := make([]*AdminModelAction, 0)
		userPerm := permissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
		for _, mAction := range amar.AdminModelActions {
			if mAction.PermName != "" && !userPerm.DoesUserHaveRightFor(mAction.PermName) {
				continue
			}
			mActions = append(mActions, mAction)
		}
		sort.Slice(mActions, func(i, j int) bool {
			return mActions[i].ActionName < mActions[j].ActionName
		})
		for _, mAction := range mActions {
			chnl <- mAction
		}
	}()
	return chnl
}

func (amar *AdminModelActionRegistry) GetModelActionByName(actionName string) (*AdminModelAction, error) {
	mAction, exists := amar.AdminModelActions[actionName]
	if !exists {
		return nil, fmt.Errorf("found no model action with name %s", actionName)
	}
	return mAction, nil
}

type RemovalTreeNode struct {
	Model            interface{}
	ModelDescription *ModelDescription
	Next             []*RemovalTreeNode
	Prev             []*RemovalTreeNode
	RawSQL           []*DeleteRowStructure
	Visited          bool
	Level            int
}

type RemovalTreeNodeStringified struct {
	Explanation template.HTML
	Level       int
}

type RemovalOrderList []*RemovalTreeNode

func TraverseRemovalTreeNode(nodeToVisit *RemovalTreeNode, removalOrderList *RemovalOrderList) {
	*removalOrderList = append(*removalOrderList, nodeToVisit)
	nodeToVisit.Visited = false
	for _, nextToRemove := range nodeToVisit.Next {
		nextToRemove.Visited = false
		*removalOrderList = append(*removalOrderList, nextToRemove)
		TraverseRemovalTreeNode(nextToRemove, removalOrderList)
	}
}

func (rtn *RemovalTreeNode) RemoveFromDatabase(database *ProjectDatabase) error {
	var removalOrder RemovalOrderList
	TraverseRemovalTreeNode(rtn, &removalOrder)
	sort.Slice(removalOrder, func(i, j int) bool {
		return removalOrder[i].Level > removalOrder[j].Level
	})
	for _, removalTreeNode := range removalOrder {
		if len(removalTreeNode.RawSQL) > 0 {
			for _, rawSQL := range removalTreeNode.RawSQL {
				res := database.Db.Exec(rawSQL.SQL, rawSQL.Values...)
				if res.Error != nil {
					return res.Error
				}
			}
		}
		res1 := database.Db.Unscoped().Delete(removalTreeNode.Model)
		if res1.Error != nil {
			return res1.Error
		}
	}
	return nil
}

func (rtn *RemovalTreeNode) BuildDeletionTreeStringified(database *ProjectDatabase) []*RemovalTreeNodeStringified {
	var removalTreeStringified []*RemovalTreeNodeStringified
	var removalOrder RemovalOrderList
	TraverseRemovalTreeNode(rtn, &removalOrder)
	//sort.Slice(removalOrder, func(i, j int) bool {
	//	return removalOrder[i].Level > removalOrder[j].Level
	//})
	removalTreeStringified = make([]*RemovalTreeNodeStringified, 0)
	for _, removalTreeNode := range removalOrder {
		database.Db.Unscoped().First(removalTreeNode.Model)
		if len(removalTreeNode.RawSQL) > 0 {
			for _, rawSQL := range removalTreeNode.RawSQL {
				removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
					Explanation: template.HTML(fmt.Sprintf("Association with %s", rawSQL.Table)),
					Level:       removalTreeNode.Level,
				})
			}
		}
		gormModelV := reflect.Indirect(reflect.ValueOf(removalTreeNode.Model))
		Idv := TransformValueForWidget(gormModelV.FieldByName(removalTreeNode.ModelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
		modelAdminPage := CurrentAdminPageRegistry.GetByModelName(removalTreeNode.ModelDescription.Statement.Schema.Name)
		if modelAdminPage != nil {
			url := fmt.Sprintf("%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, modelAdminPage.ParentPage.Slug, modelAdminPage.Slug)
			removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
				Explanation: template.HTML(fmt.Sprintf("<a target='_blank' href='%s%s/'>%s</a>", url, Idv, reflect.ValueOf(removalTreeNode.Model).MethodByName("String").Call([]reflect.Value{})[0])),
				Level:       removalTreeNode.Level,
			})
		} else {
			removalTreeStringified = append(removalTreeStringified, &RemovalTreeNodeStringified{
				Explanation: template.HTML(fmt.Sprintf("%s", reflect.ValueOf(removalTreeNode.Model).MethodByName("String").Call([]reflect.Value{})[0])),
				Level:       removalTreeNode.Level,
			})
		}
	}
	return removalTreeStringified
}

func BuildRemovalTree(database *ProjectDatabase, model interface{}, level ...int) *RemovalTreeNode {
	var realLevel int
	if len(level) == 0 {
		realLevel = 1
	} else {
		realLevel = level[0] + 1
	}
	modelInfo := ProjectModels.GetModelFromInterface(model)
	removalTreeNode := &RemovalTreeNode{
		Model:            model,
		Next:             make([]*RemovalTreeNode, 0),
		Prev:             make([]*RemovalTreeNode, 0),
		RawSQL:           make([]*DeleteRowStructure, 0),
		ModelDescription: modelInfo,
		Level:            realLevel,
	}
	for modelDescription := range ProjectModels.Iterate() {
		for _, relationShip := range modelDescription.Statement.Schema.Relationships.Relations {
			if relationShip.Type == "many_to_many" {
				database.Db.Model(model).Preload(relationShip.Name)
				for _, reference := range relationShip.References {
					if reference.PrimaryKey.Schema.Table == modelInfo.Statement.Table {
						gormModelV := reflect.Indirect(reflect.ValueOf(model))
						cond := fmt.Sprintf(
							"%s = ?",
							reference.ForeignKey.DBName,
						)
						deleteSQL := database.Adapter.BuildDeleteString(
							reference.ForeignKey.Schema.Table,
							cond,
							TransformValueForWidget(gormModelV.FieldByName(modelInfo.Statement.Schema.PrimaryFields[0].Name).Interface()),
						)
						deleteSQL.Table = reference.ForeignKey.Schema.Table
						removalTreeNode.RawSQL = append(removalTreeNode.RawSQL, deleteSQL)
					}
				}
			}
			if relationShip.Type == "belongs_to" {
				relationsString := []string{}
				foundRelation := false
				primaryKeyName := ""
				primaryStructField := ""
				for _, reference := range relationShip.References {
					if reference.PrimaryKey.Schema.Table == modelInfo.Statement.Table {
						foundRelation = true
						primaryKeyName = reference.PrimaryKey.DBName
						primaryStructField = reference.ForeignKey.Name
					}
					relationsString = append(
						relationsString,
						fmt.Sprintf(
							"%s.%s = %s.%s",
							modelDescription.Statement.Table, reference.ForeignKey.DBName, modelInfo.Statement.Table,
							reference.PrimaryKey.DBName,
						),
					)
				}
				if !foundRelation {
					continue
				}
				modelI, _ := modelDescription.GenerateModelI()
				db := database.Db.Model(modelI)
				if relationShip.Field.NotNull {
					db = db.Joins(
						fmt.Sprintf(
							"INNER JOIN %s on %s",
							modelInfo.Statement.Table, strings.Join(relationsString, " AND "),
						),
					)
				} else {
					db = db.Joins(
						fmt.Sprintf(
							"LEFT JOIN %s on %s",
							modelInfo.Statement.Table, strings.Join(relationsString, " AND "),
						),
					)
				}
				gormModelV := reflect.Indirect(reflect.ValueOf(model))
				Idv := TransformValueForWidget(gormModelV.FieldByName(modelInfo.Statement.Schema.PrimaryFields[0].Name).Interface())

				rows, _ := db.Unscoped().Preload(primaryStructField).Where(fmt.Sprintf("%s.%s = ?", modelInfo.Statement.Table, primaryKeyName), Idv).Rows()
				for rows.Next() {
					newModel1, _ := modelDescription.GenerateModelI()
					database.Db.ScanRows(rows, newModel1)
					removalTreeNode.Next = append(removalTreeNode.Next, BuildRemovalTree(database, newModel1, realLevel))
				}
				rows.Close()
			}
		}
	}
	return removalTreeNode
}

func NewAdminModelActionRegistry() *AdminModelActionRegistry {
	adminModelActions := make(map[string]*AdminModelAction)
	ret := &AdminModelActionRegistry{AdminModelActions: adminModelActions}
	if GlobalModelActionRegistry != nil {
		for adminModelAction := range GlobalModelActionRegistry.GetAllModelActions() {
			ret.AddModelAction(adminModelAction)
		}
	}
	return ret
}

func NewEmptyModelActionRegistry() *AdminModelActionRegistry {
	adminModelActions := make(map[string]*AdminModelAction)
	ret := &AdminModelActionRegistry{AdminModelActions: adminModelActions}
	return ret
}
