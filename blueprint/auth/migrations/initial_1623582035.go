package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type initial1623582035 struct {
}

func (m initial1623582035) GetName() string {
	return "auth.1623582035"
}

func (m initial1623582035) GetID() int64 {
	return 1623582035
}

func (m initial1623582035) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(core.UserAuthToken{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623582035) Down(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.Migrator().DropTable(core.UserAuthToken{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623582035) Deps() []string {
	return []string{"user.1621680132"}
}
