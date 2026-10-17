package adeebs

import (
	"adeeb_huma/cache"
	"adeeb_huma/database"
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.Adeeb_Descriptive], error) {
	var results []AdeebWithTotalCount

	err := database.Conn.Table("adeebs").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /adeebs")
		return nil, api_error
	}

	adeebs, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.Adeeb_Descriptive]{
		Data:       adeebs,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOne_Repository(ctx context.Context, input *GetOneAdeeb_Req) (*schemas.Adeeb_Descriptive, error) {
	var adeeb_res schemas.Adeeb_Descriptive

	cache_key := cache.FormatKeyByID("adeebs", input.ID)
	adeeb_res, err := cache.GetJSON[schemas.Adeeb_Descriptive](ctx, cache_key, schemas.Adeeb_Descriptive{})
	if err == nil {
		return &adeeb_res, nil
	}

	adeeb_model, err := gorm.G[database.Adeeb](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /adeebs/{id}")
		return nil, api_error
	}

	adeeb_res = DBModel_To_ResSchema(adeeb_model)

	// Adding the result to the cache service
	err = cache.SetJSON(ctx, cache_key, adeeb_res)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.SetJSON() in GET /adeebs/{id}")
	}

	return &adeeb_res, nil
}

func CreateOne_Repository(ctx context.Context, input *CreateOneAdeeb_Req) (*schemas.Adeeb_Descriptive, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.Adeeb](
		database.Conn,
		// clause.OnConflict{DoNothing: true},
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		},
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Adeeb already exists", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /adeebs")
		return nil, api_error
	}

	adeeb := DBModel_To_ResSchema(data)
	return &adeeb, nil
}

func CreateMany_Repository(ctx context.Context, input *CreateManyAdeebs_Req) (*schemas.CreateMany_Res_Body[schemas.Adeeb_Descriptive], error) {
	var CreatedItems []schemas.Adeeb_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.Adeeb](
			database.Conn,
			clause.Returning{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "name"},
					{Name: "bio"},
					{Name: "time_period"},
					{Name: "reviewed"},
				},
			},
		).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /adeebs/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_adeeb := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_adeeb)

	}

	return &schemas.CreateMany_Res_Body[schemas.Adeeb_Descriptive]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func UpdateOne_Repository(ctx context.Context, input *UpdateAdeeb_Req) error {

	adeeb_model, err := gorm.G[database.Adeeb](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /adeebs/{id}")
		return api_error
	}

	if input.Body.Name != nil {
		adeeb_model.Name = *input.Body.Name
	}
	if input.Body.Bio != nil {
		adeeb_model.Bio = input.Body.Bio
	}

	if input.Body.TimePeriod != nil {
		adeeb_model.TimePeriod = *input.Body.TimePeriod
	}

	if input.Body.Reviewed != nil {
		adeeb_model.Reviewed = *input.Body.Reviewed
	}

	err = database.Conn.Save(&adeeb_model).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Adeeb already exists", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /adeebs/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("adeebs", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /adeebs/{id}")
	}

	return nil
}

func DeleteOne_Repository(ctx context.Context, input *DeleteAdeeb_Req) error {
	_, err := gorm.G[database.Adeeb](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /adeebs/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("adeebs", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /adeebs/{id}")
	}

	return nil
}
