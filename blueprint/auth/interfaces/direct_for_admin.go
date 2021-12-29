package interfaces

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	user2 "github.com/sergeyglazyrindev/go-monolith/blueprint/user"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"time"
)

// Binding from JSON
type LoginParamsForGoMonolithAdmin struct {
	// SigninByField     string `form:"signinbyfield" json:"signinbyfield" xml:"signinbyfield"  binding:"required"`
	SigninField string `form:"signinfield" json:"signinfield" xml:"signinfield"  binding:"required"`
	Password    string `form:"password" json:"password" xml:"password" binding:"required"`
	OTP         string `form:"otp" json:"otp" xml:"otp" binding:"omitempty"`
}

type SignupParamsForGoMonolithAdmin struct {
	Username          string `form:"username" json:"username" xml:"username"  binding:"required" valid:"username-unique"`
	Email             string `form:"email" json:"email" xml:"email"  binding:"required" valid:"email,email-unique"`
	Password          string `form:"password" json:"password" xml:"password" binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password" binding:"required"`
}

type DirectAuthForAdminProvider struct {
}

func (ap *DirectAuthForAdminProvider) GetUserFromRequest(c *gin.Context) core.IUser {
	session := ap.GetSession(c)
	if session != nil {
		return session.GetUser()
	}
	return nil
}

func (ap *DirectAuthForAdminProvider) Signin(c *gin.Context) {
	var json LoginParamsForGoMonolithAdmin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	database := core.NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	var user = core.GenerateUserModel()
	directAPISigninByField := core.CurrentConfig.D.GoMonolith.DirectAPISigninByField
	db.Model(core.User{}).Where(fmt.Sprintf("%s = ?", directAPISigninByField), json.SigninField).First(&user)
	if user.GetID() == 0 {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("login_credentials_incorrect", "login credentials are incorrect."))
		return
	}
	if !user.GetActive() {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("user_inactive", "this user is inactive"))
		return
	}
	if !user.GetIsSuperUser() && !user.GetIsStaff() {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("user_has_no_access_to_admin_panel", "this user doesn't have an access to admin panel"))
		return
	}
	if !user.GetIsPasswordUsable() {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("password_is_not_configured", "this user doesn't have a password"))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(json.Password+user.GetSalt()))
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("login_credentials_incorrect", "login credentials are incorrect."))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	cookieName := core.CurrentConfig.D.GoMonolith.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	sessionDuration := time.Duration(core.CurrentConfig.D.GoMonolith.SessionDuration) * time.Second
	sessionExpirationTime := time.Now().UTC().Add(sessionDuration)
	if cookie != "" {
		sessionAdapter, _ = sessionAdapter.GetByKey(cookie)
		sessionAdapter.ExpiresOn(&sessionExpirationTime)
		c.SetCookie(core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionAdapter.GetKey(), int(core.CurrentConfig.D.GoMonolith.SessionDuration), "/", c.Request.URL.Host, core.CurrentConfig.D.GoMonolith.SecureCookie, core.CurrentConfig.D.GoMonolith.HTTPOnlyCookie)
	} else {
		sessionAdapter = sessionAdapter.Create()
		sessionAdapter.ExpiresOn(&sessionExpirationTime)
		c.SetCookie(core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionAdapter.GetKey(), int(core.CurrentConfig.D.GoMonolith.SessionDuration), "/", c.Request.URL.Host, core.CurrentConfig.D.GoMonolith.SecureCookie, core.CurrentConfig.D.GoMonolith.HTTPOnlyCookie)
	}
	sessionAdapter.SetUser(user)
	sessionAdapter.Save()
	c.JSON(http.StatusOK, getUserForGoMonolithPanel(sessionAdapter.GetUser()))
}

func (ap *DirectAuthForAdminProvider) Signup(c *gin.Context) {
	var json SignupParamsForGoMonolithAdmin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	_, err := govalidator.ValidateStruct(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	passwordValidationStruct := &user2.PasswordValidationStruct{
		Password:          json.Password,
		ConfirmedPassword: json.ConfirmedPassword,
	}
	_, err = govalidator.ValidateStruct(passwordValidationStruct)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	//if utils.Contains(config.CurrentConfig.D.Auth.Twofactor_auth_required_for_signin_adapters, ap.GetName()) {
	//	if json.OTP == "" {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "otp is required"})
	//		return
	//	}
	//	if user.GeneratedOTPToVerify != json.OTP {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "otp provided by user is wrong"})
	//		return
	//	}
	//	user.GeneratedOTPToVerify = ""
	//	db.Save(&user)
	//}
	database := core.NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	salt := core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass(json.Password, salt)
	user := core.User{
		Username:         json.Username,
		Email:            json.Email,
		Password:         hashedPassword,
		Active:           true,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db.Create(&user)
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter = sessionAdapter.Create()
	sessionAdapter.SetUser(&user)
	sessionDuration := time.Duration(core.CurrentConfig.D.GoMonolith.SessionDuration) * time.Second
	sessionExpirationTime := time.Now().UTC().Add(sessionDuration)
	sessionAdapter.ExpiresOn(&sessionExpirationTime)
	sessionAdapter.Save()
	c.SetCookie(core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionAdapter.GetKey(), int(core.CurrentConfig.D.GoMonolith.SessionDuration), "/", c.Request.URL.Host, core.CurrentConfig.D.GoMonolith.SecureCookie, core.CurrentConfig.D.GoMonolith.HTTPOnlyCookie)
	c.JSON(http.StatusOK, getUserForGoMonolithPanel(sessionAdapter.GetUser()))
}

func (ap *DirectAuthForAdminProvider) Logout(c *gin.Context) {
	var cookie string
	var err error
	var cookieName string
	cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
	cookie, err = c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, core.APIBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err == nil {
		sessionAdapter.Delete()
	}
	timeInPast := time.Now().Add(-10 * time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(""),
		MaxAge:   0,
		Path:     "/",
		Domain:   c.Request.URL.Host,
		SameSite: http.SameSiteDefaultMode,
		Secure:   core.CurrentConfig.D.GoMonolith.SecureCookie,
		HttpOnly: core.CurrentConfig.D.GoMonolith.HTTPOnlyCookie,
		Expires:  timeInPast,
	})
	c.Status(http.StatusNoContent)
}

func (ap *DirectAuthForAdminProvider) IsAuthenticated(c *gin.Context) {
	var cookieName string
	cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, core.APIBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	if sessionAdapter.IsExpired() {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("session_expired", "session expired"))
		return
	}
	c.JSON(http.StatusOK, getUserForGoMonolithPanel(sessionAdapter.GetUser()))
}

func getUserForGoMonolithPanel(user core.IUser) *gin.H {
	if user == nil {
		return &gin.H{}
	}
	return &gin.H{"name": user.GetUsername(), "id": user.GetID(), "for-gomonolith-panel": true}
}

func (ap *DirectAuthForAdminProvider) GetSession(c *gin.Context) core.ISessionProvider {
	var cookieName string
	cookieName = core.CurrentConfig.D.GoMonolith.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil
	}
	if cookie == "" {
		return nil
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		return nil
	}
	if sessionAdapter.IsExpired() {
		return nil
	}
	return sessionAdapter
}

func (ap *DirectAuthForAdminProvider) GetName() string {
	return "direct-for-admin"
}
