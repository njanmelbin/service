package userbus

import (
	"context"
	"errors"
	"fmt"
	"service/foundation/logger"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authenticaton failed")
)

type Storer interface {
	Create(ctx context.Context, usr User) error
	// Update(ctx context.Context, usr User) error
	// Delete(ctx context.Context, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

type Core struct {
	log    *logger.Logger
	storer Storer
}

func NewCore(log *logger.Logger, storer Storer) *Core {
	return &Core{
		log:    log,
		storer: storer,
	}
}

func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		Department:   nu.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}
	return usr, nil
}
