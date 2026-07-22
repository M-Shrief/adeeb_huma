package poems

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/poems",
			Summary:       "Create One",
			Description:   "Create One Poem",
			Tags:          []string{"Poems"},
			DefaultStatus: http.StatusCreated,
		},
		CreateOnePoem_Handler,
	)

}
