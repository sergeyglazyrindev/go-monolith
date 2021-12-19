package core

import (
	"fmt"
	"html/template"
	"mime/multipart"
	"strings"
)

type InlineType string

var TabularInline InlineType
var StackedInline InlineType

func init() {
	TabularInline = "tabular"
	StackedInline = "stacked"
}

type AdminPageInline struct {
	Ordering          int
	GenerateModelI    func(m interface{}) interface{}
	GetQueryset       func(adminContext IAdminContext, afo IAdminFilterObjects, model interface{}) IAdminFilterObjects
	Actions           *AdminModelActionRegistry
	EmptyValueDisplay string
	ExcludeFields     IFieldRegistry
	FieldsToShow      IFieldRegistry
	ShowAllFields     bool
	Validators        *ValidatorRegistry
	Classes           []string
	Extra             int
	MaxNum            int
	MinNum            int
	VerboseName       string
	VerboseNamePlural string
	ShowChangeLink    bool
	Template          string
	ContentType       *ContentType
	Permission        CustomPermission
	InlineType        InlineType
	Prefix            string
	ListDisplay       *ListDisplayRegistry `json:"-"`
}

func (api *AdminPageInline) RenderExampleForm(adminContext IAdminContext) template.HTML {
	type Context struct {
		AdminContext
		AdminContextInitial IAdminContext
		Inline              *AdminPageInline
	}
	c := &Context{}
	c.AdminContextInitial = adminContext
	c.Inline = api
	templateRenderer := NewTemplateRenderer("")
	func1 := make(template.FuncMap)
	path := "admin/inlineexampleform"
	templateName := CurrentConfig.GetPathToTemplate(path)
	templateRenderer.AddFuncMap("Translate", func(v interface{}) string {
		return Tf(adminContext.GetLanguage().Code, v)
	})
	return templateRenderer.RenderAsString(
		templateName,
		c, FuncMap, func1,
	)
}

func (api *AdminPageInline) GetFormForExample(adminContext IAdminContext) *FormListEditable {
	modelI := api.GenerateModelI(nil)
	form := api.ListDisplay.BuildListEditableFormForNewModel(adminContext, "toreplacewithid", modelI)
	r := NewTemplateRenderer("")
	r.AddFuncMap("Translate", func(v interface{}) string {
		return Tf(adminContext.GetLanguage().Code, v)
	})
	for _, field := range form.FieldRegistry.GetAllFields() {
		field.FieldConfig.Widget.RenderUsingRenderer(r)
	}
	//	return r.RenderAsString(templateName, data, baseFuncMap)
	return form
}

func (api *AdminPageInline) GetFormIdenForNewItems() string {
	return fmt.Sprintf("example-%s", api.Prefix)
}

func (api *AdminPageInline) GetInlineID() string {
	return PrepareStringToBeUsedForHTMLID(api.VerboseNamePlural)
}

func (api *AdminPageInline) GetAll(adminContext IAdminContext, model interface{}) <-chan *IterateAdminObjects {
	qs := api.GetQueryset(adminContext, nil, model)
	return qs.IterateThroughWholeQuerySet()
}

func (api *AdminPageInline) ProceedRequest(afo IAdminFilterObjects, f *multipart.Form, model interface{}, adminContext IAdminContext) (InlineFormListEditableCollection, error) {
	collection := make(InlineFormListEditableCollection)
	var firstEditableField *ListDisplay
	qs := api.GetQueryset(adminContext, afo, model)
	for ld := range api.ListDisplay.GetAllFields() {
		if ld.IsEditable {
			firstEditableField = ld
			break
		}
	}
	if firstEditableField == nil {
		return collection, nil
	}
	var form *FormListEditable
	err := false
	var removalError error
	for fieldName := range f.Value {
		if !strings.HasSuffix(fieldName, firstEditableField.Field.FieldConfig.Widget.GetHTMLInputName()) {
			continue
		}
		if !strings.HasPrefix(fieldName, firstEditableField.Prefix) {
			continue
		}
		if strings.Contains(fieldName, "toreplacewithid") {
			continue
		}
		removalError = nil
		form = nil
		inlineID := strings.TrimPrefix(fieldName, firstEditableField.Prefix+"-")
		inlineID = strings.TrimSuffix(inlineID, "-"+firstEditableField.Field.FieldConfig.Widget.GetHTMLInputName())
		realInlineID := strings.Split(inlineID, "_")
		modelI := api.GenerateModelI(model)
		inlineIDToRemove := f.Value[firstEditableField.Prefix+"-"+"object_id-to-remove-"+realInlineID[0]]
		isNew := false
		if !strings.Contains(inlineID, "new") {
			qs.LoadDataForModelByID(realInlineID[0], modelI)
			form = api.ListDisplay.BuildFormForListEditable(adminContext, realInlineID[0], modelI, nil)
			collection[realInlineID[0]] = form
			if len(inlineIDToRemove) > 0 {
				removalError = qs.RemoveModelPermanently(modelI)
			}
		} else {
			form = api.ListDisplay.BuildListEditableFormForNewModel(adminContext, realInlineID[0], modelI)
			collection[realInlineID[0]] = form
			isNew = true
		}
		if len(inlineIDToRemove) > 0 {
			if removalError != nil {
				form.FormError = &FormError{
					FieldError:    make(map[string]ValidationError),
					GeneralErrors: make(ValidationError, 0),
				}
				form.FormError.AddGeneralError(removalError)
			}
		} else {
			formError := form.ProceedRequest(f, modelI, adminContext)
			if removalError != nil {
				formError.AddGeneralError(formError)
			}
			if !formError.IsEmpty() {
				err = true
			} else {
				if isNew {
					error1 := afo.CreateNew(modelI)
					if error1 != nil {
						formError.AddGeneralError(error1)
						err = true
					}
				} else {
					error1 := afo.SaveModel(modelI)
					if error1 != nil {
						formError.AddGeneralError(error1)
						err = true
					}
				}
			}
			collection[realInlineID[0]].FormError = formError
		}
	}
	if err {
		return collection, NewHTTPErrorResponse("inline_validation_error", "error while validating inlines")
	}
	return collection, nil
}

func NewAdminPageInline(
	inlineIden string,
	inlineType InlineType,
	generateModelI func(m interface{}) interface{},
	getQuerySet func(adminContext IAdminContext, afo IAdminFilterObjects, model interface{}) IAdminFilterObjects,
) *AdminPageInline {
	modelI := generateModelI(nil)
	ld := NewListDisplayRegistryFromGormModelForInlines(modelI)
	ld.SetPrefix(PrepareStringToBeUsedForHTMLID(inlineIden))
	ret := &AdminPageInline{
		Actions:           NewEmptyModelActionRegistry(),
		ExcludeFields:     NewFieldRegistry(),
		FieldsToShow:      NewFieldRegistry(),
		Validators:        NewValidatorRegistry(),
		Classes:           make([]string, 0),
		InlineType:        inlineType,
		ListDisplay:       ld,
		GenerateModelI:    generateModelI,
		GetQueryset:       getQuerySet,
		VerboseNamePlural: inlineIden,
	}
	return ret
}

func NewAdminPageInlineRegistry() *AdminPageInlineRegistry {
	return &AdminPageInlineRegistry{
		Inlines: make([]*AdminPageInline, 0),
	}
}
