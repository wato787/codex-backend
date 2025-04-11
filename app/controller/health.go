package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func HealthCheck(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"status": "healthy"})
}