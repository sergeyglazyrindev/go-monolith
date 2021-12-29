package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var CurrentAdminPageRegistry *AdminPageRegistry

type AdminPagesList []*AdminPage

func (apl AdminPagesList) Len() int { return len(apl) }
func (apl AdminPagesList) Less(i, j int) bool {
	if apl[i].Ordering == apl[j].Ordering {
		return apl[i].PageName < apl[j].PageName
	}
	return apl[i].Ordering < apl[j].Ordering
}
func (apl AdminPagesList) Swap(i, j int) { apl[i], apl[j] = apl[j], apl[i] }

type AdminPageRegistry struct {
	AdminPages map[string]*AdminPage
}

func (apr *AdminPageRegistry) GetByModelName(modelName string) *AdminPage {
	for adminPage := range apr.GetAll() {
		for subPage := range adminPage.SubPages.GetAll() {
			projectModel := ProjectModels.GetModelFromInterface(subPage.Model)
			if projectModel == nil {
				continue
			}
			// spew.Dump("AAAAAAAAAA", projectModel.Statement.Schema.Name)
			if projectModel.Statement.Schema.Name == modelName {
				return subPage
			}
			for subPage1 := range subPage.SubPages.GetAll() {
				projectModel1 := ProjectModels.GetModelFromInterface(subPage1.Model)
				// spew.Dump("AAAAAAAAAA1111111111", projectModel1.Statement.Schema.Name)
				if projectModel1.Statement.Schema.Name == modelName {
					return subPage1
				}
			}
		}
		if adminPage.Model == nil {
			continue
		}
		projectModel := ProjectModels.GetModelFromInterface(adminPage.Model)
		// spew.Dump("AAAAAAAAAA222", projectModel.Statement.Schema.Name)
		if projectModel.Statement.Schema.Name == modelName {
			return adminPage
		}
	}
	return nil
}

func (apr *AdminPageRegistry) GetBySlug(slug string) (*AdminPage, error) {
	adminPage, ok := apr.AdminPages[slug]
	if !ok {
		return nil, fmt.Errorf("no admin page with alias %s", slug)
	}
	return adminPage, nil
}

func (apr *AdminPageRegistry) AddAdminPage(adminPage *AdminPage) error {
	apr.AdminPages[adminPage.Slug] = adminPage
	return nil
}

func (apr *AdminPageRegistry) GetAll() <-chan *AdminPage {
	chnl := make(chan *AdminPage)
	go func() {
		defer close(chnl)
		sortedPages := make(AdminPagesList, 0)

		for _, adminPage := range apr.AdminPages {
			sortedPages = append(sortedPages, adminPage)
		}
		sort.Slice(sortedPages, func(i int, j int) bool {
			if sortedPages[i].Ordering == sortedPages[j].Ordering {
				return sortedPages[i].PageName < sortedPages[j].PageName
			}
			return sortedPages[i].Ordering < sortedPages[j].Ordering

		})
		for _, adminPage := range sortedPages {
			chnl <- adminPage
		}

	}()
	return chnl
}

func (apr *AdminPageRegistry) PreparePagesForTemplate(permRegistry *UserPermRegistry) []byte {
	pages := make([]*AdminPage, 0)

	for page := range apr.GetAll() {
		blueprintName := page.BlueprintName
		modelName := page.ModelName
		var userPerm *UserPerm
		if modelName != "" {
			userPerm = permRegistry.GetPermissionForBlueprint(blueprintName, modelName)
			if !userPerm.HasReadPermission() {
				continue
			}
		} else {
			existsAnyPermission := permRegistry.IsThereAnyPermissionForBlueprint(blueprintName)
			if !existsAnyPermission {
				continue
			}
		}
		pages = append(pages, page)
	}
	ret, err := json.Marshal(pages)
	if err != nil {
		Trail(CRITICAL, "error while generating menu in admin", err)
	}
	return ret
}

