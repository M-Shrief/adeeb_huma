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

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/users",
			Summary:       "Get All",
			Description:   "Get All Users",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
		},
		GetAllUsers_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/users/me",
			Summary:       "Get Current User",
			Description:   "Get Current User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
		},
		GetCurrentUser_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/users/{id}",
			Summary:       "Get User by ID",
			Description:   "Get User by ID",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusOK,
		},
		GetUserByID_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/users/me",
			Summary:       "Update Current User",
			Description:   "Update Current User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusNoContent,
		},
		UpdateCurrentUser_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/users/{id}",
			Summary:       "Update User By ID",
			Description:   "Update User By ID",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusNoContent,
		},
		UpdateUserByID_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/users/{id}/ban",
			Summary:       "Ban User By ID",
			Description:   "Ban User By ID",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusNoContent,
		},
		BanUserByID_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodDelete,
			Path:          "/api/v1/users/me",
			Summary:       "Delete Current User",
			Description:   "Delete Current User",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusNoContent,
		},
		DeleteCurrentUser_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodDelete,
			Path:          "/api/v1/users/{id}",
			Summary:       "Delete User By ID",
			Description:   "Delete User By ID",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusNoContent,
		},
		DeleteUserByID_Handler,
	)

}
