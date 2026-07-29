package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/logger"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func CreateOneOrder_Handler(ctx context.Context, input *CreateOneOrder_Req) (*CreateOneOrder_Res, error) {
	data := ReqModel_To_DBModel(input.Body)

	err := gorm.G[database.Order](database.Conn).
		Create(ctx, &data)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders.")
		return nil, huma.Error400BadRequest("Bad Request creating Order.")
	}

	new_order := DBModel_To_DescriptiveSchema(data)
	return &CreateOneOrder_Res{Body: new_order, Status: http.StatusCreated}, nil
}
