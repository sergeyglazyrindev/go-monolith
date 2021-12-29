package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type SqliteAdapter struct {
	Statement *gorm.Statement
	DbType    string
	LastError error
}

func (d *SqliteAdapter) Equals(name interface{}, args ...interface{}) {
	query := d.Statement.Quote(name) + " = ?"
	clause.Expr{SQL: query, Vars: args}.Build(d.Statement)
}

func (d *SqliteAdapter) GetStringToExtractYearFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%Y\", %s)", filterOptionField)
}

func (d *SqliteAdapter) GetStringToExtractMonthFromField(filterOptionField string) string {
	return fmt.Sprintf("strftime(\"%%m\", %s)", filterOptionField)
}

// @todo analyze
func (d *SqliteAdapter) GetDb(alias string, dryRun bool) (*gorm.DB, error) {
	var aliasDatabaseSettings *DBSettings
	if alias == "default" {
		aliasDatabaseSettings = CurrentDatabaseSettings.Default
	}
	var db *gorm.DB
	var err error
	if CurrentConfig.InTests {
		if CurrentConfig.DebugTests || true {
			db, err = gorm.Open(sqlite.Dialector{DriverName: "GoMonolithSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
				DryRun:                                   dryRun,
			})
		} else {
			db, err = gorm.Open(sqlite.Dialector{DriverName: "GoMonolithSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
				DryRun:                                   dryRun,
				// Logger: logger.Default.LogMode(logger.Info),
			})
		}
	} else {
		db, err = gorm.Open(sqlite.Dialector{DriverName: "GoMonolithSqliteDriver", DSN: aliasDatabaseSettings.Name}, &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			DryRun:                                   dryRun,
			// Logger:                                   logger.Default.LogMode(logger.Info),
		})
	}
	db.Exec("PRAGMA case_sensitive_like = 1;")
	return db, err
}

