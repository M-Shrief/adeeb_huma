package users

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllUsers_Handler(ctx context.Context, input *GetAllUsers_Req) (*schemas.GetAll_Res[schemas.User_Descriptive], error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_Analytics, auth.OPEnum_Read),
	}
	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	list, err := gorm.G[database.User](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				{Name: "roles"},
			},
		}).
		Limit(input.Limit).
		Offset(input.Offset).
		Find(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /users.")
		return nil, huma.Error404NotFound("Unknown error while getting users")
	}

	users := DBModels_To_DescriptiveSchemas(list)
	res := &schemas.GetAll_Res[schemas.User_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.User_Descriptive]{
			Data:   users,
			Limit:  input.Limit,
			Offset: input.Offset,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetCurrentUser_Handler(ctx context.Context, input *GetCurrentUser_Req) (*GetOneUser_Res, error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)
	user_model, err := gorm.G[database.User](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				{Name: "roles"},
			},
		}).
		Where("id = ?", user_id).
		First(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /users/me.")
		return nil, huma.Error404NotFound("Unknown error while getting current user")
	}

	user := DBModel_To_DescriptiveSchema(user_model)
	res := &GetOneUser_Res{
		Body:   user,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetUserByID_Handler(ctx context.Context, input *GetUserByID_Req) (*GetOneUser_Res, error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_Analytics, auth.OPEnum_Read),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_model, err := gorm.G[database.User](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				{Name: "roles"},
			},
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /users/{id}.")
		return nil, huma.Error404NotFound("Unknown error while users by ID")
	}

	user := DBModel_To_DescriptiveSchema(user_model)
	res := &GetOneUser_Res{
		Body:   user,
		Status: http.StatusOK,
	}

	return res, nil
}

func Signup_Handler(ctx context.Context, input *Signup_Req) (*UserAuthorized_Res, error) {
	new_hashed_password, err := auth.HashPassword(input.Body.Password)
	if err != nil {
		logger.Error().Err(err).Msg("Error hashing passsord in POST /users/signup")
		return nil, huma.Error400BadRequest("Couldn't process the password")
	}

	roles := utils.EnsureSliceItemsAreUnique(input.Body.Roles)
	if slices.Contains(roles, database.RoleEnum_Normal) == false {
		roles = append(roles, database.RoleEnum_Normal)
	}
	user_model := database.User{
		Username: input.Body.Username,
		Password: new_hashed_password,
		Roles:    roles,
	}

	err = gorm.G[database.User](
		database.Conn,
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				// {Name: "password"},
				{Name: "roles"},
			},
		},
	).Create(ctx, &user_model)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("User already exists")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error creating User in POST /users/signup")
		return nil, huma.Error400BadRequest("Bad Request signing up user.")
	}

	token, err := auth.CreateJWT(
		time.Hour*2,
		auth.JWTUserClaim{
			ID:       user_model.ID.String(),
			Username: user_model.Username,
			Roles:    user_model.Roles,
		},
		auth.CreatePermissions(user_model.Roles),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error creating JWT token in POST /users/signup")
		return nil, huma.Error500InternalServerError("Unknown Error, try again later.")
	}

	return &UserAuthorized_Res{
		Body:   UserAuthorized_Res_Body{DBModel_To_DescriptiveSchema(user_model), token},
		Status: http.StatusCreated,
	}, nil
}

func Login_Handler(ctx context.Context, input *Login_Req) (*UserAuthorized_Res, error) {
	user_model, err := gorm.G[database.User](database.Conn).
		Where("username = ?", input.Body.Username).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("User's doesn't exist")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /poems/{id}.")
		return nil, huma.Error400BadRequest("Bad Request trying to login, try again later")
	}

	err = auth.VerifyPassword(input.Body.Password, user_model.Password)
	if err != nil {
		return nil, huma.Error401Unauthorized("Password is incorrect")
	}

	token, err := auth.CreateJWT(
		time.Hour*2,
		auth.JWTUserClaim{
			ID:       user_model.ID.String(),
			Username: user_model.Username,
			Roles:    user_model.Roles,
		},
		auth.CreatePermissions(user_model.Roles),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error creating JWT token in POST /users/signup")
		return nil, huma.Error500InternalServerError("Unknown Error, try again later.")
	}

	return &UserAuthorized_Res{
		Body:   UserAuthorized_Res_Body{DBModel_To_DescriptiveSchema(user_model), token},
		Status: http.StatusCreated,
	}, nil

}

func UpdateCurrentUser_Handler(ctx context.Context, input *UpdateCurrentUser_Req) (*schemas.Update_Res, error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Write),
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Write)
	if is_authorized == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)
	user_model, err := gorm.G[database.User](database.Conn).
		Where("id = ?", user_id).
		First(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /users/me.")
		return nil, huma.Error404NotFound("Unknown error while updating current user")
	}

	if input.Body.Username != nil {
		user_model.Username = *input.Body.Username
	}
	if input.Body.Password != nil {
		new_hashed_password, err := auth.HashPassword(*input.Body.Password)
		if err != nil {
			logger.Error().Err(err).Msg("Error hashing passsord in POST /users/signup")
			return nil, huma.Error400BadRequest("Couldn't process the password")
		}
		user_model.Password = new_hashed_password
	}
	if input.Body.Roles != nil {
		roles := utils.EnsureSliceItemsAreUnique(*input.Body.Roles)
		if slices.Contains(roles, database.RoleEnum_Normal) == false {
			roles = append(roles, database.RoleEnum_Normal)
		}
		user_model.Roles = roles
	}

	database.Conn.Save(&user_model)

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
