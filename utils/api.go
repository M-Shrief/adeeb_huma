package utils

type BaseResponse_JSONBody struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Body   BaseResponse_JSONBody
	Status int
}
