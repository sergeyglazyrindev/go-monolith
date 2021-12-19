package interfaces

import (
	"database/sql"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	user2 "github.com/sergeyglazyrindev/go-monolith/blueprint/user"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type TokenWithExpirationAuthProvider struct {
}

func (ap *TokenWithExpirationAuthProvider) GetUserFromRequest(c *gin.Context) core.IUser {
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

func (ap *TokenWithExpirationAuthProvider) Signin(c *gin.Context) {
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
			SessionExpiresAt: sql.NullInt64{
				Int64: time.Now().UTC().Add(time.Duration(core.CurrentConfig.D.GoMonolith.SessionDuration) * time.Second).Unix(),
				Valid: true,
			},
		}
		db1 = db.Db.Create(&userAuthToken)
		if db1.Error != nil {
			c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
			return
		}
	} else {
		if userAuthToken.Token == "" {
			c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
			return
		}
		userAuthToken.SessionExpiresAt = sql.NullInt64{
			Int64: time.Now().UTC().Add(time.Duration(core.CurrentConfig.D.GoMonolith.SessionDuration) * time.Second).Unix(),
			Valid: true,
		}
		db.Db.Save(userAuthToken)
	}
	if userAuthToken.Token == "" {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
		return
	}
	c.JSON(http.StatusOK, GetUserForAPI(&user))
}

func (ap *TokenWithExpirationAuthProvider) Signup(c *gin.Context) {
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
	userAuthToken.SessionExpiresAt = sql.NullInt64{
		Int64: time.Now().UTC().Add(time.Duration(core.CurrentConfig.D.GoMonolith.SessionDuration) * time.Second).Unix(),
		Valid: true,
	}
	db1 := db.Db.Create(&userAuthToken)
	if db1.Error != nil {
		c.JSON(http.StatusBadRequest, core.APIBadResponseWithCode("no_token_has_been_created", "no token has been created"))
		return
	}
	c.JSON(http.StatusOK, GetUserForAPI(&user))
}

func (ap *TokenWithExpirationAuthProvider) Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (ap *TokenWithExpirationAuthProvider) IsAuthenticated(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (ap *TokenWithExpirationAuthProvider) GetSession(c *gin.Context) core.ISessionProvider {
	return nil
}

func (ap *TokenWithExpirationAuthProvider) GetName() string {
	return "token-with-expiration"
}
