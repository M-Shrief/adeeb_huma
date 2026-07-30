package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
	"time"

	"github.com/google/uuid"
)

type GetAllOrders_Req struct {
	schemas.GetAll_Req
	schemas.AuthHeader
}

type GetOrderByID_Req struct {
	schemas.AuthHeader
	schemas.IDPath
}

type GetOrderByID_Res struct {
	Body   OneOrder_Res
	Status int
}

type CreateOneOrder_Req struct {
	Body CreateOneOrder_Req_Body
}
type CreateOneOrder_Req_Body struct {
	schemas.Order_NameField
	schemas.Order_PhoneField
	schemas.Order_AddressField
	// We assign a standard delivery schedule in the server
	// schemas.Order_DeliveryScheduleField_Optional
	// Rlations
	schemas.UserIDField_Optional

	Prints []PrintItem_Req `json:"prints"`
}

type CreateOneOrder_Res struct {
	Body   OneOrder_Res
	Status int
}

type CreateManyOrder_Req struct {
	Body []CreateOneOrder_Req_Body
}

type CreateManyOrders_Res struct {
	Body   schemas.CreateMany_Res_Body[OneOrder_Res]
	Status int
}

type OneOrder_Res struct {
	schemas.Order_Descriptive
	Prints []PrintItem_Res `json:"prints,omitempty" required:"false"`
}

type AddPrint_Req struct {
	Body PrintItem_Req
	schemas.IDPath
	schemas.AuthHeader
}

type AddPrint_Res struct {
	Body   PrintItem_Res
	Status int
}

type PrintItem_Req struct {
	schemas.Print_FontTypeField
	schemas.Print_FontColorField
	schemas.Print_OutfitColorField
	schemas.Print_OutfitTypeField

	schemas.ProseQoute_QouteField_Optional
	schemas.VersesField_Optional
	schemas.IsCoupletField_Optional

	schemas.PoemIDField_Optional
	schemas.ChosenVerseIDField_Optional
	schemas.ProseQouteIDField_Optional
}

type PrintItem_Res struct {
	schemas.IDField
	schemas.Print_FontTypeField
	schemas.Print_FontColorField
	schemas.Print_OutfitColorField
	schemas.Print_OutfitTypeField
	// Relations
	schemas.ProseQoute_QouteField_Optional
	schemas.VersesField_Optional
	schemas.IsCoupletField_Optional

	schemas.PoemIDField_Optional
	schemas.ChosenVerseIDField_Optional
	schemas.ProseQouteIDField_Optional
}

const DeliveryAfter time.Duration = time.Hour * 24 * 7

func ReqSchema_To_DBModel(order_req CreateOneOrder_Req_Body) database.Order {
	DeliverySchedule := time.Now().UTC().Add(DeliveryAfter)
	order_model := database.Order{
		Name:             order_req.Name,
		Phone:            order_req.Phone,
		Address:          order_req.Address,
		Reviewed:         false,
		IsUpdateable:     true,
		Status:           database.StatusEnum_InProgress,
		DeliverySchedule: &DeliverySchedule,
		UserID:           order_req.UserID,
	}

	for _, print := range order_req.Prints {
		new_print := NewPrintModel(print, order_req.UserID)
		order_model.Prints = append(order_model.Prints, new_print)
	}

	return order_model

}

func ReqSchemas_To_DBModels(orders_req []CreateOneOrder_Req_Body) []database.Order {
	var order_models []database.Order

	for _, order_req := range orders_req {
		order_model := ReqSchema_To_DBModel(order_req)
		order_models = append(order_models, order_model)
	}

	return order_models

}

func DBModel_To_ResSchema(order_model database.Order) OneOrder_Res {
	var order_res OneOrder_Res
	order_res.ID = order_model.ID
	order_res.Name = order_model.Name
	order_res.Phone = order_model.Phone
	order_res.Address = order_model.Address
	order_res.Reviewed = order_model.Reviewed
	order_res.IsUpdateable = order_model.IsUpdateable
	order_res.Status = order_model.Status
	order_res.DeliverySchedule = order_model.DeliverySchedule

	order_res.UserID = order_model.UserID
	for _, print_model := range order_model.Prints {
		print_res := NewPrintRes(print_model)
		order_res.Prints = append(order_res.Prints, print_res)
	}

	return order_res
}

func DBModels_To_ResSchemas(order_models []database.Order) []OneOrder_Res {
	var orders_res []OneOrder_Res

	for _, order_model := range order_models {
		order_res := DBModel_To_ResSchema(order_model)
		orders_res = append(orders_res, order_res)
	}

	return orders_res
}

type OrderWithTotalCount struct {
	database.Order
	TotalCount int64 `json:"total_count"`
}

func DistillDBModelsWithCount(models_with_count []OrderWithTotalCount) ([]schemas.Order_Descriptive, int64) {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var orders_res []schemas.Order_Descriptive
	var total_count int64
	for i, model := range models_with_count {
		var order_res schemas.Order_Descriptive
		order_res.ID = model.ID
		order_res.Name = model.Name
		order_res.Phone = model.Phone
		order_res.Address = model.Address
		order_res.Reviewed = model.Reviewed
		order_res.IsUpdateable = model.IsUpdateable
		order_res.Status = model.Status
		order_res.DeliverySchedule = model.DeliverySchedule

		order_res.UserID = model.UserID

		orders_res = append(orders_res, order_res)
		if i == 0 {
			total_count = model.TotalCount
		}
	}
	return orders_res, total_count
}

func NewPrintModel(print_req PrintItem_Req, user_id *uuid.UUID) database.Print {
	new_print := database.Print{
		FontType:      print_req.FontType,
		FontColor:     print_req.FontColor,
		OutfitType:    print_req.OutfitType,
		OutfitColor:   print_req.OutfitColor,
		Qoute:         print_req.Qoute,
		Verses:        print_req.Verses,
		IsCouplet:     print_req.IsCouplet,
		PoemID:        print_req.PoemID,
		ChosenVerseID: print_req.ChosenVerseID,
		ProseQouteID:  print_req.ProseQouteID,
		UserID:        user_id,
	}
	return new_print
}

func NewPrintRes(print_model database.Print) PrintItem_Res {
	var print_res PrintItem_Res
	print_res.ID = print_model.ID
	print_res.FontType = print_model.FontType
	print_res.FontColor = print_model.FontColor
	print_res.OutfitColor = print_model.OutfitColor
	print_res.OutfitType = print_model.OutfitType
	print_res.Qoute = print_model.Qoute
	print_res.Verses = print_model.Verses
	print_res.IsCouplet = print_model.IsCouplet
	print_res.PoemID = print_model.PoemID
	print_res.ChosenVerseID = print_model.ChosenVerseID
	print_res.ProseQouteID = print_model.ProseQouteID

	return print_res
}
