package models

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Name string `json:"name" gorm:"type:varchar(100);not null;unique"`
}