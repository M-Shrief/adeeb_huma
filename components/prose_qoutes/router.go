package prose_qoutes

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
