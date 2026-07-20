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
