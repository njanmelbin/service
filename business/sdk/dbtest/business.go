package dbtest

import (
	"service/business/domain/userbus"
	"service/business/domain/userbus/stores/userdb"
	"service/business/sdk/delegate"
	"service/foundation/logger"

	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate

	User userbus.ExtBusiness
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {

	delegate := delegate.New(log)

	userBus := userbus.NewBusiness(log, delegate, userdb.NewStore(log, db))

	return BusDomain{
		User: userBus,
	}
}
