package adeebs

import (
	"adeeb_huma/internal/auth"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/adeebs",
			Summary:       "Get All",
			Description:   "Get All Adeebs",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusOK,
		},
		GetAllAdeebs_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/adeebs/{id}",
			Summary:       "Get One",
			Description:   "Get One Adeeb",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusOK,
		},
		GetOneAdeeb_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/adeebs",
			Summary:       "Create One",
			Description:   "Create One Adeeb",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusCreated,
			Middlewares:   huma.Middlewares{auth.VerifyAdminstratorMiddleware(api)},
			Parameters: []*huma.Param{
				{
					Name:     "Authorization",
					In:       "header",
					Required: true,
					Schema: &huma.Schema{
						Type: "string",
					},
				},
			},
		},
		CreateOneAdeeb_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/adeebs/many",
			Summary:       "Create Many",
			Description:   "Create Many Adeebs",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusCreated,
			Middlewares:   huma.Middlewares{auth.VerifyAdminstratorMiddleware(api)},
			Parameters: []*huma.Param{
				{
					Name:     "Authorization",
					In:       "header",
					Required: true,
					Schema: &huma.Schema{
						Type: "string",
					},
				},
			},
		},
		CreateManyAdeeb_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/adeebs/{id}",
			Summary:       "Update One",
			Description:   "Update one Adeeb",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusNoContent,
			Middlewares:   huma.Middlewares{auth.VerifyAdminstratorMiddleware(api)},
			Parameters: []*huma.Param{
				{
					Name:     "Authorization",
					In:       "header",
					Required: true,
					Schema: &huma.Schema{
						Type: "string",
					},
				},
			},
		},
		UpdateAdeeb_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodDelete,
			Path:          "/api/v1/adeebs/{id}",
			Summary:       "Delete One",
			Description:   "Delete One Adeeb",
			Tags:          []string{"Adeebs"},
			DefaultStatus: http.StatusNoContent,
			Middlewares:   huma.Middlewares{auth.VerifyAdminstratorMiddleware(api)},
			Parameters: []*huma.Param{
				{
					Name:     "Authorization",
					In:       "header",
					Required: true,
					Schema: &huma.Schema{
						Type: "string",
					},
				},
			},
		},
		DeleteAdeeb_Handler,
	)

}
