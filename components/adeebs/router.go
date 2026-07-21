package adeebs

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

func GetAllAdeebs_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[OneAdeeb_Res], error) {

	list, err := gorm.G[database.Adeeb](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		}).
		Limit(input.Limit).
		Offset(input.Offset).
		Find(ctx)

	if err != nil {
		return nil, huma.Error404NotFound("Adeebs are not available")
	}

	adeebs := DBModels_To_ResModels(list)
	res := &schemas.GetAll_Res[OneAdeeb_Res]{
		Body: schemas.GetAll_Res_Body[OneAdeeb_Res]{
			Data:   adeebs,
			Limit:  input.Limit,
			Offset: input.Offset,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOneAdeeb_Handler(ctx context.Context, input *GetOneAdeeb_Req) (*GetOneAdeeb_Res, error) {

	adeeb_model, err := gorm.G[database.Adeeb](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("Adeeb's not found")
		} else {
			return nil, huma.Error400BadRequest("Bad Request getting Adeeb")
		}
	}

	adeeb_res := DBModel_To_ResModel(adeeb_model)
	res := &GetOneAdeeb_Res{
		Body:   adeeb_res,
		Status: http.StatusOK,
	}

	return res, nil
}

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

func CreateManyAdeeb_Handler(ctx context.Context, input *CreateManyAdeeb_Req) (*CreateManyAdeeb_Res, error) {

	var CreatedItems []OneAdeeb_Res
	var InvalidItems []InvalidItem

	new_data := ReqModels_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.Adeeb](
			database.Conn,
			clause.Returning{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "name"},
					{Name: "bio"},
					{Name: "time_period"},
					{Name: "reviewed"},
				},
			},
		).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				InvalidItems = append(InvalidItems, InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else {
				InvalidItems = append(InvalidItems, InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_adeeb := DBModel_To_ResModel(item)
		CreatedItems = append(CreatedItems, new_adeeb)

	}

	return &CreateManyAdeeb_Res{
		Body: CreateManyAdeeb_Res_Body{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}

func UpdateAdeeb_Handler(ctx context.Context, input *UpdateAdeeb_Req) (*schemas.Update_Res, error) {

	adeeb_model, err := gorm.G[database.Adeeb](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "bio"},
				{Name: "time_period"},
				{Name: "reviewed"},
			},
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("Adeeb's not found")
		} else {
			return nil, huma.Error400BadRequest("Bad Request getting Adeeb")
		}
	}

	if input.Body.Name != nil {
		adeeb_model.Name = *input.Body.Name
	}
	if input.Body.Bio != nil {
		adeeb_model.Bio = input.Body.Bio
	}

	if input.Body.TimePeriod != nil {
		adeeb_model.TimePeriod = *input.Body.TimePeriod
	}

	if input.Body.Reviewed != nil {
		adeeb_model.Reviewed = *input.Body.Reviewed
	}

	database.Conn.Save(&adeeb_model)

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
