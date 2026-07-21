package adeebs

import (
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
		},
		CreateManyAdeeb_Handler,
	)

}
