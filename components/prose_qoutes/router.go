package prose_qoutes

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

func CreateOneProseQoute_Handler(ctx context.Context, input *CreateOneProseQoute_Req) (*CreateOneProseQoute_Res, error) {
	data := ReqModel_To_DBModel(input.Body)

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
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("ProseQoute already exists")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /prose_qoutes.")
		return nil, huma.Error400BadRequest("Bad Request creating ProseQoute.")
	}

	prose_qoute := DBModel_To_DescriptiveSchema(data)
	return &CreateOneProseQoute_Res{Body: prose_qoute, Status: http.StatusCreated}, nil
}

func CreateManyProseQoutes_Handler(ctx context.Context, input *CreateManyProseQoutes_Req) (*CreateManyProseQoutes_Res, error) {

	var CreatedItems []schemas.ProseQoute_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqModels_To_DBModels(input.Body)
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
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /prose_qoutes/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_prose_qoute := DBModel_To_DescriptiveSchema(item)
		CreatedItems = append(CreatedItems, new_prose_qoute)

	}

	return &CreateManyProseQoutes_Res{
		Body: schemas.CreateMany_Res_Body[schemas.ProseQoute_Descriptive]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}
