package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api

// PingExample godoc
// @Summary liveness probe
// @Schemes
// @Description do ping
// @Tags Hello World
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /hello [get]
func HelloWorld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}