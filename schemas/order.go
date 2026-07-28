package schemas

import (
	"adeeb_huma/database"
	"time"
)

type Order_Full struct {
	IDField
	Order_NameField
	Order_PhoneField
	Order_AddressField
	ReviewedField
	Order_IsUpdateableField
	Order_StatusField
	Order_DeliveryScheduleField_Optional

	CreatedAtField
	UpdatedAtField
	// Relations
	UserIDField_Optional
}

type Order_Descriptive struct {
	IDField
	Order_NameField
	Order_PhoneField
	Order_AddressField
	ReviewedField
	Order_IsUpdateableField
	Order_StatusField
	Order_DeliveryScheduleField_Optional
	// Relations
	UserIDField_Optional
}

type Order_Minimal struct {
	IDField
	Order_StatusField
	Order_DeliveryScheduleField_Optional
	// Relations
	UserIDField_Optional
}

type Order_NameField struct {
	Name string `json:"name" doc:"User's name" minLength:"4" maxLength:"256"`
}
type Order_NameField_Optional struct {
	Name *string `json:"name" doc:"User's name" minLength:"4" maxLength:"256" required:"false"`
}

type Order_PhoneField struct {
	Phone string `json:"phone" doc:"User's Phone" minLength:"4" maxLength:"128"`
}
type Order_PhoneField_Optional struct {
	Phone *string `json:"phone" doc:"User's Phone" minLength:"4" maxLength:"128" required:"false"`
}

type Order_AddressField struct {
	Address string `json:"address" doc:"User's address" minLength:"4" maxLength:"256"`
}
type Order_AddressField_Optional struct {
	Address *string `json:"address" doc:"User's address" minLength:"4" maxLength:"256" required:"false"`
}

type Order_IsUpdateableField struct {
	IsUpdateable bool `json:"is_updateable" doc:"is it updateable?"`
}
type Order_IsUpdateableField_Optional struct {
	IsUpdateable *bool `json:"is_updateable" doc:"is it updateable?" required:"false"`
}

type Order_StatusField struct {
	Status database.StatusEnum `json:"status" enum:"in progress,aborted,completed" doc:"Order's status"`
}
type Order_StatusField_Optional struct {
	Status *database.StatusEnum `json:"status" enum:"in progress,aborted,completed" doc:"Order's status" required:"false"`
}

type Order_DeliveryScheduleField_Optional struct {
	DeliverySchedule *time.Time `json:"delivery_schedule" required:"false"`
}
