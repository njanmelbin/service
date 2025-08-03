package userapi

import (
	"service/app/sdk/apitest"
	"testing"
)

func Test_User(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_User")

	// -------------------------------------------------------------------------

	_, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}
}
