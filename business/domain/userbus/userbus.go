package userbus

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"service/business/sdk/delegate"
	"service/business/sdk/sqldb"
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
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, usr User) error
	// Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core busines logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error)
	// Update(ctx context.Context, actorID uuid.UUID, usr User, uu UpdateUser) (User, error)
	Delete(ctx context.Context, actorID uuid.UUID, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	Authenticate(ctx context.Context, email mail.Address, password string) (User, error)
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

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:      b.log,
		delegate: b.delegate,
		storer:   storer,
	}

	return &bus, nil
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

// Delete removes the specified user.
func (b *Business) Delete(ctx context.Context, actorID uuid.UUID, usr User) error {
	if err := b.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// Other domains may need to know when a user is deleted so business
	// logic can be applied. This represents a delegate call to other domains.
	if err := b.delegate.Call(ctx, ActionDeletedData(usr.ID)); err != nil {
		return fmt.Errorf("failed to execute `%s` action: %w", ActionDeleted, err)
	}

	return nil
}

// QueryByEmail finds the user by a specified user email.
func (b *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := b.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := b.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
