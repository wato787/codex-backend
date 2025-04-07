package dto

// SignupRequest はサインアップリクエスト用のDTO
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignupResponse はサインアップレスポンス用のDTO
type SignupResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// LoginRequest はログインリクエスト用のDTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse はログインレスポンス用のDTO
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse はユーザー情報レスポンス用のDTO
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

