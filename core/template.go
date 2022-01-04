package core

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

type ITemplateRenderer interface {
	AddFuncMap(funcName string, concreteFunc interface{})
	Render(ctx *gin.Context, path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap)
	RenderAsString(path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap) template.HTML
}

type IncludeContext struct {
	SiteName  string
	PageTitle string
}

type TemplateRenderer struct {
	funcMap   template.FuncMap
	pageTitle string
}

func (tr *TemplateRenderer) AddFuncMap(funcName string, concreteFunc interface{}) {
	tr.funcMap[funcName] = concreteFunc
}

func (tr *TemplateRenderer) Render(ctx *gin.Context, path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap) {
	Include := func(funcs1 template.FuncMap) func(templateName string, data1 ...interface{}) template.HTML {
		return func(templateName string, data1 ...interface{}) template.HTML {
			data2 := data
			if len(data1) == 1 {
				data2 = data1[0]
			}
			for rendererTemplateFuncName, rendererTemplateFunc := range tr.funcMap {
				funcs1[rendererTemplateFuncName] = rendererTemplateFunc
			}
			return tr.RenderAsString(CurrentConfig.GetPathToTemplate(templateName), data2, baseFuncMap, funcs1)
		}
	}
	PageTitle := func() string {
		return fmt.Sprintf("%s - %s", CurrentConfig.D.GoMonolith.SiteName, tr.pageTitle)
	}
	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = make(template.FuncMap)
		funcs1["PageTitle"] = PageTitle
	} else {
		funcs1 = funcs[0]
		funcs1["PageTitle"] = PageTitle
	}
	for rendererTemplateFuncName, rendererTemplateFunc := range tr.funcMap {
		funcs1[rendererTemplateFuncName] = rendererTemplateFunc
	}
	funcs1["Include"] = Include(funcs1)
	RenderHTML(ctx, path, data, baseFuncMap, funcs1)
}

func (tr *TemplateRenderer) RenderAsString(path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap) template.HTML {
	Include := func(funcs1 template.FuncMap) func(templateName string, data1 ...interface{}) template.HTML {
		return func(templateName string, data1 ...interface{}) template.HTML {
			data2 := data
			if len(data1) == 1 {
				data2 = data1[0]
			}
			for rendererTemplateFuncName, rendererTemplateFunc := range tr.funcMap {
				funcs1[rendererTemplateFuncName] = rendererTemplateFunc
			}
			return tr.RenderAsString(CurrentConfig.GetPathToTemplate(templateName), data2, baseFuncMap, funcs1)
		}
	}
	PageTitle := func() string {
		return fmt.Sprintf("%s - %s", CurrentConfig.D.GoMonolith.SiteName, tr.pageTitle)
	}
	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = make(template.FuncMap)
		funcs1["PageTitle"] = PageTitle
	} else {
		funcs1 = funcs[0]
		funcs1["PageTitle"] = PageTitle
	}
	for rendererTemplateFuncName, rendererTemplateFunc := range tr.funcMap {
		funcs1[rendererTemplateFuncName] = rendererTemplateFunc
	}
	funcs1["Include"] = Include(funcs1)
	templateWriter := bytes.NewBuffer([]byte{})
	RenderHTMLAsString(templateWriter, path, data, baseFuncMap, funcs1)
	return template.HTML(templateWriter.String())
}

func NewTemplateRenderer(pageTitle string) ITemplateRenderer {
	templateRenderer := TemplateRenderer{funcMap: template.FuncMap{}, pageTitle: pageTitle}
	return &templateRenderer
}

