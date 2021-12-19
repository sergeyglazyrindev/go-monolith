package core

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type IPersistenceMigrator interface {
	AutoMigrate(dst ...interface{}) error
	CurrentDatabase() string
	FullDataTypeOf(*schema.Field) clause.Expr
	CreateTable(dst ...interface{}) error
	DropTable(dst ...interface{}) error
	HasTable(dst interface{}) bool
	RenameTable(oldName interface{}, newName interface{}) error
	AddColumn(dst interface{}, field string) error
	DropColumn(dst interface{}, field string) error
	AlterColumn(dst interface{}, field string) error
	MigrateColumn(dst interface{}, field *schema.Field, columnType gorm.ColumnType) error
	HasColumn(dst interface{}, field string) bool
	RenameColumn(dst interface{}, oldName string, field string) error
	ColumnTypes(dst interface{}) ([]gorm.ColumnType, error)
	CreateView(name string, option gorm.ViewOption) error
	DropView(name string) error
	CreateConstraint(dst interface{}, name string) error
	DropConstraint(dst interface{}, name string) error
	HasConstraint(dst interface{}, name string) bool
	CreateIndex(dst interface{}, name string) error
	DropIndex(dst interface{}, name string) error
	HasIndex(dst interface{}, name string) bool
	RenameIndex(dst interface{}, oldName string, newName string) error
}

type IPersistenceAssociation interface {
	Find(out interface{}, conds ...interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() (count int64)
}

type IPersistenceIterateRow interface {
	Scan(dest ...interface{}) error
	Err() error
}

type IPersistenceIterateRows interface {
	Next() bool
	NextResultSet() bool
	Err() error
	Columns() ([]string, error)
	ColumnTypes() ([]*sql.ColumnType, error)
	Scan(dest ...interface{}) error
	Close() error
}

type IPersistenceStorage interface {
	Association(column string) IPersistenceAssociation
	Model(value interface{}) IPersistenceStorage
	Clauses(conds ...clause.Expression) IPersistenceStorage
	Table(name string, args ...interface{}) IPersistenceStorage
	Distinct(args ...interface{}) IPersistenceStorage
	Select(query interface{}, args ...interface{}) IPersistenceStorage
	Omit(columns ...string) IPersistenceStorage
	Where(query interface{}, args ...interface{}) IPersistenceStorage
	Not(query interface{}, args ...interface{}) IPersistenceStorage
	Or(query interface{}, args ...interface{}) IPersistenceStorage
	Joins(query string, args ...interface{}) IPersistenceStorage
	Group(name string) IPersistenceStorage
	Having(query interface{}, args ...interface{}) IPersistenceStorage
	Order(value interface{}) IPersistenceStorage
	Limit(limit int) IPersistenceStorage
	Offset(offset int) IPersistenceStorage
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) IPersistenceStorage
	Preload(query string, args ...interface{}) IPersistenceStorage
	Attrs(attrs ...interface{}) IPersistenceStorage
	Assign(attrs ...interface{}) IPersistenceStorage
	Unscoped() IPersistenceStorage
	Raw(sql string, values ...interface{}) IPersistenceStorage
	Migrator() IPersistenceMigrator
	AutoMigrate(dst ...interface{}) error
	Session(config *gorm.Session) IPersistenceStorage
	WithContext(ctx context.Context) IPersistenceStorage
	Debug() IPersistenceStorage
	Set(key string, value interface{}) IPersistenceStorage
	Get(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) IPersistenceStorage
	InstanceGet(key string) (interface{}, bool)
	AddError(err error) error
	DB() (*sql.DB, error)
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Use(plugin gorm.Plugin) error
	Create(value interface{}) IPersistenceStorage
	CreateInBatches(value interface{}, batchSize int) IPersistenceStorage
	Save(value interface{}) IPersistenceStorage
	First(dest interface{}, conds ...interface{}) IPersistenceStorage
	Take(dest interface{}, conds ...interface{}) IPersistenceStorage
	Last(dest interface{}, conds ...interface{}) IPersistenceStorage
	Find(dest interface{}, conds ...interface{}) IPersistenceStorage
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) IPersistenceStorage
	FirstOrInit(dest interface{}, conds ...interface{}) IPersistenceStorage
	FirstOrCreate(dest interface{}, conds ...interface{}) IPersistenceStorage
	Update(column string, value interface{}) IPersistenceStorage
	Updates(values interface{}) IPersistenceStorage
	UpdateColumn(column string, value interface{}) IPersistenceStorage
	UpdateColumns(values interface{}) IPersistenceStorage
	Delete(value interface{}, conds ...interface{}) IPersistenceStorage
	Count(count *int64) IPersistenceStorage
	Row() IPersistenceIterateRow
	Rows() (IPersistenceIterateRows, error)
	Scan(dest interface{}) IPersistenceStorage
	Pluck(column string, dest interface{}) IPersistenceStorage
	ScanRows(rows IPersistenceIterateRows, dest interface{}) error
	Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) IPersistenceStorage
	Commit() IPersistenceStorage
	Rollback() IPersistenceStorage
	SavePoint(name string) IPersistenceStorage
	RollbackTo(name string) IPersistenceStorage
	Exec(sql string, values ...interface{}) IPersistenceStorage
	GetCurrentDB() *gorm.DB
	GetLastError() error
	LoadDataForModelByID(modelI interface{}, ID string) IPersistenceStorage
}
