package core

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"html/template"
	"mime/multipart"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type WidgetType string

const UnknownInputWidgetType WidgetType = "unknown"
const TextInputWidgetType WidgetType = "text"
const DynamicInputWidgetType WidgetType = "dynamic"
const NumberInputWidgetType WidgetType = "number"
const EmailInputWidgetType WidgetType = "email"
const URLInputWidgetType WidgetType = "url"
const PasswordInputWidgetType WidgetType = "password"
const HiddenInputWidgetType WidgetType = "hidden"
const DateInputWidgetType WidgetType = "date"
const DateTimeInputWidgetType WidgetType = "datetime"
const TimeInputWidgetType WidgetType = "time"
const TextareaInputWidgetType WidgetType = "textarea"
const CheckboxInputWidgetType WidgetType = "checkbox"
const SelectWidgetType WidgetType = "select"
const ForeignKeyWidgetType WidgetType = "foreignkey"
const NullBooleanWidgetType WidgetType = "nullboolean"
const SelectMultipleWidgetType WidgetType = "selectmultiple"
const RadioSelectWidgetType WidgetType = "radioselect"
const RadioWidgetType WidgetType = "radio"
const CheckboxSelectMultipleWidgetType WidgetType = "checkboxselectmultiple"
const FileInputWidgetType WidgetType = "file"
const ClearableFileInputWidgetType WidgetType = "clearablefileinput"
const MultipleHiddenInputWidgetType WidgetType = "multiplehiddeninput"
const SplitDateTimeWidgetType WidgetType = "splitdatetime"
const SplitHiddenDateTimeWidgetType WidgetType = "splithiddendatetime"
const SelectDateWidgetType WidgetType = "selectdate"
const ChooseFromSelectWidgetType WidgetType = "choose_from_select"
const FkLinkWidgetType WidgetType = "fklink"
const JSONWidgetType WidgetType = "json"
const ArrayWidgetType WidgetType = "array"

type WidgetData map[string]interface{}
type IWidget interface {
	IDForLabel(model interface{}, F *Field) string
	GetWidgetType() WidgetType
	GetAttrs() map[string]string
	GetTemplateName() string
	SetTemplateName(templateName string)
	RenderUsingRenderer(renderer ITemplateRenderer)
	// GetValue(v interface{}, model interface{}) interface{}
	Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML
	SetValue(v interface{})
	SetName(name string)
	GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData
	SetAttr(attrName string, value string)
	SetBaseFuncMap(baseFuncMap template.FuncMap)
	InitializeAttrs()
	SetFieldDisplayName(displayName string)
	SetReadonly(readonly bool)
	GetValue() interface{}
	ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error
	SetRequired()
	SetShowOnlyHTMLInput()
	SetOutputValue(v interface{})
	GetOutputValue() interface{}
	SetErrors(validationErrors ValidationError)
	RenderForAdmin()
	SetHelpText(helpText string)
	IsValueChanged() bool
	SetPopulate(func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{})
	SetPrefix(prefix string)
	GetHTMLInputName() string
	GetPopulate() func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}
	IsReadOnly() bool
	IsValueConfigured() bool
	SetValueConfigured()
	GetRenderer() ITemplateRenderer
	GetFieldDisplayName() string
	GetName() string
	CloneAllOtherImportantSettings(widget IWidget)
}

func GetWidgetByWidgetType(widgetType string, fieldOptions IFieldFormOptions) IWidget {
	var widget IWidget
	//case 1:
	//	return &core.TextWidget{}
	//	case 2:
	//	return &core.NumberWidget{}
	//	case 3:
	//	widget := &core.NumberWidget{}
	//	widget.SetAttr("step", "0.1")
	//	return widget
	//	case 4:
	//	return &core.CheckboxWidget{}
	//	case 5:
	//	return &core.FileWidget{}
	//	case 6:
	//	widget := &core.FileWidget{}
	//	widget.SetAttr("accept", "image/*")
	//	return widget
	//	case 7:
	//	return &core.DateTimeWidget{}
	//
	switch widgetType {
	case "url":
		widget = &URLWidget{}
	case "file":
		widget = &FileWidget{}
	case "boolean":
		widget = &CheckboxWidget{}
	case "float":
		widget = &NumberWidget{}
		widget.SetAttr("step", "0.1")
	case "integer":
		widget = &NumberWidget{}
	case "string":
		widget = &TextWidget{}
	case "image":
		widget = &FileWidget{}
		widget.SetAttr("accept", "image/*")
	case "hidden":
		widget = &HiddenWidget{}
	case "password":
		widget = &PasswordWidget{}
	case "dynamic":
		widget = &DynamicWidget{}
	case "email":
		widget = &EmailWidget{}
	case "foreignkey":
		widget = &ForeignKeyWidget{}
		if fieldOptions != nil && fieldOptions.GetIsAutocomplete() {
			widget.(*ForeignKeyWidget).Autocomplete = fieldOptions.GetIsAutocomplete()
		}
	case "choose_from_select":
		widget = &ChooseFromSelectWidget{}
	case "fklink":
		widget = &FkLinkWidget{}
		widget.SetReadonly(true)
		widget.SetPopulate(func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{} {
			gormModelV := reflect.Indirect(reflect.ValueOf(renderContext.Model))
			if gormModelV.FieldByName(currentField.Name).IsZero() {
				return ""
			}
			adminPage := CurrentDashboardAdminPanel.FindPageForGormModel(gormModelV.FieldByName(currentField.Name).Interface())
			if adminPage != nil {
				link := adminPage.GenerateLinkToEditModel(gormModelV.FieldByName(currentField.Name))
				fkModel := reflect.New(reflect.TypeOf(gormModelV.FieldByName(currentField.Name).Interface()))
				fkModel.Elem().Set(reflect.ValueOf(gormModelV.FieldByName(currentField.Name).Interface()))
				stringRepresentation := fkModel.MethodByName("String").Call([]reflect.Value{})
				return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", link, stringRepresentation[0].Interface().(string))
			}
			fkModel := reflect.New(reflect.TypeOf(gormModelV.FieldByName(currentField.Name).Interface()))
			fkModel.Elem().Set(reflect.ValueOf(gormModelV.FieldByName(currentField.Name).Interface()))
			stringRepresentation := fkModel.MethodByName("String").Call([]reflect.Value{})
			return stringRepresentation
		})
	case "textarea":
		widget = &TextareaWidget{}
		widget.SetAttr("style", "width: 60%;height: 200px;")
	case "select":
		widget = &SelectWidget{}
	case "datetime":
		widget = &DateTimeWidget{}
	case "contenttypeselector":
		widget = &ContentTypeSelectorWidget{}
	}
	return widget
}

type Widget struct {
	Attrs             map[string]string
	TemplateName      string
	Renderer          ITemplateRenderer
	Value             interface{}
	Name              string
	FieldDisplayName  string
	BaseFuncMap       template.FuncMap
	ReadOnly          bool
	ShowOnlyHTMLInput bool
	Required          bool
	OutputValue       interface{}
	ValidationErrors  ValidationError
	IsForAdmin        bool
	HelpText          string
	ValueChanged      bool
	Populate          func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}
	Prefix            string
	ValueConfigured   bool
}

func (w *Widget) GetRenderer() ITemplateRenderer {
	return w.Renderer
}

func (w *Widget) GetFieldDisplayName() string {
	return w.FieldDisplayName
}

func (w *Widget) GetName() string {
	return w.Name
}

func (w *Widget) SetValueConfigured() {
	w.ValueConfigured = true
}

func (w *Widget) SetPrefix(prefix string) {
	w.Prefix = prefix
}

func (w *Widget) IsReadOnly() bool {
	return w.ReadOnly
}

func (w *Widget) IsValueConfigured() bool {
	return w.ValueConfigured
}

func (w *Widget) CloneAllOtherImportantSettings(widget IWidget) {
	widget.SetPopulate(w.Populate)
	widget.SetErrors(w.ValidationErrors)
}

func (w *Widget) IsValueChanged() bool {
	return w.ValueChanged
}

func (w *Widget) SetPopulate(pFunc func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}) {
	w.Populate = pFunc
}

func (w *Widget) GetPopulate() func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{} {
	return w.Populate
}

func (w *Widget) SetRequired() {
	w.Required = true
}

func (w *Widget) SetHelpText(helpText string) {
	w.HelpText = helpText
}

func (w *Widget) RenderForAdmin() {
	w.IsForAdmin = true
}

