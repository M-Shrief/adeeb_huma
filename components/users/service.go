package users

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"net/http"
	"time"
)

func GetAll_Service(input *GetAllUsers_Req) (*schemas.GetAll_Res_Body[schemas.User_Descriptive], error) {
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
			logger.Error().Err(api_error).Msg("Error in GET /users")
			return nil, api_error
		}
	}
	return result, nil
}

func GetCurrentUser_Service(ctx context.Context, input *GetCurrentUser_Req) (*schemas.User_Descriptive, error) {
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

	result, err := GetOne_Repository(ctx, user_id)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /users/me")
			return nil, api_error
		}
	}
	return result, nil
}

func GetOneById_Service(ctx context.Context, input *GetUserByID_Req) (*schemas.User_Descriptive, error) {
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

	result, err := GetOne_Repository(ctx, input.ID.String())
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /users/{id}")
			return nil, api_error
		}
	}
	return result, nil
}

func Signup_Service(ctx context.Context, input *Signup_Req) (*UserAuthorized_Res_Body, error) {
	user, err := Signup_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /users/signup")
			return nil, api_error
		}
	}
	token, err := auth.CreateJWT(
		time.Hour*2,
		auth.JWTUserClaim{
			ID:       user.ID.String(),
			Username: user.Username,
			Roles:    user.Roles,
		},
		auth.CreatePermissions(user.Roles),
	)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, Message: "Unknown error creating JWT token in POST /users/signup", CausedIn: "service", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /users/signup")
		return nil, api_error
	}

	return &UserAuthorized_Res_Body{*user, token}, nil
}

func Login_Service(ctx context.Context, input *Login_Req) (*UserAuthorized_Res_Body, error) {
	user, err := GetOneForLogin_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /users/login")
			return nil, api_error
		}
	}

	err = auth.VerifyPassword(input.Body.Password, user.Password)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, Message: "Password is incorrect", CausedIn: "service", Err: err}
		return nil, api_error
	}

	token, err := auth.CreateJWT(
		time.Hour*2,
		auth.JWTUserClaim{
			ID:       user.ID.String(),
			Username: user.Username,
			Roles:    user.Roles,
		},
		auth.CreatePermissions(user.Roles),
	)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, Message: "Unknown error creating JWT token in POST /users/login", CausedIn: "service", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /users/login")
		return nil, api_error
	}

	return &UserAuthorized_Res_Body{DBModel_To_ResSchema(*user), token}, nil
}

func UpdateCurrentUser_Service(ctx context.Context, input *UpdateCurrentUser_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Write)
	if is_authorized == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)

	err = UpdateCurrentUser_Repository(ctx, user_id, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/me")
			return api_error
		}
	}
	return nil
}

func UpdateUserById_Service(ctx context.Context, input *UpdateUserByID_Req) error {
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

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	err = UpdateUserById_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/{id}")
			return api_error
		}
	}
	return nil
}

func BanUserById_Service(ctx context.Context, input *BanUserByID_Req) error {
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

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	err = BanUserByID_Repository(ctx, input.ID.String())
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/{id}/ban")
			return api_error
		}
	}
	return nil
}

func DeleteCurrentUser_Service(ctx context.Context, input *DeleteCurrentUser_Req) error {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Write)
	if is_authorized == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)

	err = DeleteOne_Repository(ctx, user_id)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /users/me")
			return api_error
		}
	}
	return nil
}

func DeleteUserByID_Service(ctx context.Context, input *DeleteUserByID_Req) error {
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

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Write)
	if is_adminstrator == false {
		api_error := &errs.APIError{Status: http.StatusUnauthorized, CausedIn: "service", Err: err}
		return api_error
	}

	err = DeleteOne_Repository(ctx, input.ID.String())
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /users/{id}")
			return api_error
		}
	}
	return nil
}
