package userbus

import (
	"context"
	"errors"
	"net/mail"
	"service/business/sdk/order"
	"service/foundation/logger"

	"github.com/google/uuid"
)

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authenticaton failed")
)

type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
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
	return User{}, nil
}
