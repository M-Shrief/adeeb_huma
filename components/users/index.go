package users

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/users/signup",
			Summary:       "User Signup",
			Description:   "User signing up.",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusCreated,
		},
		Signup_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/users/login",
			Summary:       "User Login",
			Description:   "User Login.",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
		},
		Login_Handler,
	)

}
