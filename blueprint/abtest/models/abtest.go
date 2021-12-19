package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
)

// ABTestType is the type of the AB testing: model or static
type ABTestType int

// Static is used to do AB testing for static assets (images, js, css, ...)
func (ABTestType) Static() ABTestType {
	return 1
}

// Model is used to do AB testing for model values coming from database
func (ABTestType) Model() ABTestType {
	return 2
}

func HumanizeAbTestType(abTestType ABTestType) string {
	switch abTestType {
	case 1:
		return "static"
	case 2:
		return "model"
	default:
		return "unknown"
	}
}

// ABTest is a model that stores an A/B test
type ABTest struct {
	core.Model
	Name          string           `gomonolithform:"RequiredFieldOptions" gomonolith:"list"`
	Type          ABTestType       `gomonolithform:"RequiredSelectFieldOptions" gomonolith:"list"`
	StaticPath    string           `gomonolith:"list"`
	ContentType   core.ContentType `gomonolith:"list" gomonolithform:"ContentTypeFieldOptions"`
	ContentTypeID uint
	Field         string `gomonolith:"list" gomonolithform:"RequiredSelectFieldOptions"`
	PrimaryKey    uint   `gomonolith:"list" gorm:"default:0"`
	Active        bool   `gorm:"default:false" gomonolith:"list"`
	Group         string `gomonolith:"list"`
}

func (m *ABTest) String() string {
	return fmt.Sprintf("ABTest %s", m.Name)
}
