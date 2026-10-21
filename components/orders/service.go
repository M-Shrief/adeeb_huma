package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"net/http"
)

func GetAll_Service(input *GetAllOrders_Req) (*schemas.GetAll_Res_Body[schemas.Order_Descriptive], error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}
	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	result, err := GetAll_Repository(&input.GetAll_Req)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /orders")
			return nil, api_error
		}
	}
	return result, nil
}

func GetUserOrders_Service(input *GetAllOrders_Req) (*schemas.GetAll_Res_Body[schemas.Order_Descriptive], error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}
	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)

	result, err := GetUserOrders_Repository(user_id, &input.GetAll_Req)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /orders/me")
			return nil, api_error
		}
	}
	return result, nil
}

func GetOrderByID_Service(ctx context.Context, input *GetOrderByID_Req) (*OneOrder_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	order_res, err := GetOrderByID_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /orders/{id}")
			return nil, api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return nil, api_error
		}
	}

	return order_res, nil
}

func CreateOneOrder_Service(ctx context.Context, input *CreateOneOrder_Req) (*OneOrder_Res, error) {
	result, err := CreateOneOrder_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /orders")
			return nil, api_error
		}
	}
	return result, nil
}

func CreateManyOrders_Service(ctx context.Context, input *CreateManyOrders_Req) (*schemas.CreateMany_Res_Body[OneOrder_Res], error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}
	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	result, err := CreateManyOrders_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /orders/many")
			return nil, api_error
		}
	}
	return result, nil
}

func AddOnePrint_Service(ctx context.Context, input *AddOnePrint_Req) (*PrintItem_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{id}")
			return nil, api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return nil, api_error
		}
	}

	result, err := AddOnePrint_Repository(ctx, order_res, input.Body)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /orders/{id}/prints")
			return nil, api_error
		}
	}
	return result, nil
}

func AddManyPrints_Service(ctx context.Context, input *AddManyPrints_Req) (*schemas.CreateMany_Res_Body[PrintItem_Res], error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return nil, api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{id}")
			return nil, api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return nil, api_error
		}
	}

	result, err := AddManyPrints_Repository(ctx, order_res, input.Body)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /orders/{id}/prints")
			return nil, api_error
		}
	}
	return result, nil
}

func UpdateOrder_Service(ctx context.Context, input *UpdateOrder_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{id}")
			return api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return api_error
		}

		// Make sure the user can't update interanl values like: is_updateable, is_completed
		// we assign them to None, as we can exclude them later with model_dump(exclude_none=True)
		input.Body.IsUpdateable = nil
		input.Body.Status = nil
		input.Body.UserID = nil
		input.Body.Reviewed = nil
	}

	err = UpdateOrder_Repository(ctx, order_res, &input.Body)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{id}")
			return api_error
		}
	}
	return nil
}

func UpdatePrint_Service(ctx context.Context, input *UpdatePrint_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.OrderID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{order_id}/prints/{print_id}")
			return api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return api_error
		}
	}

	err = UpdatePrint_Repository(ctx, order_res, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /orders/{order_id}/prints/{print_id}")
			return api_error
		}
	}
	return nil
}

func DeleteOrder_Service(ctx context.Context, input *DeleteOrder_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /orders/{id}")
			return api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return api_error
		}
	}

	err = DeleteOrder_Repository(ctx, input.ID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /orders/{id}")
			return api_error
		}
	}
	return nil
}

func DeletePrint_Service(ctx context.Context, input *DeletePrint_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	order_res, err := GetOrderForOps_Repository(ctx, input.OrderID)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /orders/{order_id}/prints/{print_id}")
			return api_error
		}
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_res.UserID, claims) == false {
			api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
			return api_error
		}
	}

	err = DeletePrint_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /orders/{order_id}/prints/{print_id}")
			return api_error
		}
	}
	return nil
}
