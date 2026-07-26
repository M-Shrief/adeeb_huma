package prose_qoutes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {
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
}
