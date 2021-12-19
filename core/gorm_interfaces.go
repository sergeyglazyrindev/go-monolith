package core

import (
	"gorm.io/gorm"
)

// Model is the standard struct to be embedded
// in any other struct to make it a model for go-monolith
type Model struct {
	gorm.Model
}

func (m *Model) GetID() uint {
	return m.ID
}