func (w *Widget) SetShowOnlyHTMLInput() {
	w.ShowOnlyHTMLInput = true
}

func (w *Widget) SetTemplateName(templateName string) {
	w.TemplateName = templateName
}

func (w *Widget) SetOutputValue(v interface{}) {
	w.OutputValue = v
	w.ValueChanged = true
}

func (w *Widget) GetOutputValue() interface{} {
	return w.OutputValue
}

func (w *Widget) SetErrors(validationErrors ValidationError) {
	w.ValidationErrors = validationErrors
}

func (w *Widget) InitializeAttrs() {
	if w.Attrs == nil {
		w.Attrs = make(map[string]string)
	}
}

func (w *Widget) SetBaseFuncMap(baseFuncMap template.FuncMap) {
	w.BaseFuncMap = baseFuncMap
}

func (w *Widget) IDForLabel(model interface{}, F *Field) string {
	return ""
}

func (w *Widget) SetFieldDisplayName(fieldDisplayName string) {
	w.FieldDisplayName = fieldDisplayName
}

func (w *Widget) SetReadonly(readonly bool) {
	w.ReadOnly = readonly
}

func (w *Widget) GetWidgetType() WidgetType {
	return UnknownInputWidgetType
}

func (w *Widget) GetTemplateName() string {
	return ""
}

func (w *Widget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

func (w *Widget) RenderUsingRenderer(r ITemplateRenderer) {
	w.Renderer = r
}

func (w *Widget) SetAttr(attrName string, value string) {
	if w.Attrs == nil {
		w.InitializeAttrs()
	}
	w.Attrs[attrName] = value
}

func (w *Widget) SetName(name string) {
	w.Name = name
}

func (w *Widget) GetAttrs() map[string]string {
	if w.Attrs != nil {
		return w.Attrs
	}
	return make(map[string]string)
}

func (w *Widget) SetValue(v interface{}) {
	w.Value = v
}

func (w *Widget) GetValue() interface{} {
	return w.Value
}

func (w *Widget) GetHTMLInputName() string {
	if w.Prefix != "" {
		return w.Prefix + "-" + w.Name
	}
	return w.Name
}

func (w *Widget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("1", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *Widget) GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData {
	var value interface{}
	var valueStr string
	if w.Populate != nil {
		value = w.Populate(w, formRenderContext, currentField)
		valueStr = value.(string)
	} else {
		value = TransformValueForWidget(w.Value)
		if value != nil {
			valueStr = template.HTMLEscapeString(TransformValueForListDisplay(value))
		} else {
			valueStr = ""
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(), "Value": template.HTML(valueStr),
		"Name": w.GetHTMLInputName(), "FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
		"Required": w.Required, "HelpText": w.HelpText, "FormError": w.ValidationErrors,
		"FormErrorNotEmpty": len(w.ValidationErrors) > 0,
	}
}

func RenderWidget(formRenderContext *FormRenderContext, renderer ITemplateRenderer, templateName string, data map[string]interface{}, baseFuncMap template.FuncMap) template.HTML {
	if renderer == nil {
		r := NewTemplateRenderer("")
		r.AddFuncMap("Translate", func(v interface{}) string {
			if formRenderContext.Context != nil && formRenderContext.Context.GetLanguage() != nil {
				return Tf(formRenderContext.Context.GetLanguage().Code, v)
			}
			return v.(string)
		})
		return r.RenderAsString(templateName, data, baseFuncMap)
	}
	return renderer.RenderAsString(
		templateName,
		data, baseFuncMap,
	)
}

type TextWidget struct {
	Widget
}

func (tw *TextWidget) GetWidgetType() WidgetType {
	return TextInputWidgetType
}

func (tw *TextWidget) GetTemplateName() string {
	if tw.TemplateName == "" {
		path := "widgets/text"
		if tw.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(tw.TemplateName)
}

func (tw *TextWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("2", tw.FieldDisplayName)
	data := tw.Widget.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = tw.GetWidgetType()
	data["ShowOnlyHtmlInput"] = tw.ShowOnlyHTMLInput
	return RenderWidget(formRenderContext, tw.Renderer, tw.GetTemplateName(), data, tw.BaseFuncMap) // tw.Value, tw.Widget.GetAttrs()
}

type DynamicWidget struct {
	Widget
	GetRealWidget                  func(formRenderContext *FormRenderContext, currentField *Field) IWidget
	GetRealWidgetForFormProceeding func(form *multipart.Form, afo IAdminFilterObjects) IWidget
}

func (tw *DynamicWidget) GetWidgetType() WidgetType {
	return DynamicInputWidgetType
}

func (tw *DynamicWidget) GetTemplateName() string {
	panic("shouldn't be called")
}

func (tw *DynamicWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	var realWidget IWidget
	if formRenderContext.Context.GetCtx().Query("widgetType") != "" {
		realWidget = GetWidgetByWidgetType(formRenderContext.Context.GetCtx().Query("widgetType"), nil)
	} else {
		realWidget = tw.GetRealWidget(formRenderContext, currentField)
	}
	realWidget.InitializeAttrs()
	if tw.IsForAdmin {
		realWidget.RenderForAdmin()
	}
	if tw.Populate != nil {
		realWidget.SetPopulate(tw.Populate)
	}
	if tw.Required {
		realWidget.SetRequired()
	}
	if tw.HelpText != "" {
		realWidget.SetHelpText(tw.HelpText)
	}
	realWidget.RenderUsingRenderer(tw.Renderer)
	realWidget.SetName(tw.GetName())
	realWidget.SetFieldDisplayName(tw.FieldDisplayName)
	realWidget.SetValue(tw.Value)
	realWidget.SetErrors(tw.ValidationErrors)
	ret := realWidget.Render(formRenderContext, currentField)
	return ret
}

func (tw *DynamicWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	var realWidget IWidget
	if renderContext.Context.GetCtx().Query("widgetType") != "" {
		realWidget = GetWidgetByWidgetType(renderContext.Context.GetCtx().Query("widgetType"), nil)
	} else {
		realWidget = tw.GetRealWidgetForFormProceeding(form, afo)
	}
	realWidget.InitializeAttrs()
	if tw.IsForAdmin {
		realWidget.RenderForAdmin()
	}
	if tw.Populate != nil {
		realWidget.SetPopulate(tw.Populate)
	}
	if tw.Required {
		realWidget.SetRequired()
	}
	if tw.HelpText != "" {
		realWidget.SetHelpText(tw.HelpText)
	}
	realWidget.RenderUsingRenderer(tw.Renderer)
	realWidget.SetFieldDisplayName(tw.FieldDisplayName)
	realWidget.SetName(tw.GetName())
	realWidget.SetErrors(tw.ValidationErrors)
	realWidget.SetValue(tw.Value)
	return realWidget.ProceedForm(form, afo, renderContext)
}

type FkLinkWidget struct {
	Widget
	Context string
}

func (w *FkLinkWidget) GetWidgetType() WidgetType {
	return FkLinkWidgetType
}

func (w *FkLinkWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/text"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *FkLinkWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	if w.IsReadOnly() && w.Context != "edit" {
		return template.HTML(w.Populate(w, formRenderContext, currentField).(string))
	}
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap) // tw.Value, tw.Widget.GetAttrs()
}

type NumberWidget struct {
	Widget
	NumberType GoMonolithFieldType
}

func (w *NumberWidget) GetWidgetType() WidgetType {
	return NumberInputWidgetType
}

func (w *NumberWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/number"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *NumberWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("3", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *NumberWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	if !govalidator.IsInt(v[0]) {
		return NewHTTPErrorResponse("should_be_number", "should be a number")
	}
	w.SetOutputValue(w.TransformValueForOutput(v[0]))
	return nil
}

func (w *NumberWidget) TransformValueForOutput(v string) interface{} {
	switch w.NumberType {
	case PositiveIntegerFieldType:
		vI, _ := strconv.ParseUint(v, 10, 32)
		return uint(vI)
	case BigIntegerFieldType:
		vI64, _ := strconv.ParseInt(v, 10, 64)
		return vI64
	case IntegerFieldType:
		vI32, _ := strconv.ParseInt(v, 10, 32)
		return int(vI32)
	case SmallIntegerFieldType:
		vI32, _ := strconv.ParseInt(v, 10, 32)
		return int(vI32)
	case PositiveBigIntegerFieldType:
		vI, _ := strconv.ParseUint(v, 10, 64)
		return vI
	case PositiveSmallIntegerFieldType:
		vI, _ := strconv.ParseUint(v, 10, 32)
		return uint(vI)
	case DecimalFieldType:
		vI, _ := strconv.ParseFloat(v, 64)
		return vI
	case FloatFieldType:
		vI, _ := strconv.ParseFloat(v, 64)
		return vI
	}
	return nil
}

type EmailWidget struct {
	Widget
}

func (w *EmailWidget) GetWidgetType() WidgetType {
	return EmailInputWidgetType
}

func (w *EmailWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/email"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *EmailWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("4", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *EmailWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	if !govalidator.IsEmail(v[0]) {
		return NewHTTPErrorResponse("should_be_email", "should be an email")
	}
	w.SetOutputValue(v[0])
	return nil
}

