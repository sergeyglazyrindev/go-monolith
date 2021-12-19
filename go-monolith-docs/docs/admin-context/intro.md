---
sidebar_position: 1
---

# Admin context

Admin context is in use everywhere in the admin panel. It implements following interface
```go
type IAdminContext interface {
	SetSiteName(siteName string)
	SetCurrentURL(currentURL string)
	SetFullURL(fullURL *url.URL)
	SetRootAdminURL(rootAdminURL string)
	SetSessionKey(sessionKey string)
	SetRootURL(rootURL string)
	SetLanguage(language *Language)
	SetLogo(logo string)
	SetFavIcon(favicon string)
	SetLanguages(langs []Language)
	SetPageTitle(pageTitle string)
	SetUser(user string)
	SetUserExists(userExists bool)
	SetDemo()
	SetError(err string)
	SetErrorExists()
	GetLanguage() *Language
	GetRootURL() string
	SetUserPermissionRegistry(permRegistry *UserPermRegistry)
	SetForm(form IForm)
	SetCurrentQuery(currentQuery string)
	SetBreadCrumbs(breadcrumbs *AdminBreadCrumbsRegistry)
	GetSessionKey() string
	SetID(ID uint)
	GetID() uint
	SetUserObject(u *User)
	GetUserObject() *User
	SetPostForm(formD *multipart.Form)
	GetPostForm() *multipart.Form
	GetCtx() *gin.Context
	SetCtx(ctx *gin.Context)
}

```
Later on we will provide example how it could be customized for your specific administration use cases.
