package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"

	"time"
)

// ApprovalAction is a selection of approval actions
type ApprovalAction int

// Approved is an accepted change
func (ApprovalAction) Approved() ApprovalAction {
	return 1
}

// Rejected is a rejected change
func (ApprovalAction) Rejected() ApprovalAction {
	return 2
}

func HumanizeApprovalAction(approvalAction ApprovalAction) string {
	switch approvalAction {
	case 1:
		return "approved"
	case 2:
		return "rejected"
	default:
		return "unknown"
	}
}

// Approval is a model that stores approval data
type Approval struct {
	core.Model
	ApprovalAction      ApprovalAction   `gomonolith:"list" gomonolithform:"SelectFieldOptions"`
	ApprovalBy          string           `gomonolith:"list" gomonolithform:"ReadonlyField"`
	ApprovalDate        *time.Time       `gomonolith:"list" gomonolithform:"DatetimeReadonlyFieldOptions"`
	ContentType         core.ContentType `gomonolith:"list" gomonolithform:"ReadonlyField"`
	ContentTypeID       uint
	ModelPK             uint      `gomonolith:"list" gomonolithform:"ReadonlyField" gorm:"default:0"`
	ColumnName          string    `gomonolith:"list" gomonolithform:"ReadonlyField"`
	OldValue            string    `gomonolith:"list" gomonolithform:"ReadonlyField"`
	NewValue            string    `gomonolith:"list"`
	NewValueDescription string    `gomonolith:"list" gomonolithform:"ReadonlyField"`
	ChangedBy           string    `gomonolith:"list" gomonolithform:"ReadonlyField"`
	ChangeDate          time.Time `gomonolith:"list" gomonolithform:"DatetimeReadonlyFieldOptions"`
	UpdatedBy           string
}

func (a *Approval) String() string {
	return fmt.Sprintf("Approval for %s.%s %d", a.ContentType.ModelName, a.ColumnName, a.ModelPK)
}
