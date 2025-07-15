package userbus

import (
	"context"
	"errors"
	"fmt"
	"service/business/sdk/delegate"
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

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core busines logic.
type ExtBusiness interface {
	//NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error)
	// Update(ctx context.Context, actorID uuid.UUID, usr User, uu UpdateUser) (User, error)
	// Delete(ctx context.Context, actorID uuid.UUID, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	// Authenticate(ctx context.Context, email mail.Address, password string) (User, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

type Business struct {
	log      *logger.Logger
	storer   Storer
	delegate *delegate.Delegate
}

func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		log:      log,
		delegate: delegate,
		storer:   storer,
	})

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

func (b *Business) Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error) {
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

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
