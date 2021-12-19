---
sidebar_position: 1
---

# Database support

Right now gomonolith supports only sqlite, postgres databases, but it's easy to provide adapters for another databases, we just need to write implementation of the interface.
For non sqlite adapter, please build your package with corresponding tag: for postgres - with postgres.
```go
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
	BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure
	SetIsolationLevelForTests(db *gorm.DB)
	Close(db *gorm.DB)
	ClearTestDatabase()
	SetTimeZone(db *gorm.DB, timezone string)
	InitializeDatabaseForTests(databaseSettings *DBSettings)
	StartDBShell(databaseSettings *DBSettings) error
	GetLastError() error
}
```
Please also create ci-cd job for new database type you are adding to the framework.
And don't forget to handle this database type in the core.NewDbAdapter function.  
You can get instance of the Database using function:
```go
type Database struct {
	Db      *gorm.DB
	Adapter IDbAdapter
}

func (uad *Database) Close() {
	db, _ := uad.Db.DB()
	db.Close()
}
func NewDatabaseInstance(alias1 ...string) *Database {
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
	return &Database{Db: Db, Adapter: adapter}
}
DBInstance := NewDatabaseInstance()
defer DBInstance.Close()
// do whatever you want with database
```
