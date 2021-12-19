// +build postgres

package core

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // import needed for postgres database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"time"
)

type PostgresAdapter struct {
	Statement *gorm.Statement
	DbType    string
	LastError error
}

func (d *PostgresAdapter) Equals(name interface{}, args ...interface{}) {
	query := d.Statement.Quote(name) + " = ?"
	clause.Expr{SQL: query, Vars: args}.Build(d.Statement)
}

func (d *PostgresAdapter) GetStringToExtractYearFromField(filterOptionField string) string {
	return fmt.Sprintf("EXTRACT(YEAR FROM %s AT TIME ZONE 'UTC')", filterOptionField)
}

func (d *PostgresAdapter) GetStringToExtractMonthFromField(filterOptionField string) string {
	return fmt.Sprintf("EXTRACT(MONTH FROM %s AT TIME ZONE 'UTC')", filterOptionField)
}

var cachedPostgresDB *gorm.DB

// @todo analyze
func (d *PostgresAdapter) GetDb(alias string, dryRun bool) (*gorm.DB, error) {
	var aliasDatabaseSettings *DBSettings
	if alias == "default" {
		aliasDatabaseSettings = CurrentDatabaseSettings.Default
	} else {
		aliasDatabaseSettings = CurrentDatabaseSettings.Slave
	}
	host := aliasDatabaseSettings.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := aliasDatabaseSettings.Port
	if port == 0 {
		port = 5432
	}
	user := aliasDatabaseSettings.User
	if user == "" {
		user = "root"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		host,
		user,
		aliasDatabaseSettings.Password,
		aliasDatabaseSettings.Name,
		port,
	)
	var db *gorm.DB
	var err error
	if CurrentConfig.InTests {
		if cachedPostgresDB != nil {
			return cachedPostgresDB, nil
		}
		if CurrentConfig.DebugTests || true {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		} else {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}
		cachedPostgresDB = db
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		db.Exec("SET TIME ZONE UTC")
	}
	return db, err
}

func (d *PostgresAdapter) Exact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IExact(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) = UPPER(?) ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Contains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? || '%%' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IContains(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE UPPER('%%' || ? || '%%')", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) In(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s IN ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Gt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s > ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Gte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s >= ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Lt(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s < ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Lte(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s <= ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) StartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE ? || '%%' ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IStartsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE UPPER(? || '%%') ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) EndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s LIKE '%%' || ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IEndsWith(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" UPPER(%s.%s) LIKE UPPER('%%' || ?) ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Date(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s::date = ?::date ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Year(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s BETWEEN ? AND ? ", operatorContext.TableName, field.DBName)
	year := value.(int)
	startOfTheYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfTheYear := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, startOfTheYear, endOfTheYear)
}

func (d *PostgresAdapter) Month(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('month' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Day(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('day' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Week(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('week' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) WeekDay(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('dow' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Quarter(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('quarter' FROM %s.%s AT TIME ZONE 'UTC') = ? + 1", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Hour(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('hour' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Minute(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" EXTRACT('minute' FROM %s.%s AT TIME ZONE 'UTC') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Second(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" FLOOR(EXTRACT('second' FROM %s.%s AT TIME ZONE 'UTC')) = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Regex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s ~ ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IRegex(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" %s.%s ~* ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) Time(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	query := fmt.Sprintf(" to_char((%s.%s AT TIME ZONE 'UTC')::time, 'HH24:MI') = ? ", operatorContext.TableName, field.DBName)
	args := value
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query, args)
}

func (d *PostgresAdapter) IsNull(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
	isTruthyValue := IsTruthyValue(value)
	isNull := " IS NULL "
	if !isTruthyValue {
		isNull = " IS NOT NULL "
	}
	query := fmt.Sprintf(" %s.%s %s ", operatorContext.TableName, field.DBName, isNull)
	operatorContext.Tx = SQLConditionBuilder.Build(operatorContext.Tx, query)
}

func (d *PostgresAdapter) Range(operatorContext *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) {
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

func (d *PostgresAdapter) BuildDeleteString(table string, cond string, values ...interface{}) *DeleteRowStructure {
	deleteRowStructure := &DeleteRowStructure{SQL: fmt.Sprintf("DELETE FROM %s WHERE %s", table, cond), Values: values}
	return deleteRowStructure
}

func (d *PostgresAdapter) SetIsolationLevelForTests(db *gorm.DB) {
	db.Exec("set transaction isolation level repeatable read")
}

func (d *PostgresAdapter) Close(db *gorm.DB) {
	if !CurrentConfig.InTests {
		db1, _ := db.DB()
		db1.Close()
	}
}

func (d *PostgresAdapter) ClearTestDatabase() {
	cachedPostgresDB = nil
}

func (d *PostgresAdapter) SetTimeZone(db *gorm.DB, timezone string) {
	db.Exec("SET TIME ZONE UTC")
}

func (d *PostgresAdapter) InitializeDatabaseForTests(databaseSettings *DBSettings) {
	aliasDatabaseSettings := databaseSettings
	host := aliasDatabaseSettings.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := aliasDatabaseSettings.Port
	if port == 0 {
		port = 5432
	}
	user := aliasDatabaseSettings.User
	if user == "" {
		user = "root"
	}
	dsnToCreateDatabase := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=postgres sslmode=disable TimeZone=UTC",
		host,
		user,
		aliasDatabaseSettings.Password,
		port,
	)
	db1, err1 := sql.Open(databaseSettings.Type, dsnToCreateDatabase)
	if err1 != nil {
		panic(err1)
	}
	db1.Exec("create database " + aliasDatabaseSettings.Name)
	dsnToCreateDatabase = fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable TimeZone=UTC",
		host,
		user,
		aliasDatabaseSettings.Password,
		port,
		aliasDatabaseSettings.Name,
	)
	db1, err1 = sql.Open(databaseSettings.Type, dsnToCreateDatabase)
	if err1 != nil {
		panic(err1)
	}
	_, err4 := db1.Exec("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	if err4 != nil {
		panic(err4)
	}

}

func (d *PostgresAdapter) StartDBShell(databaseSettings *DBSettings) error {
	host := databaseSettings.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := databaseSettings.Port
	if port == 0 {
		port = 5432
	}
	commandToExecute := exec.Command(
		"psql", "-U", databaseSettings.User,
		"-h", host, "-p", strconv.Itoa(port), databaseSettings.Name,
	)
	commandToExecute.Env = os.Environ()
	commandToExecute.Env = append(commandToExecute.Env, "PGPASSWORD="+databaseSettings.Password)
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

func (d *PostgresAdapter) GetLastError() error {
	return d.LastError
}

func init() {
	InitializeGlobalAdapterRegistry()
	GlobalDbAdapterRegistry.RegisterAdapter("postgres", func(db *gorm.DB) IDbAdapter {
		return &PostgresAdapter{
			DbType: "postgres",
			Statement: &gorm.Statement{
				DB:      db,
				Context: context.Background(),
				Clauses: map[string]clause.Clause{},
			},
		}
	})
}
