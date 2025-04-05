package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config はデータベース接続設定を保持する構造体
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// DefaultConfig はデフォルトの接続設定を返す
func DefaultConfig() *Config {
	return &Config{
		User:     "root",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     "3306",
		DBName:   "codex",
	}
}

// DSN はデータソース名（接続文字列）を生成する
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// Connect はデータベースに接続する
func Connect(config *Config) error {
	dsn := config.DSN()
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %w", err)
	}
	
	log.Println("データベースに接続しました")
	return nil
}

// Close はデータベース接続を閉じる
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("データベース接続オブジェクト取得エラー: %v", err)
			return
		}
		
		if err := sqlDB.Close(); err != nil {
			log.Printf("データベース接続クローズエラー: %v", err)
			return
		}
		
		log.Println("データベース接続を閉じました")
	}
}