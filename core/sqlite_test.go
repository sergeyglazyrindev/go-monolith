package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

// User !
type UserTest struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	Password  string
	Email     string
	Active    bool
	Admin     bool
	Photo     string
	//Language     []Language `gorm:"many2many:user_languages" listExclude:"true"`
	LastLogin   *time.Time
	ExpiresOn   *time.Time
	OTPRequired bool
	OTPSeed     string
}

func GetDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("not able to open database: %s", "test.db"))
	}
	// Initialize system models
	modelList := []interface{}{
		UserTest{},
	}
	// Migrate schema
	for _, model := range modelList {
		db.AutoMigrate(model)
	}
	return db
}

func TestSqlite(t *testing.T) {
	db := GetDb()
	adapter := NewDbAdapter(db, "sqlite")
	adapter.Equals("admin", true)
}

func TestSqlite_gomonolith_datetime_cast_date(t *testing.T) {
	time1 := "2005-07-29 09:56:00.781963317+00:00"
	dt := sqliteGoMonolithDatetimeParse(time1, "UTC", "UTC")
	assert.Equal(t, dt.Year(), 2005)
	dt = sqliteGoMonolithDatetimeParse("", "UTC", "UTC")
	assert.Nil(t, dt)
}
