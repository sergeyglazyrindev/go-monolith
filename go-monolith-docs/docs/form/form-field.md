# Form fields

GoMonolith form field extends Gorm Schema Field, and it has following structure:
```go
type FieldConfig struct {
  // field widget
	Widget                 IWidget
  // right now not implemented
	AutocompleteURL        string
  // maybe will be removed later
	DependsOnAnotherFields []string
}

type Field struct {
	schema.Field
	ReadOnly        bool
	GoMonolithFieldType GoMonolithFieldType
	FieldConfig     *FieldConfig
	Required        bool
	DisplayName     string
	HelpText        string
	Choices         *FieldChoiceRegistry
  // list of custom validators to be executed during saving this model
	Validators      *ValidatorRegistry
	SortingDisabled bool
  // populate value for this field
	Populate        func(field *Field, m interface{}) interface{}
	Initial         interface{}
	WidgetType      string
  // the way to store value of the field properly to model, in case if there's any need for this.
  // could be used like this:
  // permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
  // permissionsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
  //   model := modelI.(*core.UserGroup)
  //   vTmp := v.([]string)
  //   var permission *core.Permission
  //   if model.ID != 0 {
  //     afo.GetDatabase().Db.Model(model).Association("Permissions").Clear()
  //     model.Permissions = make([]core.Permission, 0)
  //   }
  //   for _, ID := range vTmp {
  //     afo.GetDatabase().Db.First(&permission, ID)
  //     if permission.ID != 0 {
  //       model.Permissions = append(model.Permissions, *permission)
  //     }
  //     permission = nil
  //   }
  //   return nil
  // }  
	SetUpField      func(w IWidget, modelI interface{}, v interface{}, afo IAdminFilterObjects) error
	Ordering        int
}

```
