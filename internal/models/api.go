package models

// CronRequest models the incoming JSON body for cron control.
type CronRequest struct {
	Action string `json:"action"`
}
type SendMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
