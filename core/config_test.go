package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {
	config := NewConfig("configs/test_sqlite.yml")
	assert.Equal(t, config.D.GoMonolith.Theme, "default")
}
