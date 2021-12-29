// +build elasticsearch

package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type ElasticModelInterface interface {
	SetID(ID string)
	GetID() string
	GetIndexName() string
	String() string
}

func NewElasticSearchAdminPage(parentPage *AdminPage, modelI interface{}, generateForm func(modelI interface{}, ctx IAdminContext) *Form) *AdminPage {
	modelName := ""
	var form *Form
	var listDisplay *ListDisplayRegistry
	var searchFieldRegistry *SearchFieldRegistry
	genModelI := func() (interface{}, interface{}) { return nil, nil }
	if modelI != nil {
		modelDesc := ProjectModels.GetModelFromInterface(modelI)
		genModelI = modelDesc.GenerateModelI
		modelI4, _ := genModelI()
		if modelI4 != nil {
			modelName = reflect.TypeOf(modelI4).Name()
		}
		if modelI4 != nil {
			form = NewFormFromModelFromGinContext(&AdminContext{}, modelI4, make([]string, 0), []string{"ID"}, true, "")
			listDisplay = NewListDisplayRegistryFromGormModel(modelI4)
			searchFieldRegistry = NewSearchFieldRegistryFromGormModel(modelI4)
		}

	}
	return &AdminPage{
		Form:       form,
		SubPages:   NewAdminPageRegistry(),
		ParentPage: parentPage,
		GetQueryset: func(adminContext IAdminContext, adminPage *AdminPage, adminRequestParams *AdminRequestParams) IAdminFilterObjects {
			database := NewDatabaseInstance()
			var paginatedQuerySet IPersistenceStorage
			var perPage int
			client := NewProjectESClient()
			modelI3, _ := genModelI()
			indexName := ""
			if v1, ok := modelI3.(ElasticModelInterface); ok {
				indexName = v1.GetIndexName()
			}

			initialQuerySet := NewElasticSearchPersistenceStorage(client, indexName, genModelI)
			esQuerySet := NewElasticSearchPersistenceStorage(client, indexName, genModelI)
			paginatedQuerySet1 := NewElasticSearchPersistenceStorage(client, indexName, genModelI)
			ret := &ElasticSearchAdminFilterObjects{
				InitialESQuerySet:   initialQuerySet,
				ESQuerySet:          esQuerySet,
				PaginatedESQuerySet: paginatedQuerySet1,
				Model:               modelI3,
				DatabaseInstance:    database,
				SearchBy:            make([]*ESSearchParam, 0),
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
				for filter := range adminPage.SearchFields.GetAll() {
					filter.Search(ret, adminRequestParams.Search)
				}
				ret.SearchString = adminRequestParams.Search
				ret.SearchInIndex()
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
			return ret
		},
		Model:                   modelI,
		ModelName:               modelName,
		Validators:              NewValidatorRegistry(),
		ExcludeFields:           NewFieldRegistry(),
		FieldsToShow:            NewFieldRegistry(),
		ModelActionsRegistry:    NewAdminModelActionRegistryForElasticSearch(),
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

type ElasticSearchPersistenceStorage struct {
	ESClient  *elastic.Client
	ESSearch  *elastic.SearchService
	IndexName string
	GenModelI func() (interface{}, interface{})
	LastError error
}

func NewElasticSearchPersistenceStorage(elasticClient *elastic.Client, indexName string, modelI interface{}) *ElasticSearchPersistenceStorage {
	modelDesc := ProjectModels.GetModelFromInterface(modelI)
	return &ElasticSearchPersistenceStorage{ESClient: elasticClient, GenModelI: modelDesc.GenerateModelI, IndexName: indexName}
}

func (gps *ElasticSearchPersistenceStorage) Association(column string) IPersistenceAssociation {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Model(value interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Clauses(conds ...clause.Expression) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) GetCurrentDB() *gorm.DB {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Table(name string, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Distinct(args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Select(query interface{}, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Omit(columns ...string) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Where(query interface{}, args ...interface{}) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	if v1, ok := query.(elastic.Query); ok {
		gps.ESSearch = gps.ESSearch.Query(v1)
	} else if len(args) > 1 {
		gps.ESSearch = gps.ESSearch.Query(elastic.NewTermsQuery(strings.ToLower(query.(string)), args...))
	} else {
		gps.ESSearch = gps.ESSearch.Query(elastic.NewTermQuery(strings.ToLower(query.(string)), args[0]))
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) SavePoint(name string) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Not(query interface{}, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Or(query interface{}, args ...interface{}) IPersistenceStorage {
	return gps.Where(query, args)
}

func (gps *ElasticSearchPersistenceStorage) Joins(query string, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Group(name string) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Having(query interface{}, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Order(value interface{}) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	value1 := value.(*ESSortBy)
	sortBy := ""
	if value1.FieldName != "" {
		sortBy = value1.FieldName
	} else {
		sortBy = value1.Field.DBName
	}
	if value1.SortByKeyword {
		sortBy += ".keyword"
	}
	gps.ESSearch = gps.ESSearch.Sort(sortBy, value1.Direction)
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Limit(limit int) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	gps.ESSearch = gps.ESSearch.Size(limit)
	return gps
}

func (gps *ElasticSearchPersistenceStorage) Offset(offset int) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	gps.ESSearch = gps.ESSearch.From(offset)
	return gps
}

func (gps *ElasticSearchPersistenceStorage) Scopes(funcs ...func(*gorm.DB) *gorm.DB) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Preload(query string, args ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Attrs(attrs ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Assign(attrs ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Unscoped() IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Raw(sql string, values ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Migrator() IPersistenceMigrator {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) AutoMigrate(dst ...interface{}) error {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Session(config *gorm.Session) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) WithContext(ctx context.Context) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Debug() IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Set(key string, value interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Get(key string) (interface{}, bool) {
	return nil, false
}

func (gps *ElasticSearchPersistenceStorage) InstanceSet(key string, value interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) InstanceGet(key string) (interface{}, bool) {
	return nil, false
}

func (gps *ElasticSearchPersistenceStorage) AddError(err error) error {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) DB() (*sql.DB, error) {
	return nil, nil
}

func (gps *ElasticSearchPersistenceStorage) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Use(plugin gorm.Plugin) error {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Create(value interface{}) IPersistenceStorage {
	es := gps.ESClient.Index().Index(gps.IndexName)
	if v1, ok := value.(ElasticModelInterface); ok {
		res, err := es.BodyJson(value).Refresh("wait_for").Do(context.Background())
		if err != nil {
			gps.LastError = err
		} else {
			v1.SetID(res.Id)
		}
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) CreateInBatches(value interface{}, batchSize int) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Save(value interface{}) IPersistenceStorage {
	es := gps.ESClient.Index().Index(gps.IndexName)
	if v1, ok := value.(ElasticModelInterface); ok {
		if v1.GetID() != "" {
			es = es.Id(v1.GetID())
		}
		res, err := es.BodyJson(value).Refresh("wait_for").Do(context.Background())
		if err != nil {
			gps.LastError = err
		} else {
			if v1.GetID() == "" {
				v1.SetID(res.Id)
			}
		}
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) First(dest interface{}, conds ...interface{}) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	searchRes, err := gps.ESSearch.Query(conds[0].(elastic.Query)).Pretty(true).Do(context.Background())
	if err != nil {
		gps.LastError = err
	} else {
		for _, hit := range searchRes.Hits.Hits {
			if err := json.Unmarshal(hit.Source, dest); err != nil {
				// Handle error
			}
			if v1, ok := dest.(ElasticModelInterface); ok {
				v1.SetID(hit.Id)
			}
			// Use v
		}
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) Take(dest interface{}, conds ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Last(dest interface{}, conds ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Find(dest interface{}, conds ...interface{}) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	// Query(elastic.NewTermQuery("user", "olivere")).
	searchRes, err := gps.ESSearch.Pretty(true).Do(context.Background())
	if err != nil {
		Trail(ERROR, err)
		gps.LastError = err
	} else {
		reflectValue := reflect.ValueOf(dest)
		var (
			reflectValueType = reflectValue.Type().Elem()
			isPtr            = reflectValueType.Kind() == reflect.Ptr
		)

		if isPtr {
			reflectValueType = reflectValueType.Elem()
		}
		reflectValue = reflect.Indirect(reflectValue)
		for _, hit := range searchRes.Hits.Hits {
			modelI, _ := gps.GenModelI()
			if err = json.Unmarshal(hit.Source, modelI); err != nil {
				// Handle error
			}
			if v1, ok := modelI.(ElasticModelInterface); ok {
				v1.SetID(hit.Id)
				reflectValue = reflect.Append(reflectValue, reflect.ValueOf(v1))
			}
			// Use v
		}
		valuePtr := reflect.ValueOf(dest)
		value := valuePtr.Elem()
		value.Set(reflectValue)
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) FirstOrInit(dest interface{}, conds ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) FirstOrCreate(dest interface{}, conds ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Update(column string, value interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Updates(values interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) UpdateColumn(column string, value interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) UpdateColumns(values interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Delete(value interface{}, conds ...interface{}) IPersistenceStorage {
	if v1, ok := value.(ElasticModelInterface); ok {
		_, err := gps.ESClient.Delete().Index(gps.IndexName).Id(v1.GetID()).Pretty(true).Do(context.Background())
		if err != nil {
			gps.LastError = err
		}
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) Count(count *int64) IPersistenceStorage {
	if gps.ESSearch == nil {
		gps.ESSearch = gps.ESClient.Search().Index(gps.IndexName)
	}
	searchRes, err := gps.ESSearch.Pretty(true).Do(context.Background())
	if err != nil {
		gps.LastError = err
	} else {
		*count = searchRes.Hits.TotalHits.Value
	}
	return gps
}

func (gps *ElasticSearchPersistenceStorage) Row() IPersistenceIterateRow {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Rows() (IPersistenceIterateRows, error) {
	return nil, nil
}

func (gps *ElasticSearchPersistenceStorage) Scan(dest interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Pluck(column string, dest interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) ScanRows(rows IPersistenceIterateRows, dest interface{}) error {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Begin(opts ...*sql.TxOptions) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Commit() IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Rollback() IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchAdminFilterObjects) SavePoint(name string) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) RollbackTo(name string) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) Exec(sql string, values ...interface{}) IPersistenceStorage {
	return nil
}

func (gps *ElasticSearchPersistenceStorage) GetLastError() error {
	return gps.LastError
}

func (gps *ElasticSearchPersistenceStorage) ResetLastError() {
	gps.LastError = nil
}

func (gps *ElasticSearchPersistenceStorage) LoadDataForModelByID(modelI interface{}, ID string) IPersistenceStorage {
	cond := elastic.NewTermQuery("_id", ID)
	gps.First(modelI, cond)
	return gps
}

type ESSortBy struct {
	Field         *Field
	Direction     bool
	FieldName     string
	SortByKeyword bool
}

type ESSearchParam struct {
	FieldName string
}

type ElasticSearchAdminFilterObjects struct {
	InitialESQuerySet   IPersistenceStorage
	ESQuerySet          IPersistenceStorage
	PaginatedESQuerySet IPersistenceStorage
	Model               interface{}
	DatabaseInstance    *ProjectDatabase
	SearchBy            []*ESSearchParam
	SearchString        string
	LastError           error
}

func (afo *ElasticSearchAdminFilterObjects) ResetLastError() {
	afo.LastError = nil
}

func (afo *ElasticSearchAdminFilterObjects) GetLastError() error {
	ret := afo.LastError
	return ret
}

func (afo *ElasticSearchAdminFilterObjects) SetLastError(err error) {
	if err != nil {
		afo.LastError = err
	}
}

func (afo *ElasticSearchAdminFilterObjects) FilterByMultipleIds(field *Field, realObjectIds []string) {
	objectInterfaceIds := make([]interface{}, 0)
	for _, realObjectId := range realObjectIds {
		if realObjectId == "" {
			continue
		}
		objectInterfaceIds = append(objectInterfaceIds, realObjectId)
	}
	afo.SetFullQuerySet(afo.GetFullQuerySet().Where("_id", objectInterfaceIds...))
	afo.SetLastError(afo.ESQuerySet.GetLastError())
}

func (afo *ElasticSearchAdminFilterObjects) Search(field *Field, searchString string) {
	afo.SearchBy = append(afo.SearchBy, &ESSearchParam{
		FieldName: field.DBName,
	})
}

func (afo *ElasticSearchAdminFilterObjects) SearchInIndex() {
	fields := make([]string, 0)
	for _, search := range afo.SearchBy {
		fields = append(fields, search.FieldName)
	}
	afo.ESQuerySet.Where(elastic.NewMultiMatchQuery(afo.SearchString, fields...))
	afo.PaginatedESQuerySet.Where(elastic.NewMultiMatchQuery(afo.SearchString, fields...))
	afo.SetLastError(afo.PaginatedESQuerySet.GetLastError())
}

func (afo *ElasticSearchAdminFilterObjects) FilterQs(filterString string) {
	statement := &gorm.Statement{DB: afo.GetDatabase().Db}
	statement.Parse(afo.GetCurrentModel())
	schema1 := statement.Schema
	FilterElasticSearchModel(afo.GetFullQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	FilterElasticSearchModel(afo.GetPaginatedQuerySet(), schema1, []string{filterString}, afo.GetCurrentModel())
	afo.SetLastError(afo.PaginatedESQuerySet.GetLastError())
}

func (afo *ElasticSearchAdminFilterObjects) GetPaginatedQuerySet() IPersistenceStorage {
	return afo.PaginatedESQuerySet
}

func (afo *ElasticSearchAdminFilterObjects) GetFullQuerySet() IPersistenceStorage {
	return afo.ESQuerySet
}

func (afo *ElasticSearchAdminFilterObjects) SetFullQuerySet(storage IPersistenceStorage) {
	afo.ESQuerySet = storage
}

func (afo *ElasticSearchAdminFilterObjects) GetInitialQuerySet() IPersistenceStorage {
	return afo.InitialESQuerySet
}

func (afo *ElasticSearchAdminFilterObjects) SetInitialQuerySet(storage IPersistenceStorage) {
	afo.InitialESQuerySet = storage
}

func (afo *ElasticSearchAdminFilterObjects) GetCurrentModel() interface{} {
	return afo.Model
}

func (afo *ElasticSearchAdminFilterObjects) GetDatabase() *ProjectDatabase {
	return nil
}

func (afo *ElasticSearchAdminFilterObjects) SetPaginatedQuerySet(storage IPersistenceStorage) {
	afo.PaginatedESQuerySet = storage
}

func (afo *ElasticSearchAdminFilterObjects) WithTransaction(handler func(afo1 IAdminFilterObjects) error) error {
	err := handler(afo)
	afo.SetLastError(err)
	return afo.GetLastError()
}

func (afo *ElasticSearchAdminFilterObjects) LoadDataForModelByID(ID interface{}, model interface{}) {
	cond := elastic.NewTermQuery("_id", ID)
	afo.InitialESQuerySet.First(model, cond)
	afo.SetLastError(afo.InitialESQuerySet.GetLastError())
}

func (afo *ElasticSearchAdminFilterObjects) SortBy(field *Field, direction int) {
	directionB := true
	if direction == -1 {
		directionB = false
	}
	afo.PaginatedESQuerySet.Order(&ESSortBy{
		Field:         field,
		Direction:     directionB,
		SortByKeyword: true,
	})
}

func (afo *ElasticSearchAdminFilterObjects) SaveModel(model interface{}) error {
	afo.InitialESQuerySet.Save(model)
	afo.SetLastError(afo.InitialESQuerySet.GetLastError())
	return nil
}

func (afo *ElasticSearchAdminFilterObjects) CreateNew(model interface{}) error {
	afo.InitialESQuerySet.Create(model)
	afo.SetLastError(afo.InitialESQuerySet.GetLastError())
	return nil
}

func (afo *ElasticSearchAdminFilterObjects) RemoveModelPermanently(model interface{}) error {
	afo.InitialESQuerySet.Delete(model)
	afo.SetLastError(afo.InitialESQuerySet.GetLastError())
	return nil
}

func (afo *ElasticSearchAdminFilterObjects) GetDB() IPersistenceStorage {
	modelDesc := ProjectModels.GetModelFromInterface(afo.Model)
	return NewElasticSearchPersistenceStorage(
		afo.InitialESQuerySet.(*ElasticSearchPersistenceStorage).ESClient,
		afo.InitialESQuerySet.(*ElasticSearchPersistenceStorage).IndexName,
		modelDesc.GenerateModelI,
	)
}

func (afo *ElasticSearchAdminFilterObjects) GetPaginated() <-chan *IterateAdminObjects {
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
		_, models := modelDesc.GenerateModelI()
		afo.PaginatedESQuerySet.Find(models)
		afo.SetLastError(afo.PaginatedESQuerySet.GetLastError())
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			if v1, ok := model.(ElasticModelInterface); ok {
				ID := v1.GetID()
				yieldV := &IterateAdminObjects{
					Model:         model,
					ID:            ID,
					RenderContext: &FormRenderContext{Model: model},
				}
				chnl <- yieldV
			}
		}
	}()
	return chnl
}

func (afo *ElasticSearchAdminFilterObjects) IterateThroughWholeQuerySet() <-chan *IterateAdminObjects {
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
		_, models := modelDesc.GenerateModelI()
		afo.ESQuerySet.Find(models)
		afo.SetLastError(afo.ESQuerySet.GetLastError())
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			if v1, ok := model.(ElasticModelInterface); ok {
				ID := v1.GetID()
				yieldV := &IterateAdminObjects{
					Model:         model,
					ID:            ID,
					RenderContext: &FormRenderContext{Model: model},
				}
				chnl <- yieldV
			}
		}
	}()
	return chnl
}

var ESGlobalModelActionRegistry *AdminModelActionRegistry

func init() {
	ESGlobalModelActionRegistry = NewEmptyModelActionRegistry()
	removalModelAction := NewAdminModelAction(
		"Delete permanently", &AdminActionPlacement{
			ShowOnTheListPage: true,
		},
	)
	removalModelAction.RequiresExtraSteps = true
	removalModelAction.PermName = "delete"
	removalModelAction.Description = "Delete permanently"
	removalModelAction.Handler = func(ap *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
		type Context struct {
			AdminContext
		}
		c := &Context{}
		adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
		PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)
		user := c.GetUserObject()
		if !ap.DoesUserHavePermission(user, "delete") {
			ctx.AbortWithStatus(409)
			return false, 0
		}
		removalPlan := make([]RemovalTreeList, 0)
		removalConfirmed := ctx.PostForm("removal_confirmed")
		removalError := afo.WithTransaction(func(afo1 IAdminFilterObjects) error {
			for modelIterated := range afo.IterateThroughWholeQuerySet() {
				if v1, ok := modelIterated.Model.(ElasticModelInterface); ok {
					if removalConfirmed == "" {
						deletionStringified := []*RemovalTreeNodeStringified{{
							Explanation: template.HTML(fmt.Sprintf("Delete from index %s - %s", v1.GetIndexName(), v1.String())),
							Level:       0,
						}}
						removalPlan = append(removalPlan, deletionStringified)
					} else {
						afo.GetInitialQuerySet().Delete(modelIterated.Model)
						if afo.GetInitialQuerySet().GetLastError() != nil {
							return afo.GetInitialQuerySet().GetLastError()
						}
					}
				}
			}
			if removalConfirmed != "" {
				query := ctx.Request.URL.Query()
				query.Set("message", "Objects were removed succesfully")
				ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/?%s", CurrentConfig.D.GoMonolith.RootAdminURL, ap.ParentPage.Slug, ap.ModelName, query.Encode()))
				return nil
			}
			type Context struct {
				AdminContext
				RemovalPlan []RemovalTreeList
				AdminPage   *AdminPage
				ObjectIds   string
			}
			c := &Context{}
			adminRequestParams := NewAdminRequestParams()
			c.RemovalPlan = removalPlan
			c.AdminPage = ap
			c.ObjectIds = ctx.PostForm("object_ids")
			PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := NewTemplateRenderer(fmt.Sprintf("Remove %s ?", ap.ModelName))
			tr.Render(ctx, CurrentConfig.GetPathToTemplate("remove_objects"), c, FuncMap)
			return nil
		})
		if removalError != nil {
			return false, 0
		}
		return true, int64(len(ctx.PostForm("object_ids")))
	}
	ESGlobalModelActionRegistry.AddModelAction(removalModelAction)
}

func NewAdminModelActionRegistryForElasticSearch() *AdminModelActionRegistry {
	adminModelActions := make(map[string]*AdminModelAction)
	ret := &AdminModelActionRegistry{AdminModelActions: adminModelActions}
	if ESGlobalModelActionRegistry != nil {
		for adminModelAction := range ESGlobalModelActionRegistry.GetAllModelActions() {
			ret.AddModelAction(adminModelAction)
		}
	}
	return ret
}

func FilterElasticSearchModel(db IPersistenceStorage, schema1 *schema.Schema, filterString []string, model interface{}) {
	context := NewGormOperatorContext(db, model)
	context.Tx = db
	gormModelV := reflect.Indirect(reflect.ValueOf(model))
	for _, filter := range filterString {
		filterParams := strings.Split(filter, "=")
		filterName := filterParams[0]
		filterValue := filterParams[1]
		filterNameParams := strings.Split(filterName, "__")
		field, _ := schema1.FieldsByName[filterNameParams[0]]
		field1 := NewGoMonolithFieldFromGormField(gormModelV, field, nil, false)
		operator, _ := ProjectGormOperatorRegistry.GetOperatorByName(filterNameParams[len(filterNameParams)-1])
		filterValueTransformed := operator.TransformValue(filterValue)
		context = operator.Build(nil, context, field1, filterValueTransformed, &SQLConditionBuilder{Type: "and"})
	}
}
