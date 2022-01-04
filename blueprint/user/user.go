package user

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/user/migrations"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Blueprint struct {
	core.Blueprint
}

type PasswordValidationStruct struct {
	Password          string `valid:"password-go-monolith"`
	ConfirmedPassword string
}

type ForgotPasswordHandlerParams struct {
	Email string `form:"email" json:"email" xml:"email"  binding:"required" valid:"email"`
}

type ResetPasswordHandlerParams struct {
	Code              string `form:"code" json:"code" xml:"code"  binding:"required"`
	Password          string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

type ChangePasswordHandlerParams struct {
	OldPassword       string `form:"old_password" json:"old_password" xml:"old_password"  binding:"required"`
	Password          string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	app.GetRouter().GET("/reset-password/", func(ctx *gin.Context) {
		type Context struct {
			core.AdminContext
		}
		c := &Context{}
		core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
		tr := core.NewTemplateRenderer("Reset Password")
		tr.Render(ctx, core.CurrentConfig.GetPathToTemplate("resetpassword"), c, core.FuncMap)
	})
	group.POST("/api/forgot/", func(ctx *gin.Context) {
		var json ForgotPasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err1.Error()))
			return
		}
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		user := core.GenerateUserModel()
		db.Model(core.GenerateUserModel()).Where(&core.User{Email: json.Email}).First(user)
		if user.GetID() == 0 {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("user_not_found", "User with this email not found"))
			return
		}
		templateWriter := bytes.NewBuffer([]byte{})
		template1, err := template.ParseFS(core.CurrentConfig.TemplatesFS, core.CurrentConfig.GetPathToTemplate("email/forgot"))
		if err != nil {
			core.Trail(core.ERROR, "RenderHTML unable to parse %s. %s", core.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		type Context struct {
			Name    string
			Website string
			URL     string
		}

		c := Context{}
		c.Name = user.GetFullName()
		c.Website = core.CurrentConfig.D.GoMonolith.SiteName
		host := core.CurrentConfig.D.GoMonolith.PoweredOnSite
		// @todo, generate code to restore access
		actionExpiresAt := time.Now()
		actionExpiresAt = actionExpiresAt.Add(time.Duration(core.CurrentConfig.D.GoMonolith.ForgotCodeExpiration) * time.Minute)
		var oneTimeAction = core.OneTimeAction{
			User:       *user.(*core.User),
			ExpiresOn:  actionExpiresAt,
			Code:       core.GenerateRandomString(32, &core.OnlyLetersNumbersStringAlphabet),
			ActionType: 1,
		}

		db.Model(core.OneTimeAction{}).Save(&oneTimeAction)
		link := host + core.CurrentConfig.D.GoMonolith.RootAdminURL + "resetpassword/?key=" + oneTimeAction.Code
		c.URL = link
		err = template1.Execute(templateWriter, c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			core.Trail(core.ERROR, "RenderHTML unable to parse %s. %s", core.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			return
		}
		subject := "Password reset for admin panel on the " + core.CurrentConfig.D.GoMonolith.SiteName
		emailSender := core.ProjectEmailSenderFactory.MakeEmailSender()
		emailSender.SetFrom(core.CurrentConfig.D.GoMonolith.EmailFrom)
		emailSender.AddRecipient(user.GetEmail())
		emailSender.SetSubject(subject)
		emailSender.SetBody(templateWriter.String())
		err = emailSender.SendEmail()

		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		}
		return
	})
	group.POST("/api/reset-password/", func(ctx *gin.Context) {
		var json ResetPasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err1.Error()))
			return
		}
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		var oneTimeAction core.OneTimeAction
		db.Model(core.OneTimeAction{}).Where(&core.OneTimeAction{Code: json.Code, IsUsed: false}).Preload("User").First(&oneTimeAction)
		if oneTimeAction.ID == 0 {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("code_not_found", "No such code found"))
			return
		}
		if oneTimeAction.ExpiresOn.Before(time.Now()) {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("code_expired", "Code is expired"))
			return
		}
		passwordValidationStruct := &PasswordValidationStruct{
			Password:          json.Password,
			ConfirmedPassword: json.ConfirmedPassword,
		}
		_, err := govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		hashedPassword, err := utils2.HashPass(json.Password, oneTimeAction.User.Salt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		oneTimeAction.User.Password = hashedPassword
		oneTimeAction.User.IsPasswordUsable = true
		oneTimeAction.IsUsed = true
		db.Save(&oneTimeAction.User)
		db.Save(&oneTimeAction)
	})
	group.POST("/api/change-password/", func(ctx *gin.Context) {
		var json ChangePasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err1.Error()))
			return
		}
		passwordValidationStruct := &PasswordValidationStruct{
			Password:          json.Password,
			ConfirmedPassword: json.ConfirmedPassword,
		}
		_, err := govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		hashedPassword, err := utils2.HashPass(json.OldPassword, user.GetSalt())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
			return
		}
		// @todo, get it back once stabilize pass api
		//if hashedPassword != user.Password {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password doesn't match current one"})
		//	return
		//}
		if user.GetSalt() == "" {
			user.SetSalt(core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength))
		}
		hashedPassword, err = utils2.HashPass(json.Password, user.GetSalt())
		user.SetPassword(hashedPassword)
		user.SetIsPasswordUsable(true)
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		db.Save(reflect.ValueOf(user).Interface())
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/disable-2fa/", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.SetOTPRequired(false)
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/enable-2fa/", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.SetOTPRequired(true)
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	app.GetRouter().NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.RequestURI, "/static-inbuilt/") || strings.HasSuffix(ctx.Request.RequestURI, ".css") ||
			strings.HasSuffix(ctx.Request.RequestURI, ".js") || strings.HasSuffix(ctx.Request.RequestURI, ".map") {
			ctx.Abort()
			return
		}
		// ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		type Context struct {
			core.AdminContext
			Menu string
		}
		c := &Context{}
		core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
		//
		//if r.Form.Get("err_msg") != "" {
		//	c.ErrMsg = r.Form.Get("err_msg")
		//}
		//if code, err := strconv.ParseUint(r.Form.Get("err_code"), 10, 16); err == nil {
		//	c.ErrCode = int(code)
		//}NewAdminPage
		ctx.Status(404)
		tr := core.NewTemplateRenderer("Page not found")
		tr.Render(ctx, core.CurrentConfig.GetPathToTemplate("404"), c, core.FuncMap)
	})
	usersAdminPage := core.NewGormAdminPage(
		nil,
		nil,
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	usersAdminPage.PageName = "Users"
	usersAdminPage.Slug = "users"
	usersAdminPage.BlueprintName = "user"
	usersAdminPage.Router = app.GetRouter()
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(usersAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	var usermodelAdminPage *core.AdminPage
	usermodelAdminPage = core.NewGormAdminPage(
		usersAdminPage,
		&core.User{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			firstField := "Username"
			if core.CurrentConfig.D.GoMonolith.DirectAPISigninByField == "email" {
				firstField = "Email"
			}
			fields := []string{firstField, "FirstName", "LastName", "Email", "Active", "IsStaff", "IsSuperUser", "Password", "Photo", "LastLogin", "ExpiresOn"}
			if ctx.GetUserObject().GetIsSuperUser() {
				fields = append(fields, "UserGroups")
				fields = append(fields, "Permissions")
			}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().GetIsSuperUser() {
				usergroupsField, _ := form.FieldRegistry.GetByName("UserGroups")
				usergroupsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.User)
					vTmp := v.([]string)
					var usergroup *core.UserGroup
					if model.ID != 0 {
						afo.GetDatabase().Db.Model(model).Association("UserGroups").Clear()
						model.UserGroups = make([]core.UserGroup, 0)
					}
					for _, ID := range vTmp {
						if ID == "" {
							continue
						}
						afo.GetDatabase().Db.First(&usergroup, ID)
						if usergroup.ID != 0 {
							model.UserGroups = append(model.UserGroups, *usergroup)
						}
						usergroup = nil
					}
					return nil
				}
				userGroupsWidget := usergroupsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				userGroupsWidget.AddNewLink = fmt.Sprintf("%s/%s/usergroup/edit/%s?_to_field=id&_popup=1", core.CurrentConfig.D.GoMonolith.RootAdminURL, usersAdminPage.Slug, "new")
				userGroupsWidget.AddNewTitle = "Add another group"
				userGroupsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var groups []*core.UserGroup
					database := core.NewDatabaseInstance()
					database.Db.Find(&groups)
					ret := make([]*core.SelectOptGroup, 0)
					for _, group := range groups {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: group.GroupName,
							Value:    group.ID,
						})
					}
					database.Close()
					return ret
				}
				userGroupsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.User)
					if user.ID != 0 {
						var groups []*core.UserGroup
						database := core.NewDatabaseInstance()
						database.Db.Model(user).Association("UserGroups").Find(&groups)
						ret = make([]*core.SelectOptGroup, 0)
						for _, group := range groups {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: group.GroupName,
								Value:    group.ID,
							})
						}
						database.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["UserGroups"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						if len(IDI) > 0 {
							var groups []*core.UserGroup
							database := core.NewDatabaseInstance()
							database.Db.Find(&groups, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, group := range groups {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: group.GroupName,
									Value:    group.ID,
								})
							}
							database.Close()
							return ret
						}
					}
					return ret
				}
				userGroupsWidget.LeftSelectTitle = "Available groups"
				userGroupsWidget.LeftSelectHelp = "available_groups_left"
				userGroupsWidget.LeftSearchSelectHelp = "available_groups_search_help"
				userGroupsWidget.LeftHelpChooseAll = "Click to choose all groups at once."
				userGroupsWidget.RightSelectTitle = "Chosen groups"
				userGroupsWidget.RightSelectHelp = "chosen_groups_left"
				userGroupsWidget.RightSearchSelectHelp = ""
				userGroupsWidget.RightHelpChooseAll = "Click to remove all chosen groups at once."
				userGroupsWidget.HelpText = "group_widget_help"
				permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
				permissionsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.User)
					vTmp := v.([]string)
					var permission *core.Permission
					if model.ID != 0 {
						afo.GetDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]core.Permission, 0)
					}
					for _, ID := range vTmp {
						if ID == "" {
							continue
						}
						afo.GetDatabase().Db.First(&permission, ID)
						if permission.ID != 0 {
							model.Permissions = append(model.Permissions, *permission)
						}
						permission = nil
					}
					return nil
				}
				permissionsWidget := permissionsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var permissions []*core.Permission
					database := core.NewDatabaseInstance()
					database.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*core.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value:    permission.ID,
						})
					}
					database.Close()
					return ret
				}
				permissionsWidget.LeftSelectTitle = "Available user permissions"
				permissionsWidget.LeftSelectHelp = "available_permissions_left"
				permissionsWidget.LeftSearchSelectHelp = "available_permissions_search_help"
				permissionsWidget.LeftHelpChooseAll = "Click to choose all user permissions at once."
				permissionsWidget.RightSelectTitle = "Chosen user permissions"
				permissionsWidget.RightSelectHelp = "chosen_permissions_left"
				permissionsWidget.RightSearchSelectHelp = ""
				permissionsWidget.RightHelpChooseAll = "Click to remove all chosen user permissions at once."
				permissionsWidget.HelpText = "permission_widget_help"
				permissionsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.User)
					if user.ID != 0 {
						var permissions []*core.Permission
						database := core.NewDatabaseInstance()
						database.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*core.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						database.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["Permissions"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						var permissions []*core.Permission
						if len(IDI) > 0 {
							database := core.NewDatabaseInstance()
							database.Db.Preload("ContentType").Find(&permissions, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, permission := range permissions {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: permission.ShortDescription(),
									Value:    permission.ID,
								})
							}
							database.Close()
							return ret
						}
					}
					return ret
				}
			}
			passwordField, _ := form.FieldRegistry.GetByName("Password")
			passwordField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				user := m.(*core.User)
				vI, _ := v.(string)
				if vI != "" {
					if user.Salt == "" {
						user.Salt = core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength)
					}
					hashedPassword, _ := utils2.HashPass(vI, user.Salt)
					user.IsPasswordUsable = true
					user.Password = hashedPassword
				}

				return nil
			}
			return form
		},
	)
	usermodelAdminPage.PageName = "Users"
	usermodelAdminPage.Slug = "user"
	usermodelAdminPage.BlueprintName = "user"
	usermodelAdminPage.Router = app.GetRouter()
	listFilter := &core.ListFilter{
		URLFilteringParam: "IsSuperUser__exact",
		Title:             "Is super user ?",
	}
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "Yes", Value: true})
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "No", Value: false})
	usermodelAdminPage.ListFilter.Add(listFilter)
	err = usersAdminPage.SubPages.AddAdminPage(usermodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	usergroupsAdminPage := core.NewGormAdminPage(
		usersAdminPage,
		&core.UserGroup{},
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"GroupName"}
			if ctx.GetUserObject().GetIsSuperUser() {
				fields = append(fields, "Permissions")
			}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().GetIsSuperUser() {
				permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
				permissionsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.UserGroup)
					vTmp := v.([]string)
					var permission *core.Permission
					if model.ID != 0 {
						afo.GetDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]core.Permission, 0)
					}
					for _, ID := range vTmp {
						if ID == "" {
							continue
						}
						afo.GetDatabase().Db.First(&permission, ID)
						if permission.ID != 0 {
							model.Permissions = append(model.Permissions, *permission)
						}
						permission = nil
					}
					return nil
				}
				permissionsWidget := permissionsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var permissions []*core.Permission
					database := core.NewDatabaseInstance()
					database.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*core.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value:    permission.ID,
						})
					}
					database.Close()
					return ret
				}
				permissionsWidget.LeftSelectTitle = "Available permissions"
				permissionsWidget.LeftSelectHelp = "available_permissions_left"
				permissionsWidget.LeftSearchSelectHelp = "available_permissions_search_help"
				permissionsWidget.LeftHelpChooseAll = "Click to choose all user permissions at once."
				permissionsWidget.RightSelectTitle = "Chosen permissions"
				permissionsWidget.RightSelectHelp = "chosen_permissions_left"
				permissionsWidget.RightSearchSelectHelp = ""
				permissionsWidget.RightHelpChooseAll = "Click to remove all chosen permissions at once."
				permissionsWidget.HelpText = "permission_widget_help"
				permissionsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.UserGroup)
					if user.ID != 0 {
						var permissions []*core.Permission
						database := core.NewDatabaseInstance()
						database.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*core.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						database.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["Permissions"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						var permissions []*core.Permission
						if len(IDI) > 0 {
							database := core.NewDatabaseInstance()
							database.Db.Preload("ContentType").Find(&permissions, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, permission := range permissions {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: permission.ShortDescription(),
									Value:    permission.ID,
								})
							}
							database.Close()
							return ret
						}
					}
					return ret
				}
			}
			return form
		},
	)
	usergroupsAdminPage.PageName = "User groups"
	usergroupsAdminPage.Slug = "usergroup"
	usergroupsAdminPage.BlueprintName = "user"
	usergroupsAdminPage.Router = app.GetRouter()
	err = usersAdminPage.SubPages.AddAdminPage(usergroupsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.OneTimeAction{}, &[]*core.OneTimeAction{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.Permission{}, &[]*core.Permission{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.User{}, &[]*core.User{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.UserGroup{}, &[]*core.UserGroup{} })

	core.GoMonolithValidatorRegistry.AddValidator("username-unique", func(i interface{}, o interface{}) error {
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		var cUsers int64
		db.Model(&core.User{}).Where(&core.User{Username: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return core.NewHTTPErrorResponse("user_name_already_registered", "user with name %s is already registered", i.(string))
	})

	core.GoMonolithValidatorRegistry.AddValidator("email-unique", func(i interface{}, o interface{}) error {
		database := core.NewDatabaseInstance()
		defer database.Close()
		db := database.Db
		var cUsers int64
		db.Model(&core.User{}).Where(&core.User{Email: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return core.NewHTTPErrorResponse("user_email_already_registered", "user with email %s is already registered", i.(string))
	})

	core.GoMonolithValidatorRegistry.AddValidator("username-go-monolith", func(i interface{}, o interface{}) error {
		minLength := core.CurrentConfig.D.Auth.MinUsernameLength
		maxLength := core.CurrentConfig.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return core.NewHTTPErrorResponse("username_length_error", "length of the username has to be between %s and %s symbols", strconv.Itoa(minLength), strconv.Itoa(maxLength))
		}
		return nil
	})

	core.GoMonolithValidatorRegistry.AddValidator("password-go-monolith", func(i interface{}, o interface{}) error {
		passwordStruct := o.(PasswordValidationStruct)
		if passwordStruct.Password != passwordStruct.ConfirmedPassword {
			return core.NewHTTPErrorResponse("password_not_equal", "password doesn't equal to confirmed password")
		}
		if len(passwordStruct.Password) < core.CurrentConfig.D.Auth.MinPasswordLength {
			return core.NewHTTPErrorResponse("password_length_error", "length of the password has to be at least %d symbols", strconv.Itoa(core.CurrentConfig.D.Auth.MinPasswordLength))
		}
		return nil
	})

	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		validator, _ := core.GoMonolithValidatorRegistry.GetValidator("username-unique")
		userExists := validator(i, o)
		return userExists == nil
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		validator, _ := core.GoMonolithValidatorRegistry.GetValidator("email-unique")
		emailExists := validator(i, o)
		return emailExists == nil
	})
	govalidator.CustomTypeTagMap.Set("username-go-monolith", func(i interface{}, o interface{}) bool {
		validator, _ := core.GoMonolithValidatorRegistry.GetValidator("username-go-monolith")
		isValidUsername := validator(i, o)
		return isValidUsername == nil
	})
	govalidator.CustomTypeTagMap.Set("password-go-monolith", func(i interface{}, o interface{}) bool {
		validator, _ := core.GoMonolithValidatorRegistry.GetValidator("password-go-monolith")
		isValidPassword := validator(i, o)
		return isValidPassword == nil
	})
	core.FormCongirurableOptionInstance.AddFieldFormOptions(&core.FieldFormOptions{
		WidgetType: "image",
		Name:       "UserPhotoFormOptions",
		WidgetPopulate: func(w core.IWidget, renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
			fsStorage := core.NewFsStorage()
			photo := renderContext.Model.(core.IUser).GetPhoto()
			if photo == "" {
				return ""
			}
			return fmt.Sprintf("%s%s", fsStorage.GetUploadURL(), photo)
		},
	})
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
