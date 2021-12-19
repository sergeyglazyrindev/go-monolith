package forms

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

type UsernameFormOptions struct {
	core.FieldFormOptions
}

func NewTestForm2() *multipart.Form {
	form1 := multipart.Form{
		Value: make(map[string][]string),
	}
	return &form1
}

type FormTestSuite struct {
	gomonolith.TestSuite
}

func (s *FormTestSuite) TestFormBuilder() {
	fieldChoiceRegistry := core.FieldChoiceRegistry{}
	fieldChoiceRegistry.Choices = make([]*core.FieldChoice, 0)
	formOptions := &UsernameFormOptions{
		FieldFormOptions: core.FieldFormOptions{
			Name:        "UsernameOptions",
			Initial:     "InitialUsername",
			DisplayName: "Display name",
			Validators:  core.NewValidatorRegistry(),
			Choices:     &fieldChoiceRegistry,
			HelpText:    "help for username",
		},
	}
	core.FormCongirurableOptionInstance.AddFieldFormOptions(formOptions)
	// initial=\"test\",displayname=\"uname\",validators=\"password-go-monolith\",choices=UsernameChoices,helptext=\"HELPPPPPPPPPP\"
	user := &core.User{}
	form1 := core.NewFormFromModel(user, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin", "ExpiresOn", "OTPRequired"}, true, "")
	res := form1.Render()
	assert.Contains(s.T(), res, "<form")
	form2 := NewTestForm2()
	form2.Value["Username"] = []string{"username"}
	form2.Value["FirstName"] = []string{"first name"}
	form2.Value["LastName"] = []string{"last name"}
	form2.Value["Email"] = []string{"email@example.com"}
	form2.Value["OTPRequired"] = []string{"yes"}
	formError := form1.ProceedRequest(form2, user, &core.AdminContext{})
	assert.Equal(s.T(), user.Username, "username")
	assert.True(s.T(), formError.IsEmpty())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestForm(t *testing.T) {
	gomonolith.RunTests(t, new(FormTestSuite))
}
