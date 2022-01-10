package token_with_expiration

import (
	"bytes"
	"encoding/json"
	"fmt"
	gomonolith "github.com/sergeyglazyrindev/go-monolith"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/otp/services"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TokenWithExpirationAuthProviderTestSuite struct {
	gomonolith.TestSuite
}

type UserToken struct {
	Token string `json:"token"`
}

func (s *TokenWithExpirationAuthProviderTestSuite) TestTokenWithExpirationAuthProviderForApi() {
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
		token := &UserToken{}
		json.Unmarshal(w.Body.Bytes(), token)
		assert.NotEmpty(s.T(), token.Token)
		req1, _ := http.NewRequest("GET", "/auth/token-with-expiration/status/", nil)
		req1.Header.Set(
			"X-" + core.CurrentConfig.D.GoMonolith.APIToken,
			token.Token,
		)
		gomonolith.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "name")
			return strings.Contains(w.Body.String(), "name")
		})
		req1, _ = http.NewRequest("POST", "/auth/token-with-expiration/logout/", nil)
		req1.Header.Set(
			"X-" + core.CurrentConfig.D.GoMonolith.APIToken,
			token.Token,
		)
		session, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		session = session.Create()
		csrfToken := core.GenerateCSRFToken()
		session.Set("csrf_token", csrfToken)
		session.Save()
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.APICookieName, session.GetKey()),
		)
		tokenmasked := core.MaskCSRFToken(csrfToken)
		req1.Header.Set("CSRF-TOKEN", tokenmasked)
		gomonolith.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Equal(s.T(), w.Code, 204)
			return w.Code == 204
		})
		return w.Code == 200
	})
}

func (s *TokenWithExpirationAuthProviderTestSuite) TestSignupForApi() {
	// hashedPassword, err := utils2.HashPass(password, salt)
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "go-monolithapitest@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/token-with-expiration/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTokenWithExpirationAuthProvider(t *testing.T) {
	// gomonolith.RunTests(t, new(TokenWithExpirationAuthProviderTestSuite))
}
