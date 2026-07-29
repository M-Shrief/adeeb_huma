package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/internal/auth"
	"adeeb_huma/internal/logger"
	"adeeb_huma/internal/utils"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func GetAllOrders_Handler(ctx context.Context, input *GetAllOrders_Req) (*schemas.GetAll_Res[schemas.Order_Descriptive], error) {

	claims, err := auth.VerifyJWT(input.Auth)
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}
	authorized_list := []string{
		auth.CreatePermission(database.RoleEnum_Management, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_DBA, auth.OPEnum_Read),
		auth.CreatePermission(database.RoleEnum_Analytics, auth.OPEnum_Read),
	}
	user_permissions, err := utils.InterfaceToStringSlice(claims["permissions"])
	if err != nil {
		return nil, huma.Error401Unauthorized("Not Authorizaed")
	}

	is_authorized := auth.CheckPermissions(authorized_list, user_permissions, auth.OPEnum_Read)
	if is_authorized == false {
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
