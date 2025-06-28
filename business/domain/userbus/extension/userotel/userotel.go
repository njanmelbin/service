package userotel

import (
	"context"
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