func (d *SqliteAdapter) Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) = UPPER(?) ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE '%%' || UPPER(?) || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) In(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s IN ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s > ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s >= ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s < ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s <= ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE ? || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE UPPER(?) || '%%' ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE '%%' || UPPER(?) ESCAPE '\\' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Date(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_cast_date(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Year(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName)
	year := value.(int)
	startOfTheYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfTheYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, startOfTheYear, endOfTheYear)
}

func (d *SqliteAdapter) Month(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('month', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Day(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Week(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('week', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('week_day', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('quarter', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('hour', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('minute', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Second(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_extract('second', %s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_regex(%s.%s, ?) ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_regex(%s.%s, '(?i)' || ?) ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) Time(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" gomonolith_datetime_cast_time(%s.%s, 'UTC', 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *SqliteAdapter) IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	isTruthyValue := IsTruthyValue(value)
	isNull := " IS NULL "
	if !isTruthyValue {
		isNull = " IS NOT NULL "
	}
	query := fmt.Sprintf(" %s.%s %s ", operatorContext.TableName, field.DBName, isNull)
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query)
}

func (d *SqliteAdapter) Range(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	s := reflect.ValueOf(value)
	var f interface{}
	var second interface{}
	for i := 0; i < s.Len(); i++ {
		if i == 0 {
			f = s.Index(i).Interface()
		} else if i == 1 {
			second = s.Index(i).Interface()
			break
		}
	}
	query := fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName)
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, f, second)
}

func (d *SqliteAdapter) BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure {
	deleteRowStructure := &DeleteRowStructure{SQL: fmt.Sprintf("DELETE FROM %s WHERE %s", table, cond), Values: values}
	return deleteRowStructure
}

func (d *SqliteAdapter) SetIsolationLevelForTests(db *gorm.DB) {
}

func (d *SqliteAdapter) Close(db *gorm.DB) {
	if !CurrentConfig.InTests {
		db1, _ := db.DB()
		db1.Close()
	}
}

func (d *SqliteAdapter) ClearTestDatabase() {

}

func (d *SqliteAdapter) SetTimeZone(db *gorm.DB, timezone string) {
}

func (d *SqliteAdapter) InitializeDatabaseForTests(databaseSettings *DBSettings) {

}

func (d *SqliteAdapter) ArrayIncludes(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	d.LastError = errors.New("right now sqlite doesn't support filtering by array")
}

func (d *SqliteAdapter) JSONContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	d.LastError = errors.New("right now sqlite doesn't support searching in json")
}

func (d *SqliteAdapter) StartDBShell(databaseSettings *DBSettings) error {
	commandToExecute := exec.Command(
		"sqlite3", databaseSettings.Name,
	)
	// Sets standard output to cmd.stdout writer
	commandToExecute.Stdout = os.Stdout
	// Sets standard input to cmd.stdin reader
	commandToExecute.Stdin = os.Stdin
	var err error
	if err = commandToExecute.Start(); err != nil {
		return err
	}
	if err = commandToExecute.Wait(); err != nil {
		return err
	}
	return nil
}

func (d *SqliteAdapter) GetLastError() error {
	return d.LastError
}

func (d *SqliteAdapter) ResetLastError() {
	d.LastError = nil
}


func sqliteGoMonolithDatetimeParse(dt string, tzName string, connTzname string) *time.Time {
	if dt == "" {
		return nil
	}
	dt = strings.Replace(dt, "+00:00", "", 1)
	ret, _ := time.Parse("2006-01-_2 15:04:05", strings.Split(dt, ".")[0])
	loc, _ := time.LoadLocation(tzName)
	ret = ret.In(loc)
	return &ret
}

func sqliteGoMonolithDatetimeCastDate(dt string, tzName string, connTzname string) string {
	dtTmp := sqliteGoMonolithDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Round(0).Format(time.RFC3339)
	return res
}

func sqliteGoMonolithDatetimeCastTime(dt string, tzName string, connTzname string) string {
	dtTmp := sqliteGoMonolithDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return ""
	}
	res := dtTmp.Format("15:04")
	return res
}

func sqliteGoMonolithDatetimeExtract(extract string, dt string, tzName string, connTzname string) int {
	dtTmp := sqliteGoMonolithDatetimeParse(dt, tzName, connTzname)
	if dtTmp == nil {
		return 0
	}
	if extract == "month" {
		return int(dtTmp.Month())
	}
	if extract == "quarter" {
		return int(math.Ceil(float64(dtTmp.Month() / 3)))
	}
	if extract == "day" {
		return dtTmp.Day()
	}
	if extract == "hour" {
		return dtTmp.Hour()
	}
	if extract == "minute" {
		return dtTmp.Minute()
	}
	if extract == "second" {
		return dtTmp.Second()
	}
	if extract == "week" {
		_, isoWeek := dtTmp.ISOWeek()
		return isoWeek
	}
	if extract == "week_day" {
		return int(dtTmp.Weekday())
	}
	return 0
}

func sqliteGoMonolithRegex(reString string, rePattern string) bool {
	regex := regexp.MustCompile(rePattern)
	return regex.Find([]byte(reString)) != nil
}

//if err := conn.RegisterFunc("gomonolith_datetime_cast_year", sqlite_gomonolith_datetime_cast_year, true); err != nil {
//	return err
//}

func init() {
	sql.Register("GoMonolithSqliteDriver", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			var err error
			if err = conn.RegisterFunc("gomonolith_datetime_cast_date", sqliteGoMonolithDatetimeCastDate, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("gomonolith_datetime_extract", sqliteGoMonolithDatetimeExtract, true); err != nil {
				return err
			}
			if err = conn.RegisterFunc("gomonolith_datetime_cast_time", sqliteGoMonolithDatetimeCastTime, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("gomonolith_regex", sqliteGoMonolithRegex, true); err != nil {
				return err
			}
			for operator := range ProjectGormOperatorRegistry.GetAll() {
				err = operator.RegisterDbHandlers(conn)
				if err != nil {
					return err
				}
			}
			return nil
		},
	})
	InitializeGlobalAdapterRegistry()
	GlobalDbAdapterRegistry.RegisterAdapter("sqlite", func(db *gorm.DB) IDbAdapter {
		return &SqliteAdapter{
			DbType: "sqlite",
			Statement: &gorm.Statement{
				DB:      db,
				Context: context.Background(),
				Clauses: map[string]clause.Clause{},
			},
		}
	})
}
