package orders

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(api huma.API) {

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/orders",
			Summary:       "Get All",
			Description:   "Get All Order",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusOK,
		},
		GetAllOrders_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/orders/me",
			Summary:       "Get User orders",
			Description:   "Get User Order",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusOK,
		},
		GetUserOrders_Handler,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/api/v1/orders/{id}",
			Summary:       "Get One",
			Description:   "Get One Order",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusOK,
		},
		GetOrderByID_Handler,
	)

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

	huma.Register(
		api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/api/v1/orders/{id}/prints",
			Summary:       "Add Print",
			Description:   "Add One Print",
			Tags:          []string{"Orders"},
			DefaultStatus: http.StatusCreated,
		},
		AddPrint_Handler,
	)

}
