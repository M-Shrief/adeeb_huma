package chosen_verses

import (
	"adeeb_huma/database"
	"adeeb_huma/logger"
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
