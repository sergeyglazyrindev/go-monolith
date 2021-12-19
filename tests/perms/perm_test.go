package perms

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PermTestSuite struct {
	gomonolith.TestSuite
}

func (suite *PermTestSuite) TestIntegration() {
	db := suite.Database.Db
	contentType := core.ContentType{BlueprintName: "user", ModelName: "user"}
	db.Create(&contentType)
	permission := core.Permission{ContentType: contentType, PermissionBits: core.RevertPermBit}
	db.Create(&permission)
	g1 := core.UserGroup{GroupName: "usergroup"}
	db.Create(&g1)
	db.Model(&g1).Association("Permissions").Append(&permission)
	db.Save(&g1)
	permission = core.Permission{ContentType: contentType, PermissionBits: core.AddPermBit}
	db.Create(&permission)
	u := core.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	db.Create(&u)
	db.Model(&u).Association("Permissions").Append(&permission)
	db.Model(&u).Association("UserGroups").Append(&g1)
	db.Save(&u)
	var u1 core.User
	db.Model(&core.User{}).Where("email = 'ffsdfsd@example.com'").First(&u1)
	permRegistry := u1.BuildPermissionRegistry()
	userPerm := permRegistry.GetPermissionForBlueprint("user", "user")
	assert.True(suite.T(), userPerm.HasRevertPermission())
	assert.True(suite.T(), userPerm.HasAddPermission())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPermissionSystem(t *testing.T) {
	gomonolith.RunTests(t, new(PermTestSuite))
}
