package models

type RedisConfigStruct struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      int
}

type RedisRecord struct {
	MessageID string `json:"messageId"`
	SentAt    string `json:"sentAt"`
}
