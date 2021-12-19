package admin

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

type AdminPaginationTestSuite struct {
	gomonolith.TestSuite
}

func (suite *AdminPaginationTestSuite) SetupTestData() {
	for i := range core.GenerateNumberSequence(1, 100) {
		userModel := core.GenerateUserModel()
		userModel.SetEmail(fmt.Sprintf("admin_%d@example.com", i))
		userModel.SetUsername("admin_" + strconv.Itoa(i))
		userModel.SetFirstName("firstname_" + strconv.Itoa(i))
		userModel.SetLastName("lastname_" + strconv.Itoa(i))
		suite.Database.Db.Create(userModel)
	}
}

func (suite *AdminPaginationTestSuite) TestPagination() {
	suite.SetupTestData()
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	var users = core.GenerateBunchOfUserModels()
	adminRequestParams := core.NewAdminRequestParams()
	adminUserPage.GetQueryset(nil, adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(users)
	assert.Equal(suite.T(), reflect.Indirect(reflect.ValueOf(users)).Len(), core.CurrentConfig.D.GoMonolith.AdminPerPage)
	adminRequestParams.Paginator.Offset = 88
	adminUserPage.GetQueryset(nil, adminUserPage, adminRequestParams).GetPaginatedQuerySet().Find(users)
	assert.Greater(suite.T(), reflect.Indirect(reflect.ValueOf(users)).Len(), core.CurrentConfig.D.GoMonolith.AdminPerPage)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminPagination(t *testing.T) {
	gomonolith.RunTests(t, new(AdminPaginationTestSuite))
}
