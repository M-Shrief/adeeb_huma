package poems

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
			Path:          "/api/v1/poems",
			Summary:       "Get All",
			Description:   "Get All Poemss",
			Tags:          []string{"Poems"},
			DefaultStatus: http.StatusOK,
		},
		GetAllPoems_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/poems/{id}",
			Summary:       "Get One",
			Description:   "Get One Poemss",
			Tags:          []string{"Poems"},
			DefaultStatus: http.StatusOK,
		},
		GetOnePoem_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/poems",
			Summary:       "Create One",
			Description:   "Create One Poem",
			Tags:          []string{"Poems"},
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
		CreateOnePoem_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/poems/many",
			Summary:       "Create Many",
			Description:   "Create Many Poems",
			Tags:          []string{"Poems"},
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
		CreateManyPoems_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/poems/{id}",
			Summary:       "Update One",
			Description:   "Update One Poem",
			Tags:          []string{"Poems"},
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
		UpdatePoem_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodDelete,
			Path:          "/api/v1/poems/{id}",
			Summary:       "Delete One",
			Description:   "Delete One Poem",
			Tags:          []string{"Poems"},
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
		DeletePoemHandler,
	)

}
