package core

import (
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
	"mime/multipart"
	"reflect"
	"sort"
)

type IFieldRegistry interface {
	GetByName(name string) (*Field, error)
	AddField(field *Field)
	GetAllFields() map[string]*Field
	GetAllFieldsWithOrdering() []*Field
	GetPrimaryKey() (*Field, error)
}

type GoMonolithFieldType string

const BigIntegerFieldType GoMonolithFieldType = "biginteger"
const BinaryFieldType GoMonolithFieldType = "binary"
const BooleanFieldType GoMonolithFieldType = "boolean"
const CharFieldType GoMonolithFieldType = "char"
const DateFieldType GoMonolithFieldType = "date"
const DateTimeFieldType GoMonolithFieldType = "datetime"
const DecimalFieldType GoMonolithFieldType = "decimal"
const DurationFieldType GoMonolithFieldType = "duration"
const EmailFieldType GoMonolithFieldType = "email"
const FileFieldType GoMonolithFieldType = "file"
const FilePathFieldType GoMonolithFieldType = "filepath"
const FloatFieldType GoMonolithFieldType = "float"
const ForeignKeyFieldType GoMonolithFieldType = "foreignkey"
const ImageFieldFieldType GoMonolithFieldType = "imagefield"
const IntegerFieldType GoMonolithFieldType = "integer"
const UintFieldType GoMonolithFieldType = "uint"
const IPAddressFieldType GoMonolithFieldType = "ipaddress"
const GenericIPAddressFieldType GoMonolithFieldType = "genericipaddress"
const ManyToManyFieldType GoMonolithFieldType = "manytomany"
const NullBooleanFieldType GoMonolithFieldType = "nullboolean"
const PositiveBigIntegerFieldType GoMonolithFieldType = "positivebiginteger"
const PositiveIntegerFieldType GoMonolithFieldType = "positiveinteger"
const PositiveSmallIntegerFieldType GoMonolithFieldType = "positivesmallinteger"
const SlugFieldType GoMonolithFieldType = "slug"
const SmallIntegerFieldType GoMonolithFieldType = "smallinteger"
const TextFieldType GoMonolithFieldType = "text"
const TimeFieldType GoMonolithFieldType = "time"
const URLFieldType GoMonolithFieldType = "url"
const UUIDFieldType GoMonolithFieldType = "uuid"

type FieldConfig struct {
	Widget                 IWidget
	AutocompleteURL        string
	DependsOnAnotherFields []string
}

type Field struct {
	schema.Field
	ReadOnly            bool
	GoMonolithFieldType GoMonolithFieldType
	FieldConfig         *FieldConfig
	Required            bool
	DisplayName         string
	HelpText            string
	Choices             *FieldChoiceRegistry
	Validators          *ValidatorRegistry
	SortingDisabled     bool
	Populate            func(field *Field, m interface{}) interface{}
	Initial             interface{}
	WidgetType          string
	SetUpField          func(w IWidget, modelI interface{}, v interface{}, afo IAdminFilterObjects) error
	Ordering            int
}

func (f *Field) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) ValidationError {
	err := f.FieldConfig.Widget.ProceedForm(form, afo, renderContext)
	if err == nil {
		validationErrors := make(ValidationError, 0)
		for validator := range f.Validators.GetAllValidators() {
			validationErr := validator(f.FieldConfig.Widget.GetOutputValue(), form)
			if validationErr == nil {
				continue
			}
			validationErrors = append(validationErrors, validationErr)
		}
		if len(validationErrors) > 0 {
			f.FieldConfig.Widget.SetErrors(validationErrors)
		}
		return validationErrors
	}
	errors := ValidationError{err}
	f.FieldConfig.Widget.SetErrors(errors)
	return errors
}

type ValidationError []error