type URLWidget struct {
	Widget
	URLValid                 bool
	CurrentLabel             string
	Href                     string
	Value                    string
	ChangeLabel              string
	AppendHTTPSAutomatically bool
}

func (w *URLWidget) GetWidgetType() WidgetType {
	return URLInputWidgetType
}

func (w *URLWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/url"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *URLWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("5", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	data["UrlValid"] = w.URLValid
	if w.CurrentLabel == "" {
		data["CurrentLabel"] = "URL"
	} else {
		data["CurrentLabel"] = w.CurrentLabel
	}
	data["Href"] = w.Href
	data["Value"] = w.Widget.Value
	if w.ChangeLabel == "" {
		data["ChangeLabel"] = "Change"
	} else {
		data["ChangeLabel"] = w.ChangeLabel
	}
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *URLWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	url := v[0]
	if w.AppendHTTPSAutomatically {
		urlInitialRegex := regexp.MustCompile(`^http(s)?://.*`)
		if !urlInitialRegex.Match([]byte(v[0])) {
			url = "https://" + url
		}
	}
	if !govalidator.IsURL(url) {
		return NewHTTPErrorResponse("should_be_url", "should be an url")
	}
	w.SetOutputValue(v[0])
	return nil
}

type PasswordWidget struct {
	Widget
}

func (w *PasswordWidget) GetWidgetType() WidgetType {
	return PasswordInputWidgetType
}

func (w *PasswordWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/password"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *PasswordWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("6", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	data["DisplayName"] = w.FieldDisplayName
	data["Value"] = ""
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *PasswordWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !w.Required {
		w.SetOutputValue("")
		return nil
	}
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	if len(v[0]) < CurrentConfig.D.Auth.MinPasswordLength {
		return NewHTTPErrorResponse("password_length_error", "length of the password has to be at least %d symbols", strconv.Itoa(CurrentConfig.D.Auth.MinPasswordLength))
	}
	w.SetOutputValue(v[0])
	return nil
}

type HiddenWidget struct {
	Widget
}

func (w *HiddenWidget) GetWidgetType() WidgetType {
	return HiddenInputWidgetType
}

func (w *HiddenWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/hidden"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *HiddenWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("7", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *HiddenWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type DateWidget struct {
	Widget
	DateValue string
}

func (w *DateWidget) GetWidgetType() WidgetType {
	return DateInputWidgetType
}

