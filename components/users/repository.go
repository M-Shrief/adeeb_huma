package users

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.User_Descriptive], error) {
	var results []UserWithTotalCount

	err := database.Conn.Table("users").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /users")
		return nil, api_error
	}

	users, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.User_Descriptive]{
		Data:       users,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOne_Repository(ctx context.Context, id string) (*schemas.User_Descriptive, error) {
	user_model, err := gorm.G[database.User](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				{Name: "roles"},
			},
		}).
		Where("id = ?", id).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /users/{id} or /users/me")
		return nil, api_error
	}

	user_res := DBModel_To_ResSchema(user_model)

	return &user_res, nil
}

func Signup_Repository(ctx context.Context, input *Signup_Req) (*schemas.User_Descriptive, error) {
	new_hashed_password, err := auth.HashPassword(input.Body.Password)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, Message: "Couldn't process the password", CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /users/signup")
		return nil, api_error
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
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "User already exists", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /users/signup")
		return nil, api_error
	}

	user := DBModel_To_ResSchema(user_model)
	return &user, nil
}

func GetOneForLogin_Repository(ctx context.Context, input *Login_Req) (*database.User, error) {
	user_model, err := gorm.G[database.User](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "username"},
				{Name: "password"},
				{Name: "roles"},
			},
		}).
		Where("username = ?", input.Body.Username).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /users/login")
		return nil, api_error
	}

	return &user_model, nil
}

func UpdateCurrentUser_Repository(ctx context.Context, id string, input *UpdateCurrentUser_Req) error {

	user_model, err := gorm.G[database.User](database.Conn).
		Where("id = ?", id).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /users/me")
		return api_error
	}

	if input.Body.Username != nil {
		user_model.Username = *input.Body.Username
	}
	if input.Body.Password != nil {
		new_hashed_password, err := auth.HashPassword(*input.Body.Password)
		if err != nil {
			api_error := &errs.APIError{Status: http.StatusBadRequest, Message: "Couldn't process the password", CausedIn: "repository", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/me")
			return api_error
		}
		user_model.Password = new_hashed_password
	}

	err = database.Conn.Save(&user_model).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /users/me")
		return api_error
	}

	return nil
}

func UpdateUserById_Repository(ctx context.Context, input *UpdateUserByID_Req) error {
	user_model, err := gorm.G[database.User](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /users/{id}")
		return api_error
	}

	if input.Body.Username != nil {
		user_model.Username = *input.Body.Username
	}
	if input.Body.Password != nil {
		new_hashed_password, err := auth.HashPassword(*input.Body.Password)
		if err != nil {
			api_error := &errs.APIError{Status: http.StatusBadRequest, Message: "Couldn't process the password", CausedIn: "repository", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /users/{id}")
			return api_error
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

	err = database.Conn.Save(&user_model).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /users/{id}")
		return api_error
	}

	return nil
}

func BanUserByID_Repository(ctx context.Context, id string) error {

	// We use raw sql to make it in one query, it's fast and ensure uniqueness like we do in the server.
	_, err := database.Conn.Raw("UPDATE users set roles =  (select array_agg(distinct e) from unnest(array_append(users.roles, $2::role_enum)) e) WHERE id = $1", id, "Banned").Rows()
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /users/{id}/ban")
		return api_error
	}

	return nil
}

func DeleteOne_Repository(ctx context.Context, id string) error {
	_, err := gorm.G[database.User](database.Conn).
		Where("id = ?", id).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /users/me or /users/{id}")
		return api_error
	}

	return nil
}
