package main

import (
	"context"
	"messaging-server/internal/configs"
	"messaging-server/internal/cron"
	"messaging-server/internal/database"
	"messaging-server/internal/jobs"
	log "messaging-server/internal/logging"
	"messaging-server/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.InitLogger()

	// initialize Postgres connection
	if err := database.ConnectPostgres(); err != nil {
		log.Logger.Fatalf("failed to connect to Postgres: %v", err)
	}

	// initialize Redis connection
	if err := database.ConnectRedis(); err != nil {
		log.Logger.Fatalf("failed to connect to Redis: %v", err)
	}

	// initialize cron job
	cronJob, err := cron.NewCron(jobs.SendMessageJob)
	if err != nil {
		log.Logger.Fatalf("failed to create cron job: %v", err)
	}

	log.Logger.Infoln("Starting messaging server...")

	// initialize Gin router with all endpoints
	r := router.SetupRouter(cronJob)

	// immediately start the cron job
	cronJob.Start()

	// create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatalf("listen: %s", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server and cron job
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Logger.Infoln("Shutting down server...")

	// stop cron job
	cronJob.Stop()

	// shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.AppConfig.ServerGracePeriod)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Logger.Infoln("Server exiting")
}