func (w *DateWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/date"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *DateWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("8", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	if w.DateValue != "" {
		data["Value"] = w.DateValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *DateWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.DateValue = v[0]
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.GoMonolith.DateFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type DateTimeWidget struct {
	Widget
	DateTimeValue string
}

func (w *DateTimeWidget) GetWidgetType() WidgetType {
	return DateTimeInputWidgetType
}

func (w *DateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/datetime"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *DateTimeWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("9", w.FieldDisplayName)
	var value interface{}
	var valueStr string
	if w.Populate != nil {
		value = w.Populate(w, formRenderContext, currentField)
		valueStr = value.(string)
	} else {
		value = TransformDateTimeValueForWidget(w.Value)
		if value != nil {
			valueStr = template.HTMLEscapeString(value.(string))
		} else {
			valueStr = ""
		}
	}
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(), "Value": valueStr,
		"Name": w.GetHTMLInputName(), "FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
		"Required": w.Required, "HelpText": w.HelpText, "FormError": w.ValidationErrors,
		"FormErrorNotEmpty": len(w.ValidationErrors) > 0,
	}
	if w.DateTimeValue != "" {
		data["Value"] = w.DateTimeValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *DateTimeWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.DateTimeValue = v[0]
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.GoMonolith.DateTimeFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type TimeWidget struct {
	Widget
	TimeValue string
}

func (w *TimeWidget) GetWidgetType() WidgetType {
	return TimeInputWidgetType
}

func (w *TimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/time"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *TimeWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("10", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	if w.TimeValue != "" {
		data["Value"] = w.TimeValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *TimeWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.TimeValue = v[0]
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.GoMonolith.TimeFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type TextareaWidget struct {
	Widget
}

func (w *TextareaWidget) GetWidgetType() WidgetType {
	return TextareaInputWidgetType
}

func (w *TextareaWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/textarea"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *TextareaWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("11", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *TextareaWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type ArrayWidget struct {
	TextareaWidget
}

func (w *ArrayWidget) GetWidgetType() WidgetType {
	return ArrayWidgetType
}

func (w *ArrayWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/textarea"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ArrayWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("11", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	// spew.Dump("data debug", data, w.GetTemplateName())
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ArrayWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	if renderContext.Field.FieldType.Name() == "StringArray" {
		w.SetOutputValue(pq.StringArray(strings.Split(v[0], "\r\n")))
	} else {
		return NewHTTPErrorResponse("not_handled_type_of_array", "array of type %s not handled", renderContext.Field.FieldType.Name())
	}

	return nil
}

type JSONWidget struct {
	TextareaWidget
}

func (w *JSONWidget) GetWidgetType() WidgetType {
	return ArrayWidgetType
}

func (w *JSONWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/textarea"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *JSONWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("11", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *JSONWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue([]byte(v[0]))
	return nil
}


type CheckboxWidget struct {
	Widget
}

func (w *CheckboxWidget) GetWidgetType() WidgetType {
	return CheckboxInputWidgetType
}

func (w *CheckboxWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/checkbox"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *CheckboxWidget) SetValue(v interface{}) {
	v1 := TransformValueForOperator(v)
	w.Value = v1
}

func (w *CheckboxWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("12", w.FieldDisplayName)
	var value interface{}
	if w.Populate != nil {
		value = w.Populate(w, formRenderContext, currentField)
	} else {
		value = TransformValueForWidget(w.Value)
	}
	if valueS, ok := value.(string); ok {
		if valueS != "" && valueS != "false" {
			w.Attrs["checked"] = "checked"
		}
	}
	if valueB, ok := value.(bool); ok {
		if valueB {
			w.Attrs["checked"] = "checked"
		}
	}
	// w.Value = nil
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *CheckboxWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.Populate != nil {
		w.Populate(w, renderContext, nil)
	}
	if w.ReadOnly {
		return nil
	}
	_, ok := form.Value[w.GetHTMLInputName()]
	w.SetValue(ok == true)
	w.SetOutputValue(ok == true)
	return nil
}

type SelectOptGroup struct {
	OptLabel string
	Value    interface{}
	Selected bool
	Attrs    map[string]string
}

type SelectOptGroupStringified struct {
	OptLabel           string
	Value              string
	Selected           bool
	OptionTemplateName string
	Attrs              map[string]string
}

type SelectWidget struct {
	Widget
	OptGroups                map[string][]*SelectOptGroup
	DontValidateForExistence bool
	Multiple                 bool
}

func (w *SelectWidget) CloneAllOtherImportantSettings(widget IWidget) {
	widget1 := widget.(*SelectWidget)
	widget1.OptGroups = w.OptGroups
	widget1.Populate = w.Populate
	widget1.Multiple = w.Multiple
	widget.SetErrors(w.ValidationErrors)
}

func (w *SelectWidget) GetWidgetType() WidgetType {
	return SelectWidgetType
}

func (w *SelectWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/select"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectWidget) GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData {
	var value interface{}
	realV := make([]string, 0)
	if w.Populate != nil {
		value = w.Populate(w, formRenderContext, currentField)
	} else {
		value = TransformValueForWidget(w.Value)
	}
	if w.Multiple {
		realV = value.([]string)
	}
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/select.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			if w.Multiple {
				optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
					OptLabel:           optGroup.OptLabel,
					Value:              value1,
					Selected:           Contains(realV, value1),
					OptionTemplateName: optionTemplateName,
					Attrs:              make(map[string]string),
				})
			} else {
				optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
					OptLabel:           optGroup.OptLabel,
					Value:              value1,
					Selected:           value1 == value,
					OptionTemplateName: optionTemplateName,
					Attrs:              make(map[string]string),
				})
			}
		}
	}
	if w.Multiple {
		w.SetAttr("multiple", "true")
	} else {
		w.SetAttr("data-selected", value.(string))
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *SelectWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("13", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	if w.Populate != nil {
		w.Populate(w, renderContext, nil)
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	var notExistentValue string
	if !w.DontValidateForExistence {
		optValues := []string{}
		for _, optGroup := range w.OptGroups {
			for _, optGroupOption := range optGroup {
				optValues = append(optValues, optGroupOption.Value.(string))
			}
		}
		for _, v1 := range v {
			if !Contains(optValues, v1) {
				foundNotExistent = true
				notExistentValue = v1
				break
			}
		}
	}
	if w.Multiple {
		w.SetValue(v)
	} else {
		w.SetValue(v[0])
	}
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	if w.Multiple {
		w.SetOutputValue(v)
	} else {
		w.SetOutputValue(v[0])
	}
	return nil
}

type ForeignKeyWidget struct {
	Widget
	OptGroups                map[string][]*SelectOptGroup
	DontValidateForExistence bool
	AddNewLink               string
	GetQuerySet              func(formRenderContext *FormRenderContext) IPersistenceStorage
	GenerateModelInterface   func() (interface{}, interface{})
	Autocomplete             bool
}

func (w *ForeignKeyWidget) GetWidgetType() WidgetType {
	return ForeignKeyWidgetType
}

func (w *ForeignKeyWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/foreignkey"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ForeignKeyWidget) CloneAllOtherImportantSettings(widget IWidget) {
	widget1 := widget.(*ForeignKeyWidget)
	widget1.GetQuerySet = w.GetQuerySet
	widget1.AddNewLink = w.AddNewLink
	widget1.GenerateModelInterface = w.GenerateModelInterface
	widget1.Autocomplete = w.Autocomplete
	widget1.Populate = w.Populate
	widget.SetErrors(w.ValidationErrors)
}

func (w *ForeignKeyWidget) GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData {
	var value interface{}
	if w.Populate != nil {
		value = w.Populate(w, formRenderContext, currentField)
	} else {
		if !reflect.ValueOf(w.Value).IsZero() {
			value = GetID(reflect.Indirect(reflect.ValueOf(w.Value)))
		}
	}
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	autocompleteURL := ""
	editURL := ""
	valueID := ""
	valueHTML := ""
	valueS := TransformValueForListDisplay(value)
	if !w.Autocomplete {
		for optGroupName, optGroups := range w.OptGroups {
			optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
			for _, optGroup := range optGroups {
				value1 := TransformValueForWidget(optGroup.Value).(string)
				optionTemplateName := "widgets/select.option"
				if w.IsForAdmin {
					optionTemplateName = "admin/" + optionTemplateName
				}
				optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
					OptLabel:           optGroup.OptLabel,
					Value:              value1,
					Selected:           value1 == valueS,
					OptionTemplateName: optionTemplateName,
					Attrs:              make(map[string]string),
				})
			}
		}
		w.SetAttr("data-selected", TransformValueForListDisplay(value))
	} else {
		modelI, _ := w.GenerateModelInterface()

		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		adminPage := CurrentDashboardAdminPanel.AdminPages.GetByModelName(modelDescription.Statement.Schema.Name)
		autocompleteURL = adminPage.GenerateLinkForModelAutocompletion()
		valueID = ""
		if !reflect.ValueOf(w.Value).IsZero() {
			valueID = TransformValueForListDisplay(GetID(reflect.ValueOf(w.Value)))
			if valueID == "0" {
				valueID = ""
			}
		}
		if valueID != "" {
			queryset := w.GetDbHandler(formRenderContext)
			queryset.LoadDataForModelByID(modelI, TransformValueForListDisplay(GetID(reflect.ValueOf(w.Value))))
			valueHTML = modelI.(GoMonolithString).String()
			editURL = adminPage.GenerateLinkToEditModel(reflect.ValueOf(w.Value))
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
		"AutocompleteUrl": autocompleteURL,
		"ValueID":         valueID,
		"Value":           valueHTML,
		"EditUrl":         editURL,
	}
}

func (w *ForeignKeyWidget) GetDbHandler(formRenderContext *FormRenderContext) IPersistenceStorage {
	if w.GetQuerySet == nil {
		database := NewDatabaseInstance()
		return NewGormPersistenceStorage(database.Db)
	}
	return w.GetQuerySet(formRenderContext)
}
func (w *ForeignKeyWidget) BuildChoices(formRenderContext *FormRenderContext) {
	queryset := w.GetDbHandler(formRenderContext)
	w.OptGroups = make(map[string][]*SelectOptGroup)
	w.OptGroups[""] = make([]*SelectOptGroup, 0)
	_, models := w.GenerateModelInterface()
	queryset.Find(models)
	list := reflect.Indirect(reflect.ValueOf(models))
	for i := 0; i < list.Len(); i++ {
		modelStringified := reflect.ValueOf(list.Index(i).Interface()).MethodByName("String").Call([]reflect.Value{})[0].Interface()
		modelID := GetID(list.Index(i))
		w.OptGroups[""] = append(w.OptGroups[""], &SelectOptGroup{
			OptLabel: modelStringified.(string),
			Value:    modelID,
			Selected: modelID == w.Value,
		})
	}
}

func (w *ForeignKeyWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("13", w.FieldDisplayName)
	if !w.Autocomplete {
		w.BuildChoices(formRenderContext)
	}
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["AddNewLink"] = w.AddNewLink
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ForeignKeyWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v := make([]string, 0)
	ok := false
	if !w.Autocomplete {
		v, ok = form.Value[w.GetHTMLInputName()]
	} else {
		v, ok = form.Value[w.GetHTMLInputName()+"-value"]
	}
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	foundNotExistent := false
	var notExistentValue string
	if !w.Autocomplete {
		w.BuildChoices(renderContext)
		if !w.DontValidateForExistence {
			optValues := []string{}
			for _, optGroup := range w.OptGroups {
				for _, optGroupOption := range optGroup {
					optValues = append(optValues, TransformValueForListDisplay(optGroupOption.Value))
				}
			}
			for _, v1 := range v {
				if !Contains(optValues, v1) {
					foundNotExistent = true
					notExistentValue = v1
					break
				}
			}
		}
	}
	if v[0] == "" {
		return nil
	}
	var c int64
	queryset := w.GetDbHandler(renderContext)
	modelI, _ := w.GenerateModelInterface()
	queryset.Model(modelI).Count(&c)
	if c == 0 {
		return NewHTTPErrorResponse("object_not_found", "no object found to be used for this field")
	}
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	queryset.LoadDataForModelByID(modelI, v[0])
	w.SetOutputValue(modelI)
	return nil
}

type ContentTypeSelectorWidget struct {
	Widget
	OptGroups             map[string][]*SelectOptGroup
	LoadFieldsOfAllModels bool
}

func (w *ContentTypeSelectorWidget) GetWidgetType() WidgetType {
	return SelectWidgetType
}

func (w *ContentTypeSelectorWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "admin/widgets/contenttypeselector"
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ContentTypeSelectorWidget) PopulateSelectorOptions(afo IAdminFilterObjects) {
	w.OptGroups = make(map[string][]*SelectOptGroup)
	w.OptGroups[""] = make([]*SelectOptGroup, 0)
	var contentTypes []*ContentType
	if afo == nil {
		database := NewDatabaseInstance()
		database.Db.Find(&contentTypes)
	} else {
		afo.GetDatabase().Db.Find(&contentTypes)
	}
	defaultOption := &SelectOptGroup{
		OptLabel: "Choose content type",
		Value:    "",
	}
	w.OptGroups[""] = append(w.OptGroups[""], defaultOption)
	for _, contentType := range contentTypes {
		option := &SelectOptGroup{
			OptLabel: contentType.String(),
			Value:    strconv.Itoa(int(contentType.ID)),
			Attrs:    make(map[string]string),
		}
		option.Attrs["data-iden"] = fmt.Sprintf("%s:%s", contentType.BlueprintName, contentType.ModelName)
		w.OptGroups[""] = append(w.OptGroups[""], option)
	}
}

func (w *ContentTypeSelectorWidget) GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData {
	var value interface{}
	a := w.Value.(ContentType).ID
	value = strconv.Itoa(int(a))
	w.PopulateSelectorOptions(nil)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/select.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel:           optGroup.OptLabel,
				Value:              value1,
				Selected:           value1 == value,
				OptionTemplateName: optionTemplateName,
				Attrs:              optGroup.Attrs,
			})
		}
	}
	allFields := "{}"
	if w.LoadFieldsOfAllModels {
		database := NewDatabaseInstance()
		fD := make(map[string][]string)
		for blueprintRootAdminPage := range CurrentDashboardAdminPanel.AdminPages.GetAll() {
			for modelPage := range blueprintRootAdminPage.SubPages.GetAll() {
				if modelPage.Model == nil {
					continue
				}
				iden := fmt.Sprintf("%s:%s", modelPage.BlueprintName, modelPage.ModelName)
				fD[iden] = make([]string, 0)
				statement := &gorm.Statement{DB: database.Db}
				statement.Parse(modelPage.Model)
				for _, field := range statement.Schema.Fields {
					fD[iden] = append(fD[iden], field.Name)
				}
			}
			allFieldsB, _ := json.Marshal(fD)
			allFields = string(allFieldsB)
		}
		database.Close()
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
		"AllFields": allFields,
	}
}

