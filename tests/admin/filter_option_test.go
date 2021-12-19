package admin

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type FilterOptionTestSuite struct {
	gomonolith.TestSuite
}

func (suite *FilterOptionTestSuite) BeforeTest(suiteName string, testMethod string) {
	if testMethod == "TestFilterOptionByYear" {
	} else {
	}
}

func (suite *FilterOptionTestSuite) Test() {
	userModel := core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption@example.com")
	createdAt := time.Now().Add((-10 * 12 * 86400 * 30) * time.Second)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption1")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption1@example.com")
	createdAt = time.Now().Add((-5 * 12 * 86400 * 30) * time.Second)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption2")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption2@example.com")
	createdAt = time.Now().Add((-3 * 12 * 86400 * 30) * time.Second)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption3")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption3@example.com")
	createdAt = time.Now().Add((-1 * 12 * 86400 * 30) * time.Second)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption := core.NewFilterOption()
	newFilterOption.FetchOptions = func(afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions(nil)) > 0)
	suite.Database.Db.Unscoped().Where("1 = 1").Delete(core.GenerateUserModel())
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroptionA")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption@example.com")
	createdAt = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	userModel.SetCreatedAt(&createdAt)
	// ((-10 * 12 * 86400 * 30) * time.Second)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption1A")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption1@example.com")
	createdAt = time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption2A")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption2@example.com")
	createdAt = time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	userModel = core.GenerateUserModel()
	userModel.SetUsername("adminfilteroption3A")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminfilteroption3@example.com")
	createdAt = time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC)
	userModel.SetCreatedAt(&createdAt)
	suite.Database.Db.Create(userModel)
	adminUserBlueprintPage, _ = core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ = adminUserBlueprintPage.SubPages.GetBySlug("user")
	newFilterOption = core.NewFilterOption()
	newFilterOption.FetchOptions = func(afo core.IAdminFilterObjects) []*core.DisplayFilterOption {
		return core.FetchOptionsFromGormModelFromDateTimeField(afo, "created_at")
	}
	adminUserPage.FilterOptions = core.NewFilterOptionsRegistry()
	adminUserPage.FilterOptions.AddFilterOption(newFilterOption)
	assert.True(suite.T(), len(adminUserPage.FetchFilterOptions(nil)) > 0)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFilterOption(t *testing.T) {
	gomonolith.RunTests(t, new(FilterOptionTestSuite))
}
