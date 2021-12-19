---
sidebar_position: 1
---

# Admin filter objects

Admin filter objects implements core.IAdminFilterObjects interface.

```go
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
	GetDatabase() *Database
	GetCurrentModel() interface{}
	GetInitialQuerySet() IPersistenceStorage
	SetInitialQuerySet(IPersistenceStorage)
	GenerateModelInterface() (interface{}, interface{})
	RemoveModelPermanently(model interface{}) error
	FilterQs(filterString string)
	Search(field *Field, searchString string)
	SortBy(field *Field, direction int)
	FilterByMultipleIds(field *Field, realObjectIds []string)
	GetDB() IPersistenceStorage
	GetLastError() error
}
```
Right now go-monolith supports only objects that stored in database or elasticsearch and we use corresponding go package to interact with storage.
Later on we want to provide implementations for the objects stored in NoSQL, like Mongo, etc
