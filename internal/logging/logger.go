package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"messaging-server/internal/configs"
	"strings"
)

var Logger *log.Logger

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	// ANSI colors
	colors := map[string]string{
		"reset":  "\033[0m",
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"purple": "\033[35m",
		"cyan":   "\033[36m",
		"white":  "\033[97m",
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	// Color level
	var levelColor string
	var messageColor = colors["white"]

	switch entry.Level {
	case log.InfoLevel:
		levelColor = colors["green"]
	case log.WarnLevel:
		levelColor = colors["yellow"]
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = colors["red"]
		messageColor = colors["red"]
	case log.DebugLevel:
		levelColor = colors["purple"]
	default:
		levelColor = colors["white"]
	}

	coloredLevel := fmt.Sprintf("%s[%s]%s", levelColor, level, colors["reset"])
	coloredTime := fmt.Sprintf("%s%s%s", colors["cyan"], timestamp, colors["reset"])
	coloredMessage := fmt.Sprintf("%s%s%s", messageColor, entry.Message, colors["reset"])

	logLine := fmt.Sprintf("%s: %s - %s\n", coloredLevel, coloredTime, coloredMessage)
	return []byte(logLine), nil
}

func InitLogger() {

	Logger = log.New()
	Logger.SetFormatter(&CustomFormatter{})

	// log level
	logLevel, err := log.ParseLevel(strings.ToLower(configs.AppConfig.LogLevel))
	if err != nil {
		Logger.Errorf("Invalid log level: %s, defaulting to INFO", configs.AppConfig.LogLevel)
		logLevel = log.InfoLevel // Default to INFO if parsing fails
	}

	Logger.SetLevel(logLevel)
}
