package handler

import (
	"github.com/gin-gonic/gin"
	"messaging-server/internal/cron"
	log "messaging-server/internal/logging"
	"messaging-server/internal/models"
	"net/http"
)

// CronHandler handles both starting and stopping the cron based on the "action" field.
// @Summary      Control cron job
// @Description  Start or stop the cron based on the "action" field.
// @Tags         Cron
// @Accept       json
// @Produce      json
// @Param        payload  body      models.CronRequest  true  "start or stop"
// @Success      200        "Cron job started"
// @Success      202        "Cron job will be stopped"
// @Failure      400        "Invalid request payload"
// @Router       /api/v1/cron/control [post]
func CronHandler(cronJob *cron.Cron) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CronRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
			return
		}

		switch req.Action {
		case "start":
			cronJob.Start()
			c.JSON(http.StatusOK, gin.H{"message": "Cron job started"})
		case "stop":
			go func() {
				cronJob.Stop()
				log.Logger.Info("background cron.Stop() goroutine finished")
			}()
			c.JSON(http.StatusAccepted, gin.H{"message": "Cron job will be stopped"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "action must be 'start' or 'stop'"})
		}
	}
}
