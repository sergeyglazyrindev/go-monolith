---
sidebar_position: 1
---

# Admin page

Admin page is a structure with following fields:
```go
type AdminPage struct {
	Model                              interface{}                                               `json:"-"`
	GenerateModelI                     func() (interface{}, interface{})                         `json:"-"`
	GenerateForm                       func(modelI interface{}, ctx IAdminContext) *Form         `json:"-"`
	GetQueryset                        func(*AdminPage, *AdminRequestParams) IAdminFilterObjects `json:"-"`
	ModelActionsRegistry               *AdminModelActionRegistry                                 `json:"-"`
	FilterOptions                      *FilterOptionsRegistry                                    `json:"-"`
	ActionsSelectionCounter            bool                                                      `json:"-"`
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
}
```
The easiest way to create new admin page for Gorm is:
```go
abtestmodelAdminPage := core.NewGormAdminPage(
	// parent admin page
	abTestAdminPage,
	// func to provide interfaces for single entity and for multiple entities, used to load/save data
	func() (interface{}, interface{}) { return &abtestmodel.ABTest{}, &[]*abtestmodel.ABTest{} },
	func(modelI interface{}, ctx core.IAdminContext) *core.Form {
		fields := []string{"ContentType", "Type", "Name", "Field", "PrimaryKey", "Active", "Group", "StaticPath"}
		form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
		// here you can customize form for your model.
		return form
	},
)
err = abTestAdminPage.SubPages.AddAdminPage(abtestmodelAdminPage)
```
You may customize everything, for example, if you need customized http handlers, just specify for your admin page any of these:
```go
ListHandler, EditHandler, AddHandler, DeleteHandler.
```
if you need customized SaveModel function, you may do that. For example, like here:
```go
usermodelAdminPage.SaveModel = func(modelI interface{}, ID uint, afo core.IAdminFilterObjects) interface{} {
	user := modelI.(*core.User)
	if user.Salt == "" && user.Password != "" {
		user.Salt = utils.RandStringRunes(core.CurrentConfig.D.Auth.SaltLength)
	}
	if ID != 0 {
		userM := &core.User{}
		afo.GetDatabase().Db.First(userM, ID)
		if userM.Password != user.Password && user.Password != "" {
			// hashedPassword, err := utils2.HashPass(password, salt)
			hashedPassword, _ := utils2.HashPass(user.Password, user.Salt)
			user.IsPasswordUsable = true
			user.Password = hashedPassword
		}
	} else {
		if user.Password != "" {
			// hashedPassword, err := utils2.HashPass(password, salt)
			hashedPassword, _ := utils2.HashPass(user.Password, user.Salt)
			user.Password = hashedPassword
			user.IsPasswordUsable = true
		}
	}
	afo.GetDatabase().Db.Save(user)
	return user
}
```
