package auth

import (
	"bytes"
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/otp/services"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/sergeyglazyrindev/go-monolith/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type AuthProviderTestSuite struct {
	gomonolith.TestSuite
}

func (s *AuthProviderTestSuite) TestDirectAuthProviderForGoMonolithAdmin() {
	req, _ := http.NewRequest("GET", "/auth/direct-for-admin/status/", nil)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "http: named cookie not present")
		return strings.Contains(w.Body.String(), "http: named cookie not present")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, ""),
	)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "empty cookie passed")
		return strings.Contains(w.Body.String(), "empty cookie passed")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, "test"),
	)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "no session with key")
		return strings.Contains(w.Body.String(), "no session with key")
	})
	sessionsblueprint1, _ := s.App.BlueprintRegistry.GetByName("sessions")
	sessionAdapterRegistry := sessionsblueprint1.(sessionsblueprint.Blueprint).SessionAdapterRegistry
	defaultAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	defaultAdapter = defaultAdapter.Create()
	expiresOn := time.Now().UTC().Add(-5 * time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	// directProvider.
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, defaultAdapter.GetKey()),
	)
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "session expired")
		return strings.Contains(w.Body.String(), "session expired")
	})
	expiresOn = time.Now().UTC()
	expiresOn = expiresOn.Add(10 * time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	req.URL = &url.URL{
		Path: "/auth/direct-for-admin/status/",
	}
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Body.String(), "{}\n")
		return w.Body.String() == "{}\n"
	})
	var jsonStr = []byte(`{"signinfield":"test", "password": "123456"}`)
	req, _ = http.NewRequest("POST", "/auth/direct-for-admin/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "login credentials are incorrect")
		return strings.Contains(w.Body.String(), "login credentials are incorrect")
	})
	salt := core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass("123456", salt)
	user := core.User{
		FirstName:        "testuser-firstname",
		LastName:         "testuser-lastname",
		Username:         "test",
		Password:         hashedPassword,
		Active:           false,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db := s.Database.Db
	db.Create(&user)
	req, _ = http.NewRequest("POST", "/auth/direct-for-admin/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "this user is inactive")
		return strings.Contains(w.Body.String(), "this user is inactive")
	})
	user.Active = true
	secretString, _ := services.GenerateOTPSeed(core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod, &user)
	user.OTPSeed = secretString
	user.IsSuperUser = true
	otpPassword := services.GetOTP(user.OTPSeed, core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod)
	user.GeneratedOTPToVerify = otpPassword
	var jsonStrForSignup = []byte(fmt.Sprintf(`{"signinfield":"test", "password": "123456", "otp": "%s"}`, otpPassword))
	db.Save(&user)
	req, _ = http.NewRequest("POST", "/auth/direct-for-admin/signin/", bytes.NewBuffer(jsonStrForSignup))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "go-monolith-admin=")
		sessionKey := strings.Split(strings.Split(w.Header().Get("Set-Cookie"), ";")[0], "=")[1]
		req1, _ := http.NewRequest("GET", "/auth/direct-for-admin/status/", nil)
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionKey),
		)
		gomonolith.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "for-gomonolith-panel")
			return strings.Contains(w.Body.String(), "for-gomonolith-panel")
		})
		req, _ = http.NewRequest("GET", "/admin/profile/", bytes.NewBuffer([]byte("")))
		req.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionKey),
		)
		gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "oldPassword")
			assert.Contains(s.T(), w.Body.String(), "<form")
			return strings.Contains(w.Body.String(), "oldPassword") && strings.Contains(w.Body.String(), "<form")
		})
		return strings.Contains(w.Header().Get("Set-Cookie"), "go-monolith-admin=")
	})
}

func (s *AuthProviderTestSuite) TestSignupForGoMonolithAdmin() {
	// hashedPassword, err := utils2.HashPass(password, salt)
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "go-monolithapitest@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/direct-for-admin/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "go-monolith-admin=")
		return strings.Contains(w.Header().Get("Set-Cookie"), "go-monolith-admin=")
	})
}

