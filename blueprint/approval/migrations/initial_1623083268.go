package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/blueprint/approval/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm"
)

type initial1623083268 struct {
}

func (m initial1623083268) GetName() string {
	return "approval.1623083268"
}

func (m initial1623083268) GetID() int64 {
	return 1623083268
}

func (m initial1623083268) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(models.Approval{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623083268) Down(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.Migrator().DropTable(models.Approval{})
	if err != nil {
		return err
	}
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&models.Approval{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "approval", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	return nil
}

func (m initial1623083268) Deps() []string {
	return make([]string, 0)
}
