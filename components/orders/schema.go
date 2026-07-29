package orders

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
	"time"
)

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
	Prints []PrintItem_Res `json:"prints"`
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
	// if order_req.UserID != nil {
	// 	order_model.UserID = *order_req.UserID
	// }

	for _, print := range order_req.Prints {
		new_print := database.Print{
			FontType:    print.FontType,
			FontColor:   print.FontColor,
			OutfitType:  print.OutfitType,
			OutfitColor: print.OutfitColor,
		}
		if print.Qoute != nil {
			new_print.Qoute = print.Qoute
		}
		if print.Verses != nil {
			new_print.Verses = print.Verses
		}
		if print.IsCouplet != nil {
			new_print.IsCouplet = print.IsCouplet
		}
		if print.PoemID != nil {
			new_print.PoemID = print.PoemID
		}
		if print.ChosenVerseID != nil {
			new_print.ChosenVerseID = print.ChosenVerseID
		}
		if print.ProseQouteID != nil {
			new_print.ProseQouteID = print.ProseQouteID
		}

		if order_req.UserID != nil {
			new_print.UserID = order_req.UserID
		}

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
		var print_res PrintItem_Res
		print_res.ID = print_model.ID
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
