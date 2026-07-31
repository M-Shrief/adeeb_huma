package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func GetAllOrders_Handler(ctx context.Context, input *GetAllOrders_Req) (*schemas.GetAll_Res[schemas.Order_Descriptive], error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	var results []OrderWithTotalCount

	err = database.Conn.Table("orders").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /orders.")
		return nil, huma.Error404NotFound("Unknown error while getting orders")
	}

	orders, total_count := DistillDBModelsWithCount(results)
	res := &schemas.GetAll_Res[schemas.Order_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.Order_Descriptive]{
			Data:       orders,
			Limit:      input.Limit,
			Offset:     input.Offset,
			TotalCount: total_count,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetUserOrders_Handler(ctx context.Context, input *GetAllOrders_Req) (*schemas.GetAll_Res[schemas.Order_Descriptive], error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Normal, auth.OPEnum_Read),
	}
	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_claim := claims["user"].(map[string]interface{})
	user_id := user_claim["id"].(string)

	var results []OrderWithTotalCount
	err = database.Conn.Table("orders").
		Select("*, COUNT(*) OVER() as total_count").
		Limit(input.Limit).
		Offset(input.Offset).
		Where("user_id = ?", user_id).
		Find(&results).Error

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /orders.")
		return nil, huma.Error404NotFound("Unknown error while getting orders")
	}

	orders, total_count := DistillDBModelsWithCount(results)
	res := &schemas.GetAll_Res[schemas.Order_Descriptive]{
		Body: schemas.GetAll_Res_Body[schemas.Order_Descriptive]{
			Data:       orders,
			Limit:      input.Limit,
			Offset:     input.Offset,
			TotalCount: total_count,
		},
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOrderByID_Handler(ctx context.Context, input *GetOrderByID_Req) (*GetOrderByID_Res, error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Preload("Prints", func(db gorm.PreloadBuilder) error {
			db.Select("*")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in GET /orders/{id}.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
	}

	order_res := DBModel_To_ResSchema(order_model)
	res := &GetOrderByID_Res{
		Body:   order_res,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneOrder_Handler(ctx context.Context, input *CreateOneOrder_Req) (*CreateOneOrder_Res, error) {
	data := ReqSchema_To_DBModel(input.Body)

	err := gorm.G[database.Order](database.Conn).
		Create(ctx, &data)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders.")
		return nil, huma.Error400BadRequest("Bad Request creating Order.")
	}

	new_order := DBModel_To_ResSchema(data)
	return &CreateOneOrder_Res{Body: new_order, Status: http.StatusCreated}, nil
}

func CreateManyOrder_Handler(ctx context.Context, input *CreateManyOrder_Req) (*CreateManyOrders_Res, error) {
	var CreatedItems []OneOrder_Res
	var InvalidItems []schemas.CreateMany_Res_Body_InvalidItem

	data := ReqSchemas_To_DBModels(input.Body)

	for i, item := range data {
		err := gorm.G[database.Order](database.Conn).Create(ctx, &item)

		if err != nil {
			logger.Error().Err(err).Msg("Unknown errror in POST /orders/many.")
			InvalidItems = append(InvalidItems, schemas.CreateMany_Res_Body_InvalidItem{ItemIndex: i, Message: "Bad Request, try again later"})
			continue
		}

		new_order := DBModel_To_ResSchema(item)
		CreatedItems = append(CreatedItems, new_order)
	}

	return &CreateManyOrders_Res{
		Body: schemas.CreateMany_Res_Body[OneOrder_Res]{
			CreatedItems: CreatedItems,
			SuccessCount: len(CreatedItems),
			InvalidItems: InvalidItems,
		},
		Status: http.StatusCreated}, nil
}

func AddPrint_Handler(ctx context.Context, input *AddPrint_Req) (*AddPrint_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Preload("Prints", func(db gorm.PreloadBuilder) error {
			db.Select("*")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
	}

	print_model := NewPrintModel(input.Body, &input.ID, order_model.UserID)

	err = gorm.G[database.Print](
		database.Conn,
	).Create(ctx, &print_model)
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error creating User in POST /orders/{id}/prints")
		return nil, huma.Error400BadRequest("Bad Request adding print.")
	}

	return &AddPrint_Res{Body: NewPrintRes(print_model), Status: http.StatusCreated}, nil

}

func UpdateOrder_Handler(ctx context.Context, input *UpdateOrder_Req) (*schemas.Update_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Preload("Prints", func(db gorm.PreloadBuilder) error {
			db.Select("*")
			return nil
		}).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
		// Make sure the user can't update interanl values like: is_updateable, is_completed
		// we assign them to None, as we can exclude them later with model_dump(exclude_none=True)
		input.Body.IsUpdateable = nil
		input.Body.Status = nil
		input.Body.UserID = nil
		input.Body.Reviewed = nil
	}

	// Ensuring Data Integrity: we check order.status & order.is_updateable,

	// but because go don't allow us to directly compare them to values without pointers
	// and it raises error if we made a point to nil, so we check them for nil
	// and if it's nil, then we assign them for order's existing values
	if input.Body.IsUpdateable == nil {
		input.Body.IsUpdateable = &order_model.IsUpdateable
	}
	if input.Body.Status == nil {
		input.Body.Status = &order_model.Status
	}

	// If the order is aborted or marked as completed, then we make sure that is_updateable is False
	if *input.Body.Status == database.StatusEnum_Aborted || *input.Body.Status == database.StatusEnum_Completed {
		*input.Body.IsUpdateable = false
	} else if *input.Body.Status == database.StatusEnum_InProgress {
		*input.Body.IsUpdateable = true
	} else if *input.Body.IsUpdateable {
		// if it want to make is_updateable true, then we make sure status == "in progress".
		// We don't need to worry about the user setting is_updateable to true, as we raise Auth error if it's false above
		*input.Body.Status = database.StatusEnum_InProgress
	}

	if input.Body.Name != nil {
		order_model.Name = *input.Body.Name
	}
	if input.Body.Phone != nil {
		order_model.Phone = *input.Body.Phone
	}
	if input.Body.Address != nil {
		order_model.Address = *input.Body.Address
	}
	if input.Body.DeliverySchedule != nil {
		order_model.DeliverySchedule = input.Body.DeliverySchedule
	}
	if input.Body.Reviewed != nil {
		order_model.Reviewed = *input.Body.Reviewed
	}
	if input.Body.IsUpdateable != nil {
		order_model.IsUpdateable = *input.Body.IsUpdateable
	}
	if input.Body.Status != nil {
		order_model.Status = *input.Body.Status
	}
	if input.Body.UserID != nil {
		order_model.UserID = input.Body.UserID
	}

	err = database.Conn.Save(&order_model).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, huma.Error400BadRequest("foreign key error")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error updating order in PUt /orders/{id}")
		return nil, huma.Error400BadRequest("Bad Request updating order.")
	}

	return &schemas.Update_Res{Status: http.StatusNoContent}, nil

}

func UpdatePrint_Handler(ctx context.Context, input *UpdatePrint_Req) (*schemas.Update_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Where("id = ?", input.OrderID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
	}

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
		return nil, huma.Error400BadRequest("foreign key error")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error updating print in PUt /orders/{order_id}/prints/{print_id}")
		return nil, huma.Error400BadRequest("Bad Request updating print.")
	}

	return &schemas.Update_Res{Status: http.StatusNoContent}, nil

}

