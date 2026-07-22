package poems

import (
	"adeeb_huma/database"
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateOnePoem_Handler(ctx context.Context, input *CreateOnePoem_Req) (*CreateOnePoem_Res, error) {
	data := database.Poem{
		Intro:     input.Body.Intro,
		Verses:    input.Body.Verses,
		IsCouplet: input.Body.IsCouplet,
		Reviewed:  input.Body.Reviewed,
		AdeebID:   input.Body.AdeebID,
	}

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
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("Poem already exists")
	}
	if err != nil {
		return nil, huma.Error400BadRequest("Bad Request creating Poem.")
	}

	poem := DBModel_To_ResModel(data)
	return &CreateOnePoem_Res{Body: poem, Status: http.StatusCreated}, nil
}
