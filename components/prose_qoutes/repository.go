package prose_qoutes

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

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.ProseQoute_Descriptive], error) {
	var results []ProseQouteWithTotalCount

	err := database.Conn.Table("prose_qoutes").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes")
		return nil, api_error
	}

	prose_qoutes, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.ProseQoute_Descriptive]{
		Data:       prose_qoutes,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOne_Repository(ctx context.Context, input *GetOneProseQoute_Req) (*GetOneProseQoute_Res_Body, error) {
	var prose_qoute_res GetOneProseQoute_Res_Body

	cache_key := cache.FormatKeyByID("prose_qoutes", input.ID)
	prose_qoute_res, err := cache.GetJSON[GetOneProseQoute_Res_Body](ctx, cache_key, GetOneProseQoute_Res_Body{})
	if err == nil {
		return &prose_qoute_res, nil
	}

	prose_qoute_model, err := gorm.G[database.ProseQoute](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "qoute"},
				{Name: "source"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
			},
		}).
		Preload("Adeeb", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes/{id}")
		return nil, api_error
	}

	prose_qoute_res.ID = prose_qoute_model.ID
	prose_qoute_res.Qoute = prose_qoute_model.Qoute
	prose_qoute_res.Source = prose_qoute_model.Source
	prose_qoute_res.Tags = prose_qoute_model.Tags
	prose_qoute_res.Reviewed = prose_qoute_model.Reviewed

	prose_qoute_res.AdeebID = prose_qoute_model.AdeebID
	prose_qoute_res.Adeeb.ID = prose_qoute_model.Adeeb.ID
	prose_qoute_res.Adeeb.Name = prose_qoute_model.Adeeb.Name

	// Adding the result to the cache service
	err = cache.SetJSON(ctx, cache_key, prose_qoute_res)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.SetJSON() in GET /prose_qoutes/{id}")
	}

	return &prose_qoute_res, nil
}

func CreateOne_Repository(ctx context.Context, input *CreateOneProseQoute_Req) (*schemas.ProseQoute_Descriptive, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.ProseQoute](
		database.Conn,
		// clause.OnConflict{DoNothing: true},
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "qoute"},
				{Name: "source"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
			},
		},
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "ProseQoute already exists", CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /prose_qoutes")
		return nil, api_error
	}

	prose_qoute := DBModel_To_ResSchema(data)
	return &prose_qoute, nil
}

func CreateMany_Repository(ctx context.Context, input *CreateManyProseQoutes_Req) (*schemas.CreateMany_Res_Body[schemas.ProseQoute_Descriptive], error) {
	var CreatedItems []schemas.ProseQoute_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.ProseQoute](
			database.Conn,
			clause.Returning{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "qoute"},
					{Name: "source"},
					{Name: "tags"},
					{Name: "reviewed"},

					{Name: "adeeb_id"},
				},
			},
		).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Foreign Key Error"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /prose_qoutes/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_prose_qoute := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_prose_qoute)

	}

	return &schemas.CreateMany_Res_Body[schemas.ProseQoute_Descriptive]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func UpdateOne_Repository(ctx context.Context, input *UpdateProseQoute_Req) error {

	prose_qoute_model, err := gorm.G[database.ProseQoute](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /prose_qoutes/{id}")
		return api_error
	}

	if input.Body.Qoute != nil {
		prose_qoute_model.Qoute = *input.Body.Qoute
	}
	if input.Body.Source != nil {
		prose_qoute_model.Source = input.Body.Source
	}
	if input.Body.Tags != nil {
		prose_qoute_model.Tags = *input.Body.Tags
	}
	if input.Body.Reviewed != nil {
		prose_qoute_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.AdeebID != nil {
		prose_qoute_model.AdeebID = *input.Body.AdeebID
	}

	err = database.Conn.Save(&prose_qoute_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /prose_qoutes/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("prose_qoutes", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /prose_qoutes/{id}")
	}

	return nil
}

func DeleteOne_Repository(ctx context.Context, input *DeleteProseQoute_Req) error {
	_, err := gorm.G[database.ProseQoute](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /prose_qoutes/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("prose_qoutes", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /prose_qoutes/{id}")
	}

	return nil
}
