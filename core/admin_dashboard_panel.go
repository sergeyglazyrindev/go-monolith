package core

import (
	"bytes"
	"fmt"
	excelize1 "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/http"
	"reflect"
)

type DashboardAdminPanel struct {
	AdminPages  *AdminPageRegistry
	ListHandler func(ctx *gin.Context)
}

type AutocompleteItemResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (dap *DashboardAdminPanel) FindPageForGormModel(m interface{}) *AdminPage {
	mDescription := ProjectModels.GetModelFromInterface(m)
	for adminPage := range dap.AdminPages.GetAll() {
		for subPage := range adminPage.SubPages.GetAll() {
			modelDescription := ProjectModels.GetModelFromInterface(subPage.Model)
			if modelDescription.Statement.Table == mDescription.Statement.Table {
				return subPage
			}
		}
	}
	return nil
}

func (dap *DashboardAdminPanel) RegisterHTTPHandlers(router *gin.Engine) {
	if dap.ListHandler != nil {
		router.GET(CurrentConfig.D.GoMonolith.RootAdminURL+"/", dap.ListHandler)
	}
	for adminPage := range dap.AdminPages.GetAll() {
		router.GET(fmt.Sprintf("%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug), func(pageTitle string, adminPageRegistry *AdminPageRegistry) func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				type Context struct {
					AdminContext
					Menu        string
					CurrentPath string
				}

				c := &Context{}
				PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
				if c.GetUserObject() != nil {
					menu := string(adminPageRegistry.PreparePagesForTemplate(c.UserPermissionRegistry))
					c.Menu = menu
				}
				c.CurrentPath = ctx.Request.URL.Path
				tr := NewTemplateRenderer(pageTitle)
				tr.Render(ctx, CurrentConfig.GetPathToTemplate("home"), c, FuncMap)
			}
		}(adminPage.PageName, adminPage.SubPages))
		for subPage := range adminPage.SubPages.GetAll() {
			if subPage.RegisteredHTTPHandlers {
				continue
			}
			router.Any(fmt.Sprintf("%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					if adminPage.ListHandler != nil {
						adminPage.ListHandler(ctx)
					} else {
						type Context struct {
							AdminContext
							AdminFilterObjects       IAdminFilterObjects
							ListDisplay              *ListDisplayRegistry
							PermissionForBlueprint   *UserPerm
							ListFilter               *ListFilterRegistry
							InitialOrder             string
							InitialOrderList         []string
							Search                   string
							TotalRecords             int64
							TotalPages               int64
							ListEditableFormError    bool
							AdminModelActionRegistry *AdminModelActionRegistry
							Message                  string
							Error                    error
							CurrentAdminContext      IAdminContext
							NoPermissionToAddNew     bool
							AdminPage                *AdminPage
							NoPermissionToEdit       bool
						}

						c := &Context{}
						c.NoPermissionToAddNew = adminPage.NoPermissionToAddNew
						adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
						PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
						user := c.GetUserObject()
						if user == nil {
							ctx.Redirect(302, adminPage.GetURLToBackAfterSignin(ctx))
							return
						}
						existsAnyPermission := user.BuildPermissionRegistry().IsThereAnyPermissionForBlueprint(adminPage.BlueprintName)
						if !existsAnyPermission {
							ctx.AbortWithStatus(409)
							return
						}
						c.Message = ctx.Query("message")
						c.NoPermissionToEdit = adminPage.NoPermissionToEdit
						c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
						c.AdminFilterObjects = adminPage.GetQueryset(c, adminPage, adminRequestParams)
						c.AdminModelActionRegistry = adminPage.ModelActionsRegistry
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.BlueprintName, URL: fmt.Sprintf("%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug)})
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.ModelName, IsActive: true})
						c.AdminPage = adminPage
						if ctx.Request.Method == "POST" {
							err := c.AdminFilterObjects.WithTransaction(func(afo1 IAdminFilterObjects) error {
								postForm, _ := ctx.MultipartForm()
								ids := postForm.Value["object_id"]
								modelDesc := ProjectModels.GetModelFromInterface(c.AdminFilterObjects.GetCurrentModel())
								for _, objectID := range ids {
									objectModel, _ := modelDesc.GenerateModelI()
									afo1.LoadDataForModelByID(objectID, objectModel)
									modelI, _ := modelDesc.GenerateModelI()
									listEditableForm := NewFormListEditableFromListDisplayRegistry(c, "", objectID, modelI, adminPage.ListDisplay, nil)
									formListEditableErr := listEditableForm.ProceedRequest(postForm, objectModel, c)
									if formListEditableErr.IsEmpty() {
										dbRes := afo1.SaveModel(objectModel)
										if dbRes != nil {
											c.ListEditableFormError = true
											return dbRes
										}
									} else {
										return formListEditableErr
									}
								}
								if afo1.GetLastError() != nil {
									return afo1.GetLastError()
								}
								return nil
							})
							if err != nil {
								c.Error = err
							}
						}
						c.AdminFilterObjects.GetFullQuerySet().Count(&c.TotalRecords)
						c.TotalPages = int64(math.Ceil(float64(c.TotalRecords / int64(adminPage.Paginator.PerPage))))
						c.ListDisplay = adminPage.ListDisplay
						c.Search = adminRequestParams.Search
						c.ListFilter = adminPage.ListFilter
						c.InitialOrder = adminRequestParams.GetOrdering()
						c.InitialOrderList = adminRequestParams.Ordering
						c.CurrentAdminContext = c
						tr := NewTemplateRenderer(adminPage.PageName)
						tr.Render(ctx, CurrentConfig.GetPathToTemplate("list"), c, FuncMap)
					}
				}
			}(subPage))
			router.POST(fmt.Sprintf("%s/%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.Slug, "export"), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					type Context struct {
						AdminContext
					}
					c := &Context{}
					adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
					PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
					user := c.GetUserObject()
					if !adminPage.DoesUserHavePermission(user, "read") {
						ctx.AbortWithStatus(409)
						return
					}
					// permissionForBlueprint := c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
					adminFilterObjects := adminPage.GetQueryset(c, adminPage, adminRequestParams)
					f := excelize1.NewFile()
					i := 1
					currentColumn := 'A'
					for listDisplay := range adminPage.ListDisplay.GetAllFields() {
						f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.DisplayName)
						currentColumn++
					}
					i++
					for iterateAdminObjects := range adminFilterObjects.IterateThroughWholeQuerySet() {
						currentColumn = 'A'
						for listDisplay := range adminPage.ListDisplay.GetAllFields() {
							f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.GetValue(iterateAdminObjects.Model, true))
							currentColumn++
						}
						i++
					}
					b, _ := f.WriteToBuffer()
					downloadName := adminPage.PageName + ".xlsx"
					ctx.Header("Content-Description", "File Transfer")
					ctx.Header("Content-Disposition", "attachment; filename="+downloadName)
					ctx.Data(http.StatusOK, "application/octet-stream", b.Bytes())
				}
			}(subPage))
			if len(subPage.SearchFields.Fields) > 0 {
				router.GET(fmt.Sprintf("%s/%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.Slug, "autocomplete"), func(adminPage *AdminPage) func(ctx *gin.Context) {
					return func(ctx *gin.Context) {
						type Context struct {
							AdminContext
						}
						c := &Context{}
						adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
						PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
						user := c.GetUserObject()
						if !adminPage.DoesUserHavePermission(user, "read") {
							ctx.AbortWithStatus(409)
							return
						}
						// permissionForBlueprint := c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
						adminFilterObjects := adminPage.GetQueryset(c, adminPage, adminRequestParams)
						resp := make([]*AutocompleteItemResponse, 0)
						for iterateAdminObjects := range adminFilterObjects.GetPaginated() {
							model := iterateAdminObjects.Model.(GoMonolithString)
							resp = append(resp, &AutocompleteItemResponse{Label: model.String(), Value: iterateAdminObjects.ID})
						}
						ctx.JSON(200, resp)
					}
				}(subPage))
			}
			router.Any(fmt.Sprintf("%s/%s/%s/edit/:id/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					id := ctx.Param("id")
					type Context struct {
						AdminContext
						AdminModelActionRegistry    *AdminModelActionRegistry
						Message                     string
						Error                       string
						PermissionForBlueprint      *UserPerm
						Form                        *Form
						Model                       interface{}
						ID                          uint
						IsNew                       bool
						ListURL                     string
						AdminPageInlineRegistry     *AdminPageInlineRegistry
						AdminRequestParams          *AdminRequestParams
						CurrentAdminContext         IAdminContext
						ListEditableFormsForInlines *FormListEditableCollection
						AdminPage                   *AdminPage
						RequestMethod string
					}

					c := &Context{}
					c.ListURL = fmt.Sprintf("%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug)
					c.PageTitle = adminPage.ModelName
					c.CurrentAdminContext = c
					c.ListEditableFormsForInlines = NewFormListEditableCollection()
					modelDesc := ProjectModels.GetModelFromInterface(adminPage.Model)
					modelI, _ := modelDesc.GenerateModelI()
					if id != "new" {
						adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
						qs := adminPage.GetQueryset(c, adminPage, adminRequestParams)
						qs.LoadDataForModelByID(id, modelI)
						// qs.CloseConnection()
					}
					adminRequestParams := NewAdminRequestParams()
					c.AdminRequestParams = adminRequestParams
					PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)
					form := adminPage.GenerateForm(modelI, c)
					field, _ := form.FieldRegistry.GetByName("ID")
					ID, _ := field.FieldConfig.Widget.GetValue().(uint)
					c.ID = ID
					form.TemplateName = "admin/form_edit"
					form.RequestContext["ID"] = c.GetID()
					c.Model = modelI
					form.DontGenerateFormTag = true
					c.IsNew = true
					c.AdminPageInlineRegistry = adminPage.InlineRegistry
					c.AdminPage = adminPage
					c.RequestMethod = ctx.Request.Method
					form.ForAdminPanel = true
					user := c.GetUserObject()
					if user == nil {
						ctx.Redirect(302, adminPage.GetURLToBackAfterSignin(ctx))
						return
					}
					if ctx.Request.Method == "POST" {
						if id != "new" {
							if !subPage.DoesUserHavePermission(user, "edit") {
								ctx.AbortWithStatus(409)
								return
							}
						} else {
							if !subPage.DoesUserHavePermission(user, "add") {
								ctx.AbortWithStatus(409)
								return
							}
						}
						requestForm, _ := ctx.MultipartForm()
						var modelToSave interface{}
						if id != "new" {
							modelToSave = modelI
						} else {
							modelDesc1 := ProjectModels.GetModelFromInterface(adminPage.Model)
							modelToSave, _ = modelDesc1.GenerateModelI()
						}
						afo := adminPage.GetQueryset(c, adminPage, adminRequestParams)
						err := afo.WithTransaction(func(afo1 IAdminFilterObjects) error {
							formError := form.ProceedRequest(requestForm, modelToSave, c, afo1)
							if formError.IsEmpty() {
								if adminPage.SaveModel != nil {
									modelToSave = adminPage.SaveModel(modelToSave, ID, afo1)
									err1 := afo1.GetLastError()
									if err1 != nil {
										for inline := range adminPage.InlineRegistry.GetAll() {
											inlineListEditableCollection, _ := inline.ProceedRequest(afo1, requestForm, modelToSave, c, true)
											c.ListEditableFormsForInlines.AddForInlineWholeCollection(inline.Prefix, inlineListEditableCollection)
										}
										return err1
									}
								} else {
									afo1.GetInitialQuerySet().Model(modelToSave).Save(modelToSave)
									err1 := afo1.GetInitialQuerySet().GetLastError()
									if err1 != nil {
										afo1.SetLastError(err1)
										for inline := range adminPage.InlineRegistry.GetAll() {
											inlineListEditableCollection, _ := inline.ProceedRequest(afo1, requestForm, modelToSave, c, true)
											c.ListEditableFormsForInlines.AddForInlineWholeCollection(inline.Prefix, inlineListEditableCollection)
										}
										return err1
									}
								}
								var successfulInline error
								for inline := range adminPage.InlineRegistry.GetAll() {
									inlineListEditableCollection, formError1 := inline.ProceedRequest(afo1, requestForm, modelToSave, c)
									if formError1 != nil {
										successfulInline = formError1
									}
									c.ListEditableFormsForInlines.AddForInlineWholeCollection(inline.Prefix, inlineListEditableCollection)
								}
								if successfulInline != nil {
									return NewHTTPErrorResponse(successfulInline.Error(), successfulInline.Error())
								}
								if ctx.Query("_popup") == "1" {
									mID := GetID(reflect.ValueOf(modelToSave))
									data := make(map[string]interface{})
									data["Link"] = ctx.Request.URL.String()
									data["ID"] = mID
									data["Name"] = reflect.ValueOf(modelToSave).MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)
									htmlResponseWriter := bytes.NewBuffer(make([]byte, 0))
									AddedObjectInPopup.ExecuteTemplate(htmlResponseWriter, "addedobjectinpopup", data)
									ctx.Data(http.StatusOK, "text/html; charset=utf-8", htmlResponseWriter.Bytes())
								} else if len(requestForm.Value["save_add_another"]) > 0 {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/edit/new/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug))
								} else if len(requestForm.Value["save_continue"]) > 0 {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/edit/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug, id))
								} else {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug))
								}
								return nil
							}
							return NewHTTPErrorResponse("not_successful_form_validation", "not successful form validation")
						})
						if err != nil {
							form.FormError.GeneralErrors = append(form.FormError.GeneralErrors, err)
						} else {
							return
						}
					} else {
						if id != "new" {
							if !subPage.DoesUserHavePermission(user, "edit") {
								ctx.AbortWithStatus(409)
								return
							}
						} else {
							if !subPage.DoesUserHavePermission(user, "add") {
								ctx.AbortWithStatus(409)
								return
							}
						}
						for inline := range adminPage.InlineRegistry.GetAll() {
							if id == "new" {
								continue
							}
							for iterateAdminObjects := range inline.GetAll(c, c.Model) {
								listEditable := inline.ListDisplay.BuildFormForListEditable(c, iterateAdminObjects.ID, iterateAdminObjects.Model, nil)
								c.ListEditableFormsForInlines.AddForInline(inline.Prefix, iterateAdminObjects.ID, listEditable)
							}
						}
					}
					c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.BlueprintName, URL: fmt.Sprintf("%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug)})
					c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.ModelName, URL: fmt.Sprintf("%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug)})
					if id != "new" {
						values := reflect.ValueOf(modelI).MethodByName("String").Call([]reflect.Value{})
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{IsActive: true, Name: values[0].String()})
					} else {
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{IsActive: true, Name: "New"})
					}
					c.Form = form
					c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
					c.Message = ctx.Query("message")
					c.AdminModelActionRegistry = adminPage.ModelActionsRegistry
					tr := NewTemplateRenderer(adminPage.PageName)
					tr.Render(ctx, CurrentConfig.GetPathToTemplate("change"), c, FuncMap)
				}
			}(subPage))
			for adminModelAction := range subPage.ModelActionsRegistry.GetAllModelActions() {
				router.Any(fmt.Sprintf("%s/%s/%s/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *AdminPage, slugifiedModelActionName string) func(ctx *gin.Context) {
					return func(ctx *gin.Context) {
						adminPage.HandleModelAction(slugifiedModelActionName, ctx)
					}
				}(subPage, adminModelAction.SlugifiedActionName))
			}
			for pageInline := range subPage.InlineRegistry.GetAll() {
				for inlineAdminModelAction := range pageInline.Actions.GetAllModelActions() {
					router.Any(fmt.Sprintf("%s/%s/%s/edit/:id/%s/", CurrentConfig.D.GoMonolith.RootAdminURL, adminPage.Slug, subPage.ModelName, inlineAdminModelAction.SlugifiedActionName), func(adminPage *AdminPage, adminPageInline *AdminPageInline, slugifiedModelActionName string) func(ctx *gin.Context) {
						return func(ctx *gin.Context) {
							adminPage.HandleModelAction(slugifiedModelActionName, ctx)
						}
					}(subPage, pageInline, inlineAdminModelAction.SlugifiedActionName))
				}
			}
			subPage.RegisteredHTTPHandlers = true
		}
	}
}

