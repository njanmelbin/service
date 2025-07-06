package userapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"service/app/sdk/errs"
	"service/business/domain/userbus"
	"time"
)

// User represents information about an individual user.
type User struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppUser(usr userbus.User) User {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.String()
	}

	return User{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []User {
	app := make([]User, len(users))
	for i, usr := range users {
		app[i] = toAppUser(usr)
	}

	return app
}

// NewUser defines the data needed to add a new user.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

func toBusNewUser(app NewUser) (userbus.NewUser, error) {
	roles := make([]userbus.Role, len(app.Roles))
	for i, roleStr := range app.Roles {
		role, err := userbus.Parse(roleStr)
		if err != nil {
			return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
		}
		roles[i] = role
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		Name:       app.Name,
		Email:      *addr,
		Roles:      roles,
		Department: app.Department,
		Password:   app.Password,
	}

	return bus, nil
}

// Decode implements the decoder interface.
func (app *NewUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}
