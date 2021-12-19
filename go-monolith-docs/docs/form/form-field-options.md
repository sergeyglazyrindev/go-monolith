# Form field options

It simplifies to configure fields for forms. For example, if you want to make some field readonly, just specify it in the gorm model:
```go
type Language struct {
	Code           string `gomonolithform:"ReadonlyField"`
}
```
The list of predefined field form options is here:
```go
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
```
You can easily add your own field form option:
```go
core.FormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "YOUROWNFIELDFORMOPTIONS",
		WidgetType: "select",
		Required:   true,
	})
```
