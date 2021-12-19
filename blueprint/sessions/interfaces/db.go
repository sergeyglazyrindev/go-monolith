package interfaces

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
	"time"
)

type DbSession struct {
	session *core.Session
}

func (s *DbSession) GetName() string {
	return "db"
}

func (s *DbSession) Set(name string, value string) {
	s.session.SetData(name, value)
}

func (s *DbSession) SetUser(user core.IUser) {
	userID := user.GetID()
	s.session.UserID = &userID
	s.session.User = user.(*core.User)
}

func (s *DbSession) Get(name string) (string, error) {
	return s.session.GetData(name)
}

func (s *DbSession) ClearAll() bool {
	return s.session.ClearAll()
}

func (s *DbSession) ExpiresOn(expiresOn *time.Time) {
	s.session.ExpiresOn = expiresOn
}

func (s *DbSession) GetKey() string {
	return s.session.Key
}

func (s *DbSession) GetUser() core.IUser {
	if s.session == nil {
		return nil
	}
	if s.session.UserID == nil {
		return nil
	}
	if s.session.User == nil {
		return nil
	}
	return s.session.User
}

func (s *DbSession) GetByKey(sessionKey string) (core.ISessionProvider, error) {
	db := core.NewDatabaseInstance()
	defer db.Close()
	var session core.Session
	db.Db.Model(&core.Session{}).Where(&core.Session{Key: sessionKey}).Preload("User").First(&session)
	if session.ID == 0 {
		return nil, core.NewHTTPErrorResponse("session_not_found", "no session with key %s found", sessionKey)
	}
	return &DbSession{
		session: &session,
	}, nil
}

func (s *DbSession) Create() core.ISessionProvider {
	session := NewSession()
	db := core.NewDatabaseInstance()
	defer db.Close()
	db.Db.Create(session)
	return &DbSession{
		session: session,
	}
}

func (s *DbSession) Delete() bool {
	db := core.NewDatabaseInstance()
	defer db.Close()
	db.Db.Unscoped().Delete(&core.Session{}, s.session.ID)
	return db.Db.Error == nil
}

func (s *DbSession) IsExpired() bool {
	if s.session == nil || s.session.ExpiresOn == nil {
		return true
	}
	return s.session.ExpiresOn.Before(time.Now().UTC())
}

func (s *DbSession) Save() bool {
	db := core.NewDatabaseInstance()
	defer db.Close()
	res := db.Db.Save(s.session)
	return res.Error == nil
}