func (w *ContentTypeSelectorWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("13", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ContentTypeSelectorWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	if v[0] == "" {
		return nil
	}
	var contentType = ContentType{}
	afo.GetDatabase().Db.First(&contentType, v[0])
	w.SetValue(contentType)
	w.SetOutputValue(contentType)
	return nil
}

type NullBooleanWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}

func (w *NullBooleanWidget) GetWidgetType() WidgetType {
	return NullBooleanWidgetType
}

func (w *NullBooleanWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/nullboolean"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *NullBooleanWidget) GetOptGroups() map[string][]*SelectOptGroup {
	if w.OptGroups == nil {
		defaultOptGroups := make(map[string][]*SelectOptGroup)
		defaultOptGroups[""] = make([]*SelectOptGroup, 0)
		defaultOptGroups[""] = append(defaultOptGroups[""], &SelectOptGroup{
			OptLabel: "Yes",
			Value:    "yes",
		})
		defaultOptGroups[""] = append(defaultOptGroups[""], &SelectOptGroup{
			OptLabel: "No",
			Value:    "no",
		})
		return defaultOptGroups
	}
	return w.OptGroups
}

func (w *NullBooleanWidget) GetDataForRendering(formRenderContext *FormRenderContext, currentField *Field) WidgetData {
	value := TransformValueForWidget(w.Value)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.GetOptGroups() {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/select.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel:           optGroup.OptLabel,
				Value:              value1,
				Selected:           value1 == value,
				OptionTemplateName: optionTemplateName,
				Attrs:              make(map[string]string),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *NullBooleanWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("14", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext, currentField)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *NullBooleanWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v[0])
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type SelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}

func (w *SelectMultipleWidget) GetWidgetType() WidgetType {
	return SelectMultipleWidgetType
}

func (w *SelectMultipleWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/select"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectMultipleWidget) GetDataForRendering(formRenderContext *FormRenderContext) WidgetData {
	w.Attrs["multiple"] = "true"
	value := TransformValueForWidget(w.Value).([]string)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/select.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel:           optGroup.OptLabel,
				Value:              value1,
				Selected:           Contains(value, value1),
				OptionTemplateName: optionTemplateName,
				Attrs:              make(map[string]string),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *SelectMultipleWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("15", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectMultipleWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v)
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v)
	return nil
}

type RadioOptGroup struct {
	OptLabel string
	Value    interface{}
	Selected bool
	Label    string
}

type RadioOptGroupStringified struct {
	OptLabel           string
	Value              string
	Selected           bool
	OptionTemplateName string
	WrapLabel          bool
	ForID              string
	Label              string
	Type               string
	Name               string
	Attrs              map[string]string
	FieldDisplayName   string
	ReadOnly           bool
}

type RadioSelectWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	ID        string
	WrapLabel bool
}

func (w *RadioSelectWidget) GetWidgetType() WidgetType {
	return RadioWidgetType
}

func (w *RadioSelectWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/radioselect"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *RadioSelectWidget) SetID(ID string) {
	w.ID = ID
}

func (w *RadioSelectWidget) GetDataForRendering(formRenderContext *FormRenderContext) WidgetData {
	value := TransformValueForWidget(w.Value).(string)
	optGroupSstringified := make(map[string][]*RadioOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*RadioOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/radio.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &RadioOptGroupStringified{
				OptLabel:           optGroup.OptLabel,
				Value:              value1,
				Selected:           value == value1,
				OptionTemplateName: optionTemplateName,
				Label:              optGroup.Label,
				WrapLabel:          w.WrapLabel,
				ForID:              w.ID,
				Type:               "radio",
				Name:               w.GetHTMLInputName(),
				Attrs:              w.Widget.GetAttrs(),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified, "Id": w.ID,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *RadioSelectWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("16", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *RadioSelectWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v[0])
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type CheckboxSelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	ID        string
	WrapLabel bool
}

func (w *CheckboxSelectMultipleWidget) GetWidgetType() WidgetType {
	return CheckboxSelectMultipleWidgetType
}

func (w *CheckboxSelectMultipleWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/checkboxselectmultiple"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *CheckboxSelectMultipleWidget) SetID(ID string) {
	w.ID = ID
}

func (w *CheckboxSelectMultipleWidget) GetDataForRendering(formRenderContext *FormRenderContext) WidgetData {
	value := TransformValueForWidget(w.Value).([]string)
	optGroupSstringified := make(map[string][]*RadioOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*RadioOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optionTemplateName := "widgets/checkbox.option"
			if w.IsForAdmin {
				optionTemplateName = "admin/" + optionTemplateName
			}
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &RadioOptGroupStringified{
				OptLabel:           optGroup.OptLabel,
				Value:              value1,
				Selected:           Contains(value, value1),
				OptionTemplateName: optionTemplateName,
				Label:              optGroup.Label,
				WrapLabel:          w.WrapLabel,
				ForID:              w.ID,
				Type:               "checkbox",
				Name:               w.GetHTMLInputName(),
				Attrs:              w.Widget.GetAttrs(),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(), "OptGroups": optGroupSstringified, "Id": w.ID,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *CheckboxSelectMultipleWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("17", w.FieldDisplayName)
	data := w.GetDataForRendering(formRenderContext)
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *CheckboxSelectMultipleWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v)
	if foundNotExistent {
		return NewHTTPErrorResponse("value_invalid_for_field", "not valid value %s for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v)
	return nil
}

type FileWidget struct {
	Widget
	Storage    IStorageInterface
	UploadPath string
	Multiple   bool
}

func (w *FileWidget) GetWidgetType() WidgetType {
	return FileInputWidgetType
}

func (w *FileWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/file"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *FileWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("18", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	vI := reflect.ValueOf(w.Value)
	if w.Value != nil && !vI.IsZero() {
		data["UploadedFile"] = storage.GetUploadURL() + "/" + w.Value.(string)
		data["IsItImage"] = strings.Contains(w.Attrs["accept"], "image/")
	}
	data["Value"] = w.Value
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *FileWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	files := form.File[w.GetHTMLInputName()]
	if len(files) == 0 {
		return nil
	}
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	ret := make([]string, 0)
	var filename string
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return err
		}
		bytecontent := make([]byte, file.Size)
		_, err = f.Read(bytecontent)
		if err != nil {
			return err
		}
		filename, err = storage.Save(&FileForStorage{
			Content:           bytecontent,
			PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
			Filename:          file.Filename,
		})
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		ret = append(ret, filename)
	}
	if w.Multiple {
		w.SetOutputValue(ret)
	} else if len(ret) > 0 {
		w.SetOutputValue(ret[0])
	} else {
		w.SetOutputValue("")
	}
	return nil
}

type URLValue struct {
	URL string
}

type ClearableFileWidget struct {
	Widget
	InitialText        string
	CurrentValue       *URLValue
	Required           bool
	ID                 string
	ClearCheckboxLabel string
	InputText          string
	Storage            IStorageInterface
	UploadPath         string
	Multiple           bool
}

func (w *ClearableFileWidget) GetWidgetType() WidgetType {
	return FileInputWidgetType
}

func (w *ClearableFileWidget) SetID(ID string) {
	w.ID = ID
}

func (w *ClearableFileWidget) IsInitial() bool {
	return w.CurrentValue == nil
}

func (w *ClearableFileWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/clearablefile"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ClearableFileWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("19", w.FieldDisplayName)
	data := w.Widget.GetDataForRendering(formRenderContext, currentField)
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	vI := reflect.ValueOf(w.Value)
	if w.Value != nil && !vI.IsZero() {
		data["UploadedFile"] = storage.GetUploadURL() + w.Value.(string)
		data["IsItImage"] = strings.Contains(w.Attrs["accept"], "image/")
	}
	data["Value"] = w.Value
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	data["IsInitial"] = w.IsInitial()
	data["InitialText"] = w.InitialText
	data["CurrentValue"] = w.CurrentValue
	data["Required"] = w.Required
	data["Id"] = w.ID
	data["ClearCheckboxLabel"] = w.ClearCheckboxLabel
	data["InputText"] = w.InputText
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ClearableFileWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	files := form.File[w.GetHTMLInputName()]
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	ret := make([]string, 0)
	var err error
	var filename string
	for _, file := range files {
		f, _ := file.Open()
		bytecontent := make([]byte, file.Size)
		_, err = f.Read(bytecontent)
		filename, err = storage.Save(&FileForStorage{
			Content:           bytecontent,
			PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
			Filename:          file.Filename,
		})
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		ret = append(ret, filename)
	}
	if w.Multiple {
		w.SetOutputValue(ret)
	} else if len(ret) > 0 {
		w.SetOutputValue(ret[0])
	} else {
		w.SetOutputValue("")
		w.SetValue("")
	}
	return nil
}

type MultipleInputHiddenWidget struct {
	Widget
}

func (w *MultipleInputHiddenWidget) GetWidgetType() WidgetType {
	return MultipleHiddenInputWidgetType
}

func (w *MultipleInputHiddenWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/multipleinputhidden"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *MultipleInputHiddenWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("20", w.FieldDisplayName)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(),
	}
	data["Required"] = w.Required
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["FormError"] = w.ValidationErrors
	data["FormErrorNotEmpty"] = len(w.ValidationErrors) > 0
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).([]string)
	for _, v := range value {
		w1 := HiddenWidget{}
		w1.Name = w.GetHTMLInputName()
		w1.SetValue(v)
		w1.Attrs = make(map[string]string)
		for attrName, attrValue := range w.Attrs {
			w1.Attrs[attrName] = attrValue
		}
		vd := w1.GetDataForRendering(formRenderContext, currentField)
		vd["Type"] = w1.GetWidgetType()
		templateName := "widgets/hidden"
		if w.IsForAdmin {
			templateName = "admin/widgets/hidden"
		}
		vd["TemplateName"] = templateName
		subwidgets = append(subwidgets, vd)
	}
	data["Subwidgets"] = subwidgets
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *MultipleInputHiddenWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v)
	w.SetOutputValue(v)
	return nil
}

