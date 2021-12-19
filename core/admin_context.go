package core

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/url"
)

var PopulateTemplateContextForAdminPanel func(ctx *gin.Context, context IAdminContext, adminRequestParams *AdminRequestParams)

type IForm interface {
}

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
	SetUserObject(u IUser)
	GetUserObject() IUser
	SetPostForm(formD *multipart.Form)
	GetPostForm() *multipart.Form
	GetCtx() *gin.Context
	SetCtx(ctx *gin.Context)
	SetAdminRequestParams(rp *AdminRequestParams)
	GetAdminRequestParams() *AdminRequestParams
}

type AdminContext struct {
	Err                    string
	PageTitle              string
	ErrExists              bool
	SiteName               string
	Languages              []Language
	RootURL                string
	OTPRequired            bool
	Language               *Language
	Username               string
	Password               string
	Logo                   string
	FavIcon                string
	SessionKey             string
	RootAdminURL           string
	User                   string
	UserExists             bool
	UserObject             IUser
	Demo                   bool
	UserPermissionRegistry *UserPermRegistry
	CurrentURL             string
	CurrentQuery           string
	FullURL                *url.URL
	Form                   IForm
	BreadCrumbs            *AdminBreadCrumbsRegistry
	ID                     uint
	FormD                  *multipart.Form
	Ctx                    *gin.Context
	AdminRequestParams     *AdminRequestParams
}

func (c *AdminContext) SetID(ID uint) {
	c.ID = ID
}

func (c *AdminContext) GetCtx() *gin.Context {
	return c.Ctx
}

func (c *AdminContext) SetAdminRequestParams(rp *AdminRequestParams) {
	c.AdminRequestParams = rp
}

func (c *AdminContext) GetAdminRequestParams() *AdminRequestParams {
	return c.AdminRequestParams
}

func (c *AdminContext) SetCtx(ctx *gin.Context) {
	c.Ctx = ctx
}

func (c *AdminContext) SetUserObject(u IUser) {
	c.UserObject = u
}

func (c *AdminContext) GetUserObject() IUser {
	return c.UserObject
}

func (c *AdminContext) SetPostForm(formD *multipart.Form) {
	c.FormD = formD
}

func (c *AdminContext) GetPostForm() *multipart.Form {
	return c.FormD
}

func (c *AdminContext) GetID() uint {
	return c.ID
}

func (c *AdminContext) SetBreadCrumbs(breadcrumbs *AdminBreadCrumbsRegistry) {
	c.BreadCrumbs = breadcrumbs
}

func (c *AdminContext) GetSessionKey() string {
	return c.SessionKey
}

func (c *AdminContext) SetSiteName(siteName string) {
	c.SiteName = siteName
}

func (c *AdminContext) SetCurrentURL(currentURL string) {
	c.CurrentURL = currentURL
}

func (c *AdminContext) SetCurrentQuery(currentQuery string) {
	c.CurrentQuery = currentQuery
}

func (c *AdminContext) SetForm(form IForm) {
	c.Form = form
}

func (c *AdminContext) SetFullURL(fullURL *url.URL) {
	c.FullURL = fullURL
}

func (c *AdminContext) SetRootAdminURL(rootAdminURL string) {
	c.RootAdminURL = rootAdminURL
}

func (c *AdminContext) SetSessionKey(sessionKey string) {
	c.SessionKey = sessionKey
}

func (c *AdminContext) SetRootURL(rootURL string) {
	c.RootURL = rootURL
}

func (c *AdminContext) GetRootURL() string {
	return c.RootURL
}

func (c *AdminContext) SetLanguage(language *Language) {
	c.Language = language
}

func (c *AdminContext) SetLogo(logo string) {
	c.Logo = logo
}

func (c *AdminContext) SetFavIcon(favicon string) {
	c.FavIcon = favicon
}

func (c *AdminContext) GetLanguage() *Language {
	return c.Language
}

func (c *AdminContext) SetLanguages(langs []Language) {
	c.Languages = langs
}

func (c *AdminContext) SetUserPermissionRegistry(permRegistry *UserPermRegistry) {
	c.UserPermissionRegistry = permRegistry
}

func (c *AdminContext) SetPageTitle(pageTitle string) {
	c.PageTitle = pageTitle
}

func (c *AdminContext) SetUser(user string) {
	c.User = user
}

func (c *AdminContext) SetUserExists(userExists bool) {
	c.UserExists = userExists
}

func (c *AdminContext) SetDemo() {
	c.Demo = true
}

func (c *AdminContext) SetError(err string) {
	c.Err = err
}

func (c *AdminContext) SetErrorExists() {
	c.ErrExists = true
}
