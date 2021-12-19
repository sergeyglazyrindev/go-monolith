package migrations

import (
	abtestmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/abtest/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm"
)

type initial1623083395 struct {
}

func (m initial1623083395) GetName() string {
	return "abtest.1623083395"
}

func (m initial1623083395) GetID() int64 {
	return 1623083395
}

func (m initial1623083395) Up(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.AutoMigrate(abtestmodel.ABTest{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(abtestmodel.ABTestValue{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1623083395) Down(database *core.ProjectDatabase) error {
	db := database.Db
	err := db.Migrator().DropTable(abtestmodel.ABTestValue{})
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(abtestmodel.ABTest{})
	if err != nil {
		return err
	}
	var contentType core.ContentType
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&abtestmodel.ABTestValue{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	stmt = &gorm.Statement{DB: db}
	stmt.Parse(&abtestmodel.ABTest{})
	db.Model(&core.ContentType{}).Where(&core.ContentType{BlueprintName: "abtest", ModelName: stmt.Schema.Table}).First(&contentType)
	db.Unscoped().Where(" 1 = 1").Delete(&contentType)
	return nil
}

func (m initial1623083395) Deps() []string {
	return make([]string, 0)
}
