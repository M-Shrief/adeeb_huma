package poems

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

func GetAllPoems_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.Poem_Descriptive], error) {

	var results []PoemWithTotalCount
	err := database.Conn.Table("poems").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /poems.")
		return nil, huma.Error404NotFound("Unknown error while getting poems")
	}

	poems, total_count := DistillDBModelsWithCount(results)
	res := &schemas.GetAll_Res[schemas.Poem_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.Poem_Descriptive]{
			Data:       poems,
			Limit:      input.Limit,
			Offset:     input.Offset,
			TotalCount: total_count,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOnePoem_Handler(ctx context.Context, input *GetOnePoem_Req) (*GetOnePoem_Res, error) {

	poem_model, err := gorm.G[database.Poem](
		database.Conn,
		clause.Select{
			Columns: []clause.Column{
				{Name: "id"},
				{Name: "intro"},
				{Name: "verses"},
				{Name: "is_couplet"},
				{Name: "reviewed"},
				{Name: "adeeb_id"},
			},
		}).
		Preload("Adeeb", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name")
			return nil
		}).
		Preload("ChosenVerses", func(db gorm.PreloadBuilder) error {
			db.Select("id", "verses", "is_couplet", "poem_id")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Poem's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /poems/{id}.")
		return nil, huma.Error400BadRequest("Bad Request getting Poem")
	}

	var poem_res GetOnePoem_Res_Body
	poem_res.ID = poem_model.ID
	poem_res.Intro = poem_model.Intro
	poem_res.Verses = poem_model.Verses
	poem_res.IsCouplet = poem_model.IsCouplet
	poem_res.Reviewed = poem_model.Reviewed

	poem_res.AdeebID = poem_model.AdeebID
	poem_res.Adeeb.ID = poem_model.Adeeb.ID
	poem_res.Adeeb.Name = poem_model.Adeeb.Name

	for _, model := range poem_model.ChosenVerses {
		var item schemas.ChosenVerse_Minimal
		item.ID = model.ID
		item.Verses = model.Verses
		item.IsCouplet = model.IsCouplet

		poem_res.ChosenVerses = append(poem_res.ChosenVerses, item)
	}

	res := &GetOnePoem_Res{
		Body:   poem_res,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOnePoem_Handler(ctx context.Context, input *CreateOnePoem_Req) (*CreateOnePoem_Res, error) {
	data := ReqSchema_To_DBModel(input.Body)

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
		logger.Error().Err(err).Msg("Unknown errror in POST /poems.")
		return nil, huma.Error400BadRequest("Bad Request creating Poem.")
	}

	poem := DBModel_To_ResSchema(data)
	return &CreateOnePoem_Res{Body: poem, Status: http.StatusCreated}, nil
}

func CreateManyPoems_Handler(ctx context.Context, input *CreateManyPoems_Req) (*CreateManyPoems_Res, error) {

	var CreatedItems []schemas.Poem_Descriptive
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
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
				logger.Error().Err(err).Msg("Unknown errror in POST /poems/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_poem := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_poem)

	}

	return &CreateManyPoems_Res{
		Body: schemas.CreateMany_Res_Body[schemas.Poem_Descriptive]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil

}

func UpdatePoem_Handler(ctx context.Context, input *UpdatePoem_Req) (*schemas.Update_Res, error) {

	poem_model, err := gorm.G[database.Poem](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Poem's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /poems/{id}.")
		return nil, huma.Error400BadRequest("Bad Request update Poem")
	}

	if input.Body.Intro != nil {
		poem_model.Intro = *input.Body.Intro
	}
	if input.Body.Verses != nil {
		poem_model.Verses = *input.Body.Verses
	}
	if input.Body.IsCouplet != nil {
		poem_model.IsCouplet = *input.Body.IsCouplet
	}
	if input.Body.Reviewed != nil {
		poem_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.AdeebID != nil {
		poem_model.AdeebID = *input.Body.AdeebID
	}

	err = database.Conn.Save(&poem_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, huma.Error400BadRequest("foreign key error")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in PUT /poems/{id}.")
		return nil, huma.Error400BadRequest("Bad Request updating poem")
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeletePoemHandler(ctx context.Context, input *DeletePoem_Req) (*schemas.Delete_Res, error) {

	_, err := gorm.G[database.Poem](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in DELETE /poems/{id}.")
		return nil, huma.Error400BadRequest("Bad Request Deleting Poem")
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
