package repository

import (
	"errors"

	db "github.com/wato787/app/database"
	"github.com/wato787/app/model"
	"gorm.io/gorm"
)

// UserRepository はユーザーデータへのアクセスを提供する
type UserRepository struct{}

// NewUserRepository は新しいUserRepositoryのインスタンスを作成する
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByID はIDによるユーザー検索
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail はメールアドレスによるユーザー検索
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// ExistsByEmail はメールアドレスの重複チェック
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	result := db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Create は新規ユーザー作成
func (r *UserRepository) Create(user *model.User) error {
	return db.DB.Create(user).Error
}

// Update はユーザー情報更新
func (r *UserRepository) Update(user *model.User) error {
	return db.DB.Save(user).Error
}

// Delete はユーザー削除
func (r *UserRepository) Delete(id uint) error {
	return db.DB.Delete(&model.User{}, id).Error
}