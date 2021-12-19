# Widgets

Widgets are responsible to render field for the form/read data from  the POST data.  
Widget has to implement following interface:
```go
type WidgetData map[string]interface{}
type IWidget interface {
	IDForLabel(model interface{}, F *Field) string
	GetWidgetType() WidgetType
	GetAttrs() map[string]string
	GetTemplateName() string
	SetTemplateName(templateName string)
	RenderUsingRenderer(renderer ITemplateRenderer)
	// GetValue(v interface{}, model interface{}) interface{}
	Render(formRenderContext *FormRenderContext, currentField *Field) string
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
	SetPopulate(func(renderContext *FormRenderContext, currentField *Field) interface{})
	SetPrefix(prefix string)
	GetHTMLInputName() string
	GetPopulate() func(renderContext *FormRenderContext, currentField *Field) interface{}
	IsReadOnly() bool
	IsValueConfigured() bool
	SetValueConfigured()
	GetRenderer() ITemplateRenderer
	GetFieldDisplayName() string
	GetName() string
}
```
In most cases it's ok to extend Widget structure for your own widget, like here:
```go
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

func (tw *TextWidget) Render(formRenderContext *FormRenderContext, currentField *Field) string {
	// spew.Dump("2", tw.FieldDisplayName)
	data := tw.Widget.GetDataForRendering(formRenderContext, currentField)
	data["Type"] = tw.GetWidgetType()
	data["ShowOnlyHtmlInput"] = tw.ShowOnlyHTMLInput
	return RenderWidget(tw.Renderer, tw.GetTemplateName(), data, tw.BaseFuncMap) // tw.Value, tw.Widget.GetAttrs()
}
```
Currently we support following widgets:
```go
// widget where widget could be determined dynamically.
type DynamicWidget struct {
	Widget
	GetRealWidget                  func(formRenderContext *FormRenderContext, currentField *Field) IWidget
	GetRealWidgetForFormProceeding func(form *multipart.Form, afo IAdminFilterObjects) IWidget
}
// text widget
type TextWidget struct {
	Widget
}
// foreign key link widget
type FkLinkWidget struct {
	Widget
}
// widget to display field as number
type NumberWidget struct {
	Widget
	NumberType GoMonolithFieldType
}
// widget to display field as email
type EmailWidget struct {
	Widget
}
// widget to display field as url
type URLWidget struct {

}
// widget to display field as password
type PasswordWidget struct {
	Widget
}
// widget to display field as hidden input
type HiddenWidget struct {
	Widget
}
// widget to display field as date, currently we don't support date picker
type DateWidget struct {
	Widget
	DateValue string
}
// widget to display field as datetime, currently we don't support date picker
type DateTimeWidget struct {
	Widget
	DateTimeValue string
}
// widget to display field as time, currently we don't support time picker
type TimeWidget struct {
	Widget
	TimeValue string
}
// widget to display field as textarea
type TextareaWidget struct {
	Widget
}
// widget to display field as checkbox
type CheckboxWidget struct {
	Widget
}
// widget to display field as select
type SelectWidget struct {
	Widget
	OptGroups                map[string][]*SelectOptGroup
	DontValidateForExistence bool
}
// widget to display field as foreign key, it fetches automaitcally all options to be choosed from
type ForeignKeyWidget struct {
	Widget
	OptGroups                map[string][]*SelectOptGroup
	DontValidateForExistence bool
	AddNewLink               string
	GetQuerySet              func(formRenderContext *FormRenderContext) IPersistenceStorage
	GenerateModelInterface   func() (interface{}, interface{})
}
// widget to display field as content type, it fetches automaitcally all content types to be choosed from
type ContentTypeSelectorWidget struct {
	Widget
	OptGroups             map[string][]*SelectOptGroup
	LoadFieldsOfAllModels bool
}
// not tested in UI. To be described later.
type NullBooleanWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}
// not tested in UI. To be described later.
type SelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}
// not tested in UI. To be described later.
type RadioSelectWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	ID        string
	WrapLabel bool
}
// not tested in UI. To be described later.
type CheckboxSelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	ID        string
	WrapLabel bool
}
// it allows you to store uploaded file, by default it uploads to the FS. To the directory for uploaded files,
// but it can be easily changed to store it to S3 or any different place.
// Please check out documentation for storages to understand how to do that.
type FileWidget struct {
	Widget
	Storage    IStorageInterface
	UploadPath string
	Multiple   bool
}
// not tested in UI. To be described later.
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
// not tested in UI. To be described later.
type MultipleInputHiddenWidget struct {
	Widget
}
// implements it like: two selects(left and right side) and it allows you to move between two selects
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
// splits date and time into two separated fields, not tested in UI. To be described later.
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
// not tested in UI. To be described later.
type SplitHiddenDateTimeWidget struct {
	Widget
	DateAttrs  map[string]string
	TimeAttrs  map[string]string
	DateFormat string
	TimeFormat string
	DateValue  string
	TimeValue  string
}
// implements it like three separated selects: for year, month, day
// not tested in UI. To be described later.
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

```
