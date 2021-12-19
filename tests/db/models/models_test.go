package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type BuildRemovalTreeTestSuite struct {
	gomonolith.TestSuite
	ContentType *core.ContentType
}

func (s *BuildRemovalTreeTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.Database.Db.AutoMigrate(&UserGroupContentType{})
	s.Database.Db.AutoMigrate(&UserContentType{})
	s.Database.Db.AutoMigrate(&OneTimeActionContentType{})
	s.Database.Db.AutoMigrate(&SessionContentType{})
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &SessionContentType{}, &[]*SessionContentType{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &OneTimeActionContentType{}, &[]*OneTimeActionContentType{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &UserContentType{}, &[]*UserContentType{} })
	core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &UserGroupContentType{}, &[]*UserGroupContentType{} })
}

func (s *BuildRemovalTreeTestSuite) ConfigureData(database *core.ProjectDatabase) {
	contentType := &core.ContentType{BlueprintName: "user", ModelName: "user"}
	database.Db.Create(contentType)
	s.ContentType = contentType
	permission := &core.Permission{Name: "user_read", ContentType: *contentType}
	database.Db.Create(permission)
	usergroup := &core.UserGroup{GroupName: "test"}
	database.Db.Create(usergroup)
	database.Db.Model(usergroup).Association("Permissions").Append(permission)
	database.Db.Save(usergroup)
	user := &core.User{Email: "adminmodelstest@example.com"}
	database.Db.Create(user)
	database.Db.Model(user).Association("Permissions").Append(permission)
	database.Db.Model(user).Association("UserGroups").Append(usergroup)
	database.Db.Save(user)
	oneTimeAction := &core.OneTimeAction{User: *user, Code: "aaa"}
	database.Db.Create(oneTimeAction)
	session := &core.Session{User: user, LoginTime: time.Now(), LastLogin: time.Now()}
	database.Db.Create(session)
	sessionContentType := &SessionContentType{Session: *session, ContentType: *contentType}
	database.Db.Save(sessionContentType)
	oneTimeActionContentType := &OneTimeActionContentType{OneTimeAction: *oneTimeAction, ContentType: *contentType}
	database.Db.Save(oneTimeActionContentType)
	userContentType := &UserContentType{User: *user, ContentType: *contentType}
	database.Db.Save(userContentType)
	userGroupContentType := &UserGroupContentType{UserGroup: *usergroup, ContentType: *contentType}
	database.Db.Save(userGroupContentType)
}

func (s *BuildRemovalTreeTestSuite) TearDownSuite() {
	s.Database.Db.Migrator().DropTable(&UserGroupContentType{})
	s.Database.Db.Migrator().DropTable(&UserContentType{})
	s.Database.Db.Migrator().DropTable(&OneTimeActionContentType{})
	s.Database.Db.Migrator().DropTable(&SessionContentType{})
	s.TestSuite.TearDownSuite()
}

type UserGroupContentType struct {
	core.Model
	UserGroup     core.UserGroup
	UserGroupID   uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *UserGroupContentType) String() string {
	return fmt.Sprintf("dsadsa-usergroup-%d-%s", ugct.ID, ugct.ContentType.String())
}

type UserContentType struct {
	core.Model
	User          core.User
	UserID        uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *UserContentType) String() string {
	return fmt.Sprintf("dsadsa-user-%d-%s", ugct.ID, ugct.ContentType.String())
}

type OneTimeActionContentType struct {
	core.Model
	OneTimeAction   core.OneTimeAction
	OneTimeActionID uint
	ContentType     core.ContentType
	ContentTypeID   uint
}

func (ugct *OneTimeActionContentType) String() string {
	return fmt.Sprintf("dsadsa-onetimeaction-%d-%s", ugct.ID, ugct.ContentType.String())
}

type SessionContentType struct {
	core.Model
	Session       core.Session
	SessionID     uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *SessionContentType) String() string {
	return fmt.Sprintf("dsadsa-session-%d-%s", ugct.ID, ugct.ContentType.String())
}

func (s *BuildRemovalTreeTestSuite) TestRemovalStringified() {
	s.ConfigureData(s.Database)
	//spew.Dump("contentType", contentType.ID)
	//spew.Dump("permission", permission.ID)
	//spew.Dump("usergroup", usergroup.ID)
	//spew.Dump("user", user.ID)
	//spew.Dump("onetimeaction", oneTimeAction.ID)
	//spew.Dump("session", session.ID)
	//spew.Dump("usergroup permissions", len(usergroup.Permissions))
	//spew.Dump("user permissions", len(user.Permissions))
	//spew.Dump("user groups", len(user.UserGroups))
	removalTreeNode := core.BuildRemovalTree(s.Database, s.ContentType)
	deletionStringified := removalTreeNode.BuildDeletionTreeStringified(s.Database)
	assert.Equal(s.T(), len(deletionStringified), 15)
}

func (s *BuildRemovalTreeTestSuite) TestRemoval() {
	s.ConfigureData(s.Database)
	var c int64
	removalTreeNode := core.BuildRemovalTree(s.Database, s.ContentType)
	removalTreeNode.RemoveFromDatabase(s.Database)
	s.Database.Db.Model(&core.Permission{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	s.Database.Db.Model(&OneTimeActionContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	s.Database.Db.Model(&UserContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	s.Database.Db.Model(&UserGroupContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	s.Database.Db.Model(&SessionContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBuildRemovalTree(t *testing.T) {
	gomonolith.RunTests(t, new(BuildRemovalTreeTestSuite))
}