type ChooseFromSelectWidget struct {
	Widget
	PopulateLeftSide      func() []*SelectOptGroup
	PopulateRightSide     func() []*SelectOptGroup
	LeftSelectTitle       string
	LeftSelectHelp        string
	LeftHelpChooseAll     string
	LeftSearchSelectHelp  string
	LeftChooseAllText     string
	RightSelectTitle      string
	RightSelectHelp       string
	RightHelpChooseAll    string
	RightSearchSelectHelp string
	RightChooseAllText    string
	AddNewLink            string
	AddNewTitle           string
}

func (w *ChooseFromSelectWidget) GetWidgetType() WidgetType {
	return ChooseFromSelectWidgetType
}

func (w *ChooseFromSelectWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/choosefromselect"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ChooseFromSelectWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("21", w.FieldDisplayName)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(),
	}
	tmpOptions := w.PopulateLeftSide()
	var leftSideOptions []*SelectOptGroup
	rightSideOptions := w.PopulateRightSide()
	for _, option := range tmpOptions {
		found := false
		for _, option1 := range rightSideOptions {
			if option1.Value.(uint) == option.Value.(uint) {
				found = true
			}
		}
		if !found {
			leftSideOptions = append(leftSideOptions, option)
		}
	}
	data["HelpText"] = w.HelpText
	ValueIds := make([]string, 0)
	for _, option := range rightSideOptions {
		ValueIds = append(ValueIds, strconv.Itoa(int(option.Value.(uint))))
	}
	w.Value = strings.Join(ValueIds, ",")
	data["Value"] = w.Value
	data["AddNewLink"] = w.AddNewLink
	data["AddNewTitle"] = w.AddNewTitle
	data["FieldDisplayName"] = w.FieldDisplayName
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["FormError"] = w.ValidationErrors
	data["FormErrorNotEmpty"] = len(w.ValidationErrors) > 0
	subwidgets := make([]WidgetData, 0)
	w1 := SelectWidget{}
	w1.OptGroups = make(map[string][]*SelectOptGroup)
	w1.OptGroups[""] = leftSideOptions
	w1.Name = w.GetHTMLInputName() + "_left"
	w1.Attrs = w.GetAttrs()
	vd := w1.GetDataForRendering(formRenderContext, currentField)
	vd["Type"] = "select"
	vd["ShowOnlyHtmlInput"] = true
	vd["GenerateSelector"] = true
	vd["Multiple"] = true
	vd["SelectClass"] = "available-select"
	vd["SelectorClass"] = "left-selector"
	vd["SelectTitle"] = w.LeftSelectTitle
	vd["SelectHelp"] = w.LeftSelectHelp
	vd["SearchSelectHelp"] = w.LeftSearchSelectHelp
	vd["HelpChooseAll"] = w.LeftHelpChooseAll
	vd["ChooseAllText"] = w.LeftChooseAllText
	if vd["ChooseAllText"] == "" {
		vd["ChooseAllText"] = "Choose all"
	}
	vd["ChooseAllIsActive"] = len(w1.OptGroups[""]) > 0
	vd["SelectorGeneralClass"] = "selector-available"
	vd["ChooseAllClass"] = "selector-chooseall"
	templateName := "widgets/selectwithsearch"
	if w.IsForAdmin {
		templateName = "admin/widgets/selectwithsearch"
	}
	vd["TemplateName"] = templateName
	subwidgets = append(subwidgets, vd)
	w2 := SelectWidget{}
	w2.OptGroups = make(map[string][]*SelectOptGroup)
	w2.OptGroups[""] = rightSideOptions
	w2.Name = w.GetHTMLInputName() + "_right"
	w2.Attrs = w.GetAttrs()
	vd2 := w2.GetDataForRendering(formRenderContext, currentField)
	vd2["ShowOnlyHtmlInput"] = true
	vd2["Type"] = "select"
	vd2["TemplateName"] = templateName
	vd2["GenerateSelector"] = false
	vd2["Multiple"] = true
	vd2["SelectorClass"] = "right-selector"
	vd2["SelectTitle"] = w.RightSelectTitle
	vd2["SelectHelp"] = w.RightSelectHelp
	vd2["SearchSelectHelp"] = w.RightSearchSelectHelp
	vd2["HelpChooseAll"] = w.RightHelpChooseAll
	vd2["SelectClass"] = "chosen-select"
	vd2["ChooseAllText"] = w.RightChooseAllText
	vd2["ChooseAllClass"] = "selector-clearall"
	if vd2["ChooseAllText"] == "" {
		vd2["ChooseAllText"] = "Remove all"
	}
	vd2["ChooseAllIsActive"] = len(w2.OptGroups[""]) > 0
	vd2["SelectorGeneralClass"] = "selector-chosen related-target"
	subwidgets = append(subwidgets, vd2)
	data["Subwidgets"] = subwidgets
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ChooseFromSelectWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.GetHTMLInputName()]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetOutputValue(strings.Split(v[0], ","))
	return nil
}

type SplitDateTimeWidget struct {
	Widget
	DateAttrs  map[string]string
	TimeAttrs  map[string]string
	DateFormat string
	TimeFormat string
	DateLabel  string
	TimeLabel  string
	DateValue  string
	TimeValue  string
}

func (w *SplitDateTimeWidget) GetWidgetType() WidgetType {
	return SplitDateTimeWidgetType
}

