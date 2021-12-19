package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm"
)

type initial1621680132 struct {
}

func (m initial1621680132) GetName() string {
	return "user.1621680132"
}

func (m initial1621680132) GetID() int64 {
	return 1621680132
}

func (m initial1621680132) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(core.ContentType{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(core.UserGroup{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(core.GenerateUserModel())
	if err != nil {
		return err
	}
	err = db.AutoMigrate(core.Permission{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(core.OneTimeAction{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1621680132) Down(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.Migrator().DropTable(core.Permission{})
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(core.GenerateUserModel())
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(core.UserGroup{})
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(core.OneTimeAction{})
	if err != nil {
		return err
	}
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&core.OneTimeAction{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	stmt = &gorm.Statement{DB: db}
	stmt.Parse(core.GenerateUserModel())
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "user", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	return nil
}

func (m initial1621680132) Deps() []string {
	return []string{}
}
