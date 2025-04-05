package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wato787/app/controller"
	"github.com/wato787/app/models"
	"github.com/wato787/db"
	"github.com/wato787/docs"
)

func main() {

	// DB接続
	SetupDB()

	// ルーターのセットアップ
	route := SetupRoutes()
	docs.SwaggerInfo.BasePath = "/api"
	api := route.Group("/api")
	{
		api.GET("/hello", controller.HelloWorld)
	}
	route.GET("/doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	route.Run(":8080")
	
}

func SetupRoutes() *gin.Engine {
	route := gin.Default()
	route.GET("/", controller.HelloWorld)
	return route
}

func SetupDB() {
	if err := db.Connect(db.DefaultConfig()); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.DB.AutoMigrate(&models.User{})
}