type AdminPage struct {
	Model                              interface{}                                                              `json:"-"`
	GenerateForm                       func(modelI interface{}, ctx IAdminContext) *Form                        `json:"-"`
	GetQueryset                        func(IAdminContext, *AdminPage, *AdminRequestParams) IAdminFilterObjects `json:"-"`
	ModelActionsRegistry               *AdminModelActionRegistry                                                `json:"-"`
	FilterOptions                      *FilterOptionsRegistry                                                   `json:"-"`
	ActionsSelectionCounter            bool                                                                     `json:"-"`
	BlueprintName                      string
	EmptyValueDisplay                  string                   `json:"-"`
	ExcludeFields                      IFieldRegistry           `json:"-"`
	FieldsToShow                       IFieldRegistry           `json:"-"`
	Form                               *Form                    `json:"-"`
	ShowAllFields                      bool                     `json:"-"`
	Validators                         *ValidatorRegistry       `json:"-"`
	InlineRegistry                     *AdminPageInlineRegistry `json:"-"`
	ListDisplay                        *ListDisplayRegistry     `json:"-"`
	ListFilter                         *ListFilterRegistry      `json:"-"`
	MaxShowAll                         int                      `json:"-"`
	PreserveFilters                    bool                     `json:"-"`
	SaveAndContinue                    bool                     `json:"-"`
	SaveOnTop                          bool                     `json:"-"`
	SearchFields                       *SearchFieldRegistry     `json:"-"`
	ShowFullResultCount                bool                     `json:"-"`
	ViewOnSite                         bool                     `json:"-"`
	ListTemplate                       string                   `json:"-"`
	AddTemplate                        string                   `json:"-"`
	EditTemplate                       string                   `json:"-"`
	DeleteConfirmationTemplate         string                   `json:"-"`
	DeleteSelectedConfirmationTemplate string                   `json:"-"`
	ObjectHistoryTemplate              string                   `json:"-"`
	PopupResponseTemplate              string                   `json:"-"`
	Paginator                          *Paginator               `json:"-"`
	SubPages                           *AdminPageRegistry       `json:"-"`
	Ordering                           int
	PageName                           string
	ModelName                          string
	Slug                               string
	ToolTip                            string
	Icon                               string
	ListHandler                        func(ctx *gin.Context)                                                 `json:"-"`
	EditHandler                        func(ctx *gin.Context)                                                 `json:"-"`
	AddHandler                         func(ctx *gin.Context)                                                 `json:"-"`
	DeleteHandler                      func(ctx *gin.Context)                                                 `json:"-"`
	Router                             *gin.Engine                                                            `json:"-"`
	ParentPage                         *AdminPage                                                             `json:"-"`
	SaveModel                          func(modelI interface{}, ID uint, afo IAdminFilterObjects) interface{} `json:"-"`
	RegisteredHTTPHandlers             bool
	NoPermissionToAddNew               bool
	NoPermissionToEdit                 bool
	PermissionName                     CustomPermission
	EnhanceQuerySet                    func(afo IAdminFilterObjects)                                                                `json:"-"`
	CustomizeQuerySet                  func(adminContext IAdminContext, afo IAdminFilterObjects, requestParams *AdminRequestParams) `json:"-"`
}

type ModelActionRequestParams struct {
	ObjectIds     string `form:"object_ids" json:"object_ids" xml:"object_ids"  binding:"required"`
	RealObjectIds []string
}

func (ap *AdminPage) GenerateLinkForModelAutocompletion() string {
	return fmt.Sprintf("%s/%s/%s/autocomplete/", CurrentConfig.D.GoMonolith.RootAdminURL, ap.ParentPage.Slug, ap.Slug)
}

func (ap *AdminPage) GenerateLinkToEditModel(gormModelV reflect.Value) string {
	ID := GetID(gormModelV)
	return fmt.Sprintf("%s/%s/%s/edit/%d/", CurrentConfig.D.GoMonolith.RootAdminURL, ap.ParentPage.Slug, ap.Slug, ID)
}

func (ap *AdminPage) DoesUserHavePermission(u IUser, permissionNameL ...CustomPermission) bool {
	permissionName := ap.PermissionName
	if permissionName == "" && len(permissionNameL) > 0 {
		permissionName = permissionNameL[0]
	}
	if permissionName == "" {
		return false
	}
	userPermissions := u.BuildPermissionRegistry()
	// modelI, _ := ap.GenerateModelI()
	userPerm := userPermissions.GetPermissionForBlueprint(ap.BlueprintName, ap.ModelName)
	return userPerm.DoesUserHaveRightFor(permissionName)
}

