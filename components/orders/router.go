package orders

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllOrders_Handler(ctx context.Context, input *GetAllOrders_Req) (*schemas.GetAll_Res[schemas.Order_Descriptive], error) {
	result, err := GetAll_Service(input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /adeebs")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.GetAll_Res[schemas.Order_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetUserOrders_Handler(ctx context.Context, input *GetAllOrders_Req) (*schemas.GetAll_Res[schemas.Order_Descriptive], error) {
	result, err := GetUserOrders_Service(input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /adeebs")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	res := &schemas.GetAll_Res[schemas.Order_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOrderByID_Handler(ctx context.Context, input *GetOrderByID_Req) (*GetOrderByID_Res, error) {
	result, err := GetOrderByID_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Order's not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /users/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &GetOrderByID_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneOrder_Handler(ctx context.Context, input *CreateOneOrder_Req) (*CreateOneOrder_Res, error) {
	result, err := CreateOneOrder_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &CreateOneOrder_Res{Body: *result, Status: http.StatusCreated}, nil
}

func CreateManyOrder_Handler(ctx context.Context, input *CreateManyOrders_Req) (*CreateManyOrders_Res, error) {
	result, err := CreateManyOrders_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &CreateManyOrders_Res{
		Body:   *result,
		Status: http.StatusCreated}, nil
}

func AddOnePrint_Handler(ctx context.Context, input *AddOnePrint_Req) (*AddOnePrint_Res, error) {
	result, err := AddOnePrint_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &AddOnePrint_Res{Body: *result, Status: http.StatusCreated}, nil
}

func AddManyPrints_Handler(ctx context.Context, input *AddManyPrints_Req) (*AddManyPrints_Res, error) {
	result, err := AddManyPrints_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &AddManyPrints_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil
}

func UpdateOrder_Handler(ctx context.Context, input *UpdateOrder_Req) (*schemas.Update_Res, error) {
	err := UpdateOrder_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Foreign Key Error")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Order is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &schemas.Update_Res{Status: http.StatusNoContent}, nil
}

func UpdatePrint_Handler(ctx context.Context, input *UpdatePrint_Req) (*schemas.Update_Res, error) {
	err := UpdatePrint_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Foreign Key Error")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Print is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &schemas.Update_Res{Status: http.StatusNoContent}, nil
}

func DeleteOrder_Handler(ctx context.Context, input *DeleteOrder_Req) (*schemas.Delete_Res, error) {
	err := DeleteOrder_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Order is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &schemas.Delete_Res{Status: http.StatusNoContent}, nil
}

func DeletePrint_Handler(ctx context.Context, input *DeletePrint_Req) (*schemas.Delete_Res, error) {
	err := DeletePrint_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Order is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &schemas.Delete_Res{Status: http.StatusNoContent}, nil
}
