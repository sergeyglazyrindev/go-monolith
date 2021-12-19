package sessions

import (
	"github.com/gin-gonic/gin"
	interfaces2 "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/sessions/migrations"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"strings"
)

type Blueprint struct {
	core.Blueprint
	SessionAdapterRegistry *interfaces2.SessionProviderRegistry
}

func (b Blueprint) InitRouter(app core.IApp, group *gin.RouterGroup) {
	// function to verify CSRF
	app.GetRouter().Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if !core.CurrentConfig.RequiresCsrfCheck(c) {
				c.Next()
				return
			}
			if c.Request.Method != "POST" {
				c.Next()
				return
			}
			contentType := c.Request.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				c.Next()
				return
			}
			var serverKey string
			var csrfTokenFromRequest string
			csrfTokenFromRequest = c.Request.Header.Get("CSRF-TOKEN")
			if csrfTokenFromRequest == "" {
				csrfTokenFromRequest, _ = c.Cookie("csrf_token")
				if csrfTokenFromRequest == "" {
					csrfTokenFromRequest = c.PostForm("csrf-token")
				}
			}
			serverKey = c.Request.Header.Get("X-" + strings.ToUpper(core.CurrentConfig.D.GoMonolith.APICookieName))
			if serverKey == "" {
				if c.Query("for-gomonolith-panel") == "1" || strings.Contains(c.Request.URL.String(), core.CurrentConfig.D.GoMonolith.RootAdminURL) {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.GoMonolith.AdminCookieName)
				} else {
					serverKey, _ = c.Cookie(core.CurrentConfig.D.GoMonolith.APICookieName)
				}
			}
			if serverKey == "" {
				c.String(400, "No user session found")
				c.Abort()
				return
			}
			defaultSessionAdapter, _ := b.SessionAdapterRegistry.GetDefaultAdapter()
			session, _ := defaultSessionAdapter.GetByKey(serverKey)
			if session == nil {
				c.String(400, "No user session found")
				c.Abort()
				return
			}
			// @todo, comment it out when stabilize token
			//csrfToken, err := session.Get("csrf_token")
			//if err != nil {
			//	c.String(400, err.Error())
			//	c.Abort()
			//	return
			//}

			if len(csrfTokenFromRequest) != 64 {
				c.String(400, "Incorrect length of csrf-token")
				c.Abort()
				return
			}
			// @todo, comment it out when stabilize token
			//tokenUnmasked := utils.UnmaskCSRFToken(csrfTokenFromRequest)
			//if tokenUnmasked != csrfToken {
			//	c.String(400, "Incorrect csrf-token")
			//	c.Abort()
			//	return
			//}
			c.Next()
		}
	}())
	app.GetRouter().Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.URL.Path, core.CurrentConfig.D.GoMonolith.RootAdminURL) {
				c.Next()
				return
			}
			contentType := c.Request.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				c.Next()
				return
			}
			serverKey := c.Request.Header.Get("X-" + strings.ToUpper(core.CurrentConfig.D.GoMonolith.APICookieName))
			if serverKey == "" {
				serverKey, _ = c.Cookie(core.CurrentConfig.D.GoMonolith.AdminCookieName)
			}
			if serverKey == "" {
				c.Next()
				return
			}
			defaultSessionAdapter, _ := b.SessionAdapterRegistry.GetDefaultAdapter()
			session, _ := defaultSessionAdapter.GetByKey(serverKey)
			if session == nil {
				c.Next()
				return
			}
			if session.IsExpired() && c.Request.URL.Path != core.CurrentConfig.D.GoMonolith.RootAdminURL+"/" {
				if !strings.Contains(c.Request.URL.Path, "resetpassword") {
					c.Redirect(302, core.CurrentConfig.D.GoMonolith.RootAdminURL)
					c.Abort()
					return
				}
			}
			user := session.GetUser()
			if c.Request.URL.Path != core.CurrentConfig.D.GoMonolith.RootAdminURL+"/" && (user == nil || (!user.GetIsStaff() && !user.GetIsSuperUser())) {
				if !strings.Contains(c.Request.URL.Path, "resetpassword") {
					c.Redirect(302, core.CurrentConfig.D.GoMonolith.RootAdminURL)
					c.Abort()
					return
				}
			}
			c.Next()
		}
	}())
	core.FuncMap["CSRF"] = func(Key string) string {
		sessionAdapter, _ := ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		session, _ := sessionAdapter.GetByKey(Key)
		csrfToken, _ := session.Get("csrf_token")
		return core.MaskCSRFToken(csrfToken)
	}
}

func (b Blueprint) InitApp(app core.IApp) {
	b.SessionAdapterRegistry.RegisterNewAdapter(&interfaces2.DbSession{}, true)
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.Session{}, &[]*core.Session{} })
}

var ConcreteBlueprint = Blueprint{
	Blueprint: core.Blueprint{
		Name:              "sessions",
		Description:       "Sessions blueprint responsible to keep session data in database",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
	SessionAdapterRegistry: interfaces2.NewSessionRegistry(),
}
