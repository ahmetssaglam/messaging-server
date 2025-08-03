package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvStr(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Fatal("Error converting environment variable", key, "to int:", err)
		}
	}
	return fallback
}