func NewFieldForListDisplayFromGormField(gormField *schema.Field, fieldOptions IFieldFormOptions) *Field {
	var widget IWidget
	forcedWidgetType := ""
	if fieldOptions != nil && fieldOptions.GetListFieldWidget() != "" {
		forcedWidgetType = fieldOptions.GetListFieldWidget()
	}
	if forcedWidgetType == "" && fieldOptions != nil && fieldOptions.GetWidgetType() != "" {
		forcedWidgetType = fieldOptions.GetWidgetType()
	}
	if gormField.PrimaryKey {
		widget = GetWidgetByWidgetType("hidden", nil)
	} else if forcedWidgetType != "" {
		widget = GetWidgetByWidgetType(forcedWidgetType, fieldOptions)
	} else {
		fieldType := GetGoMonolithFieldTypeFromGormField(gormField)
		widget = GetWidgetFromFieldTypeAndGormField(fieldType, gormField)
	}
	widget.InitializeAttrs()
	widget.SetName(gormField.Name)
	if gormField.NotNull && !gormField.HasDefaultValue {
		widget.SetRequired()
	}
	if gormField.Unique {
		widget.SetRequired()
	}
	if !gormField.PrimaryKey {
		widget.SetValue(gormField.DefaultValueInterface)
	}
	field := &Field{
		Field:               *gormField,
		GoMonolithFieldType: GetGoMonolithFieldTypeFromGormField(gormField),
		FieldConfig:         &FieldConfig{Widget: widget},
		Required:            gormField.NotNull && !gormField.HasDefaultValue,
		DisplayName:         gormField.Name,
	}
	return field
}

func NewFieldFromGormField(gormField *schema.Field, fieldOptions IFieldFormOptions) *Field {
	var widget IWidget
	forcedWidgetType := ""
	if fieldOptions != nil && fieldOptions.GetWidgetType() != "" {
		forcedWidgetType = fieldOptions.GetWidgetType()
	}
	if gormField.PrimaryKey {
		widget = GetWidgetByWidgetType("hidden", nil)
	} else if forcedWidgetType != "" {
		widget = GetWidgetByWidgetType(forcedWidgetType, fieldOptions)
	} else {
		fieldType := GetGoMonolithFieldTypeFromGormField(gormField)
		widget = GetWidgetFromFieldTypeAndGormField(fieldType, gormField)
	}
	widget.InitializeAttrs()
	widget.SetName(gormField.Name)
	if gormField.NotNull && !gormField.HasDefaultValue {
		widget.SetRequired()
	}
	if gormField.Unique {
		widget.SetRequired()
	}
	if !gormField.PrimaryKey {
		widget.SetValue(gormField.DefaultValueInterface)
	}
	field := &Field{
		Field:               *gormField,
		GoMonolithFieldType: GetGoMonolithFieldTypeFromGormField(gormField),
		FieldConfig:         &FieldConfig{Widget: widget},
		Required:            gormField.NotNull && !gormField.HasDefaultValue,
		DisplayName:         gormField.Name,
	}
	return field
}

func NewGoMonolithFieldForListDisplayFromGormField(gormModelV reflect.Value, field *schema.Field, r ITemplateRenderer, renderForAdmin bool) *Field {
	formtag := field.Tag.Get("gomonolithform")
	var fieldOptions IFieldFormOptions
	var field1 *Field
	if formtag != "" {
		fieldOptions = FormCongirurableOptionInstance.GetFieldFormOptions(formtag)
		field1 = NewFieldForListDisplayFromGormField(field, fieldOptions)
	} else {
		if field.PrimaryKey {
			fieldOptions = FormCongirurableOptionInstance.GetFieldFormOptions("ReadonlyField")
			field1 = NewFieldFromGormField(field, fieldOptions)
		} else {
			field1 = NewFieldFromGormField(field, nil)
		}
	}
	field1.DisplayName = field.Name
	if renderForAdmin {
		field1.FieldConfig.Widget.RenderForAdmin()
	}
	if fieldOptions != nil {
		field1.Initial = fieldOptions.GetInitial()
		if fieldOptions.GetDisplayName() != "" {
			field1.DisplayName = fieldOptions.GetDisplayName()
		}
		if fieldOptions.GetWidgetPopulate() != nil {
			field1.FieldConfig.Widget.SetPopulate(fieldOptions.GetWidgetPopulate())
		}
		field1.Validators = fieldOptions.GetValidators()
		field1.Choices = fieldOptions.GetChoices()
		field1.HelpText = fieldOptions.GetHelpText()
		field1.WidgetType = fieldOptions.GetWidgetType()
		field1.ReadOnly = fieldOptions.GetReadOnly()
		field1.FieldConfig.Widget.SetReadonly(field1.ReadOnly)
		if fieldOptions.GetIsRequired() {
			field1.FieldConfig.Widget.SetRequired()
		}
		if fieldOptions.GetHelpText() != "" {
			field1.FieldConfig.Widget.SetHelpText(fieldOptions.GetHelpText())
		}
	}
	field1.FieldConfig.Widget.RenderUsingRenderer(r)
	field1.FieldConfig.Widget.SetFieldDisplayName(field.Name)
	isTruthyValue := IsTruthyValue(gormModelV.FieldByName(field.Name).Interface())
	if isTruthyValue {
		field1.FieldConfig.Widget.SetValue(gormModelV.FieldByName(field.Name).Interface())
	}
	return field1
}

