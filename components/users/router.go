package users

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllUsers_Handler(ctx context.Context, input *GetAllUsers_Req) (*schemas.GetAll_Res[schemas.User_Descriptive], error) {
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
			logger.Error().Err(api_error).Msg("Error in GET /users")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.GetAll_Res[schemas.User_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetCurrentUser_Handler(ctx context.Context, input *GetCurrentUser_Req) (*GetOneUser_Res, error) {
	result, err := GetCurrentUser_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("User's not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	res := &GetOneUser_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetUserByID_Handler(ctx context.Context, input *GetUserByID_Req) (*GetOneUser_Res, error) {
	result, err := GetOneById_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("User's not found")
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
	res := &GetOneUser_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func Signup_Handler(ctx context.Context, input *Signup_Req) (*UserAuthorized_Res, error) {
	result, err := Signup_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Username already exists")
			case http.StatusBadRequest:
				if api_error.Message == "" {
					api_error.Message = "Bad request, try again later"
				}
				return nil, huma.Error400BadRequest(api_error.Message)
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /users/signup")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &UserAuthorized_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil
}

func Login_Handler(ctx context.Context, input *Login_Req) (*UserAuthorized_Res, error) {
	result, err := Login_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				if api_error.Message == "" {
					api_error.Message = "Not Authorized"
				}
				return nil, huma.Error401Unauthorized(api_error.Message)
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("User is not found")
			case http.StatusBadRequest:
				if api_error.Message == "" {
					api_error.Message = "Bad request, try again later"
				}
				return nil, huma.Error400BadRequest(api_error.Message)
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /users/login")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &UserAuthorized_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil

}

func UpdateCurrentUser_Handler(ctx context.Context, input *UpdateCurrentUser_Req) (*schemas.Update_Res, error) {
	err := UpdateCurrentUser_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("User already exists")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("User is not found")
			case http.StatusBadRequest:
				if api_error.Message == "" {
					api_error.Message = "Bad request, try again later"
				}
				return nil, huma.Error400BadRequest(api_error.Message)
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/me")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func UpdateUserByID_Handler(ctx context.Context, input *UpdateUserByID_Req) (*schemas.Update_Res, error) {
	err := UpdateUserById_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("User already exists")
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("User is not found")
			case http.StatusBadRequest:
				if api_error.Message == "" {
					api_error.Message = "Bad request, try again later"
				}
				return nil, huma.Error400BadRequest(api_error.Message)
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func BanUserByID_Handler(ctx context.Context, input *BanUserByID_Req) (*schemas.Update_Res, error) {
	err := BanUserById_Service(ctx, input)
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
			logger.Error().Err(api_error).Msg("Error in PUT /users/{id}/ban")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteCurrentUser_Handler(ctx context.Context, input *DeleteCurrentUser_Req) (*schemas.Delete_Res, error) {
	err := DeleteCurrentUser_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("User is referenced in other tables")
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

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteUserByID_Handler(ctx context.Context, input *DeleteUserByID_Req) (*schemas.Delete_Res, error) {
	err := DeleteUserByID_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusUnauthorized:
				return nil, huma.Error401Unauthorized("Not Authorized")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("User is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /users/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
