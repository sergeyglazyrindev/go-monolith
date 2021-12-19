package forms

import (
	"bytes"
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type WidgetTestSuite struct {
	gomonolith.TestSuite
}

func NewTestForm() *multipart.Form {
	form1 := multipart.Form{
		Value: make(map[string][]string),
	}
	return &form1
}

func (w *WidgetTestSuite) TestTextWidget() {
	textWidget := &core.TextWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	textWidget.SetName("dsadas")
	textWidget.SetValue("dsadas")
	textWidget.SetRequired()
	renderedWidget := textWidget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	err := textWidget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test"}
	err = textWidget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestNumberWidget() {
	widget := &core.NumberWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"test"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"121"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestEmailWidget() {
	widget := &core.EmailWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"test@example.com"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestURLWidget() {
	widget := &core.URLWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"example.com"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestPasswordWidget() {
	widget := &core.PasswordWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "type=\"password\"")
	widget.SetRequired()
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"12345678901234567890"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestHiddenWidget() {
	widget := &core.HiddenWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas<>")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas&lt;&gt;\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasas"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestDateWidget() {
	widget := &core.DateWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("11/01/2021")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "datetimepicker_dsadas")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"11/02/2021"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestDateTimeWidget() {
	widget := &core.DateTimeWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("11/02/2021 10:04")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"11/02/2021 10:04\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"11/02/2021 10:04"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestTimeWidget() {
	widget := &core.TimeWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("15:05")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"15:05\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestTextareaWidget() {
	widget := &core.TextareaWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<textarea name=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestCheckboxWidget() {
	widget := &core.CheckboxWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "checked=\"checked\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), widget.GetOutputValue() == true)
}

func (w *WidgetTestSuite) TestSelectWidget() {
	widget := &core.SelectWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	widget.OptGroups = make(map[string][]*core.SelectOptGroup)
	widget.OptGroups["test"] = make([]*core.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test1",
		Value:    "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test2",
		Value:    "dsadas",
	})
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"dsadas"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestNullBooleanWidget() {
	widget := &core.NullBooleanWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.OptGroups = make(map[string][]*core.SelectOptGroup)
	widget.OptGroups["test"] = make([]*core.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test1",
		Value:    "yes",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test2",
		Value:    "no",
	})
	widget.SetName("dsadas")
	widget.SetValue("yes")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas\" data-placeholder=\"Select\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"no"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSelectMultipleWidget() {
	widget := &core.SelectMultipleWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas"})
	widget.OptGroups = make(map[string][]*core.SelectOptGroup)
	widget.OptGroups["test"] = make([]*core.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test1",
		Value:    "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.SelectOptGroup{
		OptLabel: "test2",
		Value:    "dsadas",
	})
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestRadioSelectWidget() {
	widget := &core.RadioSelectWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		ID:        "test",
		WrapLabel: true,
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	widget.OptGroups = make(map[string][]*core.RadioOptGroup)
	widget.OptGroups["test"] = make([]*core.RadioOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.RadioOptGroup{
		OptLabel: "test1",
		Value:    "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.RadioOptGroup{
		OptLabel: "test2",
		Value:    "dsadas",
	})
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<li>test<ul id=\"test_0\">")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestCheckboxSelectMultipleWidget() {
	widget := &core.CheckboxSelectMultipleWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		ID:        "test",
		WrapLabel: true,
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas"})
	widget.OptGroups = make(map[string][]*core.RadioOptGroup)
	widget.OptGroups["test"] = make([]*core.RadioOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.RadioOptGroup{
		OptLabel: "test1",
		Value:    "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &core.RadioOptGroup{
		OptLabel: "test2",
		Value:    "dsadas",
	})
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<ul id=\"test\">\n  \n  \n  \n    <li>test<ul id=\"test_0\">")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestFileWidget() {
	widget := &core.FileWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "type=\"file\"")
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	err := writer.SetBoundary("foo")
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	path := os.Getenv("GOMONOLITH_PATH") + "/tests/file_for_uploading.txt"
	file, err := os.Open(path)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("GOMONOLITH_PATH"), "upload-for-tests"), 0755)
	if err != nil {
		assert.True(w.T(), false, "Couldnt create directory for file uploading", err)
		return
	}
	defer file.Close()
	defer os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("GOMONOLITH_PATH"), "upload-for-tests"))
	part, err := writer.CreateFormFile("dsadas", filepath.Base(path))
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = writer.Close()
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	form1, _ := multipart.NewReader(bytes.NewReader(body.Bytes()), "foo").ReadForm(1000000)
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestClearableFileWidget() {
	widget := &core.ClearableFileWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		InitialText:        "test",
		Required:           true,
		ID:                 "test",
		ClearCheckboxLabel: "clear file",
		InputText:          "upload your image",
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<p class=\"file-upload\">test:")
	assert.Contains(w.T(), renderedWidget, "upload your image:")
	//     <input type="file" name="dsadas" test="test1"></p>
	widget = &core.ClearableFileWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		InitialText:        "test",
		Required:           true,
		ID:                 "test",
		ClearCheckboxLabel: "clear file",
		InputText:          "upload your image",
		CurrentValue:       &core.URLValue{URL: "https://microsoft.com"},
	}
	widget.SetName("dsadas")
	renderedWidget = widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<input type=\"file\" ")
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	err := writer.SetBoundary("foo")
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	path := os.Getenv("GOMONOLITH_PATH") + "/tests/file_for_uploading.txt"
	file, err := os.Open(path)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("GOMONOLITH_PATH"), "upload-for-tests"), 0755)
	if err != nil {
		assert.True(w.T(), false, "Couldnt create directory for file uploading", err)
		return
	}
	defer file.Close()
	defer os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("GOMONOLITH_PATH"), "upload-for-tests"))
	part, err := writer.CreateFormFile("dsadas", filepath.Base(path))
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = writer.Close()
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	form1, _ := multipart.NewReader(bytes.NewReader(body.Bytes()), "foo").ReadForm(1000000)
	err = widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestMultipleHiddenInputWidget() {
	widget := &core.MultipleInputHiddenWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas", "test1"})
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadas", "test1"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSplitDateTimeWidget() {
	widget := &core.SplitDateTimeWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		DateAttrs:  map[string]string{"test": "test1"},
		TimeAttrs:  map[string]string{"test": "test1"},
		TimeFormat: "15:04",
		DateFormat: "Mon Jan _2",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_date\"")
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_time\"")
	form1 := NewTestForm()
	form1.Value["dsadas_date"] = []string{"Mon Jan 12"}
	form1.Value["dsadas_time"] = []string{"10:20"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSplitHiddenDateTimeWidget() {
	widget := &core.SplitHiddenDateTimeWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		DateAttrs:  map[string]string{"test": "test1"},
		TimeAttrs:  map[string]string{"test": "test1"},
		TimeFormat: "15:04",
		DateFormat: "Mon Jan _2",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_date\"")
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_time\"")
	form1 := NewTestForm()
	form1.Value["dsadas_date"] = []string{"Mon Jan 12"}
	form1.Value["dsadas_time"] = []string{"10:20"}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSelectDateWidget() {
	widget := &core.SelectDateWidget{
		Widget: core.Widget{
			Attrs:       map[string]string{"test": "test1"},
			BaseFuncMap: core.FuncMap,
		},
		EmptyLabelString: "choose any",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render(core.NewFormRenderContext(), nil)
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_month\"")
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_day\"")
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_year\"")
	form1 := NewTestForm()
	form1.Value["dsadas_month"] = []string{"1"}
	form1.Value["dsadas_day"] = []string{"1"}
	form1.Value["dsadas_year"] = []string{strconv.Itoa(time.Now().Year())}
	err := widget.ProceedForm(form1, nil, nil)
	assert.True(w.T(), err == nil)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestWidget(t *testing.T) {
	gomonolith.RunTests(t, new(WidgetTestSuite))
}
