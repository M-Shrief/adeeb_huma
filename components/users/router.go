package users

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/logger"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Signup_Handler(ctx context.Context, input *Signup_Req) (*UserAuthorized_Res, error) {
	new_hashed_password, err := auth.HashPassword(input.Body.Password)
	if err != nil {
		logger.Error().Err(err).Msg("Error hashing passsord in POST /users/signup")
		return nil, huma.Error400BadRequest("Couldn't process the password")
	}

	user_model := database.User{
		Username:  input.Body.Username,
		Passsword: new_hashed_password,
		Roles:     input.Body.Roles,
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

	err = auth.VerifyPassword(input.Body.Password, user_model.Passsword)
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
