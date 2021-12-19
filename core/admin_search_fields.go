package core

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type ISearchFieldInterface interface {
	Search(afo IAdminFilterObjects, searchString string)
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

func NewSearchFieldRegistryFromGormModel(modelI interface{}) *SearchFieldRegistry {
	if modelI == nil {
		return nil
	}
	ret := &SearchFieldRegistry{Fields: make([]*SearchField, 0)}
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
	Fields []*SearchField
}

func (sfr *SearchFieldRegistry) GetAll() <-chan *SearchField {
	chnl := make(chan *SearchField)
	go func() {
		defer close(chnl)
		for _, field := range sfr.Fields {
			chnl <- field
		}

	}()
	return chnl
}

func (sfr *SearchFieldRegistry) AddField(sf *SearchField) {
	sfr.Fields = append(sfr.Fields, sf)
}

func (sfr *SearchFieldRegistry) GetFieldByName(fieldName string) *SearchField {
	for _, field := range sfr.Fields {
		if field.Field.DBName == fieldName {
			return field
		}
	}
	return nil
}
