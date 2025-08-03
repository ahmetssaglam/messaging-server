package handler

import (
	"github.com/gin-gonic/gin"
	"messaging-server/internal/database"
	"net/http"
)

// ListMessageHandler gets all messages that have been sent and returns a JSON response.
// @Summary      List sent messages
// @Description  Retrieves all messages that have been sent.
// @Tags         Messages
// @Produce      json
// @Success      200  {object} []models.Message      "Messages fetched successfully"
// @Failure      500   "failed to fetch messages"
// @Router       /api/v1/list/sent-messages [get]
func ListMessageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		msgs, err := database.PostgresConnection.FetchAllSentMessages()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Messages fetched successfully", "data": msgs})
	}
}