// RenderHTML creates a new template and applies a parsed template to the specified
// data object. For function, Tf is available by default and if you want to add functions
//to your template, just add them to funcs which will add them to the template with their
// original function names. If you added anonymous functions, they will be available in your
// templates as func1, func2 ...etc.
func RenderHTML(ctx *gin.Context, path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap) error {
	var err error

	var funcs1 template.FuncMap
	if len(funcs) == 0 {
		funcs1 = make(template.FuncMap)
	} else {
		funcs1 = funcs[0]
	}
	for k, v := range baseFuncMap {
		funcs1[k] = v
	}
	//// Check for ABTesting cookie
	//if cookie, err := ctx.Cookie("abt"); err != nil || cookie == nil {
	//	now := time.Now().AddDate(0, 0, 1)
	//	cookie1 := &http.Cookie{
	//		Name:    "abt",
	//		Value:   fmt.Sprint(now.Second()),
	//		Path:    "/",
	//		Expires: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	//	}
	//	http.SetCookie(ctx.Writer, cookie1)
	//}
	templateNameParts := strings.Split(path, "/")
	templateContent := CurrentConfig.GetTemplateContent(path)
	templateIdentifier := strings.Join(templateNameParts, "-")
	newT, err := template.New(templateIdentifier).Funcs(funcs1).Parse(
		string(templateContent),
	)
	if err != nil {
		// ctx.AbortWithStatus(500)
		ctx.String(500, err.Error())
		Trail(ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	err = newT.Execute(ctx.Writer, data)
	if err != nil {
		// ctx.AbortWithStatus(500)
		ctx.String(500, err.Error())
		Trail(ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	return nil
}

// RenderHTML creates a new template and applies a parsed template to the specified
// data object. For function, Tf is available by default and if you want to add functions
//to your template, just add them to funcs which will add them to the template with their
// original function names. If you added anonymous functions, they will be available in your
// templates as func1, func2 ...etc.
func RenderHTMLAsString(writer *bytes.Buffer, path string, data interface{}, baseFuncMap template.FuncMap, funcs ...template.FuncMap) error {
	var err error

	var funcs1 template.FuncMap
	//Include := func(funcs1 template.FuncMap) func (templateName string, data1 ...interface{}) string {
	//	return func (templateName string, data1 ...interface{}) string {
	//		templateWriter := bytes.NewBuffer([]byte{})
	//		data2 := data
	//		if len(data1) == 1 {
	//			data2 = data1[0]
	//		}
	//		err := RenderHTMLAsString(templateWriter, fsys, CurrentConfig.GetPathToTemplate(templateName), data2, baseFuncMap, funcs1)
	//		if err != nil {
	//			Trail(CRITICAL, "Error while parsing include of the template %s", templateName)
	//			panic(err)
	//		}
	//		return templateWriter.String()
	//	}
	//}
	if len(funcs) == 0 {
		funcs1 = make(template.FuncMap)
	} else {
		funcs1 = funcs[0]
	}
	for k, v := range baseFuncMap {
		funcs1[k] = v
	}
	for k, v := range FuncMap {
		_, funcExists := funcs1[k]
		if !funcExists {
			funcs1[k] = v
		}
	}
	//includeToKeep, includeToKeepExists := funcs1["IncludeToKeep"]
	//if includeToKeepExists {
	//	funcs1["Include"] = includeToKeep
	//	// delete(funcs1, "IncludeToKeep")
	//} else {
	//	funcs1["Include"] = Include(funcs1)
	//}

	//// Check for ABTesting cookie
	//if cookie, err := ctx.Cookie("abt"); err != nil || cookie == nil {
	//	now := time.Now().AddDate(0, 0, 1)
	//	cookie1 := &http.Cookie{
	//		Name:    "abt",
	//		Value:   fmt.Sprint(now.Second()),
	//		Path:    "/",
	//		Expires: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	//	}
	//	http.SetCookie(ctx.Writer, cookie1)
	//}
	templateNameParts := strings.Split(path, "/")
	templateContent := CurrentConfig.GetTemplateContent(path)
	templateIdentifier := strings.Join(templateNameParts, "-")
	newT, err := template.New(templateIdentifier).Funcs(funcs1).Parse(
		string(templateContent),
	)
	if err != nil {
		// ctx.AbortWithStatus(500)
		Trail(ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	err = newT.Execute(writer, data)
	if err != nil {
		// ctx.AbortWithStatus(500)
		Trail(ERROR, "RenderHTML unable to parse %s. %s", path, err)
		return err
	}
	return nil
}
