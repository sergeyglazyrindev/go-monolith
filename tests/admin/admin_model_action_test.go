package admin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type AdminModelActionTestSuite struct {
	gomonolith.TestSuite
}

func (suite *AdminModelActionTestSuite) TestAdminModelAction() {
	userModel := core.MakeUser()
	userModel.SetUsername("adminmodelactiontest")
	userModel.SetFirstName("firstname")
	userModel.SetLastName("lastname")
	userModel.SetIsSuperUser(true)
	userModel.SetEmail("adminmodelactiontest@gmail.com")
	suite.Database.Db.Create(userModel)
	adminUserBlueprintPage, _ := core.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	adminModelAction := core.NewAdminModelAction(
		"TurnSuperusersToNormalUsers", &core.AdminActionPlacement{},
	)
	adminModelAction.Handler = func(ap *core.AdminPage, afo core.IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
		tx := afo.GetFullQuerySet().Update("IsSuperUser", false).Commit()
		return tx.(*core.GormPersistenceStorage).Db.Error == nil, tx.(*core.GormPersistenceStorage).Db.RowsAffected
	}
	adminUserPage.ModelActionsRegistry.AddModelAction(adminModelAction)
	var jsonStr = []byte(fmt.Sprintf(`{"object_ids": "%d"}`, userModel.GetID()))
	req, _ := http.NewRequest("POST", "/admin/users/user/turnsuperuserstonormalusers/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	sessionsblueprint1, _ := suite.App.BlueprintRegistry.GetByName("sessions")
	sessionAdapterRegistry := sessionsblueprint1.(sessionsblueprint.Blueprint).SessionAdapterRegistry
	defaultAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	defaultAdapter = defaultAdapter.Create()
	expiresOn := time.Now().UTC().Add(5 * time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	defaultAdapter.SetUser(userModel)
	defaultAdapter.Save()
	// directProvider.
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", core.CurrentConfig.D.GoMonolith.AdminCookieName, defaultAdapter.GetKey()),
	)
	for adminModelAction := range adminUserPage.ModelActionsRegistry.GetAllModelActions() {
		if adminModelAction.SlugifiedActionName != "turnsuperuserstonormalusers" {
			continue
		}
		suite.App.Router.Any(fmt.Sprintf("%s/%s/%s/%s/", core.CurrentConfig.D.GoMonolith.RootAdminURL, "users", adminUserPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *core.AdminPage, slugifiedModelActionName string) func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				adminPage.HandleModelAction(slugifiedModelActionName, ctx)
			}
		}(adminUserPage, adminModelAction.SlugifiedActionName))

	}
	adminContext := &core.AdminContext{}
	userForm := core.NewFormFromModelFromGinContext(adminContext, core.MakeUser(), make([]string, 0), []string{"ID"}, true, "")
	adminUserPage.Form = userForm
	gomonolith.TestHTTPResponse(suite.T(), suite.App, req, func(w *httptest.ResponseRecorder) bool {
		db := suite.Database.Db
		var user core.User
		db.Model(&core.User{}).First(&user)
		assert.False(suite.T(), user.IsSuperUser)
		return user.IsSuperUser == false
	})
	// adminUserPage.HandleModelAction("TurnSuperusersToNormalUsers", &gin.Context{})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminModelAction(t *testing.T) {
	gomonolith.RunTests(t, new(AdminModelActionTestSuite))
}
