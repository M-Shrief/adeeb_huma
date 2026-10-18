package chosen_verses

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

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.ChosenVerse_Descriptive], error) {
	var results []ChosenVerseWithTotalCount

	err := database.Conn.Table("chosen_verses").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /chosen_verses")
		return nil, api_error
	}

	chosen_verses, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.ChosenVerse_Descriptive]{
		Data:       chosen_verses,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOne_Repository(ctx context.Context, input *GetOneChosenVerse_Req) (*GetOneChosenVerse_Res_Body, error) {
	var chosen_verse_res GetOneChosenVerse_Res_Body

	cache_key := cache.FormatKeyByID("chosen_verses", input.ID)
	chosen_verse_res, err := cache.GetJSON[GetOneChosenVerse_Res_Body](ctx, cache_key, GetOneChosenVerse_Res_Body{})
	if err == nil {
		return &chosen_verse_res, nil
	}

	chosen_verse_model, err := gorm.G[database.ChosenVerse](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "verses"},
				{Name: "is_couplet"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
				{Name: "poem_id"},
			},
		}).
		Preload("Adeeb", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name")
			return nil
		}).
		Preload("Poem", func(db gorm.PreloadBuilder) error {
			db.Select("id", "intro")
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
		logger.Error().Err(api_error).Msg("Error in GET /chosen_verses/{id}")
		return nil, api_error
	}

	chosen_verse_res.ID = chosen_verse_model.ID
	chosen_verse_res.Verses = chosen_verse_model.Verses
	chosen_verse_res.IsCouplet = chosen_verse_model.IsCouplet
	chosen_verse_res.Tags = chosen_verse_model.Tags
	chosen_verse_res.Reviewed = chosen_verse_model.Reviewed

	chosen_verse_res.AdeebID = chosen_verse_model.AdeebID
	chosen_verse_res.Adeeb.ID = chosen_verse_model.Adeeb.ID
	chosen_verse_res.Adeeb.Name = chosen_verse_model.Adeeb.Name

	chosen_verse_res.PoemID = chosen_verse_model.PoemID
	chosen_verse_res.Poem.ID = chosen_verse_model.Poem.ID
	chosen_verse_res.Poem.Intro = chosen_verse_model.Poem.Intro

	// Adding the result to the cache service
	err = cache.SetJSON(ctx, cache_key, chosen_verse_res)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.SetJSON() in GET /chosen_verses/{id}")
	}

	return &chosen_verse_res, nil
}

func CreateOne_Repository(ctx context.Context, input *CreateOneChosenVerse_Req) (*schemas.ChosenVerse_Descriptive, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.ChosenVerse](
		database.Conn,
		// clause.OnConflict{DoNothing: true},
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "verses"},
				{Name: "is_couplet"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
				{Name: "poem_id"},
			},
		},
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /chosen_verses")
		return nil, api_error
	}

	chosen_verse := DBModel_To_ResSchema(data)
	return &chosen_verse, nil
}

func CreateMany_Repository(ctx context.Context, input *CreateManyChosenVerses_Req) (*schemas.CreateMany_Res_Body[schemas.ChosenVerse_Descriptive], error) {
	var CreatedItems []schemas.ChosenVerse_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.ChosenVerse](
			database.Conn,
			clause.Returning{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "verses"},
					{Name: "is_couplet"},
					{Name: "tags"},
					{Name: "reviewed"},

					{Name: "adeeb_id"},
					{Name: "poem_id"},
				},
			},
		).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Foreign Key Error"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /chosen_verses/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_chosen_verse := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_chosen_verse)

	}

	return &schemas.CreateMany_Res_Body[schemas.ChosenVerse_Descriptive]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func UpdateOne_Repository(ctx context.Context, input *UpdateChosenVerse_Req) error {

	chosen_verse_model, err := gorm.G[database.ChosenVerse](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /chosen_verses/{id}")
		return api_error
	}

	if input.Body.Verses != nil {
		chosen_verse_model.Verses = *input.Body.Verses
	}
	if input.Body.IsCouplet != nil {
		chosen_verse_model.IsCouplet = *input.Body.IsCouplet
	}
	if input.Body.Tags != nil {
		chosen_verse_model.Tags = *input.Body.Tags
	}
	if input.Body.Reviewed != nil {
		chosen_verse_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.AdeebID != nil {
		chosen_verse_model.AdeebID = *input.Body.AdeebID
	}
	if input.Body.PoemID != nil {
		chosen_verse_model.PoemID = *input.Body.PoemID
	}

	err = database.Conn.Save(&chosen_verse_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /chosen_verses/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("chosen_verses", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /chosen_verses/{id}")
	}

	return nil
}

func DeleteOne_Repository(ctx context.Context, input *DeleteChosenVerse_Req) error {
	_, err := gorm.G[database.ChosenVerse](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /chosen_verses/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("chosen_verses", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /chosen_verses/{id}")
	}

	return nil
}
