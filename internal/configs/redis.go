package configs

import (
	"messaging-server/internal/models"
	pkgUtils "messaging-server/pkg/utils"
)

var RedisConfig = models.RedisConfigStruct{
	Host:     pkgUtils.GetEnvStr("REDIS_HOST", "localhost"),
	Port:     pkgUtils.GetEnvStr("REDIS_PORT", "6379"),
	Password: pkgUtils.GetEnvStr("REDIS_PASSWORD", ""),
	DB:       pkgUtils.GetEnvInt("REDIS_DB", 0),
	TTL:      pkgUtils.GetEnvInt("REDIS_TTL", 3600),
}
