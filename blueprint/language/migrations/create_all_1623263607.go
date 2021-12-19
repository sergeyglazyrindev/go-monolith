package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type createall1623263607 struct {
}

func (m createall1623263607) GetName() string {
	return "language.1623263607"
}

func (m createall1623263607) GetID() int64 {
	return 1623263607
}

func (m createall1623263607) Up(database *core.ProjectDatabase) error {
	langs := [][]string{
		{"English", "English", "en"},
	}
	db := database.Db
	tx := db
	for _, lang := range langs {
		l := core.Language{
			EnglishName: lang[0],
			Name:        lang[1],
			Code:        lang[2],
			Active:      false,
		}
		if l.Code == "en" {
			l.AvailableInGui = true
			l.Active = true
			l.Default = true
		}
		tx.Create(&l)
	}
	return nil
}

func (m createall1623263607) Down(database *core.ProjectDatabase) error {
	db := database.Db
	db.Unscoped().Where("1 = 1").Delete(&core.Language{Code: "en"})
	return nil
}

func (m createall1623263607) Deps() []string {
	return []string{"language.1623083053"}
}
