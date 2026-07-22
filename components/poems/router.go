package poems

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
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

func CreateManyPoems_Handler(ctx context.Context, input *CreateManyPoems_Req) (*CreateManyPoems_Res, error) {

	var CreatedItems []OnePoem_Res
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqModels_To_DBModels(input.Body)
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
			} else {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_poem := DBModel_To_ResModel(item)
		CreatedItems = append(CreatedItems, new_poem)

	}

	return &CreateManyPoems_Res{
		Body: schemas.CreateMany_Res_Body[OnePoem_Res]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}
