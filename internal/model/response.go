package model

type ApiResponse struct {
	Code    int    `json:"code" example:"200"`
	Type    string `json:"type" example:"success"`
	Message string `json:"message" example:"operation completed successfully"`
}
