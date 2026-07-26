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

func GetAllProseQoutes_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.ProseQoute_Descriptive], error) {

	list, err := gorm.G[database.ProseQoute](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "qoute"},
				{Name: "source"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
			},
		}).
		Limit(input.Limit).
		Offset(input.Offset).
		Find(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /prose_qoutes.")
		return nil, huma.Error404NotFound("Unknown error while getting prose_qoutes")
	}

	prose_qoutes := DBModels_To_DescriptiveSchemas(list)
	res := &schemas.GetAll_Res[schemas.ProseQoute_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.ProseQoute_Descriptive]{
			Data:   prose_qoutes,
			Limit:  input.Limit,
			Offset: input.Offset,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOneProseQoute_Handler(ctx context.Context, input *GetOneProseQoute_Req) (*GetOneProseQoute_Res, error) {

	prose_qoute_model, err := gorm.G[database.ProseQoute](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "qoute"},
				{Name: "source"},
				{Name: "tags"},
				{Name: "reviewed"},

				{Name: "adeeb_id"},
			},
		}).
		Preload("Adeeb", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("ProseQoute's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /prose_qoutes/{id}.")
		return nil, huma.Error400BadRequest("Bad Request getting ProseQoute")
	}

	var prose_qoute_res GetOneProseQoute_Res_Body
	prose_qoute_res.ID = prose_qoute_model.ID
	prose_qoute_res.Qoute = prose_qoute_model.Qoute
	prose_qoute_res.Source = prose_qoute_model.Source
	prose_qoute_res.Tags = prose_qoute_model.Tags
	prose_qoute_res.Reviewed = prose_qoute_model.Reviewed

	prose_qoute_res.AdeebID = prose_qoute_model.AdeebID
	prose_qoute_res.Adeeb.ID = prose_qoute_model.Adeeb.ID
	prose_qoute_res.Adeeb.Name = prose_qoute_model.Adeeb.Name

	res := &GetOneProseQoute_Res{
		Body:   prose_qoute_res,
		Status: http.StatusOK,
	}

	return res, nil
}

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

func UpdateProseQoute_Handler(ctx context.Context, input *UpdateProseQoute_Req) (*schemas.Update_Res, error) {

	prose_qoute_model, err := gorm.G[database.ProseQoute](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("ProseQoute's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /prose_qoutes/{id}.")
		return nil, huma.Error400BadRequest("Bad Request update ProseQoute")
	}

	if input.Body.Qoute != nil {
		prose_qoute_model.Qoute = *input.Body.Qoute
	}
	if input.Body.Source != nil {
		prose_qoute_model.Source = input.Body.Source
	}
	if input.Body.Tags != nil {
		prose_qoute_model.Tags = *input.Body.Tags
	}
	if input.Body.Reviewed != nil {
		prose_qoute_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.AdeebID != nil {
		prose_qoute_model.AdeebID = *input.Body.AdeebID
	}

	database.Conn.Save(&prose_qoute_model)

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
