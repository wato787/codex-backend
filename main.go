package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wato787/app/controller"
	db "github.com/wato787/app/database"
	"github.com/wato787/app/model"

	"github.com/wato787/docs"
)


func main() {
	// DB接続
	SetupDB()
	defer db.Close()

	// ルーターのセットアップ
	route := SetupRoutes()
	route.Run(":8080")
}

func SetupRoutes() *gin.Engine {
	route := gin.Default()
	
	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	route.GET("/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	
	// コントローラーとミドルウェアのインスタンス作成
	authController := controller.NewAuthController()
	// authMiddleware := middleware.NewAuthMiddleware()
	// TODO : ミドルウェア適用
	

	
	// APIグループ
	api := route.Group("/api")
	{
		// 認証関連のエンドポイント（認証不要）
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authController.Signup)
			auth.POST("/login", authController.Login)
		}
		
		// 公開エンドポイント
		api.GET("/health", controller.HealthCheck)
	}
	
	return route
}



func SetupDB() {
	if err := db.Connect(db.DefaultConfig()); err != nil {
		log.Fatal(err)
	}
	
	// マイグレーション
	if err := db.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	
	log.Println("Database migration completed")
}