func (s *AuthProviderTestSuite) TestSignupForApi() {
	// hashedPassword, err := utils2.HashPass(password, salt)
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "go-monolithapitest@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/direct/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

func (s *AuthProviderTestSuite) TestTokenAuthProviderForApi() {
	var jsonStr = []byte(`{"signinfield":"test", "password": "123456"}`)
	req, _ := http.NewRequest("POST", "/auth/token/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "login credentials are incorrect")
		return strings.Contains(w.Body.String(), "login credentials are incorrect")
	})
	salt := core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass("123456", salt)
	user := core.User{
		FirstName:        "testuser-firstname",
		LastName:         "testuser-lastname",
		Username:         "test",
		Password:         hashedPassword,
		Active:           false,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db := s.Database.Db
	db.Create(&user)
	req, _ = http.NewRequest("POST", "/auth/token/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "this user is inactive")
		return strings.Contains(w.Body.String(), "this user is inactive")
	})
	user.Active = true
	secretString, _ := services.GenerateOTPSeed(core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod, &user)
	user.OTPSeed = secretString
	otpPassword := services.GetOTP(user.OTPSeed, core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod)
	user.GeneratedOTPToVerify = otpPassword
	var jsonStrForSignup = []byte(fmt.Sprintf(`{"signinfield":"test", "password": "123456", "otp": "%s"}`, otpPassword))
	db.Save(&user)
	req, _ = http.NewRequest("POST", "/auth/token/signin/", bytes.NewBuffer(jsonStrForSignup))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

func (s *AuthProviderTestSuite) TestTokenWithExpirationAuthProviderForApi() {
	var jsonStr = []byte(`{"signinfield":"test", "password": "123456"}`)
	req, _ := http.NewRequest("POST", "/auth/token-with-expiration/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "login credentials are incorrect")
		return strings.Contains(w.Body.String(), "login credentials are incorrect")
	})
	salt := core.GenerateRandomString(core.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass("123456", salt)
	user := core.User{
		FirstName:        "testuser-firstname",
		LastName:         "testuser-lastname",
		Username:         "test",
		Password:         hashedPassword,
		Active:           false,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db := s.Database.Db
	db.Create(&user)
	req, _ = http.NewRequest("POST", "/auth/token-with-expiration/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "this user is inactive")
		return strings.Contains(w.Body.String(), "this user is inactive")
	})
	user.Active = true
	secretString, _ := services.GenerateOTPSeed(core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod, &user)
	user.OTPSeed = secretString
	otpPassword := services.GetOTP(user.OTPSeed, core.CurrentConfig.D.GoMonolith.OTPDigits, core.CurrentConfig.D.GoMonolith.OTPAlgorithm, core.CurrentConfig.D.GoMonolith.OTPSkew, core.CurrentConfig.D.GoMonolith.OTPPeriod)
	user.GeneratedOTPToVerify = otpPassword
	var jsonStrForSignup = []byte(fmt.Sprintf(`{"signinfield":"test", "password": "123456", "otp": "%s"}`, otpPassword))
	db.Save(&user)
	req, _ = http.NewRequest("POST", "/auth/token-with-expiration/signin/", bytes.NewBuffer(jsonStrForSignup))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

func (s *AuthProviderTestSuite) TestOpenAdminPage() {
	req, _ := http.NewRequest("GET", core.CurrentConfig.D.GoMonolith.RootAdminURL+"/", nil)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "Go-Monolith - Admin Login")
		assert.Equal(s.T(), w.Code, 200)
		return strings.Contains(w.Body.String(), "Go-Monolith - Admin Login")
	})
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "go-monolithapitest@example.com"}`)
	req, _ = http.NewRequest("POST", "/auth/direct-for-admin/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "go-monolith-admin=")
		sessionKey := strings.Split(strings.Split(w.Header().Get("Set-Cookie"), ";")[0], "=")[1]
		req1, _ := http.NewRequest("GET", core.CurrentConfig.D.GoMonolith.RootAdminURL+"/", nil)
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionKey),
		)
		gomonolith.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "Go-Monolith - Dashboard")
			assert.Equal(s.T(), w.Code, 200)
			return strings.Contains(w.Body.String(), "Go-Monolith - Dashboard")
		})
		return strings.Contains(w.Header().Get("Set-Cookie"), "go-monolith-admin=")
	})
}

func (s *AuthProviderTestSuite) TestForgotFunctionality() {
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "go-monolithapitest@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/direct-for-admin/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "go-monolith-admin=")
		sessionKey := strings.Split(strings.Split(w.Header().Get("Set-Cookie"), ";")[0], "=")[1]
		session, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		session, _ = session.GetByKey(sessionKey)
		token := core.GenerateCSRFToken()
		session.Set("csrf_token", token)
		session.Save()
		var jsonStr1 = []byte(`{"email": "go-monolithapitest@example.com"}`)
		req1, _ := http.NewRequest("POST", "/user/api/forgot/", bytes.NewBuffer(jsonStr1))
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionKey),
		)
		tokenmasked := core.MaskCSRFToken(token)
		req1.Header.Set("CSRF-TOKEN", tokenmasked)
		gomonolith.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			isSentEmail := utils.SentEmailsDuringTests.IsAnyEmailSentWithStringInBodyOrSubject(&utils.SentEmail{
				Subject: "Password reset for admin panel",
			})
			var oneTimeAction core.OneTimeAction
			db := s.Database.Db
			db.Model(core.OneTimeAction{}).First(&oneTimeAction)
			var jsonStr2 = []byte(fmt.Sprintf(`{"code": "%s", "password": "1234567890", "confirm_password": "1234567890"}`, oneTimeAction.Code))
			req2, _ := http.NewRequest("POST", "/user/api/reset-password/", bytes.NewBuffer(jsonStr2))
			req2.Header.Set(
				"Cookie",
				fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, sessionKey),
			)
			tokenmasked1 := core.MaskCSRFToken(token)
			req2.Header.Set("CSRF-TOKEN", tokenmasked1)
			gomonolith.TestHTTPResponse(s.T(), s.App, req2, func(w *httptest.ResponseRecorder) bool {
				var oneTimeAction1 core.OneTimeAction
				db1 := s.Database.Db
				db1.Model(core.OneTimeAction{}).First(&oneTimeAction1)
				assert.True(s.T(), oneTimeAction1.IsUsed)
				assert.Equal(s.T(), w.Code, 200)
				return w.Code == 200
			})
			assert.True(s.T(), isSentEmail)
			return isSentEmail
		})
		return strings.Contains(w.Header().Get("Set-Cookie"), "go-monolith-admin=")
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAuthAdapters(t *testing.T) {
	gomonolith.RunTests(t, new(AuthProviderTestSuite))
}
