package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wato787/app/controller"
	"github.com/wato787/docs"
)

func main() {
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