---
sidebar_position: 1
---

# Elastic search support

GoMonolith supports elastic search (please build your package with tag 'elasticsearch'). For this you need just to describe your data model, like below (model has to implement core.ElasticModelInterface interface):
```go
type Tweet struct {
	User    string `json:"user" gomonolith:"list,search"`
	Message string `json:"message" gomonolith:"list,search"`
	ID string
}

func (t *Tweet) String() string {
	return fmt.Sprintf("User %s tweeted following %s", t.User, t.Message)
}

func (t *Tweet) SetID(ID string) {
	t.ID = ID
}

func (t *Tweet) GetID() string {
	return t.ID
}

func (t *Tweet) GetIndexName() string {
	return "tweets"
}
```
don't forget to create index in data migration:
```go
// Create a client
client := core.NewESClient()
client.DeleteIndex("tweets").Do(context.Background())
// Create an index
_, err := client.CreateIndex("tweets").Do(context.Background())
if err != nil {
		// Handle error
		panic(err)
}
return nil
```
and then register admin panel for your elastic search model:
```go
tweetsAdminPage := core.NewElasticSearchAdminPage(
	nil,
	func() (interface{}, interface{}) { return nil, nil },
	func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
)
tweetsAdminPage.PageName = "Tweets"
tweetsAdminPage.Slug = "tweets"
tweetsAdminPage.BlueprintName = "tweets"
tweetsAdminPage.Router = mainRouter
err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(tweetsAdminPage)
if err != nil {
	panic(fmt.Errorf("error initializing tweets blueprint: %s", err))
}
var tweetsmodelAdminPage *core.AdminPage
tweetsmodelAdminPage = core.NewElasticSearchAdminPage(
	tweetsAdminPage,
	func() (interface{}, interface{}) { return &Tweet{}, &[]*Tweet{} },
	func(modelI interface{}, ctx core.IAdminContext) *core.Form {
		fields := []string{"User", "Message"}
		form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
		return form
	},
)
tweetsmodelAdminPage.PageName = "Tweets"
tweetsmodelAdminPage.Slug = "tweet"
tweetsmodelAdminPage.BlueprintName = "tweets"
tweetsmodelAdminPage.Router = mainRouter
tweetsmodelAdminPage.ModelName = "tweet"
IDListDisplayField, _ := tweetsmodelAdminPage.ListDisplay.GetFieldByDisplayName("ID")
IDListDisplayField.SortBy.SetSortCustomImplementation(func(afo core.IAdminFilterObjects, field *core.Field, direction int) {
	directionB := true
	if direction == -1 {
		directionB = false
	}
	afo.GetPaginatedQuerySet().Order(&core.ESSortBy{
		FieldName: "_id",
		Direction: directionB,
	})
})
err = tweetsAdminPage.SubPages.AddAdminPage(tweetsmodelAdminPage)
if err != nil {
	panic(fmt.Errorf("error initializing tweets blueprint: %s", err))
}
```
