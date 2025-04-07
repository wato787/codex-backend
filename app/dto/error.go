package dto

// ErrorResponse はエラーレスポンス用のDTO
type ErrorResponse struct {
	Error string `json:"error"`
}