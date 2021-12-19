package sessions

import (
	"github.com/sergeyglazyrindev/go-monolith"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	interfaces2 "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SessionTestSuite struct {
	gomonolith.TestSuite
}

func (s *SessionTestSuite) TestSavingSession() {
	session := interfaces2.NewSession()
	session.SetData("testkey", "testvalue")
	s.Database.Db.Create(session)
	var loadedsession core.Session
	s.Database.Db.Model(&core.Session{}).First(&loadedsession)
	val, _ := loadedsession.GetData("testkey")
	assert.Equal(s.T(), val, "testvalue")
}

func (s *SessionTestSuite) TestTransactionConsistencyInTests() {
	var loadedsession core.Session
	s.Database.Db.Model(&core.Session{}).First(&loadedsession)
	val, _ := loadedsession.GetData("testkey")
	assert.Equal(s.T(), val, "")
}

func (s *SessionTestSuite) TestDbSessionAdapter() {
	blueprint, _ := s.App.BlueprintRegistry.GetByName("sessions")
	dbadapter, _ := blueprint.(sessionsblueprint.Blueprint).SessionAdapterRegistry.GetAdapter("db")
	dbadapter = dbadapter.Create()
	assert.Nil(s.T(), dbadapter.GetUser())
	dbadapter.Set("testkey", "testvalue")
	dbadapter.Save()
	dbadapter, _ = dbadapter.GetByKey(dbadapter.GetKey())
	val, _ := dbadapter.Get("testkey")
	assert.Equal(s.T(), val, "testvalue")
	dbadapter.ClearAll()
	dbadapter.Save()
	dbadapter, _ = dbadapter.GetByKey(dbadapter.GetKey())
	_, err := dbadapter.Get("testkey")
	if err == nil {
		assert.True(s.T(), false)
	}
	expiresOn := time.Now().UTC()
	expiresOn = expiresOn.Add(10 * time.Minute)
	dbadapter.ExpiresOn(&expiresOn)
	assert.False(s.T(), dbadapter.IsExpired())
	sessionKey := dbadapter.GetKey()
	isRemoved := dbadapter.Delete()
	if !isRemoved {
		assert.True(s.T(), false)
	}
	dbadapter, err = dbadapter.GetByKey(sessionKey)
	if err == nil {
		assert.True(s.T(), false)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSessions(t *testing.T) {
	gomonolith.RunTests(t, new(SessionTestSuite))
}
