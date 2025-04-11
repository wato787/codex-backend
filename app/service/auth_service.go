package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wato787/app/model"
	"github.com/wato787/app/repository"
)

// JWTClaims はJWTのペイロード部分
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 定数
const (
	jwtSecret = "your-secret-key-here" // 本番環境では環境変数から取得すべき
	tokenTTL  = 24 * time.Hour         // トークンの有効期限
)

// カスタムエラー
var (
	ErrEmailAlreadyExists  = errors.New("このメールアドレスは既に使用されています")
	ErrInvalidCredentials  = errors.New("メールアドレスまたはパスワードが無効です")
	ErrUserNotFound        = errors.New("ユーザーが見つかりません")
	ErrTokenGenerationFail = errors.New("トークン生成に失敗しました")
	ErrInvalidToken        = errors.New("無効なトークンです")
)

// AuthService は認証関連の機能を提供する
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService は新しいAuthServiceのインスタンスを作成する
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// RegisterUser は新規ユーザーを登録する
func (s *AuthService) RegisterUser(email, password string) (*model.User, error) {
	// メールアドレスの重複チェック
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("メールアドレス確認中にエラーが発生しました: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// 新規ユーザー作成
	user := &model.User{
		Email:    email,
		Password: password,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("ユーザー作成に失敗しました: %w", err)
	}

	return user, nil
}

// LoginUser はユーザーを認証しJWTトークンを発行する
func (s *AuthService) LoginUser(email, password string) (*model.User, string, error) {
	// ユーザーを検索
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("ユーザー検索中にエラーが発生しました: %w", err)
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// パスワードを確認
	if !user.CheckPassword(password) {
		return nil, "", ErrInvalidCredentials
	}

	// JWTトークンを生成
	token, err := s.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", ErrTokenGenerationFail
	}

	return user, token, nil
}

// GenerateJWTToken はユーザーIDからJWTトークンを生成する
func (s *AuthService) GenerateJWTToken(userID uint) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(tokenTTL)

	// クレームを作成
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 署名してトークンを文字列化
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken はJWTトークンを検証し、ユーザーIDを返す
func (s *AuthService) ValidateJWTToken(tokenString string) (uint, error) {
	// トークンをパース
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方法の確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// クレームを取得
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, ErrInvalidToken
}

// GetUserByID はIDによるユーザー取得
func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索中にエラーが発生しました: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}