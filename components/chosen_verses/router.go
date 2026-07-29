package chosen_verses

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllChosenVerses_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.ChosenVerse_Descriptive], error) {

	list, err := gorm.G[database.ChosenVerse](
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
		Limit(input.Limit).
		Offset(input.Offset).
		Find(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /chosen_verses.")
		return nil, huma.Error404NotFound("Unknown error while getting chosen_verses")
	}

	chosen_verses := DBModels_To_ResSchemas(list)
	res := &schemas.GetAll_Res[schemas.ChosenVerse_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.ChosenVerse_Descriptive]{
			Data:   chosen_verses,
			Limit:  input.Limit,
			Offset: input.Offset,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOneChosenVerse_Handler(ctx context.Context, input *GetOneChosenVerse_Req) (*GetOneChosenVerse_Res, error) {

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
		return nil, huma.Error404NotFound("ChosenVerse's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /chosen_verses/{id}.")
		return nil, huma.Error400BadRequest("Bad Request getting ChosenVerse")
	}

	var chosen_verse_res GetOneChosenVerse_Res_Body
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

	res := &GetOneChosenVerse_Res{
		Body:   chosen_verse_res,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneChosenVerse_Handler(ctx context.Context, input *CreateOneChosenVerse_Req) (*CreateOneChosenVerse_Res, error) {
	data := ReqSchema_To_DBModel(input.Body)

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
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("ChosenVerse already exists")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /chosen_verses.")
		return nil, huma.Error400BadRequest("Bad Request creating ChosenVerse.")
	}

	chosen_verse := DBModel_To_ResSchema(data)
	return &CreateOneChosenVerse_Res{Body: chosen_verse, Status: http.StatusCreated}, nil
}

func CreateManyChosenVerses_Handler(ctx context.Context, input *CreateManyChosenVerses_Req) (*CreateManyChosenVerses_Res, error) {

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
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /chosen_verses/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_chosen_verse := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_chosen_verse)

	}

	return &CreateManyChosenVerses_Res{
		Body: schemas.CreateMany_Res_Body[schemas.ChosenVerse_Descriptive]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}

func UpdateChosenVerse_Handler(ctx context.Context, input *UpdateChosenVerse_Req) (*schemas.Update_Res, error) {

	chosen_verse_model, err := gorm.G[database.ChosenVerse](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("ChosenVerse's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /chosen_verses/{id}.")
		return nil, huma.Error400BadRequest("Bad Request update ChosenVerse")
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

	_ = database.Conn.Save(&chosen_verse_model)

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteChosenVerseHandler(ctx context.Context, input *DeleteChosenVerse_Req) (*schemas.Delete_Res, error) {

	_, err := gorm.G[database.ChosenVerse](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in DELETE /chosen_verses/{id}.")
		return nil, huma.Error400BadRequest("Bad Request Deleting ChosenVerse")
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
