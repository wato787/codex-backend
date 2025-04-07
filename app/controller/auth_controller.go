package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wato787/app/dto"
	"github.com/wato787/app/repository"
	"github.com/wato787/app/service"
)

// AuthController は認証関連のコントローラー
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController は新しいAuthControllerを作成する
func NewAuthController() *AuthController {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	
	return &AuthController{
		authService: authService,
	}
}

// Signup は新規ユーザー登録を行うハンドラー
// @Summary ユーザー登録API
// @Description 新規ユーザーを登録する
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "サインアップ情報"
// @Success 201 {object} dto.SignupResponse "登録成功"
// @Failure 400 {object} dto.ErrorResponse "リクエスト不正"
// @Failure 409 {object} dto.ErrorResponse "ユーザー既存"
// @Failure 500 {object} dto.ErrorResponse "サーバーエラー"
// @Router /api/auth/signup [post]
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

// Login はユーザーログインを行うハンドラー
// @Summary ユーザーログインAPI
// @Description ユーザー認証とトークン発行を行う
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "ログイン情報"
// @Success 200 {object} dto.LoginResponse "ログイン成功"
// @Failure 400 {object} dto.ErrorResponse "リクエスト不正"
// @Failure 401 {object} dto.ErrorResponse "認証失敗"
// @Failure 500 {object} dto.ErrorResponse "サーバーエラー"
// @Router /api/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "無効なリクエスト形式: " + err.Error()})
		return
	}

	// ログイン処理
	user, token, err := ac.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}