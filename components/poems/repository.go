package poems

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

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.Poem_Descriptive], error) {
	var results []PoemWithTotalCount

	err := database.Conn.Table("poems").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /poems")
		return nil, api_error
	}

	poems, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.Poem_Descriptive]{
		Data:       poems,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOne_Repository(ctx context.Context, input *GetOnePoem_Req) (*GetOnePoem_Res_Body, error) {
	cache_key := cache.FormatKeyByID("poems", input.ID)
	poem_res, err := cache.GetJSON[GetOnePoem_Res_Body](ctx, cache_key, GetOnePoem_Res_Body{})
	if err == nil {
		return &poem_res, nil
	}

	poem_model, err := gorm.G[database.Poem](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "intro"},
				{Name: "verses"},
				{Name: "is_couplet"},
				{Name: "reviewed"},
				{Name: "adeeb_id"},
			},
		}).
		Preload("Adeeb", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name")
			return nil
		}).
		Preload("ChosenVerses", func(db gorm.PreloadBuilder) error {
			db.Select("id", "verses", "is_couplet", "poem_id")
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
		logger.Error().Err(api_error).Msg("Error in GET /poems/{id}")
		return nil, api_error
	}

	poem_res.ID = poem_model.ID
	poem_res.Intro = poem_model.Intro
	poem_res.Verses = poem_model.Verses
	poem_res.IsCouplet = poem_model.IsCouplet
	poem_res.Reviewed = poem_model.Reviewed

	poem_res.AdeebID = poem_model.AdeebID
	poem_res.Adeeb.ID = poem_model.Adeeb.ID
	poem_res.Adeeb.Name = poem_model.Adeeb.Name

	for _, model := range poem_model.ChosenVerses {
		var item schemas.ChosenVerse_Minimal
		item.ID = model.ID
		item.Verses = model.Verses
		item.IsCouplet = model.IsCouplet

		poem_res.ChosenVerses = append(poem_res.ChosenVerses, item)
	}

	// Adding the result to the cache service
	err = cache.SetJSON(ctx, cache_key, poem_res)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.SetJSON() in GET /poems/{id}")
	}

	return &poem_res, nil
}

func CreateOne_Repository(ctx context.Context, input *CreateOnePoem_Req) (*schemas.Poem_Descriptive, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.Poem](
		database.Conn,
		// clause.OnConflict{DoNothing: true},
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "intro"},
				{Name: "verses"},
				{Name: "is_couplet"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
			},
		},
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Poem already exists", CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /poems")
		return nil, api_error
	}

	poem := DBModel_To_ResSchema(data)
	return &poem, nil
}

func CreateMany_Repository(ctx context.Context, input *CreateManyPoems_Req) (*schemas.CreateMany_Res_Body[schemas.Poem_Descriptive], error) {
	var CreatedItems []schemas.Poem_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.Poem](
			database.Conn,
			clause.Returning{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "intro"},
					{Name: "verses"},
					{Name: "is_couplet"},
					{Name: "reviewed"},

					{Name: "adeeb_id"},
				},
			},
		).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Foreign Key Error"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /poems/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_poem := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_poem)

	}

	return &schemas.CreateMany_Res_Body[schemas.Poem_Descriptive]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func UpdateOne_Repository(ctx context.Context, input *UpdatePoem_Req) error {

	poem_model, err := gorm.G[database.Poem](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /poems/{id}")
		return api_error
	}

	if input.Body.Intro != nil {
		poem_model.Intro = *input.Body.Intro
	}
	if input.Body.Verses != nil {
		poem_model.Verses = *input.Body.Verses
	}
	if input.Body.IsCouplet != nil {
		poem_model.IsCouplet = *input.Body.IsCouplet
	}
	if input.Body.Reviewed != nil {
		poem_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.AdeebID != nil {
		poem_model.AdeebID = *input.Body.AdeebID
	}

	err = database.Conn.Save(&poem_model).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Poem already exists", CausedIn: "repository", Err: err}
		return api_error
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /poems/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("poems", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /poems/{id}")
	}

	return nil
}

func DeleteOne_Repository(ctx context.Context, input *DeletePoem_Req) error {
	_, err := gorm.G[database.Poem](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /poems/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("poems", input.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /poems/{id}")
	}

	return nil
}
