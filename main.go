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
	// プログラム終了時にDBをクローズ
	defer db.Close()

	// ルーターのセットアップ
	route := SetupRoutes()
	docs.SwaggerInfo.BasePath = "/api"
	api := route.Group("/api")
	{
		api.GET("/health", controller.HealthCheck)
	}
	route.GET("/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	route.Run(":8080")
}

func SetupRoutes() *gin.Engine {
	route := gin.Default()
	return route
}

func SetupDB() {
	if err := db.Connect(db.DefaultConfig()); err != nil {
		log.Fatal(err)
	}
	db.DB.AutoMigrate(&model.User{})
}