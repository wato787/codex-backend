package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wato787/app/dto"
	"github.com/wato787/app/repository"
	"github.com/wato787/app/service"
)

// AuthMiddleware は認証関連のミドルウェア
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware は新しいAuthMiddlewareを作成する
func NewAuthMiddleware() *AuthMiddleware {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth は認証が必要なエンドポイント用のミドルウェア
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authorization ヘッダーが必要です"})
			c.Abort()
			return
		}

		// "Bearer "プレフィックスを確認して削除
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authorization ヘッダーの形式は 'Bearer {token}' である必要があります"})
			c.Abort()
			return
		}

		// トークンを検証
		userID, err := m.authService.ValidateJWTToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "無効または期限切れのトークンです"})
			c.Abort()
			return
		}

		// ユーザーが存在するか確認
		user, err := m.authService.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "ユーザーが見つかりません"})
			c.Abort()
			return
		}

		// コンテキストにユーザー情報を保存
		c.Set("user", user)
		c.Set("userID", userID)
		
		// 次のハンドラーに進む
		c.Next()
	}
}