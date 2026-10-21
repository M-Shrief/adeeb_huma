package orders

import (
	"adeeb_huma/cache"
	"adeeb_huma/database"
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAll_Repository(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.Order_Descriptive], error) {
	var results []OrderWithTotalCount

	err := database.Conn.Table("orders").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /orders")
		return nil, api_error
	}

	orders, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.Order_Descriptive]{
		Data:       orders,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetUserOrders_Repository(user_id string, input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.Order_Descriptive], error) {
	var results []OrderWithTotalCount

	err := database.Conn.Table("orders").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Where("user_id = ?", user_id).
		Find(&results).Error

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /orders/me")
		return nil, api_error
	}

	orders, total_count := DistillDBModelsWithCount(results)

	return &schemas.GetAll_Res_Body[schemas.Order_Descriptive]{
		Data:       orders,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalCount: total_count,
	}, nil
}

func GetOrderByID_Repository(ctx context.Context, id uuid.UUID) (*OneOrder_Res, error) {
	var order_res OneOrder_Res

	cache_key := cache.FormatKeyByID("orders", id)
	order_res, err := cache.GetJSON[OneOrder_Res](ctx, cache_key, OneOrder_Res{})
	if err == nil {
		return &order_res, nil
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Preload("Prints", func(db gorm.PreloadBuilder) error {
			db.Select("*")
			return nil
		}).
		Where("id = ?", id).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /orders/{id}")
		return nil, api_error
	}

	order_res = DBModel_To_ResSchema(order_model)

	// Adding the result to the cache service
	err = cache.SetJSON(ctx, cache_key, order_res)
	if err != nil {
		logger.Error().Err(err).Msg("Couldn't Cache.SetJSON() in GET /orders/{id}")
	}
	return &order_res, nil
}

func GetOrderForOps_Repository(ctx context.Context, id uuid.UUID) (*database.Order, error) {
	order_model, err := gorm.G[database.Order](database.Conn).
		Where("id = ?", id).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		api_error := &errs.APIError{Status: http.StatusNotFound, CausedIn: "repository", Err: err}
		return nil, api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in GET /orders/{id}")
		return nil, api_error
	}
	return &order_model, nil
}

func CreateOneOrder_Repository(ctx context.Context, input *CreateOneOrder_Req) (*OneOrder_Res, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.Order](database.Conn).
		Create(ctx, &data)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /orders")
		return nil, api_error
	}

	new_order := DBModel_To_ResSchema(data)
	return &new_order, nil
}

func CreateManyOrders_Repository(ctx context.Context, input *CreateManyOrders_Req) (*schemas.CreateMany_Res_Body[OneOrder_Res], error) {
	var CreatedItems []OneOrder_Res
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	new_data := ReqSchemas_To_DBModels(input.Body)
	for i, item := range new_data {
		err := gorm.G[database.Order](database.Conn).Create(ctx, &item)

		if err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Foreign Key Error"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /orders/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_order := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_order)
	}

	return &schemas.CreateMany_Res_Body[OneOrder_Res]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func AddOnePrint_Repository(ctx context.Context, order_model *database.Order, data PrintItem_Req) (*PrintItem_Res, error) {
	print_model := NewPrintModel(data, &order_model.ID, order_model.UserID)

	err := gorm.G[database.Print](
		database.Conn,
	).Create(ctx, &print_model)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return nil, api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in POST /orders")
		return nil, api_error
	}

	cache_key := cache.FormatKeyByID("orders", order_model.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in POST /orders/{id}/prints")
	}

	new_print := NewPrintRes(print_model)
	return &new_print, nil
}

func AddManyPrints_Repository(ctx context.Context, order_model *database.Order, data []PrintItem_Req) (*schemas.CreateMany_Res_Body[PrintItem_Res], error) {
	var CreatedItems []PrintItem_Res
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	for i, item := range data {
		print_model := NewPrintModel(item, &order_model.ID, order_model.UserID)
		err := gorm.G[database.Print](
			database.Conn,
		).Create(ctx, &print_model)

		if err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Foreign Key Error"})
			} else {
				logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints/many.")
				InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			}
			continue
		}

		new_print := NewPrintRes(print_model)
		CreatedItems = append(CreatedItems, new_print)
	}

	cache_key := cache.FormatKeyByID("orders", order_model.ID)
	err := cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in POST /orders/{id}/prints/many")
	}

	return &schemas.CreateMany_Res_Body[PrintItem_Res]{
		CreatedItems: CreatedItems,
		SuccessCount: len(CreatedItems),
		InvalidItems: InvalidItems,
	}, nil
}

