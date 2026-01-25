package models

type Response struct {
	Message string `json:"message" example:"Success"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message" example:"The request contains invalid data"`
}
