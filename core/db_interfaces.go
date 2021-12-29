package core

import (
	"fmt"
	"gorm.io/gorm"
)

type DeleteRowStructure struct {
	SQL         string
	Values      []interface{}
	Explanation string
	Table       string
	Cond        string
}

type IDbAdapter interface {
	Equals(name interface{}, args ...interface{})
	GetDb(alias string, dryRun bool) (*gorm.DB, error)
	GetStringToExtractYearFromField(filterOptionField string) string
	GetStringToExtractMonthFromField(filterOptionField string) string
	Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	In(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Range(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Date(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Year(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Month(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Day(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Week(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Time(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Second(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	ArrayIncludes(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	JSONContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder)
	BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure
	SetIsolationLevelForTests(db *gorm.DB)
	Close(db *gorm.DB)
	ClearTestDatabase()
	SetTimeZone(db *gorm.DB, timezone string)
	InitializeDatabaseForTests(databaseSettings *DBSettings)
	StartDBShell(databaseSettings *DBSettings) error
	GetLastError() error
	ResetLastError()
}

var Db *gorm.DB

type ProjectDatabase struct {
	Db      *gorm.DB
	Adapter IDbAdapter
}

func (uad *ProjectDatabase) Close() {
	uad.Adapter.Close(uad.Db)
}

func (uad *ProjectDatabase) ForcefullyClose() {
	db1, _ := uad.Db.DB()
	db1.Close()
}

var TestInstanceDatabase *ProjectDatabase

func NewDatabaseInstance(alias1 ...string) *ProjectDatabase {
	if CurrentConfig.InTests && TestInstanceDatabase != nil {
		return TestInstanceDatabase
	}
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	adapter := GetAdapterForDb(alias)
	Db, _ = adapter.GetDb(
		alias, false,
	)
	return &ProjectDatabase{Db: Db, Adapter: adapter}
}

func NewDatabaseInstanceWithoutConnection(alias1 ...string) *ProjectDatabase {
	if CurrentConfig.InTests && TestInstanceDatabase != nil {
		return TestInstanceDatabase
	}
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	adapter := GetAdapterForDb(alias)
	Db, _ = adapter.GetDb(
		alias, true,
	)
	return &ProjectDatabase{Db: Db, Adapter: adapter}
}

type Database struct {
	config    *Config
	databases map[string]*ProjectDatabase
}

func NewDatabase(config *Config) *Database {
	database := Database{}
	database.config = config
	database.databases = make(map[string]*ProjectDatabase)
	return &database
}

func (d Database) ConnectTo(alias string) *gorm.DB {
	if alias == "" {
		alias = "default"
	}
	return GetDB(alias)
}

type DatabaseSettings struct {
	Default *DBSettings
	Slave   *DBSettings
}

var CurrentDatabaseSettings *DatabaseSettings

// GetDB returns a pointer to the DB
func GetDB(alias1 ...string) *gorm.DB {
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	var err error

	// Check if there is a database config file
	dialect := GetAdapterForDb(alias)
	Db, err = dialect.GetDb(
		alias, false,
	)
	if err != nil {
		Trail(ERROR, "unable to connect to DB. %s", err)
		Db.Error = fmt.Errorf("unable to connect to DB. %s", err)
	}
	return Db
}

func GetAdapterForDb(alias1 ...string) IDbAdapter {
	var databaseConfig *DBSettings
	var alias string
	if len(alias1) == 0 {
		alias = "default"
	} else {
		alias = alias1[0]
	}
	if alias == "default" {
		databaseConfig = CurrentDatabaseSettings.Default
	} else {
		databaseConfig = CurrentDatabaseSettings.Slave
	}
	return NewDbAdapter(Db, databaseConfig.Type)
}

type DbAdapterRegistry struct {
	dbTypeToAdapter map[string]func(db *gorm.DB) IDbAdapter
}

func (dar *DbAdapterRegistry) RegisterAdapter(dbType string, createAdapterHandler func(db *gorm.DB) IDbAdapter) {
	dar.dbTypeToAdapter[dbType] = createAdapterHandler
}

var GlobalDbAdapterRegistry *DbAdapterRegistry

func InitializeGlobalAdapterRegistry() {
	if GlobalDbAdapterRegistry == nil {
		GlobalDbAdapterRegistry = &DbAdapterRegistry{
			dbTypeToAdapter: make(map[string]func(db *gorm.DB) IDbAdapter),
		}
	}
}

func NewDbAdapter(db *gorm.DB, dbType string) IDbAdapter {
	adapter, ok := GlobalDbAdapterRegistry.dbTypeToAdapter[dbType]
	if !ok {
		panic("no adapter " + dbType + " has been registered")
	}
	return adapter(db)
}
