package core

import (
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"math"
	"reflect"
	"sort"
	"strings"
)

type ListDisplayRegistry struct {
	ListDisplayFields map[string]*ListDisplay
	MaxOrdering       int
	Prefix            string
	Placement         string
}

func (ldr *ListDisplayRegistry) GetFieldsCount() int {
	return len(ldr.ListDisplayFields)
}

func (ldr *ListDisplayRegistry) SetPrefix(prefix string) {
	ldr.Prefix = prefix
	for _, ld := range ldr.ListDisplayFields {
		ld.SetPrefix(prefix)
	}
}

func (ldr *ListDisplayRegistry) ClearAllFields() {
	ldr.MaxOrdering = 0
	ldr.ListDisplayFields = make(map[string]*ListDisplay)
}

func (ldr *ListDisplayRegistry) IsThereAnyEditable() bool {
	for ld := range ldr.GetAllFields() {
		if ld.IsEditable {
			return true
		}
	}
	return false
}
func (ldr *ListDisplayRegistry) AddField(ld *ListDisplay) {
	ldr.ListDisplayFields[ld.DisplayName] = ld
	ldr.MaxOrdering = int(math.Max(float64(ldr.MaxOrdering+1), float64(ld.Ordering+1)))
	ld.Ordering = ldr.MaxOrdering
}

func (ldr *ListDisplayRegistry) BuildFormForListEditable(adminContext IAdminContext, ID string, model interface{}, formError error) *FormListEditable {
	form := NewFormListEditableFromListDisplayRegistry(adminContext, ldr.Prefix, ID, model, ldr, formError)
	return form
}

func (ldr *ListDisplayRegistry) BuildListEditableFormForNewModel(adminContext IAdminContext, ID string, model interface{}) *FormListEditable {
	form := NewFormListEditableForNewModelFromListDisplayRegistry(adminContext, ldr.Prefix, ID, model, ldr)
	return form
}

func (ldr *ListDisplayRegistry) GetAllFields() <-chan *ListDisplay {
	chnl := make(chan *ListDisplay)
	go func() {
		defer close(chnl)
		dFields := make([]*ListDisplay, 0)
		for _, dField := range ldr.ListDisplayFields {
			dFields = append(dFields, dField)
		}
		sort.Slice(dFields, func(i, j int) bool {
			if dFields[i].Ordering == dFields[j].Ordering {
				return dFields[i].DisplayName < dFields[j].DisplayName
			}
			return dFields[i].Ordering < dFields[j].Ordering
		})
		for _, dField := range dFields {
			chnl <- dField
		}
	}()
	return chnl
}

func (ldr *ListDisplayRegistry) GetFieldByDisplayName(displayName string) (*ListDisplay, error) {
	listField, exists := ldr.ListDisplayFields[displayName]
	if !exists {
		return nil, fmt.Errorf("found no display field with name %s", displayName)
	}
	return listField, nil
}

func (ldr *ListDisplayRegistry) RemoveFieldByName(fieldName string) {
	delete(ldr.ListDisplayFields, fieldName)
}

type ListDisplay struct {
	DisplayName string
	Field       *Field
	ChangeLink  bool
	Ordering    int
	SortBy      ISortBy
	Populate    func(m interface{}) string
	MethodName  string
	IsEditable  bool
	Prefix      string
}

func (ld *ListDisplay) SetPrefix(prefix string) {
	ld.Prefix = prefix
}

func (ld *ListDisplay) GetOrderingName(initialOrdering []string) string {
	for _, part := range initialOrdering {
		negativeOrdering := false
		if strings.HasPrefix(part, "-") {
			part = part[1:]
			negativeOrdering = true
		}
		if part == ld.DisplayName {
			if negativeOrdering {
				return ld.DisplayName
			}
			return "-" + ld.DisplayName
		}
	}

	return ld.DisplayName
}

func (ld *ListDisplay) IsEligibleForOrdering() bool {
	return ld.SortBy != nil
}

func (ld *ListDisplay) GetValue(m interface{}, forExportP ...bool) template.HTML {
	forExport := false
	if len(forExportP) > 0 {
		forExport = forExportP[0]
	}
	if ld.MethodName != "" {
		values := reflect.ValueOf(m).MethodByName(ld.MethodName).Call([]reflect.Value{})
		return template.HTML(values[0].String())
	}
	if ld.Populate != nil {
		return template.HTML(ld.Populate(m))
	}
	if ld.Field.FieldConfig.Widget.GetPopulate() != nil {
		return template.HTML(TransformValueForListDisplay(ld.Field.FieldConfig.Widget.GetPopulate()(ld.Field.FieldConfig.Widget, &FormRenderContext{Model: m}, ld.Field)))
	}
	if ld.Field.FieldConfig.Widget.IsValueConfigured() {
		return template.HTML(TransformValueForListDisplay(ld.Field.FieldConfig.Widget.GetValue()))
	}
	gormModelV := reflect.Indirect(reflect.ValueOf(m))
	if reflect.ValueOf(m).IsZero() || gormModelV.IsZero() { // || gormModelV.FieldByName(ld.Field.Name).IsZero()
		return ""
	}
	return template.HTML(TransformValueForListDisplay(gormModelV.FieldByName(ld.Field.Name).Interface(), forExport))
}

func NewListDisplay(field *Field) *ListDisplay {
	displayName := ""
	if field != nil {
		displayName = field.DisplayName
	}
	return &ListDisplay{
		DisplayName: displayName, Field: field, ChangeLink: true,
		SortBy: &SortBy{Field: field, Direction: 1},
	}
}

func NewListDisplayRegistryFromGormModelForInlines(modelI interface{}) *ListDisplayRegistry {
	ret := &ListDisplayRegistry{
		ListDisplayFields: make(map[string]*ListDisplay),
	}
	database := NewDatabaseInstanceWithoutConnection()
	stmt := &gorm.Statement{DB: database.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		tag := field.Tag.Get("gomonolith")
		if !strings.Contains(tag, "inline") && field.Name != "ID" {
			continue
		}
		field1 := NewGoMonolithFieldFromGormField(gormModelV, field, nil, true)
		ld := NewListDisplay(field1)
		if field.Name != "ID" {
			ld.IsEditable = true
		}
		ret.AddField(ld)
	}
	return ret
}

func NewListDisplayRegistryFromGormModel(modelI interface{}) *ListDisplayRegistry {
	if modelI == nil {
		return nil
	}
	ret := &ListDisplayRegistry{
		ListDisplayFields: make(map[string]*ListDisplay),
	}
	database := NewDatabaseInstanceWithoutConnection()
	stmt := &gorm.Statement{DB: database.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		tag := field.Tag.Get("gomonolith")
		if !strings.Contains(tag, "list") && field.Name != "ID" {
			continue
		}
		field1 := NewGoMonolithFieldForListDisplayFromGormField(gormModelV, field, nil, true)
		ret.AddField(NewListDisplay(field1))
	}
	return ret
}
