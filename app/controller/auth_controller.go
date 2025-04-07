package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wato787/app/dto"
	"github.com/wato787/app/repository"
	"github.com/wato787/app/service"
)

type AuthController struct {
	authService *service.AuthService
}


func NewAuthController() *AuthController {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	
	return &AuthController{
		authService: authService,
	}
}


func (ac *AuthController) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "無効なリクエスト形式: " + err.Error()})
		return
	}

	// ユーザー登録処理
	user, err := ac.authService.RegisterUser(req.Email, req.Password)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusCreated, dto.SignupResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}


func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "無効なリクエスト形式: " + err.Error()})
		return
	}


	user, token, err := ac.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}