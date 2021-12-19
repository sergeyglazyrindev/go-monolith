package gomonolith

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"html/template"
	"os"
	"testing"
)

func TestEmbeddingTemplates(t *testing.T) {
	app := NewApp(os.Getenv("TEST_ENVIRONMENT"))
	t1, _ := template.ParseFS(app.Config.TemplatesFS, "templates/test.html")
	templateBuffer := &bytes.Buffer{}
	t1.Execute(templateBuffer, struct {
		Title string
	}{Title: "test"})
	assert.Contains(t, templateBuffer.String(), "test")
}
