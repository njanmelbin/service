package userbus

import (
	"fmt"
	"net/mail"
	"service/app/sdk/errs"
	"time"

	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID               *uuid.UUID
	Name             *string `validate:"omitempty,min=3"`
	Email            *mail.Address
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}

// Validate can perform a check of the data against the validate tags.
func (qf *QueryFilter) Validate() error {
	if err := errs.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// WithUserID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.ID = &userID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithEmail sets the Email field of the QueryFilter value.
func (qf *QueryFilter) WithEmail(email mail.Address) {
	qf.Email = &email
}

// WithStartDateCreated sets the StartCreatedDate field of the QueryFilter value.
func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

// WithEndCreatedDate sets the EndCreatedDate field of the QueryFilter value.
func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
