package core

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
)

type FormListEditable struct {
	FieldRegistry     IFieldRegistry
	Prefix            string
	FormRenderContext *FormRenderContext
	FormError         *FormError
}

type InlineFormListEditableCollection map[string]*FormListEditable

type FormListEditableCollection struct {
	InlineCollection map[string]InlineFormListEditableCollection
}

func (flec *FormListEditableCollection) AddForInline(prefix string, ID string, formListEditable *FormListEditable) {
	if flec.InlineCollection[prefix] == nil {
		flec.InlineCollection[prefix] = make(InlineFormListEditableCollection)
	}
	flec.InlineCollection[prefix][ID] = formListEditable
}

func (flec *FormListEditableCollection) GetForInlineAndForModel(prefix string, ID string) *FormListEditable {
	return flec.InlineCollection[prefix][ID]
}

func (flec *FormListEditableCollection) GetForInlineNew(prefix string) <-chan *FormListEditable {
	chnl := make(chan *FormListEditable)
	go func() {
		defer close(chnl)
		for modelID, ret := range flec.InlineCollection[prefix] {
			if !strings.Contains(modelID, "new") {
				continue
			}
			chnl <- ret
		}
	}()
	return chnl
}

func (flec *FormListEditableCollection) AddForInlineWholeCollection(prefix string, collection InlineFormListEditableCollection) {
	if flec.InlineCollection[prefix] == nil {
		flec.InlineCollection[prefix] = make(InlineFormListEditableCollection)
	}
	flec.InlineCollection[prefix] = collection
}

func NewFormListEditableCollection() *FormListEditableCollection {
	return &FormListEditableCollection{InlineCollection: make(map[string]InlineFormListEditableCollection)}
}

func (f *FormListEditable) SetPrefix(prefix string) {
	f.Prefix = prefix
	for _, field := range f.FieldRegistry.GetAllFields() {
		field.FieldConfig.Widget.SetPrefix(prefix)
	}
}

func (f *FormListEditable) ExistsField(ld *ListDisplay) bool {
	_, err := f.FieldRegistry.GetByName(ld.Field.Name)
	return err == nil
}

func (f *FormListEditable) ProceedRequest(form *multipart.Form, gormModel interface{}, adminContext IAdminContext) *FormError {
	formError := &FormError{
		FieldError:    make(map[string]ValidationError),
		GeneralErrors: make(ValidationError, 0),
	}
	renderContext := &FormRenderContext{Context: adminContext, Model: gormModel}
	for fieldName, field := range f.FieldRegistry.GetAllFields() {
		errors1 := field.ProceedForm(form, nil, renderContext)
		if len(errors1) == 0 {
			continue
		}
		field.FieldConfig.Widget.SetErrors(errors1)
		formError.FieldError[fieldName] = errors1
	}
	if formError.IsEmpty() {
		valueOfModel := reflect.ValueOf(gormModel)
		model := valueOfModel.Elem()
		for _, field := range f.FieldRegistry.GetAllFields() {
			modelF := model.FieldByName(field.Name)
			if !field.FieldConfig.Widget.IsValueChanged() {
				continue
			}
			if !modelF.IsValid() {
				formError.AddGeneralError(
					NewHTTPErrorResponse("field_invalid", "not valid field %s for model", field.Name),
				)
				continue
			}
			if !modelF.CanSet() {
				formError.AddGeneralError(
					NewHTTPErrorResponse("cant_set_field", "can't set field %s for model", field.Name),
				)
				continue
			}
			if field.SetUpField != nil {
				err := field.SetUpField(field.FieldConfig.Widget, gormModel, field.FieldConfig.Widget.GetOutputValue(), nil)
				if err != nil {
					formError.AddGeneralError(err)
				}
				continue
			}
			err := SetUpStructField(modelF, field.FieldConfig.Widget.GetOutputValue())
			if err != nil {
				formError.AddGeneralError(err)
			}
		}
	}
	f.FormRenderContext = &FormRenderContext{Model: gormModel}
	f.FormError = formError
	return formError
}

