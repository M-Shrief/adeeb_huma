package utils

type BaseResponse struct {
	Body   BaseResponse_JSONBody
	Status int
}

type BaseResponse_JSONBody struct {
	Message string `json:"message"`
}