var CurrentDashboardAdminPanel *DashboardAdminPanel

func NewDashboardAdminPanel() *DashboardAdminPanel {
	adminPageRegistry := NewAdminPageRegistry()
	CurrentAdminPageRegistry = adminPageRegistry
	return &DashboardAdminPanel{
		AdminPages: adminPageRegistry,
	}
}

var AddedObjectInPopup *template.Template

func init() {
	AddedObjectInPopup, _ = template.New("addedobjectinpopup").Parse(`{{define "addedobjectinpopup"}}<html><head></head><body>
<script type="text/javascript">
	var link = "{{ .Link }}";
	var ID = "{{ .ID }}";
	var Name = "{{ .Name }}";
    var relatedTarget = window.opener.$("a[href='{{ .Link }}']").parent().parent().find('.related-target');
    if (relatedTarget.find('.hidden-id').length > 0) {
		var hiddenId = relatedTarget.find('.hidden-id');
		hiddenId.val(ID);
		hiddenId.prev().text(Name);
	} else {
		var newOption = window.opener.$('<select><option value=""></option></select>');
		newOption.find('option').attr('value', ID);
		newOption.find('option').text(Name);
		newOption.find('option').attr('selected', 'selected');
		var select = relatedTarget.find('select');
		select.find('option:selected').removeAttr('selected');
		select.append(newOption.html());
		select.trigger('change');
	}
	window.close();
</script>
</body></html>{{end}}
`)

}