func NewFormListEditableForNewModelFromListDisplayRegistry(adminContext IAdminContext, prefix string, ID string, model interface{}, listDisplayRegistry *ListDisplayRegistry) *FormListEditable {
	modelForm := NewFormFromModel(model, []string{}, []string{}, false, "")
	modelForm.ForAdminPanel = true
	ret := &FormListEditable{FieldRegistry: NewFieldRegistry()}
	ret.SetPrefix(prefix)
	renderer := NewTemplateRenderer("")
	renderer.AddFuncMap("Translate", func(v interface{}) string {
		return Tf(adminContext.GetLanguage().Code, v)
	})
	for ld := range listDisplayRegistry.GetAllFields() {
		if ld.IsEditable && ld.Field.Name != "ID" {
			fieldFromNewForm, _ := modelForm.FieldRegistry.GetByName(ld.Field.Name)
			name := fieldFromNewForm.FieldConfig.Widget.GetHTMLInputName()
			if ret.Prefix != "" {
				fieldFromNewForm.FieldConfig.Widget.SetPrefix(ret.Prefix)
			}
			fieldFromNewForm.FieldConfig.Widget.SetName(fmt.Sprintf("%s_%s", ID, name))
			fieldFromNewForm.FieldConfig.Widget.SetShowOnlyHTMLInput()
			fieldFromNewForm.FieldConfig.Widget.RenderForAdmin()
			fieldFromNewForm.FieldConfig.Widget.RenderUsingRenderer(renderer)
			fieldFromNewForm.SetUpField = ld.Field.SetUpField
			ld.Field.FieldConfig.Widget.CloneAllOtherImportantSettings(fieldFromNewForm.FieldConfig.Widget)
			ret.FieldRegistry.AddField(fieldFromNewForm)
		}
	}
	ret.FormRenderContext = &FormRenderContext{Model: model}
	ret.FormError = &FormError{
		FieldError:    make(map[string]ValidationError),
		GeneralErrors: make(ValidationError, 0),
	}
	return ret
}

func NewFormListEditableFromListDisplayRegistry(adminContext IAdminContext, prefix string, ID string, model interface{}, listDisplayRegistry *ListDisplayRegistry, formError error) *FormListEditable {
	modelForm := NewFormFromModel(model, []string{}, []string{}, false, "")
	modelForm.ForAdminPanel = true
	ret := &FormListEditable{FieldRegistry: NewFieldRegistry()}
	ret.SetPrefix(prefix)
	renderer := NewTemplateRenderer("")
	renderer.AddFuncMap("Translate", func(v interface{}) string {
		return Tf(adminContext.GetLanguage().Code, v)
	})
	if formError != nil {
		if formErrorT, ok := formError.(*FormError); ok {
			modelForm.FormError = formErrorT
		}
	}
	if modelForm.FormError == nil {
		ret.FormError = &FormError{
			FieldError:    make(map[string]ValidationError),
			GeneralErrors: make(ValidationError, 0),
		}
	} else {
		ret.FormError = modelForm.FormError
	}
	for ld := range listDisplayRegistry.GetAllFields() {
		if ld.IsEditable && ld.Field.Name != "ID" {
			fieldFromNewForm, _ := modelForm.FieldRegistry.GetByName(ld.Field.Name)
			name := fieldFromNewForm.FieldConfig.Widget.GetHTMLInputName()
			if ret.Prefix != "" {
				fieldFromNewForm.FieldConfig.Widget.SetPrefix(ret.Prefix)
			}
			fieldFromNewForm.FieldConfig.Widget.SetName(fmt.Sprintf("%s_%s", ID, name))
			fieldFromNewForm.FieldConfig.Widget.SetShowOnlyHTMLInput()
			fieldFromNewForm.FieldConfig.Widget.RenderForAdmin()
			fieldFromNewForm.FieldConfig.Widget.RenderUsingRenderer(renderer)
			fieldFromNewForm.SetUpField = ld.Field.SetUpField
			ld.Field.FieldConfig.Widget.CloneAllOtherImportantSettings(fieldFromNewForm.FieldConfig.Widget)
			if formError != nil {
				if formErrorT, ok := formError.(*FormError); ok {
					fieldError, ok := formErrorT.FieldError[ld.Field.Name]
					if ok {
						fieldFromNewForm.FieldConfig.Widget.SetErrors(fieldError)
					}
				}
			}
			ret.FieldRegistry.AddField(fieldFromNewForm)
		}
	}
	ret.FormRenderContext = &FormRenderContext{Model: model}
	return ret
}
