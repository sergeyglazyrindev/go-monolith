package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type initial1623082009 struct {
}

func (m initial1623082009) GetName() string {
	return "sessions.1623082009"
}

func (m initial1623082009) GetID() int64 {
	return 1623082009
}

func (m initial1623082009) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(core.Session{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623082009) Down(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.Migrator().DropTable(core.Session{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623082009) Deps() []string {
	return []string{"user.1621680132"}
}
