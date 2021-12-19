package core

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

type ModelDescription struct {
	Model          interface{}
	Statement      *gorm.Statement
	GenerateModelI func() (interface{}, interface{})
}

type ProjectModelRegistry struct {
	models map[string]*ModelDescription
}

func (pmr *ProjectModelRegistry) RegisterModel(generateModelI func() (interface{}, interface{})) {
	model, _ := generateModelI()
	database := NewDatabaseInstanceWithoutConnection()
	statement := &gorm.Statement{DB: database.Db}
	statement.Parse(model)
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	modelName := v.Type().Name()
	pmr.models[modelName] = &ModelDescription{Model: model, Statement: statement, GenerateModelI: generateModelI}
}

func (pmr *ProjectModelRegistry) Iterate() <-chan *ModelDescription {
	chnl := make(chan *ModelDescription)
	go func() {
		defer close(chnl)
		for _, modelDescription := range pmr.models {
			chnl <- modelDescription
		}
	}()
	return chnl
}

func (pmr *ProjectModelRegistry) GetModelByName(modelName string) *ModelDescription {
	model, exists := pmr.models[modelName]
	if !exists {
		panic(fmt.Errorf("no model with name %s registered in the project", modelName))
	}
	return model
}

func (pmr *ProjectModelRegistry) GetModelFromInterface(model interface{}) *ModelDescription {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	modelName := v.Type().Name()
	modelI, _ := pmr.models[modelName]
	return modelI
}

var ProjectModels *ProjectModelRegistry

func init() {
	ProjectModels = &ProjectModelRegistry{
		models: make(map[string]*ModelDescription),
	}
}

func ClearProjectModels() {
	ProjectModels = &ProjectModelRegistry{
		models: make(map[string]*ModelDescription),
	}
}