func NewGoMonolithFieldFromGormField(gormModelV reflect.Value, field *schema.Field, r ITemplateRenderer, renderForAdmin bool) *Field {
	formtag := field.Tag.Get("gomonolithform")
	var fieldOptions IFieldFormOptions
	var field1 *Field
	if formtag != "" {
		fieldOptions = FormCongirurableOptionInstance.GetFieldFormOptions(formtag)
		field1 = NewFieldFromGormField(field, fieldOptions)
	} else {
		if field.PrimaryKey {
			fieldOptions = FormCongirurableOptionInstance.GetFieldFormOptions("ReadonlyField")
			field1 = NewFieldFromGormField(field, fieldOptions)
		} else {
			field1 = NewFieldFromGormField(field, nil)
		}
	}
	field1.DisplayName = field.Name
	if renderForAdmin {
		field1.FieldConfig.Widget.RenderForAdmin()
	}
	field1.Validators = NewValidatorRegistry()
	if fieldOptions != nil {
		field1.Initial = fieldOptions.GetInitial()
		if fieldOptions.GetDisplayName() != "" {
			field1.DisplayName = fieldOptions.GetDisplayName()
		}
		if fieldOptions.GetWidgetPopulate() != nil {
			field1.FieldConfig.Widget.SetPopulate(fieldOptions.GetWidgetPopulate())
		}
		field1.Validators = fieldOptions.GetValidators()
		field1.Choices = fieldOptions.GetChoices()
		field1.HelpText = fieldOptions.GetHelpText()
		field1.WidgetType = fieldOptions.GetWidgetType()
		field1.ReadOnly = fieldOptions.GetReadOnly()
		field1.FieldConfig.Widget.SetReadonly(field1.ReadOnly)
		if fieldOptions.GetIsRequired() {
			field1.FieldConfig.Widget.SetRequired()
		}
		if fieldOptions.GetHelpText() != "" {
			field1.FieldConfig.Widget.SetHelpText(fieldOptions.GetHelpText())
		}
	}
	field1.FieldConfig.Widget.RenderUsingRenderer(r)
	field1.FieldConfig.Widget.SetFieldDisplayName(field.Name)
	isTruthyValue := IsTruthyValue(gormModelV.FieldByName(field.Name).Interface())
	if isTruthyValue {
		field1.FieldConfig.Widget.SetValue(gormModelV.FieldByName(field.Name).Interface())
	}
	return field1
}

type FieldRegistry struct {
	Fields      map[string]*Field
	MaxOrdering int
}

func (fr *FieldRegistry) GetByName(name string) (*Field, error) {
	f, ok := fr.Fields[name]
	if !ok {
		return nil, NewHTTPErrorResponse("field_not_found", "no field %s found", name)
	}
	return f, nil
}

func (fr *FieldRegistry) GetAllFields() map[string]*Field {
	return fr.Fields
}

func (fr *FieldRegistry) GetAllFieldsWithOrdering() []*Field {
	allFields := make([]*Field, 0)
	for _, field := range fr.Fields {
		allFields = append(allFields, field)
	}
	sort.Slice(allFields, func(i int, j int) bool {
		if allFields[i].Ordering == allFields[j].Ordering {
			return allFields[i].Name < allFields[j].Name
		}
		return allFields[i].Ordering < allFields[j].Ordering
	})
	return allFields
}

func (fr *FieldRegistry) GetPrimaryKey() (*Field, error) {
	for _, field := range fr.Fields {
		if field.PrimaryKey {
			return field, nil
		}
	}
	return nil, errors.New("no primary key found for model")
}

func (fr *FieldRegistry) AddField(field *Field) {
	if _, err := fr.GetByName(field.Name); err == nil {
		Trail(ERROR, fmt.Errorf("field %s already in the field registry", field.Name))
		return
	}
	fr.Fields[field.Name] = field
	ordering := fr.MaxOrdering + 1
	field.Ordering = ordering
	fr.MaxOrdering = ordering
}

func NewFieldRegistry() *FieldRegistry {
	return &FieldRegistry{Fields: make(map[string]*Field)}
}
