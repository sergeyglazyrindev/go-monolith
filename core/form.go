package core

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
	"html/template"
	"mime/multipart"
	"reflect"
	"strings"
)

type FormRenderContext struct {
	Model   interface{}
	Context IAdminContext
	Field *Field
}

func NewFormRenderContext() *FormRenderContext {
	return &FormRenderContext{}
}

type ColumnSchema struct {
	ShowLabel bool
	Fields    []*Field
}

type FormRow struct {
	Columns []*ColumnSchema
}

type IGrouppedFieldsRegistry interface {
	AddGroup(grouppedFields *GrouppedFields)
	GetGroupByName(name string) *GrouppedFields
}

type GrouppedFieldsRegistry struct {
	GrouppedFields map[string]*GrouppedFields
}

func (tfr *GrouppedFieldsRegistry) GetGroupByName(name string) (*GrouppedFields, error) {
	gf, ok := tfr.GrouppedFields[name]
	if !ok {
		return nil, NewHTTPErrorResponse("field_not_found", "no field %s found", name)
	}
	return gf, nil
}

func (tfr *GrouppedFieldsRegistry) AddGroup(grouppedFields *GrouppedFields) {
	if _, err := tfr.GetGroupByName(grouppedFields.Name); err != nil {
		panic(err)
	}
	tfr.GrouppedFields[grouppedFields.Name] = grouppedFields
}

type GrouppedFields struct {
	Rows            []*FormRow
	ExtraCSSClasses []string
	Description     string
	Name            string
}

type StaticFiles struct {
	ExtraCSS []string
	ExtraJS  []string
}

type FormError struct {
	FieldError    map[string]ValidationError
	GeneralErrors ValidationError
}

func (fe *FormError) AddGeneralError(err error) {
	fe.GeneralErrors = append(fe.GeneralErrors, err)
}
func (fe *FormError) IsEmpty() bool {
	return len(fe.FieldError) == 0 && len(fe.GeneralErrors) == 0
}

func (fe *FormError) String() string {
	if len(fe.GeneralErrors) > 0 {
		return fe.GeneralErrors[0].Error()
	}
	return "Form validation not successful"
}

func (fe *FormError) Error() string {
	if len(fe.GeneralErrors) > 0 {
		return fe.GeneralErrors[0].Error()
	}
	return "Form validation not successful"
}

func (fe *FormError) GetErrorForField(fieldName string) ValidationError {
	vE, _ := fe.FieldError[fieldName]
	return vE
}

type Form struct {
	ExcludeFields       IFieldRegistry
	FieldsToShow        IFieldRegistry
	FieldRegistry       IFieldRegistry
	GroupsOfTheFields   *GrouppedFieldsRegistry
	TemplateName        string
	FormTitle           string
	Renderer            ITemplateRenderer
	RequestContext      map[string]interface{}
	ErrorMessage        string
	ExtraStatic         *StaticFiles `json:"-"`
	ForAdminPanel       bool
	FormError           *FormError
	DontGenerateFormTag bool
	Prefix              string
	RenderContext       *FormRenderContext
	ChangesSaved        bool
	Debug               bool
}

func (f *Form) SetPrefix(prefix string) {
	f.Prefix = prefix
	for _, field := range f.FieldRegistry.GetAllFields() {
		field.FieldConfig.Widget.SetPrefix(prefix)
	}
}

