package userapi

import (
	"service/app/domain/userapp"
	"service/business/domain/userbus"
	"service/business/types/role"
	"time"
)

func toAppUser(bus userbus.User) userapp.User {
	return userapp.User{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Email:       bus.Email.Address,
		Roles:       role.ParseToString(bus.Roles),
		Department:  bus.Department,
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(bus userbus.User) *userapp.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
