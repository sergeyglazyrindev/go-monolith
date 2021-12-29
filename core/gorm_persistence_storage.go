package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func NewGormAdminPage(parentPage *AdminPage, modelI interface{}, generateForm func(modelI interface{}, ctx IAdminContext) *Form) *AdminPage {
	modelName := ""
	var form *Form
	var listDisplay *ListDisplayRegistry
	var searchFieldRegistry *SearchFieldRegistry
	generateModelI := func() (interface{}, interface{}) { return nil, nil }
	if modelI != nil {
		modelDesc := ProjectModels.GetModelFromInterface(modelI)
		generateModelI = modelDesc.GenerateModelI
		modelI4, _ := generateModelI()
		if modelI4 != nil {
			database := NewDatabaseInstanceWithoutConnection()
			stmt := &gorm.Statement{DB: database.Db}
			stmt.Parse(modelI4)
			modelName = strings.ToLower(stmt.Schema.Name)
		}
		form = NewFormFromModelFromGinContext(&AdminContext{}, modelI4, make([]string, 0), []string{"ID"}, true, "")
		listDisplay = NewListDisplayRegistryFromGormModel(modelI4)
		searchFieldRegistry = NewSearchFieldRegistryFromGormModel(modelI4)
	}
	return &AdminPage{
		Form:       form,
		SubPages:   NewAdminPageRegistry(),
		ParentPage: parentPage,
		GetQueryset: func(adminContext IAdminContext, adminPage *AdminPage, adminRequestParams *AdminRequestParams) IAdminFilterObjects {
			database := NewDatabaseInstance()
			db := database.Db
			var paginatedQuerySet IPersistenceStorage
			var perPage int
			modelI, _ := generateModelI()
			modelI1, _ := generateModelI()
			modelI2, _ := generateModelI()
			modelI3, _ := generateModelI()
			ret := &GormAdminFilterObjects{
				InitialGormQuerySet:   NewGormPersistenceStorage(db.Model(modelI)),
				GormQuerySet:          NewGormPersistenceStorage(db.Model(modelI1)),
				PaginatedGormQuerySet: NewGormPersistenceStorage(db.Model(modelI2)),
				Model:                 modelI3,
				DatabaseInstance:      database,
			}
			if adminPage.EnhanceQuerySet != nil {
				adminPage.EnhanceQuerySet(ret)
			}
			if adminRequestParams != nil && adminRequestParams.RequestURL != "" {
				url1, _ := url.Parse(adminRequestParams.RequestURL)
				queryParams, _ := url.ParseQuery(url1.RawQuery)
				for filter := range adminPage.ListFilter.Iterate() {
					filterValue := queryParams.Get(filter.URLFilteringParam)
					if filterValue != "" {
						filter.FilterQs(ret, fmt.Sprintf("%s=%s", filter.URLFilteringParam, filterValue))
					}
				}
			}
			if adminRequestParams != nil && adminRequestParams.Search != "" {
				searchFilterObjects := &GormAdminFilterObjects{
					InitialGormQuerySet:   NewGormPersistenceStorage(db),
					GormQuerySet:          NewGormPersistenceStorage(db),
					PaginatedGormQuerySet: NewGormPersistenceStorage(db),
					Model:                 modelI3,
					DatabaseInstance:      database,
				}
				for filter := range adminPage.SearchFields.GetAll() {
					filter.Search(searchFilterObjects, adminRequestParams.Search)
				}
				ret.SetPaginatedQuerySet(ret.GetPaginatedQuerySet().Where(searchFilterObjects.GetPaginatedQuerySet().GetCurrentDB()))
				ret.SetFullQuerySet(ret.GetFullQuerySet().Where(searchFilterObjects.GetFullQuerySet().GetCurrentDB()))
				searchFilterObjects.AddNeededJoinsIfNecessary(ret)
			}
			if adminRequestParams != nil && adminRequestParams.Paginator.PerPage > 0 {
				perPage = adminRequestParams.Paginator.PerPage
			} else {
				perPage = adminPage.Paginator.PerPage
			}
			if adminRequestParams != nil {
				paginatedQuerySet = ret.GetPaginatedQuerySet().Offset(adminRequestParams.Paginator.Offset)
				if adminPage.Paginator.ShowLastPageOnPreviousPage {
					var countRecords int64
					ret.GetFullQuerySet().Count(&countRecords)
					if countRecords > int64(adminRequestParams.Paginator.Offset+(2*perPage)) {
						paginatedQuerySet = paginatedQuerySet.Limit(perPage)
					} else {
						paginatedQuerySet = paginatedQuerySet.Limit(int(countRecords - int64(adminRequestParams.Paginator.Offset)))
					}
				} else {
					paginatedQuerySet = paginatedQuerySet.Limit(perPage)
				}
				ret.SetPaginatedQuerySet(paginatedQuerySet)
				for listDisplay := range adminPage.ListDisplay.GetAllFields() {
					direction := listDisplay.SortBy.GetDirection()
					if len(adminRequestParams.Ordering) > 0 {
						for _, ordering := range adminRequestParams.Ordering {
							directionSort := 1
							if strings.HasPrefix(ordering, "-") {
								directionSort = -1
								ordering = ordering[1:]
							}
							if ordering == listDisplay.DisplayName {
								direction = directionSort
								listDisplay.SortBy.Sort(ret, direction)
							}
						}
					}
				}
			}
			if adminPage.CustomizeQuerySet != nil {
				adminPage.CustomizeQuerySet(adminContext, ret, adminRequestParams)
			}
			return ret
		},
		Model:                   modelI,
		ModelName:               modelName,
		Validators:              NewValidatorRegistry(),
		ExcludeFields:           NewFieldRegistry(),
		FieldsToShow:            NewFieldRegistry(),
		ModelActionsRegistry:    NewAdminModelActionRegistry(),
		InlineRegistry:          NewAdminPageInlineRegistry(),
		ListDisplay:             listDisplay,
		ListFilter:              &ListFilterRegistry{ListFilter: make([]*ListFilter, 0)},
		SearchFields:            searchFieldRegistry,
		Paginator:               &Paginator{PerPage: CurrentConfig.D.GoMonolith.AdminPerPage, ShowLastPageOnPreviousPage: true},
		ActionsSelectionCounter: true,
		FilterOptions:           NewFilterOptionsRegistry(),
		GenerateForm:            generateForm,
	}
}

