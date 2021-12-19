---
sidebar_position: 1
---

# Admin dashboard panel

Admin dashboard panel is responsible for handling admin http requests. For model actions, though it provides one important method:

```go
func (dap *DashboardAdminPanel) FindPageForGormModel(m interface{}) *AdminPage {
}

```
This mwthod finds page that is responsible for handling admin functionality for this model.
Right now it's used internally, only in ForeignKey widget to build url that leads to the model's admin page.
