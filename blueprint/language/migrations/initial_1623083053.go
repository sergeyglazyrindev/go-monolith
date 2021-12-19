package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm"
)

type initial1623083053 struct {
}

func (m initial1623083053) GetName() string {
	return "language.1623083053"
}

func (m initial1623083053) GetID() int64 {
	return 1623083053
}

func (m initial1623083053) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(core.Language{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623083053) Down(database *core.ProjectDatabase) error {
	db := database.Db
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&core.Language{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "language", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where("1 = 1").Delete(&contentType)
	err := db.Migrator().DropTable(core.Language{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623083053) Deps() []string {
	return make([]string, 0)
}
