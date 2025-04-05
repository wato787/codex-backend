package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api

// HealthCheck godoc
// @Summary Health check endpoint
// @Schemes
// @Description Check the health of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"status": "healthy"})
}