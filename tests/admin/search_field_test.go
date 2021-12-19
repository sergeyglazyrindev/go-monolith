package admin

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"strconv"
	"testing"
)

type AdminSearchFieldTestSuite struct {
	gomonolith.TestSuite
}

func (suite *AdminSearchFieldTestSuite) SetupTestData() {
	for i := range core.GenerateNumberSequence(201, 300) {
		userModel := core.GenerateUserModel()
		userModel.SetEmail(fmt.Sprintf("admin_%d@example.com", i))
		userModel.SetUsername("admin_" + strconv.Itoa(i))
		userModel.SetFirstName("firstname_" + strconv.Itoa(i))
		userModel.SetLastName("lastname_" + strconv.Itoa(i))
		suite.Database.Db.Create(userModel)
	}
}

//func (suite *AdminSearchFieldTestSuite) TestFiltering() {
//	suite.SetupTestData()
//	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
//	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
//	var users []core.IUser
//	adminRequestParams := core.NewAdminRequestParams()
//	statement := &gorm.Statement{DB: suite.Database.Db}
//	statement.Parse(core.GenerateUserModel())
//	adminRequestParams.Search = "admin_202@example.com"
//	adminUserPage.GetQueryset(adminUserPage, adminRequestParams).GetFullQuerySet().Find(&users)
//	assert.Equal(suite.T(), len(users), 1)
//}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSearchField(t *testing.T) {
	gomonolith.RunTests(t, new(AdminSearchFieldTestSuite))
}