func (f *Form) Render() template.HTML {
	RenderFieldGroups := func(funcs1 template.FuncMap) func() template.HTML {
		return func() template.HTML {
			templateWriter := bytes.NewBuffer([]byte{})
			ret := make([]string, 0)
			for _, group := range f.GroupsOfTheFields.GrouppedFields {
				for _, row := range group.Rows {
					data2 := row
					templateWriter.Reset()
					path := "form/grouprow"
					if f.ForAdminPanel {
						path = "admin/form/grouprow"
					}
					if f.Debug {
						for _, column := range row.Columns {
							for _, field := range column.Fields {
								spew.Dump(field.Name, field.FieldConfig.Widget)
							}
						}
					}
					//err := f.Renderer.RenderAsString(CurrentConfig.GetPathToTemplate(path), data2, FuncMap, funcs1)
					//if err != nil {
					//	Trail(CRITICAL, "Error while parsing include of the template %s", "form/grouprow")
					//	return ""
					//}
					ret = append(ret, string(f.Renderer.RenderAsString(CurrentConfig.GetPathToTemplate(path), data2, FuncMap, funcs1)))
				}
			}
			return template.HTML(strings.Join(ret, "\n"))
		}
	}
	if f.GroupsOfTheFields == nil {
		f.GroupsOfTheFields = &GrouppedFieldsRegistry{}
		grouppedFields := make(map[string]*GrouppedFields)
		grouppedFields["default"] = &GrouppedFields{
			Rows:            make([]*FormRow, 0),
			ExtraCSSClasses: make([]string, 0),
			Name:            "Default",
		}
		for _, field := range f.FieldRegistry.GetAllFieldsWithOrdering() {
			formRow := &FormRow{
				Columns: make([]*ColumnSchema, 0),
			}
			formRow.Columns = append(formRow.Columns, &ColumnSchema{
				Fields: []*Field{field},
			})
			grouppedFields["default"].Rows = append(grouppedFields["default"].Rows, formRow)
		}
		f.GroupsOfTheFields.GrouppedFields = grouppedFields
	}
	FieldValue := func(fieldName string, currentField *Field) interface{} {
		field, _ := f.FieldRegistry.GetByName(fieldName)
		if field.FieldConfig.Widget.GetPopulate() != nil {
			return field.FieldConfig.Widget.GetPopulate()(field.FieldConfig.Widget, f.RenderContext, currentField)
		}
		return field.FieldConfig.Widget.GetValue()
	}
	func1 := make(template.FuncMap)
	func1["FormFieldValue"] = FieldValue
	func1["GetRenderContext"] = func() *FormRenderContext {
		return f.RenderContext
	}
	f.Renderer.AddFuncMap("Translate", func(v interface{}) string {
		return Tf(f.RenderContext.Context.GetLanguage().Code, v)
	})
	for _, field := range f.FieldRegistry.GetAllFields() {
		field.FieldConfig.Widget.RenderUsingRenderer(f.Renderer)
	}
	func1["RenderFieldGroups"] = RenderFieldGroups(func1)
	path := "form"
	if f.ForAdminPanel {
		path = "admin/form"
	}
	if f.TemplateName != "" {
		path = f.TemplateName
	}
	templateName := CurrentConfig.GetPathToTemplate(path)
	return f.Renderer.RenderAsString(
		templateName,
		f, FuncMap, func1,
	)
}

func (f *Form) ProceedRequest(form *multipart.Form, gormModel interface{}, adminContext IAdminContext, afoP ...IAdminFilterObjects) *FormError {
	var afo IAdminFilterObjects
	if len(afoP) > 0 {
		afo = afoP[0]
	}
	formError := &FormError{
		FieldError:    make(map[string]ValidationError),
		GeneralErrors: make(ValidationError, 0),
	}
	renderContext := &FormRenderContext{Context: adminContext}
	for fieldName, field := range f.FieldRegistry.GetAllFields() {
		if field.Name == "ID" {
			continue
		}
		if field.ReadOnly {
			continue
		}
		renderContext.Field = field
		errors := field.ProceedForm(form, afo, renderContext)
		if len(errors) == 0 {
			continue
		}
		formError.FieldError[fieldName] = errors
	}
	valueOfModel := reflect.ValueOf(gormModel)
	model := valueOfModel.Elem()
	for _, field := range f.FieldRegistry.GetAllFields() {
		if field.Name == "ID" {
			continue
		}
		if field.ReadOnly {
			continue
		}
		if !field.FieldConfig.Widget.IsValueChanged() {
			continue
		}
		modelF := model.FieldByName(field.Name)
		if !modelF.IsValid() {
			formError.AddGeneralError(
				NewHTTPErrorResponse("field_invalid", "not valid field %s for model", field.Name),
			)
			continue
		}
		if !formError.IsEmpty() {
			continue
		}
		if formError.IsEmpty() && field.SetUpField != nil {
			err := field.SetUpField(field.FieldConfig.Widget, gormModel, field.FieldConfig.Widget.GetOutputValue(), afo)
			if err != nil {
				formError.AddGeneralError(err)
			}
			continue
		}
		if !modelF.CanSet() {
			formError.AddGeneralError(
				NewHTTPErrorResponse("cant_set_field", "can't set field %s for model", field.Name),
			)
			continue
		}
		err := SetUpStructField(modelF, field.FieldConfig.Widget.GetOutputValue())
		if err != nil {
			formError.AddGeneralError(err)
		}
	}
	f.FormError = formError
	return formError
}

