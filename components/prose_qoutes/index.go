package prose_qoutes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/prose_qoutes",
			Summary:       "Get All",
			Description:   "Get All ProseQoutess",
			Tags:          []string{"ProseQoutes"},
			DefaultStatus: http.StatusOK,
		},
		GetAllProseQoutes_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/prose_qoutes/{id}",
			Summary:       "Get One",
			Description:   "Get One ProseQoutess",
			Tags:          []string{"ProseQoutes"},
			DefaultStatus: http.StatusOK,
		},
		GetOneProseQoute_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/prose_qoutes",
			Summary:       "Create One",
			Description:   "Create One ProseQoute",
			Tags:          []string{"ProseQoutes"},
			DefaultStatus: http.StatusCreated,
		},
		CreateOneProseQoute_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/prose_qoutes/many",
			Summary:       "Create Many",
			Description:   "Create Many ProseQoutes",
			Tags:          []string{"ProseQoutes"},
			DefaultStatus: http.StatusCreated,
		},
		CreateManyProseQoutes_Handler,
	)

}
