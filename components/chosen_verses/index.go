package chosen_verses

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/chosen_verses",
			Summary:       "Get All",
			Description:   "Get All ChosenVersess",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusOK,
		},
		GetAllChosenVerses_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/chosen_verses/{id}",
			Summary:       "Get One",
			Description:   "Get One ChosenVersess",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusOK,
		},
		GetOneChosenVerse_Handler,
	)

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

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/chosen_verses/many",
			Summary:       "Create Many",
			Description:   "Create Many ChosenVerses",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusCreated,
		},
		CreateManyChosenVerses_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPut,
			Path:          "/api/v1/chosen_verses/{id}",
			Summary:       "Update One",
			Description:   "Update One ChosenVerse",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusNoContent,
		},
		UpdateChosenVerse_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodDelete,
			Path:          "/api/v1/chosen_verses/{id}",
			Summary:       "Delete One",
			Description:   "Delete One ChosenVerse",
			Tags:          []string{"ChosenVerses"},
			DefaultStatus: http.StatusNoContent,
		},
		DeleteChosenVerseHandler,
	)

}
