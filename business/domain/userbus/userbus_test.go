package userbus_test

import (
	"context"
	"fmt"
	"net/mail"
	"testing"

	"service/business/domain/userbus"

	"service/business/sdk/dbtest"
	"service/business/sdk/unitest"
	"service/business/types/name"
	"service/business/types/role"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Test_User(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_User")

	// sd, err := insertSeedData(db.BusDomain)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
	// }

	// -------------------------------------------------------------------------

	unitest.Run(t, create(db.BusDomain), "create")

}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.AdminRole, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := unitest.User{
		User: usrs[0],
	}

	tu2 := unitest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 2, role.UserRole, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := unitest.User{
		User: usrs[0],
	}

	tu4 := unitest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Users:  []unitest.User{tu3, tu4},
		Admins: []unitest.User{tu1, tu2},
	}

	return sd, nil
}

func create(busDomain dbtest.BusDomain) []unitest.Table {
	email, _ := mail.ParseAddress("bill@ardanlabs.com")

	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: userbus.User{
				Name:       name.MustParse("Bill Kennedy"),
				Email:      *email,
				Roles:      []role.Role{role.AdminRole},
				Department: "ITO",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				nu := userbus.NewUser{
					Name:       name.MustParse("Bill Kennedy"),
					Email:      *email,
					Roles:      []role.Role{role.AdminRole},
					Department: "ITO",
					Password:   "123",
				}

				resp, err := busDomain.User.Create(ctx, uuid.UUID{}, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(userbus.User)
				if !exists {
					return "error occurred"
				}

				if err := bcrypt.CompareHashAndPassword(gotResp.PasswordHash, []byte("123")); err != nil {
					return err.Error()
				}

				expResp := exp.(userbus.User)

				expResp.ID = gotResp.ID
				expResp.PasswordHash = gotResp.PasswordHash
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
