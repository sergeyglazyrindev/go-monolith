---
sidebar_position: 1
---

# Forms

Forms is the way to build your edit functionality in go-monolith admin panel. Later it could be used for non admin too. But nowadays it's not usual to use multipart/form-data in the projects.
It has following structure.
```go
type Form struct {
	// in use only when we build form, probably could be removed later
	ExcludeFields       IFieldRegistry
	// in use only when we build form, probably could be removed later
	FieldsToShow        IFieldRegistry
	FieldRegistry       IFieldRegistry
	// in use only when we render form, probably could be removed later
	GroupsOfTheFields   *GrouppedFieldsRegistry
	// template name, if empty default would be used
	TemplateName        string
	// form title
	FormTitle           string
	// form renderer
	Renderer            ITemplateRenderer
	// request context
	RequestContext      map[string]interface{}
	// error message, maybe removed later
	ErrorMessage        string
	// could be used like below:
	// form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/gomonolith/assets/highlight.js/highlight.pack.js")
	// form.ExtraStatic.ExtraJS = append(form.ExtraStatic.ExtraJS, "/static-inbuilt/gomonolith/assets/js/initialize.highlight.js")
	// form.ExtraStatic.ExtraCSS = append(form.ExtraStatic.ExtraCSS, "/static-inbuilt/gomonolith/assets/highlight.js/styles/default.css")
	ExtraStatic         *StaticFiles `json:"-"`
	// if it's true, all fields and its widgets would be rendered as for admin panel.
	ForAdminPanel       bool
	// form error, it contains all field errors, general errors
	FormError           *FormError
	DontGenerateFormTag bool
	// used to render field names
	Prefix              string
	// render context
	RenderContext       *FormRenderContext
}
```
The easiest way to create form for Gorm models is to use function `NewFormFromModelFromGinContext`, like here:
```go
fields := []string{"ContentType", "Type", "Name", "Field", "PrimaryKey", "Active", "Group", "StaticPath"}
form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
// or
modelForm := NewFormFromModel(model, []string{}, []string{}, false, "")
```
