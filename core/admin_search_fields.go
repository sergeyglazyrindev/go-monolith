package core

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type ISearchFieldInterface interface {
	Search(afo IAdminFilterObjects, searchString string)
	GetField() *Field
	SetCustomSearch(func(afo IAdminFilterObjects, searchString string))
}

type SearchField struct {
	Field        *Field
	CustomSearch func(afo IAdminFilterObjects, searchString string)
}

func (sf *SearchField) Search(afo IAdminFilterObjects, searchString string) {
	if sf.CustomSearch != nil {
		sf.CustomSearch(afo, searchString)
	} else {
		afo.Search(sf.Field, searchString)
	}
}

func (sf *SearchField) GetField() *Field {
	return sf.Field
}

func (sf *SearchField) SetCustomSearch(customSearch func(afo IAdminFilterObjects, searchString string)) {
	sf.CustomSearch = customSearch
}

//type ArraySearchField struct {
//	SearchField
//}
//
//func (sf *ArraySearchField) Search(afo IAdminFilterObjects, searchString string) {
//	if sf.CustomSearch != nil {
//		sf.CustomSearch(afo, searchString)
//	} else {
//		operator := ArrayIncludesGormOperator{}
//		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
//		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, sf.Field, searchString, &SQLConditionBuilder{Type: "or"})
//		afo.SetFullQuerySet(gormOperatorContext.Tx)
//		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
//		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, sf.Field, searchString, &SQLConditionBuilder{Type: "or"})
//		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
//		afo.SetLastError(afo.GetPaginatedQuerySet().GetLastError())
//	}
//}

//type JSONSearchField struct {
//	SearchField
//}
//
//func (sf *JSONSearchField) Search(afo IAdminFilterObjects, searchString string) {
//	if sf.CustomSearch != nil {
//		sf.CustomSearch(afo, searchString)
//	} else {
//		operator := JSONIncludesGormOperator{}
//		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
//		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, sf.Field, searchString, &SQLConditionBuilder{Type: "or"})
//		afo.SetFullQuerySet(gormOperatorContext.Tx)
//		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
//		operator.Build(afo.GetDatabase().Adapter, gormOperatorContext, sf.Field, searchString, &SQLConditionBuilder{Type: "or"})
//		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
//		afo.SetLastError(afo.GetPaginatedQuerySet().GetLastError())
//	}
//}
//
//func (sf *JSONSearchField) GetField() *Field {
//	return sf.Field
//}
//
//func (sf *JSONSearchField) SetCustomSearch(customSearch func(afo IAdminFilterObjects, searchString string)) {
//	sf.CustomSearch = customSearch
//}

func NewSearchFieldRegistryFromGormModel(modelI interface{}) *SearchFieldRegistry {
	if modelI == nil {
		return nil
	}
	ret := &SearchFieldRegistry{Fields: make([]ISearchFieldInterface, 0)}
	database := NewDatabaseInstanceWithoutConnection()
	stmt := &gorm.Statement{DB: database.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		tag := field.Tag.Get("gomonolith")
		if !strings.Contains(tag, "search") && field.Name != "ID" {
			continue
		}
		field1 := NewGoMonolithFieldFromGormField(gormModelV, field, nil, true)
		searchField := &SearchField{
			Field: field1,
		}
		ret.AddField(searchField)
	}
	return ret
}

type SearchFieldRegistry struct {
	Fields []ISearchFieldInterface
}

func (sfr *SearchFieldRegistry) GetAll() <-chan ISearchFieldInterface {
	chnl := make(chan ISearchFieldInterface)
	go func() {
		defer close(chnl)
		for _, field := range sfr.Fields {
			chnl <- field
		}

	}()
	return chnl
}

func (sfr *SearchFieldRegistry) AddField(sf ISearchFieldInterface) {
	sfr.Fields = append(sfr.Fields, sf)
}

func (sfr *SearchFieldRegistry) GetFieldByName(fieldName string) ISearchFieldInterface {
	for _, field := range sfr.Fields {
		if field.GetField().Name == fieldName {
			return field
		}
	}
	return nil
}
