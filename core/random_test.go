package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateCSRFToken(t *testing.T) {
	token := GenerateCSRFToken()
	assert.Equal(t, len(token), 32)
	tokenmasked := MaskCSRFToken(token)
	assert.Equal(t, len(tokenmasked), 64)
	tokenunmasked := UnmaskCSRFToken(tokenmasked)
	assert.Equal(t, tokenunmasked, token)
}
