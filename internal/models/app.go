package models

type AppConfigStruct struct {
	AppName           string
	LogLevel          string
	WebhookURL        string
	ServerGracePeriod int
	MessageFetchLimit int
	CronInterval      int
	MaxConcurrentJobs int
}
