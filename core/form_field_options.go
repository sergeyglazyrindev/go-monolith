package core

type FieldFormOptions struct {
	Name            string
	Initial         interface{}
	DisplayName     string
	Validators      *ValidatorRegistry
	Choices         *FieldChoiceRegistry
	HelpText        string
	WidgetType      string
	ReadOnly        bool
	Required        bool
	WidgetPopulate  func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}
	IsFk            bool
	Autocomplete    bool
	ListFieldWidget string
}

func (ffo *FieldFormOptions) GetName() string {
	return ffo.Name
}

func (ffo *FieldFormOptions) IsItFk() bool {
	return ffo.IsFk
}

func (ffo *FieldFormOptions) GetListFieldWidget() string {
	return ffo.ListFieldWidget
}

func (ffo *FieldFormOptions) GetIsAutocomplete() bool {
	return ffo.Autocomplete
}

func (ffo *FieldFormOptions) GetWidgetPopulate() func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{} {
	return ffo.WidgetPopulate
}

func (ffo *FieldFormOptions) GetInitial() interface{} {
	return ffo.Initial
}

func (ffo *FieldFormOptions) GetDisplayName() string {
	return ffo.DisplayName
}

func (ffo *FieldFormOptions) GetValidators() *ValidatorRegistry {
	if ffo.Validators == nil {
		return NewValidatorRegistry()
	}
	return ffo.Validators
}

func (ffo *FieldFormOptions) GetChoices() *FieldChoiceRegistry {
	return ffo.Choices
}

func (ffo *FieldFormOptions) GetHelpText() string {
	return ffo.HelpText
}

func (ffo *FieldFormOptions) GetWidgetType() string {
	return ffo.WidgetType
}

func (ffo *FieldFormOptions) GetReadOnly() bool {
	return ffo.ReadOnly
}

func (ffo *FieldFormOptions) GetIsRequired() bool {
	return ffo.Required
}

type FormConfigurableOptionRegistry struct {
	Options map[string]IFieldFormOptions
}

func (c *FormConfigurableOptionRegistry) AddFieldFormOptions(formOptions IFieldFormOptions) {
	c.Options[formOptions.GetName()] = formOptions
}

func (c *FormConfigurableOptionRegistry) GetFieldFormOptions(formOptionsName string) IFieldFormOptions {
	ret, _ := c.Options[formOptionsName]
	return ret
}

var FormCongirurableOptionInstance *FormConfigurableOptionRegistry

func init() {
	FormCongirurableOptionInstance = &FormConfigurableOptionRegistry{
		Options: make(map[string]IFieldFormOptions),
	}
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "RequiredSelectFieldOptions",
		WidgetType: "select",
		Required:   true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ContentTypeFieldOptions",
		WidgetType: "contenttypeselector",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "SelectFieldOptions",
		WidgetType: "select",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ReadonlyTextareaFieldOptions",
		WidgetType: "textarea",
		ReadOnly:   true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "TextareaFieldOptions",
		WidgetType: "textarea",
	})
	fieldChoiceRegistry := FieldChoiceRegistry{}
	fieldChoiceRegistry.Choices = make([]*FieldChoice, 0)
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:        "UsernameOptions",
		Initial:     "InitialUsername",
		DisplayName: "Username",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ImageFormOptions",
		WidgetType: "image",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "OTPRequiredOptions",
		WidgetType: "hidden",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:     "ReadonlyField",
		ReadOnly: true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "PasswordOptions",
		WidgetType: "password",
		HelpText:   "To reset password, clear the field and type a new password.",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ChooseFromSelectOptions",
		WidgetType: "choose_from_select",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DateTimeFieldOptions",
		WidgetType: "datetime",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DatetimeReadonlyFieldOptions",
		WidgetType: "datetime",
		ReadOnly:   true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:     "RequiredFieldOptions",
		Required: true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkFieldOptions",
		IsFk:       true,
		WidgetType: "fklink",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkReadonlyFieldOptions",
		IsFk:       true,
		ReadOnly:   true,
		WidgetType: "fklink",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkRequiredFieldOptions",
		IsFk:       true,
		Required:   true,
		WidgetType: "fklink",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DynamicTypeFieldOptions",
		WidgetType: "dynamic",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ForeignKeyFieldOptions",
		WidgetType: "foreignkey",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:            "ForeignKeyWithAutocompleteFieldOptions",
		WidgetType:      "foreignkey",
		Autocomplete:    true,
		ListFieldWidget: "fklink",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:            "ForeignKeyReadonlyFieldOptions",
		WidgetType:      "foreignkey",
		ReadOnly:        true,
		Autocomplete:    true,
		ListFieldWidget: "fklink",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "TextareaReadonlyFieldOptions",
		WidgetType: "textarea",
		ReadOnly:   true,
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "EmailFieldOptions",
		WidgetType: "email",
	})
	FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "URLFieldOptions",
		WidgetType: "url",
	})
}
