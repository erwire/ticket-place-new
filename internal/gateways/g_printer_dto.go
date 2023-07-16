package gateways

type ResponseDTO struct {
	Message interface{} `json:"message"`
	Error   string      `json:"error"`
}

func NewResponseDTO() *ResponseDTO {
	return &ResponseDTO{}
}
