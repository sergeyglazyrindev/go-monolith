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
)

const tokenHeaderName = "AUTH-TOKEN"

type TokenAuthProvider struct {
}

func (ap *TokenAuthProvider) GetUserFromRequest(c *gin.Context) core.IUser {
	header := c.GetHeader(tokenHeaderName)
	if header == "" {
		return nil
	}
	database := core.NewDatabaseInstance()
	defer database.Close()
	token := &core.UserAuthToken{}
	database.Db.Model(token).Where(&core.UserAuthToken{Token: header}).First(token)
	return &token.User
}

// swagger:parameters tokenSignin
type TokenLoginParams struct {
	// SigninByField     string `form:"username" json:"username" xml:"username"  binding:"required"`
	SigninField string `form:"signinfield" json:"signinfield" xml:"signinfield"  binding:"required"`
	Password    string `form:"password" json:"password" xml:"password" binding:"required"`
	OTP         string `form:"otp" json:"otp" xml:"otp" binding:"omitempty"`
}

// swagger:parameters tokenSignup
type TokenSignupParams struct {
	Username          string `form:"username" json:"username" xml:"username"  binding:"required" valid:"username-unique"`
	Email             string `form:"email" json:"email" xml:"email"  binding:"required" valid:"email,email-unique"`
	Password          string `form:"password" json:"password" xml:"password" binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password" binding:"required"`
}

// A UserApiResponse is a serialized view of the user
// swagger:response userApiResponse
type UserAPIResponse struct {
	Name string
	ID   uint
}

// swagger:route POST /auth/token/signin auth token tokenSignin
//
// Signs in user with provided credentials
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       default: generalError
//       200: userApiResponse
//       400: validationError
func (ap *TokenAuthProvider) Signin(c *gin.Context) {
	var json TokenLoginParams
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponse(err.Error()))
		return
	}
	db := core.NewDatabaseInstance()
	defer db.Close()
	var user core.User
	// @todo, complete
	directAPISigninByField := core.CurrentConfig.D.GoMonolith.DirectAPISigninByField
	db.Db.Model(core.User{}).Where(fmt.Sprintf("%s = ?", directAPISigninByField), json.SigninField).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("login_credentials_incorrect", "login credentials are incorrect"))
		return
	}
	if !user.Active {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("user_inactive", "this user is inactive"))
		return
	}
	if !user.IsPasswordUsable {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("password_is_not_configured", "this user doesn't have a password"))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password+user.Salt))
	if err != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("login_credentials_incorrect", "login credentials are incorrect"))
		return
	}
	if user.GeneratedOTPToVerify != "" {
		if json.OTP == "" {
			c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("otp_required", "otp is required"))
			return
		}
		if user.GeneratedOTPToVerify != json.OTP {
			c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("otp_is_wrong", "otp provided by user is wrong"))
			return
		}
		user.GeneratedOTPToVerify = ""
		db.Db.Save(&user)
	}
	userAuthToken := &core.UserAuthToken{
		User: user,
	}
	db1 := db.Db.Model(&core.UserAuthToken{}).Where(&core.UserAuthToken{UserID: user.ID}).First(userAuthToken)
	if db1.Error != nil {
		userAuthToken = &core.UserAuthToken{
			User: user,
		}
		db1 = db.Db.Create(&userAuthToken)
		if db1.Error != nil {
			c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
			return
		}
	}
	if userAuthToken.Token == "" {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
		return
	}
	d := GetUserForAPIWithToken(&user, userAuthToken)
	c.JSON(http.StatusOK, d)
}

// swagger:route POST /auth/token/signup auth token tokenSignup
//
// Signs in user with provided credentials
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       default: generalError
//       200: userApiResponse
//       400: validationError
func (ap *TokenAuthProvider) Signup(c *gin.Context) {
	var json TokenSignupParams
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
	db := core.NewDatabaseInstance()
	defer db.Close()
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
	db.Db.Create(&user)
	userAuthToken := &core.UserAuthToken{
		User: user,
	}
	db1 := db.Db.Create(&userAuthToken)
	if db1.Error != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
		return
	}
	d := GetUserForAPIWithToken(&user, userAuthToken)
	c.JSON(http.StatusOK, d)
}

func (ap *TokenAuthProvider) Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (ap *TokenAuthProvider) IsAuthenticated(c *gin.Context) {
	token := c.GetHeader("X-" + core.CurrentConfig.D.GoMonolith.APIToken)
	if token == "" {
		c.JSON(http.StatusBadRequest, core.APIBadResponse("no token header passed"))
		return
	}
	db := core.NewDatabaseInstance()
	defer db.Close()
	tokenFromDb := &core.UserAuthToken{}
	db.Db.Model(tokenFromDb).Preload("User").Where(&core.UserAuthToken{Token: token}).First(tokenFromDb)
	if tokenFromDb.ID == 0 {
		c.JSON(http.StatusBadRequest, core.APIBadResponse("wrong token passed"))
		return
	}
	d := GetUserForAPIWithToken(&tokenFromDb.User, tokenFromDb)
	c.JSON(http.StatusOK, d)
}

func (ap *TokenAuthProvider) GetSession(c *gin.Context) core.ISessionProvider {
	var cookieName string
	cookieName = core.CurrentConfig.D.GoMonolith.APICookieName
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

func (ap *TokenAuthProvider) GetName() string {
	return "token"
}
