package configs

import (
	"messaging-server/internal/models"
	pkgUtils "messaging-server/pkg/utils"
)

// AppConfig holds the application configuration settings.
var AppConfig = models.AppConfigStruct{
	AppName:           pkgUtils.GetEnvStr("APP_NAME", "Messaging Server V1"),
	LogLevel:          pkgUtils.GetEnvStr("LOG_LEVEL", "DEBUG"),
	WebhookURL:        pkgUtils.GetEnvStr("WEBHOOK_URL", ""),
	ServerGracePeriod: pkgUtils.GetEnvInt("SERVER_GRACE_PERIOD", 30),
	MessageFetchLimit: pkgUtils.GetEnvInt("MESSAGE_FETCH_LIMIT", 2),
	CronInterval:      pkgUtils.GetEnvInt("CRON_INTERVAL", 120),
	MaxConcurrentJobs: pkgUtils.GetEnvInt("MAX_CONCURRENT_JOBS", 5),
}
