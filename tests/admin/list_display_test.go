package admin

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

type AdminListDisplayTestSuite struct {
	gomonolith.TestSuite
}

func (suite *AdminListDisplayTestSuite) TestListDisplay() {
	userModel := core.MakeUser()
	userModel.SetUsername("admin")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	listDisplayUsername, _ := adminUserPage.ListDisplay.GetFieldByDisplayName("Username")
	assert.Equal(suite.T(), listDisplayUsername.GetValue(userModel), template.HTML("admin"))
	compositeField := core.NewListDisplay(nil)
	compositeField.MethodName = "GetFullName"
	compositeField = core.NewListDisplay(nil)
	compositeField.Populate = func(m interface{}) string {
		return m.(core.IUser).GetFullName()
	}
	assert.Equal(suite.T(), compositeField.GetValue(userModel), template.HTML("firstname lastname"))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminListDisplay(t *testing.T) {
	gomonolith.RunTests(t, new(AdminListDisplayTestSuite))
}
