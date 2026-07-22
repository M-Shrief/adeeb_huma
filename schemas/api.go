package schemas

type BaseResponse struct {
	Body   BaseResponse_JSONBody
	Status int
}

type BaseResponse_JSONBody struct {
	Message string `json:"message"`
}

type GetAll_Req struct {
	Limit  int `query:"limit" default:"100" required:"true" maximum:"100" doc:"limit the number of data items."`
	Offset int `query:"offset" default:"0" required:"true"`
}

type GetAll_Res[T any] struct {
	Body   GetAll_Res_Body[T]
	Status int
}

type GetAll_Res_Body[T any] struct {
	Data   []T `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CreateMany_Res_Body[T any] struct {
	CreatedItems []T                               `json:"created_items"`
	SuccessCount int                               `json:"success_count"`
	InvalidItems []CreateMany_Res_Body_InvalidItem `json:"invalid_items"`
}

type CreateMany_Res_Body_InvalidItem struct {
	ItemIndex int    `json:"item_index"`
	Message   string `json:"message"`
}

type Update_Res struct {
	Status int
}

type Delete_Res struct {
	Status int
}