func (w *SplitDateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/splitdatetime"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SplitDateTimeWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("23", w.FieldDisplayName)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(),
	}
	data["FormError"] = w.ValidationErrors
	data["FormErrorNotEmpty"] = len(w.ValidationErrors) > 0
	data["Required"] = w.Required
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	if w.DateLabel == "" {
		data["DateLabel"] = "Date:"
	} else {
		data["DateLabel"] = w.DateLabel
	}
	if w.TimeLabel == "" {
		data["TimeLabel"] = "Time:"
	} else {
		data["TimeLabel"] = w.TimeLabel
	}
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).(*time.Time)
	w1 := DateWidget{}
	w1.Name = w.GetHTMLInputName() + "_date"
	if w.DateValue != "" {
		w1.SetValue(w.DateValue)
	} else {
		w1.SetValue(value.Format(w.DateFormat))
	}
	w1.Attrs = w.DateAttrs
	vd := w1.Widget.GetDataForRendering(formRenderContext, currentField)
	vd["Type"] = w1.GetWidgetType()
	templateName := "widgets/date"
	if w.IsForAdmin {
		templateName = "admin/widgets/date"
	}
	vd["TemplateName"] = templateName
	subwidgets = append(subwidgets, vd)
	w2 := TimeWidget{}
	w2.Name = w.GetHTMLInputName() + "_time"
	if w.TimeValue != "" {
		w2.SetValue(w.TimeValue)
	} else {
		w2.SetValue(value.Format(w.TimeFormat))
	}
	w2.Attrs = w.TimeAttrs
	vd1 := w2.Widget.GetDataForRendering(formRenderContext, currentField)
	vd1["Type"] = w2.GetWidgetType()
	templateName = "widgets/time"
	if w.IsForAdmin {
		templateName = "admin/widgets/time"
	}
	vd1["TemplateName"] = templateName
	subwidgets = append(subwidgets, vd1)
	data["Subwidgets"] = subwidgets
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SplitDateTimeWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	vDate, ok := form.Value[w.GetHTMLInputName()+"_date"]
	if !ok {
		return NewHTTPErrorResponse("no_date", "no date has been submitted for field %s", w.FieldDisplayName)
	}
	w.DateValue = vDate[0]
	vTime, ok := form.Value[w.GetHTMLInputName()+"_time"]
	if !ok {
		return NewHTTPErrorResponse("no_time", "no time has been submitted for field %s", w.FieldDisplayName)
	}
	w.TimeValue = vTime[0]
	if w.Required && (vDate[0] == "" || vTime[0] == "") {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(w.DateFormat, vDate[0])
	if err != nil {
		return err
	}
	t, err := time.Parse(w.TimeFormat, vTime[0])
	if err != nil {
		return err
	}
	newT := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	w.SetOutputValue(&newT)
	return nil
}

type SplitHiddenDateTimeWidget struct {
	Widget
	DateAttrs  map[string]string
	TimeAttrs  map[string]string
	DateFormat string
	TimeFormat string
	DateValue  string
	TimeValue  string
}

func (w *SplitHiddenDateTimeWidget) GetWidgetType() WidgetType {
	return SplitHiddenDateTimeWidgetType
}

func (w *SplitHiddenDateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/splithiddendatetime"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SplitHiddenDateTimeWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("24", w.FieldDisplayName)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(),
	}
	data["FormError"] = w.ValidationErrors
	data["FormErrorNotEmpty"] = len(w.ValidationErrors) > 0
	data["Required"] = w.Required
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).(*time.Time)
	w1 := DateWidget{}
	w1.Name = w.GetHTMLInputName() + "_date"
	if w.DateValue != "" {
		w1.SetValue(w.DateValue)
	} else {
		w1.SetValue(value.Format(w.DateFormat))
	}
	w1.Attrs = w.DateAttrs
	vd := w1.Widget.GetDataForRendering(formRenderContext, currentField)
	vd["Type"] = "hidden"
	templateName := "widgets/date"
	if w.IsForAdmin {
		templateName = "admin/widgets/date"
	}
	vd["TemplateName"] = templateName
	subwidgets = append(subwidgets, vd)
	w2 := TimeWidget{}
	w2.Name = w.GetHTMLInputName() + "_time"
	if w.TimeValue != "" {
		w2.SetValue(w.TimeValue)
	} else {
		w2.SetValue(value.Format(w.TimeFormat))
	}
	w2.Attrs = w.TimeAttrs
	vd1 := w2.Widget.GetDataForRendering(formRenderContext, currentField)
	vd1["Type"] = "hidden"
	templateName = "widgets/time"
	if w.IsForAdmin {
		templateName = "admin/widgets/time"
	}
	vd1["TemplateName"] = templateName
	subwidgets = append(subwidgets, vd1)
	data["Subwidgets"] = subwidgets
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SplitHiddenDateTimeWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	vDate, ok := form.Value[w.GetHTMLInputName()+"_date"]
	if !ok {
		return NewHTTPErrorResponse("no_date", "no date has been submitted for field %s", w.FieldDisplayName)
	}
	w.DateValue = vDate[0]
	vTime, ok := form.Value[w.GetHTMLInputName()+"_time"]
	if !ok {
		return NewHTTPErrorResponse("no_time", "no time has been submitted for field %s", w.FieldDisplayName)
	}
	w.TimeValue = vTime[0]
	if w.Required && (vDate[0] == "" || vTime[0] == "") {
		return NewHTTPErrorResponse("field_required", "field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(w.DateFormat, vDate[0])
	if err != nil {
		return err
	}
	t, err := time.Parse(w.TimeFormat, vTime[0])
	if err != nil {
		return err
	}
	newT := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	w.SetOutputValue(&newT)
	return nil
}

type SelectDateWidget struct {
	Widget
	Years            []int
	Months           []*SelectOptGroup
	EmptyLabel       []*SelectOptGroup
	EmptyLabelString string
	YearValue        string
	MonthValue       string
	DayValue         string
}

func (w *SelectDateWidget) GetWidgetType() WidgetType {
	return SelectDateWidgetType
}

func (w *SelectDateWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		path := "widgets/selectdate"
		if w.IsForAdmin {
			path = "admin/" + path
		}
		return CurrentConfig.GetPathToTemplate(path)
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectDateWidget) Render(formRenderContext *FormRenderContext, currentField *Field) template.HTML {
	// spew.Dump("25", w.FieldDisplayName)
	value := TransformValueForWidget(w.Value).(*time.Time)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name":  w.GetHTMLInputName(),
	}
	data["FormError"] = w.ValidationErrors
	data["FormErrorNotEmpty"] = len(w.ValidationErrors) > 0
	data["Required"] = w.Required
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHTMLInput
	data["Type"] = w.GetWidgetType()
	dateParts := []string{}
	for _, formatChar := range CurrentConfig.D.GoMonolith.DateFormatOrder {
		if formatChar == 'y' {
			if Contains(dateParts, "year") {
				continue
			}
			dateParts = append(dateParts, "year")
		} else if formatChar == 'd' {
			if Contains(dateParts, "day") {
				continue
			}
			dateParts = append(dateParts, "day")
		} else if formatChar == 'm' {
			if Contains(dateParts, "month") {
				continue
			}
			dateParts = append(dateParts, "month")
		}
	}
	if w.Years == nil {
		initialYear := time.Now().Year()
		w.Years = make([]int, 0)
		for i := initialYear; i <= initialYear+10; i++ {
			w.Years = append(w.Years, i)
		}
	}
	var yearNoneValue *SelectOptGroup
	var monthNoneValue *SelectOptGroup
	var dayNoneValue *SelectOptGroup
	if w.EmptyLabel == nil {
		noneValue := &SelectOptGroup{
			OptLabel: w.EmptyLabelString,
			Value:    "",
		}
		dayNoneValue = noneValue
		yearNoneValue = noneValue
		monthNoneValue = noneValue
	} else {
		if len(w.EmptyLabel) != 3 {
			panic("empty_label slice must have 3 elements.")
		}
		dayNoneValue = w.EmptyLabel[2]
		yearNoneValue = w.EmptyLabel[0]
		monthNoneValue = w.EmptyLabel[1]
	}
	if w.Months == nil {
		w.Months = MakeMonthsSelect()
		if !w.Required {
			w.Months = append(w.Months, monthNoneValue)
			copy(w.Months[1:], w.Months)
			w.Months[0] = monthNoneValue
		}
	}
	var yearChoices []*SelectOptGroup
	if !w.Required {
		yearChoices = append(yearChoices, yearNoneValue)
	}
	for _, year := range w.Years {
		yearChoices = append(yearChoices, &SelectOptGroup{
			OptLabel: strconv.Itoa(year),
			Value:    strconv.Itoa(year),
		})
	}
	var dayChoices []*SelectOptGroup
	for i := 1; i < 32; i++ {
		dayChoices = append(dayChoices, &SelectOptGroup{
			OptLabel: strconv.Itoa(i),
			Value:    strconv.Itoa(i),
		})
		if !w.Required {
			dayChoices = append(dayChoices, dayNoneValue)
			copy(dayChoices[1:], dayChoices)
			dayChoices[0] = dayNoneValue
		}
	}
	subwidgets := make([]WidgetData, 0)
	w1 := SelectWidget{}
	w1.OptGroups = make(map[string][]*SelectOptGroup)
	w1.OptGroups[""] = yearChoices
	w1.Name = w.GetHTMLInputName() + "_year"
	if w.YearValue != "" {
		w1.SetValue(w.YearValue)
	} else {
		w1.SetValue(value.Year())
	}
	w1.Attrs = w.GetAttrs()
	vd := w1.GetDataForRendering(formRenderContext, currentField)
	vd["Type"] = "select"
	templateName := "widgets/select"
	if w.IsForAdmin {
		templateName = "admin/widgets/select"
	}
	vd["TemplateName"] = templateName
	yearWd := vd
	w2 := SelectWidget{}
	w2.OptGroups = make(map[string][]*SelectOptGroup)
	w2.OptGroups[""] = w.Months
	w2.Name = w.GetHTMLInputName() + "_month"
	if w.YearValue != "" {
		w2.SetValue(w.MonthValue)
	} else {
		w2.SetValue(value.Month())
	}
	w2.Attrs = w.GetAttrs()
	vd2 := w2.GetDataForRendering(formRenderContext, currentField)
	vd2["Type"] = "select"
	vd2["TemplateName"] = templateName
	w3 := SelectWidget{}
	w3.OptGroups = make(map[string][]*SelectOptGroup)
	w3.OptGroups[""] = dayChoices
	w3.Name = w.GetHTMLInputName() + "_day"
	if w.DayValue != "" {
		w3.SetValue(w.DayValue)
	} else {
		w3.SetValue(value.Day())
	}
	w3.Attrs = w.GetAttrs()
	vd3 := w3.GetDataForRendering(formRenderContext, currentField)
	vd3["Type"] = "select"
	vd3["TemplateName"] = templateName
	dayWd := vd3
	monthWd := vd2
	for _, datePart := range dateParts {
		if datePart == "year" {
			subwidgets = append(subwidgets, yearWd)
		} else if datePart == "month" {
			subwidgets = append(subwidgets, monthWd)
		} else if datePart == "day" {
			subwidgets = append(subwidgets, dayWd)
		}
	}
	data["Subwidgets"] = subwidgets
	return RenderWidget(formRenderContext, w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectDateWidget) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) error {
	if w.ReadOnly {
		return nil
	}
	vYear, ok := form.Value[w.GetHTMLInputName()+"_year"]
	if !ok {
		return NewHTTPErrorResponse("no_year", "no year has been submitted for field %s", w.FieldDisplayName)
	}
	w.YearValue = vYear[0]
	vMonth, ok := form.Value[w.GetHTMLInputName()+"_month"]
	if !ok {
		return NewHTTPErrorResponse("no_month", "no month has been submitted for field %s", w.FieldDisplayName)
	}
	w.MonthValue = vMonth[0]
	vDay, ok := form.Value[w.GetHTMLInputName()+"_day"]
	if !ok {
		return NewHTTPErrorResponse("no_day", "no day has been submitted for field %s", w.FieldDisplayName)
	}
	w.DayValue = vDay[0]
	if w.Required && (w.YearValue == "" || w.MonthValue == "" || w.DayValue == "") {
		return NewHTTPErrorResponse("either_year_month_value_empty", "either year, month, value is empty")
	}
	day, err := strconv.Atoi(w.DayValue)
	if err != nil {
		return NewHTTPErrorResponse("incorrect_day", "incorrect day")
	}
	month, err := strconv.Atoi(w.MonthValue)
	if err != nil {
		return NewHTTPErrorResponse("incorrect_month", "incorrect month")
	}
	year, err := strconv.Atoi(w.YearValue)
	if err != nil {
		return NewHTTPErrorResponse("incorrect_year", "incorrect year")
	}
	d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	w.SetOutputValue(&d)
	return nil
}