func (ap *AdminPage) GenerateLinkToAddNewModel(inPopup bool) string {
	url := fmt.Sprintf("%s/%s/%s/edit/new/", CurrentConfig.D.GoMonolith.RootAdminURL, ap.ParentPage.Slug, ap.Slug)
	if inPopup {
		url += "?_to_field=id&_popup=1"
	}
	return url
}

func (ap *AdminPage) HandleModelAction(modelActionName string, ctx *gin.Context) {
	adminContext := &AdminContext{}
	adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
	PopulateTemplateContextForAdminPanel(ctx, adminContext, adminRequestParams)
	afo := ap.GetQueryset(adminContext, ap, adminRequestParams)
	var json1 ModelActionRequestParams
	if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
		if err := ctx.ShouldBindJSON(&json1); err != nil {
			ctx.JSON(http.StatusBadRequest, APIBadResponse(err.Error()))
			return
		}
	} else {
		if err := ctx.ShouldBind(&json1); err != nil {
			ctx.JSON(http.StatusBadRequest, APIBadResponse(err.Error()))
			return
		}
	}
	objectIds := strings.Split(json1.ObjectIds, ",")
	objectUintIds := make([]string, 0)
	for _, objectID := range objectIds {
		objectUintIds = append(objectUintIds, objectID)
	}
	json1.RealObjectIds = objectUintIds
	if len(json1.RealObjectIds) > 0 {
		primaryKeyField, _ := ap.Form.FieldRegistry.GetPrimaryKey()
		afo.FilterByMultipleIds(primaryKeyField, json1.RealObjectIds)
		modelAction, _ := ap.ModelActionsRegistry.GetModelActionByName(modelActionName)
		if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
			_, affectedRows := modelAction.Handler(ap, afo, ctx)
			ctx.JSON(http.StatusOK, gin.H{"Affected": strconv.Itoa(int(affectedRows))})
		} else {
			modelAction.Handler(ap, afo, ctx)
		}
	} else {
		ctx.Status(400)
	}
}

func (ap *AdminPage) FetchFilterOptions(ctx *gin.Context) []*DisplayFilterOption {
	adminContext := &AdminContext{}
	adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
	PopulateTemplateContextForAdminPanel(ctx, adminContext, adminRequestParams)
	afo := ap.GetQueryset(adminContext, ap, adminRequestParams)
	filterOptions := make([]*DisplayFilterOption, 0)
	for filterOption := range ap.FilterOptions.GetAll() {
		filterOptions = append(filterOptions, filterOption.FetchOptions(afo)...)
	}
	return filterOptions
}

func GetURLToBackAfterSignin(ctx *gin.Context) string {
	url1 := CurrentConfig.D.GoMonolith.RootAdminURL
	urlNew := url.URL{}
	urlNew.Path = url1
	urlNew.Host = ctx.Request.URL.Host
	urlNew.Scheme = ctx.Request.URL.Scheme
	qs := urlNew.Query()
	qs.Set("backto", ctx.Request.URL.String())
	urlNew.RawQuery = qs.Encode()
	return urlNew.String()
}

func (ap *AdminPage) GetURLToBackAfterSignin(ctx *gin.Context) string {
	return GetURLToBackAfterSignin(ctx)
}

func NewAdminPageRegistry() *AdminPageRegistry {
	return &AdminPageRegistry{
		AdminPages: make(map[string]*AdminPage),
	}
}

type AdminPageInlineRegistry struct {
	Inlines []*AdminPageInline
}

func (apir *AdminPageInlineRegistry) Add(pageInline *AdminPageInline) {
	if pageInline.Prefix == "" {
		pageInline.Prefix = ASCIIRegex.ReplaceAllLiteralString(pageInline.VerboseName, "")
	}
	apir.Inlines = append(apir.Inlines, pageInline)
}

func (apir *AdminPageInlineRegistry) GetAll() <-chan *AdminPageInline {
	chnl := make(chan *AdminPageInline)
	go func() {
		defer close(chnl)
		for _, inline := range apir.Inlines {
			chnl <- inline
		}
	}()
	return chnl
}
