package userapi

import (
	"net/http"
	"service/app/domain/userapp"
	"service/app/sdk/apitest"

	"github.com/google/go-cmp/cmp"
)

func create200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &userapp.NewUser{
				Name:            "Bill Kennedy",
				Email:           "bill@ardanlabs.com",
				Roles:           []string{"ADMIN"},
				Department:      "ITO",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &userapp.User{},
			ExpResp: &userapp.User{
				Name:       "Bill Kennedy",
				Email:      "bill@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "ITO",
				Enabled:    true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapp.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapp.User)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}
	return table
}