func NewFormFromModel(gormModel interface{}, excludeFields []string, fieldsToShow []string, buildFieldPlacement bool, formTitle string, forAdminP ...bool) *Form {
	forAdmin := false
	if len(forAdminP) > 0 {
		forAdmin = forAdminP[0]
	}
	fieldRegistry := NewFieldRegistry()
	fieldsToShowRegistry := NewFieldRegistry()
	excludeFieldsRegistry := NewFieldRegistry()
	database := NewDatabaseInstance()
	defer database.Close()
	statement := &gorm.Statement{DB: database.Db}
	statement.Parse(gormModel)
	r := NewTemplateRenderer(formTitle)
	fields := statement.Schema.Fields
	gormModelV := reflect.Indirect(reflect.ValueOf(gormModel))
	for _, field := range fields {
		if len(fieldsToShow) > 0 && !Contains(fieldsToShow, field.Name) {
			if !field.PrimaryKey {
				continue
			}
		}
		fieldToBeExcluded := Contains(excludeFields, field.Name)
		if len(excludeFields) > 0 && fieldToBeExcluded {
			continue
		}
		field1 := NewGoMonolithFieldFromGormField(gormModelV, field, r, forAdmin)
		fieldRegistry.AddField(field1)
	}
	renderContext := NewFormRenderContext()
	renderContext.Model = gormModel
	form := &Form{
		ExcludeFields: excludeFieldsRegistry,
		FieldsToShow:  fieldsToShowRegistry,
		FieldRegistry: fieldRegistry,
		Renderer:      r,
		ExtraStatic: &StaticFiles{
			ExtraCSS: make([]string, 0),
			ExtraJS:  make([]string, 0),
		},
		FormError: &FormError{
			FieldError:    make(map[string]ValidationError),
			GeneralErrors: make(ValidationError, 0),
		},
		RenderContext: renderContext,
	}
	// form.GroupsOfTheFields.GrouppedFields = grouppedFields
	return form
}

func NewFormFromModelFromGinContext(contextFromGin IAdminContext, gormModel interface{}, excludeFields []string, fieldsToShow []string, buildFieldPlacement bool, formTitle string, forAdminP ...bool) *Form {
	forAdmin := false
	if len(forAdminP) > 0 {
		forAdmin = forAdminP[0]
	}
	form := NewFormFromModel(gormModel, excludeFields, fieldsToShow, buildFieldPlacement, formTitle, forAdmin)
	form.ForAdminPanel = forAdmin
	form.RequestContext = make(map[string]interface{})
	form.RequestContext["Language"] = contextFromGin.GetLanguage()
	form.RequestContext["RootURL"] = contextFromGin.GetRootURL()
	form.RequestContext["OTPImage"] = ""
	form.RequestContext["SessionKey"] = contextFromGin.GetSessionKey()
	form.RequestContext["ID"] = contextFromGin.GetID()
	form.RenderContext = &FormRenderContext{Context: contextFromGin, Model: gormModel}
	contextFromGin.SetForm(form)
	return form
}
