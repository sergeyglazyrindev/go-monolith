---
sidebar_position: 1
---

# Admin list display

Admin list display is designed to provide developers a way to use absolutely any type of column in the list view for model.
An example of how it could be used.
```go
// here we get type list display field of the admin page.
typeListDisplay, _ := abtestmodelAdminPage.ListDisplay.GetFieldByDisplayName("Type")
// custom populate method for type list display
typeListDisplay.Populate = func(m interface{}) string {
	return abtestmodel.HumanizeAbTestType(m.(*abtestmodel.ABTest).Type)
}
```
You can add your own list display column to the model view.  
For example, like here:
```go
abTestValueInline.ListDisplay.AddField(&core.ListDisplay{
	DisplayName: "Click through rate",
	MethodName:  "ClickThroughRate",
})
```
And then value for this list display would be fetched using method: `ClickThroughRate` that is a method of the abtestValue model.  
Also you can make this field editable right on the list page. ListDisplay object does have a property:
```go
IsEditable bool
```
If you set it to true, then proper widget configured for this field would be shown on the list view and you edit it right away.  
There are helper methods that you can use to build ListDisplayRegistry from GormModel:
```go
// this function is related to inlines. Please check out corresponding part of the documentation.
func NewListDisplayRegistryFromGormModelForInlines(modelI interface{}) *ListDisplayRegistry {
}
// if you want the field to be included automatically to the list display registry, then please specify tag for the field: gomonolith: "list"
func NewListDisplayRegistryFromGormModel(modelI interface{}) *ListDisplayRegistry {
}
```
Also you can manipulate the order of the fields in the list view, there's a property:
```go
Ordering    int
```
Also you can change the way how the objects sorted against this field, you may provide your own SearchBy approach. There's a field:
```go
SortBy      ISortBy
type ISortBy interface {
	Sort(afo IAdminFilterObjects, direction int)
	GetDirection() int
}
```
Later on we will migrate it to interface as well, so it could be used easily for any type of list fields.