type GormPersistenceStorage struct {
	Db        *gorm.DB
	LastError error
}

func NewGormPersistenceStorage(db *gorm.DB) *GormPersistenceStorage {
	return &GormPersistenceStorage{Db: db}
}

func (gps *GormPersistenceStorage) Association(column string) IPersistenceAssociation {
	ret := gps.Db.Association(column)
	if ret.Error != nil {
		gps.LastError = ret.Error
	}
	return ret
}

func (gps *GormPersistenceStorage) Model(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Model(value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Clauses(conds ...clause.Expression) IPersistenceStorage {
	gps.Db = gps.Db.Clauses(conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) GetCurrentDB() *gorm.DB {
	return gps.Db
}

func (gps *GormPersistenceStorage) Table(name string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Table(name, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Distinct(args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Distinct(args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Select(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Select(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Omit(columns ...string) IPersistenceStorage {
	gps.Db = gps.Db.Omit(columns...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Where(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Where(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Not(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Not(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Or(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Or(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Joins(query string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Joins(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Group(name string) IPersistenceStorage {
	gps.Db = gps.Db.Group(name)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Having(query interface{}, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Having(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Order(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Order(value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Limit(limit int) IPersistenceStorage {
	gps.Db = gps.Db.Limit(limit)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Offset(offset int) IPersistenceStorage {
	gps.Db = gps.Db.Offset(offset)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Scopes(funcs ...func(*gorm.DB) *gorm.DB) IPersistenceStorage {
	gps.Db = gps.Db.Scopes(funcs...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Preload(query string, args ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Preload(query, args...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Attrs(attrs ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Attrs(attrs...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Assign(attrs ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Assign(attrs...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Unscoped() IPersistenceStorage {
	gps.Db = gps.Db.Unscoped()
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Raw(sql string, values ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Raw(sql, values...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Migrator() IPersistenceMigrator {
	return gps.Db.Migrator()
}

func (gps *GormPersistenceStorage) AutoMigrate(dst ...interface{}) error {
	return gps.Db.AutoMigrate(dst...)
}

func (gps *GormPersistenceStorage) Session(config *gorm.Session) IPersistenceStorage {
	gps.Db = gps.Db.Session(config)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) WithContext(ctx context.Context) IPersistenceStorage {
	gps.Db = gps.Db.WithContext(ctx)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Debug() IPersistenceStorage {
	gps.Db = gps.Db.Debug()
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Set(key string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Set(key, value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Get(key string) (interface{}, bool) {
	return gps.Db.Get(key)
}

func (gps *GormPersistenceStorage) InstanceSet(key string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.InstanceSet(key, value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) InstanceGet(key string) (interface{}, bool) {
	return gps.Db.InstanceGet(key)
}

func (gps *GormPersistenceStorage) AddError(err error) error {
	return gps.Db.AddError(err)
}

func (gps *GormPersistenceStorage) DB() (*sql.DB, error) {
	return gps.Db.DB()
}

func (gps *GormPersistenceStorage) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return gps.Db.SetupJoinTable(model, field, joinTable)
}

func (gps *GormPersistenceStorage) Use(plugin gorm.Plugin) error {
	return gps.Db.Use(plugin)
}

func (gps *GormPersistenceStorage) Create(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Create(value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) CreateInBatches(value interface{}, batchSize int) IPersistenceStorage {
	gps.Db = gps.Db.CreateInBatches(value, batchSize)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Save(value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Save(value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) First(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.First(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Take(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Take(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Last(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Last(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Find(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Find(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) IPersistenceStorage {
	gps.Db = gps.Db.FindInBatches(dest, batchSize, fc)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) FirstOrInit(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.FirstOrInit(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) FirstOrCreate(dest interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.FirstOrCreate(dest, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Update(column string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Update(column, value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Updates(values interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Updates(values)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) UpdateColumn(column string, value interface{}) IPersistenceStorage {
	gps.Db = gps.Db.UpdateColumn(column, value)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) UpdateColumns(values interface{}) IPersistenceStorage {
	gps.Db = gps.Db.UpdateColumns(values)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Delete(value interface{}, conds ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Delete(value, conds...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Count(count *int64) IPersistenceStorage {
	gps.Db = gps.Db.Count(count)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Row() IPersistenceIterateRow {
	return gps.Db.Row()
}

func (gps *GormPersistenceStorage) Rows() (IPersistenceIterateRows, error) {
	return gps.Db.Rows()
}

func (gps *GormPersistenceStorage) Scan(dest interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Scan(dest)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Pluck(column string, dest interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Pluck(column, dest)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) ScanRows(rows IPersistenceIterateRows, dest interface{}) error {
	return gps.Db.ScanRows(rows.(*sql.Rows), dest)
}

func (gps *GormPersistenceStorage) Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return gps.Db.Transaction(fc, opts...)
}

func (gps *GormPersistenceStorage) Begin(opts ...*sql.TxOptions) IPersistenceStorage {
	gps.Db = gps.Db.Begin(opts...)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Commit() IPersistenceStorage {
	gps.Db = gps.Db.Commit()
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Rollback() IPersistenceStorage {
	gps.Db = gps.Db.Rollback()
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) SavePoint(name string) IPersistenceStorage {
	gps.Db = gps.Db.SavePoint(name)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) RollbackTo(name string) IPersistenceStorage {
	gps.Db = gps.Db.RollbackTo(name)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) Exec(sql string, values ...interface{}) IPersistenceStorage {
	gps.Db = gps.Db.Exec(sql, values)
	if gps.Db.Error != nil {
		gps.LastError = gps.Db.Error
	}
	return gps
}

func (gps *GormPersistenceStorage) ResetLastError() {
	gps.LastError = nil
}

func (gps *GormPersistenceStorage) GetLastError() error {
	ret := gps.LastError
	return ret
}

func (gps *GormPersistenceStorage) LoadDataForModelByID(modelI interface{}, ID string) IPersistenceStorage {
	modelDescription := ProjectModels.GetModelFromInterface(modelI)
	gps.Db.Where(fmt.Sprintf("\"%s\" = ?", modelDescription.Statement.Schema.PrimaryFields[0].DBName), ID).Preload(clause.Associations).First(modelI)
	return gps
}

type GormAdminFilterObjects struct {
	InitialGormQuerySet   IPersistenceStorage
	GormQuerySet          IPersistenceStorage
	PaginatedGormQuerySet IPersistenceStorage
	Model                 interface{}
	DatabaseInstance      *ProjectDatabase
	LastError             error
	NeededJoins           []string
}

func (afo *GormAdminFilterObjects) SetLastError(err error) {
	if err != nil {
		afo.LastError = err
	}
}

func (afo *GormAdminFilterObjects) AddNeededJoinsIfNecessary(afo1 IAdminFilterObjects) {
	if afo.NeededJoins == nil {
		return
	}
	if len(afo.NeededJoins) == 0 {
		return
	}
	for _, join := range afo.NeededJoins {
		afo1.SetFullQuerySet(afo1.GetFullQuerySet().Joins(join))
		afo1.SetPaginatedQuerySet(afo1.GetPaginatedQuerySet().Joins(join))
	}
}

func (afo *GormAdminFilterObjects) StoreNeededJoin(join string) {
	if afo.NeededJoins == nil {
		afo.NeededJoins = make([]string, 0)
	}
	afo.NeededJoins = append(afo.NeededJoins, join)
}

func (afo *GormAdminFilterObjects) ResetLastError() {
	afo.LastError = nil
}

func (afo *GormAdminFilterObjects) GetLastError() error {
	ret := afo.LastError
	return ret
}

func (afo *GormAdminFilterObjects) FilterQs(filterString string) {
	statement := &gorm.Statement{DB: afo.GetDatabase().Db}
	statement.Parse(afo.GetCurrentModel())
	schema1 := statement.Schema
	operatorContext := FilterGormModel(afo.GetDatabase().Adapter, afo.GetFullQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	afo.SetFullQuerySet(operatorContext.Tx)
	operatorContext = FilterGormModel(afo.GetDatabase().Adapter, afo.GetPaginatedQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	afo.SetPaginatedQuerySet(operatorContext.Tx)
	afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
}

func (afo *GormAdminFilterObjects) Search(field *Field, searchString string) {
	fieldType := field.FieldType.Kind()
	if fieldType == reflect.Ptr {
		fieldType = field.FieldType.Elem().Kind()
	}
	if fieldType == reflect.Struct {
		mReflectValue := reflect.ValueOf(field.FieldConfig.Widget.GetValue())
		mInterface := mReflectValue.Interface()
		adminPage := CurrentDashboardAdminPanel.FindPageForGormModel(mInterface)
		joinModelI := ProjectModels.GetModelFromInterface(mInterface)
		model, _ := joinModelI.GenerateModelI()
		model1, _ := joinModelI.GenerateModelI()
		relation := field.Schema.Relationships
		relationsString := []string{}
		for _, relation1 := range relation.Relations {
			for _, reference := range relation1.References {
				if joinModelI.Statement.Schema.Table != reference.PrimaryKey.Schema.Table &&
					joinModelI.Statement.Schema.Table != reference.ForeignKey.Schema.Table {
					continue
				}
				relationsString = append(
					relationsString,
					fmt.Sprintf(
						"%s.%s = %s.%s",
						reference.PrimaryKey.Schema.Table, reference.PrimaryKey.DBName, reference.ForeignKey.Schema.Table,
						reference.ForeignKey.DBName,
					),
				)
			}
		}
		if field.NotNull {
			afo.StoreNeededJoin(
				fmt.Sprintf(
					"INNER JOIN %s on %s",
					joinModelI.Statement.Table, strings.Join(relationsString, " AND "),
				),
			)
		} else {
			afo.StoreNeededJoin(
				fmt.Sprintf(
					"LEFT JOIN %s on %s",
					joinModelI.Statement.Table, strings.Join(relationsString, " AND "),
				),
			)
		}
		fullGormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), model)
		paginatedGormOperatorContext := NewGormOperatorContext(afo.GetPaginatedQuerySet(), model1)
		for searchField := range adminPage.SearchFields.GetAll() {
			fieldType1 := searchField.GetField().FieldType.Kind()
			if searchField.GetField().FieldType.Kind() == reflect.Struct {
				afo.Search(searchField.GetField(), searchString)
				continue
			} else if (fieldType1 == reflect.Uint) || (fieldType1 == reflect.Uint8) || (fieldType1 == reflect.Uint16) || (fieldType1 == reflect.Uint64) || (fieldType1 == reflect.Uint32) || (fieldType1 == reflect.Int64) || (fieldType1 == reflect.Int) || (fieldType1 == reflect.Int32) || (fieldType1 == reflect.Int8)  || (fieldType1 == reflect.Int16) {
				searchID, err1 := strconv.Atoi(searchString)
				if err1 != nil {
					continue
				}
				operator := ExactGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, fullGormOperatorContext, searchField.GetField(), searchID, &SQLConditionBuilder{Type: "or"})
				operator = ExactGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, paginatedGormOperatorContext, searchField.GetField(), searchID, &SQLConditionBuilder{Type: "or"})
			} else if searchField.GetField().FieldType.Name() == "JSON" {
				operator := JSONIncludesGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, fullGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
				operator = JSONIncludesGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, paginatedGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
			} else if strings.HasSuffix(searchField.GetField().FieldType.Name(), "Array") {
				operator := ArrayIncludesGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, fullGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
				operator = ArrayIncludesGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, paginatedGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
			} else {
				operator := IContainsGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, fullGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
				operator = IContainsGormOperator{}
				operator.Build(afo.GetDatabase().Adapter, paginatedGormOperatorContext, searchField.GetField(), searchString, &SQLConditionBuilder{Type: "or"})
			}
		}
		afo.SetFullQuerySet(fullGormOperatorContext.Tx)
		afo.SetPaginatedQuerySet(fullGormOperatorContext.Tx)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
	} else if (fieldType == reflect.Uint) || (fieldType == reflect.Uint8) || (fieldType == reflect.Uint16) || (fieldType == reflect.Uint64) || (fieldType == reflect.Uint32) || (fieldType == reflect.Int64) || (fieldType == reflect.Int) || (fieldType == reflect.Int32) || (fieldType == reflect.Int8) || (fieldType == reflect.Int16) {
		searchID, err1 := strconv.Atoi(searchString)
		if err1 != nil {
			return
		}
		operator := ExactGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchID, &SQLConditionBuilder{Type: "or"})
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchID, &SQLConditionBuilder{Type: "or"})
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
	} else if field.FieldType.Name() == "JSON" {
		operator := JSONIncludesGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
	} else if strings.HasSuffix(field.FieldType.Name(), "Array") {
		operator := ArrayIncludesGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
	} else {
		operator := IContainsGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, field, searchString, &SQLConditionBuilder{Type: "or"})
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
	}
}

func (afo *GormAdminFilterObjects) GetPaginatedQuerySet() IPersistenceStorage {
	return afo.PaginatedGormQuerySet
}

func (afo *GormAdminFilterObjects) GetFullQuerySet() IPersistenceStorage {
	return afo.GormQuerySet
}

func (afo *GormAdminFilterObjects) SetFullQuerySet(storage IPersistenceStorage) {
	afo.GormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GetInitialQuerySet() IPersistenceStorage {
	return afo.InitialGormQuerySet
}

func (afo *GormAdminFilterObjects) SetInitialQuerySet(storage IPersistenceStorage) {
	afo.InitialGormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GetCurrentModel() interface{} {
	return afo.Model
}

func (afo *GormAdminFilterObjects) GetDatabase() *ProjectDatabase {
	return afo.DatabaseInstance
}

func (afo *GormAdminFilterObjects) SetPaginatedQuerySet(storage IPersistenceStorage) {
	afo.PaginatedGormQuerySet = storage
}

func (afo *GormAdminFilterObjects) GetDB() IPersistenceStorage {
	return NewGormPersistenceStorage(afo.DatabaseInstance.Db)
}

func (afo *GormAdminFilterObjects) WithTransaction(handler func(afo1 IAdminFilterObjects) error) error {
	afo.DatabaseInstance.Db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		ret := handler(&GormAdminFilterObjects{
			DatabaseInstance:    &ProjectDatabase{Db: tx, Adapter: afo.DatabaseInstance.Adapter},
			InitialGormQuerySet: NewGormPersistenceStorage(tx),
		})
		afo.SetLastError(ret)
		return ret
	})
	return afo.GetLastError()
}

func (afo *GormAdminFilterObjects) LoadDataForModelByID(ID interface{}, model interface{}) {
	modelDescription := ProjectModels.GetModelFromInterface(model)
	afo.DatabaseInstance.Db.Preload(clause.Associations).Where(fmt.Sprintf(" \"%s\" = ?", modelDescription.Statement.Schema.PrimaryFields[0].DBName), ID).First(model)
	afo.SetLastError(afo.DatabaseInstance.Db.Error)
}

func (afo *GormAdminFilterObjects) SaveModel(model interface{}) error {
	res := afo.DatabaseInstance.Db.Save(model)
	afo.SetLastError(res.Error)
	return res.Error
}

func (afo *GormAdminFilterObjects) CreateNew(model interface{}) error {
	res := afo.DatabaseInstance.Db.Model(model).Create(model)
	afo.SetLastError(res.Error)
	return res.Error
}

func (afo *GormAdminFilterObjects) FilterByMultipleIds(field *Field, realObjectIds []string) {
	afo.SetFullQuerySet(afo.GetFullQuerySet().Where(fmt.Sprintf("%s IN ?", field.DBName), realObjectIds))
	afo.SetLastError(afo.GormQuerySet.GetLastError())
}

func (afo *GormAdminFilterObjects) RemoveModelPermanently(model interface{}) error {
	res := afo.DatabaseInstance.Db.Unscoped().Delete(model)
	afo.SetLastError(res.Error)
	return res.Error
}

func (afo *GormAdminFilterObjects) SortBy(field *Field, direction int) {
	sortBy := field.DBName
	if direction == -1 {
		sortBy += " desc"
	}
	afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Order(sortBy))
	afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
}

func (afo *GormAdminFilterObjects) GetPaginated() <-chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				afo.SetLastError(errors.New("issue while iterating objects"))
				Trail(CRITICAL, "Recovering from panic in GetPaginated error is: %v \n", r)
			}
			close(chnl)
		}()
		modelDesc := ProjectModels.GetModelFromInterface(afo.Model)
		modelI, models := modelDesc.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.PaginatedGormQuerySet.Preload(clause.Associations).Find(models)
		afo.SetLastError(afo.PaginatedGormQuerySet.GetLastError())
		if afo.PaginatedGormQuerySet.GetLastError() != nil {
			panic(afo.PaginatedGormQuerySet.GetLastError())
		}
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			ID := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            ID.(string),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

func (afo *GormAdminFilterObjects) IterateThroughWholeQuerySet() <-chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				afo.SetLastError(errors.New("issue while iterating objects"))
				Trail(CRITICAL, "Recovering from panic in IterateThroughWholeQuerySet error is: %v \n", r)
			}
			close(chnl)
		}()
		modelDesc := ProjectModels.GetModelFromInterface(afo.Model)
		modelI, models := modelDesc.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.GormQuerySet.Preload(clause.Associations).Find(models)
		afo.SetLastError(afo.GormQuerySet.GetLastError())
		if afo.GormQuerySet.GetLastError() != nil {
			panic(afo.GormQuerySet.GetLastError())
		}
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			ID := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
			yieldV := &IterateAdminObjects{
				Model:         model,
				ID:            ID.(string),
				RenderContext: &FormRenderContext{Model: model},
			}
			chnl <- yieldV
		}
	}()
	return chnl
}
