package chosen_verses

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/chosen_verses",
			Summary:       "Create One",
			Description:   "Create One ChosenVerse",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusCreated,
		},
		CreateOneChosenVerse_Handler,
	)

}
