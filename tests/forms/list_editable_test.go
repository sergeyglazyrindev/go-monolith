package forms

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func NewTestForm1() *multipart.Form {
	form1 := multipart.Form{
		Value: make(map[string][]string),
	}
	return &form1
}

type ListEditableFormTestSuite struct {
	gomonolith.TestSuite
}

func (s *ListEditableFormTestSuite) TestFormBuilder() {
	// userBlueprintRegistry, _ := s.App.BlueprintRegistry.GetByName("user")
	// NewFormListEditableFromListDisplayRegistry
	adminPanel, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	userAdminPage, _ := adminPanel.SubPages.GetBySlug("user")
	ld, _ := userAdminPage.ListDisplay.GetFieldByDisplayName("Email")
	ld.IsEditable = true
	listEditableForm := core.NewFormListEditableFromListDisplayRegistry(nil, "", "10", core.MakeUser(), userAdminPage.ListDisplay, nil)
	form := NewTestForm1()
	userTest := &core.User{}
	err := listEditableForm.ProceedRequest(form, userTest, &core.AdminContext{})
	assert.False(s.T(), err.IsEmpty())
	form.Value["10_Email"] = []string{"adminlisteditable@example.com"}
	err = listEditableForm.ProceedRequest(form, userTest, &core.AdminContext{})
	assert.True(s.T(), err.IsEmpty())
	assert.Equal(s.T(), userTest.Email, "adminlisteditable@example.com")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestListEditableForm(t *testing.T) {
	gomonolith.RunTests(t, new(ListEditableFormTestSuite))
}
