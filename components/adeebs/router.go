package adeebs

import (
	"adeeb_huma/database"
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateOneAdeeb_Handler(ctx context.Context, input *CreateOneAdeeb_Req) (*CreateOneAdeeb_Res, error) {
	data := database.Adeeb{
		Name:       input.Body.Name,
		TimePeriod: input.Body.TimePeriod,
		Bio:        &input.Body.Bio,
		Reviewed:   true,
	}

	err := gorm.G[database.Adeeb](
		database.Conn,
		// clause.OnConflict{DoNothing: true},
		clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		},
	).Create(ctx, &data)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("Adeeb already exists")
	}

	if err != nil {
		return nil, huma.Error400BadRequest("Bad Request creating Adeeb.")
	}

	adeeb := DBModel_To_ResModel(data)
	return &CreateOneAdeeb_Res{Body: adeeb, Status: http.StatusCreated}, nil

}
