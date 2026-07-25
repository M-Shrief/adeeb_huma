package chosen_verses

import (
	"adeeb_huma/database"
	"adeeb_huma/logger"
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

	chosen_verses := DBModels_To_ResModels(list)
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

func CreateOneChosenVerse_Handler(ctx context.Context, input *CreateOneChosenVerse_Req) (*CreateOneChosenVerse_Res, error) {
	data := database.ChosenVerse{
		Verses:    input.Body.Verses,
		IsCouplet: input.Body.IsCouplet,
		Tags:      input.Body.Tags,
		Reviewed:  input.Body.Reviewed,
		AdeebID:   input.Body.AdeebID,
	}

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

	chosen_verse := DBModel_To_ResModel(data)
	return &CreateOneChosenVerse_Res{Body: chosen_verse, Status: http.StatusCreated}, nil
}

func CreateManyChosenVerses_Handler(ctx context.Context, input *CreateManyChosenVerses_Req) (*CreateManyChosenVerses_Res, error) {

	var CreatedItems []schemas.ChosenVerse_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqModels_To_DBModels(input.Body)
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

		new_chosen_verse := DBModel_To_ResModel(item)
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
