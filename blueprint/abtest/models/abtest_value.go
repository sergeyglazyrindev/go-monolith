package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"io/ioutil"
	"strings"
)

// ABTestValue is a model to represent a possible value of an AB test
type ABTestValue struct {
	core.Model
	ABTest      ABTest
	ABTestID    uint
	Value       string `gomonolith:"inline"`
	Active      bool   `gorm:"default:false" gomonolith:"inline"`
	Impressions int    `gomonolith:"inline" gorm:"default:0"`
	Clicks      int    `gomonolith:"inline" gorm:"default:0"`
}

func (a *ABTestValue) String() string {
	return fmt.Sprintf("ABTest Value %s", a.Value)
}

// ClickThroughRate returns the rate of click through of this value
func (a *ABTestValue) ClickThroughRate() string {
	if a.Impressions == 0 {
		return "0.0"
	}
	return fmt.Sprintf("%.2f", float64(a.Clicks)/float64(a.Impressions)*100)
}

// Preview__Form__List shows a preview of the AB test's value
func (a ABTestValue) PreviewFormList() string {
	// Check if the value is a path to a file
	if strings.HasPrefix(a.Value, "/") {
		// Check the file type
		// Image
		if strings.HasSuffix(a.Value, "png") || strings.HasSuffix(a.Value, "jpg") || strings.HasSuffix(a.Value, "gif") || strings.HasSuffix(a.Value, "jpeg") {
			return fmt.Sprintf(`<img src="%s" style="width:256px">`, a.Value)
		}
		// CSS/JS
		if strings.HasSuffix(a.Value, "css") || strings.HasSuffix(a.Value, "js") {
			buf, _ := ioutil.ReadFile("." + a.Value)
			return fmt.Sprintf(`<pre style="width:256px">%s\n%s</pre>`, a.Value, string(buf))
		}
	}
	return a.Value
}