func MakeMonthsSelect() []*SelectOptGroup {
	ret := make([]*SelectOptGroup, 0)
	ret = append(ret, &SelectOptGroup{
		Value:    "1",
		OptLabel: "January",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "2",
		OptLabel: "February",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "3",
		OptLabel: "March",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "4",
		OptLabel: "April",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "5",
		OptLabel: "May",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "6",
		OptLabel: "June",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "7",
		OptLabel: "July",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "8",
		OptLabel: "August",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "9",
		OptLabel: "September",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "10",
		OptLabel: "October",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "11",
		OptLabel: "November",
	})
	ret = append(ret, &SelectOptGroup{
		Value:    "12",
		OptLabel: "December",
	})
	return ret
}

func GetWidgetFromFieldTypeAndGormField(fieldType GoMonolithFieldType, gormField *schema.Field) IWidget {
	var widget IWidget
	switch fieldType {
	case "biginteger":
		widget = &NumberWidget{
			NumberType: BigIntegerFieldType,
		}
	case "integer":
		widget = &NumberWidget{
			NumberType: IntegerFieldType,
		}
	case "positivebiginteger":
		widget = &NumberWidget{
			NumberType: PositiveBigIntegerFieldType,
		}
	case "positiveinteger":
		widget = &NumberWidget{
			NumberType: PositiveIntegerFieldType,
		}
	case "positivesmallinteger":
		widget = &NumberWidget{
			NumberType: PositiveSmallIntegerFieldType,
		}
	case "smallinteger":
		widget = &NumberWidget{
			NumberType: SmallIntegerFieldType,
		}
	case "binary":
		widget = &TextareaWidget{}
	case "char":
		widget = &TextWidget{}
		widget.SetAttr("maxlength", "1")
	case "boolean":
		widget = &CheckboxWidget{}
	case "decimal":
		widget = &NumberWidget{
			NumberType: DecimalFieldType,
		}
		widget.SetAttr("step", "0.1")
	case "float":
		widget = &NumberWidget{
			NumberType: FloatFieldType,
		}
		widget.SetAttr("step", "0.1")
	case "email":
		widget = &EmailWidget{}
	case "file":
		widget = &FileWidget{}
	case "filepath":
		widget = &TextWidget{}
	case "text":
		widget = &TextWidget{}
	case "time":
		widget = &TimeWidget{}
	case "nullboolean":
		widget = &NullBooleanWidget{}
	case "slug":
		widget = &TextWidget{}
	case "url":
		widget = &URLWidget{}
	case "uuid":
		widget = &TextWidget{}
	case "date":
		widget = &DateWidget{}
	case "datetime":
		widget = &DateTimeWidget{}
	case "duration":
		widget = &TextWidget{}
	case "foreignkey":
		// @todo, integrate autocomplate widget
		widget = &TextWidget{}
	case "array":
		widget = &ArrayWidget{}
	case "json":
		widget = &JSONWidget{}
	case "imagefield":
		widget = &FileWidget{}
		widget.SetAttr("accept", "image/*")
	case "ipaddress":
		widget = &TextWidget{}
		widget.SetAttr("minlength", "7")
		widget.SetAttr("maxlength", "15")
		widget.SetAttr("size", "15")
		widget.SetAttr("pattern", "^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$")
	case "genericipaddress":
		widget = &TextWidget{}
		widget.SetAttr("minlength", "7")
		widget.SetAttr("maxlength", "15")
		widget.SetAttr("size", "15")
		widget.SetAttr("pattern", "^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$")
		// @todo, make sure we handle many to many widget type
		// const ManyToManyUadminFieldType GoMonolithFieldType = "manytomany"
	default:
		widget = &TextWidget{}
	}
	widget.InitializeAttrs()
	widget.SetBaseFuncMap(FuncMap)
	widget.SetName(gormField.Name)
	widget.SetValue(gormField.DefaultValueInterface)
	return widget
}
