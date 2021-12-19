package core

type IAdminFilterObjects interface {
	WithTransaction(handler func(afo1 IAdminFilterObjects) error) error
	LoadDataForModelByID(ID interface{}, model interface{})
	SaveModel(model interface{}) error
	CreateNew(model interface{}) error
	GetPaginated() <-chan *IterateAdminObjects
	IterateThroughWholeQuerySet() <-chan *IterateAdminObjects
	GetPaginatedQuerySet() IPersistenceStorage
	GetFullQuerySet() IPersistenceStorage
	SetFullQuerySet(IPersistenceStorage)
	SetPaginatedQuerySet(IPersistenceStorage)
	GetDatabase() *ProjectDatabase
	GetCurrentModel() interface{}
	GetInitialQuerySet() IPersistenceStorage
	SetInitialQuerySet(IPersistenceStorage)
	RemoveModelPermanently(model interface{}) error
	FilterQs(filterString string)
	Search(field *Field, searchString string)
	SortBy(field *Field, direction int)
	FilterByMultipleIds(field *Field, realObjectIds []string)
	GetDB() IPersistenceStorage
	GetLastError() error
}

type IterateAdminObjects struct {
	Model         interface{}
	ID            string
	RenderContext *FormRenderContext
}