func DeleteOrder_Handler(ctx context.Context, input *DeleteOrder_Req) (*schemas.Delete_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Where("id = ?", input.ID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
	}

	_, err = gorm.G[database.Print](database.Conn).
		Where("order_id = ?", input.ID).
		Delete(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Unknown error deleting order in DELETE /orders/{id}")
		return nil, huma.Error400BadRequest("Bad Request deleting order.")
	}

	_, err = gorm.G[database.Order](database.Conn).
		Where("id = ?", input.ID).
		Delete(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown error deleting order in DELETE /orders/{id}")
		return nil, huma.Error400BadRequest("Bad Request deleting order.")
	}

	return &schemas.Delete_Res{Status: http.StatusNoContent}, nil

}

func DeletePrint_Handler(ctx context.Context, input *DeletePrint_Req) (*schemas.Delete_Res, error) {
	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	order_model, err := gorm.G[database.Order](database.Conn).
		Where("id = ?", input.OrderID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("Order's not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("Unknown errror in POST /orders/{id}/prints.")
		return nil, huma.Error400BadRequest("Bad Request getting Order")
	}

	is_adminstrator := auth.CheckAdminstration(user_permissions, auth.OPEnum_Read)
	if is_adminstrator == false {
		if auth.CheckOwnership(order_model.UserID, claims) == false {
			return nil, huma.Error401Unauthorized("Not Authorizaed")
		}
	}

	_, err = gorm.G[database.Print](database.Conn).
		Where("id = ? AND order_id = ?", input.PrintID, input.OrderID).
		Delete(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Unknown error deleting order in DELETE /orders/{id}")
		return nil, huma.Error400BadRequest("Bad Request deleting order.")
	}

	return &schemas.Delete_Res{Status: http.StatusNoContent}, nil

}
