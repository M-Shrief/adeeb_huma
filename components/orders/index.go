package orders

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/orders",
			Summary:       "Create One",
			Description:   "Create One Order",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusCreated,
		},
		CreateOneOrder_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/orders/many",
			Summary:       "Create Many",
			Description:   "Create Many Order",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusCreated,
		},
		CreateManyOrder_Handler,
	)
}
