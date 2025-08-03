package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "messaging-server/docs"
	"messaging-server/internal/cron"
	"messaging-server/internal/handler"
)

// initEngine initializes the Gin engine without any routes
func initEngine() *gin.Engine {
	return gin.Default()
}

// SetupRouter configures all routes under /api/v1 and returns the engine
func SetupRouter(cronJob *cron.Cron) *gin.Engine {
	r := initEngine()

	// index endpoint
	r.GET("/", handler.Index)

	// healthcheck endpoint
	r.GET("/health", handler.Health)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// group all endpoints under /api/v1
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// cron control endpoint
			v1.POST("/cron/control", handler.CronHandler(cronJob))

			// list sent messages endpoint
			v1.GET("/list/sent-messages", handler.ListMessageHandler())
		}

	}

	return r
}
