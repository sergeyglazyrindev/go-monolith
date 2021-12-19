package utils

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"testing"
)

type TemplateTestSuite struct {
	gomonolith.TestSuite
}

func (s *TemplateTestSuite) TestRenderHTML() {
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTemplate(t *testing.T) {
	gomonolith.RunTests(t, new(TemplateTestSuite))
}
