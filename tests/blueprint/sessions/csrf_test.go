package sessions

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	interfaces2 "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CsrfTestSuite struct {
	gomonolith.TestSuite
}

func (s *CsrfTestSuite) TestSuccessfulCsrfCheck() {
	session := interfaces2.NewSession()
	token := core.GenerateCSRFToken()
	session.SetData("csrf_token", token)
	s.Database.Db.Create(session)
	req, _ := http.NewRequest("POST", "/testcsrf/", nil)
	tokenmasked := core.MaskCSRFToken(token)
	req.Header.Set("CSRF-TOKEN", tokenmasked)
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, session.Key),
	)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
	req, _ = http.NewRequest("POST", "/testcsrf/", nil)
	req.Header.Set("CSRF-TOKEN", "dsadsada")
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, session.Key),
	)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		body := w.Body.String()
		assert.Equal(s.T(), body, "Incorrect length of csrf-token")
		return strings.EqualFold(body, "Incorrect length of csrf-token")
	})
}

func (s *CsrfTestSuite) TestIgnoreCsrfCheck() {
	req, _ := http.NewRequest("POST", "/ignorecsrfcheck/", nil)
	gomonolith.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCsrf(t *testing.T) {
	gomonolith.RunTests(t, new(CsrfTestSuite))
}
