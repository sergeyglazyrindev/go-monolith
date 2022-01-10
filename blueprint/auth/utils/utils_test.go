package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPass(t *testing.T) {
	_, err := HashPass("123456", "123456")
	assert.True(t, err == nil)
}
