package dbtest

import (
	"service/business/domain/userbus"
	"service/business/domain/userbus/stores/userdb"
	"service/foundation/logger"

	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	User userbus.ExtBusiness
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))

	return BusDomain{
		User: userBus,
	}
}