func UpdateOrder_Repository(ctx context.Context, order_model *database.Order, data *UpdateOrder_Req_Body) error {
	// Ensuring Data Integrity: we check order.status & order.is_updateable,

	// but because go don't allow us to directly compare them to values without pointers
	// and it raises error if we made a point to nil, so we check them for nil
	// and if it's nil, then we assign them for order's existing values
	if data.IsUpdateable == nil {
		data.IsUpdateable = &order_model.IsUpdateable
	}
	if data.Status == nil {
		data.Status = &order_model.Status
	}

	// If the order is aborted or marked as completed, then we make sure that is_updateable is False
	if *data.Status == database.StatusEnum_Aborted || *data.Status == database.StatusEnum_Completed {
		*data.IsUpdateable = false
	} else if *data.Status == database.StatusEnum_InProgress {
		*data.IsUpdateable = true
	} else if *data.IsUpdateable {
		// if it want to make is_updateable true, then we make sure status == "in progress".
		// We don't need to worry about the user setting is_updateable to true, as we raise Auth error if it's false above
		*data.Status = database.StatusEnum_InProgress
	}

	if data.Name != nil {
		order_model.Name = *data.Name
	}
	if data.Phone != nil {
		order_model.Phone = *data.Phone
	}
	if data.Address != nil {
		order_model.Address = *data.Address
	}
	if data.DeliverySchedule != nil {
		order_model.DeliverySchedule = data.DeliverySchedule
	}
	if data.Reviewed != nil {
		order_model.Reviewed = *data.Reviewed
	}
	if data.IsUpdateable != nil {
		order_model.IsUpdateable = *data.IsUpdateable
	}
	if data.Status != nil {
		order_model.Status = *data.Status
	}
	if data.UserID != nil {
		order_model.UserID = data.UserID
	}

	err := database.Conn.Save(&order_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /orders/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("orders", order_model.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /orders/{id}")
	}

	return nil

}

func UpdatePrint_Repository(ctx context.Context, order *database.Order, input *UpdatePrint_Req) error {
	print_model, err := gorm.G[database.Print](database.Conn).
		Where("id = ? AND order_id = ?", input.PrintID, input.OrderID).
		First(ctx)

	if input.Body.FontType != nil {
		print_model.FontType = *input.Body.FontType
	}
	if input.Body.FontColor != nil {
		print_model.FontColor = *input.Body.FontColor
	}
	if input.Body.OutfitType != nil {
		print_model.OutfitType = *input.Body.OutfitType
	}
	if input.Body.OutfitColor != nil {
		print_model.OutfitColor = *input.Body.OutfitColor
	}
	if input.Body.FontColor != nil {
		print_model.FontColor = *input.Body.FontColor
	}
	if input.Body.Qoute != nil {
		print_model.Qoute = input.Body.Qoute
	}
	if input.Body.Verses != nil {
		print_model.Verses = input.Body.Verses
	}
	if input.Body.IsCouplet != nil {
		print_model.IsCouplet = input.Body.IsCouplet
	}
	if input.Body.PoemID != nil {
		print_model.PoemID = input.Body.PoemID
	}
	if input.Body.ChosenVerseID != nil {
		print_model.ChosenVerseID = input.Body.ChosenVerseID
	}
	if input.Body.ProseQouteID != nil {
		print_model.ProseQouteID = input.Body.ProseQouteID
	}

	err = database.Conn.Save(&print_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, Message: "Foreign Key Error", CausedIn: "repository", Err: err}
		return api_error
	}
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in PUT /orders/{order_id}/prints/{print_id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("orders", order.ID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in PUT /orders/{order_id}/prints/{print_id}")
	}

	return nil
}

func DeleteOrder_Repository(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[database.Print](database.Conn).
		Where("order_id = ?", id).
		Delete(ctx)
	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /orders/{id}")
		return api_error
	}

	_, err = gorm.G[database.Order](database.Conn).
		Where("id = ?", id).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /orders/{id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("orders", id)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /orders/{id}")
	}

	return nil
}

func DeletePrint_Repository(ctx context.Context, input *DeletePrint_Req) error {
	_, err := gorm.G[database.Print](database.Conn).
		Where("id = ? AND order_id = ?", input.PrintID, input.OrderID).
		Delete(ctx)

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		api_error := &errs.APIError{Status: http.StatusConflict, CausedIn: "repository", Err: err}
		return api_error
	}

	if err != nil {
		api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "repository", Err: err}
		logger.Error().Err(api_error).Msg("Error in DELETE /orders/{order_id}/prints/{print_id}")
		return api_error
	}

	cache_key := cache.FormatKeyByID("orders", input.OrderID)
	err = cache.DelKey(ctx, cache_key)
	if err != nil {
		logger.Error().Err(err).Msg("Error with Cache.DelKey() in DELETE /orders/{order_id}/prints/{print_id}")
	}

	return nil
}
