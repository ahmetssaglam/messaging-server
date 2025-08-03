package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"messaging-server/internal/configs"
	"net/http"
)

// Index     	 godoc
// @Summary      Welcome message
// @Description  Returns a welcome message including the application’s configured name.
// @Tags         Base
// @Produce      json
// @Success      200  "Welcome to App!"
// @Router       / [get]
func Index(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome to %s!", configs.AppConfig.AppName)})
}

// Health     	 godoc
// @Summary      Health check
// @Description  Simple endpoint to verify the service is running.
// @Tags         Base
// @Produce      json
// @Success      200  "Healthy!"
// @Router       /health [get]
func Health(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Healthy!"})
}
