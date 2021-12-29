package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	interfaces3 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/auth/migrations"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm/schema"
	"net/http"
	"reflect"
	"strings"
)

type Blueprint struct {
	core.Blueprint
	AuthAdapterRegistry *core.AuthProviderRegistry
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	b.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.DirectAuthForAdminProvider{})
	// register signin adapters http endpoints
	for adapter := range b.AuthAdapterRegistry.Iterate() {
		adapterGroup := group.Group("/" + adapter.GetName())
		adapterGroup.POST("/signin/", adapter.Signin)
		adapterGroup.POST("/signup/", adapter.Signup)
		adapterGroup.POST("/logout/", adapter.Logout)
		adapterGroup.GET("/status/", adapter.IsAuthenticated)
	}
	// customize root admin page list handler to open home page for administrators, could be customized to show some enhanced dashboards
	core.CurrentDashboardAdminPanel.ListHandler = func(ctx *gin.Context) {
		defaultAdapter, _ := b.AuthAdapterRegistry.GetAdapter("direct-for-admin")
		userSession := defaultAdapter.GetSession(ctx)
		if userSession == nil || userSession.GetUser() == nil {
			type Context struct {
				core.AdminContext
				PrimaryFieldForSignin string
			}
			c := &Context{}
			adminRequestParams := core.NewAdminRequestParams()
			adminRequestParams.NeedAllLanguages = true
			c.PrimaryFieldForSignin = strings.Title(core.CurrentConfig.D.GoMonolith.DirectAPISigninByField)
			core.PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := core.NewTemplateRenderer("Admin Login")
			tr.Render(ctx, core.CurrentConfig.GetPathToTemplate("login"), c, core.FuncMap)
		} else {
			type Context struct {
				core.AdminContext
				Menu        string
				CurrentPath string
			}

			c := &Context{}
			core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
			menu := string(core.CurrentDashboardAdminPanel.AdminPages.PreparePagesForTemplate(c.UserPermissionRegistry))
			c.Menu = menu
			c.CurrentPath = ctx.Request.URL.Path
			tr := core.NewTemplateRenderer("Dashboard")
			tr.Render(ctx, core.CurrentConfig.GetPathToTemplate("home"), c, core.FuncMap)
		}
	}
	// Serve static fs if upload directory is configured. Has to be created in the root of the project.
	if core.CurrentConfig.GetURLToUploadDirectory() != "" {
		app.GetRouter().StaticFS(core.CurrentConfig.GetURLToUploadDirectory(), http.Dir(fmt.Sprintf("./%s", core.CurrentConfig.GetURLToUploadDirectory())))
	}
	// profile page for admin panel
	app.GetRouter().Any(core.CurrentConfig.D.GoMonolith.RootAdminURL+"/profile/", func(ctx *gin.Context) {
		type Context struct {
			core.AdminContext
			ID           uint
			Status       bool
			IsUpdated    bool
			Notif        string
			ProfilePhoto string
			OTPImage     string
			OTPRequired  bool
			ChangesSaved bool
			DBFields     []*schema.Field
			F            *core.Form
			User         string
		}

		c := &Context{}
		core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		form1 := core.NewFormFromModelFromGinContext(c, user, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin", "ExpiresOn", "OTPRequired"}, true, "")
		form1.TemplateName = "form/profile_form"
		c.F = form1
		c.User = user.GetUsername()
		if ctx.Request.Method == "POST" {
			requestForm, _ := ctx.MultipartForm()
			formError := form1.ProceedRequest(requestForm, user, c)
			if formError.IsEmpty() {
				database := core.NewDatabaseInstance()
				defer database.Close()
				db := database.Db
				db.Save(reflect.ValueOf(user).Interface())
				if db.Error != nil {
					formError.GeneralErrors = append(formError.GeneralErrors, db.Error)
				} else {
					c.ChangesSaved = true
					form1.ChangesSaved = true
					ctx.Redirect(302, ctx.Request.URL.String())
				}
				return
			}
		}
		tr := core.NewTemplateRenderer(fmt.Sprintf("%s's Profile", c.User))
		tr.Render(ctx, core.CurrentConfig.GetPathToTemplate("profile"), c, core.FuncMap)
	})
}

func (b Blueprint) InitApp(app core.IApp) {
}

var ConcreteBlueprint = Blueprint{
	Blueprint: core.Blueprint{
		Name:              "auth",
		Description:       "blueprint for auth functionality",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	AuthAdapterRegistry: core.NewAuthProviderRegistry(),
}
