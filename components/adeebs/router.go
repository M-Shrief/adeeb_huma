package adeebs

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

func GetAllAdeebs_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.Adeeb_Descriptive], error) {
	var results []AdeebWithTotalCount

	err := database.Conn.Table("adeebs").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /adeebs")
		return nil, huma.Error400BadRequest("Uknown error getting adeebs")
	}

	adeebs, total_count := DistillDBModelsWithCount(results)
	res := &schemas.GetAll_Res[schemas.Adeeb_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.Adeeb_Descriptive]{
			Data:       adeebs,
			Limit:      input.Limit,
			Offset:     input.Offset,
			TotalCount: total_count,
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

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Adeeb's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /adeebs/{id}.")
		return nil, huma.Error400BadRequest("Bad Request getting Adeeb")
	}

	adeeb_res := DBModel_To_ResSchema(adeeb_model)
	res := &GetOneAdeeb_Res{
		Body:   adeeb_res,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneAdeeb_Handler(ctx context.Context, input *CreateOneAdeeb_Req) (*CreateOneAdeeb_Res, error) {
	data := ReqSchema_To_DBModel(input.Body)

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
		logger.Error().Err(err).Msg("Unknown errror in POST /adeebs.")
		return nil, huma.Error400BadRequest("Bad Request creating Adeeb.")
	}

	adeeb := DBModel_To_ResSchema(data)
	return &CreateOneAdeeb_Res{Body: adeeb, Status: http.StatusCreated}, nil

}

func CreateManyAdeeb_Handler(ctx context.Context, input *CreateManyAdeebs_Req) (*CreateManyAdeebs_Res, error) {

	var CreatedItems []schemas.Adeeb_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
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
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Already exists"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /adeebs/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_adeeb := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_adeeb)

	}

	return &CreateManyAdeebs_Res{
		Body: schemas.CreateMany_Res_Body[schemas.Adeeb_Descriptive]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}

func UpdateAdeeb_Handler(ctx context.Context, input *UpdateAdeeb_Req) (*schemas.Update_Res, error) {

	adeeb_model, err := gorm.G[database.Adeeb](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Adeeb's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /adeebs/{id}.")
		return nil, huma.Error400BadRequest("Bad Request updating Adeeb")
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

	err = database.Conn.Save(&adeeb_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, huma.Error400BadRequest("foreign key error")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /adeebs/{id}.")
		return nil, huma.Error400BadRequest("Bad Request updating adeeb")
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteAdeeb_Handler(ctx context.Context, input *DeleteAdeeb_Req) (*schemas.Delete_Res, error) {

	_, err := gorm.G[database.Adeeb](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in DELETE /adeebs/{id}.")
		return nil, huma.Error400BadRequest("Bad Request Deleting Adeeb")
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
