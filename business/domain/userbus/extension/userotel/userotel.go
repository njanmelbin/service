package userotel

import (
	"context"
	"net/mail"
	"service/business/domain/userbus"
	"service/foundation/otel"

	"github.com/google/uuid"
)

type Extension struct {
	bus userbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the userbus with otel.
func NewExtension() userbus.Extension {
	return func(bus userbus.ExtBusiness) userbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

// Create applies otel to the user creation process.
func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, nu userbus.NewUser) (userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.create")
	defer span.End()

	usr, err := ext.bus.Create(ctx, actorID, nu)
	if err != nil {
		return userbus.User{}, err
	}

	return usr, nil
}

// Delete applies otel to the user deletion process.
func (ext *Extension) Delete(ctx context.Context, actorID uuid.UUID, usr userbus.User) error {
	ctx, span := otel.AddSpan(ctx, "business.userbus.delete")
	defer span.End()

	if err := ext.bus.Delete(ctx, actorID, usr); err != nil {
		return err
	}

	return nil
}

// QueryByEmail applies otel to the user query by email process.
func (ext *Extension) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.querybyemail")
	defer span.End()

	return ext.bus.QueryByEmail(ctx, email)
}

// Authenticate applies otel to the user authentication process.
func (ext *Extension) Authenticate(ctx context.Context, email mail.Address, password string) (userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.authenticate")
	defer span.End()

	return ext.bus.Authenticate(ctx, email, password)
}
